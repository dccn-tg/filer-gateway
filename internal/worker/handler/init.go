package handler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dccn-tg/filer-gateway/internal/worker/config"
	"github.com/dccn-tg/filer-gateway/pkg/filer"
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
