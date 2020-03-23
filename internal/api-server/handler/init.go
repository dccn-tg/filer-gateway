package handler

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/task"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/acl"
	"github.com/go-openapi/runtime/middleware"
	"github.com/thoas/bokchoy"

	fp "github.com/Donders-Institute/tg-toolset-golang/pkg/filepath"
	log "github.com/sirupsen/logrus"
)

var (
	// PathProject is the top-leve directory in which directories of active projects are located.
	PathProject string = "/project"

	// PathProjectFreenas is the top-level mount point of project hosted on FreeNAS box.
	PathProjectFreenas string = "/project_freenas"

	// QueueSetProject is the queue name for setting project resources.
	QueueSetProject string = "tasks.setProject"

	// QueueSetUser is the queue name for setting user resources.
	QueueSetUser string = "tasks.setUser"
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

	// UserLookupError indicates error when looking up the user in the system.
	UserLookupError int64 = 104

	// MemberOfGettingError indicates error when looking up user's membership on all active projects.
	MemberOfGettingError int64 = 105

	// TaskQueueError indicates the request is failed to be added to the task queue for resource setting.
	TaskQueueError int64 = 106
)

// Common payload for the ResponseBody500.
var responseNotImplemented = models.ResponseBody500{
	ErrorMessage: "not implemented",
	ExitCode:     NotImplementedError,
}

// GetTask retrieves task status.
func GetTask(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.GetTasksTypeIDParams) middleware.Responder {
	return func(params operations.GetTasksTypeIDParams) middleware.Responder {
		id := params.ID

		var qn string

		switch t := params.Type; t {
		case "project":
			qn = QueueSetProject
		case "user":
			qn = QueueSetUser
		default:
			qn = ""
		}

		if qn == "" {
			return operations.NewGetTasksTypeIDBadRequest().WithPayload(
				&models.ResponseBody400{
					ErrorMessage: fmt.Sprintf("invalid task type: %s", params.Type),
				},
			)
		}

		// retrieve task from the queue
		task, err := bok.Queue(qn).Get(ctx, id)

		if err != nil {
			return operations.NewGetTasksTypeIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     TaskQueueError,
				},
			)
		}

		return operations.NewGetTasksTypeIDOK().WithPayload(
			&models.ResponseBodyTaskResource{
				TaskID:     models.TaskID(task.ID),
				TaskStatus: models.TaskStatus(task.StatusDisplay()),
			},
		)
	}
}

// CreateProject implements the project creation on filer systems.
//
func CreateProject(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.PostProjectsParams) middleware.Responder {
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
func UpdateProject(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.PatchProjectsIDParams) middleware.Responder {
	// Not implemented
	return func(params operations.PatchProjectsIDParams) middleware.Responder {

		// construct task data from request data
		t := task.SetProjectResource{
			ProjectID: params.ID,
			Storage: task.Storage{
				System:  *params.ProjectUpdateData.Storage.System,
				QuotaGb: *params.ProjectUpdateData.Storage.QuotaGb,
			},
			Members: make([]task.Member, 0),
		}

		for _, m := range params.ProjectUpdateData.Members {

			switch *m.Role {
			case acl.Manager.String():
			case acl.Contributor.String():
			case acl.Viewer.String():
			default:
				// only accept setting for manager,contributor and viewer roles
				return operations.NewPatchProjectsIDBadRequest().WithPayload(
					&models.ResponseBody400{
						ErrorMessage: fmt.Sprintf("invalid member role for set: %s", *m.Role),
					},
				)
			}

			t.Members = append(t.Members, task.Member{
				UserID: *m.UserID,
				Role:   *m.Role,
			})
		}

		// publish task to the queue, and set timeout to 12 hours
		// TODO: the timeout should be optimized!!
		task, err := bok.Queue(QueueSetProject).Publish(ctx, &t, bokchoy.WithTimeout(12*time.Hour))

		if err != nil {
			return operations.NewPatchProjectsIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     TaskQueueError,
				},
			)
		}

		return operations.NewPatchProjectsIDOK().WithPayload(
			&models.ResponseBodyTaskResource{
				TaskID:     models.TaskID(task.ID),
				TaskStatus: models.TaskStatus(task.StatusDisplay()),
			},
		)
	}
}

// GetUserResource implements retrival of file resource for a user (i.e. storage).
func GetUserResource() func(params operations.GetUsersIDParams) middleware.Responder {
	return func(params operations.GetUsersIDParams) middleware.Responder {
		uname := params.ID

		u, e := user.Lookup(uname)

		if e != nil {
			switch e.(type) {
			case user.UnknownUserError:
				return operations.NewGetUsersIDNotFound().WithPayload(e.Error())
			default:
				return operations.NewGetUsersIDInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: e.Error(),
						ExitCode:     UserLookupError,
					},
				)
			}
		}

		// getting storage quota on the user's home directory
		system, quota, usage, err := GetStorageQuota(u.HomeDir)

		// Return response error based on error code.
		if err != nil {
			switch err.code {
			case 404:
				return operations.NewGetUsersIDNotFound().WithPayload(err.Error())
			default:
				return operations.NewGetUsersIDInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: err.Error(),
						ExitCode:     QuotaGettingError,
					},
				)
			}
		}

		// getting user's membership on all active projects
		memberOf, err := getMemberOf(uname)
		if err != nil {
			switch err.code {
			case 404:
				return operations.NewGetUsersIDNotFound().WithPayload(err.Error())
			default:
				return operations.NewGetUsersIDInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: err.Error(),
						ExitCode:     MemberOfGettingError,
					},
				)
			}
		}

		// return 200 success with storage quota information.
		return operations.NewGetUsersIDOK().WithPayload(
			&models.ResponseBodyUserResource{
				UserID:   models.UserID(uname),
				MemberOf: memberOf,
				Storage: &models.StorageResponse{
					QuotaGb: &quota,
					System:  &system,
					UsageGb: &usage,
				},
			},
		)
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
		system, quota, usage, err := GetStorageQuota(path)
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
				Storage: &models.StorageResponse{
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

		system, quota, usage, err := GetStorageQuota(path)

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
				Storage: &models.StorageResponse{
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
		path = filepath.Join(PathProject, pid)
	}

	// evaluate symlink to its absolute path.
	return filepath.EvalSymlinks(path)
}

// getStorageSystem retrives storage system based on the suffix of the path.
func getStorageSystem(path string) string {

	// evaluate symlink to its absolute path.
	path, _ = filepath.EvalSymlinks(path)

	system := "netapp"

	if strings.HasPrefix(path, PathProjectFreenas) {
		system = "freenas"
	}

	return system
}

// GetStorageQuota retrives quota limitation and its usage on the path.
func GetStorageQuota(path string) (system string, quota, usage int64, err *responseError) {

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

// getMemberOf scans through all projects' top-level directories to find out
// the membership of the user `uid`.
func getMemberOf(uid string) ([]*models.ProjectRole, *responseError) {

	nworkers := runtime.NumCPU()

	dirs := make(chan string, nworkers*2)
	members := make(chan *models.ProjectRole)

	wg := sync.WaitGroup{}
	for i := 0; i < nworkers; i++ {
		wg.Add(1)
		go findUserMember(uid, dirs, members, &wg)
	}

	// go routine to list all directories in the /project folder
	go func(path string) {
		// close the dirs channel on exit
		defer close(dirs)

		objs, err := fp.ListDir(path)
		if err != nil {
			log.Errorf("cannot get content of path: %s", path)
			return
		}

		for _, obj := range objs {
			dirs <- obj
		}

	}(PathProject)

	// go routine to wait for all workers to complete and close the members channel.
	go func() {
		wg.Wait()
		close(members)
	}()

	// making up the output from the data in the members channel.
	memberOf := make([]*models.ProjectRole, 0)
	for member := range members {
		memberOf = append(memberOf, member)
	}

	return memberOf, nil
}

// findUserMember is a worker function that scan through each entry in `dirs`, and search for a member role
// of a given user `uid`.
func findUserMember(uid string, dirs chan string, members chan *models.ProjectRole, wg *sync.WaitGroup) {

	defer wg.Done()

	for dir := range dirs {

		log.Debugf("finding user member for %s in %s", uid, dir)

		// get all members of the dir
		runner := acl.Runner{
			RootPath:   dir,
			FollowLink: true,
			SkipFiles:  true,
			Nthreads:   1,
		}

		chanOut, err := runner.GetRoles(false)
		if err != nil {
			log.Errorf("cannot get role for path %s: %s", dir, err)
			continue
		}

		// feed members channel if the user in question is in the list.
		for o := range chanOut {
			for r, users := range o.RoleMap {
				if r == acl.System {
					continue
				}
				rstr := r.String()
				pid := filepath.Base(dir)
				for _, u := range users {
					if u == uid {
						members <- &models.ProjectRole{
							ProjectID: &pid,
							Role:      &rstr,
						}
						break
					}
					continue
				}
			}
		}
	}
}
