package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

const defaultMsg = `failed on unexpected value of result "%v"`

func TestNewScheduler(t *testing.T) {
	RegisterTestingT(t)
	expected := &Scheduler{}
	actual := NewScheduler(nil, nil, nil)
	Ω(actual).To(Equal(expected), fmt.Sprintf(defaultMsg, expected))
}

func TestScheduler_ProcessTasks(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	currentTime := time.Now().UTC().Truncate(time.Minute)
	logger.Load(config.Config.Log)

	tc := []struct {
		name    string
		payload func() (SchedulerUC, map[tasking.Regularity]SchedulerTypeUC, logger.Logger)
	}{
		{
			name: "test-case1: can't get tasks",
			payload: func() (uc SchedulerUC, uc2 map[tasking.Regularity]SchedulerTypeUC, log logger.Logger) {
				ucSchedulerMock := NewMockSchedulerUC(mockCtrl)
				err := errors.New("can't get tasks")
				ucSchedulerMock.EXPECT().GetScheduledTasks(currentTime).Return(nil, err).Times(1)
				return ucSchedulerMock, nil, nil
			},
		},
		{
			name: "test-case2: can't get usecase",
			payload: func() (uc SchedulerUC, uc2 map[tasking.Regularity]SchedulerTypeUC, log logger.Logger) {
				ucSchedulerMock := NewMockSchedulerUC(mockCtrl)
				tasks := map[tasking.Regularity][]models.Task{
					tasking.OneTime: {},
				}
				ucSchedulerMock.EXPECT().GetScheduledTasks(currentTime).Return(tasks, nil).Times(1)
				ucSchedulerMock.EXPECT().UpdateSchedulerTime(currentTime).Return(nil).Times(1)
				return ucSchedulerMock, nil, nil
			},
		},
		{
			name: "test-case3: success execution",
			payload: func() (uc SchedulerUC, uc2 map[tasking.Regularity]SchedulerTypeUC, log logger.Logger) {
				ucSchedulerMock := NewMockSchedulerUC(mockCtrl)
				ucSchedulerType := NewMockSchedulerTypeUC(mockCtrl)
				ucSchedulerType.EXPECT().Process(gomock.Any(), gomock.Any(), []models.Task{})
				tasks := map[tasking.Regularity][]models.Task{
					tasking.OneTime: {},
				}
				usecases := map[tasking.Regularity]SchedulerTypeUC{
					tasking.OneTime: ucSchedulerType,
				}
				ucSchedulerMock.EXPECT().GetScheduledTasks(currentTime).Return(tasks, nil).Times(1)
				ucSchedulerMock.EXPECT().UpdateSchedulerTime(currentTime).Return(nil).Times(1)
				return ucSchedulerMock, usecases, nil
			},
		},
		{
			name: "test-case3: can't update scheduler time",
			payload: func() (uc SchedulerUC, uc2 map[tasking.Regularity]SchedulerTypeUC, log logger.Logger) {
				ucSchedulerMock := NewMockSchedulerUC(mockCtrl)
				ucSchedulerType := NewMockSchedulerTypeUC(mockCtrl)
				ucSchedulerType.EXPECT().Process(gomock.Any(), gomock.Any(), []models.Task{})
				tasks := map[tasking.Regularity][]models.Task{
					tasking.OneTime: {},
				}
				usecases := map[tasking.Regularity]SchedulerTypeUC{
					tasking.OneTime: ucSchedulerType,
				}
				err := "some error"
				ucSchedulerMock.EXPECT().GetScheduledTasks(currentTime).Return(tasks, nil).Times(1)
				ucSchedulerMock.EXPECT().UpdateSchedulerTime(currentTime).Return(fmt.Errorf(err)).Times(1)
				return ucSchedulerMock, usecases, nil
			},
		},
	}

	for _, test := range tc {
		tasksUC, schedulerUC, _ := test.payload()
		scheduler := NewScheduler(tasksUC, schedulerUC, logger.Log)
		ctx := context.Background()
		scheduler.ProcessTasks(ctx)
	}
}
