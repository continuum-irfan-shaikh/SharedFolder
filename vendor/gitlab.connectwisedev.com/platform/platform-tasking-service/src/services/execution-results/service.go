// Package taskExecutionResults represents a mocked service for UI requests intended to fill the Task Execution Results Table
// see the screenshot https://continuum.atlassian.net/wiki/spaces/C2E/pages/223371940/Juno+Scripting+The+Aftermath?preview=/223371940/223387401/Confirmation_3.png
// The request includes PartnerID and ManagedEndpointID and expects to retrieve
package taskExecutionResults

import (
	"net/http"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/gorilla/mux"
)

const defaultHistoryCount = 10

// TaskResultsService represents
// a Task Results Service incorporated Task Execution Results Repository
type TaskResultsService struct {
	executionResultViewRepo models.ExecutionResultViewPersistence
	executionResultRepo     models.ExecutionResultPersistence
	userService             user.Service
	httpClient              *http.Client
}

// NewTaskResultsService returns an instance of TaskResultsService initialized with the repository
func NewTaskResultsService(
	executionResultViewRepo models.ExecutionResultViewPersistence,
	executionResultPersistence models.ExecutionResultPersistence,
	userService user.Service,
	httpClient *http.Client,
) TaskResultsService {
	return TaskResultsService{
		executionResultViewRepo: executionResultViewRepo,
		executionResultRepo:     executionResultPersistence,
		httpClient:              httpClient,
		userService:             userService,
	}
}

// Get returns an execution results of all tasks running on the specified Endpoint for the specified Partner
func (taskResultsService TaskResultsService) Get(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)["partnerID"]
	ctx := r.Context()

	endpointID, err := common.ExtractUUID("TaskResultsService.Get", w, r, "managedEndpointID")
	if err != nil {
		logger.Log.WarnfCtx(ctx, "Results.Get: cant parse endpointID from req, err: %v", err)
		return
	}
	count, _, err := common.ExtractOptionalCount("TaskResultsService.Get", w, r)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "Results.Get: cant parse count from req, err: %v", err)
		return
	}

	currentUser := taskResultsService.userService.GetUser(r, taskResultsService.httpClient)

	executionResultsView, err := taskResultsService.executionResultViewRepo.Get(ctx, partnerID, endpointID, count, currentUser.HasNOCAccess())
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskExecutionResultsForManagedEndpoint, "TaskResultsService.Get: can't get all task execution results for endpoint (UUID=%s) and partner (ID=%s). err=%v", endpointID, partnerID, err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskExecutionResultsForManagedEndpoint)
		return
	}

	logger.Log.DebugfCtx(ctx, "TaskResultsService.Get: successfully returned list of task execution results for endpoint (UUID=%s) and partner (ID=%s)", endpointID, partnerID)
	common.RenderJSON(w, executionResultsView)
}

const managedEndpointIDKey = "managedEndpointID"
const methodName   =  "TaskResultsService.History"
// History returns a history of execution results on the specified Task for the specified Partner
func (taskResultsService TaskResultsService) History(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)["partnerID"]
	taskID, err := common.ExtractUUID(methodName, w, r, "taskID")
	if err != nil {
		logger.Log.WarnfCtx(r.Context(), "TaskResultsService.History: can't parse Task ID (UUID=%s),  partner (ID=%s). err=%v", taskID, partnerID, err)
		return
	}
	endpointID, err := common.ExtractUUID(methodName, w, r, managedEndpointIDKey)
	if err != nil {
		logger.Log.WarnfCtx(r.Context(), "TaskResultsService.History: can't parse endpoint ID (UUID=%s),  partner (ID=%s). err=%v", endpointID, partnerID, err)
		return
	}
	count, isSpecified, err := common.ExtractOptionalCount(methodName, w, r)
	if err != nil {
		return
	}
	if !isSpecified {
		count = defaultHistoryCount
	}

	currentUser := taskResultsService.userService.GetUser(r, taskResultsService.httpClient)

	executionResultsView, err := taskResultsService.executionResultViewRepo.History(r.Context(), partnerID, taskID, endpointID, count, currentUser.HasNOCAccess())
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskExecutionResults, "TaskResultsService.History: can't get all task execution results for Task ID (UUID=%s), endpointID (UUID=%s) and partner (ID=%s). err=%v", taskID, endpointID, partnerID, err)
		switch err.(type) {
		case models.TaskNotFoundError:
			common.SendNotFound(w, r, errorcode.ErrorCantGetTaskExecutionResults)
		default:
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskExecutionResults)
		}
		return
	}

	logger.Log.DebugfCtx(r.Context(), "TaskResultsService.History: successfully returned a history of last %v execution results for Task ID (UUID=%s) and partner (ID=%s)", count, taskID, partnerID)
	common.RenderJSON(w, executionResultsView)
}

//GetTaskExecutionStdoutLogs returns stdout logs about task execution results on specified Endpoint and Task instance for the specified Partner
func (taskResultsService TaskResultsService) GetTaskExecutionStdoutLogs(responseWriter http.ResponseWriter, request *http.Request) {
	getTaskExecutionLog("GetTaskExecutionStdoutLogs", taskResultsService.executionResultRepo)(responseWriter, request)
}

//GetTaskExecutionStderrLogs returns stderr logs about task execution results on specified Endpoint and Task instance for the specified Partner
func (taskResultsService TaskResultsService) GetTaskExecutionStderrLogs(responseWriter http.ResponseWriter, request *http.Request) {
	getTaskExecutionLog("GetTaskExecutionStderrLogs", taskResultsService.executionResultRepo)(responseWriter, request)
}

const taskResultsService  = "TaskResultsService."

// getTaskExecutionLog returns function to receive task execution logs
func getTaskExecutionLog(method string, executionResultPersistence models.ExecutionResultPersistence) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		endpointID, err := common.ExtractUUID(taskResultsService+method, responseWriter, request, managedEndpointIDKey)
		if err != nil {
			logger.Log.WarnfCtx(request.Context(), "getTaskExecutionLog: can't parse endpointID, err: %v", err)
			return
		}
		taskInstanceID, err := common.ExtractUUID(taskResultsService+method, responseWriter, request, "taskInstanceID")
		if err != nil {
			logger.Log.WarnfCtx(request.Context(), "getTaskExecutionLog: can't parse taskInstanceID, err: %v", err)
			return
		}
		logger.Log.DebugfCtx(request.Context(), "getTaskExecutionLog: taskInstanceID %v, endpointID %v", taskInstanceID, endpointID)

		executionResults, err := executionResultPersistence.GetByTargetAndTaskInstanceIDs(endpointID, taskInstanceID)
		if err != nil {
			logger.Log.ErrfCtx(request.Context(), errorcode.ErrorCantDecodeInputData,taskResultsService+method+": can't get script execution result by endpoint ID (UUID=%s) and task instance ID (UUID=%s). err=%v", endpointID, taskInstanceID, err)
			common.SendInternalServerError(responseWriter, request, errorcode.ErrorCantDecodeInputData)
			return
		}

		var output string
		if len(executionResults) > 0 {
			if method == "GetTaskExecutionStderrLogs" {
				output = executionResults[0].StdErr
			} else {
				output = executionResults[0].StdOut
			}
		}
		common.RenderJSON(responseWriter, output)
	}
}
