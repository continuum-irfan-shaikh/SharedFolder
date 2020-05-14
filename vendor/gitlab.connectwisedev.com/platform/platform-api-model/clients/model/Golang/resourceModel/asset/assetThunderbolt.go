package asset

// AssetThunderbolt is the struct definition of /asset/assetThunderbolt
type AssetThunderbolt struct {
	VendorName          string `json:"vendorName,omitempty" cql:"vendor_name"`
	DeviceName          string `json:"deviceName,omitempty" cql:"device_name"`
	UID                 string `json:"UID,omitempty" cql:"uid"`
	FirmwareVersion     string `json:"firmwareVersion,omitempty" cql:"firmware_version"`
	PortStatus          string `json:"portStatus,omitempty" cql:"port_status"`
	PortLink            string `json:"portLink,omitempty" cql:"port_link"`
	PortFirmwareVersion string `json:"portFirmwareVersion,omitempty" cql:"port_firmware_version"`
}
