package dg

// DynamicGroupChange contains data for message about dynamic group that has been changed
type DynamicGroupChange struct {
	PartnerID         string `json:"partner_id"`
	ClientID          string `json:"client_id"`
	SiteID            string `json:"site_id"`
	ManagedEndpointID string `json:"endpoint_id"`
	DynamicGroupID    string `json:"dynamic_group_id"`
	Type              string `json:"type"`
}

// MonitoringDG structure contains data about dg that should be been started/stopped monitoring
type MonitoringDG struct {
	PartnerID          string `json:"partner_id"`
	DynamicGroupID     string `json:"dynamic_group_id"`
	ServiceID          string `json:"service_id"   description:"Identifier of a service/client requesting to monitor changes for this group"`
	Operation          string `json:"operation"    description:"One of: START-MONITORING,STOP-MONITORING"`
}
