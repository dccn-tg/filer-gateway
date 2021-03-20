package handler

import (
	"context"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/api-server/config"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	fp "github.com/Donders-Institute/tg-toolset-golang/pkg/filepath"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

type projectResource struct {
	storage *models.StorageResponse
	members []*models.Member
}

// ProjectResourceCache is an in-memory store for caching `projectResource` of all existing projects
// on the filer.
type ProjectResourceCache struct {
	Config  config.Configuration
	Context context.Context
	store   map[string]*projectResource
	mutex   sync.Mutex
}

// init initializes the cache with first reload.
func (c *ProjectResourceCache) Init() {

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
			case <-c.Context.Done():
				log.Infof("stopping cache refresh")
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

	// clean up the store
	c.mutex.Lock()
	c.store = make(map[string]*projectResource)
	c.mutex.Unlock()

	// merge resources into internal store
	for r := range resources {
		c.mutex.Lock()
		c.store[r.pnumber] = r.resource
		c.mutex.Unlock()
	}
}

// getProjectResource finds and returns project resource from the cache.
// An error is returned if the project doesn't exist in cache.
func (c *ProjectResourceCache) getProjectResource(pnumber string) (*projectResource, error) {
	if r, ok := c.store[pnumber]; !ok {

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
