package agent

import "time"

// ProvisioningData is the struct definition of agent provisioning data structure
type ProvisioningData struct {
	EndpointMapping
	SysInfo               SystemInfo `json:"sysInfo,omitempty"`
	SysInfoAPI            SystemInfo `json:"sysInfoAPI,omitempty"`
	PublicIPAddress       string     `json:"publicIPAddress,omitempty"`
	Token                 string     `json:"token,omitempty"`
	AgentInstalledVersion string     `json:"agentInstalledVersion,omitempty"`
	DcTimestampUTC        time.Time  `json:"dcTimestampUTC,omitempty"`
	DcAgentInstalledTSUTC time.Time  `json:"dcAgentInstalledTSUTC,omitempty"`
	CollectionMethod      string     `json:"collectionMethod,omitempty"`
	AgentPublicKey        string     `json:"agentPublicKey,omitempty"`
	SBAUUID               string     `json:"sbauuid,omitempty"`
}

//EndpointMapping struct provides mapping between endpoint and its partner,client,site,agent.
//This will be returned back to agent as response to successful registration.
type EndpointMapping struct {
	EndpointID    string `json:"endpointID,omitempty"`
	AgentID       string `json:"agentID,omitempty"`
	AgentPublicID string `json:"agentPublicID,omitempty"`
	PartnerID     string `json:"partnerID,omitempty"`
	SiteID        string `json:"siteID,omitempty"`
	ClientID      string `json:"clientID,omitempty"`
	LegacyRegID   string `json:"legacyRegID,omitempty"`
	Installed     bool   `json:"installed,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

//MigrationData stores the old and new endpoint mapping values. Used in case of migration.
type MigrationData struct {
	Newmapping EndpointMapping `json:"newmapping,omitempty"`
	Oldmapping EndpointMapping `json:"oldmapping,omitempty"`
}

//EndpointIDPair contains pair of private endpoint ID and public endpoint ID
type EndpointIDPair struct {
	PrivateID string `json:"privateID,omitempty"`
	PublicID  string `json:"publicID,omitempty"`
}

//PublicPrivateMapping contains pair of private endpoint ID and public endpoint ID and install status
type PublicPrivateMapping struct {
	EndpointIDPair
	Active bool `json:"active"`
}

//PrivateEndpointMapping contain mapping of private endpointId
type PrivateEndpointMapping struct {
	EndpointID string `json:"privateID"`
	Active     bool   `json:"active"`
}
