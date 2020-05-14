package scheduler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/scheduler"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/zookeeper"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/memcache"
	s "gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/scheduler"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
)

// Distributed job names
const (
	processTasks              = "processTasks"
	checkForRetainedData      = "checkForRetainedData"
	checkForExpiredExecutions = "checkForExpiredExecutions"
	recalculateTasks          = "recalculateTasks"
	loadTriggersToCache       = "loadTriggersToCache"
	activeTriggersReopening   = "activeTriggersReopening"
)

// timeout between zookeeper re-connections
const reconnTimeoutSec = 5

// Job implementing `scheduler.DistributedJob` interface
type Job struct {
	name         string
	log          logger.Logger
	callbackFunc func(context.Context)
}

// Callback is necessary for implementing `scheduler.DistributedJob` interface
func (j *Job) Callback(i ...interface{}) {
	ctx, ok := i[0].(context.Context)
	if !ok {
		j.log.ErrfCtx(ctx,errorcode.ErrorCantProcessData, "can't get Ctx from distributed job")
		return
	}
	j.callbackFunc(ctx)
}

// GetName is necessary for implementing `scheduler.DistributedJob` interface
func (j *Job) GetName() string {
	return j.name
}

// LoadDTO describes DTO for scheduler
type LoadDTO struct {
	Ctx       context.Context
	Log       logger.Logger
	Conf      config.Configuration
	WG        *sync.WaitGroup
	DS        scheduler.Interface // distributed scheduler
	Service   s.Service
	Scheduler *s.Scheduler
	Loader    *memcache.Loader
	Trigger   trigger.Usecase
}

// Load scheduler service.
func Load(l LoadDTO) error {
	jobs := []scheduler.DistributedJob{
		&Job{
			name:         checkForRetainedData,
			log:          l.Log,
			callbackFunc: l.Service.CheckForRetainedData,
		},
		&Job{
			name:         checkForExpiredExecutions,
			log:          l.Log,
			callbackFunc: l.Service.CheckForExpiredExecutions,
		},
		&Job{
			name:         recalculateTasks,
			log:          l.Log,
			callbackFunc: l.Service.RecalculateTasks,
		},
		// New Scheduler Implementation
		&Job{
			name:         processTasks,
			log:          l.Log,
			callbackFunc: l.Scheduler.ProcessTasks,
		},
		&Job{
			name:         loadTriggersToCache,
			log:          l.Log,
			callbackFunc: l.Loader.LoadTriggersToCache,
		},
		&Job{
			name:         activeTriggersReopening,
			log:          l.Log,
			callbackFunc: l.Trigger.ActiveTriggersReopening,
		},
	}

	var scheduledJobs = make([]scheduler.ScheduledJob, 0, len(config.Config.ScheduledJobs))
	for _, job := range l.Conf.ScheduledJobs {
		scheduledJobs = append(scheduledJobs, job)
	}

	return do(l, scheduledJobs, jobs)
}

func do(l LoadDTO, schJobs []scheduler.ScheduledJob, distJobs []scheduler.DistributedJob) error {
	err := zookeeper.Init(strings.Split(l.Conf.ZookeeperHosts, ","), l.Conf.ZookeeperBasePath)
	if err != nil {
		return fmt.Errorf("failed to connect to zookeeper: %s", err)
	}

	go monitor(l, schJobs, distJobs)
	return nil
}

// listen for events and re-init zk connection in case of errors
func monitor(l LoadDTO, schJobs []scheduler.ScheduledJob, distJobs []scheduler.DistributedJob) {
	for i := 0; ; i++ {
		var (
			ctx, cancel = context.WithCancel(l.Ctx)
			kill        = func() {
				cancel()
				zookeeper.Client.Close()
			}
		)

		if i != 0 {
			time.Sleep(time.Second * reconnTimeoutSec)
			err := zookeeper.Init(strings.Split(l.Conf.ZookeeperHosts, ","), l.Conf.ZookeeperBasePath)
			if err != nil {
				l.Log.WarnfCtx(ctx, "failed to connect to zookeeper, reconnecting..: %s", err)
				kill()
				continue
			}
		}

		err := run(ctx, l, schJobs, distJobs)
		if err != nil {
			l.Log.WarnfCtx(ctx, "failed to run scheduled jobs, retrying..: %s", err)
			kill()
			continue
		}

		for {
			e := <-zookeeper.Client.Events()
			if e.State == zk.StateDisconnected || e.State == zk.StateAuthFailed || e.State == zk.StateExpired {
				break
			}
		}

		l.Log.WarnfCtx(ctx, "lost connection with zookeeper, reconnecting..")
		kill()
	}
}

func run(ctx context.Context, l LoadDTO, schJobs []scheduler.ScheduledJob, distJobs []scheduler.DistributedJob) error {
	err := l.DS.DistributedScheduler(ctx, l.WG, schJobs, l.Conf.JobSchedulingInterval)
	if err != nil {
		l.Log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "can't perform InitDistributedScheduling function, err: %v", err)
		return err
	}

	err = l.DS.DistributedJobListener(ctx, l.WG, distJobs, l.Conf.JobListeningInterval)
	if err != nil {
		l.Log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "can't perform InitDistributedTasksListener function, err: %v", err)
		return err
	}

	return nil
}
