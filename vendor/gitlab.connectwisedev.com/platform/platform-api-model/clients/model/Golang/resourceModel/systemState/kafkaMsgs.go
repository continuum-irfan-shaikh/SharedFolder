package systemState

import "time"

//SystemStateChange is the struct definition of kafka message for systemstate_change topic
type SystemStateChange struct {
	PartnerID     string    `json:"partner_id"`
	SiteID        string    `json:"site_id"`
	ClientID      string    `json:"client_id"`
	LegacyRegID   string    `json:"legacyReg_id"`
	EndpointID    string    `json:"endpoint_id"`
	EventType     string    `json:"type"`
	BootUpTimeUTC time.Time `json:"bootUpTime"`
}

