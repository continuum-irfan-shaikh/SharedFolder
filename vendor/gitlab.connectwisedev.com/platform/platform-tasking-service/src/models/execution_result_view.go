package models

//go:generate mockgen -destination=../mocks/mocks-gomock/executionResultViewPersistence_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models ExecutionResultViewPersistence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"
	agentModel "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	taskModel "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

type (
	// ExecutionResultView is defined by UI requirements.
	// See: https://continuum.atlassian.net/wiki/spaces/C2E/pages/223371940/Juno+Scripting+The+Aftermath?preview=/223371940/223387401/Confirmation_3.png7
	ExecutionResultView struct {
		PartnerID         string                      `json:"partnerId"`
		ManagedEndpointID gocql.UUID                  `json:"managedEndpointId"`
		TaskID            gocql.UUID                  `json:"taskId"`
		ExecutionID       gocql.UUID                  `json:"executionId"`
		OriginID          gocql.UUID                  `json:"originId"`
		TaskName          string                      `json:"taskName"`
		Type              string                      `json:"type"`
		Description       string                      `json:"description"`
		Regularity        taskModel.Regularity        `json:"regularity"`
		InitiatedBy       string                      `json:"initiatedBy"`
		Status            statuses.TaskState          `json:"status"`
		LastRunTime       time.Time                   `json:"lastRunTime"`
		LastRunStatus     statuses.TaskInstanceStatus `json:"lastRunStatus"`
		LastRunStdOut     string                      `json:"lastRunStdOut"`
		LastRunStdErr     string                      `json:"lastRunStdErr"`
		ResultMessage     string                      `json:"resultMessage"`
		DeviceCount       int                         `json:"deviceCount"`
		ModifiedBy        string                      `json:"modifiedBy"`
		CanBePostponed    bool                        `json:"canBePostponed"`
		CanBeCanceled     bool                        `json:"canBeCanceled"`
		NextRunTime       time.Time                   `json:"nextRunTime"`
	}

	// ExecutionResultViewPersistence represents a repository with TaskExecutionResults
	ExecutionResultViewPersistence interface {
		Get(ctx context.Context, partnerID string, endpointID gocql.UUID, count int, isNOCUser bool) (executionResultsView []*ExecutionResultView, err error)
		History(ctx context.Context, partnerID string, taskID, endpointID gocql.UUID, count int, isNOCUser bool) (executionResultsView []*ExecutionResultView, err error)
	}

	// ExecutionResultKafkaMessage structure of the Kafka message with Script execution results
	ExecutionResultKafkaMessage struct {
		agentModel.BrokerEnvelope
		Message taskModel.ScriptPluginReturnMessage `json:"message"`
	}

	// ExecutionResultViewRepoCassandra is a realisation of ExecutionResultViewPersistence interface for Cassandra
	ExecutionResultViewRepoCassandra struct{}
)

var (
	// ExecutionResultViewPersistenceInstance is an instance implemented ExecutionResultViewPersistence interface
	ExecutionResultViewPersistenceInstance ExecutionResultViewPersistence = ExecutionResultViewRepoCassandra{}
)

// Get returns executionResults for the specified Partner and endpoint ID
func (e ExecutionResultViewRepoCassandra) Get(ctx context.Context, partnerID string, endpointID gocql.UUID, count int, isNOCUser bool) ([]*ExecutionResultView, error) {
	listOfTasks, err := TaskPersistenceInstance.GetByPartnerAndManagedEndpointID(ctx, partnerID, endpointID, count)
	if err != nil {
		return nil, err
	}

	var lastTaskInstanceIDs = make([]gocql.UUID, len(listOfTasks))
	for i, task := range listOfTasks {
		lastTaskInstanceIDs[i] = task.LastTaskInstanceID
	}

	var (
		executionResultsView                   []*ExecutionResultView
		executionResultsErr, tasksInstancesErr error
		executionResultsByTaskInstanceIDMap    = make(map[gocql.UUID]ExecutionResult)
		taskInstancesByIDMap                   = make(map[gocql.UUID]TaskInstance)
		wg                                     = &sync.WaitGroup{}
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		var executionResults []ExecutionResult
		executionResults, executionResultsErr = ExecutionResultPersistenceInstance.GetByTargetAndTaskInstanceIDs(endpointID, lastTaskInstanceIDs...)
		if executionResultsErr != nil {
			return
		}

		for _, result := range executionResults {
			executionResultsByTaskInstanceIDMap[result.TaskInstanceID] = result
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var tasksInstances []TaskInstance
		tasksInstances, tasksInstancesErr = TaskInstancePersistenceInstance.GetByIDs(ctx, lastTaskInstanceIDs...)
		if len(tasksInstances) < 1 {
			return
		}

		for _, taskInstance := range tasksInstances {
			taskInstancesByIDMap[taskInstance.ID] = taskInstance
		}
	}()
	wg.Wait()

	if executionResultsErr != nil {
		logger.Log.WarnfCtx(ctx,"GetByTargetAndTaskInstanceIDs: error during retrieving Execution Results by EndpointID %v. Err: %v", endpointID, executionResultsErr)
	}
	if tasksInstancesErr != nil {
		logger.Log.WarnfCtx(ctx, "GetByIDs: error during retrieving TaskInstances by lastTaskInstanceIDs %v. Err: %v", lastTaskInstanceIDs, tasksInstancesErr)
	}
	for _, task := range listOfTasks {
		executionResultView := e.getResultByTask(isNOCUser, task, taskInstancesByIDMap, endpointID, executionResultsByTaskInstanceIDMap)
		if executionResultView == nil {
			continue
		}

		executionResultsView = appendExecutionResultsView(ctx, executionResultsView, executionResultView, task.OriginID, partnerID)
	}
	return executionResultsView, nil
}

func (e ExecutionResultViewRepoCassandra) getResultByTask(isNOCUser bool, task Task, taskInstancesByIDMap map[gocql.UUID]TaskInstance, endpointID gocql.UUID, executionResultsByTaskInstanceIDMap map[gocql.UUID]ExecutionResult) *ExecutionResultView {
	// NOC user accesses checking
	if !isNOCUser && task.IsRequireNOCAccess {
		return nil
	}

	var executionResult ExecutionResult
	var lastTaskInstance TaskInstance
	var ok bool

	if lastTaskInstance, ok = taskInstancesByIDMap[task.LastTaskInstanceID]; !ok {
		return nil
	}

	if status, ok := lastTaskInstance.Statuses[endpointID]; ok {
		executionResult.ExecutionStatus = status
	}

	if result, ok := executionResultsByTaskInstanceIDMap[task.LastTaskInstanceID]; ok {
		executionResult = result
	}

	executionResultView := newExecutionResultView(&task, &executionResult, endpointID)
	if ti, ok := taskInstancesByIDMap[task.LastTaskInstanceID]; ok {
		executionResultView.DeviceCount = len(ti.Statuses)
	}

	if task.State == statuses.TaskStateActive &&
		(executionResult.ExecutionStatus == statuses.TaskInstanceScheduled ||
			executionResult.ExecutionStatus == statuses.TaskInstancePostponed) && !task.IsTrigger() && !task.IsTaskAndTriggerNotActivated() {
		executionResultView.CanBePostponed = true
	}

	if task.State == statuses.TaskStateActive &&
		executionResult.ExecutionStatus == statuses.TaskInstanceScheduled &&
		!task.IsTrigger() {
		executionResultView.CanBeCanceled = true
	}
	return executionResultView
}

// History returns a history of execution results on the specified Task for the specified Partner
func (ExecutionResultViewRepoCassandra) History(ctx context.Context, partnerID string, taskID, endpointID gocql.UUID, count int, isNOCUser bool) ([]*ExecutionResultView, error) {
	tasks, err := TaskPersistenceInstance.GetByIDAndManagedEndpoints(ctx, partnerID, taskID, endpointID)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("task not found, %v id", taskID)
	}

	//NOC user accesses checking
	if !isNOCUser && tasks[0].IsRequireNOCAccess {
		return []*ExecutionResultView{}, nil
	}

	taskForEndpoint := tasks[0]
	taskInstances, err := TaskInstancePersistenceInstance.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	instancesCount := 0
	taskInstancesIDs := make([]gocql.UUID, 0, len(taskInstances))
	for _, taskInstance := range taskInstances {
		if instancesCount >= count && count != common.UnlimitedCount {
			break
		}

		if _, ok := taskInstance.Statuses[endpointID]; !ok {
			continue
		}

		taskInstancesIDs = append(taskInstancesIDs, taskInstance.ID)
		instancesCount++
	}

	executionResults, err := ExecutionResultPersistenceInstance.GetByTargetAndTaskInstanceIDs(endpointID, taskInstancesIDs...)
	if err != nil {
		return nil, fmt.Errorf("ExecutionResultViewRepoCassandra.History: Error while trying to retrive Script Execution Results by Task Instance ID and Managed Endpoint ID %v. Err: %v",
			taskID, err)
	}

	var executionResultsView = make([]*ExecutionResultView, 0)
	for _, executionResult := range executionResults {
		executionResultView := newExecutionResultView(&taskForEndpoint, &executionResult, endpointID)
		executionResultsView = appendExecutionResultsView(ctx, executionResultsView, executionResultView, taskForEndpoint.OriginID, partnerID)
	}

	return executionResultsView, nil
}

// GetResultMessage returns the result message as a string
func GetResultMessage(ctx context.Context, partnerID string, originID gocql.UUID, status statuses.TaskInstanceStatus) (string, error) {
	template, err := TemplateCacheInstance.GetByOriginID(ctx, partnerID, originID, true)
	if err != nil {
		return "", err
	}

	switch status {
	case statuses.TaskInstanceFailed:
		return template.FailureMessage, nil
	case statuses.TaskInstanceSuccess:
		return template.SuccessMessage, nil
	default:
		return "", nil
	}
}

func appendExecutionResultsView(ctx context.Context, executionResultsView []*ExecutionResultView, executionResultView *ExecutionResultView, originID gocql.UUID, partnerID string) []*ExecutionResultView {
	resultMessage, err := GetResultMessage(ctx, partnerID, originID, executionResultView.LastRunStatus)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "can't find templates in cache by originID %v. Err: %v", originID, err)
	}
	executionResultView.ResultMessage = resultMessage
	return append(executionResultsView, executionResultView)
}

func newExecutionResultView(task *Task, execResult *ExecutionResult, endpointID gocql.UUID) *ExecutionResultView {
	return &ExecutionResultView{
		PartnerID:         task.PartnerID,
		ManagedEndpointID: endpointID,
		TaskID:            task.ID,
		ExecutionID:       task.LastTaskInstanceID,
		OriginID:          task.OriginID,
		TaskName:          task.Name,
		Type:              task.Type,
		Description:       task.Description,
		Regularity:        task.Schedule.Regularity,
		InitiatedBy:       task.CreatedBy,
		ModifiedBy:        task.ModifiedBy,
		Status:            task.State,
		NextRunTime:       task.RunTimeUTC,
		LastRunTime:       execResult.UpdatedAt,
		LastRunStatus:     execResult.ExecutionStatus,
		LastRunStdOut:     execResult.StdOut,
		LastRunStdErr:     execResult.StdErr,
	}
}
