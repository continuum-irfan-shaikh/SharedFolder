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
	"reflect"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/services/model"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/web"
	cwmock "gitlab.connectwisedev.com/platform/platform-common-lib/src/web/mock"
	"github.com/golang/mock/gomock"
)

func TestDownloadGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger.Update(logger.Config{Destination: logger.DISCARD})
	mockRsp := httptest.NewRecorder()
	mockRc := cwmock.NewMockRequestContext(ctrl)
	mockRc.EXPECT().GetResponse().Return(mockRsp).AnyTimes()

	downloadResource{
		fileInfo: model.DownloadFileInfo{
			Name:        "Swagger.yaml",
			Path:        "abcd",
			ContentType: "application/octet-stream",
		},
	}.Get(mockRc)

	if mockRsp.Code == http.StatusOK {
		t.Errorf("Unexpected error code : %v", mockRsp.Code)
	}
}

func TestCreateDownloadRouteConfig(t *testing.T) {
	logger.Update(logger.Config{Destination: logger.DISCARD})
	fileInfo := model.DownloadFileInfo{
		Name:        "Swagger.yaml",
		Path:        "abcd",
		ContentType: "application/octet-stream",
	}

	route := web.RouteConfig{
		URLPathSuffix: "/download",
		URLVars:       []string{},
		Res: downloadResource{
			fileInfo: fileInfo,
		},
	}

	rout := CreateDownloadRouteConfig(fileInfo, "/download")

	if reflect.DeepEqual(route, rout) {
		t.Errorf("Expected same but got Different %v : %v", route, rout)
	}
}
