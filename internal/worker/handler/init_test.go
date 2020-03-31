package handler

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/filer"
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
		api, err := getFilerAPI(sys, pathCfg)
		if err != nil {
			t.Errorf("%s", err)
		}

		assertEqual(t, reflect.TypeOf(api), reflect.TypeOf(tobj[sys]), "")
	}
}
