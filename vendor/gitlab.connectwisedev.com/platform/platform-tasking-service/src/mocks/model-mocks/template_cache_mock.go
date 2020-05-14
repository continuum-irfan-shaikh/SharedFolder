package modelMocks

import (
	"context"
	"fmt"
	"reflect"

	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

const (
	defaultExpectedExecutionTimeSec = 300
)

// TemplateCacheMock represents a Mock for a TemplateDetails Repository
type TemplateCacheMock struct {
	Repo          map[string][]models.TemplateDetails
	RepoTemplates map[string][]models.Template
}

var (
	categories = []string{"Management", "Security"}
	// TestPartnerID - predefined test partner ID
	TestPartnerID = "00000001"
	// AnotherPartnerID - predefined another one test partner ID
	AnotherPartnerID = "00000002"
	// TestPartnerIDWithoutTemplates - predefined test partner ID without Templates
	TestPartnerIDWithoutTemplates = "00000003"
	// TemplateType - predefined template type
	TemplateType = "script"
	emptyUUID    gocql.UUID
	// ValidJSONSchema - predefined valid json schema
	ValidJSONSchema = "{\"properties\":{\"body\":{\"title\":\"PowerShell Script\",\"type\":\"string\"}},\"required\":[\"body\"],\"type\":\"object\",\"additionalProperties\": false}"
)

// DefaultTemplates is an array of predefined templates
var DefaultTemplates = []models.Template{
	{PartnerID: TestPartnerID, OriginID: str2uuid("10000000-0000-0000-0000-000000000001"), Name: "name1", Description: "description1", Type: TemplateType, Categories: categories},
	{PartnerID: TestPartnerID, OriginID: str2uuid("20000000-0000-0000-0000-000000000011"), Name: "name2", Description: "description2", Type: TemplateType, Categories: categories},
	{PartnerID: TestPartnerID, OriginID: str2uuid("30000000-0000-0000-0000-000000000011"), Name: "name3", Description: "description3", Type: TemplateType, Categories: categories},
	{PartnerID: TestPartnerID, OriginID: str2uuid("40000000-0000-0000-0000-000000000011"), Name: "name4", Description: "description4", Type: TemplateType, Categories: categories},
	{PartnerID: TestPartnerID, OriginID: str2uuid("50000000-0000-0000-0000-000000000011"), Name: "name5", Description: "description5", Type: TemplateType, Categories: categories},
	{PartnerID: AnotherPartnerID, OriginID: str2uuid("60000000-0000-0000-0000-000000000011"), Name: "name6", Description: "description6", Type: TemplateType, Categories: categories},
}

// DefaultTemplatesDetails is an array of predefined templates
var DefaultTemplatesDetails = []models.TemplateDetails{
	{PartnerID: TestPartnerID, OriginID: str2uuid("10000000-0000-0000-0000-000000000001"), Name: "name1", Description: "description1", Type: TemplateType, Categories: categories, JSONSchema: ValidJSONSchema},
	{PartnerID: TestPartnerID, OriginID: str2uuid("20000000-0000-0000-0000-000000000011"), Name: "name2", Description: "description2", Type: TemplateType, Categories: categories},
	{PartnerID: TestPartnerID, OriginID: str2uuid("30000000-0000-0000-0000-000000000011"), Name: "name3", Description: "description3", Type: TemplateType, Categories: categories},
	{PartnerID: TestPartnerID, OriginID: str2uuid("40000000-0000-0000-0000-000000000011"), Name: "name4", Description: "description4", Type: TemplateType, Categories: categories},
	{PartnerID: TestPartnerID, OriginID: str2uuid("50000000-0000-0000-0000-000000000011"), Name: "name5", Description: "description5", Type: TemplateType, Categories: categories, JSONSchema: ValidJSONSchema},
	{PartnerID: AnotherPartnerID, OriginID: str2uuid("60000000-0000-0000-0000-000000000011"), Name: "name6", Description: "description6", Type: TemplateType, Categories: categories},
}

// NewTemplateCacheMock creates a Mock for a TemplateDetails Cache.
// The Mock could be empty or filled with the predefined data.
func NewTemplateCacheMock(fill bool) TemplateCacheMock {
	mock := TemplateCacheMock{}
	mock.Repo = make(map[string][]models.TemplateDetails)
	mock.RepoTemplates = make(map[string][]models.Template)
	if fill {
		mock.Repo[TestPartnerID] = DefaultTemplatesDetails[0:5]
		mock.Repo[AnotherPartnerID] = DefaultTemplatesDetails[5:]
		mock.Repo[TestPartnerIDWithoutTemplates] = nil
		mock.RepoTemplates[TestPartnerID] = DefaultTemplates[0:5]
		mock.RepoTemplates[AnotherPartnerID] = DefaultTemplates[5:]
		mock.RepoTemplates[TestPartnerIDWithoutTemplates] = nil
	}

	return mock
}

// ClearMock removes all data from the Mocked Repository
func (mock *TemplateCacheMock) ClearMock() {
	mock.Repo = nil
}

// ExistsWithName ..
func (mock TemplateCacheMock) ExistsWithName(ctx context.Context, partnerID, scriptName string) bool {
	return false
}

// GetAllTemplates returns all templates form the Mocked Repository by partnerID
func (mock TemplateCacheMock) GetAllTemplates(ctx context.Context, partnerID string, hasNOCAccess bool) ([]models.Template, error) {
	fmt.Printf("TemplateCacheMock.GetAllTemplatesDetails method called, used RequestID: %v, partnerID: %v\n", transactionID.FromContext(ctx), partnerID)

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cache and Scripting MS is down")
	}

	templates := make([]models.Template, 0, len(mock.Repo))
	for _, template := range mock.RepoTemplates[partnerID] {
		if partnerID == template.PartnerID || template.PartnerID == emptyUUID.String() {
			templates = append(templates, template)
		}
	}
	return templates, nil
}

// CalculateExpectedExecutionTimeSec ..
func (mock TemplateCacheMock) CalculateExpectedExecutionTimeSec(arg0 context.Context, arg1 models.Task) int {
	return defaultExpectedExecutionTimeSec
}

// GetAllTemplatesDetails returns all templates form the Mocked Repository by partnerID
func (mock TemplateCacheMock) GetAllTemplatesDetails(ctx context.Context, partnerID string) ([]models.TemplateDetails, error) {
	fmt.Printf("TemplateCacheMock.GetAllTemplatesDetails method called, used RequestID: %v, partnerID: %v\n", transactionID.FromContext(ctx), partnerID)

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cache and Scripting MS is down")
	}

	templates := make([]models.TemplateDetails, 0, len(mock.Repo))
	for _, template := range mock.Repo[partnerID] {
		if partnerID == template.PartnerID || template.PartnerID == emptyUUID.String() {
			templates = append(templates, template)
		}
	}
	return templates, nil
}

// GetByType returns templates by TaskType
func (mock TemplateCacheMock) GetByType(ctx context.Context, partnerID string, taskType string, hasNOCAccess bool) ([]models.Template, error) {
	fmt.Printf("TemplateCacheMock.GetByType method called, used RequestID: %v, partnerID: %v, taskType: %v\n",
		transactionID.FromContext(ctx), partnerID, taskType)

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cache and Scripting MS is down")
	}
	templates := make([]models.Template, 0, len(mock.RepoTemplates))
	for _, template := range mock.RepoTemplates[partnerID] {
		if taskType == template.Type {
			templates = append(templates, template)
		}
	}
	return templates, nil
}

// GetByOriginID returns TemplateDetails founded by Origin ID and partner ID
func (mock TemplateCacheMock) GetByOriginID(ctx context.Context, partnerID string, originID gocql.UUID, hasNOCAccess bool) (models.TemplateDetails, error) {
	fmt.Printf("TemplateCacheMock.GetByOriginID method called, used RequestID: %v, partnerID: %v, originID: %v",
		transactionID.FromContext(ctx), partnerID, originID.String())

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return models.TemplateDetails{}, errors.New("Cache and Scripting MS is down")
	}
	for _, template := range mock.Repo[partnerID] {
		if template.OriginID == originID && template.PartnerID == partnerID {
			return template, nil
		}
	}
	return models.TemplateDetails{}, models.TemplateNotFoundError{
		OriginID:  originID,
		PartnerID: partnerID,
	}
}

// CompareTemplates - compares a slices of task definition templates.
func CompareTemplates(expectedTemplates, actualTemplates []models.TemplateDetails) error {
	if len(expectedTemplates) != len(actualTemplates) {
		return fmt.Errorf("Wrong number of templates, expected: %v, actual: %v", len(expectedTemplates), len(actualTemplates))
	}
	for _, expectedTemplate := range expectedTemplates {
		for _, actualTemplate := range actualTemplates {
			if expectedTemplate.OriginID == actualTemplate.OriginID {
				if !reflect.DeepEqual(actualTemplate, expectedTemplate) {
					return fmt.Errorf("Templates are not equals: %v, actual: %v", expectedTemplate, actualTemplate)
				}
			}
		}
	}
	return nil
}

// CompareTemplatesInfo - compares a slices of task definition templates.
func CompareTemplatesInfo(expectedTemplates, actualTemplates []models.Template) error {
	if len(expectedTemplates) != len(actualTemplates) {
		return fmt.Errorf("Wrong number of templates, expected: %v, actual: %v", len(expectedTemplates), len(actualTemplates))
	}
	for _, expectedTemplate := range expectedTemplates {
		for _, actualTemplate := range actualTemplates {
			if expectedTemplate.OriginID == actualTemplate.OriginID {
				if !reflect.DeepEqual(actualTemplate, expectedTemplate) {
					return fmt.Errorf("Templates are not equals: %v, actual: %v", expectedTemplate, actualTemplate)
				}
			}
		}
	}
	return nil
}
