package statuses

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TaskInstanceStatus type is used for status definition
type TaskInstanceStatus int

// These constants describe status of the Execution Result
const (
	_ TaskInstanceStatus = iota
	TaskInstanceRunning
	TaskInstanceSuccess
	TaskInstanceFailed
	TaskInstanceScheduled
	TaskInstanceDisabled
	TaskInstanceSomeFailures
	TaskInstancePending
	TaskInstanceStopped
	TaskInstancePostponed
	TaskInstanceCanceled
)

// These are instances statuses
const (
	TaskInstanceRunningText      = "Running"
	TaskInstanceSuccessText      = "Success"
	TaskInstanceFailedText       = "Failed"
	TaskInstanceScheduledText    = "Scheduled"
	TaskInstanceDisabledText     = "Disabled"
	TaskInstanceSomeFailuresText = "Some Failures"
	TaskInstancePendingText      = "Pending"
	TaskInstanceStoppedText      = "Stopped"
	TaskInstancePostponedText    = "Postponed"
	TaskInstanceCanceledText     = "Canceled"
)

// taskInstanceStatusText contains the text representation of the TaskInstanceStatus
var taskInstanceStatusText = map[TaskInstanceStatus]string{
	TaskInstanceStatus(0):    "",
	TaskInstanceRunning:      TaskInstanceRunningText,
	TaskInstanceSuccess:      TaskInstanceSuccessText,
	TaskInstanceFailed:       TaskInstanceFailedText,
	TaskInstanceScheduled:    TaskInstanceScheduledText,
	TaskInstanceDisabled:     TaskInstanceDisabledText,
	TaskInstanceSomeFailures: TaskInstanceSomeFailuresText,
	TaskInstancePending:      TaskInstancePendingText,
	TaskInstanceStopped:      TaskInstanceStoppedText,
	TaskInstancePostponed:    TaskInstancePostponedText,
	TaskInstanceCanceled:     TaskInstanceCanceledText,
}

// TaskInstanceStatusText converts a value of the TaskInstanceStatus to its string representation
// It returns an error if there is no such status.
func TaskInstanceStatusText(status TaskInstanceStatus) (statusText string, err error) {
	statusText, ok := taskInstanceStatusText[status]
	if !ok {
		err = fmt.Errorf("incorrect Task Instance Status: %v", status)
	}
	return
}

// taskInstanceStatuses contains the mapping of the text representation of the TaskInstanceStatus and its value
var taskInstanceStatuses = map[string]TaskInstanceStatus{
	TaskInstanceRunningText:      TaskInstanceRunning,
	TaskInstanceSuccessText:      TaskInstanceSuccess,
	TaskInstanceFailedText:       TaskInstanceFailed,
	TaskInstanceScheduledText:    TaskInstanceScheduled,
	TaskInstanceDisabledText:     TaskInstanceDisabled,
	TaskInstanceSomeFailuresText: TaskInstanceSomeFailures,
	TaskInstancePendingText:      TaskInstancePending,
	TaskInstanceStoppedText:      TaskInstanceStopped,
	TaskInstanceCanceledText:     TaskInstanceCanceled,
	TaskInstancePostponedText:    TaskInstancePostponed,
}

// TaskInstanceStatusFromText converts a string value of the Task Instance Status to the TaskInstanceStatus.
// It returns an error if there is no status for the statusText.
func TaskInstanceStatusFromText(statusText string) (status TaskInstanceStatus, err error) {
	if len(statusText) <= 2 {
		return status, fmt.Errorf("status can't be empty, and should have more than 2 symbols")
	}

	statusText = strings.ToLower(statusText)
	// first case must be upper as this is used on UI as well
	statusText = fmt.Sprintf("%s%s", strings.ToUpper(statusText[0:1]), statusText[1:])

	status, ok := taskInstanceStatuses[statusText]
	if !ok {
		err = fmt.Errorf("incorrect Task Instance Status: %s", statusText)
	}
	return
}

// UnmarshalJSON used to convert string Status representation to TaskInstanceStatus type
func (status *TaskInstanceStatus) UnmarshalJSON(byteResult []byte) error {
	var (
		stringValue string
		err         error
	)
	if err = json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}

	*status, err = TaskInstanceStatusFromText(stringValue)
	return err
}

// MarshalJSON custom marshal method for field TaskInstanceStatus
func (status TaskInstanceStatus) MarshalJSON() ([]byte, error) {
	statusString, err := TaskInstanceStatusText(status)
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(statusString)
}

// CalculateStatusForNotStartedTask returns TaskInstanceStatus based on the taskState
func CalculateStatusForNotStartedTask(taskState TaskState) TaskInstanceStatus {
	switch taskState {
	case TaskStateDisabled:
		return TaskInstanceDisabled
	case TaskStateActive:
		return TaskInstanceScheduled
	}
	return TaskInstanceStatus(0)
}

// CalculateForStartedTask returns TaskInstanceStatus based on the info about specific execution of the task
func CalculateForStartedTask(deviceCount, successCount, failureCount int) TaskInstanceStatus {
	// calculations work as defined in https://continuum.atlassian.net/wiki/spaces/C2E/pages/301400183/Task+Global+and+Individual+Statuses
	switch deviceCount {
	case failureCount:
		return TaskInstanceFailed
	case successCount:
		return TaskInstanceSuccess
	case failureCount + successCount:
		return TaskInstanceSomeFailures
	default:
		return TaskInstanceRunning
	}
}
