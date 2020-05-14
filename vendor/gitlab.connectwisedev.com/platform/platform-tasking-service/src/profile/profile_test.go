package profile_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/profile"
)

const (
	defaultMsg = `failed on unexpected value of result "%v"`
)

func TestHandler_ServeHTTP(t *testing.T) {
	RegisterTestingT(t)

	type payload struct {
		handlerName string
		url         string
	}

	type expected struct {
		statusCode int
		headers    map[string]string
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "Status not found",
			payload: payload{
				handlerName: "abc",
				url:         "/debug/profile",
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				headers: map[string]string{
					"X-Content-Type-Options": "nosniff",
					"Content-Type":           "text/plain; charset=utf-8",
					"X-Go-Pprof":             "1",
				},
			},
		},
		{
			name: "Status ok for heap",
			payload: payload{
				handlerName: "heap",
				url:         "/debug/profile/heap?gc=1",
			},
			expected: expected{
				statusCode: http.StatusOK,
				headers: map[string]string{
					"X-Content-Type-Options": "nosniff",
					"Content-Type":           "text/plain; charset=utf-8",
				},
			},
		},
		{
			name: "Status ok for goroutine",
			payload: payload{
				handlerName: "goroutine",
				url:         "/debug/profile/goroutine",
			},
			expected: expected{
				statusCode: http.StatusOK,
				headers: map[string]string{
					"X-Content-Type-Options": "nosniff",
					"Content-Type":           "text/plain; charset=utf-8",
				},
			},
		},
	}

	for _, test := range tc {
		h := profile.Handler(test.payload.handlerName)
		r, err := http.NewRequest(http.MethodGet, test.payload.url, nil)
		Ω(err).To(BeNil(), "failed on unexpected error with request")
		w := &httptest.ResponseRecorder{}

		h.ServeHTTP(w, r)
		Ω(w.Code).To(Equal(test.expected.statusCode), fmt.Sprintf(defaultMsg, test.name))
		for k, v := range test.expected.headers {
			Ω(w.Header().Get(k)).To(Equal(v), fmt.Sprintf(defaultMsg, test.name))
		}
	}
}
