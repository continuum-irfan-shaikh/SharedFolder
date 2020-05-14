package tasking

import (
	"encoding/json"
	"fmt"
)

// TargetType is used for targets definition
type TargetType int

// These constants describe the type of entities the task was run on
const (
	_ TargetType = iota
	ManagedEndpointType
	DynamicGroup
	Site
	DynamicSite
)

const (
	ManagedEndpointStr = "MANAGED_ENDPOINT"
	DynamicGroupStr    = "DYNAMIC_GROUP"
	SiteStr            = "SITE"
	Unknown            = "Unknown"
	DynamicSiteStr     = "DYNAMIC_SITE"
)

// UnmarshalJSON used to convert the string representation of the type of the Targets to TargetType type
func (targetType *TargetType) UnmarshalJSON(byteResult []byte) error {
	var stringValue string
	if err := json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}

	types := map[string]TargetType{
		"":                 TargetType(0),
		ManagedEndpointStr: ManagedEndpointType,
		DynamicGroupStr:    DynamicGroup,
		SiteStr:            Site,
		DynamicSiteStr:     DynamicSite,
	}

	var ok bool
	*targetType, ok = types[stringValue]
	if !ok {
		return fmt.Errorf("incorrect Type of Targets: %s", stringValue)
	}

	return nil
}

// MarshalJSON custom marshal method for TargetType type
func (targetType TargetType) MarshalJSON() ([]byte, error) {
	types := map[TargetType]string{
		TargetType(0):       "",
		ManagedEndpointType: ManagedEndpointStr,
		DynamicGroup:        DynamicGroupStr,
		Site:                SiteStr,
		DynamicSite:         DynamicSiteStr,
	}

	stringValue, ok := types[targetType]
	if !ok {
		return []byte{}, fmt.Errorf("incorrect Type of ManagedEndpointID: %v", targetType)
	}
	return json.Marshal(stringValue)
}

// Parse is used to Parse string TargetType to TargetType type
func (targetType *TargetType) Parse(s string) error {
	switch s {
	case "":
		*targetType = TargetType(0)
	case ManagedEndpointStr:
		*targetType = ManagedEndpointType
	case DynamicGroupStr:
		*targetType = DynamicGroup
	case SiteStr:
		*targetType = Site
	case DynamicSiteStr:
		*targetType = DynamicSite
	default:
		return fmt.Errorf("incorrect targetType: %s", s)
	}
	return nil
}

// String returns string representation
func (targetType TargetType) String() string {
	types := map[TargetType]string{
		TargetType(0):       "",
		ManagedEndpointType: ManagedEndpointStr,
		DynamicGroup:        DynamicGroupStr,
		Site:                SiteStr,
		DynamicSite:         DynamicSiteStr,
	}

	stringValue, ok := types[targetType]
	if !ok {
		return Unknown
	}
	return stringValue
}
