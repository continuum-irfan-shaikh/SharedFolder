package tasks

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// DeleteExecutions gets particular task's instance and its results by TaskID and TaskInstanceID and deletes it from Cassandra
func (taskService TaskService) DeleteExecutions(w http.ResponseWriter, r *http.Request) {
	var (
		emptyUUID   gocql.UUID
		taskIDStr   = mux.Vars(r)["taskID"]
		partnerID   = mux.Vars(r)["partnerID"]
		ctx         = r.Context()
		currentUser = taskService.userService.GetUser(r, taskService.httpClient)
	)

	taskID, err := gocql.ParseUUID(taskIDStr)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, "taskID: %v has bad format. err=%v", taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	taskInstanceID, err := common.ExtractUUID("DeleteExecutions", w, r, "taskInstanceID")
	if err != nil || taskInstanceID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskInstanceIDHasBadFormat, "taskInstanceID: %v has bad format. err=%v", taskInstanceID, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskInstanceIDHasBadFormat)
		return
	}

	internalTasks, err := taskService.taskPersistence.GetByIDs(ctx, taskService.cache, partnerID, true, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "Can't get internal tasks by task ID: %v. err=%v", taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if len(internalTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID,"Can't get internal tasks by task ID: %v. err=%v", taskID, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if currentUser.HasNOCAccess() != internalTasks[0].IsRequireNOCAccess {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorAccessDenied,"Current user %s is not authorized to delete task's results with task instance ID %v", currentUser.UID(), taskInstanceID)
		common.SendForbidden(w, r, errorcode.ErrorAccessDenied)
		return
	}

	// get taskInstance info for deleting itself and its execution results
	taskInstances, err := taskService.taskInstancePersistence.GetByIDs(ctx, taskInstanceID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, "Can't get task instances by task instance ID %v. err=%v", taskInstanceID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	var taskInstanceIDs []gocql.UUID
	if len(taskInstances) != 0 {
		taskInstanceIDs = []gocql.UUID{taskInstances[0].ID}
	}

	//get executionResults by task instance for their deleting
	executionResults, err := taskService.resultPersistence.GetByTaskInstanceIDs(taskInstanceIDs)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskExecutionResults, "Can't get execution results by task instance IDs %v. err=%v", taskInstanceIDs, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskExecutionResults)
		return
	}

	err = taskService.performDeleteExecutions(ctx, taskInstances, executionResults)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDeleteTaskInstance, "Can't delete task instance and its results. err=%v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantDeleteTaskInstance)
		return
	}

	logger.Log.InfofCtx(r.Context(), "Successfully deleted taskInstance and its results with taskInstanceID = %v", taskInstanceID)
	common.SendNoContent(w)
}

func (taskService TaskService) performDeleteExecutions(
	ctx context.Context,
	taskInstances []models.TaskInstance,
	executionResults []models.ExecutionResult,
) (errTotal error) {

	var (
		wg   = sync.WaitGroup{}
		errs = make(chan error, 2)
		done = make(chan struct{})
	)

	go func() {
		for err := range errs {
			if err != nil {
				if errTotal != nil {
					errTotal = fmt.Errorf("%v;\n%v", errTotal, err)
				} else {
					errTotal = err
				}
			}
		}
		done <- struct{}{}
	}()

	wg.Add(2)
	go func() {
		defer wg.Done()
		taskInstanceErr := taskService.taskInstancePersistence.DeleteBatch(ctx, []models.TaskInstance{taskInstances[0]})
		errs <- taskInstanceErr
	}()

	go func() {
		defer wg.Done()
		resultsErr := taskService.resultPersistence.DeleteBatch(ctx, executionResults)
		errs <- resultsErr
	}()

	wg.Wait()

	close(errs)
	<-done
	close(done)
	return
}
