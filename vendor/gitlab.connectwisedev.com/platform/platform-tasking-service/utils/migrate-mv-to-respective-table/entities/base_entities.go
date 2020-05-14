package entities

import (
	"time"

	"github.com/gocql/gocql"
)

type TaskInstance struct {
	PartnerID     string
	ID            gocql.UUID
	TaskID        gocql.UUID
	OriginID      gocql.UUID
	Targets       []string
	StartedAt     time.Time
	LastRunTime   time.Time
	Statuses      map[gocql.UUID]TaskInstanceStatus
	OverallStatus TaskInstanceStatus
	FailureCount  int
	SuccessCount  int
	TriggeredBy   string
	Name          string
	Status        int
}

type TaskInstanceStatus int

type ExecutionResult struct {
	ManagedEndpointID gocql.UUID
	TaskInstanceID    gocql.UUID
	UpdatedAt         time.Time
	ExecutionStatus   TaskInstanceStatus
	StdErr            string
	StdOut            string
}
