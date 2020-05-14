package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	legacymigration "gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/legacy"
)

const (
	getByPartnerURL     = "/{partnerID}/scripts"
	getByScriptURL      = "/{partnerID}/scripts/{scriptID}"
	insertScriptInfoURL = "/{partnerID}/scripts"

	getByJobsPartnerURL = "/{partnerID}/job"
	getByJobURL         = "/{partnerID}/job/{jobID}"
	insertJobInfoURL    = "/{partnerID}/job"

	legacyPartner  = "1"
	legacyScriptID = "2"
	legacyJobID    = "2"
)

func init()  {
	logger.Load(config.Config.Log)
}

func TestGetByPartner(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		URL          string
		mock         func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC)
		expectedCode int
	}{
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetByPartner(legacyPartner).Return([]models.LegacyScriptInfo{}, errors.New("err"))
				return nil, legacy
			},
			expectedCode: http.StatusBadRequest,
			URL:          "/" + legacyPartner + "/scripts",
		},
		{
			name: "testCase 2 - ok",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetByPartner(legacyPartner)
				return nil, legacy
			},
			expectedCode: http.StatusOK,
			URL:          "/" + legacyPartner + "/scripts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, legacy := tc.mock(mc)
			s := *NewLegacyMigration(legacy, logger.Log)
			router := getLegacyMigrationRouter(s)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestLegacyMigration_GetJobInfoByPartner(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		URL          string
		mock         func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC)
		expectedCode int
	}{
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetJobsByPartner(legacyPartner).Return([]models.LegacyJobInfo{}, errors.New("err"))
				return nil, legacy
			},
			expectedCode: http.StatusBadRequest,
			URL:          "/" + legacyPartner + "/job",
		},
		{
			name: "testCase 2 - ok",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetJobsByPartner(legacyPartner)
				return nil, legacy
			},
			expectedCode: http.StatusOK,
			URL:          "/" + legacyPartner + "/job",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, legacy := tc.mock(mc)
			s := *NewLegacyMigration(legacy, logger.Log)
			router := getLegacyMigrationRouter(s)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestGetByScript(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		URL          string
		mock         func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC)
		expectedCode int
	}{
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {

				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetByScriptID(legacyPartner, legacyScriptID).
					Return(models.LegacyScriptInfo{}, errors.New("err"))
				return nil, legacy
			},
			expectedCode: http.StatusBadRequest,
			URL:          "/" + legacyPartner + "/scripts/" + legacyScriptID,
		},
		{
			name: "testCase 2 - ok",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetByScriptID(legacyPartner, legacyScriptID)
				return nil, legacy
			},
			expectedCode: http.StatusOK,
			URL:          "/" + legacyPartner + "/scripts/" + legacyScriptID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, legacy := tc.mock(mc)
			s := *NewLegacyMigration(legacy, logger.Log)
			router := getLegacyMigrationRouter(s)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestLegacyMigration_GetByLegacyJob(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		URL          string
		mock         func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC)
		expectedCode int
	}{
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetByJobID(legacyPartner, legacyJobID).
					Return(models.LegacyJobInfo{}, errors.New("err"))
				return nil, legacy
			},
			expectedCode: http.StatusBadRequest,
			URL:          "/" + legacyPartner + "/job/" + legacyJobID,
		},
		{
			name: "testCase 2 - ok",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().GetByJobID(legacyPartner, legacyJobID)
				return nil, legacy
			},
			expectedCode: http.StatusOK,
			URL:          "/" + legacyPartner + "/job/" + legacyJobID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, legacy := tc.mock(mc)
			s := *NewLegacyMigration(legacy, logger.Log)
			router := getLegacyMigrationRouter(s)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestLegacyMigration_InsertLegacyInfo(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}

	validBody := models.LegacyScriptInfo{
		LegacyScriptID: legacyScriptID,
		OriginID:       legacyScriptID,
		DefinitionID:   legacyScriptID,
		IsSequence:     false,
		DefinitionDetails: models.TaskDefinitionDetails{
			UserParameters: `{"hey":"it's me"}`,
		},
	}

	testCases := []struct {
		name         string
		URL          string
		mock         func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC)
		body         interface{}
		expectedCode int
	}{

		{
			name: "testCase 0 - can't parse - wrong body",
			URL:  "/" + legacyPartner + "/scripts",
			body: models.TaskInstance{
				TriggeredBy: "this is not struct we're supposed to get",
			},
			mock: func(ctl *gomock.Controller) (i logger.Logger, uc legacymigration.MigrationUC) {
				return nil, nil
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().InsertScriptInfo(legacyPartner, validBody).
					Return(errorcode.NewInternalServerErr("e", "rr"))
				return nil, legacy
			},
			body:         validBody,
			expectedCode: http.StatusInternalServerError,
			URL:          "/" + legacyPartner + "/scripts",
		},
		{
			name: "testCase 2 - uc returned not found err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().InsertScriptInfo(legacyPartner, validBody).
					Return(errorcode.NewNotFoundErr("e", "rr"))
				return nil, legacy
			},
			body:         validBody,
			expectedCode: http.StatusNotFound,
			URL:          "/" + legacyPartner + "/scripts",
		},
		{
			name: "testCase 3 - ok",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().InsertScriptInfo(legacyPartner, validBody)
				return nil, legacy
			},
			body:         validBody,
			expectedCode: http.StatusCreated,
			URL:          "/" + legacyPartner + "/scripts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, legacy := tc.mock(mc)
			s := *NewLegacyMigration(legacy, logger.Log)
			router := getLegacyMigrationRouter(s)

			b, err := json.Marshal(&tc.body)
			if err != nil {
				t.Error("can't marshal object")
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, tc.URL, bytes.NewBuffer(b))
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestLegacyMigration_InsertJobInfo(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}

	validBody := models.LegacyJobInfo{
		LegacyJobID:      legacyScriptID,
		LegacyScriptID:   legacyScriptID,
		LegacyTemplateID: legacyScriptID,
		Type:             "script",
		DefinitionID:     legacyScriptID,
		OriginID:         legacyScriptID,
		TaskID:           "id",
	}

	testCases := []struct {
		name         string
		URL          string
		mock         func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC)
		body         interface{}
		expectedCode int
	}{

		{
			name: "testCase 0 - can't parse - wrong body",
			URL:  "/" + legacyPartner + "/job",
			body: models.TaskInstance{
				TriggeredBy: "this is not struct we're supposed to get",
			},
			mock: func(ctl *gomock.Controller) (i logger.Logger, uc legacymigration.MigrationUC) {
				return nil, nil
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().InsertJobInfo(legacyPartner, validBody).
					Return(errorcode.NewInternalServerErr("e", "rr"))
				return nil, legacy
			},
			body:         validBody,
			expectedCode: http.StatusInternalServerError,
			URL:          "/" + legacyPartner + "/job",
		},
		{
			name: "testCase 2 - uc returned not found err",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {

				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().InsertJobInfo(legacyPartner, validBody).
					Return(errorcode.NewNotFoundErr("e", "rr"))
				return nil, legacy
			},
			body:         validBody,
			expectedCode: http.StatusNotFound,
			URL:          "/" + legacyPartner + "/job",
		},
		{
			name: "testCase 3 - ok",
			mock: func(ctl *gomock.Controller) (logger.Logger, legacymigration.MigrationUC) {
				legacy := mocks.NewMockMigrationUC(ctl)
				legacy.EXPECT().InsertJobInfo(legacyPartner, validBody)
				return nil, legacy
			},
			body:         validBody,
			expectedCode: http.StatusCreated,
			URL:          "/" + legacyPartner + "/job",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, legacy := tc.mock(mc)
			s := *NewLegacyMigration(legacy, logger.Log)
			router := getLegacyMigrationRouter(s)

			b, err := json.Marshal(&tc.body)
			if err != nil {
				t.Error("can't marshal object")
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, tc.URL, bytes.NewBuffer(b))
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func getLegacyMigrationRouter(service LegacyMigration) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(insertScriptInfoURL, service.InsertLegacyInfo).Methods(http.MethodPost)
	router.HandleFunc(getByPartnerURL, service.GetByPartner).Methods(http.MethodGet)
	router.HandleFunc(getByScriptURL, service.GetByLegacyScript).Methods(http.MethodGet)

	router.HandleFunc(insertJobInfoURL, service.InsertJobInfo).Methods(http.MethodPost)
	router.HandleFunc(getByJobsPartnerURL, service.GetJobInfoByPartner).Methods(http.MethodGet)
	router.HandleFunc(getByJobURL, service.GetByLegacyJob).Methods(http.MethodGet)
	return
}
