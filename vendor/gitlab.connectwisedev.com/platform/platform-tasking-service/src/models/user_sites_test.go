package models_test

import (
	"context"
	"errors"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func TestUser_SaveEndpoints(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		user        = models.User{}
	)

	MonkeyPatchCassandraSession(sessionMock)
	logger.Load(config.Config.Log)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)
		user.SaveEndpoints(context.Background(), []entities.Endpoints{{}})
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))
		user.SaveEndpoints(context.Background(), []entities.Endpoints{{}})
	})
}

func TestUser_EndpointsByPartner(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		user        = models.User{}
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := user.EndpointsByPartner(context.Background(), "")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := user.EndpointsByPartner(context.Background(), "")
		if err == nil {
			t.Fatal("error shouldn't be nil")
		}
	})
}

func TestUser_Endpoints(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		user        = models.User{}
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := user.Endpoints(context.Background(), "", []string{})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := user.Endpoints(context.Background(), "", []string{})
		if err == nil {
			t.Fatal("error shouldn't be nil")
		}
	})
}

func TestUser_GetSitesByTaskID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		user        = models.User{}

		taskID, _ = gocql.RandomUUID()
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any())
		iterMock.EXPECT().Close().Return(nil)

		_, err := user.GetSitesByTaskID(context.Background(), taskID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any())
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := user.GetSitesByTaskID(context.Background(), taskID)
		if err == nil {
			t.Fatal("error cannot be nil")
		}
	})
}

func TestUser_Insert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		user        = models.User{}

		taskID, _ = gocql.RandomUUID()
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("InsertSitesByTaskID", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := user.InsertSitesByTaskID(context.Background(), taskID, []string{})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("InsertUserSites", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := user.InsertUserSites(context.Background(), "", "", []int64{})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestUser_Sites(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		user        = models.User{}
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any())
		iterMock.EXPECT().Close().Return(nil)

		_, err := user.Sites(context.Background(), "", "")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any())
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := user.Sites(context.Background(), "", "")
		if err == nil {
			t.Fatal("error cannot be nil")
		}
	})
}
