package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	pathCfg := os.Getenv("FILER_GATEWAY_APISERVER_CONFIG")

	cfg, err := LoadConfig(pathCfg)

	if err != nil {
		t.Errorf("%s", err)
	}

	t.Logf("config data: %+v", cfg)

	if cfg.ApiKey != "demo" {
		t.Errorf("fail loading configuration: %s", pathCfg)
	}
}
