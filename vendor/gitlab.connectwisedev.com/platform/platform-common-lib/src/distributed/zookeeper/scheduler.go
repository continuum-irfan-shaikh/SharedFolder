package zookeeper

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/scheduler"
)

var (
	schedulerCron *cron.Cron
	schedulerInit = false
	scheduledJobs []scheduler.ScheduledJob
)

// Job is a struct defining the actual scheduled job
// Implementing `Job` interface from cron package and `ScheduledJob` from scheduler package
type Job struct {
	Name     string
	Task     string
	Schedule string
}

// GetName returns job name
func (j Job) GetName() string {
	return j.Name
}

// GetTask returns job task
func (j Job) GetTask() string {
	return j.Task
}

// GetSchedule returns job schedule
func (j Job) GetSchedule() string {
	return j.Schedule
}

// Run initial entry point of a Job
func (j Job) Run() {
	Logger().Info(defaultTransaction, "Scheduling job `%s` for execution", j.GetName())
	_, err := Queue.CreateItem(nil, j.Task)
	if err != nil {
		Logger().Error(defaultTransaction, "Queue.CreateItemFailed", "Scheduler. Couldn't run a distributed job %v, err: %v", j.GetName(), err)
	}
}

// DistributedScheduler initializes distributed scheduler
func (schedulerImpl) DistributedScheduler(ctx context.Context, wg *sync.WaitGroup, jobs []scheduler.ScheduledJob, interval int) error {
	scheduledJobs = jobs
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(interval) * time.Second):
				sPeerID, leader, err := LeaderElector.BecomeALeader()
				if err != nil {
					Logger().Error(defaultTransaction, "LeaderElector.BecomeALeaderFailed", "become leader got error: %v", err)
				}
				if leader && !schedulerInit {
					// I'm a new leader
					startScheduler()
				}
				if sPeerID == undefined && schedulerInit {
					stopScheduler()
				}
			}
		}
	}()
	return nil
}

func startScheduler() {
	Logger().Info(defaultTransaction, "I'm a new leader. Initializing scheduler...")
	schedulerInit = true
	schedulerCron = cron.New()

	for _, sj := range scheduledJobs {
		job := Job{
			Name:     sj.GetName(),
			Task:     sj.GetTask(),
			Schedule: sj.GetSchedule(),
		}
		err := schedulerCron.AddJob(job.GetSchedule(), job)
		if err != nil {
			Logger().Error(defaultTransaction, "schedulerCron.AddJobFailed", "Couldn't add job %v with schedule %v, err: ", job.GetName(), job.GetSchedule(), err)
			continue
		}
	}
	schedulerCron.Start()
}

func stopScheduler() {
	if schedulerInit {
		Logger().Info(defaultTransaction, "Stopping scheduler...")
		schedulerInit = false
		schedulerCron.Stop()
	}
}
