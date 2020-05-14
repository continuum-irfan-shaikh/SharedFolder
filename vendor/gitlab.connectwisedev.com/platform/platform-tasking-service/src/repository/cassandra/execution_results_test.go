package cassandra

import (
	"errors"
	"testing"

	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func TestScriptExecutionResults_GetLastExecutions(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		sxr         = NewScriptExecutionResults(sessionMock)

		partnerID        = "123"
		endpointID       = "1"
		endpointIDsSlice = []string{endpointID}
		endpointIDsMap   = map[string]struct{}{
			endpointID: {},
		}
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, endpointIDsSlice).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Release()
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := sxr.GetLastExecutions(partnerID, endpointIDsMap)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, endpointIDsSlice).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Release()
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := sxr.GetLastExecutions(partnerID, endpointIDsMap)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}

func TestScriptExecutionResults_GetLastResultByEndpointID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		sxr         = NewScriptExecutionResults(sessionMock)

		endpointID = "123"
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), endpointID).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Release()
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)

		_, err := sxr.GetLastResultByEndpointID(endpointID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), endpointID).Return(queryMock)
		queryMock.EXPECT().SetConsistency(gocql.One)
		queryMock.EXPECT().Release()
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := sxr.GetLastResultByEndpointID(endpointID)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}
