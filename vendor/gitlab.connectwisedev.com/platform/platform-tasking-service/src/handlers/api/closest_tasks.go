package api

//go:generate mockgen -destination=../../mocks/mock-usecases/tasks_mock.go -package=mockusecases -source=./closest_tasks.go

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/types"
)

const (
	partnerIDKey = "partnerID"
)

//TasksInteractor - interface of Closest Tasks use case
type TasksInteractor interface {
	GetClosestTasks(ctx context.Context, endpoints entities.EndpointsInput) (tasks entities.EndpointsClosestTasks, err error)
	GetScheduledTasks(ctx context.Context) (tasks []entities.ScheduledTasks, err error)
	DeleteScheduledTasks(ctx context.Context, taskIDs entities.TaskIDs) error
	GetTasksHistory(ctx context.Context, from, to time.Time) ([]entities.ScheduledTasks, error)
}

//NewClosestTasks - returns the new instance of Closest Tasks request handler
func NewClosestTasks(uc TasksInteractor, log logger.Logger) *ClosestTasks {
	return &ClosestTasks{
		uc:  uc,
		log: log,
	}
}

//ClosestTasks - Closest Tasks request handler
type ClosestTasks struct {
	uc  TasksInteractor
	log logger.Logger
}

//ServeHTTP - handles closest tasks request
func (ct *ClosestTasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, err := ct.getURLParams(r)
	if err != nil {
		ct.log.ErrfCtx(r.Context(), errorcode.ErrorTimeFrameHasBadFormat, "ClosestTask: can't getURLparams %s", err)
		common.SendBadRequest(w, r, errorcode.ErrorTimeFrameHasBadFormat)
		return
	}

	endpoints, err := ct.parseRequest(r)
	if err != nil {
		ct.log.ErrfCtx(r.Context(), errorcode.ErrorCantDecodeInputData, "ClosestTask: can't parseRequest %+v", err)
		common.SendBadRequest(w, r, err.Error())
		return
	}

	tasks, err := ct.uc.GetClosestTasks(ctx, endpoints)
	if err != nil {
		ct.log.ErrfCtx(r.Context(), errorcode.ErrorCantGetClosestTasks, "GetClosestTasks: %+v", err)

		switch err.(type) {
		case *types.MultiError:
		default:
			common.SendInternalServerError(w, r, errorcode.ErrorCantGetClosestTasks)
			return
		}
	}
	common.RenderJSON(w, tasks)
}

func (ct *ClosestTasks) getURLParams(r *http.Request) (ctx context.Context, err error) {
	ctx = context.Background()

	params := mux.Vars(r)
	partnerID := params[partnerIDKey]
	ctx = context.WithValue(ctx, config.PartnerIDKeyCTX, partnerID)

	ct.log.DebugfCtx(r.Context(), "closest tasks request have been received with partnerID:%s", partnerID)

	return ctx, nil
}

func (ct *ClosestTasks) parseRequest(r *http.Request) (endpoints entities.EndpointsInput, err error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return endpoints, errors.Wrap(err, errorcode.ErrorCantDecodeInputData)
	}

	err = json.Unmarshal(data, &endpoints)
	if err != nil {
		return endpoints, errors.Wrap(err, errorcode.ErrorCantDecodeInputData)
	}
	return
}
