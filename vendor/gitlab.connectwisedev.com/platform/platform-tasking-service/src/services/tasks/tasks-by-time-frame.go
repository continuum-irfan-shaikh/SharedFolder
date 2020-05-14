package tasks

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// LastTasks returns aggregated data about Task and its Statuses
func (t *TaskService) LastTasks(w http.ResponseWriter, r *http.Request) {
	isScheduledFrame, err := strconv.ParseBool(r.URL.Query().Get("isScheduled"))
	if err != nil { // to pass linter
		isScheduledFrame = false
	}

	var from, to time.Time
	if !isScheduledFrame {
		from, to, err = fetchAndValidateTimeFrame(r)
		if err != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTimeFrameHasBadFormat, "TaskService.LastTasks: cannot fetch time frame. err=%v", err)
			common.SendBadRequest(w, r, errorcode.ErrorTimeFrameHasBadFormat)
			return
		}
		logger.Log.DebugfCtx(r.Context(), "TaskService.LastTasks: got from %v to %v", from, to)
	}

	aggregatedResults, err := t.getLastTasksResults(r, w, isScheduledFrame, from, to)
	if err != nil {
		return
	}
	common.RenderJSON(w, aggregatedResults)
}

func (t *TaskService) getLastTasksResults(r *http.Request, w http.ResponseWriter, isScheduledFrame bool, from, to time.Time) (aggregatedResults []models.TaskDetailsWithStatuses, err error) {
	var (
		currentUser        = t.userService.GetUser(r, t.httpClient)
		ctx                = r.Context()
		taskInstances      []models.TaskInstance
		taskPersistenceErr error
		tasks              []models.Task
	)

	var (
		tasksByIDMap      = make(map[gocql.UUID]models.Task)
		taskInstancesMap  map[gocql.UUID][]models.TaskInstance
		canBePostponedMap = make(map[gocql.UUID]bool)
	)

	if !isScheduledFrame {
		// fetches all instances for the defined period
		taskInstances, err = t.taskInstancePersistence.GetByStartedAtAfter(
			ctx,
			currentUser.PartnerID(),
			from,
			to,
		)
		if err != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances,"TaskService.GetByStartedAtAfter: cannot get TaskInstances for the last 48h. err=%v", err)
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
			return []models.TaskDetailsWithStatuses{}, err
		}

		tasksIDs, filteredTaskInstances := getFilteredRecentTaskIDsAndInstances(taskInstances)
		wg := &sync.WaitGroup{}
		// concurrent calculation BEGIN ----------
		wg.Add(2)
		go func() {
			defer wg.Done()
			tasks, taskPersistenceErr = t.taskPersistence.GetByIDs(
				ctx,
				t.cache,
				currentUser.PartnerID(),
				true,
				tasksIDs...,
			)
			tasksByIDMap = filteredByNOCMap(tasks, currentUser.HasNOCAccess())
		}()

		go func() {
			defer wg.Done()
			taskInstancesMap = models.GroupTaskInstancesByTaskID(filteredTaskInstances)
		}()
		wg.Wait()
		// concurrent calculation END ----------

		if taskPersistenceErr != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "TaskService.LastTasks: cannot get Tasks by TaskID. err=%v", taskPersistenceErr)
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
			return []models.TaskDetailsWithStatuses{}, taskPersistenceErr
		}
	} else {
		// scheduled frame functionality
		tasksByIDMap, taskInstancesMap, canBePostponedMap, err = t.getScheduledTasksData(ctx, currentUser.HasNOCAccess(), currentUser.PartnerID())
		if err != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID,"TaskService.LastTasks: cannot get Tasks by TaskID. err=%v", err)
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
			return
		}
	}
	return t.buildResults(tasksByIDMap, taskInstancesMap, isScheduledFrame, canBePostponedMap), nil
}

func (t *TaskService) buildResults(tasksByIDMap map[gocql.UUID]models.Task, taskInstancesMap map[gocql.UUID][]models.TaskInstance, isScheduledFrame bool, canBePostponedMap map[gocql.UUID]bool) (aggregatedResults []models.TaskDetailsWithStatuses) {
	// getting taskInstance by TaskID -> getting taskInstanceID from taskInstancesMap ->
	// -> getting statusCountsMap by taskInstanceID from statusCountsMap map
	for taskID, task := range tasksByIDMap {
		if task.ExternalTask {
			continue
		}

		for _, instance := range taskInstancesMap[taskID] {
			result := t.buildRecentTaskDetails(instance, task, isScheduledFrame, canBePostponedMap)
			if result == nil {
				continue
			}

			aggregatedResults = append(aggregatedResults, *result)
		}
	}
	return aggregatedResults
}

func (t *TaskService) buildRecentTaskDetails(instance models.TaskInstance, task models.Task, isScheduledFrame bool, canBePostponedMap map[gocql.UUID]bool) *models.TaskDetailsWithStatuses {
	timeNow := time.Now().UTC()
	allStatuses, err := instance.CalculateStatuses()
	if err != nil {
		return nil
	}

	var devicesCount = len(instance.Statuses)
	canBePostponed, canBeCanceled := t.canBePostponedCanceled(isScheduledFrame, canBePostponedMap, task, allStatuses, devicesCount)

	//--------------------------------- filtering results start ---------------------
	// skipping canceled Recurrent Task on recent panel before canceled execution
	skipCanceled := !isScheduledFrame &&
		allStatuses[statuses.TaskInstanceCanceledText] == devicesCount && instance.LastRunTime.IsZero()

	//skipping OneTime task before Canceled execution
	skipCanceledOneTime := !isScheduledFrame &&
		allStatuses[statuses.TaskInstanceCanceledText] == devicesCount &&
		task.Schedule.Regularity == apiModels.OneTime &&
		task.RunTimeUTC.After(timeNow)

	// cancel skipping on recent logic
	if skipCanceled && task.Schedule.Regularity != apiModels.OneTime {
		return nil
	}

	// cancel skipping on recent logic for OneTimeTask
	if skipCanceledOneTime {
		return nil
	}

	successCount := allStatuses[statuses.TaskInstanceSuccessText]
	failedCount := allStatuses[statuses.TaskInstanceFailedText]
	runningCount := allStatuses[statuses.TaskInstanceRunningText]
	disabledCount := allStatuses[statuses.TaskInstanceDisabledText]
	stoppedCount := allStatuses[statuses.TaskInstanceStoppedText]
	postponedCount := allStatuses[statuses.TaskInstancePostponedText]
	scheduledCount := allStatuses[statuses.TaskInstanceScheduledText]

	// last run of a task witch has disabled or stopped devices
	hasFutureStoppedDevices := hasFutureStoppedServices(isScheduledFrame, disabledCount, stoppedCount, failedCount, successCount, runningCount)
	if hasFutureStoppedDevices {
		return nil
	}

	// if in instance all devices finished their running or currently running
	scheduledAndFinished := isScheduledFrame &&
		(successCount+failedCount+runningCount) == len(allStatuses) &&
		len(instance.Statuses) != 0
	if scheduledAndFinished {
		return nil
	}

	// if all devices has postponed status - we won't show them on recent panel
	allPostponed := !isScheduledFrame &&
		(postponedCount == devicesCount || scheduledCount != 0)
	if allPostponed { // postpone 1 + scheduled (localtime)
		return nil
	}
	//--------------------------------- filtering results end ---------------------
	if instance.TriggeredBy == "" && isScheduledFrame {
		instance.TriggeredBy = strings.Join(task.Schedule.TriggerTypes, ",")
	}

	return &models.TaskDetailsWithStatuses{
		Task:           task,
		TaskInstance:   instance,
		Statuses:       allStatuses,
		CanBePostponed: canBePostponed,
		CanBeCanceled:  canBeCanceled,
	}
}

func hasFutureStoppedServices(isScheduledFrame bool, disabledCount int, stoppedCount int, failedCount int, successCount int, runningCount int) bool {
	return isScheduledFrame &&
		(disabledCount != 0 || stoppedCount != 0) &&
		(failedCount != 0 || successCount != 0 ||
			runningCount != 0)
}

func (t *TaskService) canBePostponedCanceled(isScheduledFrame bool, canBePostponedMap map[gocql.UUID]bool, task models.Task, allStatuses map[string]int, devicesCount int) (canBePostponed bool, canBeCanceled bool) {
	if isScheduledFrame {
		// additional checking by task instance statuses
		if canBePostponedMap[task.ID] {
			canBePostponed = t.CanBePostponed(allStatuses)
			if task.IsTrigger() || task.IsTaskAndTriggerNotActivated() {
				// its not possible to postpone triggers and trigger task that is not activated yet
				canBePostponed = false
			}
		}
		canBeCanceled = t.CanBeCanceled(allStatuses, devicesCount)
		// currently its not possible to cancel triggers
		if task.IsTrigger() {
			canBeCanceled = false
		}
	}
	return
}

func (t *TaskService) getHistory(ctx context.Context, currentUser user.User, from time.Time, to time.Time) (map[gocql.UUID]models.Task, map[gocql.UUID][]models.TaskInstance, error) {
	var (
		taskPersistenceErr error
		tasksByIDMap       = make(map[gocql.UUID]models.Task)
		taskInstancesMap   = make(map[gocql.UUID][]models.TaskInstance)
		tasks              []models.Task
	)

	// fetches all instances for the defined period
	taskInstances, err := t.taskInstancePersistence.GetByStartedAtAfter(
		ctx,
		currentUser.PartnerID(),
		from,
		to,
	)
	if err != nil {
		return nil, nil, err
	}

	tasksIDs, filteredTaskInstances := getFilteredRecentTaskIDsAndInstances(taskInstances)

	wg := &sync.WaitGroup{}
	// concurrent calculation BEGIN ----------
	wg.Add(2)

	go func() {
		defer wg.Done()
		tasks, taskPersistenceErr = t.taskPersistence.GetByIDs(
			ctx,
			t.cache,
			currentUser.PartnerID(),
			true,
			tasksIDs...,
		)
		tasksByIDMap = filteredByNOCMap(tasks, currentUser.HasNOCAccess())

	}()

	go func() {
		defer wg.Done()
		taskInstancesMap = models.GroupTaskInstancesByTaskID(filteredTaskInstances)
	}()

	wg.Wait()
	if taskPersistenceErr != nil {
		return nil, nil, taskPersistenceErr
	}
	return tasksByIDMap, taskInstancesMap, nil
}

func (t *TaskService) getScheduledTasksData(
	ctx context.Context,
	hasNOCAccess bool,
	partnerID string) (
	map[gocql.UUID]models.Task,
	map[gocql.UUID][]models.TaskInstance,
	map[gocql.UUID]bool,
	error) {

	//get tasks by runTime starting from current time
	tasks, err := t.taskPersistence.GetByPartnerAndTime(ctx, partnerID, time.Now())
	if err != nil {
		return nil, nil, nil, err
	}

	var tasksByIDMap = make(map[gocql.UUID]models.Task)

	for _, task := range tasks {
		//ignore inactive tasks + NOC check
		if (task.IsRequireNOCAccess && !hasNOCAccess) || (task.State == statuses.TaskStateInactive) {
			continue
		}

		if task.HasPostponedTime() {
			task.RunTimeUTC = task.PostponedRunTime
		}

		if existedTask, ok := tasksByIDMap[task.ID]; (ok && task.RunTimeUTC.Before(existedTask.RunTimeUTC)) || !ok {
			tasksByIDMap[task.ID] = task
		}
	}

	taskInstancesMap, canBePostponedByTaskIDs := t.processTaskData(ctx, tasksByIDMap)
	return tasksByIDMap, taskInstancesMap, canBePostponedByTaskIDs, nil
}

func (t *TaskService) processTaskData(ctx context.Context, tasksByIDMap map[gocql.UUID]models.Task) (map[gocql.UUID][]models.TaskInstance, map[gocql.UUID]bool) {
	taskInstancesMap := make(map[gocql.UUID][]models.TaskInstance)
	canBePostponedByTaskIDs := make(map[gocql.UUID]bool)

	var err error
	for tID, task := range tasksByIDMap {
		var instances []models.TaskInstance
		if task.IsTrigger() {
			// for trigger needed instance is in LTI
			instances, err = t.taskInstancePersistence.GetByIDs(ctx, task.LastTaskInstanceID)
			if err != nil {
				//continue to work with instances we got
				logger.Log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskInstances, "LastTasks.getScheduledTasksData: Error while getting TaskInstances, Err: %v", err)
				continue
			}
		} else {
			instances, err = t.taskInstancePersistence.GetTopInstancesByTaskID(ctx, tID)
			if err != nil {
				//continue to work with instances we got
				logger.Log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskInstances, "LastTasks.getScheduledTasksData: Error while getting top TaskInstances, Err: %v", err)
				continue
			}
		}

		//  empty struct
		if len(instances) == 0 {
			continue
		}

		canBePostponedByTaskIDs[instances[0].TaskID] = t.canBePostponed(instances)
		taskInstancesMap[instances[0].TaskID] = append(taskInstancesMap[instances[0].TaskID], instances[0])
	}
	return taskInstancesMap, canBePostponedByTaskIDs
}

func (t *TaskService) canBePostponed(tis []models.TaskInstance) bool {
	for _, ti := range tis {
		for _, status := range ti.Statuses {
			if status == statuses.TaskInstanceScheduled {
				break
			}

			if status == statuses.TaskInstancePending {
				return false
			}
		}
	}
	return true
}

func fetchAndValidateTimeFrame(r *http.Request) (from, to time.Time, err error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err = time.Parse(time.RFC3339Nano, fromStr)
	if err != nil {
		return
	}

	to, err = time.Parse(time.RFC3339Nano, toStr)
	if err != nil {
		return
	}

	var (
		now                  = time.Now().In(to.Location())
		threeMonthsBeforeNow = now.Add(-3 * (time.Hour * 24 * 30)).Add(-1 * 24 * time.Hour) // TBD
	)

	if from.IsZero() ||
		to.IsZero() ||
		to.Before(from) ||
		to.After(now.Add(24*time.Hour)) || // TBD
		from.Before(threeMonthsBeforeNow) ||
		from.After(now) {

		err = fmt.Errorf("invalid time frame: from [%v] to [%v]", from, to)
		return
	}
	return
}

func getFilteredRecentTaskIDsAndInstances(taskInstances []models.TaskInstance) (taskIDsSlice []gocql.UUID, filteredTaskInstances []models.TaskInstance) {
	taskIDs := make(map[gocql.UUID]bool, len(taskInstances))

	for _, instance := range taskInstances {
		if instance.IsScheduled() {
			continue
		}

		taskIDs[instance.TaskID] = true
		filteredTaskInstances = append(filteredTaskInstances, instance)
	}

	for taskID := range taskIDs {
		taskIDsSlice = append(taskIDsSlice, taskID)
	}

	return
}

// CanBePostponed this function checks if the task can be postponed
func (*TaskService) CanBePostponed(allStatuses map[string]int) (canBePostponed bool) {

	if _, ok := allStatuses[statuses.TaskInstanceCanceledText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceDisabledText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstancePendingText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceStoppedText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceSuccessText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceFailedText]; ok {
		return false
	}
	return true
}

// CanBeCanceled this function checks if the task can be postponed
func (*TaskService) CanBeCanceled(allStatuses map[string]int, devicesCount int) (canBeCanceled bool) {
	if cnt, ok := allStatuses[statuses.TaskInstanceCanceledText]; ok && cnt == devicesCount {
		// we can make cancel all here only if not all devices are canceled
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceDisabledText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstancePendingText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceStoppedText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstancePostponedText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceSuccessText]; ok {
		return false
	}

	if _, ok := allStatuses[statuses.TaskInstanceFailedText]; ok {
		return false
	}
	return true
}
