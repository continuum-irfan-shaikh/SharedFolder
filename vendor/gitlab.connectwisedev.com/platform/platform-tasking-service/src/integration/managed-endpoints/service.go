package managedEndpoints

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	taskCounterCassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/task-counter-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

const (
	dgDelimiter   = "%20OR%20"
	siteDelimiter = "%2C"

	dgPathPattern    = "%s/partners/%s/dynamic-groups/managed-endpoints/set?expression=%s"
	assetPathPattern = "%s/partner/%s/sites/%s/summary"
)

// ManagedEndpoint describes a member of DynamicGroup or Site
type ManagedEndpoint struct {
	// ManagedEndpointID of DynamicGroup
	ID gocql.UUID `json:"id"`
	// SiteID of DynamicGroup
	SiteID string `json:"site"`
	// ManagedEndpointID of Site
	EndpointID gocql.UUID `json:"endpointID"`
}

// GetManagedEndpointsFromTargets returns managedEndpointIDs based on the targets.Type
func GetManagedEndpointsFromTargets(ctx context.Context, task models.Task, httpClient integration.HTTPClient) (managedEndpointIDsMap map[gocql.UUID]models.TargetType, err error) {
	var managedEndpoints []ManagedEndpoint
	managedEndpointIDsMap = make(map[gocql.UUID]models.TargetType)

	if len(task.Targets.IDs) != 0 {
		if task.TargetsByType == nil {
			task.TargetsByType = make(models.TargetsByType)
		}
		task.TargetsByType[task.Targets.Type] = task.Targets.IDs
	}

	if err = task.TargetsByType.Validate(); err != nil {
		return
	}

	if targets, ok := task.TargetsByType[models.DynamicGroup]; ok && len(targets) > 0 {
		if err = getDynamicGroupsData(ctx, task, targets, managedEndpoints, httpClient, managedEndpointIDsMap); err != nil {
			return
		}
	}

	if targets, ok := task.TargetsByType[models.Site]; ok && len(targets) > 0 {
		if err = getSitesData(ctx, task, targets, models.Site, httpClient, managedEndpointIDsMap); err != nil {
			return
		}
	}

	if targets, ok := task.TargetsByType[models.DynamicSite]; ok && len(targets) > 0 {
		if err = getSitesData(ctx, task, targets, models.DynamicSite, httpClient, managedEndpointIDsMap); err != nil {
			return
		}
	}
	return managedEndpointIDsMap, nil
}

func getSitesData(ctx context.Context, task models.Task, targets []string, siteType models.TargetType, httpClient integration.HTTPClient, managedEndpointIDsMap map[gocql.UUID]models.TargetType) (err error) {
	var managedEndpoints []ManagedEndpoint
	url := fmt.Sprintf(assetPathPattern, config.Config.AssetMsURL, task.PartnerID, strings.Join(targets, siteDelimiter))
	if err = integration.GetDataByURL(ctx, &managedEndpoints, httpClient, url, "", true); err != nil {
		return
	}

	for _, managedEndpoint := range managedEndpoints {
		managedEndpointIDsMap[managedEndpoint.EndpointID] = siteType
	}
	return
}

func getDynamicGroupsData(ctx context.Context, task models.Task, targets []string, managedEndpoints []ManagedEndpoint, httpClient integration.HTTPClient, managedEndpointIDsMap map[gocql.UUID]models.TargetType) (err error) {
	url := fmt.Sprintf(dgPathPattern, config.Config.DynamicGroupsMsURL, task.PartnerID, strings.Join(targets, dgDelimiter))
	if err = integration.GetDataByURL(ctx, &managedEndpoints, httpClient, url, "", true); err != nil {
		return
	}

	userSites, err := models.UserSitesPersistenceInstance.Sites(context.Background(), task.PartnerID, task.CreatedBy)
	if err != nil {
		return fmt.Errorf("error while getting user sites fron Cassandra, err: %v", err)
	}

	for _, managedEndpoint := range managedEndpoints {
		siteID, err := strconv.Atoi(managedEndpoint.SiteID)
		if err != nil {
			return fmt.Errorf("siteID(%s) is not an integer, err = %v", managedEndpoint.SiteID, err)
		}

		if containsSiteID(userSites.SiteIDs, int64(siteID)) {
			managedEndpointIDsMap[managedEndpoint.ID] = models.DynamicGroup
		}
	}
	return
}

func containsSiteID(siteIDs []int64, siteID int64) bool {
	for _, id := range siteIDs {
		if id == siteID {
			return true
		}
	}
	return false
}

// GetDifference receives old tasksToRunGroup, disabledTasksGroup and newManagedEndpoints list
// It detects deleted managed endpoints from the list and updates their internal Tasks to Inactive state
// It detects added managed endpoints to the list and creates internal Tasks for them
// It calculates oldActiveTasks which is tasksToRunGroup / (deprecatedTasks U disabledTasksGroup)
// tasksToRunGroup is expected to be not empty slice
func GetDifference(ctx context.Context, tasksToRunGroup, disabledTasksGroup []models.Task, newManagedEndpointsMap, activeTasksEndpointsMap map[gocql.UUID]struct{}) (
	newTasks []models.Task,
	deprecatedTasks []models.Task,
	oldActiveTasks []models.Task,
) {
	taskTemplate := tasksToRunGroup[0]
	oldManagedEndpoints := make([]gocql.UUID, 0, len(tasksToRunGroup)+len(disabledTasksGroup))

	for _, task := range tasksToRunGroup {
		taskExec := task
		oldManagedEndpoints = append(oldManagedEndpoints, taskExec.ManagedEndpointID)
		if _, ok := newManagedEndpointsMap[task.ManagedEndpointID]; ok {
			if task.Schedule.Regularity == apiModels.Recurrent ||
				task.Schedule.Regularity == apiModels.OneTime {
				oldActiveTasks = append(oldActiveTasks, taskExec)
			}
			continue
		}

		// Seems like endpoint has been removed from DG
		taskExec.State = statuses.TaskStateInactive
		deprecatedTasks = append(deprecatedTasks, taskExec)
	}

	// Populate oldManagedEndpoints by managedEndpoints for which the task is disabled
	for _, task := range disabledTasksGroup {
		oldManagedEndpoints = append(oldManagedEndpoints, task.ManagedEndpointID)
	}

	newTasks, err := processDiff(ctx, newManagedEndpointsMap, oldManagedEndpoints, activeTasksEndpointsMap, taskTemplate)
	if err != nil {
		return
	}
	return
}

func processDiff(ctx context.Context, newManagedEndpointsMap map[gocql.UUID]struct{}, oldManagedEndpoints []gocql.UUID, activeTasksEndpointsMap map[gocql.UUID]struct{}, taskTemplate models.Task) (newTasks []models.Task, err error) {
	// newManagedEndpoints / oldManagedEndpoints = newTasks
	for managedEndpoint := range newManagedEndpointsMap {

		if common.UUIDSliceContainsElement(oldManagedEndpoints, managedEndpoint) {
			continue
		}

		if _, ok := activeTasksEndpointsMap[managedEndpoint]; ok {
			continue
		}
		// Seems like endpoint has been added to DG
		// As far as we don't use Machine Local Time for DG tasks
		// we don't need to recalculate RunTime based on ManagedEndpointID Location
		copiedTask := taskTemplate.CopyWithRunTime(managedEndpoint)
		newTasks = append(newTasks, *copiedTask)

		ti, err := models.TaskInstancePersistenceInstance.GetNearestInstanceAfter(copiedTask.ID, copiedTask.LastTaskInstanceID.Time())
		if err != nil {
			go insertNewDeviceIntoInstance(taskTemplate.ManagedEndpointID, managedEndpoint, taskTemplate.LastTaskInstanceID)
		}

		ti.Statuses[managedEndpoint] = ti.Statuses[taskTemplate.ManagedEndpointID]
		err = models.TaskInstancePersistenceInstance.Insert(context.Background(), ti)
		if err != nil {
			logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "GetDifference: error while updating taskInstance: %v", err)
			return newTasks, nil
		}

		go updateCounters(ctx, *copiedTask)
	}
	return newTasks, nil
}

func insertNewDeviceIntoInstance(oldEndpointID, newEndpointID gocql.UUID, lastTaskInstanceID gocql.UUID) {
	ctx := context.Background()
	ti, err := models.TaskInstancePersistenceInstance.GetByIDs(ctx, lastTaskInstanceID)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, "insertNewDeviceIntoInstance: error while updating taskInstance: %v", err)
		return
	}

	if len(ti) == 0 {
		return
	}

	ti[0].Statuses[newEndpointID] = ti[0].Statuses[oldEndpointID]
	err = models.TaskInstancePersistenceInstance.Insert(ctx, ti[0])
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "insertNewDeviceIntoInstance: error while updating taskInstance: %v", err)
		return
	}
}

func updateCounters(ctx context.Context, task models.Task) {
	counters := getCounterForTask(task)
	taskCounterDAO := taskCounterCassandra.New(config.Config.CassandraBatchSize)

	if err := taskCounterDAO.IncreaseCounter(task.PartnerID, counters, false); err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "TaskService.Create: error while trying to increase counter: ", err)
	}
}

func getCounterForTask(task models.Task) []models.TaskCount {
	if task.ExternalTask {
		return []models.TaskCount{}
	}

	counter := models.TaskCount{
		ManagedEndpointID: task.ManagedEndpointID,
		Count:             1,
	}

	counters := make([]models.TaskCount, 0)
	counters = append(counters, counter)
	return counters
}
