package migrateModels

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

// TaskState type is used for states definition
type TaskState int

// RunTime alias for JSON custom marshaling/unmarshaling
type RunTime int64

// These constants describe states of the task
const (
	_ TaskState = iota
	TaskStateScheduled
	TaskStateRunning
	TaskStateCompleted
)

// Regularity type is used for regularity definition
type Regularity int

// These constants describe regularity of the task
const (
	_ Regularity = iota
	RunNow
	OneTime
	Recurrent
)

type (
	//Task contains common data for all tasks
	Task struct {
		ID           gocql.UUID `json:"id"`
		Name         string     `json:"name"`
		Description  string     `json:"description"`
		Schedule     string     `json:"schedule"`
		CreatedAt    time.Time  `json:"createdAt"`
		CreatedBy    string     `json:"createdBy"`
		PartnerID    string     `json:"partnerId"`
		OriginID     gocql.UUID `json:"originId"`
		State        TaskState  `json:"state"`
		Regularity   Regularity `json:"regularity"`
		RunTime      time.Time  `json:"runTime"`
		StartRunTime time.Time  `json:"startRunTime"`
		EndRunTime   time.Time  `json:"endRunTime"`
		Trigger      string     `json:"trigger"`
		Type         string     `json:"type"`
		Parameters   string     `json:"parameters"`
	}

	//ObsoleteTask contains tasks data in obsolete format
	ObsoleteTask struct {
		Task
		TaskTargets map[string]bool `json:"taskTargets"`
	}

	//NewTask contains tasks data in new format
	NewTask struct {
		Task
		Target           string    `json:"targets"`
		RunTimeUnix      time.Time `json:"run_time_unix"`
		TaskLocation     string    `json:"location,omitempty"`
		RunTimeInSeconds int64     `json:"-"`
	}
)

// UnmarshalJSON used to convert string State representation to State type
func (state *TaskState) UnmarshalJSON(byteResult []byte) error {
	var stringValue string
	if err := json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}
	switch stringValue {
	case "":
		*state = 0
	case "Running":
		*state = TaskStateRunning
	case "Scheduled":
		*state = TaskStateScheduled
	case "Completed":
		*state = TaskStateCompleted
	case "Inactive":
		*state = TaskStateRunning
	case "Active":
		*state = TaskStateScheduled
	case "Disabled":
		*state = TaskStateCompleted
	default:
		return fmt.Errorf("incorrect state: %s", stringValue)
	}

	return nil
}

// UnmarshalJSON used to convert string Regularity representation to Regularity type
func (regularity *Regularity) UnmarshalJSON(byteResult []byte) error {
	var stringValue string
	if err := json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}
	switch stringValue {
	case "":
		*regularity = 0
	case "RunNow":
		*regularity = RunNow
	case "OneTime":
		*regularity = OneTime
	case "Recurrent":
		*regularity = Recurrent
	default:
		return fmt.Errorf("incorrect regularity: %s", stringValue)
	}

	return nil
}
