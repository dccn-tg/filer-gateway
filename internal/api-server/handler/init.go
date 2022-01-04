package handler

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Donders-Institute/filer-gateway/internal/api-server/config"
	"github.com/Donders-Institute/filer-gateway/internal/task"
	"github.com/Donders-Institute/filer-gateway/pkg/filer"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/Donders-Institute/tg-toolset-golang/project/pkg/acl"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hurngchunlee/bokchoy"

	fp "github.com/Donders-Institute/tg-toolset-golang/pkg/filepath"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

var (
	// PathProject is the top-leve directory in which directories of active projects are located.
	PathProject string = "/project"

	// QueueSetProject is the queue name for setting project resources.
	QueueSetProject string = "tasks.setProject"

	// QueueSetUser is the queue name for setting user resources.
	QueueSetUser string = "tasks.setUser"

	// DefaultHomeStorageSystem is the name of the default storage system used for user's home directory
	DefaultHomeStorageSystem = "netapp"
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

// GetPing returns dummy string for health check, including the authentication.
func GetPing(cfg config.Configuration) func(params operations.GetPingParams, principle *models.Principle) middleware.Responder {
	return func(params operations.GetPingParams, principle *models.Principle) middleware.Responder {
		return operations.NewGetPingOK().WithPayload("pong")
	}
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
			if err == bokchoy.ErrTaskNotFound { // task not found
				return operations.NewGetTasksTypeIDNotFound()
			}
			return operations.NewGetTasksTypeIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     TaskQueueError,
				},
			)
		}

		taskStatus := task.StatusDisplay()
		taskRslt := ""
		taskErr := ""
		if task.Result != nil {
			taskRslt = fmt.Sprintf("%s", task.Result)
		}
		if task.Error != nil {
			taskErr = fmt.Sprintf("%s", task.Error)
		}

		tid := models.TaskID(task.ID)

		return operations.NewGetTasksTypeIDOK().WithPayload(
			&models.ResponseBodyTaskResource{
				TaskID: &tid,
				TaskStatus: &models.TaskStatus{
					Status: &taskStatus,
					Result: &taskRslt,
					Error:  &taskErr,
				},
			},
		)
	}
}

// CreateProject handles the project creation on filer by formulating it into an.
// asynchronous task.
//
// task configuration:
// - canceled if running more than 12 hours.
// - no retry.
// - result is kept for 7 days.
func CreateProject(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.PostProjectsParams, principle *models.Principle) middleware.Responder {
	return func(params operations.PostProjectsParams, principle *models.Principle) middleware.Responder {

		// construct task data from request data
		t := task.SetProjectResource{
			ProjectID: string(*params.ProjectProvisionData.ProjectID),
			Storage: task.Storage{
				System:  *params.ProjectProvisionData.Storage.System,
				QuotaGb: *params.ProjectProvisionData.Storage.QuotaGb,
			},
			Members: make([]task.Member, 0),
		}

		for _, m := range params.ProjectProvisionData.Members {

			switch *m.Role {
			case acl.Manager.String():
			case acl.Contributor.String():
			case acl.Viewer.String():
			case "none":
			default:
				// only accept setting for manager,contributor and viewer roles
				return operations.NewPostProjectsBadRequest().WithPayload(
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
		task, err := bok.Queue(QueueSetProject).Publish(ctx, &t,
			bokchoy.WithTimeout(12*time.Hour),
			bokchoy.WithMaxRetries(0),
			bokchoy.WithTTL(7*24*time.Hour))

		if err != nil {
			return operations.NewPostProjectsInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     TaskQueueError,
				},
			)
		}

		taskStatus := task.StatusDisplay()

		tid := models.TaskID(task.ID)

		return operations.NewPostProjectsOK().WithPayload(
			&models.ResponseBodyTaskResource{
				TaskID: &tid,
				TaskStatus: &models.TaskStatus{
					Status: &taskStatus,
					Result: nil,
					Error:  nil,
				},
			},
		)
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
// task configuration:
// - canceled if running more than 12 hours.
// - no retry.
// - result is kept for 7 days.
func UpdateProject(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.PatchProjectsIDParams, principle *models.Principle) middleware.Responder {
	// Not implemented
	return func(params operations.PatchProjectsIDParams, principle *models.Principle) middleware.Responder {

		if (params.ProjectUpdateData.Members == nil || len(params.ProjectUpdateData.Members) == 0) && params.ProjectUpdateData.Storage == nil {
			// return 204 No Content if both `members` and `storage` are empty
			return operations.NewPatchProjectsIDNoContent()
		}

		// construct task data from request data
		t := task.SetProjectResource{
			ProjectID: params.ID,
			Storage: task.Storage{
				System:  "none",
				QuotaGb: -1,
			},
			Members: make([]task.Member, 0),
		}

		if params.ProjectUpdateData.Storage != nil {
			t.Storage.QuotaGb = *params.ProjectUpdateData.Storage.QuotaGb
		}

		for _, m := range params.ProjectUpdateData.Members {

			switch *m.Role {
			case acl.Manager.String():
			case acl.Contributor.String():
			case acl.Viewer.String():
			case "none":
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
		task, err := bok.Queue(QueueSetProject).Publish(ctx, &t,
			bokchoy.WithTimeout(12*time.Hour),
			bokchoy.WithMaxRetries(0),
			bokchoy.WithTTL(7*24*time.Hour))

		if err != nil {
			return operations.NewPatchProjectsIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     TaskQueueError,
				},
			)
		}

		taskStatus := task.StatusDisplay()

		tid := models.TaskID(task.ID)

		return operations.NewPatchProjectsIDOK().WithPayload(
			&models.ResponseBodyTaskResource{
				TaskID: &tid,
				TaskStatus: &models.TaskStatus{
					Status: &taskStatus,
					Result: nil,
					Error:  nil,
				},
			},
		)
	}
}

// CreateUserResource handles the request for creating user home space on the filer by
// formulating the request into a asynchronous task.
//
// task configuration:
// - canceled if running more than 1 hour.
// - no retry.
// - result is kept for 7 days.
func CreateUserResource(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.PostUsersParams, principle *models.Principle) middleware.Responder {
	return func(params operations.PostUsersParams, principle *models.Principle) middleware.Responder {
		// construct task data from request data
		t := task.SetUserResource{
			UserID: string(*params.UserProvisionData.UserID),
			Storage: task.Storage{
				System:  *params.UserProvisionData.Storage.System,
				QuotaGb: *params.UserProvisionData.Storage.QuotaGb,
			},
		}

		// publish task to the queue, and set timeout to 12 hours
		// TODO: the timeout should be optimized!!
		task, err := bok.Queue(QueueSetUser).Publish(ctx, &t,
			bokchoy.WithTimeout(1*time.Hour),
			bokchoy.WithMaxRetries(0),
			bokchoy.WithTTL(7*24*time.Hour))

		if err != nil {
			return operations.NewPostUsersInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     TaskQueueError,
				},
			)
		}

		taskStatus := task.StatusDisplay()

		tid := models.TaskID(task.ID)

		return operations.NewPostUsersOK().WithPayload(
			&models.ResponseBodyTaskResource{
				TaskID: &tid,
				TaskStatus: &models.TaskStatus{
					Status: &taskStatus,
					Result: nil,
					Error:  nil,
				},
			},
		)
	}
}

// UpdateUserResource handles the request for updating user home space quota on the filer by
// formulating the request into a asynchronous task.
//
// task configuration:
// - canceled if running more than 1 hour.
// - no retry.
// - result is kept for 7 days.
func UpdateUserResource(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.PatchUsersIDParams, principle *models.Principle) middleware.Responder {
	return func(params operations.PatchUsersIDParams, principle *models.Principle) middleware.Responder {

		if params.UserUpdateData.Storage == nil {
			// return 204 No Content if `storage` is not provided
			return operations.NewPatchUsersIDNoContent()
		}

		// construct task data from request data
		t := task.SetUserResource{
			UserID: params.ID,
			Storage: task.Storage{
				System:  "none",
				QuotaGb: *params.UserUpdateData.Storage.QuotaGb,
			},
		}
		// publish task to the queue, and set timeout to 12 hours
		// TODO: the timeout should be optimized!!
		task, err := bok.Queue(QueueSetUser).Publish(ctx, &t,
			bokchoy.WithTimeout(1*time.Hour),
			bokchoy.WithMaxRetries(0),
			bokchoy.WithTTL(7*24*time.Hour))

		if err != nil {
			return operations.NewPostUsersInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     TaskQueueError,
				},
			)
		}

		taskStatus := task.StatusDisplay()

		tid := models.TaskID(task.ID)

		return operations.NewPostUsersOK().WithPayload(
			&models.ResponseBodyTaskResource{
				TaskID: &tid,
				TaskStatus: &models.TaskStatus{
					Status: &taskStatus,
					Result: nil,
					Error:  nil,
				},
			},
		)
	}
}

// GetProjects implements retrival of resources of all system users with UID >= 1000.
func GetUsers(ucache *UserResourceCache, pcache *ProjectResourceCache) func(params operations.GetUsersParams) middleware.Responder {
	return func(params operations.GetUsersParams) middleware.Responder {

		// max. 4 concurrent workers (because we are already getting data from cache)
		nworkers := 4
		if nworkers > runtime.NumCPU() {
			nworkers = runtime.NumCPU()
		}

		// list all directories in `handler.PathProject`
		usernames := make(chan string, nworkers*2)
		resources := make(chan struct {
			username string
			resource *userResource
		})

		wg := sync.WaitGroup{}
		// start concurrent workers to get project resources from the cache.
		for i := 0; i < nworkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for uname := range usernames {
					if r, err := ucache.getResource(uname, false); err == nil {
						resources <- struct {
							username string
							resource *userResource
						}{
							uname,
							r,
						}
					} else {
						log.Errorf("%s", err)
					}
				}
			}()
		}

		// go routine to get all system users.
		go func() {
			// close the dirs channel on exit
			defer close(usernames)
			for _, uname := range getSystemUsers() {
				usernames <- uname
			}
		}()

		// go routine to wait for all workers to complete and close the resources channel.
		go func() {
			wg.Wait()
			close(resources)
		}()

		users := make([]*models.ResponseBodyUserResource, 0)
		for r := range resources {
			// getting user's membership on all active projects from the cache
			var memberOf = make([]*models.ProjectRole, 0)
			for k, v := range pcache.store {
				pid := k // should reassign the value of `k` for assigning the string pointer of `ProjectID`
				for _, m := range v.members {
					if *m.UserID == r.username {
						memberOf = append(memberOf, &models.ProjectRole{
							ProjectID: &pid,
							Role:      m.Role,
						})
						break
					}
				}
			}

			uid := models.UserID(r.username)
			users = append(users, &models.ResponseBodyUserResource{
				UserID:   &uid,
				Storage:  r.resource.storage,
				MemberOf: memberOf,
			})
		}

		return operations.NewGetUsersOK().WithPayload(
			&models.ResponseBodyUsers{
				Users: users,
			},
		)
	}
}

// GetUserResource implements retrival of file resource for a user (i.e. storage).
func GetUserResource(ucache *UserResourceCache, pcache *ProjectResourceCache) func(params operations.GetUsersIDParams) middleware.Responder {
	return func(params operations.GetUsersIDParams) middleware.Responder {
		uname := params.ID

		ur, err := ucache.getResource(uname, false)

		if err != nil {
			switch err.(type) {
			case user.UnknownUserError:
				return operations.NewGetUsersIDNotFound().WithPayload(err.Error())
			default:
				return operations.NewGetUsersIDInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: err.Error(),
						ExitCode:     UserLookupError,
					},
				)
			}
		}

		// getting user's membership on all active projects from the cache
		var memberOf = make([]*models.ProjectRole, 0)
		for k, v := range pcache.store {
			pid := k // should reassign the value of `k` for assigning the string pointer of `ProjectID`
			for _, m := range v.members {
				if *m.UserID == uname {
					memberOf = append(memberOf, &models.ProjectRole{
						ProjectID: &pid,
						Role:      m.Role,
					})
					break
				}
			}
		}

		// return 200 success with user resource information.
		uid := models.UserID(uname)
		return operations.NewGetUsersIDOK().WithPayload(
			&models.ResponseBodyUserResource{
				UserID:   &uid,
				MemberOf: memberOf,
				Storage:  ur.storage,
			},
		)
	}
}

// GetProjects implements retrival of resources of all projects implemented on the filer, under path of
// `handler.PathProject`.
func GetProjects(cache *ProjectResourceCache) func(params operations.GetProjectsParams) middleware.Responder {
	return func(params operations.GetProjectsParams) middleware.Responder {

		// max. 4 concurrent workers (because we are already getting data from cache)
		nworkers := 4
		if nworkers > runtime.NumCPU() {
			nworkers = runtime.NumCPU()
		}

		// list all directories in `handler.PathProject`
		pnumbers := make(chan string, nworkers*2)
		resources := make(chan struct {
			pnumber  string
			resource *projectResource
		})

		wg := sync.WaitGroup{}
		// start concurrent workers to get project resources from the cache.
		for i := 0; i < nworkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for pnumber := range pnumbers {
					if r, err := cache.getResource(pnumber, false); err == nil {
						resources <- struct {
							pnumber  string
							resource *projectResource
						}{
							pnumber,
							r,
						}
					} else {
						log.Errorf("%s", err)
					}
				}
			}()
		}

		// go routine to get all project numbers in the `PathProject` directory
		go func(path string) {
			// close the dirs channel on exit
			defer close(pnumbers)
			objs, err := fp.ListDir(path)
			if err != nil {
				log.Errorf("cannot get content of path: %s", path)
				return
			}
			for _, obj := range objs {
				pnumbers <- filepath.Base(obj)
			}
		}(PathProject)

		// go routine to wait for all workers to complete and close the resources channel.
		go func() {
			wg.Wait()
			close(resources)
		}()

		projects := make([]*models.ResponseBodyProjectResource, 0)
		for r := range resources {
			pid := models.ProjectID(r.pnumber)
			projects = append(projects, &models.ResponseBodyProjectResource{
				ProjectID: &pid,
				Storage:   r.resource.storage,
				Members:   r.resource.members,
			})
		}

		return operations.NewGetProjectsOK().WithPayload(
			&models.ResponseBodyProjects{
				Projects: projects,
			},
		)
	}
}

// GetProjectResource implements retrival of project resource (i.e. storage and members).
func GetProjectResource(cache *ProjectResourceCache) func(params operations.GetProjectsIDParams) middleware.Responder {
	return func(params operations.GetProjectsIDParams) middleware.Responder {

		r, err := cache.getResource(params.ID, false)

		// Return response error based on error code.
		if err != nil {
			switch err.(*ResponseError).code {
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

		pid := models.ProjectID(params.ID)

		return operations.NewGetProjectsIDOK().WithPayload(
			&models.ResponseBodyProjectResource{
				ProjectID: &pid,
				Storage:   r.storage,
				Members:   r.members,
			},
		)
	}
}

// ResponseError is an internal error type for the API handler function to
// determine which response error should be returned to the API client.
type ResponseError struct {
	code int
	err  string
}

func (e *ResponseError) Error() string {
	return e.err
}

// getStorageSystem retrives storage system based on the suffix of the path.
func getStorageSystem(cfg config.Configuration, path string) string {

	// evaluate symlink to its absolute path.
	path, _ = filepath.EvalSymlinks(path)

	switch true {
	case strings.HasPrefix(path, filer.New("freenas", cfg.FreeNas).GetProjectRoot()):
		return "freenas"
	case strings.HasPrefix(path, filer.New("cephfs", cfg.CephFs).GetProjectRoot()):
		return "cephfs"
	case strings.HasPrefix(path, filer.New("netapp", cfg.NetApp).GetProjectRoot()):
		return "netapp"
	default:
		return "netapp"
	}
}

// getStorageQuota retrives quota limitation and its usage on the path.  The boolean argument `isHomePath` is used
// to indicate whether the path is referring to a home directory (when the value is `true`) or a project directory
// (when the value is `false`).
func getStorageQuota(cfg config.Configuration, path string, isHomePath bool) (system string, quota, usage int64, err error) {

	fi, e := os.Stat(path)

	if e != nil {
		err = &ResponseError{code: 500, err: e.Error()}
		return
	}
	if !fi.Mode().IsDir() {
		err = &ResponseError{code: 500, err: fmt.Sprintf("Not a directory: %s", path)}
		return
	}

	system = getStorageSystem(cfg, path)

	// for cephfs
	var f filer.Filer
	switch system {
	case "cephfs":
		f = filer.New("cephfs", cfg.CephFs)
	case "netapp":
		f = filer.New("netapp", cfg.NetApp)
	case "freenas":
		f = filer.New("freenas", cfg.FreeNas)
	default:
		err = &ResponseError{code: 500, err: fmt.Sprintf("unsupported storage system: %s", system)}
		return
	}

	if isHomePath {
		uname := filepath.Base(path)
		group := filepath.Base(filepath.Dir(path))
		quota, usage, err = f.GetHomeQuotaInBytes(group, uname)
	} else {
		quota, usage, err = f.GetProjectQuotaInBytes(filepath.Base(path))
	}

	if err != nil {
		err = &ResponseError{code: 500, err: err.Error()}
	}

	// if system == "cephfs" {
	// 	cephfs := filer.New("cephfs", cfg.CephFs)
	// 	quota, err = cephfs.(filer.CephFs).GetQuotaInBytes(path)
	// 	usage, err = cephfs.(filer.CephFs).GetUsageInBytes(path)
	// 	return
	// }

	// var filer filer.Filer
	// if system == "netapp"

	// 	// getting storage quota and usage from the filer's API
	// 	var filer filer.Filer
	// 	if system == "netapp"

	// 	// Caution: the code below uses Linux system call to get quota and used space!!
	// 	var stat syscall.Statfs_t
	// 	syscall.Statfs(path, &stat)

	// 	quota = int64(stat.Blocks * uint64(stat.Bsize))
	// 	usage = int64(math.Round(float64((stat.Blocks - stat.Bfree) * uint64(stat.Bsize))))
	// }

	log.Debugf("path: %s, quota: %d bytes, usage: %d bytes", path, quota, usage)
	return
}
