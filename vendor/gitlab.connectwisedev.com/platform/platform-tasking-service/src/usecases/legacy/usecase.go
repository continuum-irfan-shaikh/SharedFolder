package legacy

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
)

//go:generate mockgen -destination=../../mocks/mocks-gomock/legacy_usecase_mock.go  -package=mocks -source=./usecase.go

// MigrationUC represents usecase interface to provide legacy mapping
type MigrationUC interface {
	InsertScriptInfo(partnerID string, data models.LegacyScriptInfo) (err error)
	GetByPartner(partnerID string) ([]models.LegacyScriptInfo, error)
	GetByScriptID(partnerID, scriptID string) (models.LegacyScriptInfo, error)

	InsertJobInfo(partnerID string, data models.LegacyJobInfo) (err error)
	GetJobsByPartner(partnerID string) ([]models.LegacyJobInfo, error)
	GetByJobID(partnerID, scriptID string) (models.LegacyJobInfo, error)
}

// Usecase to work with legacy
type Usecase struct {
	migrateRepo     repository.LegacyMigration
	taskDefinitions models.TaskDefinitionPersistence
	log             logger.Logger
}

// NewMigrationUsecase returns new legacy usecase
func NewMigrationUsecase(migrateRepo repository.LegacyMigration, taskDefinitions models.TaskDefinitionPersistence, log logger.Logger) *Usecase {
	return &Usecase{migrateRepo: migrateRepo, taskDefinitions: taskDefinitions, log: log}
}
