package handler

import (
	"encoding/json"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/task"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/bokchoy"
)

// SetProjectResourceHandler implements `bokchoy.Handler` for applying update on project resource.
type SetProjectResourceHandler struct {
}

// Handle performs project resource update based on the request payload.
func (h *SetProjectResourceHandler) Handle(r *bokchoy.Request) error {

	log.Debugf("payload: %+v", r.Task.Payload)

	res, err := json.Marshal(r.Task.Payload)
	if err != nil {
		log.Errorf("Marshal error: %s", err)
		return err
	}

	var data task.SetProjectResource

	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Errorf("Unmarshal error: %s", err)
		return err
	}

	log.Debugf("payload data: %+v", data)

	// sleep for 5 minutes to test if on-demand timeout works
	time.Sleep(5 * time.Minute)

	// Put payload data as result.
	r.Task.Result = &data

	return nil
}

// SetUserResourceHandler implements `bokchoy.Handler` for applying update on user resource.
type SetUserResourceHandler struct {
}

// Handle performs user resource update based on the request payload.
func (h *SetUserResourceHandler) Handle(r *bokchoy.Request) error {
	return nil
}
