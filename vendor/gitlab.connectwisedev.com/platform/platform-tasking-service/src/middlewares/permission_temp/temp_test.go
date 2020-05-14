package permission_temp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mockLoggerTasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestNew(t *testing.T) {
	r := RoleSuperUser
	r = RoleAdmin
	r = RolePrimarySuperUser
	r = RoleManager
	r = RoleTechnician
	r = userSitesKeyPattern
	r = partnerSitesKeyPattern
	r = partnerSitesURLPattern
	r = userSitesURLPattern
	r = tokenHeader
	fmt.Println(r)

	NewPermissionMiddleware(nil, nil, nil, nil)
}

func TestPermission_ServeHTTP(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller
	var server *httptest.Server

	type payload struct {
		role        string
		logger      func() logger.Logger
		userService func() user.Service
		httpClient  func() *http.Client
		loadConfig  func()
	}

	type expected struct {
		code int
		body string
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "Success_#1",
			payload: payload{
				role: "someRole",
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", true)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
					return server.Client()
				},
				loadConfig: func() {},
			},
			expected: expected{
				code: 200,
				body: "Success",
			},
		},
		{
			name: "Success_#2",
			payload: payload{
				role: RoleAdmin,
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", false)
					return userServiceMock
				},
				httpClient: func() *http.Client {
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
					return server.Client()
				},
				loadConfig: func() {},
			},
			expected: expected{
				code: 200,
				body: "Success",
			},
		},
		{
			name: "Success_#3",
			payload: payload{
				role: "someRole",
				logger: func() logger.Logger {
					loggerMock := mockLoggerTasking.NewMockLogger(mockCtrlr)
					return loggerMock
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", false)
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
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
				},
			},
			expected: expected{
				code: 200,
				body: "Success",
			},
		},
		{
			name: "unauthorized",
			payload: payload{
				role: "someRole",
				logger: func() logger.Logger {
					return nil
				},
				userService: func() user.Service {
					userServiceMock := user.NewMock("name", "partnerID", "uid", "token", false)
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
				loadConfig: func() {
					config.Config.AssetCacheEnabled = false
					config.Config.SitesMsURL = server.URL
					config.Config.SitesNoTokenURL = server.URL
				},
			},
			expected: expected{
				code: 401,
				body: `{"message":"Forbidden request: role - \"someRole\"","errorCode":"Forbidden request: role - \"someRole\""}`,
			},
		},
	}

	logger.Load(config.Config.Log)
	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		md := Permission{
			log:         logger.Log,
			userService: test.payload.userService(),
			httpClient:  test.payload.httpClient(),
		}
		test.loadConfig()
		err := translation.Load()
		Ω(err).To(BeNil())

		req, err := http.NewRequest(http.MethodPost, "/tasking/v1/partners/123456/tasks", nil)
		Ω(err).To(BeNil())
		req.Header.Set("role", test.payload.role)
		req.Header.Set("Accept-Language", "en-US")

		w := httptest.NewRecorder()

		md.ServeHTTP(w, req, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Success"))
			Ω(err).To(BeNil())
		})

		mockCtrlr.Finish()

		Ω(w.Body.String()).To(Equal(test.expected.body), fmt.Sprintf(defaultMsg, test.name))
		Ω(w.Code).To(Equal(test.expected.code), fmt.Sprintf(defaultMsg, test.name))
	}
}
