package scheduler

import (
	"context"
	"time"

	t "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

//go:generate mockgen -destination=./scheduler_uc_mock_test.go -package=scheduler -source=./scheduler.go

//SchedulerUC - represents scheduler usecases
type SchedulerUC interface {
	GetScheduledTasks(currentTime time.Time) (map[t.Regularity][]models.Task, error)
	UpdateSchedulerTime(currentTime time.Time) error
}

//SchedulerTypeUC - represents scheduler type usecase for different type of tasks
type SchedulerTypeUC interface {
	Process(ctx context.Context, currentTime time.Time, tasks []models.Task)
}

//NewScheduler - returns a new instance of the handler (scheduler)
func NewScheduler(s SchedulerUC, st map[t.Regularity]SchedulerTypeUC, log logger.Logger) *Scheduler {
	return &Scheduler{
		scheduler:        s,
		schedulerTypeUCs: st,
		log:              log,
	}
}

//Scheduler - represents a scheduler handler
type Scheduler struct {
	scheduler        SchedulerUC
	schedulerTypeUCs map[t.Regularity]SchedulerTypeUC
	log              logger.Logger
}

//ProcessTasks - gets tasks, run suitable usecase
func (s *Scheduler) ProcessTasks(ctx context.Context) {
	currentTime := time.Now().UTC().Truncate(time.Minute)
	ctx = transactionID.NewContext()

	groupedTasks, err := s.scheduler.GetScheduledTasks(currentTime)
	if err != nil {
		s.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, "ProcessTasks: can't get scheduled tasks: %s", err)
		return
	}

	for key, val := range groupedTasks {
		uc, ok := s.schedulerTypeUCs[key]
		if !ok {
			s.log.WarnfCtx(ctx, "ProcessTasks: can't get usecase for: %v", key)
			continue
		}
		uc.Process(ctx, currentTime, val)
	}

	if err = s.scheduler.UpdateSchedulerTime(currentTime); err != nil {
		s.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "can't UpdateSchedulerTime, err: %v", err)
	}
}
