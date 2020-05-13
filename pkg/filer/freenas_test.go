package filer

import (
	"os"
	"testing"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

var (
	freenas          Filer
	freenasProjectID string
)

func init() {
	freenasProjectID = "3010000.04"

	filerCfg := FreeNasConfig{
		ApiURL:           os.Getenv("FREENAS_API_SERVER"),
		ApiUser:          os.Getenv("FREENAS_API_USERNAME"),
		ApiPass:          os.Getenv("FREENAS_API_PASSWORD"),
		ProjectUser:      "project",
		ProjectGroup:     "project_g",
		ProjectRoot:      "/project_freenas",
		ZfsDatasetPrefix: "zpool001/project",
	}

	freenas = New("freenas", filerCfg)

	logCfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Debug,
	}

	// initialize logger
	log.NewLogger(logCfg, log.InstanceLogrusLogger)
}

func TestFreeNasGetProject(t *testing.T) {
	d, err := freenas.(FreeNas).getProjectDataset(freenasProjectID)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	t.Logf("dataset: %+v\n", d)
}

func TestFreeNasCreateProject(t *testing.T) {

	freenasProjectID = "3010000.05"

	if err := freenas.CreateProject(freenasProjectID, 10); err != nil {
		t.Errorf("%s\n", err)
	}

	// check if the dataset is actually created.
	d, err := freenas.(FreeNas).getProjectDataset(freenasProjectID)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	t.Logf("dataset: %+v\n", d)
}

func TestFreeNasSetProjectQuota(t *testing.T) {

	freenasProjectID = "3010000.05"

	if err := freenas.SetProjectQuota(freenasProjectID, 20); err != nil {
		t.Errorf("%s\n", err)
	}

	// check if the dataset quota is actually set to 20 GiB.
	d, err := freenas.(FreeNas).getProjectDataset(freenasProjectID)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	if d.RefQuota.Parsed>>30 != 20 {
		t.Errorf("quota not set to the targeting size.")
	}
	t.Logf("dataset: %+v\n", d)
}
