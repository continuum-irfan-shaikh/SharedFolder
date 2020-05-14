package tasks

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/validator"
)

const (
	thirtyMinutesDuration        = time.Duration(time.Minute * 30)
	hourDuration                 = time.Duration(time.Hour)
	threeHoursDuration           = time.Duration(time.Hour * 3)
	dayDuration                  = time.Duration(time.Hour * 24)
	oneHundredEightyDaysDuration = time.Duration(time.Hour * 24 * 180)
	postponeTemplate             = "TaskService.PostponeNearestExecution: "
	postponeDeviceTemplate       = "TaskService.PostponeDeviceNearestExecution: "
	errTaskIsInactive            = "task with ID %v is inactive"
	errDeviceIDBadFormat         = "managedEndpointID %v has bad format. err: %v"
	errZeroTasksByDevice         = "got zero internal tasks for Task ID %v and device ID %v"
	errCantUpdateTask            = "can't update pending task with ID %v with status %v"
	errCantUpdateTaskDevice      = "can't update pending task with ID %v and deviceID %v with status %v"
)

type postponeBody struct {
	DurationString string `json:"duration"`
}

// PostponeDeviceNearestExecution is a handler func to postpone single device in task
func (t *TaskService) PostponeDeviceNearestExecution(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		emptyUUID     gocql.UUID
		taskIDParam   = mux.Vars(r)["taskID"]
		deviceIDParam = mux.Vars(r)["managedEndpointID"]
		partnerID     = mux.Vars(r)["partnerID"]
		currentUser   = t.userService.GetUser(r, t.httpClient)
	)

	deviceID, err := gocql.ParseUUID(deviceIDParam)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorEndpointIDHasBadFormat, postponeDeviceTemplate+errDeviceIDBadFormat, deviceIDParam, err)
		common.SendBadRequest(w, r, errorcode.ErrorEndpointIDHasBadFormat)
		return
	}

	taskID, err := gocql.ParseUUID(taskIDParam)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, postponeTemplate+errTaskIDBadFormat, taskIDParam, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	durationTimeToPostpone, err := t.validatePostponeInput(r)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTimeFrameHasBadFormat, postponeTemplate+errTaskIDBadFormat, taskIDParam, err)
		common.SendBadRequest(w, r, errorcode.ErrorTimeFrameHasBadFormat)
		return
	}

	internalTasks, err := t.taskPersistence.GetByIDAndManagedEndpoints(ctx, partnerID, taskID, deviceID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, postponeDeviceTemplate+errCantGetTasks, taskID, err)
		switch err.(type) {
		case models.TaskNotFoundError:
			common.SendNotFound(w, r, errorcode.ErrorCantGetTaskByTaskID)
		default:
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		}
		return
	}

	if len(internalTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, postponeDeviceTemplate+errZeroTasksByDevice, taskID, deviceID)
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if internalTasks[0].IsTrigger() || internalTasks[0].IsTaskAndTriggerNotActivated() {
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	if currentUser.HasNOCAccess() != internalTasks[0].IsRequireNOCAccess {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeDeviceTemplate+errWrongNOC, taskIDParam)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	task, err := t.processPostponeDeviceTasks(ctx, processPostponedDeviceDTO{
		currentUser:            currentUser,
		internalTasks:          internalTasks,
		durationTimeToPostpone: durationTimeToPostpone,
		deviceID:               deviceID,
	}, r, w)
	if err != nil {
		return
	}

	go t.SendTaskUpdateEventToKafka(ctx, task.ID, partnerID)
	common.RenderJSON(w, task)
}

type processPostponedDeviceDTO struct {
	currentUser            user.User
	internalTasks          []models.Task
	durationTimeToPostpone time.Duration
	deviceID               gocql.UUID
}

func (t *TaskService) processPostponeDeviceTasks(ctx context.Context, o processPostponedDeviceDTO, r *http.Request, w http.ResponseWriter) (task models.Task, err error) {
	taskID := o.internalTasks[0].ID
	taskInstances, err := t.taskInstancePersistence.GetTopInstancesByTaskID(ctx, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, postponeTemplate+errCantGetTopInstances, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	if len(taskInstances) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, postponeDeviceTemplate+errZeroInstances, taskID)
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskInstances)
		return models.Task{}, fmt.Errorf("no instances")
	}

	if len(taskInstances) > 1 {
		for _, s := range taskInstances[1].Statuses {
			if s == statuses.TaskInstancePending {
				logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeTemplate+errCantUpdateTask, taskID.String(), statuses.TaskInstancePending)
				common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
				return models.Task{}, fmt.Errorf("has pending")
			}
		}
	}

	updatedTasks, err := t.postponeTasksByDuration(o.currentUser.UID(), o.internalTasks, o.durationTimeToPostpone)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeDeviceTemplate+errCommon, taskID.String(), err)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	if len(updatedTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeDeviceTemplate+errTaskIsInactive, taskID.String())
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return models.Task{}, fmt.Errorf("no tasks")
	}

	switch taskInstances[0].Statuses[o.deviceID] {
	case statuses.TaskInstancePending, statuses.TaskInstanceCanceled, statuses.TaskInstanceRunning:
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeDeviceTemplate+errCantUpdateTaskDevice, taskID.String(), o.deviceID, taskInstances[0].Statuses[o.deviceID])
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return models.Task{}, fmt.Errorf("has %v status", taskInstances[0].Statuses[o.deviceID])
	}

	taskInstances[0].Statuses[o.deviceID] = statuses.TaskInstancePostponed
	if err = t.taskPersistence.InsertOrUpdate(ctx, updatedTasks...); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskToDB, postponeDeviceTemplate+errCommon, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskToDB)
		return
	}

	//updating taskInstance
	if err = t.taskInstancePersistence.Insert(ctx, taskInstances[0]); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTaskInstances, postponeDeviceTemplate+errCommon, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantUpdateTaskInstances)
		return
	}
	return updatedTasks[0], nil
}

// PostponeNearestExecution is a handler function to postpone nearest execution time of a task
func (t *TaskService) PostponeNearestExecution(w http.ResponseWriter, r *http.Request) {
	var (
		ctx          = r.Context()
		emptyUUID    gocql.UUID
		updatedTasks []models.Task
		partnerID    = mux.Vars(r)["partnerID"]
		taskIDParam  = mux.Vars(r)["taskID"]
		currentUser  = t.userService.GetUser(r, t.httpClient)
	)

	taskID, err := gocql.ParseUUID(taskIDParam)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, postponeTemplate+errTaskIDBadFormat, taskIDParam, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	durationTimeToPostpone, err := t.validatePostponeInput(r)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTimeFrameHasBadFormat, postponeTemplate+errTaskIDBadFormat, taskIDParam, err)
		common.SendBadRequest(w, r, errorcode.ErrorTimeFrameHasBadFormat)
		return
	}

	internalTasks, err := t.taskPersistence.GetByIDs(ctx, nil, partnerID, false, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, postponeTemplate+errCommon, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if len(internalTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, postponeTemplate+errZeroTasks, taskID)
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if internalTasks[0].IsTrigger() || internalTasks[0].IsTaskAndTriggerNotActivated() {
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	if currentUser.HasNOCAccess() != internalTasks[0].IsRequireNOCAccess {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeTemplate+errWrongNOC, taskIDParam)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	updatedTasks, err = t.postponeTasksByDuration(currentUser.UID(), internalTasks, durationTimeToPostpone)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeTemplate+errCommon, taskIDParam, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	if len(updatedTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeTemplate+errTaskIsInactive, taskIDParam)
		common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
		return
	}

	taskInstances, err := t.taskInstancePersistence.GetTopInstancesByTaskID(ctx, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, postponeTemplate+errCantGetTopInstances, taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	if err := t.processPostponedTasks(ctx, taskInstances, r, w, updatedTasks); err != nil {
		return
	}

	go t.SendTaskUpdateEventToKafka(ctx, updatedTasks[0].ID, partnerID)
	common.RenderJSON(w, updatedTasks[0])
}

func (t *TaskService) processPostponedTasks(ctx context.Context, taskInstances []models.TaskInstance, r *http.Request, w http.ResponseWriter, updatedTasks []models.Task) error {
	taskID := updatedTasks[0].ID
	if len(taskInstances) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, postponeDeviceTemplate+errZeroInstances, taskID)
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskInstances)
		return fmt.Errorf(errorcode.ErrorCantGetTaskInstances)
	}

	if len(taskInstances) > 1 {
		for _, s := range taskInstances[1].Statuses {
			if s == statuses.TaskInstancePending {
				logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeTemplate+errCantUpdateTask, taskID.String(), statuses.TaskInstancePending)
				common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
				return fmt.Errorf(errorcode.ErrorCantUpdateTask)
			}
		}
	}

	for deviceID := range taskInstances[0].Statuses {
		//additional check for OneTime Tasks
		switch taskInstances[0].Statuses[deviceID] {
		case statuses.TaskInstancePending, statuses.TaskInstanceCanceled, statuses.TaskInstanceRunning:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTask, postponeTemplate+errCantUpdateTaskDevice, taskID.String(), taskInstances[0].Statuses[deviceID])
			common.SendBadRequest(w, r, errorcode.ErrorCantUpdateTask)
			return fmt.Errorf(errorcode.ErrorCantUpdateTask)
		}
		taskInstances[0].Statuses[deviceID] = statuses.TaskInstancePostponed
	}

	if err := t.taskInstancePersistence.Insert(ctx, taskInstances[0]); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTaskInstances, postponeTemplate+errCommon, taskID.String(), err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantUpdateTaskInstances)
		return err
	}

	if err := t.taskPersistence.InsertOrUpdate(ctx, updatedTasks...); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskToDB, postponeTemplate+errCommon, taskID.String(), err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskToDB)
		return err
	}
	return nil
}

func (t *TaskService) validatePostponeInput(r *http.Request) (durationTimeToPostpone time.Duration, err error) {
	var body postponeBody
	err = validator.ExtractStructFromRequest(r, &body)
	if err != nil {
		return durationTimeToPostpone, fmt.Errorf("cant extract body :%v", err)
	}

	duration, err := strconv.Atoi(body.DurationString)
	if err != nil {
		return durationTimeToPostpone, err
	}

	if duration <= 0 {
		return durationTimeToPostpone, fmt.Errorf("duration is zero")
	}

	durationTimeToPostpone = time.Duration(duration) * time.Minute
	if durationTimeToPostpone > oneHundredEightyDaysDuration {
		return durationTimeToPostpone, fmt.Errorf("more than 180 days")
	}
	return
}

func (t *TaskService) postponeTasksByDuration(modifiedBy string, tasks []models.Task, postponeTime time.Duration) (
	updatedTasks []models.Task, err error) {
	modifiedAt := time.Now().Truncate(time.Minute).UTC()

	for _, task := range tasks {
		if task.State == statuses.TaskStateInactive {
			continue
		}

		switch task.Schedule.Regularity {
		case apiModels.OneTime:
			task = t.getPostponedOneTimeTask(task, postponeTime)
		case apiModels.Recurrent:
			task, err = t.getPostponedRecurrentTask(task, postponeTime)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("wrong Regularity to postpone at task with taskID %v, regularity %v", task.ID, task.Schedule.Regularity)
		}

		task.ModifiedBy = modifiedBy
		task.ModifiedAt = modifiedAt
		if task.Schedule.EndRunTime.Before(task.RunTimeUTC.Add(postponeTime)) && !task.Schedule.EndRunTime.IsZero() {
			task.Schedule.EndRunTime = task.RunTimeUTC.Add(postponeTime)
		}
		updatedTasks = append(updatedTasks, task)
	}
	return
}

func (t *TaskService) getPostponedOneTimeTask(task models.Task, postponeDuration time.Duration) models.Task {
	limitDate, overrides := checkForLimiter(task, postponeDuration, oneHundredEightyDaysDuration)
	if overrides {
		task.RunTimeUTC = limitDate
		return task
	}

	//we'll save only original run time for OneTime
	if !task.HasOriginalRunTime() {
		task.OriginalNextRunTime = task.RunTimeUTC // store original runTime here
	}
	task.RunTimeUTC = task.RunTimeUTC.Add(postponeDuration)
	return task
}

func (t *TaskService) getPostponedRecurrentTask(task models.Task, postponeDuration time.Duration) (postponedTask models.Task, err error) {
	switch task.Schedule.Repeat.Frequency {
	case apiModels.Hourly, apiModels.Daily, apiModels.Weekly, apiModels.Monthly:
		t.updateTaskByPostponedNextRunTime(&task, postponeDuration)
		return task, nil
	default:
		return models.Task{}, fmt.Errorf("wrong repeat frequency: %v for taskID: %v", task.Schedule.Repeat.Frequency, task.ID)
	}
}

// UpdateTaskByPostponedNextRunTime updates task.RunTimeInSeconds and OriginalNextRunTime (if needed) by next cron time for Recurrent tasks
func (*TaskService) updateTaskByPostponedNextRunTime(task *models.Task, postponeDuration time.Duration) {
	if !task.HasPostponedTime() {
		postponedTime := task.RunTimeUTC.Add(postponeDuration).Truncate(time.Minute)
		if !task.OriginalNextRunTime.IsZero() {
			if postponedTime.Before(task.OriginalNextRunTime) {
				task.RunTimeUTC = postponedTime
				return
			}

			task.PostponedRunTime = postponedTime
			task.RunTimeUTC = task.OriginalNextRunTime
			task.OriginalNextRunTime = time.Time{}
		} else {
			task.PostponedRunTime = postponedTime
		}
		return
	}

	limitDate, overrides := checkForLimiter(*task, postponeDuration, oneHundredEightyDaysDuration)
	if overrides {
		task.PostponedRunTime = limitDate
		return
	}
	task.PostponedRunTime = task.PostponedRunTime.Add(postponeDuration).Truncate(time.Minute)
	return
}

// details at RMM-40696 description
func checkForLimiter(task models.Task, postponeTime time.Duration, limiter time.Duration) (limitDate time.Time, overrides bool) {
	// we can't override 180 days with postpone
	if task.Schedule.Regularity == apiModels.Recurrent {
		limitDate = task.RunTimeUTC.Add(limiter)
		return limitDate, task.HasPostponedTime() && task.PostponedRunTime.Add(postponeTime).After(limitDate)
	}

	if !task.HasOriginalRunTime() {
		return time.Time{}, false
	}

	if task.OriginalNextRunTime.After(time.Now()) {
		limitDate = task.OriginalNextRunTime.Add(oneHundredEightyDaysDuration)
		if task.RunTimeUTC.Add(postponeTime).After(limitDate) {
			return limitDate, true
		}
	}

	if task.OriginalNextRunTime.Before(time.Now()) {
		limitDate = time.Now().Add(oneHundredEightyDaysDuration)
		if task.RunTimeUTC.Add(postponeTime).After(limitDate) {
			return limitDate, true
		}
	}
	return
}
