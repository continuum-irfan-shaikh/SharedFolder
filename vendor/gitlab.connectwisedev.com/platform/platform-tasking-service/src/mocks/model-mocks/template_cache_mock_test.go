package modelMocks

import (
	"context"

	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

const ID = "00000000-0000-0000-0000-000000000000"
const userHasNOCAccess = true

var repoMock TemplateCacheMock

func TestGetAll(testContext *testing.T) {
	startUp(true)
	templates := DefaultTemplatesDetails[:5]
	actualTemplates, err := repoMock.GetAllTemplatesDetails(getContextWithTransactionID(testContext), TestPartnerID)

	if err != nil {
		testContext.Errorf("Expected Ok but got error: %v", err)
	}

	if err = CompareTemplates(templates, actualTemplates); err != nil {
		testContext.Errorf(err.Error())
	}

	tireDown()
}

func TestTemplateCacheMock_GetAllTemplates(t *testing.T) {

	startUp(true)
	defer tireDown()

	templates, err := repoMock.GetAllTemplates(getContextWithTransactionID(t), TestPartnerID, userHasNOCAccess)
	if err != nil {
		t.Fatal("GetAllTemplates: got error: ", err)
	}

	if templates[0].PartnerID != TestPartnerID {
		t.Fatalf("%#v", templates[0])
	}

	ctx := context.WithValue(getContextWithTransactionID(t), IsNeedError, true)
	errMsgWant := "Cache and Scripting MS is down"
	_, err = repoMock.GetAllTemplates(ctx, TestPartnerID, userHasNOCAccess)
	if err == nil || err.Error() != errMsgWant {
		t.Fatalf("GetAllTemplates: expect an error: %s, but got: %s", errMsgWant, err)
	}
}

func TestGetByType(t *testing.T) {
	startUp(true)
	defer tireDown()
	var tests = []struct {
		partnerID         string
		taskType          string
		expectedTemplates []models.Template
	}{
		{TestPartnerID, `script`, DefaultTemplates[:5]},
		{AnotherPartnerID, `script`, []models.Template{DefaultTemplates[5]}},
		{AnotherPartnerID, ``, []models.Template{}},
	}
	for iteration, test := range tests {
		actualSliceOfTemplates, _ := repoMock.GetByType(getContextWithTransactionID(t), test.partnerID, test.taskType, userHasNOCAccess)
		if !reflect.DeepEqual(test.expectedTemplates, actualSliceOfTemplates) {
			t.Errorf("Iteration #%v, Expected %v, but Actual is %v", iteration, test.expectedTemplates, actualSliceOfTemplates)
		}
	}
}

func TestGetByTypeWithError(t *testing.T) {
	startUp(true)
	defer tireDown()

	ctx := context.WithValue(getContextWithTransactionID(t), IsNeedError, true)
	errMsgWant := "Cache and Scripting MS is down"

	_, err := repoMock.GetByType(ctx, TestPartnerID, `script`, userHasNOCAccess)
	if err == nil || err.Error() != errMsgWant {
		t.Fatalf("GetAllTemplates: expect an error: %s, but got: %s", errMsgWant, err)
	}
}

func TestGetByOriginID(testContext *testing.T) {
	startUp(true)
	var tests = []struct {
		partnerID        string
		originID         gocql.UUID
		expectedTemplate models.TemplateDetails
	}{
		{TestPartnerID, str2uuid("50000000-0000-0000-0000-000000000011"), DefaultTemplatesDetails[4]},
		{TestPartnerID, str2uuid("50000000-0000-0000-0000-000000000022"), models.TemplateDetails{}},
		{AnotherPartnerID, str2uuid("50000000-0000-0000-0000-000000000022"), models.TemplateDetails{}},
	}

	for iteration, test := range tests {
		actualTemplate, _ := repoMock.GetByOriginID(getContextWithTransactionID(testContext), test.partnerID, test.originID, userHasNOCAccess)
		if !reflect.DeepEqual(test.expectedTemplate, actualTemplate) {
			testContext.Errorf("Iteration #%v, Expected %v, but Actual is %v", iteration, test.expectedTemplate, actualTemplate)
		}
	}

	tireDown()
}

func TestGetByOriginIDWithError(t *testing.T) {
	startUp(true)
	defer tireDown()

	ctx := context.WithValue(getContextWithTransactionID(t), IsNeedError, true)
	errMsgWant := "Cache and Scripting MS is down"

	_, err := repoMock.GetByOriginID(ctx, TestPartnerID, str2uuid("50000000-0000-0000-0000-000000000022"), userHasNOCAccess)
	if err == nil || err.Error() != errMsgWant {
		t.Fatalf("GetAllTemplates: expect an error: %s, but got: %s", errMsgWant, err)
	}
}

func TestGetAllTemplatesNegativeWithNeedErr(t *testing.T) {
	startUp(true)
	_, err := repoMock.GetAllTemplatesDetails(getContextWithNeedErr(t, true), TestPartnerID)
	if err == nil {
		t.Error(err)
	}
	tireDown()
}

func TestCompareTemplatesNotEqualSize(t *testing.T) {
	startUp(true)
	err := CompareTemplates([]models.TemplateDetails{
		{
			Name: "Name1",
		},
	}, []models.TemplateDetails{})
	if err == nil {
		t.Error(err)
	}
	tireDown()
}

func TestCompareTemplatesNotEqual(t *testing.T) {
	startUp(true)
	err := CompareTemplates([]models.TemplateDetails{
		{
			Name:     "Name1",
			OriginID: originID,
		},
	}, []models.TemplateDetails{
		{
			Name:     "Name2",
			OriginID: originID,
		},
	})
	if err == nil {
		t.Error(err)
	}
	tireDown()
}

func TestCompareTemplatesInfo(t *testing.T) {
	startUp(true)
	defer tireDown()

	err := CompareTemplatesInfo([]models.Template{
		{
			Name:     "Name1",
			OriginID: originID,
		},
	}, []models.Template{
		{
			Name:     "Name2",
			OriginID: originID,
		},
	})

	if err == nil || !strings.HasPrefix(err.Error(), "Templates are not equals") {
		t.Fatal("CompareTemplatesInfo: must be an error!")
	}

	err = CompareTemplatesInfo([]models.Template{
		{
			Name:     "Name1",
			OriginID: originID,
		},
	}, []models.Template{
		{
			Name:     "Name2",
			OriginID: originID,
		},
		{
			Name:     "Name3",
			OriginID: originID,
		},
	})

	if err == nil || !strings.HasPrefix(err.Error(), "Wrong number of templates") {
		t.Fatal("CompareTemplatesInfo: must be an error!")
	}

	err = CompareTemplatesInfo([]models.Template{
		{
			Name:     "Name1",
			OriginID: originID,
		},
	}, []models.Template{
		{
			Name:     "Name1",
			OriginID: originID,
		},
	})

	if err != nil {
		t.Fatal("CompareTemplatesInfo: got an error: ", err)
	}
}

func startUp(fillMock bool) {
	repoMock = NewTemplateCacheMock(fillMock)
}

func tireDown() {
	repoMock.ClearMock()
}

func getContextWithTransactionID(t *testing.T) (ctx context.Context) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx = req.Context()
	ctx = context.WithValue(ctx, transactionID.Key, ID)
	return ctx
}

func getContextWithNeedErr(t *testing.T, isNeedError bool) (ctx context.Context) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx = req.Context()
	ctx = context.WithValue(ctx, IsNeedError, isNeedError)
	return ctx
}
