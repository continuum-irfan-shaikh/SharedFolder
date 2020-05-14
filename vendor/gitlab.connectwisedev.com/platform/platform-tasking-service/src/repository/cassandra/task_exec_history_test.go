package cassandra

import (
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"github.com/golang/mock/gomock"
)

func TestTaskExecutionHistory_Insert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		teh         = NewTaskExecutionHistory(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := teh.Insert(entities.TaskExecHistory{})
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}
