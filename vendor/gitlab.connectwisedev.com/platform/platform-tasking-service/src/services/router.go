package services

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	nr "github.com/urfave/negroni"
	commonlibRest "gitlab.connectwisedev.com/platform/platform-common-lib/src/web/rest"

	accessControl "gitlab.connectwisedev.com/platform/platform-tasking-service/src/access-control"
	apiDefinitions "gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/api"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/permission_temp"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/profile"
	taskExecutionResults "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/execution-results"
	executionResultsUpdate "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/execution-results-update"
	taskCounters "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/task-counters"
	taskDefinitions "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/task-definitions"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/templates"
)

// RouterDTO is a dto structure that countains router data
type RouterDTO struct {
	TasksService                 tasks.TaskService
	TemplateService              templates.TemplateService
	TaskResultsService           taskExecutionResults.TaskResultsService
	ExecutionResultUpdateService executionResultsUpdate.ExecutionResultUpdateService
	TaskDefinitionService        taskDefinitions.TaskDefinitionService
	TaskCounterService           taskCounters.Service
	UserMD                       *user.User
	Md                           *permission_temp.Permission
	Handlers                     map[string]func(http.ResponseWriter, *http.Request)
}

// NewRouter creates a router for URL-to-service mapping
func NewRouter(
	o RouterDTO,
) *mux.Router {
	var (
		basePath       = "/tasking/v1"
		managementPath = "/tasking/v1/management"
		partnerPath    = "/partners/{partnerID}"
		router         = mux.NewRouter()
		apiRouter      = mux.NewRouter()
		api            = apiRouter.PathPrefix(basePath + partnerPath).Subrouter()
		apiManagement  = apiRouter.PathPrefix(managementPath + partnerPath).Subrouter()
		clientSite     = apiRouter.PathPrefix(basePath + partnerPath + "/clients/{clientID}/sites/{siteID}").Subrouter()
		debugRouter    = router.PathPrefix("/tasking").Subrouter()
	)

	groupMD := nr.New(o.Md)
	group := nr.New(o.UserMD, o.Md)

	router.HandleFunc("/tasking/version", commonlibRest.HandlerVersion).Methods(http.MethodGet)
	router.HandleFunc("/tasking/health", commonlibRest.HandlerHealth).Methods(http.MethodGet)
	router.HandleFunc("/tasking/v1/recalculate-counters", o.TaskCounterService.RecalculateAllCounters).Methods(http.MethodPost)

	registerProfilers(debugRouter)

	api.Handle("/tasks", group.With(nr.WrapFunc(o.TasksService.Create))).Methods(http.MethodPost)
	api.HandleFunc("/tasks/count", o.TaskCounterService.GetCountersByPartner).Methods(http.MethodGet)
	// tasksService.GetTasksSummaryData is deprecated
	api.HandleFunc("/tasks/data", o.TasksService.GetTasksSummaryData).Methods(http.MethodGet)

	// Path handled by predefined handler
	api.HandleFunc("/tasks/data/closest", o.Handlers[apiDefinitions.TasksClosestHandler]).Methods(http.MethodPost)
	api.Handle("/tasks/data/scheduled", group.With(nr.WrapFunc(o.Handlers[apiDefinitions.TasksScheduledHandler]))).Methods(http.MethodGet, http.MethodDelete) // 2.1 UI scheduled panel
	api.Handle("/tasks/data/history", group.With(nr.WrapFunc(o.Handlers[apiDefinitions.TasksHistoryHandler]))).Methods(http.MethodGet)

	api.Handle("/tasks/data/recent", groupMD.With(nr.WrapFunc(o.TasksService.LastTasks))).Methods(http.MethodGet)
	api.Handle("/tasks/data/script/{originID}", groupMD.With(nr.WrapFunc(o.TasksService.GetByOriginID))).Methods(http.MethodGet)

	var tasksID = "/tasks/{taskID}" // sonar
	api.Handle(tasksID, group.With(nr.WrapFunc(o.TasksService.Edit))).Methods(http.MethodPut)
	api.HandleFunc(tasksID, o.TasksService.GetByPartnerAndID).Methods(http.MethodGet)
	apiManagement.HandleFunc(tasksID, o.TasksService.Delete).Methods(http.MethodDelete)

	api.HandleFunc("/tasks/{taskID}/task-instances/{taskInstanceID}", o.TasksService.DeleteExecutions).Methods(http.MethodDelete)
	api.Handle("/tasks/{taskID}/data", groupMD.With(nr.WrapFunc(o.TasksService.SubRecent))).Methods(http.MethodGet)
	api.Handle("/tasks/{taskID}/enable", groupMD.With(nr.WrapFunc(o.TasksService.EnableTaskForAllTargets))).Methods(http.MethodPut)
	api.HandleFunc("/tasks/{taskID}/enable/targets", o.TasksService.EnableTaskForSelectedTargets).Methods(http.MethodPut)
	api.HandleFunc("/tasks/{taskID}/postpone", o.TasksService.PostponeNearestExecution).Methods(http.MethodPut)
	api.HandleFunc("/tasks/{taskID}/managed-endpoints/{managedEndpointID}/postpone", o.TasksService.PostponeDeviceNearestExecution).Methods(http.MethodPut)

	api.HandleFunc("/tasks/{taskID}/cancel", o.TasksService.CancelNearestExecution).Methods(http.MethodPut)
	api.HandleFunc("/tasks/{taskID}/managed-endpoint/{managedEndpointID}/cancel", o.TasksService.CancelNearestExecutionForEndpoint).Methods(http.MethodPut)

	api.HandleFunc("/tasks/{taskID}/task-instances/count", o.TasksService.TaskInstancesCountByTaskID).Methods(http.MethodGet)

	api.HandleFunc("/tasks/task-instances/stop", o.TasksService.StopTaskInstanceExecutions).Methods(http.MethodPut)

	api.HandleFunc("/tasks/managed-endpoints/{managedEndpointID}/count", o.TaskCounterService.GetCountersByPartnerAndEndpoint).Methods(http.MethodGet)
	api.HandleFunc("/tasks/managed-endpoints/{managedEndpointID}", o.TasksService.GetByPartnerAndManagedEndpointID).Methods(http.MethodGet)

	api.HandleFunc("/task-execution-results/tasks/{taskID}/managed-endpoints/{managedEndpointID}/history", o.TaskResultsService.History).Methods(http.MethodGet)
	api.HandleFunc("/task-execution-results/task-instances/{taskInstanceID}", o.ExecutionResultUpdateService.UpdateTaskAndTaskInstanceStatuses).Methods(http.MethodPost)
	api.Handle("/task-execution-results/managed-endpoints/{managedEndpointID}", groupMD.With(nr.WrapFunc(o.TaskResultsService.Get))).Methods(http.MethodGet)
	api.HandleFunc("/task-execution-results/managed-endpoints/{managedEndpointID}/task-instances/{taskInstanceID}/logs/stdout", o.TaskResultsService.GetTaskExecutionStdoutLogs).Methods(http.MethodGet)
	api.HandleFunc("/task-execution-results/managed-endpoints/{managedEndpointID}/task-instances/{taskInstanceID}/logs/stderr", o.TaskResultsService.GetTaskExecutionStderrLogs).Methods(http.MethodGet)

	var definitionsID = "/task-definitions/{definitionID}" // sonar
	api.Handle("/task-definitions", group.With(nr.WrapFunc(o.TaskDefinitionService.Create))).Methods(http.MethodPost)
	api.Handle("/task-definitions", groupMD.With(nr.WrapFunc(o.TaskDefinitionService.GetByPartnerID))).Methods(http.MethodGet)
	api.Handle(definitionsID, groupMD.With(nr.WrapFunc(o.TaskDefinitionService.GetByID))).Methods(http.MethodGet)
	api.Handle(definitionsID, group.With(nr.WrapFunc(o.TaskDefinitionService.DeleteByID))).Methods(http.MethodDelete)
	api.Handle(definitionsID, group.With(nr.WrapFunc(o.TaskDefinitionService.UpdateByID))).Methods(http.MethodPut)

	api.HandleFunc("/tasks-definition-templates", o.TemplateService.GetAll).Methods(http.MethodGet)
	api.HandleFunc("/tasks-definition-templates/{type}", o.TemplateService.GetByType).Methods(http.MethodGet)
	api.HandleFunc("/tasks-definition-templates/{type}/{originID}", o.TemplateService.GetByOriginID).Methods(http.MethodGet)

	// these urls used in legacy portal migration
	api.HandleFunc("/legacy-info/script", o.Handlers[apiDefinitions.GetLegacyScriptByPartnerHandler]).Methods(http.MethodGet)
	api.HandleFunc("/legacy-info/script", o.Handlers[apiDefinitions.InsertScriptInfoHandler]).Methods(http.MethodPost)
	api.HandleFunc("/legacy-info/script/{scriptID}", o.Handlers[apiDefinitions.GetLegacyScriptByScriptHandler]).Methods(http.MethodGet)

	api.HandleFunc("/legacy-info/job", o.Handlers[apiDefinitions.JobGetByPartnerHandler]).Methods(http.MethodGet)
	api.HandleFunc("/legacy-info/job", o.Handlers[apiDefinitions.JobInsertHandler]).Methods(http.MethodPost)
	api.HandleFunc("/legacy-info/job/{jobID}", o.Handlers[apiDefinitions.JobGetByJobIDHandler]).Methods(http.MethodGet)

	api.HandleFunc("/triggers", o.TasksService.GetTriggersList).Methods(http.MethodGet)
	router.HandleFunc("/tasking/v1/triggers", o.TasksService.UploadTriggerDefinitions).Methods(http.MethodPut)

	clientSite.HandleFunc("/endpoints/{endpointID}/triggers/{triggerType}/execute-trigger", o.TasksService.ExecuteTrigger).Methods(http.MethodPost)

	corsOptions := cors.Options{
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
	}

	commonMD := nr.New(
		cors.New(corsOptions),
		nr.HandlerFunc(accessControl.IsPartnerAuthorizedForRoute),
		nr.Wrap(apiRouter),
	)
	router.PathPrefix(basePath).Handler(commonMD)

	return router
}

func registerProfilers(router *mux.Router) {
	router.Handle("/debug/pprof/heap", profile.Handler("heap"))
	router.Handle("/debug/pprof/goroutine", profile.Handler("goroutine"))
	router.Handle("/debug/pprof/allocs", profile.Handler("allocs"))
	router.Handle("/debug/pprof/block", profile.Handler("block"))
	router.Handle("/debug/pprof/cmdline", profile.Handler("cmdline"))
	router.Handle("/debug/pprof/mutex", profile.Handler("mutex"))
	router.Handle("/debug/pprof/profile", profile.Handler("profile"))
	router.Handle("/debug/pprof/threadcreate", profile.Handler("threadcreate"))
	router.Handle("/debug/pprof/trace", profile.Handler("trace"))
}
