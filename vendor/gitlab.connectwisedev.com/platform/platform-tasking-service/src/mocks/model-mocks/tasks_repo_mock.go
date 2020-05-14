package modelMocks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

var (
	someState    = statuses.TaskStateActive
	originID, _  = gocql.RandomUUID()
	someTime     = time.Now().UTC()
	createdByStr = "Admin"
	// PartnerID is allowed access partnerID
	PartnerID = "0"
)

//key type is unexported to prevent collisions with context keys defined in
//other packages.
type key string

//IsNeedError key for Context which shows whether RepoMock should return error or not.
const IsNeedError key = "isNeedError"

// DefaultTasks is an array of predefined tasks
var DefaultTasks = []models.Task{
	{ID: ExistedTaskID, Name: "name1", Description: "description1", CreatedAt: someTime, CreatedBy: createdByStr, PartnerID: PartnerID, OriginID: originID, State: someState, RunTimeUTC: someTime, Type: models.TaskTypeScript, ManagedEndpointID: ExistedManagedEndpointID, Schedule: apiModels.Schedule{Regularity: apiModels.RunNow}},
	{ID: str2uuid("22222222-2222-2222-2222-222222222222"), Name: "name2", Description: "description2", CreatedAt: someTime, CreatedBy: createdByStr, PartnerID: PartnerID, OriginID: originID, State: someState, RunTimeUTC: someTime, Type: models.TaskTypeScript, ManagedEndpointID: ExistedManagedEndpointID, Schedule: apiModels.Schedule{Regularity: apiModels.RunNow}},
	{ID: str2uuid("00000000-0000-0000-0000-000000000000"), Name: "name0", Description: "description0", CreatedAt: someTime, CreatedBy: createdByStr, PartnerID: PartnerID, OriginID: originID, State: someState, RunTimeUTC: someTime, Type: models.TaskTypeScript, ManagedEndpointID: ExistedManagedEndpointID, Schedule: apiModels.Schedule{Regularity: apiModels.RunNow}},
	{ID: str2uuid("33333333-3333-3333-3333-333333333333"), Name: "name3", Description: "description3", CreatedAt: someTime, CreatedBy: createdByStr, PartnerID: PartnerID, OriginID: originID, State: someState, RunTimeUTC: someTime, Type: models.TaskTypeScript, ManagedEndpointID: ExistedManagedEndpointID, Schedule: apiModels.Schedule{Regularity: apiModels.RunNow}},
	{ID: str2uuid("44444444-4444-4444-4444-444444444444"), Name: "name4", Description: "description4", CreatedAt: someTime, CreatedBy: createdByStr, PartnerID: PartnerID, OriginID: originID, State: someState, RunTimeUTC: someTime.Add(time.Minute), Type: models.TaskTypeScript, ManagedEndpointID: ExistedManagedEndpointID, Schedule: apiModels.Schedule{Regularity: apiModels.OneTime}},
	{ID: str2uuid("55555555-5555-5555-5555-555555555555"), Name: "name5", Description: "description5", CreatedAt: someTime, CreatedBy: createdByStr, PartnerID: PartnerID, OriginID: originID, State: someState, RunTimeUTC: someTime.Add(time.Minute), Type: models.TaskTypeScript, ManagedEndpointID: ExistedManagedEndpointID, Schedule: apiModels.Schedule{Regularity: apiModels.Recurrent}},
	{ID: str2uuid("66666666-5555-5555-5555-555555555555"), Name: "name7", Description: "description7", CreatedAt: someTime, CreatedBy: createdByStr, PartnerID: PartnerID, OriginID: originID, State: someState, RunTimeUTC: someTime, Type: models.TaskTypeScript, ManagedEndpointID: TargetID, Schedule: apiModels.Schedule{Regularity: apiModels.RunNow}},
}

func str2uuid(stringUUID string) gocql.UUID {
	uuid, err := gocql.ParseUUID(stringUUID)
	if err != nil {
		return uuid
	}
	return uuid
}

// NewTaskRepoMock creates a Mock for a Task Repository.
// The Mock could be empty or filled with the predefined data.
func NewTaskRepoMock(fill bool) TaskRepoMock {
	mock := TaskRepoMock{}
	mock.Repo = make(map[gocql.UUID]models.Task)
	if fill {
		for _, task := range DefaultTasks {
			mock.Repo[task.ID] = task
		}
	}
	return mock
}

// TaskRepoMock represents a Mock for a Task Repository.
type TaskRepoMock struct {
	Repo map[gocql.UUID]models.Task
}

// GetByPartner returns an array of tasks by PartnerID from the mocked Repository
func (mock TaskRepoMock) GetByPartner(ctx context.Context, partnerID string) ([]models.Task, error) {
	fmt.Println("TaskRepoMock.GetByPartner method called, used RequestID: ", transactionID.FromContext(ctx))
	var resultTasks = make([]models.Task, 0)

	if IsNeedError, _ := ctx.Value(IsNeedError).(bool); IsNeedError {
		return resultTasks, errors.New("cassandra is down")
	}

	for _, task := range mock.Repo {
		if task.PartnerID == partnerID {
			resultTasks = append(resultTasks, task)
		}
	}

	return resultTasks, nil
}

// GetTargetTypeByEndpoint ..
func (mock TaskRepoMock) GetTargetTypeByEndpoint(partnerID string, taskID, endpointID gocql.UUID, external bool) (models.TargetType, error) {
	return models.DynamicGroup, nil
}

// GetByIDs returns an array of Tasks specified by its TaskIDs and Partner form the Mocked Repository.
func (mock TaskRepoMock) GetByIDs(ctx context.Context, mCache persistency.Cache, partnerID string, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]models.Task, error) {
	fmt.Println("TaskRepoMock.GetByID method called, used RequestID: ", transactionID.FromContext(ctx))
	fmt.Printf("(GetByID) task id: %v", taskIDs)

	var resultTasks = make([]models.Task, 0)

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return resultTasks, errors.New("cassandra is down")
	}

	for _, taskID := range taskIDs {
		if task, ok := mock.Repo[taskID]; ok && task.PartnerID == partnerID {
			resultTasks = append(resultTasks, task)
		}
	}

	return resultTasks, nil
}

// GetCountForAllTargets returns a map with task count of each target for specific partner
func (mock TaskRepoMock) GetCountForAllTargets(ctx context.Context, partnerID string) ([]models.TaskCount, error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("cassandra is down")
	}
	taskCountMap := make(map[gocql.UUID]int)
	for _, task := range mock.Repo {
		taskCountMap[task.ManagedEndpointID]++
	}
	return models.ConvertMapToTaskCountArray(taskCountMap), nil
}

// GetCountByManagedEndpointID returns a number of Tasks by target_id form the Mocked Repository
func (mock TaskRepoMock) GetCountByManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID) (taskCount models.TaskCount, err error) {
	fmt.Println("TaskRepoMock.GetCountByTargetID method called, used RequestID: ", transactionID.FromContext(ctx))
	taskCount.ManagedEndpointID = managedEndpointID

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return taskCount, errors.New("cassandra is down")
	}

	for _, task := range mock.Repo {
		if task.ManagedEndpointID == managedEndpointID && task.PartnerID == partnerID {
			taskCount.Count++
		}
	}
	return
}

// GetCountsByPartner returns a number of Tasks by Partner form the Mocked Repository
func (mock TaskRepoMock) GetCountsByPartner(ctx context.Context, partnerID string) (taskCounts []models.TaskCount, err error) {
	fmt.Println("TaskRepoMock.GetCountsByPartner method called, used RequestID: ", transactionID.FromContext(ctx))

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return taskCounts, errors.New("cassandra is down")
	}

	taskCountMap := make(map[gocql.UUID]int)
	for _, task := range mock.Repo {
		if task.PartnerID == partnerID && !task.ExternalTask {
			taskCountMap[task.ManagedEndpointID]++
		}
	}

	return models.ConvertMapToTaskCountArray(taskCountMap), err
}

// GetByPartnerAndManagedEndpointID returns Tasks by target_id and partner form the Mocked Repository
func (mock TaskRepoMock) GetByPartnerAndManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID, count int) ([]models.Task, error) {
	fmt.Println("TaskRepoMock.GetByPartnerAndManagedEndpointID method called, used RequestID: ", transactionID.FromContext(ctx))
	fmt.Printf("(GetByPartnerAndManagedEndpointID) task target_id: %v", managedEndpointID)

	var resultTasks []models.Task

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return resultTasks, errors.New("cassandra is down")
	}

	for _, task := range mock.Repo {
		if task.ManagedEndpointID == managedEndpointID && task.PartnerID == partnerID {
			resultTasks = append(resultTasks, task)
		}
	}
	return resultTasks, nil
}

// GetByRunTimeRange returns a Task specified by next Run Time.
func (mock TaskRepoMock) GetByRunTimeRange(ctx context.Context, runTimes []time.Time) ([]models.Task, error) {
	fmt.Println("TaskRepoMock.GetByRunTime method called, used RequestID: ", transactionID.FromContext(ctx))
	var resultTasks []models.Task

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return resultTasks, errors.New("cassandra is down")
	}

	runTimesStrSlice := make([]interface{}, len(runTimes))
	for i, v := range runTimes {
		runTimesStrSlice[i] = v.Format(common.CassandraTimeFormat)
	}
	for _, task := range mock.Repo {

		if contains(runTimesStrSlice, task.RunTimeUTC.Format(common.CassandraTimeFormat)) && task.Schedule.Regularity != apiModels.RunNow {
			resultTasks = append(resultTasks, task)
		}
	}
	return resultTasks, nil
}

// GetByPartnerAndTime ..
func (mock TaskRepoMock) GetByPartnerAndTime(ctx context.Context, partnerID string, time time.Time) ([]models.Task, error) {
	return []models.Task{}, nil
}

// GetRecentTasksByIDs ..
func (mock TaskRepoMock) GetRecentTasksByIDs(ctx context.Context, partnerID string, hasNocAccess bool, taskIDs ...gocql.UUID) (map[gocql.UUID]models.Task, error) {
	return nil, nil
}

// GetExecutionResultTaskData ..
func (mock TaskRepoMock) GetExecutionResultTaskData(partnerID string, taskID, endpointID gocql.UUID) (models.ExecutionResultTaskData, error) {
	return models.ExecutionResultTaskData{}, nil
}

// GetManagedEndpointIDsOfActiveTasks ..
func (mock TaskRepoMock) GetManagedEndpointIDsOfActiveTasks(ctx context.Context, partnerID string, taskID gocql.UUID) (map[gocql.UUID]struct{}, error) {
	return nil, nil
}

// GetByIDAndManagedEndpoints returns Task founded by taskID and Targets for specific partner
func (mock TaskRepoMock) GetByIDAndManagedEndpoints(ctx context.Context, partnerID string, taskID gocql.UUID, managedEndpointIDs ...gocql.UUID) ([]models.Task, error) {
	var resultTasks []models.Task

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return resultTasks, errors.New("Cassandra is down")
	}

	targetsInterfaces := make([]interface{}, len(managedEndpointIDs))
	for i, v := range managedEndpointIDs {
		targetsInterfaces[i] = v
	}

	for _, task := range mock.Repo {
		if task.ID == taskID && task.PartnerID == partnerID && contains(targetsInterfaces, task.ManagedEndpointID) {
			resultTasks = append(resultTasks, task)
		}
	}
	return resultTasks, nil
}

// UpdateTask updated a Task with id and partner_id into Mocked Repository.
func (mock TaskRepoMock) UpdateTask(ctx context.Context, inputStruct interface{}, partnerID string, taskID gocql.UUID) (err error) {
	fmt.Println("TaskRepoMock.UpdateTask method called, used RequestID: ", transactionID.FromContext(ctx))

	tasks, err := mock.GetByIDs(ctx, nil, partnerID, false, taskID)
	if err != nil {
		return
	}

	if len(tasks) == 0 {
		return models.CantUpdateTaskError{}
	}
	//Update task here based on inputStruct

	err = mock.InsertOrUpdate(ctx, tasks[0])
	return

}

// InsertOrUpdate places a Task into Mocked Repository assigning a new TaskID for it.
func (mock TaskRepoMock) InsertOrUpdate(ctx context.Context, tasks ...models.Task) error {
	fmt.Println("TaskRepoMock.InsertOrUpdate method called, used RequestID: ", transactionID.FromContext(ctx))

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return errors.New("Cassandra is down")
	}

	for _, task := range tasks {
		if task.OriginID == emptyUUID {
			return errors.New("bad input data")
		}

		mock.Repo[task.ID] = task
		fmt.Printf("(Insert) task id: %v", task.ID)
	}
	return nil
}

// UpdateSchedulerFields ..
func (mock TaskRepoMock) UpdateSchedulerFields(ctx context.Context, internalTasks ...models.Task) error {
	return nil
}

// UpdateModifiedFieldsByMEs ..
func (mock TaskRepoMock) UpdateModifiedFieldsByMEs(ctx context.Context, task models.Task, managedEndpoints ...gocql.UUID) error {
	return nil
}

// GetByLastTaskInstanceIDs gets Tasks by Last TaskInstances
func (mock TaskRepoMock) GetByLastTaskInstanceIDs(
	ctx context.Context,
	partnerID string,
	endpointID gocql.UUID,
	lastTaskInstanceIDs ...gocql.UUID,
) (map[gocql.UUID]models.Task, error) {

	var resultTasks = make(map[gocql.UUID]models.Task)
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return resultTasks, errors.New("Cassandra is down")
	}

	for _, task := range DefaultTasks {
		resultTasks[task.ID] = task
	}

	return resultTasks, nil
}

func contains(slice []interface{}, value interface{}) bool {
	for _, sliceValue := range slice {
		fmt.Println(value, sliceValue)
		if sliceValue == value {
			return true
		}
	}
	return false
}

// GetTasksFilteredWithTime ..
func (mock TaskRepoMock) GetTasksFilteredWithTime(ctx context.Context, partnerID string, taskID gocql.UUID, excludedEndpoints map[gocql.UUID]models.Task, after time.Time) (map[gocql.UUID][]models.Task, error) {
	return nil, nil
}

// Delete ...
func (mock TaskRepoMock) Delete(ctx context.Context, tasks []models.Task) error {
	return nil
}

// TaskCustomizableMock ...
type TaskCustomizableMock struct {
	InsertOrUpdateF            func(ctx context.Context, tasks ...models.Task) error
	DeleteF                    func(ctx context.Context, tasks []models.Task) error
	UpdateTaskF                func(ctx context.Context, inputStruct interface{}, partnerID string, taskID gocql.UUID) error
	UpdateSchedulerFieldsF     func(ctx context.Context, internalTasks ...models.Task) error
	UpdateModifiedFieldsByMEsF func(ctx context.Context, task models.Task, managedEndpoints ...gocql.UUID) error

	GetByPartnerF                       func(ctx context.Context, partnerID string) ([]models.Task, error)
	GetByIDsF                           func(ctx context.Context, cache persistency.Cache, partnerID string, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]models.Task, error)
	GetByIDsAndMEsF                     func(ctx context.Context, cache persistency.Cache, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]models.Task, error)
	GetByLastTaskInstanceIDsF           func(ctx context.Context, partnerID string, endpointID gocql.UUID, lastTaskInstanceIDs ...gocql.UUID) (tasksByIDMap map[gocql.UUID]models.Task, err error)
	GetByIDAndManagedEndpointsF         func(ctx context.Context, partnerID string, taskID gocql.UUID, managedEndpointIDs ...gocql.UUID) ([]models.Task, error)
	GetByRunTimeRangeF                  func(ctx context.Context, runTimeRange []time.Time) ([]models.Task, error)
	GetByPartnerAndManagedEndpointIDF   func(ctx context.Context, partnerID string, managedEndpointID gocql.UUID, count int) ([]models.Task, error)
	GetCountByManagedEndpointIDF        func(ctx context.Context, partnerID string, managedEndpointID gocql.UUID) (models.TaskCount, error)
	GetRecentTasksByIDsF                func(ctx context.Context, partnerID string, hasNocAccess bool, taskIDs ...gocql.UUID) (map[gocql.UUID]models.Task, error)
	GetCountsByPartnerF                 func(ctx context.Context, partnerID string) ([]models.TaskCount, error)
	GetTasksFilteredWithTimeF           func(ctx context.Context, partnerID string, taskID gocql.UUID, excludedEndpoints map[gocql.UUID]models.Task, after time.Time) (map[gocql.UUID][]models.Task, error)
	GetByPartnerAndTimeF                func(ctx context.Context, partnerID string, time time.Time) ([]models.Task, error)
	GetManagedEndpointIDsOfActiveTasksF func(ctx context.Context, partnerID string, taskID gocql.UUID) (map[gocql.UUID]struct{}, error)
	GetExecutionResultTaskDataF         func(partnerID string, taskID, endpointID gocql.UUID) (models.ExecutionResultTaskData, error)
	GetTargetTypeByEndpointF            func(partnerID string, taskID, endpointID gocql.UUID, external bool) (models.TargetType, error)
}

// GetExecutionResultTaskData ..
func (t TaskCustomizableMock) GetExecutionResultTaskData(partnerID string, taskID, endpointID gocql.UUID) (models.ExecutionResultTaskData, error) {
	return t.GetExecutionResultTaskDataF(partnerID, taskID, endpointID)
}

// GetTargetTypeByEndpoint ..
func (t TaskCustomizableMock) GetTargetTypeByEndpoint(partnerID string, taskID, endpointID gocql.UUID, external bool) (models.TargetType, error) {
	return 0, nil
}

// GetManagedEndpointIDsOfActiveTasks ..
func (t TaskCustomizableMock) GetManagedEndpointIDsOfActiveTasks(ctx context.Context, partnerID string, taskID gocql.UUID) (map[gocql.UUID]struct{}, error) {
	return t.GetManagedEndpointIDsOfActiveTasksF(ctx, partnerID, taskID)
}

// InsertOrUpdate ...
func (t TaskCustomizableMock) InsertOrUpdate(ctx context.Context, tasks ...models.Task) error {
	return t.InsertOrUpdateF(ctx, tasks...)
}

// Delete ...
func (t TaskCustomizableMock) Delete(ctx context.Context, tasks []models.Task) error {
	return t.DeleteF(ctx, tasks)
}

// GetByPartnerAndTime ...
func (t TaskCustomizableMock) GetByPartnerAndTime(ctx context.Context, partnerID string, time time.Time) ([]models.Task, error) {
	return t.GetByPartnerAndTimeF(ctx, partnerID, time)
}

// UpdateTask ...
func (t TaskCustomizableMock) UpdateTask(ctx context.Context, inputStruct interface{}, partnerID string, taskID gocql.UUID) error {
	return t.UpdateTaskF(ctx, inputStruct, partnerID, taskID)
}

// UpdateModifiedFieldsByMEs ..
func (t TaskCustomizableMock) UpdateModifiedFieldsByMEs(ctx context.Context, task models.Task, managedEndpoints ...gocql.UUID) error {
	return t.UpdateModifiedFieldsByMEsF(ctx, task, managedEndpoints...)
}

// UpdateRunTimeAndSchedule ...
func (t TaskCustomizableMock) UpdateSchedulerFields(ctx context.Context, internalTasks ...models.Task) error {
	return t.UpdateSchedulerFieldsF(ctx, internalTasks...)
}

// GetByPartner ...
func (t TaskCustomizableMock) GetByPartner(ctx context.Context, partnerID string) ([]models.Task, error) {
	return t.GetByPartnerF(ctx, partnerID)
}

// GetByIDs ...
func (t TaskCustomizableMock) GetByIDs(ctx context.Context, cache persistency.Cache, partnerID string, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]models.Task, error) {
	return t.GetByIDsF(ctx, cache, partnerID, isCommonFieldsNeededOnly, taskIDs...)
}

// GetByIDsAndMEs ...
func (t TaskCustomizableMock) GetByIDsAndMEs(ctx context.Context, cache persistency.Cache, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]models.Task, error) {
	return t.GetByIDsAndMEsF(ctx, cache, isCommonFieldsNeededOnly, taskIDs...)
}

// GetByLastTaskInstanceIDs ...
func (t TaskCustomizableMock) GetByLastTaskInstanceIDs(ctx context.Context, partnerID string, endpointID gocql.UUID, lastTaskInstanceIDs ...gocql.UUID) (tasksByIDMap map[gocql.UUID]models.Task, err error) {
	return t.GetByLastTaskInstanceIDsF(ctx, partnerID, endpointID, lastTaskInstanceIDs...)
}

// GetRecentTasksByIDs ..
func (t TaskCustomizableMock) GetRecentTasksByIDs(ctx context.Context, partnerID string, hasNocAccess bool, taskIDs ...gocql.UUID) (map[gocql.UUID]models.Task, error) {
	return t.GetRecentTasksByIDsF(ctx, partnerID, hasNocAccess, taskIDs...)
}

// GetByIDAndManagedEndpoints ..
func (t TaskCustomizableMock) GetByIDAndManagedEndpoints(ctx context.Context, partnerID string, taskID gocql.UUID, managedEndpointIDs ...gocql.UUID) ([]models.Task, error) {
	return t.GetByIDAndManagedEndpointsF(ctx, partnerID, taskID, managedEndpointIDs...)
}

// GetByRunTimeRange ...
func (t TaskCustomizableMock) GetByRunTimeRange(ctx context.Context, runTimeRange []time.Time) ([]models.Task, error) {
	return t.GetByRunTimeRangeF(ctx, runTimeRange)
}

// GetByPartnerAndManagedEndpointID ...
func (t TaskCustomizableMock) GetByPartnerAndManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID, count int) ([]models.Task, error) {
	return t.GetByPartnerAndManagedEndpointIDF(ctx, partnerID, managedEndpointID, count)
}

// GetCountByManagedEndpointID ...
func (t TaskCustomizableMock) GetCountByManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID) (models.TaskCount, error) {
	return t.GetCountByManagedEndpointIDF(ctx, partnerID, managedEndpointID)
}

// GetCountsByPartner ...
func (t TaskCustomizableMock) GetCountsByPartner(ctx context.Context, partnerID string) ([]models.TaskCount, error) {
	return t.GetCountsByPartnerF(ctx, partnerID)
}

// GetTasksFilteredWithTime ...
func (t TaskCustomizableMock) GetTasksFilteredWithTime(ctx context.Context, partnerID string, taskID gocql.UUID, excludedEndpoints map[gocql.UUID]models.Task, after time.Time) (map[gocql.UUID][]models.Task, error) {
	return t.GetTasksFilteredWithTimeF(ctx, partnerID, taskID, excludedEndpoints, after)
}

// TaskDefaultMock ...
type TaskDefaultMock struct{}

// InsertOrUpdate ...
func (TaskDefaultMock) InsertOrUpdate(ctx context.Context, tasks ...models.Task) error {
	return nil
}

// UpdateTask ...
func (TaskDefaultMock) UpdateTask(ctx context.Context, inputStruct interface{}, partnerID string, taskID gocql.UUID) error {
	return nil
}

// GetByPartner ...
func (TaskDefaultMock) GetByPartner(ctx context.Context, partnerID string) ([]models.Task, error) {
	return nil, nil
}

// GetByIDs ...
func (TaskDefaultMock) GetByIDs(ctx context.Context, cache persistency.Cache, partnerID string, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]models.Task, error) {
	return nil, nil
}

// GetByLastTaskInstanceIDs ...
func (TaskDefaultMock) GetByLastTaskInstanceIDs(ctx context.Context, partnerID string, endpointID gocql.UUID, lastTaskInstanceIDs ...gocql.UUID) (tasksByIDMap map[gocql.UUID]models.Task, err error) {
	return nil, nil
}

// GetByIDAndManagedEndpoints ...
func (TaskDefaultMock) GetByIDAndManagedEndpoints(ctx context.Context, partnerID string, taskID gocql.UUID, managedEndpointIDs ...gocql.UUID) ([]models.Task, error) {
	return nil, nil
}

// GetByRunTimeRange ...
func (TaskDefaultMock) GetByRunTimeRange(ctx context.Context, runTimeRange []time.Time) ([]models.Task, error) {
	return nil, nil
}

// GetByPartnerAndManagedEndpointID ...
func (TaskDefaultMock) GetByPartnerAndManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID, count int) ([]models.Task, error) {
	return nil, nil
}

// GetCountByManagedEndpointID ...
func (TaskDefaultMock) GetCountByManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID) (models.TaskCount, error) {
	return models.TaskCount{}, nil
}

// GetCountsByPartner ...
func (TaskDefaultMock) GetCountsByPartner(ctx context.Context, partnerID string) ([]models.TaskCount, error) {
	return nil, nil
}

// GetTasksFilteredWithTime ...
func (TaskDefaultMock) GetTasksFilteredWithTime(ctx context.Context, partnerID string, taskID gocql.UUID, excludedEndpoints map[gocql.UUID]struct{}, after time.Time) (map[gocql.UUID][]models.Task, error) {
	return nil, nil
}

// GetManagedEndpointIDsOfActiveTasks ..
func (TaskDefaultMock) GetManagedEndpointIDsOfActiveTasks(ctx context.Context, partnerID string, taskID gocql.UUID) (map[gocql.UUID]struct{}, error) {
	return nil, nil
}
