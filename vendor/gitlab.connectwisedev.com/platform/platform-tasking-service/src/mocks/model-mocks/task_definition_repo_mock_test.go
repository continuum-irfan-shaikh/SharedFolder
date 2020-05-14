package modelMocks

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"golang.org/x/net/context"
	"reflect"
	"testing"
	"time"
)

func TestTaskDefRepoMock_Create(t *testing.T) {

	repoMock := NewTaskDefRepoMock(true)
	taskDef := models.TaskDefinition{
		ID:         gocql.TimeUUID(),
		PartnerID:  TestPartnerID,
		Name:       "Test9",
		Categories: []string{"Custom"},
		Type:       "script",
		OriginID:   gocql.TimeUUID(),
	}
	taskDefDetails := models.TaskDefinitionDetails{
		TaskDefinition: taskDef,
		CreatedAt:      time.Now().UTC().Truncate(time.Millisecond),
		CreatedBy:      "admin",
		UserParameters: "aaaaa",
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, IsNeedError, true)

	if err := repoMock.Upsert(ctx, taskDefDetails); err == nil {
		t.Fatal("Expected error during mock insertion")
	}

	if err := repoMock.Upsert(context.Background(), taskDefDetails); err != nil {
		t.Fatalf("Error: %v", err)
	}

	if (len(DefaultTaskDefs) + 1) != len(repoMock.repo) {
		t.Fatalf("Expected mock size %v but got %v", len(DefaultTaskDefs)+1, len(repoMock.repo))
	}
}

func TestTaskDefRepoMock_GetByID(t *testing.T) {

	repoMock := NewTaskDefRepoMock(true)

	ctx := context.Background()
	ctx = context.WithValue(ctx, IsNeedError, true)
	if _, err := repoMock.GetByID(ctx, DefaultTaskDefs[2].PartnerID, DefaultTaskDefs[2].ID); err == nil {
		t.Fatal("Expected error during mock lookup")
	}

	if _, err := repoMock.GetByID(context.Background(), "UnknownPartner", DefaultTaskDefs[2].ID); err == nil {
		t.Fatal("Expected error during mock lookup")
	}

	actualTaskDef, err := repoMock.GetByID(context.Background(), DefaultTaskDefs[2].PartnerID, DefaultTaskDefs[2].ID)
	if err != nil {
		t.Fatalf("Error during mock lookup: %v", err)
	}

	if !reflect.DeepEqual(actualTaskDef, DefaultTaskDefs[2]) {
		t.Fatalf("Expected %v but got %v", DefaultTaskDefs[2], actualTaskDef)
	}
}

func TestTaskDefRepoMock_GetAllByPartnerID(t *testing.T) {

	repoMock := NewTaskDefRepoMock(true)

	ctx := context.Background()
	ctx = context.WithValue(ctx, IsNeedError, true)
	if _, err := repoMock.GetAllByPartnerID(ctx, DefaultTaskDefs[1].PartnerID); err == nil {
		t.Fatal("Expected error during mock lookup")
	}

	actualTaskDefs, err := repoMock.GetAllByPartnerID(context.Background(), DefaultTaskDefs[1].PartnerID)
	if err != nil {
		t.Fatalf("Error during mock lookup: %v", err)
	}

	if !reflect.DeepEqual(len(actualTaskDefs), 2) {
		t.Fatalf("Expected %v but got %v", 2, len(actualTaskDefs))
	}
}
