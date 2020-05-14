package appLoader

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gopkg.in/jarcoal/httpmock.v1"
)

var (
	engine           = "VBA"
	categories       = []string{"Management", "Security"}
	someTime         = time.Now().UTC()
	generalPartnerID = "00000000-0000-0000-0000-000000000000"
)

// DefaultScripts is an array of predefined scripts
var DefaultScripts = []apiModels.Script{
	{ID: "00000000-0000-0000-0000-000000000000", PartnerID: generalPartnerID, Name: "name0", Description: "description0", Content: "content0", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema0", UISchema: "UISchema0"},
	{ID: "11111111-1111-1111-1111-111111111111", PartnerID: generalPartnerID, Name: "name1", Description: "description1", Content: "content1", Engine: engine, Categories: categories, CreatedAt: someTime, JSONSchema: "JSONSchema1", UISchema: "UISchema1"},
}

func TestLoadApplicationServices(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerScriptingResponder()

	testLoadingOfAppLoaderServices(t, true)
}

func testLoadingOfAppLoaderServices(t *testing.T, isTest bool) {

	AppLoaderService = nil

	LoadApplicationServices(isTest)

	if AppLoaderService == nil {
		t.Errorf("AppLoaderServices [isTest=%t] were not loaded properly", true)
	}
}

func registerScriptingResponder() {
	httpmock.RegisterResponder("GET", config.Config.ScriptingMsURL+"/scripts",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, DefaultScripts)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
			}
			return resp, nil
		},
	)
	fmt.Println("Scripting responder was successfully registered!")
}
