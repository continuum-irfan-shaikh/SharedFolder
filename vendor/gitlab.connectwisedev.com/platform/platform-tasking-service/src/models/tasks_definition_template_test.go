package models

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"github.com/gocql/gocql"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var (
	engine           = "VBA"
	categories       = []string{"Management", "Security"}
	someTime         = time.Now().UTC()
	partnerID        = "123456789"
	generalPartnerID = "00000000-0000-0000-0000-000000000000"
	newScript        = apiModels.Script{ID: "55555555-5555-5555-5555-555555555555", PartnerID: generalPartnerID, Name: "name5", Description: "description5", Content: "content5", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema4", UISchema: "UISchema4"}
)

// DefaultScripts is an array of predefined scripts
var defaultScripts = []apiModels.Script{
	{ID: "00000000-0000-0000-0000-000000000000", PartnerID: generalPartnerID, Name: "name0", Description: "description0", Content: "content0", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema0", UISchema: "UISchema0", NOCVisibleOnly: true},
	{ID: "00000000-0000-0000-0000-000000000000", PartnerID: generalPartnerID, Name: "name0", Description: "description0", Content: "content0", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema0", UISchema: "UISchema0"},
	{ID: "11111111-1111-1111-1111-111111111111", PartnerID: generalPartnerID, Name: "name1", Description: "description1", Content: "content1", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema1", UISchema: "UISchema1"},
	{ID: "22222222-2222-2222-2222-222222222222", PartnerID: generalPartnerID, Name: "name2", Description: "description2", Content: "content2", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema2", UISchema: "UISchema2"},
	{ID: "33333333-3333-3333-3333-333333333333", PartnerID: generalPartnerID, Name: "name3", Description: "description3", Content: "content3", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema3", UISchema: "UISchema3"},
	{ID: "44444444-4444-4444-4444-444444444444", PartnerID: partnerID, Name: "name4", Description: "description4", Content: "content4", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema4", UISchema: "UISchema4"},
}

func init() {
	config.Load()
	logger.Load(config.Config.Log)
}

func TestLoadTemplatesCache(t *testing.T) {
	expectedScripts := defaultScripts
	config.Config.TDTCacheSettings.ReloadIntervalSec = 1
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// attempt number 1
	registerScriptingResponder(expectedScripts)
	LoadTemplatesCache()
	getAndCountTemplates(t, len(expectedScripts), partnerID)
	// attempt number 2 after add new one script
	expectedScripts = append(expectedScripts, newScript)
	registerScriptingResponder(expectedScripts)
	time.Sleep(2 * time.Second) // We are waiting for cache reloading
	getAndCountTemplates(t, len(expectedScripts), partnerID)
}

func TestGetAllTemplates(t *testing.T) {
	expectedScripts := defaultScripts
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerScriptingResponder(expectedScripts)
	LoadTemplatesCache()

	var tests = []struct {
		partner        string
		templatesCount int
	}{
		{
			partner:        partnerID,
			templatesCount: 6,
		},
	}

	for _, test := range tests {
		getAndCountTemplates(t, test.templatesCount, test.partner)
	}
}

func TestGetByOriginID(t *testing.T) {
	expectedScripts := defaultScripts
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerScriptingResponder(expectedScripts)
	LoadTemplatesCache()

	originID, _ := gocql.ParseUUID(expectedScripts[0].ID)
	originIDScript2, _ := gocql.ParseUUID(expectedScripts[1].ID)
	newScriptOriginID, _ := gocql.ParseUUID(newScript.ID)

	var tests = []struct {
		partner          string
		originID         gocql.UUID
		expectedTemplate TemplateDetails
		hasNoc           bool
		err              bool
		errMsg           string
	}{
		{
			partner: generalPartnerID,
			hasNoc:  true,
			expectedTemplate: TemplateDetails{
				Categories:         expectedScripts[0].Categories,
				Description:        expectedScripts[0].Description,
				Name:               expectedScripts[0].Name,
				OriginID:           originID,
				PartnerID:          expectedScripts[0].PartnerID,
				Type:               "script",
				CreatedAt:          expectedScripts[0].CreatedAt,
				JSONSchema:         expectedScripts[0].JSONSchema,
				UISchema:           expectedScripts[0].UISchema,
				SuccessMessage:     expectedScripts[0].SuccessMessage,
				FailureMessage:     expectedScripts[0].FailureMessage,
				IsParameterized:    true,
				Engine:             expectedScripts[0].Engine,
				IsRequireNOCAccess: true,
			},
		},
		{
			partner: generalPartnerID,
			hasNoc:  false,
			expectedTemplate: TemplateDetails{
				Categories:      expectedScripts[0].Categories,
				Description:     expectedScripts[0].Description,
				Name:            expectedScripts[0].Name,
				OriginID:        originID,
				PartnerID:       expectedScripts[0].PartnerID,
				Type:            "script",
				CreatedAt:       expectedScripts[0].CreatedAt,
				JSONSchema:      expectedScripts[0].JSONSchema,
				UISchema:        expectedScripts[0].UISchema,
				SuccessMessage:  expectedScripts[0].SuccessMessage,
				FailureMessage:  expectedScripts[0].FailureMessage,
				IsParameterized: true,
				Engine:          expectedScripts[0].Engine,
			},
		},
		{
			partner:  "NewPartner",
			hasNoc:   true,
			originID: newScriptOriginID,
			err:      true,
			errMsg:   fmt.Sprintf(`Template with OriginID=55555555-5555-5555-5555-555555555555 not found for partner NewPartner`),
		},
		{
			partner:  generalPartnerID,
			originID: originIDScript2,
			hasNoc:   true,
			expectedTemplate: TemplateDetails{
				Categories:         expectedScripts[1].Categories,
				Description:        expectedScripts[1].Description,
				Name:               expectedScripts[1].Name,
				OriginID:           originIDScript2,
				PartnerID:          expectedScripts[1].PartnerID,
				Type:               "script",
				CreatedAt:          expectedScripts[1].CreatedAt,
				JSONSchema:         expectedScripts[1].JSONSchema,
				UISchema:           expectedScripts[1].UISchema,
				SuccessMessage:     expectedScripts[1].SuccessMessage,
				FailureMessage:     expectedScripts[1].FailureMessage,
				IsParameterized:    true,
				Engine:             expectedScripts[0].Engine,
				IsRequireNOCAccess: true,
			},
			err:    true,
			errMsg: `Template with OriginID=55555555-5555-5555-5555-555555555555 not found for partner 00000000-0000-0000-0000-000000000000`,
		},
	}

	for _, test := range tests {
		actualTemplate, err := TemplateCacheInstance.GetByOriginID(context.TODO(), test.partner, test.originID, test.hasNoc)
		if err != nil && test.err && test.errMsg != err.Error() {
			t.Fatalf("Expected error message %v but got %v", test.errMsg, err.Error())
		}
		if err != nil && !test.err {
			t.Fatalf("Couldn't retrieve template by ID %v, err: %v", test.originID, err)
		}
		if !reflect.DeepEqual(actualTemplate, test.expectedTemplate) {
			t.Fatalf("Expected %v but got %v", test.expectedTemplate, actualTemplate)
		}
	}
}

func registerScriptingResponder(scripts []apiModels.Script) {
	httpmock.RegisterResponder("GET", config.Config.ScriptingMsURL+"/scripts",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, scripts)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	fmt.Println("Scripting responder was successfully registered!")
}

func getAndCountTemplates(t *testing.T, defaultScriptsSize int, partner string) {
	templates, err := TemplateCacheInstance.GetAllTemplatesDetails(context.TODO(), partner)
	if err != nil {
		t.Fatalf("Couldn't retrieve templates, err: %v", err)
	}
	if len(templates) != defaultScriptsSize {
		t.Fatalf("Expected %d templates but got %d", defaultScriptsSize, len(templates))
	}
}

func TestGetExtraTime(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Minute)

	gotTime := getExtraTime(Task{
		Schedule: apiModels.Schedule{
			Repeat: apiModels.Repeat{
				RunTime: now,
			},
			BetweenEndTime: now.Add(time.Minute * 5),
		},
	})

	expected := 5 * 60
	if expected != gotTime {
		t.Fatalf("expected %v but got %v", expected, gotTime)
	}
}
