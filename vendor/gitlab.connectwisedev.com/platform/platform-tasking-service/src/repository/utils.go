package repository

import (
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

// DatabaseRepositories is a DTO object that contains repositories to work with DB
type DatabaseRepositories struct {
	ExecutionHistory    TaskExecutionHistoryRepo
	ExecutionExpiration ExecutionExpirationRepo
	ExecutionResults    ExecutionResultsRepo
	TaskInstance        InstancesRepo
	Scheduler           SchedulerRepo
	LegacyMigration     LegacyMigration
	Triggers            TriggersRepo
	Targets             TargetsRepo
	Profiles            Profiles
	Task                Task
}

// ExecutionExpirationRepo ..
type ExecutionExpirationRepo interface {
	Insert(ee entities.ExecutionExpiration, ttl int) error
}

// TargetsRepo represents interface to work with targets table
type TargetsRepo interface {
	GetTargetsByTaskID(partnerID string, taskID gocql.UUID) (ids []gocql.UUID, err error)
	Insert(partnerID string, taskID gocql.UUID, targets models.Target) error
}

// InstancesRepo - task instances repository interface
type InstancesRepo interface {
	GetInstancesForScheduled(IDs []string) ([]entities.TaskInstance, error)
	GetTopInstancesForScheduledByTaskIDs(taskIDs []string) ([]entities.TaskInstance, error)
	GetMinimalInstanceByID(id string) (entities.TaskInstance, error)
	GetByStartedAtAfter(partnerID string, from, to time.Time) ([]entities.TaskInstance, error)
}

// TaskExecutionHistoryRepo - interface to perform actions with task execution history database
type TaskExecutionHistoryRepo interface {
	Insert(entities.TaskExecHistory) error
}

// Task - tasks repository interface
type Task interface {
	GetName(partnerID string, id string) (string, error)
	GetNext(partnerID string) ([]entities.Task, error)
	GetScheduledTasks(partnerID string) ([]entities.ScheduledTasks, error)
}

// SchedulerRepo - represents interface for scheduler repository
type SchedulerRepo interface {
	GetLastUpdate() (time.Time, error)
	UpdateScheduler(time.Time) error
}

//ExecutionResultsRepo - script execution results repository interface
type ExecutionResultsRepo interface {
	GetLastResultByEndpointID(endpointID string) (res entities.ExecutionResult, err error)
	GetLastExecutions(partnerID string, endpointIDs map[string]struct{}) (res []entities.LastExecution, err error)
}
