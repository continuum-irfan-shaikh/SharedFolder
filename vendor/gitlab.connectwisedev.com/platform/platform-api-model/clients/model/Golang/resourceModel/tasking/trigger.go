package tasking

// These are alerts type for the same alert id
const (
	AlertTypeDelete = "DELETE"
	AlertTypeCreate = "CREATE"
)

// TriggerExecutionPayload is a general payload that required for trigger execution
type TriggerExecutionPayload struct {
	DynamicGroupID string `json:"dynamicGroupId,omitempty"  valid:"-"`
	AlertType      string `json:"alertType,omitempty"       valid:"-"`
}
