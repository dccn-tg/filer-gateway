// Package filer defines the interfaces for provisioning or updating
// a storage space on DCCN storage systems (a.k.a. filer) for a user
// or a project.
package filer

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

func init() {

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
	}

	// initialize logger
	log.NewLogger(cfg, log.InstanceLogrusLogger)
}

// New function returns the corresponding File implementation based on the
// `system` name.
func New(system string, config Config) Filer {
	switch system {
	case "netapp":
		return NetApp{config: config.(NetAppConfig)}
	case "freenas":
		return FreeNas{config: config.(FreeNasConfig)}
	default:
		return nil
	}
}

// Config defines interfaces for retriving configuration parameters that are
// common across different filer systems.
type Config interface {
	GetApiURL() string
	GetApiUser() string
	GetApiPass() string
	GetProjectRoot() string
}

// Filer defines the interfaces for provisioning and setting storage space
// for a project and a personal home directory.
type Filer interface {
	CreateProject(projectID string, quotaGiB int) error
	CreateHome(username, groupname string, quotaGiB int) error
	SetProjectQuota(projectID string, quotaGiB int) error
	SetHomeQuota(username, groupname string, quotaGiB int) error
	GetProjectQuotaInBytes(projectID string) (int64, error)
	GetHomeQuotaInBytes(username, groupname string) (int64, error)
}

// newHTTPSClient initiate a HTTPS client.
func newHTTPSClient(timeout time.Duration, insecure bool) (client *http.Client) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
	}

	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client = &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	return
}
