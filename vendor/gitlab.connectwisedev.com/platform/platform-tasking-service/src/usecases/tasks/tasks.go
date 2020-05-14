package tasks

//go:generate mockgen -destination=./tasks_repo_mock_test.go  -package=tasks -source=./tasks_mock.go

import (
	"context"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/types"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
	"golang.org/x/sync/semaphore"
)

const (
	contextParameterError = "can't get parameter from context: %s"
)

//ExecutionResultsRepo - script execution results repository interface
type ExecutionResultsRepo interface {
	GetLastResultByEndpointID(endpointID string) (res entities.ExecutionResult, err error)
	GetLastExecutions(partnerID string, endpointIDs map[string]struct{}) (res []entities.LastExecution, err error)
}

//Repo - tasks repository interface
type Repo interface {
	GetName(partnerID string, id string) (string, error)
	GetNext(partnerID string) ([]entities.Task, error)
	GetScheduledTasks(partnerID string) ([]entities.ScheduledTasks, error)
}

//LegacyRepo - represents legacy task repository
type LegacyRepo interface {
	GetByIDs(ctx context.Context, cache persistency.Cache, partnerID string, isCommonFieldsNeededOnly bool, taskIDs ...gocql.UUID) ([]models.Task, error)
	UpdateSchedulerFields(ctx context.Context, tasks ...models.Task) error
}

//InstancesRepo - task instances repository interface
type InstancesRepo interface {
	GetInstancesForScheduled(IDs []string) ([]entities.TaskInstance, error)
	GetTopInstancesForScheduledByTaskIDs(taskIDs []string) ([]entities.TaskInstance, error)
	GetMinimalInstanceByID(id string) (entities.TaskInstance, error)
	GetByStartedAtAfter(partnerID string, from, to time.Time) ([]entities.TaskInstance, error)
}

//NewTasks - returns new Tasks use case instance
func NewTasks(tasksRepo Repo, legacyRepo LegacyRepo, taskInstanceRepo InstancesRepo, executionResultsRepo ExecutionResultsRepo, cache persistency.Cache, tr trigger.Usecase, l logger.Logger) *Tasks {
	return &Tasks{
		tasksRepo:        tasksRepo,
		legacyRepo:       legacyRepo,
		taskInstanceRepo: taskInstanceRepo,
		execResultsRepo:  executionResultsRepo,
		cache:            cache,
		tr:               tr,
		log:              l,
	}
}

//Tasks - Tasks use case
type Tasks struct {
	tasksRepo        Repo
	legacyRepo       LegacyRepo
	taskInstanceRepo InstancesRepo
	execResultsRepo  ExecutionResultsRepo
	cache            persistency.Cache
	tr               trigger.Usecase
	log              logger.Logger
}

type taskData struct {
	id         string
	instanceID string
	name       string
	run        time.Time
	status     statuses.TaskInstanceStatus
}

type fullTaskData struct {
	*taskData
	endpointID string
}

type instanceCache struct {
	m         sync.RWMutex
	instances map[string]entities.TaskInstance
}

func (i *instanceCache) get(key string) (val entities.TaskInstance, ok bool) {
	i.m.RLock()
	val, ok = i.instances[key]
	i.m.RUnlock()
	return
}
func (i *instanceCache) put(key string, val entities.TaskInstance) {
	i.m.Lock()
	i.instances[key] = val
	i.m.Unlock()
}

//GetClosestTasks - returns last executed and next scheduled Tasks for every endpoint of partner
func (t *Tasks) GetClosestTasks(ctx context.Context, endpointsInput entities.EndpointsInput) (tasks entities.EndpointsClosestTasks, err error) {
	tasks = entities.EndpointsClosestTasks{}

	endpoints := make(map[string]struct{})
	for _, e := range endpointsInput {
		endpoints[e] = struct{}{}
	}

	partnerID, err := t.readCtx(ctx)
	if err != nil {
		return tasks, err
	}
	prev := make(map[string]*taskData)
	multiErr := types.NewMultiError()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		tasks, err = t.mapNextRuns(partnerID, endpoints)
		if err != nil {
			multiErr = append(multiErr, err)
		}
	}()

	go func() {
		defer wg.Done()
		prev, err = t.findLastRuns(ctx, partnerID, endpoints)
		if err != nil {
			multiErr = append(multiErr, err)
		}
	}()

	wg.Wait()

	tasks, err = t.mapLastRuns(prev, tasks)
	if err != nil {
		return
	}
	return tasks, multiErr.ToError()
}

func (t *Tasks) findLastRuns(ctx context.Context, partnerID string, endpoints map[string]struct{}) (map[string]*taskData, error) {
	endpointInstances := make(map[string]*taskData)
	dataChan := make(chan fullTaskData, len(endpoints))
	multiErr := types.NewMultiError()

	unhandled, err := t.getLastExecutions(partnerID, endpoints, dataChan)
	if err != nil {
		multiErr = append(multiErr, err)
	}

	if len(unhandled) > 0 {
		err = t.retrieveLastExecutionsByResults(ctx, unhandled, partnerID, dataChan)
		if err != nil {
			multiErr = append(multiErr, err)
		}
	}
	close(dataChan)

	for data := range dataChan {
		if data.taskData != nil {
			endpointInstances[data.endpointID] = data.taskData
		}
	}
	return endpointInstances, multiErr.ToError()
}

func (t *Tasks) retrieveLastExecutionsByResults(ctx context.Context, unhandled map[string]struct{}, partnerID string, dataChan chan fullTaskData) (processingErr error) {
	instanceCache := instanceCache{instances: make(map[string]entities.TaskInstance)}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(config.Config.ClosestTasksWorkersTimeoutSec)*time.Second)
	defer cancel()
	sem := semaphore.NewWeighted(int64(config.Config.CassandraConcurrentCallNumber))
	wg := sync.WaitGroup{}
	for eid := range unhandled {
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}
		wg.Add(1)
		go func(endpointID string) {
			defer wg.Done()
			defer sem.Release(1)

			data, err := t.getFullLastTaskDataForEndpoint(partnerID, endpointID, &instanceCache)
			if err != nil {
				cancel()
				processingErr = err
				return
			}

			dataChan <- fullTaskData{
				taskData:   data,
				endpointID: endpointID,
			}
		}(eid)
	}
	wg.Wait()
	return
}

func (t *Tasks) getLastExecutions(partnerID string, endpoints map[string]struct{}, dataChan chan<- fullTaskData) (map[string]struct{}, error) {
	unhandled := make(map[string]struct{})
	for k, v := range endpoints {
		unhandled[k] = v
	}
	lastExecs, err := t.execResultsRepo.GetLastExecutions(partnerID, endpoints)

	for _, ex := range lastExecs {
		if ex.Name == "" {
			continue
		}
		delete(unhandled, ex.EndpointID)
		dataChan <- fullTaskData{
			taskData: &taskData{
				name:   ex.Name,
				run:    ex.RunTime,
				status: ex.Status,
			},
			endpointID: ex.EndpointID,
		}
	}

	return unhandled, err
}

func (t *Tasks) getFullLastTaskDataForEndpoint(partnerID, endpointID string, instances *instanceCache) (data *taskData, err error) {
	var lastRes entities.ExecutionResult
	lastRes, err = t.execResultsRepo.GetLastResultByEndpointID(endpointID)
	if err != nil {
		return
	}
	if lastRes.ExecutionStatus == 0 {
		return
	}

	ti, ok := instances.get(lastRes.TaskInstanceID)
	if !ok {
		ti, err = t.taskInstanceRepo.GetMinimalInstanceByID(lastRes.TaskInstanceID)
		if err != nil {
			return
		}
		if ti.TaskName == "" {
			var name string
			name, err = t.tasksRepo.GetName(partnerID, ti.TaskID)
			if err != nil {
				return
			}
			if name == "" {
				return
			}
			ti.TaskName = name
		}

		instances.put(lastRes.TaskInstanceID, ti)
	}

	if ti.PartnerID != partnerID {
		return
	}

	data = &taskData{
		id:         ti.TaskID,
		instanceID: lastRes.TaskInstanceID,
		name:       ti.TaskName,
		run:        lastRes.UpdatedAt,
		status:     lastRes.ExecutionStatus,
	}
	return
}

func (t *Tasks) mapNextRuns(partnerID string, endpoints map[string]struct{}) (entities.EndpointsClosestTasks, error) {
	closest := make(entities.EndpointsClosestTasks)
	next, err := t.tasksRepo.GetNext(partnerID)

	for _, task := range next {
		_, ok := endpoints[task.ManagedEndpointID]
		if !ok || task.Name == "" || task.State == statuses.TaskStateDisabled {
			continue
		}

		runDate := task.RunTimeUTC.Unix()
		if !task.PostponedRunTime.IsZero() {
			runDate = task.PostponedRunTime.Unix() // For Postponed Recurrent tasks next time will be a Postponed time
		} // and for Postponed OneTime tasks - RunTimeUTC

		if current := closest[task.ManagedEndpointID].Next; current != nil && current.RunDate < runDate {
			continue
		}

		closest[task.ManagedEndpointID] = entities.ClosestTasks{
			Next: &entities.ClosestTask{
				ID:      task.ID,
				Name:    task.Name,
				RunDate: runDate,
			},
		}
	}

	return closest, err
}

func (t *Tasks) mapLastRuns(prev map[string]*taskData, tasks entities.EndpointsClosestTasks) (entities.EndpointsClosestTasks, error) {
	for endpointID, last := range prev {
		if last == nil {
			continue
		}

		task := tasks[endpointID]

		status, err := statuses.TaskInstanceStatusText(last.status)
		if err != nil {
			return tasks, errors.Wrap(err, "invalid status")
		}

		task.Previous = &entities.ClosestTask{
			ID:      last.id,
			Name:    last.name,
			RunDate: last.run.Unix(),
			Status:  status,
		}
		tasks[endpointID] = task
	}
	return tasks, nil
}

func (t *Tasks) readCtx(ctx context.Context) (partnerID string, err error) {
	partnerID, ok := ctx.Value(config.PartnerIDKeyCTX).(string)
	if !ok {
		err = errors.Errorf(contextParameterError, "partnerID")
		return
	}
	return
}
