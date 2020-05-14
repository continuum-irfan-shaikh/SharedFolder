package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	api "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	triggerUseCaseLogPrefix             = "usecases/trigger/handler/dynamic-group-based-triggers.go: "
	canNotGetDGBasedTriggerParamsErrMsg = "can't get dynamic group based trigger parameters from execution payload. err: %v"
	canNotGetSitesByUserErrMsg          = "can't get sites  user : %s has access to. err=%s"
	canNotGetInternalTasksErrMsg        = "can't get internal tasks by task ID: %v, endpointID: %v. err: %v"
)

// DynamicBasedTrigger is a struct that represents DynamicBasedTriggers handling
type DynamicBasedTrigger struct {
	dgClient  integration.DynamicGroups
	tasksRepo m.TaskPersistence
	log       logger.Logger
	cache     persistency.Cache
	sitesRepo m.UserSitesPersistence
}

// PreExecution supposed to be empty
func (tr *DynamicBasedTrigger) PreExecution(task m.Task) error {
	return nil
}

// PostExecution supposed to be empty
func (tr *DynamicBasedTrigger) PostExecution(task m.Task) error {
	return nil
}

// NewDynamicBasedTrigger returns new DynamicBasedTrigger triggers handler usecase
func NewDynamicBasedTrigger(
	dg integration.DynamicGroups,
	sitesRepo m.UserSitesPersistence,
	tasksRepo m.TaskPersistence,
	c persistency.Cache,
	l logger.Logger) *DynamicBasedTrigger {
	return &DynamicBasedTrigger{
		dgClient:  dg,
		log:       l,
		cache:     c,
		tasksRepo: tasksRepo,
		sitesRepo: sitesRepo,
	}
}

// Activate implements activating for DG based triggers
func (tr *DynamicBasedTrigger) Activate(ctx context.Context, _ string, tasks []m.Task) error {
	t := tasks[0]
	if len(t.TargetsByType[m.DynamicGroup]) == 0 {
		return fmt.Errorf("Activate: wrong target type for DG based trigger - targetsMap %v", t.TargetsByType)
	}

	if err := tr.dgClient.StartMonitoringGroups(ctx, t.PartnerID, t.TargetsByType[m.DynamicGroup], t.ID); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorKafka,"DynamicBasedTrigger.Activate: %v", err)
		return err
	}
	return nil
}

// Deactivate implements Deactivating for DG based triggers
func (tr *DynamicBasedTrigger) Deactivate(ctx context.Context, _ string, tasks []m.Task) error {
	t := tasks[0]
	if len(t.TargetsByType[m.DynamicGroup]) == 0 {
		return fmt.Errorf("Deactivate: wrong target type for DG based trigger - targetsMap %v", t.TargetsByType)
	}

	// start monitor
	if err := tr.dgClient.StopGroupsMonitoring(ctx, t.PartnerID, t.TargetsByType[m.DynamicGroup], t.ID); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorKafka,"DynamicBasedTrigger.Deactivate: %v", err)
		return err
	}
	return nil
}

// Update implements Deactivating for DG based triggers
func (tr *DynamicBasedTrigger) Update(ctx context.Context,triggerType string, tasks []m.Task) error {
	return nil
}

// GetTask returns task to process
func (tr *DynamicBasedTrigger) GetTask(ctx context.Context, taskID gocql.UUID) (task m.Task, err error) {
	var (
		partnerID  = ctx.Value(config.PartnerIDKeyCTX).(string)
		endpointID = ctx.Value(config.EndpointIDKeyCTX).(gocql.UUID)
	)

	internalTasks, err := tr.tasksRepo.GetByIDs(ctx, tr.cache, partnerID, true, taskID)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, triggerUseCaseLogPrefix+canNotGetInternalTasksErrMsg, taskID, endpointID, err)
		return task, err
	}
	if len(internalTasks) < 1 {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, triggerUseCaseLogPrefix+canNotGetInternalTasksErrMsg, taskID, endpointID, err)
		return task, fmt.Errorf(triggerUseCaseLogPrefix+canNotGetInternalTasksErrMsg, taskID, endpointID, err)
	}
	return internalTasks[0], nil
}

// IsApplicable checks if the trigger with that DG id is applicable for this task
func (tr *DynamicBasedTrigger) IsApplicable(ctx context.Context, task m.Task, payload api.TriggerExecutionPayload) bool {
	var emptyUUID gocql.UUID
	dynamicGroupID, err := gocql.ParseUUID(payload.DynamicGroupID)
	if err != nil || dynamicGroupID == emptyUUID {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantDecodeInputData, canNotGetDGBasedTriggerParamsErrMsg, err)
		return false
	}

	if !task.TargetsByType.Contains(m.DynamicGroup, dynamicGroupID.String()) {
		return false
	}

	var (
		partnerID = ctx.Value(config.PartnerIDKeyCTX).(string)
		siteID    = ctx.Value(config.SiteIDKeyCTX).(string)
	)

	uSites, err := tr.sitesRepo.Sites(context.Background(), partnerID, task.CreatedBy)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetUserSites, errors.Wrapf(err, canNotGetSitesByUserErrMsg, task.CreatedBy, err.Error()).Error())
		return false
	}

	siteIDint, err := strconv.Atoi(siteID)
	if err != nil {
		tr.log.WarnfCtx(ctx, triggerUseCaseLogPrefix+canNotGetSitesByUserErrMsg, task.CreatedBy, err.Error())
		return false
	}

	siteIDint64 := int64(siteIDint)
	for _, id := range uSites.SiteIDs {
		if id == siteIDint64 {
			return true
		}
	}
	return false
}
