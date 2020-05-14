package vmware

// VMState represents virtual machine power state
type VMState string

// HostState represents virtual machine power state
type HostState string

const (
	// StateOff means that VM is powered off
	StateOff = VMState("poweredOff")

	// StateRunning means that VM is powered on
	StateRunning = VMState("poweredOn")

	// StateSuspended means that VM is suspended
	StateSuspended = VMState("suspended")

	// StateOff means that Host is powered off
	HostStateOff = HostState("poweredOff")

	// StateRunning means that Host is powered on
	HostStateRunning = HostState("poweredOn")

	// StateSuspended means that Host is suspended
	HostStateSuspended = HostState("suspended")

	// HostStateUnknown means that Host is not responding
	HostStateUnknown = HostState("unknown")
)

// VirtualMachine contains data about specific VM
type VirtualMachine struct {
	*EndpointInfo

	// UUID is machine unique ID
	UUID string `json:"uuid"`

	// Name is virtual machine name
	Name string `json:"name"`

	// Type represents virtual machine type (e.g "Some Linux x86")
	Type string `json:"type"`

	// LocationID is location ID
	LocationID string `json:"locationId"`

	// InstanceUUID is instance ID
	InstanceUUID string `json:"instanceUuid"`

	// State represents virtual machine state
	State VMState `json:"state"`

	// Details contains additional VM information
	//
	// Is nil if VM is powered off or doesn't have VMWare Tools installed
	Details *GuestInfo `json:"details,omitempty"`

	//name of vm in Self object. Need to map VM to host
	SelfName string `json:"-"`
}

// NIC represents vm network adapter
type NIC struct {
	// Name is network name
	Name string `json:"name"`

	// MAC is MAC address
	MAC string `json:"mac"`

	// Connected represents if adapter is connected to network
	Connected bool `json:"connected"`

	// IP is list of IPs assigned to interface (usually IPv6/IPv6 pair)
	IP []string `json:"ip,omitempty"`
}

// GuestInfo represents additional VM information provided by VMWare Tools
type GuestInfo struct {
	// OS is guest operating system name
	OS string `json:"os"`

	// Hostname is obviously guest VM hostname
	Hostname string `json:"hostname"`

	// IPAddress is IP addr
	IPAddress string `json:"ip"`

	// Network is list of network adapters
	Network []NIC `json:"network,omitempty"`
}
