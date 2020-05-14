package performance

import (
	"time"
)

//SensorsInfo snesors Data
type Sensors struct {
	PartnerID           string        `json:"partnerID,omitempty"  cql:"partner_id"`
	ClientID            string        `json:"clientID,omitempty"  cql:"client_id"`
	SiteID              string        `json:"siteID,omitempty"  cql:"site_id"`
	EndpointID          string        `json:"endpointID,omitempty"  cql:"endpoint_id"`
	LastAccessedDateUTC time.Time     `json:"lastAccessedDateUTC"`
	CreateTimeUTC       time.Time     `json:"createTimeUTC"`
	CreatedBy           string        `json:"createdBy" cql:"created_by"`
	Info                []SensorsInfo `json:"sensorsInfo" cql:"info"`
}

//SensorsInfo stores the info about the sensors
type SensorsInfo struct {
	Name    string           `json:"sensorName" cql:"name"`
	Adapter string           `json:"adapter" cql:"adapter"`
	Type    []SensorInfoType `json:"type" cql:"types"`
}

//SensorType type of the sensors
type SensorInfoType struct {
	Name        string `json:"@typeName" cql:"name"`
	TempVol     string `json:"@tempVol" cql:"temp_vol"`
	Range       string `json:"@range" cql:"range"`
	AlarmStatus string `json:"@alarmStatus"  cql:"alarm_status"`
}
