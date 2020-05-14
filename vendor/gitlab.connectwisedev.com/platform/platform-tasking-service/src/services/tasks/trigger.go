package tasks

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/validator"
)

const (
	jsonPostfix                   = ".json"
	executeTriggerLogPrefix       = "TaskService.ExecuteTrigger: "
	getTriggerLogPrefix           = "TaskService.GetTriggersList: "
	canNotDecodeInputDataErrMsg   = "can not decode input data : %s"
	canNotGetActiveTriggersErrMsg = "can not get active triggers by type %s and partnerID %s: %s"
	saveAndProcessTaskErrMsg      = "error while saving and processing internal task %v: . err : %v"
	postExecutionErrMsg           = "error while making postExecution %v: . err : %v"
	closingFileError              = "error while closing the file: %v"
	notApplicable                 = "trigger is not applicable"
	triggersPackFileName          = "triggers-pack-file"
	contextError                  = "can't get parameter %s from context"
)

// ExecuteTrigger is an endpoint that Automation engine MS uses as a hook for policies events
func (t *TaskService) ExecuteTrigger(w http.ResponseWriter, r *http.Request) {
	var payload apiModels.TriggerExecutionPayload
	if err := validator.ExtractStructFromRequest(r, &payload); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, executeTriggerLogPrefix+canNotDecodeInputDataErrMsg, err.Error())
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	ctx, err := t.buildExecuteTriggerCtx(r)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, executeTriggerLogPrefix, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantDecodeInputData)
		return
	}

	if err := t.processExecuteTriggerRequest(ctx, payload); err != nil {
		switch err.(type) {
		case errorcode.InternalServerErr:
			connErr := err.(errorcode.InternalServerErr)
			logger.Log.ErrfCtx(r.Context(), connErr.ErrorCode, connErr.LogMessage)
			common.SendInternalServerError(w, r, connErr.ErrorCode)
			return
		default:
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorInternalServerError, err.Error())
			common.SendInternalServerError(w, r, errorcode.ErrorInternalServerError)
			return
		}
	}
	common.SendCreated(w, r, errorcode.CodeUpdated)
}

func (t *TaskService) processExecuteTriggerRequest(ctx context.Context, payload apiModels.TriggerExecutionPayload) error {
	// RMM-47378. Alerting sends alert types for deleting the alert and for creating. if it's deleting we don't execute tasks
	if payload.AlertType == apiModels.AlertTypeDelete {
		return nil
	}

	triggerType, partnerID, endpointID, err := t.extractCtx(ctx)
	if err != nil {
		return err
	}

	activeTriggers, err := t.trigger.GetActiveTriggers(ctx)
	if err != nil {
		return errorcode.InternalServerErr{
			ErrorCode:  errorcode.ErrorCantGetActiveTriggers,
			LogMessage: fmt.Sprintf(executeTriggerLogPrefix+canNotGetActiveTriggersErrMsg, triggerType, partnerID, err),
		}
	}

	for _, at := range activeTriggers {
		if err = t.processActiveTriggers(at, ctx, payload, endpointID, triggerType); err != nil {
			logger.Log.ErrfCtx(ctx, errorcode.ErrorCantExecuteTrigger, fmt.Sprintf("error during processing active trigger %v", err))
		}
	}

	return nil
}

func (t *TaskService) processActiveTriggers(at e.ActiveTrigger, ctx context.Context, payload apiModels.TriggerExecutionPayload, endpointID gocql.UUID, triggerType string) (err error) {
	if !t.triggerHasValidTimeFrame(at) {
		return
	}

	task, err := t.trigger.GetTask(ctx, at.TaskID)
	if err != nil {
		return
	}

	if !t.trigger.IsApplicable(ctx, task, payload) {
		return
	}

	if err = t.trigger.PreExecution(ctx.Value(config.TriggerTypeIDKeyCTX).(string), task); err != nil {
		return
	}

	// if task trigger is disabled
	if task.State == statuses.TaskStateDisabled || task.State == statuses.TaskStateInactive {
		return
	}

	task.Schedule.Regularity = apiModels.RunNow // we won't save this task so we mark it as RunNow to send on execution
	task.ManagedEndpointID = endpointID
	task.Schedule.TriggerTypes = []string{triggerType}
	task.Schedule.TriggerFrames = []apiModels.TriggerFrame{{TriggerType: triggerType}}
	task.Schedule.EndRunTime = time.Time{} // to prevent RunNow Canceled till

	if errMsg, err := t.saveAndProcessTask(ctx, &task, []models.Task{task}, false, false); err != nil {
		return errorcode.InternalServerErr{
			ErrorCode:  errMsg,
			LogMessage: fmt.Sprintf(executeTriggerLogPrefix+saveAndProcessTaskErrMsg, task, err),
		}
	}

	if err = t.trigger.PostExecution(triggerType, task); err != nil {
		return
	}
	return
}

func (t *TaskService) extractCtx(ctx context.Context) (string, string, gocql.UUID, error) {
	triggerType, ok := ctx.Value(config.TriggerTypeIDKeyCTX).(string)
	if !ok {
		return "", "", gocql.UUID{}, errors.Errorf(contextError, config.TriggerTypeIDKeyCTX)
	}

	partnerID, ok := ctx.Value(config.PartnerIDKeyCTX).(string)
	if !ok {
		return "", "", gocql.UUID{}, errors.Errorf(contextError, config.PartnerIDKeyCTX)
	}

	endpointID, ok := ctx.Value(config.EndpointIDKeyCTX).(gocql.UUID)
	if !ok {
		return "", "", gocql.UUID{}, errors.Errorf(contextError, config.EndpointIDKeyCTX)
	}
	return triggerType, partnerID, endpointID, nil
}

func (t *TaskService) triggerHasValidTimeFrame(trigger e.ActiveTrigger) bool {
	if trigger.StartTimeFrame.IsZero() {
		return true
	}

	startTime := t.resetTimeToCurrentDay(trigger.StartTimeFrame)
	endTime := t.resetTimeToCurrentDay(trigger.EndTimeFrame)
	currentTime := time.Now()
	return currentTime.After(startTime) && currentTime.Before(endTime)
}

func (t *TaskService) resetTimeToCurrentDay(day time.Time) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), day.Hour(), day.Minute(), day.Second(), day.Nanosecond(), day.Location())
}

func (t *TaskService) buildExecuteTriggerCtx(r *http.Request) (context.Context, error) {
	var (
		ctx         = r.Context()
		params      = mux.Vars(r)
		partnerID   = params["partnerID"]
		siteID      = params["siteID"]
		clientID    = params["clientID"]
		endpoint    = params["endpointID"]
		triggerType = params["triggerType"]
		emptyUUID   gocql.UUID
	)

	endpointID, err := gocql.ParseUUID(endpoint)
	if err != nil || endpointID == emptyUUID {
		return ctx, fmt.Errorf("bad endpointID %v, err:%v", endpoint, err)
	}

	ctx = context.WithValue(ctx, config.PartnerIDKeyCTX, partnerID)
	ctx = context.WithValue(ctx, config.SiteIDKeyCTX, siteID)
	ctx = context.WithValue(ctx, config.ClientIDKeyCTX, clientID)
	ctx = context.WithValue(ctx, config.EndpointIDKeyCTX, endpointID)
	ctx = context.WithValue(ctx, config.TriggerTypeIDKeyCTX, triggerType)
	return ctx, nil
}

// GetTriggersList returns all trigger types stored in DB
func (t *TaskService) GetTriggersList(w http.ResponseWriter, r *http.Request) {
	types, err := t.triggerDefinition.GetTriggerTypes()
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTaskByTaskID, getTriggerLogPrefix+err.Error())
		common.SendInternalServerError(w, r, errorcode.ErrorCantGetTaskByTaskID)
		return
	}
	common.RenderJSON(w, types)
}

// UploadTriggerDefinitions returns all trigger types stored in DB
func (t *TaskService) UploadTriggerDefinitions(w http.ResponseWriter, r *http.Request) {
	triggerDefs, err := t.fetchTriggers(r)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "cannot fetch trigger definitions from the request: %v", err)
		common.SendBadRequest(w, r, fmt.Sprint("cannot fetch trigger definitions from the request: ", err))
		return
	}

	if len(triggerDefs) < 1 {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "trigger definitions length cannot be less than 1")
		common.SendBadRequest(w, r, "trigger definitions length cannot be less than 1")
		return
	}

	if err = t.triggerDefinition.ImportExternalTriggers(triggerDefs); err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantProcessTriggerDefinitions, "UploadTriggerDefinitions: insert error %v", err)
		common.SendInternalServerError(w, r, errorcode.ErrorCantProcessTriggerDefinitions)
		return
	}
}

func (t *TaskService) fetchTriggers(r *http.Request) (triggerDefs []e.TriggerDefinition, err error) {
	file, _, err := r.FormFile(triggersPackFileName)
	if err != nil {
		return
	}
	zipRawData, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	defer func() {
		if err = r.Body.Close(); err != nil {
			logger.Log.WarnfCtx(r.Context(), closingFileError, err)
		}
	}()

	zipReader, err := zip.NewReader(bytes.NewReader(zipRawData), int64(len(zipRawData)))
	if err != nil {
		return
	}

	for _, file := range zipReader.File {
		if !strings.Contains(file.Name, jsonPostfix) {
			// ignore all non json files
			continue
		}

		triggerDef, err := t.extractSingleTrigger(r.Context(), file)
		if err != nil {
			return triggerDefs, err
		}

		triggerDefs = append(triggerDefs, triggerDef)
	}
	return
}

func (t *TaskService) extractSingleTrigger(ctx context.Context, file *zip.File) (def e.TriggerDefinition, err error) {
	currentFileRC, err := file.Open()
	if err != nil {
		return
	}
	defer func() {
		if err = currentFileRC.Close(); err != nil {
			logger.Log.WarnfCtx(ctx, closingFileError, err)
		}
	}()

	var rawScriptDefinitionData []byte
	rawScriptDefinitionData, err = ioutil.ReadAll(currentFileRC)
	if err != nil {
		return
	}

	if err = json.Unmarshal(rawScriptDefinitionData, &def); err != nil {
		return def, errors.Wrap(err, file.Name)
	}
	return
}
