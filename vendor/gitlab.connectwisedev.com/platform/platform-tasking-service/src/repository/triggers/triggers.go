package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gocql/gocql"
)

// These keys represents triggers cache key
const (
	TriggerIDKeyPrefix          = "TKS_TRIGGER_ID_"
	TriggerPartnersPrefix       = "TKS_TRIGGER_PARTNERS"
	TriggerPartnersByTypePrefix = "TKS_TRIGGER_PARTNERS_BY_TYPE_"
)

const (
	alertTriggerCategory = "alert"
	shortTTL             = 60
)

type activeTriggersMap map[string][]e.ActiveTrigger
type activeTriggersPartners map[string]struct{}

// Repo is a triggers repository implementation for cassandra
type Repo struct {
	conn     cassandra.ISession
	memCache memcached.Cache
}

// NewTriggersRepo returns new cassandra repo for triggers
func NewTriggersRepo(conn cassandra.ISession, mem memcached.Cache) *Repo {
	return &Repo{
		conn:     conn,
		memCache: mem,
	}
}

// LoadTriggersToCache every minute loads triggers data to local cache
func (t *Repo) LoadTriggersToCache() (err error) {
	trs, err := t.GetAll()
	if err != nil {
		return
	}

	partners := make(activeTriggersPartners)
	partnersByType := make(map[string]activeTriggersPartners)

	groupedTriggers := t.groupTriggers(trs)
	for partnerID, triggersByType := range groupedTriggers {
		partners[partnerID] = struct{}{}
		for triggerType, activeTriggers := range triggersByType {
			if _, ok := partnersByType[triggerType]; !ok {
				partnersByType[triggerType] = make(activeTriggersPartners)
			}
			partnersByType[triggerType][partnerID] = struct{}{}

			b, err := json.Marshal(activeTriggers)
			if err != nil {
				return err
			}

			//third cache level
			if err = t.memCache.Set(&memcache.Item{
				Key:        TriggerIDKeyPrefix + partnerID + triggerType,
				Value:      b,
				Expiration: int32(time.Now().Unix() + int64(config.Config.Memcached.DefaultDataTTLSec)),
			}); err != nil {
				return err
			}
		}
	}

	return t.setToFirstAndSecondCacheLevel(partners, partnersByType)
}

func (t *Repo) setToFirstAndSecondCacheLevel(activeTriggersPartners activeTriggersPartners, partnersByType map[string]activeTriggersPartners) error {
	b, err := json.Marshal(activeTriggersPartners)
	if err != nil {
		return err
	}

	//first cache level
	if err = t.memCache.Set(&memcache.Item{
		Key:        TriggerPartnersPrefix,
		Value:      b,
		Expiration: int32(time.Now().Unix() + int64(shortTTL)),
	}); err != nil {
		return err
	}

	b, err = json.Marshal(partnersByType)
	if err != nil {
		return err
	}

	for triggerType, partners := range partnersByType {
		b, err := json.Marshal(partners)
		if err != nil {
			return err
		}

		//second cache level
		if err = t.memCache.Set(&memcache.Item{
			Key:        TriggerPartnersByTypePrefix + triggerType,
			Value:      b,
			Expiration: int32(time.Now().Unix() + int64(shortTTL)),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (t *Repo) groupTriggers(triggers []e.ActiveTrigger) (groupedTriggers map[string]activeTriggersMap) {
	// first group by partner
	triggersByPartner := make(activeTriggersMap)
	for _, tr := range triggers {
		triggersByPartner[tr.PartnerID] = append(triggersByPartner[tr.PartnerID], tr)
	}

	// then group by partner+type
	groupedTriggers = make(map[string]activeTriggersMap)
	for partner, trs := range triggersByPartner {
		mapByType := make(activeTriggersMap)
		for _, trigger := range trs {
			mapByType[trigger.Type] = append(mapByType[trigger.Type], trigger)
		}
		groupedTriggers[partner] = mapByType
	}
	return
}

// GetAll returns all active triggers
func (t *Repo) GetAll() ([]e.ActiveTrigger, error) {
	var (
		trigger  e.ActiveTrigger
		trs      = make([]e.ActiveTrigger, 0)
		iterator = t.conn.Query(`SELECT 
									type,
									partner_id, 
									task_id,
									start_time_frame, 
									end_time_frame
			                     FROM active_triggers`).Iter()
	)

	for iterator.Scan(
		&trigger.Type,
		&trigger.PartnerID,
		&trigger.TaskID,
		&trigger.StartTimeFrame,
		&trigger.EndTimeFrame,
	) {
		trs = append(trs, trigger)
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf("GetAll: error while working with found entities: %v", err)
	}
	return trs, nil
}

// GetAllByType returns all active triggers by partner and trigger type
func (t *Repo) GetAllByType(ctx context.Context, typeTrigger string, partnerID string, fromCache bool) ([]e.ActiveTrigger, error) {
	// first try to get it from cache
	if fromCache {
		activeTriggers, err := t.getFromCache(ctx, partnerID, typeTrigger)
		if err != nil {
			return activeTriggers, err
		}

		if len(activeTriggers) != 0 {
			return activeTriggers, nil
		}
	}

	return t.getByTypeAndPartnerFromDB(partnerID, typeTrigger)
}

func (t *Repo) getByTypeAndPartnerFromDB(partnerID string, typeTrigger string) ([]e.ActiveTrigger, error) {
	var (
		trigger  e.ActiveTrigger
		trs      = make([]e.ActiveTrigger, 0)
		iterator = t.conn.Query(`SELECT 
									type,
									partner_id,
									task_id,
									start_time_frame, 
									end_time_frame
							    FROM active_triggers WHERE partner_id = ? AND type = ?`, partnerID, typeTrigger).Iter()
	)

	for iterator.Scan(
		&trigger.Type,
		&trigger.PartnerID,
		&trigger.TaskID,
		&trigger.StartTimeFrame,
		&trigger.EndTimeFrame,
	) {
		trs = append(trs, trigger)
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf("GetAllByType: error while working with found entities: %v", err)
	}
	return trs, nil
}

func (t *Repo) getFromCache(ctx context.Context, partnerID string, typeTrigger string) ([]e.ActiveTrigger, error) {
	//check if partner has active triggers
	bin, err := t.memCache.Get(TriggerPartnersPrefix)
	if err == nil {
		var partners activeTriggersPartners
		if err = json.Unmarshal(bin.Value, &partners); err != nil {
			return nil, err
		}

		if _, ok := partners[partnerID]; !ok {
			return []e.ActiveTrigger{}, nil
		}

	} else {
		logger.Log.WarnfCtx(ctx, "GetAllByType Couldn't get list of partners that has active triggers from memcache, err: %v", err)
	}

	//check if partner has active triggers of particular type
	bin, err = t.memCache.Get(TriggerPartnersByTypePrefix + typeTrigger)
	if err == nil {
		var partners activeTriggersPartners
		if err = json.Unmarshal(bin.Value, &partners); err != nil {
			return nil, err
		}

		if _, ok := partners[partnerID]; !ok {
			return []e.ActiveTrigger{}, nil
		}
	} else {
		logger.Log.WarnfCtx(ctx, "GetAllByType Couldn't get list of partners that has active triggers of %s type from memcache, err: %v", typeTrigger, err)
	}

	bin, err = t.memCache.Get(TriggerIDKeyPrefix + partnerID + typeTrigger)
	if err == nil {
		var activeTriggers []e.ActiveTrigger
		if err = json.Unmarshal(bin.Value, &activeTriggers); err != nil {
			return nil, err
		}

		return activeTriggers, err
	}

	logger.Log.WarnfCtx(ctx, "GetAllByType Couldn't get active triggers from memcache, err: %v", err)
	return nil, nil
}

// GetAllByTaskID returns active triggers by partner and taskID
func (t *Repo) GetAllByTaskID(partnerID string, taskID gocql.UUID) ([]e.ActiveTrigger, error) {
	var (
		trigger  e.ActiveTrigger
		trs      = make([]e.ActiveTrigger, 0)
		iterator = t.conn.Query(`SELECT type, partner_id, task_id, start_time_frame, end_time_frame FROM active_triggers WHERE partner_id = ? AND task_id = ? ALLOW FILTERING`, partnerID, taskID).Iter()
	)

	for iterator.Scan(
		&trigger.Type,
		&trigger.PartnerID,
		&trigger.TaskID,
		&trigger.StartTimeFrame,
		&trigger.EndTimeFrame,
	) {
		trs = append(trs, trigger)
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf("GetAllByType: error while working with found entities: %v", err)
	}
	return trs, nil
}

// Insert implements trigger insert for cassandra
func (t *Repo) Insert(trigger e.ActiveTrigger) error {
	return t.conn.Query(`INSERT INTO active_triggers (
							type, 
							partner_id,
							task_id,
							start_time_frame,
							end_time_frame
						) VALUES (?, ?, ?, ?, ?)`,
		trigger.Type,
		trigger.PartnerID,
		trigger.TaskID,
		trigger.StartTimeFrame,
		trigger.EndTimeFrame,
	).Exec()
}

// Delete implements trigger delete for cassandra
func (t *Repo) Delete(trigger e.ActiveTrigger) error {
	return t.conn.Query(`DELETE FROM active_triggers WHERE 
			type = ? AND 
			partner_id = ? AND 
			task_id = ?`,
		trigger.Type,
		trigger.PartnerID,
		trigger.TaskID,
	).Exec()
}

// InsertDefinitions inserts triggers definitions into corresponding table in cassandra
func (t *Repo) InsertDefinitions(defs []e.TriggerDefinition) error {
	for _, def := range defs {
		if def.TriggerCategory == alertTriggerCategory {
			def.ID = triggers.AlertTypePrefix + def.ID
		} else {
			def.ID = triggers.GenericTypePrefix + def.ID
		}

		bytes, err := json.Marshal(def)
		if err != nil {
			return err
		}

		fields := []interface{}{
			def.ID,
			def.DisplayName,
			def.Description,
			string(bytes),
		}
		if err = t.conn.Query(`INSERT INTO trigger_definitions (
								id, 
								name,
								description,
								data
							 ) VALUES (?, ?, ?, ?)`, fields...).
			Exec(); err != nil {
			return err
		}
	}
	return nil
}

// TruncateDefinitions truncate trigger definitions cassandra table
func (t *Repo) TruncateDefinitions() error {
	return t.conn.Query(`TRUNCATE trigger_definitions`).Exec()
}

// GetAllDefinitionsNamesAndIDs return all trigger types stored in cassandra
func (t *Repo) GetAllDefinitionsNamesAndIDs() ([]e.TriggerDefinition, error) {
	var (
		trigger  e.TriggerDefinition
		trs      = make([]e.TriggerDefinition, 0)
		iterator = t.conn.Query(`SELECT 
									id, 
									name 
								FROM trigger_definitions`).Iter()
	)

	for iterator.Scan(&trigger.ID, &trigger.DisplayName) {
		trs = append(trs, trigger)
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf("GetAllDefinitionsNamesAndIDs: error while working with found entities: %v", err)
	}
	return trs, nil
}

// GetAllDefinitions return all trigger types stored in cassandra
func (t *Repo) GetAllDefinitions() ([]e.TriggerDefinition, error) {
	var (
		rawData  = ""
		trs      = make([]e.TriggerDefinition, 0)
		iterator = t.conn.Query(`SELECT data FROM trigger_definitions`).Iter()
	)

	for iterator.Scan(&rawData) {
		var trigger e.TriggerDefinition
		if err := json.Unmarshal([]byte(rawData), &trigger); err != nil {
			return []e.TriggerDefinition{}, err
		}
		trs = append(trs, trigger)
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf("GetAllDefinitions: error while working with found entities: %v", err)
	}
	return trs, nil
}

// GetDefinition definition by triggerType
func (t *Repo) GetDefinition(triggerType string) (e.TriggerDefinition, error) {
	var (
		rawData = ""
		trigger = e.TriggerDefinition{}
		err     = t.conn.Query(`SELECT data FROM trigger_definitions WHERE id = ?`, triggerType).Scan(&rawData)
	)
	if err != nil {
		return e.TriggerDefinition{}, err
	}
	if err = json.Unmarshal([]byte(rawData), &trigger); err != nil {
		return e.TriggerDefinition{}, err
	}
	return trigger, nil
}

// GetTriggerCounterByType returns trigger counter for a given type
func (t *Repo) GetTriggerCounterByType(triggerType string) (counter e.TriggerCounter, err error) {
	err = t.conn.Query(`SELECT
							trigger_id,
							policy_id,
							count 
						FROM trigger_policy_counter WHERE trigger_id = ?`, triggerType).
		Scan(
			&counter.TriggerID,
			&counter.PolicyID,
			&counter.Count,
		)
	if err != nil && err == gocql.ErrNotFound {
		return e.TriggerCounter{}, nil
	}
	return
}

// IncreaseCounter increases trigger counter for a giver trigger
func (t *Repo) IncreaseCounter(counter e.TriggerCounter) error {
	return t.conn.Query(`UPDATE trigger_policy_counter SET count = count + 1 
							WHERE trigger_id = ? AND policy_id = ?`, counter.TriggerID, counter.PolicyID).Exec()
}

// DecreaseCounter decreases trigger counter for a giver trigger
func (t *Repo) DecreaseCounter(counter e.TriggerCounter) error {
	return t.conn.Query(`UPDATE trigger_policy_counter SET count = count - 1 
							WHERE trigger_id = ? AND policy_id = ?`, counter.TriggerID, counter.PolicyID).Exec()
}
