package taskDefinitions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/validator"
)

const customCategory = "Custom"

//EncryptionService - provide us interface for encryption of credentials
type EncryptionService interface {
	Encrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
	Decrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
}

// TaskDefinitionService represents Task Definition Service
type TaskDefinitionService struct {
	taskDefinitionPersistence models.TaskDefinitionPersistence
	templatesCache            models.TemplateCache
	httpClient                *http.Client
	userService               user.Service
	encryptionService         EncryptionService
}

// NewTaskDefinitionService creates new service with cache and persistence instances.
func NewTaskDefinitionService(persistence models.TaskDefinitionPersistence, cache models.TemplateCache, httpClient *http.Client, userService user.Service, encryptionService EncryptionService) TaskDefinitionService {
	return TaskDefinitionService{
		taskDefinitionPersistence: persistence,
		templatesCache:            cache,
		httpClient:                httpClient,
		userService:               userService,
		encryptionService:         encryptionService,
	}
}

// Create saves new Task Definition.
func (td TaskDefinitionService) Create(w http.ResponseWriter, r *http.Request) {
	initiatedBy := r.Header.Get(common.InitiatedByHeader)

	taskDefDetails, err := td.prepareUpsert(w, r, initiatedBy)
	if err != nil {
		return
	}

	if taskDefDetails.Credentials != nil {
		encrypted, err := td.encryptionService.Encrypt(*taskDefDetails.Credentials)
		if err != nil {
			err = fmt.Errorf("can't encrypt credentials. err: %s", err.Error())
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantEcryptCredentials,"TaskDefinitionService.Create: %s", err.Error())
			common.SendInternalServerError(w, r, errorcode.ErrorCantEcryptCredentials)
			return
		}
		taskDefDetails.Credentials = &encrypted
	}

	if exists := td.taskDefinitionPersistence.Exists(r.Context(), taskDefDetails.PartnerID, taskDefDetails.Name); exists {
		err = fmt.Errorf(alreadyExistsFmt, taskDefDetails.Name, taskDefDetails.PartnerID)
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionExists,"TaskDefinitionService.Create: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskDefinitionExists)
		return
	}

	taskDefDetails.ID = gocql.TimeUUID()
	taskDefDetails.CreatedBy = initiatedBy
	taskDefDetails.CreatedAt = time.Now().Truncate(time.Second).UTC()

	err = td.taskDefinitionPersistence.Upsert(r.Context(), taskDefDetails)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(),  errorcode.ErrorCantSaveTaskDefinitionToDB, "TaskDefinitionService.Insert: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskDefinitionToDB)
		return
	}

	logger.Log.InfofCtx(r.Context(), "New task definition with ID=%v is created.", taskDefDetails.ID)
	common.RenderJSONCreated(w, taskDefDetails)
}

const currentUserFmt  = "Current user: %+v"

// GetByID returns Task Definition Details by partner and ID.
func (td TaskDefinitionService) GetByID(w http.ResponseWriter, r *http.Request) {
	taskDefDetails, err := td.getByID(w, r)
	if err != nil {
		return
	}

	currentUser := td.userService.GetUser(r, td.httpClient)

	template, err := td.templatesCache.GetByOriginID(r.Context(), taskDefDetails.PartnerID, taskDefDetails.OriginID, currentUser.HasNOCAccess())
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskDefinitionTemplate,"TaskDefinitionService.GetByID: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskDefinitionTemplate)
		return
	}

	taskDefDetails.JSONSchema = template.JSONSchema
	taskDefDetails.UISchema = template.UISchema
	taskDefDetails.Engine = template.Engine

	logger.Log.DebugfCtx(r.Context(), "Task definition by ID=%v is found.", taskDefDetails.ID)
	common.RenderJSON(w, taskDefDetails)
}

// GetByPartnerID returns all Task Definitions by partner ID.
func (td TaskDefinitionService) GetByPartnerID(w http.ResponseWriter, r *http.Request) {
	curUser := td.userService.GetUser(r, td.httpClient)
	ctx := r.Context()

	defs, err := td.taskDefinitionPersistence.GetAllByPartnerID(r.Context(), curUser.PartnerID())
	if err != nil {
		logger.Log.ErrfCtx(r.Context(),  errorcode.ErrorTaskDefinitionByPartnerNotFound,"TaskDefinitionService.GetByPartnerID: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorTaskDefinitionByPartnerNotFound)
		return
	}

	results := make([]models.TaskDefinitionDetails, 0, len(defs))
	//retrieving engine by originID
	for i := range defs {
		template, err := td.templatesCache.GetByOriginID(ctx, curUser.PartnerID(), defs[i].OriginID, curUser.HasNOCAccess())
		if err != nil {
			continue
		}

		defs[i].Engine = template.Engine
		results = append(results, defs[i])
	}

	logger.Log.DebugfCtx(r.Context(), "Successfully returned task definitions list for partner with ID %v", curUser.PartnerID())
	common.RenderJSON(w, defs)
}

// DeleteByID removes Task Definition by ID.
func (td TaskDefinitionService) DeleteByID(w http.ResponseWriter, r *http.Request) {
	deletedBy := r.Header.Get(common.InitiatedByHeader)
	if deletedBy == "" {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorUIDHeaderIsEmpty,"TaskDefinitionService.DeleteByID: uid header is empty")
		common.SendBadRequest(w, r, errorcode.ErrorUIDHeaderIsEmpty)
		return
	}

	taskDefs, err := td.getByIDs(w, r)
	if err != nil {
		return
	}

	for _, taskDef := range taskDefs {
		taskDef.UpdatedBy = deletedBy
		taskDef.UpdatedAt = time.Now().Truncate(time.Second).UTC()
		taskDef.Deleted = true

		err = td.taskDefinitionPersistence.Upsert(r.Context(), taskDef)
		if err != nil {
			logger.Log.ErrfCtx(r.Context(),  errorcode.ErrorCantDeleteTaskDefinition,"TaskDefinitionService.DeleteByID: Can't get Delete Task Definition by ID, err: %v", err)
			common.SendInternalServerError(w, r, errorcode.ErrorCantDeleteTaskDefinition)
			return
		}

		logger.Log.DebugfCtx(r.Context(), "Task definition with ID=%v is deleted successfully.", taskDef.ID)
	}

	common.SendNoContent(w)
}

const alreadyExistsFmt  = "task template name '%v' already exists for partner %v"

// UpdateByID updates Task Definition Details by partner and ID.
func (td TaskDefinitionService) UpdateByID(w http.ResponseWriter, r *http.Request) {
	initiatedBy := r.Header.Get(common.InitiatedByHeader)
	oldTaskDef, err := td.getByID(w, r)
	if err != nil {
		return
	}

	taskDefDetails, err := td.prepareUpsert(w, r, initiatedBy)
	if err != nil {
		logger.Log.WarnfCtx(r.Context(), "TaskDefinitionService.UpdateByID: couldn't collect task definitions details, err=%v", err)
		return
	}

	canBeUpdated, err := td.taskDefinitionPersistence.CanBeUpdated(r.Context(), taskDefDetails.PartnerID, taskDefDetails.Name, oldTaskDef.ID)
	if !canBeUpdated {
		err = fmt.Errorf(alreadyExistsFmt, taskDefDetails.Name, taskDefDetails.PartnerID)
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionExists,"TaskDefinitionService.UpdateByID: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskDefinitionExists)
		return
	}

	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskDefinitionToDB,"TaskDefinitionService.UpdateByID:  CanBeUpdated db %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskDefinitionToDB)
		return
	}

	taskDefDetails.ID = oldTaskDef.ID
	taskDefDetails.CreatedBy = oldTaskDef.CreatedBy
	taskDefDetails.CreatedAt = oldTaskDef.CreatedAt
	taskDefDetails.UpdatedBy = initiatedBy
	taskDefDetails.UpdatedAt = time.Now().Truncate(time.Second).UTC()

	err = td.taskDefinitionPersistence.Upsert(r.Context(), taskDefDetails)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantSaveTaskDefinitionToDB, "TaskDefinitionService.UpdateByID: Upsert %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantSaveTaskDefinitionToDB)
		return
	}

	logger.Log.DebugfCtx(r.Context(), "Task definition with ID=%v is updated, err=%v", taskDefDetails.ID, err)
	common.RenderJSONCreated(w, taskDefDetails)
}

func (td TaskDefinitionService) prepareUpsert(w http.ResponseWriter, r *http.Request, initiatedBy string) (taskDefDetails models.TaskDefinitionDetails, err error) {
	currentUser := td.userService.GetUser(r, td.httpClient)

	if initiatedBy == "" {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "TaskDefinitionService.prepareUpsert: uid header is empty")
		common.SendBadRequest(w, r, errorcode.ErrorUIDHeaderIsEmpty)
		return taskDefDetails, errors.New("uid header is empty")
	}

	err = validator.ExtractStructFromRequest(r, &taskDefDetails)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "TaskDefinitionService.prepareUpsert: ExtractStructFromRequest %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	template, err := td.templatesCache.GetByOriginID(r.Context(), currentUser.PartnerID(), taskDefDetails.OriginID, currentUser.HasNOCAccess())
	if err != nil {
		switch err.(type) {
		case models.TemplateNotFoundError:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionByPartnerNotFound,"TaskDefinitionService.prepareUpsert: %v", err)
			common.SendBadRequest(w, r, errorcode.ErrorCantGetTaskDefinitionTemplate)
			return
		default:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskDefinitionTemplate, "TaskDefinitionService.prepareUpsert: GetByOriginID %v", err)
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskDefinitionTemplate)
			return
		}
	}

	if err = validator.ValidateParametersField(template.JSONSchema, taskDefDetails.UserParameters, false); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData,"TaskDefinitionService.prepareUpsert: %v for originID %v", err, taskDefDetails.OriginID)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	if len(taskDefDetails.Categories) == 0 {
		taskDefDetails.Categories = append(taskDefDetails.Categories, customCategory)
	}

	taskDefDetails.PartnerID = currentUser.PartnerID()
	taskDefDetails.JSONSchema = template.JSONSchema
	taskDefDetails.UISchema = template.UISchema

	if exists := td.templatesCache.ExistsWithName(r.Context(), taskDefDetails.PartnerID, taskDefDetails.Name); exists {
		err = fmt.Errorf(alreadyExistsFmt, taskDefDetails.Name, taskDefDetails.PartnerID)
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionExists,"TaskDefinitionService.Create: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskDefinitionExists)
		return taskDefDetails, fmt.Errorf("already exists")
	}

	return taskDefDetails, err
}

func (td TaskDefinitionService) getByID(w http.ResponseWriter, r *http.Request) (taskDefDetails models.TaskDefinitionDetails, err error) {
	partnerID := mux.Vars(r)["partnerID"]

	definitionID, err := common.ExtractUUID("TaskDefinitionService.getByID", w, r, "definitionID")
	if err != nil {
		return
	}

	taskDefDetails, err = td.taskDefinitionPersistence.GetByID(r.Context(), partnerID, definitionID)

	if err != nil {
		switch err.(type) {
		case models.TaskDefNotFoundError:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionNotFound,"TaskDefinitionService.getByID: Task Definition with ID %v was not found for Partner: %v.",
				definitionID.String(), partnerID)
			common.SendNotFound(w, r, errorcode.ErrorTaskDefinitionNotFound)
			return
		default:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionNotFound,"TaskDefinitionService.getByID: Can't get Task Definition by ID, err: %v", err)
			common.SendInternalServerError(w, r, errorcode.ErrorTaskDefinitionNotFound)
			return
		}
	}

	return taskDefDetails, err
}

func (td TaskDefinitionService) getByIDs(w http.ResponseWriter, r *http.Request) (taskDefs []models.TaskDefinitionDetails, err error) {
	partnerID := mux.Vars(r)["partnerID"]

	definitionIDs, err := common.ExtractUUIDs(r, "definitionID")
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionIDHasBadFormat,"TaskDefinitionService.getByIDs: can't extract uuids from request. err: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorTaskDefinitionIDHasBadFormat)
		return
	}
	logger.Log.DebugfCtx(r.Context(), "TaskDefinitionService.getByIDs: Deleting '%v' for partner %v", definitionIDs, partnerID)

	for _, definitionID := range definitionIDs {
		taskDefDetails, err := td.taskDefinitionPersistence.GetByID(r.Context(), partnerID, definitionID)
		if err != nil {
			switch err.(type) {
			case models.TaskDefNotFoundError:
				logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionNotFound,"TaskDefinitionService.getByIDs: Task Definition with ID %v was not found for Partner: %v.",
					definitionID.String(), partnerID)
				common.SendNotFound(w, r, errorcode.ErrorTaskDefinitionNotFound)
				return nil, err
			default:
				logger.Log.ErrfCtx(r.Context(), errorcode.ErrorTaskDefinitionNotFound, "TaskDefinitionService.getByIDs: Can't get Task Definition by ID, err: %v", err)
				common.SendInternalServerError(w, r, errorcode.ErrorTaskDefinitionNotFound)
				return nil, err
			}
		}

		taskDefs = append(taskDefs, taskDefDetails)
	}

	return taskDefs, err
}
