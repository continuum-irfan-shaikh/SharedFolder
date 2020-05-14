package entities

import (
	"time"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

// Here are described originIDs of custom scripts
const (
	CustomPowershell = "51a74346-e19b-11e7-9809-0800279505d9"
	CustomCMD        = "e3d2c26b-c5ba-49cf-a089-7637f6de949e"
	CustomBash       = "37f7f19f-40e8-11e9-a643-e0d55e1ce78a"
)

// ScheduledTasks represents scheduled tasks data returned by scheduled tasks request handler
type ScheduledTasks struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	OverallStatus      statuses.OverallStatus `json:"status"`
	LastRunTime        time.Time              `json:"lastRunTime"`
	RunTimeUTC         time.Time              `json:"nextRunTime"`
	Description        string                 `json:"description"`
	CreatedBy          string                 `json:"createdBy"`
	CreatedAt          time.Time              `json:"createdAt"`
	ModifiedBy         string                 `json:"modifiedBy"`
	ModifiedAt         time.Time              `json:"modifiedAt"`
	ExecutionInfo      ExecutionInfo          `json:"executionInfo"`
	TaskType           string                 `json:"taskType"`
	Regularity         tasking.Regularity     `json:"regularity"`
	CanBeCanceled      bool                   `json:"canBeCanceled"` // this field tells if the task can be canceled or not
	TriggerFrames      []tasking.TriggerFrame `json:"-"`
	OriginID           string                 `json:"-"`
	State              statuses.TaskState     `json:"-"`
	IsNOC              bool                   `json:"-"`
	LastTaskInstanceID string                 `json:"-"`
	TriggerTypes       []string               `json:"-"`
	PostponedTime      time.Time              `json:"-"`
}

// ExecutionInfo  represents instance execution info
type ExecutionInfo struct {
	DeviceCount  int `json:"deviceCount"`
	SuccessCount int `json:"successCount"`
	FailedCount  int `json:"failedCount"`
}

// TaskIDs represents scheduled tasks ids list
type TaskIDs struct {
	IDs []string `json:"ids"`
}
