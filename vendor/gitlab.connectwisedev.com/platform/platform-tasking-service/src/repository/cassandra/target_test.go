package cassandra

import (
	"errors"
	"testing"

	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func TestTargets_Insert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		targets     = NewTargets(sessionMock)

		partnerID = "1"
		taskID, _ = gocql.RandomUUID()
		target    = models.Target{}
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := targets.Insert(partnerID, taskID, target)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := targets.Insert(partnerID, taskID, target)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestTargets_GetTargetsByTaskID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		targets     = NewTargets(sessionMock)

		partnerID = "1"
		taskID, _ = gocql.RandomUUID()

		targetIndex   = 0
		targetIDIndex = 1
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, taskID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			*(p[targetIndex].(*models.TargetType)) = models.ManagedEndpoint
			*(p[targetIDIndex].(*[]string)) = []string{taskID.String()}
		}).Return(nil)

		_, err := targets.GetTargetsByTaskID(partnerID, taskID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, taskID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := targets.GetTargetsByTaskID(partnerID, taskID)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, taskID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(gocql.ErrNotFound)

		_, err := targets.GetTargetsByTaskID(partnerID, taskID)
		if err != gocql.ErrNotFound {
			t.Fatalf("err should be eq to gocql.ErrNotFound")
		}
	})

	t.Run("negative_3", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, taskID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			*(p[targetIndex].(*models.TargetType)) = models.DynamicGroup
			*(p[targetIDIndex].(*[]string)) = []string{taskID.String()}
		}).Return(nil)

		_, err := targets.GetTargetsByTaskID(partnerID, taskID)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})

	t.Run("negative_4", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, taskID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
			*(p[targetIndex].(*models.TargetType)) = models.ManagedEndpoint
			*(p[targetIDIndex].(*[]string)) = []string{"invalid"}
		}).Return(nil)

		_, err := targets.GetTargetsByTaskID(partnerID, taskID)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}
