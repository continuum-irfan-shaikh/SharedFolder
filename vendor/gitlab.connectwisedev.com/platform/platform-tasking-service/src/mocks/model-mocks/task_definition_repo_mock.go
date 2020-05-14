package modelMocks

import (
	"context"
	"errors"
	"fmt"
	"time"

	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
)

var (
	defaultTime = time.Now().UTC().Truncate(time.Millisecond)
	// ValidUserParams - valid UserParameters according to variable validJSONSchema
	ValidUserParams = "{\"body\":\"test\"}"
	// InvalidUserParams - predefined invalid UserParameters
	InvalidUserParams = "test"
	// InvalidTypeUserParams - predefined invalid UserParameters with wrong type of parameter according to variable validJSONSchema
	InvalidTypeUserParams = "{\"body\":true}"
)

// DefaultTaskDefs is an array of predefined task definitions
var DefaultTaskDefs = []models.TaskDefinitionDetails{
	{TaskDefinition: models.TaskDefinition{OriginID: str2uuid("10000000-0000-0000-0000-000000000001"), Name: "name1", Type: TemplateType, Categories: categories, Description: "description1",},  UserParameters: ValidUserParams},
	{TaskDefinition: models.TaskDefinition{ID: str2uuid("00000000-0000-0000-0000-000000000001"), PartnerID: TestPartnerID, OriginID: str2uuid("10000000-0000-0000-0000-000000000001"), Name: "name1", Type: TemplateType, Categories: categories, Description: "description1", }, CreatedBy: "admin", CreatedAt: defaultTime, UserParameters: ""},
	{TaskDefinition: models.TaskDefinition{ID: str2uuid("00000000-0000-0000-0000-000000000002"), PartnerID: TestPartnerID, OriginID: str2uuid("20000000-0000-0000-0000-000000000002"), Name: "name2", Type: TemplateType, Categories: categories, Description: "description2",},  CreatedBy: "admin", CreatedAt: defaultTime, UserParameters: ""},
	{TaskDefinition: models.TaskDefinition{ID: str2uuid("00000000-0000-0000-0000-000000000003"), PartnerID: "YetAnotherPartner", OriginID: str2uuid("30000000-0000-0000-0000-000000000003"), Name: "name3", Type: TemplateType, Categories: categories, Description: "description3", }, CreatedBy: "admin", CreatedAt: defaultTime, UserParameters: ""},
	{TaskDefinition: models.TaskDefinition{OriginID: str2uuid("50000000-0000-0000-0000-000000000011"), Name: "name5", Type: TemplateType, Categories: categories,  Description: "description5",}, UserParameters: InvalidUserParams},
}

// TaskDefRepoMock represents a Mock for a Task Repository.
type TaskDefRepoMock struct {
	repo map[gocql.UUID]models.TaskDefinitionDetails
}

// NewTaskDefRepoMock creates a Mock for a Task Definition Repository.
// The Mock could be empty or filled with the predefined data.
func NewTaskDefRepoMock(fill bool) TaskDefRepoMock {
	var emptyUUID gocql.UUID
	mock := TaskDefRepoMock{}
	mock.repo = make(map[gocql.UUID]models.TaskDefinitionDetails)
	if fill {
		for _, definition := range DefaultTaskDefs {
			if definition.ID == emptyUUID {
				definition.ID = gocql.TimeUUID()
			}
			mock.repo[definition.TaskDefinition.ID] = definition
		}
	}
	return mock
}

// GetByID returns a TaskDefinitionDetails specified by its ID form the Mocked Repository.
func (mock TaskDefRepoMock) GetByID(ctx context.Context, partnerID string, id gocql.UUID) (models.TaskDefinitionDetails, error) {
	fmt.Println("TaskDefRepoMock.GetByID method called, used RequestID: ", transactionID.FromContext(ctx))
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return models.TaskDefinitionDetails{}, errors.New("Cassandra is down")
	}
	if taskDefinition, ok := mock.repo[id]; ok && taskDefinition.PartnerID == partnerID {
		return taskDefinition, nil
	}
	return models.TaskDefinitionDetails{}, models.TaskDefNotFoundError{
		ID:        id,
		PartnerID: partnerID,
	}
}

// GetAllByPartnerID returns a slice of Task Definitions by Partner ID form the Mocked Repository.
func (mock TaskDefRepoMock) GetAllByPartnerID(ctx context.Context, partnerID string) ([]models.TaskDefinition, error) {
	fmt.Println("TaskRepoMock.GetAllByPartnerID method called, used RequestID: ",
		transactionID.FromContext(ctx))

	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, errors.New("Cassandra is down")
	}

	taskDefinitions := make([]models.TaskDefinition, 0, len(mock.repo))
	for _, taskDefDetails := range mock.repo {
		if taskDefDetails.PartnerID == partnerID {
			taskDefinitions = append(taskDefinitions, taskDefDetails.TaskDefinition)
		}
	}
	return taskDefinitions, nil
}

// Upsert a TaskDefinition into Mocked Repository.
func (mock TaskDefRepoMock) Upsert(ctx context.Context, taskDefinition models.TaskDefinitionDetails) error {
	fmt.Println("TaskDefinitionRepoMock.Insert method called, used RequestID: ", transactionID.FromContext(ctx))
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return errors.New("Cassandra is down")
	}
	mock.repo[taskDefinition.ID] = taskDefinition
	return nil
}

func (mock TaskDefRepoMock)	Exists(context.Context, string, string) bool{
	return false
}