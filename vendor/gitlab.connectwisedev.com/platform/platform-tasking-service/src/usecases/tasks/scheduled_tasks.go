package tasks

import (
	"context"
	"fmt"
	"sort"
	"time"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

const (
	invalidUUIDError = "uuid is invalid: %s"
)

// GetScheduledTasks returns scheduled tasks by given partnerID
func (t *Tasks) GetScheduledTasks(ctx context.Context) (scheduledTasks []e.ScheduledTasks, err error) {
	partnerID, user, err := t.extractScheduledTasksCtx(ctx)
	if err != nil {
		return nil, errorcode.NewBadRequestErr(errorcode.ErrorCantDecodeInputData, err.Error())
	}

	tasks, err := t.tasksRepo.GetScheduledTasks(partnerID)
	if err != nil {
		return nil, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskByTaskID,
			fmt.Errorf("GetScheduledTasks returned error %v", err).Error())
	}
	return t.getFilteredScheduledTaskData(ctx, tasks, user), nil
}

func (t *Tasks) extractScheduledTasksCtx(ctx context.Context) (partnerID string, us e.User, err error) {
	us, ok := ctx.Value(config.UserKeyCTX).(e.User)
	if !ok {
		err = fmt.Errorf("can't get current user from ctx. ctx: %v", ctx)
		return
	}

	partnerID, ok = ctx.Value(config.PartnerIDKeyCTX).(string)
	if !ok {
		err = errors.Errorf(contextParameterError, "partnerID")
		return
	}
	return
}

func (t *Tasks) getFilteredScheduledTaskData(ctx context.Context, internalTasks []e.ScheduledTasks, user e.User) []e.ScheduledTasks {
	tasksByIDMap, canBeCanceled := t.filter(internalTasks, user.IsNOCAccess)
	tisMap, canBeCanceled := t.getInstancesForTasks(ctx, tasksByIDMap, canBeCanceled)
	return t.buildTasks(tasksByIDMap, tisMap, canBeCanceled)
}

func (t *Tasks) buildTasks(tasksByID map[string]e.ScheduledTasks, instances map[string]e.TaskInstance, canBeCanceledMap map[string]bool) []e.ScheduledTasks {
	tasks := make([]e.ScheduledTasks, 0)

	for taskID, task := range tasksByID {
		var (
			overallStatus statuses.OverallStatus
			lastRunTime   = time.Time{}
			successCount  int
			failedCount   int
			tiID          string
		)

		ti, ok := instances[task.ID]
		if ok {
			overallStatus = ti.CalculateOverallStatus()
			failedCount = ti.StatusesCount[statuses.TaskInstanceFailed]
			successCount = ti.StatusesCount[statuses.TaskInstanceSuccess]
			if successCount != 0 || failedCount != 0 {
				lastRunTime = ti.LastRunTime
			}
			tiID = ti.ID
		} else {
			tiID = task.LastTaskInstanceID
			overallStatus = statuses.OverallNew
		}

		canBeCanceled := canBeCanceledMap[taskID]
		if len(task.TriggerTypes) != 0 {
			canBeCanceled = false
		}

		nextRunTime := time.Time{}
		if task.Regularity != tasking.Trigger {
			nextRunTime = task.RunTimeUTC
		}

		tasks = append(tasks, e.ScheduledTasks{
			ID:                 taskID,
			LastTaskInstanceID: tiID,
			LastRunTime:        lastRunTime,
			Name:               task.Name,
			RunTimeUTC:         nextRunTime,
			Description:        task.Description,
			CreatedBy:          task.CreatedBy,
			CreatedAt:          task.CreatedAt,
			ModifiedBy:         task.ModifiedBy,
			ModifiedAt:         task.ModifiedAt,
			OverallStatus:      overallStatus,
			TaskType:           task.TaskType,
			Regularity:         task.Regularity,
			ExecutionInfo: e.ExecutionInfo{
				DeviceCount:  len(ti.Statuses),
				SuccessCount: successCount,
				FailedCount:  failedCount,
			},
			CanBeCanceled: canBeCanceled,
		})
	}
	// sorting tasks by created At
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.UnixNano() > tasks[j].CreatedAt.UnixNano()
	})
	return tasks
}

func (t *Tasks) filter(tasks []e.ScheduledTasks, isNOCUser bool) (map[string]e.ScheduledTasks, map[string]bool) {
	tasksByIDMap := make(map[string]e.ScheduledTasks)
	canBeCanceledMap := make(map[string]bool)
	for _, task := range tasks {
		//ignore inactive tasks + NOC check
		if (task.IsNOC && !isNOCUser) || task.State == statuses.TaskStateInactive {
			continue
		}

		// getting internal task with the nearest runTime
		if existedTask, ok := tasksByIDMap[task.ID]; (ok && task.RunTimeUTC.Before(existedTask.RunTimeUTC)) || !ok {
			tasksByIDMap[task.ID] = task
		}

		if canBeCanceled, ok := canBeCanceledMap[task.ID]; !ok || canBeCanceled {
			// we're going here only in case task in allowed to be canceled to make additional check
			canBeCanceledMap[task.ID] = task.PostponedTime.IsZero()
		}
	}
	return tasksByIDMap, canBeCanceledMap
}

func (t *Tasks) getInstancesForTasks(ctx context.Context, tasks map[string]e.ScheduledTasks, canBeCanceledMap map[string]bool) (map[string]e.TaskInstance, map[string]bool) {
	instanceIDs := make([]string, 0)
	taskIDs := make([]string, 0)

	for id, task := range tasks {
		if task.LastTaskInstanceID == "" {
			continue
		}

		instanceIDs = append(instanceIDs, task.LastTaskInstanceID)
		if task.Regularity == tasking.OneTime {
			continue
		}

		// we get for all others 2 instances to get canBeCanceled info + trigger info
		taskIDs = append(taskIDs, id)
	}

	var taskInstances []e.TaskInstance
	var err error
	if len(instanceIDs) != 0 {
		taskInstances, err = t.taskInstanceRepo.GetInstancesForScheduled(instanceIDs)
		if err != nil {
			//continue to work with instances we got
			t.log.WarnfCtx(ctx, "getInstancesForTasks: Error while getting TaskInstances, err: %v", err)
		}
	}

	var topTaskInstances []e.TaskInstance
	if len(taskIDs) != 0 {
		topTaskInstances, err = t.taskInstanceRepo.GetTopInstancesForScheduledByTaskIDs(taskIDs)
		if err != nil {
			//continue to work with instances we got
			t.log.WarnfCtx(ctx, "getInstancesForTasks: Error while getting TaskInstances, err: %v", err)
		}
	}
	return t.mergeLastInstances(taskInstances, topTaskInstances, canBeCanceledMap)
}

func (t *Tasks) mergeLastInstances(regular []e.TaskInstance, topInstances []e.TaskInstance, canBeCanceledMap map[string]bool) (map[string]e.TaskInstance, map[string]bool) {
	tisMap := make(map[string]e.TaskInstance)
	for _, inst := range regular {
		if canBeCanceled, ok := canBeCanceledMap[inst.TaskID]; !ok || canBeCanceled {
			// we're going here only in case task in allowed to be canceled to make additional check
			canBeCanceledMap[inst.TaskID] = t.canBeCanceledByInstance(inst)
		}
		tisMap[inst.TaskID] = inst
	}

	for _, inst := range topInstances {
		if lastTi, ok := tisMap[inst.TaskID]; ok && inst.LastRunTime.Before(lastTi.LastRunTime) {
			if inst.LastRunTime.IsZero() {
				canBeCanceledMap[inst.TaskID] = t.canBeCanceledByInstance(inst)
			}
			continue
		}

		if canBeCanceled, ok := canBeCanceledMap[inst.TaskID]; !ok || canBeCanceled {
			// we're going here only in case task in allowed to be canceled to make additional check
			canBeCanceledMap[inst.TaskID] = t.canBeCanceledByInstance(inst)
		}
		tisMap[inst.TaskID] = inst
	}
	return tisMap, canBeCanceledMap
}

// returns true if there are no stopped or pending tasks
func (t *Tasks) canBeCanceledByInstance(inst e.TaskInstance) bool {
	if len(inst.Statuses) < 1 {
		return false
	}
	statusesCounts, err := inst.CalculateStatuses()
	if err != nil {
		return false
	}

	if statusesCounts[statuses.TaskInstancePendingText] != 0 || statusesCounts[statuses.TaskInstanceStoppedText] != 0 ||
		statusesCounts[statuses.TaskInstanceCanceledText] == len(inst.Statuses) {
		return false
	}
	return true
}

// DeleteScheduledTasks makes scheduled tasks inactive
func (t *Tasks) DeleteScheduledTasks(ctx context.Context, ids e.TaskIDs) (err error) {
	partnerID, ok := ctx.Value(config.PartnerIDKeyCTX).(string)
	if !ok {
		return errors.Errorf(contextParameterError, config.PartnerIDKeyCTX)
	}

	uuids, err := t.mapToUUIDs(ids)
	if err != nil {
		return err
	}

	tasks, err := t.legacyRepo.GetByIDs(ctx, nil, partnerID, false, uuids...)
	if err != nil {
		return errors.Wrap(err, errorcode.ErrorCantGetTaskByTaskID)
	}

	if len(tasks) == 0 {
		return fmt.Errorf("got zero tasks")
	}
	tasksToDeactivate := make(map[gocql.UUID][]models.Task)

	for i, task := range tasks {
		tasks[i].State = statuses.TaskStateInactive
		tasks[i].RunTimeUTC = time.Now().Truncate(time.Minute)
		task.Schedule.EndRunTime = time.Now()
		if task.Schedule.Regularity == tasking.Trigger || (task.Schedule.Regularity == tasking.Recurrent && len(task.Schedule.TriggerTypes) > 0) {
			tasksToDeactivate[task.ID] = append(tasksToDeactivate[task.ID], task)
		}
	}

	err = t.legacyRepo.UpdateSchedulerFields(ctx, tasks...)
	if err != nil {
		return errors.Wrap(err, errorcode.ErrorCantUpdateTask)
	}

	for _, toDeactivate := range tasksToDeactivate {
		go func(ctx context.Context, tasks []models.Task) {
			if err = t.tr.Deactivate(ctx, tasks); err != nil {
				t.log.WarnfCtx(ctx, "DeleteScheduledTask: error while deactivating task %v", tasks[0].ID)
			}
		}(ctx, toDeactivate)
	}

	return
}

func (t *Tasks) mapToUUIDs(ids e.TaskIDs) (uuids []gocql.UUID, err error) {
	uuids = make([]gocql.UUID, 0)
	for _, id := range ids.IDs {
		uuid, err := gocql.ParseUUID(id)
		if err != nil {
			return uuids, errors.Errorf(invalidUUIDError, id)
		}
		uuids = append(uuids, uuid)
	}
	return
}
