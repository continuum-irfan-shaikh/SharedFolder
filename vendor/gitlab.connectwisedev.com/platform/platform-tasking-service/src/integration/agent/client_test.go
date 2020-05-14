package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	common "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-integration"
)

func TestNew(t *testing.T) {
	cli := NewClient(nil, "", nil)
	if cli == nil {
		t.Fatal("can't be nil")
	}
}

func Test_client_Encrypt(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type args struct {
		endpointID  gocql.UUID
		credentials agent.Credentials
	}
	tests := []struct {
		name          string
		httpClient    func() integration.HTTPClient
		args          args
		wantEncrypted agent.Credentials
		wantErr       bool
	}{
		{
			name: "Success",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rw{},
				}, nil).Times(3)
				return c
			},
			args: args{
				endpointID: gocql.UUID{},
				credentials: agent.Credentials{
					UseCurrentUser: false,
					Username:       "test",
					Domain:         "test",
					Password:       "test",
				},
			},
			wantEncrypted: agent.Credentials{
				UseCurrentUser: false,
				Username:       "test",
				Domain:         "test",
				Password:       "test",
			},
		},
		{
			name: "Failed to encrypt pwd",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail")).Times(1)
				return c
			},
			args: args{
				endpointID: gocql.UUID{},
				credentials: agent.Credentials{
					UseCurrentUser: false,
					Username:       "test",
					Domain:         "test",
					Password:       "test",
				},
			},
			wantErr: true,
		},
		{
			name: "Failed to encrypt uname",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rw{},
				}, nil).Times(1)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail")).Times(1)
				return c
			},
			args: args{
				endpointID: gocql.UUID{},
				credentials: agent.Credentials{
					UseCurrentUser: false,
					Username:       "test",
					Domain:         "test",
					Password:       "test",
				},
			},
			wantErr: true,
		},
		{
			name: "Failed to encrypt domain",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rw{},
				}, nil).Times(2)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail")).Times(1)
				return c
			},
			args: args{
				endpointID: gocql.UUID{},
				credentials: agent.Credentials{
					UseCurrentUser: false,
					Username:       "test",
					Domain:         "test",
					Password:       "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		c := &client{
			httpClient: tt.httpClient(),
			agentELB:   "",
			log:        logger.Log,
		}
		gotEncrypted, err := c.Encrypt(context.TODO(), tt.args.endpointID, tt.args.credentials)

		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("encrypt() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(gotEncrypted).To(gomega.Equal(tt.wantEncrypted), fmt.Sprintf("encrypt() name = %s, gotResult = %v, want %v", tt.name, gotEncrypted, tt.wantEncrypted))
		}
	}
}

func Test_client_encrypt(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name       string
		httpClient func() integration.HTTPClient
		body       payload
		wantResult string
		wantErr    bool
	}{
		{
			name: "Success",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rw{},
				}, nil)
				return c
			},
			body: payload{
				EndpointID: gocql.UUID{},
				Data:       "test",
			},
			wantResult: "test",
		},
		{
			name: "Empty request",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				return c
			},
			body: payload{
				EndpointID: gocql.UUID{},
			},
		},
		{
			name: "Invalid request",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				return c
			},
			body: payload{
				EndpointID: gocql.UUID{},
			},
		},
		{
			name: "Request failed",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			body: payload{
				EndpointID: gocql.UUID{},
				Data:       "test",
			},
			wantErr: true,
		},
		{
			name: "Read and Close failed",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rwFail{},
				}, nil)
				return c
			},
			body: payload{
				EndpointID: gocql.UUID{},
				Data:       "test",
			},
			wantErr: true,
		},
		{
			name: "Response code !200",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       &rw{},
				}, nil)
				return c
			},
			body: payload{
				EndpointID: gocql.UUID{},
				Data:       "test",
			},
			wantErr: true,
		},
		{
			name: "Invalid response body",
			httpClient: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rwEmpty{},
				}, nil)
				return c
			},
			body: payload{
				EndpointID: gocql.UUID{},
				Data:       "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		c := &client{
			httpClient: tt.httpClient(),
			agentELB:   "",
			log:        logger.Log,
		}
		gotResult, err := c.encrypt(context.TODO(), tt.body)

		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("encrypt() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
		gomega.Expect(gotResult).To(gomega.Equal(tt.wantResult), fmt.Sprintf("encrypt() name = %s, gotResult = %v, want %v", tt.name, gotResult, tt.wantResult))
	}
}

type rw struct{}

func (*rw) Read(p []byte) (n int, err error) {
	s := payload{
		EndpointID: gocql.UUID{},
		Data:       "test",
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
