package triggers

import (
	"errors"
	"fmt"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	mocks_cassandra "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestNew(t *testing.T) {
	c := TriggerIDKeyPrefix
	c = TriggerPartnersPrefix
	c = TriggerPartnersByTypePrefix
	c = alertTriggerCategory
	b := shortTTL
	fmt.Println(c, b)
}

func TestTestMethods(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
	)

	t.Run("negative", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New("err"))

		trDef := entities.ActiveTrigger{}
		if err := tr.Insert(trDef); err == nil {
			t.Fatalf("there should be an error")
		}
	})

	t.Run("negative Delete", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New("err"))

		trDef := entities.ActiveTrigger{}
		if err := tr.Delete(trDef); err == nil {
			t.Fatalf("there should be an error")
		}
	})

}

func TestGetAll(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
	)

	t.Run("positive GetAll", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(false).Times(1)
		iterMock.EXPECT().Close()

		//trDef := entities.ActiveTrigger{}
		if _, err := tr.GetAll(); err != nil {
			t.Fatalf("there should not be an error")
		}
	})

	t.Run("negative GetAll", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(false).AnyTimes()
		iterMock.EXPECT().Close().Return(errors.New("err"))

		//trDef := entities.ActiveTrigger{}
		if _, err := tr.GetAll(); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestGetAllByTaskID(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
	)

	t.Run("positive GetAllByTaskID", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(false).Times(1)
		iterMock.EXPECT().Close()

		//trDef := entities.ActiveTrigger{}
		if _, err := tr.GetAllByTaskID("1", gocql.TimeUUID()); err != nil {
			t.Fatalf("there should not be an error")
		}
	})

	t.Run("negative GetAllByTaskID", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(false).AnyTimes()
		iterMock.EXPECT().Close().Return(errors.New("err"))

		//trDef := entities.ActiveTrigger{}
		if _, err := tr.GetAllByTaskID("1", gocql.TimeUUID()); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestGetAllDefinitionsNamesAndIDs(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
	)

	t.Run("positive GetAllDefinitionsNamesAndIDs", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(2)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(mocks.AnyN(2)...).Return(false).Times(1)
		iterMock.EXPECT().Close()

		if _, err := tr.GetAllDefinitionsNamesAndIDs(); err != nil {
			t.Fatalf("there should not be an error")
		}
	})

	t.Run("negative GetAllDefinitionsNamesAndIDs", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(2)...).Return(false).AnyTimes()
		iterMock.EXPECT().Close().Return(errors.New("err"))

		if _, err := tr.GetAllDefinitionsNamesAndIDs(); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestGetByTypeAndPartnerFromDB(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
	)

	t.Run("positive GetAllDefinitionsNamesAndIDs", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(false).Times(1)
		iterMock.EXPECT().Close()

		if _, err := tr.getByTypeAndPartnerFromDB("", ""); err != nil {
			t.Fatalf("there should not be an error")
		}
	})

	t.Run("negative GetAllDefinitionsNamesAndIDs", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(5)...).Return(false).AnyTimes()
		iterMock.EXPECT().Close().Return(errors.New("err"))

		if _, err := tr.getByTypeAndPartnerFromDB("", ""); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestTestGroupTriggers(t *testing.T) {
	RegisterTestingT(t)

	tr := NewTriggersRepo(nil, nil)
	triggers := []entities.ActiveTrigger{{
		Type:      "alert",
		PartnerID: "1",
	}}

	expected := make(map[string]activeTriggersMap)
	expected["1"] = map[string][]entities.ActiveTrigger{"alert": triggers}

	gotMap := tr.groupTriggers(triggers)
	Ω(gotMap).To(Equal(expected), "wrong data", expected)
}

func TestGetDefinition(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
	)

	t.Run("positive but can't unmarshall GetDefinition", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(mocks.AnyN(1)...)

		if _, err := tr.GetDefinition("type"); err == nil {
			t.Fatalf("there should not be an error")
		}
	})

	t.Run("negative GetDefinition", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(mocks.AnyN(1)...).Return(errors.New("err"))

		if _, err := tr.GetDefinition("type"); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestGetTriggerCounterByType(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
	)
	t.Run("positive  GetTriggerCounterByType", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(mocks.AnyN(3)...).AnyTimes()

		if _, err := tr.GetTriggerCounterByType("type"); err != nil {
			t.Fatalf("there should not be an error")
		}
	})
	ctrl.Finish()

	ctrl = gomock.NewController(t)
	sessionMock = mocks_cassandra.NewMockISession(ctrl)
	queryMock = mocks_cassandra.NewMockIQuery(ctrl)
	tr = NewTriggersRepo(sessionMock, nil)
	t.Run("negative GetTriggerCounterByType", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(mocks.AnyN(3)...).Return(errors.New("err")).AnyTimes()

		if _, err := tr.GetTriggerCounterByType("type"); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestGetTriggerCounterByTypeNegative(t *testing.T) {
	ctrl := gomock.NewController(t)
	sessionMock := mocks_cassandra.NewMockISession(ctrl)
	queryMock := mocks_cassandra.NewMockIQuery(ctrl)
	tr := NewTriggersRepo(sessionMock, nil)

	t.Run("negative GetTriggerCounterByType", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Scan(mocks.AnyN(3)...).Return(gocql.ErrNotFound).Times(1)

		if _, err := tr.GetTriggerCounterByType("type"); err != nil {
			t.Fatalf("there should not be an error")
		}
	})
}

func TestIncreaseDecrease(t *testing.T) {
	ctrl := gomock.NewController(t)
	sessionMock := mocks_cassandra.NewMockISession(ctrl)
	queryMock := mocks_cassandra.NewMockIQuery(ctrl)
	tr := NewTriggersRepo(sessionMock, nil)

	t.Run("negative IncreaseCounter", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec()

		if err := tr.IncreaseCounter(entities.TriggerCounter{}); err != nil {
			t.Fatalf("there should not be an error")
		}
	})

	t.Run("negative IncreaseCounter", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec()

		if err := tr.DecreaseCounter(entities.TriggerCounter{}); err != nil {
			t.Fatalf("there should not be an error")
		}
	})
}

func TestGetAllDefinitions(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
		iterMock    = mocks_cassandra.NewMockIIter(ctrl)
	)

	t.Run("negative can't parse GetAllDefinitions", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(1)...).Return(true).Times(1)
		iterMock.EXPECT().Scan(mocks.AnyN(1)...).Return(false).Times(1)
		iterMock.EXPECT().Close()

		if _, err := tr.GetAllDefinitions(); err == nil {
			t.Fatalf("there should be an error")
		}
	})

	ctrl = gomock.NewController(t)
	sessionMock = mocks_cassandra.NewMockISession(ctrl)
	queryMock = mocks_cassandra.NewMockIQuery(ctrl)
	tr = NewTriggersRepo(sessionMock, nil)
	iterMock = mocks_cassandra.NewMockIIter(ctrl)

	t.Run("negative GetAllDefinitions", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Iter().Return(iterMock)
		iterMock.EXPECT().Scan(mocks.AnyN(1)...).Return(false).AnyTimes()
		iterMock.EXPECT().Close().Return(errors.New("err"))

		if _, err := tr.GetAllDefinitions(); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}

func TestInsertDefinitions(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		sessionMock = mocks_cassandra.NewMockISession(ctrl)
		queryMock   = mocks_cassandra.NewMockIQuery(ctrl)
		tr          = NewTriggersRepo(sessionMock, nil)
	)

	trDef := []entities.TriggerDefinition{
		{
			ID:              "id1",
			TriggerCategory: alertTriggerCategory,
		},
		{
			ID:              "id2",
			TriggerCategory: "generic",
		},
	}
	t.Run("positive InsertDefinitions", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock).Times(len(trDef))
		queryMock.EXPECT().Exec().Times(len(trDef))

		if err := tr.InsertDefinitions(trDef); err != nil {
			t.Fatalf("there should not be an error")
		}
	})

	t.Run("negative InsertDefinitions", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New("err"))

		if err := tr.InsertDefinitions(trDef); err == nil {
			t.Fatalf("there should be an error")
		}
	})

	t.Run("negative TruncateDefinitions", func(t *testing.T) {
		sessionMock.EXPECT().Query(gomock.Any()).Return(queryMock)
		queryMock.EXPECT().Exec().Return(errors.New("err"))

		if err := tr.TruncateDefinitions(); err == nil {
			t.Fatalf("there should be an error")
		}
	})
}
