package repository

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

//go:generate mockgen -destination=../mocks/mocks-gomock/repository_interfaces_mock.go  -package=mocks -source=./intefaces.go

// LegacyMigration interface to work with migration data
type LegacyMigration interface {
	GetAllScriptInfoByPartner(partnerID string) (data []models.LegacyScriptInfo, err error)
	GetByScriptID(partnerID, scriptID string) (data models.LegacyScriptInfo, err error)
	InsertScriptInfo(data models.LegacyScriptInfo) error

	GetAllJobsInfoByPartner(partnerID string) (data []models.LegacyJobInfo, err error)
	GetByJobID(partnerID, jobID string) (data models.LegacyJobInfo, err error)
	InsertJobInfo(data models.LegacyJobInfo) error
}
