package tasks

import (
	"net/http"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/validator"
	"github.com/gorilla/mux"
)

func (taskService TaskService) enableTask(w http.ResponseWriter, r *http.Request, inputStructPtr interface{}) {
	taskID, err := common.ExtractUUID("TaskService.enableTask", w, r, "taskID")
	if err != nil {
		return
	}
	partnerID := mux.Vars(r)["partnerID"]

	if err = validator.ExtractStructFromRequest(r, inputStructPtr); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData,"TaskService.enableTask: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	// update task
	err = taskService.taskPersistence.UpdateTask(r.Context(), inputStructPtr, partnerID, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask,"TaskService.enableTask: can't update an existing task. err=%v", err)
		switch err.(type) {
		case models.TaskIsExpiredError, models.CantUpdateTaskError:
			common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		case models.TaskNotFoundError:
			common.SendNotFound(w, r, errorcode.ErrorCantUpdateTask)
		default:
			common.SendInternalServerError(w, r, errorcode.ErrorCantUpdateTask)
		}
		return
	}

	logger.Log.DebugfCtx(r.Context(), "TaskService.enableTask: updated task with ID = %v", taskID)
	common.SendStatusOkWithMessage(w, r, errorcode.CodeUpdated)
}

// EnableTaskForAllTargets enables/disables future task execution for all targets
func (taskService TaskService) EnableTaskForAllTargets(w http.ResponseWriter, r *http.Request) {
	var allTargetsEnable models.AllTargetsEnable
	taskService.enableTask(w, r, &allTargetsEnable)
}

// EnableTaskForSelectedTargets enables/disables future task execution for specified targets
func (taskService TaskService) EnableTaskForSelectedTargets(w http.ResponseWriter, r *http.Request) {
	var selectedTargetsEnable models.SelectedManagedEndpointEnable
	taskService.enableTask(w, r, &selectedTargetsEnable)
}
