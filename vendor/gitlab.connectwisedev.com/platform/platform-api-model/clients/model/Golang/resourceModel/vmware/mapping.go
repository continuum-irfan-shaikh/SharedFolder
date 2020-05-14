package vmware

// EndpointInfo contains partner, site, client and endpoint information
type EndpointInfo struct {
	// PartnerID is partner ID
	PartnerID string `json:"partnerId,omitempty"`

	// SiteID is site ID
	SiteID string `json:"siteId,omitempty"`

	// ClientID is client ID
	ClientID string `json:"clientId,omitempty"`

	// EndpointID is endpoint ID
	EndpointID string `json:"endpointId,omitempty"`

	// HasAgent is has agent status
	HasAgent bool `json:"hasAgent"`

	// Monitored is monitor status
	Monitored bool `json:"monitored"`
}
