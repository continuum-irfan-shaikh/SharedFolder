package autoupdate

import (
	"time"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
)

//Status is for package installation status
type Status string

//InstallVariables is for package installation variables
type InstallVariables map[string]string

const (
	//Success for successful installation of package
	Success Status = "SUCCESS"
	//Failure for failed installation of package
	Failure Status = "FAILURE"
	//PartialSuccess for partial installation of package
	PartialSuccess Status = "PARTIAL_SUCCESS"
	//Created for creation of installation status
	Created Status = "CREATED"
	//MailboxFailed for failure of posting mailbox message
	MailboxFailed Status = "MAILBOX_FAILED"
	//AgentUpdateStarted for successful start of processing auto-update message
	AgentUpdateStarted Status = "AGENT_UPDATE_STARTED"
	//ManifestRetrieved for successful call of GET agent manifest
	ManifestRetrieved Status = "MANIFEST_RETRIEVED"
	//Downloading for downloading in progress
	Downloading Status = "DOWNLOADING"
	//DownloadSuccess for successful download
	DownloadSuccess Status = "DOWNLOAD_SUCCESS"
	//DownloadFailure for failed download
	DownloadFailure Status = "DOWNLOAD_FAILURE"
	//Installing for installationg in progress
	Installing Status = "INSTALLING"
	//InstallSuccess for successful install
	InstallSuccess Status = "INSTALL_SUCCESS"
	//InstallFailure for failed install
	InstallFailure Status = "INSTALL_FAILURE"
	//UninstallSuccess for successful uninstall
	UninstallSuccess Status = "UNINSTALL_SUCCESS"
	//UninstallFailure for failed uninstall
	UninstallFailure Status = "UNINSTALL_FAILURE"
	//InProgress for in progress installation packages
	InProgress Status = "INPROGRESS"
	//RestoreFailure status for when installation manager fails to restore old package
	RestoreFailure Status = "RESTORE_FAILURE"
	//Invalid status when package is not meant to be installed on the endpoint.
	Invalid Status = "INVALID"
	//InvokingIM for IM invocation
	InvokingIM Status = "INVOKING_IM"
	//InstallationManager app name
	InstallationManager string = "InstallationManager"
	// StatusCompleteIM - Mark installation manager complete
	StatusCompleteIM int = 1
	// StatusInProgressIM - Mark installation manager in-progress
	StatusInProgressIM int = 2
)

//ManifestStatus is a struct defining status for manifest and package installation
type ManifestStatus struct {
	PartnerID      string    `json:"partnerID,omitempty" cql:"partner_id"`
	ClientID       string    `json:"clientID,omitempty" cql:"client_id"`
	SiteID         string    `json:"siteID,omitempty" cql:"site_id"`
	EndpointID     string    `json:"endpointID,omitempty" cql:"endpoint_id"`
	AgentID        string    `json:"agentID,omitempty" cql:"agent_id"`
	RegID          string    `json:"regID,omitempty" cql:"reg_id"`
	OSName         string    `json:"osName,omitempty" cql:"os_name"`
	OSType         string    `json:"osType,omitempty" cql:"os_type"`
	OSVersion      string    `json:"osVersion,omitempty" cql:"os_version"`
	OSArch         string    `json:"osArch,omitempty" cql:"os_arch"`
	DCTimestampUTC time.Time `json:"dcTimestampUTC,omitempty" cql:"dc_timestamp_utc"`
	InstallationStatus
}

//InstallationStatus to store manifest status in agentcore sqlite database
type InstallationStatus struct {
	Version           string                      `json:"version,omitempty" gorm:"primary_key" cql:"version"`
	Status            Status                      `json:"status,omitempty" cql:"status"`
	ErrorCode         string                      `json:"errorCode,omitempty" cql:"error_code"`
	SubErrorCode      string                      `json:"subErrorCode,omitempty" cql:"sub_error_code"`
	StatusMessage     string                      `json:"statusMessage,omitempty" cql:"status_message"`
	MessageID         string                      `json:"messageID,omitempty" cql:"message_id"`
	Originator        agent.Originator            `json:"originator,omitempty" cql:"originator"`
	PackageStatus     []PackageInstallationStatus `json:"packageStatus,omitempty" cql:"installation_status" gorm:"foreignkey:ManifestVersion;association_foreignkey:Version"`
	AgentTimestampUTC time.Time                   `json:"agentTimestampUTC,omitempty" cql:"agent_timestamp_utc"`
	InstallRetryCount int                         `json:"installRetryCount,omitempty" cql:"install_retry_count"`
	TransactionID     string                      `json:"transactionID,omitempty" cql:"transaction_id"`
	ForceUpdate       bool                        `json:"forceUpdate" cql:"force_update"`
}

//PackageInstallationStatus is a struct defining the Installation status of a package on an endpoint
type PackageInstallationStatus struct {
	ManifestVersion       string           `json:"manifestversion,omitempty" gorm:"primary_key"`
	Name                  string           `json:"name,omitempty" cql:"name" gorm:"primary_key"`
	Type                  string           `json:"type,omitempty" cql:"type"`
	Status                Status           `json:"status,omitempty" cql:"status"`
	Version               string           `json:"version,omitempty" cql:"version"`
	ErrorCode             string           `json:"errorCode,omitempty" cql:"error_code"`
	SubErrorCode          string           `json:"subErrorCode,omitempty" cql:"sub_error_code"`
	StackTrace            string           `json:"stackTrace,omitempty" cql:"stack_trace"`
	InstallationVariables InstallVariables `json:"installationVariables,omitempty" cql:"install_variables" gorm:"type:blob"`
	TimestampUTC          time.Time        `json:"timestampUTC,omitempty" cql:"timestamp_utc"`
	InstallerPath         string           `json:"installerPath,omitempty"`
	SourceURL             string           `json:"sourceURL,omitempty" cql:"source_url"`
	AvailableURLs         string           `json:"availableURLs,omitempty" cql:"available_urls"` //pipe `|` separated URLs
	Operation             string           `json:"operation,omitempty" cql:"operation" gorm:"default:install"`
	DownloadRetryCount    int              `json:"downloadRetryCount,omitempty" cql:"download_retry_count"`
	ChecksumType          string           `json:"checksumType,omitempty" cql:"checksum_type"`
	ChecksumValue         string           `json:"checksumValue,omitempty" cql:"checksum_value"`
}

//ManifestStatusDetails is defining manifest staus details
type ManifestStatusDetails struct {
	ManifestStatuses []ManifestStatus `json:"manifestStatuses,omitempty"`
}

//InstallationManagerStatus struct to keep track of installation manager status
type InstallationManagerStatus struct {
	ApplicationName string `json:"appname" gorm:"primary_key"`
	InProgress      int    `json:"inprogress"`
}
