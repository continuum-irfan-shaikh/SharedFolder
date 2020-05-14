package models_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func TestTaskSummaryRepoCassandra_UpdateTaskInstanceStatusCount(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		repo        = models.TaskSummaryRepoCassandra{}

		taskInstanceID, _ = gocql.RandomUUID()
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := repo.UpdateTaskInstanceStatusCount(context.Background(), taskInstanceID, 2, 2)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := repo.UpdateTaskInstanceStatusCount(context.Background(), taskInstanceID, 2, 2)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}

func TestTaskSummaryRepoCassandra_GetStatusCountsByIDs(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		cacheMock   = mocks.NewMockCache(ctrl)
		repo        = models.TaskSummaryRepoCassandra{}

		taskInstanceID, _    = gocql.RandomUUID()
		taskInstancesMapByID = map[gocql.UUID]models.TaskInstance{taskInstanceID: {}}
		taskInstanceIDs      = []gocql.UUID{taskInstanceID}

		cacheKey = []byte("TKS_STATUS_COUNT_BY_ID_" + taskInstanceID.String())
	)

	MonkeyPatchCassandraSession(sessionMock)

	t.Run("positive_1", func(t *testing.T) {
		config.Config.AssetCacheEnabled = true
		taskInstStatusCache := models.TaskInstanceStatusCount{
			TaskInstanceID: taskInstanceID,
			SuccessCount:   10,
			FailureCount:   10,
		}

		b, err := json.Marshal(taskInstStatusCache)
		if err != nil {
			t.Fatalf(err.Error())
		}

		cacheMock.EXPECT().Get(cacheKey).Return(b, nil)

		_, err = repo.GetStatusCountsByIDs(context.Background(), cacheMock, taskInstancesMapByID, taskInstanceIDs)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("positive_2", func(t *testing.T) {
		config.Config.AssetCacheEnabled = true
		b := []byte("invalid")
		cacheMock.EXPECT().Get(cacheKey).Return(b, nil)
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(p ...interface{}) {
			uuid, _ := gocql.RandomUUID()
			*(p[0]).(*gocql.UUID) = uuid
			*(p[1]).(*int) = 2
			*(p[2]).(*int) = 2
		}).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		taskInstancesMapByID = map[gocql.UUID]models.TaskInstance{
			taskInstanceID: {
				PartnerID:     "",
				ID:            gocql.UUID{},
				TaskID:        gocql.UUID{},
				Name:          "",
				OriginID:      gocql.UUID{},
				StartedAt:     time.Time{},
				LastRunTime:   time.Time{},
				Statuses:      nil,
				OverallStatus: 0,
				FailureCount:  0,
				SuccessCount:  0,
				TriggeredBy:   "",
			},
		}

		cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any())
		_, err := repo.GetStatusCountsByIDs(context.Background(), cacheMock, taskInstancesMapByID, taskInstanceIDs)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("positive_3", func(t *testing.T) {
		config.Config.AssetCacheEnabled = false
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(p ...interface{}) {
			uuid, _ := gocql.RandomUUID()
			*(p[0]).(*gocql.UUID) = uuid
			*(p[1]).(*int) = 2
			*(p[2]).(*int) = 2
		}).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := repo.GetStatusCountsByIDs(context.Background(), cacheMock, taskInstancesMapByID, taskInstanceIDs)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		config.Config.AssetCacheEnabled = false
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := repo.GetStatusCountsByIDs(context.Background(), cacheMock, taskInstancesMapByID, taskInstanceIDs)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		config.Config.AssetCacheEnabled = true
		b := []byte("invalid")
		cacheMock.EXPECT().Get(cacheKey).Return(b, nil)
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		taskInstancesMapByID = map[gocql.UUID]models.TaskInstance{
			taskInstanceID: {
				PartnerID:     "",
				ID:            gocql.UUID{},
				TaskID:        gocql.UUID{},
				Name:          "",
				OriginID:      gocql.UUID{},
				StartedAt:     time.Time{},
				LastRunTime:   time.Time{},
				Statuses:      nil,
				OverallStatus: 0,
				FailureCount:  0,
				SuccessCount:  0,
				TriggeredBy:   "",
			},
		}

		cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any())
		_, err := repo.GetStatusCountsByIDs(context.Background(), cacheMock, taskInstancesMapByID, taskInstanceIDs)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}
