package request_info

import (
	"net/http"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// MiddlewareFromLogger  struct add TransactionID instead of Request-Id
type MiddlewareFromLogger struct {
	negronilogrus.Middleware
}

//ServeHTTP where TransactionID is added
func (middleWare *MiddlewareFromLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if middleWare.Before == nil {
		middleWare.Before = negronilogrus.DefaultBefore
	}

	if middleWare.After == nil {
		middleWare.After = negronilogrus.DefaultAfter
	}

	for _, url := range middleWare.ExcludedURLs() {
		if r.URL.Path == url {
			return
		}
	}

	start := time.Now()

	// Try to get the real IP
	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		remoteAddr = realIP
	}

	msgStartedHandling := "started handling r method=%v remote=%v r=%v"
	logger.Log.InfofCtx(r.Context(), msgStartedHandling, r.Method, remoteAddr, r.RequestURI)

	next(rw, r)

	latency := time.Since(start)
	resStatus := rw.(negroni.ResponseWriter).Status()

	msgCompletedHandling := "completed handling r measure#%s.latency=%v method=%v remote=%v r=%v status=%v text_status-%v took=%v"
	v := []interface{}{
		middleWare.Name,
		latency.Nanoseconds(),
		r.Method,
		remoteAddr,
		r.URL.Path,
		resStatus,
		http.StatusText(resStatus),
		latency,
	}
	logger.Log.InfofCtx(r.Context(), msgCompletedHandling, v...)
}

//NewMiddlewareFromLogger method to create TransactionIdMiddleware
func NewMiddlewareFromLogger(name string) *MiddlewareFromLogger {
	var log = new(logrus.Logger)

	middleWare := &MiddlewareFromLogger{
		Middleware: negronilogrus.Middleware{
			Logger: log,
			Name:   name,
			Before: negronilogrus.DefaultBefore,
			After:  negronilogrus.DefaultAfter,
		},
	}
	middleWare.SetLogStarting(true)
	return middleWare
}
