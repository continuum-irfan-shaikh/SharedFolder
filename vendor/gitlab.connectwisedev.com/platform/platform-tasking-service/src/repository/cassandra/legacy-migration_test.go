package cassandra

import (
	"errors"
	"testing"

	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/golang/mock/gomock"
)

func TestLegacyMigration_InsertScriptInfo(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		lm          = NewLegacyMigration(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := lm.InsertScriptInfo(models.LegacyScriptInfo{})
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := lm.InsertScriptInfo(models.LegacyScriptInfo{})
		if err == nil {
			t.Fatalf("error shouldn't be eq to nil")
		}
	})
}

func TestLegacyMigration_GetByScriptID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		lm          = NewLegacyMigration(sessionMock)

		partnerID = "123"
		scriptID  = "321"
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, scriptID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)

		_, err := lm.GetByScriptID(partnerID, scriptID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, scriptID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := lm.GetByScriptID(partnerID, scriptID)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}

func TestLegacyMigration_GetAllScriptInfoByPartner(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		lm          = NewLegacyMigration(sessionMock)

		partnerID = "123"
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := lm.GetAllScriptInfoByPartner(partnerID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := lm.GetAllScriptInfoByPartner(partnerID)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}

func TestLegacyMigration_InsertJobInfo(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		lm          = NewLegacyMigration(sessionMock)
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(nil)

		err := lm.InsertJobInfo(models.LegacyJobInfo{})
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New(""))

		err := lm.InsertJobInfo(models.LegacyJobInfo{})
		if err == nil {
			t.Fatalf("error shouldn't be eq to nil")
		}
	})
}

func TestLegacyMigration_GetByJobID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		lm          = NewLegacyMigration(sessionMock)

		partnerID = "123"
		jobID     = "321"
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, jobID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(nil)

		_, err := lm.GetByJobID(partnerID, jobID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID, jobID).Return(queryMock)
		queryMock.EXPECT().Scan(gomock.Any()).Return(errors.New(""))

		_, err := lm.GetByJobID(partnerID, jobID)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}

func TestLegacyMigration_GetAllJobsInfoByPartner(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
		lm          = NewLegacyMigration(sessionMock)

		partnerID = "123"
	)

	t.Run("positive", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(nil)

		_, err := lm.GetAllJobsInfoByPartner(partnerID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), partnerID).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(gomock.Any()).Return(true).Times(1)
		iterMock.EXPECT().Scan(gomock.Any()).Return(false).Times(1)
		iterMock.EXPECT().Close().Return(errors.New(""))

		_, err := lm.GetAllJobsInfoByPartner(partnerID)
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}
