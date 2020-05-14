package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/utils"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/sites"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
)

const (
	executionTimeFixedDelay = 120 * time.Second
)

// RecurrentTaskUsecases ..
type RecurrentTaskUsecases interface {
	Process(ctx context.Context, currentTime time.Time, tasks []models.Task)
}

type tasksBucket struct {
	tasks    []models.Task
	instance models.TaskInstance
}

// RecurrentTaskProcessor ..
type RecurrentTaskProcessor struct {
	conf                    config.Configuration
	log                     logger.Logger
	triggerUC               trigger.Usecase
	assetCli                integration.Asset
	taskInstanceRepo        TaskInstanceRepo
	targetsRepo             TargetsRepo
	dgRepo                  DynamicGroupRepo
	sitesRepo               SiteRepo
	taskExecutionRepo       TaskExecutionRepo
	executionResultRepo     ExecutionResultRepo
	taskRepo                TaskRepo
	executionExpirationRepo ExecutionExpirationRepo
	cacheRepo               CacheRepo
	assetsClient            AssetsClient
	encryptionService       EncryptionService
	agentEncryptionService  AgentEncryptionService
	httpClient              *http.Client
}

// New creates new service of recurrent task logic executor
func New(conf config.Configuration, log logger.Logger, triggerUC trigger.Usecase,
	taskInstanceRepo TaskInstanceRepo, targetsRepo TargetsRepo, dgRepo DynamicGroupRepo, sitesRepo SiteRepo, taskExecutionRepo TaskExecutionRepo,
	executionResultRepo ExecutionResultRepo, taskRepo TaskRepo, executionExpirationRepo ExecutionExpirationRepo, cacheRepo CacheRepo,
	assetsClient AssetsClient, encryptionService EncryptionService, agentEncryptionService AgentEncryptionService, assetCli integration.Asset, http *http.Client) *RecurrentTaskProcessor {
	return &RecurrentTaskProcessor{
		conf:       conf,
		log:        log,
		httpClient: http,

		triggerUC: triggerUC,

		taskInstanceRepo:        taskInstanceRepo,
		targetsRepo:             targetsRepo,
		dgRepo:                  dgRepo,
		sitesRepo:               sitesRepo,
		taskExecutionRepo:       taskExecutionRepo,
		executionResultRepo:     executionResultRepo,
		taskRepo:                taskRepo,
		executionExpirationRepo: executionExpirationRepo,

		assetCli:               assetCli,
		cacheRepo:              cacheRepo,
		assetsClient:           assetsClient,
		encryptionService:      encryptionService,
		agentEncryptionService: agentEncryptionService,
	}
}

// Process processes all recurrent task that must be processed by current time
func (pr *RecurrentTaskProcessor) Process(ctx context.Context, current time.Time, tasks []models.Task) {
	start := time.Now()

	tasksForRunning, tasksForUpdate := pr.processTriggers(ctx, current, tasks)
	go pr.saveTasks(ctx, tasksForUpdate)

	//processing buckets for next running
	nextBuckets := make(chan *tasksBucket)
	go func(buckets chan *tasksBucket) {
		for b := range buckets {
			err := pr.saveNextBucket(b)
			if err != nil {
				pr.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "can't save next bucket %v. err: %v", b, err)
			}
		}
	}(nextBuckets)

	//processing buckets for current  running
	currentBuckets := make(chan *tasksBucket)
	go func(ctx context.Context, buckets chan *tasksBucket) {
		for b := range buckets {
			go func(instance models.TaskInstance, tasks []models.Task) {
				if err := pr.executeTasks(ctx, instance, tasks); err != nil {
					pr.log.ErrfCtx(ctx, errorcode.ErrorCantExecuteTasks, "Process.executeTasks: %s", err.Error())
					pr.sendExecutionResultsErr(ctx, time.Now().UTC(), err, tasks)
				}
			}(b.instance, b.tasks)

			var err error
			b.tasks, err = pr.recalculateTasks(ctx, b.tasks)
			if err != nil {
				pr.log.ErrfCtx(ctx, errorcode.ErrorCantRecalculate, "Process.recalculateTasks: %s", err)
			}

			if len(b.tasks) == 0 {
				continue
			}

			if err = pr.saveCurrentBucket(b); err != nil {
				pr.log.ErrfCtx(ctx, errorcode.ErrorCantUpdateTask, "can save current bucket %v. err: %v", b, err)
			}

			if err = pr.saveTimeExpiration(tasks[0], b.instance); err != nil {
				pr.log.ErrfCtx(ctx, errorcode.ErrorCantSaveExecutionExpiration, "can save time expiration. err: %v", err)
			}
		}
	}(ctx, currentBuckets)

	tasksByIDs := pr.filterInactiveAndGroupByID(tasksForRunning)
	for _, tasks := range tasksByIDs {
		tasksByTarget := pr.groupByTarget(tasks)
		for targetType, ts := range tasksByTarget {
			currentBucket, nextBucket, err := pr.processInstances(ctx, ts)
			if err != nil {
				pr.log.ErrfCtx(ctx, "Process.processInstances: %s", err.Error())
				continue
			}

			switch targetType {
			case models.DynamicGroup, models.Site, models.DynamicSite:
				currentBucket, nextBucket, err = pr.actualizeDynamicTargets(ctx, currentBucket, nextBucket)
				if err != nil {
					pr.log.WarnfCtx(ctx, "can't actualize dynamic targets. err: %v", err)
				}
			}

			currentBuckets <- currentBucket
			nextBuckets <- nextBucket
		}
	}

	close(currentBuckets)
	close(nextBuckets)
	pr.log.InfofCtx(ctx, "RecurrentTaskProcessor: finished work after %v seconds", time.Since(start).Seconds())
}

func (pr *RecurrentTaskProcessor) processTriggers(ctx context.Context, currentTime time.Time, tasks []models.Task) (tasksToRun, tasksForUpdate []models.Task) {
	tasksActivate := make(map[gocql.UUID][]models.Task)
	tasksDeactivate := make(map[gocql.UUID][]models.Task)

	// grouping tasks by task ID and by process type
	for _, task := range tasks {
		if len(task.Schedule.TriggerTypes) < 1 {
			tasksToRun = append(tasksToRun, task)
			continue
		}

		if task.Schedule.StartRunTime.UTC().Equal(task.RunTimeUTC) {
			tasksActivate[task.ID] = append(tasksActivate[task.ID], task)

			if task.Schedule.StartRunTime.UTC() == task.RunTimeUTC && !task.OriginalNextRunTime.IsZero() {
				task.RunTimeUTC = task.OriginalNextRunTime
				task.OriginalNextRunTime = time.Time{}
				tasksForUpdate = append(tasksForUpdate, task)
				continue
			}

			tasksToRun = append(tasksToRun, task)
			continue
		}

		if task.Schedule.EndRunTime.UTC().Equal(task.RunTimeUTC) {
			tasksDeactivate[task.ID] = append(tasksDeactivate[task.ID], task)
			tasksToRun = append(tasksToRun, task)
			continue
		}
		tasksToRun = append(tasksToRun, task)
	}

	pr.activateTrigger(ctx, tasksActivate)
	pr.deactivateTrigger(ctx, tasksDeactivate)

	return tasksToRun, tasksForUpdate
}

func (pr *RecurrentTaskProcessor) deactivateTrigger(ctx context.Context, tasksDeactivate map[gocql.UUID][]models.Task) {
	for _, groupedTasks := range tasksDeactivate {
		// sending by UCs
		go func(ctx context.Context, tasks []models.Task) {
			if err := pr.triggerUC.Deactivate(ctx, tasks); err != nil {
				pr.log.WarnfCtx(ctx, "Process: deactivate err for taskID %v", tasks[0].ID)
				return
			}
		}(ctx, groupedTasks)
	}
}

func (pr *RecurrentTaskProcessor) activateTrigger(ctx context.Context, tasksActivate map[gocql.UUID][]models.Task) {
	for _, groupedTasks := range tasksActivate {
		// sending by UCs
		go func(ctx context.Context, tasks []models.Task) {
			if err := pr.triggerUC.Activate(ctx, tasks); err != nil {
				pr.log.WarnfCtx(ctx, "Process: activate err for taskID %v", tasks[0].ID)
				return
			}
		}(ctx, groupedTasks)
	}
}

func (pr *RecurrentTaskProcessor) processInstances(ctx context.Context, tasks []models.Task) (current *tasksBucket, next *tasksBucket, err error) {
	if len(tasks) == 0 {
		return nil, nil, fmt.Errorf("processInstances failed because of empty tasks list")
	}

	//because here we have tasks of particular ID and targetType. TaskInstance for running / next TaskInstance will be the same for all tasks
	t := tasks[0]
	currentInstance, err := pr.getRunningInstance(ctx, t)
	if err != nil {
		return nil, nil, fmt.Errorf("processInstances failed for task %v. can't get instance for running. err: %v", t, err)
	}

	currentInstance = pr.updateInstanceStatuses(ctx, tasks, currentInstance)

	tasks = pr.updateLastTaskInstanceID(tasks, currentInstance.ID)

	nextInst, err := pr.getNextInstance(t, currentInstance)
	if err != nil {
		return nil, nil, fmt.Errorf("processInstances failed for task %v. can't get next instance by current %v. err: %v", t, currentInstance, err)
	}

	for id, status := range nextInst.Statuses {
		if status == statuses.TaskInstancePostponed && currentInstance.Statuses[id] == statuses.TaskInstanceRunning {
			nextInst.Statuses[id] = statuses.TaskInstanceScheduled
		}
	}

	return &tasksBucket{tasks, currentInstance}, &tasksBucket{make([]models.Task, 0), nextInst}, nil
}

// checks if target type for an internal task the same as the given one
func (pr *RecurrentTaskProcessor) isGivenTargetType(ctx context.Context, partnerID string, taskID, endpointID gocql.UUID, targetType models.TargetType, external bool) (isGiven bool, isFound bool) {
	gotTargetType, err := pr.taskRepo.GetTargetTypeByEndpoint(partnerID, taskID, endpointID, external)
	if err != nil {
		if err == gocql.ErrNotFound {
			return false, false
		}
		pr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "isGivenTargetType: can't fetch task by target type cause: %v", err)
		return false, false
	}

	return targetType == gotTargetType, true
}

func (pr *RecurrentTaskProcessor) actualizeDynamicTargets(ctx context.Context, currentBucket *tasksBucket, nextBucket *tasksBucket) (*tasksBucket, *tasksBucket, error) {
	t := currentBucket.tasks[0]

	endpoints, err := pr.getActualEndpoints(ctx, currentBucket.tasks) //this is the actual set of endpoints for concrete target type of concrete task
	if err != nil {
		return currentBucket, nextBucket, fmt.Errorf("getActualEndpoints failed. err: %v", err)
	}

	if endpoints, err = pr.filterEndpointsByResourceType(ctx, t.PartnerID, endpoints, t.ResourceType); err != nil { // TODO pass here ctx from upper lvl
		return currentBucket, nextBucket, fmt.Errorf("getCurrentEndpoints:  can't filter endpoints by resourceType: %v, and ids: %v. cause: %s", t.TargetType, t.TargetsByType[t.TargetType], err.Error())
	}

	//actualizing of next Instance
	different, added, removed := pr.diff(nextBucket.instance.Statuses, endpoints)
	if different {
		for _, e := range added {
			_, isFound := pr.isGivenTargetType(ctx, t.PartnerID, t.ID, e, t.TargetType, t.ExternalTask)
			if isFound { // if we found, that means this endpoint is already presented
				continue
			}

			nextBucket.instance.Statuses[e] = statuses.TaskInstanceScheduled
		}

		toBeRemoved := make([]gocql.UUID, 0)
		for _, e := range removed {
			if isGivenType, _ := pr.isGivenTargetType(ctx, t.PartnerID, t.ID, e, t.TargetType, t.ExternalTask); isGivenType {
				toBeRemoved = append(toBeRemoved, e)
				delete(nextBucket.instance.Statuses, e)
			}
		}
		go pr.removeDeactivatedEndpoints(ctx, nextBucket.instance, toBeRemoved...)
	}

	triggerActivate := make(map[gocql.UUID][]models.Task)
	triggerDeactivate := make(map[gocql.UUID][]models.Task)

	//actualizing of current Instance
	different, added, removed = pr.diff(currentBucket.instance.Statuses, endpoints)
	if different {
		toBeAppended := make(map[gocql.UUID]statuses.TaskInstanceStatus)
		for _, e := range added {
			// if it's already presented - we don't need to create one more internal task with different target type
			if _, isFound := pr.isGivenTargetType(ctx, t.PartnerID, t.ID, e, t.TargetType, t.ExternalTask); isFound {
				continue
			}

			newTask := *(t.CopyWithRunTime(e))
			if len(newTask.Schedule.TriggerTypes) > 0 {
				triggerActivate[newTask.ID] = append(triggerActivate[newTask.ID], newTask)
			}

			if t.Schedule.Location != "" {
				currentBucket.instance.Statuses[e] = statuses.TaskInstanceRunning
				toBeAppended[e] = statuses.TaskInstanceRunning
				currentBucket.tasks = append(currentBucket.tasks, newTask)
				continue
			}

			actualLocation, err := pr.assetsClient.GetLocationByEndpointID(ctx, t.PartnerID, e)
			if err != nil {
				pr.log.WarnfCtx(ctx, "recalculate:  can't get location for managed endpoint id: %v, task.Schedule.Location: %s; cause err:%s", t.ManagedEndpointID, t.Schedule.Location, err.Error())
				currentBucket.instance.Statuses[e] = statuses.TaskInstanceRunning
				toBeAppended[e] = statuses.TaskInstanceRunning
				currentBucket.tasks = append(currentBucket.tasks, newTask)
				continue
			}

			newTask.Schedule.StartRunTime = common.AddLocationToTime(newTask.Schedule.StartRunTime, actualLocation)
			if !newTask.Schedule.EndRunTime.IsZero() {
				newTask.Schedule.EndRunTime = common.AddLocationToTime(newTask.Schedule.EndRunTime, actualLocation)
			}

			next, err := common.CalcFirstNextRunTime(newTask.RunTimeUTC.Add(-time.Minute), newTask.Schedule)
			if err != nil {
				if err.Error() == common.NextRunTimeExceedsEndRunTime {
					newTask.RunTimeUTC = t.Schedule.EndRunTime.UTC()
					newTask.State = statuses.TaskStateInactive
					nextBucket.tasks = append(nextBucket.tasks, newTask)
					continue
				}

				pr.log.WarnfCtx(ctx, "%s; next run time for task %v haven't been calculated. cause: %s", err, newTask.ID, err.Error())
				continue
			}

			if next == newTask.RunTimeUTC {
				currentBucket.instance.Statuses[e] = statuses.TaskInstanceRunning
				toBeAppended[e] = statuses.TaskInstanceRunning
				currentBucket.tasks = append(currentBucket.tasks, newTask)
				continue
			}

			newTask.RunTimeUTC = next
			nextBucket.tasks = append(nextBucket.tasks, newTask) //to save task which in other location
		}

		go pr.appendNewEndpoints(ctx, currentBucket.instance, toBeAppended)

		toBeRemoved := make([]gocql.UUID, 0)
		for _, e := range removed {
			if isGivenType, _ := pr.isGivenTargetType(ctx, t.PartnerID, t.ID, e, t.TargetType, t.ExternalTask); !isGivenType {
				continue
			}

			currentBucket = pr.deactivateTask(currentBucket, e)
			toBeRemoved = append(toBeRemoved, e)

			if len(t.Schedule.TriggerTypes) > 0 {
				triggerDeactivate[t.ID] = append(triggerDeactivate[t.ID], t)
			}
		}
		go pr.removeDeactivatedEndpoints(ctx, currentBucket.instance, toBeRemoved...)

		if len(currentBucket.instance.Statuses) == 0 {
			defaultEID, err := gocql.ParseUUID(models.DefaultEndpointUID)
			if err != nil {
				pr.log.WarnfCtx(ctx, "reopenTasks: can't parse UUID: %v", err)
			}
			t.State = statuses.TaskStateActive
			newTask := *(t.CopyWithRunTime(defaultEID))
			currentBucket.tasks = append(currentBucket.tasks, newTask)
		}

	}

	// activating
	pr.activateTrigger(ctx, triggerActivate)
	// deactivating
	pr.deactivateTrigger(ctx, triggerDeactivate)

	return currentBucket, nextBucket, nil
}

func (pr *RecurrentTaskProcessor) removeDeactivatedEndpoints(ctx context.Context, ti models.TaskInstance, endpoints ...gocql.UUID) {
	if len(endpoints) == 0 {
		return
	}
	if err := pr.taskInstanceRepo.RemoveInactiveEndpoints(ti, endpoints...); err != nil {
		pr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "actualizeDynamicTargets: can't remove inactive targets, cause: %v", err)
	}
}

func (pr *RecurrentTaskProcessor) appendNewEndpoints(ctx context.Context, ti models.TaskInstance, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) {
	if len(endpoints) == 0 {
		return
	}
	if err := pr.taskInstanceRepo.AppendNewEndpoints(ti, endpoints); err != nil {
		pr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "actualizeDynamicTargets: can't append new targets, cause: %v", err)
	}
}

func (pr *RecurrentTaskProcessor) deactivateTask(bucket *tasksBucket, endpointID gocql.UUID) *tasksBucket {
	for i, t := range bucket.tasks {
		if t.ManagedEndpointID == endpointID {
			t.State = statuses.TaskStateInactive
			bucket.tasks[i] = t
		}
	}

	delete(bucket.instance.Statuses, endpointID)

	return bucket
}

func (pr *RecurrentTaskProcessor) groupByTarget(tasks []models.Task) map[models.TargetType][]models.Task {
	tasksByTarget := make(map[models.TargetType][]models.Task)

	for _, task := range tasks {
		tasksByTarget[task.TargetType] = append(tasksByTarget[task.TargetType], task)
	}
	return tasksByTarget
}

func (pr *RecurrentTaskProcessor) filterInactiveAndGroupByID(tasks []models.Task) map[gocql.UUID][]models.Task {
	tasksByID := make(map[gocql.UUID][]models.Task)

	for _, task := range tasks {
		if task.State == statuses.TaskStateInactive {
			continue
		}

		tasksByID[task.ID] = append(tasksByID[task.ID], task)
	}
	return tasksByID
}

func (pr *RecurrentTaskProcessor) getActualEndpoints(ctx context.Context, tasks []models.Task) (endpoints []gocql.UUID, err error) {
	if len(tasks) == 0 {
		return nil, fmt.Errorf("getActualEndpoints failed: empty tasks slice")
	}

	t := tasks[0]
	switch t.TargetType {
	case models.DynamicGroup:
		endpoints, err = pr.dgRepo.GetEndpointsByGroupIDs(ctx, t.TargetsByType[models.DynamicGroup], t.CreatedBy, t.PartnerID, t.IsRequireNOCAccess)
	case models.Site:
		endpoints, err = pr.sitesRepo.GetEndpointsBySiteIDs(ctx, t.PartnerID, t.TargetsByType[models.Site])
	case models.DynamicSite:
		siteIDs, err := sites.GetSiteIDs(ctx, pr.httpClient, t.PartnerID, config.Config.SitesMsURL, "")
		if err != nil {
			return nil, fmt.Errorf("can't get sites By Partner %v, err: %v", t.PartnerID, err)
		}
		endpoints, err = pr.sitesRepo.GetEndpointsBySiteIDs(ctx, t.PartnerID, utils.Int64SliceToStringSlice(siteIDs))
	}

	if err != nil {
		return nil, fmt.Errorf("getActualEndpoints failed: can get actual endpoints for task %v. err: %v", t, err)
	}
	return
}

func (pr *RecurrentTaskProcessor) filterEndpointsByResourceType(ctx context.Context, partnerID string, endpoints []gocql.UUID, resType integration.ResourceType) (filteredEndpoints []gocql.UUID, err error) {
	if resType.IsAllResources() {
		return endpoints, nil
	}

	for _, e := range endpoints {
		gotType, err := pr.assetCli.GetResourceTypeByEndpointID(ctx, partnerID, e)
		if err != nil {
			return endpoints, err
		}

		if gotType == resType {
			filteredEndpoints = append(filteredEndpoints, e)
		}
	}
	return
}

func (pr *RecurrentTaskProcessor) getNextInstance(t models.Task, inst models.TaskInstance) (nextInst models.TaskInstance, err error) {
	instanceForCreation := inst
	for {
		nextInst, err = pr.taskInstanceRepo.GetNearestInstanceAfter(t.ID, inst.StartedAt)
		if err == gocql.ErrNotFound {
			nextInst, err = pr.createNewInstance(t, instanceForCreation)
			return
		}

		if err != nil {
			return
		}

		// if this is trigger + reccurent task it might have already executed instances ahead. that's why we need to find the nearest scheduled
		// if it exists
		if nextInst.TriggeredBy == "" {
			return
		}
		inst = nextInst
	}
}

func (pr *RecurrentTaskProcessor) executeTasks(ctx context.Context, instance models.TaskInstance, tasks []models.Task) error {
	pattern := "%s/partners/%s/task-execution-results/task-instances/%s"
	endpoints := make([]apiModels.ManagedEndpoint, 0, len(tasks))

	for _, t := range tasks {
		if !pr.shouldBeExecuted(t) || instance.Statuses[t.ManagedEndpointID] != statuses.TaskInstanceRunning {
			continue
		}

		endpoints = append(endpoints, apiModels.ManagedEndpoint{
			ID:          t.ManagedEndpointID.String(),
			NextRunTime: t.RunTimeUTC,
		})
	}

	if len(endpoints) == 0 {
		return nil
	}

	taskTemplate := tasks[0]
	webHookURL := fmt.Sprintf(pattern, pr.conf.TaskingMsURL, taskTemplate.PartnerID, instance.ID)
	payload := apiModels.ExecutionPayload{
		ExecutionID:              instance.ID.String(),
		OriginID:                 instance.OriginID.String(),
		ManagedEndpoints:         endpoints,
		Parameters:               taskTemplate.Parameters,
		TaskID:                   taskTemplate.ID,
		WebhookURL:               webHookURL,
		Credentials:              taskTemplate.Credentials,
		ExpectedExecutionTimeSec: pr.cacheRepo.CalculateExpectedExecutionTimeSec(context.Background(), taskTemplate),
	}

	if taskTemplate.RunTimeUTC.Add(executionTimeFixedDelay * time.Duration(payload.ExpectedExecutionTimeSec)).Before(time.Now().UTC()) {
		pr.log.WarnfCtx(ctx, "executeTasks: task can't be run as it currently out of schedule. RunTimeUTC: %v ExpectedExecutionTimeSec: %v", taskTemplate.RunTimeUTC, payload.ExpectedExecutionTimeSec)
		return nil
	}

	if taskTemplate.IsRunAsUserApplied() {
		var (
			err      error
			totalErr error
		)

		decrypted, err := pr.encryptionService.Decrypt(*taskTemplate.Credentials)
		if err != nil {
			return fmt.Errorf("executeTasks: could't decrypt credentials for task with ID %v . err: %s", taskTemplate.ID, err.Error())
		}

		taskTemplate.Credentials = &decrypted

		for i := range endpoints {
			payload.ManagedEndpoints = []apiModels.ManagedEndpoint{endpoints[i]}
			endpointID, err := gocql.ParseUUID(endpoints[i].ID)
			if err != nil {
				totalErr = fmt.Errorf("%v err: executeTasks: invalid format of endpoint ID %s. err: %s ", totalErr, endpoints[i].ID, err.Error())
				continue
			}

			if taskTemplate.Credentials != nil {
				encrypted, err := pr.agentEncryptionService.Encrypt(ctx, endpointID, *taskTemplate.Credentials)
				if err != nil {
					totalErr = fmt.Errorf("%v err: executeTasks: could't encrypt credentials for task with ID %v and endpointID %v. err: %s", totalErr, taskTemplate.ID, endpointID, err.Error())
					continue
				}
				payload.Credentials = &encrypted
			}

			if err = pr.taskExecutionRepo.ExecuteTasks(ctx, payload, taskTemplate.PartnerID, taskTemplate.Type); err != nil {
				if totalErr == nil {
					totalErr = fmt.Errorf("executeTasks: error during execution recurrent tasks")
				}
				totalErr = fmt.Errorf("%s err: %s", totalErr, err)
			}
		}

		return totalErr
	}

	return pr.taskExecutionRepo.ExecuteTasks(ctx, payload, taskTemplate.PartnerID, taskTemplate.Type)
}

func (pr *RecurrentTaskProcessor) recalculateTasks(ctx context.Context, tasks []models.Task) ([]models.Task, error) {
	out := make([]models.Task, 0, len(tasks))
	var totalErr error

	for _, t := range tasks {
		var (
			currentRunTimeUTC = t.RunTimeUTC
			actualLocation    *time.Location
			err               error
		)

		if t.HasOriginalRunTime() {
			if currentRunTimeUTC != t.OriginalNextRunTime {
				t.RunTimeUTC = t.OriginalNextRunTime
			}
			t.OriginalNextRunTime = time.Time{}
			out = append(out, t)
			continue
		}

		if t.Schedule.Location != "" {
			actualLocation, err = time.LoadLocation(t.Schedule.Location)
		} else {
			actualLocation, err = pr.assetsClient.GetLocationByEndpointID(context.Background(), t.PartnerID, t.ManagedEndpointID)
		}

		if err != nil {
			pr.log.WarnfCtx(ctx, "recalculateTasks:  can't get location for managed endpoint id: %v, task.Schedule.Location: %s; cause err:%s", t.ManagedEndpointID, t.Schedule.Location, err.Error())
			actualLocation = time.UTC
			err = nil
		}

		next := t.RunTimeUTC
		nowUTC := time.Now().Truncate(time.Minute).UTC()
		for !next.Truncate(time.Minute).UTC().After(nowUTC) {
			next, err = common.CalcNextRunTime(next, t.Schedule, *actualLocation)
			if err != nil && err.Error() == common.NextRunTimeExceedsEndRunTime {
				t.RunTimeUTC = t.Schedule.EndRunTime.UTC()
				t.State = statuses.TaskStateInactive
				out = append(out, t)
				break
			}

			if err != nil {
				t.RunTimeUTC = getUnexpectedRunTime(t.RunTimeUTC, t.Schedule)
				out = append(out, t)
				pr.log.ErrfCtx(ctx, errorcode.ErrorCantRecalculate, "recalculateTasks:  can't recalculate for %v and endp %v, err :%v", t.ID, t.ManagedEndpointID, err)
				totalErr = fmt.Errorf("%v; next run time for task %v haven't been calculated. cause: %s", totalErr, t.ID, err.Error())
				break
			}
		}

		if err != nil {
			continue
		}

		t.RunTimeUTC = next.UTC()
		if t.HasPostponedTime() &&
			(next.After(t.PostponedRunTime) ||
				next == t.PostponedRunTime) {
			t.RunTimeUTC = t.PostponedRunTime
			t.PostponedRunTime = time.Time{}
			t.OriginalNextRunTime = next.UTC()
			out = append(out, t)
			continue
		}
		out = append(out, t)
	}
	return out, totalErr
}

// getUnexpectedRunTime returns new run time if unexpected error in recalculation happens so we dont lose task in the past
func getUnexpectedRunTime(t time.Time, schedule apiModels.Schedule) time.Time {
	switch schedule.Repeat.Frequency {
	case apiModels.Hourly:
		return t.Add(time.Hour * time.Duration(schedule.Repeat.Every))
	case apiModels.Daily:
		return t.Add(time.Hour * time.Duration(schedule.Repeat.Every) * 24)
	case apiModels.Weekly:
		return t.Add(time.Hour * 24 * time.Duration(schedule.Repeat.Every) * 7)
	case apiModels.Monthly:
		return t.Add(time.Hour * 24 * 7 * 30)
	}
	return t.Add(time.Hour * time.Duration(schedule.Repeat.Every))
}

func (pr *RecurrentTaskProcessor) saveTasks(ctx context.Context, tasks []models.Task) {
	for _, task := range tasks {
		if err := pr.taskRepo.UpdateSchedulerFields(ctx, task); err != nil {
			pr.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, err.Error())
		}
	}
}

func (pr *RecurrentTaskProcessor) shouldBeExecuted(t models.Task) bool {
	if t.State == statuses.TaskStateInactive || t.State == statuses.TaskStateDisabled {
		return false
	}
	if t.HasPostponedTime() {
		return false
	}
	return true
}

func (pr *RecurrentTaskProcessor) updateInstanceStatuses(ctx context.Context, tasks []models.Task, inst models.TaskInstance) models.TaskInstance {
	inst.LastRunTime = time.Now().UTC()
	toBeAppended := make(map[gocql.UUID]statuses.TaskInstanceStatus)

	for _, t := range tasks {
		status, _ := inst.Statuses[t.ManagedEndpointID]
		switch status {
		case statuses.TaskInstanceStopped, statuses.TaskInstanceDisabled, statuses.TaskInstanceCanceled:
			continue
		case statuses.TaskInstancePostponed:
			if t.HasPostponedTime() {
				continue
			}
			fallthrough
		default:
			inst.Statuses[t.ManagedEndpointID] = statuses.TaskInstanceRunning
			toBeAppended[t.ManagedEndpointID] = statuses.TaskInstanceRunning
			inst.OverallStatus = statuses.TaskInstanceRunning
		}
	}

	go pr.appendNewEndpoints(ctx, inst, toBeAppended)

	return inst
}

func (pr *RecurrentTaskProcessor) createNewInstance(t models.Task, currInst models.TaskInstance) (models.TaskInstance, error) {
	endpoints := make([]gocql.UUID, 0, len(currInst.Statuses))
	for e := range currInst.Statuses {
		endpoints = append(endpoints, e)
	}

	instanceStatuses := make(map[gocql.UUID]statuses.TaskInstanceStatus)

	for _, e := range endpoints {

		prevEndpointStatus, ok := currInst.Statuses[e]
		if ok {
			switch prevEndpointStatus {
			case statuses.TaskInstanceDisabled, statuses.TaskInstancePostponed:
				instanceStatuses[e] = prevEndpointStatus
			default:
				instanceStatuses[e] = statuses.TaskInstanceScheduled
			}
			continue
		}
		instanceStatuses[e] = statuses.TaskInstanceScheduled
	}

	return models.TaskInstance{
		PartnerID:     t.PartnerID,
		ID:            gocql.TimeUUID(),
		Name:          t.Name,
		TaskID:        t.ID,
		OriginID:      t.OriginID,
		StartedAt:     time.Now().UTC(),
		Statuses:      instanceStatuses,
		OverallStatus: statuses.TaskInstanceScheduled,
	}, nil
}

func (pr *RecurrentTaskProcessor) getRunningInstance(ctx context.Context, t models.Task) (inst models.TaskInstance, err error) {
	inst, err = pr.taskInstanceRepo.GetInstance(t.LastTaskInstanceID)
	if err != nil {
		return inst, fmt.Errorf("can't get instance by LastTaskInstanceID, cause: %s", err.Error())
	}

	var loc *time.Location
	if t.Schedule.Location == "" {
		loc, err = pr.assetCli.GetLocationByEndpointID(ctx, t.PartnerID, t.ManagedEndpointID)
		if err != nil {
			pr.log.WarnfCtx(ctx, "can't load location by task and enpdoint, cause: %s", err.Error())
			loc = t.Schedule.StartRunTime.Location()
		}
	} else {
		loc = t.Schedule.StartRunTime.Location()
	}

	if pr.IsFirstRunning(ctx, t, inst, loc) || pr.IsPostponedNotEqualToScheduleExecution(t) {
		return
	}

	for {
		gotTI, err := pr.taskInstanceRepo.GetNearestInstanceAfter(t.ID, inst.StartedAt)
		if err != nil {
			return gotTI, fmt.Errorf("can't get nearest instance for running. partnerID : %v, started_at : %v, cause: %s", t.PartnerID, inst.StartedAt, err.Error())
		}

		if gotTI.TriggeredBy == "" {
			return gotTI, nil
		}
		inst = gotTI
	}
}

func (pr *RecurrentTaskProcessor) saveCurrentBucket(bucket *tasksBucket) error {
	err := pr.taskRepo.InsertOrUpdate(context.Background(), bucket.tasks...)
	if err != nil {
		return fmt.Errorf("saveBucket failed for tasks %v. err: %v", bucket.tasks, err)
	}

	return nil
}

func (pr *RecurrentTaskProcessor) saveNextBucket(bucket *tasksBucket) error {
	ttl := pr.conf.DataRetentionIntervalDay * expirationTime
	err := pr.taskInstanceRepo.Insert(bucket.instance, ttl)
	if err != nil {
		return fmt.Errorf("saveBucket failed for instance %v . err: %v", bucket.instance, err)
	}

	err = pr.taskRepo.InsertOrUpdate(context.Background(), bucket.tasks...)
	if err != nil {
		return fmt.Errorf("saveBucket failed for tasks %v. err: %v", bucket.tasks, err)
	}

	return nil
}

func (pr *RecurrentTaskProcessor) sendExecutionResultsErr(ctx context.Context, time time.Time, err error, tasks []models.Task) {
	for _, task := range tasks {
		executionResult := apiModels.ExecutionResultKafkaMessage{
			Message: apiModels.ScriptPluginReturnMessage{
				ExecutionID:  task.LastTaskInstanceID.String(),
				TimestampUTC: time,
				Status:       statuses.TaskInstanceFailedText,
				Stderr:       err.Error(),
			},
			BrokerEnvelope: agent.BrokerEnvelope{
				EndpointID: task.ManagedEndpointID.String(),
				PartnerID:  task.PartnerID,
			},
		}

		if errPublish := pr.executionResultRepo.Publish(executionResult); errPublish != nil {
			pr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, err.Error())
		}
	}
}

func (pr *RecurrentTaskProcessor) saveTimeExpiration(task models.Task, ti models.TaskInstance) error {
	meIDs := make([]gocql.UUID, 0)
	hasRunning := false
	for meID, status := range ti.Statuses {
		if status == statuses.TaskInstanceRunning {
			hasRunning = true
			meIDs = append(meIDs, meID)
		}
	}

	if !hasRunning {
		return nil
	}

	timeExpiration := pr.cacheRepo.CalculateExpectedExecutionTimeSec(context.Background(), task)
	timeDuration := time.Duration(timeExpiration+pr.conf.HTTPClientResultsTimeoutSec) * time.Second
	ex := entities.ExecutionExpiration{
		ExpirationTimeUTC:  time.Now().UTC().Add(timeDuration).Truncate(time.Minute),
		PartnerID:          ti.PartnerID,
		TaskInstanceID:     ti.ID,
		ManagedEndpointIDs: meIDs,
	}

	ttl := pr.conf.DataRetentionIntervalDay * expirationTime
	if err := pr.executionExpirationRepo.Insert(ex, ttl); err != nil {
		return fmt.Errorf("saveTimeExpiration failed. err :%v", err)
	}

	return nil
}

// IsFirstRunning says if this is the first running of a task or not
func (pr *RecurrentTaskProcessor) IsFirstRunning(ctx context.Context, task models.Task, instance models.TaskInstance, loc *time.Location) bool {
	var creationTime time.Time
	creationTime = task.CreatedAt.Truncate(time.Minute)
	if !task.ModifiedAt.IsZero() {
		creationTime = task.ModifiedAt.Truncate(time.Minute)
	}

	task.Schedule.StartRunTime = common.AddLocationToTime(task.Schedule.StartRunTime, loc)
	if !task.Schedule.EndRunTime.IsZero() {
		task.Schedule.EndRunTime = common.AddLocationToTime(task.Schedule.EndRunTime, loc)
	}

	nextExpectedRunTime, err := common.CalcFirstNextRunTime(creationTime, task.Schedule)
	if err != nil {
		pr.log.WarnfCtx(ctx, "IsFirstRunning: err during calc firstNextRunTime %v", err)
		nextExpectedRunTime = common.AddLocationToTime(task.Schedule.Repeat.RunTime, task.Schedule.StartRunTime.Location())
	}

	return instance.TriggeredBy == "" &&
		(task.LastTaskInstanceID == gocql.UUID{} || task.Schedule.StartRunTime.UTC() == task.RunTimeUTC || nextExpectedRunTime.UTC() == task.RunTimeUTC)
}

func (pr *RecurrentTaskProcessor) diff(inst map[gocql.UUID]statuses.TaskInstanceStatus, actual []gocql.UUID) (different bool, add []gocql.UUID, rem []gocql.UUID) {
	add = make([]gocql.UUID, 0)
	rem = make([]gocql.UUID, 0)

	intersect := make(map[gocql.UUID]struct{})

	if len(inst) == 0 { // case when instance has zero endpoints
		return len(actual) != 0, actual, rem
	}

	for _, e := range actual {
		if _, ok := inst[e]; !ok {
			add = append(add, e)
			different = true
			continue
		}
		intersect[e] = struct{}{}
	}

	for e := range inst {
		if _, ok := intersect[e]; !ok {
			rem = append(rem, e)
			different = true
		}
	}
	return different, add, rem
}

// IsPostponedNotEqualToScheduleExecution returns if this is postponed execution that must be run in LTI, or in scheduled instance
func (pr *RecurrentTaskProcessor) IsPostponedNotEqualToScheduleExecution(task models.Task) bool {
	isPostponedExecution := !task.HasPostponedTime() && task.HasOriginalRunTime()
	return isPostponedExecution && task.OriginalNextRunTime != task.RunTimeUTC
}

func (pr *RecurrentTaskProcessor) updateLastTaskInstanceID(in []models.Task, instanceID gocql.UUID) (out []models.Task) {
	for _, t := range in {
		t.LastTaskInstanceID = instanceID
		out = append(out, t)
	}
	return
}
