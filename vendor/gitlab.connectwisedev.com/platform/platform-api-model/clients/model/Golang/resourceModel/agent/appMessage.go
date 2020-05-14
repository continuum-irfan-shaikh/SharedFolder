package agent

import "time"

//AppManagement is a struct defining the actual app management task
type AppManagement struct {
	Action          string `json:"action"`
	PackageName     string `json:"packageName"`
	ForceUpdate     bool   `json:"forceUpdate"`
	ManifestVersion string `json:"manifestVersion"`
}

//AppMessage is the struct definition of /resources/agent/appMessage
type AppMessage struct {
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	SubType      string     `json:"subType"`
	Version      string     `json:"version"`
	TimestampUTC time.Time  `json:"timestampUTC"`
	Path         string     `json:"path"`
	Originator   Originator `json:"originator"`
	AppManagement
	MessageID     string
	TransactionID string `json:"transactionID"`
}
