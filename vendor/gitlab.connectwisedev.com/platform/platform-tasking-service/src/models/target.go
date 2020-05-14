package models

import (
	"encoding/json"
	"fmt"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

type (
	// TargetData represents the count and type of entities the task was run on (ManagedEndpoint, DynamicGroup etc)
	TargetData struct {
		Count int        `json:"count"`
		Type  TargetType `json:"-"` // deprecated
	}

	// ManagedEndpointDetailed struct is used to display detailed info for each ManagedEndpoint in TaskOutput
	ManagedEndpointDetailed struct {
		apiModels.ManagedEndpoint
		State statuses.TaskState `json:"state"`
	}

	// Target struct is used to describe Targets of the Task and their type
	Target struct {
		IDs  []string   `json:"ids"     valid:"requiredUniqueTargetIDs"`
		Type TargetType `json:"type"`
	}

	// TargetType is used for targets definition
	TargetType int
)

// These constants describe the type of entities the task was run on
const (
	_ TargetType = iota
	ManagedEndpoint
	DynamicGroup
	Site
	DynamicSite
)

const (
	endpoint    = apiModels.ManagedEndpointStr
	dg          = apiModels.DynamicGroupStr
	site        = apiModels.SiteStr
	dynamicSite = apiModels.DynamicSiteStr
)

// UnmarshalJSON used to convert the string representation of the type of the Targets to TargetType type
func (targetType *TargetType) UnmarshalJSON(byteResult []byte) error {
	var stringValue string
	if err := json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}

	types := map[string]TargetType{
		"":          TargetType(0),
		endpoint:    ManagedEndpoint,
		dg:          DynamicGroup,
		site:        Site,
		dynamicSite: DynamicSite,
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
		TargetType(0):   "",
		ManagedEndpoint: endpoint,
		DynamicGroup:    dg,
		Site:            site,
		DynamicSite:     dynamicSite,
	}

	stringValue, ok := types[targetType]
	if !ok {
		return []byte{}, fmt.Errorf("incorrect Type of ManagedEndpointID: %v", targetType)
	}
	return json.Marshal(stringValue)
}
