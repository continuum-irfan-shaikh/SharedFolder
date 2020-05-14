package tasks

import (
	"net/http"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// InstancesCount ..
type InstancesCount struct {
	InstancesCount int `json:"instancesCount"`
}

// TaskInstancesCountByTaskID returns aggregated data about Task and its Statuses
func (taskService TaskService) TaskInstancesCountByTaskID(w http.ResponseWriter, r *http.Request) {
	var (
		emptyUUID gocql.UUID
		ctx       = r.Context()
		taskIDStr = mux.Vars(r)["taskID"]
		partnerID = mux.Vars(r)["partnerID"]
	)

	logger.Log.DebugfCtx(r.Context(), "TaskInstancesCountByTaskID: partnerID %v and taskID %v and endpointID %v", partnerID, taskIDStr)

	taskID, err := gocql.ParseUUID(taskIDStr)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, "taskID: %v has bad format. err=%v", taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	taskInstancesCount, err := taskService.taskInstancePersistence.GetInstancesCountByTaskID(ctx, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstanceCountByTaskID, "taskID: %v has bad format. err=%v", taskIDStr, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstanceCountByTaskID)
		return
	}

	common.RenderJSON(w, InstancesCount{
		InstancesCount: taskInstancesCount,
	})
}
