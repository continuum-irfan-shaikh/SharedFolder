package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	cancelMsgTemplate             = "TaskService.CancelNearestExecution: "
	errMEIDBadFormat              = "cannot ParseUUID of EndpointID [%s]: %v"
	successfullyCanceledMsg       = "successfully canceled task with ID = %v for endpointID = %v"
	errTaskIDBadFormat            = "task ID %v has bad format. err: %v"
	errTaskIsTrigger              = "task ID %v is trigger and can't be canceled"
	errCantGetTasks               = "can not get internal tasks by Task ID %v. err: %v"
	errZeroTasks                  = "got zero internal tasks for Task ID %v. err: %v"
	errWrongNOC                   = "wrong NOC access task by task ID %v"
	errCommon                     = "task ID %v err: %v"
	errCantGetTopInstances        = "can't get top task instances by taskID %v. err: %v"
	errZeroInstances              = "got zero instance by ID %v"
	errTaskIsAlreadyExecuted      = "can't update already executed task with ID  %v"
	errTaskIsAlreadyCanceled      = "can't update already canceled task with ID  %v"
	errCantUpdateTaskInstances    = "error while updating Instances. TI ID: %v. err: %v"
	errCantUpdateTasks            = "error while updating internal tasks. Task ID: %v. err: %v"
	errDeviceStatusIsNotScheduled = "error while updating device status. Task ID: %v, deviceID: %v"
)

// CancelNearestExecution is a handler function to cancel nearest execution of a task for all devices
func (t *TaskService) CancelNearestExecution(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		emptyUUID   gocql.UUID
		taskIDStr   = mux.Vars(r)["taskID"]
		partnerID   = mux.Vars(r)["partnerID"]
		currentUser = t.userService.GetUser(r, t.httpClient)
	)

	logger.Log.DebugfCtx(r.Context(), "CancelNearestExecution: received cancel req for partnerID %v and taskID %v", partnerID, taskIDStr)

	// ---- start of getting and validating all needed data
	taskID, err := gocql.ParseUUID(taskIDStr)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, cancelMsgTemplate+errTaskIDBadFormat, taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	internalTasks, err := t.taskPersistence.GetByIDs(ctx, nil, partnerID, true, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, cancelMsgTemplate+errCantGetTasks, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if len(internalTasks) == 0 {
		logger.Log.WarnfCtx(r.Context(), cancelMsgTemplate+errZeroTasks, taskID, err)
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if internalTasks[0].IsTrigger() {
		logger.Log.WarnfCtx(r.Context(), cancelMsgTemplate+errTaskIsTrigger, taskID, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	if currentUser.HasNOCAccess() != internalTasks[0].IsRequireNOCAccess {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelMsgTemplate+errWrongNOC, taskIDStr)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}
	// ---- end of getting and validating all needed data

	updatedTask, err := t.updateAndValidateTaskFields(currentUser.UID(), internalTasks)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(),  errorcode.ErrorCantUpdateTask, cancelMsgTemplate+errCommon, taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	taskInstances, err := t.taskInstancePersistence.GetTopInstancesByTaskID(ctx, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, cancelMsgTemplate+errCantGetTopInstances, taskID.String(), err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	if len(taskInstances) == 0 {
		logger.Log.ErrfCtx(r.Context(),  errorcode.ErrorCantGetTaskInstances, cancelMsgTemplate+errZeroInstances, taskID.String())
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	if err := t.processCanceledData(ctx, taskInstances, r, taskID, w, updatedTask); err != nil {
		return
	}

	go t.SendTaskUpdateEventToKafka(ctx, taskID, partnerID)
	common.RenderJSON(w, updatedTask)
}

func (t *TaskService) processCanceledData(ctx context.Context, taskInstances []models.TaskInstance, r *http.Request, taskID gocql.UUID, w http.ResponseWriter, updatedTask models.Task) error {
	if len(taskInstances) > 1 {
		for _, s := range taskInstances[1].Statuses {
			if s == statuses.TaskInstancePending {
				logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelMsgTemplate+errCantUpdateTask, taskID.String(), statuses.TaskInstancePending)
				common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
				return fmt.Errorf("has penging")
			}
		}
	}

	var managedEndpoints []gocql.UUID
	canceledCount := 0
	for deviceID := range taskInstances[0].Statuses {
		// additional check for OneTime Tasks
		status := taskInstances[0].Statuses[deviceID]
		switch status {
		case statuses.TaskInstanceScheduled:
			taskInstances[0].Statuses[deviceID] = statuses.TaskInstanceCanceled
			managedEndpoints = append(managedEndpoints, deviceID)
			continue
		case statuses.TaskInstanceCanceled:
			managedEndpoints = append(managedEndpoints, deviceID)
			canceledCount++
			continue
		case statuses.TaskInstancePending, statuses.TaskInstanceStopped, statuses.TaskInstanceFailed,
			statuses.TaskInstanceSuccess, statuses.TaskInstanceRunning:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelMsgTemplate+errTaskIsAlreadyExecuted, taskID.String())
			common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
			return fmt.Errorf("has %v", status)
		}
	}

	if canceledCount == len(taskInstances[0].Statuses) {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelMsgTemplate+errTaskIsAlreadyCanceled, taskID.String())
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return fmt.Errorf("already canceled")
	}

	if err := t.taskInstancePersistence.Insert(ctx, taskInstances[0]); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTaskInstances, cancelMsgTemplate+errCantUpdateTaskInstances, taskInstances[0].ID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantUpdateTaskInstances)
		return err
	}

	if err := t.taskPersistence.UpdateModifiedFieldsByMEs(ctx, updatedTask, managedEndpoints...); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskToDB, cancelMsgTemplate+errCantUpdateTasks, taskID.String(), err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskToDB)
		return err
	}
	return nil
}

func (*TaskService) updateAndValidateTaskFields(modifiedBy string, tasks []models.Task) (updatedTask models.Task, err error) {
	if len(tasks) == 0 {
		return models.Task{}, errors.New("got empty slice of tasks")
	}
	modifiedAt := time.Now().Truncate(time.Minute).UTC()
	tasks[0].ModifiedBy = modifiedBy
	tasks[0].ModifiedAt = modifiedAt
	return tasks[0], nil
}
