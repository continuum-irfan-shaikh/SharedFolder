package modelMocks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// DefaultTaskInstances is an array of predefined task instances
var (
	knownTargetID        = "target_id"
	DefaultTaskInstances = []models.TaskInstance{
		{ID: str2uuid("00000000-0000-0000-0000-000000000000"), Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{}, TaskID: str2uuid("00000000-0000-0000-0000-000000000000"), OriginID: str2uuid("00000000-0000-0000-0000-000000000000"), StartedAt: someTime},
		{ID: str2uuid("11111111-1111-1111-1111-111111111111"), Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{}, TaskID: str2uuid("11111111-1111-1111-1111-111111111111"), OriginID: str2uuid("11111111-1111-1111-1111-111111111111"), StartedAt: someTime},
		{ID: str2uuid("22222222-2222-2222-2222-222222222222"), Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{}, TaskID: str2uuid("22222222-2222-2222-2222-222222222222"), OriginID: str2uuid("22222222-2222-2222-2222-222222222222"), StartedAt: someTime},
		{ID: str2uuid("33333333-3333-3333-3333-333333333333"), Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{}, TaskID: str2uuid("33333333-3333-3333-3333-333333333333"), OriginID: str2uuid("33333333-3333-3333-3333-333333333333"), StartedAt: someTime},
		{ID: str2uuid("44444444-4444-4444-4444-444444444444"), Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{}, TaskID: str2uuid("44444444-4444-4444-4444-444444444444"), OriginID: str2uuid("44444444-4444-4444-4444-444444444444"), StartedAt: someTime},
	}
)

// NewTaskInstanceRepoMock creates a Mock for a Task Repository.
// The Mock could be empty or filled with the predefined data.
func NewTaskInstanceRepoMock(fill bool) TaskInstanceRepoMock {
	mock := TaskInstanceRepoMock{}
	mock.repo = make(map[gocql.UUID]models.TaskInstance)
	if fill {
		for _, taskInstance := range DefaultTaskInstances {
			mock.repo[taskInstance.ID] = taskInstance
		}
	}
	return mock
}

// TaskInstanceRepoMock represents a Mock for a Task Instance Repository
type TaskInstanceRepoMock struct {
	repo map[gocql.UUID]models.TaskInstance
}

// GetByTaskID returns a slice of TaskInstances found by TaskID in mocked repo
func (mock TaskInstanceRepoMock) GetByTaskID(ctx context.Context, taskID gocql.UUID) ([]models.TaskInstance, error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cassandra is down")
	}
	var results = make([]models.TaskInstance, 0)
	for _, v := range mock.repo {
		if v.TaskID == taskID {
			results = append(results, v)
		}
	}

	return results, nil
}

// GetByTaskIDs returns a slice of TaskInstances found by TaskID in mocked repo
func (mock TaskInstanceRepoMock) GetByTaskIDs(ctx context.Context, IDs ...gocql.UUID) ([]models.TaskInstance, error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cassandra is down")
	}
	var results = make([]models.TaskInstance, 0)
	for _, v := range mock.repo {
		for _, taskID := range IDs {
			if v.TaskID == taskID {
				results = append(results, v)
			}
		}
	}
	return results, nil
}

// GetInstancesCountByTaskID ...
func (mock TaskInstanceRepoMock) GetInstancesCountByTaskID(ctx context.Context, taskID gocql.UUID) (instancesCount int, err error) {
	return
}

// GetByIDs returns task instances by PartnerID from the Mocked Repository
func (mock TaskInstanceRepoMock) GetByIDs(ctx context.Context, IDs ...gocql.UUID) ([]models.TaskInstance, error) {
	fmt.Println("TaskInstanceRepoMock.GetByIDs method called, used RequestID: ", transactionID.FromContext(ctx))

	var resultInstances = make([]models.TaskInstance, 0)

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cassandra is down")
	}

	for _, id := range IDs {
		if instance, ok := mock.repo[id]; ok {
			resultInstances = append(resultInstances, instance)
		}
	}
	resultInstances = append(resultInstances, models.TaskInstance{ID: gocql.TimeUUID(), Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{gocql.TimeUUID(): statuses.TaskInstanceCanceled}})
	return resultInstances, nil
}

// GetByIDsNoTargets ..
func (mock TaskInstanceRepoMock) GetByIDsNoTargets(ctx context.Context, IDs ...gocql.UUID) ([]models.TaskInstance, error) {
	return nil, nil
}

// GetNearestInstanceAfter ..
func (mock TaskInstanceRepoMock) GetNearestInstanceAfter(taskID gocql.UUID, sinceDate time.Time) (models.TaskInstance, error) {
	return models.TaskInstance{Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{gocql.TimeUUID(): statuses.TaskInstanceCanceled}}, nil
}

// GetByIDsAndMEs returns task instances by PartnerID from the Mocked Repository
func (mock TaskInstanceRepoMock) GetByIDsAndMEs(ctx context.Context, IDs ...gocql.UUID) ([]models.TaskInstance, error) {
	fmt.Println("TaskInstanceRepoMock.GetByIDsAndMEs method called, used RequestID: ", transactionID.FromContext(ctx))

	var resultInstances = make([]models.TaskInstance, 0)

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cassandra is down")
	}

	for _, id := range IDs {
		if instance, ok := mock.repo[id]; ok {
			resultInstances = append(resultInstances, instance)
		}
	}
	return resultInstances, nil
}

// GetTopInstancesByTaskID ..
func (mock TaskInstanceRepoMock) GetTopInstancesByTaskID(ctx context.Context, taskID gocql.UUID) ([]models.TaskInstance, error) {
	return nil, nil
}

// InsertNoTargets places a Task into Mocked Repository assigning a new TaskID for it
func (mock TaskInstanceRepoMock) InsertNoTargets(ctx context.Context, taskInstance models.TaskInstance) error {
	return nil
}

// Insert places a Task into Mocked Repository assigning a new TaskID for it
func (mock TaskInstanceRepoMock) Insert(ctx context.Context, taskInstance models.TaskInstance) error {
	fmt.Println("TaskInstanceRepoMock.Insert method called, used RequestID: ", transactionID.FromContext(ctx))
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return errors.New("cassandra is down")
	}
	var emptyUUID gocql.UUID
	var err error
	if taskInstance.TaskID == emptyUUID {
		return errors.New("bad input data")
	}

	if taskInstance.ID == emptyUUID {
		taskInstance.ID, err = gocql.RandomUUID()
		if err != nil {
			logger.Log.ErrfCtx(context.TODO(), errorcode.ErrorCantProcessData, "Insert: error while getting random UUID: ", err)
		}
	}

	mock.repo[taskInstance.TaskID] = taskInstance
	return nil
}

// InsertBatch a batch of TaskInstances in repository with parameters from task instances
func (mock TaskInstanceRepoMock) InsertBatch(ctx context.Context, taskInstances []models.TaskInstance) error {
	return nil
}

// DeleteBatch ...
func (mock TaskInstanceRepoMock) DeleteBatch(ctx context.Context, taskInstances []models.TaskInstance) error {
	return nil
}

// GetByStartedAtAfter returns slice TaskInstance found by PartnerID and timestamp (started_at)
func (mock TaskInstanceRepoMock) GetByStartedAtAfter(ctx context.Context, partnerID string, from, to time.Time) ([]models.TaskInstance, error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cassandra is down")
	}

	return DefaultTaskInstances, nil
}

// Delete ...
func (mock TaskInstanceRepoMock) Delete(ctx context.Context, taskInstance models.TaskInstance) error {
	return nil
}

// GetByEndpointIDInDescOrder ...
func (mock TaskInstanceRepoMock) GetByEndpointIDInDescOrder(ctx context.Context, endpointID string, limit int) ([]models.TaskInstance, error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cassandra is down")
	}
	return DefaultTaskInstances, nil
}

// GetMinimalInstanceByID ...
func (mock TaskInstanceRepoMock) GetMinimalInstanceByID(ctx context.Context, id gocql.UUID) (models.TaskInstance, error) {
	return models.TaskInstance{}, nil
}

// UpdateStatuses ...
func (mock TaskInstanceRepoMock) UpdateStatuses(ctx context.Context, taskInstance models.TaskInstance) (err error) {
	return nil
}

// TaskInstanceCustomizableMock provides ability to mock any method of TaskInstancePersistence interface
type TaskInstanceCustomizableMock struct {
	InsertF                     func(context.Context, models.TaskInstance) error
	InsertBatchF                func(context.Context, []models.TaskInstance) error
	GetInstancesCountByTaskIDF  func(ctx context.Context, taskID gocql.UUID) (instancesCount int, err error)
	DeleteBatchF                func(context.Context, []models.TaskInstance) error
	DeleteF                     func(context.Context, models.TaskInstance) error
	GetByIDsF                   func(context.Context, ...gocql.UUID) ([]models.TaskInstance, error)
	GetByIDsAndMEsF             func(context.Context, ...gocql.UUID) ([]models.TaskInstance, error)
	GetByTaskIDAndEndpointIDF   func(context.Context, gocql.UUID, string, int) ([]models.TaskInstance, error)
	GetByEndpointIDInDescOrderF func(ctx context.Context, endpointID string, limit int) ([]models.TaskInstance, error)
	GetByTaskIDF                func(context.Context, gocql.UUID) ([]models.TaskInstance, error)
	GetByStartedAtAfterF        func(ctx context.Context, partnerID string, from, to time.Time) ([]models.TaskInstance, error)
	GetTopInstanceByTaskIDF     func(ctx context.Context, taskID gocql.UUID) ([]models.TaskInstance, error)
	GetNearestInstanceAfterF    func(taskID gocql.UUID, sinceDate time.Time) (models.TaskInstance, error)
	GetByTaskIDsF               func(ctx context.Context, IDs []gocql.UUID) ([]models.TaskInstance, error)
	GetMinimalInstanceByIDF     func(ctx context.Context, id gocql.UUID) (models.TaskInstance, error)
	UpdateStatusesF             func(ctx context.Context, taskInstance models.TaskInstance) (err error)
}

// Insert ...
func (ti TaskInstanceCustomizableMock) Insert(ctx context.Context, taskInstance models.TaskInstance) error {
	return ti.InsertF(ctx, taskInstance)
}

// InsertNoTargets ...
func (ti TaskInstanceCustomizableMock) InsertNoTargets(ctx context.Context, taskInstance models.TaskInstance) error {
	return ti.InsertF(ctx, taskInstance)
}

// GetByIDsNoTargets ..
func (ti TaskInstanceCustomizableMock) GetByIDsNoTargets(ctx context.Context, IDs ...gocql.UUID) ([]models.TaskInstance, error) {
	return nil, nil
}

// GetInstancesCountByTaskID ...
func (ti TaskInstanceCustomizableMock) GetInstancesCountByTaskID(ctx context.Context, taskID gocql.UUID) (instancesCount int, err error) {
	return ti.GetInstancesCountByTaskIDF(ctx, taskID)
}

// GetNearestInstanceAfter  ..
func (ti TaskInstanceCustomizableMock) GetNearestInstanceAfter(taskID gocql.UUID, sinceDate time.Time) (models.TaskInstance, error) {
	return ti.GetNearestInstanceAfterF(taskID, sinceDate)
}

// InsertBatch ...
func (ti TaskInstanceCustomizableMock) InsertBatch(ctx context.Context, taskInstances []models.TaskInstance) error {
	return ti.InsertBatchF(ctx, taskInstances)
}

// GetByIDs ...
func (ti TaskInstanceCustomizableMock) GetByIDs(ctx context.Context, ids ...gocql.UUID) ([]models.TaskInstance, error) {
	return ti.GetByIDsF(ctx, ids...)
}

// GetTopInstancesByTaskID ..
func (ti TaskInstanceCustomizableMock) GetTopInstancesByTaskID(ctx context.Context, taskID gocql.UUID) ([]models.TaskInstance, error) {
	return ti.GetTopInstanceByTaskIDF(ctx, taskID)
}

// GetByIDsAndMEs ...
func (ti TaskInstanceCustomizableMock) GetByIDsAndMEs(ctx context.Context, ids ...gocql.UUID) ([]models.TaskInstance, error) {
	return ti.GetByIDsAndMEsF(ctx, ids...)
}

// GetByTaskID ...
func (ti TaskInstanceCustomizableMock) GetByTaskID(ctx context.Context, taskIDs gocql.UUID) ([]models.TaskInstance, error) {
	return ti.GetByTaskIDF(ctx, taskIDs)
}

// GetByTaskIDs ...
func (ti TaskInstanceCustomizableMock) GetByTaskIDs(ctx context.Context, IDs ...gocql.UUID) ([]models.TaskInstance, error) {
	return ti.GetByTaskIDsF(ctx, IDs)
}

// GetByStartedAtAfter ...
func (ti TaskInstanceCustomizableMock) GetByStartedAtAfter(ctx context.Context, partnerID string, from, to time.Time) ([]models.TaskInstance, error) {
	return ti.GetByStartedAtAfterF(ctx, partnerID, from, to)
}

// DeleteBatch ...
func (ti TaskInstanceCustomizableMock) DeleteBatch(ctx context.Context, taskInstances []models.TaskInstance) error {
	return ti.DeleteBatchF(ctx, taskInstances)
}

// Delete ...
func (ti TaskInstanceCustomizableMock) Delete(ctx context.Context, taskInstance models.TaskInstance) error {
	return ti.DeleteF(ctx, taskInstance)
}

// GetMinimalInstanceByID ...
func (ti TaskInstanceCustomizableMock) GetMinimalInstanceByID(ctx context.Context, id gocql.UUID) (models.TaskInstance, error) {
	return ti.GetMinimalInstanceByIDF(ctx, id)
}

// UpdateStatuses ...
func (ti TaskInstanceCustomizableMock) UpdateStatuses(ctx context.Context, taskInstance models.TaskInstance) (err error) {
	return ti.UpdateStatusesF(ctx, taskInstance)
}
