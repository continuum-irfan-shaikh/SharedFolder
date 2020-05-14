package ticket

// AlertClosure represent Alertclosure object from cherwell to alerting
type AlertClosure struct {
	TransactionID string `json:"TransactionID"`
	IncidentID    string `json:"IncidentID"`
	AlertID       string `json:"AlertID"`
	UpdatedOn     string `json:"UpdatedOn"`
}
