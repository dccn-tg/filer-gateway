package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Donders-Institute/filer-gateway/pkg/filer"
	tscfg "github.com/Donders-Institute/tg-toolset-golang/pkg/config"
	"github.com/spf13/viper"
)

// Configuration is the data structure for marshaling the
// config.yml file using the viper configuration framework.
type Configuration struct {
	NetApp  filer.NetAppConfig
	FreeNas filer.FreeNasConfig
	CephFs  filer.CephFsConfig
	Smtp    tscfg.SMTPConfiguration
	Pdb     tscfg.PDBConfiguration
}

// LoadConfig reads configuration file `cpath` and returns the
// `Configuration` data structure.
func LoadConfig(cpath string) (Configuration, error) {

	var conf Configuration

	// load configuration
	cfg, err := filepath.Abs(cpath)
	if err != nil {
		return conf, fmt.Errorf("cannot resolve config path: %s", cpath)
	}

	if _, err := os.Stat(cfg); err != nil {
		return conf, fmt.Errorf("cannot load config: %s", cfg)
	}

	viper.SetConfigFile(cfg)
	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return conf, nil
}
