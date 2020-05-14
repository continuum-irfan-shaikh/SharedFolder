package cassandra

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

// NewLegacyMigration  returns new Scheduler repo client
func NewLegacyMigration(conn cassandra.ISession) *LegacyMigration {
	return &LegacyMigration{
		conn: conn,
	}
}

// LegacyMigration  repo client to work with legacy_migration_info table
type LegacyMigration struct {
	conn cassandra.ISession
}

// GetAllJobsInfoByPartner gets all job info by partner
func (l *LegacyMigration) GetAllJobsInfoByPartner(partnerID string) (data []models.LegacyJobInfo, err error) {
	var (
		query = `SELECT 
				legacy_job_id, 
				legacy_script_id, 
				legacy_template_id, 
				origin_id,
				definition_id, 
				type, 
				task_id,
				reason
			FROM job_migration_info WHERE partner_id = ?`
		iter = l.conn.Query(query, partnerID).Iter()

		info   = models.LegacyJobInfo{}
		params = []interface{}{
			&info.LegacyJobID,
			&info.LegacyScriptID,
			&info.LegacyTemplateID,
			&info.OriginID,
			&info.DefinitionID,
			&info.Type,
			&info.TaskID,
			&info.ErrorReason,
		}
	)

	for iter.Scan(params...) {
		migrationData := models.LegacyJobInfo{
			PartnerID:        partnerID,
			LegacyTemplateID: info.LegacyTemplateID,
			LegacyScriptID:   info.LegacyScriptID,
			OriginID:         info.OriginID,
			DefinitionID:     info.DefinitionID,
			TaskID:           info.TaskID,
			Type:             info.Type,
			ErrorReason:      info.ErrorReason,
		}
		data = append(data, migrationData)
	}

	if err = iter.Close(); err != nil {
		return nil, err
	}
	return
}

// GetByJobID get job info by  jobID
func (l *LegacyMigration) GetByJobID(partnerID, jobID string) (data models.LegacyJobInfo, err error) {
	var (
		query = `SELECT 
				legacy_job_id, 
				legacy_script_id, 
				legacy_template_id, 
				origin_id, 
				definition_id,
				type,
				task_id,
				reason
			FROM job_migration_info WHERE partner_id = ? AND legacy_job_id = ?`
		params = []interface{}{
			&data.LegacyJobID,
			&data.LegacyScriptID,
			&data.LegacyTemplateID,
			&data.OriginID,
			&data.DefinitionID,
			&data.Type,
			&data.TaskID,
			&data.ErrorReason,
		}
	)

	if err = l.conn.Query(query, partnerID, jobID).Scan(params...); err != nil {
		return data, err
	}

	data.PartnerID = partnerID
	return
}

// InsertJobInfo inserts migration job info
func (l *LegacyMigration) InsertJobInfo(data models.LegacyJobInfo) error {
	query := `INSERT INTO job_migration_info ( 
				partner_id,
				legacy_job_id,
				legacy_script_id,
				legacy_template_id,
				origin_id,
				definition_id,
				type,
				task_id,
				reason
			) VALUES (?,?,?,?,?,?,?,?,?)`
	params := []interface{}{
		data.PartnerID,
		data.LegacyJobID,
		data.LegacyScriptID,
		data.LegacyTemplateID,
		data.OriginID,
		data.DefinitionID,
		data.Type,
		data.TaskID,
		data.ErrorReason,
	}

	if err := l.conn.Query(query, params...).Exec(); err != nil {
		return fmt.Errorf("can't insert legacy info profile: %v", err)
	}
	return nil
}

// GetAllScriptInfoByPartner returns migration Data by partner
func (l *LegacyMigration) GetAllScriptInfoByPartner(partnerID string) (data []models.LegacyScriptInfo, err error) {
	query := `SELECT 
					legacy_id,
					legacy_template_id, 
					origin_id,
					definition_id,
					is_sequence,
					is_parametrized,
					reason
			FROM script_migration_info WHERE partner_id = ?`
	iter := l.conn.Query(query, partnerID).Iter()
	info := models.LegacyScriptInfo{}
	params := []interface{}{
		&info.LegacyScriptID,
		&info.LegacyTemplateID,
		&info.OriginID,
		&info.DefinitionID,
		&info.IsSequence,
		&info.IsParametrized,
		&info.ErrorReason,
	}

	for iter.Scan(params...) {
		migrationData := models.LegacyScriptInfo{
			PartnerID:        partnerID,
			LegacyTemplateID: info.LegacyTemplateID,
			LegacyScriptID:   info.LegacyScriptID,
			OriginID:         info.OriginID,
			DefinitionID:     info.DefinitionID,
			IsSequence:       info.IsSequence,
			IsParametrized:   info.IsParametrized,
			ErrorReason:      info.ErrorReason,
		}
		data = append(data, migrationData)
	}

	if err = iter.Close(); err != nil {
		return nil, err
	}
	return
}

// GetByScriptID returns migration Data by partner
func (l *LegacyMigration) GetByScriptID(partnerID, scriptID string) (data models.LegacyScriptInfo, err error) {
	query := `SELECT origin_id,
				definition_id,
				legacy_template_id, 
				is_parametrized,
				is_sequence,
				reason
			FROM script_migration_info WHERE partner_id = ? AND legacy_id = ?`
	params := []interface{}{
		&data.OriginID,
		&data.DefinitionID,
		&data.LegacyTemplateID,
		&data.IsParametrized,
		&data.IsSequence,
		&data.ErrorReason,
	}

	if err = l.conn.Query(query, partnerID, scriptID).Scan(params...); err != nil {
		return data, err
	}

	data.PartnerID = partnerID
	data.LegacyScriptID = scriptID
	return
}

// InsertScriptInfo inserts profile and task IDs to cassandra
func (l *LegacyMigration) InsertScriptInfo(data models.LegacyScriptInfo) error {
	query := `INSERT INTO script_migration_info (
				partner_id,
				legacy_id,
				origin_id,
				definition_id,
				is_sequence,
				is_parametrized, 
				reason, 
				legacy_template_id
			) VALUES (?,?,?,?,?,?,?,?)`
	params := []interface{}{
		data.PartnerID,
		data.LegacyScriptID,
		data.OriginID,
		data.DefinitionID,
		data.IsSequence,
		data.IsParametrized,
		data.ErrorReason,
		data.LegacyTemplateID,
	}

	if err := l.conn.Query(query, params...).Exec(); err != nil {
		return fmt.Errorf("can't insert legacy info profile: %v", err)
	}
	return nil
}
