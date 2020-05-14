// Package executionResultsUpdate receives execution result from scripting service
package executionResultsUpdate

//go:generate mockgen -destination=../../mocks/mocks-repository/task_exec_history_mock.go -package=mockrepositories -source=./service.go TaskExecutionHistoryRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	taskResultWebhook "gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/task-result-webhook"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/validator"
)

// TaskExecutionHistoryRepo - interface to perform actions with task execution history database
type TaskExecutionHistoryRepo interface {
	Insert(entities.TaskExecHistory) error
}

const (
	monthFormat = "01"
	dateFormat  = "2006-01-02"

	failedByTimeout = "Failed by timeout"
)

// ExecutionResultUpdateService represents Execution Results Update Service
type ExecutionResultUpdateService struct {
	executionResultPersistence models.ExecutionResultPersistence
	taskRepo                   models.TaskPersistence
	taskInstanceRepo           models.TaskInstancePersistence
	taskExecHistory            TaskExecutionHistoryRepo
	assetService               integration.Asset
}

// NewExecutionResultUpdateService creates new Execution Results Update Service with
func NewExecutionResultUpdateService(er models.ExecutionResultPersistence, tp models.TaskPersistence, ti models.TaskInstancePersistence, th TaskExecutionHistoryRepo, asset integration.Asset) (service ExecutionResultUpdateService) {
	return ExecutionResultUpdateService{
		executionResultPersistence: er,
		taskRepo:                   tp,
		taskInstanceRepo:           ti,
		taskExecHistory:            th,
		assetService:               asset,
	}
}

// UpdateTaskAndTaskInstanceStatuses updates task and taskInstances statuses
func (s *ExecutionResultUpdateService) UpdateTaskAndTaskInstanceStatuses(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)["partnerID"]
	taskInstanceID, err := common.ExtractUUID("UpdateTaskAndTaskInstanceStatuses", w, r, "taskInstanceID")
	if err != nil {
		logger.Log.WarnfCtx(r.Context(), "UpdateTaskAndTaskInstanceStatuses: cant parse taskInstanceID from req, err: %v", err)
		return
	}

	var results []tasking.ExecutionResult
	if results, err = s.ExtractFromRequest(r); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "UpdateTaskAndTaskInstanceStatuses: Error while retrieving execution results. Err: %v", err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	ctx := r.Context()
	err = s.ProcessExecutionResults(ctx, partnerID, taskInstanceID, results...)
	if err != nil {
		switch err.(type) {
		case errorcode.BadRequestErr:
			e := err.(errorcode.BadRequestErr)
			logger.Log.ErrfCtx(ctx, e.ErrorCode, e.LogMessage)
			common.SendBadRequest(w, r, e.ErrorCode)
		case errorcode.InternalServerErr:
			e := err.(errorcode.InternalServerErr)
			logger.Log.ErrfCtx(ctx, e.ErrorCode, e.LogMessage)
			common.SendInternalServerError(w, r, e.ErrorCode)
		}
		return
	}

	logger.Log.DebugfCtx(r.Context(), "UpdateTaskAndTaskInstanceStatuses: statuses updated successfully.")
	common.SendCreated(w, r, errorcode.CodeUpdated)
}

// ExtractFromRequest Unmarshal Request body to slice
func (s *ExecutionResultUpdateService) ExtractFromRequest(r *http.Request) ([]tasking.ExecutionResult, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read input data: %v, err: %v", string(b), err)
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logger.Log.WarnfCtx(r.Context(), "cannot close request.Body err: %s", err)
		}
	}()

	logger.Log.DebugfCtx(r.Context(), "ExecutionResultUpdateService: received payload '%s' from %v", string(b), r.URL.String())

	// decode input data
	var results []tasking.ExecutionResult
	if err = json.Unmarshal(b, &results); err != nil {
		// try to decode into old format
		var res tasking.ExecutionResult
		if err = json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("can't decode input data: %v, err: %v", string(b), err)
		}
		results = append(results, res)
	}

	for _, res := range results {
		if err = validator.ValidateByCustomValidators(res); err != nil {
			return nil, fmt.Errorf("can't validate input data: %v, err: %v", string(b), err)
		}
	}
	return results, nil
}

func (s *ExecutionResultUpdateService) createTaskExecHistory(ctx context.Context, pData *processData, execResult models.ExecutionResult) error {
	if pData == nil {
		return errors.New("not enough data to process task execution history")
	}

	siteID, _, err := s.assetService.GetSiteIDByEndpointID(ctx, pData.PartnerID, execResult.ManagedEndpointID)
	if err != nil {
		return errors.Wrap(err, "error while getting siteID from asset service")
	}

	machineName, err := s.assetService.GetMachineNameByEndpointID(ctx, pData.PartnerID, execResult.ManagedEndpointID)
	if err != nil {
		return errors.Wrap(err, "error while getting machineName from asset service")

	}
	startedAt := pData.TaskInstance.StartedAt
	execStatus, err := statuses.TaskInstanceStatusText(execResult.ExecutionStatus)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "Error with task instance status: %v", err)
	}

	output := execResult.StdOut

	if execResult.StdErr != "" {
		output = execResult.StdErr
	}

	t := entities.TaskExecHistory{
		ExecYear:      strconv.Itoa(startedAt.Year()),
		ExecMonth:     startedAt.Format(monthFormat),
		ExecDate:      startedAt.Format(dateFormat),
		ExecTime:      startedAt,
		EndpointID:    execResult.ManagedEndpointID.String(),
		ScriptName:    pData.Name,
		ScriptID:      pData.TaskInstance.OriginID.String(),
		CompletedTime: execResult.UpdatedAt,
		ExecStatus:    execStatus,
		PartnerID:     pData.PartnerID,
		SiteID:        siteID,
		MachineName:   machineName,
		ExecBy:        pData.CreatedBy,
		Output:        output,
	}

	err = s.taskExecHistory.Insert(t)
	if err != nil {
		return errors.Wrap(err, "error while inserting history")
	}
	return nil
}

type processData struct {
	PartnerID     string
	TaskInstance  models.TaskInstance
	Name          string
	CreatedBy     string
	ResultWebhook string
	TaskID        gocql.UUID
	OriginID      gocql.UUID
}

// ProcessExecutionResults updates task and task instance accordingly to execution results.
func (s *ExecutionResultUpdateService) ProcessExecutionResults(
	ctx context.Context,
	partnerID string,
	taskInstanceID gocql.UUID,
	results ...tasking.ExecutionResult,
) error {
	var (
		executionResults   = make([]models.ExecutionResult, 0, len(results))
		managedEndpointIDs = make([]gocql.UUID, 0, len(results))
		pData              = processData{}
	)

	isFailedByTimeout := results[0].ErrorDetails == failedByTimeout

	taskInstance, err := getTaskInstance(ctx, taskInstanceID, isFailedByTimeout)
	if err != nil {
		return errorcode.NewBadRequestErr(errorcode.ErrorCantGetTaskInstances,
			fmt.Sprintf("cannot get TaskInstance: %v", err))
	}
	pData.PartnerID = partnerID
	pData.TaskInstance = taskInstance

	if taskInstance.Statuses == nil {
		taskInstance.Statuses = make(map[gocql.UUID]statuses.TaskInstanceStatus)
	}

	executionResults, managedEndpointIDs = s.processResults(ctx, results, &taskInstance)
	//update lastRunTime for instance
	taskInstance.LastRunTime = time.Now().UTC()

	if len(executionResults) == 0 {
		// nothing to do here
		return nil
	}

	err = s.executionResultPersistence.Upsert(ctx, partnerID, taskInstance.Name, executionResults...)
	if err != nil {
		return errorcode.NewInternalServerErr(errorcode.ErrorCantInsertData,
			fmt.Sprintf("UpdateTaskAndTaskInstanceStatuses error saving in scripting results table: %s", err.Error()))
	}

	err = models.TaskInstancePersistenceInstance.UpdateStatuses(ctx, taskInstance)
	if err != nil {
		return errorcode.NewInternalServerErr(errorcode.ErrorCantUpdateTaskInstances,
			fmt.Sprintf("error while updating taskInstance [%v]: %v", taskInstanceID, err))
	}
	// TBD remove this? if executionResults len is not 0, than len of managedEndpointIDs cant be 0 !
	if len(managedEndpointIDs) == 0 {
		// nothing to do here
		return nil
	}

	taskData, err := models.TaskPersistenceInstance.GetExecutionResultTaskData(partnerID, taskInstance.TaskID, managedEndpointIDs[0])
	if err != nil {
		return errorcode.NewBadRequestErr(errorcode.ErrorCantGetTaskExecutionResults,
			fmt.Sprintf("UpdateTasks: Error while updating the Task by ID [%v] and PartnerID [%v]. Err: %v", taskInstance.TaskID, partnerID, err))
	}

	pData.ResultWebhook = taskData.ResultWebHook
	pData.CreatedBy = taskData.CreatedBy
	pData.Name = taskData.Name
	pData.TaskID = taskInstance.TaskID
	pData.OriginID = taskInstance.OriginID

	go func(ctx context.Context, pData *processData, results []models.ExecutionResult) {
		for _, result := range results {
			if err := s.createTaskExecHistory(ctx, pData, result); err != nil {
				logger.Log.WarnfCtx(ctx, "createTaskExecHistory: %+v", err)
			}
		}
	}(ctx, &pData, executionResults)

	handleResultWebhooks(ctx, partnerID, pData, executionResults)

	return nil
}

func (s *ExecutionResultUpdateService) processResults(ctx context.Context, results []tasking.ExecutionResult, taskInstance *models.TaskInstance) (executionResults []models.ExecutionResult, managedEndpointIDs []gocql.UUID) {
	taskInstanceID := taskInstance.ID
	deviceStatuses := make(map[gocql.UUID]statuses.TaskInstanceStatus)

	for _, result := range results {
		resultEndpointUUID, err := gocql.ParseUUID(result.EndpointID)
		if err != nil {
			logger.Log.ErrfCtx(ctx, "ProcessExecutionResults: can not parse endpointID %s to uuid in ExecutionResult. Err: %s",
				result.EndpointID, err)
			continue
		}

		if taskInstance.Statuses[resultEndpointUUID] == statuses.TaskInstanceSuccess ||
			taskInstance.Statuses[resultEndpointUUID] == statuses.TaskInstanceFailed ||
			taskInstance.Statuses[resultEndpointUUID] == statuses.TaskInstanceSomeFailures ||
			taskInstance.Statuses[resultEndpointUUID] == statuses.TaskInstanceDisabled {
			// nothing to do here
			// because we've already got results previously
			continue
		}

		status, err := statuses.TaskInstanceStatusFromText(result.CompletionStatus)
		if err != nil {
			logger.Log.WarnfCtx(ctx, "UpdateTaskAndTaskInstanceStatuses: %s", err)
			continue
		}

		executionResult := models.ExecutionResult{
			TaskInstanceID:    taskInstanceID,
			ExecutionStatus:   status,
			StdErr:            result.ErrorDetails,
			StdOut:            result.ResultDetails,
			UpdatedAt:         result.UpdateTime,
			ManagedEndpointID: resultEndpointUUID,
		}
		executionResults = append(executionResults, executionResult)
		managedEndpointIDs = append(managedEndpointIDs, resultEndpointUUID)

		deviceStatuses[resultEndpointUUID] = status
		if status == statuses.TaskInstanceSuccess {
			taskInstance.SuccessCount++
		} else {
			taskInstance.FailureCount++
		}
	}
	taskInstance.Statuses = deviceStatuses
	return executionResults, managedEndpointIDs
}

func getTaskInstance(ctx context.Context, taskInstanceID gocql.UUID, isFailedByTimeout bool) (models.TaskInstance, error) {
	var emptyUUID gocql.UUID
	if !isFailedByTimeout {
		ti, err := models.TaskInstancePersistenceInstance.GetMinimalInstanceByID(ctx, taskInstanceID)
		if err != nil {
			return ti, err
		}
		if ti.ID == emptyUUID {
			return ti, fmt.Errorf("no TaskInstance found by TaskInstanceID %v", taskInstanceID)
		}
		return ti, nil
	}

	taskInstances, err := models.TaskInstancePersistenceInstance.GetByIDs(ctx, taskInstanceID)
	if err != nil {
		return models.TaskInstance{}, err
	}
	if len(taskInstances) == 0 {
		return models.TaskInstance{}, fmt.Errorf("no TaskInstance found by TaskInstanceID %v", taskInstanceID)
	}
	return taskInstances[0], nil
}

func handleResultWebhooks(
	ctx context.Context,
	partner string,
	taskData processData,
	executionResults []models.ExecutionResult,
) {
	var hookedTasks []taskResultWebhook.TaskResult
	for _, res := range executionResults {
		if len(taskData.ResultWebhook) == 0 {
			continue
		}

		resultMessage, err := models.GetResultMessage(ctx, partner, taskData.OriginID, res.ExecutionStatus)
		if err != nil {
			logger.Log.WarnfCtx(ctx, "handleResultWebhooks: error during getResultMessage for Task ID %v. Err: %v", taskData.TaskID, err)
		}

		hookedTasks = append(hookedTasks, taskResultWebhook.TaskResult{
			ID:            taskData.TaskID,
			ResultWebhook: taskData.ResultWebhook,
			StdErr:        res.StdErr,
			StdOut:        res.StdOut,
			Success:       res.ExecutionStatus == statuses.TaskInstanceSuccess,
			ResultMessage: resultMessage,
		})
	}

	if len(hookedTasks) > 0 {
		go taskResultWebhook.CallWebhooksFor(ctx, hookedTasks)
	}
}
