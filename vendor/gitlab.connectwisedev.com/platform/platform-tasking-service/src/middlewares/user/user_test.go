package user

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

const defaultMsg = `failed on unexpected value of result "%v"`

func TestUser_ServeHTTP(t *testing.T) {
	RegisterTestingT(t)
	config.Config.DefaultLanguage = "en-US"
	if err := translation.Load(); err != nil {
		t.Error("can't load translator")
	}
	if err := logger.Load(config.Config.Log); err != nil {
		t.Error("can't load logger")
	}

	var check bool
	var ctx context.Context

	type payload struct {
		params  map[string]string
		headers map[string]string
		next    func(rw http.ResponseWriter, r *http.Request)
		check   bool
	}

	type expected struct {
		user entities.User
		body string
		code int
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "success: get user",
			payload: payload{
				params: map[string]string{"partnerID": "50016646"},
				headers: map[string]string{
					"uid":   "uid",
					"realm": "/activedirectory",
				},
				next: func(rw http.ResponseWriter, r *http.Request) {
					check = true
					ctx = r.Context()
				},
				check: true,
			},
			expected: expected{
				user: entities.User{
					PartnerID:   "50016646",
					UID:         "uid",
					IsNOCAccess: true,
				},
			},
		},
		{
			name: "error: can't validate user (partnerID)",
			payload: payload{
				params: map[string]string{},
				headers: map[string]string{
					"uid":             "uid",
					"Accept-Language": "en-US",
				},
				check: false,
			},
			expected: expected{
				user: entities.User{
					PartnerID: "50016646",
					UID:       "uid",
				},
				code: http.StatusBadRequest,
				body: `{"message":"Can't validate user parameters","errorCode":"error_cant_validate_user_parameters"}`,
			},
		},
		{
			name: "error: can't validate user (UID)",
			payload: payload{
				params: map[string]string{"partnerID": "50016646"},
				headers: map[string]string{
					"Accept-Language": "en-US",
				},
				check: false,
			},
			expected: expected{
				user: entities.User{
					PartnerID: "50016646",
					UID:       "uid",
				},
				code: http.StatusBadRequest,
				body: `{"message":"Can't validate user parameters","errorCode":"error_cant_validate_user_parameters"}`,
			},
		},
	}

	for _, test := range tc {
		check = false
		ctx = context.Background()

		r := httptest.NewRequest("GET", "http://someurl.com/tasking/v1/partners/50016646/tasks", nil)
		r = mux.SetURLVars(r, test.params)
		for k, v := range test.payload.headers {
			r.Header.Add(k, v)
		}
		rw := httptest.NewRecorder()
		md := NewUser(logger.Log)
		md.ServeHTTP(rw, r, test.next)

		Ω(check).To(Equal(test.check), fmt.Sprintf(defaultMsg, "check next md"))
		if test.check {
			actual := ctx.Value(config.UserKeyCTX).(entities.User)
			Ω(actual).To(Equal(test.expected.user), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(rw.Code).To(Equal(test.code), fmt.Sprintf(defaultMsg, test.name))
		Ω(rw.Body.String()).To(Equal(test.expected.body), fmt.Sprintf(defaultMsg, test.name))
	}
}
