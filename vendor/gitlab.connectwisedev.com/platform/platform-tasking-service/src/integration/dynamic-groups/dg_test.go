package dynamicGroups

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
	m "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/dg"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	common "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func TestNew(t *testing.T) {
	a := NewDynamicGroupsClient(nil, nil, nil)
	if a == nil {
		t.Fatal("can't be nil")
	}
}

func Test_GetMachineNameByEndpointID(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type infra struct {
		repo func() models.UserSitesPersistence
		cli  func() integration.HTTPClient
	}
	tests := []struct {
		name    string
		infra   infra
		wantIDs []gocql.UUID
		wantErr bool
	}{
		{
			name: "Success w/o NOC access",
			infra: infra{
				repo: func() models.UserSitesPersistence {
					r := mocks.NewMockUserSitesPersistence(ctrl)
					r.EXPECT().Sites(gomock.Any(), gomock.Any(), gomock.Any()).Return(entities.UserSites{
						PartnerID: "test",
						UserID:    "test",
						SiteIDs:   []int64{0},
					}, nil)
					return r
				},
				cli: func() integration.HTTPClient {
					c := common.NewMockHTTPClient(ctrl)
					c.EXPECT().Do(gomock.Any()).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       &rw{},
					}, nil)
					return c
				},
			},
			wantIDs: []gocql.UUID{{}},
		},
		{
			name: "Failed to get sites",
			infra: infra{
				repo: func() models.UserSitesPersistence {
					r := mocks.NewMockUserSitesPersistence(ctrl)
					r.EXPECT().Sites(gomock.Any(), gomock.Any(), gomock.Any()).Return(entities.UserSites{}, errors.New("fail"))
					return r
				},
				cli: func() integration.HTTPClient {
					c := common.NewMockHTTPClient(ctrl)
					c.EXPECT().Do(gomock.Any()).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       &rw{},
					}, nil)
					return c
				},
			},
			wantErr: true,
		},
		{
			name: "Failed to call DG",
			infra: infra{
				repo: func() models.UserSitesPersistence {
					r := mocks.NewMockUserSitesPersistence(ctrl)
					r.EXPECT().Sites(gomock.Any(), gomock.Any(), gomock.Any()).Return(entities.UserSites{
						PartnerID: "test",
						UserID:    "test",
						SiteIDs:   []int64{0},
					}, nil)
					return r
				},
				cli: func() integration.HTTPClient {
					c := common.NewMockHTTPClient(ctrl)
					c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
					return c
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Client{
			userRepo:   tt.infra.repo(),
			httpClient: tt.infra.cli(),
		}
		ids, err := s.GetEndpointsByGroupIDs(context.Background(), []string{"test"}, "test", "test", false)
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("GetEndpointsByGroupIDs() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(ids).To(gomega.Equal(tt.wantIDs), fmt.Sprintf("GetEndpointsByGroupIDs() name = %s, gotResult = %v, want %v", tt.name, ids, tt.wantIDs))
		}
	}
}

func Test_StartMonitoringGroups(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		client  func() Pusher
		wantErr bool
	}{
		{
			name: "Success",
			client: func() Pusher {
				c := common.NewMockPusher(ctrl)
				c.EXPECT().Push(gomock.Any(), m.MonitoringDG{
					PartnerID:      "test",
					DynamicGroupID: "test",
					ServiceID:      TaskingServiceIDPrefix + "00000000-0000-0000-0000-000000000000",
					Operation:      MessageTypeDynamicGroupStartMonitoring,
				}).Return(nil)
				return c
			},
		},
		{
			name: "Fail",
			client: func() Pusher {
				c := common.NewMockPusher(ctrl)
				c.EXPECT().Push(gomock.Any(), gomock.Any()).Return(errors.New("fail"))
				return c
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Client{
			client: tt.client(),
		}
		err := s.StartMonitoringGroups(context.TODO(), "test", []string{"test"}, gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("StartMonitoringGroups() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("StartMonitoringGroups() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_StopGroupsMonitoring(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		client  func() Pusher
		wantErr bool
	}{
		{
			name: "Success",
			client: func() Pusher {
				c := common.NewMockPusher(ctrl)
				c.EXPECT().Push(gomock.Any(), m.MonitoringDG{
					PartnerID:      "test",
					DynamicGroupID: "test",
					ServiceID:      TaskingServiceIDPrefix + "00000000-0000-0000-0000-000000000000",
					Operation:      MessageTypeDynamicGroupStopMonitoring,
				}).Return(nil)
				return c
			},
		},
		{
			name: "Fail",
			client: func() Pusher {
				c := common.NewMockPusher(ctrl)
				c.EXPECT().Push(gomock.Any(), gomock.Any()).Return(errors.New("fail"))
				return c
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Client{
			client: tt.client(),
		}
		err := s.StopGroupsMonitoring(context.TODO(), "test", []string{"test"}, gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("StopGroupsMonitoring() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("StopGroupsMonitoring() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

type rw struct{}

func (*rw) Read(p []byte) (n int, err error) {
	s := []dgResponse{{
		ID:     gocql.UUID{},
		SiteID: "0",
	}}
	b, _ := json.Marshal(&s)
	for i := 0; i < len(b); i++ {
		p[i] = b[i]
	}
	return len(b), io.EOF
}
func (*rw) Close() error { return nil }
