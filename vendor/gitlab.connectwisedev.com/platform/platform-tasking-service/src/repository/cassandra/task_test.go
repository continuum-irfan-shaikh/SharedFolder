package cassandra

import (
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

var partnerID = strconv.Itoa(rand.Int())

func TestNewTask(t *testing.T) {
	ti := NewTask(mocks_cassandra.NewMockISession(gomock.NewController(t)))
	if reflect.TypeOf(ti) != reflect.TypeOf(&Task{}) {
		t.Fatalf("NewTask didn't create *Task")
	}
}

func TestTask_GetName(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		task        = Task{Conn: sessionMock}
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Iter().Return(iterMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)
		queryMock.EXPECT().Release()

		_, err := task.GetName(partnerID, "")
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Iter().Return(iterMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))
		queryMock.EXPECT().Release()

		_, err := task.GetName(partnerID, "")
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}

func TestTask_GetNext(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		task        = Task{Conn: sessionMock}
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Iter().Return(iterMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)
		queryMock.EXPECT().Release()

		_, err := task.GetNext(partnerID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Iter().Return(iterMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))
		queryMock.EXPECT().Release()

		_, err := task.GetNext(partnerID)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}

func TestTask_GetScheduledTasks(t *testing.T) {
	const scheduleStringIndex = 9
	var ctrl *gomock.Controller
	var sessionMock *mocks_cassandra.MockISession
	var queryMock *mocks_cassandra.MockIQuery
	var iterMock *mocks_cassandra.MockIIter
	var task Task

	t.Run("positive", func(t *testing.T) {
		ctrl = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock = mocks_cassandra.NewMockIIter(ctrl)
		task = Task{sessionMock}

		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)

		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			schedule := apiModels.Schedule{
				Regularity: 1,
				Repeat: apiModels.Repeat{
					RunTime: time.Now().Add(time.Hour),
				},
			}

			b, err := json.Marshal(schedule)
			if err != nil {
				t.Fatalf(err.Error())
			}

			// set &scheduleString to pass validation
			*(p[scheduleStringIndex].(*string)) = string(b)
		}).Return(nil).Times(1)

		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)
		_, err := task.GetScheduledTasks(partnerID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		ctrl = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock = mocks_cassandra.NewMockIIter(ctrl)
		task = Task{sessionMock}

		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)

		sessionMock.EXPECT().Query(gomock.Any(), partnerID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			// set &scheduleString to pass validation
			*(p[scheduleStringIndex].(*string)) = "invalid json"
		}).Return(nil)

		_, err := task.GetScheduledTasks(partnerID)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		ctrl = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock = mocks_cassandra.NewMockIIter(ctrl)
		task = Task{sessionMock}

		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)

		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			schedule := apiModels.Schedule{
				Regularity: 1,
				Repeat: apiModels.Repeat{
					RunTime: time.Now().Add(time.Hour),
				},
			}

			b, err := json.Marshal(schedule)
			if err != nil {
				t.Fatalf(err.Error())
			}

			// set &scheduleString to pass validation
			*(p[scheduleStringIndex].(*string)) = string(b)
		}).Return(nil).Times(1)

		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New("error"))
		_, err := task.GetScheduledTasks(partnerID)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}
