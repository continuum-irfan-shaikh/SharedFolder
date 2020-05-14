package acms

import (
	"encoding/json"
	"time"

	"github.com/gocql/gocql"
)

// Profile contains set of deltas and conditions on which it applies to an endpoint
type Profile struct {
	PartnerID             string               `json:"-"`
	ID                    gocql.UUID           `json:"id"`
	Description           string               `json:"description,omitempty"`
	Tag                   string               `json:"tag, valid:"required~profile tag should not be empty""`
	Targets               []ExtendedEndpointID `json:"targets,omitempty" cql:"targets"`
	Condition             string               `json:"condition,omitempty"` //see: https://github.com/Knetic/govaluate/README.md and profile_test.go
	Deltas                []Delta              `json:"configurations,omitempty" cql:"deltas"`
	ConfigurationTemplate string               `json:"configurationTemplate,omitempty"`
	CreatedTime           time.Time            `json:"createdTime,omitempty" cql:"created_at"`
	CreatedBy             string               `json:"createdBy,omitempty" cql:"created_by"`
	ModifiedTime          time.Time            `json:"modifiedTime,omitempty" cql:"modified_at"`
	ModifiedBy            string               `json:"modifiedBy,omitempty" cql:"modified_by"`
}

// Delta contains patch that should be applied to a file on an endpoint
type Delta struct {
	PackageName         string          `json:"packageName" cql:"package_name" valid:"required~configuration package name is mandatory"`
	FileName            string          `json:"fileName,omitempty" cql:"file_name"`
	MinSupportedVersion string          `json:"minSupportedVersion,omitempty" cql:"min_supported_version"`
	MaxSupportedVersion string          `json:"maxSupportedVersion,omitempty" cql:"max_supported_version"`
	Operation           string          `json:"operation,omitempty" cql:"operation"`
	Patch               json.RawMessage `json:"patch,omitempty" cql:"patch"`
}

// ExtendedEndpointID contains partner id, client id, site id and endpoint id
type ExtendedEndpointID struct {
	PartnerID  string      `json:"partnerID,omitempty" cql:"partner_id"`
	ClientID   string      `json:"clientID,omitempty" cql:"client_id"`
	SiteID     string      `json:"siteID,omitempty" cql:"site_id"`
	EndpointID *gocql.UUID `json:"endpointID,omitempty" cql:"endpoint_id"`
}
