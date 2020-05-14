package models

// DataBaseConnectors represents DTO object that contains all interfaces implemented in models package
type DataBaseConnectors struct {
	TemplateCache       TemplateCache
	Task                TaskPersistence
	UserSites           UserSitesPersistence
	TaskSummary         TaskSummaryPersistence
	TaskInstance        TaskInstancePersistence
	TaskDefinition      TaskDefinitionPersistence
	ExecutionResult     ExecutionResultPersistence
	ExecResultView      ExecutionResultViewPersistence
	ExecutionExpiration ExecutionExpirationPersistence
}
