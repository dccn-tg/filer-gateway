package handler

import (
	"context"
	"sync"
	"time"

	"github.com/dccn-tg/filer-gateway/internal/api-server/config"
	"github.com/dccn-tg/filer-gateway/pkg/filer"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

type systemInfo struct {
	totalGiB int64
	usedGiB  int64
}

// SystemInfoCache is an in-memory store for caching `systemInfo` of storage systems.
type SystemInfoCache struct {

	// Config is the general API server configuration.
	Config config.Configuration

	// Context is the API server context.
	Context context.Context

	// IsStopped indicates whether the cache service is stopped.
	IsStopped bool
	store     map[string]*systemInfo
	mutex     sync.RWMutex
}

// init initializes the cache with first reload.
func (c *SystemInfoCache) Init() {

	c.IsStopped = false

	// first refresh
	c.refresh()

	// every 10 minutes??
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Infof("refreshing system info cache")
				c.refresh()
				log.Infof("system info cache refreshed")

			case <-c.Context.Done():
				log.Infof("system info cache refresh stopped")
				c.IsStopped = true
				return
			}
		}
	}()

	log.Infof("system info cache initalized")
}

// refresh update the cache with up-to-data project resources.
func (c *SystemInfoCache) refresh() {

	// new data map
	d := make(map[string]*systemInfo)

	systems := map[string]filer.Filer{
		"netapp": filer.New("netapp", c.Config.NetApp),
		"cephfs": filer.New("cephfs", c.Config.CephFs),
	}

	// netapp system
	for s, f := range systems {
		total, used, err := f.GetSystemSpaceInBytes()
		if err != nil {
			log.Errorf("fail refreshing storage system (%s) info: %s\n", s, err)

			// copy current system info
			if cur, ok := c.store[s]; ok {
				d[s] = cur
			}
		} else {
			// update with new system info
			d[s] = &systemInfo{
				totalGiB: total << 30,
				usedGiB:  used << 30,
			}
		}
	}

	// set store to new data map
	c.mutex.Lock()
	c.store = d
	c.mutex.Unlock()
}

// getAllSystems gets information of all storage systems from the cache.
// If `force` argument is `true`, the cache is refreshed (by updating data from the filers) before
// it returns.
func (c *SystemInfoCache) getAllSystems(force bool) map[string]*systemInfo {
	if force {
		c.refresh()
	}
	return c.store
}
