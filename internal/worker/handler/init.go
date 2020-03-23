package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	hapi "github.com/Donders-Institute/filer-gateway/internal/api-server/handler"
	"github.com/Donders-Institute/filer-gateway/internal/task"
	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/acl"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/bokchoy"
)

// SetProjectResourceHandler implements `bokchoy.Handler` for applying update on project resource.
type SetProjectResourceHandler struct {
}

// Handle performs project resource update based on the request payload.
func (h *SetProjectResourceHandler) Handle(r *bokchoy.Request) error {

	log.Debugf("payload: %+v", r.Task.Payload)

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
	// 1. create project namespace
	ppath := filepath.Join(hapi.PathProject, data.ProjectID)
	if _, err := os.Stat(ppath); os.IsNotExist(err) {
		// call filer API to create project volume and/or namespace
		log.Debug("creating project storage on %s, path %s", data.Storage.System, ppath)
	}

	// 2. update project quota
	_, quota, _, err := hapi.GetStorageQuota(ppath)
	if err != nil {
		return err
	}
	if data.Storage.QuotaGb != quota {
		// call filer API to set the new quota
		log.Debug("setting project storage quota from %d Gb to %d Gb", quota, data.Storage.QuotaGb)
	}

	// 3. set project roles
	managers := make([]string, 0)
	contributors := make([]string, 0)
	viewers := make([]string, 0)

	for _, m := range data.Members {
		switch m.Role {
		case acl.Manager.String():
			managers = append(managers, m.UserID)
		case acl.Contributor.String():
			contributors = append(contributors, m.UserID)
		case acl.Viewer.String():
			viewers = append(viewers, m.UserID)
		}
	}

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

	// Put payload data as result.
	//r.Task.Result = &data

	return nil
}

// SetUserResourceHandler implements `bokchoy.Handler` for applying update on user resource.
type SetUserResourceHandler struct {
}

// Handle performs user resource update based on the request payload.
func (h *SetUserResourceHandler) Handle(r *bokchoy.Request) error {
	return nil
}
