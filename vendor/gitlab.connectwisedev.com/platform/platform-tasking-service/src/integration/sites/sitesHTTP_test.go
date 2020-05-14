package sites

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
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	common "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-integration"
)

func TestNew(t *testing.T) {
	a := NewClient(nil)
	if a == nil {
		t.Fatal("can't be nil")
	}
}

func TestClient_GetEndpointsBySiteIDs(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		cli     func() integration.HTTPClient
		wantIDs []gocql.UUID
		wantErr bool
	}{
		{
			name: "Success",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rwSites{},
				}, nil)
				return c
			},
			wantIDs: []gocql.UUID{{}},
		},
		{
			name: "Not found",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusNotFound,
					Body:       &rwEmpty{},
				}, nil)
				return c
			},
			wantIDs: nil,
		},
		{
			name: "Failed make call",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			wantErr: true,
		},
		{
			name: "Failed read and close response body",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rwFail{},
				}, nil)
				return c
			},
			wantErr: true,
		},
		{
			name: "Failed unmarshal response body",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       &rwEmpty{},
				}, nil)
				return c
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)

		c := &Client{
			cli: tt.cli(),
		}
		ids, err := c.GetEndpointsBySiteIDs(context.Background(), "test", []string{"testSite"})
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("getPartnerSites() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(ids).To(gomega.Equal(tt.wantIDs), fmt.Sprintf("GetEndpointsByGroupIDs() name = %s, gotResult = %v, want %v", tt.name, ids, tt.wantIDs))
		} else {
			gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("getPartnerSites() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

type rwSites struct{}

func (*rwSites) Read(p []byte) (n int, err error) {
	s := []sitesResponse{{EndpointID: gocql.UUID{}}}
	b, _ := json.Marshal(&s)
	for i := 0; i < len(b); i++ {
		p[i] = b[i]
	}
	return len(b), io.EOF
}
func (*rwSites) Close() error { return nil }
