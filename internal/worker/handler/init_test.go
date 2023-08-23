package handler

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/dccn-tg/filer-gateway/pkg/filer"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestGetFilerAPI(t *testing.T) {

	pathCfg := os.Getenv("FILER_GATEWAY_WORKER_CONFIG")

	tobj := map[string]filer.Filer{
		"netapp":  filer.NetApp{},
		"freenas": filer.FreeNas{},
	}

	for _, sys := range []string{"netapp", "freenas"} {
		api, err := getFilerAPIBySystem(sys, pathCfg)
		if err != nil {
			t.Errorf("%s", err)
		}

		assertEqual(t, reflect.TypeOf(api), reflect.TypeOf(tobj[sys]), "")
	}
}

func TestNotifyProjectProvisioned(t *testing.T) {

	pathCfg := os.Getenv("FILER_GATEWAY_WORKER_CONFIG")

	h := SetProjectResourceHandler{
		ConfigFile: pathCfg,
	}

	managers := []string{"honlee"}
	projectID := "3010000.01"

	err := h.notifyProjectProvisioned(projectID, managers)
	if err != nil {
		t.Errorf("%s", err)
	}
}
