package tasks

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

func (taskService TaskService) createNewInternalTasks(ctx context.Context, task *models.Task, hasNOCAccess bool) (intTasks []models.Task, siteIDs []string, sErr *serviceError) {
	var (
		err      error
		template models.TemplateDetails
		eIDs     map[gocql.UUID]models.TargetType
	)

	template, sErr = taskService.getTemplateAndValidateParameters(ctx, task, hasNOCAccess)
	if sErr != nil {
		return
	}

	eIDs, siteIDs, err = taskService.getEndpoints(ctx, *task)
	if err != nil {
		sErr = &serviceError{
			errCode: http.StatusInternalServerError,
			errMsg:  errorcode.ErrorCantOpenTargets,
			err:     err,
		}
		return
	}

	if !task.ResourceType.IsAllResources() {
		eIDs, err = taskService.filterByResourceType(ctx, eIDs, task.PartnerID, task.ResourceType)
		if err != nil {
			sErr = &serviceError{
				errCode: http.StatusInternalServerError,
				errMsg:  errorcode.ErrorCantOpenTargets,
				err:     err,
			}
			return
		}
	}

	intTasks, sErr = taskService.createInternalTasks(ctx, eIDs, task, template)
	return
}

func (taskService TaskService) createInternalTasks(ctx context.Context, eIDs map[gocql.UUID]models.TargetType, task *models.Task, template models.TemplateDetails) (intTasks []models.Task, sErr *serviceError) {
	var err error
	timeZonesByEndpoint := make(map[gocql.UUID]*time.Location, len(eIDs))

	if len(eIDs) > 0 {
		timeZonesByEndpoint = getLocationsMap(ctx, *task, eIDs)
	} else {
		if task.Schedule.Regularity == apiModels.RunNow {
			sErr = &serviceError{
				errCode: http.StatusInternalServerError,
				errMsg:  errorcode.ErrorNoEndpointsForTargets,
				err:     errors.New("Unable to create RunNow task without endpoints"),
			}
			return
		}
		defaultEID, err := gocql.ParseUUID(models.DefaultEndpointUID)
		if err != nil {
			logger.Log.WarnfCtx(ctx, "can't parse UUID: %v", err)
		}

		eIDs = make(map[gocql.UUID]models.TargetType)
		for targetType := range task.TargetsByType {
			eIDs[defaultEID] = targetType
		}

		loc := time.UTC
		if len(task.Schedule.Location) > 0 {
			loc, err = time.LoadLocation(task.Schedule.Location)
			if err != nil {
				logger.Log.WarnfCtx(ctx, "can't LoadLocation [%+v]: %v", task.Schedule.Location, err)
				loc = time.UTC
			}
		}
		timeZonesByEndpoint[defaultEID] = loc
	}

	intTasks, err = taskService.buildInternalTasks(task, template, timeZonesByEndpoint, eIDs)
	if err == nil {
		return
	}

	if _, ok := err.(*common.RunTimeInvalidError); ok {
		sErr = &serviceError{
			errCode: http.StatusBadRequest,
			errMsg:  errorcode.ErrorCantCreateNewTask,
			err:     err,
		}
		return
	}

	sErr = &serviceError{
		errCode: http.StatusInternalServerError,
		errMsg:  errorcode.ErrorCantCreateNewTask,
		err:     err,
	}
	return
}

func getLocationsMap(ctx context.Context, task models.Task, endpointIDs map[gocql.UUID]models.TargetType) map[gocql.UUID]*time.Location {
	var (
		timeZoneMapByEndpointID = make(map[gocql.UUID]*time.Location, len(endpointIDs))
		semaphore               = make(chan struct{}, config.Config.ConcurrentRESTCalls)
		wg                      = &sync.WaitGroup{}
		mu                      = &sync.Mutex{}
	)

	if len(task.Schedule.Location) > 0 {
		loc, err := time.LoadLocation(task.Schedule.Location)
		if err != nil {
			logger.Log.WarnfCtx(ctx, "can't LoadLocation [%+v]: %v", task.Schedule.Location, err)
			loc = time.UTC
		}

		for endpointID := range endpointIDs {
			timeZoneMapByEndpointID[endpointID] = loc
		}
		return timeZoneMapByEndpointID
	}

	if task.Schedule.Regularity == apiModels.RunNow {
		for endpointID := range endpointIDs {
			timeZoneMapByEndpointID[endpointID] = time.UTC
		}
		return timeZoneMapByEndpointID
	}

	wg.Add(len(endpointIDs))

	for meID := range endpointIDs {
		managedEndpointID := meID
		semaphore <- struct{}{}

		go func(ctx context.Context, endpointID gocql.UUID) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			loc, err := asset.ServiceInstance.GetLocationByEndpointID(ctx, task.PartnerID, endpointID)
			if err != nil {
				logger.Log.WarnfCtx(ctx, "can't get location for ManagedEndpointID [%v]: %v", endpointID, err)
				loc = time.UTC
			}

			mu.Lock()
			timeZoneMapByEndpointID[endpointID] = loc
			mu.Unlock()
		}(ctx, managedEndpointID)
	}

	wg.Wait()
	close(semaphore)

	return timeZoneMapByEndpointID
}

// buildInternalTasks creates internal Task for each ManagedEndpointID and calculates nextRunTimeUTC
// This function updates some fields of the task
func (taskService TaskService) buildInternalTasks(task *models.Task, template models.TemplateDetails, timeZonesByEndpoint map[gocql.UUID]*time.Location, eIDs map[gocql.UUID]models.TargetType) ([]models.Task, error) {
	var (
		runTimeUTC               time.Time
		schedule                 apiModels.Schedule
		err                      error
		emptyUUID                gocql.UUID
		internalTasks            = make([]models.Task, 0, len(timeZonesByEndpoint))
		managedEndpointsDetailed = make([]models.ManagedEndpointDetailed, 0, len(timeZonesByEndpoint))
		taskName                 = getTaskName(task.Name, template.Name)
	)

	if task.ID == emptyUUID {
		task.ID = gocql.TimeUUID()
	}

	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now().Truncate(time.Millisecond).UTC()
	}

	for endpointID, loc := range timeZonesByEndpoint {
		runTimeUTC, schedule, err = parseTaskTimeInLoc(task.Schedule, loc)
		if err != nil {
			return internalTasks, err
		}

		var originalRunTimeUTC time.Time
		if schedule.Regularity == apiModels.Recurrent && len(task.Schedule.TriggerTypes) != 0 &&
			runTimeUTC.Truncate(time.Minute) != schedule.StartRunTime.UTC().Truncate(time.Minute) {
			originalRunTimeUTC = runTimeUTC
			runTimeUTC = schedule.StartRunTime.UTC()
		}

		managedEndpointsDetailed = append(managedEndpointsDetailed, models.ManagedEndpointDetailed{
			ManagedEndpoint: apiModels.ManagedEndpoint{
				ID:          endpointID.String(),
				NextRunTime: runTimeUTC,
			},
			State: statuses.TaskStateActive,
		})

		internalTasks = append(internalTasks, models.Task{
			ID:                  task.ID,
			Name:                taskName,
			Description:         template.Description,
			CreatedAt:           task.CreatedAt,
			CreatedBy:           task.CreatedBy,
			PartnerID:           task.PartnerID,
			OriginID:            task.OriginID,
			State:               statuses.TaskStateActive,
			Type:                task.Type,
			Parameters:          task.Parameters,
			RunTimeUTC:          runTimeUTC,
			ExternalTask:        task.ExternalTask,
			ResultWebhook:       task.ResultWebhook,
			ManagedEndpointID:   endpointID,
			TargetType:          eIDs[endpointID],
			TargetsByType:       task.TargetsByType,
			ResourceType:        task.ResourceType,
			Schedule:            schedule,
			IsRequireNOCAccess:  task.IsRequireNOCAccess,
			ModifiedBy:          task.ModifiedBy,
			ModifiedAt:          task.ModifiedAt,
			DefinitionID:        task.DefinitionID,
			OriginalNextRunTime: originalRunTimeUTC,
			Credentials:         task.Credentials,
		})
	}

	task.Name = taskName
	task.Description = template.Description
	task.State = statuses.TaskStateActive
	task.ManagedEndpoints = managedEndpointsDetailed
	if len(internalTasks) > 0 {
		internalTasks[0].ManagedEndpoints = managedEndpointsDetailed
	}
	return internalTasks, nil
}

func getTaskName(taskName, scriptName string) string {
	if len(strings.TrimSpace(taskName)) == 0 {
		return scriptName
	}
	return taskName
}

func parseTaskTimeInLoc(initialSchedule apiModels.Schedule, location *time.Location) (runTime time.Time, schedule apiModels.Schedule, err error) {
	schedule = initialSchedule

	switch schedule.Regularity {
	case apiModels.RunNow:
		// End time for offline endpoints
		if !schedule.EndRunTime.IsZero() {
			schedule.EndRunTime = common.AddLocationToTime(schedule.EndRunTime, location)
		}
	case apiModels.OneTime:
		schedule.StartRunTime = common.AddLocationToTime(schedule.StartRunTime, location)

		// validate time
		if !common.ValidFutureTime(schedule.StartRunTime) {
			runTime = schedule.StartRunTime.Add(24 * time.Hour).UTC()
			return
		}
		runTime = schedule.StartRunTime.UTC()
	case apiModels.Recurrent:
		schedule.StartRunTime = common.AddLocationToTime(schedule.StartRunTime, location)
		if !schedule.EndRunTime.IsZero() {
			schedule.EndRunTime = common.AddLocationToTime(schedule.EndRunTime, location)
		}

		runTime, err = common.CalcFirstNextRunTime(time.Now().UTC().Truncate(time.Minute), schedule)
		if err != nil {
			return
		}
	case apiModels.Trigger:
		schedule.StartRunTime = common.AddLocationToTime(schedule.StartRunTime, location)
		if !schedule.EndRunTime.IsZero() {
			schedule.EndRunTime = common.AddLocationToTime(schedule.EndRunTime, location)
		}
		runTime = schedule.StartRunTime.UTC()
	case apiModels.Regularity(0):
		return runTime, initialSchedule, errors.New("wrong Regularity")
	}
	return
}

//getPeriod calculate initial period of task
func getPeriod(frequency apiModels.Frequency, startRunTime time.Time, location *time.Location) (period int) {
	if location != nil && !common.ValidFutureTime(startRunTime) {
		startRunTime = time.Now().In(location).Truncate(time.Minute)
	}

	switch frequency {
	case apiModels.Daily:
		period = startRunTime.YearDay()
	case apiModels.Weekly:
		_, period = startRunTime.ISOWeek()
	case apiModels.Monthly:
		_, month, _ := startRunTime.Date()
		period = int(month)
	}

	return
}
