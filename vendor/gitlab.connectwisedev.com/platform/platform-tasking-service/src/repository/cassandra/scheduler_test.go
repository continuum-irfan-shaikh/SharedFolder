package cassandra

import (
	"errors"
	"testing"
	"time"

	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"github.com/golang/mock/gomock"
)

func TestScheduler_UpdateLastExpiredExecutionCheck(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		sch         = NewScheduler(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := sch.UpdateLastExpiredExecutionCheck(time.Now())
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := sch.UpdateLastExpiredExecutionCheck(time.Now())
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestScheduler_GetLastExpiredExecutionCheck(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		sch         = NewScheduler(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any())

		_, err := sch.GetLastExpiredExecutionCheck()
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := sch.GetLastExpiredExecutionCheck()
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestScheduler_UpdateScheduler(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		sch         = NewScheduler(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := sch.UpdateScheduler(time.Now())
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := sch.UpdateScheduler(time.Now())
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestScheduler_GetLastUpdate(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		sch         = NewScheduler(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any())

		_, err := sch.GetLastUpdate()
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := sch.GetLastUpdate()
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}
