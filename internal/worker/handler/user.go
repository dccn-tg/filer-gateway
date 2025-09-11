package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	hapi "github.com/dccn-tg/filer-gateway/internal/api-server/handler"
	"github.com/dccn-tg/filer-gateway/internal/task"
	"github.com/dccn-tg/filer-gateway/pkg/filer"

	"github.com/go-redis/redis/v8"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
	"github.com/hurngchunlee/bokchoy"
)

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
