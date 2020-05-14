package tasking

import "time"

// TaskSummaryData represents task summary data in TasksAndSequences page
type TaskSummaryData struct {
	TaskID             string            `json:"taskID"`
	Name               string            `json:"name"`
	Type               string            `json:"type"`
	CreatedAt          time.Time         `json:"createdAt"`
	RunOn              RunOnData         `json:"runOn"`
	Regularity         string            `json:"regularity"`
	InitiatedBy        string            `json:"initiatedBy"`
	Status             string            `json:"status"`
	LastRunTime        time.Time         `json:"lastRunTime"`
	LastRunStatus      LastRunStatusData `json:"lastRunStatus"`
	NextRunTime        time.Time         `json:"nextRunTime"`
	ModifiedBy         string            `json:"modifiedBy"`
	ModifiedAt         time.Time         `json:"modifiedAt"`
	NearestNextRunTime time.Time         `json:"nearestNextRunTime"`
	TriggerTypes       []string          `json:"triggerTypes"`
}
