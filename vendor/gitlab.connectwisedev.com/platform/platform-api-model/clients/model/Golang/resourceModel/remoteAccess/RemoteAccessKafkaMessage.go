package remoteAccess

import "time"

//KafkaMessageXLMI contains Remote Access Kafka Message FOR LMI
type KafkaMessageXLMI struct {
	Status             int    `json:"status"`
	Detail             string `json:"detail"`
	RemoteAccessHostID uint64 `json:"remoteAccessHostID"`
	OSLMI              uint64 `json:"oslmi"`
	OSSpec             uint64 `json:"osSpec"`
	TimeZoneIndex      uint64 `json:"timeZoneIndex"`
	CurrentPort        uint64 `json:"currentPort"`
	InstalledBuild     string `json:"installedBuild"`
	LicenseType        uint64 `json:"licenseType"`
	ProductBuild       uint64 `json:"productBuild"`
	ProductType        uint64 `json:"productType"`
	ProductVersion     string `json:"productVersion"`
	InstalledVersion   string `json:"installedVersion"`
	License            string `json:"license"`
}

//KafkaMessageLMI contains Remote Access Kafka Message FOR LMI
type KafkaMessageLMI struct {
	RemoteAccessHostID uint64 `json:"remoteAccessHostID"`
	OSLMI              uint64 `json:"oslmi"`
	OSSpec             uint64 `json:"osSpec"`
	TimeZoneIndex      uint64 `json:"timeZoneIndex"`
	CurrentPort        uint64 `json:"currentPort"`
	InstalledBuild     string `json:"installedBuild"`
	LicenseType        uint64 `json:"licenseType"`
	ProductBuild       uint64 `json:"productBuild"`
	ProductType        uint64 `json:"productType"`
	ProductVersion     string `json:"productVersion"`
	InstalledVersion   string `json:"installedVersion"`
	License            string `json:"license"`
}

//KafkaMessageCONTROL contains Remote Access Kafka Message FOR CONTROL
type KafkaMessageCONTROL struct {
	SessionID string `json:"sessionID"`
}

//KafkaMessage contains Remote Access Kafka Message FOR LMI and CONTROL
type KafkaMessage struct {
	Status     int                 `json:"status"`
	Detail     string              `json:"detail"`
	MsgLMI     KafkaMessageLMI     `json:"msgLMI"`
	MsgControl KafkaMessageCONTROL `json:"msgControl"`
}

//PackageStatusData contains Package Status Message
type PackageStatusData struct {
	Manifestversion string    `json:"manifestversion"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Version         string    `json:"version"`
	ErrorCode       string    `json:"errorCode"`
	StackTrace      string    `json:"stackTrace"`
	TimestampUTC    time.Time `json:"timestampUTC"`
	InstallerPath   string    `json:"installerPath"`
	SourceURL       string    `json:"sourceURL"`
	Operation       string    `json:"operation"`
}

//ManifestStatuses contains Manifest Status Message
type ManifestStatuses struct {
	PartnerID         string              `json:"partnerID"`
	ClientID          string              `json:"clientID"`
	SiteID            string              `json:"siteID"`
	EndpointID        string              `json:"endpointID"`
	AgentID           string              `json:"agentID"`
	RegID             string              `json:"regID"`
	OsName            string              `json:"osName"`
	OsType            string              `json:"osType"`
	OsVersion         string              `json:"osVersion"`
	OsArch            string              `json:"osArch"`
	DcTimestampUTC    time.Time           `json:"dcTimestampUTC"`
	Version           string              `json:"version"`
	Status            string              `json:"status"`
	MessageID         string              `json:"messageID"`
	PackageStatus     []PackageStatusData `json:"packageStatus"`
	AgentTimestampUTC time.Time           `json:"agentTimestampUTC"`
	InstallRetryCount int                 `json:"installRetryCount"`
}
