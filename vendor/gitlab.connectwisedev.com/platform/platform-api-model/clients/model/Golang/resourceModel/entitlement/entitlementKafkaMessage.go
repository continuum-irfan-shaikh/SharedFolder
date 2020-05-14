package entitlement

import (
	"time"
)

// Action constants
const (
	CreateAction = "CREATE"
	DeleteAction = "DELETE"
	RemoveAction = "REMOVE"
)

// FeatureChange is the struct definition of /resources/entitlement/entitlementMessege
type FeatureChange struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Action        string                 `json:"action"`
	Entitlements  []EndpointRelationship `json:"entitlements"`
	Timestamp     int64                  `json:"timestamp"`
	UID           string                 `json:"uid"`
	TransactionID string                 `json:"transactionID"`
}

// ProductAssignmentKafkaMsgStructure - kafka message structure
type ProductAssignmentKafkaMsgStructure struct {
	Action        string    `json:"action"`
	ActionDate    time.Time `json:"actionDate"`
	PartnerID     string    `json:"partnerId"`
	ClientID      string    `json:"clientId"`
	SiteID        string    `json:"siteId"`
	EndpointID    string    `json:"endpointId"`
	SKU           string    `json:"productSku"`
	TransactionID string    `json:"transactionID"`
}

// ProductAssignmentKafkaMsg message to sent into kafka when product assigned to client
type ProductAssignmentKafkaMsg struct {
	Msg string `json:"Message"`
}
