package scheduler

import (
	"context"
	"errors"
	"time"

	ts "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

//NewScheduler - initialize usecase scheduler struct
func NewScheduler(schedulerRepo SchedulerRepo, taskRepo TaskRepo) *Scheduler {
	return &Scheduler{
		schedulerRepo: schedulerRepo,
		taskRepo:      taskRepo,
	}
}

//Scheduler - represents scheduler usecase
type Scheduler struct {
	schedulerRepo SchedulerRepo
	taskRepo      TaskRepo
}

//GetScheduledTasks - returns scheduled tasks, sorted by Regularity
func (s *Scheduler) GetScheduledTasks(currentTime time.Time) (map[ts.Regularity][]models.Task, error) {
	lastUpdate, err := s.schedulerRepo.GetLastUpdate()
	if err != nil {
		// means that no record was created
		err := s.schedulerRepo.UpdateScheduler(currentTime)
		if err != nil {
			return nil, err
		}
		lastUpdate = currentTime
	}

	timeRange, err := s.getTimeRange(lastUpdate, currentTime)
	if err != nil {
		return nil, err
	}

	ctx := context.Background() // Deprecated dependency in repository
	tasks, err := s.taskRepo.GetByRunTimeRange(ctx, timeRange)
	if err != nil {
		return nil, err
	}

	return s.sortTasks(tasks), nil
}

//UpdateSchedulerTime - update scheduler time
func (s *Scheduler) UpdateSchedulerTime(currentTime time.Time) error {
	return s.schedulerRepo.UpdateScheduler(currentTime)
}

func (*Scheduler) sortTasks(tasks []models.Task) map[ts.Regularity][]models.Task {
	sortedTasks := make(map[ts.Regularity][]models.Task, 0)
	for _, task := range tasks {
		r := task.Schedule.Regularity
		if _, ok := sortedTasks[r]; !ok {
			sortedTasks[r] = make([]models.Task, 0)
		}
		sortedTasks[r] = append(sortedTasks[r], task)
	}
	return sortedTasks
}

// getTimeRange - time range by minutes
func (*Scheduler) getTimeRange(lastUpdate, currentTime time.Time) ([]time.Time, error) {
	var timeRange []time.Time
	if lastUpdate.Equal(currentTime) {
		return timeRange, errors.New("timestamps are equal")
	}

	// t - new time stamp
	t := lastUpdate.Add(time.Minute)
	for !t.After(currentTime) {
		timeRange = append(timeRange, t)
		t = t.Add(time.Minute)
	}
	return timeRange, nil
}
