package scheduler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mockLoggerTasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestNewShutdownTrigger(t *testing.T) {
	RegisterTestingT(t)
	expected := &Trigger{}
	actual := NewTrigger(nil, nil)
	Ω(actual).To(Equal(expected), fmt.Sprintf(defaultMsg, expected))
}

func TestTrigger_Process(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	currentTime := time.Now()
	triggerType := "SomeType"
	type args struct {
		currentTime time.Time
		tasks       []models.Task
	}

	type mock struct {
		log       func() logger.Logger
		triggerUC func() trigger.Usecase
	}

	type test struct {
		name string
		args args
		mock mock
		init func(test)
	}

	id := gocql.TimeUUID()
	tests := []test{
		{
			name: "test-case1: can't match time, currentTime",
			mock: mock{
				log: func() logger.Logger {
					log := mockLoggerTasking.NewMockLogger(mockCtrl)
					return log
				},
				triggerUC: func() trigger.Usecase {
					uc := mocks.NewMockUsecase(mockCtrl)
					uc.EXPECT().Activate(gomock.Any(), gomock.Any()).AnyTimes()
					uc.EXPECT().Deactivate(gomock.Any(), gomock.Any()).AnyTimes()
					return uc
				},
			},
			args: args{
				tasks: []models.Task{
					{
						Schedule: apiModels.Schedule{
							StartRunTime: currentTime.Add(1 * time.Minute),
							TriggerTypes: []string{triggerType},
						},
					},
				},
				currentTime: currentTime,
			},
		},
		{
			name: "test-case2: ok",
			mock: mock{
				log: func() logger.Logger {
					log := mockLoggerTasking.NewMockLogger(mockCtrl)
					return log
				},
				triggerUC: func() trigger.Usecase {
					uc := mocks.NewMockUsecase(mockCtrl)
					uc.EXPECT().Activate(gomock.Any(), gomock.Any()).AnyTimes()
					uc.EXPECT().Deactivate(gomock.Any(), gomock.Any()).AnyTimes()
					return uc
				},
			},
			args: args{
				tasks: []models.Task{
					{
						Schedule: apiModels.Schedule{
							StartRunTime: currentTime,
							TriggerTypes: []string{triggerType},
						},
					},
					{
						Schedule: apiModels.Schedule{
							EndRunTime:   currentTime,
							TriggerTypes: []string{triggerType},
						},
					},
				},
				currentTime: currentTime,
			},
		},
		{
			name: "test-case3: UC returned err",
			mock: mock{
				log: func() logger.Logger {
					log := mockLoggerTasking.NewMockLogger(mockCtrl)
					return log
				},
				triggerUC: func() trigger.Usecase {
					uc := mocks.NewMockUsecase(mockCtrl)
					uc.EXPECT().Activate(gomock.Any(), gomock.Any()).Return(errors.New("err")).AnyTimes()
					uc.EXPECT().Deactivate(gomock.Any(), gomock.Any()).Return(errors.New("err")).AnyTimes()
					return uc
				},
			},
			args: args{
				tasks: []models.Task{
					{
						ID: id,
						Schedule: apiModels.Schedule{
							StartRunTime: currentTime,
							TriggerTypes: []string{"wrongType"},
						},
					},
					{
						ID: id,
						Schedule: apiModels.Schedule{
							EndRunTime:   currentTime,
							TriggerTypes: []string{"wrongType"},
						},
					},
				},
				currentTime: currentTime,
			},
		},
	}

	for _, t := range tests {
		fmt.Println(t.name)
		tr := Trigger{
			log:       logger.Log,
			triggerUC: t.mock.triggerUC(),
		}
		tr.Process(context.Background(), t.args.currentTime, t.args.tasks)
		time.Sleep(time.Second / 2)
	}
}
