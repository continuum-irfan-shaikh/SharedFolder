package entities

import (
	"time"
)

// TaskExecHistory - data structure to hold task execution history results
type TaskExecHistory struct {
	ExecYear      string
	ExecMonth     string
	ExecDate      string
	ExecTime      time.Time
	EndpointID    string
	ScriptName    string
	ScriptID      string
	CompletedTime time.Time
	ExecStatus    string
	PartnerID     string
	SiteID        string
	MachineName   string
	ExecBy        string
	Output        string
}
