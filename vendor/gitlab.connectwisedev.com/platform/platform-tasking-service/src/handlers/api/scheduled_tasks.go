package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const schedulerTasksLogFormat = "ScheduledTasks.Api: usecase returned error: %v"

// ScheduledTasks request handler
type ScheduledTasks struct {
	uc  TasksInteractor
	log logger.Logger
}

// NewScheduledTasksApi returns new scheduled Tasks handler
func NewScheduledTasksApi(uc TasksInteractor, log logger.Logger) *ScheduledTasks {
	return &ScheduledTasks{uc: uc, log: log}
}

// ServeHTTP handles scheduled tasks requests
func (st *ScheduledTasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		partnerID = mux.Vars(r)[partnerIDKey]
		ctx       = context.WithValue(r.Context(), config.PartnerIDKeyCTX, partnerID)
	)

	switch r.Method {
	case http.MethodGet:
		data, err := st.uc.GetScheduledTasks(ctx)
		if err != nil {
			st.handleError(w, r, err)
			return
		}
		common.RenderJSON(w, data)

	case http.MethodDelete:
		tasks, err := st.parseRequest(r)
		if err != nil {
			st.log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "ScheduledTasks.Api DELETE: can't parse request %+v", err)
			common.SendBadRequest(w, r, err.Error())
			return
		}

		err = st.uc.DeleteScheduledTasks(ctx, tasks)
		if err != nil {
			st.log.ErrfCtx(r.Context(), errorcode.ErrorUsecaseProcessing, "ScheduledTasks.Api DELETE: usecase returned error %+v", err)
			common.SendInternalServerError(w, r, err.Error())
			return
		}
		common.SendNoContent(w)
	default:
		common.SendStatusCodeWithMessage(w, r, http.StatusMethodNotAllowed, "unsupported method")
	}
}

func (st *ScheduledTasks) parseRequest(r *http.Request) (tasks entities.TaskIDs, err error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return tasks, errors.Wrap(err, errorcode.ErrorCantDecodeInputData)
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return tasks, errors.Wrap(err, errorcode.ErrorCantDecodeInputData)
	}

	st.log.DebugfCtx(r.Context(), "ScheduledTasks.parseRequest: payload %v", string(data))
	return
}

// handleError  handles usecase errors with corresponding http status code
func (st *ScheduledTasks) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch err.(type) {
	case errorcode.BadRequestErr:
		e := err.(errorcode.BadRequestErr)
		st.log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, schedulerTasksLogFormat, e.LogMessage)
		common.SendBadRequest(w, r, e.ErrorCode)
	case errorcode.InternalServerErr:
		e := err.(errorcode.InternalServerErr)
		st.log.ErrfCtx(r.Context(), e.ErrorCode, schedulerTasksLogFormat, e.LogMessage)
		common.SendInternalServerError(w, r, e.ErrorCode)
	default:
		st.log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTasksSummaryData, schedulerTasksLogFormat, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantGetTasksSummaryData)
	}
}
