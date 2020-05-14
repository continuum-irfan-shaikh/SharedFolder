package models

//go:generate mockgen -destination=../mocks/mocks-gomock/taskPersistance_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models TaskPersistence

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

const (
	baseTable                               = "tasks"
	taskByRuntimeUnixTableName              = "task_by_runtime_unix_mv"
	tasksOrderByLastTaskInstanceIDTableName = "tasks_order_by_last_task_instance_id_mv"
	tasksByRuntimeTableName                 = "tasks_by_runtime_mv"
	tasksByIDTableName                      = "tasks_by_id_mv"

	selectTaskBeforeUpdate = `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                                                          run_time_unix, run_time, state, target_type FROM tasks
                                                      WHERE partner_id = ? AND id = ? AND managed_endpoint_id = ? and external_task = ?`

	escape          = "\""
	anyTimesReplace = -1
)

// TaskTypeScript is the constant for Task with type `script`
const TaskTypeScript = "script"

// TaskTypeAction is constant used for UI representation of continuum defined scripts
const TaskTypeAction = "action"

const (
	// LogoffTriggerType is a logoff trigger type constant that stored in task.Trigger
	LogoffTriggerType = "logoff"
)

// DefaultEndpointUID aka EndpointID for task with empty Site/DG
const DefaultEndpointUID = "10000000-0000-0000-0000-000000000001"

type (
	// Task contains fields which describe task
	// Fields explanation:
	// - ExternalTask: if true - created by external MS
	//   and should not be displayed in responses to UI
	Task struct {
		ID                  gocql.UUID                `json:"id"                     valid:"unsettableByUsers"`
		Name                string                    `json:"name"                   valid:"optional"`
		Description         string                    `json:"description"            valid:"-"`
		CreatedAt           time.Time                 `json:"createdAt"              valid:"unsettableByUsers"`
		CreatedBy           string                    `json:"createdBy"              valid:"unsettableByUsers"`
		PartnerID           string                    `json:"partnerId"              valid:"unsettableByUsers"`
		OriginID            gocql.UUID                `json:"originId"               valid:"requiredForUsers"` // script or patch ID
		State               statuses.TaskState        `json:"state"                  valid:"unsettableByUsers"`
		Type                string                    `json:"type"                   valid:"validType"`
		Parameters          string                    `json:"parameters"             valid:"-"`
		ParametersObject    map[string]interface{}    `json:"parametersObject,omitempty" valid:"-"`             // same as parameters but used only for request. data from this field should always be in the parameters field
		RunTimeUTC          time.Time                 `json:"nextRunTime"            valid:"unsettableByUsers"` // special field for scheduler
		PostponedRunTime    time.Time                 `json:"postponedTime"          valid:"unsettableByUsers"`
		OriginalNextRunTime time.Time                 `json:"originalNextRunTime"    valid:"unsettableByUsers"`
		ExternalTask        bool                      `json:"externalTask"           valid:"-"`
		ResultWebhook       string                    `json:"resultWebhook"          valid:"url"`
		Targets             Target                    `json:"targets"                valid:"validTargets"`
		TargetType          TargetType                `json:"targetType"             valid:"unsettableByUsers"` // target type per internal task/endpointID
		TargetsByType       TargetsByType             `json:"targetsByType"          valid:"-"`                 // targets related data is stored here
		ManagedEndpoints    []ManagedEndpointDetailed `json:"managedEndpoints"       valid:"unsettableByUsers"`
		LastTaskInstanceID  gocql.UUID                `json:"-"                      valid:"-"`
		ManagedEndpointID   gocql.UUID                `json:"-"                      valid:"-"`
		Schedule            apiModels.Schedule        `json:"schedule"               valid:"required,validatorDynamicGroup,optionalTriggerTypes,recurrentDGTriggerTarget"`
		IsRequireNOCAccess  bool                      `json:"isRequireNOCAccess"     valid:"-"`
		ModifiedBy          string                    `json:"modifiedBy"             valid:"unsettableByUsers"`
		ModifiedAt          time.Time                 `json:"modifiedAt"             valid:"unsettableByUsers"`
		DefinitionID        gocql.UUID                `json:"definitionID"           valid:"optional"`
		Credentials         *agentModels.Credentials  `json:"credentials,omitempty"  valid:"validCreds"`
		ResourceType        integration.ResourceType  `json:"resourceType"           valid:"validResourceType"`
	}

	// TargetsByType  represents DTO of targets grouped by type
	TargetsByType map[TargetType][]string

	// ExecutionResultTaskData  represents DTO of task data for execution results
	ExecutionResultTaskData struct {
		ResultWebHook string
		CreatedBy     string
		IsNOC         bool
		Name          string
	}

	// TaskCount struct is used to display Task count per target
	TaskCount struct {
		ManagedEndpointID gocql.UUID `json:"managedEndpointId"`
		Count             int        `json:"count"`
	}

	// AllTargetsEnable contains the state which should be applied for all targets in task
	AllTargetsEnable struct {
		Active bool `json:"active"`
	}

	// SelectedManagedEndpointEnable contains the map of specific targets and their states to update in task
	SelectedManagedEndpointEnable struct {
		ManagedEndpoints map[string]bool `json:"targets"   valid:"requiredForUsers"`
	}

	// TaskNotFoundError returns in case of Task was not found in the DB
	TaskNotFoundError struct {
		ErrorParameters string
	}

	// TaskIsExpiredError returns in case of Task won't be run in future (it expired)
	TaskIsExpiredError struct {
		TaskID            gocql.UUID
		PartnerID         string
		ManagedEndpointID gocql.UUID
	}

	// CantUpdateTaskError returns in case when there is no internal task which could be updated found
	CantUpdateTaskError struct {
		TaskID               gocql.UUID
		RequirementsToUpdate interface{}
	}

	// TaskDetailsWithStatuses represents aggregated info about Task and its Statuses
	TaskDetailsWithStatuses struct {
		Task           Task
		TaskInstance   TaskInstance
		Statuses       map[string]int
		OverallStatus  string
		CanBePostponed bool // represents possibility to postpone a task
		CanBeCanceled  bool // represents possibility to cancel a task
	}

	// TaskDetailsWithExecutionResult represents aggregated info about Task and ExecutionResults` output
	TaskDetailsWithExecutionResult struct {
		Task            Task
		TaskInstance    TaskInstance
		ExecutionResult ExecutionResult
	}
)

// HasEndpointTypeOnly - returns true if DG and Site targets are presented ONLY
func (t TargetsByType) HasEndpointTypeOnly() bool {
	eIds, e := t[ManagedEndpoint]
	dIds, dg := t[DynamicGroup]
	sIds, site := t[Site]
	dynamicSiteIDs, dSite := t[DynamicSite]
	return (e && len(eIds) > 0) && (!dg || len(dIds) == 0) && (!site || len(sIds) == 0) && (!dSite || len(dynamicSiteIDs) == 0)
}

// IsValid - checks if all targetTypes are valid
func (t TargetsByType) Validate() error {
	for k := range t {
		if k != ManagedEndpoint && k != DynamicGroup && k != Site && k != DynamicSite {
			return fmt.Errorf("invalid type of targets: %v", t)
		}
	}
	return nil
}

// Contains - checks if current target ID is appeared for current targetType
func (t TargetsByType) Contains(tType TargetType, tID string) bool {
	targets, ok := t[tType]
	if !ok {
		return false
	}

	for _, target := range targets {
		if target == tID {
			return true
		}
	}
	return false
}

// UnmarshalJSON used to convert the string representation of the type of the Targets to TargetType type in TargetsByType key
func (t *TargetsByType) UnmarshalJSON(byteResult []byte) error {
	input := make(map[string][]string)

	if err := json.Unmarshal(byteResult, &input); err != nil {
		return err
	}

	if len(input) == 0 {
		return nil
	}

	targets := make(map[TargetType][]string)
	for k, v := range input {
		tt := TargetType(0)

		if err := json.Unmarshal([]byte(escape+k+escape), &tt); err != nil {
			return err
		}

		targets[tt] = v
	}

	*t = targets

	return nil
}

// MarshalJSON custom marshal method for TargetsByType type
func (t TargetsByType) MarshalJSON() ([]byte, error) {
	output := make(map[string][]string)

	for k, v := range t {
		bType, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}

		tt := strings.Replace(string(bType), escape, "", anyTimesReplace)
		output[tt] = v
	}

	return json.Marshal(output)
}

// Error - error interface implementation for TaskNotFoundError type
func (err TaskNotFoundError) Error() string {
	return fmt.Sprintf("No Task found. %s", err.ErrorParameters)
}

// Error - error interface implementation for CantUpdateTaskError type
func (err CantUpdateTaskError) Error() string {
	return fmt.Sprintf("Task (ID = %s) can't be updated with such requirements (%v).", err.TaskID, err.RequirementsToUpdate)
}

// Error - error interface implementation for TaskIsExpiredError type
func (err TaskIsExpiredError) Error() string {
	return fmt.Sprintf("Task (ID = %s) created by partner(%s) for ManagedEndpointID (%s) is expired", err.TaskID, err.PartnerID, err.ManagedEndpointID)
}

// IsRecurringlyScheduled checks if the task is still scheduled
func (task *Task) IsRecurringlyScheduled() bool {
	return task.Schedule.Regularity == apiModels.Recurrent && common.ValidTime(task.Schedule.EndRunTime.UTC(), task.RunTimeUTC)
}

// IsScheduled checks if the task should be run by scheduler
func (task *Task) IsScheduled() bool {
	if task.State != statuses.TaskStateActive {
		return false
	}
	return task.IsRecurringlyScheduled() || task.Schedule.Regularity == apiModels.OneTime
}

// IsTaskAndTriggerNotActivated checks if the task has trigger and recurrence executions
func (task *Task) IsTaskAndTriggerNotActivated() bool {
	return task.Schedule.Regularity == apiModels.Recurrent && len(task.Schedule.TriggerTypes) != 0 && !task.OriginalNextRunTime.IsZero() &&
		task.RunTimeUTC == task.Schedule.StartRunTime.UTC()
}

// IsTrigger checks if the task should be run by scheduler
func (task *Task) IsTrigger() bool {
	switch task.Schedule.Regularity {
	case apiModels.Trigger:
		return true
	default:
		return false
	}
}

// IsActivatedTrigger returns true if trigger has been activated by scheduler (we can undestand this by comparing RunTime and EndRunTime)
func (task *Task) IsActivatedTrigger() bool {
	now := time.Now().UTC()
	if task.Schedule.EndRunTime.IsZero() { // its 'Never' task trigger
		return task.RunTimeUTC.After(now)
	}
	return task.Schedule.EndRunTime.Equal(task.RunTimeUTC) // when trigger is activated they are equal
}

// IsDynamicGroupBasedTrigger says that task is dynamicGroupEnter or Exit trigger
func (task *Task) IsDynamicGroupBasedTrigger() bool {
	return task.IsDynamicGroupEnterTrigger() || task.IsDynamicGroupExitTrigger()
}

// IsDynamicGroupEnterTrigger checks if this is DynamicGroupEnter trigger
func (task *Task) IsDynamicGroupEnterTrigger() bool {
	switch task.Schedule.Regularity {
	case apiModels.Trigger:
		return contains(triggers.DynamicGroupEnterTrigger, task.Schedule.TriggerTypes)
	}
	return false
}

// IsDynamicGroupExitTrigger checks if this is DynamicGroupExit trigger
func (task *Task) IsDynamicGroupExitTrigger() bool {
	switch task.Schedule.Regularity {
	case apiModels.Trigger:
		return contains(triggers.DynamicGroupExitTrigger, task.Schedule.TriggerTypes)
	}
	return false
}

// IsRunAsUserApplied checks if run as user option applied
func (task *Task) IsRunAsUserApplied() bool {
	return nil != task.Credentials
}

func contains(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// IsExpired checks if the task is not scheduled for the future
func (task *Task) IsExpired() bool {
	switch task.Schedule.Regularity {
	case apiModels.OneTime, apiModels.Trigger:
		return common.ValidTime(time.Now().UTC(), task.RunTimeUTC)
	case apiModels.Recurrent:
		return !common.ValidTime(task.Schedule.EndRunTime.UTC(), task.RunTimeUTC)
	}
	return true
}

// HasOriginalRunTime checks if the task is postponed
func (task *Task) HasOriginalRunTime() bool {
	return !task.OriginalNextRunTime.IsZero()
}

// HasPostponedTime checks if the task is postponed
func (task *Task) HasPostponedTime() bool {
	return !task.PostponedRunTime.IsZero()
}

// Enable updates the task state to Active for Disabled Task
func (task *Task) Enable(ctx context.Context) (bool, error) {
	if task.State != statuses.TaskStateDisabled {
		return false, nil
	}

	// Check if the task is not expired
	if task.IsExpired() {
		return false, TaskIsExpiredError{TaskID: task.ID, PartnerID: task.PartnerID, ManagedEndpointID: task.ManagedEndpointID}
	}

	// Do not update RunTimeInSeconds if it still valid
	if common.ValidTime(task.RunTimeUTC, time.Now().UTC().Truncate(time.Minute)) {
		task.State = statuses.TaskStateActive
		return true, nil
	}

	// Try to update task.RunTimeInSeconds for Recurrent task
	task.RunTimeUTC = time.Now().UTC().Truncate(time.Minute)
	if _, err := UpdateTaskByNextRunTime(ctx, task); err != nil {
		return false, err
	}

	// Check if the task doesn't become expired after time update
	if task.IsExpired() {
		return false, TaskIsExpiredError{TaskID: task.ID, PartnerID: task.PartnerID, ManagedEndpointID: task.ManagedEndpointID}
	}

	task.State = statuses.TaskStateActive
	return true, nil
}

// Disable updates the task state to Disable for Active Task
func (task *Task) Disable() (bool, error) {
	// try to disable Task
	if task.State != statuses.TaskStateActive {
		return false, nil
	}
	if task.IsExpired() {
		return false, TaskIsExpiredError{TaskID: task.ID, PartnerID: task.PartnerID, ManagedEndpointID: task.ManagedEndpointID}
	}
	task.State = statuses.TaskStateDisabled
	return true, nil
}

// CopyWithRunTime creates a copy of Internal Task with the same RunTimeInSeconds value for specific managedEndpoint
func (task *Task) CopyWithRunTime(managedEndpoint gocql.UUID) *Task {
	return &Task{
		ID:                  task.ID,
		Name:                task.Name,
		Description:         task.Description,
		CreatedAt:           task.CreatedAt,
		CreatedBy:           task.CreatedBy,
		PartnerID:           task.PartnerID,
		OriginID:            task.OriginID,
		State:               task.State,
		RunTimeUTC:          task.RunTimeUTC,
		LastTaskInstanceID:  task.LastTaskInstanceID,
		OriginalNextRunTime: task.OriginalNextRunTime,
		PostponedRunTime:    task.PostponedRunTime,
		Type:                task.Type,
		Parameters:          task.Parameters,
		ExternalTask:        task.ExternalTask,
		ResultWebhook:       task.ResultWebhook,
		ManagedEndpointID:   managedEndpoint,
		TargetType:          task.TargetType,
		Schedule:            task.Schedule,
		IsRequireNOCAccess:  task.IsRequireNOCAccess,
		ModifiedBy:          task.ModifiedBy,
		ModifiedAt:          task.ModifiedAt,
		DefinitionID:        task.DefinitionID,
		Credentials:         task.Credentials,
		TargetsByType:       task.TargetsByType,
		ResourceType:        task.ResourceType,
	}
}

// TaskPersistence interface to perform actions with Task database
type TaskPersistence interface {
	InsertOrUpdate(ctx context.Context, tasks ...Task) error
	Delete(ctx context.Context, tasks []Task) error
	UpdateTask(ctx context.Context, inputStruct interface{}, partnerID string, taskID gocql.UUID) error
	UpdateSchedulerFields(ctx context.Context, tasks ...Task) error
	UpdateModifiedFieldsByMEs(ctx context.Context, task Task, managedEndpoints ...gocql.UUID) error

	GetByPartner(ctx context.Context, partnerID string) ([]Task, error)
	GetByIDs(ctx context.Context, cache persistency.Cache, partnerID string, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]Task, error)
	GetByLastTaskInstanceIDs(ctx context.Context, partnerID string, endpointID gocql.UUID, lastTaskInstanceIDs ...gocql.UUID) (tasksByIDMap map[gocql.UUID]Task, err error)
	GetByIDAndManagedEndpoints(ctx context.Context, partnerID string, taskID gocql.UUID, managedEndpointIDs ...gocql.UUID) ([]Task, error)
	GetExecutionResultTaskData(partnerID string, taskID, endpointID gocql.UUID) (ExecutionResultTaskData, error)
	GetByRunTimeRange(ctx context.Context, runTimeRange []time.Time) ([]Task, error)
	GetByPartnerAndTime(ctx context.Context, partnerID string, timeToSearchFrom time.Time) ([]Task, error)
	GetByPartnerAndManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID, count int) ([]Task, error)
	GetCountByManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID) (TaskCount, error)
	GetCountsByPartner(ctx context.Context, partnerID string) ([]TaskCount, error)
	GetTasksFilteredWithTime(ctx context.Context, partnerID string, taskID gocql.UUID, excludedEndpoints map[gocql.UUID]Task, after time.Time) (map[gocql.UUID][]Task, error)
	GetManagedEndpointIDsOfActiveTasks(ctx context.Context, partnerID string, taskID gocql.UUID) (map[gocql.UUID]struct{}, error)
	GetTargetTypeByEndpoint(partnerID string, taskID, endpointID gocql.UUID, external bool) (TargetType, error)
}

// TaskRepoCassandra is a realisation of TaskPersistence interface for Cassandra
type TaskRepoCassandra struct{}

// GetExecutionResultTaskData returns execution results Data
func (repo TaskRepoCassandra) GetExecutionResultTaskData(partnerID string, taskID, endpointID gocql.UUID) (ExecutionResultTaskData, error) {

	queryInternal := `SELECT name, created_by, result_webhook, require_noc_access FROM tasks WHERE id = ? 
				AND partner_id = ? AND managed_endpoint_id = ? and external_task = false`

	queryExternal := `SELECT name, created_by, result_webhook, require_noc_access FROM tasks WHERE id = ? 
				AND partner_id = ? AND managed_endpoint_id = ? and external_task = true`

	data, err := repo.getTaskExecResultData(queryInternal, taskID, partnerID, endpointID)
	if err != nil {
		return ExecutionResultTaskData{}, err
	}
	if data != nil {
		return *data, nil
	}

	data, err = repo.getTaskExecResultData(queryExternal, taskID, partnerID, endpointID)
	if err != nil {
		return ExecutionResultTaskData{}, err
	}
	if data != nil {
		return *data, nil
	}
	return ExecutionResultTaskData{}, fmt.Errorf("not found task by id %v", taskID)
}

func (repo TaskRepoCassandra) getTaskExecResultData(query string, values ...interface{}) (task *ExecutionResultTaskData, err error) {
	var res ExecutionResultTaskData
	if err = cassandra.QueryCassandra(context.TODO(),
		query, values...).Scan(
		&res.Name,
		&res.CreatedBy,
		&res.ResultWebHook,
		&res.IsNOC); err == nil {
		return &res, nil
	}

	if err == gocql.ErrNotFound {
		return nil, nil
	}
	return nil, fmt.Errorf(errWorkingWithEntities, err)

}

var (
	// TaskPersistenceInstance is an instance presented TaskRepoCassandra
	TaskPersistenceInstance TaskPersistence = TaskRepoCassandra{}
)

type insertTaskDTO struct {
	task              Task
	ctx               context.Context
	tasksBeforeUpdate map[string]Task
	batch             cassandra.IBatch
	tasksToBeUpserted map[string]struct{}
	batchTasks2       cassandra.IBatch
	tasksForUpdate    map[string][]interface{}
}

// InsertOrUpdate inserts or updates task in database
func (repo TaskRepoCassandra) InsertOrUpdate(ctx context.Context, internalTasks ...Task) error {
	batch := cassandra.Session.NewBatch(gocql.UnloggedBatch)
	// new batch for tasks2 table
	batchTasks2 := cassandra.Session.NewBatch(gocql.UnloggedBatch)
	tasksBeforeUpdate := make(map[string]Task)
	tasksForUpdate := make(map[string][]interface{})
	tasksToBeUpserted := make(map[string]struct{})

	dto := &insertTaskDTO{
		ctx:               ctx,
		tasksBeforeUpdate: tasksBeforeUpdate,
		batch:             batch,
		tasksToBeUpserted: tasksToBeUpserted,
		batchTasks2:       batchTasks2,
		tasksForUpdate:    tasksForUpdate,
	}

	for i, task := range internalTasks {
		dto.task = task
		if err := repo.insertTask(dto); err != nil {
			return err
		}

		// no more than 30 in 1 batch or all if it's the last iteration
		if (i+1)%config.Config.CassandraBatchSize != 0 && i+1 != len(internalTasks) {
			continue
		}

		if err := cassandra.Session.ExecuteBatch(dto.batch); err != nil {
			return err
		}

		if batchTasks2.Size() > 0 {
			err := cassandra.Session.ExecuteBatch(batchTasks2)
			if err != nil {
				return err
			}
		}

		if err := repo.mutateMVTables(dto.tasksForUpdate, dto.tasksBeforeUpdate); err != nil {
			return err
		}

		dto.tasksForUpdate = make(map[string][]interface{})
		dto.tasksBeforeUpdate = make(map[string]Task)
		dto.batch = cassandra.Session.NewBatch(gocql.UnloggedBatch)
		dto.batchTasks2 = cassandra.Session.NewBatch(gocql.UnloggedBatch)
	}

	return nil
}

func (repo TaskRepoCassandra) insertTask(o *insertTaskDTO) error {
	task := o.task
	for targetType, ids := range task.TargetsByType { // done for rollback purposes (resource selector)
		task.Targets.IDs = ids
		task.Targets.Type = targetType
		break
	}

	// targets ids data is redundant because IDs is stored in ManagedEndpointsID field and its copied for each internal Task
	// and we faced performance issues
	if task.TargetType == ManagedEndpoint {
		task.Targets.IDs = []string{}
	}

	taskCompositeKey := task.PartnerID + ":" + task.ID.String() + ":" + task.ManagedEndpointID.String() + ":" + strconv.FormatBool(task.ExternalTask)
	oldTasks, err := selectTasks(o.ctx, selectTaskBeforeUpdate,
		task.PartnerID, task.ID, task.ManagedEndpointID, task.ExternalTask)

	if err != nil {
		return fmt.Errorf("TaskRepoCassandra.InsertOrUpdate: Error while trying to retrieve task by partner_id=%s, id=%v, managed_endpoint_id=%v and external_task=%v : %v",
			task.PartnerID, task.ID, task.ManagedEndpointID, task.ExternalTask, err)
	}

	if len(oldTasks) != 0 {
		o.tasksBeforeUpdate[taskCompositeKey] = oldTasks[0]
	}

	taskTTL := 0

	if task.State != statuses.TaskStateActive {
		taskTTL = config.Config.DataRetentionIntervalDay * secondsInDay
	}

	scheduleBytes, err := json.Marshal(task.Schedule)
	if err != nil {
		return err
	}

	if task.RunTimeUTC.IsZero() {
		task.RunTimeUTC = task.CreatedAt
	}

	taskFields := []interface{}{
		task.DefinitionID,
		task.ID,
		task.Name,
		task.Description,
		task.ManagedEndpointID,
		task.CreatedAt,
		task.CreatedBy,
		task.PartnerID,
		task.OriginID,
		task.State,
		task.RunTimeUTC,
		task.Type,
		task.Parameters,
		task.ExternalTask,
		task.ResultWebhook,
		task.LastTaskInstanceID,
		task.IsRequireNOCAccess,
		task.ModifiedBy,
		task.ModifiedAt,
		task.OriginalNextRunTime,
		task.PostponedRunTime,
		task.Targets.IDs,
		task.TargetType,
		string(scheduleBytes),
		task.Credentials,
		taskTTL,
	}

	o.batch.Query(`INSERT INTO tasks (definition_id, id, name, description, 
			                            managed_endpoint_id, created_at, created_by, partner_id, origin_id, 
			                            state, run_time_unix, type, parameters, 
			                            external_task, result_webhook, last_task_instance_id, require_noc_access, modified_by, 
			                            modified_at, original_next_run_time, run_time, targets, target_type, 
			                            schedule, credentials)
                                        VALUES  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) USING TTL ?`, taskFields...)

	taskID := task.ID.String()
	// new tasks fields for new tasks2 table
	if _, ok := o.tasksToBeUpserted[taskID]; !ok {
		targets := repo.getTargetsFromTask(task)
		task2Fields := []interface{}{
			task.DefinitionID,
			task.ID,
			task.Name,
			task.Description,
			task.CreatedAt,
			task.CreatedBy,
			task.PartnerID,
			task.OriginID,
			task.Type,
			task.Parameters,
			task.ExternalTask,
			task.ResultWebhook,
			task.IsRequireNOCAccess,
			task.ModifiedBy,
			task.ModifiedAt,
			targets,
			string(scheduleBytes),
			task.Credentials,
			task.ResourceType,
		}

		o.batchTasks2.Query(`INSERT INTO tasks2 (definition_id, id, name, description, 
			                                     created_at, created_by,
												 partner_id, origin_id, 
			                                     type, parameters, external_task, result_webhook,
                                                 require_noc_access, modified_by, 
												 modified_at, targets,
                                                 schedule, credentials, resources_type)
                                            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, task2Fields...)
		o.tasksToBeUpserted[taskID] = struct{}{}
	}

	o.tasksForUpdate[taskCompositeKey] = taskFields
	return nil
}

func (repo TaskRepoCassandra) getTargetsFromTask(task Task) map[TargetType][]string {
	targets := task.TargetsByType
	if len(targets) == 0 {
		if targets == nil {
			targets = make(TargetsByType)
		}
		targets[task.TargetType] = task.Targets.IDs
	}

	for k, v := range targets {
		if v == nil {
			targets[k] = make([]string, 0)
		}
	}

	return targets
}

func (repo TaskRepoCassandra) mutateMVTables(newTasks map[string][]interface{}, oldTasks map[string]Task) error {
	var (
		totalErr error
		errCh    = make(chan error)
		doneErr  = make(chan struct{})
		tables   = []string{
			taskByRuntimeUnixTableName,
			tasksOrderByLastTaskInstanceIDTableName,
			tasksByRuntimeTableName,
			tasksByIDTableName,
		}
	)

	go func() {
		for err := range errCh {
			totalErr = fmt.Errorf("%v:%v", totalErr, err)
		}
		doneErr <- struct{}{}
	}()

	var wg sync.WaitGroup
	for _, tableName := range tables {
		wg.Add(1)
		go func(tableName string, newTasks map[string][]interface{}, oldTasks map[string]Task) {
			defer wg.Done()
			for taskCompositeKey, newTask := range newTasks {
				repo.updateTaskPerCompositeKey(oldTasks, taskCompositeKey, tableName, errCh, newTask)
			}
		}(tableName, newTasks, oldTasks)
	}

	wg.Wait()
	close(errCh)
	<-doneErr
	return totalErr
}

func (repo TaskRepoCassandra) updateTaskPerCompositeKey(oldTasks map[string]Task, taskCompositeKey string, tableName string, errCh chan error, newTask []interface{}) {
	batch := cassandra.Session.NewBatch(gocql.LoggedBatch)
	if oldTask, ok := oldTasks[taskCompositeKey]; ok {
		switch tableName {
		case taskByRuntimeUnixTableName, tasksByRuntimeTableName, tasksByIDTableName:
			query := `DELETE
						            FROM ` + tableName + ` USING TIMESTAMP ? WHERE partner_id = ? 
						            AND id = ? AND managed_endpoint_id = ? 
						            AND external_task = ? AND run_time_unix = ?`
			params := []interface{}{
				time.Now().UnixNano() / int64(time.Microsecond),
				oldTask.PartnerID,
				oldTask.ID,
				oldTask.ManagedEndpointID,
				oldTask.ExternalTask,
				oldTask.RunTimeUTC,
			}
			batch.Query(query, params...)
		case tasksOrderByLastTaskInstanceIDTableName:
			query := `DELETE
						            FROM tasks_order_by_last_task_instance_id_mv 
						            USING TIMESTAMP ? WHERE partner_id = ? AND id = ? 
						            AND managed_endpoint_id = ? AND external_task = ? AND last_task_instance_id = ?`
			params := []interface{}{
				time.Now().UnixNano() / int64(time.Microsecond),
				oldTask.PartnerID,
				oldTask.ID,
				oldTask.ManagedEndpointID,
				oldTask.ExternalTask,
				oldTask.LastTaskInstanceID,
			}
			batch.Query(query, params...)
		default:
			errCh <- fmt.Errorf("where is no handlers for %s table", tableName)
			return
		}
	}
	newTask = append(newTask, time.Now().Add(1*time.Microsecond).UnixNano()/int64(time.Microsecond))
	batch.Query(fmt.Sprintf(`INSERT 
					                        INTO %s 
					                           (definition_id, id, name, description, 
					                            managed_endpoint_id, created_at, created_by, partner_id, origin_id, 
					                            state, run_time_unix, type, parameters, 
					                            external_task, result_webhook, last_task_instance_id, require_noc_access, modified_by, 
					                            modified_at, original_next_run_time, run_time, targets, target_type, 
					                            schedule, credentials) 
					                        VALUES  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
					                        USING TTL ?
					                        AND TIMESTAMP ?`,
		tableName), newTask...)

	if err := cassandra.Session.ExecuteBatch(batch); err != nil {
		errCh <- err
		return
	}
}

func updateSelectedTaskTargets(ctx context.Context, tasks []Task, selectedTargetsEnable *SelectedManagedEndpointEnable) ([]Task, error) {
	var updatedInternalTasks = make([]Task, 0)
	var (
		ok  bool
		err error
	)
	for i, task := range tasks {
		if selectedTargetsEnable.ManagedEndpoints[task.ManagedEndpointID.String()] {
			ok, err = tasks[i].Enable(ctx)
		} else {
			ok, err = tasks[i].Disable()
		}
		if err != nil {
			return nil, err
		}
		if ok {
			updatedInternalTasks = append(updatedInternalTasks, tasks[i])
		}
	}
	return updatedInternalTasks, nil
}

func updateAllTaskTargets(ctx context.Context, tasks []Task, allTargetsEnable *AllTargetsEnable) ([]Task, error) {
	var updatedInternalTasks = make([]Task, 0)
	var (
		ok  bool
		err error
	)
	for i := range tasks {
		if allTargetsEnable.Active {
			ok, err = tasks[i].Enable(ctx)
		} else {
			ok, err = tasks[i].Disable()
		}
		if err != nil {
			return nil, err
		}
		if ok {
			updatedInternalTasks = append(updatedInternalTasks, tasks[i])
		}
	}
	return updatedInternalTasks, nil
}

// Delete deletes Tasks from database
func (TaskRepoCassandra) Delete(ctx context.Context, tasks []Task) error {
	var (
		totalErr error
		wg       sync.WaitGroup
		errCh    = make(chan error)
		doneErr  = make(chan struct{})
		tables   = []string{
			baseTable,
			taskByRuntimeUnixTableName,
			tasksOrderByLastTaskInstanceIDTableName,
			tasksByRuntimeTableName,
			tasksByIDTableName,
		}
	)

	go func() {
		for err := range errCh {
			totalErr = fmt.Errorf("%v:%v", totalErr, err)
		}
		close(doneErr)
	}()

	for _, tableName := range tables {
		wg.Add(1)
		go func(tasks []Task, tableName string) {
			defer wg.Done()
			for _, task := range tasks {
				query := ""
				var params []interface{}
				switch tableName {
				case taskByRuntimeUnixTableName, tasksByRuntimeTableName, tasksByIDTableName:
					query = `DELETE FROM ` + tableName + ` 
						            WHERE partner_id = ? 
						            AND id = ? 
						            AND managed_endpoint_id = ? 
						            AND external_task = ? 
						            AND run_time_unix = ? IF EXISTS`
					params = []interface{}{
						task.PartnerID,
						task.ID,
						task.ManagedEndpointID,
						task.ExternalTask,
						task.RunTimeUTC,
					}
				case baseTable:
					query = `DELETE FROM ` + tableName + ` 
						            WHERE partner_id = ? 
						            AND id = ? 
						            AND managed_endpoint_id = ? 
						            AND external_task = ? IF EXISTS`
					params = []interface{}{
						task.PartnerID,
						task.ID,
						task.ManagedEndpointID,
						task.ExternalTask,
					}
				case tasksOrderByLastTaskInstanceIDTableName:
					query = `DELETE FROM tasks_order_by_last_task_instance_id_mv
						            WHERE partner_id = ? 
						            AND id = ? 
						            AND managed_endpoint_id = ? 
						            AND external_task = ? 
						            AND last_task_instance_id = ? IF EXISTS`
					params = []interface{}{
						task.PartnerID,
						task.ID,
						task.ManagedEndpointID,
						task.ExternalTask,
						task.LastTaskInstanceID,
					}
				default:
					errCh <- fmt.Errorf("no handler for %s", tableName)
					return
				}

				if err := cassandra.QueryCassandra(ctx, query, params...).Exec(); err != nil {
					errCh <- fmt.Errorf("error during deleting task from %v err: %v", tableName, err)
					return
				}
			}
		}(tasks, tableName)
	}
	wg.Wait()

	close(errCh)
	<-doneErr
	return totalErr
}

// UpdateTask retrieves a task by id and partner_id and update it based on inputStruct
func (t TaskRepoCassandra) UpdateTask(ctx context.Context, inputStruct interface{}, partnerID string, taskID gocql.UUID) (err error) {
	var tasks []Task
	// update the task
	switch v := inputStruct.(type) {
	case *SelectedManagedEndpointEnable:
		if tasks, err = t.getTasksBySelectedEndpoints(v, ctx, partnerID, taskID); err != nil {
			return err
		}
	case *AllTargetsEnable:
		tasks, err = TaskPersistenceInstance.GetByIDs(ctx, nil, partnerID, false, taskID)
		if err != nil {
			return
		}
		if tasks, err = updateAllTaskTargets(ctx, tasks, v); err != nil {
			return
		}
	default:
		return errors.New("wrong input type for update")
	}

	if len(tasks) == 0 {
		return CantUpdateTaskError{taskID, inputStruct}
	}

	// Update only scheduled instance
	currentTaskInstances, err := TaskInstancePersistenceInstance.GetByIDs(ctx, tasks[0].LastTaskInstanceID)
	if err != nil {
		return err
	}

	if len(currentTaskInstances) == 0 {
		return CantUpdateTaskError{taskID, inputStruct}
	}

	task := tasks[0]
	inst, err := t.getInstanceToDisable(task, currentTaskInstances[0])
	if err != nil {
		return err
	}

	inst.Statuses = t.setInstanceStatuses(tasks, inst)

	if err = TaskInstancePersistenceInstance.Insert(ctx, inst); err != nil {
		return err
	}

	return TaskPersistenceInstance.InsertOrUpdate(ctx, tasks...)
}

func (t TaskRepoCassandra) getInstanceToDisable(task Task, currentTaskInstance TaskInstance) (TaskInstance, error) {
	inst := currentTaskInstance
	for {
		gotTI, err := TaskInstancePersistenceInstance.GetNearestInstanceAfter(task.ID, inst.StartedAt)
		if err == gocql.ErrNotFound {
			break
		}

		if err != nil {
			return TaskInstance{}, fmt.Errorf("can't get nearest instance for running. partnerID : %v, started_at : %v, cause: %s", task.PartnerID, inst.StartedAt, err.Error())
		}

		if gotTI.TriggeredBy == "" {
			inst.PartnerID = currentTaskInstance.PartnerID
			inst.OverallStatus = currentTaskInstance.OverallStatus
			inst = gotTI
			break
		}
		inst = gotTI
	}
	return inst, nil
}

func (TaskRepoCassandra) setInstanceStatuses(tasks []Task, inst TaskInstance) map[gocql.UUID]statuses.TaskInstanceStatus {
	for _, t := range tasks {
		if t.State == statuses.TaskStateDisabled {
			inst.Statuses[t.ManagedEndpointID] = statuses.TaskInstanceDisabled
		}
		if t.State == statuses.TaskStateActive {
			if t.HasPostponedTime() || t.HasOriginalRunTime() {
				inst.Statuses[t.ManagedEndpointID] = statuses.TaskInstancePostponed
				continue
			}
			inst.Statuses[t.ManagedEndpointID] = statuses.TaskInstanceScheduled
		}
	}
	return inst.Statuses
}

func (TaskRepoCassandra) getTasksBySelectedEndpoints(v *SelectedManagedEndpointEnable, ctx context.Context, partnerID string, taskID gocql.UUID) ([]Task, error) {
	managedEndpointsToUpdateUUID, err := convertSelectedManagedEndpointEnableToListOfManagedEndpointIDs(v)
	if err != nil {
		return nil, err
	}

	tasks, err := TaskPersistenceInstance.GetByIDAndManagedEndpoints(ctx, partnerID, taskID, managedEndpointsToUpdateUUID...)
	if err != nil {
		return nil, err
	}

	if tasks, err = updateSelectedTaskTargets(ctx, tasks, v); err != nil {
		return nil, err
	}
	return tasks, nil
}

type updateSchedulerFieldsDTO struct {
	ctx               context.Context
	tasksBeforeUpdate map[string]Task
	tasksToBeUpserted map[string]struct{}
	batch             cassandra.IBatch
	batchTasks2       cassandra.IBatch
	tasksForUpdate    map[string][]interface{}
	lenOfTasks        int
}

// UpdateSchedulerFields inserts or updates scheduler specific fields in database
func (repo TaskRepoCassandra) UpdateSchedulerFields(ctx context.Context, internalTasks ...Task) error {
	batch := cassandra.Session.NewBatch(gocql.UnloggedBatch)
	// new batch for tasks2 table
	batchTasks2 := cassandra.Session.NewBatch(gocql.UnloggedBatch)
	tasksForUpdate := make(map[string][]interface{})
	tasksBeforeUpdate := make(map[string]Task)
	tasksToBeUpserted := make(map[string]struct{})
	dto := &updateSchedulerFieldsDTO{
		ctx:               ctx,
		tasksBeforeUpdate: tasksBeforeUpdate,
		tasksToBeUpserted: tasksToBeUpserted,
		batch:             batch,
		batchTasks2:       batchTasks2,
		tasksForUpdate:    tasksForUpdate,
		lenOfTasks:        len(internalTasks),
	}

	for i, task := range internalTasks {
		if err := repo.updateSchedulerFieldsTask(task, dto, i); err != nil {
			return err
		}
	}
	return nil
}

func (repo TaskRepoCassandra) updateSchedulerFieldsTask(task Task, o *updateSchedulerFieldsDTO, i int) error {
	taskCompositeKey := task.PartnerID + ":" + task.ID.String() + ":" + task.ManagedEndpointID.String() + ":" + strconv.FormatBool(task.ExternalTask)

	oldTasks, err := selectTasks(o.ctx, selectTaskBeforeUpdate,
		task.PartnerID, task.ID, task.ManagedEndpointID, task.ExternalTask)

	if err != nil {
		return fmt.Errorf("TaskRepoCassandra.UpdateSchedulerFields: Error while trying to retrieve task by partner_id=%s, id=%v, managed_endpoint_id=%v and external_task=%v : %v",
			task.PartnerID, task.ID, task.ManagedEndpointID, task.ExternalTask, err)
	}

	if len(oldTasks) == 0 {
		return nil
	}

	o.tasksBeforeUpdate[taskCompositeKey] = oldTasks[0]
	oldTask := oldTasks[0]

	taskTTL := 0

	if task.State != statuses.TaskStateActive {
		taskTTL = config.Config.DataRetentionIntervalDay * secondsInDay
	}

	scheduleBytes, err := json.Marshal(task.Schedule)
	if err != nil {
		return fmt.Errorf("marshall schedule err: %v", err)
	}

	taskFields := []interface{}{
		taskTTL,
		task.RunTimeUTC,
		task.OriginalNextRunTime,
		task.PostponedRunTime,
		string(scheduleBytes),
		task.LastTaskInstanceID,
		task.State,
		task.TargetType,
		task.PartnerID,
		task.ID,
		task.ManagedEndpointID,
		task.ExternalTask,
	}

	o.batch.Query(`UPDATE tasks USING TTL ? SET run_time_unix = ?, original_next_run_time = ?, run_time = ?, schedule = ?,
					 last_task_instance_id = ?, state = ?, target_type = ? WHERE partner_id = ? AND id = ? AND managed_endpoint_id = ? AND external_task = ?`, taskFields...)

	taskID := task.ID.String()
	// new task fields for tasks2 table
	if _, ok := o.tasksToBeUpserted[taskID]; !ok {
		task2Fields := []interface{}{
			string(scheduleBytes),
			task.PartnerID,
			task.ID,
		}
		o.batchTasks2.Query(`UPDATE tasks2 SET schedule = ?
					                 WHERE partner_id = ? AND id = ?`, task2Fields...)

		o.tasksToBeUpserted[taskID] = struct{}{}
	}

	if oldTask.TargetType == ManagedEndpoint {
		oldTask.Targets.IDs = nil
	}

	o.tasksForUpdate[taskCompositeKey] = []interface{}{
		oldTask.DefinitionID,
		task.ID,
		oldTask.Name,
		oldTask.Description,
		task.ManagedEndpointID,
		oldTask.CreatedAt,
		oldTask.CreatedBy,
		task.PartnerID,
		oldTask.OriginID,
		task.State,
		task.RunTimeUTC,
		oldTask.Type,
		oldTask.Parameters,
		task.ExternalTask,
		oldTask.ResultWebhook,
		task.LastTaskInstanceID,
		oldTask.IsRequireNOCAccess,
		oldTask.ModifiedBy,
		oldTask.ModifiedAt,
		task.OriginalNextRunTime,
		task.PostponedRunTime,
		oldTask.Targets.IDs,
		oldTask.TargetType,
		string(scheduleBytes),
		oldTask.Credentials,
		taskTTL,
	}

	// no more than 30 in 1 batch or all if it's the last iteration
	if (i+1)%config.Config.CassandraBatchSize != 0 && i+1 != o.lenOfTasks {
		return nil
	}

	err = cassandra.Session.ExecuteBatch(o.batch)
	if err != nil {
		return err
	}

	if o.batchTasks2.Size() > 0 {
		err = cassandra.Session.ExecuteBatch(o.batchTasks2)
		if err != nil {
			return err
		}
	}

	err = repo.mutateMVTables(o.tasksForUpdate, o.tasksBeforeUpdate)
	if err != nil {
		return err
	}

	o.tasksForUpdate = make(map[string][]interface{})
	o.tasksBeforeUpdate = make(map[string]Task)
	o.batch = cassandra.Session.NewBatch(gocql.UnloggedBatch)
	o.batchTasks2 = cassandra.Session.NewBatch(gocql.UnloggedBatch)

	return nil
}

// UpdateModifiedFieldsByMEs updated modified fields by managedEndpointIDs
func (repo TaskRepoCassandra) UpdateModifiedFieldsByMEs(ctx context.Context, task Task, managedEndpoints ...gocql.UUID) error {
	tasksBeforeUpdate := make(map[string]Task)
	tasksForUpdate := make(map[string][]interface{})

	taskFields := []interface{}{
		task.ModifiedAt,
		task.ModifiedBy,
		task.PartnerID,
		task.ID,
		task.ExternalTask,
	}

	// new fields for tasks2 table
	task2Fields := []interface{}{
		task.ModifiedAt,
		task.ModifiedBy,
		task.PartnerID,
		task.ID,
	}

	for _, id := range managedEndpoints {
		taskFields = append(taskFields, id)
		oldTasks, err := selectTasks(ctx, `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                                                          run_time_unix, run_time, state, target_type
                                                      FROM tasks
                                                      WHERE partner_id = ? AND id = ? AND managed_endpoint_id = ? and external_task = ?`,
			task.PartnerID, task.ID, id, task.ExternalTask)

		if err != nil {
			return fmt.Errorf("TaskRepoCassandra.UpdateSchedulerFields: Error while trying to retrieve task by partner_id=%s, id=%v, managed_endpoint_id=%v and external_task=%v : %v",
				task.PartnerID, task.ID, id, task.ExternalTask, err)
		}

		if len(oldTasks) == 0 {
			continue
		}

		oldTask := oldTasks[0]
		taskCompositeKey := oldTask.PartnerID + ":" + oldTask.ID.String() + ":" + oldTask.ManagedEndpointID.String() + ":" + strconv.FormatBool(oldTask.ExternalTask)

		tasksBeforeUpdate[taskCompositeKey] = oldTask

		taskTTL := 0

		if task.State != statuses.TaskStateActive {
			taskTTL = config.Config.DataRetentionIntervalDay * secondsInDay
		}

		scheduleBytes, err := json.Marshal(oldTask.Schedule)
		if err != nil {
			return err
		}

		tasksForUpdate[taskCompositeKey] = []interface{}{
			oldTask.DefinitionID,
			task.ID,
			oldTask.Name,
			oldTask.Description,
			oldTask.ManagedEndpointID,
			oldTask.CreatedAt,
			oldTask.CreatedBy,
			task.PartnerID,
			oldTask.OriginID,
			oldTask.State,
			oldTask.RunTimeUTC,
			oldTask.Type,
			oldTask.Parameters,
			task.ExternalTask,
			oldTask.ResultWebhook,
			oldTask.LastTaskInstanceID,
			oldTask.IsRequireNOCAccess,
			task.ModifiedBy,
			task.ModifiedAt,
			oldTask.OriginalNextRunTime,
			oldTask.PostponedRunTime,
			oldTask.Targets.IDs,
			oldTask.TargetType,
			string(scheduleBytes),
			oldTask.Credentials,
			taskTTL,
		}
	}

	query := fmt.Sprintf(`UPDATE tasks SET modified_at = ?, modified_by = ? WHERE partner_id = ? AND id = ? AND external_task = ?
							and managed_endpoint_id IN (%s)`, common.GetQuestionMarkString(len(managedEndpoints)))

	// new query for tasks2 table
	queryTask2 := fmt.Sprintf(`UPDATE tasks2 SET modified_at = ?, modified_by = ? WHERE partner_id = ? AND id = ?`)

	if err := cassandra.QueryCassandra(ctx, query, taskFields...).Exec(); err != nil {
		return err
	}
	if err := cassandra.QueryCassandra(ctx, queryTask2, task2Fields...).Exec(); err != nil {
		return err
	}

	tasksForUpdateForBatch := make(map[string][]interface{})
	tasksBeforeUpdateForBatch := make(map[string]Task)

	i := 0
	for key := range tasksForUpdate {
		if err := repo.batchInsertUpdateModifiedBy(key, tasksBeforeUpdate, tasksBeforeUpdateForBatch, tasksForUpdateForBatch, tasksForUpdate, i); err != nil {
			return err
		}
		i++
	}
	return nil
}

func (repo TaskRepoCassandra) batchInsertUpdateModifiedBy(key string, tasksBeforeUpdate map[string]Task, tasksBeforeUpdateForBatch map[string]Task, tasksForUpdateForBatch map[string][]interface{}, tasksForUpdate map[string][]interface{}, i int) error {
	if task, ok := tasksBeforeUpdate[key]; ok {
		tasksBeforeUpdateForBatch[key] = task
	}

	tasksForUpdateForBatch[key] = tasksForUpdate[key]
	// no more than 30 in 1 batch or all if it's the last iteration
	if (i+1)%config.Config.CassandraBatchSize == 0 || i+1 == len(tasksForUpdateForBatch) {
		if err := repo.mutateMVTables(tasksForUpdateForBatch, tasksBeforeUpdateForBatch); err != nil {
			return err
		}
		tasksForUpdateForBatch = make(map[string][]interface{})
		tasksBeforeUpdateForBatch = make(map[string]Task)
	}
	return nil
}

// GetTargetTypeByEndpoint  returns task by endpoint and target type
func (TaskRepoCassandra) GetTargetTypeByEndpoint(partnerID string, taskID, endpointID gocql.UUID, external bool) (TargetType, error) {
	var (
		targetType TargetType
		query      = `SELECT target_type FROM tasks WHERE id = ? AND managed_endpoint_id = ? and partner_id = ? and external_task = ?`
	)

	err := cassandra.QueryCassandra(context.Background(), query, taskID, endpointID, partnerID, external).
		Scan(&targetType)
	return targetType, err
}

// GetByPartner returns Tasks found by taskID
func (TaskRepoCassandra) GetByPartner(ctx context.Context, partnerID string) ([]Task, error) {
	selectStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                       run_time_unix, run_time, state, target_type
                   FROM tasks
                   WHERE partner_id = ? AND external_task = false`
	return selectTasks(ctx, selectStmt, partnerID)
}

// GetByIDs returns Tasks found by taskID
func (repo TaskRepoCassandra) GetByIDs(ctx context.Context, cache persistency.Cache, partnerID string, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]Task, error) {
	selectStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                       run_time_unix, run_time, state, target_type
                   FROM tasks_by_id_mv
                   WHERE partner_id = ? AND id = ?`

	if isCommonFieldsNeededOnly {
		selectStmt += " LIMIT 1"
	}

	tasks := make([]Task, 0, len(taskIDs))
	results := make(chan []Task, config.Config.CassandraConcurrentCallNumber)
	limit := make(chan struct{}, config.Config.CassandraConcurrentCallNumber)
	done := make(chan struct{})
	var resultErr error

	go func() {
		for task := range results {
			tasks = append(tasks, task...)
		}
		done <- struct{}{}
	}()

	var wg sync.WaitGroup

	for _, tID := range taskIDs {
		taskID := tID
		limit <- struct{}{}
		wg.Add(1)

		go func(taskID gocql.UUID) {
			defer func() {
				<-limit
				wg.Done()
			}()

			if config.Config.AssetCacheEnabled && cache != nil && isCommonFieldsNeededOnly {
				err := repo.getFromCache(ctx, taskID, cache, selectStmt, partnerID, results)
				if err != nil {
					resultErr = fmt.Errorf("%v : %v", resultErr, err)
					return
				}
				return
			}

			tasks, err := selectTasks(ctx, selectStmt, partnerID, taskID)
			if err != nil {
				resultErr = fmt.Errorf("%v : %v", resultErr, err)
				return
			}

			if len(tasks) > 0 {
				results <- tasks
				return
			}
			logger.Log.InfofCtx(ctx, "%v : instance with id [%s] not found", resultErr, taskID.String())
		}(taskID)
	}

	wg.Wait()
	close(results)
	<-done
	return tasks, resultErr
}

func (repo TaskRepoCassandra) getFromCache(ctx context.Context, taskID gocql.UUID, cache persistency.Cache, selectStmt string, partnerID string, results chan []Task) error {
	keyForCache := []byte("TKS_TASKS_BY_ID_" + taskID.String())

	taskBin, err := cache.Get(keyForCache)
	task := Task{}
	if err != nil || json.Unmarshal(taskBin, &task) != nil {
		tasksFromCassandra, err := selectTasks(ctx, selectStmt, partnerID, taskID)
		if err != nil {
			return err
		}

		if len(tasksFromCassandra) == 0 {
			return nil
		}

		results <- tasksFromCassandra
		for _, task := range tasksFromCassandra {
			taskBytes, err := json.Marshal(task)
			if err != nil {
				logger.Log.WarnfCtx(ctx,"couldn't marshal  task  for partner with partnerID=%s, err:%s", partnerID, err.Error())
				continue
			}

			err = cache.Set(keyForCache, taskBytes, 0)
			if err != nil {
				logger.Log.WarnfCtx(ctx,"couldn't set task  with id=%s, err:", partnerID, err.Error())
			}
		}
		return nil
	}

	results <- []Task{task}
	return nil
}

// GetByIDAndManagedEndpoints returns Task founded by taskID and Targets for specific partner
func (TaskRepoCassandra) GetByIDAndManagedEndpoints(ctx context.Context, partnerID string, taskID gocql.UUID, managedEndpointIDs ...gocql.UUID) ([]Task, error) {
	tasks := make([]Task, 0)
	selectInternalTasksStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                                    run_time_unix, run_time, state, target_type
                                FROM tasks
                                WHERE partner_id = ? AND id = ? AND managed_endpoint_id = ? AND external_task = false`
	selectExternalTasksStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                                    run_time_unix, run_time, state, target_type
                                FROM tasks
                                WHERE partner_id = ? AND id = ? AND managed_endpoint_id = ? AND external_task = true`
	for _, endpointID := range managedEndpointIDs {
		ts, err := selectTasks(ctx, selectInternalTasksStmt, partnerID, taskID, endpointID)

		if err != nil {
			return nil, err
		}

		if len(ts) == 0 {
			ts, err = selectTasks(ctx, selectExternalTasksStmt, partnerID, taskID, endpointID)
			if err != nil {
				return nil, err
			}

			if len(ts) == 0 {
				continue
			}

		}

		task := ts[0]
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil, TaskNotFoundError{ErrorParameters: fmt.Sprintf("TaskID: %s, PartnerID: %s, ManagedEndpoints: %v", taskID, partnerID, managedEndpointIDs)}
	}
	return tasks, nil
}

// GetByRunTimeRange returns Tasks found by next run time range truncated to minutes
func (TaskRepoCassandra) GetByRunTimeRange(ctx context.Context, runTimeRange []time.Time) ([]Task, error) {
	runTimeSlice := make([]interface{}, len(runTimeRange))
	for i, runTime := range runTimeRange {
		runTimeSlice[i] = runTime
	}

	selectStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                        run_time_unix, run_time, state, target_type
                   FROM tasks_by_runtime_mv
                   WHERE run_time_unix IN (%s)`
	return selectTasks(ctx, fmt.Sprintf(selectStmt, common.GetQuestionMarkString(len(runTimeRange))), runTimeSlice...)
}

// GetByPartnerAndTime if a function to get tasks selected by nextRunTime and PartnerID
func (TaskRepoCassandra) GetByPartnerAndTime(ctx context.Context, partnerID string, timeToSearchFrom time.Time) ([]Task, error) {
	selectStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                       run_time_unix, run_time, state, target_type
                   FROM task_by_runtime_unix_mv
                   WHERE partner_id = ? AND external_task = false AND run_time_unix > ?`
	return selectTasks(ctx, selectStmt, partnerID, timeToSearchFrom)
}

// GetByPartnerAndManagedEndpointID returns Tasks found by partnerID and targetID
func (TaskRepoCassandra) GetByPartnerAndManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID, count int) (tasks []Task, err error) {
	selectStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                       run_time_unix, run_time, state, target_type
                   FROM tasks_order_by_last_task_instance_id_mv
                   WHERE partner_id = ? AND managed_endpoint_id = ? AND external_task = false`
	if count != common.UnlimitedCount {
		selectStmt += fmt.Sprintf(" LIMIT %v", count)
	}
	return selectTasks(ctx, selectStmt, partnerID, managedEndpointID)
}

// GetManagedEndpointIDsOfActiveTasks retrieve  managed_endpoint_id of all active tasks by taskID
func (TaskRepoCassandra) GetManagedEndpointIDsOfActiveTasks(ctx context.Context, partnerID string, taskID gocql.UUID) (map[gocql.UUID]struct{}, error) {
	var (
		selectStmt         = `SELECT managed_endpoint_id, state FROM tasks_by_id_mv WHERE partner_id = ? AND id = ?`
		query              = cassandra.QueryCassandra(ctx, selectStmt, partnerID, taskID)
		iter               = query.Iter()
		managedEndpointID  gocql.UUID
		state              statuses.TaskState
		managedEndpointIDs = make(map[gocql.UUID]struct{})
	)

	for iter.Scan(&managedEndpointID, &state) {
		if state != statuses.TaskStateInactive {
			managedEndpointIDs[managedEndpointID] = struct{}{}
		}
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return managedEndpointIDs, nil
}

// GetCountByManagedEndpointID returns TaskCount by ManagedEndpointID ID
func (TaskRepoCassandra) GetCountByManagedEndpointID(ctx context.Context, partnerID string, managedEndpointID gocql.UUID) (tasksCount TaskCount, err error) {
	tasksCount.ManagedEndpointID = managedEndpointID
	selectStmt := `SELECT COUNT(*) FROM tasks WHERE partner_id = ? AND managed_endpoint_id = ? AND external_task = false`
	cassandraQuery := cassandra.QueryCassandra(ctx, selectStmt, partnerID, managedEndpointID)
	if err := cassandraQuery.Scan(&tasksCount.Count); err != nil {
		return tasksCount, fmt.Errorf("TaskRepoCassandra.Count: Error while trying to retrive count for "+
			"(partner ID=%s target ID=%v): %v", partnerID, managedEndpointID, err)
	}
	return
}

// GetCountsByPartner returns a map with task count of each target for specific partner
func (TaskRepoCassandra) GetCountsByPartner(ctx context.Context, partnerID string) ([]TaskCount, error) {
	var (
		selectStmt        = `SELECT managed_endpoint_id FROM tasks WHERE partner_id = ? AND external_task = false`
		cassandraQuery    = cassandra.QueryCassandra(ctx, selectStmt, partnerID)
		iter              = cassandraQuery.Iter()
		taskCountMap      = make(map[gocql.UUID]int)
		managedEndpointID gocql.UUID
	)

	for iter.Scan(&managedEndpointID) {
		taskCountMap[managedEndpointID]++
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	tasksCount := ConvertMapToTaskCountArray(taskCountMap)
	return tasksCount, nil
}

// GetByLastTaskInstanceIDs returns map of Tasks found by Last TaskInstances
func (TaskRepoCassandra) GetByLastTaskInstanceIDs(
	ctx context.Context,
	partnerID string,
	endpointID gocql.UUID,
	lastTaskInstanceIDs ...gocql.UUID,
) (map[gocql.UUID]Task, error) {

	valuesInterface := make([]interface{}, 0, len(lastTaskInstanceIDs)+2) // for next 2 lines
	valuesInterface = append(valuesInterface, partnerID)
	valuesInterface = append(valuesInterface, endpointID)
	selectStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                       run_time_unix, run_time, state
                   FROM tasks_order_by_last_task_instance_id_mv
                   WHERE partner_id = ? AND external_task = false AND managed_endpoint_id = ? AND last_task_instance_id IN (%s)`

	for _, taskID := range lastTaskInstanceIDs {
		valuesInterface = append(valuesInterface, taskID)
	}

	return selectTasksMapWithoutTargets(
		ctx,
		fmt.Sprintf(selectStmt, common.GetQuestionMarkString(len(lastTaskInstanceIDs))),
		valuesInterface...)
}

// GetTasksFilteredWithTime returns map of Tasks found by Last TaskInstances
func (TaskRepoCassandra) GetTasksFilteredWithTime(
	ctx context.Context,
	partnerID string,
	taskID gocql.UUID,
	excludedEndpoints map[gocql.UUID]Task,
	currentPoint time.Time,
) (map[gocql.UUID][]Task, error) {
	selectStmt := `SELECT partner_id, external_task, managed_endpoint_id, id, last_task_instance_id, original_next_run_time,
                       run_time_unix, run_time, state, target_type
                   FROM tasks_by_id_mv
                   WHERE partner_id = ? AND id = ?`
	tasks, err := selectTasks(ctx, selectStmt, partnerID, taskID)
	if err != nil {
		return nil, err
	}

	tasksGroupByIDMap := make(map[gocql.UUID][]Task)

	for _, task := range tasks {
		if !task.RunTimeUTC.After(currentPoint) {
			// skip previous tasks
			continue
		}
		if _, ok := excludedEndpoints[task.ManagedEndpointID]; ok {
			// skip currently running tasks
			continue
		}
		tasksGroupByIDMap[task.ID] = append(tasksGroupByIDMap[task.ID], task)
	}

	return tasksGroupByIDMap, nil
}

// ConvertMapToTaskCountArray converts map[targetID]count into a slice of TaskCount objects
func ConvertMapToTaskCountArray(taskCountMap map[gocql.UUID]int) []TaskCount {
	var (
		tasksCount = make([]TaskCount, len(taskCountMap))
		i          = 0
	)
	for managedEndpointID, count := range taskCountMap {
		tasksCount[i].ManagedEndpointID = managedEndpointID
		tasksCount[i].Count = count
		i++
	}
	return tasksCount
}

// Returns set of found rows. All select queries should use this function.
// selectQuery is the select template in which the values are inserted
func selectTasks(ctx context.Context, selectQuery string, values ...interface{}) ([]Task, error) {
	var target Task
	tasksByID := make(map[string][]Task)
	q := cassandra.QueryCassandra(ctx, selectQuery, values...)
	iter := q.Iter()

	params := []interface{}{
		&target.PartnerID,
		&target.ExternalTask,
		&target.ManagedEndpointID,
		&target.ID,
		&target.LastTaskInstanceID,
		&target.OriginalNextRunTime,
		&target.RunTimeUTC,
		&target.PostponedRunTime,
		&target.State,
		&target.TargetType,
	}

	for iter.Scan(params...) {
		task := target
		tasks, ok := tasksByID[task.ID.String()]
		if !ok {
			commonData, err := selectTaskCommonData(ctx, task.PartnerID, task.ID)
			if err != nil {
				return nil, err
			}
			task = setCommonFieldsToTask(task, commonData)
		} else {
			commonData := tasks[0]
			task = setCommonFieldsToTask(task, commonData)
		}
		tasks = append(tasks, task)
		tasksByID[task.ID.String()] = tasks
	}

	tasks := make([]Task, 0)
	for _, tasksGroup := range tasksByID {
		tasks = append(tasks, tasksGroup...)
	}

	return tasks, nil
}

func selectTaskCommonData(ctx context.Context, partnerID string, taskID gocql.UUID) (Task, error) {
	var task Task
	selectStmt := `SELECT
                       partner_id, external_task, id, created_at, created_by, credentials, definition_id, description,
                       modified_at, modified_by, name, origin_id, parameters, require_noc_access, result_webhook, schedule,
                       targets, type, resources_type
                   FROM tasks2
                   WHERE partner_id = ? AND id = ?`

	var scheduleString string
	parameters := []interface{}{
		&task.PartnerID,
		&task.ExternalTask,
		&task.ID,
		&task.CreatedAt,
		&task.CreatedBy,
		&task.Credentials,
		&task.DefinitionID,
		&task.Description,
		&task.ModifiedAt,
		&task.ModifiedBy,
		&task.Name,
		&task.OriginID,
		&task.Parameters,
		&task.IsRequireNOCAccess,
		&task.ResultWebhook,
		&scheduleString,
		&task.TargetsByType,
		&task.Type,
		&task.ResourceType,
	}

	q := cassandra.QueryCassandra(ctx, selectStmt, partnerID, taskID)
	err := q.Scan(parameters...)
	if err != nil {
		return Task{}, err
	}

	err = json.Unmarshal([]byte(scheduleString), &task.Schedule)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func setCommonFieldsToTask(task, common Task) Task {
	task.CreatedAt = common.CreatedAt
	task.CreatedBy = common.CreatedBy
	task.Credentials = common.Credentials
	task.DefinitionID = common.DefinitionID
	task.Description = common.Description
	task.ModifiedAt = common.ModifiedAt
	task.ModifiedBy = common.ModifiedBy
	task.Name = common.Name
	task.OriginID = common.OriginID
	task.Parameters = common.Parameters
	task.IsRequireNOCAccess = common.IsRequireNOCAccess
	task.ResultWebhook = common.ResultWebhook
	task.Schedule = common.Schedule
	task.TargetsByType = common.TargetsByType
	task.Type = common.Type
	task.ResourceType = common.ResourceType
	return task
}

func selectTasksMapWithoutTargets(ctx context.Context, selectQuery string, values ...interface{}) (map[gocql.UUID]Task, error) {
	var (
		cassandraQuery         = cassandra.QueryCassandra(ctx, selectQuery, values...)
		tasksByIDMap           = make(map[gocql.UUID]Task)
		iterator               = cassandraQuery.Iter()
		task                   Task
		scheduleString         string
		selectCommonFieldsStmt = `SELECT
                                      partner_id, external_task, id, created_at, created_by, credentials, definition_id, description,
                                      modified_at, modified_by, name, origin_id, parameters, require_noc_access, result_webhook, schedule, type
                                  FROM tasks2
                                  WHERE partner_id = ? AND id = ?`
	)

	for iterator.Scan(
		&task.PartnerID,
		&task.ExternalTask,
		&task.ManagedEndpointID,
		&task.ID,
		&task.LastTaskInstanceID,
		&task.OriginalNextRunTime,
		&task.RunTimeUTC,
		&task.PostponedRunTime,
		&task.State,
	) {
		q := cassandra.QueryCassandra(ctx, selectCommonFieldsStmt, task.PartnerID, task.ID)
		err := q.Scan(
			&task.PartnerID,
			&task.ExternalTask,
			&task.ID,
			&task.CreatedAt,
			&task.CreatedBy,
			&task.Credentials,
			&task.DefinitionID,
			&task.Description,
			&task.ModifiedAt,
			&task.ModifiedBy,
			&task.Name,
			&task.OriginID,
			&task.Parameters,
			&task.IsRequireNOCAccess,
			&task.ResultWebhook,
			&scheduleString,
			&task.Type,
		)
		if err != nil {
			return nil, err
		}

		schedule := apiModels.Schedule{}
		err = json.Unmarshal([]byte(scheduleString), &schedule)
		if err != nil {
			logger.Log.ErrfCtx(ctx, "wrong format of schedule string %s. Err=%s", scheduleString, err.Error())
			continue
		}

		tasksByIDMap[task.ID] = Task{
			ID:                  task.ID,
			Name:                task.Name,
			DefinitionID:        task.DefinitionID,
			Description:         task.Description,
			CreatedAt:           task.CreatedAt,
			CreatedBy:           task.CreatedBy,
			Credentials:         task.Credentials,
			PartnerID:           task.PartnerID,
			OriginID:            task.OriginID,
			State:               task.State,
			RunTimeUTC:          task.RunTimeUTC,
			OriginalNextRunTime: task.OriginalNextRunTime,
			PostponedRunTime:    task.PostponedRunTime,
			Type:                task.Type,
			Parameters:          task.Parameters,
			ExternalTask:        task.ExternalTask,
			ResultWebhook:       task.ResultWebhook,
			ManagedEndpointID:   task.ManagedEndpointID,
			LastTaskInstanceID:  task.LastTaskInstanceID,
			Schedule:            schedule,
			IsRequireNOCAccess:  task.IsRequireNOCAccess,
			ModifiedBy:          task.ModifiedBy,
			ModifiedAt:          task.ModifiedAt,
		}
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf(errWorkingWithEntities, err)
	}

	return tasksByIDMap, nil
}

// UpdateTaskByNextRunTime updates task.RunTimeInSeconds by next cron time for Recurrent tasks
func UpdateTaskByNextRunTime(ctx context.Context, task *Task) (bool, error) {
	if !task.IsRecurringlyScheduled() {
		return false, nil
	}
	// Get time in location to dial with daylight saving time
	location, err := GetLocation(ctx, task.PartnerID, task.Schedule.Location, task.ManagedEndpointID)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "scheduler.UpdateTaskByNextRunTime: can't load location for task %v. Err=%v", task.ID, err)
		location = time.UTC
	}

	runTimeLoc := task.RunTimeUTC.In(location)
	nextRunTimeLoc, newSchedule, err := common.GetNextRunTime(runTimeLoc, task.Schedule)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantRecalculate, "scheduler.UpdateTaskByNextRunTime: can't get next run time for task %v. Err=%v", task.ID, err)
		return false, err
	}

	emptyTime := time.Time{}
	currentRunTimeUTC := task.RunTimeUTC

	task.RunTimeUTC = nextRunTimeLoc.UTC()
	// postpone logic here
	if task.HasOriginalRunTime() {
		if currentRunTimeUTC != task.OriginalNextRunTime {
			task.RunTimeUTC = task.OriginalNextRunTime
		}
		task.OriginalNextRunTime = emptyTime
		task.Schedule = newSchedule
		return true, err
	}

	if task.HasPostponedTime() &&
		(nextRunTimeLoc.UTC().After(task.PostponedRunTime) ||
			nextRunTimeLoc.UTC() == task.PostponedRunTime) {
		task.RunTimeUTC = task.PostponedRunTime
		task.PostponedRunTime = emptyTime
		task.OriginalNextRunTime = nextRunTimeLoc.UTC()
		// we don't update schedule here, we'll update it in next pickup
		return true, err
	}

	task.Schedule = newSchedule
	return true, err
}

// GetLocation return the location for internal Task
func GetLocation(ctx context.Context, partnerID, locationName string, endpointID gocql.UUID) (location *time.Location, err error) {
	if len(locationName) > 0 {
		return time.LoadLocation(locationName)
	}

	return asset.ServiceInstance.GetLocationByEndpointID(ctx, partnerID, endpointID)
}

// GetGlobalTaskState calculates global task state based on task state per target
// https://continuum.atlassian.net/wiki/spaces/C2E/pages/301400183/Task+Global+and+Individual+Statuses
func GetGlobalTaskState(tasksInternal []Task) statuses.TaskState {
	containsInactive := false
	for _, internalTask := range tasksInternal {
		switch internalTask.State {
		case statuses.TaskStateActive:
			// If at least one internalTask is active, then global TaskState is active
			return statuses.TaskStateActive
		case statuses.TaskStateInactive:
			containsInactive = true
		}
	}

	// If no active internalTask left, then the task is Inactive or Disabled
	// If at least one internalTask is inactive and no active tasks left, then global TaskState is inactive
	if containsInactive {
		return statuses.TaskStateInactive
	}

	// If there is no inactive and active internalTask, then global TaskState is disabled
	return statuses.TaskStateDisabled
}

// NewTaskOutput builds Output Task based on a group of internal Tasks
func NewTaskOutput(ctx context.Context, internalTaskGroup []Task) (*Task, error) {
	if len(internalTaskGroup) == 0 {
		return nil, errors.New("NewTaskOutput can't build a task from an empty list of internal tasks")
	}

	var (
		managedEndpointsDetailed = make([]ManagedEndpointDetailed, len(internalTaskGroup))
		tasksByModifiedBy        = make(map[time.Time][]Task)
		latestUpdateTime         time.Time
	)

	// filter task by last updated one's
	for _, task := range internalTaskGroup {
		updateTime := task.ModifiedAt
		if latestUpdateTime.IsZero() {
			latestUpdateTime = updateTime
		}

		tasksByModifiedBy[updateTime] = append(tasksByModifiedBy[updateTime], task)
		if updateTime.After(latestUpdateTime) {
			latestUpdateTime = updateTime
		}
	}

	// user will get only the last updated one's
	internalTaskGroup = tasksByModifiedBy[latestUpdateTime]
	for i, task := range internalTaskGroup {
		managedEndpointsDetailed[i] = ManagedEndpointDetailed{
			ManagedEndpoint: apiModels.ManagedEndpoint{
				ID:          task.ManagedEndpointID.String(),
				NextRunTime: task.RunTimeUTC.UTC(),
			},
			State: task.State,
		}
	}

	globalTaskState := GetGlobalTaskState(internalTaskGroup)
	taskTemplate := internalTaskGroup[0]
	for k, v := range taskTemplate.TargetsByType {
		taskTemplate.Targets.IDs = v
		taskTemplate.Targets.Type = k
		break
	}

	taskOutput := &Task{
		ID:                  taskTemplate.ID,
		Name:                taskTemplate.Name,
		Description:         taskTemplate.Description,
		CreatedAt:           taskTemplate.CreatedAt,
		CreatedBy:           taskTemplate.CreatedBy,
		PartnerID:           taskTemplate.PartnerID,
		OriginID:            taskTemplate.OriginID,
		State:               globalTaskState,
		RunTimeUTC:          taskTemplate.RunTimeUTC,
		Type:                taskTemplate.Type,
		Targets:             taskTemplate.Targets,
		Parameters:          taskTemplate.Parameters,
		ExternalTask:        taskTemplate.ExternalTask,
		ResultWebhook:       taskTemplate.ResultWebhook,
		ManagedEndpoints:    managedEndpointsDetailed,
		Schedule:            taskTemplate.Schedule,
		DefinitionID:        taskTemplate.DefinitionID,
		ModifiedAt:          latestUpdateTime,
		ModifiedBy:          taskTemplate.ModifiedBy,
		OriginalNextRunTime: taskTemplate.OriginalNextRunTime,
		ResourceType:        taskTemplate.ResourceType,
		TargetsByType:       taskTemplate.TargetsByType,
	}

	if taskTemplate.Schedule.Location == "" {
		return taskOutput, nil
	}
	// Get time in location to represent data as the user has specified it
	location, err := GetLocation(ctx, taskTemplate.PartnerID, taskTemplate.Schedule.Location, taskTemplate.ManagedEndpointID)
	if err != nil {
		// If err then return time fields in UTC (when location = "")
		logger.Log.WarnfCtx(ctx, "NewTaskOutput: can't load location for task %v. Err=%v", taskTemplate.ID, err)
		return taskOutput, nil
	}

	taskOutput.Schedule.StartRunTime = taskTemplate.Schedule.StartRunTime.In(location)
	if taskOutput.Schedule.EndRunTime.UTC().IsZero() { // RMM-48594
		taskOutput.Schedule.EndRunTime = taskTemplate.Schedule.EndRunTime.UTC()
		return taskOutput, nil
	}

	taskOutput.Schedule.EndRunTime = taskTemplate.Schedule.EndRunTime.In(location)
	return taskOutput, nil
}

func convertSelectedManagedEndpointEnableToListOfManagedEndpointIDs(v *SelectedManagedEndpointEnable) ([]gocql.UUID, error) {
	var (
		managedEndpointIDs = make([]gocql.UUID, len(v.ManagedEndpoints))
		i                  = 0
	)
	for target := range v.ManagedEndpoints {
		managedEndpointID, err := gocql.ParseUUID(target)
		if err != nil {
			return nil, err
		}
		managedEndpointIDs[i] = managedEndpointID
		i++
	}
	return managedEndpointIDs, nil
}
