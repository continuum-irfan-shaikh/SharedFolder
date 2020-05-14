package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/kafka"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/sites"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/validator"
)

const (
	updateEventKafkaMsg = "TaskService.SendTaskUpdateEventToKafka: can't send update message to Kafka. err: %v"
	partnerIDKey        = "partnerID"
	currentUserFmt      = "Current user: %+v"
)

//go:generate mockgen -destination=../../mocks/mocks-gomock/user_uc_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks UserUC
// UserUC users usecase interface
type UserUC interface {
	SaveEndpoints(ctx context.Context, ep []entities.Endpoints)
	EndpointsFromAsset(ctx context.Context, sites []string) ([]entities.Endpoints, error)
}

//go:generate mockgen -destination=../../mocks/mocks-gomock/targets_repo_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks TargetsRepo
// TargetsRepo represents interface to work with targets table
type TargetsRepo interface {
	Insert(partnerID string, taskID gocql.UUID, targets models.Target) error
}

//go:generate mockgen -destination=../../mocks/mock-usecases/execution_results_update_uc_mock.go -package=mockusecases gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks ExecutionResultUpdateUC
// ExecutionResultUpdateUC - interface for execution results update service
type ExecutionResultUpdateUC interface {
	ProcessExecutionResults(ctx context.Context, partnerID string, taskInstanceID gocql.UUID, results ...apiModels.ExecutionResult) error
}

//EncryptionService - provide us interface for encryption of credentials
type EncryptionService interface {
	Encrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
	Decrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
}

//AgentEncryptionService - provide interface for encryption credentials by key stored for particular endpoint
type AgentEncryptionService interface {
	Encrypt(ctx context.Context, endpointID gocql.UUID, credentials agentModels.Credentials) (encrypted agentModels.Credentials, err error)
}

type KafkaService interface {
	Push(message interface{}, msgType string) error
}

//go:generate mockgen -destination=../../mocks/mocks-gomock/site_repo_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks SiteRepo
// SiteRepo ..
type SiteRepo interface {
	GetEndpointsBySiteIDs(ctx context.Context, partnerID string, siteIDs []string) (ids []gocql.UUID, err error)
}

//go:generate mockgen -destination=../../mocks/mocks-gomock/dg_repo_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks DynamicGroupRepo
// DynamicGroupRepo ..
type DynamicGroupRepo interface {
	GetEndpointsByGroupIDs(ctx context.Context, targetIDs []string, createdBy, partnerID string, hasNOCAccess bool) (ids []gocql.UUID, err error)
}

type serviceError struct {
	errCode int
	errMsg  string
	err     error
}

// TaskService represents a Tasking Service instance
type TaskService struct {
	taskDefinition                 models.TaskDefinitionPersistence
	taskPersistence                models.TaskPersistence
	taskInstancePersistence        models.TaskInstancePersistence
	templateCache                  models.TemplateCache
	resultPersistence              models.ExecutionResultPersistence
	taskSummaryPersistence         models.TaskSummaryPersistence
	executionExpirationPersistence models.ExecutionExpirationPersistence
	userSitesPersistence           models.UserSitesPersistence
	userService                    user.Service
	assetsService                  integration.Asset
	siteRepo                       SiteRepo
	dgRepo                         DynamicGroupRepo
	taskCounterRepo                repository.TaskCounter
	cache                          persistency.Cache
	kafkaService                   KafkaService
	httpClient                     *http.Client
	userUC                         UserUC
	trigger                        trigger.Usecase
	triggerDefinition              trigger.DefinitionUseCase
	targetsRepo                    TargetsRepo
	executionResultUpdateService   ExecutionResultUpdateUC
	encryptionService              EncryptionService
	agentEncryptionService         AgentEncryptionService
}

// New creates a new Tasking Service initialized with the Tasking Repository
func New(
	taskDefinition models.TaskDefinitionPersistence,
	taskRepo models.TaskPersistence,
	taskInstanceRepo models.TaskInstancePersistence,
	templateCache models.TemplateCache,
	resultRepo models.ExecutionResultPersistence,
	taskInstancesStatusCountRepo models.TaskSummaryPersistence,
	executionExpirationRepo models.ExecutionExpirationPersistence,
	userSitesRepo models.UserSitesPersistence,
	userService user.Service,
	assetsService integration.Asset,
	siteRepo SiteRepo,
	dgRepo DynamicGroupRepo,
	taskCounterRepo repository.TaskCounter,
	cache persistency.Cache,
	kafkaService KafkaService,
	httpClient *http.Client,
	userUC UserUC,
	tr trigger.Usecase,
	triggerDefinition trigger.DefinitionUseCase,
	targetsRepo TargetsRepo,
	executionResultUpdate ExecutionResultUpdateUC,
	encryptionService EncryptionService,
	agentEncryptionService AgentEncryptionService,
) (service TaskService) {
	return TaskService{
		taskDefinition:                 taskDefinition,
		taskPersistence:                taskRepo,
		taskInstancePersistence:        taskInstanceRepo,
		templateCache:                  templateCache,
		resultPersistence:              resultRepo,
		taskSummaryPersistence:         taskInstancesStatusCountRepo,
		executionExpirationPersistence: executionExpirationRepo,
		userSitesPersistence:           userSitesRepo,
		userService:                    userService,
		assetsService:                  assetsService,
		siteRepo:                       siteRepo,
		dgRepo:                         dgRepo,
		taskCounterRepo:                taskCounterRepo,
		cache:                          cache,
		kafkaService:                   kafkaService,
		httpClient:                     httpClient,
		userUC:                         userUC,
		trigger:                        tr,
		triggerDefinition:              triggerDefinition,
		targetsRepo:                    targetsRepo,
		executionResultUpdateService:   executionResultUpdate,
		encryptionService:              encryptionService,
		agentEncryptionService:         agentEncryptionService,
	}
}

// SendTaskUpdateEventToKafka send update task message to taskingEvent kafka topic to update task
func (t *TaskService) SendTaskUpdateEventToKafka(ctx context.Context, taskID gocql.UUID, partnerID string) {
	msg := models.TaskIsUpdatedMessage{TaskID: taskID, PartnerID: partnerID}

	err := t.kafkaService.Push(msg, kafka.MsgTypeTaskUpdated)
	if err != nil {
		logger.Log.WarnfCtx(ctx, updateEventKafkaMsg, err)
		return
	}
}

// GetByPartnerAndManagedEndpointID returns a list of Tasks by target_id and partner_id
func (taskService TaskService) GetByPartnerAndManagedEndpointID(w http.ResponseWriter, r *http.Request) {
	endpointID, err := common.ExtractUUID("TaskService.GetByPartnerAndManagedEndpointID", w, r, "managedEndpointID")
	if err != nil {
		return
	}

	partnerID := mux.Vars(r)[partnerIDKey]
	currentUser := taskService.userService.GetUser(r, taskService.httpClient)

	listOfInternalTasksByTarget, err := taskService.taskPersistence.GetByPartnerAndManagedEndpointID(r.Context(), partnerID, endpointID, common.UnlimitedCount)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetListOfTasksByManagedEndpoint, "TaskService.GetByPartnerAndManagedEndpointID: can not get list of Tasks by ManagedEndpointID. err=%v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetListOfTasksByManagedEndpoint)
		return
	}

	taskIDs := make([]gocql.UUID, len(listOfInternalTasksByTarget))
	for i, internalTask := range listOfInternalTasksByTarget {
		taskIDs[i] = internalTask.ID
	}

	listOfInternalTasksByIDs, err := taskService.taskPersistence.GetByIDs(r.Context(), nil, partnerID, false, taskIDs...)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetListOfTasksByManagedEndpoint, "TaskService.GetByPartnerAndManagedEndpointID: can not get list of Tasks by IDs(%v). err=%v", taskIDs, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetListOfTasksByManagedEndpoint)
		return
	}

	// Filter by NOC
	filteredTaskGroup := filterTasksByNOCAccsess(listOfInternalTasksByIDs, currentUser.HasNOCAccess())

	// Group tasks by taskIDs
	tasksGroup := make(map[gocql.UUID][]models.Task)
	for _, task := range filteredTaskGroup {
		tasksGroup[task.ID] = append(tasksGroup[task.ID], task)
	}

	listOfOutputTasks := make([]models.Task, 0)
	for _, internalTaskGroup := range tasksGroup {
		taskOutputPtr, err := models.NewTaskOutput(r.Context(), internalTaskGroup)
		if err != nil {
			logger.Log.WarnfCtx(
				r.Context(),
				"TaskService.GetByPartnerAndManagedEndpointID: failed to build task output from internal tasks %v. Err: %s",
				internalTaskGroup,
				err,
			)
			continue
		}
		listOfOutputTasks = append(listOfOutputTasks, *taskOutputPtr)
	}

	logger.Log.DebugfCtx(r.Context(), "TaskService.GetByPartnerAndManagedEndpointID: successfully returned list of Tasks by ManagedEndpointID.")
	common.RenderJSON(w, listOfOutputTasks)
}

// GetTasksSummaryData returns a Tasks Summary Data
func (taskService TaskService) GetTasksSummaryData(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)[partnerIDKey]
	currentUser := taskService.userService.GetUser(r, taskService.httpClient)

	summaryPageData, err := taskService.taskSummaryPersistence.GetTasksSummaryData(
		r.Context(),
		currentUser.HasNOCAccess(),
		taskService.cache,
		partnerID,
	)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTasksSummaryData, "TaskService.GetTasksSummaryData: can't get tasks summary data for partner=%s, err=%v.", currentUser.PartnerID(), err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTasksSummaryData)
		return
	}

	logger.Log.DebugfCtx(r.Context(), "TaskService.GetTasksSummaryData: successfully returned tasksSummaryData")
	common.RenderJSON(w, summaryPageData)
}

// GetByPartnerAndID returns a Task by task_id and partner_id
func (taskService TaskService) GetByPartnerAndID(w http.ResponseWriter, r *http.Request) {
	taskID, err := common.ExtractUUID("TaskService.GetByPartnerAndID", w, r, "taskID")
	if err != nil {
		return
	}

	partnerID := mux.Vars(r)[partnerIDKey]
	currentUser := taskService.userService.GetUser(r, taskService.httpClient)

	internalTasks, err := taskService.taskPersistence.GetByIDs(r.Context(), nil, partnerID, false, taskID)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "TaskService.GetByPartnerAndID: can not get a Task by Task ID %v. err=%v", taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	filteredInternalTasks := filterTasksByNOCAccsess(internalTasks, currentUser.HasNOCAccess())

	taskOutput, err := models.NewTaskOutput(r.Context(), filteredInternalTasks)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "TaskService.GetByPartnerAndID: Task with ID %v and PartnerID %s not found. Err : %s", taskID, partnerID, err.Error())
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	logger.Log.InfofCtx(r.Context(), "TaskService.GetByPartnerAndID: successfully returned a Task by Task ID.")
	common.RenderJSON(w, taskOutput)
}

// Edit sets all internal tasks by TaskID as inactive and creates new internal tasks with new values of Targets, Schedule or Parameters fields
func (t *TaskService) Edit(w http.ResponseWriter, r *http.Request) {
	var (
		emptyUUID   gocql.UUID
		taskIDStr   = mux.Vars(r)["taskID"]
		partnerID   = mux.Vars(r)[partnerIDKey]
		ctx         = r.Context()
		currentUser = t.userService.GetUser(r, t.httpClient)
		modifiedAt  = time.Now().Truncate(time.Millisecond).UTC()
	)

	taskID, err := gocql.ParseUUID(taskIDStr)
	if err != nil || taskID == emptyUUID {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskIDHasBadFormat, "TaskService.Edit: task ID(UUID=%s) has bad format or empty. err=%v", taskIDStr, err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskIDHasBadFormat)
		return
	}

	inputTask, err := t.extractPostTaskPayload(r, w)
	if err != nil {
		return
	}

	inputTask.PartnerID = partnerID
	internalTasks, err := t.taskPersistence.GetByIDs(ctx, nil, inputTask.PartnerID, false, taskID)
	if err != nil && err != gocql.ErrNotFound {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "TaskService.Edit: can not get internal tasks by Task ID %v. err=%v", taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	if len(internalTasks) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, "TaskService.Edit: can not get internal tasks by Task ID %v. err=%v", taskID, err)
		common.SendNotFound(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}

	task := internalTasks[0]
	task.Schedule = inputTask.Schedule
	task.Schedule.StartRunTime = task.Schedule.StartRunTime.Truncate(time.Minute)
	task.Schedule.EndRunTime = task.Schedule.EndRunTime.Truncate(time.Minute)
	task.ModifiedBy = currentUser.UID()
	task.ModifiedAt = modifiedAt
	task.DefinitionID = inputTask.DefinitionID
	task.OriginID = inputTask.OriginID
	task.PartnerID = partnerID

	task.TargetsByType = inputTask.TargetsByType
	if task.TargetsByType == nil {
		task.TargetsByType = make(models.TargetsByType)
	}

	if inputTask.Targets.Type != 0 {
		task.TargetsByType[inputTask.Targets.Type] = inputTask.Targets.IDs
	}

	for targetType, targets := range task.TargetsByType {
		task.Targets.Type = targetType
		task.Targets.IDs = targets
	}

	if len(inputTask.Parameters) > 0 {
		task.Parameters = inputTask.Parameters
	}

	for i := range internalTasks {
		internalTasks[i].OriginalNextRunTime = time.Time{}
		if internalTasks[i].State != statuses.TaskStateDisabled {
			internalTasks[i].State = statuses.TaskStateInactive
		}
		internalTasks[i].ModifiedBy = currentUser.UID()
	}

	t.processEditReq(ctx, internalTasks, r, w, currentUser, task)
}

func (t *TaskService) processEditReq(ctx context.Context, internalTasks []models.Task, r *http.Request, w http.ResponseWriter, currentUser user.User, task models.Task) {
	taskID := internalTasks[0].ID
	partnerID := task.PartnerID

	err := t.taskPersistence.InsertOrUpdate(ctx, internalTasks...)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskToDB, "TaskService.Edit: error while updating internal tasks. Edited task ID=%v. err=%v", taskID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskToDB)
		return
	}

	if internalTasks[0].IsTrigger() && internalTasks[0].IsActivatedTrigger() { // if it's not activated - don't do anything
		if err = t.trigger.Deactivate(ctx, internalTasks); err != nil { // sending old ones only
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskToDB, "TaskService.Edit: error while updating internal tasks. Edited task ID=%v. err=%v", taskID, err)
			common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskToDB)
			return
		}
	}

	siteIDs, err := sites.GetSiteIDs(ctx, t.httpClient, partnerID, config.Config.SitesMsURL, currentUser.Token())
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Edit: User-Site restrictions for dynamic groups, err: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		return
	}

	if ids, ok := task.TargetsByType[models.Site]; ok && len(ids) > 0 && !currentUser.HasNOCAccess() {
		site, ok := isUserSites(ids, siteIDs)
		if !ok {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: User doesn't have access to site with id: %s", site)
			common.SendForbidden(w, r, errorcode.ErrorCantCreateNewTask)
			return
		}
	}

	if ids, ok := task.TargetsByType[models.DynamicSite]; ok && len(ids) > 0 && !currentUser.HasNOCAccess() {
		site, ok := isUserSites(ids, siteIDs)
		if !ok {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: User doesn't have access to site with id: %s", site)
			common.SendForbidden(w, r, errorcode.ErrorCantCreateNewTask)
			return
		}
	}

	err = t.userSitesPersistence.InsertUserSites(r.Context(), partnerID, currentUser.UID(), siteIDs)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Edit: User-Site restrictions for dynamic groups, err: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		return
	}

	newInternalTasks, err := t.createEditedTask(ctx, task, currentUser, r, w, internalTasks)
	if err != nil {
		return
	}

	logger.Log.InfofCtx(r.Context(), "TaskService.Edit: edited task with ID = %v", taskID)
	if len(newInternalTasks) == 0 {
		common.RenderJSON(w, task)
		return
	}
	common.RenderJSON(w, newInternalTasks[0])
}

func (t *TaskService) createEditedTask(ctx context.Context, task models.Task, currentUser user.User, r *http.Request, w http.ResponseWriter, internalTasks []models.Task) ([]models.Task, error) {
	taskID := internalTasks[0].ID
	partnerID := task.PartnerID
	newInternalTasks, taskSiteIDs, sErr := t.createNewInternalTasks(ctx, &task, currentUser.HasNOCAccess())
	if sErr != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Edit: error while creating new internal tasks. err=%v", sErr.err)

		switch sErr.errCode {
		case http.StatusBadRequest:
			common.SendBadRequest(w, r, sErr.errMsg)
		case http.StatusNotFound:
			common.SendNotFound(w, r, sErr.errMsg)
		case http.StatusInternalServerError:
			common.SendInternalServerError(w, r, sErr.errMsg)
		default:
			common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		}
		return nil, fmt.Errorf("%v", sErr)
	}

	if err := t.userSitesPersistence.InsertSitesByTaskID(ctx, taskID, taskSiteIDs); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Edit: error while saving taskSiteIDs[%v] of tasks: %v", taskSiteIDs, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		return nil, err
	}

	errMsg, err := t.saveAndProcessTask(ctx, &task, newInternalTasks, true, true)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Edit: error while saving and processing internal tasks. err=%v", err)
		common.SendInternalServerError(w, r, errMsg)
		return nil, err
	}

	if err = t.targetsRepo.Insert(task.PartnerID, task.ID, task.Targets); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: error while saving targets for task: id: %v, cause err: %s", task.ID, err.Error())
		common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		return nil, err
	}

	if !currentUser.HasNOCAccess() {
		// update counters for tasks in separate goroutine
		go t.updateTaskCounters(ctx, internalTasks, newInternalTasks)
	}

	go t.SendTaskUpdateEventToKafka(ctx, task.ID, partnerID)
	return newInternalTasks, nil
}

func (t *TaskService) updateTaskCounters(ctx context.Context, oldTasks, newTasks []models.Task) {
	oldCounters := getCountersForInternalTasks(oldTasks)
	newCounters := getCountersForInternalTasks(newTasks)

	if err := t.taskCounterRepo.DecreaseCounter(oldTasks[0].PartnerID, oldCounters, false); err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "TaskService.Edit: error while trying to decrease counter: ", err)
	}

	if err := t.taskCounterRepo.IncreaseCounter(oldTasks[0].PartnerID, newCounters, false); err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "TaskService.Edit: error while trying to increase counter: ", err)
	}
}

// Create inserts a received Task into Cassandra with new TaskID
func (t *TaskService) Create(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	currentUser, err := common.UserFromCtx(ctx)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "TaskService.Create: %s", "can't get current user. err: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}
	logger.Log.InfofCtx(r.Context(), currentUserFmt, currentUser)

	task, err := t.extractPostTaskPayload(r, w)
	if err != nil {
		return
	}

	task.CreatedBy = currentUser.UID
	task.PartnerID = currentUser.PartnerID
	task.IsRequireNOCAccess = currentUser.IsNOCAccess
	task.Schedule.StartRunTime = task.Schedule.StartRunTime.Truncate(time.Minute)
	task.Schedule.EndRunTime = task.Schedule.EndRunTime.Truncate(time.Minute)
	task.ID = gocql.TimeUUID()
	if task.TargetsByType == nil {
		task.TargetsByType = make(models.TargetsByType)
	}

	if task.Targets.Type != 0 {
		task.TargetsByType[task.Targets.Type] = task.Targets.IDs
	}

	for targetType, targets := range task.TargetsByType {
		task.Targets.Type = targetType
		task.Targets.IDs = targets
	}

	siteIDs, err := sites.GetSiteIDs(ctx, t.httpClient, currentUser.PartnerID, config.Config.SitesMsURL, currentUser.Token)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: User-Site restrictions for dynamic groups, err: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		return
	}

	if ids, ok := task.TargetsByType[models.Site]; ok && len(ids) > 0 && !currentUser.IsNOCAccess {
		site, ok := isUserSites(ids, siteIDs)
		if !ok {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: User doesn't have access to site with id: %s", site)
			common.SendForbidden(w, r, errorcode.ErrorCantCreateNewTask)
			return
		}
	}

	if ids, ok := task.TargetsByType[models.DynamicSite]; ok && len(ids) > 0 && !currentUser.IsNOCAccess {
		site, ok := isUserSites(ids, siteIDs)
		if !ok {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: User doesn't have access to site with id: %s", site)
			common.SendForbidden(w, r, errorcode.ErrorCantCreateNewTask)
			return
		}
	}

	err = t.userSitesPersistence.InsertUserSites(r.Context(), currentUser.PartnerID, currentUser.UID, siteIDs)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: User-Site restrictions for dynamic groups, err: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		return
	}

	if task.ExternalTask {
		go t.createTaskBackground(ctx, task, currentUser.IsNOCAccess, r, siteIDs)
		common.RenderJSONCreated(w, task)
		return
	}

	t.createTaskFlow(ctx, task, currentUser, r, w, siteIDs)
}

func (t *TaskService) extractPostTaskPayload(r *http.Request, w http.ResponseWriter) (task models.Task, err error) {
	err = validator.ExtractStructFromRequest(r, &task)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "extractPostTaskPayload: %s", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return models.Task{}, err
	}

	if len(task.ParametersObject) != 0 {
		b, err := json.Marshal(task.ParametersObject)
		if err != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "extractPostTaskPayload: can't marshall parameters object. err: %s", err.Error())
			common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
			return models.Task{}, err
		}
		task.Parameters = string(b)
	}

	if task.Credentials != nil {
		encrypted, err := t.encryptionService.Encrypt(*task.Credentials)
		if err != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantEcryptCredentials, "extractPostTaskPayload: can't encrypt credentials. err: %s", err.Error())
			common.SendInternalServerError(w, r, errorcode.ErrorCantEcryptCredentials)
			return models.Task{}, err
		}
		task.Credentials = &encrypted
	}

	// RMM-57339 backward compatibility
	if task.Schedule.Repeat.WeekDay != nil && len(task.Schedule.Repeat.WeekDays) == 0 {
		task.Schedule.Repeat.WeekDays = append(task.Schedule.Repeat.WeekDays, *task.Schedule.Repeat.WeekDay)
	}

	return task, nil
}

func (t *TaskService) createTaskFlow(ctx context.Context, task models.Task, currentUser entities.User, r *http.Request, w http.ResponseWriter, siteIDs []int64) {
	internalTasks, taskSiteIDs, sErr := t.createNewInternalTasks(ctx, &task, currentUser.IsNOCAccess)
	if sErr != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: error while creating new internal tasks. err=%v", sErr.err)

		switch sErr.errCode {
		case http.StatusBadRequest:
			common.SendBadRequest(w, r, sErr.errMsg)
		case http.StatusNotFound:
			common.SendNotFound(w, r, sErr.errMsg)
		case http.StatusInternalServerError:
			common.SendInternalServerError(w, r, sErr.errMsg)
		}
		return
	}

	err := t.userSitesPersistence.InsertSitesByTaskID(ctx, task.ID, taskSiteIDs)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantInsertData, "TaskService.Create: error while saving taskSiteIDs[%v] of tasks: %v", taskSiteIDs, err)
		common.SendInternalServerError(w, r, err.Error())
		return
	}

	errMsg, err := t.saveAndProcessTask(ctx, &task, internalTasks, true, true)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: error while saving and processing internal tasks: %v", err)
		common.SendInternalServerError(w, r, errMsg)
		return
	}

	// this call could be removed when resource selector build will be deployed
	if err = t.targetsRepo.Insert(task.PartnerID, task.ID, task.Targets); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.Create: error while saving targets for task: id: %v, cause err: %s", task.ID, err.Error())
		common.SendInternalServerError(w, r, errorcode.ErrorCantCreateNewTask)
		return
	}

	if !currentUser.IsNOCAccess {
		// update counters for tasks in separate goroutine
		go func(ctx context.Context, iTasks []models.Task) {
			counters := getCountersForInternalTasks(iTasks)

			err := t.taskCounterRepo.IncreaseCounter(task.PartnerID, counters, false)
			if err != nil {
				logger.Log.WarnfCtx(ctx, "TaskService.Create: error while trying to increase counter: ", err)
			}
		}(ctx, internalTasks)
	}

	logger.Log.DebugfCtx(r.Context(), "TaskService.Create: created task with ID = %v", task.ID)
	if len(internalTasks) == 0 {
		common.RenderJSONCreated(w, task)
		return
	}

	siteIDsStr := make([]string, 0)
	for _, id := range siteIDs {
		siteIDsStr = append(siteIDsStr, strconv.Itoa(int(id)))
	}

	outputTask := internalTasks[0]
	outputTask.TargetType = 0
	common.RenderJSONCreated(w, outputTask)
}

func (t *TaskService) createTaskBackground(ctx context.Context, task models.Task, hasNOCAccess bool, r *http.Request, siteIDs []int64) {
	internalTasks, taskSiteIDs, sErr := t.createNewInternalTasks(ctx, &task, hasNOCAccess)
	if sErr != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.createTaskBackground: error while creating new internal tasks. err=%v", sErr.err)
		return
	}

	err := t.userSitesPersistence.InsertSitesByTaskID(ctx, task.ID, taskSiteIDs)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantInsertData, "TaskService.createTaskBackground: error while saving taskSiteIDs[%v] of tasks: %v", taskSiteIDs, err)
		return
	}

	_, err = t.saveAndProcessTask(ctx, &task, internalTasks, true, true)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantCreateNewTask, "TaskService.createTaskBackground: error while saving and processing internal tasks: %v", err)
		return
	}

	// this call could be removed when resource selector build will be deployed
	if err = t.targetsRepo.Insert(task.PartnerID, task.ID, task.Targets); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantInsertData, "TaskService.createTaskBackground: error while saving targets for task: id: %v, cause err: %s", task.ID, err.Error())
		return
	}

	if !hasNOCAccess {
		// update counters for tasks in separate goroutine
		go func(ctx context.Context, iTasks []models.Task) {
			counters := getCountersForInternalTasks(iTasks)

			err := t.taskCounterRepo.IncreaseCounter(task.PartnerID, counters, false)
			if err != nil {
				logger.Log.WarnfCtx(ctx, "TaskService.createTaskBackground: error while trying to increase counter: ", err)
			}
		}(ctx, internalTasks)
	}

	logger.Log.DebugfCtx(r.Context(), "TaskService.createTaskBackground: created task with ID = %v", task.ID)
	if len(internalTasks) == 0 {
		return
	}

	siteIDsStr := make([]string, 0)
	for _, id := range siteIDs {
		siteIDsStr = append(siteIDsStr, strconv.Itoa(int(id)))
	}
}

// used for Create because client sends targets data in target list
func (taskService TaskService) getManagedEndpointIDsMapFromTargets(targets []string) []gocql.UUID {
	managedEndpointIDs := make([]gocql.UUID, 0)
	for _, id := range targets {
		endpointID, err := gocql.ParseUUID(id)
		if err != nil {
			continue
		}
		managedEndpointIDs = append(managedEndpointIDs, endpointID)
	}

	return managedEndpointIDs
}

func (taskService TaskService) getTemplateAndValidateParameters(ctx context.Context, task *models.Task, hasNOCAccess bool) (template models.TemplateDetails, sErr *serviceError) {
	var err error
	// Default template for Tasks which are different from type Script
	template = models.TemplateDetails{Name: "Scheduled task"}

	if task.Type != config.ScriptTaskType {
		template.Description = task.Description
		if len(task.Parameters) == 0 {
			return
		}

		var testJSON map[string]interface{}
		err = json.Unmarshal([]byte(task.Parameters), &testJSON)
		if err != nil {
			sErr = &serviceError{
				errCode: http.StatusBadRequest,
				errMsg:  errorcode.ErrorCantDecodeInputData,
				err:     errors.Wrap(err, "parameters must be in JSON format"),
			}
			return
		}
		return
	}

	template, err = taskService.templateCache.GetByOriginID(ctx, task.PartnerID, task.OriginID, hasNOCAccess)
	if err != nil {
		switch err.(type) {
		case models.TemplateNotFoundError:
			sErr = &serviceError{
				errCode: http.StatusNotFound,
				errMsg:  errorcode.ErrorCantGetTaskDefinitionTemplate,
				err:     err,
			}
		default:
			sErr = &serviceError{
				errCode: http.StatusInternalServerError,
				errMsg:  errorcode.ErrorCantGetTemplatesForExecutionMS,
				err:     err,
			}
		}
		return
	}

	err = validator.ValidateParametersField(template.JSONSchema, task.Parameters, true)
	if err != nil {
		sErr = &serviceError{
			errCode: http.StatusBadRequest,
			errMsg:  errorcode.ErrorCantDecodeInputData,
			err:     err,
		}
		return
	}
	// user can now define custom description even for continuum defined scripts
	if task.Description != "" {
		template.Description = task.Description
	}
	return
}

func (taskService TaskService) getEndpoints(ctx context.Context, task models.Task) (eIDs map[gocql.UUID]models.TargetType, siteIDs []string, err error) {
	var eIDList []gocql.UUID
	siteIDs = make([]string, 0)
	eIDs = make(map[gocql.UUID]models.TargetType)
	uniqueSites := make(map[string]struct{})

	if targets, ok := task.TargetsByType[models.DynamicGroup]; ok && len(targets) > 0 {
		err = taskService.extractDynamicGroupTargets(ctx, targets, task, eIDs, uniqueSites)
		if err != nil {
			return
		}
	}

	if targets, ok := task.TargetsByType[models.Site]; ok && len(targets) > 0 {
		err = taskService.extractSitesFromTargets(ctx, models.Site, task, targets, eIDs, uniqueSites)
		if err != nil {
			return
		}
	}

	if targets, ok := task.TargetsByType[models.DynamicSite]; ok && len(targets) > 0 {
		err = taskService.extractSitesFromTargets(ctx, models.DynamicSite, task, targets, eIDs, uniqueSites)
		if err != nil {
			return
		}
	}

	if targets, ok := task.TargetsByType[models.ManagedEndpoint]; ok && len(targets) > 0 {
		eIDList = taskService.getManagedEndpointIDsMapFromTargets(targets)

		for _, e := range eIDList {
			eIDs[e] = models.ManagedEndpoint
		}

		uniqueSites = taskService.getTaskSites(ctx, eIDList, uniqueSites, task.PartnerID)
	}

	for site := range uniqueSites {
		siteIDs = append(siteIDs, site)
	}
	return
}

func (taskService TaskService) extractSitesFromTargets(ctx context.Context, siteType models.TargetType, task models.Task, targets []string, eIDs map[gocql.UUID]models.TargetType, uniqueSites map[string]struct{}) error {
	eIDList, err := taskService.siteRepo.GetEndpointsBySiteIDs(ctx, task.PartnerID, targets)
	if err != nil {
		return err
	}

	for _, e := range eIDList {
		eIDs[e] = siteType
	}

	for _, site := range targets {
		uniqueSites[site] = struct{}{}
	}
	return nil
}

func (taskService TaskService) extractDynamicGroupTargets(ctx context.Context, targets []string, task models.Task, eIDs map[gocql.UUID]models.TargetType, uniqueSites map[string]struct{}) error {
	eIDList, err := taskService.dgRepo.GetEndpointsByGroupIDs(ctx, targets, task.CreatedBy, task.PartnerID, task.IsRequireNOCAccess)
	if err != nil {
		return err
	}

	for _, e := range eIDList {
		eIDs[e] = models.DynamicGroup
	}

	uniqueSites = taskService.getTaskSites(ctx, eIDList, uniqueSites, task.PartnerID)
	return nil
}

func (taskService TaskService) saveAndProcessTask(ctx context.Context, inputTask *models.Task, internalTasks []models.Task, updateTask bool, createdByTask bool) (errMsg string, err error) {
	taskInstance := models.NewTaskInstance(internalTasks, createdByTask)

	if inputTask.Schedule.Regularity != apiModels.RunNow {
		err = taskService.taskInstancePersistence.Insert(ctx, taskInstance)
		if err != nil {
			errMsg = errorcode.ErrorCantUpdateTask
			return
		}
	}

	if updateTask {
		for i := range internalTasks {
			internalTasks[i].LastTaskInstanceID = taskInstance.ID
		}

		internalTasks = taskService.deactivateRunNowTasks(internalTasks)
		if err = taskService.taskPersistence.InsertOrUpdate(ctx, internalTasks...); err != nil {
			errMsg = errorcode.ErrorCantSaveTaskToDB
			return
		}
	}

	if inputTask.Schedule.Regularity == apiModels.RunNow {
		for endpointID := range taskInstance.Statuses {
			taskInstance.Statuses[endpointID] = statuses.TaskInstanceRunning
		}

		err = taskService.SendTaskOnExecution(ctx, internalTasks, taskInstance)
		if err != nil {
			errMsg = errorcode.ErrorCantPrepareTaskForSendingOnExecution
			return
		}
	}

	taskService.setCreatedTaskToCache(ctx, taskInstance, inputTask, internalTasks)
	return errMsg, nil
}

func (taskService TaskService) setCreatedTaskToCache(ctx context.Context, taskInstance models.TaskInstance, inputTask *models.Task, internalTasks []models.Task) {
	go func(taskInstance models.TaskInstance, inputTask *models.Task) {
		allStatuses, err := taskInstance.CalculateStatuses()
		if err != nil {
			return
		}

		taskService.kafkaService.Push(models.KafkaMessage{
			TaskID:              inputTask.ID,
			PartnerID:           inputTask.PartnerID,
			UID:                 inputTask.CreatedBy,
			IsRequiredNOCAccess: inputTask.IsRequireNOCAccess,
			Entity: models.TaskDetailsWithStatuses{
				Task:         *inputTask,
				TaskInstance: taskInstance,
				Statuses:     allStatuses,
			},
		}, kafka.MsgTypeTaskCreated) // during creating and editing tasks
	}(taskInstance, inputTask)

	if config.Config.AssetCacheEnabled && taskService.cache != nil {
		t := internalTasks[0]
		t.ManagedEndpoints = nil
		if err := taskService.setTaskToCache(t); err != nil {
			logger.Log.WarnfCtx(ctx, "TaskService.saveAndProcessTask: couldn't marshal task while inserting to cache:%s", err.Error())
		}
	}
}

func (taskService *TaskService) deactivateRunNowTasks(tasks []models.Task) []models.Task {
	for i, task := range tasks {
		if task.Schedule.Regularity == apiModels.RunNow {
			tasks[i].State = statuses.TaskStateInactive
		}
	}
	return tasks
}

// getting siteIDs for ManagedEndpoints
func (taskService TaskService) getTaskSites(ctx context.Context, managedEndpoints []gocql.UUID, siteIDs map[string]struct{}, partnerID string) map[string]struct{} {
	if len(managedEndpoints) == 0 {
		return siteIDs
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(len(managedEndpoints))

	for _, endpointID := range managedEndpoints {
		go func(ctx context.Context, endpointID gocql.UUID) {
			defer wg.Done()

			siteID, _, err := taskService.assetsService.GetSiteIDByEndpointID(ctx, partnerID, endpointID)
			if err != nil {
				logger.Log.WarnfCtx(ctx, "cannot get siteID for ManagedEndpoint[%v]: %v", endpointID, err)
				return
			}

			mu.Lock()
			siteIDs[siteID] = struct{}{}
			mu.Unlock()
		}(ctx, endpointID)
	}

	wg.Wait()
	return siteIDs
}

func getCountersForInternalTasks(internalTasks []models.Task) []models.TaskCount {
	counters := make([]models.TaskCount, len(internalTasks))
	for i, task := range internalTasks {
		if task.ExternalTask {
			// we don't need to have count of Tasks created by external services (for instance, Sequence MS)
			return []models.TaskCount{}
		}

		counters[i] = models.TaskCount{
			ManagedEndpointID: task.ManagedEndpointID,
			Count:             1,
		}
	}
	return counters
}

func buildManagedEndpoints(tasks []models.Task, taskInstance models.TaskInstance) []apiModels.ManagedEndpoint {
	managedEndpointIDs := make([]apiModels.ManagedEndpoint, 0, len(tasks))
	for _, task := range tasks {
		// skip Stopped, Disabled and other Endpoints which should not run
		if status, ok := taskInstance.Statuses[task.ManagedEndpointID]; ok && status == statuses.TaskInstanceRunning {
			// add only running
			managedEndpointIDs = append(managedEndpointIDs, apiModels.ManagedEndpoint{
				ID:          task.ManagedEndpointID.String(),
				NextRunTime: task.RunTimeUTC.UTC(),
			})
		}
	}
	return managedEndpointIDs
}

// SendTaskOnExecution creates execution payload message for executing Task(Tasking/Patching/etc) on target MS
// tasks is expected not to be an empty slice
func (taskService TaskService) SendTaskOnExecution(ctx context.Context, tasks []models.Task, taskInstance models.TaskInstance) error {
	if len(tasks) < 1 {
		return fmt.Errorf("empty tasks")
	}

	var (
		executionErr, taskInstanceErr error
		wg                            = &sync.WaitGroup{}
		tasksCommonFields             = tasks[0]
		managedEndpointIDs            = buildManagedEndpoints(tasks, taskInstance)
	)

	payload := apiModels.ExecutionPayload{
		ExecutionID:      taskInstance.ID.String(),
		ManagedEndpoints: managedEndpointIDs,
		OriginID:         taskInstance.OriginID.String(),
		Parameters:       tasksCommonFields.Parameters,
		WebhookURL: fmt.Sprintf("%s/partners/%s/task-execution-results/task-instances/%s",
			config.Config.TaskingMsURL, tasksCommonFields.PartnerID, taskInstance.ID),
		TaskID:                   tasksCommonFields.ID,
		ExpectedExecutionTimeSec: taskService.templateCache.CalculateExpectedExecutionTimeSec(ctx, tasksCommonFields),
	}

	wg.Add(2)

	// save data about Tasks' expiration
	go func(cx context.Context, task models.Task, ti models.TaskInstance, payload apiModels.ExecutionPayload) {
		defer wg.Done()
		executionErr = taskService.saveExecutionExpiration(cx, task.PartnerID, ti, payload.ExpectedExecutionTimeSec)
	}(ctx, tasksCommonFields, taskInstance, payload)

	// update TaskInstance
	go func(ctx context.Context, ti models.TaskInstance) {
		defer wg.Done()
		// ti.StartedAt = time.Now().UTC() - do not do this! it breaks scheduler functionality
		ti.LastRunTime = time.Now().UTC()
		ti.OverallStatus = statuses.TaskInstanceRunning
		taskInstanceErr = taskService.taskInstancePersistence.Insert(ctx, ti)
	}(ctx, taskInstance)
	wg.Wait()

	if executionErr != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantSaveExecutionExpiration, "saveExecutionExpiration: err: %v", executionErr)
		return executionErr
	}
	if taskInstanceErr != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "error while updating taskInstance [%v]: %v", taskInstance.ID, taskInstanceErr)
		return taskInstanceErr
	}

	executionURL, err := getExecutionURL(tasksCommonFields.Type, tasksCommonFields.PartnerID)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "SendTaskOnExecution: could't get correct execution URL for partnerID=%v and taskType=%v, err=%v",
			tasksCommonFields.PartnerID, tasksCommonFields.Type, err)
		return err
	}

	if tasksCommonFields.IsRunAsUserApplied() {
		if err = taskService.sendTaskWithCredentials(ctx, tasksCommonFields, managedEndpointIDs, payload, executionURL, taskInstance); err != nil {
			return err
		}
		return nil
	}

	go taskService.sendTaskOnExecutionREST(ctx, executionURL, payload, taskInstance, tasksCommonFields)
	return nil
}

func (taskService TaskService) sendTaskWithCredentials(ctx context.Context, tasksCommonFields models.Task, managedEndpointIDs []apiModels.ManagedEndpoint, payload apiModels.ExecutionPayload, executionURL string, taskInstance models.TaskInstance) error {
	decrypted, err := taskService.encryptionService.Decrypt(*tasksCommonFields.Credentials)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantDecrypt, "SendTaskOnExecution: could't decrypt credentials for task with ID %v . err: %s", tasksCommonFields.ID, err.Error())
		return err
	}
	tasksCommonFields.Credentials = &decrypted

	for i := range managedEndpointIDs {
		payload.ManagedEndpoints = []apiModels.ManagedEndpoint{managedEndpointIDs[i]}
		endpointID, err := gocql.ParseUUID(managedEndpointIDs[i].ID)
		if err != nil {
			logger.Log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "SendTaskOnExecution: invalid format of endpointID %v. err: %s", endpointID, err.Error())
			return err
		}

		if tasksCommonFields.Credentials != nil {
			encrypted, err := taskService.agentEncryptionService.Encrypt(ctx, endpointID, *tasksCommonFields.Credentials)
			if err != nil {
				logger.Log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "SendTaskOnExecution: could't encrypt credentials for task with ID %v and endpointID %v. err: %s", tasksCommonFields.ID, endpointID, err.Error())
				return err
			}
			payload.Credentials = &encrypted
		}

		go taskService.sendTaskOnExecutionREST(ctx, executionURL, payload, taskInstance, tasksCommonFields)
	}
	return nil
}

func (taskService TaskService) sendTaskOnExecutionREST(ctx context.Context, executionURL string, payload apiModels.ExecutionPayload, taskInstance models.TaskInstance, tasksCommonFields models.Task) {
	requestBody, err := json.Marshal(&payload)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantMarshall, "sendTaskOnExecutionREST: error marshaling payload: ", err)
		return
	}

	logger.Log.DebugfCtx(ctx, "Sending payload on execution %v", string(requestBody))
	response, err := common.HTTPRequestWithRetry(ctx, taskService.httpClient, http.MethodPost, executionURL, requestBody)
	if err != nil {
		var status int
		if response != nil {
			status = response.StatusCode
		}
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantExecuteTasks, "taskService.sendTaskOnExecutionREST: Error while sending task on execution err=%v, status code:%v", err, status)

		taskService.updateExecutionResultWithErr(ctx, taskInstance, payload, tasksCommonFields.PartnerID, executionURL)
	}

	err = common.CloseRespBody(response)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "SendTaskOnExecution: error while closing response body: ", err)
	}
}

func (taskService TaskService) saveExecutionExpiration(
	ctx context.Context,
	partnerID string,
	taskInstance models.TaskInstance,
	expectedExecutionTimeSec int,
) error {
	var expectedExecDuration = time.Duration(expectedExecutionTimeSec+config.Config.HTTPClientResultsTimeoutSec) * time.Second

	managedEndpointIDs := make([]gocql.UUID, 0)
	for me, status := range taskInstance.Statuses {
		if status == statuses.TaskInstanceRunning {
			managedEndpointIDs = append(managedEndpointIDs, me)
		}
	}

	executionExpiration := models.ExecutionExpiration{
		ExpirationTimeUTC:  time.Now().Add(expectedExecDuration).Truncate(time.Minute),
		PartnerID:          partnerID,
		TaskInstanceID:     taskInstance.ID,
		ManagedEndpointIDs: managedEndpointIDs,
	}

	if err := taskService.executionExpirationPersistence.InsertExecutionExpiration(ctx, executionExpiration); err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantSaveExecutionExpiration, "taskService.saveExecutionExpiration: can't store new execution expiration, err=%v", err)
		return err
	}
	return nil
}

func getExecutionURL(taskType, partnerID string) (string, error) {
	if serviceELB, ok := config.Config.TaskTypes[taskType]; ok {
		return fmt.Sprintf("%s/partners/%s/executions", serviceELB, partnerID), nil
	}
	return "", fmt.Errorf("unknown execution MS type: %s", taskType)
}

func (taskService TaskService) updateExecutionResultWithErr(
	ctx context.Context,
	taskInstance models.TaskInstance,
	payload apiModels.ExecutionPayload,
	partnerID string,
	executionURL string,
) {
	executionResults := make([]apiModels.ExecutionResult, 0, len(taskInstance.Statuses))
	for _, target := range payload.ManagedEndpoints {
		executionResult := apiModels.ExecutionResult{
			EndpointID:       target.ID,
			CompletionStatus: "Failed",
			ErrorDetails:     "Error while sending task on execution to " + executionURL,
			UpdateTime:       time.Now().UTC().Truncate(time.Millisecond),
		}
		executionResults = append(executionResults, executionResult)
	}

	err := taskService.executionResultUpdateService.ProcessExecutionResults(ctx, partnerID, taskInstance.ID, executionResults...)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't update execution results, err: %v", err)
	}
}

// StopTaskInstanceExecutions stops future executions for task instances on managed endpoints with scheduled or pending status
func (taskService TaskService) StopTaskInstanceExecutions(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		InstanceIDs []gocql.UUID `json:"instanceIDs"`
	}

	if err := validator.ExtractStructFromRequest(r, &payload); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "TaskService.StopTaskInstanceExecutions: error while unmarshaling request body. Err=%v", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	ctx := r.Context()
	taskInstances, err := taskService.taskInstancePersistence.GetByIDs(ctx, payload.InstanceIDs...)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, "TaskService.StopTaskInstanceExecutions: can't get a Task Instances by TaskInstanceIDs %v. Err=%v", payload.InstanceIDs, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	if len(taskInstances) == 0 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskInstances, "TaskService.StopTaskInstanceExecutions: error can't get Task Instances by id. Err=%v", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantGetTaskInstances)
		return
	}

	for i := range taskInstances {
		// isStopped describes if future execution was stopped at least on one device
		isStopped := false
		for deviceID := range taskInstances[i].Statuses {
			if taskInstances[i].Statuses[deviceID] == statuses.TaskInstancePending {
				taskInstances[i].Statuses[deviceID] = statuses.TaskInstanceStopped
				isStopped = true
			}
		}

		if isStopped {
			taskInstances[i].OverallStatus = statuses.TaskInstanceStopped
		}
	}

	for _, instance := range taskInstances {
		if err = taskService.taskInstancePersistence.Insert(ctx, instance); err != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantUpdateTaskInstances, "TaskService.StopTaskInstanceExecutions: error while updating Task Instance (ID: %v). Err=%v", instance.ID, err)
			common.SendInternalServerError(w, r, errorcode.ErrorCantUpdateTaskInstances)
			return
		}
	}

	logger.Log.InfofCtx(r.Context(), "TaskService.StopTaskInstanceExecutions: Tasks' executions are successfully stopped")
	common.RenderJSON(w, struct {
		Status string
	}{Status: "Complete"})
}

// filterTasksByNOCAccsess filters the tasks by NOC accesses
func filterTasksByNOCAccsess(tasks []models.Task, isNOCUser bool) (filteredTasks []models.Task) {
	for _, task := range tasks {
		if !isNOCUser && task.IsRequireNOCAccess {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}
	return
}

// filteredByNOCMap returns filtered tasks map by taskID and NOC
func filteredByNOCMap(tasks []models.Task, isNOCUser bool) (filteredTasks map[gocql.UUID]models.Task) {
	filteredTasks = make(map[gocql.UUID]models.Task)
	for _, task := range tasks {
		if !isNOCUser && task.IsRequireNOCAccess {
			continue
		}

		gotTask, ok := filteredTasks[task.ID]
		if !ok {
			filteredTasks[task.ID] = task
			continue
		}

		// find the nearest ModifiedAt
		if task.ModifiedAt.After(gotTask.ModifiedAt) {
			filteredTasks[task.ID] = task
		}
	}
	return
}

func (taskService TaskService) setTaskToCache(task models.Task) (err error) {
	keyForCache := []byte("TKS_TASKS_BY_ID_" + task.ID.String())
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return
	}
	return taskService.cache.Set(keyForCache, taskBytes, 0)
}

func (taskService *TaskService) filterByResourceType(ctx context.Context, eIDs map[gocql.UUID]models.TargetType, partnerID string, requiredType integration.ResourceType) (filtered map[gocql.UUID]models.TargetType, err error) {
	var actualType integration.ResourceType
	filtered = make(map[gocql.UUID]models.TargetType)

	for id, tt := range eIDs {
		actualType, err = taskService.assetsService.GetResourceTypeByEndpointID(ctx, partnerID, id)
		if err != nil {
			return
		}

		if actualType == requiredType {
			filtered[id] = tt
		}
	}
	return
}

func isUserSites(targets []string, userSites []int64) (string, bool) {
	var has bool
	for _, id := range targets {
		has = false
		for _, usID := range userSites {
			if strconv.FormatInt(usID, 10) == id {
				has = true
				break
			}
		}

		if !has {
			return id, false
		}

	}
	return "", true

}
