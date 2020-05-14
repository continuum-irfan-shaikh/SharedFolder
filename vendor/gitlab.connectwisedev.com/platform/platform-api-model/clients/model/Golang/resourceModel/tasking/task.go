package tasking

import (
	"encoding/json"
	"strings"
	"time"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"github.com/gocql/gocql"
)

type Task struct {
	Name               string                 `json:"name"                   valid:"optional"`
	Description        string                 `json:"description"            valid:"-"`
	ResultWebhook      string                 `json:"resultWebhook"          valid:"url"`
	IsRequireNOCAccess bool                   `json:"isRequireNOCAccess"     valid:"-"`
	Targets            Target                 `json:"targets"                valid:"validTargets"` // deprecated but used field
	TargetsByType      TargetsByType          `json:"targetsByType"          valid:"-"`            // targets related data is stored here
	ExternalTask       bool                   `json:"externalTask"           valid:"-"`
	Credentials        *agent.Credentials     `json:"credentials,omitempty"  valid:"validCreds"`
	Type               string                 `json:"type"                   valid:"validType"`
	Schedule           Schedule               `json:"schedule"               valid:"required,validatorDynamicGroup,optionalTriggerTypes,recurrentDGTriggerTarget"`
	DefinitionID       gocql.UUID             `json:"definitionID"           valid:"optional"`
	OriginID           gocql.UUID             `json:"originId"               valid:"requiredForUsers"` // script or patch ID
	ResourceType       ResourceType           `json:"resourceType"           valid:"validResourceType"`
	ParametersObject   map[string]interface{} `json:"parametersObject,omitempty" valid:"-"` // same as parameters but used only for request. data from this field should always be in the parameters field
	Parameters         string                 `json:"parameters"             valid:"-"`

	// unsettableByUsers fields
	ID                  gocql.UUID                `json:"id"                     valid:"unsettableByUsers"`
	CreatedAt           time.Time                 `json:"createdAt"              valid:"unsettableByUsers"`
	CreatedBy           string                    `json:"createdBy"              valid:"unsettableByUsers"`
	PartnerID           string                    `json:"partnerId"              valid:"unsettableByUsers"`
	State               TaskState                 `json:"state"                  valid:"unsettableByUsers"`
	RunTimeUTC          time.Time                 `json:"nextRunTime"            valid:"unsettableByUsers"` // special field for scheduler
	PostponedRunTime    time.Time                 `json:"postponedTime"          valid:"unsettableByUsers"`
	OriginalNextRunTime time.Time                 `json:"originalNextRunTime"    valid:"unsettableByUsers"`
	TargetType          TargetType                `json:"targetType"             valid:"unsettableByUsers"` // target type per internal task/endpointID
	ManagedEndpoints    []ManagedEndpointDetailed `json:"managedEndpoints"       valid:"unsettableByUsers"`
	ModifiedBy          string                    `json:"modifiedBy"             valid:"unsettableByUsers"`
	ModifiedAt          time.Time                 `json:"modifiedAt"             valid:"unsettableByUsers"`
}

// ManagedEndpointDetailed struct is used to display detailed info for each ManagedEndpoint in TaskOutput
type ManagedEndpointDetailed struct {
	ManagedEndpoint
	State TaskState `json:"state"`
}

// TargetsByType  represents DTO of targets grouped by type
type TargetsByType map[TargetType][]string

// Target struct is used to describe Targets of the Task and their type
type Target struct {
	IDs  []string   `json:"ids"     valid:"requiredUniqueTargetIDs"`
	Type TargetType `json:"type"`
}

// ResourceType represents endpoint resource type
type ResourceType string

const (
	// Desktop is a Desktop type
	Desktop ResourceType = "Desktop"
	// Server is a Server type
	Server ResourceType = "Server"
)

const anyTimesReplace = -1
const escape = "\""

// UnmarshalJSON used to convert the string representation of the type of the Targets to TargetType type in TargetsByType key
func (t *TargetsByType) UnmarshalJSON(byteResult []byte) error {
	input := make(map[string][]string)

	if err := json.Unmarshal(byteResult, &input); err != nil {
		return err
	}

	if len(input) == 0 {
		return nil
	}

	targets := make(map[TargetType][]string)
	for k, v := range input {
		tt := TargetType(0)

		if err := json.Unmarshal([]byte(escape+k+escape), &tt); err != nil {
			return err
		}

		targets[tt] = v
	}

	*t = targets

	return nil
}

// MarshalJSON custom marshal method for TargetsByType type
func (t TargetsByType) MarshalJSON() ([]byte, error) {
	output := make(map[string][]string)

	for k, v := range t {
		bType, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}

		tt := strings.Replace(string(bType), escape, "", anyTimesReplace)
		output[tt] = v
	}

	return json.Marshal(output)
}
