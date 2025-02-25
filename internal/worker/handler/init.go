package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	hapi "github.com/dccn-tg/filer-gateway/internal/api-server/handler"
	"github.com/dccn-tg/filer-gateway/internal/task"
	"github.com/dccn-tg/filer-gateway/internal/worker/config"
	"github.com/dccn-tg/filer-gateway/pkg/filer"

	"github.com/go-redis/redis/v8"

	ufp "github.com/dccn-tg/tg-toolset-golang/pkg/filepath"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
	"github.com/dccn-tg/tg-toolset-golang/pkg/mailer"
	"github.com/dccn-tg/tg-toolset-golang/project/pkg/acl"
	"github.com/dccn-tg/tg-toolset-golang/project/pkg/pdb"
	"github.com/hurngchunlee/bokchoy"
)

// getFilerAPIBySystem
func getFilerAPIBySystem(system, configFile string) (filer.Filer, error) {

	// load filer config and panic out if there is a problem loading it.
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("fail to laod filer configuration %s: %s", configFile, err)
	}

	var fConfig filer.Config
	switch system {
	case "netapp":
		fConfig = cfg.NetApp
	case "cephfs":
		fConfig = cfg.CephFs
	default:
		return nil, fmt.Errorf("unknown filer system name: %s", system)
	}

	// initiate filer API instances
	return filer.New(system, fConfig), nil
}

// getFilerAPIByPath
func getFilerAPIByPath(path, configFile string) (filer.Filer, error) {

	// load filer config and panic out if there is a problem loading it.
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("fail to laod filer configuration %s: %s", configFile, err)
	}

	for _, api := range []filer.Filer{filer.New("netapp", cfg.NetApp), filer.New("cephfs", cfg.CephFs)} {
		if strings.HasPrefix(path, filepath.Clean(api.GetProjectRoot())+"/") {
			return api, nil
		}
	}

	return nil, fmt.Errorf("unknown filer system for path: %s", path)
}

// TaskResults defines the output structure of the task
type TaskResults struct {
	Error error  `json:"errors"`
	Info  string `json:"info"`
}

// SetProjectResourceHandler implements `bokchoy.Handler` for applying update on project resource.
type SetProjectResourceHandler struct {
	// Configuration file for the worker
	ConfigFile        string
	ApiNotifierClient *redis.Client
}

// Handle performs project resource update based on the request payload.
func (h *SetProjectResourceHandler) Handle(r *bokchoy.Request) error {

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

	// `lpath` is the logical project path under the de-facto top-level directory for all presented projects defined by `hapi.PathProject`.
	// It may be a physical directory or a symbolic link to a physical directory on a different storage system.
	lpath := filepath.Join(hapi.PathProject, data.ProjectID)

	// determine the physical project directory and the corresponding storage api.
	var ppath string
	var api filer.Filer
	if data.Storage.System == "none" {
		// storage system name is not provided, assuming the project storage has been provisioned at `spath`
		// - resolving `ppath` from the locgical path of `lpath`
		// - resolving `api` from the physical path
		if ppath, err = filepath.EvalSymlinks(lpath); err != nil {
			log.Errorf("[%s] fail to resolve storage system path: %s", r.Task.ID, err)
			return err
		}
		if api, err = getFilerAPIByPath(ppath, h.ConfigFile); err != nil {
			log.Errorf("[%s] fail to resolve storage API: %s", r.Task.ID, err)
			return err
		}
	} else {
		// storage sytem name is provided
		// - get `api` corresponding to the storage system name
		// - derive `ppath` from the `api`, assuming that projects are organized in sub-directories with project id as the path
		api, err = getFilerAPIBySystem(data.Storage.System, h.ConfigFile)
		if err != nil {
			log.Errorf("[%s] fail to load filer api: %s", r.Task.ID, err)
			return err
		}
		ppath = filepath.Join(api.GetProjectRoot(), data.ProjectID)
	}

	// only performs storage quota update when
	// - requested quota >= 0
	isNewProject := false
	if data.Storage.QuotaGb >= 0 {
		// create project namespace or update project quota depending on whether the project directory exists.
		if _, err := os.Stat(ppath); os.IsNotExist(err) {
			isNewProject = true
			// call filer API to create project volume and/or namespace
			if err := api.CreateProject(data.ProjectID, int(data.Storage.QuotaGb)); err != nil {
				log.Errorf("[%s] fail to create space for project %s: %s", r.Task.ID, data.ProjectID, err)
				return err
			}
			log.Debugf("[%s] project space created on %s at path %s with quota %d GB", r.Task.ID, data.Storage.System, ppath, data.Storage.QuotaGb)
		} else {
			// call filer API to update project quota
			if err := api.SetProjectQuota(data.ProjectID, int(data.Storage.QuotaGb)); err != nil {
				log.Errorf("[%s] fail to set quota for project %s: %s", r.Task.ID, data.ProjectID, err)
				return err
			}
			log.Debugf("[%s] project space quota on %s at path %s updated to %d GB", r.Task.ID, data.Storage.System, ppath, data.Storage.QuotaGb)
		}
	}

	// create symlink under the `api-server.ProjectRoot` directory for the case that the
	// project is not created on the main filer system (i.e. the filer mounted directly
	// under the `api-server.ProjectRoot`).
	if _, err := os.Stat(lpath); os.IsNotExist(err) {
		// make symlink to ppath
		if err := os.Symlink(ppath, lpath); err != nil {
			log.Errorf("[%s] fail to make symlink %s -> %s: %s", r.Task.ID, lpath, ppath, err)
		}
	}

	// perform ACL settings for updating members.
	managers := make([]string, 0)
	contributors := make([]string, 0)
	viewers := make([]string, 0)
	udelete := make([]string, 0)

	for _, m := range data.Members {
		switch m.Role {
		case acl.Manager.String():
			managers = append(managers, m.UserID)
		case acl.Contributor.String():
			contributors = append(contributors, m.UserID)
		case acl.Viewer.String():
			viewers = append(viewers, m.UserID)
		case "none":
			udelete = append(udelete, m.UserID)
		}
	}

	// set members
	if len(managers)+len(contributors)+len(viewers) > 0 {

		switch data.Recursion {
		case true:
			// use runner to set ACL recursively
			log.Debugf("[%s] setting ACL for members with recursion: %s", r.Task.ID, data.ProjectID)

			runner := acl.Runner{
				Managers:     strings.Join(managers, ","),
				Contributors: strings.Join(contributors, ","),
				Viewers:      strings.Join(viewers, ","),
				RootPath:     ppath,
				Traverse:     false,
				Force:        false,
				FollowLink:   false,
				SkipFiles:    false,
				Silence:      true,
				Nthreads:     4,
			}

			if ec, err := runner.SetRoles(); ec != 0 || err != nil {
				err = fmt.Errorf("fail to set member roles (ec=%d): %s", ec, err)
				log.Errorf("[%s] %s", r.Task.ID, err)
				return err
			}
		case false:

			// use roler to set top-level ACL
			log.Debugf("[%s] setting ACL for members without recursion: %s", r.Task.ID, data.ProjectID)

			// state the physical path to get the mode
			ppathInfo, err := os.Stat(ppath)
			if err != nil {
				err = fmt.Errorf("fail to get info of %s: %s", ppath, err)
				log.Errorf("[%s] %s", r.Task.ID, err)
				return err
			}

			spathMode := ufp.FilePathMode{
				Path: ppath,
				Mode: ppathInfo.Mode(),
			}

			roler := acl.GetRoler(spathMode)

			// construct acl.RoleMap
			rmap := acl.RoleMap{
				acl.Manager:     managers,
				acl.Contributor: contributors,
				acl.Viewer:      viewers,
			}

			if _, err := roler.SetRoles(spathMode, rmap, false, false); err != nil {
				err = fmt.Errorf("fail to set member roles %s", err)
				log.Errorf("[%s] %s", r.Task.ID, err)
				return err
			}
		}
	}

	// delete members
	if len(udelete) > 0 {
		switch data.Recursion {
		case true:
			// delete ACL recursively
			log.Debugf("[%s] deleting ACL for members with recursion: %s", r.Task.ID, data.ProjectID)

			udelstr := strings.Join(udelete, ",")
			runner := acl.Runner{
				RootPath:     ppath,
				Managers:     udelstr,
				Contributors: udelstr,
				Viewers:      udelstr,
				Traversers:   udelstr,
				FollowLink:   false,
				SkipFiles:    false,
				Nthreads:     4,
				Silence:      true,
				Traverse:     false,
				Force:        false,
			}

			if ec, err := runner.RemoveRoles(); ec != 0 || err != nil {
				err = fmt.Errorf("fail to remove member roles (ec=%d): %s", ec, err)
				log.Errorf("[%s] %s", r.Task.ID, err)
				return err
			}

		case false:
			// delete ACL only on top-level project directory
			log.Debugf("[%s] deleting ACL for members without recursion: %s", r.Task.ID, data.ProjectID)

			// state the physical path to get the mode
			ppathInfo, err := os.Stat(ppath)
			if err != nil {
				err = fmt.Errorf("fail to get info of %s: %s", ppath, err)
				log.Errorf("[%s] %s", r.Task.ID, err)
				return err
			}

			spathMode := ufp.FilePathMode{
				Path: ppath,
				Mode: ppathInfo.Mode(),
			}

			roler := acl.GetRoler(spathMode)

			// construct acl.RoleMap
			rmap := acl.RoleMap{
				acl.Manager:     udelete,
				acl.Contributor: udelete,
				acl.Viewer:      udelete,
			}

			if _, err := roler.DelRoles(spathMode, rmap, false, false); err != nil {
				err = fmt.Errorf("fail to remove member roles %s", err)
				log.Errorf("[%s] %s", r.Task.ID, err)
				return err
			}
		}
	}

	// notify api server to update cache for the project
	p := task.UpdateProjectPayload{
		ProjectID: data.ProjectID,
	}

	if m, err := json.Marshal(p); err == nil {
		h.ApiNotifierClient.Publish(context.Background(), "api_pcache_update", string(m))
	}

	// notify managers if new project is created
	if isNewProject {
		err := h.notifyProjectProvisioned(data.ProjectID, managers, contributors)
		if err != nil {
			log.Errorf("[%s] %s", r.Task.ID, err)
		}
	}

	return nil
}

func (h *SetProjectResourceHandler) notifyProjectProvisioned(projectID string, managers, contributors []string) error {
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

	var template string
	var mailTemplate string

	recipients := managers

	switch p.Kind {
	case pdb.Dataset:
		recipients = append(recipients, contributors...)
		mailTemplate = "/etc/filer-gateway/new-dataset-project.template"
		template = `Storage of your dataset project {{.ProjectID}} has been initalized!

Dear {{.RecipientName}},

The storage of your dataset project {{.ProjectID}} with title

	{{.ProjectTitle}}

has been initialised.

Dataset builders specified on the PPM form are given the contributor (i.e. read-write) permission for building the dataset. They may access the storage via the following paths:

	* on Windows desktop: R:\{{.ProjectID}}
	* in the cluster: /project_cephfs/{{.ProjectID}}

Please note that data in this project SHOULD NOT be unique as the storage is NOT backed up.

After you complete the dataset build process, please notify the DCCN data support (datasupport@donders.ru.nl).  The data steward will conduct a review on the dataset content followed by an adjustment to the data access control by the TG to enable the data sharing/reuse.

For procedures related to managing the dataset project, please refer to:

	https://intranet.donders.ru.nl/index.php?id=dataset

Should you have any questions, please don't hesitate to contact the DCCN datasupport <datasupport@donders.ru.nl>.

Best regards, {{.SenderName}}
		`

	default:
		mailTemplate = "/etc/filer-gateway/new-research-project.template"
		template = `Storage of your research project {{.ProjectID}} has been initalized!

Dear {{.RecipientName}},

The storage of your research project {{.ProjectID}} with title

	{{.ProjectTitle}}

has been initialised.

You may now access the storage via the following paths:

	* on Windows desktop: P:\{{.ProjectID}}
	* in the cluster: /project/{{.ProjectID}}

For managing data access permission for project collaborators, please follow the guide:

	http://hpc.dccn.nl/docs/project_storage/access_management.html

For more information about the project storage, please refer to the intranet page:

	https://intranet.donders.ru.nl/index.php?id=4733

Should you have any questions, please don't hesitate to contact the TG helpdesk <helpdesk@donders.ru.nl>.

Best regards, {{.SenderName}}
		`
	}

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
				log.Errorf("cannot compose message to notify user %s for project %s: %s", u.Email, projectID, err)
				continue
			}

			err = m.SendMail("helpdesk@donders.ru.nl", subject, body, []string{u.Email})
			if err != nil {
				log.Errorf("cannot notify user %s for project %s: %s", u.Email, projectID, err)
			} else {
				log.Infof("user %s notified for provisioned project %s", recipient, projectID)
			}
		}
	}
	return nil
}

// SetUserResourceHandler implements `bokchoy.Handler` for applying update on user resource.
type SetUserResourceHandler struct {
	// Configuration file for the worker
	ConfigFile        string
	ApiNotifierClient *redis.Client
}

// Handle performs user resource update based on the request payload.
func (h *SetUserResourceHandler) Handle(r *bokchoy.Request) error {

	res, err := json.Marshal(r.Task.Payload)
	if err != nil {
		log.Errorf("[%s] fail to serialize (marshal) payload: %s", r.Task.ID, err)
		return err
	}

	var data task.SetUserResource

	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Errorf("[%s] fail to de-serialize (unmarshal) payload: %s", r.Task.ID, err)
		return err
	}

	// check if user exists on the system.
	u, err := user.Lookup(data.UserID)
	if err != nil {
		log.Errorf("[%s] fail to find user %s: %s", r.Task.ID, data.UserID, err)
		return err
	}

	// get user's primary group.
	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		log.Errorf("[%s] fail to get user's primary group id %s: %s", r.Task.ID, u.Gid, err)
		return err
	}

	// check if user's home dir is under the group namespace
	gdir := filepath.Join("/home", g.Name)

	if !strings.Contains(u.HomeDir, gdir) {
		err = fmt.Errorf("user home dir %s not in group dir %s", u.HomeDir, gdir)
		log.Errorf("[%s] %s", r.Task.ID, err)
		return err
	}

	// skip quota setup for user's home space.
	if data.Storage.QuotaGb < 0 {
		log.Warnf("[%s] skip setting quota to %d Gb for home space %s", r.Task.ID, data.Storage.QuotaGb, u.HomeDir)
		return nil
	}

	// setting project resource
	ssystem := data.Storage.System
	if ssystem == "none" {
		ssystem = hapi.DefaultHomeStorageSystem
	}
	api, err := getFilerAPIBySystem(ssystem, h.ConfigFile)
	if err != nil {
		log.Errorf("[%s] fail to load filer api: %s", r.Task.ID, err)
		return err
	}

	// create home if user's home dir doesn't exist
	if _, err := os.Stat(u.HomeDir); os.IsNotExist(err) {
		// call filer API to create qtree for user's home
		if err := api.CreateHome(u.Username, g.Name, int(data.Storage.QuotaGb)); err != nil {
			log.Errorf("[%s] fail to create home space for user %s: %s", r.Task.ID, u.Username, err)
			return err
		}
		log.Debugf("[%s] home space created on %s at path %s", r.Task.ID, data.Storage.System, u.HomeDir)
	} else {
		log.Warnf("[%s] skip home space creation as path already exists: %s", r.Task.ID, u.HomeDir)
		// update storage quota
		if err := api.SetHomeQuota(u.Username, g.Name, int(data.Storage.QuotaGb)); err != nil {
			log.Errorf("[%s] fail to set home space quota for %s: %s", r.Task.ID, u.HomeDir, err)
			return err
		}
	}

	// notify api server to update cache for the user
	p := task.UpdateUserPayload{
		UserID: data.UserID,
	}

	if m, err := json.Marshal(p); err == nil {
		h.ApiNotifierClient.Publish(context.Background(), "api_ucache_update", string(m))
	}

	return nil
}

// DelUserResourceHandler implements `bokchoy.Handler` for applying deletion on user home dir qtree.
type DelUserResourceHandler struct {
	// Configuration file for the worker
	ConfigFile        string
	ApiNotifierClient *redis.Client
}

// Handle performs user resource update based on the request payload.
func (h *DelUserResourceHandler) Handle(r *bokchoy.Request) error {

	res, err := json.Marshal(r.Task.Payload)
	if err != nil {
		log.Errorf("[%s] fail to serialize (marshal) payload: %s", r.Task.ID, err)
		return err
	}

	var data task.DelUserResource

	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Errorf("[%s] fail to de-serialize (unmarshal) payload: %s", r.Task.ID, err)
		return err
	}

	// check if user exists on the system.
	u, err := user.Lookup(data.UserID)
	if err != nil {
		log.Errorf("[%s] fail to find user %s: %s", r.Task.ID, data.UserID, err)
		return err
	}

	// get user's primary group.
	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		log.Errorf("[%s] fail to get user's primary group id %s: %s", r.Task.ID, u.Gid, err)
		return err
	}

	// check if user's home dir is under the group namespace
	gdir := filepath.Join("/home", g.Name)

	if !strings.Contains(u.HomeDir, gdir) {
		err = fmt.Errorf("user home dir %s not in group dir %s", u.HomeDir, gdir)
		log.Errorf("[%s] %s", r.Task.ID, err)
		return err
	}

	// home directory doesn't exist, skip qtree deletion
	if _, err := os.Stat(u.HomeDir); os.IsNotExist(err) {
		log.Debugf("[%s] home dir %s doesn't exist, skip deletion", r.Task.ID, u.HomeDir)
	} else {
		// check if the homedir contains data
		if empty, _ := filer.IsDirEmpty(u.HomeDir); !empty {
			log.Errorf("[%s] home dir %s not empty", r.Task.ID, u.HomeDir)
			return err
		}
		// perform deletion
		// get filer API, for home qtree, the system is always "netapp"
		api, err := getFilerAPIBySystem("netapp", h.ConfigFile)
		if err != nil {
			log.Errorf("[%s] fail to load filer api: %s", r.Task.ID, err)
			return err
		}
		if err := api.DeleteHome(u.Username, g.Name); err != nil {
			log.Errorf("[%s] fail to delete user qtree, user: %s group:%s", u.Username, g.Name)
			return err
		}
	}

	// notify api server to update cache for the user
	p := task.UpdateUserPayload{
		UserID: data.UserID,
	}

	if m, err := json.Marshal(p); err == nil {
		h.ApiNotifierClient.Publish(context.Background(), "api_ucache_update", string(m))
	}

	return nil
}
