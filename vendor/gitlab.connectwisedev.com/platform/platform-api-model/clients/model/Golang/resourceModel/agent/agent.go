package agent

import "time"


// CryptoPassword is the password will be going to use in agent authentication
const CryptoPassword string = "v1(m$p@55\\/\\/0rd"

//Agent is the struct definition of /resources/agent/agent
type Agent struct {
	TimeStampUTC time.Time   `json:"timeStampUTC"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Components   []Component `json:"components"`
	Plugins      []Plugin    `json:"plugins"`
}
