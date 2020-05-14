package automationEngine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	common "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-integration"
)

const triggerDefinition = `{
    "id": "11d43b69-36bf-4467-8c81-e0e8f4d1990f",
    "displayName": "Endpoint Entering Dynamic Group",
	"internal": false,
    "description": "A new Endpoint enters one of dynamic groups selected in the Task Targets list",
    "triggerCategory": "dynamic-group",
    "eventDetails": {
        "topic": "dynamic_group_change",
        "messageIdentifier": {
            "type": "ENDPOINT-ADDED"
        },
        "endpointIdentifiers": {
            "partnerID": {
                "key": "partner_id"
            },
            "clientID": {
                "key": "client_id"
            },
            "siteID": {
                "key": "site_id"
            },
            "endpointID": {
                "key": "endpoint_id"
            }
        },
        "payloadIdentifiers": {
            "dynamicGroupID": {
                "key": "dynamic_group_id"
            }
        }
    }
}`

func TestClient_GeneratePolicyFiles(t *testing.T) {
	var trDef entities.TriggerDefinition
	err := json.Unmarshal([]byte(triggerDefinition), &trDef)
	if err != nil {
		t.Fatalf("Can't unmarshall triggerDefinition %v", err)
	}

	cli := New(nil, nil, "", "http://127.0.0.1:12121/tasking/v1")
	tempFolderName, _, err := cli.generatePolicyFiles([]entities.TriggerDefinition{trDef})
	if err != nil {
		t.Fatalf("Can't generate template %v", err)
	}

	expectedBytes, err := ioutil.ReadFile("11d43b69-36bf-4467-8c81-e0e8f4d1990fExpected.policy")
	if err != nil {
		t.Fatalf("Can't open expected file")
	}

	gotBytes, err := ioutil.ReadFile(tempFolderName + "/" + "11d43b69-36bf-4467-8c81-e0e8f4d1990f.policy")
	if err != nil {
		t.Fatalf("Can't open received file")
	}

	if !bytes.Equal(expectedBytes, gotBytes) {
		t.Fatalf("Files are not equal")
	}

	if err = cli.removeCreatedFiles(tempFolderName); err != nil {
		t.Fatalf("Can't generate template %v", err)
	}
	fmt.Println(alertPolicyTemplate)
}

func TestClient_GeneratePolicyFilesAlerting(t *testing.T) {
	const triggerDefinition = `{
  "id": "alert-27025",
  "displayName": "Low Server disk space",
  "description": "Emergency low disk space issue observed On Server",
  "internal": false,
  "triggerCategory": "alert",
  "eventDetails": {
    "topic": "alerting",
    "messageIdentifier": {
      "condition_id": "27025"
    },
    "endpointIdentifiers": {
      "partnerID": {
        "key": "partner_id"
      },
      "clientID": {
        "key": "client_id"
      },
      "siteID": {
        "key": "site_id"
      },
      "endpointID": {
        "key": "endpoint_id"
      }
    }
  }
}`
	var trDef entities.TriggerDefinition
	err := json.Unmarshal([]byte(triggerDefinition), &trDef)
	if err != nil {
		t.Fatalf("Can't unmarshall triggerDefinition %v", err)
	}

	cli := New(nil, nil, "", "http://127.0.0.1:12121/tasking/v1")
	tempFolderName, _, err := cli.generatePolicyFiles([]entities.TriggerDefinition{trDef})
	if err != nil {
		t.Fatalf("Can't generate template %v", err)
	}

	expectedBytes, err := ioutil.ReadFile("alert-27025Expected.policy")
	if err != nil {
		t.Fatalf("Can't open expected file")
	}

	gotBytes, err := ioutil.ReadFile(tempFolderName + "/" + "alert-27025.policy")
	if err != nil {
		t.Fatalf("Can't open received file")
	}

	if !bytes.Equal(expectedBytes, gotBytes) {
		t.Fatalf("Files are not equal")
	}

	if err = cli.removeCreatedFiles(tempFolderName); err != nil {
		t.Fatalf("Can't generate template %v", err)
	}
}

func Test_RemovePolicy(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name      string
		cli       func() integration.HTTPClient
		wantError bool
	}{
		{
			name: "Already removed",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
				}, nil)
				return c
			},
		},
		{
			name: "Removed successfully",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusNoContent,
				}, nil)
				return c
			},
		},
		{
			name: "Failed to remove",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
				}, nil)
				return c
			},
			wantError: true,
		},
		{
			name: "Failed",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Client{
			log: logger.Log,
			cli: tt.cli(),
		}
		err := s.RemovePolicy(context.TODO(), map[string]interface{}{})
		ctrl.Finish()

		if !tt.wantError {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("RemovePolicy() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantError))
		}
	}
}

func Test_UpdateRemotePolicies(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name         string
		cli          func() integration.HTTPClient
		wantPolicyID string
		wantError    bool
	}{
		{
			name: "Success",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       &rw{},
				}, nil)
				return c
			},
			wantPolicyID: "test",
		},
		{
			name: "Invalid response status",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       &rw{},
				}, nil)
				return c
			},
			wantError: true,
		},
		{
			name: "Failed to send request",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			wantError: true,
		},
		{
			name: "Failed to read and close body",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       &rwFail{},
				}, nil)
				return c
			},
			wantError: true,
		},
		{
			name: "Failed to unmarshal body",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       &rwEmpty{},
				}, nil)
				return c
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		var trDef entities.TriggerDefinition
		err := json.Unmarshal([]byte(triggerDefinition), &trDef)
		gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("Can't unmarshall triggerDefinition %v", err))

		ctrl = gomock.NewController(t)
		s := &Client{
			log:           logger.Log,
			cli:           tt.cli(),
			taskingDomain: "/",
		}
		policyID, err := s.UpdateRemotePolicies(context.TODO(), []entities.TriggerDefinition{trDef})
		ctrl.Finish()

		if !tt.wantError {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("UpdateRemotePolicies() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantError))
			gomega.Expect(policyID).To(gomega.Equal(tt.wantPolicyID), fmt.Sprintf("UpdateRemotePolicies() name = %s, policyID = %v, wantPolicyID %v", tt.name, policyID, tt.wantPolicyID))
		} else {
			gomega.Expect(err).NotTo(gomega.BeNil(), fmt.Sprintf("UpdateRemotePolicies() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantError))
		}
	}
}

type rw struct{}

func (*rw) Read(p []byte) (n int, err error) {
	s := []automationEngineResponse{{
		ID: "test",
	}}
	b, _ := json.Marshal(&s)
	for i := 0; i < len(b); i++ {
		p[i] = b[i]
	}
	return len(b), io.EOF
}
func (*rw) Close() error { return nil }

type rwEmpty struct{}

func (*rwEmpty) Read(p []byte) (n int, err error) {
	return 10, io.EOF
}
func (*rwEmpty) Close() error { return nil }

type rwFail struct{}

func (*rwFail) Read(_ []byte) (n int, err error) { return 0, errors.New("fail") }
func (*rwFail) Close() error                     { return errors.New("fail") }
