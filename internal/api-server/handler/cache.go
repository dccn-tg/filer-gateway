package handler

import (
	"context"
	"encoding/json"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/api-server/config"
	"github.com/Donders-Institute/filer-gateway/internal/task"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	fp "github.com/Donders-Institute/tg-toolset-golang/pkg/filepath"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type projectResource struct {
	storage *models.StorageResponse
	members []*models.Member
}

// ProjectResourceCache is an in-memory store for caching `projectResource` of all existing projects
// on the filer.
type ProjectResourceCache struct {

	// Config is the general API server configuration.
	Config config.Configuration

	// Context is the API server context.
	Context context.Context

	// Notifier is the redis channel subscription via which
	// a refresh on a given project can be triggered on-demand.
	Notifier <-chan *redis.Message

	// IsStopped indicates whether the cache service is stopped.
	IsStopped bool
	store     map[string]*projectResource
	mutex     sync.Mutex
}

// init initializes the cache with first reload.
func (c *ProjectResourceCache) Init() {

	c.IsStopped = false

	// first refresh
	c.refresh()

	// every 10 minutes??
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Infof("refreshing cache")
				c.refresh()
				log.Infof("cache refreshed")
			case m := <-c.Notifier:
				// interpret request payload
				p := task.UpdatePayload{}
				if err := json.Unmarshal([]byte(m.Payload), &p); err != nil {
					log.Errorf("unknown update payload: %s", m.Payload)
					continue
				}
				// perform cache update
				if _, err := c.getProjectResource(p.ProjectNumber, true); err != nil {
					log.Errorf("fail to update cache for project %s: %s", p.ProjectNumber, err)
					continue
				}
				log.Infof("cache updated for project: %s", p.ProjectNumber)

			case <-c.Context.Done():
				log.Infof("cache refresh stopped")
				c.IsStopped = true
				return
			}
		}
	}()

	log.Infof("cache initalized")
}

// refresh update the cache with up-to-data project resources.
func (c *ProjectResourceCache) refresh() {

	nworkers := runtime.NumCPU()

	pnumbers := make(chan string, nworkers*2)
	resources := make(chan struct {
		pnumber  string
		resource *projectResource
	})

	wg := sync.WaitGroup{}
	// start concurrent workers to get project resources from the filer.
	for i := 0; i < nworkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pnumber := range pnumbers {
				storage, members, err := getProjectResource(pnumber, c.Config)
				if err != nil {
					log.Errorf("cannot get filer resource for %s: %s", pnumber, err)
				}
				resources <- struct {
					pnumber  string
					resource *projectResource
				}{
					pnumber,
					&projectResource{
						storage: storage,
						members: members,
					},
				}
			}
		}()
	}

	// go routine to list all directories in the /project folder
	go func(path string) {
		// close the dirs channel on exit
		defer close(pnumbers)

		objs, err := fp.ListDir(path)
		if err != nil {
			log.Errorf("cannot get content of path: %s", path)
			return
		}

		for _, obj := range objs {
			pnumbers <- filepath.Base(obj)
		}

	}(PathProject)

	// go routine to wait for all workers to complete and close the resources channel.
	go func() {
		wg.Wait()
		close(resources)
	}()

	// new data map
	d := make(map[string]*projectResource)

	// merge resources into internal store
	for r := range resources {
		d[r.pnumber] = r.resource
	}

	// set store to new data map
	c.mutex.Lock()
	c.store = d
	c.mutex.Unlock()
}

// getProjectResource finds and returns project resource from the cache.
// An error is returned if the project doesn't exist in cache.
func (c *ProjectResourceCache) getProjectResource(pnumber string, force bool) (*projectResource, error) {
	if r, ok := c.store[pnumber]; !ok || force {

		// not found in cache, try fetch from the filer.
		storage, members, err := getProjectResource(pnumber, c.Config)
		if err != nil {
			return nil, err
		}

		// update cache with data retrieved from the filer.
		c.mutex.Lock()
		c.store[pnumber] = &projectResource{
			storage: storage,
			members: members,
		}
		c.mutex.Unlock()

		// return the data from cache
		return c.store[pnumber], nil
	} else {
		return r, nil
	}
}
