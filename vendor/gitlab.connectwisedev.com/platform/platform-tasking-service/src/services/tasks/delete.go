package tasks

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

type deleteOneDTO struct {
	taskInstances    []models.TaskInstance
	executionResults []models.ExecutionResult
	triggers         []entities.ActiveTrigger
	execExpirations  []models.ExecutionExpiration
}

// Delete gets all internal tasks, task instances and execution results by TaskID and deletes it from Cassandra
func (t *TaskService) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		emptyUUID   gocql.UUID
		taskIDStr   = mux.Vars(r)["taskID"]
		partnerID   = mux.Vars(r)["partnerID"]
		ctx         = r.Context()
		currentUser = t.userService.GetUser(r, t.httpClient)
	)

	taskID, err := gocql.ParseUUID(taskIDStr)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, "TaskService.Delete: task ID(UUID=%s) has bad format or empty. err=%v", taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	internalTasks, err := t.taskPersistence.GetByIDs(ctx, nil, partnerID, false, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "TaskService.Delete: can't get internal tasks by task ID %v. err=%v", taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if len(internalTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIsNotFoundByTaskID, "TaskService.Delete: task with ID %v not found.", taskID)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIsNotFoundByTaskID)
		return
	}

	commonTaskData := internalTasks[0]
	if currentUser.HasNOCAccess() != commonTaskData.IsRequireNOCAccess {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorAccessDenied, "TaskService.Delete: current user %s is not authorized to delete task with ID %v for partnerID %v", currentUser.UID(), commonTaskData.ID, commonTaskData.PartnerID)
		common.SendForbidden(w, r, errorcode.ErrorAccessDenied)
		return
	}

	dto, err := t.getDataToDelete(ctx, taskID, r, w, partnerID)
	if err != nil {
		return
	}

	dto.tasks = internalTasks
	if err = t.executeDeleting(dto); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDeleteTask, "TaskService.Delete: can't process deleting of the task. err=%v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantDeleteTask)
		return
	}

	if !currentUser.HasNOCAccess() {
		// update counters for tasks in separate goroutine
		go func(ctx context.Context, iTasks []models.Task) {
			counters := getCountersForInternalTasks(iTasks)
			if len(counters) == 0 {
				return
			}

			err := t.taskCounterRepo.DecreaseCounter(commonTaskData.PartnerID, counters, false)
			if err != nil {
				logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "Delete: error while trying to increase counter: ", err)
			}
		}(ctx, internalTasks)
	}

	logger.Log.InfofCtx(r.Context(), "TaskService.Delete: successfully deleted task with ID = %v", taskID)
	common.SendNoContent(w)
}

func (t *TaskService) getDataToDelete(ctx context.Context, taskID gocql.UUID, r *http.Request, w http.ResponseWriter, partnerID string) (deleteDTO, error) {
	taskInstances, err := t.taskInstancePersistence.GetByTaskID(ctx, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, "TaskService.Delete: can't get task instances by task ID %v. err=%v", taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return deleteDTO{}, err
	}

	taskInstanceIDs := make([]gocql.UUID, 0)
	for _, ti := range taskInstances {
		taskInstanceIDs = append(taskInstanceIDs, ti.ID)
	}

	executionResults, err := t.resultPersistence.GetByTaskInstanceIDs(taskInstanceIDs)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskExecutionResults, "TaskService.Delete: can't get execution results by task instance IDs %v. err=%v", taskInstanceIDs, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskExecutionResults)
		return deleteDTO{}, err
	}

	triggers, err := t.trigger.GetActiveTriggersByTaskID(partnerID, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetActiveTriggers, "TaskService.Delete: can't get active triggers by taskID err:", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetActiveTriggers)
		return deleteDTO{}, err
	}

	execExpirations, err := t.executionExpirationPersistence.GetByTaskInstanceIDs(partnerID, taskInstanceIDs)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDeleteExecutionResults, "TaskService.Delete: can't get active triggers by taskID err:", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantDeleteExecutionResults)
		return deleteDTO{}, err
	}
	return deleteDTO{ctx, nil, taskInstances, executionResults, triggers, execExpirations}, nil
}

type deleteDTO struct {
	ctx              context.Context
	tasks            []models.Task
	taskInstances    []models.TaskInstance
	executionResults []models.ExecutionResult
	triggers         []entities.ActiveTrigger
	execExpirations  []models.ExecutionExpiration
}

func (t *TaskService) executeDeleting(
	o deleteDTO,
) (errTotal error) {

	wg := sync.WaitGroup{}
	errs := make(chan error)
	done := make(chan struct{})

	go func() {
		for e := range errs {
			if e != nil {
				if errTotal != nil {
					errTotal = fmt.Errorf("%v;\n%v", errTotal, e)
				} else {
					errTotal = e
				}
			}
		}
		done <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tErr := t.taskPersistence.Delete(o.ctx, o.tasks)
		errs <- tErr
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tiErr := t.taskInstancePersistence.DeleteBatch(o.ctx, o.taskInstances)
		errs <- tiErr
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		erErr := t.resultPersistence.DeleteBatch(o.ctx, o.executionResults)
		errs <- erErr
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		erErr := t.trigger.DeleteActiveTriggers(o.triggers)
		errs <- erErr
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, ex := range o.execExpirations {
			erErr := t.executionExpirationPersistence.Delete(ex)
			errs <- erErr
		}
	}()

	wg.Wait()

	close(errs)
	<-done
	close(done)
	return
}
