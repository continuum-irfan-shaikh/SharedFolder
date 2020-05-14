package tasks

import (
	"context"
	"fmt"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

//GetTasksHistory ...
func (t *Tasks) GetTasksHistory(ctx context.Context, from, to time.Time) ([]entities.ScheduledTasks, error) {
	partnerID, user, err := t.extractScheduledTasksCtx(ctx)
	if err != nil {
		return nil, errorcode.NewBadRequestErr(errorcode.ErrorCantDecodeInputData, err.Error())
	}

	var scheduledTasks = make([]entities.ScheduledTasks, 0)

	instances, err := t.taskInstanceRepo.GetByStartedAtAfter(partnerID, from, to)
	if err != nil {
		return nil, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskInstances,
			fmt.Errorf("GetTasksHistory.GetByStartedAtAfter returned error - %v", err).Error())
	}

	filtered, sortedInstMap := filterByLastRunTime(instances)

	tasks, err := t.legacyRepo.GetByIDs(ctx, t.cache, partnerID, true, t.getTaskIDs(ctx, filtered)...)
	if err != nil {
		return nil, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskByTaskID,
			fmt.Errorf("GetTasksHistory.GetByIDs returned error - %v", err).Error())
	}

	tasksByID := sortTasks(tasks, user.IsNOCAccess)

	for taskID, task := range tasksByID {
		if task.ExternalTask {
			continue
		}

		var lastRunTime = time.Time{}

		instance := sortedInstMap[taskID]
		instance.FillStatusCount()

		overallStatus := instance.CalculateOverallStatus()
		failedCount := instance.StatusesCount[statuses.TaskInstanceFailed]
		successCount := instance.StatusesCount[statuses.TaskInstanceSuccess]
		if task.Type != models.TaskTypeScript {
			failedCount += instance.StatusesCount[statuses.TaskInstanceSomeFailures]
		}

		if successCount != 0 || failedCount != 0 {
			lastRunTime = instance.LastRunTime
		}

		scheduledTasks = append(scheduledTasks, entities.ScheduledTasks{
			ID:            taskID,
			Name:          task.Name,
			OverallStatus: overallStatus,
			LastRunTime:   lastRunTime,
			Description:   task.Description,
			CreatedBy:     task.CreatedBy,
			CreatedAt:     task.CreatedAt,
			ModifiedBy:    task.ModifiedBy,
			ModifiedAt:    task.ModifiedAt,
			ExecutionInfo: entities.ExecutionInfo{
				DeviceCount:  len(instance.Statuses),
				SuccessCount: successCount,
				FailedCount:  failedCount,
			},
			TaskType:           task.Type,
			LastTaskInstanceID: instance.ID,
		})
	}
	return scheduledTasks, err
}

func sortTasks(tasks []models.Task, isNOCUser bool) map[string]models.Task {
	var sortedTask = make(map[string]models.Task)
	for _, t := range tasks {
		if t.IsRequireNOCAccess && !isNOCUser {
			continue
		}
		sortedTask[t.ID.String()] = t
	}
	return sortedTask
}

func filterByLastRunTime(instances []entities.TaskInstance) (keys []entities.TaskInstance, filterMap map[string]entities.TaskInstance) {
	filterMap = make(map[string]entities.TaskInstance, len(instances))

	for _, inst := range instances {
		if _, ok := filterMap[inst.TaskID]; !ok {
			filterMap[inst.TaskID] = inst
			continue
		}

		if filterMap[inst.TaskID].LastRunTime.Before(inst.LastRunTime) {
			filterMap[inst.TaskID] = inst
		}
	}

	keys = make([]entities.TaskInstance, 0, len(filterMap))
	for _, v := range filterMap {
		keys = append(keys, v)
	}
	return
}

func (t *Tasks) getTaskIDs(ctx context.Context, instances []entities.TaskInstance) []gocql.UUID {
	var res = make([]gocql.UUID, 0, len(instances))
	for _, i := range instances {
		taskID, err := gocql.ParseUUID(i.TaskID)
		if err != nil {
			t.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "GetTasksHistory.getTaskIDs: can't parse task ID %v, err -  %s", i.TaskID, err.Error())
			continue
		}
		res = append(res, taskID)
	}
	return res
}
