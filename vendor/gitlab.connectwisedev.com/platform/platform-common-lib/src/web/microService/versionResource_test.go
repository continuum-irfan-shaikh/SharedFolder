// Package microservice implements the some common resources.
//
// Deprecated: microservice is old implementation of common-resources and should not be used
// except for compatibility with legacy systems.
//
// Use src/web/rest for all common handlers
// This package is frozen and no new functionality will be added.
package microservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"reflect"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/services"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/services/model"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/web"
	cwmock "gitlab.connectwisedev.com/platform/platform-common-lib/src/web/mock"
	"github.com/golang/mock/gomock"
)

type MockVersionDependenciesImpl struct {
	services.VersionFactoryImpl
}

func TestVersionGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listerMock := MockVersionDependenciesImpl{}
	mockRsp := httptest.NewRecorder()
	mockRc := cwmock.NewMockRequestContext(ctrl)
	mockRc.EXPECT().GetResponse().Return(mockRsp).AnyTimes()

	versionResource{
		f: listerMock,
		version: model.Version{
			SolutionName:    "SolutionName",
			ServiceName:     "ServiceName",
			ServiceProvider: "ContinuumLLC",
			Major:           "1",
			Minor:           "1",
			Patch:           "11",
		},
	}.Get(mockRc)

	if mockRsp.Code != http.StatusOK {
		t.Errorf("Unexpected error code : %v", mockRsp.Code)
	}
}

func TestCreateVersionRouteConfig(t *testing.T) {
	version := model.Version{
		SolutionName:    "SolutionName",
		ServiceName:     "ServiceName",
		ServiceProvider: "ContinuumLLC",
		Major:           "1",
		Minor:           "1",
		Patch:           "11",
	}

	listerMock := MockVersionDependenciesImpl{}
	route := web.RouteConfig{
		URLPathSuffix: "/version",
		URLVars:       []string{},
		Res: versionResource{
			version: version,
			f:       listerMock,
		},
	}

	rout := CreateVersionRouteConfig(version, listerMock)

	if reflect.DeepEqual(route, rout) {
		t.Errorf("Expected same but got Different %v : %v", route, rout)
	}
}
