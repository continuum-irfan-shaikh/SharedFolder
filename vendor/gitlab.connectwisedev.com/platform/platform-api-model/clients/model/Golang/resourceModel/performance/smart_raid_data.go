package performance

import (
	"time"
)

type SMARTValues struct {
	ID         int    `json:"id" cql:"id"`
	Flag       string `json:"flag" cql:"flag"`
	Value      int    `json:"value" cql:"value"`
	Worst      int    `json:"worst" cql:"worst"`
	Threshold  int    `json:"threshold" cql:"threshold"`
	Stype      string `json:"type" cql:"type"`
	Updated    string `json:"updated" cql:"updated"`
	WhenFailed string `json:"whenFailed" cql:"when_failed"`
	RawValue   string `json:"rawValue" cql:"raw_value"`
}

type SMARTRaidDeviceData struct {
	DeviceType           string      `json:"deviceType" cql:"device_type"`
	DeviceName           string      `json:"deviceName" cql:"device_name"`
	SerialNumber         string      `json:"serialNumber" cql:"serial_number"`
	ErrorCode            string      `json:"errorCode" cql:"error_code"`
	ErrorMessage         string      `json:"errorMessage" cql:"error_message"`
	Status               string      `json:"status" cql:"status"`
	Protocol             string      `json:"protocol" cql:"protocol"`
	SMARTAvailable       string      `json:"smartAvailable" cql:"smartavailable"`
	SMARTEnabled         string      `json:"smartEnabled" cql:"smartenabled"`
	RawReadErrorRate     SMARTValues `json:"rawReadErrorRate" cql:"raw_read_error_rate"`
	ReallocatedSectorCt  SMARTValues `json:"reallocatedSectorCt" cql:"reallocated_sector_ct"`
	SeekErrorRate        SMARTValues `json:"seekErrorRate" cql:"seek_error_rate"`
	SpinRetryCount       SMARTValues `json:"spinRetryCount" cql:"spin_retry_count"`
	CommandTimeout       SMARTValues `json:"commandTimeout" cql:"command_timeout"`
	TemperatureCelsius   SMARTValues `json:"temperatureCelsius" cql:"temperature_celsius"`
	HardwareEccRecovered SMARTValues `json:"hardwareEccRecovered" cql:"hardware_ecc_recovered"`
	CurrentPendingSector SMARTValues `json:"currentPendingSector" cql:"current_pending_sector"`
	OfflineUncorrectable SMARTValues `json:"offlineUncorrectable" cql:"offline_uncorrectable"`
	UdmaCrcErrorCount    SMARTValues `json:"udmaCrcErrorCount"  cql:"udma_crc_error_count"`
}

//SMARTRaidData harware raid
type SMARTRaidData struct {
	PartnerID           string                `json:"partnerID,omitempty"  cql:"partner_id"`
	SiteID              string                `json:"siteID,omitempty"  cql:"site_id"`
	ClientID            string                `json:"clientID,omitempty"  cql:"client_id"`
	EndpointID          string                `json:"endpointID,omitempty"  cql:"endpoint_id"`
	CreatedBy           string                `json:"createdBy" cql:"created_by"`
	LastAccessedDateUTC time.Time             `json:"lastAccessedDateUTC"`
	CreateTimeUTC       time.Time             `json:"createTimeUTC"`
	SMARTData           []SMARTRaidDeviceData `json:"smartData" cql:"smart_data"`
}
