package tasks

import (
	"net/http"
	"sort"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	doneStatus      = "Done"
	failedStatus    = "Failed"
	runningStatus   = "Running"
	partialFailed   = "Partially Failed"
	activeJobsCount = 5
)

// GetByOriginID returns task details by originID
func (t *TaskService) GetByOriginID(w http.ResponseWriter, r *http.Request) {
	originID, err := gocql.ParseUUID(mux.Vars(r)["originID"])
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData,"TaskService.GetByOriginID: can not parse Origin ID %v. err:%v", originID, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	to := time.Now()
	from := to.Add(time.Hour * (-48))

	ar, err := t.getLastTasksResults(r, w, false, from, to)
	if err != nil {
		return
	}

	sort.Slice(ar, func(i, j int) bool { return ar[i].TaskInstance.StartedAt.After(ar[j].TaskInstance.StartedAt) })

	//filter them here
	filteredResults := make([]models.TaskDetailsWithStatuses, 0)
	for _, r := range ar {
		if len(filteredResults) >= activeJobsCount {
			break
		}

		if r.Task.OriginID != originID && r.Task.DefinitionID != originID {
			continue
		}

		r.OverallStatus = calculateJobStatus(r.Statuses)

		//skipping canceled, postponed, disabled
		if r.OverallStatus == "" {
			continue
		}

		filteredResults = append(filteredResults, r)
	}
	common.RenderJSON(w, filteredResults)
}

func calculateJobStatus(allStatuses map[string]int) string {
	if allStatuses[statuses.TaskInstanceRunningText] != 0 {
		return runningStatus
	}

	if allStatuses[statuses.TaskInstanceSuccessText] == len(allStatuses) {
		return doneStatus
	}

	if allStatuses[statuses.TaskInstanceFailedText] == len(allStatuses) {
		return failedStatus
	}

	if (allStatuses[statuses.TaskInstanceFailedText] + allStatuses[statuses.TaskInstanceSuccessText]) == len(allStatuses) {
		return partialFailed
	}
	// if there are success + failed + something else
	if allStatuses[statuses.TaskInstanceFailedText] != 0 && allStatuses[statuses.TaskInstanceSuccessText] != 0 {
		return partialFailed
	}
	// if none of above happened - trying these 2 options
	if allStatuses[statuses.TaskInstanceFailedText] != 0 {
		return failedStatus
	}
	if allStatuses[statuses.TaskInstanceSuccessText] != 0 {
		return doneStatus
	}
	return ""
}
