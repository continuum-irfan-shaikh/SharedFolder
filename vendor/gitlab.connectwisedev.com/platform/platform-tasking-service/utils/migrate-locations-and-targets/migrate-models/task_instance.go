package migrateModels

import (
	"github.com/gocql/gocql"
	"time"
)

// TaskInstanceStatus type is used for status definition
type TaskInstanceStatus int

type (
	// TaskInstance represents task info for particular managed endpoint ID and particular start time.
	TaskInstance struct {
		ID        gocql.UUID         `json:"id"`
		TaskID    gocql.UUID         `json:"taskId"`
		OriginID  gocql.UUID         `json:"originId"`
		Targets   []string           `json:"targets"`
		StartedAt time.Time          `json:"startedAt"`
		Status    TaskInstanceStatus `json:"status"`
	}
)
