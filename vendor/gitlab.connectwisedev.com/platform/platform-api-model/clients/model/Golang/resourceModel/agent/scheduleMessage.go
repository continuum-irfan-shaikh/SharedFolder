package agent

import "time"

const (
	//CRON for cron schedule message type
	CRON = "cron"
	//EVENT for event schedule message type
	EVENT = "event"
)

//ScheduleTask is a struct defining the actual schedule task
type ScheduleTask struct {
	Task             string       `json:"task,omitempty"`
	TaskInput        string       `json:"taskInput,omitempty"`
	ExecuteNow       string       `json:"executeNow,omitempty"`
	ExecuteOnStartup bool         `json:"executeOnStartup,omitempty"`
	Schedule         string       `json:"schedule,omitempty"`
	Event            string       `json:"event,omitempty"`
	TimeoutInSeconds int          `json:"timeout,omitempty"`
	Credentials      *Credentials `json:"credentials,omitempty"`
}

// Credentials is details to run task with given credentials
type Credentials struct {
	UseCurrentUser bool   `cql:"use_current_user" json:"use_current_user,omitempty"`
	Username       string `cql:"username" json:"username,omitempty"`
	Domain         string `cql:"domain" json:"domain,omitempty"`
	Password       string `cql:"password" json:"password,omitempty"`
}

//ScheduleMessage is the struct definition of /resources/agent/scheduleMessage
type ScheduleMessage struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Version       string    `json:"version"`
	TimestampUTC  time.Time `json:"timestampUTC"`
	Path          string    `json:"path"`
	TransactionID string    `json:"transactionID"`
	ScheduleTask
}
