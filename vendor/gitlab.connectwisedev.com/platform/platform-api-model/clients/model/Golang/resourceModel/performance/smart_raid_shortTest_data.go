package performance

import (
	"time"
)

//SMARTShortTestAttributeValues contains SMART attributes
type SMARTShortTestAttributeValues struct {
	ID         int    `json:"id" cql:"id"`
	Name       string `json:"name" cql:"name"`
	Flag       string `json:"flag" cql:"flag"`
	Value      int    `json:"value" cql:"value"`
	Worst      int    `json:"worst" cql:"worst"`
	Threshold  int    `json:"threshold" cql:"threshold"`
	Stype      string `json:"type" cql:"type"`
	Updated    string `json:"updated" cql:"updated"`
	WhenFailed string `json:"whenFailed" cql:"when_failed"`
	RawValue   string `json:"rawValue" cql:"raw_value"`
}

//SMARTShortTestResults contains short test values
type SMARTShortTestResults struct {
	Number          string `json:"number" cql:"number"`
	Description     string `json:"description" cql:"description"`
	Status          string `json:"status" cql:"status"`
	Remaining       string `json:"remaining" cql:"remaining"`
	LifeTime        string `json:"lifeTime" cql:"life_time"`
	LbaOfFirstError string `json:"lbaOfFirstError" cql:"lba_of_first_error"`
}

//SMARTRaidShortTestDeviceData contains device info
type SMARTRaidShortTestDeviceData struct {
	DeviceType       string                          `json:"deviceType" cql:"device_type"`
	DeviceName       string                          `json:"deviceName" cql:"device_name"`
	SerialNumber     string                          `json:"serialNumber" cql:"serial_number"`
	DeviceModel      string                          `json:"deviceModel" cql:"device_model"`
	FirmwareVersion  string                          `json:"firmwareVersion" cql:"firmware_version"`
	HealthStatus     string                          `json:"healthStatus" cql:"health_status"`
	SMARTAvailable   string                          `json:"smartAvailable" cql:"smartavailable"`
	SMARTEnabled     string                          `json:"smartEnabled" cql:"smartenabled"`
	SMARTTestResults SMARTShortTestResults           `json:"smarttestresults" cql:"smarttestresults"`
	SMARTAttributes  []SMARTShortTestAttributeValues `json:"smartAttributes" cql:"smartattributes"`
}

//SMARTRaidShortTestData harware raid ShortTest data
type SMARTRaidShortTestData struct {
	PartnerID           string                         `json:"partnerID,omitempty"  cql:"partner_id"`
	SiteID              string                         `json:"siteID,omitempty"  cql:"site_id"`
	ClientID            string                         `json:"clientID,omitempty"  cql:"client_id"`
	EndpointID          string                         `json:"endpointID,omitempty"  cql:"endpoint_id"`
	CreatedBy           string                         `json:"createdBy" cql:"created_by"`
	LastAccessedDateUTC time.Time                      `json:"lastAccessedDateUTC"`
	CreateTimeUTC       time.Time                      `json:"createTimeUTC"`
	SMARTShortTestData  []SMARTRaidShortTestDeviceData `json:"smartShortTestData" cql:"smart_short_test_data"`
}
