package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/sites"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

const (
	canNotGetInternalTaskErrMsg   = "can't get internal tasks by task ID: %v."
	canNotFoundInternalTaskErrMsg = "can't found internal tasks by task ID: %v."
	canNotGetUserSitesErrMsg      = "can't get user sites by partner ID: %v."
	canNotGetSiteFromCtxErrMsg    = "can't get site from context."
	canNotParseIntErrMsg          = "can't parse siteID to int64. siteID: %v."
)

// FirstCheckInHandler represent FirstCheckIn handler for trigger
type FirstCheckInHandler struct {
	*DefaultTriggerHandler
	sitesRepo m.UserSitesPersistence
}

// NewFirstCheckIn returns FirstCheckInHandler
func NewFirstCheckIn(tRepo m.TaskPersistence, log logger.Logger, sitesRepo m.UserSitesPersistence) *FirstCheckInHandler {
	return &FirstCheckInHandler{
		DefaultTriggerHandler: &DefaultTriggerHandler{
			taskRepo: tRepo,
			log:      log,
		},
		sitesRepo: sitesRepo,
	}
}

// GetTask returns task to process
func (tr *FirstCheckInHandler) GetTask(ctx context.Context, taskID gocql.UUID) (task m.Task, err error) {
	var partnerID = ctx.Value(config.PartnerIDKeyCTX).(string)

	internalTasks, err := tr.taskRepo.GetByIDs(ctx, nil, partnerID, false, taskID)
	if err != nil {
		switch err.(type) {
		case m.TaskNotFoundError:
			return task, errors.Wrapf(err, canNotFoundInternalTaskErrMsg, taskID)
		default:
			tr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, errors.Wrapf(err, canNotGetInternalTaskErrMsg, taskID).Error())
			return task, err
		}
	}

	if len(internalTasks) == 0 {
		return task, errors.Errorf(canNotFoundInternalTaskErrMsg, taskID)
	}

	return internalTasks[0], nil
}

// IsApplicable check whether the trigger is set on Site target
func (tr *FirstCheckInHandler) IsApplicable(ctx context.Context, task m.Task, _ apiModels.TriggerExecutionPayload) bool {
	siteIDctx, ok := ctx.Value(config.SiteIDKeyCTX).(string)
	if !ok {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, canNotGetSiteFromCtxErrMsg)
		return false
	}

	for _, site := range task.TargetsByType[m.Site] {
		if siteIDctx == site {
			return true
		}
	}

	if len(task.TargetsByType[m.DynamicSite]) != 0 {
		id, err  := strconv.Atoi(siteIDctx)
		if err != nil {
			tr.log.WarnfCtx(ctx, errorcode.ErrorCantGetUserSites, "error during convertion site for partner %v, err %v", task.PartnerID, err)
			return false
		}
		siteID := int64(id)

		siteIDs, err := sites.GetSiteIDs(ctx, http.DefaultClient, task.PartnerID, config.Config.SitesMsURL, "")
		if err != nil {
			tr.log.WarnfCtx(ctx, errorcode.ErrorCantGetUserSites, "error during retreiving sites for partner %v, err %v", task.PartnerID, err)
			return false
		}

		found := false
		for _, id := range siteIDs {
			if id == siteID {
				found = true
			}
		}

		if !found {
			return false
		}

		for _, site := range task.TargetsByType[m.DynamicSite] {
			if siteIDctx == site {
				return true
			}
		}
	}
	return false
}
