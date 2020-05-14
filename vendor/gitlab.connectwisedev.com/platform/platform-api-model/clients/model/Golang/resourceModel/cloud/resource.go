package cloud

//Resource represents a cloud resource
type Resource struct {
	ClientID          string              `json:"clientid"`
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	ServiceType       string              `json:"service"`
	ServiceName       string              `json:"servicename"`
	Location          string              `json:"location"`
	AvailabilityState string              `json:"availabilitystatus"`
	Hierarchies       []ResourceHierarchy `json:"hierarchies"`
	HealthURL         string              `json:"healthURL"`
}

//ResourceHierarchy represents a Resource hierarchy
type ResourceHierarchy struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	ID    string `json:"id"`
	Level int    `json:"level"`
}

//ResourceConfiguration stores configuration information for the Resource
type ResourceConfiguration struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Location      string              `json:"location"`
	ServiceType   string              `json:"service"`
	Hierarchies   []ResourceHierarchy `json:"hierarchies"`
	ResProperties interface{}         `json:"properties,omitempty"`
	//VirtualMachine  *VirtualMachine
	//ResourceChild ResourceChildDependent `json:"resourcechild,omitempty"`
	//NetworkInterfaces *NetworkInterfaces `bson:"ni" json:"properties,omitempty"`
}

//VirtualMachine stores the virtual machine configuration information
type VirtualMachine struct {
	OperatingSystem string `json:"operatingsystem"`
	Size            string `json:"size"`
	Tags            string `json:"tags"`
}

//NetworkInterfaces stores the network interface configuration information
type NetworkInterfaces struct {
	PrivateIPAddress     string            `json:"privateipaddress"`
	NetworkSecurityGroup string            `json:"networksecuritygroup"`
	PublicIPAddress      ResourceReference `json:"publicipaddress"`
}

//ResourceReference stores ID and Value for child elements
type ResourceReference struct {
	ID    string
	Value string
}

//PublicIPAddress stores the public IPAddress configuration information
type PublicIPAddress struct {
	PublicIPAddressVersion   string            `json:"publicipaddressversion"`
	PublicIPAllocationMethod string            `json:"publicipallocationmethod"`
	IPConfiguration          ResourceReference `json:"ipconfiguration"`
}

//ResourceInfo represents all mandatory fields for resource
type ResourceInfo struct {
	VendorID          string `json:"vendorid"`
	ResourceID        string `json:"resourceid"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Location          string `json:"location"`
	ResourceGroupName string `json:"resourcegroupname"`
}

//ResourceChildDependent stores the common configuration information
// type ResourceDependent struct{
// 	ManagedDisk ManagedDisk `json:"manageddisk,omitempty"`
// 	NetworkProfile []NetworkInterfaces `json:"networkProfile,omitempty"`
// 	PublicIPAddress PublicIPAddress `json:"publicIPAddress,omitempty"`
// 	Subnet Subnet `json:"subnet,omitempty"`
// 	NetworkSecurityGroup NetworkSecurityGroup `json:"networkSecurityGroup,omitempty"`
// }

// //ManagedDisk ...
// type ManagedDisk struct{
// 	StorageAccountType string `json:"storageAccountType,omitempty"`
// 	ID	string `json:"id,omitempty"`
// }

// //NetworkInterfaces ...
// type NetworkInterfaces struct{
// 	ID []string `json:"id,omitempty"`
// }

// //PublicIPAddress ...
// type PublicIPAddress struct{
// 	ID	string `json:"id,omitempty"`
// }

// //Subnet ...
// type Subnet struct{
// 	ID	string `json:"id,omitempty"`
// }

// //NetworkSecurityGroup ...
// type NetworkSecurityGroup struct{
// 	ID	string `json:"id,omitempty"`
// }
