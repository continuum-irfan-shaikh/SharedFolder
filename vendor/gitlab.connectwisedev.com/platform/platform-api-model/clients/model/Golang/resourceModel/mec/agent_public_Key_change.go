package mec

const (
	//HeaderAgentPublicKeyChange is the type of message
	HeaderAgentPublicKeyChange string = "AGENTPUBLICKEYCHANGE"
)

//AgentPublicKeyChange is the message that will get published when an agent's public key gets changed
type AgentPublicKeyChange struct {
	EndpointID string `json:"endpointID,omitempty"`
	PartnerID  string `json:"partnerID,omitempty"`
	SiteID     string `json:"siteID,omitempty"`
	ClientID   string `json:"clientID,omitempty"`
}
