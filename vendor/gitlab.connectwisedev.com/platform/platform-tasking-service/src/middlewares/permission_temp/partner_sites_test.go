package permission_temp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mockLoggerTasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
)

const defaultMsg = `failed on unexpected value of result "%v"`

func init() {
	logger.Load(config.Config.Log)
}

func TestPermission_PartnerSitesCheck(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller
	var server *httptest.Server

	type payload struct {
		logger      func() logger.Logger
		userService func() user.Service
		httpClient  func() *http.Client
		cache       func() persistency.Cache
		loadConfig  func()
	}

	tc := []struct {
		name     string
		expected error
		payload
	}{
		{
			name: "Success_#1",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					return &http.Client{}
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					b := []byte(`["someSiteID"]`)
					cacheMock.EXPECT().Get(key).Return(b, nil).Times(1)

					partnerKey := []byte("TKS_SITES_BY_PARTNER_partnerID")
					bytes := []byte(`["someSiteID"]`)
					cacheMock.EXPECT().Get(partnerKey).Return(bytes, nil).Times(1)

					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
				},
			},
		},
		{
			name: "Success_#2",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`{"outdata":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
				},
			},
		},
		{
			name: "Success_#3",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					cacheMock.EXPECT().Get(key).Return(nil, errors.New("some_error")).Times(1)

					cacheMock.EXPECT().Set(key, []byte(`["123"]`), 60).Return(nil).Times(1)

					partnerKey := []byte("TKS_SITES_BY_PARTNER_partnerID")
					bytes := []byte(`["123"]`)
					cacheMock.EXPECT().Get(partnerKey).Return(bytes, nil).Times(1)

					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
					config.Config.SitesMsURL = server.URL
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
		},
		{
			name: "error while creating request",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					cacheMock.EXPECT().Get(key).Return(nil, errors.New("some_error")).Times(1)

					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
					config.Config.SitesMsURL = ":http://"
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
			expected: errors.New(`parse ":http:///partner/partnerID/sites": missing protocol scheme`),
		},
		{
			name: "error while performing request",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					cacheMock.EXPECT().Get(key).Return(nil, errors.New("some_error")).Times(1)

					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
					config.Config.SitesMsURL = ""
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
			expected: errors.New(`Get "/partner/partnerID/sites": unsupported protocol scheme ""`),
		},
		{
			name: "error while unmarshal",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`}"siteDetailList":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					cacheMock.EXPECT().Get(key).Return(nil, errors.New("some_error")).Times(1)

					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
					config.Config.SitesMsURL = server.URL
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
			expected: errors.New("invalid character '}' looking for beginning of value"),
		},
		{
			name: "error while setting to cache",
			payload: payload{
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`{"outdata":[{"siteId":123}]}`)

					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}

						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					cacheMock.EXPECT().Get(key).Return(nil, errors.New("some_error")).Times(1)

					cacheMock.EXPECT().Set(key, []byte(`["123"]`), 60).Return(errors.New("some_error")).
						Times(1)

					partnerKey := []byte("TKS_SITES_BY_PARTNER_partnerID")
					cacheMock.EXPECT().Get(partnerKey).Return(nil, errors.New("some_error")).Times(1)
					cacheMock.EXPECT().Set(partnerKey, []byte(`["123"]`), 60).Return(errors.New("some_error")).
						Times(1)

					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
		},
		{
			name: "error while getting user sites",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)

					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = ""
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
			expected: errors.New(`Get "/partner/partnerID/sites": unsupported protocol scheme ""`),
		},
		{
			name: "error while getting from cache",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`{"outdata":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					cacheMock.EXPECT().Get(key).Return([]byte("{}"), nil).Times(1)
					cacheMock.EXPECT().Set(key, []byte(`["123"]`), 60).Return(nil).
						Times(1)

					partnerKey := []byte("TKS_SITES_BY_PARTNER_partnerID")
					cacheMock.EXPECT().Get(partnerKey).Return([]byte("{}"), nil).Times(1)
					cacheMock.EXPECT().Set(partnerKey, []byte(`["123"]`), 60).Return(nil).
						Times(1)
					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
		},
		{
			name: "error with new request",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`{"outdata":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = ":http://"
				},
			},
			expected: errors.New(`parse ":http:///partners/partnerID/sites?Operation=activesites": missing protocol scheme`),
		},
		{
			name: "error with performing request",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`{"outdata":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = ""
				},
			},
			expected: errors.New(`Get "/partners/partnerID/sites?Operation=activesites": unsupported protocol scheme ""`),
		},
		{
			name: "error with performing request",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`}"outdata":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
				},
			},
			expected: errors.New("invalid character '}' looking for beginning of value"),
		},
		{
			name: "error while getting partner site",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`}"outdata":[{"siteId":123}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					key := []byte("TKS_SITES_BY_PARTNER_partnerID_USER_uid")
					partnerKey := []byte("TKS_SITES_BY_PARTNER_partnerID")

					cacheMock.EXPECT().Get(key).Return(nil, errors.New("some_error")).Times(1)
					cacheMock.EXPECT().Set(key, []byte(`["123"]`), 60).Return(nil).
						Times(1)
					cacheMock.EXPECT().Get(partnerKey).Return(nil, errors.New("some_error")).Times(1)

					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = true
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = ":http://"
					config.Config.PartnerSitesCacheExpiration = 60
				},
			},
			expected: errors.New(`parse ":http:///partners/partnerID/sites?Operation=activesites": missing protocol scheme`),
		},
		{
			name: "error with wrong number of sites",
			payload: payload{
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					userSites := []byte(`{"siteDetailList":[{"siteId":123}]}`)
					partnerSites := []byte(`{"outdata":[{"siteId":123},{"siteId":321}]}`)
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.String() == "/partner/partnerID/sites" {
							_, err := w.Write(userSites)
							Ω(err).To(BeNil())
							return
						}
						if r.URL.String() == "/partners/partnerID/sites?Operation=activesites" {
							_, err := w.Write(partnerSites)
							Ω(err).To(BeNil())
							return
						}
					}))
					return server.Client()
				},
				cache: func() persistency.Cache {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					return cacheMock
				},
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
				},
			},
			expected: errors.New("number of user sites != partner sites for user name and partner partnerID"),
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		md := Permission{
			log:         logger.Log,
			userService: test.payload.userService(),
			httpClient:  test.payload.httpClient(),
			cache:       test.payload.cache(),
		}
		test.loadConfig()

		u := md.userService.GetUser(nil, nil)
		err := md.PartnerSitesCheck(context.TODO(), u)
		mockCtrlr.Finish()

		if test.expected == nil {
			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
		} else {
			Ω(err.Error()).To(Equal(test.expected.Error()), fmt.Sprintf(defaultMsg, test.name))
		}

	}
}
