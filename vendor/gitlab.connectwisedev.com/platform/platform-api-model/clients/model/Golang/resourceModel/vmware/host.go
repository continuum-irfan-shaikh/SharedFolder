package vmware

import "time"

// HostType represents VMWare API host type
type HostType string

// EventType represents type of Event: Reboot, Shutdown, ConnectionLost
type EventType string

const (
	// VCenterHost is vCenter host type
	VCenterHost = HostType("vCenter")

	// ESXiHost is ESXi host type
	ESXiHost = HostType("ESXi")

	// HypervHost is the Hyper-V host type
	HypervHost = HostType("Hyper-V")

	// Reboot is the EventType of event that comes from vCenter
	Reboot = EventType("Reboot")

	// Shutdown is the EventType of event that comes from vCenter
	Shutdown = EventType("Shutdown")

	// ConnectionLost is the EventType of event that comes from vCenter
	ConnectionLost = EventType("ConnectionLost")
)

// NetworkMap is map with network interfaces.
// Key is NIC MAC address and value is NIC IP's
//
// Example:
//
//		NetworkMap{
//			"e0:d5:5e:1d:92:4e": [
// 				"10.31.39.212",
// 				"fe80::e2d5:5eff:fe1d:924e"
// 			],
//		}
//
type NetworkMap map[string][]string

// Servers is an array of servers monitored by vmware plugin.
type Servers []Server

// ProductInfo represents VMWare product information.
type ProductInfo struct {
	Type    HostType `json:"type,omitempty"`
	Version string   `json:"version,omitempty"`
}

// Server represents VMWare API host used for monitoring by the plugin.
//
// It can be vCenter or ESXi itself.
type Server struct {
	ProductInfo

	// UUID is uuid assigned by Server Agent
	UUID string `json:"uuid"`

	// GUID is uuid assigned by Server Agent
	GUID string `json:"guid"`

	// Address is network address used by the plugin for monitoring
	Address string `json:"address"`

	// Hosts is array of ESXi/vCenter hosts managed by VMWare server
	//
	// If server is vCenter, it will contain array of managed ESXi hosts.
	// If server is standalone ESXi instance - it will contain information about itself.
	Hosts []HostSystem `json:"hosts,omitempty"`
}

// HostInfo represents host system information
type HostInfo struct {
	// Hostname is host name
	Hostname string `json:"hostname"`

	// Network is map of network
	Network NetworkMap `json:"network,omitempty"`
}

// HostSystem represents server (usually ESXi) managed by vCenter
type HostSystem struct {
	*EndpointInfo
	ProductInfo

	// Details contains information about the host
	Details *HostInfo `json:"details,omitempty"`
	//PoweredOn or PoweredOff state
	State HostState `json:"state"`
	UUID  string    `json:"uuid"`

	// NodeName is the Hyper-V cluster node name field
	NodeName string `json:"nodeName,omitempty"`

	// Children contains children hosts managed by the host.
	//
	// Applicable if the host is vCenter server.
	Children []HostSystem `json:"children,omitempty"`

	// VirtualMachines is a list of virtual machines running on the host.
	VirtualMachines         []VirtualMachine `json:"virtualMachines,omitempty"`
	VirtualMachineSelfNames []string         `json:"-"`
}

type HostStatus struct {
	Type              HostType  `json:"type"`
	EventType         EventType `json:"eventType"`
	UUID              string    `json:"uuid"`
	Hostname          string    `json:"hostname"`
	State             HostState `json:"state"`
	EventTimestampUtc time.Time `json:"eventTimestampUtc"`
}
