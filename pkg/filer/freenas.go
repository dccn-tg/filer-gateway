package filer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

const (
	// apiNsFreenasDataset is the API namespace for FreeNAS ZFS datasets.
	apiNsFreenasDataset string = "/pool/dataset"

	// apiNsFreenasNfsShare is the API namespace for FreeNAS NFS sharing.
	apiNsFreenasNfsShare string = "/sharing/nfs"
)

// FreeNasConfig implements the `Config` interface and extends it with configurations
// that are specific to the FreeNas filer.
type FreeNasConfig struct {
	// ApiURL is the server URL of the OnTAP APIs.
	ApiURL string
	// ApiUser is the username for the basic authentication of the OnTAP API.
	ApiUser string
	// ApiPass is the password for the basic authentication of the OnTAP API.
	ApiPass string
	// ProjectRoot specifies the top-level NAS path in which projects are located.
	ProjectRoot string

	// ProjectUser specifies the system username for the owner of the project directory.
	ProjectUser string
	// ProjectGID specifies the system groupname for the owner of the project directory.
	ProjectGroup string

	// ZfsDatasetPrefix specifies the dataset prefix. It is usually started with the
	// zfs pool name followed by a top-level dataset name.  E.g. /zpool001/project.
	ZfsDatasetPrefix string
}

// GetAPIURL returns the server URL of the OnTAP API.
func (c FreeNasConfig) GetAPIURL() string {
	return strings.TrimSuffix(c.ApiURL, "/")
}

// GetAPIUser returns the username for the API basic authentication.
func (c FreeNasConfig) GetAPIUser() string { return c.ApiUser }

// GetAPIPass returns the password for the API basic authentication.
func (c FreeNasConfig) GetAPIPass() string { return c.ApiPass }

// FreeNas implements `Filer` for FreeNAS system.
type FreeNas struct {
	config FreeNasConfig
}

// GetProjectRoot returns the root path in which projects are hosted on the FreeNas system.
func (filer FreeNas) GetProjectRoot() string {
	return filer.config.ProjectRoot
}

// CreateProject creates a new dataset on the FreeNAS system with the dataset size
// specified by `quotaGiB`.
func (filer FreeNas) CreateProject(projectID string, quotaGiB int) error {

	// data structure of the POST data
	d := datasetUpdate{
		Name:            strings.Join([]string{filer.config.ZfsDatasetPrefix, projectID}, "/"),
		Comments:        fmt.Sprintf("project %s", projectID),
		RefQuota:        int64(quotaGiB << 30),
		RecordSize:      "128K",
		Type:            "FILESYSTEM",
		Sync:            "STANDARD",
		Compression:     "LZ4",
		Atime:           "ON",
		Exec:            "ON",
		ReadOnly:        "OFF",
		Deduplication:   "OFF",
		Copies:          1,
		RefReservation:  0,
		Reservation:     0,
		Snapdir:         "HIDDEN",
		ShareType:       "GENERIC",
		CaseSensitivity: "SENSITIVE",
	}

	if err := filer.createObject(&d, apiNsFreenasDataset); err != nil {
		return err
	}

	// set permission of the created dataset
	p := permission{
		User:  filer.config.ProjectUser,
		Group: filer.config.ProjectGroup,
		Mode:  "0750",
		Options: permissionOptions{
			Traverse:  false,
			Resursive: true,
			StripACL:  true,
		},
	}

	ns := strings.Join([]string{
		apiNsFreenasDataset,
		"id",
		filer.encodeProjectDatasetID(projectID),
		"permission",
	}, "/")

	if err := filer.createObject(&p, ns); err != nil {
		return err
	}

	// create NFS sharing of the project dataset
	s := nfs{
		AllDirs:      false,
		ReadOnly:     false,
		Quiet:        false,
		MapRootUser:  "root",
		MapRootGroup: filer.config.ProjectGroup,
		Security:     []string{"SYS"},
		Networks:     []string{"131.174.44.0/23"},
		Paths:        []string{filepath.Join("/mnt", filer.config.ZfsDatasetPrefix, projectID)},
	}

	if err := filer.createObject(&s, apiNsFreenasNfsShare); err != nil {
		return err
	}

	return nil
}

// CreateHome is not supported on FreeNAS and therefore it always returns an error.
func (filer FreeNas) CreateHome(username, groupname string, quotaGiB int) error {
	return fmt.Errorf("user home on FreeNAS is not supported")
}

// SetProjectQuota updates the size of the dataset for the specific dataset.
func (filer FreeNas) SetProjectQuota(projectID string, quotaGiB int) error {

	c, err := filer.getProjectDataset(projectID)
	if err != nil {
		return fmt.Errorf("cannot get dataset for project %s: %s", projectID, err)
	}
	if int(c.RefQuota.Parsed>>30) == quotaGiB {
		log.Warnf("quota of project %s already in right size: %d", projectID, quotaGiB)
		return nil
	}

	ns := strings.Join([]string{
		apiNsFreenasDataset,
		"id",
		filer.encodeProjectDatasetID(projectID),
	}, "/")

	d := datasetUpdate{
		RefQuota: int64(quotaGiB << 30),
	}
	if err := filer.updateObject(&d, ns); err != nil {
		return err
	}

	return nil
}

// SetHomeQuota is not supported on FreeNAS and therefore it always returns an error.
func (filer FreeNas) SetHomeQuota(username, groupname string, quotaGiB int) error {
	return fmt.Errorf("user home on FreeNAS is not supported")
}

// GetProjectQuotaInBytes returns the size of the dataset for a specific project in
// the unit of byte.
func (filer FreeNas) GetProjectQuotaInBytes(projectID string) (int64, int64, error) {

	d, err := filer.getProjectDataset(projectID)

	if err != nil {
		return 0, 0, fmt.Errorf("cannot get dataset for project %s: %s", projectID, err)
	}

	return d.RefQuota.Parsed, d.Used.Parsed, nil
}

// GetHomeQuotaInBytes is not supported on FreeNAS and therefore it always returns an error.
func (filer FreeNas) GetHomeQuotaInBytes(username, groupname string) (int64, int64, error) {
	return 0, 0, fmt.Errorf("user home on FreeNAS is not supported")
}

// encodeProjectDatasetID constructs the dataset ID from the given `projectID`, and returns an
// URL encoded id that can be used for API call.
func (filer FreeNas) encodeProjectDatasetID(projectID string) string {
	return url.PathEscape(strings.Join([]string{filer.config.ZfsDatasetPrefix, projectID}, "/"))
}

// getProjectDataset retrieves a structured dataset from the API.
func (filer FreeNas) getProjectDataset(projectID string) (*dataset, error) {

	d := dataset{}

	ns := strings.Join([]string{
		apiNsFreenasDataset,
		"id",
		filer.encodeProjectDatasetID(projectID),
	}, "/")

	if err := filer.getObject(ns, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (filer FreeNas) getObject(nsAPI string, object interface{}) error {

	c := newHTTPSClient(30*time.Second, true)

	filer.config.GetAPIURL()

	href := strings.Join([]string{filer.config.GetAPIURL(), nsAPI}, "")

	log.Debugf("href: %s", href)

	// create request
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return err
	}

	// set request header for basic authentication
	req.SetBasicAuth(filer.config.GetAPIUser(), filer.config.GetAPIPass())
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	// read response body
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("status: %s (%d), message: %s", res.Status, res.StatusCode, string(httpBodyBytes))

	// expect status to be 200 (OK)
	if res.StatusCode != 200 {
		return fmt.Errorf("response not ok: %s", msg)
	}

	// unmarshal response body to object structure
	if err := json.Unmarshal(httpBodyBytes, object); err != nil {
		return err
	}

	return nil
}

// createObject creates given object under the specified API namespace.
func (filer FreeNas) createObject(object interface{}, nsAPI string) error {
	c := newHTTPSClient(10*time.Second, true)

	href := strings.Join([]string{filer.config.GetAPIURL(), nsAPI}, "")

	data, err := json.Marshal(object)

	if err != nil {
		return fmt.Errorf("fail to convert to json data: %+v, %s", object, err)
	}

	log.Debugf("object creation input: %s", string(data))

	// create request
	req, err := http.NewRequest("POST", href, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// set request header for basic authentication
	req.SetBasicAuth(filer.config.GetAPIUser(), filer.config.GetAPIPass())
	req.Header.Set("content-type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	// read response body as accepted job
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %s", err)
	}

	msg := fmt.Sprintf("status: %s (%d), message: %s", res.Status, res.StatusCode, string(httpBodyBytes))
	log.Debugf("%s", msg)

	// expect status to be 200
	if res.StatusCode != 200 {
		return fmt.Errorf("response not ok: %s", msg)
	}
	return nil
}

func (filer FreeNas) updateObject(object interface{}, nsAPI string) error {
	c := newHTTPSClient(10*time.Second, true)

	href := strings.Join([]string{filer.config.GetAPIURL(), nsAPI}, "")

	data, err := json.Marshal(object)

	if err != nil {
		return fmt.Errorf("fail to convert to json data: %+v, %s", object, err)
	}

	log.Debugf("object creation input: %s", string(data))

	// create request
	req, err := http.NewRequest("PUT", href, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// set request header for basic authentication
	req.SetBasicAuth(filer.config.GetAPIUser(), filer.config.GetAPIPass())
	req.Header.Set("content-type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	// read response body as accepted job
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %s", err)
	}

	msg := fmt.Sprintf("status: %s (%d), message: %s", res.Status, res.StatusCode, string(httpBodyBytes))
	log.Debugf("%s", msg)

	// expect status to be 200
	if res.StatusCode != 200 {
		return fmt.Errorf("response not ok: %s", msg)
	}
	return nil
}

// dataset defines the JSON data structure of dataset retrieved from the API.
type dataset struct {
	ID          string     `json:"id"`
	Pool        string     `json:"pool"`
	Type        string     `json:"type"`
	SharedType  string     `json:"share_type"`
	Compression valueStr   `json:"compression"`
	RefQuota    valueInt64 `json:"refquota"`
	RecordSize  valueInt   `json:"recordsize"`
	Used        valueInt64 `json:"used"`
}

// valueStr defines general JSON structure of a string value retrieved from the API.
type valueStr struct {
	Value  string `json:"value,omitempty"`
	Parsed string `json:"parsed,omitempty"`
}

// valueInt64 defines general JSON structure of a int64 value retrieved from the API.
type valueInt64 struct {
	Value  string `json:"value,omitempty"`
	Parsed int64  `json:"parsed,omitempty"`
}

// valueInt defines general JSON structure of a int value retrieved from the API.
type valueInt struct {
	Value  string `json:"value,omitempty"`
	Parsed int    `json:"parsed,omitempty"`
}

// datasetUpdate defines the JSON data structure used to update a dataset.
type datasetUpdate struct {
	Name            string `json:"name,omitempty"`
	Type            string `json:"type,omitempty"`
	Sync            string `json:"sync,omitempty"`
	Comments        string `json:"comments,omitempty"`
	RefQuota        int64  `json:"refquota,omitempty"`
	Compression     string `json:"compression,omitempty"`
	Atime           string `json:"atime,omitempty"`
	Exec            string `json:"exec,omitempty"`
	Reservation     int    `json:"reservation,omitempty"`
	RefReservation  int    `json:"refreservation,omitempty"`
	Copies          int    `json:"copies,omitempty"`
	Snapdir         string `json:"snapdir,omitempty"`
	Deduplication   string `json:"deduplication,omitempty"`
	ReadOnly        string `json:"readonly,omitempty"`
	RecordSize      string `json:"recordsize,omitempty"`
	CaseSensitivity string `json:"casesensitivity,omitempty"`
	ShareType       string `json:"share_type,omitempty"`
}

// permission defines the JSON data structure for setting a dataset permission.
type permission struct {
	User    string            `json:"user"`
	Group   string            `json:"group"`
	Mode    string            `json:"mode"`
	Options permissionOptions `json:"options"`
}

type permissionOptions struct {
	Resursive bool `json:"recursive"`
	Traverse  bool `json:"traverse"`
	StripACL  bool `json:"stripacl"`
}

// {
// 	"alldirs": false,
// 	"ro": false,
// 	"quiet": false,
// 	"maproot_user": "root",
// 	"maproot_group": "project_g",
// 	"mapall_user": null,
// 	"mapall_group": null,
// 	"security": ["SYS"],
// 	"paths": [],
// 	"networks": [
// 	 "131.174.44.0/24",
// 	 "131.174.45.0/24"
// 	]
//   }
// nfs defines the JSON data structure for setting a NFS sharing.
type nfs struct {
	AllDirs      bool     `json:"alldirs"`
	ReadOnly     bool     `json:"ro"`
	Quiet        bool     `json:"quiet"`
	MapRootUser  string   `json:"maproot_user"`
	MapRootGroup string   `json:"maproot_group"`
	Security     []string `json:"security"`
	Paths        []string `json:"paths"`
	Networks     []string `json:"networks"`
}
