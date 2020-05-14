package sites

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	common "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-integration"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"gopkg.in/jarcoal/httpmock.v1"
)

type sitesData struct {
	SiteList []SiteList `json:"siteDetailList"`
}

type SiteList struct {
	ID int64 `json:"siteId"`
}

func TestGetSiteIDs(t *testing.T) {

	var (
		httpClient          = http.DefaultClient
		partnerID           = "1"
		iPlanetDirectoryPro = "IPlanetDirectoryPro"
		sitesELB            = "localhost:8080"
		url                 = fmt.Sprintf("%s/partner/%s/sites", sitesELB, partnerID)
		sites               = sitesData{SiteList: []SiteList{{ID: 1}, {ID: 2}}}
	)

	sitesPayload, err := json.Marshal(&sites)
	if err != nil {
		t.Fatalf("cant marshall json %v", err)
	}

	testCases := []struct {
		name           string
		isNeedHTTPMock bool
		isNeedStatusOk bool
		sitesPayload   []byte
		expectedError  bool
		expectedSites  []int64
	}{
		{
			name:          "testCase 0 - no response from sites",
			expectedError: true,
		},
		{
			name:           "testCase 1 - Status not ok from httpRequest",
			isNeedHTTPMock: true,
			expectedError:  true,
		},
		{
			name:           "testCase 2 - Cannot unmarshal received body ",
			isNeedHTTPMock: true,
			isNeedStatusOk: true,
			sitesPayload:   []byte{byte('s')},
			expectedError:  true,
		},
		{
			name:           "testCase 3 - Ok ",
			isNeedHTTPMock: true,
			isNeedStatusOk: true,
			sitesPayload:   sitesPayload,
			expectedSites:  []int64{1, 2},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tc.isNeedHTTPMock {
				httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
					if tc.isNeedStatusOk {
						return httpmock.NewBytesResponse(http.StatusOK, tc.sitesPayload), nil
					}
					return httpmock.NewJsonResponse(http.StatusInternalServerError, nil)
				})
			}

			gotSites, gotErr := GetSiteIDs(context.Background(), httpClient, partnerID, sitesELB, iPlanetDirectoryPro)
			if tc.expectedError && gotErr == nil {
				t.Fatalf("Wanted err but got nil")
			}

			if !reflect.DeepEqual(tc.expectedSites, gotSites) && !tc.expectedError {
				t.Fatalf("Wanted %v but got %v", tc.expectedSites, gotSites)
			}

		})
	}
}

func Test_getPartnerSites(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		cli     func() integration.HTTPClient
		wantIDs []int64
		wantErr bool
	}{
		{
			name: "Success",

			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: &rw{},
				}, nil)
				return c
			},
			wantIDs: []int64{0},
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
					Body: &rwFail{},
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
					Body: &rwEmpty{},
				}, nil)
				return c
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)

		ids, err := getPartnerSites(context.TODO(), tt.cli(), "test")
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("getPartnerSites() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(ids).To(gomega.Equal(tt.wantIDs), fmt.Sprintf("GetEndpointsByGroupIDs() name = %s, gotResult = %v, want %v", tt.name, ids, tt.wantIDs))
		} else {
			gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("getPartnerSites() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

type rw struct{}

func (*rw) Read(p []byte) (n int, err error) {
	s := struct {
		Sites []struct {
			ID int64 `json:"siteId"`
		} `json:"outdata"`
	}{
		[]struct {
			ID int64 `json:"siteId"`
		}{{ID: 0}},
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
