package handlers

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/acl"
	"github.com/go-openapi/runtime/middleware"

	log "github.com/sirupsen/logrus"
)

// Error code definitions.
var (
	// NotImplementedError indicates the implementation of the handler is not implemented yet.
	NotImplementedError int64 = 100

	// RoleSettingError indicates error when setting the Role/ACL on the underlying filer.
	RoleSettingError int64 = 101

	// RoleGettingError indicates error when getting the Role/ACL from the underlying filer.
	RoleGettingError int64 = 102

	// QuotaGettingError indicates error when getting the quota limitation and usage from the underlying filer.
	QuotaGettingError int64 = 103
)

// Common payload for the ResponseBody500.
var responseNotImplemented = models.ResponseBody500{
	ErrorMessage: "not implemented",
	ExitCode:     NotImplementedError,
}

// CreateProject implements the project creation on filer systems.
//
func CreateProject() func(params operations.PostProjectsParams) middleware.Responder {
	// Not implemented
	return func(params operations.PostProjectsParams) middleware.Responder {
		return operations.NewPostProjectsInternalServerError().WithPayload(&responseNotImplemented)
	}
}

// UpdateProject implements the project update on filer systems.  Those updates can be one of
// the following:
//
// - update project quota.
// - set project members.
//
// The corresponding project directory on the filer should exist in advance.
//
func UpdateProject() func(params operations.PatchProjectsIDParams) middleware.Responder {
	// Not implemented
	return func(params operations.PatchProjectsIDParams) middleware.Responder {
		return operations.NewPatchProjectsIDInternalServerError().WithPayload(&responseNotImplemented)
	}
}

// GetProjectResource implements retrival of project resource (i.e. storage and members).
func GetProjectResource() func(params operations.GetProjectsIDParams) middleware.Responder {
	return func(params operations.GetProjectsIDParams) middleware.Responder {
		pid := params.ID
		path, e := pid2path(pid)
		if e != nil {
			return operations.NewGetProjectsIDNotFound().WithPayload(e.Error())
		}

		// Get Storage Resource
		system, quota, usage, err := getStorageQuota(path)
		// Return response error based on error code.
		if err != nil {
			switch err.code {
			case 404:
				return operations.NewGetProjectsIDNotFound().WithPayload(err.Error())
			default:
				return operations.NewGetProjectsIDInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: err.Error(),
						ExitCode:     QuotaGettingError,
					},
				)
			}
		}

		// Get project memeber and roles.
		members, err := getMemberRoles(path)
		// Return response error based on error code.
		if err != nil {
			switch err.code {
			case 404:
				return operations.NewGetProjectsIDNotFound().WithPayload(err.Error())
			default:
				return operations.NewGetProjectsIDInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: err.Error(),
						ExitCode:     QuotaGettingError,
					},
				)
			}
		}

		// return 200 success with storage quota information.
		return operations.NewGetProjectsIDOK().WithPayload(
			&models.ResponseBodyProjectResource{
				ProjectID: models.ProjectID(pid),
				Storage: &models.Storage{
					QuotaGb: &quota,
					System:  &system,
					UsageGb: &usage,
				},
				Members: models.Members(members),
			},
		)
	}
}

// GetProjectStorage implements retrival of project storage information.
func GetProjectStorage() func(params operations.GetProjectsIDStorageParams) middleware.Responder {
	// implementation
	return func(params operations.GetProjectsIDStorageParams) middleware.Responder {
		pid := params.ID
		path, e := pid2path(pid)
		if e != nil {
			return operations.NewGetProjectsIDStorageNotFound().WithPayload(e.Error())
		}

		log.Debugf("get storage quota on %s\n", path)

		system, quota, usage, err := getStorageQuota(path)

		// Return response error based on error code.
		if err != nil {
			switch err.code {
			case 404:
				return operations.NewGetProjectsIDStorageNotFound().WithPayload(err.Error())
			default:
				return operations.NewGetProjectsIDStorageInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: err.Error(),
						ExitCode:     QuotaGettingError,
					},
				)
			}
		}

		// return 200 success with storage quota information.
		return operations.NewGetProjectsIDStorageOK().WithPayload(
			&models.ResponseBodyProjectStorage{
				ProjectID: models.ProjectID(pid),
				Storage: &models.Storage{
					QuotaGb: &quota,
					System:  &system,
					UsageGb: &usage,
				},
			},
		)
	}
}

// GetProjectMembers implements retrival of project members and their roles from the project directory
// on the filer.
//
// The corresponding project directory on the filer should exist in advance.
//
func GetProjectMembers() func(params operations.GetProjectsIDMembersParams) middleware.Responder {

	// implementation
	return func(params operations.GetProjectsIDMembersParams) middleware.Responder {

		pid := params.ID

		path, e := pid2path(pid)
		if e != nil {
			return operations.NewGetProjectsIDMembersNotFound().WithPayload(e.Error())
		}

		log.Debugf("get project memebers on %s\n", path)

		members, err := getMemberRoles(path)

		// Return response error based on error code.
		if err != nil {
			switch err.code {
			case 404:
				return operations.NewGetProjectsIDMembersNotFound().WithPayload(err.Error())
			default:
				return operations.NewGetProjectsIDMembersInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: err.Error(),
						ExitCode:     RoleGettingError,
					},
				)
			}
		}

		// Return 200 and list of members as response body.
		return operations.NewGetProjectsIDMembersOK().WithPayload(
			&models.ResponseBodyProjectMembers{
				ProjectID: models.ProjectID(pid),
				Members:   models.Members(members),
			},
		)
	}
}

// responseError is an internal error type for the API handler function to
// determine which response error should be returned to the API client.
type responseError struct {
	code int
	err  string
}

func (e *responseError) Error() string {
	return e.err
}

// pid2path converts project id to file system path.
func pid2path(pid string) (string, error) {
	var path string
	if matched, _ := regexp.MatchString("^[0-9]{7,}", pid); matched {
		// input pid is a project number
		path = filepath.Join("/project", pid)
	}

	// evaluate symlink to its absolute path.
	return filepath.EvalSymlinks(path)
}

// getStorageSystem retrives storage system based on the suffix of the path.
func getStorageSystem(path string) string {

	// evaluate symlink to its absolute path.
	path, _ = filepath.EvalSymlinks(path)

	system := "netapp"

	if strings.HasPrefix(path, "/project_freenas/") {
		system = "freenas"
	}

	return system
}

// getStorageQuota retrives quota limitation and its usage on the path.
func getStorageQuota(path string) (system string, quota, usage int64, err *responseError) {

	fi, e := os.Stat(path)

	if e != nil {
		err = &responseError{code: 500, err: err.Error()}
		return
	}
	if !fi.Mode().IsDir() {
		err = &responseError{code: 500, err: fmt.Sprintf("Not a directory: %s", path)}
		return
	}

	// Caution: the code below uses Linux system call to get quota and used space!!
	var stat syscall.Statfs_t
	syscall.Statfs(path, &stat)

	gib := 1024. * 1024 * 1024

	quota = int64((stat.Blocks * uint64(stat.Bsize)) >> 30)
	usage = int64(math.Round(float64(((stat.Blocks - stat.Bfree) * uint64(stat.Bsize))) / gib))
	system = getStorageSystem(path)

	log.Debugf("path: %s, quota: %d GiB, usage: %d GiB", path, quota, usage)

	return
}

// getMemberRoles retrives member roles applied on the path.
func getMemberRoles(path string) ([]*models.Member, *responseError) {

	members := make([]*models.Member, 0)

	runner := acl.Runner{
		RootPath:   path,
		FollowLink: true,
		SkipFiles:  true,
		Nthreads:   1,
	}

	chanOut, err := runner.GetRoles(false)

	// we know it's path not found error because this is the only case the runner.GetRoles returns an error.
	// TODO: maybe the runner should return an explicit error type.
	if err != nil {
		return members, &responseError{code: 500, err: fmt.Sprintf("cannot get role: %s", path)}
	}

	// only one object is expected from the channel as the recursion is disabled on the runner function.
	for o := range chanOut {
		log.Debugf("found project memebers on %s, %+v\n", o.Path, o.RoleMap)
		for r, users := range o.RoleMap {
			// exclude the system role.
			if r == acl.System {
				continue
			}
			rname := r.String()
			for i := range users {
				m := models.Member{
					UserID: &users[i],
					Role:   &rname,
				}
				members = append(members, &m)
			}
		}
	}

	return members, nil
}
