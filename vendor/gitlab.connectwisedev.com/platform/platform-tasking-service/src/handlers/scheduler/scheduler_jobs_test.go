package scheduler

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	tcs := NewMockCounterService(nil)
	e := NewMockExecutionResultRepo(nil)
	eru := NewMockExecutionResultUpdateUC(nil)
	sr := NewMockSchedulerRepo(nil)

	expected := Service{taskCounterService: tcs, executionResultsRepo: e, execResultUpdateService: eru, schedulerRepo: sr}
	actual := New(tcs, e, eru, sr)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("New() = %v, want %v", actual, expected)
	}
}

func TestService_CheckForExpiredExecutions(t *testing.T) {
	var ctrl *gomock.Controller
	logger.Load(config.Config.Log)

	type fields struct {
		taskCounterService      CounterService
		executionResultsRepo    ExecutionResultRepo
		execResultUpdateService ExecutionResultUpdateUC
		schedulerRepo           SchedulerRepo
	}
	tests := []struct {
		name   string
		fields func() fields
	}{
		{
			name: "Success",
			fields: func() fields {
				t := time.Now().Add(-2 * time.Minute)

				eepiMock := mocks.NewMockExecutionExpirationPersistence(ctrl)
				eepiMock.EXPECT().GetByExpirationTime(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{}, nil).MinTimes(1)
				models.ExecutionExpirationPersistenceInstance = eepiMock

				sr := NewMockSchedulerRepo(ctrl)
				sr.EXPECT().GetLastExpiredExecutionCheck().Return(t, nil).MinTimes(1)
				sr.EXPECT().UpdateLastExpiredExecutionCheck(gomock.Any()).Return(nil).MinTimes(1)

				return fields{
					schedulerRepo: sr,
				}
			},
		},
		{
			name: "Success with sending",
			fields: func() fields {
				t := time.Now().Add(-2 * time.Minute)
				config.Config.ScriptingMsURL = ""

				eepiMock := mocks.NewMockExecutionExpirationPersistence(ctrl)
				eepiMock.EXPECT().GetByExpirationTime(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{{
					ExpirationTimeUTC:  time.Time{},
					PartnerID:          "test",
					TaskInstanceID:     gocql.UUID{},
					ManagedEndpointIDs: []gocql.UUID{{}},
				}}, nil).MinTimes(1)
				models.ExecutionExpirationPersistenceInstance = eepiMock

				sr := NewMockSchedulerRepo(ctrl)
				sr.EXPECT().GetLastExpiredExecutionCheck().Return(t, nil).MinTimes(1)
				sr.EXPECT().UpdateLastExpiredExecutionCheck(gomock.Any()).Return(nil).MinTimes(1)

				eru := NewMockExecutionResultUpdateUC(ctrl)
				eru.EXPECT().ProcessExecutionResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(1)

				return fields{
					schedulerRepo:           sr,
					execResultUpdateService: eru,
				}
			},
		},
		{
			name: "Failed to get last execution time",
			fields: func() fields {
				t := time.Now().Add(-2 * time.Minute)

				sr := NewMockSchedulerRepo(ctrl)
				sr.EXPECT().GetLastExpiredExecutionCheck().Return(t, errors.New("fail")).MinTimes(1)

				return fields{
					schedulerRepo: sr,
				}
			},
		},
		{
			name: "NotFound execution time",
			fields: func() fields {
				t := time.Now().Add(-2 * time.Minute)

				eepiMock := mocks.NewMockExecutionExpirationPersistence(ctrl)
				eepiMock.EXPECT().GetByExpirationTime(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{}, nil).MinTimes(1)
				models.ExecutionExpirationPersistenceInstance = eepiMock

				sr := NewMockSchedulerRepo(ctrl)
				sr.EXPECT().GetLastExpiredExecutionCheck().Return(t, gocql.ErrNotFound).MinTimes(1)
				sr.EXPECT().UpdateLastExpiredExecutionCheck(gomock.Any()).Return(nil).MinTimes(1)

				return fields{
					schedulerRepo: sr,
				}
			},
		},
		{
			name: "Failed to find exec expirations",
			fields: func() fields {
				t := time.Now().Add(-2 * time.Minute)

				eepiMock := mocks.NewMockExecutionExpirationPersistence(ctrl)
				eepiMock.EXPECT().GetByExpirationTime(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{}, errors.New("fail")).MinTimes(1)
				models.ExecutionExpirationPersistenceInstance = eepiMock

				sr := NewMockSchedulerRepo(ctrl)
				sr.EXPECT().GetLastExpiredExecutionCheck().Return(t, nil).MinTimes(1)

				return fields{
					schedulerRepo: sr,
				}
			},
		},
		{
			name: "Failed to update last execs",
			fields: func() fields {
				t := time.Now().Add(-2 * time.Minute)

				eepiMock := mocks.NewMockExecutionExpirationPersistence(ctrl)
				eepiMock.EXPECT().GetByExpirationTime(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{}, nil).MinTimes(1)
				models.ExecutionExpirationPersistenceInstance = eepiMock

				sr := NewMockSchedulerRepo(ctrl)
				sr.EXPECT().GetLastExpiredExecutionCheck().Return(t, nil).MinTimes(1)
				sr.EXPECT().UpdateLastExpiredExecutionCheck(gomock.Any()).Return(errors.New("fail")).MinTimes(1)

				return fields{
					schedulerRepo: sr,
				}
			},
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fields()
			s := New(f.taskCounterService, f.executionResultsRepo, f.execResultUpdateService, f.schedulerRepo)

			s.CheckForExpiredExecutions(context.Background())
			ctrl.Finish()
		})
	}
}

func TestService_updateExecutionResultsByExpiredExecutions(t *testing.T) {
	var ctrl *gomock.Controller
	logger.Load(config.Config.Log)

	type fields struct {
		execResultUpdateService ExecutionResultUpdateUC
	}
	tests := []struct {
		name   string
		fields func() fields
	}{
		{
			name: "Success",
			fields: func() fields {
				eru := NewMockExecutionResultUpdateUC(ctrl)
				eru.EXPECT().ProcessExecutionResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					execResultUpdateService: eru,
				}
			},
		},
		{
			name: "Failed to process results",
			fields: func() fields {
				eru := NewMockExecutionResultUpdateUC(ctrl)
				eru.EXPECT().ProcessExecutionResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("fail"))

				return fields{
					execResultUpdateService: eru,
				}
			},
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, nil, tt.fields().execResultUpdateService, nil)
			s.updateExecutionResultsByExpiredExecutions(context.Background(), "test", tasking.ExpiredExecution{ManagedEndpointIDs: []gocql.UUID{{}}})
			ctrl.Finish()
		})
	}
}

func TestService_RecalculateTasks(t *testing.T) {
	var ctrl *gomock.Controller
	logger.Load(config.Config.Log)

	type fields struct {
		taskCounterService CounterService
	}
	tests := []struct {
		name   string
		fields func() fields
	}{
		{
			name: "Success",
			fields: func() fields {
				cs := NewMockCounterService(ctrl)
				cs.EXPECT().RecalculateAllCounters(gomock.Any(), gomock.Any())

				return fields{
					taskCounterService: cs,
				}
			},
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields().taskCounterService, nil, nil, nil)
			s.RecalculateTasks(nil)
			ctrl.Finish()
		})
	}
}
