package templates

import (
	"net/http"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// TemplateService represents a TemplateDetails Service instance
type TemplateService struct {
	cache       models.TemplateCache
	userService user.Service
	httpClient  *http.Client
}

// NewTemplateService creates a new TemplateDetails Service initialized with the Templates Repository
func NewTemplateService(cache models.TemplateCache, userService user.Service, httpClient *http.Client) TemplateService {
	return TemplateService{
		cache:       cache,
		userService: userService,
		httpClient:  httpClient,
	}
}

const (
	currentUserLog = "Current user: %+v"
	partnerIDKey   = "partnerID"
)

// GetAll returns all task definition templates by partner ID
func (templateService TemplateService) GetAll(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)[partnerIDKey]
	currentUser := templateService.userService.GetUser(r, templateService.httpClient)

	templates, err := templateService.cache.GetAllTemplates(r.Context(), partnerID, currentUser.HasNOCAccess())
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTemplatesForExecutionMS, "TemplateService.GetAll: can't get templates, err: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTemplatesForExecutionMS)
		return
	}

	logger.Log.DebugfCtx(r.Context(), "TemplateService.Read: successfully returned list of templates.")
	common.RenderJSON(w, templates)
}

// GetByType returns all task definition templates by partner ID and task type
func (templateService TemplateService) GetByType(w http.ResponseWriter, r *http.Request) {
	var (
		partnerID   = mux.Vars(r)[partnerIDKey]
		taskType    = mux.Vars(r)["type"]
		currentUser = templateService.userService.GetUser(r, templateService.httpClient)
	)

	if _, ok := config.Config.TaskTypes[taskType]; !ok {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorWrongTemplateType, "TemplateService.Read: Wrong template type %v.", taskType)
		common.SendBadRequest(w, r, errorcode.ErrorWrongTemplateType)
		return
	}

	templates, err := templateService.cache.GetByType(r.Context(), partnerID, taskType, currentUser.HasNOCAccess())
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTemplatesForExecutionMS, "TemplateService.GetByType: can't get templates, err: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTemplatesForExecutionMS)
		return
	}
	logger.Log.DebugfCtx(r.Context(), "TemplateService.Read: successfully returned list of templates")
	common.RenderJSON(w, templates)
}

// GetByOriginID returns task definition template by TaskType and Origin ID
func (templateService TemplateService) GetByOriginID(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)[partnerIDKey]
	currentUser := templateService.userService.GetUser(r, templateService.httpClient)

	originID, err := gocql.ParseUUID(mux.Vars(r)["originID"])
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTemplatesForExecutionMS, "TemplateService.GetByOriginID: can't get template, err: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantGetTemplatesForExecutionMS)
		return
	}

	template, err := templateService.cache.GetByOriginID(r.Context(), partnerID, originID, currentUser.HasNOCAccess())
	if err != nil {
		switch err.(type) {
		case models.TemplateNotFoundError:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskDefinitionTemplate, "TemplateService.GetByOriginID: Template with ID %v is not found", originID)
			common.SendNotFound(w, r, errorcode.ErrorCantGetTaskDefinitionTemplate)
			return
		default:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTemplatesForExecutionMS,"TemplateService.GetByOriginID: can't get template, err: %v", err)
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetTemplatesForExecutionMS)
			return
		}
	}
	logger.Log.DebugfCtx(r.Context(), "TemplateService.Read: successfully returned template")
	common.RenderJSON(w, template)
}
