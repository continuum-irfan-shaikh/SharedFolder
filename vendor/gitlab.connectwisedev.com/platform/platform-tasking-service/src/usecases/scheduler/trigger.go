package scheduler

import (
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
	"golang.org/x/net/context"
)

// Trigger is a struct that represents shutdown and logoff triggers handling
type Trigger struct {
	log       logger.Logger
	triggerUC trigger.Usecase
}

// NewTrigger returns new Triggers usecase
func NewTrigger(log logger.Logger, triggersUC trigger.Usecase) *Trigger {
	return &Trigger{
		log:       log,
		triggerUC: triggersUC,
	}
}

//Process - processes all trigger for scheduler and activates/deactivates it
func (tr *Trigger) Process(ctx context.Context, currentTime time.Time, tasks []m.Task) {
	tasksActivate := make(map[gocql.UUID][]m.Task)
	tasksDeactivate := make(map[gocql.UUID][]m.Task)

	// grouping tasks by task ID and by process type
	for _, task := range tasks {
		if task.Schedule.EndRunTime.UTC().Equal(currentTime) || currentTime.After(task.Schedule.EndRunTime.UTC()) {
			tasksDeactivate[task.ID] = append(tasksDeactivate[task.ID], task)
			continue
		}

		tasksActivate[task.ID] = append(tasksActivate[task.ID], task)
	}

	// activating
	for _, groupedTasks := range tasksActivate {
		// sending by UCs
		go func(ctx context.Context, tasks []m.Task) {
			if err := tr.triggerUC.Activate(ctx, tasks); err != nil {
				tr.log.WarnfCtx(ctx,"Process: activate err for taskID %v", tasks[0].ID)
				return
			}
		}(ctx, groupedTasks)
	}

	// deactivating
	for _, groupedTasks := range tasksDeactivate {
		// sending by UCs
		go func(ctx context.Context, tasks []m.Task) {
			if err := tr.triggerUC.Deactivate(ctx, tasks); err != nil {
				tr.log.WarnfCtx(ctx, "Process: deactivate err for taskID %v", tasks[0].ID)
				return
			}
		}(ctx, groupedTasks)
	}
}
