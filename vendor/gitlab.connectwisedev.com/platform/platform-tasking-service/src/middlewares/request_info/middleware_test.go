package request_info

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/urfave/negroni"

	rLog "gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
)

var (
	logString string
)

func init() {
	config.Load()
	config.Config.Log.FileName = "logs_test.log"
	config.Config.Log.LogLevel = rLog.DEBUG
	logger.Load(config.Config.Log)
}

func TestNewMiddlewareFromLogger(t *testing.T) {
	middlewareName := "name"
	middleware := NewMiddlewareFromLogger(middlewareName)

	if middleware.Name != middlewareName || middleware.Before == nil || middleware.After == nil {
		t.Fatalf("NewMiddlewareFromLogger() failed")
	}
}

func TestServeHTTP(t *testing.T) {
	if _, err := os.Create(config.Config.Log.FileName); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(config.Config.Log.FileName)
	middlewareEmpty := NewMiddlewareFromLogger("name")
	middlewareEmpty.Before = nil //not nil by default
	middlewareEmpty.After = nil  //not nil by default

	request := httptest.NewRequest("GET", "http://url", nil)
	request.Header.Set("X-Real-IP", "1.2.3.4")
	rw := negroni.NewResponseWriter(httptest.NewRecorder())
	var nextCalled bool
	middlewareEmpty.ServeHTTP(rw, request, func(_ http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	time.Sleep(time.Second)
	b, err := ioutil.ReadFile(config.Config.Log.FileName)
	if err != nil {
		t.Fatalf("Error while reading logfile: %v", err)
	}
	logString = string(b)

	if middlewareEmpty.Before == nil || middlewareEmpty.After == nil {
		t.Fatal("Before func or After func was not set")
	}

	if !nextCalled {
		t.Fatal("next() handler must be called")
	}
}
