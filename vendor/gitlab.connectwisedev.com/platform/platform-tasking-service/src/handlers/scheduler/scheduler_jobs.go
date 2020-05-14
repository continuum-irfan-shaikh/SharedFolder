package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gocql/gocql"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

//go:generate mockgen -destination=./scheduler_deprecated_uc_mock_test.go -package=scheduler -source=./scheduler_deprecated.go

// ExecutionResultUpdateUC - interface for execution results update service
type ExecutionResultUpdateUC interface {
	ProcessExecutionResults(ctx context.Context, partnerID string, taskInstanceID gocql.UUID, results ...apiModels.ExecutionResult) error
}

//SchedulerRepo privide interface for get/update last expired execution time check
type SchedulerRepo interface {
	GetLastExpiredExecutionCheck() (lastUpdate time.Time, err error)
	UpdateLastExpiredExecutionCheck(time time.Time) (err error)
}

//ExecutionResultRepo - represents execution result repository
type ExecutionResultRepo interface {
	Publish(msg apiModels.ExecutionResultKafkaMessage) error
}

//ExecutionResultRepo - represents execution result repository
type CounterService interface {
	GetCountersByPartner(w http.ResponseWriter, r *http.Request)
	GetCountersByPartnerAndEndpoint(w http.ResponseWriter, r *http.Request)
	RecalculateAllCounters(w http.ResponseWriter, r *http.Request)
}

// Service is a scheduler Service struct
type Service struct {
	taskCounterService      CounterService
	executionResultsRepo    ExecutionResultRepo
	execResultUpdateService ExecutionResultUpdateUC
	schedulerRepo           SchedulerRepo
}

// New is a function to return New scheduler Service
func New(
	taskCounterService CounterService,
	e ExecutionResultRepo,
	executionResultUpdate ExecutionResultUpdateUC,
	schedulerRepo SchedulerRepo,
) Service {
	return Service{
		taskCounterService:      taskCounterService,
		executionResultsRepo:    e,
		execResultUpdateService: executionResultUpdate,
		schedulerRepo:           schedulerRepo,
	}
}

// CheckForRetainedData ..
func (s Service) CheckForRetainedData(ctx context.Context) {
	models.DeleteOldTaskInstanceCounts(ctx)
}

// CheckForExpiredExecutions checks and updates tasks by expired executions
func (s Service) CheckForExpiredExecutions(ctx context.Context) {
	currentTime := time.Now().UTC().Truncate(time.Minute)
	ctx = transactionID.NewContext()

	lastUpdate, err := s.schedulerRepo.GetLastExpiredExecutionCheck()
	if err != nil {
		if err != gocql.ErrNotFound {
			logger.Log.ErrfCtx(ctx, "scheduler.checkForExpiredExecutions: error while getting last time where expired execution have been checked. err: %v", err.Error())
			return
		}

		lastUpdate = currentTime.Add(-time.Minute)
	}

	for lastUpdate.Before(currentTime) {
		lastUpdate = lastUpdate.Add(time.Minute)
		if err = s.processExpiredExecutionByUpdateTime(ctx, lastUpdate); err != nil {
			return
		}
	}
}

func (s Service) processExpiredExecutionByUpdateTime(ctx context.Context, lastUpdate time.Time) error {
	expiredTasks, err := models.ExecutionExpirationPersistenceInstance.GetByExpirationTime(ctx, lastUpdate)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantGetTaskExecutionResults, "scheduler.checkForExpiredExecutions: error while getting execution expirations by expiration time. err: %v", err)
		return err
	}

	for _, exp := range expiredTasks {
		ee := apiModels.ExpiredExecution{
			TaskInstanceID:     exp.TaskInstanceID,
			ManagedEndpointIDs: exp.ManagedEndpointIDs,
		}

		s.updateExecutionResultsByExpiredExecutions(ctx, exp.PartnerID, ee)

		response, err := sendExpiredExecutionsToScripting(ctx, exp.PartnerID, ee)
		if err != nil {
			logger.Log.ErrfCtx(ctx, errorcode.ErrorCantPerformRequest, "scheduler.checkForExpiredExecutions: error while making request to ScriptingMS. err: %v", err)
		}

		if err = common.CloseRespBody(response); err != nil {
			logger.Log.WarnfCtx(ctx, "scheduler.checkForExpiredExecutions: error while closing body. err: %v", err)
		}
	}

	if err = s.schedulerRepo.UpdateLastExpiredExecutionCheck(lastUpdate); err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "scheduler.checkForExpiredExecutions: error while updating last expired execution check time. err: %v", err)
	}
	return nil
}

func (s Service) updateExecutionResultsByExpiredExecutions(ctx context.Context, partnerID string, expiredExecution apiModels.ExpiredExecution) {
	executionResults := make([]apiModels.ExecutionResult, 0)
	for _, managedEndpointID := range expiredExecution.ManagedEndpointIDs {
		completionStatus, err := statuses.TaskInstanceStatusText(statuses.TaskInstanceFailed)
		if err != nil {
			logger.Log.WarnfCtx(ctx, "scheduler.updateExecutionResultsByExpiredExecutions: error while converting TaskInstanceFailed status to text, err: %v", err)
			continue
		}

		executionResults = append(executionResults, apiModels.ExecutionResult{
			CompletionStatus: completionStatus,
			EndpointID:       managedEndpointID.String(),
			UpdateTime:       time.Now().UTC(),
			ErrorDetails:     "Failed by timeout",
		})
	}

	if err := s.execResultUpdateService.ProcessExecutionResults(ctx, partnerID, expiredExecution.TaskInstanceID, executionResults...); err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessTaskExecutionResults, "scheduler: processing expired executions by TaskInstance ID [%v]. Err: %v", expiredExecution.TaskInstanceID, err)
	}
}

func sendExpiredExecutionsToScripting(ctx context.Context, partnerID string, ee apiModels.ExpiredExecution) (response *http.Response, err error) {
	scriptingURL := fmt.Sprintf("%s/partners/%s/mailbox-messages", config.Config.ScriptingMsURL, partnerID)
	body, err := json.Marshal([]apiModels.ExpiredExecution{ee})
	if err != nil {
		return
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
			MaxIdleConns:        2 * config.Config.HTTPClientMaxIdleConnPerHost,
			MaxIdleConnsPerHost: config.Config.HTTPClientMaxIdleConnPerHost,
			DisableKeepAlives:   false,
		},
	}

	return common.HTTPRequestWithRetry(ctx, httpClient, http.MethodPost, scriptingURL, body)
}

const recalculateTasksLog = "recalculateTasks: err: %v"

// RecalculateTasks is a scheduler Job to recalculate tasks
func (s Service) RecalculateTasks(ctx context.Context) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		logger.Log.WarnfCtx(ctx, recalculateTasksLog, err)
	}

	ctx = transactionID.NewContext()
	s.taskCounterService.RecalculateAllCounters(w, r.WithContext(ctx))

	resp := w.Result()
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		logger.Log.WarnfCtx(ctx, recalculateTasksLog, err)
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Log.WarnfCtx(ctx, recalculateTasksLog, err)
	}
}
