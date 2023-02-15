package filer

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

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

// CephFs implement the filer interface specific for the Ceph filesystem.
type CephFs struct {
	config CephFsConfig
}

// GetProjectRoot returns the root path in which projects are hosted on the CephFS system.
func (filer CephFs) GetProjectRoot() string {
	return filer.config.ProjectRoot
}

// CreateProject creates a new project directory on the Ceph filesystem mounted under
// `CephFsConfig.GetProjectRoot()`.
func (filer CephFs) CreateProject(projectID string, quotaGiB int) error {

	// create project directory with permission 770
	ppath := filepath.Join(filer.GetProjectRoot(), projectID)

	if _, err := os.Stat(ppath); os.IsNotExist(err) {
		if err := os.Mkdir(ppath, 0770); err != nil {
			return err
		}
	}

	// state the dir again and make sure it is a directory.
	if s, _ := os.Stat(ppath); !s.IsDir() {
		return fmt.Errorf("not a valid directory: %s", ppath)
	}

	// change owner of the project directory
	u, err := user.Lookup(filer.config.ProjectUser)
	if err != nil {
		return err
	}
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)
	if err := os.Chown(ppath, uid, gid); err != nil {
		return fmt.Errorf("cannot set owner of %s: %s", ppath, err)
	}

	// set ACL with `m::rwx,d:m::rwx,g::--,d:g::--` implying
	// - set mask to "rwx" and make it inherited by sub-directories.
	// - remove group access
	if err := setfacl(ppath, []string{
		"-m",
		"m::rwx,d:m::rwx,g::--,d:g::--",
	}); err != nil {
		return err
	}

	// set quota by setting the directory's extended attribute `ceph.quota.max_bytes`.
	if err := setfattr(ppath, []string{
		"-n", "ceph.quota.max_bytes",
		"-v", fmt.Sprintf("%d", quotaGiB<<30),
	}); err != nil {
		return err
	}

	return nil
}

// SetProjectQuota sets or updates the project quota by setting a extended attribute
// `cephfs.quota.max_bytes` of the project directory on the Ceph filesystem.
//
// See [here](https://docs.ceph.com/docs/master/cephfs/quota/) for more detail.
func (filer CephFs) SetProjectQuota(projectID string, quotaGiB int) error {

	ppath := filepath.Join(filer.GetProjectRoot(), projectID)

	// state the dir again and make sure it is a directory.
	if s, _ := os.Stat(ppath); !s.IsDir() {
		return fmt.Errorf("not a valid directory: %s", ppath)
	}

	// set quota by setting the directory's extended attribute `ceph.quota.max_bytes`.
	// NOTE: the value `0` should be replaced with a small value (e.g. 1024 bytes) as
	//       the value `0` literally means "remove the quota limit".
	qb := 1024
	if quotaGiB > 0 {
		qb = quotaGiB << 30
	}

	if err := setfattr(ppath, []string{
		"-n", "ceph.quota.max_bytes",
		"-v", fmt.Sprintf("%d", qb),
	}); err != nil {
		return err
	}

	return nil
}

// GetProjectQuotaInBytes gets the quota of the project `projectID` hosted on the Ceph filesystem.
func (filer CephFs) GetProjectQuotaInBytes(projectID string) (int64, int64, error) {

	ppath := filepath.Join(filer.GetProjectRoot(), projectID)

	quota, err := filer.GetQuotaInBytes(ppath)
	if err != nil {
		return -1, -1, err
	}

	usage, err := filer.GetUsageInBytes(ppath)
	if err != nil {
		return -1, -1, err
	}

	return quota, usage, nil
}

// GetQuotaInBytes gets the quota specified on the given `path` on the Ceph filesystem.
// It is done by retrieving the value of the extended attribute `cephfs.quota.max_bytes`
// from a path.
//
// See [here](https://docs.ceph.com/docs/master/cephfs/quota/) for more detail.
func (filer CephFs) GetQuotaInBytes(path string) (int64, error) {
	// state the dir again and make sure it is a directory.
	if s, _ := os.Stat(path); !s.IsDir() {
		return -1, fmt.Errorf("not a valid directory: %s", path)
	}

	out, err := getfattr(path, []string{
		"--only-values",
		"-n", "ceph.quota.max_bytes",
	})

	if err != nil {
		return -1, nil
	}

	qbytes, err := strconv.Atoi(string(out))
	if err != nil {
		return -1, fmt.Errorf("cannot parse quota value: %s", out)
	}

	return int64(qbytes), nil
}

// GetUsageInBytes gets the storage usage under the given `path` on the Ceph filesystem.
// It is done by retrieving the extended attribute `ceph.dir.rbytes` from the path.
//
// See [here](https://blog.widodh.nl/2015/04/playing-with-cephfs-recursive-statistics/) for more detail.
func (filer CephFs) GetUsageInBytes(path string) (int64, error) {
	if s, _ := os.Stat(path); !s.IsDir() {
		return -1, fmt.Errorf("not a valid directory: %s", path)
	}

	out, err := getfattr(path, []string{
		"--only-values",
		"-n", "ceph.dir.rbytes",
	})

	if err != nil {
		return -1, nil
	}

	ubytes, err := strconv.Atoi(string(out))
	if err != nil {
		return -1, fmt.Errorf("cannot parse usage value: %s", out)
	}

	return int64(ubytes), nil
}

// CreateHome always returns an error with "not supported" message, given that
// Ceph filesystem is not used for personal home directory.
func (filer CephFs) CreateHome(username, groupname string, quotaGiB int) error {
	return fmt.Errorf("not supported")
}

// SetHomeQuota always returns an error with "not supported" message, given that
// Ceph filesystem is not used for personal home directory.
func (filer CephFs) SetHomeQuota(username, groupname string, quotaGiB int) error {
	return fmt.Errorf("not supported")
}

// GetHomeQuotaInBytes always returns an error with "not supported" message, given that
// Ceph filesystem is not used for personal home directory.
func (filer CephFs) GetHomeQuotaInBytes(username, groupname string) (int64, int64, error) {
	return 0, 0, fmt.Errorf("not supported")
}

// setfacl is a command wrapper for executing `setfacl`.
func setfacl(path string, args []string) error {

	cmd := exec.Command("setfacl", append(args, path)...)

	stdout, err := cmd.Output()
	log.Debugf("setfacl stdout: %s", string(stdout))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("setfacl stderr: %s", string(ee.Stderr))
		}
	}

	return nil
}

// setfattr is a command wrapper for executing `setfattr`.
func setfattr(path string, args []string) error {

	cmd := exec.Command("setfattr", append(args, path)...)

	stdout, err := cmd.Output()
	log.Debugf("setfattr stdout: %s", string(stdout))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("setfattr stderr: %s", string(ee.Stderr))
		}
	}

	return nil
}

// getfattr is a command wrapper for executing `getfattr`.
func getfattr(path string, args []string) ([]byte, error) {

	cmd := exec.Command("getfattr", append(args, path)...)

	stdout, err := cmd.Output()
	log.Debugf("setfattr stdout: %s", string(stdout))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return stdout, fmt.Errorf("setfattr stderr: %s", string(ee.Stderr))
		}
	}

	return stdout, nil
}
