package models

//go:generate mockgen -destination=../mocks/mocks-gomock/taskSummaryPersistence_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models TaskSummaryPersistence

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

type (
	// TaskSummaryData stores info for tasking summary page data:
	// Name, state, regularity, created_by are related to particular Task
	// RunOn represents the amount of devices the task was run on
	// LastRunTime represents a time LastTaskInstance was started at for particular task
	TaskSummaryData struct {
		Name               string               `json:"name"`
		TaskID             gocql.UUID           `json:"taskID"`
		Type               string               `json:"type"`
		RunOn              TargetData           `json:"runOn"`
		Regularity         apiModels.Regularity `json:"regularity"`
		InitiatedBy        string               `json:"initiatedBy"`
		LastRunTime        time.Time            `json:"lastRunTime"`
		CreatedAt          time.Time            `json:"createdAt"`
		Status             statuses.TaskState   `json:"status"`
		ModifiedAt         time.Time            `json:"modifiedAt"`
		ModifiedBy         string               `json:"modifiedBy"`
		NearestNextRunTime time.Time            `json:"nearestNextRunTime"`
		TriggerTypes       []string             `json:"triggerTypes"`
	}

	// TaskInstanceSummary stores info about Task Instance and a slice of ManagedEndpointID Summaries
	TaskInstanceSummary struct {
		ID              gocql.UUID      `json:"taskInstanceID"`
		RunTime         time.Time       `json:"runTime"`
		TargetSummaries []TargetSummary `json:"targetSummaries"`
		RunStatuses     DeviceStatuses  `json:"deviceStatuses"`
	}

	// DeviceStatuses stores devices count and map with statuses
	DeviceStatuses struct {
		DeviceCount int            `json:"deviceCount"`
		Statuses    map[string]int `json:"statuses"`
	}

	// LastRunStatusData stores info about amount of devices and its succeeded or failed execution results for task
	LastRunStatusData struct {
		Status       statuses.TaskInstanceStatus `json:"status"`
		DeviceCount  int                         `json:"deviceCount"`  // Total count of the devices the task was run on
		SuccessCount int                         `json:"successCount"` // Total count of succeeded results among all devices
		FailureCount int                         `json:"failureCount"` // Total count of failed results among all devices
	}

	// TargetSummary stores summary info for particular endpoint
	TargetSummary struct {
		InternalTaskState statuses.TaskState          `json:"internalTaskState"`
		EndpointID        gocql.UUID                  `json:"endpointID"`
		RunStatus         statuses.TaskInstanceStatus `json:"runStatus"`
		StatusDetails     string                      `json:"statusDetails"`
		Output            string                      `json:"output"`
		OriginID          string                      `json:"originID"`       // represents ID for any external entity
		NextRunTime       time.Time                   `json:"nextRunTime"`    // represents time when scheduled task will be running the next time
		LastRunTime       time.Time                   `json:"lastRunTime"`    // represents time when scheduled task was run last time on particular device
		CanBePostponed    bool                        `json:"canBePostponed"` // represents possibility to postpone a task
		CanBeCanceled     bool                        `json:"canBeCanceled"`  // represents possibility to cancel a task
		PostponedTime     time.Time                   `json:"postponedTime"`
	}

	// TaskSummaryDetails stores Task Summary for the particular Task and the its Instance Summaries
	TaskSummaryDetails struct {
		TaskSummary       TaskSummaryData       `json:"taskSummary"`
		InstanceSummaries []TaskInstanceSummary `json:"instanceSummary"`
	}
)

// TaskSummaryPersistence interface to perform actions to get TaskingSummaryPage data
type TaskSummaryPersistence interface {
	GetTasksSummaryData(context.Context, bool, persistency.Cache, string, ...gocql.UUID) ([]TaskSummaryData, error)
	UpdateTaskInstanceStatusCount(context.Context, gocql.UUID, int, int) error
	GetStatusCountsByIDs(ctx context.Context, cache persistency.Cache, taskInstancesMapByID map[gocql.UUID]TaskInstance, taskInstanceIDs []gocql.UUID) (map[gocql.UUID]TaskInstanceStatusCount, error)
}

// TaskSummaryRepoCassandra is a realisation of TaskSummaryPersistence interface for Cassandra
type TaskSummaryRepoCassandra struct{}

var (
	// TaskSummaryPersistenceInstance is an instance presented TaskSummaryRepoCassandra
	TaskSummaryPersistenceInstance TaskSummaryPersistence = TaskSummaryRepoCassandra{}
)

// GetTasksSummaryData get the tasking summary page data
func (TaskSummaryRepoCassandra) GetTasksSummaryData(ctx context.Context, isNOCUser bool, cache persistency.Cache, partnerID string, taskIDs ...gocql.UUID) ([]TaskSummaryData, error) {
	var (
		err              error
		tasks            []Task
		tasksSummaryData = make([]TaskSummaryData, 0)
	)

	if len(taskIDs) == 0 {
		tasks, err = TaskPersistenceInstance.GetByPartner(ctx, partnerID)
	} else {
		tasks, err = TaskPersistenceInstance.GetByIDs(ctx, nil, partnerID, false, taskIDs...)
	}
	if err != nil {
		return tasksSummaryData, err
	}

	var filteredTasks []Task
	// filter by NOC accesses
	for _, task := range tasks {
		if !isNOCUser && task.IsRequireNOCAccess {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}

	// Extract LastTaskInstanceIDs from tasks into slice
	// to use them later for getting all the latest TaskInstances created for these tasks
	// Also LastTaskInstanceIDs are used for getting all TaskInstancesStatusCounts
	lastTaskInstanceIDs := getLastTaskInstanceIDs(filteredTasks)

	// Building all Last Task Instances related to all Tasks by partner into slice
	lastTaskInstances, err := TaskInstancePersistenceInstance.GetByIDs(ctx, lastTaskInstanceIDs...)
	if err != nil {
		return tasksSummaryData, err
	}

	// Grouping TaskInstances by its' IDs to the map
	// so that we can easily get LastRunTime
	// by getting StartedAt field's value from specific TaskInstance
	mapTaskInstancesByID := groupTaskInstancesByID(lastTaskInstances)

	taskInstancesMapByID := GroupTaskInstancesByTaskInstanceID(lastTaskInstances)

	// Getting the slice of all Task Instance Status Counts by all lastTaskInstanceIDs
	taskInstanceStatusCountsMap, err := TaskSummaryPersistenceInstance.GetStatusCountsByIDs(ctx, cache, taskInstancesMapByID, lastTaskInstanceIDs)
	if err != nil {
		return tasksSummaryData, err
	}

	// Grouping Tasks by TaskID to map
	// in order to be able to get GlobalTaskState overall each group of tasks
	mapTasksByID := groupTasksByID(filteredTasks)

	// Calculating and grouping global TaskState by taskID
	mapGlobalTaskStateByTasksID := groupGlobalTaskStateByTasksID(mapTasksByID)

	// Grouping LastTaskInstances by TaskID to map
	// in order to get all statusCounts related to specific taskID
	mapLastTaskInstancesByTaskID := groupLastTaskInstancesByTaskID(filteredTasks, mapTaskInstancesByID)

	// Grouping LastRunStatusData by TaskID to map
	// using intersections of TaskInstances and TaskInstanceStatusCounts by Task Instance IDs
	mapLastRunStatusDataByTaskID := groupLastRunStatusDataByTaskID(
		taskInstanceStatusCountsMap,
		mapTasksByID,
		mapGlobalTaskStateByTasksID,
		mapLastTaskInstancesByTaskID,
	)

	for taskID /*lastRunStatusData*/ := range mapLastRunStatusDataByTaskID {
		// Name, Regularity, InitiatedBy and LastTaskInstanceID fields are equal among all the group of tasks
		sharedTaskData := mapTasksByID[taskID][0]

		tasksSummaryData = append(tasksSummaryData, TaskSummaryData{
			Name:   sharedTaskData.Name,
			TaskID: sharedTaskData.ID,
			Type:   sharedTaskData.Type,
			RunOn: TargetData{
				Count: len(filteredTasks),
			},
			Regularity:  sharedTaskData.Schedule.Regularity,
			InitiatedBy: sharedTaskData.CreatedBy,
			Status:      mapGlobalTaskStateByTasksID[taskID],
			LastRunTime: mapTaskInstancesByID[sharedTaskData.LastTaskInstanceID].LastRunTime,
			CreatedAt:   sharedTaskData.CreatedAt,
			ModifiedBy:  sharedTaskData.ModifiedBy,
			ModifiedAt:  sharedTaskData.ModifiedAt,
		})
	}

	return tasksSummaryData, nil
}

// Extract all Last Task Instance IDs from previously selected tasks to a slice of interfaces
// in order to get Task Instances and Task Instance Status Counts from Cassandra
func getLastTaskInstanceIDs(tasks []Task) (lastTaskInstanceIDs []gocql.UUID) {
	var instanceIDs = make(map[gocql.UUID]struct{})
	for _, task := range tasks {
		if task.LastTaskInstanceID.String() == "00000000-0000-0000-0000-000000000000" {
			continue
		}
		instanceIDs[task.LastTaskInstanceID] = struct{}{}
	}
	for id := range instanceIDs {
		lastTaskInstanceIDs = append(lastTaskInstanceIDs, id)
	}
	return lastTaskInstanceIDs
}

// Grouping Last Task Instances by their IDs to map
// in order to get the field StartedAt as LastRunTime
// by passing a LastTaskInstanceID as a key
// Returns map[TaskInstanceID]TaskInstance
func groupTaskInstancesByID(taskInstances []TaskInstance) map[gocql.UUID]TaskInstance {
	var mapTaskInstancesByID = make(map[gocql.UUID]TaskInstance)

	for _, taskInstance := range taskInstances {
		mapTaskInstancesByID[taskInstance.ID] = taskInstance
	}
	return mapTaskInstancesByID
}

// Grouping slices of Task Internal to the map by their Task IDs
// in order to calculate Global Task State of the tasks with the same ID
// by passing the Task ID to this map as a key
// Returns map[TaskID][]Task
func groupTasksByID(tasks []Task) (mapTaskByID map[gocql.UUID][]Task) {
	mapTaskByID = make(map[gocql.UUID][]Task)
	for _, task := range tasks {
		mapTaskByID[task.ID] = append(mapTaskByID[task.ID], task)
	}
	return
}

// Calculating and grouping global TaskState by taskID to the map
func groupGlobalTaskStateByTasksID(mapTaskByID map[gocql.UUID][]Task) map[gocql.UUID]statuses.TaskState {
	mapGlobalTaskStateByTasksID := make(map[gocql.UUID]statuses.TaskState)
	for taskID, tasksGroup := range mapTaskByID {
		mapGlobalTaskStateByTasksID[taskID] = GetGlobalTaskState(tasksGroup)
	}
	return mapGlobalTaskStateByTasksID
}

// Grouping Last Task Instances by Task ID to map
// in order to get grouped Last Run Status Data by Task ID
func groupLastTaskInstancesByTaskID(tasks []Task, taskInstancesByID map[gocql.UUID]TaskInstance) map[gocql.UUID][]TaskInstance {
	var (
		mapUniqueTaskInstanceIDsByTaskID = make(map[gocql.UUID][]gocql.UUID)
		mapTaskInstancesByTaskID         = make(map[gocql.UUID][]TaskInstance)
		emptyUUID                        gocql.UUID
	)
	for _, task := range tasks {
		if task.LastTaskInstanceID == emptyUUID || common.UUIDSliceContainsElement(mapUniqueTaskInstanceIDsByTaskID[task.ID], task.LastTaskInstanceID) {
			continue
		}
		mapUniqueTaskInstanceIDsByTaskID[task.ID] = append(mapUniqueTaskInstanceIDsByTaskID[task.ID], task.LastTaskInstanceID)
		mapTaskInstancesByTaskID[task.ID] = append(mapTaskInstancesByTaskID[task.ID], taskInstancesByID[task.LastTaskInstanceID])
	}
	return mapTaskInstancesByTaskID
}

// Grouping Last Run Status Data by Task ID to map
// in order to iterate through it and to form the slice of Task Summary Data
// Returns map[TaskID]LastRunStatusData
func groupLastRunStatusDataByTaskID(
	statusCountsByTaskInstanceID map[gocql.UUID]TaskInstanceStatusCount,
	tasksByID map[gocql.UUID][]Task,
	mapGlobalTaskStateByTasksID map[gocql.UUID]statuses.TaskState,
	mapLastTaskInstancesByTaskID map[gocql.UUID][]TaskInstance,
) map[gocql.UUID]LastRunStatusData {

	var results = make(map[gocql.UUID]LastRunStatusData)

	for id := range tasksByID {
		results[id] = buildLastRunStatusData(statusCountsByTaskInstanceID, mapGlobalTaskStateByTasksID[id], mapLastTaskInstancesByTaskID[id])
	}
	return results
}

func buildLastRunStatusData(
	statusCounts map[gocql.UUID]TaskInstanceStatusCount,
	taskState statuses.TaskState,
	lastTaskInstances []TaskInstance,
) (lastRunStatusData LastRunStatusData) {
	// So the task has not been run yet
	if len(lastTaskInstances) == 0 {
		lastRunStatusData.Status = statuses.CalculateStatusForNotStartedTask(taskState)
		return
	}

	// lastTaskInstances contains only unique elements
	for _, taskInstance := range lastTaskInstances {
		lastRunStatusData.DeviceCount += len(taskInstance.Statuses)
		lastRunStatusData.SuccessCount += statusCounts[taskInstance.ID].SuccessCount
		lastRunStatusData.FailureCount += statusCounts[taskInstance.ID].FailureCount
	}

	lastRunStatusData.Status = statuses.CalculateForStartedTask(
		lastRunStatusData.DeviceCount,
		lastRunStatusData.SuccessCount,
		lastRunStatusData.FailureCount,
	)
	return
}

// GroupTaskInstancesByTaskID groups Task instances by taskID to map
// Return map[TaskID]TaskInstance
func GroupTaskInstancesByTaskID(taskInstances []TaskInstance) map[gocql.UUID][]TaskInstance {
	result := make(map[gocql.UUID][]TaskInstance)
	for _, inst := range taskInstances {
		result[inst.TaskID] = append(result[inst.TaskID], inst)
	}

	return result
}

// GroupTaskInstancesByTaskInstanceID groups Task instances by taskInstanceID to map
// Return map[TaskID]TaskInstance
func GroupTaskInstancesByTaskInstanceID(taskInstances []TaskInstance) map[gocql.UUID]TaskInstance {
	result := make(map[gocql.UUID]TaskInstance)
	for _, inst := range taskInstances {
		result[inst.ID] = inst
	}
	return result
}
