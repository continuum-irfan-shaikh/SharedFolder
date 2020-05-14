package entities

import (
	"fmt"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	. "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

//EndpointsInput - list of partner endpoints
type EndpointsInput []string

//EndpointsClosestTasks - closest tasks by endpoint data representation
type EndpointsClosestTasks map[string]ClosestTasks

//ClosestTasks - closest tasks container
type ClosestTasks struct {
	Previous *ClosestTask `json:"previous,omitempty"`
	Next     *ClosestTask `json:"next,omitempty"`
}

//ClosestTask - closest task details
type ClosestTask struct {
	ID      string `json:"-"`
	Name    string `json:"name,omitempty"`
	RunDate int64  `json:"runDate,omitempty"`
	Status  string `json:"status,omitempty"`
}

//Task - Task database data representation
type Task struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	RunTimeUTC        time.Time          `json:"runTimeUTC"`
	PostponedRunTime  time.Time          `json:"postponedTime"`
	ManagedEndpointID string             `json:"managedEndpointID"`
	State             TaskState          `json:"state"`
	Schedule          apiModels.Schedule `json:"schedule"`
}

//TaskInstance - TaskInstance database data representation
type TaskInstance struct {
	ID            string                        `json:"-"`
	TaskID        string                        `json:"taskId"`
	StartedAt     time.Time                     `json:"startedAt"`
	LastRunTime   time.Time                     `json:"lastRunTime"`
	Statuses      map[string]TaskInstanceStatus `json:"statuses"`
	TaskName      string                        `json:"name"`
	PartnerID     string                        `json:"-"`
	FailureCount  int                           `json:"-"`
	SuccessCount  int                           `json:"-"`
	StatusesCount map[TaskInstanceStatus]int    `json:"-"`
}

// PreparePendingStatuses - makes scheduled endpoints pending if running endpoints are existent
func (ti *TaskInstance) PreparePendingStatuses() {
	var hasRunning bool
	for _, s := range ti.Statuses {
		if s == TaskInstanceRunning || s == TaskInstanceSuccess || s == TaskInstanceFailed {
			hasRunning = true
			break
		}
	}

	if !hasRunning {
		return
	}

	for id, s := range ti.Statuses {
		if s == TaskInstanceScheduled {
			ti.Statuses[id] = TaskInstancePending
		}
	}
}

// FillStatusCount fills status counts map with given statuses
func (ti *TaskInstance) FillStatusCount() {
	var sc = make(map[TaskInstanceStatus]int)
	for _, stat := range ti.Statuses {
		sc[stat]++
	}
	ti.StatusesCount = sc
}

// CalculateOverallStatus returns a common status for the entire instance. isLastTaskInstance flag tells if we need
// to check if it first execution or the regular one
func (ti *TaskInstance) CalculateOverallStatus() (os OverallStatus) {
	ti.FillStatusCount()

	if len(ti.Statuses) < 1 {
		return OverallNew
	}

	deviceCount := len(ti.Statuses)
	if ti.StatusesCount[TaskInstanceRunning] != 0 || ti.StatusesCount[TaskInstancePending] != 0 {
		return OverallRunning
	}

	if ti.StatusesCount[TaskInstanceScheduled] == deviceCount {
		return OverallNew
	}

	if ti.StatusesCount[TaskInstanceSuccess] != 0 && ti.StatusesCount[TaskInstanceFailed] != 0 {
		return OverallPartialFailed
	}

	if ti.StatusesCount[TaskInstancePostponed] == deviceCount ||
		ti.StatusesCount[TaskInstanceDisabled] == deviceCount ||
		ti.StatusesCount[TaskInstanceCanceled] == deviceCount {
		return OverallSuspended
	}

	if ti.StatusesCount[TaskInstanceSuccess] == deviceCount {
		return OverallSuccess
	}

	if ti.StatusesCount[TaskInstanceFailed] == deviceCount || ti.StatusesCount[TaskInstanceSomeFailures] == deviceCount {
		return OverallFailed
	}

	// if there are some other statuses and there is at least failed one - it failed
	if ti.StatusesCount[TaskInstanceFailed] != 0 || ti.StatusesCount[TaskInstanceSomeFailures] != 0 {
		return OverallPartialFailed
	}

	if ti.StatusesCount[TaskInstanceCanceled] != 0 {
		return OverallSuspended
	}

	return ""
}

//ExecutionResult - ScriptExecutionResult database table representation
type ExecutionResult struct {
	ManagedEndpointID string             `json:"managedEndpointID"`
	TaskInstanceID    string             `json:"taskInstanceID"`
	ExecutionStatus   TaskInstanceStatus `json:"executionStatus"`
	UpdatedAt         time.Time          `json:"updatedAt"`
}

//LastExecution - last_task_executions database table representation
type LastExecution struct {
	PartnerID  string
	EndpointID string
	RunTime    time.Time
	Name       string
	Status     TaskInstanceStatus
}

// CalculateStatuses returns a common status for the entire instance
func (ti *TaskInstance) CalculateStatuses() (map[string]int, error) {
	if len(ti.Statuses) < 1 {
		return nil, fmt.Errorf("empty Statuses")
	}

	var statusCounts = make(map[string]int)
	for _, stat := range ti.Statuses {
		statusStr, err := TaskInstanceStatusText(stat)
		if err != nil {
			return nil, err
		}
		statusCounts[statusStr]++
	}

	return statusCounts, nil
}
