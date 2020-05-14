package templates

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	appLoader "gitlab.connectwisedev.com/platform/platform-tasking-service/src/app-loader"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	partnerID        = "1"
	originID         = "10000000-0000-0000-0000-000000000001"
	userHasNOCAccess = true
)

func init() {
	appLoader.LoadApplicationServices(true)
}

type TemplateTestCase struct {
	name                 string
	templateCacheMock    mock.TemplateCacheConf
	method               string
	expectedErrorMessage string
	expectedCode         int
	URL                  string
	userMock             user.Service
	expectedBody         []models.Template
}

func TestTemplateService_GetAll(t *testing.T) {

	testCases := []TemplateTestCase{
		{
			name: "testCase 1 - Could not get templates, server err",
			URL:  "/" + partnerID + "/task-definition-templates",
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetAllTemplates(gomock.Any(), partnerID, userHasNOCAccess).
					Return(nil, errors.New("err"))
				return tc
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTemplatesForExecutionMS,
		},
		{
			name: "testCase 2 - Successes",
			URL:  "/" + partnerID + "/task-definition-templates",
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetAllTemplates(gomock.Any(), partnerID, userHasNOCAccess).
					Return([]models.Template{{PartnerID: partnerID, Type: "script"}}, nil)
				return tc
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: []models.Template{
				{PartnerID: partnerID, Type: "script"},
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.templateCacheMock, tc.userMock)
			router := getTemplateRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)

			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
			// error message checking
			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
			// expected body checking
			if len(tc.expectedBody) > 0 {
				expectedBody, err := json.Marshal(tc.expectedBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if !bytes.Equal(expectedBody, w.Body.Bytes()) {
					t.Errorf("Wanted %v but got %v", string(expectedBody), w.Body.String())
				}
			}
		})
	}
}

func TestTemplateService_GetByType(t *testing.T) {

	testCases := []TemplateTestCase{
		{
			name:                 "testCase 1 - Cuold not parse type",
			URL:                  "/" + partnerID + "/task-definition-templates/badTypeToSend",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorWrongTemplateType,
		},
		{
			name: "testCase 2 - Could not get template by type,server error",
			URL:  "/" + partnerID + "/task-definition-templates/script",
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByType(gomock.Any(), partnerID, "script", userHasNOCAccess).
					Return(nil, errors.New("error"))
				return tc
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTemplatesForExecutionMS,
		},
		{
			name: "testCase 3 - Successes",
			URL:  "/" + partnerID + "/task-definition-templates/script",
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByType(gomock.Any(), partnerID, "script", userHasNOCAccess).
					Return([]models.Template{{PartnerID: partnerID, Type: "script"}}, nil)
				return tc
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: []models.Template{
				{PartnerID: partnerID, Type: "script"},
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.templateCacheMock, tc.userMock)
			router := getTemplateRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
			// error message checking
			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
			// expected body checking
			if len(tc.expectedBody) > 0 {
				expectedBody, err := json.Marshal(tc.expectedBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if !bytes.Equal(expectedBody, w.Body.Bytes()) {
					t.Errorf("Wanted %v but got %v", string(expectedBody), w.Body.String())
				}
			}
		})
	}
}

func TestTemplateService_GetByOriginID(t *testing.T) {

	type TemplateTestCase struct {
		name                 string
		templateCacheMock    mock.TemplateCacheConf
		method               string
		expectedErrorMessage string
		expectedCode         int
		URL                  string
		userMock             user.Service
		expectedBody         models.TemplateDetails
	}

	testCases := []TemplateTestCase{
		{
			name:                 "testCase 1 - Could not parse originID",
			URL:                  "/" + partnerID + "/task-definition-templates/script/badOriginID",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantGetTemplatesForExecutionMS,
		},
		{
			name: "testCase 2 - Could not get template by id,server error",
			URL:  "/" + partnerID + "/task-definition-templates/script/" + originID,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, gomock.Any(), userHasNOCAccess).
					Return(models.TemplateDetails{}, errors.New("error"))
				return tc
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTemplatesForExecutionMS,
		},
		{
			name: "testCase 3 - Could not found template",
			URL:  "/" + partnerID + "/task-definition-templates/script/" + originID,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, gomock.Any(), userHasNOCAccess).
					Return(models.TemplateDetails{}, models.TemplateNotFoundError{})
				return tc
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name: "testCase 4 - Successes",
			URL:  "/" + partnerID + "/task-definition-templates/script/" + originID,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, gomock.Any(), userHasNOCAccess).
					Return(models.TemplateDetails{Name: "TemplateDetails", PartnerID: partnerID}, nil)
				return tc
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: models.TemplateDetails{PartnerID: partnerID, Name: "TemplateDetails"},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.templateCacheMock, tc.userMock)
			router := getTemplateRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)

			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
			// expected error message checking
			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
			// expected body checking
			if tc.expectedErrorMessage == "" {
				expectedBody, err := json.Marshal(tc.expectedBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if !bytes.Equal(expectedBody, w.Body.Bytes()) {
					t.Errorf("Wanted %v but got %v", string(expectedBody), w.Body.String())
				}
			}
		})
	}
}

// getMockedService is a function that returns configured mocked TaskDefinitionService
func getMockedService(mockController *gomock.Controller, tc mock.TemplateCacheConf, userMock user.Service) TemplateService {
	mockCache := mock.NewMockTemplateCache(mockController)
	if tc != nil {
		mockCache = tc(mockCache)
	}

	userServiceMock := userMock
	if userServiceMock == nil {
		userServiceMock = user.NewMock("", partnerID, "", "", true)
	}
	return NewTemplateService(mockCache, userServiceMock, http.DefaultClient)
}

func getTemplateRouter(service TemplateService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{partnerID}/task-definition-templates", service.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/task-definition-templates/{type}", service.GetByType).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/task-definition-templates/{type}/{originID}", service.GetByOriginID).Methods(http.MethodGet)
	return router
}
