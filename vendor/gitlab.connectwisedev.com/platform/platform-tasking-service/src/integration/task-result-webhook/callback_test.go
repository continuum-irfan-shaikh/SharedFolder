package taskResultWebhook

import (
	"context"
	"net/http"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"fmt"
)

func init() {
	config.Load()
	logger.Load(config.Config.Log)
}

func TestCallWebhooksFor(t *testing.T) {

	var tests = []struct {
		name        string
		callBackURL string
		statusCode  int
		needError   bool
	}{
		{
			name:        "test 1",
			callBackURL: "http://localhost:8181/task1",
			statusCode:  http.StatusOK,
		},
		{
			name:        "test 2",
			callBackURL: "http://localhost:8182/something/else",
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:        "test 3",
			callBackURL: "http://localhost:8183/something/else",
			statusCode:  http.StatusInternalServerError,
			needError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			registerResponder(test.statusCode, test.callBackURL, test.needError, t)
			taskRes := TaskResult{ResultWebhook: test.callBackURL}
			CallWebhooksFor(context.Background(), []TaskResult{taskRes})
		})
	}
}

func registerResponder(status int, url string, needError bool, t *testing.T) {

	t.Logf("Registered HTTP responder on URL: %s", url)
	httpmock.RegisterResponder(http.MethodPost, url,
		func(req *http.Request) (*http.Response, error) {
			var err error
			if needError {
				err = fmt.Errorf("Http request error")
			}
			return httpmock.NewStringResponse(status, "ok"), err
		},
	)
}
