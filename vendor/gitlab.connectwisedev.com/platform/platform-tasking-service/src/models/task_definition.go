package models

//go:generate mockgen -destination=../mocks/mocks-gomock/taskDefinitionPersistance_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models TaskDefinitionPersistence

import (
	"context"
	"fmt"
	"time"

	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/gocql/gocql"
)

const (
	cmdOriginID        = "e3d2c26b-c5ba-49cf-a089-7637f6de949e"
	powershellOriginID = "51a74346-e19b-11e7-9809-0800279505d9"
	bashOriginID       = "37f7f19f-40e8-11e9-a643-e0d55e1ce78a"
	// CMD is a cmd type
	CMD = "cmd"
	// Powershell is a powershell type
	Powershell = "powershell"
	// Bash is bash script type
	Bash = "bash"

	taskDefinitionDetailFields    = "id, partner_id, origin_id, name, description, type, categories, created_at, created_by, updated_at, updated_by, user_parameters, credentials"
	insertTaskDefinitionDetailCql = "INSERT INTO task_definitions (" + taskDefinitionDetailFields + ", deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) USING TTL ?"
	secondsInDay                  = 86400
)

// TaskDefinition structure contains general data about user defined task definition
type TaskDefinition struct {
	ID          gocql.UUID `json:"id"             valid:"unsettableByUsers"`
	PartnerID   string     `json:"partnerId"      valid:"unsettableByUsers"`
	OriginID    gocql.UUID `json:"originId"       valid:"requiredForUsers"`
	Name        string     `json:"name"           valid:"requiredForUsers"`
	Type        string     `json:"type"           valid:"validType"`
	Categories  []string   `json:"categories"     valid:"validCategories"`
	Deleted     bool       `json:"-"`
	Description string     `json:"description"    valid:"optional"`
	Engine      string     `json:"engine"         valid:"unsettableByUsers"`
}

// TaskDefinitionDetails structure contains detailed data about particular user defined task definition
type TaskDefinitionDetails struct {
	TaskDefinition
	CreatedAt      time.Time                `json:"createdAt"      valid:"unsettableByUsers"`
	CreatedBy      string                   `json:"createdBy"      valid:"unsettableByUsers"`
	UpdatedAt      time.Time                `json:"updatedAt"      valid:"unsettableByUsers"`
	UpdatedBy      string                   `json:"updatedBy"      valid:"unsettableByUsers"`
	UserParameters string                   `json:"userParameters" valid:"json"`
	UISchema       string                   `json:"UISchema"       valid:"unsettableByUsers"`
	JSONSchema     string                   `json:"JSONSchema"     valid:"unsettableByUsers"`
	Credentials    *agentModels.Credentials `json:"credentials,omitempty"  valid:"validCreds"`
}

// TaskDefinitionPersistence interface to perform actions with task_definitions table
type TaskDefinitionPersistence interface {
	Upsert(context.Context, TaskDefinitionDetails) error
	GetByID(context.Context, string, gocql.UUID) (TaskDefinitionDetails, error)
	GetAllByPartnerID(context.Context, string) ([]TaskDefinitionDetails, error)
	Exists(context.Context, string, string) bool
	CanBeUpdated(ctx context.Context, partnerID, name string, id gocql.UUID) (canBeUpdated bool, err error)
}

// TaskDefinitionRepoCassandra is a realisation of TaskDefinitionPersistence interface for Cassandra
type TaskDefinitionRepoCassandra struct{}

// TaskDefNotFoundError returns in case of Task Definition with particular ID was not found in the DB for particular partner
type TaskDefNotFoundError struct {
	ID        gocql.UUID
	PartnerID string
}

// Error - error interface implementation for TaskDefNotFoundError type
func (err TaskDefNotFoundError) Error() string {
	return fmt.Sprintf(`Task Definition with ID=%s was not found for partner %s`, err.ID, err.PartnerID)
}

// TaskDefinitionPersistenceInstance is an instance presented TaskDefinitionRepoCassandra
var (
	TaskDefinitionPersistenceInstance TaskDefinitionPersistence = TaskDefinitionRepoCassandra{}
)

// Upsert creates new Task Definition or updates existed in repository
func (TaskDefinitionRepoCassandra) Upsert(ctx context.Context, taskDefinition TaskDefinitionDetails) error {

	taskDefTTL := 0
	if taskDefinition.Deleted {
		taskDefTTL = config.Config.DataRetentionIntervalDay * secondsInDay
	}

	taskDefinitionFields := []interface{}{
		taskDefinition.ID,
		taskDefinition.PartnerID,
		taskDefinition.OriginID,
		taskDefinition.Name,
		taskDefinition.Description,
		taskDefinition.Type,
		taskDefinition.Categories,
		taskDefinition.CreatedAt,
		taskDefinition.CreatedBy,
		taskDefinition.UpdatedAt,
		taskDefinition.UpdatedBy,
		taskDefinition.UserParameters,
		taskDefinition.Credentials,
		taskDefinition.Deleted,
		taskDefTTL,
	}

	cassandraQuery := cassandra.QueryCassandra(ctx, insertTaskDefinitionDetailCql, taskDefinitionFields...)

	return cassandraQuery.Exec()
}

// GetByID returns Task Definition found by ID
func (td TaskDefinitionRepoCassandra) GetByID(ctx context.Context, partnerID string, id gocql.UUID) (definition TaskDefinitionDetails, err error) {
	var (
		query = `SELECT 
					id, 
					partner_id, 
					origin_id,
					name, 
					description, 
					type, 
					categories,
					created_at, 
					created_by,
					updated_at,
					updated_by,
					user_parameters,
					deleted,
					credentials
					FROM task_definitions WHERE partner_id = ? AND id = ?`
		details        []TaskDefinitionDetails
		taskDefDetails TaskDefinitionDetails
		cassandraQuery = cassandra.QueryCassandra(ctx, query, partnerID, id)
		iter           = cassandraQuery.Iter()
	)

	for iter.Scan(
		&taskDefDetails.ID,
		&taskDefDetails.PartnerID,
		&taskDefDetails.OriginID,
		&taskDefDetails.Name,
		&taskDefDetails.Description,
		&taskDefDetails.Type,
		&taskDefDetails.Categories,
		&taskDefDetails.CreatedAt,
		&taskDefDetails.CreatedBy,
		&taskDefDetails.UpdatedAt,
		&taskDefDetails.UpdatedBy,
		&taskDefDetails.UserParameters,
		&taskDefDetails.Deleted,
		&taskDefDetails.Credentials,
	) {
		if taskDefDetails.Deleted {
			continue
		}

		def := TaskDefinitionDetails{
			TaskDefinition: TaskDefinition{
				ID:          taskDefDetails.ID,
				PartnerID:   taskDefDetails.PartnerID,
				OriginID:    taskDefDetails.OriginID,
				Name:        taskDefDetails.Name,
				Type:        taskDefDetails.Type,
				Categories:  taskDefDetails.Categories,
				Description: taskDefDetails.Description,
				Engine:      td.getEngineByOriginID(taskDefDetails.OriginID),
			},
			CreatedAt:      taskDefDetails.CreatedAt,
			CreatedBy:      taskDefDetails.CreatedBy,
			UpdatedAt:      taskDefDetails.UpdatedAt,
			UpdatedBy:      taskDefDetails.UpdatedBy,
			UserParameters: taskDefDetails.UserParameters,
			Credentials:    taskDefDetails.Credentials,
		}
		details = append(details, def)
	}

	if err = iter.Close(); err != nil {
		return definition, fmt.Errorf("error while working with found entities: %v", err)
	}

	if len(details) > 0 {
		return details[0], err
	}
	return definition, TaskDefNotFoundError{id, partnerID}
}

// Exists returns if the task definition already exists or not
func (td TaskDefinitionRepoCassandra) Exists(ctx context.Context, partnerID, name string) bool {
	query := `SELECT name FROM task_definitions WHERE partner_id = ? and name = ? and deleted=false LIMIT 1 ALLOW FILTERING`

	var gotName string
	if err := cassandra.QueryCassandra(ctx, query, partnerID, name).Scan(&gotName); err != nil {
		return false
	}

	if gotName == name {
		return true
	}
	return false
}

// CanBeUpdated returns if the task definition already exists or not with for partner and name, but this exists
// used in updating existing TD, and we're cheking if there is no another one with the same name
func (td TaskDefinitionRepoCassandra) CanBeUpdated(ctx context.Context, partnerID, name string, id gocql.UUID) (bool, error) {
	// retrieving already existing taskDef name by ID
	query := `SELECT name FROM task_definitions WHERE partner_id = ? AND id = ? AND deleted=false ALLOW FILTERING`
	oldTaskDefName := ""
	if err := cassandra.QueryCassandra(ctx, query, partnerID, id).Scan(&oldTaskDefName); err != nil {
		return false, err
	}
	query = `SELECT name, id FROM task_definitions WHERE partner_id = ? and name = ? and deleted=false LIMIT 2 ALLOW FILTERING`

	var (
		gotName       string
		gotID         gocql.UUID
		nameIsChanged = name != oldTaskDefName
		iter          = cassandra.QueryCassandra(ctx, query, partnerID, name).Iter()
	)

	for iter.Scan(&gotName, &gotID) {
		// if we're changing the name from old to new one and this name already exists - we won't update it
		if id != gotID && nameIsChanged {
			return false, nil
		}
	}
	// but if we dont change the name even if it already exists - we can update the script
	return true, iter.Close()
}

// GetAllByPartnerID returns all Task Definition by partner ID
func (td TaskDefinitionRepoCassandra) GetAllByPartnerID(ctx context.Context, partnerID string) ([]TaskDefinitionDetails, error) {
	query := `SELECT 
				id,
				partner_id,
				origin_id,
				name,
				description, 
				type,
				categories, 
				created_at, 
				created_by, 
				updated_at, 
				updated_by, 
				user_parameters,
				deleted,
				credentials
			FROM task_definitions WHERE partner_id = ?`
	var (
		definitions    []TaskDefinitionDetails
		taskDefDetails TaskDefinitionDetails
		params         = []interface{}{
			&taskDefDetails.ID,
			&taskDefDetails.PartnerID,
			&taskDefDetails.OriginID,
			&taskDefDetails.Name,
			&taskDefDetails.Description,
			&taskDefDetails.Type,
			&taskDefDetails.Categories,
			&taskDefDetails.CreatedAt,
			&taskDefDetails.CreatedBy,
			&taskDefDetails.UpdatedAt,
			&taskDefDetails.UpdatedBy,
			&taskDefDetails.UserParameters,
			&taskDefDetails.Deleted,
			&taskDefDetails.Credentials,
		}
		cassandraQuery = cassandra.QueryCassandra(ctx, query, partnerID)
		iter           = cassandraQuery.Iter()
	)

	for iter.Scan(params...) {
		if taskDefDetails.Deleted {
			continue
		}

		def := TaskDefinitionDetails{
			TaskDefinition: TaskDefinition{
				ID:          taskDefDetails.ID,
				PartnerID:   taskDefDetails.PartnerID,
				OriginID:    taskDefDetails.OriginID,
				Name:        taskDefDetails.Name,
				Type:        taskDefDetails.Type,
				Categories:  taskDefDetails.Categories,
				Description: taskDefDetails.Description,
				Engine:      td.getEngineByOriginID(taskDefDetails.OriginID),
			},
			CreatedAt:      taskDefDetails.CreatedAt,
			CreatedBy:      taskDefDetails.CreatedBy,
			UpdatedAt:      taskDefDetails.UpdatedAt,
			UpdatedBy:      taskDefDetails.UpdatedBy,
			UserParameters: taskDefDetails.UserParameters,
		}
		definitions = append(definitions, def)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error while working with found entities: %v", err)
	}
	return definitions, nil
}

func (TaskDefinitionRepoCassandra) getEngineByOriginID(originID gocql.UUID) (engine string) {
	switch originID.String() {
	case bashOriginID:
		return Bash
	case powershellOriginID:
		return Powershell
	case cmdOriginID:
		return CMD
	default:
		return ""
	}
}
