package asset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	apiModel "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	common "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var (
	mockInstance memcached.MClientMock
)

func init() {
	config.Load()
	err := logger.Load(config.Config.Log)
	if err != nil {
		fmt.Println("init_function: error while loading logger: ", err)
	}

	mockInstance = memcached.MClientMock{
		Cache: make(map[string]*memcache.Item, 10),
	}

	// Need to set default http client instead of tuned in real asset service because of HTTP DefaultMock library couldn't work with it,
	// looks like dirty hack but I believe unit tests are not checking performance, so we could use this piece of...
	ServiceInstance = NewAssetsService(mockInstance, http.DefaultClient)
}

func TestNewAssetsService(t *testing.T) {
	expectedService := Service{
		mCache:     memcached.MClientMock{},
		httpClient: defaultHTTPClient,
	}
	gotService := NewAssetsService(memcached.MClientMock{}, nil)

	if !reflect.DeepEqual(expectedService, gotService) {
		t.Fatalf("Got %v, wanted %v", gotService, expectedService)
	}
}

func TestLoad(t *testing.T) {
	temp := ServiceInstance
	memcached.MemCacheInstance = mockInstance
	Load()
	if reflect.DeepEqual(temp, ServiceInstance) {
		t.Fatalf("Must not be equal")
	}
	ServiceInstance = NewAssetsService(mockInstance, http.DefaultClient)
}

func TestServiceGetEndpointTimeZoneOffsetByEndpointID(t *testing.T) {
	type system struct {
		TimeZone string `json:"timeZone"`
	}

	type assetResponse struct {
		SiteID string `json:"siteID"`
		System system `json:"system"`
	}

	testCases := []struct {
		name                   string
		memcachedMock          memcached.MockCacheConf
		partnerID              string
		endpointID             gocql.UUID
		httpMockResponseStatus int
		httpMockResponseBody   interface{}
		httpMockResponseError  error
		expectedLocationName   string
		assetCacheEnabled      bool
		isError                bool
	}{
		{
			name:              "testCase1 - ok",
			assetCacheEnabled: true,
			memcachedMock: func(mc *memcached.MockCache) *memcached.MockCache {
				mc.EXPECT().
					Get(gomock.Any()).
					Return(&memcache.Item{Value: []byte("+0200")}, nil)
				return mc
			},
			expectedLocationName: "UTC+0200",
			isError:              false,
		},
		{
			name:              "testCase2 - ok, memcached Get method returns error",
			assetCacheEnabled: true,
			memcachedMock: func(mc *memcached.MockCache) *memcached.MockCache {
				mc.EXPECT().
					Get(gomock.Any()).
					Return(&memcache.Item{}, errors.New("error"))
				mc.EXPECT().
					Set(gomock.Any()).
					Return(nil).
					Times(5)
				return mc
			},
			httpMockResponseStatus: http.StatusOK,
			httpMockResponseBody: assetResponse{
				SiteID: "1",
				System: system{TimeZone: "-0400"},
			},
			expectedLocationName: "UTC-0400",
			isError:              false,
		},
		{
			name:              "testCase3 - ok, assetCache is disabled",
			assetCacheEnabled: false,
			memcachedMock: func(mc *memcached.MockCache) *memcached.MockCache {
				return mc
			},
			httpMockResponseStatus: http.StatusOK,
			httpMockResponseBody: assetResponse{
				SiteID: "1",
				System: system{TimeZone: "-0900"},
			},
			expectedLocationName: "UTC-0900",
			isError:              false,
		},
		{
			name:              "testCase4 - ok, memcached is unavailable",
			assetCacheEnabled: true,
			memcachedMock: func(mc *memcached.MockCache) *memcached.MockCache {
				mc.EXPECT().
					Get(gomock.Any()).
					Return(&memcache.Item{}, errors.New("error"))
				mc.EXPECT().
					Set(gomock.Any()).
					Return(errors.New("error")).
					Times(5)
				return mc
			},
			httpMockResponseStatus: http.StatusOK,
			httpMockResponseBody: assetResponse{
				SiteID: "1",
				System: system{TimeZone: "-1000"},
			},
			expectedLocationName: "UTC-1000",
			isError:              false,
		},
		{
			name:              "testCase5 - error, assetMS is unavailable",
			assetCacheEnabled: true,
			memcachedMock: func(mc *memcached.MockCache) *memcached.MockCache {
				mc.EXPECT().
					Get(gomock.Any()).
					Return(&memcache.Item{}, errors.New("error")).
					Times(1)
				return mc
			},
			httpMockResponseError: errors.New("error"),
			isError:               true,
		},
		{
			name:              "testCase6 - error, assetMS is unavailable, assetCache is disabled",
			assetCacheEnabled: false,
			memcachedMock: func(mc *memcached.MockCache) *memcached.MockCache {
				return mc
			},
			httpMockResponseError: errors.New("error"),
			isError:               true,
		},
		{
			name:              "testCase7 - ok, assetCache is disabled and AssetMS returns invalid offset",
			assetCacheEnabled: false,
			memcachedMock: func(mc *memcached.MockCache) *memcached.MockCache {
				return mc
			},
			httpMockResponseStatus: http.StatusOK,
			httpMockResponseBody: assetResponse{
				SiteID: "1",
				System: system{TimeZone: "00900"},
			},
			expectedLocationName: "UTC",
			isError:              false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Config.AssetCacheEnabled = tc.assetCacheEnabled

			mockController := gomock.NewController(t)
			cacheMock := memcached.NewMockCache(mockController)
			cacheMock = tc.memcachedMock(cacheMock)

			httpmock.Activate()
			defer func() {
				mockController.Finish()
				httpmock.DeactivateAndReset()
			}()

			assetMsURL := fmt.Sprintf("%s/partner/%s/endpoints/%v?field=system", config.Config.AssetMsURL, tc.partnerID, tc.endpointID)
			httpmock.RegisterResponder(http.MethodGet, assetMsURL, func(req *http.Request) (*http.Response, error) {
				body, _ := json.Marshal(tc.httpMockResponseBody)
				return httpmock.NewBytesResponse(tc.httpMockResponseStatus, body), tc.httpMockResponseError
			})

			service := NewAssetsService(cacheMock, http.DefaultClient)
			location, err := service.GetLocationByEndpointID(context.Background(), tc.partnerID, tc.endpointID)
			if err == nil && location.String() != tc.expectedLocationName {
				t.Fatalf("%s: expected location.String() == %s, got result == %s",
					tc.name, tc.expectedLocationName, location.String())
			}

			if (err != nil) != tc.isError {
				t.Fatalf("%s: expected isError == %t, got err == %v",
					tc.name, tc.isError, err)
			}
		})
	}
}

func TestParseTimeZoneOffset(t *testing.T) {
	testCases := []struct {
		name           string
		offsetStr      string
		expectedResult time.Duration
		isErr          bool
	}{
		{
			name:           "testCase1 - ok",
			offsetStr:      "-0500",
			expectedResult: -5 * time.Hour,
		},
		{
			name:           "testCase2 - ok",
			offsetStr:      "+0200",
			expectedResult: 2 * time.Hour,
		},
		{
			name:           "testCase3 - ok",
			offsetStr:      "",
			expectedResult: 0,
		},
		{
			name:      "testCase4 - error, invalid offset",
			offsetStr: "invalid offset",
			isErr:     true,
		},
		{
			name:      "testCase5 - error, invalid offset",
			offsetStr: "0100",
			isErr:     true,
		},
		{
			name:      "testCase6 - error, invalid offset",
			offsetStr: "-2400",
			isErr:     true,
		},
		{
			name:      "testCase7 - error, invalid offset",
			offsetStr: "-2360",
			isErr:     true,
		},
		{
			name:      "testCase8 - error, invalid offset",
			offsetStr: "-01000",
			isErr:     true,
		},
		{
			name:      "testCase9 - error, invalid offset",
			offsetStr: "01000",
			isErr:     true,
		},
		{
			name:      "testCase10 - error, invalid offset",
			offsetStr: "+hh30",
			isErr:     true,
		},
		{
			name:      "testCase11 - error, invalid offset",
			offsetStr: "+00mm",
			isErr:     true,
		},
	}

	for _, tc := range testCases {
		result, err := parseTimeZoneOffset(tc.offsetStr)
		if (err != nil) != tc.isErr {
			t.Fatalf("%s: expected error == %v, but got %v", tc.name, tc.isErr, err)
			return
		}

		if result != tc.expectedResult {
			t.Fatalf("%s: expected result == %v, but got %v", tc.name, tc.expectedResult, result)
		}
	}
}

func Test_GetSiteIDByEndpointID(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type args struct {
		partnerID    string
		endpointID   gocql.UUID
		cacheEnabled bool
	}
	tests := []struct {
		name         string
		cli          func() integration.HTTPClient
		cache        func() memcached.Cache
		args         args
		wantSiteID   string
		wantClientID string
		wantErr      bool
	}{
		{
			name: "Success (from asset)",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: &rw{},
				}, nil)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: false,
			},
			wantSiteID:   "testSite",
			wantClientID: "testClient",
		},
		{
			name: "Success (from cache)",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(&memcache.Item{Value: []byte("testSite")}, nil)
				c.EXPECT().Get(gomock.Any()).Return(&memcache.Item{Value: []byte("testClient")}, nil)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: true,
			},
			wantSiteID:   "testSite",
			wantClientID: "testClient",
		},
		{
			name: "Cache error, get from asset",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: &rw{},
				}, nil)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(nil, errors.New("fail"))
				c.EXPECT().Get(gomock.Any()).Return(&memcache.Item{Value: []byte("testClient")}, nil)
				c.EXPECT().Set(gomock.Any()).Return(nil).Times(5)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: true,
			},
			wantSiteID:   "testSite",
			wantClientID: "testClient",
		},
		{
			name: "Failed to get data from asset",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Service{
			httpClient: tt.cli(),
			mCache:     tt.cache(),
		}
		config.Config.AssetCacheEnabled = tt.args.cacheEnabled
		siteID, clientID, err := s.GetSiteIDByEndpointID(context.TODO(),tt.args.partnerID, tt.args.endpointID)
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("GetSiteIDByEndpointID() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(siteID).To(gomega.Equal(tt.wantSiteID), fmt.Sprintf("GetSiteIDByEndpointID() name = %s, gotResult = %v, want %v", tt.name, siteID, tt.wantSiteID))
			gomega.Expect(clientID).To(gomega.Equal(tt.wantClientID), fmt.Sprintf("GetSiteIDByEndpointID() name = %s, gotResult = %v, want %v", tt.name, clientID, tt.wantClientID))
		} else {
			gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("GetSiteIDByEndpointID() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}

	}
}

func Test_GetResourceTypeByEndpointID(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type args struct {
		partnerID    string
		endpointID   gocql.UUID
		cacheEnabled bool
	}
	tests := []struct {
		name        string
		cli         func() integration.HTTPClient
		cache       func() memcached.Cache
		args        args
		wantResType integration.ResourceType
		wantErr     bool
	}{
		{
			name: "Success (from asset)",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: &rw{},
				}, nil)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: false,
			},
			wantResType: "testType",
		},
		{
			name: "Success (from cache)",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(&memcache.Item{Value: []byte("testType")}, nil)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: true,
			},
			wantResType: "testType",
		},
		{
			name: "Cache error, get from asset",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       &rw{},
				}, nil)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(nil, errors.New("fail"))
				c.EXPECT().Set(gomock.Any()).Return(nil).Times(5)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: true,
			},
			wantResType: "testType",
		},
		{
			name: "Failed to get data from asset",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Service{
			httpClient: tt.cli(),
			mCache:     tt.cache(),
		}
		config.Config.AssetCacheEnabled = tt.args.cacheEnabled
		resType, err := s.GetResourceTypeByEndpointID(context.TODO(),tt.args.partnerID, tt.args.endpointID)
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("GetResourceTypeByEndpointID() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(resType).To(gomega.Equal(tt.wantResType), fmt.Sprintf("GetResourceTypeByEndpointID() name = %s, gotResult = %v, want %v", tt.name, resType, tt.wantResType))
		} else {
			gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("GetResourceTypeByEndpointID() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}

	}
}

func Test_GetMachineNameByEndpointID(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type args struct {
		partnerID    string
		endpointID   gocql.UUID
		cacheEnabled bool
	}
	tests := []struct {
		name            string
		cli             func() integration.HTTPClient
		cache           func() memcached.Cache
		args            args
		wantMachineName string
		wantErr         bool
	}{
		{
			name: "Success (from asset)",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: &rw{},
				}, nil)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: false,
			},
			wantMachineName: "testName",
		},
		{
			name: "Success (from cache)",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(&memcache.Item{Value: []byte("testName")}, nil)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: true,
			},
			wantMachineName: "testName",
		},
		{
			name: "Cache error, get from asset",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       &rw{},
				}, nil)
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(nil, errors.New("fail"))
				c.EXPECT().Set(gomock.Any()).Return(nil).Times(5)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: true,
			},
			wantMachineName: "testName",
		},
		{
			name: "Failed to get data from asset",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("fail"))
				return c
			},
			cache: func() memcached.Cache {
				c := memcached.NewMockCache(ctrl)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Service{
			httpClient: tt.cli(),
			mCache:     tt.cache(),
		}
		config.Config.AssetCacheEnabled = tt.args.cacheEnabled
		machineName, err := s.GetMachineNameByEndpointID(context.TODO(),tt.args.partnerID, tt.args.endpointID)
		ctrl.Finish()

		if !tt.wantErr {
			gomega.Expect(err).To(gomega.BeNil(), fmt.Sprintf("GetMachineNameByEndpointID() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			gomega.Expect(machineName).To(gomega.Equal(tt.wantMachineName), fmt.Sprintf("GetMachineNameByEndpointID() name = %s, gotResult = %v, want %v", tt.name, machineName, tt.wantMachineName))
		} else {
			gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("GetMachineNameByEndpointID() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}

	}
}

func Test_getDataFromAssetMS(t *testing.T) {
	var ctrl *gomock.Controller
	gomega.RegisterTestingT(t)
	logger.Load(config.Config.Log)

	type args struct {
		partnerID    string
		endpointID   gocql.UUID
		cacheEnabled bool
	}
	tests := []struct {
		name string
		cli  func() integration.HTTPClient
		args args
	}{
		{
			name: "Failed to read and close response body",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: &rwFail{},
				}, nil)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: false,
			},
		},
		{
			name: "Failed to unmarshal response",
			cli: func() integration.HTTPClient {
				c := common.NewMockHTTPClient(ctrl)
				c.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: &rwEmpty{},
				}, nil)
				return c
			},
			args: args{
				partnerID:    "test",
				endpointID:   gocql.UUID{},
				cacheEnabled: true,
			},
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		s := &Service{
			httpClient: tt.cli(),
		}
		config.Config.AssetCacheEnabled = tt.args.cacheEnabled
		_, err := s.getDataFromAssetMS(context.Background(), tt.args.partnerID, tt.args.endpointID)
		ctrl.Finish()

		gomega.Expect(err).ToNot(gomega.BeNil(), fmt.Sprintf("getDataFromAssetMS() name = %s, error = %v, wantErr %v", tt.name, err, true))
	}
}

func Test_IsAllResources(t *testing.T) {
	gomega.RegisterTestingT(t)
	gomega.Expect(integration.Desktop.IsAllResources()).To(gomega.BeFalse())
	gomega.Expect(integration.Server.IsAllResources()).To(gomega.BeFalse())
	gomega.Expect(integration.ResourceType("test").IsAllResources()).To(gomega.BeFalse())
	gomega.Expect(integration.ResourceType("").IsAllResources()).To(gomega.BeTrue())
}

type rw struct{}

func (*rw) Read(p []byte) (n int, err error) {
	s := apiModel.AssetCollection{
		SiteID:       "testSite",
		ClientID:     "testClient",
		EndpointType: "testType",
		System: &apiModel.AssetSystem{
			SystemName: "testName",
		},
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
