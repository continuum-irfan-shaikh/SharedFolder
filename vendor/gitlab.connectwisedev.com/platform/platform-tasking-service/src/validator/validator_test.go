package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ContinuumLLC/godep-govalidator"
	"github.com/gocql/gocql"
	. "github.com/onsi/gomega"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/model-mocks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
)

var (
	testTime       = time.Now().AddDate(0, 0, 2).UTC()
	runTime        = testTime.Truncate(time.Minute)
	testUUID       = gocql.TimeUUID()
	targetEndpoint = models.Target{
		IDs:  []string{testUUID.String()},
		Type: models.ManagedEndpoint,
	}
	targetDG = models.Target{
		IDs:  []string{testUUID.String()},
		Type: models.DynamicGroup,
	}
	badTargetDG = models.Target{
		IDs:  []string{testUUID.String(), "", testUUID.String()},
		Type: models.DynamicGroup,
	}
	schedulePattern = regexp.MustCompile(`^\s*@every\s*[1-9][0-9]*[h,d,w,m]$`)
)

const defaultMsg = `failed on unexpected value of result "%v"`

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("can't read request body")
}

func init() {
	config.Load()
	logger.Load(config.Config.Log)
}

func getTaskTemplateManagedEndpointOk() models.Task {
	return models.Task{
		Targets:     targetEndpoint,
		OriginID:    testUUID,
		Type:        config.ScriptTaskType,
		Schedule:    apiModels.Schedule{},
		Credentials: &agentModels.Credentials{true, "", "", ""},
	}
}

func getTaskTemplateDynamicGroupsOk() models.Task {
	return models.Task{
		Targets:     targetDG,
		OriginID:    testUUID,
		Type:        config.ScriptTaskType,
		Schedule:    apiModels.Schedule{Location: "UTC"},
		Credentials: &agentModels.Credentials{false, "username", "domain", "password"},
	}
}

func getTaskTemplateRunNowOk() models.Task {
	task := getTaskTemplateManagedEndpointOk()
	task.Schedule.Regularity = apiModels.RunNow
	task.Schedule.EndRunTime = time.Now().AddDate(100, 0, 0)
	task.Credentials = &agentModels.Credentials{false, "username", "", "password"}
	return task
}

func getTaskTemplateOneTimeOk() models.Task {
	task := getTaskTemplateDynamicGroupsOk()
	task.Schedule.Regularity = apiModels.OneTime
	task.Schedule.StartRunTime = runTime
	return task
}

func getTaskTemplateRecurrentlyOk() models.Task {
	task := getTaskTemplateManagedEndpointOk()
	task.Schedule = apiModels.Schedule{Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Every: 5, Frequency: apiModels.Daily, RunTime: testTime}}
	task.Schedule.StartRunTime = testTime
	task.Schedule.EndRunTime = testTime.Add(time.Hour)
	return task
}

func loadTasksOk() []models.Task {
	return []models.Task{
		getTaskTemplateRunNowOk(),
		getTaskTemplateOneTimeOk(),
		getTaskTemplateRecurrentlyOk(),
	}
}

func getTaskTemplateInvalid() models.Task {
	return models.Task{
		Type:       "script",
		OriginID:   testUUID,
		Schedule:   apiModels.Schedule{Regularity: apiModels.RunNow, EndRunTime: testTime, StartRunTime: testTime},
		Parameters: "",
	}
}

func getTaskTemplateCredsInvalid() models.Task {
	task := getTaskTemplateInvalid()
	task.Credentials = &agentModels.Credentials{true, "username", "domain", "password"}
	return task
}

func getTaskTemplateCredsInvalid1() models.Task {
	task := getTaskTemplateInvalid()
	task.Credentials = &agentModels.Credentials{true, "username", "", "password"}
	return task
}

func getTaskTemplateCredsInvalid2() models.Task {
	task := getTaskTemplateInvalid()
	task.Credentials = &agentModels.Credentials{false, "username", "domain", ""}
	return task
}

func getTaskTemplateCredsInvalid3() models.Task {
	task := getTaskTemplateInvalid()
	task.Credentials = &agentModels.Credentials{false, "", "domain", "password"}
	return task
}

func getTaskTemplateRunNowInvalid() models.Task {
	task := getTaskTemplateInvalid()
	task.RunTimeUTC = runTime
	task.Schedule = apiModels.Schedule{}
	task.Schedule.StartRunTime = testTime
	task.Schedule.EndRunTime = testTime.Add(time.Hour)
	task.Schedule.Regularity = apiModels.RunNow
	task.Targets = badTargetDG
	return task
}

func getTaskTemplateOneTimeInvalid() models.Task {
	task := getTaskTemplateInvalid()
	task.Schedule = apiModels.Schedule{}
	task.Schedule.Regularity = apiModels.OneTime
	task.Schedule.StartRunTime = testTime
	task.Schedule.EndRunTime = testTime.Add(time.Hour)
	task.Targets = badTargetDG
	return task
}

func getTaskTemplateRecurrentlyInvalid() models.Task {
	task := getTaskTemplateInvalid()
	task.Schedule.Regularity = apiModels.Recurrent
	task.RunTimeUTC = runTime
	return task
}

func getTaskTemplateRecurrentlyInvalidData() models.Task {
	task := getTaskTemplateManagedEndpointOk()
	task.Schedule = apiModels.Schedule{}
	task.Schedule.Regularity = apiModels.Recurrent
	task.Schedule.StartRunTime = testTime
	task.Schedule.EndRunTime = testTime.Add(-1 * time.Hour)
	task.Parameters = "invalid JSON"
	task.Schedule.Location = "invalid location"
	return task
}

func loadTasksInvalid() []models.Task {
	return []models.Task{
		getTaskTemplateInvalid(),
		getTaskTemplateRunNowInvalid(),
		getTaskTemplateOneTimeInvalid(),
		getTaskTemplateRecurrentlyInvalid(),
		getTaskTemplateRecurrentlyInvalidData(),
		getTaskTemplateCredsInvalid(),
		getTaskTemplateCredsInvalid1(),
		getTaskTemplateCredsInvalid2(),
		getTaskTemplateCredsInvalid3(),
	}
}

func TestValidTask(t *testing.T) {
	tasks := loadTasksOk()
	SetupCustomValidators()
	for _, task := range tasks {
		if _, err := govalidator.ValidateStruct(task); err != nil {
			t.Fatalf("err: %v", err)
		}
	}
}

func TestInvalidTask(t *testing.T) {
	tasks := loadTasksInvalid()
	SetupCustomValidators()
	for _, task := range tasks {
		if _, err := govalidator.ValidateStruct(task); err == nil {
			t.Fatalf("Got no error while validating invalid %T: %[1]v", task)
		}
	}
}

func TestDefaultStruct(t *testing.T) {
	type Default struct {
		RequiredForUsers     string   `valid:"requiredForUsers"`
		UnsettableByUsers    string   `valid:"unsettableByUsers"`
		RequiredForOneTime   string   `valid:"requiredOnlyForOneTime"`
		RequiredForRecurrent string   `valid:"requiredOnlyForRecurrent"`
		OptionalForRecurrent string   `valid:"optionalOnlyForRecurrent"`
		TaskType             string   `valid:"validType"`
		Targets              []string `valid:"requiredUniqueTargetIDs"`
		Category             []string `valid:"validCategories"`
		Location             string   `valid:"validLocation"`
	}
	defaultStruct := Default{RequiredForUsers: "not empty string"}
	SetupCustomValidators()
	if _, err := govalidator.ValidateStruct(defaultStruct); err == nil {
		t.Fatalf("Got no error while validating invalid %T: %[1]v", defaultStruct)
	}
}

func TestValidSelectedTargetsEnable(t *testing.T) {
	selectedTargetsEnable := models.SelectedManagedEndpointEnable{
		ManagedEndpoints: modelMocks.GenerateTargetsMock(),
	}
	SetupCustomValidators()
	if _, err := govalidator.ValidateStruct(selectedTargetsEnable); err != nil {
		t.Fatalf("Error while validating %T: %[1]v, err: %v", selectedTargetsEnable, err)
	}
}

func TestInvalidSelectedTargetsEnable(t *testing.T) {
	var selectedTargetsEnable models.SelectedManagedEndpointEnable
	SetupCustomValidators()
	if _, err := govalidator.ValidateStruct(selectedTargetsEnable); err == nil {
		t.Fatalf("Got no error while validating invalid %T: %[1]v", selectedTargetsEnable)
	}
}

func TestValidateTaskDefinition(t *testing.T) {
	SetupCustomValidators()
	var tests = []struct {
		name      string
		input     models.TaskDefinitionDetails
		expectErr bool
	}{
		{
			name: "Test valid Task Definition",
			input: models.TaskDefinitionDetails{
				TaskDefinition: models.TaskDefinition{
					Name:        "Name",
					OriginID:    gocql.TimeUUID(),
					Type:        "script",
					Categories:  []string{"category1"},
					Description: "Description",
				},
				UserParameters: "{\"key1\": 233}",
			},
			expectErr: false,
		},
		{
			name: "Test invalid JSON in Task Definition",
			input: models.TaskDefinitionDetails{
				TaskDefinition: models.TaskDefinition{
					Name:        "Name",
					OriginID:    modelMocks.DefaultTaskDefs[0].OriginID,
					Type:        "script",
					Categories:  []string{"category1", "category1"},
					Description: "Description",
				},
				UserParameters: "{\"key1\": 233}",
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		if test.expectErr {
			if _, err := govalidator.ValidateStruct(test.input); err == nil {
				t.Fatalf("%s Expect error but haven't got it: %v", test.name, err)
			}
		} else {
			if _, err := govalidator.ValidateStruct(test.input); err != nil {
				t.Fatalf("%s error: %v", test.name, err)
			}
		}
	}
}

func TestExtractTask(t *testing.T) {
	makeRequest := func(task models.Task) *http.Request {
		js, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("makeRequest() error = %sv", err.Error())
		}
		return httptest.NewRequest(`GET`, `http://www.localhost.ua`, bytes.NewReader(js))
	}

	someFutureTime := time.Now().UTC().AddDate(1, 0, 0)

	tests := []struct {
		name     string
		r        *http.Request
		wantTask models.Task
		wantErr  bool
	}{
		{
			name: `bad Invalid request body`,
			r: httptest.NewRequest(`GET`, `http://www.localhost.ua`, bytes.NewReader(
				[]byte(`{
				"id":"00000000-0000-0000-0000-0000
				"description":"",
				"targets":null,
				"schedule":"",
				"created.....regdretgreg r
				dfgdf gdfgd
				 dafg a2354 0i&897*(&*^&*% &*`),
			)),
			wantTask: models.Task{},
			wantErr:  true,
		},
		{
			name: `good`,
			r: makeRequest(models.Task{
				Targets:  targetEndpoint,
				Schedule: apiModels.Schedule{Regularity: apiModels.Recurrent, StartRunTime: testTime, Repeat: apiModels.Repeat{Every: 5, Frequency: apiModels.Hourly, RunTime: time.Date(2018, 10, 10, 4, 5, 0, 0, time.UTC)}},
				OriginID: testUUID,
				Type:     "script",
			}),
			wantTask: models.Task{
				Targets:  targetEndpoint,
				Schedule: apiModels.Schedule{Regularity: apiModels.Recurrent, StartRunTime: testTime, Repeat: apiModels.Repeat{Every: 5, Frequency: apiModels.Hourly, RunTime: time.Date(2018, 10, 10, 4, 5, 0, 0, time.UTC)}},
				OriginID: testUUID,
				Type:     "script",
			},
			wantErr: false,
		},
		{
			name: `bad Invalid Task`,
			r: makeRequest(models.Task{
				CreatedBy: "test",
				Targets:   targetEndpoint,
				ID:        testUUID,
				Schedule:  apiModels.Schedule{Regularity: apiModels.Recurrent, StartRunTime: testTime, Repeat: apiModels.Repeat{Every: 5, Frequency: apiModels.Monthly, RunTime: time.Date(2018, 10, 10, 5, 4, 0, 0, time.UTC)}},
				OriginID:  testUUID,
				Type:      "script",
			}),
			wantTask: models.Task{
				CreatedBy: "test",
				Targets:   targetEndpoint,
				ID:        testUUID,
				Schedule:  apiModels.Schedule{Regularity: apiModels.Recurrent, StartRunTime: testTime, Repeat: apiModels.Repeat{Every: 5, Frequency: apiModels.Monthly, RunTime: time.Date(2018, 10, 10, 5, 4, 0, 0, time.UTC)}},
				OriginID:  testUUID,
				Type:      `script`,
			},
			wantErr: true,
		},
		{
			name: `location validator: negative test`,
			r: makeRequest(models.Task{
				Targets:  targetEndpoint,
				OriginID: testUUID,
				Schedule: apiModels.Schedule{Regularity: apiModels.RunNow, Location: "invalidLocation"},
				Type:     "script",
			}),
			wantTask: models.Task{
				Targets:  targetEndpoint,
				OriginID: testUUID,
				Schedule: apiModels.Schedule{Regularity: apiModels.RunNow, Location: "invalidLocation"},
				Type:     "script",
			},
			wantErr: true,
		},
		{
			name: `location validator: positive test`,
			r: makeRequest(models.Task{
				Targets:  targetEndpoint,
				OriginID: testUUID,
				Schedule: apiModels.Schedule{Regularity: apiModels.OneTime, Location: "Europe/Kiev", StartRunTime: someFutureTime},
				Type:     "script",
			}),
			wantTask: models.Task{
				Targets:  targetEndpoint,
				OriginID: testUUID,
				Schedule: apiModels.Schedule{Regularity: apiModels.OneTime, Location: "Europe/Kiev", StartRunTime: someFutureTime},
				Type:     "script",
			},
			wantErr: false,
		},
		{
			name:     `bad: read request body error`,
			r:        httptest.NewRequest(`GET`, `http://www.localhost.ua`, errReader(0)),
			wantTask: models.Task{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotTask models.Task
			err := ExtractStructFromRequest(tt.r, &gotTask)
			if (err != nil) != tt.wantErr {
				t.Fatalf("extractTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTask, tt.wantTask) {
				t.Fatalf("extractTask() = %v, want %v", gotTask, tt.wantTask)
			}
		})
	}
}

func TestValidCategories(t *testing.T) {
	SetupCustomValidators()
	testCases := []struct {
		name      string
		input     models.TaskDefinition
		expectErr bool
	}{
		{
			name: "TestCase1",
			input: models.TaskDefinition{
				OriginID: modelMocks.DefaultTemplates[0].OriginID,
				Name:     "Name1",
				Type:     "script",
				Categories: []string{
					"test1",
					"test2",
					"test3",
				},
			},
			expectErr: false,
		},
		{
			name: "TestCase2",
			input: models.TaskDefinition{
				OriginID: modelMocks.DefaultTemplates[0].OriginID,
				Name:     "Name1",
				Type:     "script",
				Categories: []string{
					"test1",
					"test2",
					"test2",
				},
			},
			expectErr: true,
		},
		{
			name: "TestCase3",
			input: models.TaskDefinition{
				OriginID: modelMocks.DefaultTemplates[0].OriginID,
				Name:     "Name1",
				Type:     "script",
				Categories: []string{
					"test1",
					"",
					"test3",
				},
			},
			expectErr: true,
		},
	}

	for _, testCase := range testCases {
		if _, err := govalidator.ValidateStruct(testCase.input); (err != nil) != testCase.expectErr {
			t.Fatalf("%s expect error but haven't got it: %v", testCase.name, err)
		}
	}
}

func TestExtractAllTargetsEnable(t *testing.T) {
	tests := []struct {
		name                 string
		r                    *http.Request
		wantAllTargetsEnable models.AllTargetsEnable
		valid                bool
	}{
		{
			name: `bad`,
			r: httptest.NewRequest(`GET`, `http://www.localhost.ua`, bytes.NewReader(
				[]byte(`{}`),
			)),
			wantAllTargetsEnable: models.AllTargetsEnable{},
			valid:                false,
		},
		{
			name: `good`,
			r: httptest.NewRequest(`GET`, `http://www.localhost.ua`, bytes.NewReader(
				[]byte(`{"active":true}`),
			)),
			wantAllTargetsEnable: models.AllTargetsEnable{
				Active: true,
			},
			valid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotAllTargetsEnable models.AllTargetsEnable
			err := ExtractStructFromRequest(tt.r, &gotAllTargetsEnable)
			if (err == nil) != tt.valid {
				t.Fatalf("extractAllTargetsEnable() error = %v, but validation of model should return %v", err, tt.valid)
				return
			}
			if !reflect.DeepEqual(gotAllTargetsEnable, tt.wantAllTargetsEnable) {
				t.Fatalf("extractAllTargetsEnable(): got %v, want %v", gotAllTargetsEnable, tt.wantAllTargetsEnable)
			}
		})
	}
}

func TestValidateParametersField(t *testing.T) {
	userParamsWithRedundantParam := "{\"body\":\"test\",\"unnecessaryParam\":true}"
	emptyJSONSchema := ""
	jsonSchemaWithoutRequiredParams := "{\"properties\":{\"body\":{\"title\":\"PowerShell Script\",\"type\":\"string\"}},\"additionalProperties\": false}"
	jsonSchemaWithRequiredParams := "{\"properties\":{\"body\":{\"title\":\"PowerShell Script\",\"type\":\"string\"}},\"required\":[\"body\"],\"additionalProperties\": false}"
	jsonSchemaWithNestedRequired := `{"type":"object","properties":{"action":{"type":"string","title":"Action","enum":["create","update","delete"],"enumNames":["Insert","Update","Delete"]},"powerPlan":{"type":"string","title":"Power Plan name"}},"required":["action","powerPlan"],"patternProperties":{"^basePowerPlan$":{},"^turnOffDisplayAC$":{},"^turnOffDisplayDC$":{},"^setActive$":{}},"additionalProperties":false,"dependencies":{"action":{"oneOf":[{"properties":{"action":{"enum":["create"]},"basePowerPlan":{"type":"string","title":"Base Power plan","enum":["Balanced","High Performance","Power saver"]},"turnOffDisplayAC":{"type":"integer","minimum":0,"title":"Turn off display (plugged in)"},"turnOffDisplayDC":{"type":"integer","minimum":0,"title":"Turn off display (on battery)"},"setActive":{"type":"boolean","title":"Set active"}},"required":["basePowerPlan","turnOffDisplayAC","turnOffDisplayDC"]},{"properties":{"action":{"enum":["update"]},"turnOffDisplayAC":{"type":"integer","minimum":0,"title":"Turn off display (plugged in)"},"turnOffDisplayDC":{"type":"integer","minimum":0,"title":"Turn off display (on battery)"},"setActive":{"type":"boolean","title":"Set active"}},"required":["turnOffDisplayAC","turnOffDisplayDC"]},{"properties":{"action":{"enum":["delete"]}}}]}}}`

	userParamsWithCustomScript := "{\"body\":\"$Path = \\\"C:\\Users\\$env:UserName\\AppData\\Local\\Microsoft\\Windows\"\"}"
	testCases := []struct {
		name          string
		jsonSchema    string
		parameters    string
		isExpectedErr bool
		expectedErr   string
		strict        bool
	}{
		{
			name:          "testCase1",
			jsonSchema:    modelMocks.ValidJSONSchema,
			parameters:    modelMocks.ValidUserParams,
			isExpectedErr: false,
			expectedErr:   "",
			strict:        true,
		},
		{
			name:          "testCase1.1",
			jsonSchema:    modelMocks.ValidJSONSchema,
			parameters:    modelMocks.ValidUserParams,
			isExpectedErr: false,
			expectedErr:   "",
			strict:        false,
		},
		{
			name:          "testCase2",
			jsonSchema:    modelMocks.ValidJSONSchema,
			parameters:    modelMocks.InvalidUserParams,
			isExpectedErr: true,
			expectedErr:   "parameters validation error: invalid character 'e' in literal true (expecting 'r')",
			strict:        true,
		},
		{
			name:          "testCase3",
			jsonSchema:    modelMocks.ValidJSONSchema,
			parameters:    modelMocks.InvalidTypeUserParams,
			isExpectedErr: true,
			expectedErr:   "parameters are not valid, err: [body: Invalid type. Expected: string, given: boolean]",
			strict:        true,
		},
		{
			name:          "testCase3.1",
			jsonSchema:    modelMocks.ValidJSONSchema,
			parameters:    modelMocks.InvalidTypeUserParams,
			isExpectedErr: true,
			expectedErr:   "parameters are not valid, err: [body: Invalid type. Expected: string, given: boolean]",
			strict:        false,
		},
		{
			name:          "testCase4",
			jsonSchema:    modelMocks.ValidJSONSchema,
			parameters:    userParamsWithRedundantParam,
			isExpectedErr: true,
			expectedErr:   "parameters are not valid, err: [unnecessaryParam: Additional property unnecessaryParam is not allowed]",
			strict:        true,
		},
		{
			name:          "testCase5",
			jsonSchema:    emptyJSONSchema,
			parameters:    modelMocks.ValidUserParams,
			isExpectedErr: true,
			expectedErr:   `parameters for non-parameterized script should be empty, but got: {"body":"test"}`,
			strict:        true,
		},
		{
			name:          "testCase6",
			jsonSchema:    emptyJSONSchema,
			parameters:    emptyJSONSchema,
			isExpectedErr: false,
			expectedErr:   "",
			strict:        true,
		},
		{
			name:          "testCase7",
			jsonSchema:    modelMocks.ValidJSONSchema,
			parameters:    userParamsWithCustomScript,
			isExpectedErr: false,
			expectedErr:   "",
			strict:        true,
		},
		{
			name:          "testCase8",
			jsonSchema:    jsonSchemaWithoutRequiredParams,
			parameters:    "",
			isExpectedErr: false,
			expectedErr:   "",
			strict:        true,
		},
		{
			name:          "testCase9",
			jsonSchema:    jsonSchemaWithRequiredParams,
			parameters:    "",
			isExpectedErr: false,
			expectedErr:   "",
			strict:        false,
		},
		{
			name:          "testCase10",
			jsonSchema:    jsonSchemaWithNestedRequired,
			parameters:    "",
			isExpectedErr: false,
			expectedErr:   "",
			strict:        false,
		},
	}

	for _, testCase := range testCases {
		err := ValidateParametersField(testCase.jsonSchema, testCase.parameters, testCase.strict)
		if testCase.isExpectedErr && err.Error() != testCase.expectedErr {
			t.Fatalf("%s with expected %v, got %v", testCase.name, testCase.expectedErr, err.Error())
		}
		if !testCase.isExpectedErr && testCase.expectedErr != "" {
			t.Fatalf("%s with expected %v, got %v", testCase.name, testCase.expectedErr, err)
		}
	}
}

func TestValidInterval(t *testing.T) {
	tests := []struct {
		testNumber string
		interval   string
		result     bool
	}{
		{
			testNumber: "1",
			interval:   "@every 5m",
			result:     true,
		},
		{
			testNumber: "2",
			interval:   "@every 122d",
			result:     true,
		},
		{
			testNumber: "3",
			interval:   "@every 0m",
			result:     false,
		},
		{
			testNumber: "4",
			interval:   "@every 10r",
			result:     false,
		},
		{
			testNumber: "5",
			interval:   "@every 10mr",
			result:     false,
		},
	}
	for _, test := range tests {
		actualRes := validInterval(test.interval)
		if test.result != actualRes {
			t.Fatalf("Test case %v failed, Expected: %v, but got %v", test.testNumber, test.result, actualRes)
		}
	}
}

func TestValidatorTargets(t *testing.T) {
	var res bool
	testCases := []struct {
		name   string
		target models.Target
		result bool
	}{
		{
			name:   "testCase1 - managedEndpoint target, valid id ",
			target: targetEndpoint,
			result: true,
		},
		{
			name: "testCase2 - managedEndpoint target, invalid id ",
			target: models.Target{
				IDs:  []string{"50040040", "50040020"},
				Type: models.ManagedEndpoint,
			},
			result: false,
		},
		{
			name:   "testCase3 - dynamicGroup target, valid id ",
			target: targetDG,
			result: true,
		},
		{
			name:   "testCase4 - dynamicGroup target, invalid id ",
			target: badTargetDG,
			result: false,
		},
		{
			name: "testCase5 - site target, valid id ",
			target: models.Target{
				IDs:  []string{"50040040", "50040020"},
				Type: models.Site,
			},
			result: true,
		},
		{
			name: "testCase6 - site target, invalid id ",
			target: models.Target{
				IDs:  []string{"50040040", "50040040"},
				Type: models.Site,
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		res = validatorTargets(tc.target.IDs, tc.target)
		if res != tc.result {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.result, res)
		}
	}
}

func validInterval(interval string) bool {
	return schedulePattern.MatchString(strings.Trim(interval, " "))
}

func TestValidatorRequiredWeekly(t *testing.T) {

	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract field, i is nil interface",
			expectBool: false,
			ctx:        apiModels.Repeat{},
		},
		{
			name:       "not weekly",
			expectBool: false,
			ctx:        apiModels.Repeat{Frequency: apiModels.Daily},
			i:          []int{1},
		},
		{
			name:       "not weekly",
			expectBool: true,
			ctx:        apiModels.Repeat{Frequency: apiModels.Weekly},
			i:          []int{1},
		},
		{
			name:       "empty weekdays",
			expectBool: false,
			ctx:        apiModels.Repeat{Frequency: apiModels.Weekly},
			i:          []int{},
		},
		{
			name:       "not valid weekday",
			expectBool: false,
			ctx:        apiModels.Repeat{Frequency: apiModels.Weekly},
			i:          []int{0, 1, 5, 7},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validatorRequiredWeekly(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorResourceType(t *testing.T) {
	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract i, is not resType",
			expectBool: false,
			i:          models.Task{},
			ctx:        models.Task{},
		},
		{
			name:       "res type all",
			expectBool: true,
			i:          integration.ResourceType(""),
			ctx:        models.Task{},
		},
		{
			name:       "res Desktop",
			expectBool: true,
			i:          integration.Desktop,
			ctx:        models.Task{},
		},
		{
			name:       "res invalid",
			expectBool: false,
			i:          integration.ResourceType("invalid"),
			ctx:        models.Task{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validResourceType(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorOptionalTriggerTypes(t *testing.T) {
	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract i, is not schedule",
			expectBool: false,
			i:          models.Task{},
			ctx:        models.Task{},
		},
		{
			name:       "new target false, dg",
			expectBool: false,
			i: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				TriggerTypes: []string{triggers.FirstCheckInTrigger},
			},
			ctx: models.Task{
				TargetsByType: map[models.TargetType][]string{
					models.DynamicGroup: {"id"},
				},
			},
		},
		{
			name:       "new target true",
			expectBool: true,
			i: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				TriggerTypes: []string{triggers.FirstCheckInTrigger},
			},
			ctx: models.Task{
				TargetsByType: map[models.TargetType][]string{
					models.Site: {"id"},
				},
			},
		},
		{
			name:       "new target true",
			expectBool: true,
			i: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				TriggerTypes: []string{triggers.FirstCheckInTrigger},
			},
			ctx: models.Task{
				TargetsByType: map[models.TargetType][]string{
					models.DynamicSite: {"id"},
				},
			},
		},
		{
			name:       "old target true",
			expectBool: true,
			i: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				TriggerTypes: []string{triggers.FirstCheckInTrigger},
			},
			ctx: models.Task{
				Targets: models.Target{
					Type: models.Site,
				},
			},
		},
		{
			name:       "old target true",
			expectBool: true,
			i: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				TriggerTypes: []string{triggers.FirstCheckInTrigger},
			},
			ctx: models.Task{
				Targets: models.Target{
					Type: models.DynamicSite,
				},
			},
		},
		{
			name:       "old target false",
			expectBool: false,
			i: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				TriggerTypes: []string{triggers.FirstCheckInTrigger},
			},
			ctx: models.Task{
				Targets: models.Target{
					Type: models.DynamicGroup,
				},
			},
		},
		{
			name:       "old target true",
			expectBool: true,
			i: apiModels.Schedule{
				Regularity: apiModels.Recurrent,
			},
			ctx: models.Task{
				Targets: models.Target{
					Type: models.DynamicGroup,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validatorOptionalTriggerTypes(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorDynamicGroup(t *testing.T) {
	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract i, is not schedule",
			expectBool: false,
			i:          models.Task{},
			ctx:        models.Task{},
		},
		{
			name:       "new target true",
			expectBool: true,
			i: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				TriggerTypes: []string{triggers.DynamicGroupEnterTrigger},
			},
			ctx: models.Task{
				TargetsByType: map[models.TargetType][]string{
					models.DynamicGroup: {"id"},
				},
			},
		},
		{
			name:       "new target false , site",
			expectBool: false,
			i:          apiModels.Schedule{},
			ctx: models.Task{
				Schedule: apiModels.Schedule{
					Regularity:   apiModels.Trigger,
					TriggerTypes: []string{triggers.DynamicGroupEnterTrigger},
				},
				TargetsByType: map[models.TargetType][]string{
					models.Site: {"id"},
				},
			},
		},
		{
			name:       "new target true , site+dg",
			expectBool: true,
			i:          apiModels.Schedule{},
			ctx: models.Task{
				Schedule: apiModels.Schedule{
					Regularity:   apiModels.Trigger,
					TriggerTypes: []string{triggers.DynamicGroupEnterTrigger},
				},
				TargetsByType: map[models.TargetType][]string{
					models.Site:         {"id"},
					models.DynamicGroup: {"id"},
				},
			},
		},
		{
			name:       "new target true , not trigger",
			expectBool: true,
			i:          apiModels.Schedule{},
			ctx: models.Task{
				Schedule: apiModels.Schedule{
					Regularity: apiModels.Recurrent,
				},
				TargetsByType: map[models.TargetType][]string{
					models.Site:         {"id"},
					models.DynamicGroup: {"id"},
				},
			},
		},
		{
			name:       "new target false , site",
			expectBool: false,
			i:          apiModels.Schedule{},
			ctx: models.Task{
				Schedule: apiModels.Schedule{
					Regularity:   apiModels.Trigger,
					TriggerTypes: []string{triggers.DynamicGroupEnterTrigger},
				},
				TargetsByType: map[models.TargetType][]string{
					models.Site: {"id"},
				},
			},
		},
		{
			name:       "old target false , site",
			expectBool: false,
			i:          apiModels.Schedule{},
			ctx: models.Task{
				Schedule: apiModels.Schedule{
					Regularity:   apiModels.Trigger,
					TriggerTypes: []string{triggers.DynamicGroupEnterTrigger},
				},
				Targets: models.Target{
					Type: models.Site,
				},
			},
		},
		{
			name:       "old target true",
			expectBool: true,
			i:          apiModels.Schedule{},
			ctx: models.Task{
				Schedule: apiModels.Schedule{
					Regularity:   apiModels.Trigger,
					TriggerTypes: []string{triggers.DynamicGroupEnterTrigger},
				},
				Targets: models.Target{
					Type: models.DynamicGroup,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validatorDynamicGroup(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorTargetStruct(t *testing.T) {
	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract i, is not targets",
			expectBool: false,
			i:          models.Task{},
			ctx:        models.Task{},
		},
		{
			name:       "res type all",
			expectBool: false,
			i: models.Target{
				IDs:  []string{"notempty"},
				Type: models.DynamicGroup,
			},
			ctx: models.Task{TargetsByType: map[models.TargetType][]string{models.DynamicGroup: []string{"notempty"}}},
		},
		{
			name:       "ok",
			expectBool: true,
			i: models.Target{
				IDs:  []string{"notempty"},
				Type: models.DynamicGroup,
			},
			ctx: models.Task{},
		},
		{
			name:       "both empty",
			expectBool: false,
			i:          models.Target{},
			ctx:        models.Task{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validTargets(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorTargetMap(t *testing.T) {
	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract i, is not targets",
			expectBool: false,
			i:          models.Task{},
			ctx:        models.Task{},
		},
		{
			name:       "has only targets",
			expectBool: true,
			i: models.Target{
				IDs:  []string{"notempty"},
				Type: models.DynamicGroup,
			},
			ctx: models.Task{
				TargetsByType: make(models.TargetsByType),
			},
		},
		{
			name:       "has  both",
			expectBool: false,
			i: models.Target{
				IDs:  []string{"notempty"},
				Type: models.DynamicGroup,
			},
			ctx: models.Task{
				TargetsByType: map[models.TargetType][]string{models.DynamicGroup: []string{"notempty"}},
			},
		},
		{
			name:       "has nothing",
			expectBool: false,
			i:          models.Target{},
			ctx: models.Task{
				TargetsByType: map[models.TargetType][]string{},
			},
		},
		{
			name:       "not equal targets",
			expectBool: false,
			i:          models.Target{},
			ctx: models.Task{
				TargetsByType: map[models.TargetType][]string{models.DynamicGroup: []string{"notempty", "notempty"}},
			},
		},
		{
			name:       "ok",
			expectBool: true,
			i:          models.Target{},
			ctx: models.Task{
				TargetsByType: models.TargetsByType{models.DynamicGroup: []string{"notempty"}},
			},
		},
		{
			name:       "cant have both dynamic site and regular",
			expectBool: false,
			i:          models.Target{},
			ctx: models.Task{
				TargetsByType: models.TargetsByType{
					models.DynamicSite: []string{"notempty"},
					models.Site:        []string{"notempty"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validTargets(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorRequiredMonthly(t *testing.T) {

	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract field, i is nil interface",
			expectBool: false,
			ctx:        apiModels.Repeat{},
		},
		{
			name:       "not Monthly",
			expectBool: false,
			ctx:        apiModels.Repeat{Frequency: apiModels.Daily},
			i:          []int{1},
		},
		{
			name:       "Monthly",
			expectBool: true,
			ctx:        apiModels.Repeat{Frequency: apiModels.Monthly},
			i:          []int{1},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validatorRequiredMonthly(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorOptionalBetween(t *testing.T) {
	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract field, i is nil interface",
			expectBool: false,
			ctx:        apiModels.Schedule{},
		},
		{
			name:       "beetwen is zero",
			expectBool: true,
			ctx:        apiModels.Schedule{Repeat: apiModels.Repeat{}},
			i:          time.Time{},
		},
		{
			name:       "beetwen after run time",
			expectBool: true,
			ctx:        apiModels.Schedule{Repeat: apiModels.Repeat{RunTime: time.Now().UTC()}},
			i:          time.Now().UTC().Add(time.Minute * 5),
		},
		{
			name:       "beetwen less run time",
			expectBool: false,
			ctx:        apiModels.Schedule{Repeat: apiModels.Repeat{RunTime: time.Now().UTC()}},
			i:          time.Now().UTC().Add(time.Minute * -5),
		},
		{
			name:       "hourly",
			expectBool: false,
			ctx:        apiModels.Schedule{Repeat: apiModels.Repeat{Frequency: apiModels.Hourly, RunTime: time.Now().UTC()}},
			i:          time.Now().UTC().Add(time.Minute * -5),
		},
		{
			name:       "run now",
			expectBool: false,
			ctx:        apiModels.Schedule{Regularity: apiModels.RunNow, Repeat: apiModels.Repeat{Frequency: apiModels.Hourly, RunTime: time.Now().UTC()}},
			i:          time.Now().UTC().Add(time.Minute * -5),
		},
		{
			name:       "beetwen after run time more than 24h",
			expectBool: false,
			ctx:        apiModels.Schedule{Repeat: apiModels.Repeat{RunTime: time.Now().UTC()}},
			i:          time.Now().UTC().Add(time.Hour * 26),
		},
		{
			name:       "beetwen after run time less 24h",
			expectBool: true,
			ctx:        apiModels.Schedule{Repeat: apiModels.Repeat{RunTime: time.Now().UTC()}},
			i:          time.Now().UTC().Add(time.Hour * 23),
		},
		{
			name:       "beetwen after run time less 24h",
			expectBool: true,
			ctx:        apiModels.Schedule{Regularity: apiModels.OneTime, StartRunTime: time.Now().UTC()},
			i:          time.Now().UTC().Add(time.Hour * 23),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validOptionalBetween(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorCreds(t *testing.T) {

	userCurrentUserCreds := &agentModels.Credentials{UseCurrentUser: true}
	runAsUser := &agentModels.Credentials{Password: "pwd", Username: "usr"}

	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract field",
			expectBool: false,
			ctx:        models.TaskDefinitionDetails{},
		},
		{
			name:       "userCurrentUserCreds",
			expectBool: true,
			ctx:        models.TaskDefinitionDetails{},
			i:          userCurrentUserCreds,
		},
		{
			name:       "Run as user",
			expectBool: true,
			ctx:        models.TaskDefinitionDetails{},
			i:          runAsUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validatorCreds(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorRequiredRecurrentAndOneTime(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Kiev")
	yesterdayInLoc := time.Now().Add(time.Hour * -24).In(loc)

	testCases := []struct {
		name       string
		expectBool bool
		i          interface{}
		ctx        interface{}
	}{
		{
			name:       "cant extract structInstance, ctx is nil",
			expectBool: false,
		},
		{
			name:       "cant extract field, i is nil interface",
			expectBool: false,
			ctx:        apiModels.Schedule{},
		},
		{
			name:       "location is zero, RunNow with Time",
			expectBool: false,
			ctx: apiModels.Schedule{
				Regularity: apiModels.RunNow,
			},
			i: time.Now(),
		},
		{
			name:       "location is zero, OneTime with Time",
			expectBool: true,
			ctx: apiModels.Schedule{
				Regularity: apiModels.OneTime,
			},
			i: time.Now(),
		},
		{
			name:       "location is non zero, OneTime with Time in the past",
			expectBool: false,
			ctx: apiModels.Schedule{
				Regularity: apiModels.OneTime,
				Location:   "Europe/Kiev",
			},
			i: yesterdayInLoc,
		},
		{
			name:       "location is non zero, invalid, OneTime with Time in the past",
			expectBool: false,
			ctx: apiModels.Schedule{
				Regularity: apiModels.OneTime,
				Location:   "Europsse/Kiev",
			},
			i: yesterdayInLoc,
		},
		{
			name:       "location is non zero, valid",
			expectBool: true,
			ctx: apiModels.Schedule{
				Regularity: apiModels.OneTime,
				Location:   "America/Phoenix",
			},
			i: time.Now().In(loc).Add(time.Minute * 15),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validatorRequiredRecurrendAndOneTime(tc.i, tc.ctx)
			if gotBool != tc.expectBool {
				t.Errorf("Wanted %v but got %v", tc.expectBool, gotBool)
			}
		})
	}
}

func TestValidatorOptionalMonthly(t *testing.T) {
	RegisterTestingT(t)

	type payload struct {
		i   interface{}
		ctx interface{}
	}

	type expected struct {
		result bool
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "can't cast ctx to needed type",
			payload: payload{
				ctx: apiModels.Schedule{},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is not monthly",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(2),
				},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "repeat frequency is not monthly and days of month are present",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(2),
					DaysOfMonth: []int{2, 5},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is not monthly and weekday is present",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(2),
					WeekDays: []apiModels.WeekDay{
						{
							Day:   time.Weekday(5),
							Index: apiModels.Index(2),
						},
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is not monthly and weekday and days of month are present",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(2),
					DaysOfMonth: []int{2, 5},
					WeekDays: []apiModels.WeekDay{
						{
							Day:   time.Weekday(5),
							Index: apiModels.Index(2),
						},
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is monthly and weekdays and days of month are empty",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(4),
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is monthly and weekdays and days of month are present",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(4),
					DaysOfMonth: []int{1, 7, 9},
					WeekDays: []apiModels.WeekDay{
						{
							Day:   2,
							Index: 2,
						},
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is monthly and weekday present",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(4),
					WeekDays: []apiModels.WeekDay{
						{
							Day:   time.Weekday(5),
							Index: apiModels.Index(2),
						},
					},
				},
				i: []apiModels.WeekDay{
					{
						Day:   time.Weekday(5),
						Index: apiModels.Index(2),
					},
				},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "repeat frequency is monthly and weekday present and day are not valid",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(4),
					WeekDays: []apiModels.WeekDay{
						{
							Day:   time.Weekday(8),
							Index: apiModels.Index(2),
						},
					},
				},
				i: []apiModels.WeekDay{
					{
						Day:   time.Weekday(8),
						Index: apiModels.Index(2),
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is monthly and weekday present and day are more than 4",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(4),
					WeekDays: []apiModels.WeekDay{
						{
							Day:   time.Weekday(5),
							Index: apiModels.Index(2),
						},
						{
							Day:   time.Weekday(5),
							Index: apiModels.Index(1),
						},
						{
							Day:   time.Weekday(2),
							Index: apiModels.Index(0),
						},
						{
							Day:   time.Weekday(2),
							Index: apiModels.Index(2),
						},
						{
							Day:   time.Weekday(3),
							Index: apiModels.Index(4),
						},
					},
				},
				i: []apiModels.WeekDay{
					{
						Day:   time.Weekday(5),
						Index: apiModels.Index(2),
					},
					{
						Day:   time.Weekday(5),
						Index: apiModels.Index(1),
					},
					{
						Day:   time.Weekday(2),
						Index: apiModels.Index(0),
					},
					{
						Day:   time.Weekday(2),
						Index: apiModels.Index(2),
					},
					{
						Day:   time.Weekday(3),
						Index: apiModels.Index(4),
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is monthly and weekday is empty",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(4),
					DaysOfMonth: []int{1, 8},
					WeekDays:    []apiModels.WeekDay{},
				},
				i: []apiModels.WeekDay{},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "repeat frequency is monthly and daysOfMonth are present",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(4),
					DaysOfMonth: []int{1, 25},
					WeekDays:    []apiModels.WeekDay{},
				},
				i: []int{1, 25},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "repeat frequency is monthly and daysOfMonth are not valid",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(4),
					DaysOfMonth: []int{1, 32},
					WeekDays:    []apiModels.WeekDay{},
				},
				i: []int{1, 32},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "repeat frequency is monthly and daysOfMonth are empty",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(4),
					DaysOfMonth: []int{},
					WeekDays: []apiModels.WeekDay{
						{},
					},
				},
				i: []int{},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "repeat frequency is monthly and field type is not required",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency:   apiModels.Frequency(4),
					DaysOfMonth: []int{},
					WeekDays: []apiModels.WeekDay{
						{},
					},
				},
				i: []int64{},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "old week day + new exists",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(4),
					WeekDay: &apiModels.WeekDay{
						Day:   time.Weekday(5),
						Index: apiModels.Index(2),
					},
					WeekDays: []apiModels.WeekDay{
						{
							Day:   time.Weekday(5),
							Index: apiModels.Index(2),
						},
					},
				},
				i: []apiModels.WeekDay{
					{
						Day:   time.Weekday(5),
						Index: apiModels.Index(2),
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "old week day only",
			payload: payload{
				ctx: apiModels.Repeat{
					Frequency: apiModels.Frequency(4),
					WeekDay: &apiModels.WeekDay{
						Day:   time.Weekday(5),
						Index: apiModels.Index(2),
					},
				},
				i: []apiModels.WeekDay{
					{},
				},
			},
			expected: expected{
				result: true,
			},
		},
	}

	for _, test := range tc {
		result := validatorOptionalMonthly(test.payload.i, test.payload.ctx)
		Ω(result).To(Equal(test.expected.result), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestValidatorRequiredForTriggers(t *testing.T) {
	RegisterTestingT(t)

	startTime, _ := time.Parse(time.RFC3339, "2030-05-16T13:40:00Z")
	startTimeAnotherDay, _ := time.Parse(time.RFC3339, "2030-05-18T13:40:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2030-05-16T15:40:00Z")

	type payload struct {
		i   interface{}
		ctx interface{}
	}

	type expected struct {
		result bool
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name:    "can't cast ctx to needed type",
			payload: payload{},
			expected: expected{
				result: false,
			},
		},
		{
			name: "can't cast field to needed type",
			payload: payload{
				ctx: apiModels.Schedule{},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "regularity is not trigger and length of field is more than 0",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(2),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType: "wrongType",
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "regularity is not trigger and length of field is equal 0",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(2),
				},
				i: []apiModels.TriggerFrame{},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "regularity is recurrent and length of field is equal 0",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(3),
				},
				i: []apiModels.TriggerFrame{},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "validation failed while empty triggerFrames with Trigger regularity",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(4),
				},
				i: []apiModels.TriggerFrame{},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "validation failed while empty trigger type",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(4),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType:    "logOn",
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
					{
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "validation failed while startTimeFrame is after endTimeFrame",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(4),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType:    "logOn",
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: endTime,
						EndTimeFrame:   startTime,
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "validation successful while days are different but hours are correct",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(4),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType:    "logOn",
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: startTimeAnotherDay,
						EndTimeFrame:   endTime,
					},
				},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "validation failed while StartTimeFrame is empty",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(4),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType:  "logOn",
						EndTimeFrame: endTime,
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "validation failed while EndTimeFrame is empty",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(4),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType:    "logOn",
						StartTimeFrame: startTime,
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "validation successful",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(4),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType: "logOn",
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
				},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "validation successful with recurrent regularity",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(3),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType: "logOn",
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: startTime,
						EndTimeFrame:   endTime,
					},
				},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "validation failed with recurrent regularity #1",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(3),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType: "logOn",
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: endTime,
						EndTimeFrame:   startTime,
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "validation failed with recurrent regularity #2",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(3),
				},
				i: []apiModels.TriggerFrame{
					{
						StartTimeFrame: endTime,
						EndTimeFrame:   startTime,
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: endTime,
						EndTimeFrame:   startTime,
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "validation failed with unknown regularity",
			payload: payload{
				ctx: apiModels.Schedule{
					Regularity: apiModels.Regularity(6),
				},
				i: []apiModels.TriggerFrame{
					{
						TriggerType: "logOn",
					},
					{
						TriggerType:    "logOff",
						StartTimeFrame: endTime,
						EndTimeFrame:   startTime,
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
	}

	for _, test := range tc {
		result := validatorRequiredForTriggers(test.payload.i, test.payload.ctx)
		Ω(result).To(Equal(test.expected.result), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestValidatorTriggerTypeOptional(t *testing.T) {
	type payload struct {
		i   interface{}
		ctx interface{}
	}

	testCases := []struct {
		name         string
		inputData    payload
		expectedBool bool
	}{
		{
			name: "1. ctx is not schedule",
			inputData: payload{
				i:   apiModels.Schedule{},
				ctx: []string{},
			},
			expectedBool: false,
		},
		{
			name: "2. triggerTypes has wrong type",
			inputData: payload{
				i:   apiModels.Schedule{},
				ctx: apiModels.Schedule{},
			},
			expectedBool: false,
		},
		{
			name: "3. triggerTypes is empty",
			inputData: payload{
				i:   []string{},
				ctx: apiModels.Schedule{},
			},
			expectedBool: true,
		},
		{
			name: "3. triggerTypes is empty while regularity is Trigger ",
			inputData: payload{
				i: []string{},
				ctx: apiModels.Schedule{
					Regularity: apiModels.Trigger,
				},
			},
			expectedBool: false,
		},
		{
			name: "4. triggerTypes != triggerFrames len",
			inputData: payload{
				i: []string{"triggerType1"},
				ctx: apiModels.Schedule{
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: "triggerType1",
						},
						{
							TriggerType: "triggerType2",
						},
					},
				},
			},
			expectedBool: false,
		},
		{
			name: "5. triggerTypes has 2 same triggerTypes",
			inputData: payload{
				i: []string{"triggerType1", "triggerType1", "triggerType2"},
				ctx: apiModels.Schedule{
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: "triggerType1",
						},
						{
							TriggerType: "triggerType1",
						},
						{
							TriggerType: "triggerType2",
						},
					},
				},
			},
			expectedBool: false,
		},
		{
			name: "6. triggerTypes has trigger that is not presented in TriggerFrames but has the same lenght",
			inputData: payload{
				i: []string{"triggerType1", "triggerType2", "triggerType3"},
				ctx: apiModels.Schedule{
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: "triggerType1",
						},
						{
							TriggerType: "triggerType2",
						},
						{
							TriggerType: "triggerType4",
						},
					},
				},
			},
			expectedBool: false,
		},
		{
			name: "7. triggerTypes has unique triggerTypes",
			inputData: payload{
				i: []string{"triggerType1", "triggerType2"},
				ctx: apiModels.Schedule{
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: "triggerType1",
						},
						{
							TriggerType: "triggerType2",
						},
					},
				},
			},
			expectedBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool := validatorTriggerTypeOptional(tc.inputData.i, tc.inputData.ctx)
			if gotBool != tc.expectedBool {
				t.Errorf("Want %v but got %v", tc.expectedBool, gotBool)
			}
		})
	}
}

func TestValidPositive(t *testing.T) {
	if validPositive(models.Task{}, models.Task{}) {
		t.Fatal("wrong object")
	}

	if !validPositive(0, models.Task{}) {
		t.Fatal("wrong object")
	}

	if validPositive(-20, models.Task{}) {
		t.Fatal("wrong object")
	}
}

func TestValidateRecurrentDGTriggerTarget(t *testing.T) {
	RegisterTestingT(t)

	type payload struct {
		i   interface{}
		ctx interface{}
	}

	type expected struct {
		result bool
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name:    "can't cast ctx to needed type",
			payload: payload{},
			expected: expected{
				result: false,
			},
		},
		{
			name: "can't cast field to needed type",
			payload: payload{
				ctx: models.Task{},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "fail with regularity is Recurrent and has DG EnterTrigger and TargetType is Site",
			payload: payload{
				ctx: models.Task{
					TargetsByType: map[models.TargetType][]string{
						models.Site: []string{"id"},
					},
				},
				i: apiModels.Schedule{
					Regularity: apiModels.Recurrent,
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: triggers.DynamicGroupEnterTrigger,
						},
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "fail with regularity is Recurrent and has DG ExitTrigger and TargetType is Site",
			payload: payload{
				ctx: models.Task{
					TargetsByType: map[models.TargetType][]string{
						models.Site: []string{"id"},
					},
				},
				i: apiModels.Schedule{
					Regularity: apiModels.Recurrent,
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: triggers.DynamicGroupExitTrigger,
						},
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "fail with regularity is Recurrent and has DG EnterTrigger and TargetType is Endpoint",
			payload: payload{
				ctx: models.Task{
					TargetsByType: map[models.TargetType][]string{
						models.ManagedEndpoint: []string{"id"},
					},
				},
				i: apiModels.Schedule{
					Regularity: apiModels.Recurrent,
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: triggers.DynamicGroupEnterTrigger,
						},
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "fail with regularity is Recurrent and has DG ExitTrigger and TargetType is Endpoint",
			payload: payload{
				ctx: models.Task{
					TargetsByType: map[models.TargetType][]string{
						models.ManagedEndpoint: []string{"id"},
					},
				},
				i: apiModels.Schedule{
					Regularity: apiModels.Recurrent,
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: triggers.DynamicGroupExitTrigger,
						},
					},
				},
			},
			expected: expected{
				result: false,
			},
		},
		{
			name: "succeed with regularity is Recurrent and has DG EnterTrigger and TargetType is DG",
			payload: payload{
				ctx: models.Task{
					TargetsByType: map[models.TargetType][]string{
						models.DynamicGroup: []string{"id"},
					},
				},
				i: apiModels.Schedule{
					Regularity: apiModels.Recurrent,
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: triggers.DynamicGroupEnterTrigger,
						},
					},
				},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "succeed with regularity is Recurrent and has DG ExitTrigger and TargetType is DG",
			payload: payload{
				ctx: models.Task{
					TargetsByType: map[models.TargetType][]string{
						models.DynamicGroup: []string{"id"},
					},
				},
				i: apiModels.Schedule{
					Regularity: apiModels.Recurrent,
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: triggers.StartupTrigger,
						},
						{
							TriggerType: triggers.DynamicGroupExitTrigger,
						},
					},
				},
			},
			expected: expected{
				result: true,
			},
		},
		{
			name: "succeed with regularity is Trigger and has DG EnterTrigger and TargetType is DG",
			payload: payload{
				ctx: models.Task{
					TargetsByType: map[models.TargetType][]string{
						models.DynamicGroup: []string{"id"},
					},
				},
				i: apiModels.Schedule{
					Regularity: apiModels.Trigger,
					TriggerFrames: []apiModels.TriggerFrame{
						{
							TriggerType: triggers.DynamicGroupEnterTrigger,
						},
					},
				},
			},
			expected: expected{
				result: true,
			},
		},
	}
	for _, test := range tc {
		result := validateRecurrentDGTriggerTarget(test.payload.i, test.payload.ctx)
		Ω(result).To(Equal(test.expected.result), fmt.Sprintf(defaultMsg, test.name))
	}
}
