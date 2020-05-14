package tasks

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func TestFindLastRuns(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrl *gomock.Controller

	type args struct {
		endpoints                map[string]struct{}
		ctx                      func() context.Context
		cassandraCallNumberLimit int
		workersTimeout           int
	}
	type repos struct {
		execResultsRepo  func() ExecutionResultsRepo
		taskInstanceRepo func() InstancesRepo
		tasksRepo        func() Repo
	}
	tests := []struct {
		name    string
		args    args
		repos   repos
		want    map[string]*taskData
		wantErr bool
		errMsg  string
	}{
		{
			name: "successfully find last runs",
			args: args{
				endpoints: map[string]struct{}{
					"e1": {},
					"e2": {},
				},
				ctx: func() context.Context {
					return context.Background()
				},
				cassandraCallNumberLimit: 1,
				workersTimeout:           10000,
			},
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastResultByEndpointID("e2").Return(entities.ExecutionResult{
						TaskInstanceID:  "i2",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastExecutions(gomock.Any(), gomock.Any()).Return([]entities.LastExecution{}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
							TaskName:  "task1",
						}, nil)

					repo.EXPECT().GetMinimalInstanceByID("i2").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t2",
						}, nil)

					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetName("p1", "t2").
						Return("task2", nil)
					return repo
				},
			},
			want: map[string]*taskData{
				"e1": {
					id:         "t1",
					instanceID: "i1",
					name:       "task1",
					run:        time.Unix(1, 0),
					status:     statuses.TaskInstanceSuccess,
				},
				"e2": {
					id:         "t2",
					instanceID: "i2",
					name:       "task2",
					run:        time.Unix(1, 0),
					status:     statuses.TaskInstanceSuccess,
				},
			},
		},
		{
			name: "successfully find last runs (from last executions table)",
			args: args{
				endpoints: map[string]struct{}{
					"e1": {},
					"e2": {},
				},
				ctx: func() context.Context {
					return context.Background()
				},
				cassandraCallNumberLimit: 1,
				workersTimeout:           10000,
			},
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastExecutions(gomock.Any(), map[string]struct{}{"e1": {}, "e2": {}}).
						Return([]entities.LastExecution{
							{
								EndpointID: "e1",
								RunTime:    time.Unix(1, 0),
								Name:       "task1",
								Status:     statuses.TaskInstanceSuccess,
							},
							{
								EndpointID: "e2",
								RunTime:    time.Unix(1, 0),
								Name:       "task2",
								Status:     statuses.TaskInstanceSuccess,
							},
						}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			want: map[string]*taskData{
				"e1": {
					name:   "task1",
					run:    time.Unix(1, 0),
					status: statuses.TaskInstanceSuccess,
				},
				"e2": {
					name:   "task2",
					run:    time.Unix(1, 0),
					status: statuses.TaskInstanceSuccess,
				},
			},
		},
		{
			name: "successfully find last runs (not all last executions found)",
			args: args{
				endpoints: map[string]struct{}{
					"e1": {},
					"e2": {},
				},
				ctx: func() context.Context {
					return context.Background()
				},
				cassandraCallNumberLimit: 1,
				workersTimeout:           10000,
			},
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e2").Return(entities.ExecutionResult{
						TaskInstanceID:  "i2",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastExecutions(gomock.Any(), map[string]struct{}{"e1": {}, "e2": {}}).
						Return([]entities.LastExecution{
							{
								EndpointID: "e1",
								RunTime:    time.Unix(1, 0),
								Name:       "task1",
								Status:     statuses.TaskInstanceSuccess,
							},
						}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i2").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskName:  "task2",
							TaskID:    "t2",
						}, nil)
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			want: map[string]*taskData{
				"e1": {
					name:   "task1",
					run:    time.Unix(1, 0),
					status: statuses.TaskInstanceSuccess,
				},
				"e2": {
					id:         "t2",
					instanceID: "i2",
					name:       "task2",
					run:        time.Unix(1, 0),
					status:     statuses.TaskInstanceSuccess,
				},
			},
		},
		{
			name: "successfully find last runs (last execution without name)",
			args: args{
				endpoints: map[string]struct{}{
					"e1": {},
					"e2": {},
				},
				ctx: func() context.Context {
					return context.Background()
				},
				cassandraCallNumberLimit: 1,
				workersTimeout:           10000,
			},
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastResultByEndpointID("e2").Return(entities.ExecutionResult{
						TaskInstanceID:  "i2",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastExecutions(gomock.Any(), map[string]struct{}{"e1": {}, "e2": {}}).
						Return([]entities.LastExecution{
							{
								EndpointID: "e1",
								RunTime:    time.Unix(1, 0),
								Status:     statuses.TaskInstanceSuccess,
							},
						}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
							TaskName:  "task1",
						}, nil)

					repo.EXPECT().GetMinimalInstanceByID("i2").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskName:  "task2",
							TaskID:    "t2",
						}, nil)
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			want: map[string]*taskData{
				"e1": {
					id:         "t1",
					instanceID: "i1",
					name:       "task1",
					run:        time.Unix(1, 0),
					status:     statuses.TaskInstanceSuccess,
				},
				"e2": {
					id:         "t2",
					instanceID: "i2",
					name:       "task2",
					run:        time.Unix(1, 0),
					status:     statuses.TaskInstanceSuccess,
				},
			},
		},
		{
			name: "failed while results retrieving",
			args: args{
				endpoints: map[string]struct{}{
					"e1": {},
					"e2": {},
				},
				ctx: func() context.Context {
					return context.Background()
				},
				cassandraCallNumberLimit: 1,
				workersTimeout:           10000,
			},
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID(gomock.Any()).
						Return(entities.ExecutionResult{}, errors.New("fail")).
						AnyTimes()
					repo.EXPECT().GetLastExecutions(gomock.Any(), gomock.Any()).Return([]entities.LastExecution{}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			wantErr: true,
			errMsg:  "fail",
		},
		{
			name: "failed while last executions retrieving",
			args: args{
				endpoints: map[string]struct{}{
					"e1": {},
					"e2": {},
				},
				ctx: func() context.Context {
					return context.Background()
				},
				cassandraCallNumberLimit: 1,
				workersTimeout:           10000,
			},
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastExecutions(gomock.Any(), gomock.Any()).
						Return([]entities.LastExecution{}, errors.New("fail"))
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastResultByEndpointID("e2").Return(entities.ExecutionResult{
						TaskInstanceID:  "i2",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)

					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
							TaskName:  "task1",
						}, nil)

					repo.EXPECT().GetMinimalInstanceByID("i2").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskName:  "task2",
							TaskID:    "t2",
						}, nil)

					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			wantErr: true,
			errMsg:  "fail",
		},
	}

	for _, tt := range tests {
		mockCtrl = gomock.NewController(t)

		config.Config.CassandraConcurrentCallNumber = tt.args.cassandraCallNumberLimit
		config.Config.ClosestTasksWorkersTimeoutSec = tt.args.workersTimeout

		tasks := &Tasks{
			tasksRepo:        tt.repos.tasksRepo(),
			taskInstanceRepo: tt.repos.taskInstanceRepo(),
			execResultsRepo:  tt.repos.execResultsRepo(),
		}
		got, err := tasks.findLastRuns(tt.args.ctx(), "p1", tt.args.endpoints)

		mockCtrl.Finish()

		if tt.wantErr {
			Ω(err).NotTo(BeNil(), tt.name)
			Ω(err.Error()).Should(ContainSubstring(tt.errMsg), tt.name)
		} else {
			Ω(err).To(BeNil(), tt.name)
			Ω(got).To(Equal(tt.want), tt.name)
		}
	}
}

func TestGetFullLastTaskDataForEndpoint(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrl *gomock.Controller

	type repos struct {
		execResultsRepo  func() ExecutionResultsRepo
		taskInstanceRepo func() InstancesRepo
		tasksRepo        func() Repo
	}
	tests := []struct {
		name    string
		repos   repos
		want    *taskData
		wantErr bool
		errMsg  string
	}{
		{
			name: "successfully find last runs",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
							TaskName:  "task1",
						}, nil)
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			want: &taskData{
				id:         "t1",
				instanceID: "i1",
				name:       "task1",
				run:        time.Unix(1, 0),
				status:     statuses.TaskInstanceSuccess,
			},
		},
		{
			name: "successfully find last runs (instance without name)",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
						}, nil)

					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetName("p1", "t1").
						Return("task1", nil)
					return repo
				},
			},
			want: &taskData{
				id:         "t1",
				instanceID: "i1",
				name:       "task1",
				run:        time.Unix(1, 0),
				status:     statuses.TaskInstanceSuccess,
			},
		},
		{
			name: "wrong partnerID)",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p2",
							TaskID:    "t1",
							TaskName:  "task1",
						}, nil)

					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			want: nil,
		},
		{
			name: "failed to fetch results",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").
						Return(entities.ExecutionResult{}, errors.New("fail"))
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			wantErr: true,
			errMsg:  "fail",
		},
		{
			name: "no results found",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			want: nil,
		},
		{
			name: "failed to retrieve instances",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{}, errors.New("fail"))
					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					return repo
				},
			},
			wantErr: true,
			errMsg:  "fail",
		},
		{
			name: "failed to get task name",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
						}, nil)

					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetName("p1", "t1").
						Return("", errors.New("fail"))
					return repo
				},
			},
			wantErr: true,
			errMsg:  "fail",
		},
		{
			name: "empty task name",
			repos: repos{
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
						}, nil)

					return repo
				},
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetName("p1", "t1").
						Return("", nil)
					return repo
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		mockCtrl = gomock.NewController(t)

		tasks := &Tasks{
			tasksRepo:        tt.repos.tasksRepo(),
			taskInstanceRepo: tt.repos.taskInstanceRepo(),
			execResultsRepo:  tt.repos.execResultsRepo(),
		}
		got, err := tasks.getFullLastTaskDataForEndpoint("p1", "e1", &instanceCache{instances: map[string]entities.TaskInstance{}})

		mockCtrl.Finish()

		if tt.wantErr {
			Ω(errors.Cause(err)).Should(MatchError(tt.errMsg), tt.name)
		} else {
			Ω(err).To(BeNil(), tt.name)
			Ω(got).To(Equal(tt.want), tt.name)
		}
	}
}

func TestMapNextRuns(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	type repos struct {
		tasksRepo func() Repo
	}
	tests := []struct {
		name    string
		repos   repos
		want    entities.EndpointsClosestTasks
		wantErr bool
		errMsg  string
	}{
		{
			name: "successfully map next runs",
			repos: repos{
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetNext("p1").Return([]entities.Task{
						{
							ID:                "id0",
							RunTimeUTC:        time.Unix(1, 0),
							Name:              "task0",
							ManagedEndpointID: "e0",
						},
						{
							ID:                "id1",
							RunTimeUTC:        time.Unix(2, 0),
							ManagedEndpointID: "e1",
						},
						{
							ID:                "id2",
							RunTimeUTC:        time.Unix(3, 0),
							Name:              "task2",
							ManagedEndpointID: "e2",
						},
						{
							ID:                "id3",
							RunTimeUTC:        time.Unix(4, 0),
							Name:              "task3",
							ManagedEndpointID: "e3",
						},
					}, nil)
					return repo
				},
			},
			want: entities.EndpointsClosestTasks{
				"e0": {
					Next: &entities.ClosestTask{
						ID:      "id0",
						Name:    "task0",
						RunDate: 1,
					},
				},
				"e2": {
					Next: &entities.ClosestTask{
						ID:      "id2",
						Name:    "task2",
						RunDate: 3,
					},
				},
				"e3": {
					Next: &entities.ClosestTask{
						ID:      "id3",
						Name:    "task3",
						RunDate: 4,
					},
				},
			},
		},
		{
			name: "successfully map next runs with disabled task",

			repos: repos{
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetNext("p1").Return([]entities.Task{
						{
							ID:                "id0",
							RunTimeUTC:        time.Unix(1, 0),
							Name:              "task0",
							ManagedEndpointID: "e0",
						},
						{
							ID:                "id1",
							RunTimeUTC:        time.Unix(2, 0),
							Name:              "task1",
							ManagedEndpointID: "e1",
							State:             statuses.TaskStateDisabled,
						},
					}, nil)
					return repo
				},
			},
			want: entities.EndpointsClosestTasks{
				"e0": {
					Next: &entities.ClosestTask{
						ID:      "id0",
						Name:    "task0",
						RunDate: 1,
					},
				},
			},
		},
		{
			name: "successfully map next runs with Postponed Recurrent task",
			repos: repos{
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetNext("p1").Return([]entities.Task{
						{
							ID:                "id0",
							RunTimeUTC:        time.Unix(1, 0),
							PostponedRunTime:  time.Unix(2, 0),
							Name:              "task0",
							ManagedEndpointID: "e0",
						},
						{
							ID:                "id1",
							RunTimeUTC:        time.Unix(3, 0),
							Name:              "task1",
							ManagedEndpointID: "e0",
						},
					}, nil)
					return repo
				},
			},
			want: entities.EndpointsClosestTasks{
				"e0": {
					Next: &entities.ClosestTask{
						ID:      "id0",
						Name:    "task0",
						RunDate: 2,
					},
				},
			},
		},
		{
			name: "failed to get next tasks",
			repos: repos{
				tasksRepo: func() Repo {
					repo := NewMockRepo(mockCtrl)
					repo.EXPECT().GetNext("p1").Return([]entities.Task{}, errors.New("fail"))
					return repo
				},
			},
			wantErr: true,
			errMsg:  "fail",
		},
	}

	for _, tt := range tests {
		mockCtrl = gomock.NewController(t)

		tasks := &Tasks{
			tasksRepo: tt.repos.tasksRepo(),
		}
		got, err := tasks.mapNextRuns("p1", map[string]struct{}{
			"e0": {},
			"e1": {},
			"e2": {},
			"e3": {},
		})

		mockCtrl.Finish()

		if tt.wantErr {
			Ω(errors.Cause(err)).Should(MatchError(tt.errMsg), tt.name)
		} else {
			Ω(err).To(BeNil(), tt.name)
			Ω(got).To(Equal(tt.want), tt.name)
		}
	}
}

func TestMapLastRuns(t *testing.T) {
	RegisterTestingT(t)
	type args struct {
		prev  map[string]*taskData
		tasks entities.EndpointsClosestTasks
	}
	tests := []struct {
		name         string
		args         args
		closestTasks entities.EndpointsClosestTasks
		wantErr      bool
		errMsg       string
	}{
		{
			name: "successfully map last runs",
			args: args{
				prev: map[string]*taskData{
					"e1": {
						id:     "id1",
						run:    time.Unix(2, 0),
						status: statuses.TaskInstanceSuccess,
						name:   "task1",
					},
					"e3": {
						id:     "id1",
						run:    time.Unix(2, 0),
						status: statuses.TaskInstanceFailed,
						name:   "task1",
					},
				},
				tasks: entities.EndpointsClosestTasks{
					"e0": {
						Next: &entities.ClosestTask{
							ID:      "id0",
							Name:    "task0",
							RunDate: 1,
						},
					},
					"e2": {
						Next: &entities.ClosestTask{
							ID:      "id2",
							Name:    "task2",
							RunDate: 3,
						},
					},
					"e3": {
						Next: &entities.ClosestTask{
							ID:      "id3",
							Name:    "task3",
							RunDate: 4,
						},
					},
				},
			},
			closestTasks: entities.EndpointsClosestTasks{
				"e0": {
					Next: &entities.ClosestTask{
						ID:      "id0",
						Name:    "task0",
						RunDate: 1,
					},
				},
				"e1": {
					Previous: &entities.ClosestTask{
						ID:      "id1",
						Name:    "task1",
						RunDate: 2,
						Status:  statuses.TaskInstanceSuccessText,
					},
				},
				"e2": {
					Next: &entities.ClosestTask{
						ID:      "id2",
						Name:    "task2",
						RunDate: 3,
					},
				},
				"e3": {
					Next: &entities.ClosestTask{
						ID:      "id3",
						Name:    "task3",
						RunDate: 4,
					},
					Previous: &entities.ClosestTask{
						ID:      "id1",
						Name:    "task1",
						RunDate: 2,
						Status:  statuses.TaskInstanceFailedText,
					},
				},
			},
		},
		{
			name: "ignore nil task values",
			args: args{
				prev: map[string]*taskData{
					"e1": nil,
				},
				tasks: entities.EndpointsClosestTasks{
					"e1": {
						Next: &entities.ClosestTask{
							ID:      "id1",
							Name:    "task1",
							RunDate: 1,
						},
					},
				},
			},
			closestTasks: entities.EndpointsClosestTasks{
				"e1": {
					Next: &entities.ClosestTask{
						ID:      "id1",
						Name:    "task1",
						RunDate: 1,
					},
				},
			},
		},
		{
			name: "error while status is invalid",
			args: args{
				prev: map[string]*taskData{
					"e1": {
						id:     "id1",
						run:    time.Unix(2, 0),
						status: math.MaxInt32,
						name:   "task1",
					},
				},
				tasks: entities.EndpointsClosestTasks{
					"e1": {
						Next: &entities.ClosestTask{
							ID:      "id1",
							Name:    "task1",
							RunDate: 1,
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "incorrect Task Instance Status: 2147483647",
		},
	}

	for _, tt := range tests {

		tasks := &Tasks{}

		closestTasks, err := tasks.mapLastRuns(tt.args.prev, tt.args.tasks)

		if tt.wantErr {
			Ω(errors.Cause(err)).Should(MatchError(tt.errMsg), tt.name)
		} else {
			Ω(err).To(BeNil(), tt.name)
			Ω(closestTasks).To(Equal(tt.closestTasks), tt.name)
		}
	}
}

func TestGetClosestTasks(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)

	type dependencies struct {
		tasksRepo        func() Repo
		taskInstanceRepo func() InstancesRepo
		execResultsRepo  func() ExecutionResultsRepo
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       dependencies
		args         args
		closestTasks entities.EndpointsClosestTasks
		wantErr      bool
		errMsg       string
	}{
		{
			name: "successfully get closest tasks",
			fields: dependencies{
				tasksRepo: func() Repo {
					tasksRepo := NewMockRepo(mockCtrl)
					tasksRepo.EXPECT().
						GetNext(gomock.Any()).
						Return([]entities.Task{
							{
								ID:                "id1",
								RunTimeUTC:        time.Unix(2, 0),
								ManagedEndpointID: "e1",
							},
							{
								ID:                "id2",
								RunTimeUTC:        time.Unix(3, 0),
								Name:              "task2",
								ManagedEndpointID: "e2",
							},
						}, nil).
						Times(1)
					return tasksRepo
				},

				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID("e1").Return(entities.ExecutionResult{
						TaskInstanceID:  "i1",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastResultByEndpointID("e2").Return(entities.ExecutionResult{
						TaskInstanceID:  "i2",
						ExecutionStatus: statuses.TaskInstanceSuccess,
						UpdatedAt:       time.Unix(1, 0),
					}, nil)
					repo.EXPECT().GetLastExecutions(gomock.Any(), gomock.Any()).Return([]entities.LastExecution{}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID("i1").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t1",
							TaskName:  "task1",
						}, nil)

					repo.EXPECT().GetMinimalInstanceByID("i2").
						Return(entities.TaskInstance{
							PartnerID: "p1",
							TaskID:    "t2",
							TaskName:  "task2",
						}, nil)

					return repo
				},
			},
			args: args{
				ctx: defaultContext(),
			},
			closestTasks: func() entities.EndpointsClosestTasks {
				t := entities.EndpointsClosestTasks{
					"e1": {
						Previous: &entities.ClosestTask{
							Name:    "task1",
							RunDate: 1,
							Status:  statuses.TaskInstanceSuccessText,
						},
					},
					"e2": {
						Previous: &entities.ClosestTask{
							Name:    "task2",
							RunDate: 1,
							Status:  statuses.TaskInstanceSuccessText,
						},
						Next: &entities.ClosestTask{
							Name:    "task2",
							RunDate: 3,
						},
					},
				}
				return t
			}(),
		},
		{
			name: "error while getting next runs",
			fields: dependencies{
				tasksRepo: func() Repo {
					tasksRepo := NewMockRepo(mockCtrl)
					tasksRepo.EXPECT().
						GetNext(gomock.Any()).
						Return([]entities.Task{}, errors.New("fail")).
						Times(1)

					return tasksRepo
				},
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID(gomock.Any()).
						Return(entities.ExecutionResult{}, nil).AnyTimes()
					repo.EXPECT().GetLastExecutions(gomock.Any(), gomock.Any()).Return([]entities.LastExecution{}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID(gomock.Any()).
						Return(entities.TaskInstance{}, nil).AnyTimes()
					return repo
				},
			},
			args: args{
				ctx: defaultContext(),
			},
			wantErr: true,
			errMsg:  "fail",
		},
		{
			name: "error while getting previous runs",
			fields: dependencies{
				tasksRepo: func() Repo {
					tasksRepo := NewMockRepo(mockCtrl)
					tasksRepo.EXPECT().GetNext(gomock.Any()).
						Return([]entities.Task{}, nil).AnyTimes()
					return tasksRepo
				},
				execResultsRepo: func() ExecutionResultsRepo {
					repo := NewMockExecutionResultsRepo(mockCtrl)
					repo.EXPECT().GetLastResultByEndpointID(gomock.Any()).
						Return(entities.ExecutionResult{}, errors.New("fail")).MinTimes(1)
					repo.EXPECT().GetLastExecutions(gomock.Any(), gomock.Any()).Return([]entities.LastExecution{}, nil)
					return repo
				},
				taskInstanceRepo: func() InstancesRepo {
					repo := NewMockInstancesRepo(mockCtrl)
					repo.EXPECT().GetMinimalInstanceByID(gomock.Any()).
						Return(entities.TaskInstance{}, nil).AnyTimes()
					return repo
				},
			},
			args: args{
				ctx: defaultContext(),
			},
			wantErr: true,
			errMsg:  "fail",
		},
		{
			name: "error while getting invalid partnerID",
			fields: dependencies{
				tasksRepo: func() Repo {
					tasksRepo := NewMockRepo(mockCtrl)
					return tasksRepo
				},

				taskInstanceRepo: func() InstancesRepo {
					taskInstanceRepo := NewMockInstancesRepo(mockCtrl)
					return taskInstanceRepo
				},
				execResultsRepo: func() ExecutionResultsRepo {
					e := NewMockExecutionResultsRepo(mockCtrl)
					return e
				},
			},
			args: args{
				ctx: func() context.Context {
					return context.Background()
				}(),
			},
			wantErr: true,
			errMsg:  "can't get parameter from context: partnerID",
		},
	}

	for _, tt := range tests {
		mockCtrl = gomock.NewController(t)

		config.Config.CassandraConcurrentCallNumber = 1
		config.Config.ClosestTasksWorkersTimeoutSec = 1000

		tasks := NewTasks(
			tt.fields.tasksRepo(),
			nil,
			tt.fields.taskInstanceRepo(),
			tt.fields.execResultsRepo(),
			nil,
			nil,
			nil)

		closestTasks, err := tasks.GetClosestTasks(tt.args.ctx, entities.EndpointsInput{"e1", "e2"})

		mockCtrl.Finish()

		if tt.wantErr {
			Ω(err).NotTo(BeNil(), tt.name)
			Ω(err.Error()).Should(ContainSubstring(tt.errMsg), tt.name)
		} else {
			Ω(err).To(BeNil(), tt.name)
			for _, v := range closestTasks {
				if v.Previous != nil {
					v.Previous.ID = ""
				}
				if v.Next != nil {
					v.Next.ID = ""
				}
			}

			Ω(closestTasks).To(Equal(tt.closestTasks), tt.name)
		}
	}
}

// Benchmark usecase test
func BenchmarkGetClosestTasks(t *testing.B) { // nolint
	RegisterTestingT(t)
	endpoints, nextTasks, results := generateTestData()

	tasksRepo := taskRepoMock{nextTasks: nextTasks}
	resultsRepo := resRepoMock{results: results}
	instanceRepo := instRepoMock{inst: entities.TaskInstance{
		TaskID:   "taskID",
		TaskName: "PREV_TEST_TASK",
	}}
	taskUC := NewTasks(&tasksRepo, nil, &instanceRepo, &resultsRepo, nil, nil, nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, config.PartnerIDKeyCTX, "p1")

	config.Config.CassandraConcurrentCallNumber = 30
	config.Config.ClosestTasksWorkersTimeoutSec = 10
	for i := 0; i < 1000; i++ {
		_, err := taskUC.GetClosestTasks(ctx, endpoints)
		fmt.Println(i)
		Ω(err).Should(BeNil())
	}

}

func defaultContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, config.PartnerIDKeyCTX, "p1")
	return ctx
}
