package audit

import "time"

type status string

const (
	// NotApplicable - Status Not Applicable
	NotApplicable status = "N/A"

	// InProgress - Status In-Progress
	InProgress status = "IN-PROGRESS"

	// Success - Status Success
	Success status = "SUCCESS"

	// Failure - Status Failure
	Failure status = "FAILURE"
)

// Event - Audit Event to be published
type Event struct {
	ID                 string    `json:"id,omitempty"`
	Type               string    `json:"type,omitempty"`
	Description        string    `json:"description,omitempty"`
	PartnerID          string    `json:"partnerID,omitempty"`
	ClientID           string    `json:"clientID,omitempty"`
	SiteID             string    `json:"siteID,omitempty"`
	AgentID            string    `json:"agentID,omitempty"`
	EndpointID         string    `json:"endpointID,omitempty"`
	Address            string    `json:"address,omitempty"`
	MessageCode        string    `json:"messageCode,omitempty"`
	MessageDescription string    `json:"messageDescription,omitempty"`
	Object             string    `json:"object,omitempty"`
	StartTime          time.Time `json:"startTime,omitempty"`
	EndTime            time.Time `json:"endTime,omitempty"`
	Status             status    `json:"status,omitempty"`
	StatusCode         string    `json:"statusCode,omitempty"`
	StatusDescription  string    `json:"statusDescription,omitempty"`
	ErrorCode          string    `json:"errorCode,omitempty"`
	User               string    `json:"user,omitempty"`
	TransactionID      string    `json:"transactionID,omitempty"`
	Service            string    `json:"service,omitempty"`
	CreatedOn          time.Time `json:"createdOn,omitempty"`
}
