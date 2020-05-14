package accessControl

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apiModel "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/entitlement"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	"github.com/gorilla/mux"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func init() {
	logger.Load(config.Config.Log)
	config.Load()
	config.Config.FeaturesForRoutes = map[string]string{
		"/tasking/v1/partners/goodID/tasks/data":               "TASKING_TASKS_HOME_PAGE",
		"/tasking/v1/partners/{partnerID}/tasks/data/xlsx":     "TASKING_TASKS_HOME_PAGE",
		"/tasking/v1/partners/{partnerID}/tasks/{taskID}/data": "TASKING_TASKS_HOME_PAGE",
	}
	translation.MockTranslations()
	// just for test coverage
	Load(nil)
	Load(http.DefaultClient)

	router.HandleFunc(pathWithFeature, entitlementHandler).Methods("GET")
	router.HandleFunc(defaultPath, entitlementHandler).Methods("GET")
}

var (
	nextWasInvoked  bool
	router          = mux.NewRouter()
	defaultPath     = "/tasking/v1/partners/goodID"
	pathWithFeature = "/tasking/v1/partners/goodID/tasks/data"
	emptyPayload    []apiModel.Feature
	fullPayload     = []apiModel.Feature{
		{Name: TaskingDynamicGroupsFeature},
		{Name: TaskingSitesFeature},
	}
	entitlementHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		IsPartnerAuthorizedForRoute(w, r, func(_ http.ResponseWriter, _ *http.Request) {
			nextWasInvoked = true
		})
	})
)

func TestAccessControlMiddlewareWithPartnerID(t *testing.T) {
	var testCases = []struct {
		name              string
		URL               string
		partnerID         string
		payload           []apiModel.Feature
		nextMustBeInvoked bool
	}{
		{
			name:              "Default access",
			URL:               defaultPath,
			partnerID:         "goodID",
			payload:           fullPayload,
			nextMustBeInvoked: true,
		},
		{
			name:              "Unauthorized",
			URL:               pathWithFeature,
			partnerID:         "goodID",
			payload:           fullPayload,
			nextMustBeInvoked: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			registerResponder(tc.partnerID, tc.payload, t)
			recorder := httptest.NewRecorder()
			nextWasInvoked = false

			request := httptest.NewRequest("GET", tc.URL, nil).WithContext(context.TODO())
			router.ServeHTTP(recorder, request)

			if nextWasInvoked != tc.nextMustBeInvoked {
				t.Errorf("\nTestCase[%s]: Invokation error: nextWasInvoked %t, but should be %t", tc.name, nextWasInvoked, tc.nextMustBeInvoked)
			}
			if !nextWasInvoked && (recorder.Code != http.StatusForbidden) {
				t.Errorf("\nTestCase[%s]: Bad status returned: want %d but got %d", tc.name, http.StatusForbidden, recorder.Code)
			}
		})
	}
}

func TestIsPartnerAuthorizedToRunTask(t *testing.T) {
	var testCases = []struct {
		name           string
		task           *models.Task
		payload        []apiModel.Feature
		expectedResult bool
	}{
		{
			name:           "Partner ID is authorized to run task according to TaskingBasicFeature rule",
			payload:        fullPayload,
			expectedResult: true,
			task: &models.Task{
				PartnerID: "goodPartner",
			},
		},
		{
			name:           "Partner ID is authorized to run task according to TaskingDynamicGroupsFeature rule",
			payload:        fullPayload,
			expectedResult: true,
			task: &models.Task{
				PartnerID: "goodPartner",
				Targets: models.Target{
					Type: models.DynamicGroup,
					IDs:  []string{"nonempty"},
				},
			},
		},
		{
			name:           "Partner ID is authorized to run task according to TaskingSitesFeature rule",
			payload:        fullPayload,
			expectedResult: true,
			task: &models.Task{
				PartnerID: "goodPartner",
				Targets: models.Target{
					Type: models.Site,
					IDs:  []string{"nonempty"},
				},
			},
		},
		{
			name:           "Partner ID is not authorized to run task",
			payload:        emptyPayload,
			expectedResult: false,
			task: &models.Task{
				PartnerID: "badPartner",
				Targets: models.Target{
					Type: models.DynamicGroup,
					IDs:  []string{"nonempty"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			registerResponder(tc.task.PartnerID, tc.payload, t)
			actualResult := IsPartnerAuthorizedToRunTask(tc.task)
			if actualResult != tc.expectedResult {
				t.Errorf("\nTestCase[%s]: expected result: %t, but got %t", tc.name, tc.expectedResult, actualResult)
			}
		})
	}
}

func registerResponder(partnerID string, payload []apiModel.Feature, t *testing.T) {
	entitlementURL := fmt.Sprintf("%s/partners/%s/features", config.Config.EntitlementMsURL, partnerID)

	t.Logf("Registered HTTP responder on URL: %s\n", entitlementURL)
	httpmock.RegisterResponder("GET", entitlementURL,
		func(req *http.Request) (*http.Response, error) {
			var (
				resp *http.Response
				err  error
			)
			resp, err = httpmock.NewJsonResponse(http.StatusOK, payload)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
			}
			return resp, nil
		},
	)
}
