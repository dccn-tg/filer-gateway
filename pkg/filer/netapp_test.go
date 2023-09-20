package filer

import (
	"encoding/json"
	"os"
	"testing"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	netapp          Filer
	netappProjectID string
)

const (
	groupname string = "tg"
	username  string = "honlee"
)

func init() {

	netappProjectID = "2420084.01"

	filerCfg := NetAppConfig{
		ApiURL:              os.Getenv("NETAPP_API_SERVER"),
		ApiUser:             os.Getenv("NETAPP_API_USERNAME"),
		ApiPass:             os.Getenv("NETAPP_API_PASSWORD"),
		Vserver:             os.Getenv("NETAPP_VSERVER"),
		ProjectGID:          1010,
		ProjectUID:          1010,
		ProjectRoot:         "/project",
		ProjectMode:         os.Getenv("NETAPP_PROJECT_MODE"),
		VolumeProjectQtrees: "project",
		ExportPolicyHome:    os.Getenv("NETAPP_EXPORT_POLICY_HOME"),
		ExportPolicyProject: os.Getenv("NETAPP_EXPORT_POLICY_PROJECT"),
	}

	netapp = New("netapp", filerCfg)

	logCfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Debug,
	}

	// initialize logger
	log.NewLogger(logCfg, log.InstanceLogrusLogger)
}

func TestUnmarshal(t *testing.T) {
	data := []byte(`{
			"uuid": "27c77b57-a06c-4af5-8c15-1c625e628f64",
			"name": "tg",
			"_links": {
			  "self": {
				"href": "/api/storage/volumes/27c77b57-a06c-4af5-8c15-1c625e628f64"
			  }
			}
		  }`)

	record := Record{}

	json.Unmarshal(data, &record)

	t.Logf("%+v", record)

	records := Records{}
	data = []byte(`{
		"records": [
		  {
			"uuid": "27c77b57-a06c-4af5-8c15-1c625e628f64",
			"name": "tg",
			"_links": {
			  "self": {
				"href": "/api/storage/volumes/27c77b57-a06c-4af5-8c15-1c625e628f64"
			  }
			}
		  }
		],
		"num_records": 1,
		"_links": {
		  "self": {
			"href": "/api/storage/volumes?name=tg"
		  }
		}
	  }`)

	json.Unmarshal(data, &records)
	t.Logf("%+v", records)
}

func TestGetDefaultQuotaRule(t *testing.T) {

	r, err := netapp.(NetApp).getDefaultQuotaRule(groupname)

	if err != nil {
		t.Errorf("fail to get default quota rule: %s\n", err)
	}

	if r == nil {
		t.Errorf("default quota rule for %s not found\n", groupname)
	} else {
		t.Logf("quota rule for %s: %+v\n", groupname, r.Space)
	}
}

func TestGetSystemSpaceInBytes(t *testing.T) {
	if total, used, err := netapp.GetSystemSpaceInBytes(); err != nil {
		t.Errorf("fail to get system space: %s", err)
	} else {
		t.Logf("total: %d, used: %d", total, used)
	}
}

func TestCreateProject(t *testing.T) {
	if err := netapp.CreateProject(netappProjectID, 10); err != nil {
		t.Errorf("fail to create project volume: %s", err)
	}
}

func TestSetProjectQuota(t *testing.T) {
	if err := netapp.SetProjectQuota(netappProjectID, 20); err != nil {
		t.Errorf("fail to update quota for project %s: %s", netappProjectID, err)
	}
}

func TestCreateHome(t *testing.T) {
	if err := netapp.CreateHome(username, groupname, 10); err != nil {
		t.Errorf("%s\n", err)
	}
}

func TestSetHomeQuota(t *testing.T) {
	if err := netapp.SetHomeQuota(username, groupname, 50); err != nil {
		t.Errorf("%s\n", err)
	}
}

func TestGetHomeQuota(t *testing.T) {
	if quota, usage, err := netapp.GetHomeQuotaInBytes(username, groupname); err != nil {
		t.Errorf("%s\n", err)
	} else {
		t.Logf("quota: %d, usage: %d\n", quota, usage)
	}
}

func TestGetProjectQuota(t *testing.T) {
	if quota, usage, err := netapp.GetProjectQuotaInBytes(netappProjectID); err != nil {
		t.Errorf("%s\n", err)
	} else {
		t.Logf("quota: %d, usage: %d\n", quota, usage)
	}
}

func TestGetQuotaReport(t *testing.T) {
	if reports, err := netapp.(NetApp).GetVolumeQuotaReports("project"); err != nil {
		t.Errorf("%s\n", err)
	} else {
		t.Logf("%d reports, first: %+v\n", len(reports), reports[0])
	}
}
