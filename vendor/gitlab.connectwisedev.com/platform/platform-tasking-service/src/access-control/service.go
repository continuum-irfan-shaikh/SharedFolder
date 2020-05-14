package accessControl

import (
	"context"
	"net/http"

	commonLibEntitlement "gitlab.connectwisedev.com/platform/platform-common-lib/src/entitlement"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/gorilla/mux"
)

// This group of constants describes features for Tasking/Scripting solution
const (
	TaskingDynamicGroupsFeature = "TASKING_DYNAMIC_GROUPS"
	TaskingSitesFeature         = "TASKING_SITES"
	CustomTasksFeature          = "AUTOMATION_CUSTOM_TASKS"
)

var (
	entitlementService commonLibEntitlement.Service
	entitlementRules   = map[string]func(*models.Task, string) bool{
		TaskingDynamicGroupsFeature: isPartnerAuthorizedToRunTaskOnDynamicGroup,
		TaskingSitesFeature:         isPartnerAuthorizedToRunTaskOnSite,
	}
)

// Load creates new entitlementService
func Load(httpClient *http.Client) {
	entitlementService = commonLibEntitlement.NewEntitlementService(
		httpClient,
		config.Config.EntitlementMsURL,
		config.Config.EntitlementCacheSettings.DataTTLSec,
		config.Config.EntitlementCacheSettings.Size,
	)
}

func isPartnerAuthorizedToRunTaskOnDynamicGroup(task *models.Task, featureName string) bool {
	if len(task.Targets.IDs) == 0 && len(task.TargetsByType[models.DynamicGroup]) == 0 {
		return true
	}

	if len(task.Targets.IDs) != 0 && task.Targets.Type != models.DynamicGroup {
		return true
	}

	return entitlementService.IsPartnerAuthorized(task.PartnerID, featureName)
}

func isPartnerAuthorizedToRunTaskOnSite(task *models.Task, featureName string) bool {
	if len(task.Targets.IDs) == 0 &&
		(len(task.TargetsByType[models.Site]) == 0 && len(task.TargetsByType[models.DynamicSite]) == 0) {
		return true
	}

	if len(task.Targets.IDs) != 0 && task.Targets.Type != models.Site && task.Targets.Type != models.DynamicSite {
		return true
	}

	return entitlementService.IsPartnerAuthorized(task.PartnerID, featureName)
}

// IsPartnerAuthorizedToRunTask checks if the Partner has access to run (or create) specified task
func IsPartnerAuthorizedToRunTask(task *models.Task) bool {
	for featureName, rule := range entitlementRules {
		if !rule(task, featureName) {
			return false
		}
	}
	return true
}

// IsPartnerAuthorizedForRoute checks if the Partner which presented in the request is entitled to use Feature for route
// ATTENTION! DO NOT USE IT BEFORE router.Handle APPLIED
func IsPartnerAuthorizedForRoute(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	partnerID := mux.Vars(r)["partnerID"]
	// NOTE: Panic when use before the router.Handle applied
	route, err := mux.CurrentRoute(r).GetPathTemplate()
	if err != nil {
		logger.Log.WarnfCtx(r.Context(), "IsPartnerAuthorizedForRoute: cannot GetPathTemplate: %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	featureName, ok := config.Config.FeaturesForRoutes[route]
	if ok && !entitlementService.IsPartnerAuthorized(partnerID, featureName) {
		// according to https://continuum.atlassian.net/browse/RMM-30976
		logger.Log.WarnfCtx(r.Context(), "Access Control IsPartnerAuthorizedForRoute: partnerID [ID=%s] has no feature [%s]. Access denied", partnerID, featureName)
		common.SendForbidden(w, r, errorcode.ErrorAccessDenied)
		return
	}

	logger.Log.DebugfCtx(r.Context(), "Access Control IsPartnerAuthorizedForRoute: partnerID [ID=%s] is in ACL. Access allowed", partnerID)
	next(w, r)
}

// CustomScriptsMiddleware md for custom scripts
type CustomScriptsMiddleware struct{}

// NewCustomScriptsMiddlware returns new CustomScriptsMiddleware
func NewCustomScriptsMiddlware() *CustomScriptsMiddleware {
	return &CustomScriptsMiddleware{}
}

// ServeHTTP transfers PartnerID from URL to Context
func (md *CustomScriptsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	partnerID := mux.Vars(r)["partnerID"]
	if entitlementService.IsPartnerAuthorized(partnerID, CustomTasksFeature) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, CustomTasksFeature, true)
		r = r.WithContext(ctx)
	}
	logger.Log.DebugfCtx(r.Context(), "Access Control CustomScriptsMiddleware: partnerID [ID=%s] is in ACL. Access allowed", partnerID)
	next(w, r)
}
