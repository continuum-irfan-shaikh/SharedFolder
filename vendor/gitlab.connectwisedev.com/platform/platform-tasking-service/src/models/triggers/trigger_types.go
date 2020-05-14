package triggers

const (
	// GenericTypePrefix says about generic trigger type (not alerts)
	GenericTypePrefix = "generic-"
	// AlertTypePrefix says about alert trigger type
	AlertTypePrefix = "alert-"

	// LogoutTrigger represents id of this trigger
	LogoutTrigger = GenericTypePrefix + "79dc260c-ad43-11e9-a2a3-2a2ae2dbcce4"

	// LoginTrigger represents id of this trigger
	LoginTrigger = GenericTypePrefix + "90763c24-ad41-11e9-a2a3-2a2ae2dbcce4"

	// StartupTrigger represents id of this trigger
	StartupTrigger = GenericTypePrefix + "76300c9e-a9f9-11e9-b79e-e4e74935e7a7"
	// DynamicGroupExitTrigger represents id of this trigger
	DynamicGroupExitTrigger = GenericTypePrefix + "5599cd6d-70c4-41ce-8e95-b4ae3eb745dd"

	// DynamicGroupEnterTrigger represents id of this trigger
	DynamicGroupEnterTrigger = GenericTypePrefix + "11d43b69-36bf-4467-8c81-e0e8f4d1990f"

	// ShutdownTrigger represents id of this trigger
	ShutdownTrigger = "ShutdownTrigger"

	// FirstCheckInTrigger represents id of this trigger
	FirstCheckInTrigger = GenericTypePrefix + "5981fe70-ca69-11e9-afc4-acde48001122"

	// MockGeneric trigger represents generic mock trigger type to retrieve mocked handlers
	MockGeneric = GenericTypePrefix + "mock"

	// MockAlerting trigger represents alerting mock trigger type to retrieve mocked handlers
	MockAlerting = AlertTypePrefix + "mock"
)
