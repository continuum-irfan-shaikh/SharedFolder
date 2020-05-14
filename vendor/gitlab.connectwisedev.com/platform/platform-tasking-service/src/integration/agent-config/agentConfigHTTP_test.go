package agentConfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	common "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-integration"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
)

func TestNew(t *testing.T) {
	a := NewAgentConfClient(nil, "", nil)
	if a == nil {
		t.Fatal("can't be nil")
	}
}

func Test_agentConfClient_Activate(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type args struct {
		content             entities.Rule
		managedEndpointsIDs map[string]entities.Endpoints
	}
	tests := []struct {
		name          string
		cli           func() integration.HTTPClient
		args          args
		wantProfileID gocql.UUID
		wantErr       bool
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
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{},
			},
			wantProfileID: gocql.UUID{},
		},
		{
			name: "Wrong response status",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rw{},
				}, nil)
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{},
			},
			wantErr: true,
		},
		{
			name: "Failed to send request",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{},
			},
			wantErr: true,
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
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{"test": {}},
			},
			wantErr: true,
		},
		{
			name: "Failed to unmarshal response",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       &rwEmpty{},
				}, nil)
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{"test": {}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &agentConfClient{
			cli:               tt.cli(),
			log:               logger.Log,
			agentConfigDomain: "",
		}
		gotProfileID, err := s.Activate(context.Background(), tt.args.content, tt.args.managedEndpointsIDs, "test")
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("Activate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(gotProfileID).To(gomega.Equal(tt.wantProfileID), fmt.Sprintf("Activate() name = %s, gotResult = %v, want %v", tt.name, gotProfileID, tt.wantProfileID))
		}
	}
}

func Test_agentConfClient_Deactivate(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		cli     func() integration.HTTPClient
		wantErr bool
	}{
		{
			name: "Success",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusNoContent,
					Body:       &rw{},
				}, nil)
				return c
			},
		},
		{
			name: "Wrong response status",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rw{},
				}, nil)
				return c
			},
			wantErr: true,
		},
		{
			name: "Failed to send request",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			wantErr: true,
		},
		{
			name: "Failed to read and close body",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusNoContent,
					Body:       &rwFail{},
				}, nil)
				return c
			},
			wantErr: true,
		},

		{
			name: "Failed to unmarshal response",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusNoContent,
					Body:       &rwEmpty{},
				}, nil)
				return c
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &agentConfClient{
			cli:               tt.cli(),
			log:               logger.Log,
			agentConfigDomain: "",
		}
		err := s.Deactivate(context.Background(), gocql.UUID{}, "test")
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("Deactivate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_agentConfClient_Update(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type args struct {
		content             entities.Rule
		managedEndpointsIDs map[string]entities.Endpoints
		partnerID           string
		profileID           gocql.UUID
	}
	tests := []struct {
		name    string
		cli     func() integration.HTTPClient
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rw{},
				}, nil)
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{},
			},
		},
		{
			name: "Wrong response status",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusNoContent,
					Body:       &rw{},
				}, nil)
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{},
			},
			wantErr: true,
		},
		{
			name: "Failed to send request",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{},
			},
			wantErr: true,
		},
		{
			name: "Failed to read and close body",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rwFail{},
				}, nil)
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{"test": {}},
			},
			wantErr: true,
		},
		{
			name: "Failed to unmarshal response",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rwEmpty{},
				}, nil)
				return c
			},
			args: args{
				content:             entities.Rule{},
				managedEndpointsIDs: map[string]entities.Endpoints{"test": {}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &agentConfClient{
			cli:               tt.cli(),
			log:               logger.Log,
			agentConfigDomain: "",
		}
		err := s.Update(context.Background(), tt.args.content, tt.args.managedEndpointsIDs, tt.args.partnerID, tt.args.profileID)
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("Update() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

type rw struct{}

func (*rw) Read(p []byte) (n int, err error) {
	s := entities.AgentActivateResp{
		ProfileID: gocql.UUID{},
	}
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
