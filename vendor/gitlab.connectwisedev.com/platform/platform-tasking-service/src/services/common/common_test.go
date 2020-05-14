package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gocql/gocql"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	commonlibRest "gitlab.connectwisedev.com/platform/platform-common-lib/src/web/rest"
	"gopkg.in/jarcoal/httpmock.v1"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
	testURL    = "http://localhost:12121/test-path"
)

const defaultMsg = `failed on unexpected value of result "%v"`

func init() {
	config.Load()
	translation.MockTranslations()
}

func TestSendStatuses(t *testing.T) {
	recorder := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:12121/tasking/version", nil)
	SendBadRequest(recorder, r, "error_cant_decode_input_data")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusBadRequest)
	}

	recorder = httptest.NewRecorder()
	SendInternalServerError(recorder, r, "error_cant_decode_input_data")
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusInternalServerError)
	}

	recorder = httptest.NewRecorder()
	SendNoContent(recorder)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusNoContent)
	}

	recorder = httptest.NewRecorder()
	SendNotFound(recorder, r, "hello")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusNoContent)
	}

	recorder = httptest.NewRecorder()
	SendStatusCodeWithMessage(recorder, r, http.StatusConflict, "error_cant_decode_input_data")
	if recorder.Code != http.StatusConflict {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusConflict)
	}

	recorder = httptest.NewRecorder()
	SendStatusOkWithMessage(recorder, r, "OK msg")
	if recorder.Code != http.StatusOK {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusOK)
	}

	recorder = httptest.NewRecorder()
	SendForbidden(recorder, r, "test")
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusForbidden)
	}

	recorder = httptest.NewRecorder()
	SendCreated(recorder, r, "test")
	if recorder.Code != http.StatusCreated {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusCreated)
	}

	r.Header.Set("Accept-Language", "0011")
	config.Config.DefaultLanguage = ""
	recorder = httptest.NewRecorder()
	SendInternalServerError(recorder, r, "test")
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Got %v, want %v", recorder.Code, http.StatusInternalServerError)
	}
	config.Config.DefaultLanguage = "en-US"
}

func TestGenerateTargetIDs(t *testing.T) {
	targetIDs := GenerateTargetIDs()
	if len(targetIDs) == 0 {
		t.Fatal("The target ids can't be empty")
	}
}

func TestRenderJSON(t *testing.T) {
	channelToBreak := make(chan int)
	defer close(channelToBreak)
	if err := logger.Load(config.Config.Log); err != nil {
		t.Fatal("can't initialize logger")
	}
	tests := []struct {
		name           string
		response       interface{}
		expectedStatus int
	}{
		{"0", channelToBreak, http.StatusInternalServerError},

		{"1", "hello world", http.StatusOK},

		{"2", 123, http.StatusOK},

		{"3", []byte("slice of bytes"), http.StatusOK},

		{"4", 1.2, http.StatusOK},

		{"5", []interface{}{}, http.StatusOK},
	}
	for _, tt := range tests {
		recorder := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			RenderJSON(recorder, tt.response)
			if recorder.Code != tt.expectedStatus {
				t.Fatalf("Wrong Status returned")
			} else if recorder.Body == nil {
				t.Fatalf("Got %v, want %v", recorder.Body, tt.response)
			}
		})
	}

	resp := http.Response{}
	CloseRespBody(&resp)
	stringReader := strings.NewReader("shiny!")
	stringReadCloser := ioutil.NopCloser(stringReader)
	resp.Body = stringReadCloser
	CloseRespBody(&resp)
}

func TestUserCtx(t *testing.T) {
	UserFromCtx(context.Background())
	ctx := context.WithValue(context.Background(), config.UserKeyCTX, entities.AgentConfigPayload{})
	UserFromCtx(ctx)
	ctx = context.WithValue(context.Background(), config.UserKeyCTX, entities.User{})
	UserFromCtx(ctx)

	UsersEndpointsFromCtx(context.Background())
	ctx = context.WithValue(context.Background(), config.UserEndPointsKeyCTX, entities.AgentConfigPayload{})
	UsersEndpointsFromCtx(ctx)
	ctx = context.WithValue(context.Background(), config.UserEndPointsKeyCTX, []string{})
	UsersEndpointsFromCtx(ctx)
	SliceToMap([]string{"temporal"})
}

func TestRenderJSONFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		response []byte
	}{
		{"1", []byte("slice of bytes")},
	}
	for _, tt := range tests {
		recorder := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			RenderJSONFromBytes(recorder, tt.response)
			if recorder.Code != http.StatusOK {
				t.Fatalf("Wrong Status returned")
			} else if !reflect.DeepEqual(recorder.Body.Bytes(), tt.response) {
				t.Fatalf("Got %T, want %T", recorder.Body, tt.response)
			}
		})
	}
}

func TestRenderJSONCreated(t *testing.T) {
	channelToBreak := make(chan int)
	defer close(channelToBreak)

	tests := []struct {
		name           string
		response       interface{}
		expectedStatus int
	}{
		{"0", channelToBreak, http.StatusInternalServerError},

		{"1", "hello world", http.StatusCreated},

		{"2", 123, http.StatusCreated},

		{"3", []byte("slice of bytes"), http.StatusCreated},

		{"4", 1.2, http.StatusCreated},
	}
	for _, tt := range tests {
		recorder := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			RenderJSONCreated(recorder, tt.response)
			if recorder.Code != tt.expectedStatus {
				t.Fatalf("Wrong Status returned")
			} else if recorder.Body == nil {
				t.Fatalf("Got %v, want %v", recorder.Body, tt.response)
			}
		})
	}
}

func TestGetGeneralInfo(t *testing.T) {
	tests := []struct {
		name string
		want commonlibRest.GeneralInfo
	}{
		{`non-empty`,
			commonlibRest.GeneralInfo{
				time.Now(),
				"serviceName",
				"serviceProvider",
				serviceVersion,
				"solutionName",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			if got := GetGeneralInfo(); !reflect.DeepEqual(got.ServiceVersion, tt.want.ServiceVersion) {
				t.Fatalf("GetGeneralInfo() ServiceVersion = %v, want %v", got.ServiceVersion, tt.want.ServiceVersion)

			} else if got.TimeStampUTC.Before(tt.want.TimeStampUTC) {
				t.Fatalf("Invalid TimeStampUTC returned")

			}
		})
	}
}

func TestFetchStartIndexAndNumberOfRows(t *testing.T) {
	tests := []struct {
		name                     string
		vals                     url.Values
		startIndex, numberOfRows int
	}{
		{"-1", map[string][]string{"startIndex": {"-1"}, "numberOfRows": {"-1"}}, 1, 1},
		{"12345", map[string][]string{"startIndex": {"12345"}, "numberOfRows": {"12345"}}, 12345, 12345},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotRows := FetchStartIndexAndNumberOfRows(tt.vals)

			if gotIndex < StartIndexDefault {
				t.Fatalf("Got %v, want not more than %v", gotIndex, StartIndexDefault)
			}
			if gotRows < NumberOfRowsDefault || gotRows > NumberOfRowsMaximum {
				t.Fatalf("Got %v, want between %v and %v", gotRows, NumberOfRowsDefault, NumberOfRowsMaximum)
			}
		})
	}
}

func TestGetLanguage(t *testing.T) {
	request, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := "en-US"
	request.Header[acceptLanguageHeader] = []string{expected}
	if GetLanguage(request) != expected {
		t.Fatalf("Got %v, want %v", GetLanguage(request), expected)
	}
}

func TestGetQuestionMarkString(t *testing.T) {
	var tests = []struct {
		count          int
		expectedString string
	}{
		{
			count:          -1,
			expectedString: "",
		},
		{
			count:          0,
			expectedString: "",
		},
		{
			count:          1,
			expectedString: "?",
		},
		{
			count:          5,
			expectedString: "?, ?, ?, ?, ?",
		},
	}
	for _, test := range tests {
		actualResult := GetQuestionMarkString(test.count)
		if actualResult != test.expectedString {
			t.Fatalf("Actual string %v doesn't equal to expected %v", actualResult, test.expectedString)
		}

	}
}

func TestExtractUUID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:12121/ts/v1/partners/1d4400c0", nil)

	tt := []struct {
		name    string
		uuidKey string
		error   string
	}{
		{
			name:    "EndpointHasBadFormatNoArgs",
			uuidKey: "managedEndpointID",
			error:   errorcode.ErrorEndpointIDHasBadFormat,
		},
		{
			name:    "TaskIDHasBadFormat",
			uuidKey: "taskID",
			error:   errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name:    "TaskInstanceIDHasBadFormat",
			uuidKey: "taskInstanceID",
			error:   errorcode.ErrorTaskInstanceIDHasBadFormat,
		},
		{
			name:    "CantDecodeInputData",
			uuidKey: "CantDecodeInputData",
			error:   errorcode.ErrorCantDecodeInputData,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			_, err := ExtractUUID("TaskResultsService.Get", w, r, test.uuidKey)
			if err == nil {
				t.Fatalf("TestExtractUUID(): must be an error: %v", test.error)
			}
			bodyBytes, _ := ioutil.ReadAll(w.Body)
			bodyString := string(bodyBytes)

			if w.Code != http.StatusBadRequest || !strings.Contains(bodyString, test.error) {
				t.Fatalf("must consist: %s , but got body: %s", test.error, bodyString)
			}
		})
	}
}

func TestExtractOptionalCount(t *testing.T) {
	var (
		w       = httptest.NewRecorder()
		baseURL = "localhost:12121/ts/v1/partners/1d4400c0//task-execution-results/managed-endpoints/f7b6a7df-522f-43d7-bee0-b4d2075d9ce3"
		tests   = []struct {
			name          string
			url           string
			isSpecified   bool
			expextedCount int
			errorExpected bool
		}{
			{
				name:          "no count specified",
				url:           baseURL,
				isSpecified:   false,
				expextedCount: 0,
				errorExpected: false,
			},
			{
				name:          "invalid count",
				url:           baseURL + "?count=a",
				isSpecified:   true,
				errorExpected: true,
			},
			{
				name:          "valid count",
				url:           baseURL + "?count=10",
				isSpecified:   true,
				expextedCount: 10,
				errorExpected: false,
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", tt.url, nil)
			actualCount, isSpecified, err := ExtractOptionalCount(tt.name, w, r)
			if !tt.errorExpected {
				if isSpecified != tt.isSpecified {
					t.Fatalf("TestExtractOptionalCount(): expected that count specified = %t but got %t\n", tt.isSpecified, isSpecified)
				}
				if actualCount != tt.expextedCount {
					t.Fatalf("TestExtractOptionalCount(): expected %v but got %v\n", tt.expextedCount, actualCount)
				} else {
					return
				}
			}

			if err == nil {
				t.Fatal("TestExtractOptionalCount(): should return an error")
			}
			bodyBytes, _ := ioutil.ReadAll(w.Body)
			bodyString := string(bodyBytes)

			if w.Code != http.StatusBadRequest || !strings.Contains(bodyString, errorcode.ErrorCountVarHasBadFormat) {
				t.Fatalf("must consist: %s , but got body: %s", errorcode.ErrorCountVarHasBadFormat, bodyString)
			}
		})
	}
}

func TestHTTPRequestWithRetry(t *testing.T) {
	config.Config.RetryStrategy.MaxNumberOfRetries = 2
	config.Config.RetryStrategy.RetrySleepIntervalSec = 1
	logger.Load(config.Config.Log)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	_, err := HTTPRequestWithRetry(context.TODO(), httpClient, http.MethodPost, testURL, []byte("test body"))
	if err == nil {
		t.Fatal("Must be an error after all retires!")
	}

	registerResponder()

	response, newErr := HTTPRequestWithRetry(context.TODO(), httpClient, http.MethodPost, testURL, []byte("test body"))
	if newErr != nil {
		t.Fatalf("Error: %v", err)
	}

	defer response.Body.Close()
	respBody, newErr := ioutil.ReadAll(response.Body)

	if newErr != nil {
		t.Fatalf("Error: %v", err)
	}

	if string(respBody) != "ok" {
		t.Fatalf("Excpected body %v, but actual is %v", "ok", string(respBody))
	}
}

func registerResponder() {
	httpmock.RegisterResponder("POST", testURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, "ok"), nil
		},
	)
}

func TestGetUTCTimeInLocation(t *testing.T) {
	locationName := "America/Toronto"
	location, _ := time.LoadLocation(locationName)
	startTime, err := time.Parse(time.RFC3339, "2017-12-19T09:12:05-05:00")
	if err != nil {
		t.Errorf("No error expexted. Got: %s", err)
	}
	expectedDateTime := startTime.In(location)
	actualStartTime := GetTimeInLocation(startTime, location)

	if !reflect.DeepEqual(actualStartTime, expectedDateTime) {
		t.Fatalf("Actual time %v isn't equal to expected %v", actualStartTime, expectedDateTime)
	}
}

func TestGetNextRunTime(t *testing.T) {
	// =========================== HOURLY ========================================
	runTimeHourly1, _ := time.Parse(time.RFC3339, "2090-10-13T09:30:00Z")
	runTimeHourlyExpected1, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")

	runTimeHourly2, _ := time.Parse(time.RFC3339, "2090-10-13T22:22:00Z")
	runTimeHourlyExpected2, _ := time.Parse(time.RFC3339, "2090-10-14T05:22:00Z")

	runTimeHourly3, _ := time.Parse(time.RFC3339, "2090-10-13T11:00:00Z")
	runTimeHourlyExpected3, _ := time.Parse(time.RFC3339, "2090-10-13T13:00:00Z")

	runTimeHourly4, _ := time.Parse(time.RFC3339, "2090-10-13T11:00:00Z")
	runTimeHourlyExpected4, _ := time.Parse(time.RFC3339, "2090-10-14T17:00:00Z")

	// =========================== DAILY =========================================

	runTimeDaily1, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeDailyExpected1, _ := time.Parse(time.RFC3339, "2090-10-23T11:35:00Z")

	runTimeDaily2, _ := time.Parse(time.RFC3339, "2088-02-10T11:35:00Z")
	runTimeDailyExpected2, _ := time.Parse(time.RFC3339, "2088-03-21T11:35:00Z")

	runTimeDaily3, _ := time.Parse(time.RFC3339, "2088-02-11T11:35:00Z")
	runTimeDailyExpected3, _ := time.Parse(time.RFC3339, "2088-02-12T11:35:00Z")

	runTimeDaily4, _ := time.Parse(time.RFC3339, "2088-02-11T11:35:00Z")
	runTimeDailyExpected4, _ := time.Parse(time.RFC3339, "2088-02-13T11:35:00Z")

	runTimeDaily5, _ := time.Parse(time.RFC3339, "2088-12-30T11:35:00Z")
	runTimeDailyExpected5, _ := time.Parse(time.RFC3339, "2088-12-31T11:35:00Z")

	runTimeDaily6, _ := time.Parse(time.RFC3339, "2088-12-30T11:35:00Z")
	runTimeDailyExpected6, _ := time.Parse(time.RFC3339, "2089-01-01T11:35:00Z")

	runTimeDaily7, _ := time.Parse(time.RFC3339, "2088-12-30T11:35:00Z")
	runTimeDailyExpected7, _ := time.Parse(time.RFC3339, "2089-01-01T00:00:00Z")

	runTimeDaily8, _ := time.Parse(time.RFC3339, "2088-12-30T00:01:00Z")
	runTimeDailyExpected8, _ := time.Parse(time.RFC3339, "2089-01-01T00:01:00Z")

	runTimeDaily9, _ := time.Parse(time.RFC3339, "2088-12-30T00:00:00Z")
	runTimeDailyExpected9, _ := time.Parse(time.RFC3339, "2089-01-01T00:00:00Z")

	runTimeDaily10, _ := time.Parse(time.RFC3339, "2088-12-30T23:59:00Z")
	runTimeDailyExpected10, _ := time.Parse(time.RFC3339, "2089-01-01T23:59:00Z")

	// =========================== WEEKLY =========================================

	runTimeWeekly1, _ := time.Parse(time.RFC3339, "2090-10-09T11:35:00Z")
	runTimeWeeklyExpected1, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")

	runTimeWeekly2, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeWeeklyExpected2, _ := time.Parse(time.RFC3339, "2090-10-16T11:35:00Z")

	runTimeWeekly3, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeWeeklyExpected3, _ := time.Parse(time.RFC3339, "2090-10-23T11:35:00Z")

	runTimeWeekly4, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeWeeklyExpected4, _ := time.Parse(time.RFC3339, "2091-01-01T11:35:00Z")

	runTimeWeekly5, _ := time.Parse(time.RFC3339, "2090-01-03T22:35:29.104Z")
	runTimeWeeklyExpected5, _ := time.Parse(time.RFC3339, "2091-01-09T21:10:00Z")

	runTimeWeekly6, _ := time.Parse(time.RFC3339, "2090-01-03T22:35:29.10Z")
	runTimeWeeklyExpected6, _ := time.Parse(time.RFC3339, "2091-01-09T00:00:00Z")

	runTimeWeekly7, _ := time.Parse(time.RFC3339, "2090-01-03T00:00:00.00Z")
	runTimeWeeklyExpected7, _ := time.Parse(time.RFC3339, "2091-01-09T00:00:00Z")

	runTimeWeekly8, _ := time.Parse(time.RFC3339, "2090-01-03T00:01:00.00Z")
	runTimeWeeklyExpected8, _ := time.Parse(time.RFC3339, "2091-01-09T00:01:00Z")

	runTimeWeekly9, _ := time.Parse(time.RFC3339, "2090-01-03T23:59:00.00Z")
	runTimeWeeklyExpected9, _ := time.Parse(time.RFC3339, "2091-01-09T23:59:00Z")

	// =========================== MONTHLY =========================================

	runTimeMonthly1, _ := time.Parse(time.RFC3339, "2090-10-13T10:30:00Z")
	runTimeMonthlyExpected1, _ := time.Parse(time.RFC3339, "2090-10-21T10:30:00Z")

	runTimeMonthly2, _ := time.Parse(time.RFC3339, "2090-10-21T10:30:00Z")
	runTimeMonthlyExpected2, _ := time.Parse(time.RFC3339, "2091-05-01T10:30:00Z")

	runTimeMonthly3, _ := time.Parse(time.RFC3339, "2090-01-30T10:30:00Z")
	runTimeMonthlyExpected3, _ := time.Parse(time.RFC3339, "2090-02-28T10:30:00Z")

	runTimeMonthly4, _ := time.Parse(time.RFC3339, "2090-03-25T20:22:00Z")
	runTimeMonthlyExpected4, _ := time.Parse(time.RFC3339, "2091-11-25T20:22:00Z")

	runTimeMonthly5, _ := time.Parse(time.RFC3339, "2090-03-25T20:22:00Z")
	runTimeMonthlyExpected5, _ := time.Parse(time.RFC3339, "2090-12-02T20:22:00Z")

	runTimeMonthly6, _ := time.Parse(time.RFC3339, "2090-03-25T20:22:00Z")
	runTimeMonthlyExpected6, _ := time.Parse(time.RFC3339, "2090-04-02T00:00:00Z")

	runTimeMonthly7, _ := time.Parse(time.RFC3339, "2090-01-01T00:00:00Z")
	runTimeMonthlyExpected7, _ := time.Parse(time.RFC3339, "2090-02-01T00:00:00Z")

	runTimeMonthly8, _ := time.Parse(time.RFC3339, "2090-01-01T09:09:09Z")
	runTimeMonthlyExpected8, _ := time.Parse(time.RFC3339, "2090-02-01T00:01:00Z")

	runTimeMonthly9, _ := time.Parse(time.RFC3339, "2090-01-01T23:59:00Z")
	runTimeMonthlyExpected9, _ := time.Parse(time.RFC3339, "2090-02-01T23:59:00Z")

	runTimeMonthly10, _ := time.Parse(time.RFC3339, "2090-01-01T00:01:00Z")
	runTimeMonthlyExpected10, _ := time.Parse(time.RFC3339, "2090-02-01T00:01:00Z")

	var tests = []struct {
		testName         string
		lastRunTime      time.Time
		schedule         apiModels.Schedule
		expectedRunTime  time.Time
		expectedSchedule apiModels.Schedule
		expectedErr      bool
	}{
		// ==================== HOURLY ==============================
		{
			testName:         "hourly 1",
			lastRunTime:      runTimeHourly1,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Hourly}},
			expectedRunTime:  runTimeHourlyExpected1,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Hourly}},
		},
		{
			testName:         "hourly 2",
			lastRunTime:      runTimeHourly2,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Hourly}},
			expectedRunTime:  runTimeHourlyExpected2,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Hourly}},
		},
		{
			testName:         "hourly 3",
			lastRunTime:      runTimeHourly3,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Hourly}},
			expectedRunTime:  runTimeHourlyExpected3,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Hourly}},
		},
		{
			testName:         "hourly 4",
			lastRunTime:      runTimeHourly4,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 30, Frequency: apiModels.Hourly}},
			expectedRunTime:  runTimeHourlyExpected4,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 30, Frequency: apiModels.Hourly}},
		},
		// ==================== DAILY ==============================

		{
			testName:         "daily 1",
			lastRunTime:      runTimeDaily1,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 13, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 283}},
			expectedRunTime:  runTimeDailyExpected1,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 13, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 296}},
		},

		{
			testName:         "daily 2",
			lastRunTime:      runTimeDaily2,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 40, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 41}},
			expectedRunTime:  runTimeDailyExpected2,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 40, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 81}},
		},

		{
			testName:         "daily 3",
			lastRunTime:      runTimeDaily3,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 42}},
			expectedRunTime:  runTimeDailyExpected3,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 43}},
		},

		{
			testName:         "daily 4",
			lastRunTime:      runTimeDaily4,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 42}},
			expectedRunTime:  runTimeDailyExpected4,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 44}},
		},

		{
			testName:         "daily 5",
			lastRunTime:      runTimeDaily5,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 365}},
			expectedRunTime:  runTimeDailyExpected5,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 366}},
		},

		{
			testName:         "daily 6",
			lastRunTime:      runTimeDaily6,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 365}},
			expectedRunTime:  runTimeDailyExpected6,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), Period: 1}},
		},

		{
			testName:         "daily 7",
			lastRunTime:      runTimeDaily7,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), Period: 365}},
			expectedRunTime:  runTimeDailyExpected7,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), Period: 1}},
		},

		{
			testName:         "daily 8",
			lastRunTime:      runTimeDaily8,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), Period: 365}},
			expectedRunTime:  runTimeDailyExpected8,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), Period: 1}},
		},

		{
			testName:         "daily 9",
			lastRunTime:      runTimeDaily9,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), Period: 365}},
			expectedRunTime:  runTimeDailyExpected9,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), Period: 1}},
		},

		{
			testName:         "daily 10",
			lastRunTime:      runTimeDaily10,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 23, 59, 0, 0, time.UTC), Period: 365}},
			expectedRunTime:  runTimeDailyExpected10,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily, RunTime: time.Date(0, 0, 0, 23, 59, 0, 0, time.UTC), Period: 1}},
		},

		// ==================== WEEKLY =============================
		{
			testName:         "weekly 1",
			lastRunTime:      runTimeWeekly1,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 41}},
			expectedRunTime:  runTimeWeeklyExpected1,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 41}},
		},

		{
			testName:         "weekly 2",
			lastRunTime:      runTimeWeekly2,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 41}},
			expectedRunTime:  runTimeWeeklyExpected2,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 42}},
		},

		{
			testName:         "weekly 3",
			lastRunTime:      runTimeWeekly3,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 41}},
			expectedRunTime:  runTimeWeeklyExpected3,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 43}},
		},

		{
			testName:         "weekly 4",
			lastRunTime:      runTimeWeekly4,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 12, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 41}},
			expectedRunTime:  runTimeWeeklyExpected4,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 12, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 11, 35, 0, 0, time.UTC), DaysOfWeek: []int{1, 2}, Period: 1}},
		},

		{
			testName:         "weekly 5",
			lastRunTime:      runTimeWeekly5,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 21, 10, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 1}},
			expectedRunTime:  runTimeWeeklyExpected5,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 21, 10, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 2}},
		},

		{
			testName:         "weekly 6",
			lastRunTime:      runTimeWeekly6,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 1}},
			expectedRunTime:  runTimeWeeklyExpected6,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 2}},
		},

		{
			testName:         "weekly 7",
			lastRunTime:      runTimeWeekly7,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 1}},
			expectedRunTime:  runTimeWeeklyExpected7,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 2}},
		},

		{
			testName:         "weekly 8",
			lastRunTime:      runTimeWeekly8,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 1}},
			expectedRunTime:  runTimeWeeklyExpected8,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 2}},
		},

		{
			testName:         "weekly 9",
			lastRunTime:      runTimeWeekly9,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 23, 59, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 1}},
			expectedRunTime:  runTimeWeeklyExpected9,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, RunTime: time.Date(0, 0, 0, 23, 59, 0, 0, time.UTC), DaysOfWeek: []int{2}, Period: 2}},
		},

		// ==================== MONTHLY ============================
		{
			testName:         "monthly 1",
			lastRunTime:      runTimeMonthly1,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC), DaysOfMonth: []int{1, 7, 13, 21}, Period: 10}},
			expectedRunTime:  runTimeMonthlyExpected1,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC), DaysOfMonth: []int{1, 7, 13, 21}, Period: 10}},
		},
		{
			testName:         "monthly 2",
			lastRunTime:      runTimeMonthly2,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC), DaysOfMonth: []int{1, 7, 13, 21}, Period: 10}},
			expectedRunTime:  runTimeMonthlyExpected2,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC), DaysOfMonth: []int{1, 7, 13, 21}, Period: 5}},
		},
		{
			testName:         "monthly 3",
			lastRunTime:      runTimeMonthly3,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC), DaysOfMonth: []int{30}, Period: 1}},
			expectedRunTime:  runTimeMonthlyExpected3,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC), DaysOfMonth: []int{30}, Period: 2}},
		},
		{
			testName:         "monthly 4",
			lastRunTime:      runTimeMonthly4,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 20, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 20, 22, 0, 0, time.UTC), DaysOfMonth: []int{25}, Period: 3}},
			expectedRunTime:  runTimeMonthlyExpected4,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 20, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 20, 22, 0, 0, time.UTC), DaysOfMonth: []int{25}, Period: 11}},
		},
		{
			testName:         "monthly 5",
			lastRunTime:      runTimeMonthly5,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 9, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 20, 22, 0, 0, time.UTC), DaysOfMonth: []int{2, 25}, Period: 3}},
			expectedRunTime:  runTimeMonthlyExpected5,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 9, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 20, 22, 0, 0, time.UTC), DaysOfMonth: []int{2, 25}, Period: 12}},
		},
		{
			testName:         "monthly 6",
			lastRunTime:      runTimeMonthly6,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfMonth: []int{2, 25}, Period: 3}},
			expectedRunTime:  runTimeMonthlyExpected6,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfMonth: []int{2, 25}, Period: 4}},
		},
		{
			testName:         "monthly 7",
			lastRunTime:      runTimeMonthly7,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 1}},
			expectedRunTime:  runTimeMonthlyExpected7,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 2}},
		},
		{
			testName:         "monthly 8",
			lastRunTime:      runTimeMonthly8,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 1}},
			expectedRunTime:  runTimeMonthlyExpected8,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 2}},
		},

		{
			testName:         "monthly 9",
			lastRunTime:      runTimeMonthly9,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 23, 59, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 1}},
			expectedRunTime:  runTimeMonthlyExpected9,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 23, 59, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 2}},
		},

		{
			testName:         "monthly 10",
			lastRunTime:      runTimeMonthly10,
			schedule:         apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 1}},
			expectedRunTime:  runTimeMonthlyExpected10,
			expectedSchedule: apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, RunTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC), DaysOfMonth: []int{1}, Period: 2}},
		},
		{
			testName:    "error, empty time, cannot parse cron string",
			schedule:    apiModels.Schedule{},
			expectedErr: true,
		},
	}

	for _, test := range tests {
		actualRunTime, actualSchedule, err := GetNextRunTime(test.lastRunTime, test.schedule)
		if err != nil && !test.expectedErr {
			t.Fatalf("Test %s, Error: %v", test.testName, err)
		}

		if actualRunTime != test.expectedRunTime && !test.expectedErr {
			t.Fatalf("Test %s failed, expected runtime %v but got %v", test.testName, test.expectedRunTime, actualRunTime)
		}

		if !assert.ObjectsAreEqualValues(actualSchedule, test.expectedSchedule) && !test.expectedErr {
			t.Fatalf("Test %s failed, expected schedule %v but got %v", test.testName, test.expectedSchedule, actualSchedule)
		}
	}
}

func TestCalcNextRunTime(t *testing.T) {
	RegisterTestingT(t)

	// =========================== HOURLY ========================================
	runTimeHourly1, _ := time.Parse(time.RFC3339, "2090-10-13T09:30:00Z")
	runTimeHourlyExpected1, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")

	runTimeHourly2, _ := time.Parse(time.RFC3339, "2090-10-13T22:22:00Z")
	runTimeHourlyExpected2, _ := time.Parse(time.RFC3339, "2090-10-14T05:22:00Z")

	runTimeHourly3, _ := time.Parse(time.RFC3339, "2090-10-13T11:00:00Z")
	runTimeHourlyExpected3, _ := time.Parse(time.RFC3339, "2090-10-13T13:00:00Z")

	runTimeHourly4, _ := time.Parse(time.RFC3339, "2090-10-13T11:00:00Z")
	runTimeHourlyExpected4, _ := time.Parse(time.RFC3339, "2090-10-14T17:00:00Z")

	// =========================== DAILY =========================================

	runTimeDaily1, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeDailyExpected1, _ := time.Parse(time.RFC3339, "2090-10-23T11:35:00Z")

	runTimeDaily2, _ := time.Parse(time.RFC3339, "2088-02-10T11:35:00Z")
	runTimeDailyExpected2, _ := time.Parse(time.RFC3339, "2088-03-21T11:35:00Z")

	runTimeDaily3, _ := time.Parse(time.RFC3339, "2088-02-11T11:35:00Z")
	runTimeDailyExpected3, _ := time.Parse(time.RFC3339, "2088-02-12T11:35:00Z")

	runTimeDaily4, _ := time.Parse(time.RFC3339, "2088-02-11T11:35:00Z")
	runTimeDailyExpected4, _ := time.Parse(time.RFC3339, "2088-02-13T11:35:00Z")

	runTimeDaily5, _ := time.Parse(time.RFC3339, "2088-12-30T11:35:00Z")
	runTimeDailyExpected5, _ := time.Parse(time.RFC3339, "2088-12-31T11:35:00Z")

	runTimeDaily6, _ := time.Parse(time.RFC3339, "2088-12-30T11:35:00Z")
	runTimeDailyExpected6, _ := time.Parse(time.RFC3339, "2089-01-01T11:35:00Z")

	runTimeDaily7, _ := time.Parse(time.RFC3339, "2088-12-30T11:35:00Z")
	runTimeDailyExpected7, _ := time.Parse(time.RFC3339, "2089-01-01T11:35:00Z")

	runTimeDaily8, _ := time.Parse(time.RFC3339, "2088-12-30T00:01:00Z")
	runTimeDailyExpected8, _ := time.Parse(time.RFC3339, "2089-01-01T00:01:00Z")

	runTimeDaily9, _ := time.Parse(time.RFC3339, "2088-12-30T00:00:00Z")
	runTimeDailyExpected9, _ := time.Parse(time.RFC3339, "2089-01-01T00:00:00Z")

	runTimeDaily10, _ := time.Parse(time.RFC3339, "2088-12-30T23:59:00Z")
	runTimeDailyExpected10, _ := time.Parse(time.RFC3339, "2089-01-01T23:59:00Z")

	// =========================== WEEKLY =========================================

	runTimeWeekly1, _ := time.Parse(time.RFC3339, "2090-10-09T11:35:00Z")
	runTimeWeeklyExpected1, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")

	runTimeWeekly2, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeWeeklyExpected2, _ := time.Parse(time.RFC3339, "2090-10-16T11:35:00Z")

	runTimeWeekly3, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeWeeklyExpected3, _ := time.Parse(time.RFC3339, "2090-10-23T11:35:00Z")

	runTimeWeekly4, _ := time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
	runTimeWeeklyExpected4, _ := time.Parse(time.RFC3339, "2091-01-01T11:35:00Z")

	runTimeWeekly5, _ := time.Parse(time.RFC3339, "2090-01-03T22:35:29.104Z")
	runTimeWeeklyExpected5, _ := time.Parse(time.RFC3339, "2091-01-09T22:35:29.104Z")

	runTimeWeekly6, _ := time.Parse(time.RFC3339, "2090-01-03T22:35:29.10Z")
	runTimeWeeklyExpected6, _ := time.Parse(time.RFC3339, "2091-01-09T22:35:29.10Z")

	runTimeWeekly7, _ := time.Parse(time.RFC3339, "2090-01-03T00:00:00.00Z")
	runTimeWeeklyExpected7, _ := time.Parse(time.RFC3339, "2091-01-09T00:00:00Z")

	runTimeWeekly8, _ := time.Parse(time.RFC3339, "2090-01-03T00:01:00.00Z")
	runTimeWeeklyExpected8, _ := time.Parse(time.RFC3339, "2091-01-09T00:01:00Z")

	runTimeWeekly9, _ := time.Parse(time.RFC3339, "2090-01-03T23:59:00.00Z")
	runTimeWeeklyExpected9, _ := time.Parse(time.RFC3339, "2091-01-09T23:59:00Z")

	// =========================== MONTHLY =========================================

	runTimeMonthly1, _ := time.Parse(time.RFC3339, "2090-10-13T10:30:00Z")
	runTimeMonthlyExpected1, _ := time.Parse(time.RFC3339, "2090-10-21T10:30:00Z")

	runTimeMonthly2, _ := time.Parse(time.RFC3339, "2090-10-21T10:30:00Z")
	runTimeMonthlyExpected2, _ := time.Parse(time.RFC3339, "2091-05-01T10:30:00Z")

	runTimeMonthly3, _ := time.Parse(time.RFC3339, "2090-01-30T10:30:00Z")
	runTimeMonthlyExpected3, _ := time.Parse(time.RFC3339, "2090-02-28T10:30:00Z")

	runTimeMonthly4, _ := time.Parse(time.RFC3339, "2090-03-25T20:22:00Z")
	runTimeMonthlyExpected4, _ := time.Parse(time.RFC3339, "2091-11-25T20:22:00Z")

	runTimeMonthly5, _ := time.Parse(time.RFC3339, "2090-03-25T20:22:00Z")
	runTimeMonthlyExpected5, _ := time.Parse(time.RFC3339, "2090-12-02T20:22:00Z")

	runTimeMonthly6, _ := time.Parse(time.RFC3339, "2090-03-25T20:22:00Z")
	runTimeMonthlyExpected6, _ := time.Parse(time.RFC3339, "2090-04-02T20:22:00Z")

	runTimeMonthly7, _ := time.Parse(time.RFC3339, "2090-01-01T00:00:00Z")
	runTimeMonthlyExpected7, _ := time.Parse(time.RFC3339, "2090-02-01T00:00:00Z")

	runTimeMonthly8, _ := time.Parse(time.RFC3339, "2090-01-01T09:09:00Z")
	runTimeMonthlyExpected8, _ := time.Parse(time.RFC3339, "2090-02-01T09:09:00Z")

	runTimeMonthly9, _ := time.Parse(time.RFC3339, "2090-01-01T23:59:00Z")
	runTimeMonthlyExpected9, _ := time.Parse(time.RFC3339, "2090-02-01T23:59:00Z")

	runTimeMonthly10, _ := time.Parse(time.RFC3339, "2090-01-01T00:01:00Z")
	runTimeMonthlyExpected10, _ := time.Parse(time.RFC3339, "2090-02-01T00:01:00Z")

	runTimeMonthly11, _ := time.Parse(time.RFC3339, "2040-01-15T12:48:00Z")
	runTimeMonthlyExpected11, _ := time.Parse(time.RFC3339, "2040-01-28T12:48:00Z")

	runTimeMonthly12, _ := time.Parse(time.RFC3339, "2040-01-28T12:48:00Z")
	runTimeMonthlyExpected12, _ := time.Parse(time.RFC3339, "2040-02-01T12:48:00Z")

	runTimeMonthly13, _ := time.Parse(time.RFC3339, "2039-01-31T12:48:00Z")
	runTimeMonthlyExpected13, _ := time.Parse(time.RFC3339, "2039-02-28T12:48:00Z")

	runTimeMonthly14, _ := time.Parse(time.RFC3339, "2040-01-31T12:48:00Z")
	runTimeMonthlyExpected14, _ := time.Parse(time.RFC3339, "2040-02-29T12:48:00Z")

	runTimeMonthly15, _ := time.Parse(time.RFC3339, "2039-12-31T12:48:00Z")
	runTimeMonthlyExpected15, _ := time.Parse(time.RFC3339, "2040-02-29T12:48:00Z")

	runTimeMonthly16, _ := time.Parse(time.RFC3339, "2039-02-27T12:48:00Z")
	runTimeMonthlyExpected16, _ := time.Parse(time.RFC3339, "2039-02-28T12:48:00Z")

	runTimeMonthly17, _ := time.Parse(time.RFC3339, "2039-02-28T12:48:00Z")
	runTimeMonthlyExpected17, _ := time.Parse(time.RFC3339, "2039-03-01T12:48:00Z")

	runTimeMonthly18, _ := time.Parse(time.RFC3339, "2040-02-28T12:48:00Z")
	runTimeMonthlyExpected18, _ := time.Parse(time.RFC3339, "2040-02-29T12:48:00Z")

	runTimeMonthly19, _ := time.Parse(time.RFC3339, "2040-02-27T12:48:00Z")
	runTimeMonthlyExpected19, _ := time.Parse(time.RFC3339, "2040-02-29T12:48:00Z")

	runTimeMonthly20, _ := time.Parse(time.RFC3339, "2040-02-27T12:48:00Z")
	runTimeMonthlyExpected20, _ := time.Parse(time.RFC3339, "2040-02-29T12:48:00Z")

	runTimeMonthly21, _ := time.Parse(time.RFC3339, "2039-02-27T12:48:00Z")
	runTimeMonthlyExpected21, _ := time.Parse(time.RFC3339, "2039-02-28T12:48:00Z")

	runTimeMonthly22, _ := time.Parse(time.RFC3339, "2039-02-28T12:48:00Z")
	runTimeMonthlyExpected22, _ := time.Parse(time.RFC3339, "2039-03-10T12:48:00Z")

	runTimeMonthly23, _ := time.Parse(time.RFC3339, "2040-02-28T12:48:00Z")
	runTimeMonthlyExpected23, _ := time.Parse(time.RFC3339, "2040-02-29T12:48:00Z")

	runTimeMonthly24, _ := time.Parse(time.RFC3339, "2038-11-30T12:48:00Z")
	runTimeMonthlyExpected24, _ := time.Parse(time.RFC3339, "2039-02-10T12:48:00Z")

	runTimeMonthly25, _ := time.Parse(time.RFC3339, "2038-11-30T12:48:00Z")
	runTimeMonthlyExpected25, _ := time.Parse(time.RFC3339, "2039-02-28T12:48:00Z")

	runTimeMonthly26, _ := time.Parse(time.RFC3339, "2038-11-30T12:48:00Z")
	runTimeMonthlyExpected26, _ := time.Parse(time.RFC3339, "2038-12-06T12:48:00Z")

	runTimeMonthly27, _ := time.Parse(time.RFC3339, "2038-11-30T12:48:00Z")
	runTimeMonthlyExpected27, _ := time.Parse(time.RFC3339, "2039-02-15T12:48:00Z")

	runTimeMonthly28, _ := time.Parse(time.RFC3339, "2038-11-18T12:48:00Z")
	runTimeMonthlyExpected28, _ := time.Parse(time.RFC3339, "2038-11-26T12:48:00Z")

	runTimeMonthly29, _ := time.Parse(time.RFC3339, "2038-11-10T12:48:00Z")
	runTimeMonthlyExpected29, _ := time.Parse(time.RFC3339, "2038-11-11T12:48:00Z")

	runTimeMonthly30, _ := time.Parse(time.RFC3339, "2038-11-10T12:48:00Z")
	runTimeMonthlyExpected30, _ := time.Parse(time.RFC3339, "2038-11-19T12:48:00Z")

	runTimeMonthly31, _ := time.Parse(time.RFC3339, "2038-11-26T12:48:00Z")
	runTimeMonthlyExpected31, _ := time.Parse(time.RFC3339, "2039-02-03T12:48:00Z")

	var tests = []struct {
		testName        string
		lastRunTime     time.Time
		schedule        apiModels.Schedule
		expectedRunTime time.Time
		expectedErr     string
	}{
		// ==================== HOURLY ==============================
		{
			testName:        "hourly 1",
			lastRunTime:     runTimeHourly1,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Hourly}},
			expectedRunTime: runTimeHourlyExpected1,
		},
		{
			testName:        "hourly 2",
			lastRunTime:     runTimeHourly2,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Hourly}},
			expectedRunTime: runTimeHourlyExpected2,
		},
		{
			testName:        "hourly 3",
			lastRunTime:     runTimeHourly3,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Hourly}},
			expectedRunTime: runTimeHourlyExpected3,
		},
		{
			testName:        "hourly 4",
			lastRunTime:     runTimeHourly4,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 30, Frequency: apiModels.Hourly}},
			expectedRunTime: runTimeHourlyExpected4,
		},
		// //==================== DAILY ==============================

		{
			testName:        "daily 1",
			lastRunTime:     runTimeDaily1,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 13, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected1,
		},
		{
			testName:        "daily 2",
			lastRunTime:     runTimeDaily2,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 40, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected2,
		},

		{
			testName:        "daily 3",
			lastRunTime:     runTimeDaily3,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected3,
		},

		{
			testName:        "daily 4",
			lastRunTime:     runTimeDaily4,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected4,
		},

		{
			testName:        "daily 5",
			lastRunTime:     runTimeDaily5,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected5,
		},

		{
			testName:        "daily 6",
			lastRunTime:     runTimeDaily6,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected6,
		},

		{
			testName:        "daily 7",
			lastRunTime:     runTimeDaily7,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected7,
		},

		{
			testName:        "daily 8",
			lastRunTime:     runTimeDaily8,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected8,
		},

		{
			testName:        "daily 9",
			lastRunTime:     runTimeDaily9,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected9,
		},

		{
			testName:        "daily 10",
			lastRunTime:     runTimeDaily10,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Daily}},
			expectedRunTime: runTimeDailyExpected10,
		},

		// ==================== WEEKLY =============================
		{
			testName:        "weekly 1",
			lastRunTime:     runTimeWeekly1,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Weekly, DaysOfWeek: []int{1, 2}}},
			expectedRunTime: runTimeWeeklyExpected1,
		},

		{
			testName:        "weekly 2",
			lastRunTime:     runTimeWeekly2,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Weekly, DaysOfWeek: []int{1, 2}}},
			expectedRunTime: runTimeWeeklyExpected2,
		},

		{
			testName:        "weekly 3",
			lastRunTime:     runTimeWeekly3,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Weekly, DaysOfWeek: []int{1, 2}}},
			expectedRunTime: runTimeWeeklyExpected3,
		},

		{
			testName:        "weekly 4",
			lastRunTime:     runTimeWeekly4,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 12, Frequency: apiModels.Weekly, DaysOfWeek: []int{1, 2}}},
			expectedRunTime: runTimeWeeklyExpected4,
		},

		{
			testName:        "weekly 5",
			lastRunTime:     runTimeWeekly5,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, DaysOfWeek: []int{2}}},
			expectedRunTime: runTimeWeeklyExpected5,
		},

		{
			testName:        "weekly 6",
			lastRunTime:     runTimeWeekly6,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, DaysOfWeek: []int{2}}},
			expectedRunTime: runTimeWeeklyExpected6,
		},

		{
			testName:        "weekly 7",
			lastRunTime:     runTimeWeekly7,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, DaysOfWeek: []int{2}}},
			expectedRunTime: runTimeWeeklyExpected7,
		},

		{
			testName:        "weekly 8",
			lastRunTime:     runTimeWeekly8,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, DaysOfWeek: []int{2}}},
			expectedRunTime: runTimeWeeklyExpected8,
		},

		{
			testName:        "weekly 9",
			lastRunTime:     runTimeWeekly9,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 53, Frequency: apiModels.Weekly, DaysOfWeek: []int{2}}},
			expectedRunTime: runTimeWeeklyExpected9,
		},

		// ==================== MONTHLY ============================
		{
			testName:        "monthly 1",
			lastRunTime:     runTimeMonthly1,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 7, 13, 21}}},
			expectedRunTime: runTimeMonthlyExpected1,
		},
		{
			testName:        "monthly 2",
			lastRunTime:     runTimeMonthly2,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 7, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 7, 13, 21}}},
			expectedRunTime: runTimeMonthlyExpected2,
		},
		{
			testName:        "monthly 3",
			lastRunTime:     runTimeMonthly3,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{30}}},
			expectedRunTime: runTimeMonthlyExpected3,
		},
		{
			testName:        "monthly 4",
			lastRunTime:     runTimeMonthly4,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 20, Frequency: apiModels.Monthly, DaysOfMonth: []int{25}}},
			expectedRunTime: runTimeMonthlyExpected4,
		},
		{
			testName:        "monthly 5",
			lastRunTime:     runTimeMonthly5,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 9, Frequency: apiModels.Monthly, DaysOfMonth: []int{2, 25}}},
			expectedRunTime: runTimeMonthlyExpected5,
		},
		{
			testName:        "monthly 6",
			lastRunTime:     runTimeMonthly6,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{2, 25}}},
			expectedRunTime: runTimeMonthlyExpected6,
		},
		{
			testName:        "monthly 7",
			lastRunTime:     runTimeMonthly7,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1}}},
			expectedRunTime: runTimeMonthlyExpected7,
		},
		{
			testName:        "monthly 8",
			lastRunTime:     runTimeMonthly8,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1}}},
			expectedRunTime: runTimeMonthlyExpected8,
		},
		{
			testName:        "monthly 9",
			lastRunTime:     runTimeMonthly9,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1}}},
			expectedRunTime: runTimeMonthlyExpected9,
		},
		{
			testName:        "monthly 10",
			lastRunTime:     runTimeMonthly10,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1}}},
			expectedRunTime: runTimeMonthlyExpected10,
		},
		{
			testName:        "monthly 11",
			lastRunTime:     runTimeMonthly11,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 15, 28}}},
			expectedRunTime: runTimeMonthlyExpected11,
		},
		{
			testName:        "monthly 12",
			lastRunTime:     runTimeMonthly12,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 15, 28}}},
			expectedRunTime: runTimeMonthlyExpected12,
		},
		{
			testName:        "monthly 13",
			lastRunTime:     runTimeMonthly13,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{30, 31}}},
			expectedRunTime: runTimeMonthlyExpected13,
		},
		{
			testName:        "monthly 14",
			lastRunTime:     runTimeMonthly14,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{30, 31}}},
			expectedRunTime: runTimeMonthlyExpected14,
		},
		{
			testName:        "monthly 15",
			lastRunTime:     runTimeMonthly15,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 2, Frequency: apiModels.Monthly, DaysOfMonth: []int{30, 31}}},
			expectedRunTime: runTimeMonthlyExpected15,
		},
		{
			testName:        "monthly 16",
			lastRunTime:     runTimeMonthly16,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 15, 28, 31}}},
			expectedRunTime: runTimeMonthlyExpected16,
		},
		{
			testName:        "monthly 17",
			lastRunTime:     runTimeMonthly17,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 15, 28, 31}}},
			expectedRunTime: runTimeMonthlyExpected17,
		},
		{
			testName:        "monthly 18",
			lastRunTime:     runTimeMonthly18,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 15, 28, 31}}},
			expectedRunTime: runTimeMonthlyExpected18,
		},
		{
			testName:        "monthly 19",
			lastRunTime:     runTimeMonthly19,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{1, 15, 27, 31}}},
			expectedRunTime: runTimeMonthlyExpected19,
		},
		{
			testName:        "monthly 20",
			lastRunTime:     runTimeMonthly20,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{10, 15, 27, 29}}},
			expectedRunTime: runTimeMonthlyExpected20,
		},
		{
			testName:        "monthly 21",
			lastRunTime:     runTimeMonthly21,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{10, 15, 27, 29}}},
			expectedRunTime: runTimeMonthlyExpected21,
		},
		{
			testName:        "monthly 22",
			lastRunTime:     runTimeMonthly22,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{10, 15, 28, 29}}},
			expectedRunTime: runTimeMonthlyExpected22,
		},
		{
			testName:        "monthly 23",
			lastRunTime:     runTimeMonthly23,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, DaysOfMonth: []int{10, 15, 28, 29}}},
			expectedRunTime: runTimeMonthlyExpected23,
		},
		{
			testName:        "monthly 24",
			lastRunTime:     runTimeMonthly24,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, DaysOfMonth: []int{10, 30}}},
			expectedRunTime: runTimeMonthlyExpected24,
		},
		{
			testName:        "monthly 25",
			lastRunTime:     runTimeMonthly25,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, DaysOfMonth: []int{30}}},
			expectedRunTime: runTimeMonthlyExpected25,
		},
		{
			testName:        "monthly 26",
			lastRunTime:     runTimeMonthly26,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Monthly, WeekDays: []apiModels.WeekDay{{Day: 1, Index: 0}}}},
			expectedRunTime: runTimeMonthlyExpected26,
		},
		{
			testName:        "monthly 27",
			lastRunTime:     runTimeMonthly27,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, WeekDays: []apiModels.WeekDay{{Day: 2, Index: 2}}}},
			expectedRunTime: runTimeMonthlyExpected27,
		},
		{
			testName:        "monthly 28",
			lastRunTime:     runTimeMonthly28,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, WeekDays: []apiModels.WeekDay{{Day: 5, Index: 3}}}},
			expectedRunTime: runTimeMonthlyExpected28,
		},
		{
			testName:        "monthly 29",
			lastRunTime:     runTimeMonthly29,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, WeekDays: []apiModels.WeekDay{{Day:5, Index:2},{Day: 4, Index: 1}}}},
			expectedRunTime: runTimeMonthlyExpected29,
		},
		{
			testName:        "monthly 30",
			lastRunTime:     runTimeMonthly30,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, WeekDays: []apiModels.WeekDay{{Day:5, Index:2},{Day: 4, Index: 0}}}},
			expectedRunTime: runTimeMonthlyExpected30,
		},
		{
			testName:        "monthly 31",
			lastRunTime:     runTimeMonthly31,
			schedule:        apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, WeekDays: []apiModels.WeekDay{{Day:5, Index:2},{Day: 4, Index: 0}}}},
			expectedRunTime: runTimeMonthlyExpected31,
		},
	}

	for _, test := range tests {
		actualRunTime, err := CalcNextRunTime(test.lastRunTime, test.schedule, time.Location{})

		if err != nil {
			Ω(err.Error()).To(Equal(test.expectedErr), fmt.Sprintf(defaultMsg, test.testName))
			Ω(actualRunTime).To(Equal(time.Time{}), fmt.Sprintf(defaultMsg, test.testName))
			continue
		}

		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.testName))
		Ω(actualRunTime).To(Equal(test.expectedRunTime), fmt.Sprintf(defaultMsg, test.testName))
	}
}

func TestCalcFirstNextRunTime(t *testing.T) {
	RegisterTestingT(t)

	now := time.Now().UTC()

	// =========================== HOURLY ========================================
	runTimeHourlyExpected1, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")
	startRunTime1, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")
	endRunTime1, _ := time.Parse(time.RFC3339, "2090-10-13T18:30:00Z")
	schedule1 := apiModels.Schedule{
		StartRunTime: startRunTime1,
		EndRunTime:   endRunTime1,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 2, Frequency: apiModels.Hourly},
	}

	runTimeHourlyExpected2 := time.Time{}
	startRunTime2 := now.Add(-5 * time.Minute)
	endRunTime2 := now.Add(20 * time.Minute)
	schedule2 := apiModels.Schedule{
		StartRunTime: startRunTime2,
		EndRunTime:   endRunTime2,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 2, Frequency: apiModels.Hourly},
	}

	runTimeHourlyExpected3 := now.Add(55 * time.Minute)
	startRunTime3 := now.Add(-5 * time.Minute)
	endRunTime3 := now.Add(150 * time.Minute)
	schedule3 := apiModels.Schedule{
		StartRunTime: startRunTime3,
		EndRunTime:   endRunTime3,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Hourly},
	}

	// =========================== DAILY =========================================

	runTimeDailyExpected1, _ := time.Parse(time.RFC3339, "2090-10-13T17:30:00Z")
	startRunTimeDaily1, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")
	endRunTimeDaily1, _ := time.Parse(time.RFC3339, "2090-11-13T20:30:00Z")
	runTimeDaily1, _ := time.Parse(time.RFC3339, "2090-10-13T17:30:00Z")
	scheduleDaily1 := apiModels.Schedule{
		StartRunTime: startRunTimeDaily1,
		EndRunTime:   endRunTimeDaily1,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Daily, RunTime: runTimeDaily1},
	}

	runTimeDailyExpected2 := time.Time{}
	startRunTimeDaily2, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")
	endRunTimeDaily2, _ := time.Parse(time.RFC3339, "2090-10-14T10:30:00Z")
	runTimeDaily2, _ := time.Parse(time.RFC3339, "2090-10-13T12:30:00Z")
	scheduleDaily2 := apiModels.Schedule{
		StartRunTime: startRunTimeDaily2,
		EndRunTime:   endRunTimeDaily2,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Daily, RunTime: runTimeDaily2},
	}

	runTimeDailyExpected3, _ := time.Parse(time.RFC3339, "2090-10-14T12:30:00Z")
	startRunTimeDaily3, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")
	endRunTimeDaily3, _ := time.Parse(time.RFC3339, "2090-10-14T14:30:00Z")
	runTimeDaily3, _ := time.Parse(time.RFC3339, "2090-10-13T12:30:00Z")
	scheduleDaily3 := apiModels.Schedule{
		StartRunTime: startRunTimeDaily3,
		EndRunTime:   endRunTimeDaily3,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Daily, RunTime: runTimeDaily3},
	}

	runTimeDailyExpected4, _ := time.Parse(time.RFC3339, "2091-01-01T14:30:00Z")
	startRunTimeDaily4, _ := time.Parse(time.RFC3339, "2090-12-31T16:30:00Z")
	endRunTimeDaily4, _ := time.Parse(time.RFC3339, "2091-03-14T14:30:00Z")
	runTimeDaily4, _ := time.Parse(time.RFC3339, "2090-10-13T14:30:00Z")
	scheduleDaily4 := apiModels.Schedule{
		StartRunTime: startRunTimeDaily4,
		EndRunTime:   endRunTimeDaily4,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Daily, RunTime: runTimeDaily4},
	}

	runTimeDailyExpected5 := now.AddDate(0, 0, 1).Add(-10 * time.Minute).Truncate(time.Minute)
	startRunTimeDaily5 := now.Add(-5 * time.Minute)
	endRunTimeDaily5 := now.AddDate(0, 0, 10)
	runTimeDaily5 := now.Add(-10 * time.Minute)
	scheduleDaily5 := apiModels.Schedule{
		StartRunTime: startRunTimeDaily5,
		EndRunTime:   endRunTimeDaily5,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Daily, RunTime: runTimeDaily5},
	}

	runTimeDailyExpected6, _ := time.Parse(time.RFC3339, "2090-12-31T16:30:00Z")
	startRunTimeDaily6, _ := time.Parse(time.RFC3339, "2090-12-31T16:30:00Z")
	endRunTimeDaily6, _ := time.Parse(time.RFC3339, "2091-03-14T14:30:00Z")
	runTimeDaily6, _ := time.Parse(time.RFC3339, "2090-10-13T16:30:00Z")
	scheduleDaily6 := apiModels.Schedule{
		StartRunTime: startRunTimeDaily6,
		EndRunTime:   endRunTimeDaily6,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Daily, RunTime: runTimeDaily6},
	}

	// =========================== WEEKLY =========================================

	runTimeWeeklyExpected1, _ := time.Parse(time.RFC3339, "2031-12-17T17:30:00Z")
	startRunTimeWeekly1, _ := time.Parse(time.RFC3339, "2031-12-17T16:30:00Z")
	endRunTimeWeekly1, _ := time.Parse(time.RFC3339, "2032-10-13T20:30:00Z")
	runTimeWeekly1, _ := time.Parse(time.RFC3339, "2031-10-06T17:30:00Z")
	scheduleWeekly1 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly1,
		EndRunTime:   endRunTimeWeekly1,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly1, DaysOfWeek: []int{3, 5}},
	}

	runTimeWeeklyExpected2, _ := time.Parse(time.RFC3339, "2039-12-21T14:30:00Z")
	startRunTimeWeekly2, _ := time.Parse(time.RFC3339, "2039-12-18T16:30:00Z")
	endRunTimeWeekly2, _ := time.Parse(time.RFC3339, "2040-10-13T20:30:00Z")
	runTimeWeekly2, _ := time.Parse(time.RFC3339, "2039-10-06T14:30:00Z")
	scheduleWeekly2 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly2,
		EndRunTime:   endRunTimeWeekly2,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly2, DaysOfWeek: []int{3, 5}},
	}

	runTimeWeeklyExpected3, _ := time.Parse(time.RFC3339, "2039-12-23T14:30:00Z")
	startRunTimeWeekly3, _ := time.Parse(time.RFC3339, "2039-12-21T16:30:00Z")
	endRunTimeWeekly3, _ := time.Parse(time.RFC3339, "2040-10-13T20:30:00Z")
	runTimeWeekly3, _ := time.Parse(time.RFC3339, "2039-10-06T14:30:00Z")
	scheduleWeekly3 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly3,
		EndRunTime:   endRunTimeWeekly3,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly3, DaysOfWeek: []int{3, 5}},
	}

	runTimeWeeklyExpected4, _ := time.Parse(time.RFC3339, "2039-12-25T14:30:00Z")
	startRunTimeWeekly4, _ := time.Parse(time.RFC3339, "2039-12-21T16:30:00Z")
	endRunTimeWeekly4, _ := time.Parse(time.RFC3339, "2040-10-13T20:30:00Z")
	runTimeWeekly4, _ := time.Parse(time.RFC3339, "2039-10-06T14:30:00Z")
	scheduleWeekly4 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly4,
		EndRunTime:   endRunTimeWeekly4,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly4, DaysOfWeek: []int{0, 1}},
	}

	runTimeWeeklyExpected5, _ := time.Parse(time.RFC3339, "2039-12-25T14:30:00Z")
	startRunTimeWeekly5, _ := time.Parse(time.RFC3339, "2039-12-22T16:30:00Z")
	endRunTimeWeekly5, _ := time.Parse(time.RFC3339, "2040-10-13T20:30:00Z")
	runTimeWeekly5, _ := time.Parse(time.RFC3339, "2039-10-06T14:30:00Z")
	scheduleWeekly5 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly5,
		EndRunTime:   endRunTimeWeekly5,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly5, DaysOfWeek: []int{0, 1}},
	}

	runTimeWeeklyExpected6, _ := time.Parse(time.RFC3339, "2040-01-03T14:30:00Z")
	startRunTimeWeekly6, _ := time.Parse(time.RFC3339, "2039-12-30T16:30:00Z")
	endRunTimeWeekly6, _ := time.Parse(time.RFC3339, "2040-10-13T20:30:00Z")
	runTimeWeekly6, _ := time.Parse(time.RFC3339, "2039-10-06T14:30:00Z")
	scheduleWeekly6 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly6,
		EndRunTime:   endRunTimeWeekly6,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly6, DaysOfWeek: []int{5, 2}},
	}

	runTimeWeeklyExpected7 := time.Time{}
	startRunTimeWeekly7, _ := time.Parse(time.RFC3339, "2039-12-21T16:30:00Z")
	endRunTimeWeekly7, _ := time.Parse(time.RFC3339, "2039-12-22T20:30:00Z")
	runTimeWeekly7, _ := time.Parse(time.RFC3339, "2039-10-06T14:30:00Z")
	scheduleWeekly7 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly7,
		EndRunTime:   endRunTimeWeekly7,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly7, DaysOfWeek: []int{2, 5}},
	}

	runTimeWeeklyExpected8, _ := time.Parse(time.RFC3339, "2091-01-02T16:30:00Z")
	startRunTimeWeekly8, _ := time.Parse(time.RFC3339, "2091-01-02T16:30:00Z")
	endRunTimeWeekly8, _ := time.Parse(time.RFC3339, "2091-10-13T20:30:00Z")
	runTimeWeekly8, _ := time.Parse(time.RFC3339, "2019-10-06T16:30:00Z")
	scheduleWeekly8 := apiModels.Schedule{
		StartRunTime: startRunTimeWeekly8,
		EndRunTime:   endRunTimeWeekly8,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Weekly, RunTime: runTimeWeekly8, DaysOfWeek: []int{5, 2}},
	}

	// =========================== MONTHLY =========================================

	runTimeMonthlyExpected1, _ := time.Parse(time.RFC3339, "2090-10-10T17:30:00Z")
	startRunTimeMonthly1, _ := time.Parse(time.RFC3339, "2090-10-06T16:30:00Z")
	endRunTimeMonthly1, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly1, _ := time.Parse(time.RFC3339, "2090-10-06T17:30:00Z")
	scheduleMonthly1 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly1,
		EndRunTime:   endRunTimeMonthly1,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly1, DaysOfMonth: []int{10, 15, 30}},
	}

	runTimeMonthlyExpected2, _ := time.Parse(time.RFC3339, "2090-10-15T17:30:00Z")
	startRunTimeMonthly2, _ := time.Parse(time.RFC3339, "2090-10-11T16:30:00Z")
	endRunTimeMonthly2, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly2, _ := time.Parse(time.RFC3339, "2090-10-06T17:30:00Z")
	scheduleMonthly2 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly2,
		EndRunTime:   endRunTimeMonthly2,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly2, DaysOfMonth: []int{10, 15, 30}},
	}

	runTimeMonthlyExpected3, _ := time.Parse(time.RFC3339, "2090-10-30T17:30:00Z")
	startRunTimeMonthly3, _ := time.Parse(time.RFC3339, "2090-10-16T16:30:00Z")
	endRunTimeMonthly3, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly3, _ := time.Parse(time.RFC3339, "2090-10-06T17:30:00Z")
	scheduleMonthly3 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly3,
		EndRunTime:   endRunTimeMonthly3,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly3, DaysOfMonth: []int{10, 15, 30}},
	}

	runTimeMonthlyExpected4, _ := time.Parse(time.RFC3339, "2090-10-15T15:30:00Z")
	startRunTimeMonthly4, _ := time.Parse(time.RFC3339, "2090-10-10T16:30:00Z")
	endRunTimeMonthly4, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly4, _ := time.Parse(time.RFC3339, "2090-10-06T15:30:00Z")
	scheduleMonthly4 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly4,
		EndRunTime:   endRunTimeMonthly4,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly4, DaysOfMonth: []int{10, 15, 30}},
	}

	runTimeMonthlyExpected5, _ := time.Parse(time.RFC3339, "2090-02-28T15:30:00Z")
	startRunTimeMonthly5, _ := time.Parse(time.RFC3339, "2090-02-20T16:30:00Z")
	endRunTimeMonthly5, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly5, _ := time.Parse(time.RFC3339, "2090-10-06T15:30:00Z")
	scheduleMonthly5 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly5,
		EndRunTime:   endRunTimeMonthly5,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly5, DaysOfMonth: []int{10, 15, 30}},
	}

	runTimeMonthlyExpected6, _ := time.Parse(time.RFC3339, "2092-02-29T15:30:00Z")
	startRunTimeMonthly6, _ := time.Parse(time.RFC3339, "2092-02-20T16:30:00Z")
	endRunTimeMonthly6, _ := time.Parse(time.RFC3339, "2092-12-13T20:30:00Z")
	runTimeMonthly6, _ := time.Parse(time.RFC3339, "2090-10-06T15:30:00Z")
	scheduleMonthly6 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly6,
		EndRunTime:   endRunTimeMonthly6,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly6, DaysOfMonth: []int{10, 15, 30}},
	}

	runTimeMonthlyExpected7, _ := time.Parse(time.RFC3339, "2090-03-10T15:30:00Z")
	startRunTimeMonthly7, _ := time.Parse(time.RFC3339, "2090-02-21T16:30:00Z")
	endRunTimeMonthly7, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly7, _ := time.Parse(time.RFC3339, "2090-10-06T15:30:00Z")
	scheduleMonthly7 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly7,
		EndRunTime:   endRunTimeMonthly7,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly7, DaysOfMonth: []int{10, 15, 20}},
	}

	runTimeMonthlyExpected8, _ := time.Parse(time.RFC3339, "2090-03-31T15:30:00Z")
	startRunTimeMonthly8, _ := time.Parse(time.RFC3339, "2090-02-28T16:30:00Z")
	endRunTimeMonthly8, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly8, _ := time.Parse(time.RFC3339, "2090-10-06T15:30:00Z")
	scheduleMonthly8 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly8,
		EndRunTime:   endRunTimeMonthly8,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly8, DaysOfMonth: []int{31}},
	}

	runTimeMonthlyExpected9, _ := time.Parse(time.RFC3339, "2090-02-28T15:30:00Z")
	startRunTimeMonthly9, _ := time.Parse(time.RFC3339, "2090-02-28T14:30:00Z")
	endRunTimeMonthly9, _ := time.Parse(time.RFC3339, "2090-12-13T20:30:00Z")
	runTimeMonthly9, _ := time.Parse(time.RFC3339, "2090-10-06T15:30:00Z")
	scheduleMonthly9 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly9,
		EndRunTime:   endRunTimeMonthly9,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly9, DaysOfMonth: []int{31}},
	}

	runTimeMonthlyExpected10, _ := time.Parse(time.RFC3339, "2040-01-20T22:21:00Z")
	startRunTimeMonthly10, _ := time.Parse(time.RFC3339, "2039-12-27T22:21:00Z")
	endRunTimeMonthly10, _ := time.Parse(time.RFC3339, "2040-02-27T16:21:00Z")
	runTimeMonthly10, _ := time.Parse(time.RFC3339, "2039-11-27T22:21:00Z")
	scheduleMonthly10 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly10,
		EndRunTime:   endRunTimeMonthly10,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly10, WeekDays: []apiModels.WeekDay{{Day: 5, Index: 2}}},
	}

	runTimeMonthlyExpected11 := time.Time{}
	startRunTimeMonthly11, _ := time.Parse(time.RFC3339, "2039-12-27T22:21:00Z")
	endRunTimeMonthly11, _ := time.Parse(time.RFC3339, "2040-01-15T16:21:00Z")
	runTimeMonthly11, _ := time.Parse(time.RFC3339, "2039-11-27T22:21:00Z")
	scheduleMonthly11 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly11,
		EndRunTime:   endRunTimeMonthly11,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly11, WeekDays: []apiModels.WeekDay{{Day: 5, Index: 2}}},
	}

	runTimeMonthlyExpected12, _ := time.Parse(time.RFC3339, "2039-11-08T22:21:00Z")
	startRunTimeMonthly12, _ := time.Parse(time.RFC3339, "2039-11-07T22:21:00Z")
	endRunTimeMonthly12, _ := time.Parse(time.RFC3339, "2040-01-15T16:21:00Z")
	runTimeMonthly12, _ := time.Parse(time.RFC3339, "2039-11-27T22:21:00Z")
	scheduleMonthly12 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly12,
		EndRunTime:   endRunTimeMonthly12,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly12, WeekDays: []apiModels.WeekDay{{Day: 2, Index: 1}}},
	}

	runTimeMonthlyExpected13, _ := time.Parse(time.RFC3339, "2039-12-13T22:21:00Z")
	startRunTimeMonthly13, _ := time.Parse(time.RFC3339, "2039-11-09T22:21:00Z")
	endRunTimeMonthly13, _ := time.Parse(time.RFC3339, "2040-01-15T16:21:00Z")
	runTimeMonthly13, _ := time.Parse(time.RFC3339, "2039-11-27T22:21:00Z")
	scheduleMonthly13 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly13,
		EndRunTime:   endRunTimeMonthly13,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly13, WeekDays: []apiModels.WeekDay{{Day: 2, Index: 1}}},
	}

	runTimeMonthlyExpected14, _ := time.Parse(time.RFC3339, "2039-11-11T22:21:00Z")
	startRunTimeMonthly14, _ := time.Parse(time.RFC3339, "2039-11-09T22:21:00Z")
	endRunTimeMonthly14, _ := time.Parse(time.RFC3339, "2040-11-15T16:21:00Z")
	runTimeMonthly14, _ := time.Parse(time.RFC3339, "2039-11-27T22:21:00Z")
	scheduleMonthly14 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly14,
		EndRunTime:   endRunTimeMonthly14,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly14, WeekDays: []apiModels.WeekDay{{Day:4, Index:2}, {Day: 5, Index: 1}}},
	}

	runTimeMonthlyExpected15, _ := time.Parse(time.RFC3339, "2039-12-09T22:21:00Z")
	startRunTimeMonthly15, _ := time.Parse(time.RFC3339, "2039-11-29T22:21:00Z")
	endRunTimeMonthly15, _ := time.Parse(time.RFC3339, "2040-11-15T16:21:00Z")
	runTimeMonthly15, _ := time.Parse(time.RFC3339, "2039-11-27T22:21:00Z")
	scheduleMonthly15 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly15,
		EndRunTime:   endRunTimeMonthly15,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly15, WeekDays: []apiModels.WeekDay{{Day:4, Index:2}, {Day: 5, Index: 1}}},
	}

	runTimeMonthlyExpected16, _ := time.Parse(time.RFC3339, "2039-11-17T22:21:00Z")
	startRunTimeMonthly16, _ := time.Parse(time.RFC3339, "2039-11-16T22:21:00Z")
	endRunTimeMonthly16, _ := time.Parse(time.RFC3339, "2040-11-15T16:21:00Z")
	runTimeMonthly16, _ := time.Parse(time.RFC3339, "2039-11-27T22:21:00Z")
	scheduleMonthly16 := apiModels.Schedule{
		StartRunTime: startRunTimeMonthly16,
		EndRunTime:   endRunTimeMonthly16,
		Regularity:   apiModels.Recurrent,
		Repeat:       apiModels.Repeat{Every: 3, Frequency: apiModels.Monthly, RunTime: runTimeMonthly16, WeekDays: []apiModels.WeekDay{{Day:4, Index:2}, {Day: 5, Index: 1}}},
	}

	var tests = []struct {
		testName        string
		schedule        apiModels.Schedule
		expectedRunTime time.Time
		expectedErr     string
	}{
		// ==================== HOURLY ==============================
		{
			testName:        "hourly 1",
			schedule:        schedule1,
			expectedRunTime: runTimeHourlyExpected1,
		},
		{
			testName:        "hourly 2",
			schedule:        schedule2,
			expectedRunTime: runTimeHourlyExpected2,
			expectedErr:     taskWillNeverRunError,
		},
		{
			testName:        "hourly 3",
			schedule:        schedule3,
			expectedRunTime: runTimeHourlyExpected3,
		},

		// ==================== DAILY ==============================
		{
			testName:        "daily 1",
			schedule:        scheduleDaily1,
			expectedRunTime: runTimeDailyExpected1,
		},
		{
			testName:        "daily 2",
			schedule:        scheduleDaily2,
			expectedRunTime: runTimeDailyExpected2,
			expectedErr:     taskWillNeverRunError,
		},
		{
			testName:        "daily 3",
			schedule:        scheduleDaily3,
			expectedRunTime: runTimeDailyExpected3,
			expectedErr:     taskWillNeverRunError,
		},
		{
			testName:        "daily 4",
			schedule:        scheduleDaily4,
			expectedRunTime: runTimeDailyExpected4,
			expectedErr:     taskWillNeverRunError,
		},
		{
			testName:        "daily 5",
			schedule:        scheduleDaily5,
			expectedRunTime: runTimeDailyExpected5,
			expectedErr:     taskWillNeverRunError,
		},
		{
			testName:        "daily 6",
			schedule:        scheduleDaily6,
			expectedRunTime: runTimeDailyExpected6,
		},

		// ==================== WEEKLY ==============================
		{
			testName:        "weekly 1",
			schedule:        scheduleWeekly1,
			expectedRunTime: runTimeWeeklyExpected1,
		},
		{
			testName:        "weekly 2",
			schedule:        scheduleWeekly2,
			expectedRunTime: runTimeWeeklyExpected2,
		},
		{
			testName:        "weekly 3",
			schedule:        scheduleWeekly3,
			expectedRunTime: runTimeWeeklyExpected3,
		},
		{
			testName:        "weekly 4",
			schedule:        scheduleWeekly4,
			expectedRunTime: runTimeWeeklyExpected4,
		},
		{
			testName:        "weekly 5",
			schedule:        scheduleWeekly5,
			expectedRunTime: runTimeWeeklyExpected5,
		},
		{
			testName:        "weekly 6",
			schedule:        scheduleWeekly6,
			expectedRunTime: runTimeWeeklyExpected6,
		},
		{
			testName:        "weekly 7",
			schedule:        scheduleWeekly7,
			expectedRunTime: runTimeWeeklyExpected7,
			expectedErr:     taskWillNeverRunError,
		},
		{
			testName:        "weekly 8",
			schedule:        scheduleWeekly8,
			expectedRunTime: runTimeWeeklyExpected8,
		},

		// ==================== MONTHLY ==============================
		{
			testName:        "monthly 1",
			schedule:        scheduleMonthly1,
			expectedRunTime: runTimeMonthlyExpected1,
		},
		{
			testName:        "monthly 2",
			schedule:        scheduleMonthly2,
			expectedRunTime: runTimeMonthlyExpected2,
		},
		{
			testName:        "monthly 3",
			schedule:        scheduleMonthly3,
			expectedRunTime: runTimeMonthlyExpected3,
		},
		{
			testName:        "monthly 4",
			schedule:        scheduleMonthly4,
			expectedRunTime: runTimeMonthlyExpected4,
		},
		{
			testName:        "monthly 5",
			schedule:        scheduleMonthly5,
			expectedRunTime: runTimeMonthlyExpected5,
		},
		{
			testName:        "monthly 6",
			schedule:        scheduleMonthly6,
			expectedRunTime: runTimeMonthlyExpected6,
		},
		{
			testName:        "monthly 7",
			schedule:        scheduleMonthly7,
			expectedRunTime: runTimeMonthlyExpected7,
		},
		{
			testName:        "monthly 8",
			schedule:        scheduleMonthly8,
			expectedRunTime: runTimeMonthlyExpected8,
		},
		{
			testName:        "monthly 9",
			schedule:        scheduleMonthly9,
			expectedRunTime: runTimeMonthlyExpected9,
		},
		{
			testName:        "monthly 10",
			schedule:        scheduleMonthly10,
			expectedRunTime: runTimeMonthlyExpected10,
		},
		{
			testName:        "monthly 11",
			schedule:        scheduleMonthly11,
			expectedRunTime: runTimeMonthlyExpected11,
			expectedErr:     taskWillNeverRunError,
		},
		{
			testName:        "monthly 12",
			schedule:        scheduleMonthly12,
			expectedRunTime: runTimeMonthlyExpected12,
		},
		{
			testName:        "monthly 13",
			schedule:        scheduleMonthly13,
			expectedRunTime: runTimeMonthlyExpected13,
		},
		{
			testName:        "monthly 14",
			schedule:        scheduleMonthly14,
			expectedRunTime: runTimeMonthlyExpected14,
		},
		{
			testName:        "monthly 15",
			schedule:        scheduleMonthly15,
			expectedRunTime: runTimeMonthlyExpected15,
		},
		{
			testName:        "monthly 16",
			schedule:        scheduleMonthly16,
			expectedRunTime: runTimeMonthlyExpected16,
		},
	}

	for _, test := range tests {
		actualRunTime, err := CalcFirstNextRunTime(now.UTC(), test.schedule)

		if err != nil {
			Ω(err.Error()).To(Equal(test.expectedErr), fmt.Sprintf(defaultMsg, test.testName))
			Ω(actualRunTime).To(Equal(time.Time{}), fmt.Sprintf(defaultMsg, test.testName))
			continue
		}

		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.testName))
		Ω(actualRunTime).To(Equal(test.expectedRunTime), fmt.Sprintf(defaultMsg, test.testName))
	}
}

func TestConvertUUIDsToInterfaces(t *testing.T) {
	var (
		uuid1 = gocql.TimeUUID()
		uuid2 = gocql.TimeUUID()
		tests = []struct {
			name string
			in   []gocql.UUID
			out  []interface{}
		}{
			{
				name: "good case",
				in:   []gocql.UUID{uuid1, uuid2},
				out:  []interface{}{uuid1, uuid2},
			},
		}
	)

	for _, test := range tests {
		actualOut := ConvertUUIDsToInterfaces(test.in)
		if !reflect.DeepEqual(actualOut, test.out) {
			t.Fatalf("Test '%s': expected %v but got %v", test.name, test.out, actualOut)
		}
	}
}

func TestConvertStringsToUUIDs(t *testing.T) {
	var (
		uuid1       = gocql.TimeUUID()
		goodUUIDStr = uuid1.String()
		badUUIDStr  = "badUUIDStr"
		tests       = []struct {
			name          string
			in            []string
			out           []gocql.UUID
			expectedError bool
		}{
			{
				name:          "good case",
				in:            []string{goodUUIDStr, goodUUIDStr, goodUUIDStr},
				out:           []gocql.UUID{uuid1, uuid1, uuid1},
				expectedError: false,
			},
			{
				name:          "with error",
				in:            []string{goodUUIDStr, goodUUIDStr, badUUIDStr},
				out:           nil,
				expectedError: true,
			},
		}
	)

	for _, test := range tests {
		actualOut, err := ConvertStringsToUUIDs(test.in)
		if !reflect.DeepEqual(actualOut, test.out) {
			t.Fatalf("Test %s expected %v but got %v", test.name, test.out, actualOut)
		}
		if (err != nil) != test.expectedError {
			t.Fatalf("Test '%s': expectedError %t but got err %v", test.name, test.expectedError, err)
		}
	}
}

func TestUUIDSliceContainsElement(t *testing.T) {
	var (
		uuid1 = gocql.TimeUUID()
		uuid2 = gocql.TimeUUID()
		tests = []struct {
			name            string
			uniqueUUIDSlice []gocql.UUID
			element         gocql.UUID
			expectedResult  bool
		}{
			{
				name:            "slice contains element",
				uniqueUUIDSlice: []gocql.UUID{uuid1, uuid2},
				element:         uuid1,
				expectedResult:  true,
			},
			{
				name:            "slice doesn't contain element",
				uniqueUUIDSlice: []gocql.UUID{uuid1, uuid2},
				element:         gocql.TimeUUID(),
				expectedResult:  false,
			},
		}
	)

	for _, test := range tests {
		actualResult := UUIDSliceContainsElement(test.uniqueUUIDSlice, test.element)
		if actualResult != test.expectedResult {
			t.Fatalf("Test '%s': expected %t but got %t", test.name, test.expectedResult, actualResult)
		}
	}
}

func TestAddLocationToTime(t *testing.T) {
	kievLocation, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		t.Fatal("Error while loading location")
	}

	UTCLocation, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatal("Error while loading location")
	}

	testCases := []struct {
		name           string
		timeToSend     time.Time
		locationToSend *time.Location
		expectedTime   time.Time
	}{
		{
			name:         "testCase 1 - Location is nil",
			expectedTime: time.Time{},
		},
		{
			name:           "testCase 2 - Ok",
			locationToSend: kievLocation,
			timeToSend:     time.Date(2016, 12, 1, 1, 1, 1, 1, UTCLocation),
			expectedTime:   time.Date(2016, 12, 1, 1, 1, 1, 1, kievLocation),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotTime := AddLocationToTime(tc.timeToSend, tc.locationToSend)
			if gotTime != tc.expectedTime {
				t.Errorf("Want %v but got %v", tc.expectedTime, gotTime)
			}
		})
	}
}
