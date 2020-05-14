package models_test

import (
	"context"
	"errors"
	"testing"

	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func AnyN(n int) []interface{} {
	m := make([]interface{}, n)
	for i := 0; i < n; i++ {
		m[i] = gomock.Any()
	}

	return m
}

func TestExecutionResultRepoCassandra_DeleteBatch(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		batchMock   = mocks_cassandra.NewMockIBatch(ctrl)
		repo        = models.ExecutionResultRepoCassandra{}
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock)
		batchMock.EXPECT().Query(gomock.Any(), AnyN(2)...).Times(2)
		sessionMock.EXPECT().ExecuteBatch(batchMock).Return(nil)
		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock)

		err := repo.DeleteBatch(context.Background(), []models.ExecutionResult{{}})
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock)
		batchMock.EXPECT().Query(gomock.Any(), AnyN(2)...).Times(2)
		sessionMock.EXPECT().ExecuteBatch(batchMock).Return(errors.New(""))

		err := repo.DeleteBatch(context.Background(), []models.ExecutionResult{{}})
		if err == nil {
			t.Fatalf("error should't be nil")
		}
	})
}

func TestExecutionResultRepoCassandra_Upsert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		batchMock   = mocks_cassandra.NewMockIBatch(ctrl)
		repo        = models.ExecutionResultRepoCassandra{}

		partnerID = "pid"
		taskName  = "task"
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("full", func(t *testing.T) {
		err1, err2 := errors.New("err1"), errors.New("err2")

		sessionMock.EXPECT().NewBatch(gocql.UnloggedBatch).Return(batchMock)
		batchMock.EXPECT().Query(gomock.Any(), AnyN(7)...)
		batchMock.EXPECT().Query(gomock.Any(), AnyN(7)...)
		batchMock.EXPECT().Query(gomock.Any(), AnyN(5)...)
		sessionMock.EXPECT().ExecuteBatch(batchMock).Return(err1)
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(7)...).Return(queryMock)
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(7)...).Return(queryMock)
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(5)...).Return(queryMock)
		queryMock.EXPECT().Exec().Return(err2).Times(3)

		err := repo.Upsert(context.Background(), partnerID, taskName, models.ExecutionResult{})
		if err == nil {
			t.Fatal("error shouldn't be nil")
		}
	})
}

func TestExecutionResultRepoCassandra_GetByTargetAndTaskInstanceIDs(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		repo        = models.ExecutionResultRepoCassandra{}

		targetID, _   = gocql.RandomUUID()
		instanceID, _ = gocql.RandomUUID()
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := repo.GetByTargetAndTaskInstanceIDs(targetID, instanceID)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestExecutionResultRepoCassandra_GetByTaskInstanceIDs(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		repo        = models.ExecutionResultRepoCassandra{}

		instanceID, _ = gocql.RandomUUID()
		instanceIDs   = []gocql.UUID{instanceID}
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := repo.GetByTaskInstanceIDs(instanceIDs)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := repo.GetByTaskInstanceIDs(instanceIDs)
		if err == nil {
			t.Fatal("error cannot be nil")
		}
	})
}
