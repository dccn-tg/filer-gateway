package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	hapi "github.com/dccn-tg/filer-gateway/internal/api-server/handler"
	"github.com/dccn-tg/filer-gateway/internal/task"
	"github.com/dccn-tg/filer-gateway/internal/worker/config"
	"github.com/dccn-tg/filer-gateway/pkg/filer"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
	"github.com/dccn-tg/tg-toolset-golang/pkg/mailer"
	"github.com/dccn-tg/tg-toolset-golang/project/pkg/pdb"
	"github.com/go-redis/redis/v8"
	"github.com/hurngchunlee/bokchoy"
)

// `getFilerAPI` returns NetApp filer API with configurations customized for `rrd4project` FlexVol
func getFilerAPI(configFile string) (filer.Filer, error) {

	// load filer config and panic out if there is a problem loading it.
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("fail to laod filer configuration %s: %s", configFile, err)
	}

	fConfig := cfg.NetApp
	fConfig.ProjectRoot = "rrd4project"
	fConfig.VolumeProjectQtrees = "rrd4project"

	// initiate filer API instances
	return filer.New("netapp", fConfig), nil
}

// SetProjectRrdResourceHandler implements `bokchoy.Handler` for applying update on project resource.
type SetProjectRrdResourceHandler struct {
	// Configuration file for the worker
	ConfigFile        string
	ApiNotifierClient *redis.Client
}

// Handle performs project resource update based on the request payload.
func (h *SetProjectRrdResourceHandler) Handle(r *bokchoy.Request) error {

	res, err := json.Marshal(r.Task.Payload)
	if err != nil {
		log.Errorf("[%s] fail to serialize (marshal) payload: %s", r.Task.ID, err)
		return err
	}

	var data task.SetProjectResource

	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Errorf("[%s] fail to de-serialize (unmarshal) payload: %s", r.Task.ID, err)
		return err
	}

	// modify `data.ProjectID` to add prefix `rrd` for the qtree name
	pid := data.ProjectID
	data.ProjectID = fmt.Sprintf("rrd%s", pid)

	api, err := getFilerAPI(h.ConfigFile)
	if err != nil {
		log.Errorf("[%s] fail to construct filer API: %s", r.Task.ID, err)
		return err
	}

	ppath := filepath.Join(hapi.PathProjectRrd, data.ProjectID)

	if err = setProjectResource(r.Task.ID, data, api, ppath, ppath, h.notifyProjectRrdProvisioned); err != nil {
		log.Errorf("[%s] %s", r.Task.ID, err)
		return err
	}

	// notify api server to update cache for the project
	p := task.UpdateProjectPayload{
		ProjectID: pid,
	}

	if m, err := json.Marshal(p); err == nil {
		h.ApiNotifierClient.Publish(context.Background(), "api_rrdcache_update", string(m))
	}

	return nil
}

func (h *SetProjectRrdResourceHandler) notifyProjectRrdProvisioned(projectID string, managers, contributors []string) error {
	cfg, err := config.LoadConfig(h.ConfigFile)
	if err != nil {
		return fmt.Errorf("cannot read config for mailer: %s", err)
	}

	// get project detail
	ipdb, err := pdb.New(cfg.Pdb)
	if err != nil {
		return fmt.Errorf("cannot read config for pdb: %s", err)
	}

	p, err := ipdb.GetProject(projectID)
	if err != nil {
		return fmt.Errorf("cannot get information of project %s: %s", projectID, err)
	}

	// send email to project managers
	m := mailer.New(cfg.Smtp)
	localRelay := false
	if cfg.Smtp.Host == "localhost" || strings.HasSuffix(cfg.Smtp.Host, "dccn.nl") {
		localRelay = true
	}

	data := mailer.ProjectAlertTemplateData{
		ProjectID:    projectID,
		ProjectTitle: p.Name,
		SenderName:   "DCCN TG Helpdesk",
	}

	recipients := append(managers, contributors...)
	mailTemplate := "/etc/filer-gateway/new-project-rrd.template"
	template := `Research-related data storage of your project {{.ProjectID}} has been initalized!

Dear {{.RecipientName}},

The research-related data storage of your project {{.ProjectID}} with title

	{{.ProjectTitle}}

has been initialised.

The hosting PI and the project owner specified on the PPM form are given the contributor (i.e. read-write) permission for the research-related data storage. They may access the storage via the following paths:

	* on Windows desktop: Q:\{{.ProjectID}}

Should you have any questions, please don't hesitate to contact the DCCN datasupport <datasupport@donders.ru.nl>.

Best regards, {{.SenderName}}
		`

	// validity check on the template file
	if state, err := os.Stat(mailTemplate); errors.Is(err, os.ErrNotExist) { // template file doesn't exist
		mailTemplate = ""
	} else if !state.Mode().IsRegular() { // template file is not a regular file
		mailTemplate = ""
	} else if state.Size() < 100 && state.Size() > 10240 { // template file size is less than 100 or larger than 10K bytes (suspicious content)
		mailTemplate = ""
	}

	for _, recipient := range recipients {
		if u, err := ipdb.GetUser(recipient); err != nil {
			log.Errorf("cannot get information of recipient %s: %s, skip notification", recipient, err)
		} else {
			// if DCCN local email relay is used, only the email in form of `{account_name}@localhos`
			// is allowed.
			if localRelay {
				u.Email = fmt.Sprintf("%s@localhost", recipient)
			}

			data.RecipientName = u.DisplayName()

			var subject string
			var body string
			var err error

			if mailTemplate != "" {
				if subject, body, err = mailer.ComposeMessageFromTemplateFile(mailTemplate, data); err != nil {
					log.Errorf("cannot compose message from template file %s: %s", mailTemplate, err)

					// fallback to compose mail with built-in template
					subject, body, err = mailer.ComposeMessageFromTemplate(template, data)
				}
			} else {
				subject, body, err = mailer.ComposeMessageFromTemplate(template, data)
			}

			if err != nil {
				log.Errorf("cannot compose message to notify user %s for project rrd %s: %s", u.Email, projectID, err)
				continue
			}

			err = m.SendMail("helpdesk@donders.ru.nl", subject, body, []string{u.Email})
			if err != nil {
				log.Errorf("cannot notify user %s for project rrd %s: %s", u.Email, projectID, err)
			} else {
				log.Infof("user %s notified for provisioned project rrd %s", recipient, projectID)
			}
		}
	}
	return nil
}
