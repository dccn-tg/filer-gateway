// Package filer defines the interfaces for provisioning or updating
// a storage space on DCCN storage systems (a.k.a. filer) for a user
// or a project.
package filer

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
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

// Common function to check if given directory is empty (i.e. contains no files or sub-directories).
func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// New function returns the corresponding File implementation based on the
// `system` name.
func New(system string, config Config) Filer {
	switch system {
	case "netapp":
		return NetApp{config: config.(NetAppConfig)}
	case "freenas":
		return FreeNas{config: config.(FreeNasConfig)}
	case "cephfs":
		return CephFs{config: config.(CephFsConfig)}
	default:
		return nil
	}
}

// Config defines interfaces for retriving configuration parameters that are
// common across different filer systems.
type Config interface {
}

// Filer defines the interfaces for provisioning and setting storage space
// for a project and a personal home directory.
type Filer interface {
	CreateProject(projectID string, quotaGiB int) error
	CreateHome(username, groupname string, quotaGiB int) error
	DeleteHome(username, groupname string) error
	SetProjectQuota(projectID string, quotaGiB int) error
	SetHomeQuota(username, groupname string, quotaGiB int) error
	GetProjectQuotaInBytes(projectID string) (int64, int64, error)
	GetHomeQuotaInBytes(username, groupname string) (int64, int64, error)
	GetSystemSpaceInBytes() (int64, int64, error)
	GetProjectRoot() string
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
