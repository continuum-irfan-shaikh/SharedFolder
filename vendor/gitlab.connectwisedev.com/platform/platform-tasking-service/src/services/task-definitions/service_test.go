package taskDefinitions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/app-loader"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	partnerID        = "50016364"
	definitionID     = "10000000-0000-0000-0000-000000000001"
	userHasNOCAccess = true
)

var (
	definitionUUID, _ = gocql.ParseUUID(definitionID)

	validTaskDefinitionDetails = models.TaskDefinitionDetails{
		TaskDefinition: models.TaskDefinition{
			ID:          definitionUUID,
			OriginID:    definitionUUID,
			Name:        "name2",
			Categories:  []string{"NotCustom", customCategory},
			Type:        "script",
			Description: "description2",
		},
		UserParameters: "",
	}

	validTaskDefinitionDetailsWithNotCustomCategory = models.TaskDefinitionDetails{
		TaskDefinition: models.TaskDefinition{
			OriginID:    definitionUUID,
			Name:        "name2",
			Categories:  []string{"NotCustom"},
			Type:        "script",
			Description: "description2",
		},
		UserParameters: "",
	}

	validTaskDefinitionDetailsWithUserParameters = models.TaskDefinitionDetails{
		TaskDefinition: models.TaskDefinition{
			OriginID:    definitionUUID,
			Name:        "name2",
			Categories:  []string{customCategory},
			Type:        "script",
			Description: "description2",
		},
		UserParameters: `{"s":"s"}`,
	}
)

func init() {
	appLoader.LoadApplicationServices(true)
}

func TestTaskDefinitionService_GetByID(t *testing.T) {

	testCases := []struct {
		name                 string
		userMock             user.Service
		taskDefinitionMock   mock.TaskDefinitionConf
		templateCacheMock    mock.TemplateCacheConf
		URL                  string
		method               string
		expectedErrorMessage string
		expectedCode         int
		expectedBody         interface{}
	}{
		{
			name: "testCase 0 - success",
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{TaskDefinition: models.TaskDefinition{
						PartnerID: partnerID,
						OriginID:  definitionUUID},
					}, nil)
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{Type: "script", Name: "Template", PartnerID: partnerID, OriginID: definitionUUID}, nil)
				return tc
			},
			URL:          "/" + partnerID + "/task-definitions/" + definitionID,
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: models.TaskDefinitionDetails{
				TaskDefinition: models.TaskDefinition{PartnerID: partnerID, OriginID: definitionUUID},
			},
		},
		{
			name: "testCase 1 - could not get template from cache by originID",
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{TaskDefinition: models.TaskDefinition{
						PartnerID: partnerID,
						OriginID:  definitionUUID},
					}, nil)
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, errors.New("some error occurred"))
				return tc
			},
			URL:                  "/" + partnerID + "/task-definitions/" + definitionID,
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name: "testCase 2 - could not get taskDefDetails by definition and partner id's",
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{}, errors.New("some error occurred"))
				return td
			},
			URL:                  "/" + partnerID + "/task-definitions/" + definitionID,
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionNotFound,
		},
		{
			name: "testCase 3 - could not get taskDefDetails, TaskDefNotFoundError error",
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{},
						models.TaskDefNotFoundError{
							ID:        definitionUUID,
							PartnerID: partnerID,
						})
				return td
			},
			URL:                  "/" + partnerID + "/task-definitions/" + definitionID,
			method:               http.MethodGet,
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionNotFound,
		},
		{
			name:                 "testCase 4 - could not get parse UUID",
			URL:                  "/" + partnerID + "/task-definitions/badUUID",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionIDHasBadFormat,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service := getMockedService(mockController, tc.taskDefinitionMock, tc.templateCacheMock, tc.userMock)
			router := getTaskDefinitionRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
			//expected error message checking
			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}

			// error expected body checking
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

func TestTaskDefinitionService_GetByPartnerID(t *testing.T) {

	testCases := []struct {
		name                 string
		userMock             user.Service
		taskDefinitionMock   mock.TaskDefinitionConf
		templateCacheMock    mock.TemplateCacheConf
		URL                  string
		method               string
		expectedErrorMessage string
		expectedCode         int
		expectedBody         []models.TaskDefinitionDetails
	}{
		{
			name: "testCase 0 - Successes",
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetAllByPartnerID(gomock.Any(), partnerID).
					Return([]models.TaskDefinitionDetails{{TaskDefinition: models.TaskDefinition{PartnerID: partnerID, OriginID: definitionUUID}}}, nil)
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, definitionUUID, userHasNOCAccess)
				return tc
			},
			URL:          "/" + partnerID + "/task-definitions",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: []models.TaskDefinitionDetails{{TaskDefinition: models.TaskDefinition{PartnerID: partnerID, OriginID: definitionUUID}}},
		},
		{
			name: "testCase 0.5 - Couldn't get Template",
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetAllByPartnerID(gomock.Any(), partnerID).
					Return([]models.TaskDefinitionDetails{{TaskDefinition: models.TaskDefinition{PartnerID: partnerID, OriginID: definitionUUID}}}, nil)
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, definitionUUID, userHasNOCAccess).
					Return(models.TemplateDetails{}, errors.New("err"))
				return tc
			},
			URL:          "/" + partnerID + "/task-definitions",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
		{
			name: "testCase 1 - Couldn't get taskDefinitions",
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetAllByPartnerID(gomock.Any(), partnerID).
					Return(nil, errors.New("some error occurred"))
				return td
			},
			URL:                  "/" + partnerID + "/task-definitions",
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionByPartnerNotFound,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.taskDefinitionMock, tc.templateCacheMock, tc.userMock)
			router := getTaskDefinitionRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)

			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
			// expected body checking
			if tc.expectedErrorMessage == "" {
				var expectedBody []byte
				var err error

				if tc.expectedBody != nil {
					expectedBody, err = json.Marshal(tc.expectedBody)
					if err != nil {
						t.Errorf("Cannot parse expected body")
					}
				}

				if !bytes.Equal(expectedBody, w.Body.Bytes()) && len(expectedBody) != 0 {
					t.Errorf("Wanted %v but \n got %v", string(expectedBody), w.Body.String())
				}
			}
		})
	}
}

func TestTaskDefinitionService_Create(t *testing.T) {

	validTaskDefinitionDetails = models.TaskDefinitionDetails{
		TaskDefinition: models.TaskDefinition{
			OriginID:    definitionUUID,
			Name:        "name2",
			Categories:  []string{"NotCustom", customCategory},
			Type:        "script",
			Description: "description2",
		},
		UserParameters: "",
	}

	testCases := []struct {
		name                 string
		URL                  string
		userMock             user.Service
		taskDefinitionMock   mock.TaskDefinitionConf
		templateCacheMock    mock.TemplateCacheConf
		bodyToSend           models.TaskDefinitionDetails
		method               string
		expectedCode         int
		expectedErrorMessage string
		expectedBody         models.TaskDefinitionDetails
		hasUIDHeader         bool
	}{
		{
			name:                 "testCase 1 - no uid header",
			method:               http.MethodPost,
			expectedErrorMessage: errorcode.ErrorUIDHeaderIsEmpty,
			URL:                  "/" + partnerID + "/task-definitions",
			expectedCode:         http.StatusBadRequest,
		},
		{
			name: "testCase 2 - could not validate taskDetails (no name,no description)",
			URL:  "/" + partnerID + "/task-definitions",
			bodyToSend: models.TaskDefinitionDetails{
				TaskDefinition: models.TaskDefinition{
					OriginID:    definitionUUID,
					Name:        "",
					Categories:  []string{customCategory},
					Type:        "script",
					Description: "",
				},
			},
			hasUIDHeader:         true,
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:         "testCase 3 - Could not get template by origin id",
			URL:          "/" + partnerID + "/task-definitions",
			hasUIDHeader: true,
			bodyToSend:   validTaskDefinitionDetails,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, errors.New("some error occurred"))
				return tc
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name:         "testCase 4 - Could not get template by origin id, TemplateNotFoundError",
			URL:          "/" + partnerID + "/task-definitions",
			hasUIDHeader: true,
			bodyToSend:   validTaskDefinitionDetails,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), gomock.Any(), definitionUUID, userHasNOCAccess).
					Return(models.TemplateDetails{}, models.TemplateNotFoundError{})
				return tc
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name:         "testCase 4.2 - Could not get template by origin id, TemplateNotFoundError",
			URL:          "/" + partnerID + "/task-definitions",
			hasUIDHeader: true,
			bodyToSend:   validTaskDefinitionDetails,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), gomock.Any(), definitionUUID, userHasNOCAccess)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(true)
				return tc
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionExists,
		},
		{
			name:         "testCase 4.5 - Could not get template by origin id, TemplateNotFoundError",
			URL:          "/" + partnerID + "/task-definitions",
			hasUIDHeader: true,
			bodyToSend:   validTaskDefinitionDetails,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), gomock.Any(), definitionUUID, userHasNOCAccess)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(false)
				return tc
			},
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					Exists(gomock.Any(), partnerID, validTaskDefinitionDetails.Name).
					Return(true)
				return td
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionExists,
		},
		{
			name:         "testCase 5 - Upsert error",
			URL:          "/" + partnerID + "/task-definitions",
			hasUIDHeader: true,
			bodyToSend:   validTaskDefinitionDetails,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, nil)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(false)
				return tc
			},
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					Exists(gomock.Any(), partnerID, validTaskDefinitionDetails.Name).
					Return(false)
				td.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(errors.New(""))
				return td
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskDefinitionToDB,
		},
		{
			name:         "testCase 6 - Successes",
			URL:          "/" + partnerID + "/task-definitions",
			hasUIDHeader: true,
			bodyToSend:   validTaskDefinitionDetails,
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, nil)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(false)
				return tc
			},
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					Exists(gomock.Any(), partnerID, validTaskDefinitionDetails.Name).
					Return(false)
				td.EXPECT().
					Upsert(gomock.Any(), gomock.Any())
				return td
			},
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			expectedBody: models.TaskDefinitionDetails{
				TaskDefinition: models.TaskDefinition{
					Description: validTaskDefinitionDetails.Description,
					OriginID:    validTaskDefinitionDetails.OriginID,
					Type:        validTaskDefinitionDetails.Type,
				},
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.taskDefinitionMock, tc.templateCacheMock, tc.userMock)
			router := getTaskDefinitionRouter(service)

			reqBody, err := json.Marshal(tc.bodyToSend)
			if err != nil {
				t.Errorf("Couldn't marshal request body: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, bytes.NewBuffer(reqBody))

			if tc.hasUIDHeader {
				r.Header.Add("uid", "Admin")
			}

			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
			// expected body checking
			if tc.expectedErrorMessage == "" {
				var gotBody models.TaskDefinitionDetails
				err := json.Unmarshal(w.Body.Bytes(), &gotBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if gotBody.OriginID != tc.expectedBody.OriginID {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.OriginID, gotBody.OriginID)
				}
				if gotBody.Description != tc.expectedBody.Description {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.Description, gotBody.Description)
				}
				if gotBody.Type != tc.expectedBody.Type {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.Type, gotBody.Type)
				}
			}
		})
	}
}

func TestTaskDefinitionService_UpdateByID(t *testing.T) {
	var (
		definitionUUID, _ = gocql.ParseUUID(definitionID)

		validTaskDefinitionDetails = models.TaskDefinitionDetails{
			TaskDefinition: models.TaskDefinition{
				OriginID:    definitionUUID,
				Name:        "name2",
				Categories:  []string{customCategory},
				Type:        "script",
				Description: "description2",
			},
			UserParameters: "",
		}

		validTaskDefinitionDetailsWithNotCustomCategory = models.TaskDefinitionDetails{
			TaskDefinition: models.TaskDefinition{
				OriginID:    definitionUUID,
				Name:        "name2",
				Categories:  []string{"NotCustom"},
				Type:        "script",
				Description: "description2",
			},
			UserParameters: "",
		}

		validTaskDefinitionDetailsWithUserParameters = models.TaskDefinitionDetails{
			TaskDefinition: models.TaskDefinition{
				OriginID:    definitionUUID,
				Name:        "name2",
				Categories:  []string{customCategory},
				Type:        "script",
				Description: "description2",
			},
			UserParameters: `{"s":"s"}`,
		}
	)

	testCases := []struct {
		name                 string
		URL                  string
		method               string
		userMock             user.Service
		taskDefinitionMock   mock.TaskDefinitionConf
		templateCacheMock    mock.TemplateCacheConf
		bodyToSend           models.TaskDefinitionDetails
		expectedCode         int
		expectedErrorMessage string
		expectedBody         models.TaskDefinitionDetails
	}{
		{
			name:                 "testCase 0 - Could not receive old task definition,bad uuid",
			URL:                  "/1/task-definitions/badUUID",
			expectedErrorMessage: errorcode.ErrorTaskDefinitionIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
			method:               http.MethodPut,
		},
		{
			name: "testCase 1 - Could not parse bodyToSend from request in while preparing upsert",
			URL:  "/" + partnerID + "/task-definitions/" + definitionID,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{TaskDefinition: models.TaskDefinition{
						PartnerID: partnerID,
						OriginID:  definitionUUID},
					}, nil)
				return td
			},
			method:               http.MethodPut,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:       "testCase 2 - Could not get templates",
			URL:        "/" + partnerID + "/task-definitions/" + definitionID,
			bodyToSend: validTaskDefinitionDetails,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{TaskDefinition: models.TaskDefinition{
						PartnerID: partnerID,
						OriginID:  definitionUUID},
					}, nil)
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, errors.New("some error occurred"))
				return tc
			},
			method:               http.MethodPut,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name:       "testCase 2.5 - Could not upsert, task with the name exists",
			URL:        "/" + partnerID + "/task-definitions/" + definitionID,
			bodyToSend: validTaskDefinitionDetails,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(validTaskDefinitionDetails, nil)
				td.EXPECT().
					CanBeUpdated(gomock.Any(), partnerID, validTaskDefinitionDetails.Name, validTaskDefinitionDetails.ID).
					Return(false, nil)
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, nil)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(false)
				return tc
			},
			method:               http.MethodPut,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionExists,
		},
		{
			name:       "testCase 2.6 - Could not upsert, CanBeUpdated returned err",
			URL:        "/" + partnerID + "/task-definitions/" + definitionID,
			bodyToSend: validTaskDefinitionDetails,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(validTaskDefinitionDetails, nil)
				td.EXPECT().
					CanBeUpdated(gomock.Any(), partnerID, validTaskDefinitionDetails.Name, validTaskDefinitionDetails.ID).
					Return(true, errors.New("err"))
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, nil)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(false)
				return tc
			},
			method:               http.MethodPut,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskDefinitionToDB,
		},
		{
			name:       "testCase 3 - Could not upsert",
			URL:        "/" + partnerID + "/task-definitions/" + definitionID,
			bodyToSend: validTaskDefinitionDetails,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(validTaskDefinitionDetails, nil)
				td.EXPECT().
					CanBeUpdated(gomock.Any(), partnerID, validTaskDefinitionDetails.Name, validTaskDefinitionDetails.ID).
					Return(true, nil)
				td.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(errors.New("some error occurred"))
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, nil)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(false)
				return tc
			},
			method:               http.MethodPut,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskDefinitionToDB,
		},
		{
			name:       "testCase 4 - Could update, userParameters are not nil but templatesJson schema is empty",
			URL:        "/" + partnerID + "/task-definitions/" + definitionID,
			bodyToSend: validTaskDefinitionDetailsWithUserParameters,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{TaskDefinition: models.TaskDefinition{
						PartnerID: partnerID,
						OriginID:  definitionUUID},
					}, nil)
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, nil)
				return tc
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:       "testCase 5 - Updated",
			URL:        "/" + partnerID + "/task-definitions/" + definitionID,
			bodyToSend: validTaskDefinitionDetailsWithNotCustomCategory,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(validTaskDefinitionDetails, nil)
				td.EXPECT().
					CanBeUpdated(gomock.Any(), partnerID, validTaskDefinitionDetails.Name, validTaskDefinitionDetails.ID).
					Return(true, nil)
				td.EXPECT().
					Upsert(gomock.Any(), gomock.Any())
				return td
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTaskDefinitionDetails.OriginID, userHasNOCAccess).
					Return(models.TemplateDetails{}, nil)
				tc.EXPECT().ExistsWithName(gomock.Any(), gomock.Any(), validTaskDefinitionDetails.Name).Return(false)
				return tc
			},
			method:       http.MethodPut,
			expectedCode: http.StatusCreated,
			expectedBody: models.TaskDefinitionDetails{
				TaskDefinition: models.TaskDefinition{
					PartnerID:  partnerID,
					OriginID:   validTaskDefinitionDetails.OriginID,
					Categories: validTaskDefinitionDetails.Categories,
					Name:       validTaskDefinitionDetails.Name,
					Type:       validTaskDefinitionDetails.Type,
				},
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.taskDefinitionMock, tc.templateCacheMock, tc.userMock)
			router := getTaskDefinitionRouter(service)

			reqBody, err := json.Marshal(tc.bodyToSend)
			if err != nil {
				t.Errorf("Couldn't marshal request body: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, bytes.NewBuffer(reqBody))
			r.Header.Add("uid", "Admin")

			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}

			// expected body checking
			if tc.expectedErrorMessage == "" {
				var gotBody models.TaskDefinitionDetails
				err := json.Unmarshal(w.Body.Bytes(), &gotBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if gotBody.OriginID != tc.expectedBody.OriginID {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.OriginID, gotBody.OriginID)
				}

				if gotBody.PartnerID != tc.expectedBody.PartnerID {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.PartnerID, gotBody.PartnerID)
				}

				if gotBody.Type != tc.expectedBody.Type {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.PartnerID, gotBody.PartnerID)
				}
			}
		})
	}
}

func TestTaskDefinitionService_DeleteByID(t *testing.T) {
	testCases := []struct {
		name                 string
		initiatedByHeader    string
		userMock             user.Service
		taskDefinitionMock   mock.TaskDefinitionConf
		templateCacheMock    mock.TemplateCacheConf
		URL                  string
		method               string
		expectedCode         int
		expectedErrorMessage string
	}{
		{
			name:                 "testCase0 - empty initiatedByHeader header",
			initiatedByHeader:    "",
			URL:                  "/" + partnerID + "/task-definitions/" + definitionID,
			method:               http.MethodDelete,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorUIDHeaderIsEmpty,
		},
		{
			name:              "testCase 1 - Could not get taskDefs by id",
			initiatedByHeader: "Admin",
			URL:               "/" + partnerID + "/task-definitions/" + definitionID,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(models.TaskDefinitionDetails{}, errors.New("some error occurred"))
				return td
			},
			method:               http.MethodDelete,
			expectedErrorMessage: errorcode.ErrorTaskDefinitionNotFound,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name:              "testCase 2 - Could not upsert taskDefs",
			initiatedByHeader: "Admin",
			URL:               "/" + partnerID + "/task-definitions/" + definitionID,
			method:            http.MethodDelete,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(validTaskDefinitionDetails, nil)
				td.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(errors.New("some error occurred"))
				return td
			},
			expectedErrorMessage: errorcode.ErrorCantDeleteTaskDefinition,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name:              "testCase 3 - Deleted",
			initiatedByHeader: "Admin",
			URL:               "/" + partnerID + "/task-definitions/" + definitionID,
			method:            http.MethodDelete,
			taskDefinitionMock: func(td *mock.MockTaskDefinitionPersistence) *mock.MockTaskDefinitionPersistence {
				td.EXPECT().
					GetByID(gomock.Any(), partnerID, definitionUUID).
					Return(validTaskDefinitionDetails, nil)
				td.EXPECT().
					Upsert(gomock.Any(), gomock.Any())
				return td
			},
			expectedCode: http.StatusNoContent,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.taskDefinitionMock, tc.templateCacheMock, tc.userMock)
			router := getTaskDefinitionRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			r.Header.Add("uid", tc.initiatedByHeader)

			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

// getMockedService is a function that returns configured mocked TaskDefinitionService
func getMockedService(mockController *gomock.Controller, td mock.TaskDefinitionConf, tc mock.TemplateCacheConf, userMock user.Service) TaskDefinitionService {
	mockTaskDefinition := mock.NewMockTaskDefinitionPersistence(mockController)
	mockCache := mock.NewMockTemplateCache(mockController)
	userServiceMock := userMock
	if td != nil {
		mockTaskDefinition = td(mockTaskDefinition)
	}
	if tc != nil {
		mockCache = tc(mockCache)
	}

	if userServiceMock == nil {
		userServiceMock = user.NewMock("", partnerID, "", "", true)
	}
	return NewTaskDefinitionService(mockTaskDefinition, mockCache, http.DefaultClient, userServiceMock, &mockrepositories.EncryptionServiceMock{})
}

func getTaskDefinitionRouter(service TaskDefinitionService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{partnerID}/task-definitions/{definitionID}", service.GetByID).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/task-definitions", service.GetByPartnerID).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/task-definitions", service.Create).Methods(http.MethodPost)
	router.HandleFunc("/{partnerID}/task-definitions/{definitionID}", service.DeleteByID).Methods(http.MethodDelete)
	router.HandleFunc("/{partnerID}/task-definitions/{definitionID}", service.UpdateByID).Methods(http.MethodPut)
	return router
}
