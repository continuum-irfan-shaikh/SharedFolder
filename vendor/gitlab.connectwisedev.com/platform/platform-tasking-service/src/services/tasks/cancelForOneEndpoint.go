package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	cancelOneMsgTemp = "TaskService.CancelNearestExecutionForEndpoint: "
)

// CancelNearestExecutionForEndpoint is a handler function for canceling nearest execution of task for one particular device
func (t *TaskService) CancelNearestExecutionForEndpoint(w http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		emptyUUID         gocql.UUID
		taskIDStr         = mux.Vars(r)["taskID"]
		partnerID         = mux.Vars(r)["partnerID"]
		managedEndpointID = mux.Vars(r)["managedEndpointID"]
		currentUser       = t.userService.GetUser(r, t.httpClient)
	)
	logger.Log.DebugfCtx(r.Context(), "CancelNearestExecution: partnerID %v and taskID %v and endpointID %v", partnerID, taskIDStr, managedEndpointID)

	taskID, err := gocql.ParseUUID(taskIDStr)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, cancelOneMsgTemp+errTaskIDBadFormat, taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	endpointID, err := gocql.ParseUUID(managedEndpointID)
	if err != nil || endpointID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorEndpointIDHasBadFormat, cancelOneMsgTemp+errMEIDBadFormat, managedEndpointID, err)
		common.SendBadRequest(w, r, errorcode.ErrorEndpointIDHasBadFormat)
		return
	}

	internalTasks, err := t.taskPersistence.GetByIDs(ctx, nil, partnerID, true, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, cancelOneMsgTemp+errCantGetTasks, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if len(internalTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, cancelOneMsgTemp+errZeroTasks, taskID, err)
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if internalTasks[0].IsTrigger() {
		logger.Log.WarnfCtx(r.Context(), cancelOneMsgTemp+errTaskIDBadFormat, taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	if currentUser.HasNOCAccess() != internalTasks[0].IsRequireNOCAccess {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelOneMsgTemp+errWrongNOC, taskIDStr)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	updatedTask, err := t.updateAndValidateTaskFields(currentUser.UID(), internalTasks)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelOneMsgTemp+errCommon, taskID, endpointID)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	taskInstances, err := t.taskInstancePersistence.GetTopInstancesByTaskID(ctx, taskID)
	if err != nil || len(taskInstances) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, cancelOneMsgTemp+errCantGetTopInstances, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	if err := t.processCanceledOne(ctx, taskInstances, r, w, endpointID, updatedTask); err != nil {
		return
	}

	go t.SendTaskUpdateEventToKafka(ctx, taskID, partnerID)
	logger.Log.InfofCtx(r.Context(), cancelOneMsgTemp+successfullyCanceledMsg, taskID, endpointID)
	common.RenderJSON(w, updatedTask)
}

func (t *TaskService) processCanceledOne(ctx context.Context, taskInstances []models.TaskInstance, r *http.Request, w http.ResponseWriter, endpointID gocql.UUID, updatedTask models.Task) error {
	taskID := updatedTask.ID
	if len(taskInstances) > 1 {
		for _, s := range taskInstances[1].Statuses {
			if s == statuses.TaskInstancePending {
				logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelOneMsgTemp+errCantUpdateTask, taskID, statuses.TaskInstancePending)
				common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
				return fmt.Errorf("has pending")
			}
		}
	}

	if taskInstances[0].Statuses[endpointID] != statuses.TaskInstanceScheduled {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, cancelOneMsgTemp+errDeviceStatusIsNotScheduled, taskID, endpointID)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return fmt.Errorf("is not scheduled")
	}

	taskInstances[0].Statuses[endpointID] = statuses.TaskInstanceCanceled
	managedEndpoints := make([]gocql.UUID, 0)

	for deviceID := range taskInstances[0].Statuses {
		managedEndpoints = append(managedEndpoints, deviceID)
	}

	err := t.taskInstancePersistence.Insert(ctx, taskInstances[0])
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTaskInstances, cancelOneMsgTemp+errCantUpdateTaskInstances, taskInstances[0].ID, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantUpdateTaskInstances)
		return err
	}

	err = t.taskPersistence.UpdateModifiedFieldsByMEs(ctx, updatedTask, managedEndpoints...)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskToDB, cancelOneMsgTemp+errCantUpdateTasks, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskToDB)
		return err
	}
	return nil
}
