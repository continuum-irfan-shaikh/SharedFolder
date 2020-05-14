package models

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

var timeUUID = gocql.TimeUUID()

// AnyN returns slice on gomock.Any's
// example: batchMock.EXPECT().Query(gomock.Any(), AnyN(2)...).Times(2)
func AnyN(n int) []interface{} {
	m := make([]interface{}, n)
	for i := 0; i < n; i++ {
		m[i] = gomock.Any()
	}

	return m
}

func TestNew(t *testing.T) {

	ti := NewTaskInstance([]Task{{ID: timeUUID, State: statuses.TaskStateActive}}, false)
	if ti.TaskID != timeUUID {
		t.Fatalf("expected %v but got %v", timeUUID, ti.ID)
	}

	NewTaskInstance([]Task{{ID: timeUUID, State: statuses.TaskStateInactive}}, false)
	NewTaskInstance([]Task{{ID: timeUUID, State: statuses.TaskStateDisabled}}, false)
	NewTaskInstance([]Task{{ID: timeUUID, State: statuses.TaskStateActive, PostponedRunTime: someTime}}, false)
	ti = NewTaskInstance([]Task{}, false)
	if ti.TaskID.String() != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("expected %v but got %v", "", ti.TaskID)
	}
}

func TestTaskInstance_CalculateStatuses(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		emptyTI := TaskInstance{}
		gotMap, _ := emptyTI.CalculateStatuses()
		if len(gotMap) > 0 {
			t.Fatalf("expected empty but got %v", gotMap)
		}
	})

	t.Run("2", func(t *testing.T) {
		ti := TaskInstance{
			Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
				timeUUID: statuses.TaskInstanceScheduled,
			},
		}

		expected := map[string]int{statuses.TaskInstanceScheduledText: 1}
		got, err := ti.CalculateStatuses()
		if err != nil {
			t.Fatalf("error is not expected")
		}

		if !reflect.DeepEqual(expected, got) {
			t.Fatalf("expected %v but got %v", expected, got)
		}
	})

	t.Run("err", func(t *testing.T) {
		ti := TaskInstance{
			Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
				timeUUID: 9999999999,
			},
		}

		_, err := ti.CalculateStatuses()
		if err == nil {
			t.Fatalf("error is expected")
		}
	})
}

func TestTaskInstance_IsScheduled(t *testing.T) {
	t.Run("err", func(t *testing.T) {
		ti := TaskInstance{
			Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
				timeUUID: statuses.TaskInstanceFailed,
			},
		}

		if ti.IsScheduled() {
			t.Fatalf("false is expected")
		}
	})

	t.Run("err", func(t *testing.T) {
		ti := TaskInstance{
			Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
				timeUUID:         statuses.TaskInstanceScheduled,
				gocql.TimeUUID(): statuses.TaskInstanceDisabled,
			},
		}

		if !ti.IsScheduled() {
			t.Fatalf("true is expected")
		}
	})
}

func MonkeyPatchCassandraSession(s cassandra.ISession) {
	cassandra.Session = s
}

func TestGetNearestInstanceAfter(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		repo        = TaskInstanceRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		taskID := gocql.TimeUUID()
		sessionMock.EXPECT().Query(gomock.Any(), taskID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(AnyN(11)...).Return(nil)

		_, err := repo.GetNearestInstanceAfter(taskID, time.Now())
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		taskID := gocql.TimeUUID()
		sessionMock.EXPECT().Query(gomock.Any(), taskID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(AnyN(11)...).Return(errors.New("err"))

		_, err := repo.GetNearestInstanceAfter(taskID, time.Now())
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		taskID := gocql.TimeUUID()
		sessionMock.EXPECT().Query(gomock.Any(), taskID, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(AnyN(11)...).Return(gocql.ErrNotFound)

		_, err := repo.GetNearestInstanceAfter(taskID, time.Now())
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestGetTopInstancesByTaskID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		repo        = TaskInstanceRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		taskID := gocql.TimeUUID()
		sessionMock.EXPECT().Query(gomock.Any(), taskID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := repo.GetTopInstancesByTaskID(context.Background(), taskID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		taskID := gocql.TimeUUID()
		sessionMock.EXPECT().Query(gomock.Any(), taskID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New("err"))

		_, err := repo.GetTopInstancesByTaskID(context.Background(), taskID)
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestGetByIDs(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		repo        = TaskInstanceRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		taskID := gocql.TimeUUID()
		sessionMock.EXPECT().Query(gomock.Any(), taskID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := repo.GetByIDs(context.Background(), taskID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		taskID := gocql.TimeUUID()
		sessionMock.EXPECT().Query(gomock.Any(), taskID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(AnyN(13)...).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New("err"))

		_, err := repo.GetByIDs(context.Background(), taskID)
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestInsert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		repo        = TaskInstanceRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(14)...).Return(queryMock).AnyTimes()
		queryMock.EXPECT().Exec().Times(4)

		if err := repo.Insert(context.Background(), TaskInstance{}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(14)...).Return(queryMock).AnyTimes()
		queryMock.EXPECT().Exec().Times(4)

		if err := repo.Insert(context.Background(), TaskInstance{}); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDeleteBatch(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		repo        = TaskInstanceRepoCassandra{}
	)

	t.Run("positive", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(4)...).Return(queryMock).AnyTimes()
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(3)...).Return(queryMock).AnyTimes()
		queryMock.EXPECT().Exec().Times(4)

		if err := repo.DeleteBatch(context.Background(), []TaskInstance{{PartnerID: "Id"}}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("negative", func(t *testing.T) {
		MonkeyPatchCassandraSession(sessionMock)
		sessionMock.EXPECT().Query(gomock.Any(), AnyN(3)...).Return(queryMock).AnyTimes()
		queryMock.EXPECT().Exec().Return(errors.New(""))

		if err := repo.DeleteBatch(context.Background(), []TaskInstance{{PartnerID: "Id"}}); err == nil {
			t.Fatal("extected err")
		}
	})
}
