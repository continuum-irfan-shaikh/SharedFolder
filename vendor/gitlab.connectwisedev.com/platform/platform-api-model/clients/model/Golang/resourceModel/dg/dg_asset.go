package dg

//DGAsset is the definition of /resources/dynamicGroups/example/managed_endpoint_message_example_asset.json
type DGAsset struct {
	OriginDomain string `json:"originDomain,omitempty"`
	ID           string `json:"id,omitempty"`
	Client       string `json:"client,omitempty"`
	Partner      string `json:"partner,omitempty"`
	Site         string `json:"site,omitempty"`
	AssetData    Asset  `json:"asset,omitempty"`
	LegacyRegID  string `json:"legacy_regid,omitempty"`
}

//Asset is the object to represent Managed Endpoint Change(MEC) message
type Asset struct {
	OS                    string                     `json:"os,omitempty"`
	ServicePack           string                     `json:"service_pack,omitempty"`
	OSVersion             string                     `json:"os_version,omitempty"`
	BaseboardManufacturer string                     `json:"baseboard_manufacturer,omitempty"`
	VirtualMachine        string                     `json:"virtual_machine,omitempty"`
	InstalledSoftware     []InstalledSoftwareMessage `json:"installed_software,omitempty"`
	RAM                   uint64                     `json:"ram,omitempty"`
	//These new fields are added to standardize the MEC message Context and the changes are triggered by change request in RMM-43590
	FriendlyName string   `json:"machine_friendly_name,omitempty"`
	EndpointType string   `json:"endpoint_type,omitempty"`
	IPv4List     []string `json:"ipv4_list,omitempty"`
	SystemName   string   `json:"machine_name,omitempty"`
}
