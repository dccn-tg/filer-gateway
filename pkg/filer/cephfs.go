package filer

import "fmt"

// CephFsConfig implements the `Config` interface and extends it with configurations
// that are specific to the CephFS filesystem.
type CephFsConfig struct {

	// ProjectRoot specifies the top-level CephFS path in which projects are located.
	ProjectRoot string

	// ProjectUser specifies the system username of the project directory owner.
	ProjectUser string

	// ProjectGroup specifies the system groupname of the project directory owner.
	ProjectGroup string
}

// GetAPIURL is a dummy implementation fo the Config interface given that
// operations on CephFS doesn't use the API interface.
func (c CephFsConfig) GetAPIURL() string {
	return ""
}

// GetAPIUser is a dummy implementation fo the Config interface given that
// operations on CephFS doesn't use the API interface.
func (c CephFsConfig) GetAPIUser() string { return "" }

// GetAPIPass is a dummy implementation fo the Config interface given that
// operations on CephFS doesn't use the API interface.
func (c CephFsConfig) GetAPIPass() string { return "" }

// GetProjectRoot returns the filesystem root path in which directories of projects are located.
func (c CephFsConfig) GetProjectRoot() string { return c.ProjectRoot }

// CephFs implement the filer interface specific for the Ceph filesystem.
type CephFs struct {
	config CephFsConfig
}

// CreateProject creates a new project directory on the Ceph filesystem mounted under
// `CephFsConfig.GetProjectRoot()`.
func (c CephFs) CreateProject(projectID string, quotaGiB int) error {
	return fmt.Errorf("not implemented")
}

// SetProjectQuota sets or updates the project quota by setting a extended attribute
// `cephfs.quota.max_bytes` of the project directory on the Ceph filesystem.
//
// See [here](https://docs.ceph.com/docs/master/cephfs/quota/) for more detail.
func (c CephFs) SetProjectQuota(projectID string, quotaGiB int) error {
	return fmt.Errorf("not implemented")
}

// GetProjectQuotaInBytes retrieves the value of the extended attribute `cephfs.quota.max_bytes`
// from the project directory on the Ceph filesystem.
//
// See [here](https://docs.ceph.com/docs/master/cephfs/quota/) for more detail.
func (c CephFs) GetProjectQuotaInBytes(projectID string) (int64, error) {
	return -1, fmt.Errorf("not implemented")
}

// CreateHome always returns an error with "not supported" message, given that
// Ceph filesystem is not used for personal home directory.
func (c CephFs) CreateHome(username, groupname string, quotaGiB int) error {
	return fmt.Errorf("not supported")
}

// SetHomeQuota always returns an error with "not supported" message, given that
// Ceph filesystem is not used for personal home directory.
func (c CephFs) SetHomeQuota(username, groupname string, quotaGiB int) error {
	return fmt.Errorf("not supported")
}

// GetHomeQuotaInBytes always returns an error with "not supported" message, given that
// Ceph filesystem is not used for personal home directory.
func (c CephFs) GetHomeQuotaInBytes(username, groupname string) (int64, error) {
	return -1, fmt.Errorf("not supported")
}
