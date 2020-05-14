package taskCounters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	modelMocks "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/model-mocks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func init() {
	config.Load()
	logger.Load(config.Config.Log)
	translation.Load()
}

func TestUpdateCounters(t *testing.T) {
	var (
		defaultCtx      = context.TODO()
		defaultPartner  = "1"
		defaultEndpoint = gocql.TimeUUID()
		defaultCount    = 5
	)

	testCases := []struct {
		name           string
		taskCounterDAO repository.TaskCounter
		expectedErr    error
	}{
		{
			name: "UpdateCounters - Good case",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return nil, fmt.Errorf("it's OK anyway :)")
				},
				IncreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
				DecreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
			},
		},
		{
			name: "UpdateCounters - IncreaseCounter error",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return nil, nil
				},
				IncreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return fmt.Errorf("IncreaseCounterErr")
				},
				DecreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
			},
			expectedErr: fmt.Errorf("IncreaseCounterErr"),
		},
		{
			name: "UpdateCounters - DecreaseCounter error",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return nil, nil
				},
				IncreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
				DecreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return fmt.Errorf("DecreaseCounterErr")
				},
			},
			expectedErr: fmt.Errorf("DecreaseCounterErr"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			service := New(tc.taskCounterDAO)

			err := service.updateCounters(
				defaultCtx,
				defaultPartner,
				defaultEndpoint,
				defaultCount,
			)

			if err != nil && tc.expectedErr == nil {
				t.Fatalf("%s: unexpected error: %v\n", tc.name, err)
			}

			if err == nil && tc.expectedErr != nil {
				t.Fatalf("%s: expected error [%v] but got <nil>\n", tc.name, tc.expectedErr)
			}

			if err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
				t.Fatalf("%s: expected error [%v] but got [%v]\n", tc.name, tc.expectedErr, err)
			}
		})
	}
}

func TestGetCountersByPartner(t *testing.T) {
	defaultPartner := "1"
	defaultUUID := gocql.TimeUUID()

	testCases := []struct {
		name             string
		taskCounterDAO   repository.TaskCounter
		r                *http.Request
		expectedHTTPCode int
		expectedResult   interface{}
	}{
		{
			name: "GetCountersByPartner - good",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return []models.TaskCount{
						{
							ManagedEndpointID: defaultUUID,
							Count:             1,
						},
					}, nil
				},
			},
			r: httptest.NewRequest(
				http.MethodGet,
				`/partners/`+defaultPartner+`/tasks/count`,
				nil,
			),
			expectedHTTPCode: http.StatusOK,
			expectedResult: []models.TaskCount{
				{
					ManagedEndpointID: defaultUUID,
					Count:             1,
				},
			},
		},
		{
			name: "GetCountersByPartner - with err",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return []models.TaskCount{
						{
							ManagedEndpointID: defaultUUID,
							Count:             0,
						},
					}, fmt.Errorf("err")
				},
			},
			r: httptest.NewRequest(
				http.MethodGet,
				`/partners/`+defaultPartner+`/tasks/count`,
				nil,
			),
			expectedHTTPCode: http.StatusOK,
			expectedResult: []models.TaskCount{
				{
					ManagedEndpointID: defaultUUID,
					Count:             0,
				},
			},
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			var (
				router  = mux.NewRouter()
				service = New(tc.taskCounterDAO)
				w       = httptest.NewRecorder()
			)

			router.HandleFunc(
				`/partners/{partnerID}/tasks/count`,
				http.HandlerFunc(service.GetCountersByPartner)).Methods(http.MethodGet)

			router.ServeHTTP(w, test.r)

			if w.Code != test.expectedHTTPCode {
				t.Fatalf("expected code [%d] but got [%d]", test.expectedHTTPCode, w.Code)
			}

			expectedBody, err := json.Marshal(test.expectedResult)
			if err != nil {
				t.Fatalf("%s: json marshal failed: %v", test.name, err)
			}

			if w.Body.String() != string(expectedBody) {
				t.Fatalf("%s: \nexpected result: \n%s\nbut got: \n%s\n", test.name, string(expectedBody), w.Body.String())
			}

		})
	}
}

func TestGetCountersByPartnerAndEndpoint(t *testing.T) {
	defaultPartner := "1"
	defaultUUID := gocql.TimeUUID()

	testCases := []struct {
		name                 string
		taskCounterDAO       repository.TaskCounter
		r                    *http.Request
		expectedHTTPCode     int
		expectedResult       interface{}
		expectedErrorMessage string
	}{
		{
			name: "GetCountersByPartnerAndEndpoint - good",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return []models.TaskCount{
						{
							ManagedEndpointID: defaultUUID,
							Count:             2,
						},
					}, nil
				},
			},
			r: httptest.NewRequest(
				http.MethodGet,
				`/partners/`+defaultPartner+`/tasks/managed-endpoints/`+defaultUUID.String()+`/count`,
				nil,
			),
			expectedHTTPCode: http.StatusOK,
			expectedResult: []models.TaskCount{
				{
					ManagedEndpointID: defaultUUID,
					Count:             2,
				},
			},
		},
		{
			name: "GetCountersByPartnerAndEndpoint - no counter and err",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return nil, fmt.Errorf("err")
				},
			},
			r: httptest.NewRequest(
				http.MethodGet,
				`/partners/`+defaultPartner+`/tasks/managed-endpoints/`+defaultUUID.String()+`/count`,
				nil,
			),
			expectedHTTPCode: http.StatusOK,
			expectedResult: []models.TaskCount{
				{
					ManagedEndpointID: defaultUUID,
					Count:             0,
				},
			},
		},
		{
			name: "GetCountersByPartnerAndEndpoint - bad endpoint",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return nil, nil
				},
			},
			r: httptest.NewRequest(
				http.MethodGet,
				`/partners/`+defaultPartner+`/tasks/managed-endpoints/bad_endpoint/count`,
				nil,
			),
			expectedHTTPCode:     http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			var (
				router  = mux.NewRouter()
				service = New(tc.taskCounterDAO)
				w       = httptest.NewRecorder()
			)

			router.HandleFunc(`/partners/{partnerID}/tasks/managed-endpoints/{managedEndpointID}/count`,
				http.HandlerFunc(service.GetCountersByPartnerAndEndpoint)).Methods(http.MethodGet)

			router.ServeHTTP(w, test.r)

			if w.Code != test.expectedHTTPCode {
				t.Fatalf("expected code [%d] but got [%d]", test.expectedHTTPCode, w.Code)
			}

			expectedBody, err := json.Marshal(test.expectedResult)
			if err != nil {
				t.Fatalf("%s: json marshal failed: %v", test.name, err)
			}

			if w.Body.String() != string(expectedBody) && tc.expectedResult != nil {
				t.Fatalf("%s: \nexpected result: \n%s\nbut got: \n%s\n", test.name, string(expectedBody), w.Body.String())
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func TestRecalculateAllCounters(t *testing.T) {
	type Errors struct {
		HasErrors bool `json:"HasErrors"`
	}

	var (
		defaultPartner = "1"
		defaultUUID    = gocql.TimeUUID()
	)

	testCases := []struct {
		name             string
		taskCounterDAO   repository.TaskCounter
		taskDAO          models.TaskPersistence
		expectedHTTPCode int
		expectedResult   interface{}
	}{
		{
			name: "RecalculateAllCounters - good",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetAllPartnersF: func(ctx context.Context) (partnerIDs map[string]struct{}, err error) {
					return map[string]struct{}{
						defaultPartner: {},
					}, nil
				},
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return nil, nil
				},
				IncreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
				DecreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
			},
			taskDAO: modelMocks.TaskCustomizableMock{
				GetCountsByPartnerF: func(ctx context.Context, partnerID string) ([]models.TaskCount, error) {
					return []models.TaskCount{
						{
							ManagedEndpointID: defaultUUID,
							Count:             3,
						},
					}, nil
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedResult: struct {
				HasErrors bool
			}{
				HasErrors: false,
			},
		},

		{
			name: "RecalculateAllCounters - GetAllPartners err",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetAllPartnersF: func(ctx context.Context) (partnerIDs map[string]struct{}, err error) {
					return nil, fmt.Errorf("err")
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedResult:   Errors{HasErrors: true},
		},

		{
			name: "RecalculateAllCounters - Increase/Decrease errors",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetAllPartnersF: func(ctx context.Context) (partnerIDs map[string]struct{}, err error) {
					return map[string]struct{}{
						defaultPartner: {},
					}, nil
				},
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					if partnerID == defaultPartner {
						return []models.TaskCount{
							{
								ManagedEndpointID: defaultUUID,
								Count:             4,
							},
						}, nil
					}
					return nil, nil
				},
				IncreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return fmt.Errorf("IncreaseCounterErr")
				},
				DecreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return fmt.Errorf("DecreaseCounterErr")
				},
			},
			taskDAO: modelMocks.TaskCustomizableMock{
				GetCountsByPartnerF: func(ctx context.Context, partnerID string) ([]models.TaskCount, error) {
					return []models.TaskCount{
						{
							ManagedEndpointID: defaultUUID,
							Count:             3,
						},
					}, nil
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedResult: struct {
				HasErrors bool
			}{
				HasErrors: true,
			},
		},

		{
			name: "RecalculateAllCounters - GetCountsByPartner errors",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetAllPartnersF: func(ctx context.Context) (partnerIDs map[string]struct{}, err error) {
					return map[string]struct{}{
						defaultPartner: {},
					}, nil
				},
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					if partnerID == defaultPartner {
						return []models.TaskCount{
							{
								ManagedEndpointID: defaultUUID,
								Count:             5,
							},
						}, nil
					}
					return nil, nil
				},
				IncreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
				DecreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
			},
			taskDAO: modelMocks.TaskCustomizableMock{
				GetCountsByPartnerF: func(ctx context.Context, partnerID string) ([]models.TaskCount, error) {
					return nil, fmt.Errorf("GetCountsByPartnerErr")
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedResult: struct {
				HasErrors bool
			}{
				HasErrors: true,
			},
		},

		{
			name: "RecalculateAllCounters - GetCounters errors",
			taskCounterDAO: modelMocks.TaskCounterCustomizableMock{
				GetAllPartnersF: func(ctx context.Context) (partnerIDs map[string]struct{}, err error) {
					return map[string]struct{}{
						defaultPartner: {},
					}, nil
				},
				GetCountersF: func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error) {
					return nil, fmt.Errorf("GetCountersErr")
				},
				IncreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
				DecreaseCounterF: func(partnerID string, counters []models.TaskCount, isExternal bool) error {
					return nil
				},
			},
			taskDAO: modelMocks.TaskCustomizableMock{
				GetCountsByPartnerF: func(ctx context.Context, partnerID string) ([]models.TaskCount, error) {
					return []models.TaskCount{
						{
							ManagedEndpointID: defaultUUID,
							Count:             5,
						},
					}, nil
				},
			},
			expectedHTTPCode: http.StatusOK,
			expectedResult: struct {
				HasErrors bool
			}{
				HasErrors: true,
			},
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			models.TaskPersistenceInstance = tc.taskDAO
			var (
				router  = mux.NewRouter()
				service = New(tc.taskCounterDAO)
				w       = httptest.NewRecorder()
				r       = httptest.NewRequest(
					http.MethodPost,
					`/tasking/v1/recalculate-counters`,
					nil,
				)
			)
			router.HandleFunc(`/tasking/v1/recalculate-counters`,
				http.HandlerFunc(service.RecalculateAllCounters)).Methods(http.MethodPost)

			router.ServeHTTP(w, r)

			if w.Code != test.expectedHTTPCode {
				t.Fatalf("expected code [%d] but got [%d]", test.expectedHTTPCode, w.Code)
			}

			expectedBody, err := json.Marshal(test.expectedResult)
			if err != nil {
				t.Fatalf("%s: json marshal failed: %v", test.name, err)
			}

			if w.Body.String() != string(expectedBody) {
				t.Fatalf("%s: \nexpected result: \n%s\nbut got: \n%s\n", test.name, string(expectedBody), w.Body.String())
			}
		})
	}
}
