package permission

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

const defaultMsg = `failed on unexpected value of result "%v"`

func TestPermission_ServeHTTP(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	config.Config.DefaultLanguage = "en-US"
	if err := translation.Load(); err != nil {
		t.Error("can't load translator")
	}
	if err := logger.Load(config.Config.Log); err != nil {
		t.Error("can't load logger")
	}

	var check bool
	ctx := context.Background()
	userNOCCTX := context.WithValue(ctx, config.UserKeyCTX, entities.User{IsNOCAccess: true})
	userCTX := context.WithValue(ctx, config.UserKeyCTX, entities.User{})
	var ctxResult context.Context

	type payload struct {
		ctx   context.Context
		next  func(rw http.ResponseWriter, r *http.Request)
		check bool
		uc    func() userUC
	}

	type expected struct {
		endpoints []string
		body      string
		code      int
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "err: can't validate user",
			payload: payload{
				ctx:   ctx,
				check: false,
				uc: func() userUC {
					return nil
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"message":"Can't get user","errorCode":"error_cant_get_user"}`,
			},
		},
		{
			name: "success: noc user",
			payload: payload{
				ctx:   userNOCCTX,
				check: true,
				uc: func() userUC {
					return nil
				},
				next: func(rw http.ResponseWriter, r *http.Request) {
					check = true
				},
			},
			expected: expected{
				endpoints: []string{},
			},
		},
		{
			name: "err: can't get user sites",
			payload: payload{
				ctx:   userCTX,
				check: false,
				uc: func() userUC {
					err := errors.New("can't get sites")
					ucMock := NewMockuserUC(mockCtrl)
					ucMock.EXPECT().Sites(userCTX).Return(nil, err).Times(1)
					return ucMock
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"message":"Can't get user sites","errorCode":"error_cant_get_user_sites"}`,
			},
		},
		{
			name: "err: can't get user endpoints",
			payload: payload{
				ctx:   userCTX,
				check: false,
				uc: func() userUC {
					err := errors.New("can't get endpoints")
					sites := []string{"1", "2", "3"}
					ucMock := NewMockuserUC(mockCtrl)
					ucMock.EXPECT().Sites(userCTX).Return(sites, nil).Times(1)
					ucMock.EXPECT().Endpoints(userCTX, sites).Return(nil, err).Times(1)
					return ucMock
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"message":"Can't get user end points by sites","errorCode":"error_cant_get_user_end_points_by_sites"}`,
			},
		},
		{
			name: "success: get endpoints by user",
			payload: payload{
				ctx:   userCTX,
				check: true,
				uc: func() userUC {
					sites := []string{"1", "2", "3"}
					endpoints := []string{"1", "2", "3"}
					ucMock := NewMockuserUC(mockCtrl)
					ucMock.EXPECT().Sites(userCTX).Return(sites, nil).Times(1)
					ucMock.EXPECT().Endpoints(userCTX, sites).Return(endpoints, nil).Times(1)
					return ucMock
				},
				next: func(rw http.ResponseWriter, r *http.Request) {
					check = true
					ctxResult = r.Context()
				},
			},
			expected: expected{
				endpoints: []string{"1", "2", "3"},
			},
		},
	}

	for _, test := range tc {
		check = false
		ctxResult = context.Background()

		r := httptest.NewRequest("GET", "http://someurl.com/tasking/v1/partners/50016646/tasks", nil)
		r = r.WithContext(test.payload.ctx)
		rw := httptest.NewRecorder()
		md := NewPermission(test.payload.uc(), logger.Log)
		md.ServeHTTP(rw, r, test.next)

		Ω(check).To(Equal(test.check), fmt.Sprintf(defaultMsg, "check next md"))
		if test.check {
			if len(test.expected.endpoints) == 0 {
				continue
			}
			actual := ctxResult.Value(config.UserEndPointsKeyCTX).([]string)
			Ω(actual).To(Equal(test.expected.endpoints), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(rw.Code).To(Equal(test.code), fmt.Sprintf(defaultMsg, test.name))
		Ω(rw.Body.String()).To(Equal(test.expected.body), fmt.Sprintf(defaultMsg, test.name))
	}
}
