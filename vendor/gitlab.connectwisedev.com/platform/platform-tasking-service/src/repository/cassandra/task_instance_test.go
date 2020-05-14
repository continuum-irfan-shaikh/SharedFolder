package cassandra

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func init() {
	config.Config.CassandraConcurrentCallNumber = 10
}

func TestTaskInstance_GetMinimalInstanceByID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		taskInst    = NewTaskInstance(sessionMock)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)

		id = "1"
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), id).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)
		queryMock.EXPECT().Release()

		_, err := taskInst.GetMinimalInstanceByID(id)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), id).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))
		queryMock.EXPECT().Release()

		_, err := taskInst.GetMinimalInstanceByID(id)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestTaskInstance_GetNearestInstanceAfter(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		taskInst    = NewTaskInstance(sessionMock)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)

		id, _ = gocql.RandomUUID()
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), id, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)

		_, err := taskInst.GetNearestInstanceAfter(id, time.Now())
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), id, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := taskInst.GetNearestInstanceAfter(id, time.Now())
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), id, gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(gocql.ErrNotFound)

		_, err := taskInst.GetNearestInstanceAfter(id, time.Now())
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestTaskInstance_Insert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		taskInst    = NewTaskInstance(sessionMock)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock).AnyTimes()
		queryMock.EXPECT().Exec().Return(nil).Times(3)

		err := taskInst.Insert(models.TaskInstance{}, 0)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock).AnyTimes()
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := taskInst.Insert(models.TaskInstance{}, 0)
		if err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestTaskInstance_GetInstance(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		taskInst    = NewTaskInstance(sessionMock)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)

		id, _ = gocql.RandomUUID()
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), id).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)

		_, err := taskInst.GetInstance(id)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), id).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := taskInst.GetInstance(id)
		if err == nil {
			t.Fatalf("there should be error")
		}
	})
}

func TestTaskInstance_GetTopInstancesForScheduledByTaskIDs(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		taskInst    = NewTaskInstance(sessionMock)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)

		instanceID1 = "1"
		instanceIDs = []string{instanceID1}
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), instanceID1).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := taskInst.GetTopInstancesForScheduledByTaskIDs(instanceIDs)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}

func TestTaskInstance_GetByStartedAtAfter(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		taskInst    = NewTaskInstance(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := taskInst.GetByStartedAtAfter("", time.Now(), time.Now())
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := taskInst.GetByStartedAtAfter("", time.Now(), time.Now())
		if err == nil {
			t.Fatalf("there should be error")
		}
	})
}

func TestTaskInstance_GetInstancesForScheduled(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		taskInst    = NewTaskInstance(sessionMock)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)

		instanceID1 = "1"
		instanceIDs = []string{instanceID1}
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), instanceID1).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := taskInst.GetInstancesForScheduled(instanceIDs)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), instanceID1).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := taskInst.GetInstancesForScheduled(instanceIDs)
		if err == nil {
			t.Fatalf("there should be error")
		}
	})
}

func TestNewTaskInstance(t *testing.T) {
	ti := NewTaskInstance(mocks_cassandra.NewMockISession(gomock.NewController(t)))
	if reflect.TypeOf(ti) != reflect.TypeOf(&TaskInstance{}) {
		t.Fatalf("NewTaskInstance didn't create *TaskInstance")
	}
}
