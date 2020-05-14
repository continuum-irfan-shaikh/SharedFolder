package agent

import "time"

//Originator is a flow which triggers a mailbox message
type Originator string

const (
	//AssetCollectionChangeAutoupdate for asset collection change triggered flow of auto-update
	AssetCollectionChangeAutoupdate Originator = "AssetCollectionChangeAutoupdate"
	//EnableManifestAutoupdate for enable manifest flow of auto-update
	EnableManifestAutoupdate Originator = "EnableManifestAutoupdate"
	//NewEndpointAutoupdate for registering new endpoint flow of auto-update
	NewEndpointAutoupdate Originator = "NewEndpointAutoupdate"
	//SiteMigrationAutoupdate for site-migration triggered flow of auto-update
	SiteMigrationAutoupdate Originator = "SiteMigrationAutoupdate"
	//ProfileChangeAutoupdate for confioguration profile change triggered flow of auto-update
	ProfileChangeAutoupdate Originator = "ProfileChangeAutoupdate"
	//VersionRecheckAutoupdate for changes made by re-audit triggered by received version message
	VersionRecheckAutoupdate Originator = "VersionRecheckAutoupdate"
	//GatewayAutoupdate for when site's gateway is updated
	GatewayAutoupdate Originator = "GatewayAutoupdate"
	//ExternalReAuditTrigger for when audit is triggered manually via HTTP API
	ExternalReAuditTrigger Originator = "ExternalReAuditTrigger"
)

//MailboxMessage is the struct definition of /resources/agent/mailboxMessage
type MailboxMessage struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	SubType       string    `json:"subType"`
	Version       string    `json:"version"`
	TimestampUTC  time.Time `json:"timestampUTC"`
	Path          string    `json:"path"`
	Message       string    `json:"message"`
	MessageID     string
	Originator    Originator `json:"originator"`
	TransactionID string     `json:"transactionID"`
}

//Enum representing status of a mailbox message
const (
	//MailboxMsgStatusPending
	MailboxMsgStatusPending = 0
	//MailboxMsgStatusSeverSent message has been sent down to client
	MailboxMsgStatusSeverSent = 1
	//MailboxMsgStatusAgentProcessedSuccess denotes client has succesfully processed the message
	MailboxMsgStatusAgentProcessedSuccess = 7
	//MailboxMsgStatusAgentProcessedFailure denotes processing failed at client
	MailboxMsgStatusAgentProcessedFailure = 8
)

//MailboxMessageStatus is the struct definition of /resources/agent/mailboxMessage
type MailboxMessageStatus struct {
	StatusCode int `json:"statusCode"`
}
