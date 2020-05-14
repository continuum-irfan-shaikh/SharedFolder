package trigger

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	managedEndpoints "gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/managed-endpoints"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

// ActiveTriggersReopening reopens active triggers
func (tr *Service) ActiveTriggersReopening(ctx context.Context) {
	ctx = transactionID.NewContext()
	activeTriggers, err := tr.triggerRepo.GetAll()
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetActiveTriggers, "ActiveTriggersReopening: can't get active triggers. err - %v", err)
		return
	}

	if len(activeTriggers) < 1 {
		return
	}

	handled := make(map[gocql.UUID]struct{})
	for _, trigger := range activeTriggers {
		if _, ok := handled[trigger.TaskID]; ok {
			continue
		}

		handled[trigger.TaskID] = struct{}{}
		//go keyword can be added to parallelize job. But no needed for now.
		tr.reopenTrigger(ctx, trigger.PartnerID, trigger.TaskID)
	}

	tr.log.InfofCtx(ctx, "ActiveTriggersReopening: successfully reopened active triggers")
	return
}

func (tr *Service) reopenTrigger(ctx context.Context, partnerID string, taskID gocql.UUID) {
	tasks, err := tr.taskRepo.GetByIDs(ctx, nil, partnerID, false, taskID)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, "ActiveTriggersReopening: error while getting task, err %v", err)
		return
	}

	if len(tasks) == 0 {
		tr.log.WarnfCtx(ctx, "ActiveTriggersReopening: there is no tasks found, id %v", taskID)
		return
	}

	if tasks[0].TargetsByType.HasEndpointTypeOnly() || tasks[0].Targets.Type == models.ManagedEndpoint {
		return
	}

	tasks = tr.filterInactive(tasks)
	if len(tasks) == 0 {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "ActiveTriggersReopening: there is no active tasks found")
		return
	}

	newEndpoints, err := tr.getNewEndpoints(ctx, tasks[0])
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "ActiveTriggersReopening: %s", err.Error())
	}

	//defaultEID is marker that shows that task's targets is empty
	defaultEID, _ := gocql.ParseUUID(models.DefaultEndpointUID)
	//if there is no newEndpoints and task is already created on dummy endpoint - nothing to update here
	if len(newEndpoints) < 1 && tasks[0].ManagedEndpointID == defaultEID {
		return
	}

	tasks = tr.createAndDeactivateTasks(tasks, newEndpoints)

	err = tr.taskRepo.InsertOrUpdate(ctx, tasks...)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "ActiveTriggersReopening: can't update tasks in DB: %v", err)
		return
	}

	instance, err := tr.getScheduledInstance(ctx, tasks)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, "ActiveTriggersReopening: can't get scheduled instance for task %v. err : %s", tasks[0], err.Error())
		return
	}

	activeTasks := tr.filterInactive(tasks)
	filteredTasks, filteredEndpoints := tr.filterDefaults(activeTasks)
	instance = tr.updateScheduledInstance(instance, filteredEndpoints)

	err = tr.taskInstanceRepo.Insert(ctx, instance)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "ActiveTriggersReopening: can't update task instance %c in DB: %s", instance, err.Error())
	}

	//update triggers for new endpoints
	err = tr.updateNewEndpoints(ctx, filteredTasks)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "ActiveTriggersReopening: can't update policy for new endpoints: %v", err)
		return
	}
}

func (tr *Service) filterInactive(oldTasks []models.Task) (activeTasks []models.Task) {
	for _, task := range oldTasks {
		if task.State != statuses.TaskStateInactive {
			activeTasks = append(activeTasks, task)
		}
	}
	return
}

//filterDefaults filters only those tasks and endpoints, that are not default ones
func (tr *Service) filterDefaults(activeTasks []models.Task) (filteredTasks []models.Task, filteredEndpoints []gocql.UUID) {
	//defaultEID is marker that shows that task's targets is empty
	defaultEID, _ := gocql.ParseUUID(models.DefaultEndpointUID)
	for _, task := range activeTasks {
		if task.ManagedEndpointID != defaultEID {
			filteredTasks = append(filteredTasks, task)
			filteredEndpoints = append(filteredEndpoints, task.ManagedEndpointID)
		}
	}
	return
}

func (tr *Service) getNewEndpoints(ctx context.Context, task models.Task) (new map[gocql.UUID]models.TargetType, err error) {
	//new endpoints we got from DG MS
	new, err = managedEndpoints.GetManagedEndpointsFromTargets(ctx, task, tr.httpClient)
	if err != nil {
		return nil, fmt.Errorf("getNewEndpoints: can't get managed endpoints for targets. err : %s", err.Error())
	}
	return
}

func (tr *Service) createAndDeactivateTasks(tasks []models.Task, newEndpoints map[gocql.UUID]models.TargetType) (result []models.Task) {
	newTask := tasks[0]
	newTask.State = statuses.TaskStateActive

	//if Site or DG is empty
	if len(newEndpoints) == 0 {
		result = tr.deactivateTasks(tasks)
		//defaultEID is marker that shows that task's targets is empty
		defaultEID, _ := gocql.ParseUUID(models.DefaultEndpointUID)
		//create new default task
		result = append(result, *newTask.CopyWithRunTime(defaultEID))
		return
	}

	for _, task := range tasks {
		if _, ok := newEndpoints[task.ManagedEndpointID]; ok {
			delete(newEndpoints, task.ManagedEndpointID)
			result = append(result, task)
			continue
		}
		task.State = statuses.TaskStateInactive
		result = append(result, task)
	}

	for e, tt := range newEndpoints {
		newTask.TargetType = tt
		result = append(result, *newTask.CopyWithRunTime(e))
	}
	return
}

func (tr *Service) deactivateTasks(tasks []models.Task) (result []models.Task) {
	for _, task := range tasks {
		task.State = statuses.TaskStateInactive
		result = append(result, task)
	}
	return
}

func (tr *Service) updateScheduledInstance(instance models.TaskInstance, newEndpoints []gocql.UUID) models.TaskInstance {
	statuses := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	for _, e := range newEndpoints {
		if status, ok := instance.Statuses[e]; !ok {
			statuses[e] = tr.getOverallStatus(instance)
		} else {
			statuses[e] = status
		}
	}

	instance.Statuses = statuses

	return instance
}

func (tr *Service) getScheduledInstance(ctx context.Context, tasks []models.Task) (instance models.TaskInstance, err error) {
	if len(tasks) == 0 {
		return instance, fmt.Errorf("getScheduledInstance: tasks slice is empty")
	}
	instances, err := tr.taskInstanceRepo.GetByIDs(ctx, tasks[0].LastTaskInstanceID)
	if err != nil {
		return instance, fmt.Errorf("getScheduledInstance: can't retrieve last task instance. err : %s", err.Error())
	}

	if len(instances) == 0 {
		return instance, fmt.Errorf("getScheduledInstance: instance with id : %v not found", tasks[0].LastTaskInstanceID)
	}

	instance = instances[0]

	if instance.StartedAt.Sub(tasks[0].CreatedAt).Minutes() <= 1 {
		return instance, nil
	}

	for {
		gotTI, err := tr.taskInstanceRepo.GetNearestInstanceAfter(instance.TaskID, instance.StartedAt)
		if err == gocql.ErrNotFound {
			return instance, nil
		}

		if err != nil {
			return instance, fmt.Errorf("getScheduledInstance: can't get nearest instance. partnerID : %v, started_at : %v, err: %s", instance.PartnerID, instance.StartedAt, err.Error())
		}

		if gotTI.TriggeredBy == "" {
			return gotTI, nil
		}
		instance = gotTI
	}
}

func (tr *Service) updateNewEndpoints(ctx context.Context, tasks []models.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	var errTotal error

	t := tasks[0]

	for _, triggerType := range t.Schedule.TriggerTypes {
		handler := tr.getTriggerHandler(triggerType)
		if err := handler.Update(ctx, triggerType, tasks); err != nil {
			errTotal = fmt.Errorf("%v : triggerType : %v, err: %v", errTotal, triggerType, err)
		}
	}

	if errTotal != nil {
		return fmt.Errorf("updateNewEndpoints: erorr during update endpoints. err : %v", errTotal)
	}

	return nil
}

func (tr *Service) fillStatusCount(instance models.TaskInstance) map[statuses.TaskInstanceStatus]int {
	var sc = make(map[statuses.TaskInstanceStatus]int)
	for _, stat := range instance.Statuses {
		sc[stat]++
	}
	statusesCount := sc
	return statusesCount
}

func (tr *Service) getOverallStatus(instance models.TaskInstance) (status statuses.TaskInstanceStatus) {

	statusesCount := tr.fillStatusCount(instance)

	if len(instance.Statuses) < 1 {
		return statuses.TaskInstanceScheduled
	}

	deviceCount := len(instance.Statuses)

	if statusesCount[statuses.TaskInstanceDisabled] == deviceCount {
		return statuses.TaskInstanceDisabled
	}

	if statusesCount[statuses.TaskInstancePostponed] == deviceCount {
		return statuses.TaskInstancePostponed
	}

	if statusesCount[statuses.TaskInstanceDisabled] == deviceCount {
		return statuses.TaskInstanceDisabled
	}
	if statusesCount[statuses.TaskInstanceCanceled] == deviceCount {
		return statuses.TaskInstanceCanceled
	}

	return statuses.TaskInstanceScheduled

}
