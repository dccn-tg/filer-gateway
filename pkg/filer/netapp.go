package filer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

const (
	// apiNsNetappSvms is the API namespace for OnTAP SVM items.
	apiNsNetappSvms string = "/svm/svms"
	// apiNsNetappJobs is the API namespace for OnTAP cluster job items.
	// apiNsNetappJobs string = "/cluster/jobs"
	// apiNsNetappVolumes is the API namespace for OnTAP volume items.
	apiNsNetappVolumes string = "/storage/volumes"
	// apiNsNetappAggregates is the API namespace for OnTAP aggregate items.
	apiNsNetappAggregates string = "/storage/aggregates"
	// apiNsNetappQtrees is the API namespace for OnTAP qtree items.
	apiNsNetappQtrees string = "/storage/qtrees"
	// apiNsNetappQuotaRules is the API namespace for OnTAP quota rule items.
	apiNsNetappQuotaRules string = "/storage/quota/rules"
	// apiNsNetappQuotaReport is the API namespace for OnTAP quota report.
	apiNsNetappQuotaReport string = "/storage/quota/reports"
)

// NetAppConfig implements the `Config` interface and extends it with configurations
// that are specific to the NetApp filer.
type NetAppConfig struct {
	// ApiURL is the server URL of the OnTAP APIs.
	ApiURL string
	// ApiUser is the username for the basic authentication of the OnTAP API.
	ApiUser string
	// ApiPass is the password for the basic authentication of the OnTAP API.
	ApiPass string
	// ProjectRoot specifies the top-level NAS path in which projects are located.
	ProjectRoot string

	// ProjectMode specifies how the project space is allocated. Valid modes are
	// "volume" and "qtree".
	ProjectMode string

	// VolumeProjectQtrees specifies the (FlexGroup) volume name in which project
	// qtrees are located.
	VolumeProjectQtrees string

	// Vserver specifies the name of OnTAP SVM on which the filer APIs will perform.
	Vserver string
	// ProjectUID specifies the system UID of user `project`
	ProjectUID int
	// ProjectGID specifies the system GID of group `project_g`
	ProjectGID int

	// ExportPolicyHome specifies the export policy name of the user home
	ExportPolicyHome string

	// ExportPolicyProject specifies the export policy name of the project
	ExportPolicyProject string
}

// GetAPIURL returns the server URL of the OnTAP API.
func (c NetAppConfig) GetAPIURL() string { return c.ApiURL }

// GetAPIUser returns the username for the API basic authentication.
func (c NetAppConfig) GetAPIUser() string { return c.ApiUser }

// GetAPIPass returns the password for the API basic authentication.
func (c NetAppConfig) GetAPIPass() string { return c.ApiPass }

// NetApp implements Filer interface for NetApp OnTAP cluster.
type NetApp struct {
	config NetAppConfig
}

// volName converts project identifier to the OnTAP volume name.
//
// e.g. 3010000.01 -> project_3010000_01
func (filer NetApp) volName(projectID string) string {
	return strings.Join([]string{
		"project",
		strings.ReplaceAll(projectID, ".", "_"),
	}, "_")
}

// GetProjectRoot returns the root path in which projects are hosted on the NetApp filer.
func (filer NetApp) GetProjectRoot() string {
	return filer.config.ProjectRoot
}

// CreateProject provisions a project space on the filer with the given quota.
func (filer NetApp) CreateProject(projectID string, quotaGiB int) error {

	switch filer.config.ProjectMode {
	case "volume":
		// check if volume with the same name doee not exist.
		qry := url.Values{}
		qry.Set("name", filer.volName(projectID))
		records, err := filer.getRecordsByQuery(qry, apiNsNetappVolumes)
		if err != nil {
			return fmt.Errorf("fail to check volume %s: %s", projectID, err)
		}
		if len(records) != 0 {
			return fmt.Errorf("project volume already exists: %s", projectID)
		}

		// determine which aggregate should be used for creating the new volume.
		quota := int64(quotaGiB << 30)
		svm := SVM{}
		if err := filer.getObjectByName(filer.config.Vserver, apiNsNetappSvms, &svm); err != nil {
			return fmt.Errorf("fail to get SVM %s: %s", filer.config.Vserver, err)
		}
		avail := int64(0)

		var theAggr *Aggregate
		for _, record := range svm.Aggregates {
			aggr := Aggregate{}
			href := strings.Join([]string{
				"/api",
				apiNsNetappAggregates,
				record.UUID,
			}, "/")
			if err := filer.getObjectByHref(href, &aggr); err != nil {
				log.Errorf("ignore aggregate %s: %s", record.Name, err)
			}
			if aggr.State == "online" && aggr.Space.BlockStorage.Available > avail && aggr.Space.BlockStorage.Available > quota {
				theAggr = &aggr
			}
		}

		if theAggr == nil {
			return fmt.Errorf("cannot find aggregate for creating volume")
		}
		log.Debugf("selected aggreate for project volume: %+v", *theAggr)

		// create project volume with given quota.
		vol := Volume{
			Name: filer.volName(projectID),
			Aggregates: []Record{
				{Name: theAggr.Name},
			},
			Size:  quota,
			Svm:   Record{Name: filer.config.Vserver},
			State: "online",
			Style: "flexvol",
			Type:  "rw",
			Nas: Nas{
				UID:             filer.config.ProjectUID,
				GID:             filer.config.ProjectGID,
				Path:            filepath.Join(filer.GetProjectRoot(), projectID),
				SecurityStyle:   "unix",
				UnixPermissions: 750,
				ExportPolicy:    ExportPolicy{Name: filer.config.ExportPolicyProject},
			},
			QoS: &QoS{
				Policy: QoSPolicy{MaxIOPS: 6000},
			},
			SnapshotPolicy: Record{Name: "none"},
			Space: &Space{
				Snapshot: &SnapshotConfig{ReservePercent: 0},
			},
			Autosize: &Autosize{Mode: "off"},
		}

		// blocking operation to create the volume.
		if err := filer.createObject(&vol, apiNsNetappVolumes); err != nil {
			return err
		}

	case "qtree":
		// blocking operation to create the qtree.
		if err := filer.createQtree(projectID, filer.config.VolumeProjectQtrees, 750, filer.config.ExportPolicyProject); err != nil {
			return err
		}

		// set owner of the project directory
		ppath := filepath.Join(filer.GetProjectRoot(), projectID)

		// wait for ppath to appear up to 5 minutes.
		t := time.Now()
		for {
			if _, err := os.Stat(ppath); os.IsNotExist(err) && time.Since(t) < 5*time.Minute {
				log.Debugf("waiting for path to become available: %s", ppath)
				time.Sleep(time.Second)
				continue
			}
			break
		}

		if err := os.Chown(ppath, filer.config.ProjectUID, filer.config.ProjectGID); err != nil {
			return fmt.Errorf("cannot set owner of %s: %s", ppath, err)
		}

		// set quota
		return filer.SetProjectQuota(projectID, quotaGiB)

	default:
		return fmt.Errorf("unsupported project mode: %s", filer.config.ProjectMode)
	}

	return nil
}

// CreateHome creates a home directory as qtree `username` under the volume `groupname`,
// and assigned the given `quotaGiB` to the qtree.
func (filer NetApp) CreateHome(username, groupname string, quotaGiB int) error {
	// blocking operation to create the qtree.
	if err := filer.createQtree(username, groupname, 700, filer.config.ExportPolicyHome); err != nil {
		return err
	}

	// set owner of the home directory to the uid/gid of user `username`, assuming the home directory
	// is `/home/groupname/username`.
	u, err := user.Lookup(username)
	if err != nil {
		return err
	}
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)
	hpath := filepath.Join("/home", groupname, username)

	// wait for hpath to appear up to 5 minutes.
	t := time.Now()
	for {
		if _, err := os.Stat(hpath); os.IsNotExist(err) && time.Since(t) < 5*time.Minute {
			log.Debugf("waiting for path to become available: %s", hpath)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	if err := os.Chown(hpath, uid, gid); err != nil {
		return fmt.Errorf("cannot set owner of %s: %s", hpath, err)
	}

	return nil
}

// SetProjectQuota updates the quota of a project space.
func (filer NetApp) SetProjectQuota(projectID string, quotaGiB int) error {

	qn, _, err := filer.GetProjectQuotaInBytes(projectID)
	if err != nil {
		return fmt.Errorf("cannot get current quota for project %s: %s", projectID, err)
	}

	if int(qn>>30) == quotaGiB {
		log.Warnf("quota of project %s already in right size: %d", projectID, quotaGiB)
		return nil
	}

	switch filer.config.ProjectMode {
	case "volume":
		// check if volume with the same name already exists.
		qry := url.Values{}
		qry.Set("name", filer.volName(projectID))
		records, err := filer.getRecordsByQuery(qry, apiNsNetappVolumes)
		if err != nil {
			return fmt.Errorf("fail to check volume %s: %s", projectID, err)
		}
		if len(records) != 1 {
			return fmt.Errorf("project volume doesn't exist: %s", projectID)
		}

		// resize the volume to the given quota.
		data := []byte(fmt.Sprintf(`{"name":"%s", "size":%d}`, filer.volName(projectID), quotaGiB<<30))

		if err := filer.patchObject(records[0], data); err != nil {
			return err
		}

	case "qtree":
		return filer.setQtreeQuota(projectID, filer.config.VolumeProjectQtrees, quotaGiB)

	default:
		return fmt.Errorf("unsupported project mode: %s", filer.config.ProjectMode)
	}

	return nil
}

// SetHomeQuota updates the quota of a home directory.
func (filer NetApp) SetHomeQuota(username, groupname string, quotaGiB int) error {
	return filer.setQtreeQuota(username, groupname, quotaGiB)
}

// GetProjectQuotaInBytes retrieves quota of a project in bytes.
func (filer NetApp) GetProjectQuotaInBytes(projectID string) (int64, int64, error) {
	switch filer.config.ProjectMode {
	case "volume":
		// check if volume with the same name already exists.
		vol := Volume{}

		if err := filer.getObjectByName(filer.volName(projectID), apiNsNetappVolumes, &vol); err != nil {
			return 0, 0, fmt.Errorf("cannot get project volume %s: %s", projectID, err)
		}

		return vol.Size, vol.Space.Used, nil

	case "qtree":

		r, err := filer.getQtreeQuotaReport(filer.config.VolumeProjectQtrees, projectID)

		if err != nil {
			return 0, 0, err
		}

		return r.Space.HardLimit, r.Space.Used.Total, nil

	default:
		return 0, 0, fmt.Errorf("unsupported project mode: %s", filer.config.ProjectMode)
	}
}

// GetHomeQuotaInBytes retrieves quota of a user home space in bytes.
func (filer NetApp) GetHomeQuotaInBytes(username, groupname string) (int64, int64, error) {

	r, err := filer.getQtreeQuotaReport(groupname, username)

	if err != nil {
		return 0, 0, err
	}

	return r.Space.HardLimit, r.Space.Used.Total, nil

	// qry := url.Values{}
	// qry.Set("volume.name", groupname)
	// qry.Set("qtree.name", username)

	// records, err := filer.getRecordsByQuery(qry, apiNsNetappQuotaRules)
	// if err != nil {
	// 	return 0, fmt.Errorf("fail to check quota rule for volume %s qtree %s: %s", groupname, username, err)
	// }

	// if len(records) == 0 { // no qtree specific quota rule, try to get the default rule
	// 	rule, err := filer.getDefaultQuotaRule(groupname)
	// 	if err != nil {
	// 		log.Errorf("cannot get default quota rule for volume %s", groupname)
	// 	}
	// 	if rule == nil {
	// 		return 0, fmt.Errorf("quota rule for volume %s qtree %s doesn't exist", groupname, username)
	// 	}
	// 	return rule.Space.HardLimit, nil
	// }

	// // with qtree specific quota rule
	// rule := QuotaRule{}
	// if err := filer.getObjectByHref(records[0].Link.Self.Href, &rule); err != nil {
	// 	return 0, fmt.Errorf("cannot get quota rule for volume %s qtree %s", groupname, username)
	// }

	// return rule.Space.HardLimit, nil
}

// createQtree implements the generic logic of creating a qtree in a volume, with given
// filesystem permission and export policy.
func (filer NetApp) createQtree(name, volume string, permission int, exportPolicy string) error {
	// check if qtree with "username" already exists.
	qry := url.Values{}
	qry.Set("name", name)
	qry.Set("volume.name", volume)
	records, err := filer.getRecordsByQuery(qry, apiNsNetappQtrees)
	if err != nil {
		return fmt.Errorf("fail to check qtree %s of volume %s: %s", name, volume, err)
	}
	if len(records) != 0 {
		return fmt.Errorf("qtree %s of volume %s already exists", name, volume)
	}

	// create qtree within the volume.
	qtree := QTree{
		Name:            name,
		SVM:             Record{Name: filer.config.Vserver},
		Volume:          Record{Name: volume},
		SecurityStyle:   "unix",
		UnixPermissions: permission,
		ExportPolicy:    ExportPolicy{Name: exportPolicy},
	}

	// blocking operation to create the qtree.
	if err := filer.createObject(&qtree, apiNsNetappQtrees); err != nil {
		return err
	}

	return nil
}

// setQtreeQuota implements the generic logic of setting quota rule on a given volume.
func (filer NetApp) setQtreeQuota(name, volume string, quotaGiB int) error {

	// check if the qtree exists
	qry := url.Values{}
	qry.Set("name", name)
	qry.Set("volume.name", volume)
	recQtrees, err := filer.getRecordsByQuery(qry, apiNsNetappQtrees)
	if err != nil {
		return fmt.Errorf("fail to check qtree %s of volume %s: %s", name, volume, err)
	}
	if len(recQtrees) == 0 {
		return fmt.Errorf("qtree %s of volume %s doesn't exit", name, volume)
	}

	// get volume record from the qtree
	qtree := QTree{}
	if err := filer.getObjectByHref(recQtrees[0].Link.Self.Href, &qtree); err != nil {
		return fmt.Errorf("fail to retrieve volume %s: %s", volume, err)
	}
	volRecord := qtree.Volume

	// check if the user-specific quota rule exists
	qry = url.Values{}
	qry.Set("volume.name", volume)
	qry.Set("qtree.name", name)
	recRules, err := filer.getRecordsByQuery(qry, apiNsNetappQuotaRules)
	if err != nil {
		return fmt.Errorf("fail to check quota rule for volume %s qtree %s: %s", volume, name, err)
	}

	// unexpected number of quota rules for the specific volume and qtree.
	if len(recRules) > 1 {
		return fmt.Errorf("more than one quota rule for volume %s qtree %s (%d)", volume, name, len(recRules))
	}

	// get qtree specific rule
	if len(recRules) == 1 {
		qtreeRule := QuotaRule{}
		if err := filer.getObjectByHref(recRules[0].Link.Self.Href, &qtreeRule); err != nil {
			return fmt.Errorf("fail to get quota rule for volume %s qtree %s", volume, name)
		}

		if quotaGiB == int(qtreeRule.Space.HardLimit>>30) {
			log.Debugf("quota target is already the current quota limitation, volume %s qtree %s", volume, name)
			return nil
		}
	}

	// get the default quota on volume level.
	volRule, err := filer.getDefaultQuotaRule(volume)
	if err != nil {
		return err
	}

	// default quota rule is available and the hard limit is identical to the quota target.
	if volRule != nil && quotaGiB == int(volRule.Space.HardLimit>>30) {
		log.Debugf("quota target is already the default quota limitation, volume %s qtree %s", volume, name)
		// already a default rule, should remove the user specific rule if it exists.
		if len(recRules) == 1 {
			log.Debugf("remove specific quota policy for qtree %s, volume %s", name, volume)
			if err := filer.delObjectByHref(recRules[0].Link.Self.Href); err != nil {
				return fmt.Errorf("cannot delete user-specific quota rule for %s volume %s: %s", name, volume, err)
			}
		}
		return nil
	}

	// switch off and on the volume quota is needed if there is no default quota rule applied on the volume.
	// TODO: is it really necessary???
	if volRule == nil {
		// switch off volume quota
		log.Debugf("turn off quota on volume %s", volume)
		if err := filer.patchObject(volRecord, []byte(`{"quota":{"enabled":false}}`)); err != nil {
			return err
		}

		// ensure the volume quota will be switched on before this function is returned.
		defer func() {
			log.Debugf("turn on quota on volume %s", volume)
			if err := filer.patchObject(volRecord, []byte(`{"quota":{"enabled":true}}`)); err != nil {
				log.Errorf("cannot turn on quota for volume %s: %s", volume, err)
			}
		}()
	}

	if len(recRules) == 0 {
		// create new user-specific quota rule.
		qrule := QuotaRule{
			SVM:    Record{Name: filer.config.Vserver},
			Volume: Record{Name: volume},
			QTree:  &Record{Name: name},
			Type:   "tree",
			Space:  &QuotaLimit{HardLimit: int64(quotaGiB << 30)},
		}
		if err := filer.createObject(&qrule, apiNsNetappQuotaRules); err != nil {
			return err
		}

		// retrieve the successfully created new rule
		if recRules, err = filer.getRecordsByQuery(qry, apiNsNetappQuotaRules); err != nil {
			return err
		}
	}

	// update corresponding quota rule for the qtree
	// NOTE: this is a redundent call in case of creating a new quota rule. But it might
	//       ensure the new quota is always applied.
	if len(recRules) > 0 {
		data := []byte(fmt.Sprintf(`{"space":{"hard_limit":%d}}`, quotaGiB<<30))
		if err := filer.patchObject(recRules[0], data); err != nil {
			return err
		}
	}

	return nil
}

// getDefaultQuotaRule returns the default quota rule on a volume as `QuotaRule`.
func (filer NetApp) getDefaultQuotaRule(volume string) (*QuotaRule, error) {

	var rule QuotaRule

	qry := url.Values{}
	qry.Set("volume.name", volume)

	records, err := filer.getRecordsByQuery(qry, apiNsNetappQuotaRules)
	if err != nil {
		return &rule, fmt.Errorf("fail to check quota rule for volume %s: %s", volume, err)
	}
	if len(records) == 0 {
		return &rule, nil
	}

	for _, rec := range records {
		err := filer.getObjectByHref(rec.Link.Self.Href, &rule)
		if err != nil {
			log.Errorf("cannot retrieve quota rule, %s: %s", rec.Link.Self.Href, err)
			continue
		}
		if rule.QTree.Name == "" {
			return &rule, nil
		}
	}

	return nil, nil
}

// getQtreeQuotaReport returns the quota assignment and usage for a specific qtree in a volume.
func (filer NetApp) getQtreeQuotaReport(volume string, qtree string) (report *QuotaReport, err error) {

	qry := url.Values{}
	qry.Set("volume.name", volume)
	qry.Set("type", "tree")
	qry.Set("qtree.name", qtree)

	records, err := filer.getRecordsByQuery(qry, apiNsNetappQuotaReport)

	if err != nil {
		return
	}

	if len(records) == 0 {
		err = fmt.Errorf("quota report not available for volume %s, qtree %s", volume, qtree)
		return
	}

	err = filer.getObjectByHref(records[0].Link.Self.Href, &report)

	return
}

// GetVolumeQuotaReports returns a list of `QuotaReport`, each refers to the quota limit and usage
// per qtree.
func (filer NetApp) GetVolumeQuotaReports(volume string) (reports []*QuotaReport, err error) {

	qry := url.Values{}
	qry.Set("volume.name", volume)
	qry.Set("type", "tree")
	qry.Set("fields", "svm,volume,qtree,space")

	var r VolumeQuotaReport

	err = filer.getObjectByQuery(qry, apiNsNetappQuotaReport, &r)

	if err != nil {
		err = fmt.Errorf("fail to check quota report for volume %s: %s", volume, err)
		return
	}

	reports = r.Records
	return
}

// getObjectByName retrives the named object from the given API namespace.
func (filer NetApp) getObjectByName(name, nsAPI string, object interface{}) error {

	query := url.Values{}
	query.Set("name", name)

	records, err := filer.getRecordsByQuery(query, nsAPI)
	if err != nil {
		return err
	}

	if len(records) != 1 {
		return fmt.Errorf("more than 1 object found: %d", len(records))
	}

	if err := filer.getObjectByHref(records[0].Link.Self.Href, object); err != nil {
		return err
	}

	return nil
}

// getRecordsByQuery retrives the object from the given API namespace using a specific URL query.
func (filer NetApp) getRecordsByQuery(query url.Values, nsAPI string) ([]Record, error) {

	records := make([]Record, 0)

	c := newHTTPSClient(30*time.Second, true)

	href := strings.Join([]string{filer.config.GetAPIURL(), "api", nsAPI}, "/")

	// create request
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return records, err
	}

	// guarantee that record's _links field is returned
	query.Set("fields", "_links")

	req.URL.RawQuery = query.Encode()

	// set request header for basic authentication
	req.SetBasicAuth(filer.config.GetAPIUser(), filer.config.GetAPIPass())
	req.Header.Set("accept", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return records, err
	}

	// expect status to be 200 (OK)
	if res.StatusCode != 200 {
		return records, fmt.Errorf("response not ok: %s (%d)", res.Status, res.StatusCode)
	}

	// read response body
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return records, err
	}

	// unmarshal response body to object structure
	rec := Records{}
	if err := json.Unmarshal(httpBodyBytes, &rec); err != nil {
		return records, err
	}

	return rec.Records, nil
}

// delObjectByHref deletes the object at the given API namespace `href`.
func (filer NetApp) delObjectByHref(href string) error {
	c := newHTTPSClient(10*time.Second, true)

	// create request
	req, err := http.NewRequest("DELETE", strings.Join([]string{filer.config.GetAPIURL(), href}, "/"), nil)
	if err != nil {
		return err
	}

	// set request header for basic authentication
	req.SetBasicAuth(filer.config.GetAPIUser(), filer.config.GetAPIPass())

	res, _ := c.Do(req)

	// expect status to be 202 (Accepted)
	if res.StatusCode != 202 {
		// try to get the error code returned as the body
		var apiErr APIError
		if httpBodyBytes, err := ioutil.ReadAll(res.Body); err == nil {
			json.Unmarshal(httpBodyBytes, &apiErr)
		}
		return fmt.Errorf("response not ok: %s (%d), error: %+v", res.Status, res.StatusCode, apiErr)
	}

	// read response body as accepted job
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %s", err)
	}

	job := APIJob{}
	// unmarshal response body to object structure
	if err := json.Unmarshal(httpBodyBytes, &job); err != nil {
		return fmt.Errorf("cannot get job id: %s", err)
	}

	log.Debugf("job data: %+v", job)

	if err := filer.waitJob(&job); err != nil {
		return err
	}

	if job.Job.State != "success" {
		return fmt.Errorf("API job failed: %s", job.Job.Message)
	}

	return nil
}

// getObjectByHref retrives the object from the given API namespace `href`.
func (filer NetApp) getObjectByHref(href string, object interface{}) error {

	c := newHTTPSClient(10*time.Second, true)

	// create request
	req, err := http.NewRequest("GET", strings.Join([]string{filer.config.GetAPIURL(), href}, "/"), nil)
	if err != nil {
		return err
	}

	// set request header for basic authentication
	req.SetBasicAuth(filer.config.GetAPIUser(), filer.config.GetAPIPass())
	// NOTE: adding "Accept: application/json" to header can causes the API server
	//       to not returning "_links" attribute containing API href to the object.
	//       Therefore, it is not set here.
	//req.Header.Set("accept", "application/json")

	res, err := c.Do(req)

	if err != nil {
		return fmt.Errorf("%s request error: %s", req.URL, err)
	}

	// expect status to be 200 (OK)
	if res.StatusCode != 200 {
		return fmt.Errorf("response not ok: %s (%d)", res.Status, res.StatusCode)
	}

	// read response body
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// unmarshal response body to object structure
	if err := json.Unmarshal(httpBodyBytes, object); err != nil {
		return err
	}

	return nil
}

// getObjectByQuery retrives the object from the given API namespace with query variables.
func (filer NetApp) getObjectByQuery(query url.Values, nsAPI string, object interface{}) error {

	c := newHTTPSClient(10*time.Second, true)

	href := strings.Join([]string{filer.config.GetAPIURL(), "api", nsAPI}, "/")

	// create request
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return err
	}

	// set request header for basic authentication
	req.SetBasicAuth(filer.config.GetAPIUser(), filer.config.GetAPIPass())
	// NOTE: adding "Accept: application/json" to header can causes the API server
	//       to not returning "_links" attribute containing API href to the object.
	//       Therefore, it is not set here.
	//req.Header.Set("accept", "application/json")

	req.URL.RawQuery = query.Encode()

	res, err := c.Do(req)

	if err != nil {
		return fmt.Errorf("%s request error: %s", req.URL, err)
	}

	// expect status to be 200 (OK)
	if res.StatusCode != 200 {
		return fmt.Errorf("response not ok: %s (%d)", res.Status, res.StatusCode)
	}

	// read response body
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// unmarshal response body to object structure
	if err := json.Unmarshal(httpBodyBytes, object); err != nil {
		return err
	}

	return nil
}

// createObject creates given object under the specified API namespace.
func (filer NetApp) createObject(object interface{}, nsAPI string) error {
	c := newHTTPSClient(10*time.Second, true)

	href := strings.Join([]string{filer.config.GetAPIURL(), "api", nsAPI}, "/")

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

	// expect status to be 202 (Accepted)
	if res.StatusCode != 202 {
		// try to get the error code returned as the body
		var apiErr APIError
		if httpBodyBytes, err := ioutil.ReadAll(res.Body); err == nil {
			json.Unmarshal(httpBodyBytes, &apiErr)
		}
		return fmt.Errorf("response not ok: %s (%d), error: %+v", res.Status, res.StatusCode, apiErr)
	}

	// read response body as accepted job
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %s", err)
	}

	job := APIJob{}
	// unmarshal response body to object structure
	if err := json.Unmarshal(httpBodyBytes, &job); err != nil {
		return fmt.Errorf("cannot get job id: %s", err)
	}

	log.Debugf("job data: %+v", job)

	if err := filer.waitJob(&job); err != nil {
		return err
	}

	if job.Job.State != "success" {
		return fmt.Errorf("API job failed: %s", job.Job.Message)
	}

	return nil
}

// patchObject patches given object `Record` with provided setting specified by `data`.
func (filer NetApp) patchObject(object Record, data []byte) error {

	c := newHTTPSClient(10*time.Second, true)

	href := strings.Join([]string{filer.config.GetAPIURL(), object.Link.Self.Href}, "/")

	// create request
	req, err := http.NewRequest("PATCH", href, bytes.NewBuffer(data))
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

	// expect status to be 202 (Accepted)
	if res.StatusCode != 202 {
		// try to get the error code returned as the body
		var apiErr APIError
		if httpBodyBytes, err := ioutil.ReadAll(res.Body); err == nil {
			json.Unmarshal(httpBodyBytes, &apiErr)
		}
		return fmt.Errorf("response not ok: %s (%d), error: %+v", res.Status, res.StatusCode, apiErr)
	}

	// read response body as accepted job
	httpBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %s", err)
	}

	job := APIJob{}
	// unmarshal response body to object structure
	if err := json.Unmarshal(httpBodyBytes, &job); err != nil {
		return fmt.Errorf("cannot get job id: %s", err)
	}

	log.Debugf("job data: %+v", job)

	if err := filer.waitJob(&job); err != nil {
		return err
	}

	if job.Job.State != "success" {
		return fmt.Errorf("API job failed: %s", job.Job.Message)
	}

	return nil
}

// waitJob polls the status of the api job unti it if finished; and reports the job's final state.
func (filer NetApp) waitJob(job *APIJob) error {

	var err error

	href := job.Job.Link.Self.Href

waitLoop:
	for {
		if e := filer.getObjectByHref(href, &(job.Job)); err != nil {
			err = fmt.Errorf("cannot poll job %s: %s", job.Job.UUID, e)
			break
		}

		log.Debugf("job status: %s", job.Job.State)

		switch job.Job.State {
		case "success":
			break waitLoop
		case "failure":
			break waitLoop
		default:
			time.Sleep(3 * time.Second)
			continue waitLoop
		}
	}

	return err
}

// APIJob of the API request.
type APIJob struct {
	Job Job `json:"job"`
}

// Job detail of the API request.
type Job struct {
	Link    *Link  `json:"_links"`
	UUID    string `json:"uuid"`
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// APIError of the API request.
type APIError struct {
	Error struct {
		Target    string `json:"target"`
		Arguments struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"arguments"`
	} `json:"error"`
}

// Records of the items within an API namespace.
type Records struct {
	NumberOfRecords int      `json:"num_records"`
	Records         []Record `json:"records"`
}

// Record of an item within an API namespace.
type Record struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
	Link *Link  `json:"_links,omitempty"`
}

// Link of an item for retriving the detail.
type Link struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

// Volume of OnTAP.
type Volume struct {
	UUID           string    `json:"uuid,omitempty"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	State          string    `json:"state"`
	Size           int64     `json:"size"`
	Style          string    `json:"style"`
	Space          *Space    `json:"space,omitempty"`
	Svm            Record    `json:"svm"`
	Aggregates     []Record  `json:"aggregates"`
	Nas            Nas       `json:"nas"`
	SnapshotPolicy Record    `json:"snapshot_policy"`
	QoS            *QoS      `json:"qos,omitempty"`
	Autosize       *Autosize `json:"autosize,omitempty"`
	Link           *Link     `json:"_links,omitempty"`
}

// QoS contains a Qolity-of-Service policy.
type QoS struct {
	Policy QoSPolicy `json:"policy"`
}

// QoSPolicy defines the data structure of the QoS policy.
type QoSPolicy struct {
	MaxIOPS int    `json:"max_throughput_iops,omitempty"`
	MaxMBPS int    `json:"max_throughput_mbps,omitempty"`
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name,omitempty"`
}

// Autosize defines the volume autosizing mode
type Autosize struct {
	Mode string `json:"mode"`
}

// Nas related attribute of OnTAP.
type Nas struct {
	Path            string       `json:"path,omitempty"`
	UID             int          `json:"uid,omitempty"`
	GID             int          `json:"gid,omitempty"`
	SecurityStyle   string       `json:"security_style,omitempty"`
	UnixPermissions int          `json:"unix_permissions,omitempty"`
	ExportPolicy    ExportPolicy `json:"export_policy,omitempty"`
}

// ExportPolicy defines the export policy for a volume or a qtree.
type ExportPolicy struct {
	Name string `json:"name"`
}

// Space information of a OnTAP volume.
type Space struct {
	Size      int64           `json:"size,omitempty"`
	Available int64           `json:"available,omitempty"`
	Used      int64           `json:"used,omitempty"`
	Snapshot  *SnapshotConfig `json:"snapshot,omitempty"`
}

// SnapshotConfig of a OnTAP volume.
type SnapshotConfig struct {
	ReservePercent int `json:"reserve_percent"`
}

// SVM of OnTAP
type SVM struct {
	UUID       string   `json:"uuid"`
	Name       string   `json:"name"`
	State      string   `json:"state"`
	Aggregates []Record `json:"aggregates"`
}

// Aggregate of OnTAP
type Aggregate struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	State string `json:"state"`
	Space struct {
		BlockStorage Space `json:"block_storage"`
	} `json:"space"`
}

// QTree of OnTAP
type QTree struct {
	ID              int          `json:"id,omitempty"`
	Name            string       `json:"name"`
	Path            string       `json:"path,omitempty"`
	SVM             Record       `json:"svm"`
	Volume          Record       `json:"volume"`
	ExportPolicy    ExportPolicy `json:"export_policy"`
	SecurityStyle   string       `json:"security_style"`
	UnixPermissions int          `json:"unix_permissions"`
	Link            *Link        `json:"_links,omitempty"`
}

// QuotaRule of OnTAP
type QuotaRule struct {
	SVM    Record      `json:"svm"`
	Volume Record      `json:"volume"`
	QTree  *Record     `json:"qtree,omitempty"`
	Users  *Record     `json:"users,omitempty"`
	Group  *Record     `json:"group,omitempty"`
	Type   string      `json:"type"`
	Space  *QuotaLimit `json:"space,omitempty"`
	Files  *QuotaLimit `json:"files,omitempty"`
}

type QuotaReport struct {
	SVM    Record      `json:"svm"`
	Volume Record      `json:"volume"`
	QTree  *Record     `json:"qtree,omitempty"`
	Space  *QuotaUsage `json:"space,omitempty"`
}

type VolumeQuotaReport struct {
	Records []*QuotaReport `json:"records"`
}

// QuotaLimit defines the quota limitation.
type QuotaLimit struct {
	HardLimit int64 `json:"hard_limit,omitempty"`
}

// QuotaUsage defines the quota limitation.
type QuotaUsage struct {
	HardLimit int64 `json:"hard_limit,omitempty"`
	Used      struct {
		Total int64 `json:total`
	} `json:"used,omitempty"`
}
