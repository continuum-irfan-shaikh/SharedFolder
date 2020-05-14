package tasks

import (
	"log"
	"strconv"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"github.com/gocql/gocql"
)

type resRepoMock struct {
	results map[string]entities.ExecutionResult
}

func (r *resRepoMock) GetLastResultByEndpointID(endpointID string) (res entities.ExecutionResult, err error) {
	return r.results[endpointID], nil
}

func (r *resRepoMock) GetLastExecutions(partnerID string, endpointIDs map[string]struct{}) (res []entities.LastExecution, err error) {
	return res, nil
}

type taskRepoMock struct {
	nextTasks []entities.Task
}

func (t *taskRepoMock) GetByIDs(cache persistency.Cache, partnerID string, taskIDs ...string) ([]entities.ScheduledTasks, error) {
	panic("implement me")
}

func (t *taskRepoMock) GetName(partnerID string, id string) (string, error) {
	panic("implement me")
}

func (t *taskRepoMock) GetNext(partnerID string) ([]entities.Task, error) {
	return t.nextTasks, nil
}

func (t *taskRepoMock) GetScheduledTasks(partnerID string) ([]entities.ScheduledTasks, error) {
	panic("implement me")
}

type instRepoMock struct {
	inst entities.TaskInstance
}

func (i *instRepoMock) GetByStartedAtAfter(partnerID string, from, to time.Time) ([]entities.TaskInstance, error) {
	panic("implement me")
}

func (i *instRepoMock) GetInstancesForScheduled(IDs []string) ([]entities.TaskInstance, error) {
	panic("implement me")
}

func (i *instRepoMock) GetTopInstancesForScheduledByTaskIDs(taskIDs []string) ([]entities.TaskInstance, error) {
	panic("implement me")
}

func (i *instRepoMock) GetMinimalInstanceByID(id string) (entities.TaskInstance, error) {
	return i.inst, nil
}

func generateTestData() ([]string, []entities.Task, map[string]entities.ExecutionResult) {
	endpoints := make([]string, 0)
	nextTasks := make([]entities.Task, 0)
	results := make(map[string]entities.ExecutionResult)
	for i := 0; i < 1000; i++ {
		uuid, err := gocql.RandomUUID()
		if err != nil {
			log.Fatal(err)
		}
		eid := uuid.String()
		endpoints = append(endpoints, eid)

		nextTask := entities.Task{
			Name:              "TEST_TASK_" + strconv.Itoa(i),
			RunTimeUTC:        time.Now().Add(24 * time.Hour),
			ManagedEndpointID: eid,
			State:             statuses.TaskStateActive,
		}
		nextTasks = append(nextTasks, nextTask)

		results[eid] = entities.ExecutionResult{
			ManagedEndpointID: eid,
			TaskInstanceID:    eid,
			ExecutionStatus:   statuses.TaskInstanceSuccess,
			UpdatedAt:         time.Unix(int64(i), 0),
		}
	}
	return endpoints, nextTasks, results
}
