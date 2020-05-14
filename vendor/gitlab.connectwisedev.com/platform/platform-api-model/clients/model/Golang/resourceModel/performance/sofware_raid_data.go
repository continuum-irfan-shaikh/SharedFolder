package performance

import (
	"time"
)

//SoftwareRaidData is data model to Hold software raid data
type SoftwareRaidData struct {
	PartnerID           string    `json:"partnerID,omitempty"  cql:"partner_id"`
	ClientID            string    `json:"clientID,omitempty"  cql:"client_id"`
	SiteID              string    `json:"siteID,omitempty"  cql:"site_id"`
	EndpointID          string    `json:"endpointID,omitempty"  cql:"endpoint_id"`
	Status              int       `json:"status" cql:"status"`
	RebootCount         int       `json:"rebootCount" cql:"reboot_count"`
	Type                string    `json:"type,omitempty"`
	Uptime              string    `json:"uptime,omitempty" cql:"uptime"`
	Data                string    `json:"data,omitempty" cql:"data"`
	LastAccessedDateUTC time.Time `json:"lastAccessedDateUTC,omitempty" `
	CreateTimeUTC       time.Time `json:"createTimeUTC,omitempty"`
	CreatedBy           string    `json:"createdBy,omitempty" cql:"created_by"`
}
