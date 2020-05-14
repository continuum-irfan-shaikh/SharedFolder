package v2

// TargetType represents target alias
type TargetType string

const (
	// ManagedEndpoint represents ME target type
	ManagedEndpoint TargetType = "MANAGED_ENDPOINT"

	// DynamicGroup represents DynamicGroup target type
	DynamicGroup TargetType = "DYNAMIC_GROUP"

	// Site represents Site target type
	Site TargetType = "SITE"
)

// Target represents struct with set of targets with same target type
type Target struct {
	Type TargetType `json:"type"    valid:"validTargetType"`
	IDs  []string   `json:"ids"     valid:"requiredUniqueTargetIDs,requiredForUsers"`
}
