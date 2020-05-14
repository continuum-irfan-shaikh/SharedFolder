package models_test

import (
	"context"
	"errors"
	"testing"
	"time"

	cas "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func MonkeyPatchCassandraSession(s cassandra.ISession) {
	cassandra.Session = s
}

func TestExecutionExpirationRepoCassandra_Delete(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = cas.NewMockISession(ctrl)
		queryMock   = cas.NewMockIQuery(ctrl)
		repo        = models.ExecutionExpirationRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)

		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := repo.Delete(models.ExecutionExpiration{})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestExecutionExpirationRepoCassandra_GetByTaskInstanceIDs(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = cas.NewMockISession(ctrl)
		queryMock   = cas.NewMockIQuery(ctrl)
		iterMock    = cas.NewMockIIter(ctrl)
		repo        = models.ExecutionExpirationRepoCassandra{}

		partnerID     = "pid"
		instanceID, _ = gocql.RandomUUID()
		instanceIDs   = []gocql.UUID{instanceID}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)

		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := repo.GetByTaskInstanceIDs(partnerID, instanceIDs)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)

		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := repo.GetByTaskInstanceIDs(partnerID, instanceIDs)
		if err == nil {
			t.Fatal("error cannot be nil")
		}
	})
}

func TestExecutionExpirationRepoCassandra_GetByExpirationTime(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = cas.NewMockISession(ctrl)
		queryMock   = cas.NewMockIQuery(ctrl)
		iterMock    = cas.NewMockIIter(ctrl)
		repo        = models.ExecutionExpirationRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)

		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := repo.GetByExpirationTime(context.Background(), time.Now())
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)

		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := repo.GetByExpirationTime(context.Background(), time.Now())
		if err == nil {
			t.Fatal("error cannot be nil")
		}
	})
}

func TestExecutionExpirationRepoCassandra_InsertExecutionExpiration(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = cas.NewMockISession(ctrl)
		queryMock   = cas.NewMockIQuery(ctrl)
		repo        = models.ExecutionExpirationRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)

		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := repo.InsertExecutionExpiration(context.Background(), models.ExecutionExpiration{})
		if err != nil {
			t.Fatal(err)
		}
	})
}
