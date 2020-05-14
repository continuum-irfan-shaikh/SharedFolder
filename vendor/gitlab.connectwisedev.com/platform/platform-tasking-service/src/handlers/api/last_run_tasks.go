package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

const (
	lastRunTasksLogFormat = "LastRunTasks.Api: useCase returned error: %v"
	taskHistoryErr        = "GetTasksHistory.fetchAndValidateTimeFrame returned error - %v"
)

//LastRunTasksApi request handler
type LastRunTasksApi struct {
	uc  TasksInteractor
	log logger.Logger
}

//NewLastRunTasksApi returns new Tasks' history handler
func NewLastRunTasksApi(uc TasksInteractor, log logger.Logger) *LastRunTasksApi {
	return &LastRunTasksApi{uc: uc, log: log}
}

func (t *LastRunTasksApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		partnerID = mux.Vars(r)[partnerIDKey]
		ctx       = context.WithValue(r.Context(), config.PartnerIDKeyCTX, partnerID)
	)

	switch r.Method {
	case http.MethodGet:
		from, to, err := fetchAndValidateTimeFrame(r)
		if err != nil {
			t.historyTabHandleError(w, r, err)
			return
		}
		t.log.DebugfCtx(r.Context(), "LastRunTasksApi: get data from %v and to %v", from, to)

		data, err := t.uc.GetTasksHistory(ctx, from, to)
		if err != nil {
			t.historyTabHandleError(w, r, err)
			return
		}
		common.RenderJSON(w, data)
	default:
		common.SendStatusCodeWithMessage(w, r, http.StatusMethodNotAllowed, "unsupported method")
	}
}

// historyTabHandleError  handles useCase errors with corresponding http status code
func (t *LastRunTasksApi) historyTabHandleError(w http.ResponseWriter, r *http.Request, err error) {
	switch err.(type) {
	case errorcode.BadRequestErr:
		e := err.(errorcode.BadRequestErr)
		t.log.ErrfCtx(r.Context(), e.ErrorCode, lastRunTasksLogFormat, e.LogMessage)
		common.SendBadRequest(w, r, e.ErrorCode)
	case errorcode.InternalServerErr:
		e := err.(errorcode.InternalServerErr)
		t.log.ErrfCtx(r.Context(), e.ErrorCode, lastRunTasksLogFormat, e.LogMessage)
		common.SendInternalServerError(w, r, e.ErrorCode)
	default:
		t.log.ErrfCtx(r.Context(), errorcode.ErrorCantGetTasksSummaryData, lastRunTasksLogFormat, err)
		common.SendBadRequest(w, r, errorcode.ErrorCantGetTasksSummaryData)
	}
}

func fetchAndValidateTimeFrame(r *http.Request) (from, to time.Time, err error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err = time.Parse(time.RFC3339Nano, fromStr)
	if err != nil {
		return from, to, errorcode.NewBadRequestErr(errorcode.ErrorCantDecodeInputData,
			fmt.Errorf(taskHistoryErr, err).Error())
	}

	to, err = time.Parse(time.RFC3339Nano, toStr)
	if err != nil {
		return from, to, errorcode.NewBadRequestErr(errorcode.ErrorCantDecodeInputData,
			fmt.Errorf(taskHistoryErr, err).Error())
	}

	var (
		now                  = time.Now().In(to.Location())
		threeMonthsBeforeNow = now.Add(-3 * (time.Hour * 24 * 30)).Add(-1 * 24 * time.Hour)
	)

	if from.IsZero() ||
		to.IsZero() ||
		to.Before(from) ||
		to.After(now.Add(24*time.Hour)) ||
		from.Before(threeMonthsBeforeNow) ||
		from.After(now) {

		err = fmt.Errorf("invalid time frame: from [%v] to [%v]", from, to)
		return from, to, errorcode.NewBadRequestErr(errorcode.ErrorCantDecodeInputData,
			fmt.Errorf(taskHistoryErr, err).Error())
	}

	return
}
