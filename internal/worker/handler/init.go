package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	hapi "github.com/Donders-Institute/filer-gateway/internal/api-server/handler"
	"github.com/Donders-Institute/filer-gateway/internal/task"
	"github.com/Donders-Institute/filer-gateway/internal/worker/config"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/acl"
	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/filer"
	"github.com/thoas/bokchoy"
)

// getFilerAPI
func getFilerAPI(system, configFile string) (filer.Filer, error) {

	// load filer config and panic out if there is a problem loading it.
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("fail to laod filer configuration %s: %s", configFile, err)
	}

	var fConfig filer.Config
	switch system {
	case "netapp":
		fConfig = cfg.NetApp
	case "freenas":
		fConfig = cfg.FreeNas
	default:
		return nil, fmt.Errorf("unknown filer system: %s", system)
	}

	// initiate filer API instances
	return filer.New(system, fConfig), nil
}

// TaskResults defines the output structure of the task
type TaskResults struct {
	Error error  `json:"errors"`
	Info  string `json:"info"`
}

// SetProjectResourceHandler implements `bokchoy.Handler` for applying update on project resource.
type SetProjectResourceHandler struct {
	// Configuration file for the worker
	ConfigFile string
}

// Handle performs project resource update based on the request payload.
func (h *SetProjectResourceHandler) Handle(r *bokchoy.Request) error {

	res, err := json.Marshal(r.Task.Payload)
	if err != nil {
		log.Errorf("Marshal error: %s", err)
		return err
	}

	var data task.SetProjectResource

	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Errorf("Unmarshal error: %s", err)
		return err
	}

	log.Debugf("payload data: %+v", data)

	// setting project resource
	api, err := getFilerAPI(data.Storage.System, h.ConfigFile)
	if err != nil {
		return err
	}

	// 1. create project namespace
	ppath := filepath.Join(hapi.PathProject, data.ProjectID)
	if _, err := os.Stat(ppath); os.IsNotExist(err) {
		// call filer API to create project volume and/or namespace
		if err := api.CreateProject(data.ProjectID, int(data.Storage.QuotaGb)); err != nil {
			return fmt.Errorf("fail to create space for project %s: %s", data.ProjectID, err)
		}
		log.Debugf("project space created on %s at path %s", data.Storage.System, ppath)
	} else {
		log.Warnf("skip project space creation as project path already exists: %s", ppath)
	}

	// 2. update project quota
	_, quota, _, err := hapi.GetStorageQuota(ppath)
	if err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	if data.Storage.QuotaGb != quota {
		// call filer API to set the new quota
		if err := api.SetProjectQuota(data.ProjectID, int(data.Storage.QuotaGb)); err != nil {
			return fmt.Errorf("fail to set space quota for project %s: %s", data.ProjectID, err)
		}
		log.Debugf("project storage quota set from %d Gb to %d Gb", quota, data.Storage.QuotaGb)
	} else {
		log.Warnf("skip setting project space quota as the quota is already in right size: project %s quota %d", data.ProjectID, quota)
	}

	// 3. set/delete project roles
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

	if len(managers)+len(contributors)+len(viewers) > 0 {
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
			return fmt.Errorf("fail setting roles (ec=%d): %s", ec, err)
		}
	}

	// 4. delete user
	if len(udelete) > 0 {
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
			return fmt.Errorf("fail removing roles (ec=%d): %s", ec, err)
		}
	}

	// Put payload data as result.
	//r.Task.Result = &data

	return nil
}

// SetUserResourceHandler implements `bokchoy.Handler` for applying update on user resource.
type SetUserResourceHandler struct {
	// Configuration file for the worker
	ConfigFile string
}

// Handle performs user resource update based on the request payload.
func (h *SetUserResourceHandler) Handle(r *bokchoy.Request) error {
	return nil
}
