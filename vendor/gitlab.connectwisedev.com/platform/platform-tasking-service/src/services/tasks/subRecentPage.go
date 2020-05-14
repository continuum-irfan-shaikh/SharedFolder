package tasks

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// SubRecent returns data by taskID
func (taskService TaskService) SubRecent(w http.ResponseWriter, r *http.Request) {
	var (
		taskIDStr   = mux.Vars(r)["taskID"]
		currentUser = taskService.userService.GetUser(r, taskService.httpClient)
	)

	tasks, mapTasksByEndpointID, canBeCanceledMap, canBePostponedMap, err := taskService.getTasks(taskIDStr, r, w, currentUser)
	if err != nil {
		return
	}

	// panic is checked above
	taskID := tasks[0].ID

	var (
		taskInstanceErr, taskTemplateErr error
		taskTemplate                     models.TemplateDetails
		uniqueTaskInstanceIDs            []gocql.UUID
		taskCommonData                   = tasks[0]
		wg                               = &sync.WaitGroup{}
		taskInstancesByIDMap             = make(map[gocql.UUID]models.TaskInstance)                   // key: TaskInstanceID
		execResByTaskInstIDAndMEIDMap    = make(map[gocql.UUID]map[gocql.UUID]models.ExecutionResult) // keys: TaskInstanceID and ManagedEndpointID
	)

	// concurrency START
	wg.Add(2)

	// TaskInstances BEGIN
	go func() {
		defer wg.Done()
		taskInstancesByIDMap, uniqueTaskInstanceIDs, taskInstanceErr = taskService.getInstances(r.Context(), taskID)
	}()
	//  TaskInstances END

	// TaskTemplate START
	go func() {
		defer wg.Done()

		// case for external services like Sequence of Patching which doesn't have Templates
		if taskCommonData.Type != models.TaskTypeScript {
			taskTemplate = models.TemplateDetails{
				OriginID:       taskCommonData.OriginID,
				SuccessMessage: "Executed successfully",
				FailureMessage: "Failed",
			}
			return
		}

		taskTemplate, taskTemplateErr = taskService.templateCache.GetByOriginID(
			r.Context(),
			currentUser.PartnerID(),
			taskCommonData.OriginID,
			currentUser.HasNOCAccess(),
		)
	}()
	// TaskTemplate END

	wg.Wait()
	// concurrency END

	if taskInstanceErr != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, "TaskService.SubRecent: cannot get TaskInstances from DB: %s", taskInstanceErr)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}
	if taskTemplateErr != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskDefinitionTemplate, "TaskService.SubRecent: cannot get TaskTemplate from Cache: %s", taskTemplateErr)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskDefinitionTemplate)
		return
	}

	executionResults, err := taskService.resultPersistence.GetByTaskInstanceIDs(uniqueTaskInstanceIDs)
	if err != nil {
		logger.Log.WarnfCtx(r.Context(), "TaskService.SubRecent: cannot get ExecutionResults from DB: ", err)
	}

	if len(executionResults) > 0 {
		for _, eRes := range executionResults {
			if _, ok := execResByTaskInstIDAndMEIDMap[eRes.TaskInstanceID]; !ok {
				execResByTaskInstIDAndMEIDMap[eRes.TaskInstanceID] = make(map[gocql.UUID]models.ExecutionResult)
			}
			execResByTaskInstIDAndMEIDMap[eRes.TaskInstanceID][eRes.ManagedEndpointID] = eRes
		}
	}

	dto := resultsDTO{
		taskInstancesByIDMap:          taskInstancesByIDMap,
		execResByTaskInstIDAndMEIDMap: execResByTaskInstIDAndMEIDMap,
		mapTasksByEndpointID:          mapTasksByEndpointID,
		taskTemplate:                  taskTemplate,
		taskCommonData:                taskCommonData,
		canBePostponedMap:             canBePostponedMap,
		canBeCanceledMap:              canBeCanceledMap,
		r:                             r,
	}

	taskSummaryDetails := taskService.agregateSubrecentResults(dto)
	logger.Log.InfofCtx(r.Context(), "TaskService.SubRecent: Tasks' summary aggregated successful for Partner [%s] and TaskID [%s]",
		currentUser.PartnerID(), taskIDStr)
	common.RenderJSON(w, taskSummaryDetails)
}

func (taskService TaskService) getInstances(ctx context.Context, taskID gocql.UUID) (map[gocql.UUID]models.TaskInstance, []gocql.UUID, error) {
	taskInstancesByIDMap := make(map[gocql.UUID]models.TaskInstance)
	taskInstances, taskInstanceErr := taskService.taskInstancePersistence.GetByTaskID(ctx, taskID)
	if taskInstanceErr != nil {
		return nil, nil, taskInstanceErr
	}

	uniqueTaskInstanceIDs := getAllUniqueTaskInstanceIDs(taskInstances)

	for _, ti := range taskInstances {
		taskInstancesByIDMap[ti.ID] = ti
	}
	return taskInstancesByIDMap, uniqueTaskInstanceIDs, nil
}

func (taskService TaskService) getTasks(taskIDStr string, r *http.Request, w http.ResponseWriter, currentUser user.User) ([]models.Task, map[gocql.UUID]models.Task, map[gocql.UUID]bool, map[gocql.UUID]bool, error) {
	taskID, err := gocql.ParseUUID(taskIDStr)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "TaskService.SubRecent: cannot parse UUID: %s", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return nil, nil, nil, nil, err
	}

	tasks, err := taskService.taskPersistence.GetByIDs(
		r.Context(),
		nil,
		currentUser.PartnerID(),
		false,
		taskID,
	)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "TaskService.SubRecent: cannot get Tasks from DB: ", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return nil, nil, nil, nil, err
	}

	if len(tasks) < 1 {
		logger.Log.WarnfCtx(r.Context(), "TaskService.SubRecent: no Tasks with TaskID [%s] for Partner [%s]",
			taskIDStr, currentUser.PartnerID())
		common.SendNotFound(w, r, errorcode.ErrorTaskIsNotFoundByTaskID)
		return nil, nil, nil, nil, fmt.Errorf("not found")
	}

	mapTasksByEndpointID := make(map[gocql.UUID]models.Task, len(tasks))
	//key - taskID
	canBeCanceledMap := make(map[gocql.UUID]bool, len(tasks))
	canBePostponedMap := make(map[gocql.UUID]bool, len(tasks))
	for _, task := range tasks {
		mapTasksByEndpointID[task.ManagedEndpointID] = task
		canBeCanceledMap[task.ID] = !task.IsTrigger() && !task.IsTaskAndTriggerNotActivated()
		canBePostponedMap[task.ID] = !task.IsTrigger() && task.IsTaskAndTriggerNotActivated()
	}
	return tasks, mapTasksByEndpointID, canBeCanceledMap, canBePostponedMap, nil
}

type resultsDTO struct {
	taskInstancesByIDMap          map[gocql.UUID]models.TaskInstance
	execResByTaskInstIDAndMEIDMap map[gocql.UUID]map[gocql.UUID]models.ExecutionResult
	mapTasksByEndpointID          map[gocql.UUID]models.Task
	taskTemplate                  models.TemplateDetails
	taskCommonData                models.Task
	canBePostponedMap             map[gocql.UUID]bool
	canBeCanceledMap              map[gocql.UUID]bool
	r                             *http.Request
}

func (taskService TaskService) agregateSubrecentResults(o resultsDTO) models.TaskSummaryDetails {
	var (
		taskInstanceSummaries = make([]models.TaskInstanceSummary, 0)
		nearestNextRunTime    time.Time
		ts                    models.TargetSummary
	)

	for taskInstID, taskInstance := range o.taskInstancesByIDMap {
		targetSummary := make([]models.TargetSummary, 0, len(taskInstance.Statuses))

		for endpointID, status := range taskInstance.Statuses {
			nearestNextRunTime, ts = taskService.buildTargetSummary(o, taskInstID, endpointID, status, nearestNextRunTime, taskInstance)
			targetSummary = append(targetSummary, ts)
		}

		allStatuses, err := taskInstance.CalculateStatuses()
		if err != nil {
			logger.Log.WarnfCtx(o.r.Context(), "TaskService.SubRecent: cannot CalculateStatuses: %s", err)
			continue
		}

		taskInstanceSummaries = append(taskInstanceSummaries, models.TaskInstanceSummary{
			ID:              taskInstID,
			RunTime:         taskInstance.LastRunTime,
			TargetSummaries: targetSummary,
			RunStatuses: models.DeviceStatuses{
				DeviceCount: len(taskInstance.Statuses),
				Statuses:    allStatuses,
			},
		})
	}

	taskSummaryDetails := models.TaskSummaryDetails{
		TaskSummary: models.TaskSummaryData{
			Name:   o.taskCommonData.Name,
			TaskID: o.taskCommonData.ID,
			Type:   o.taskCommonData.Type,
			RunOn: models.TargetData{
				Count: len(o.mapTasksByEndpointID), //target's count
			},
			Regularity:         o.taskCommonData.Schedule.Regularity,
			InitiatedBy:        o.taskCommonData.CreatedBy,
			Status:             o.taskCommonData.State,
			LastRunTime:        o.taskInstancesByIDMap[o.taskCommonData.ID].LastRunTime,
			CreatedAt:          o.taskCommonData.CreatedAt,
			ModifiedBy:         o.taskCommonData.ModifiedBy,
			ModifiedAt:         o.taskCommonData.ModifiedAt,
			NearestNextRunTime: nearestNextRunTime,
			TriggerTypes:       o.taskCommonData.Schedule.TriggerTypes,
		},
		InstanceSummaries: taskInstanceSummaries,
	}
	return taskSummaryDetails
}

func (taskService TaskService) buildTargetSummary(o resultsDTO, taskInstID gocql.UUID, endpointID gocql.UUID, status statuses.TaskInstanceStatus, nearestNextRunTime time.Time, taskInstance models.TaskInstance) (time.Time, models.TargetSummary) {
	var (
		statusMessage      string
		outputMessage      string
		currentNextRunTime time.Time
		emptyTime          time.Time
		currentExecResult  = o.execResByTaskInstIDAndMEIDMap[taskInstID][endpointID]
	)

	if o.mapTasksByEndpointID[endpointID].PostponedRunTime != emptyTime {
		currentNextRunTime = o.mapTasksByEndpointID[endpointID].PostponedRunTime
	} else {
		currentNextRunTime = o.mapTasksByEndpointID[endpointID].RunTimeUTC
	}

	switch status {
	case statuses.TaskInstanceSuccess:
		statusMessage = o.taskTemplate.SuccessMessage
		outputMessage = currentExecResult.StdOut
	case statuses.TaskInstanceFailed:
		statusMessage = o.taskTemplate.FailureMessage
		outputMessage = currentExecResult.StdErr
	}

	if nearestNextRunTime == emptyTime {
		nearestNextRunTime = currentNextRunTime
	}

	// find the nearest NextRunTime
	if currentNextRunTime.Before(nearestNextRunTime) && currentNextRunTime.After(time.Now().UTC()) {
		nearestNextRunTime = currentNextRunTime
	}

	return nearestNextRunTime, models.TargetSummary{
		InternalTaskState: o.mapTasksByEndpointID[endpointID].State,
		EndpointID:        endpointID,
		RunStatus:         status,
		StatusDetails:     statusMessage,
		Output:            outputMessage,
		OriginID:          taskInstance.OriginID.String(),
		NextRunTime:       o.mapTasksByEndpointID[endpointID].RunTimeUTC,
		LastRunTime:       currentExecResult.UpdatedAt,
		PostponedTime:     o.mapTasksByEndpointID[endpointID].PostponedRunTime,
		CanBePostponed:    (status == statuses.TaskInstanceScheduled || status == statuses.TaskInstancePostponed) && o.canBePostponedMap[taskInstance.TaskID],
		CanBeCanceled:     status == statuses.TaskInstanceScheduled && o.canBeCanceledMap[taskInstance.TaskID],
	}
}

func getAllUniqueTaskInstanceIDs(taskInstances []models.TaskInstance) (uniqueTaskInstanceIDs []gocql.UUID) {
	var (
		instanceIDs = make(map[gocql.UUID]struct{})
	)

	for _, ti := range taskInstances {
		if _, ok := instanceIDs[ti.ID]; !ok {
			uniqueTaskInstanceIDs = append(uniqueTaskInstanceIDs, ti.ID)
			instanceIDs[ti.ID] = struct{}{}
		}
	}

	return
}
