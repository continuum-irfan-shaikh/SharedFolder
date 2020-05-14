package scheduler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	ts "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/gomega"
)

const defaultMsg = `failed on unexpected value of result "%v"`

func TestNewTasks(t *testing.T) {
	RegisterTestingT(t)
	expected := &Scheduler{}
	actual := NewScheduler(nil, nil)
	Ω(actual).To(Equal(expected), fmt.Sprintf(defaultMsg, expected))
}

func TestTasks_GetScheduledTasks(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	currentTime := time.Now().UTC().Truncate(time.Minute)
	ctx := context.Background()

	type expected struct {
		err  string
		data map[ts.Regularity][]models.Task
	}

	tc := []struct {
		name     string
		payload  func() (SchedulerRepo, TaskRepo)
		expected expected
	}{
		{
			name: "test-case1: error with getting lastUpdate",
			payload: func() (repo SchedulerRepo, repo2 TaskRepo) {
				schedulerRepo := NewMockSchedulerRepo(mockCtrl)
				err := errors.New("some error")
				schedulerRepo.EXPECT().GetLastUpdate().Return(time.Time{}, err).Times(1)
				schedulerRepo.EXPECT().UpdateScheduler(currentTime).Return(err)
				return schedulerRepo, nil
			},
			expected: expected{
				err: "some error",
			},
		},
		{
			name: "test-case2: can't get range of time",
			payload: func() (repo SchedulerRepo, repo2 TaskRepo) {
				schedulerRepo := NewMockSchedulerRepo(mockCtrl)
				schedulerRepo.EXPECT().GetLastUpdate().Return(currentTime, nil).Times(1)
				return schedulerRepo, nil
			},
			expected: expected{
				err: "timestamps are equal",
			},
		},
		{
			name: "test-case3: can't GetByRunTimeRange",
			payload: func() (repo SchedulerRepo, repo2 TaskRepo) {
				schedulerRepo := NewMockSchedulerRepo(mockCtrl)
				tasksRepo := NewMockTaskRepo(mockCtrl)
				err := errors.New("some error")
				schedulerRepo.EXPECT().GetLastUpdate().Return(currentTime.Add(-3*time.Minute), nil).Times(1)
				tasksRepo.EXPECT().GetByRunTimeRange(ctx, gomock.Any()).Return(nil, err).Times(1)
				return schedulerRepo, tasksRepo
			},
			expected: expected{
				err: "some error",
			},
		},
		{
			name: "test-case4: successful result",
			payload: func() (repo SchedulerRepo, repo2 TaskRepo) {
				schedulerRepo := NewMockSchedulerRepo(mockCtrl)
				tasksRepo := NewMockTaskRepo(mockCtrl)
				tasks := []models.Task{{}}
				schedulerRepo.EXPECT().GetLastUpdate().Return(currentTime.Add(-3*time.Minute), nil).Times(1)
				tasksRepo.EXPECT().GetByRunTimeRange(ctx, gomock.Any()).Return(tasks, nil).Times(1)
				return schedulerRepo, tasksRepo
			},
			expected: expected{
				err: "",
				data: map[ts.Regularity][]models.Task{
					0: {{}},
				},
			},
		},
	}

	for _, test := range tc {
		schedulerRepo, TaskRepo := test.payload()
		scheduler := NewScheduler(schedulerRepo, TaskRepo)
		actual, err := scheduler.GetScheduledTasks(currentTime)
		if test.expected.err != "" {
			Ω(err).NotTo(BeNil(), fmt.Sprintf(defaultMsg, test.name))
			Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
		Ω(actual).To(Equal(test.expected.data), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestScheduler_UpdateSchedulerTime(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	schedulerRepo := NewMockSchedulerRepo(mockCtrl)
	scheduler := &Scheduler{
		schedulerRepo: schedulerRepo,
	}
	currentTime := time.Now()

	schedulerRepo.EXPECT().UpdateScheduler(currentTime).Return(nil).Times(1)
	err := scheduler.UpdateSchedulerTime(currentTime)
	Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, err))
}
