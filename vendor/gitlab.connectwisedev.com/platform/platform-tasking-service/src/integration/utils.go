package integration

// ExternalClients is a DTO object that contains all external communication clients
type ExternalClients struct {
	Asset            Asset
	Sites            SitesConnector
	AgentConfig      AgentConfig
	DynamicGroups    DynamicGroups
	AutomationEngine AutomationEngine
	HTTP             HTTPClient
	AgentEncryption  AgentEncryptionService
}
