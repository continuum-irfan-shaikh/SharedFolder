package scheduler

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

//go:generate mockgen -destination=./scheduler_repo_mock_test.go  -package=scheduler -source=./usecase.go

//SchedulerRepo - represents interface for scheduler repository
type SchedulerRepo interface {
	GetLastUpdate() (time.Time, error)
	UpdateScheduler(time.Time) error
}

//TaskInstanceRepo - represents task instance repository
type TaskInstanceRepo interface {
	GetInstance(id gocql.UUID) (m.TaskInstance, error)
	Insert(ti m.TaskInstance, ttl int) error
	GetNearestInstanceAfter(taskID gocql.UUID, sinceDate time.Time) (m.TaskInstance, error)
	AppendNewEndpoints(ti m.TaskInstance, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) (err error)
	RemoveInactiveEndpoints(ti m.TaskInstance, endpoints ...gocql.UUID) (err error)
}

//TaskRepo - represents task repository
type TaskRepo interface {
	GetByRunTimeRange(ctx context.Context, timeRange []time.Time) ([]m.Task, error)
	UpdateSchedulerFields(ctx context.Context, tasks ...m.Task) error
	InsertOrUpdate(ctx context.Context, tasks ...m.Task) error
	GetTargetTypeByEndpoint(partnerID string, taskID, endpointID gocql.UUID, external bool) (m.TargetType, error)
}

//TaskExecutionRepo - represents task execution repository
type TaskExecutionRepo interface {
	ExecuteTasks(ctx context.Context, payload apiModels.ExecutionPayload, partnerID, taskType string) error
}

//ExecutionResultRepo - represents execution result repository
type ExecutionResultRepo interface {
	Publish(msg tasking.ExecutionResultKafkaMessage) error
}

// CacheRepo ..
type CacheRepo interface {
	GetByOriginID(ctx context.Context, partnerID string, templateID gocql.UUID, isNoc bool) (m.TemplateDetails, error)
	CalculateExpectedExecutionTimeSec(ctx context.Context, task m.Task) int
}

// ExecutionExpirationRepo ..
type ExecutionExpirationRepo interface {
	Insert(ee entities.ExecutionExpiration, ttl int) error
}

// TargetsRepo ..
type TargetsRepo interface {
	GetTargetsByTaskID(partnerID string, taskID gocql.UUID) (ids []gocql.UUID, err error)
}

// DynamicGroupRepo ..
type DynamicGroupRepo interface {
	GetEndpointsByGroupIDs(ctx context.Context, targetIDs []string, createdBy, partnerID string, hasNOCAccess bool) (ids []gocql.UUID, err error)
}

// SiteRepo ..
type SiteRepo interface {
	GetEndpointsBySiteIDs(ctx context.Context, partnerID string, siteIDs []string) (ids []gocql.UUID, err error)
}

// AssetsClient ..
type AssetsClient interface {
	GetLocationByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (location *time.Location, err error)
}

//EncryptionService - provide us interface for encryption of credentials
type EncryptionService interface {
	Encrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
	Decrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
}

//AgentEncryptionService - provide interface for encryption credentials by key stored for particular endpoint
type AgentEncryptionService interface {
	Encrypt(ctx context.Context, endpointID gocql.UUID, credentials agentModels.Credentials) (encrypted agentModels.Credentials, err error)
}
