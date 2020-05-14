package version

import (
	"time"
)

//AgentVersion holds version details of an Installed Agent on endpoint
type AgentVersion struct {
	CreatedBy         string      `json:"createdBy,omitempty"`
	Name              string      `json:"name"`
	Type              string      `json:"type"`
	PartnerID         string      `json:"partnerId,omitempty"`
	ClientID          string      `json:"clientId,omitempty"`
	SiteID            string      `json:"siteId,omitempty"`
	AgentID           string      `json:"agentId,omitempty"`
	EndpointID        string      `json:"endpointId,omitempty"`
	ProductVersion    string      `json:"productVersion,omitempty"`
	ManifestVersion   string      `json:"manifestVersion"`
	AgentTimestampUTC time.Time   `json:"agentTimestampUTC"`
	DCTimestampUTC    time.Time   `json:"dcTimestampUTC,omitempty"`
	Components        []Component `json:"components"`
	Files             []File      `json:"files"`
}

//Component holds version details of individula binaries on endpoint
type Component struct {
	Name           string    `json:"name" cql:"component_name"`
	Version        string    `json:"version" cql:"component_version"`
	LastModifiedOn time.Time `json:"lastModifiedOn" cql:"lastmodifiedon"`
}

// File holds file path relative to continuum installation folder
type File struct {
	Path     string `json:"path" cql:"file_path"`
	Checksum string `json:"checksum" cql:"file_checksum"`
}
