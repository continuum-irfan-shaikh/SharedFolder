package models_test

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

func monkeyPatchCassandraSession(s cassandra.ISession) {
	cassandra.Session = s
}

func TestTaskDefinitionRepoCassandra_GetAllByPartnerID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		target      = models.TaskDefinitionRepoCassandra{}

		partnerID = "pid"
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		logger.Load(config.Config.Log)
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := target.GetAllByPartnerID(context.Background(), partnerID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := target.GetAllByPartnerID(context.Background(), partnerID)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			*(p[12].(*bool)) = true
		}).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := target.GetAllByPartnerID(context.Background(), partnerID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}

func TestTaskDefinitionRepoCassandra_CanBeUpdated(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		target      = models.TaskDefinitionRepoCassandra{}

		partnerID = "pid"
		name      = "name"
		id, _     = gocql.RandomUUID()
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p interface{}) {
			*(p.(*string)) = name
		}).Return(nil)
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(2)...).Do(func(p1 *string, p2 *gocql.UUID) {
			*p1 = name
			*p2 = id
		}).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(2)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := target.CanBeUpdated(context.Background(), partnerID, name, id)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := target.CanBeUpdated(context.Background(), partnerID, name, id)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		invalidID, _ := gocql.RandomUUID()
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p interface{}) {
			*(p.(*string)) = "not name"
		}).Return(nil)
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(2)...).Do(func(p1 *string, p2 *gocql.UUID) {
			*p1 = name
			*p2 = invalidID
		}).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(2)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		b, err := target.CanBeUpdated(context.Background(), partnerID, name, id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if b {
			t.Fatalf("result should be false")
		}
	})
}

func TestTaskDefinitionRepoCassandra_Exists(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		target      = models.TaskDefinitionRepoCassandra{}

		partnerID = "pid"
		name      = "name"
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p interface{}) {
			*(p.(*string)) = name
		}).Return(nil)

		b := target.Exists(context.Background(), partnerID, name)
		if !b {
			t.Fatalf("result should be true")
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		b := target.Exists(context.Background(), partnerID, name)
		if b {
			t.Fatalf("result should be false")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p interface{}) {
			*(p.(*string)) = "not name"
		}).Return(nil)

		b := target.Exists(context.Background(), partnerID, name)
		if b {
			t.Fatalf("result should be false")
		}
	})
}

func TestTaskDefinitionRepoCassandra_GetByID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		target      = models.TaskDefinitionRepoCassandra{}

		partnerID = "pid"
		id, _     = gocql.RandomUUID()
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan()
		iterMock.EXPECT().Scan(AnyN(14)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := target.GetByID(context.Background(), partnerID, id)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan()
		iterMock.EXPECT().Scan(AnyN(14)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := target.GetByID(context.Background(), partnerID, id)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan()
		iterMock.EXPECT().Scan(AnyN(14)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := target.GetByID(context.Background(), partnerID, id)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan()
		iterMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			*(p[12].(*bool)) = true
		}).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(14)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := target.GetByID(context.Background(), partnerID, id)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}

func TestTaskDefNotFoundError_Error(t *testing.T) {
	models.TaskDefNotFoundError{}.Error()
}

func TestTaskDefinitionRepoCassandra_Upsert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		target      = models.TaskDefinitionRepoCassandra{}

		taskDefinition = models.TaskDefinitionDetails{
			TaskDefinition: models.TaskDefinition{
				Deleted: true,
			},
		}
	)

	monkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := target.Upsert(context.Background(), taskDefinition)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}
