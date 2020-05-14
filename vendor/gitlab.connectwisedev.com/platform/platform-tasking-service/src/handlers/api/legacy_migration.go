package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/legacy"
)

const (
	legacyMigrationLogFormat = "LegacyMigration.Api: usecase returned error: %v"
	scriptIDKey              = "scriptID"
	jobIDKey                 = "jobID"
)

// LegacyMigration is a legacy api handler
type LegacyMigration struct {
	uc  legacy.MigrationUC
	log logger.Logger
}

// NewLegacyMigration returns new LegacyMigration api handler struct
func NewLegacyMigration(uc legacy.MigrationUC, log logger.Logger) *LegacyMigration {
	return &LegacyMigration{uc: uc, log: log}
}

// GetByPartner returns migration info by partnerID
func (l *LegacyMigration) GetByPartner(w http.ResponseWriter, r *http.Request) {
	data, err := l.uc.GetByPartner(mux.Vars(r)[partnerIDKey])
	if err != nil {
		l.handleError(w, r, err)
		return
	}
	common.RenderJSON(w, data)
}

// GetByLegacyScript returns migration info by scriptID
func (l *LegacyMigration) GetByLegacyScript(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)[partnerIDKey]
	scriptID := mux.Vars(r)[scriptIDKey]

	data, err := l.uc.GetByScriptID(partnerID, scriptID)
	if err != nil {
		l.handleError(w, r, err)
		return
	}
	common.RenderJSON(w, data)
}

// InsertLegacyInfo returns migration info by scriptID
func (l *LegacyMigration) InsertLegacyInfo(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)[partnerIDKey]
	info, err := l.extractScriptRequest(r)
	if err != nil {
		l.handleError(w, r, errorcode.NewBadRequestErr(errorcode.ErrorCantDecodeInputData, err.Error()))
		return
	}

	if err = l.uc.InsertScriptInfo(partnerID, info); err != nil {
		l.handleError(w, r, err)
		return
	}
	common.SendCreated(w, r, errorcode.CodeCreated)
}

// GetJobInfoByPartner returns migration info by partnerID
func (l *LegacyMigration) GetJobInfoByPartner(w http.ResponseWriter, r *http.Request) {
	data, err := l.uc.GetJobsByPartner(mux.Vars(r)[partnerIDKey])
	if err != nil {
		l.handleError(w, r, err)
		return
	}
	common.RenderJSON(w, data)
}

// GetByLegacyJob returns migration info by scriptID
func (l *LegacyMigration) GetByLegacyJob(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)[partnerIDKey]
	jobID := mux.Vars(r)[jobIDKey]

	data, err := l.uc.GetByJobID(partnerID, jobID)
	if err != nil {
		l.handleError(w, r, err)
		return
	}
	common.RenderJSON(w, data)
}

// InsertJobInfo returns migration info by scriptID
func (l *LegacyMigration) InsertJobInfo(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)[partnerIDKey]
	info, err := l.extractJobRequest(r)
	if err != nil {
		l.handleError(w, r, errorcode.NewBadRequestErr(errorcode.ErrorCantDecodeInputData, err.Error()))
		return
	}

	if err = l.uc.InsertJobInfo(partnerID, info); err != nil {
		l.handleError(w, r, err)
		return
	}
	common.SendCreated(w, r, errorcode.CodeCreated)
}

func (l *LegacyMigration) extractScriptRequest(r *http.Request) (info models.LegacyScriptInfo, err error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return info, err
	}

	if err = json.Unmarshal(data, &info); err != nil {
		return info, err
	}

	if info.LegacyScriptID == "" {
		return info, errors.New("script id cannot be empty")
	}
	return
}

func (l *LegacyMigration) extractJobRequest(r *http.Request) (info models.LegacyJobInfo, err error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return info, err
	}

	if err = json.Unmarshal(data, &info); err != nil {
		return info, err
	}

	if info.OriginID == "" || info.LegacyJobID == "" || info.LegacyScriptID == "" || info.LegacyTemplateID == "" {
		return info, fmt.Errorf("fields cannot be empty for %v", info)
	}
	return
}

// handleError  handles usecase errors with corresponding http status code
func (l *LegacyMigration) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch err.(type) {
	case errorcode.NotFoundErr:
		e := err.(errorcode.NotFoundErr)
		l.log.ErrfCtx(r.Context(), e.ErrorCode, legacyMigrationLogFormat, e.LogMessage)
		common.SendNotFound(w, r, e.ErrorCode)
	case errorcode.InternalServerErr:
		e := err.(errorcode.InternalServerErr)
		l.log.ErrfCtx(r.Context(), e.ErrorCode, legacyMigrationLogFormat, e.LogMessage)
		common.SendInternalServerError(w, r, e.ErrorCode)
	default:
		l.log.ErrfCtx(r.Context(), errorcode.ErrorWrongTemplateType, legacyMigrationLogFormat, err)
		common.SendBadRequest(w, r, errorcode.ErrorWrongTemplateType)
	}
}
