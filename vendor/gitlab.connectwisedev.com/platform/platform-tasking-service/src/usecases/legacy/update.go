package legacy

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// InsertScriptInfo inserts legacy script info
func (u Usecase) InsertScriptInfo(partnerID string, data models.LegacyScriptInfo) (err error) {
	data.PartnerID = partnerID
	if err = u.migrateRepo.InsertScriptInfo(data); err != nil {
		return errorcode.NewInternalServerErr(errorcode.ErrorInternalServerError, fmt.Sprintf("can't insert migrate info %v", err))
	}
	return
}

// InsertJobInfo inserts legacy job info
func (u Usecase) InsertJobInfo(partnerID string, data models.LegacyJobInfo) (err error) {
	data.PartnerID = partnerID
	if err = u.migrateRepo.InsertJobInfo(data); err != nil {
		return errorcode.NewInternalServerErr(errorcode.ErrorInternalServerError, fmt.Sprintf("can't insert migrate info %v", err))
	}
	return
}
