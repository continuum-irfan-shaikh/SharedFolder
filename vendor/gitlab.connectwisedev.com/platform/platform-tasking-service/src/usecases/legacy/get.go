package legacy

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// GetByPartner returns data by partner
func (u Usecase) GetByPartner(partnerID string) ([]models.LegacyScriptInfo, error) {
	data, err := u.migrateRepo.GetAllScriptInfoByPartner(partnerID)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, errorcode.NewNotFoundErr(errorcode.ErrorNotFound, fmt.Sprintf("GetByPartner: nothing found by partner %v", partnerID))
		}
		return nil, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskDefinitionTemplate, fmt.Sprintf("GetByPartner: partner %v :%v", partnerID, err))
	}

	for i := range data {
		if data[i].IsSequence || data[i].DefinitionID == "" {
			continue
		}

		if data[i].DefinitionID == "" {
			continue
		}

		definitionUUID, err := gocql.ParseUUID(data[i].DefinitionID)
		if err != nil {
			// transaction id is not required
			u.log.ErrfCtx(context.TODO(), errorcode.ErrorCantProcessData,  "LegacyMigrationUC.GetByPartner: can't parse definition uuid %v: %v", data[i].DefinitionID, err)
			continue
		}

		data[i].DefinitionDetails, err = u.taskDefinitions.GetByID(context.TODO(), partnerID, definitionUUID)
		if err != nil {
			// id is not required
			u.log.ErrfCtx(context.TODO(), errorcode.ErrorCantProcessData, "LegacyMigrationUC.GetByPartner: can't get definition by id %v: %v", data[i].DefinitionID, err)
			continue
		}
	}
	return data, nil
}

// GetByScriptID returns data by script
func (u Usecase) GetByScriptID(partnerID, scriptID string) (info models.LegacyScriptInfo, err error) {
	scriptInfo, err := u.migrateRepo.GetByScriptID(partnerID, scriptID)
	if err != nil {
		if err == gocql.ErrNotFound {
			return info, errorcode.NewNotFoundErr(errorcode.ErrorNotFound, fmt.Sprintf("GetByScriptID: nothing found by partner %v", partnerID))
		}
		return info, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskDefinitionTemplate, fmt.Sprintf("GetByScriptID: partner %v :%v", partnerID, err))
	}

	if scriptInfo.DefinitionID == "" {
		return scriptInfo, nil
	}

	definitionUUID, err := gocql.ParseUUID(scriptInfo.DefinitionID)
	if err != nil {
		return info, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskDefinitionTemplate,
			fmt.Sprintf("LegacyMigrationUC.GetByScriptID: can't parse definition uuid %v: %v", scriptInfo.DefinitionID, err))
	}

	scriptInfo.DefinitionDetails, err = u.taskDefinitions.GetByID(context.TODO(), partnerID, definitionUUID)
	if err != nil {
		return info, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskDefinitionTemplate,
			fmt.Sprintf("LegacyMigrationUC.GetByScriptID: can't get definition by id %v: %v", scriptInfo.DefinitionID, err))
	}

	return scriptInfo, nil
}

// GetJobsByPartner returns  all jobs info by partner
func (u Usecase) GetJobsByPartner(partnerID string) ([]models.LegacyJobInfo, error) {
	data, err := u.migrateRepo.GetAllJobsInfoByPartner(partnerID)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, errorcode.NewNotFoundErr(errorcode.ErrorNotFound, fmt.Sprintf("GetJobsByPartner: nothing found by partner %v", partnerID))
		}
		return nil, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskDefinitionTemplate, fmt.Sprintf("GetJobsByPartner: partner %v :%v", partnerID, err))
	}
	return data, nil
}

// GetByJobID returns job info by job ID
func (u Usecase) GetByJobID(partnerID, scriptID string) (models.LegacyJobInfo, error) {
	jobInfo, err := u.migrateRepo.GetByJobID(partnerID, scriptID)
	if err != nil {
		if err == gocql.ErrNotFound {
			return jobInfo, errorcode.NewNotFoundErr(errorcode.ErrorNotFound, fmt.Sprintf("GetByJobID: nothing found by partner %v", partnerID))
		}
		return jobInfo, errorcode.NewInternalServerErr(errorcode.ErrorCantGetTaskDefinitionTemplate, fmt.Sprintf("GetByJobID: partner %v :%v", partnerID, err))
	}
	return jobInfo, nil
}
