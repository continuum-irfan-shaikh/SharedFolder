package cassandra

import (
	"errors"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"github.com/golang/mock/gomock"
)

func TestExecutionExpiration_Insert(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		ee          = NewExecutionExpiration(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := ee.Insert(entities.ExecutionExpiration{}, 0)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := ee.Insert(entities.ExecutionExpiration{}, 0)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}
