package trigger

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gocql/gocql"
	api "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mockusecases "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-usecases"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	c "gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger/handlers"
)

const hundredYears = time.Duration(time.Hour * 24 * 365 * 100)

// Service represents triggers usecase
type Service struct {
	taskRepo         m.TaskPersistence
	taskInstanceRepo m.TaskInstancePersistence
	sitesRepo        m.UserSitesPersistence
	triggerRepo      repository.TriggersRepo
	dynamicGroups    integration.DynamicGroups
	automationEngine integration.AutomationEngine
	cache            persistency.Cache
	memCache         memcached.Cache
	log              logger.Logger
	agentConf        integration.AgentConfig
	profilesRepo     repository.Profiles
	assetsService    integration.Asset
	httpClient       integration.HTTPClient
}

// New returns new trigger Usecase
func New(c persistency.Cache, mc memcached.Cache, modelsDTO m.DataBaseConnectors, db repository.DatabaseRepositories, log logger.Logger, externalClients integration.ExternalClients) *Service {
	return &Service{
		taskRepo:         modelsDTO.Task,
		taskInstanceRepo: modelsDTO.TaskInstance,
		sitesRepo:        modelsDTO.UserSites,
		triggerRepo:      db.Triggers,
		profilesRepo:     db.Profiles,
		httpClient:       externalClients.HTTP,
		assetsService:    externalClients.Asset,
		agentConf:        externalClients.AgentConfig,
		dynamicGroups:    externalClients.DynamicGroups,
		automationEngine: externalClients.AutomationEngine,
		cache:            c,
		memCache:         mc,
		log:              log,
	}
}

// Activate activates trigger
func (tr *Service) Activate(ctx context.Context, tasks []m.Task) error {
	if len(tasks) == 0 {
		tr.log.WarnfCtx(ctx, "Activate: zero tasks") // TODO retries here?
		return fmt.Errorf("zero tasks")
	}

	for _, triggerType := range tasks[0].Schedule.TriggerTypes {
		err := tr.activateConcreteTrigger(ctx, triggerType, tasks)
		if err != nil {
			return err
		}
	}

	//For recurrent task combined with triggers updating RunTimeUTC it is responsibility of RecurrentTaskProcessor
	if tasks[0].Schedule.Regularity == api.Recurrent {
		return nil
	}

	for i, task := range tasks {
		if !task.Schedule.EndRunTime.IsZero() { // we don't need to reschedule 'never' trigger
			tasks[i].RunTimeUTC = task.Schedule.EndRunTime.UTC()
			continue
		}
		tasks[i].RunTimeUTC = tasks[i].RunTimeUTC.Add(hundredYears) // never task will never end
	}

	if err := tr.taskRepo.UpdateSchedulerFields(ctx, tasks...); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantUpdateTask, "Activate: UpdateSchedulerFields %v", err) // TODO retries here?
		return err
	}
	return nil
}

func (tr *Service) activateConcreteTrigger(ctx context.Context, triggerType string, tasks []m.Task) error {
	t := tasks[0]

	handler := tr.getTriggerHandler(triggerType)
	if err := handler.Activate(ctx, triggerType, tasks); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "Activate: %v", err)
		return err
	}

	var currentTriggerFrame api.TriggerFrame
	for _, frame := range t.Schedule.TriggerFrames {
		if frame.TriggerType == triggerType {
			currentTriggerFrame = frame
			break
		}
	}

	trigger := e.ActiveTrigger{
		TaskID:         t.ID,
		Type:           triggerType,
		PartnerID:      t.PartnerID,
		StartTimeFrame: currentTriggerFrame.StartTimeFrame,
		EndTimeFrame:   currentTriggerFrame.EndTimeFrame,
	}

	if err := tr.triggerRepo.Insert(trigger); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "Activate: trigger repo insert %v", err)
		return err
	}

	if err := tr.increaseTriggerCounter(ctx, triggerType); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "Activate: %v", err)
		return err
	}
	tr.activateCache(ctx, triggerType, t.PartnerID)
	return nil
}

func (tr *Service) activateCache(ctx context.Context, triggerType, partnerID string) {
	at, err := tr.triggerRepo.GetAllByType(ctx, triggerType, partnerID, false)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantGetActiveTriggers, "activateCache: GetAllByType %v", err)
		return
	}
	tr.setActiveTriggersToCache(ctx, at, partnerID, triggerType)
}

func (tr *Service) deactivateCache(ctx context.Context, triggerType, partnerID string) {
	at, err := tr.triggerRepo.GetAllByType(ctx, triggerType, partnerID, false)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "deactivateCache: GetAllByType %v", err)
		return
	}

	// this was the last task with this trigger type for a partner. remove from cache info by this key
	if len(at) == 0 {
		if err = tr.memCache.Delete(c.TriggerIDKeyPrefix + partnerID + triggerType); err != nil {
			tr.log.ErrfCtx(ctx, errorcode.ErrorCache, "deactivateCache delete: %v", err)
			return
		}

		bin, err := tr.memCache.Get(c.TriggerPartnersByTypePrefix + triggerType)
		if err != nil {
			tr.log.ErrfCtx(ctx, errorcode.ErrorCache, "deactivateCache Couldn't get list of partners that has active triggers of %s type from memcache, err: %v", triggerType, err)
			return
		}

		var partners map[string]struct{}
		if err = json.Unmarshal(bin.Value, &partners); err != nil {
			tr.log.ErrfCtx(ctx, errorcode.ErrorCache, "deactivateCache umarshal: %v", err)
			return
		}

		delete(partners, partnerID)
		if len(partners) != 0 {
			return
		}

		if err = tr.memCache.Delete(c.TriggerPartnersByTypePrefix + triggerType); err != nil {
			tr.log.ErrfCtx(ctx, errorcode.ErrorCache, "deactivateCache delete from cache: %v", err)
			return
		}
		return
	}

	tr.setActiveTriggersToCache(ctx, at, partnerID, triggerType)
}

func (tr *Service) setActiveTriggersToCache(ctx context.Context, at []e.ActiveTrigger, partnerID, triggerType string) {
	b, err := json.Marshal(at)
	if err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantMarshall, "setActiveTriggersToCache marshal: %v", err)
		return
	}
	if err = tr.memCache.Set(&memcache.Item{
		Key:   c.TriggerIDKeyPrefix + partnerID + triggerType,
		Value: b,
	}); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCache, "setActiveTriggersToCache: %v", err)
		return
	}
}

// Deactivate deactivates trigger
func (tr *Service) Deactivate(ctx context.Context, tasks []m.Task) error {
	if len(tasks) == 0 {
		tr.log.WarnfCtx(ctx, "Deactivate: zero tasks") // TODO retries here?
		return fmt.Errorf("zero tasks")
	}

	t := tasks[0]
	for _, triggerType := range tasks[0].Schedule.TriggerTypes {
		handler := tr.getTriggerHandler(triggerType)
		if err := handler.Deactivate(ctx, triggerType, tasks); err != nil {
			tr.log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "Deactivate: handler %v", err)
			return err
		}

		trigger := e.ActiveTrigger{
			TaskID:    t.ID,
			Type:      triggerType,
			PartnerID: t.PartnerID,
		}
		if err := tr.triggerRepo.Delete(trigger); err != nil {
			tr.log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "Deactivate: Delete %v", err) // TODO retries here?
			return err
		}

		if err := tr.decreaseTriggerCounter(ctx, triggerType); err != nil {
			tr.log.WarnfCtx(ctx, "Deactivate: decreaseTriggerCounter %v", err) // TODO retries here?
			// return is not missed here
		}
		tr.deactivateCache(ctx, triggerType, t.PartnerID)
	}

	for i := range tasks {
		tasks[i].State = statuses.TaskStateInactive
	}

	if err := tr.taskRepo.UpdateSchedulerFields(ctx, tasks...); err != nil {
		tr.log.ErrfCtx(ctx, errorcode.ErrorCantUpdateTask, "Deactivate: UpdateSchedulerFields %v", err) // TODO retries here?
		return err
	}
	return nil
}

// PreExecution runs preExecution for exact trigger type
func (tr *Service) PreExecution(triggerType string, task m.Task) error {
	handler := tr.getTriggerHandler(triggerType)
	return handler.PreExecution(task)
}

// PostExecution runs post execution operations
func (tr *Service) PostExecution(triggerType string, task m.Task) error {
	handler := tr.getTriggerHandler(triggerType)
	return handler.PostExecution(task)
}

// GetActiveTriggers returns all active triggers found by partner and type
func (tr *Service) GetActiveTriggers(ctx context.Context) ([]e.ActiveTrigger, error) {
	if err := tr.validateContext(ctx); err != nil {
		return nil, err
	}
	triggerType := ctx.Value(config.TriggerTypeIDKeyCTX).(string)
	partnerID := ctx.Value(config.PartnerIDKeyCTX).(string)
	return tr.triggerRepo.GetAllByType(ctx, triggerType, partnerID, true)
}

// GetTask GetTasks returns task with exact trigger type implementation
func (tr *Service) GetTask(ctx context.Context, taskID gocql.UUID) (m.Task, error) {
	if err := tr.validateContext(ctx); err != nil {
		return m.Task{}, err
	}
	triggerType := ctx.Value(config.TriggerTypeIDKeyCTX).(string)
	handler := tr.getTriggerHandler(triggerType)
	return handler.GetTask(ctx, taskID)
}

// IsApplicable checks if the task is applicable for a trigger
func (tr *Service) IsApplicable(ctx context.Context, task m.Task, payload api.TriggerExecutionPayload) bool {
	if err := tr.validateContext(ctx); err != nil {
		return false
	}
	triggerType := ctx.Value(config.TriggerTypeIDKeyCTX).(string)
	handler := tr.getTriggerHandler(triggerType)
	return handler.IsApplicable(ctx, task, payload)
}

// DeleteActiveTriggers deletes all given active triggers from db
func (tr *Service) DeleteActiveTriggers(triggers []e.ActiveTrigger) error {
	for _, t := range triggers {
		if err := tr.triggerRepo.Delete(t); err != nil {
			return err
		}
	}
	return nil
}

// GetActiveTriggersByTaskID returns all active triggers by partner and taskIDs
func (tr *Service) GetActiveTriggersByTaskID(partnerID string, taskID gocql.UUID) ([]e.ActiveTrigger, error) {
	return tr.triggerRepo.GetAllByTaskID(partnerID, taskID)
}

// getTriggerHandler returns TriggerHandler of particular implementation based on triggerType
func (tr *Service) getTriggerHandler(triggerType string) Handler {
	switch triggerType {
	case triggers.DynamicGroupEnterTrigger, triggers.DynamicGroupExitTrigger:
		return handlers.NewDynamicBasedTrigger(tr.dynamicGroups, tr.sitesRepo, tr.taskRepo, tr.cache, tr.log)
	case triggers.MockAlerting, triggers.MockGeneric:
		return mockusecases.MockTriggerHandler{}
	case triggers.LogoutTrigger, triggers.LoginTrigger:
		return handlers.NewLogOnOffTrigger(tr.agentConf, tr.taskRepo, tr.triggerRepo, tr.log, tr.sitesRepo, tr.profilesRepo, tr.assetsService)
	case triggers.FirstCheckInTrigger:
		return handlers.NewFirstCheckIn(tr.taskRepo, tr.log, tr.sitesRepo)
	default:
		return handlers.NewDefaultTrigger(tr.taskRepo, tr.triggerRepo, tr.log)
	}
}

// increaseTriggerCounter increases trigger counter and activates policy (if needed)
func (tr *Service) increaseTriggerCounter(ctx context.Context, triggerType string) (err error) {
	counter, err := tr.triggerRepo.GetTriggerCounterByType(triggerType)
	if err != nil {
		return fmt.Errorf("GetTriggerCounterByType: %v", err)
	}

	if counter.Count != 0 {
		if err = tr.triggerRepo.IncreaseCounter(counter); err != nil {
			err = fmt.Errorf("error during increasing counter %v", err)
		}
		return
	}

	// if there is no counter for such type - activate policy for it
	def, err := tr.triggerRepo.GetDefinition(triggerType)
	if err != nil {
		return fmt.Errorf("error during GetDefinition: %v", err)
	}

	policyID, err := tr.automationEngine.UpdateRemotePolicies(ctx, []e.TriggerDefinition{def})
	if err != nil {
		return fmt.Errorf("error during UpdateRemotePolicies: %v", err)
	}

	counter.TriggerID = triggerType
	counter.PolicyID = policyID
	if err = tr.triggerRepo.IncreaseCounter(counter); err != nil {
		tr.log.WarnfCtx(ctx, "Activate: error during increasing counter %v", err)
	}
	return nil
}

// decreaseTriggerCounter decreases trigger counter and deactivates policy (if needed)
func (tr *Service) decreaseTriggerCounter(ctx context.Context, triggerType string) (err error) {
	counter, err := tr.triggerRepo.GetTriggerCounterByType(triggerType)
	if err != nil {
		return fmt.Errorf("error during GetTriggerCounterByType:%v", err)
	}

	if counter.Count == 0 {
		return fmt.Errorf("there is no counter for this triggerType")
	}

	if counter.Count > 1 {
		if err = tr.triggerRepo.DecreaseCounter(counter); err != nil {
			err = fmt.Errorf("error during increasing counter %v", err)
		}
		return
	}

	triggerDef, err := tr.triggerRepo.GetDefinition(triggerType)
	if err != nil {
		return fmt.Errorf("GetDefinition: %v", err)
	}

	// if it's the last counter - deactivate it
	if err = tr.automationEngine.RemovePolicy(ctx, triggerDef.EventDetails.MessageIdentifier); err != nil {
		return
	}

	if err = tr.triggerRepo.DecreaseCounter(counter); err != nil {
		tr.log.WarnfCtx(ctx, "Deactivate: error during decreasing counter %v", err)
	}
	return nil
}

func (tr *Service) isGenericType(triggerType string) bool {
	return strings.Contains(triggerType, triggers.GenericTypePrefix)
}

func (tr *Service) isAlertType(triggerType string) bool {
	return strings.Contains(triggerType, triggers.AlertTypePrefix)
}

func (tr *Service) validateContext(ctx context.Context) error {
	if _, ok := ctx.Value(config.TriggerTypeIDKeyCTX).(string); !ok {
		return fmt.Errorf("TriggerTypeIDKeyCTX is empty")
	}
	if _, ok := ctx.Value(config.PartnerIDKeyCTX).(string); !ok {
		return fmt.Errorf("PartnerIDKeyCTX is empty")
	}
	if _, ok := ctx.Value(config.EndpointIDKeyCTX).(gocql.UUID); !ok {
		return fmt.Errorf("EndpointIDKeyCTX value is empty")
	}
	if _, ok := ctx.Value(config.SiteIDKeyCTX).(string); !ok {
		return fmt.Errorf("SiteIDKeyCTX value is empty")
	}
	if _, ok := ctx.Value(config.ClientIDKeyCTX).(string); !ok {
		return fmt.Errorf("ClientIDKeyCTX value is empty")
	}
	return nil
}
