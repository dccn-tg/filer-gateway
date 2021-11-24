package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	hapi "github.com/Donders-Institute/filer-gateway/internal/api-server/handler"
	"github.com/Donders-Institute/filer-gateway/internal/task"
	"github.com/Donders-Institute/filer-gateway/internal/worker/config"
	"github.com/Donders-Institute/filer-gateway/pkg/filer"

	"github.com/go-redis/redis/v8"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/acl"
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
	case "freenas":
		fConfig = cfg.FreeNas
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

	for _, api := range []filer.Filer{filer.New("netapp", cfg.NetApp), filer.New("freenas", cfg.FreeNas), filer.New("cephfs", cfg.CephFs)} {
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

	// `spath` is the logical project path under the de-facto top-level directory for all presented projects defined by `hapi.PathProject`.
	// It may be a physical directory or a symbolic link to a physical directory on a different storage system.
	spath := filepath.Join(hapi.PathProject, data.ProjectID)

	// determine the physical project directory and the corresponding storage api.
	var ppath string
	var api filer.Filer
	if data.Storage.System == "none" {
		// storage system name is not provided, assuming the project storage has been provisioned at `spath`
		// - resolving `ppath` from the physical path of `spath`
		// - resolving `api` from the physical path
		if ppath, err = filepath.EvalSymlinks(spath); err != nil {
			log.Errorf("fail to resolve storage system path: %s", err)
			return err
		}
		if api, err = getFilerAPIByPath(ppath, h.ConfigFile); err != nil {
			log.Errorf("fail to resolve storage API: %s", err)
			return err
		}
	} else {
		// storage sytem name is provided
		// - get `api` corresponding to the storage system name
		// - derive `ppath` from the `api`, assuming that projects are organized in sub-directories with project id as the path
		api, err = getFilerAPIBySystem(data.Storage.System, h.ConfigFile)
		if err != nil {
			log.Errorf("cannot load filer api: %s", err)
			return err
		}
		ppath = filepath.Join(api.GetProjectRoot(), data.ProjectID)
	}

	// only performs storage quota update when
	// - requested quota >= 0
	if data.Storage.QuotaGb >= 0 {
		// create project namespace or update project quota depending on whether the project directory exists.
		if _, err := os.Stat(ppath); os.IsNotExist(err) {
			// call filer API to create project volume and/or namespace
			if err := api.CreateProject(data.ProjectID, int(data.Storage.QuotaGb)); err != nil {
				log.Errorf("fail to create space for project %s: %s", data.ProjectID, err)
				return err
			}
			log.Debugf("project space created on %s at path %s with quota %d GB", data.Storage.System, ppath, data.Storage.QuotaGb)
		} else {
			// call filer API to update project quota
			if err := api.SetProjectQuota(data.ProjectID, int(data.Storage.QuotaGb)); err != nil {
				log.Errorf("fail to set quota for project %s: %s", data.ProjectID, err)
				return err
			}
			log.Debugf("project space quota on %s at path %s updated to %d GB", data.Storage.System, ppath, data.Storage.QuotaGb)
		}
	}

	// wait the ppath to present on the filesystem up to 1 minute.  Sometimes it appears immeidately; but it can be that
	// there is a significant delay between the volume creation and the path's availability to the host on which
	// the filer-gateway api server is running.
	tick := time.Now()
	for {
		if time.Since(tick) > time.Minute {
			log.Errorf("timeout waiting for file path to be available: ", ppath)
			break
		}
		if _, err := os.Stat(ppath); os.IsNotExist(err) {
			log.Debugf("wait for file path to be available: ", ppath)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	// make sure ppath is presented under the `api-server.ProjectRoot` directory.
	if _, err := os.Stat(spath); os.IsNotExist(err) {
		// make symlink to ppath
		if err := os.Symlink(ppath, spath); err != nil {
			log.Errorf("cannot make symlink %s -> %s: %s", spath, ppath, err)
		}
	}

	// ACL setting on the filesystem path of the project storage.
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

	// set user acl
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
			err = fmt.Errorf("fail setting roles (ec=%d): %s", ec, err)
			log.Errorf("%s", err)
			return err
		}
	}

	// delete user acl
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
			err = fmt.Errorf("fail removing roles (ec=%d): %s", ec, err)
			log.Errorf("%s", err)
			return err
		}
	}

	// notify api server to update cache for the project
	p := task.UpdateProjectPayload{
		ProjectID: data.ProjectID,
	}

	if m, err := json.Marshal(p); err == nil {
		h.ApiNotifierClient.Publish(context.Background(), "api_pcache_update", string(m))
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
		log.Errorf("Marshal error: %s", err)
		return err
	}

	var data task.SetUserResource

	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Errorf("Unmarshal error: %s", err)
		return err
	}

	log.Debugf("payload data: %+v", data)

	// check if user exists on the system.
	u, err := user.Lookup(data.UserID)
	if err != nil {
		log.Errorf("cannot find user %s: %s", data.UserID, err)
		return err
	}

	// get user's primary group.
	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		log.Errorf("cannot get user's primary group id %s: %s", u.Gid, err)
		return err
	}

	// check if user's home dir is under the group namespace
	gdir := filepath.Join("/home", g.Name)

	if !strings.Contains(u.HomeDir, gdir) {
		err = fmt.Errorf("user home dir %s not in group dir %s", u.HomeDir, gdir)
		log.Errorf("%s", err)
		return err
	}

	// skip quota setup for user's home space.
	if data.Storage.QuotaGb < 0 {
		log.Warnf("skip setting quota to %d Gb for home space %s", data.Storage.QuotaGb, u.HomeDir)
		return nil
	}

	// setting project resource
	ssystem := data.Storage.System
	if ssystem == "none" {
		ssystem = hapi.DefaultHomeStorageSystem
	}
	api, err := getFilerAPIBySystem(ssystem, h.ConfigFile)
	if err != nil {
		log.Errorf("cannot load filer api: %s", err)
		return err
	}

	// create home if user's home dir doesn't exist
	if _, err := os.Stat(u.HomeDir); os.IsNotExist(err) {
		// call filer API to create qtree for user's home
		if err := api.CreateHome(u.Username, g.Name, int(data.Storage.QuotaGb)); err != nil {
			log.Errorf("fail to create home space for user %s: %s", u.Username, err)
			return err
		}
		log.Debugf("home space created on %s at path %s", data.Storage.System, u.HomeDir)

		// change owner and group for the homedir
		nuid, _ := strconv.Atoi(u.Uid)
		ngid, _ := strconv.Atoi(u.Gid)

		// wait the home directory to present on the filesystem up to 1 minute.  Sometimes it appears immeidately; but it can be that
		// there is a significant delay between the qtree creation and the path's availability to the host on which
		// the filer-gateway api server is running.
		tick := time.Now()
		for {
			if time.Since(tick) > time.Minute {
				log.Errorf("timeout waiting for file path to be available: ", u.HomeDir)
				break
			}
			if _, err := os.Stat(u.HomeDir); os.IsNotExist(err) {
				log.Debugf("wait for file path to be available: ", u.HomeDir)
				time.Sleep(3 * time.Second)
			} else {
				break
			}
		}

		if err := os.Chown(u.HomeDir, nuid, ngid); err != nil {
			log.Errorf("fail to set owner of home space %s: %s", u.HomeDir, err)
			return err
		}

	} else {
		log.Warnf("skip home space creation as path already exists: %s", u.HomeDir)
	}

	// update storage quota
	if err := api.SetHomeQuota(u.Username, g.Name, int(data.Storage.QuotaGb)); err != nil {
		log.Errorf("fail to set home space quota for %s: %s", u.HomeDir, err)
		return err
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
