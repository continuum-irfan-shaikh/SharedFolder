package taskCounterCassandra

import (
	"context"
	"errors"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func init() {
	_ = logger.Load(config.Config.Log)
}

func monkeyPatchCassandraSession(s cassandra.ISession) {
	cassandra.Session = s
}

func TestTaskCounterCassandra_GetAllPartners(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		target      = New(2)
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock).Times(2)
		queryMock.EXPECT().Iter().Return(iterMock).Times(2)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false)
		iterMock.EXPECT().Close().Return(nil).Times(2)

		_, err := target.GetAllPartners(context.Background())
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock).Times(2)
		queryMock.EXPECT().Iter().Return(iterMock).Times(2)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := target.GetAllPartners(context.Background())
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}

func TestTaskCounterCassandra_DecreaseCounter(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		batchMock   = mocks_cassandra.NewMockIBatch(ctrl)
		target      = New(2)

		counters  = []models.TaskCount{{}}
		partnerID = "pid"
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock).Times(2)
		batchMock.EXPECT().Query(gomock.Any(), gomock.Any(), partnerID, gomock.Any())
		sessionMock.EXPECT().ExecuteBatch(batchMock).Return(nil)

		err := target.DecreaseCounter(partnerID, counters, true)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock).Times(1)
		batchMock.EXPECT().Query(gomock.Any(), gomock.Any(), partnerID, gomock.Any())
		sessionMock.EXPECT().ExecuteBatch(batchMock).Return(errors.New(""))

		err := target.DecreaseCounter(partnerID, counters, false)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}

func TestTaskCounterCassandra_IncreaseCounter(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		batchMock   = mocks_cassandra.NewMockIBatch(ctrl)
		target      = New(2)

		counters  = []models.TaskCount{{}}
		partnerID = "pid"
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock).Times(2)
		batchMock.EXPECT().Query(gomock.Any(), gomock.Any(), partnerID, gomock.Any())
		sessionMock.EXPECT().ExecuteBatch(batchMock).Return(nil)

		err := target.IncreaseCounter(partnerID, counters, true)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock).Times(1)
		batchMock.EXPECT().Query(gomock.Any(), gomock.Any(), partnerID, gomock.Any())
		sessionMock.EXPECT().ExecuteBatch(batchMock).Return(errors.New(""))

		err := target.IncreaseCounter(partnerID, counters, false)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}

func TestTaskCounterCassandra_GetCounters(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		target      = New(2)

		partnerID     = "pid"
		endpointID, e = gocql.RandomUUID()
		empty         = gocql.UUID{}
	)

	if e != nil {
		t.Fatalf(e.Error())
	}

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(nil, nil, gomock.Any()).Return(nil)

		_, err := target.GetCounters(context.Background(), partnerID, endpointID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("positive_2", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), nil, gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), nil, gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := target.GetCounters(context.Background(), partnerID, empty)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(nil, nil, gomock.Any()).Return(errors.New(""))

		_, err := target.GetCounters(context.Background(), partnerID, endpointID)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), nil, gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), nil, gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := target.GetCounters(context.Background(), partnerID, empty)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}
