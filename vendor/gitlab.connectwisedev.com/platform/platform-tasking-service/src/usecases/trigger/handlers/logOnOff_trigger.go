package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
)

const (
	collector     = "winlog"
	description   = "Logon/Logoff Events"
	facility      = 16
	source        = "Microsoft-Windows-Security-Auditing"
	logOnEventID  = "4624"
	logOffEventID = "4634"
	severity      = 5
)

//logOnOffTrigger is a struct for logon trigger
type logOnOffTrigger struct {
	*DefaultTriggerHandler
	agentConf     integration.AgentConfig
	userSites     m.UserSitesPersistence
	profilesRepo  repository.Profiles
	assetsService integration.Asset
}

//NewLogOnOffTrigger returns new logOnOffTrigger handler usecase
func NewLogOnOffTrigger(agentConf integration.AgentConfig, tasksRepo m.TaskPersistence, triggersRepo repository.TriggersRepo, log logger.Logger, us m.UserSitesPersistence, prepo repository.Profiles, assetsService integration.Asset) *logOnOffTrigger {
	return &logOnOffTrigger{
		DefaultTriggerHandler: &DefaultTriggerHandler{
			taskRepo:     tasksRepo,
			triggersRepo: triggersRepo,
			log:          log,
		},
		agentConf:     agentConf,
		userSites:     us,
		profilesRepo:  prepo,
		assetsService: assetsService,
	}
}

// Activate activate logonOff trigger
func (tr *logOnOffTrigger) Activate(ctx context.Context, triggerType string, tasks []m.Task) error {
	endpointsMap, rule, err := tr.buildRule(ctx, triggerType, tasks)
	if err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantDecodeInputData, "buildRule err  %v", err)
		return err
	}

	if len(tasks) < 1 {
		tr.log.WarnfCtx(ctx, "Activate: zero tasks")
		return fmt.Errorf("zero tasks")
	}

	profileID, err := tr.agentConf.Activate(ctx, rule, endpointsMap, tasks[0].PartnerID)
	if err != nil {
		return err
	}
	if err = tr.profilesRepo.Insert(tasks[0].ID, profileID); err != nil {
		err = fmt.Errorf("LogOnOff.Activate Insert profile: %v", err)
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, err.Error())
		return err
	}
	return nil
}

// Deactivate deactivates logonoff trigger
func (tr *logOnOffTrigger) Deactivate(ctx context.Context, triggerType string, tasks []m.Task) error {
	profileID, err := tr.profilesRepo.GetByTaskID(tasks[0].ID)
	if err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, "LogOnOff.Deactivate.GetByTaskID: %v", err)
		return err
	}

	if err = tr.agentConf.Deactivate(ctx, profileID, tasks[0].PartnerID); err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "Deactivate err: %v", err)
		return err
	}

	if err = tr.profilesRepo.Delete(tasks[0].ID); err != nil {
		err = fmt.Errorf("LogOnOff.Deactivate: %v", err)
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantDeleteTask, err.Error())
		return err
	}
	return nil
}

// IsApplicable checks if the trigger is applicable for this task
func (tr *logOnOffTrigger) IsApplicable(ctx context.Context, task m.Task, payload apiModels.TriggerExecutionPayload) bool {
	for _, types := range task.Schedule.TriggerTypes {
		if types == triggers.LoginTrigger || types == triggers.LogoutTrigger {
			return true
		}
	}
	return false
}

func (tr *logOnOffTrigger) buildRule(ctx context.Context, triggerType string, tasks []m.Task) (endpointsMap map[string]entities.Endpoints, rule entities.Rule, err error) {
	if len(tasks) < 1 {
		return nil, entities.Rule{}, errors.New("logOnOffTrigger.buildRule: got nil tasks")
	}

	rule, err = ruleBuilder(triggerType)
	if err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, err.Error())
		return
	}

	managedEndpointsIDs := make([]string, 0)
	for i := range tasks {
		managedEndpointsIDs = append(managedEndpointsIDs, tasks[i].ManagedEndpointID.String())
	}

	endpointsMap, err = tr.getEndpointsSiteIDsMap(ctx, tasks[0].PartnerID, managedEndpointsIDs)
	if err != nil {
		err = fmt.Errorf("LogOnOff.buildRule: %v", err)
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, err.Error())
		return
	}
	return
}

func (tr *logOnOffTrigger) buildUpdatedRule(ctx context.Context, triggerType string, tasks []m.Task) (endpointsMap map[string]entities.Endpoints, rule entities.Rule, err error) {
	rule, err = ruleBuilder(triggerType)
	if err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, err.Error())
		return
	}

	targets := make([]string, 0, len(tasks))
	for _, t := range tasks {
		targets = append(targets, t.ManagedEndpointID.String())
	}

	endpointsMap, err = tr.getEndpointsSiteIDsMap(ctx, tasks[0].PartnerID, targets)
	if err != nil {
		err = fmt.Errorf("LogOnOff.buildRule: %v", err)
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, err.Error())
		return
	}
	return
}

func ruleBuilder(triggerType string) (rule entities.Rule, err error) {
	rule = entities.Rule{
		TriggerID:   triggerType,
		Collector:   collector,
		Description: description,
		EventDetails: entities.RuleEventDetails{
			Facility: facility,
			Source:   source,
			Severity: []int{severity},
		},
	}

	switch triggerType {
	case triggers.LoginTrigger:
		rule.EventDetails.EventIDs = []string{logOnEventID}
	case triggers.LogoutTrigger:
		rule.EventDetails.EventIDs = []string{logOffEventID}
	default:
		return entities.Rule{}, errors.New("wrong type of trigger")
	}
	return rule, nil
}

func (tr *logOnOffTrigger) getEndpointsSiteIDsMap(ctx context.Context, partnerID string, endpointsIDs []string) (map[string]entities.Endpoints, error) {
	filteredEndpoints := make(map[string]entities.Endpoints)
	for _, endpointID := range endpointsIDs {
		meUUID, err := gocql.ParseUUID(endpointID)
		if err != nil {
			tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "logOnOff.getEndpointsSiteIDsMap err = %v", err)
			continue
		}

		siteID, clientID, err := tr.assetsService.GetSiteIDByEndpointID(ctx, partnerID, meUUID)
		if err != nil {
			tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetUserSites, "logOnOff can't get site's by endpoint from Asset %v", err)
			continue
		}
		endpoints := entities.Endpoints{
			PartnerID: partnerID,
			SiteID:    siteID,
			ClientID:  clientID,
		}
		filteredEndpoints[endpointID] = endpoints
	}
	return filteredEndpoints, nil
}

// Update updates logonOff trigger
func (tr *logOnOffTrigger) Update(ctx context.Context, triggerType string, tasks []m.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	endpointsMap, rule, err := tr.buildUpdatedRule(ctx, triggerType, tasks)
	if err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "logOnOffTrigger.Update: build rule for triggerType %s and taskID %v  failed. err : %s  ", triggerType, tasks[0].ID, err.Error())
		return err
	}

	profileID, err := tr.profilesRepo.GetByTaskID(tasks[0].ID)
	if err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskByTaskID, "logOnOffTrigger.Update: can't get task by id %v. err: %s", tasks[0].ID, err.Error())
		return err
	}

	err = tr.agentConf.Update(ctx, rule, endpointsMap, tasks[0].PartnerID, profileID)
	if err != nil {
		tr.DefaultTriggerHandler.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "logOnOffTrigger.Update: can't update rule %v profileID %v for endpoints, partner %v. err: %s", rule, profileID, tasks[0].PartnerID, err.Error())
		return err
	}
	return nil
}
