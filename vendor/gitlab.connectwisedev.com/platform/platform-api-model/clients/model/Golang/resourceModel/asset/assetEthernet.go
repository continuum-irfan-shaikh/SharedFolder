package asset

//AssetEthernet is the struct definition of /resources/asset/assetEthernet
type AssetEthernet struct {
	BroadCom        string `json:"broadCom ,omitempty" cql:"broad_Com"`
	Name            string `json:"name,omitempty" cql:"name"`
	Type            string `json:"type,omitempty" cql:"type"`
	Bus             string `json:"bus,omitempty" cql:"bus"`
	EthrnetVendorID string `json:"ethrnetVendorID,omitempty" cql:"ethrnet_vendor_id"`
	DeviceID        string `json:"deviceID,omitempty" cql:"device_id"`
	SubSysVendorID  string `json:"subSysVendorid,omitempty" cql:"sub_sys_vendor_id"`
	SubSysID        string `json:"subSysID,omitempty" cql:"sub_sys_id"`
	RevisionID      string `json:"revisionID,omitempty" cql:"revision_id"`
	LinkWidth       string `json:"linkWidth,omitempty" cql:"link_width"`
	BsdName         string `json:"bsdName,omitempty" cql:"bsd_name"`
	KextName        string `json:"kextName,omitempty" cql:"kext_name"`
	FirmwareVersion string `json:"firmwareVersion,omitempty" cql:"firmware_version"`
	Location        string `json:"location,omitempty" cql:"location"`
	Version         string `json:"version,omitempty" cql:"version"`
}
