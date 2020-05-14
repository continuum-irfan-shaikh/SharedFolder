package relations

//ChildType ...
type ChildType struct {
	HostName    string        `cql:"host_name" json:"hostname,omitempty"`
	Type        string        `cql:"type" json:"type,omitempty"`
	HasAgent    bool          `cql:"has_agent" json:"hasAgent,omitempty"`
	IsMonitored bool          `cql:"is_monitered" json:"isMonitored,omitempty"`
	Networks    []NetworkType `cql:"networks" json:"networks,omitempty"`
	EndpointID  string        `cql:"endpoint_id" json:"endpointId,omitempty"`
	Name        string        `cql:"name" json:"name,omitempty"`
	State       string        `cql:"state" json:"state,omitempty"`
	OS          string        `cql:"os" json:"os,omitempty"`
	UUID        string        `cql:"child_uuid" json:"uuid,omitempty"`
}

//ParentType ...
type ParentType struct {
	HostName    string        `cql:"host_name" json:"hostname,omitempty"`
	Type        string        `cql:"type" json:"type,omitempty"`
	HasAgent    bool          `cql:"has_agent" json:"hasAgent,omitempty"`
	IsMonitored bool          `cql:"is_monitered" json:"isMonitored,omitempty"`
	Networks    []NetworkType `cql:"networks" json:"networks,omitempty"`
	EndpointID  string        `cql:"endpoint_id" json:"endpointId,omitempty"`
	UUID        string        `cql:"parent_uuid" json:"uuid,omitempty"`
	GUID        string        `cql:"parent_guid" json:"guid,omitempty"`
	Version     string        `cql:"version" json:"version,omitempty"`
}

// NetworkType ...
type NetworkType struct {
	Name       string `cql:"name" json:"name,omitempty"`
	MacAddress string `cql:"mac_address" json:"macAddress,omitempty"`
	Connected  bool   `cql:"connected" json:"connected,omitempty"`
	IPv4       string `cql:"ipv4" json:"ipv4,omitempty"`
	IPv6       string `cql:"ipv6" json:"ipv6,omitempty"`
}
