package eventlog

import (
	"time"
)

//EventLogMessage is the struct definition of /resources/eventLog/eventLogMessage
type EventLogMessage struct {
	Hostname     string    `json:"hostname"`
	Source       string    `json:"source"`
	Channel      string    `json:"channel"`
	Facility     int       `json:"facility"`
	Severity     int       `json:"severity"`
	EventID      int       `json:"eventID"`
	Message      string    `json:"message"`
	XMLView      string    `json:"xmlView"`
	CreatedAt    time.Time `json:"createdAt"`
	Duplications int       `json:"duplications"`
}
