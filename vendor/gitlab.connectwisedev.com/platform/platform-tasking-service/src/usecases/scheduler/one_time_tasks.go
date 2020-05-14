package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	s "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const expirationTime = 24 * 60 * 60 // hours * minutes * seconds

//OneTimeTasksBuilder - builder for OneTimeTasks
type OneTimeTasksBuilder struct {
	Conf                    config.Configuration
	TaskExecutionRepo       TaskExecutionRepo
	TaskInstanceRepo        TaskInstanceRepo
	ExecutionResultRepo     ExecutionResultRepo
	TaskRepo                TaskRepo
	ExecutionExpirationRepo ExecutionExpirationRepo
	CacheRepo               CacheRepo
	Log                     logger.Logger
	EncryptionService       EncryptionService
	AgentEncryptionService  AgentEncryptionService
}

//Build - initialize OneTimeTasks struct
func (o *OneTimeTasksBuilder) Build() *OneTimeTasks {
	return &OneTimeTasks{
		conf:                    o.Conf,
		taskExecutionRepo:       o.TaskExecutionRepo,
		taskInstanceRepo:        o.TaskInstanceRepo,
		executionResultRepo:     o.ExecutionResultRepo,
		taskRepo:                o.TaskRepo,
		executionExpirationRepo: o.ExecutionExpirationRepo,
		cacheRepo:               o.CacheRepo,
		log:                     o.Log,
		encryptionService:       o.EncryptionService,
		agentEncryptionService:  o.AgentEncryptionService,
	}
}

//OneTimeTasks - represents one time task processing
type OneTimeTasks struct {
	conf                    config.Configuration
	taskExecutionRepo       TaskExecutionRepo
	taskInstanceRepo        TaskInstanceRepo
	executionResultRepo     ExecutionResultRepo
	taskRepo                TaskRepo
	executionExpirationRepo ExecutionExpirationRepo
	cacheRepo               CacheRepo
	log                     logger.Logger
	encryptionService       EncryptionService
	agentEncryptionService  AgentEncryptionService
}

//Process - handle one time tasks
func (o *OneTimeTasks) Process(ctx context.Context, currentTime time.Time, tasks []m.Task) {
	groupedTasks := o.groupByTaskInstanceId(tasks)

	for tiID, tasks := range groupedTasks {
		ti, err := o.taskInstanceRepo.GetInstance(tiID)
		if err != nil {
			o.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskInstances, err.Error())
			o.recalculateRunTime(tasks...)
			continue
		}

		go o.processByTaskInstanceId(ctx, currentTime, ti, tasks)
	}
}

func (o *OneTimeTasks) processByTaskInstanceId(ctx context.Context, time time.Time, ti m.TaskInstance, tasks []m.Task) {
	tasksForExecution := make([]m.Task, 0)
	tasksForDeactivation := make([]m.Task, 0)
	for _, task := range tasks {
		if o.shouldBeRun(time, ti, task) {
			tasksForExecution = append(tasksForExecution, task)
		}
		tasksForDeactivation = append(tasksForDeactivation, task)
	}
	// this is needed for case with editing
	if len(tasksForExecution) != 0 {
		ti = o.updateTaskInstance(ctx, time, ti, tasksForExecution...)
		go o.executeTasks(ctx, time, ti, tasksForExecution...)
	}
	go o.deactivateTasks(tasksForDeactivation...)
}

//executeTasks - execute tasks by task instance id
func (o *OneTimeTasks) executeTasks(ctx context.Context, currentRunTime time.Time, ti m.TaskInstance, tasks ...m.Task) {
	var (
		te             []m.Task //Tasks for execution
		meForExecution []apiModels.ManagedEndpoint
		taskTemplate   m.Task
	)

	tasksByManagedEndpointID := make(map[string]m.Task)

	for _, task := range tasks {
		te = append(te, task)
		me := apiModels.ManagedEndpoint{
			ID:          task.ManagedEndpointID.String(),
			NextRunTime: task.RunTimeUTC.UTC(),
		}
		meForExecution = append(meForExecution, me)
		tasksByManagedEndpointID[task.ManagedEndpointID.String()] = task
	}

	if len(te) == 0 {
		return //nothing to execute
	}

	taskTemplate = te[0]

	pattern := "%s/partners/%s/task-execution-results/task-instances/%s"
	webHookURL := fmt.Sprintf(pattern, o.conf.TaskingMsURL, taskTemplate.PartnerID, ti.ID)
	payload := apiModels.ExecutionPayload{
		ExecutionID:              ti.ID.String(),
		OriginID:                 ti.OriginID.String(),
		ManagedEndpoints:         meForExecution,
		Parameters:               taskTemplate.Parameters, //Common parameters for all tasks inside task instance
		TaskID:                   taskTemplate.ID,
		WebhookURL:               webHookURL,
		Credentials:              taskTemplate.Credentials,
		ExpectedExecutionTimeSec: o.cacheRepo.CalculateExpectedExecutionTimeSec(context.Background(), taskTemplate),
	}

	if taskTemplate.RunTimeUTC.Add(time.Second * time.Duration(payload.ExpectedExecutionTimeSec)).Before(time.Now().UTC()) {
		o.log.WarnfCtx(ctx, "executeTasks: task can't be run as it currently out of schedule, taskID: %v", taskTemplate.ID)
		return
	}

	if taskTemplate.IsRunAsUserApplied() {
		o.executeAsUser(ctx, te, meForExecution, payload, tasksByManagedEndpointID, currentRunTime, ti)
		return
	}

	if err := o.taskExecutionRepo.ExecuteTasks(ctx, payload, taskTemplate.PartnerID, taskTemplate.Type); err != nil {
		o.log.ErrfCtx(ctx, errorcode.ErrorCantExecuteTasks, fmt.Sprintf("executeTask: error during execution one time task. err : %s", err.Error()))
		if currentRunTime.Equal(taskTemplate.Schedule.EndRunTime) {
			o.sendExecutionResultsErr(ctx, currentRunTime, te...)
			return
		}
		o.recalculateRunTime(te...)
	}

	o.saveTimeExpiration(ctx, ti, payload.ExpectedExecutionTimeSec)
}

func (o *OneTimeTasks) executeAsUser(ctx context.Context, te []m.Task, meForExecution []apiModels.ManagedEndpoint, payload apiModels.ExecutionPayload, tasksByManagedEndpointID map[string]m.Task, time time.Time, ti m.TaskInstance) {
	var totalErr error
	var taskTemplate = te[0]

	decrypted, err := o.encryptionService.Decrypt(*taskTemplate.Credentials)
	if err != nil {
		o.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "executeTasks: could't decrypt credentials for task with ID %v . err: %s", taskTemplate.ID, err.Error())
		o.recalculateRunTime(te...)
		return
	}
	taskTemplate.Credentials = &decrypted

	for _, me := range meForExecution {
		totalErr = o.executeForME(ctx, payload, me, taskTemplate, tasksByManagedEndpointID, totalErr, time)
	}

	if totalErr != nil {
		o.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, totalErr.Error())
	}

	go o.saveTimeExpiration(ctx, ti, payload.ExpectedExecutionTimeSec)
	return
}

func (o *OneTimeTasks) executeForME(ctx context.Context, payload apiModels.ExecutionPayload, me apiModels.ManagedEndpoint, taskTemplate m.Task, tasksByManagedEndpointID map[string]m.Task, totalErr error, time time.Time) error {
	payload.ManagedEndpoints = []apiModels.ManagedEndpoint{me}

	if taskTemplate.Credentials != nil {
		encrypted, err := o.agentEncryptionService.Encrypt(ctx, tasksByManagedEndpointID[me.ID].ManagedEndpointID, *taskTemplate.Credentials)
		if err != nil {
			totalErr = fmt.Errorf("%v err: executeTasks: could't encryt credentials for task with ID %v and endpointID %v. err: %s", totalErr, taskTemplate.ID, me.ID, err.Error())
			o.recalculateRunTime(tasksByManagedEndpointID[me.ID])
			return totalErr
		}

		payload.Credentials = &encrypted
	}
	var err error
	if err = o.taskExecutionRepo.ExecuteTasks(ctx, payload, taskTemplate.PartnerID, taskTemplate.Type); err == nil {
		return totalErr
	}

	if time.Equal(taskTemplate.Schedule.EndRunTime) {
		o.sendExecutionResultsErr(ctx, time, tasksByManagedEndpointID[me.ID])
		return totalErr
	}

	if totalErr == nil {
		totalErr = fmt.Errorf("executeAsUser: error during execution one time tasks.")
	} else {
		totalErr = fmt.Errorf("%s err: %s", totalErr.Error(), err.Error())
	}

	o.recalculateRunTime(tasksByManagedEndpointID[me.ID])
	return totalErr
}

func (o *OneTimeTasks) deactivateTasks(tasks ...m.Task) {
	ctx := context.Background() // Deprecated dependency in repository
	for i := range tasks {
		tasks[i].State = s.TaskStateInactive
	}

	if err := o.taskRepo.UpdateSchedulerFields(ctx, tasks...); err != nil {
		o.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "", err.Error())
	}
}

func (o *OneTimeTasks) shouldBeRun(time time.Time, ti m.TaskInstance, t m.Task) bool {
	status, ok := ti.Statuses[t.ManagedEndpointID]
	if !ok {
		return false
	}

	switch status {
	case s.TaskInstanceCanceled,
		s.TaskInstanceDisabled,
		s.TaskInstanceStopped:
		return false
	}

	if t.State == s.TaskStateInactive {
		return false
	}

	if t.Schedule.EndRunTime.IsZero() {
		return true
	}

	if time.After(t.Schedule.EndRunTime) {
		return false
	}

	return true
}

func (o *OneTimeTasks) groupByTaskInstanceId(tasks []m.Task) map[gocql.UUID][]m.Task {
	groupedTasks := make(map[gocql.UUID][]m.Task, 0)
	for _, task := range tasks {
		//id - task instance ID
		id := task.LastTaskInstanceID
		if _, ok := groupedTasks[id]; !ok {
			groupedTasks[id] = make([]m.Task, 0)
		}
		groupedTasks[id] = append(groupedTasks[id], task)
	}
	return groupedTasks
}

func (o *OneTimeTasks) recalculateRunTime(tasks ...m.Task) {
	ctx := context.Background() // Deprecated dependency in repository
	for _, task := range tasks {
		addTime := time.Duration(o.conf.RecalculateTime) * time.Second
		task.RunTimeUTC.Add(addTime)
		if err := o.taskRepo.UpdateSchedulerFields(ctx, task); err != nil {
			o.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, err.Error())
		}
	}
}

func (o *OneTimeTasks) updateTaskInstance(ctx context.Context, time time.Time, ti m.TaskInstance, tasks ...m.Task) m.TaskInstance {
	statuses := make(map[gocql.UUID]s.TaskInstanceStatus)
	IDs := o.getMEIDs(tasks...)
	for id, status := range ti.Statuses {
		if _, ok := IDs[id]; ok {
			statuses[id] = s.TaskInstanceRunning
			continue
		}
		if status == s.TaskInstanceScheduled {
			statuses[id] = s.TaskInstancePending
			continue
		}
		statuses[id] = status
	}

	ti.Statuses = statuses
	ti.LastRunTime = time
	ttl := o.conf.DataRetentionIntervalDay * expirationTime
	if err := o.taskInstanceRepo.Insert(ti, ttl); err != nil {
		o.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, err.Error())
	}
	return ti
}

func (o *OneTimeTasks) sendExecutionResultsErr(ctx context.Context, time time.Time, tasks ...m.Task) {
	for _, task := range tasks {
		executionResult := tasking.ExecutionResultKafkaMessage{
			Message: tasking.ScriptPluginReturnMessage{
				ExecutionID:  task.LastTaskInstanceID.String(),
				TimestampUTC: time,
				Status:       s.TaskInstanceFailedText,
				Stderr:       "failed by time out",
			},
			BrokerEnvelope: agent.BrokerEnvelope{
				EndpointID: task.ManagedEndpointID.String(),
				PartnerID:  task.PartnerID,
			},
		}
		if err := o.executionResultRepo.Publish(executionResult); err != nil {
			o.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData,  err.Error())
		}
	}
}

//getMEIDs - return list management end point IDs
func (o *OneTimeTasks) getMEIDs(tasks ...m.Task) map[gocql.UUID]struct{} {
	IDs := make(map[gocql.UUID]struct{})
	for _, task := range tasks {
		IDs[task.ManagedEndpointID] = struct{}{}
	}
	return IDs
}

func (o *OneTimeTasks) saveTimeExpiration(ctx context.Context, ti m.TaskInstance, timeExpiration int) {
	meIDs := make([]gocql.UUID, 0)
	for meID, status := range ti.Statuses {
		if status == s.TaskInstanceRunning {
			meIDs = append(meIDs, meID)
		}
	}

	timeDuration := time.Duration(timeExpiration+o.conf.HTTPClientResultsTimeoutSec) * time.Second
	ex := entities.ExecutionExpiration{
		ExpirationTimeUTC:  time.Now().Add(timeDuration).Truncate(time.Minute),
		PartnerID:          ti.PartnerID,
		TaskInstanceID:     ti.ID,
		ManagedEndpointIDs: meIDs,
	}

	ttl := o.conf.DataRetentionIntervalDay * expirationTime
	if err := o.executionExpirationRepo.Insert(ex, ttl); err != nil {
		o.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData,  err.Error())
	}
}
