package cassandra

import (
	"errors"
	"testing"

	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func TestProfilesRepo_Delete(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		pr          = NewProfilesRepo(sessionMock)

		taskID, _ = gocql.RandomUUID()
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), taskID).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := pr.Delete(taskID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), taskID).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := pr.Delete(taskID)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestProfilesRepo_Insert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		pr          = NewProfilesRepo(sessionMock)

		taskID, _    = gocql.RandomUUID()
		profileID, _ = gocql.RandomUUID()
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), taskID, profileID).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := pr.Insert(taskID, profileID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), taskID, profileID).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := pr.Insert(taskID, profileID)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestProfilesRepo_GetByTaskID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		pr          = NewProfilesRepo(sessionMock)

		taskID, _ = gocql.RandomUUID()
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)

		_, err := pr.GetByTaskID(taskID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := pr.GetByTaskID(taskID)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}
