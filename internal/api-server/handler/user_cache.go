package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/api-server/config"
	"github.com/Donders-Institute/filer-gateway/internal/task"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type userResource struct {
	storage *models.StorageResponse
}

// UserResourceCache is an in-memory store for caching `projectResource` of all existing projects
// on the filer.
type UserResourceCache struct {

	// Config is the general API server configuration.
	Config config.Configuration

	// Context is the API server context.
	Context context.Context

	// Notifier is the redis channel subscription via which
	// a refresh on a given project can be triggered on-demand.
	Notifier <-chan *redis.Message

	// IsStopped indicates whether the cache service is stopped.
	IsStopped bool
	store     map[string]*userResource
	mutex     sync.RWMutex
}

// init initializes the cache with first reload.
func (c *UserResourceCache) Init() {

	c.IsStopped = false

	// first refresh
	c.refresh()

	// every 10 minutes??
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Infof("refreshing user cache")
				c.refresh()
				log.Infof("user cache refreshed")
			case m := <-c.Notifier:
				// interpret request payload
				p := task.UpdateUserPayload{}
				if err := json.Unmarshal([]byte(m.Payload), &p); err != nil {
					log.Errorf("unknown update payload: %s", m.Payload)
					continue
				}
				// perform cache update
				if _, err := c.getResource(p.UserID, true); err != nil {
					log.Errorf("fail to update cache for user %s: %s", p.UserID, err)
					continue
				}
				log.Infof("cache updated for user: %s", p.UserID)

			case <-c.Context.Done():
				log.Infof("user cache refresh stopped")
				c.IsStopped = true
				return
			}
		}
	}()

	log.Infof("user cache initalized")
}

// refresh update the cache with up-to-data project resources.
func (c *UserResourceCache) refresh() {

	nworkers := runtime.NumCPU()

	usernames := make(chan string, nworkers*2)
	resources := make(chan struct {
		username string
		resource *userResource
	})

	// TODO: cache the netapp volume quota report

	wg := sync.WaitGroup{}
	// start concurrent workers to get project resources from the filer.
	for i := 0; i < nworkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for uname := range usernames {
				storage, err := getUserStorageResource(uname, c.Config)
				if err != nil {
					log.Warnf("cannot get filer resource for %s: %s", uname, err)
				}
				resources <- struct {
					username string
					resource *userResource
				}{
					uname,
					&userResource{
						storage: storage,
					},
				}
			}
		}()
	}

	// go routine to list all directories in the /project folder
	go func() {
		// close the dirs channel on exit
		defer close(usernames)

		for _, uname := range getSystemUsers() {
			usernames <- uname
		}
	}()

	// go routine to wait for all workers to complete and close the resources channel.
	go func() {
		wg.Wait()
		close(resources)
	}()

	// new data map
	d := make(map[string]*userResource)

	// merge resources into internal store
	for r := range resources {
		d[r.username] = r.resource
	}

	// set store to new data map
	c.mutex.Lock()
	c.store = d
	c.mutex.Unlock()
}

// getResource gets resource information for the specific user.  It tries to get it from
// the cache.  If not available, it will retrieve up-to-date information from the storage
// (either via the filesystem or the storage's API) and update the cache accordingly.
func (c *UserResourceCache) getResource(username string, force bool) (*userResource, error) {

	// try to get resource from the cache
	c.mutex.RLock()
	r, ok := c.store[username]
	c.mutex.Unlock()

	// try to retrieve the resource from upstream filer/storage
	if !ok || force {

		storage, err := getUserStorageResource(username, c.Config)

		if err != nil {
			return nil, err
		}

		c.mutex.Lock()
		c.store[username] = &userResource{
			storage: storage,
		}

		r = c.store[username]
		c.mutex.Unlock()
	}

	return r, nil
}

// getSystemUsers get a list of usernames from the `getent passwd` system call, and filter out
// users with UID <= 1000.
func getSystemUsers() []string {
	usernames := make([]string, 0)

	cmd := exec.Command("/usr/bin/getent", "passwd")

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("%s", err)
		return usernames
	}

	sout := bufio.NewScanner(out)
	go func() {
		for sout.Scan() {
			// expected first three columns of getent output
			//
			// `{username}:*:{userid}`
			//
			line := sout.Text()
			data := strings.Split(line, ":")

			if len(data) < 3 {
				log.Warnf("unexpected getent output: %s", line)
				continue
			}

			// get user id
			uid, err := strconv.Atoi(data[2])
			if err != nil {
				log.Warnf("%s", err)
				continue
			}

			// keep username with uid >= 1000
			if uid >= 1000 {
				usernames = append(usernames, data[0])
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Errorf("%s", err)
		return usernames
	}

	err = cmd.Wait()
	if err != nil {
		log.Errorf("%s", err)
		return usernames
	}

	return usernames
}

// getUserStorageResource retrieves storage resource information of the user home directory.
func getUserStorageResource(username string, cfg config.Configuration) (*models.StorageResponse, error) {
	u, err := user.Lookup(username)

	if err != nil {
		return nil, err
	}

	// not found in cache, try fetch from the filer.
	system, quota, usage, err := getStorageQuota(cfg, u.HomeDir, true)
	if err != nil {
		return nil, err
	}

	// update cache with data retrieved from the filer.
	quotaGb := quota >> 30
	usageMb := usage >> 20

	return &models.StorageResponse{
		System:  &system,
		QuotaGb: &quotaGb,
		UsageMb: &usageMb,
	}, nil
}
