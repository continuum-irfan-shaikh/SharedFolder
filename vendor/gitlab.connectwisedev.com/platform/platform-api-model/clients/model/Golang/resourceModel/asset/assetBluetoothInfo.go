package asset

// BluetoothHardware is the definition of /asset/bluetoothHardware
type BluetoothHardware struct {
	Version      string `json:"version,omitempty" cql:"version"`
	Address      string `json:"address,omitempty" cql:"address"`
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	Name         string `json:"name,omitempty" cql:"name"`
	FirmwareVer  string `json:"firmwareVersion,omitempty" cql:"firmware_version"`
	Power        string `json:"power,omitempty" cql:"power"`
	Discoverable string `json:"discoverable,omitempty" cql:"discoverable"`
	VendorID     string `json:"vendorID,omitempty" cql:"vendor_id"`
	ProductID    string `json:"productID,omitempty" cql:"product_id"`
	HCIVer       string `json:"hciVer,omitempty" cql:"hci_ver"`
	HCIRev       string `json:"hciRev,omitempty" cql:"hci_rev"`
	LMPVer       string `json:"lmpVer,omitempty" cql:"lmp_ver"`
	LMPSubVer    string `json:"lmpSubVer,omitempty" cql:"lmp_sub_ver"`
	CmpltDevType string `json:"cmpltDevType,omitempty" cql:"cmplt_dev_type"`
	CompstClass  string `json:"compstClass,omitempty" cql:"compst_class"`
	ServCls      string `json:"serviceClass,omitempty" cql:"service_class"`
	MJRDevType   string `json:"majorDeviceType,omitempty" cql:"major_device_type"`
	MjrDevCls    string `json:"majorDeviceClass,omitempty" cql:"major_device_class"`
	MinrDevCls   string `json:"minorDeviceClass,omitempty" cql:"minor_device_class"`
}

// BluetoothDevice is the definition of /asset/bluetoothDevice
type BluetoothDevice struct {
	DeviceName      string `json:"deviceName,omitempty" cql:"device_name"`
	Address         string `json:"address,omitempty" cql:"address"`
	Type            string `json:"type,omitempty" cql:"type"`
	FirmwareVersion string `json:"firmwareVersion,omitempty" cql:"firmware_version"`
	Services        string `json:"services,omitempty" cql:"services"`
	Paired          string `json:"paired,omitempty" cql:"paired"`
	Favorite        string `json:"favorite,omitempty" cql:"favorite"`
	Connected       string `json:"connected,omitempty" cql:"connected"`
	Manufacturer    string `json:"manufacturer,omitempty" cql:"manufacturer"`
	VendorID        string `json:"vendorID,omitempty" cql:"vendor_id"`
	ProductID       string `json:"productID,omitempty" cql:"product_id"`
	MajorType       string `json:"majorType,omitempty" cql:"major_type"`
	DeviceClass     string `json:"deviceClass,omitempty" cql:"device_class"`
	MinorType       string `json:"minorType,omitempty" cql:"minor_type"`
	EDRSupported    string `json:"edrsSupported,omitempty" cql:"edrs_supported"`
	ESCOsupported   string `json:"escoSupported,omitempty" cql:"esco_supported"`
}

// BluetoothService is the definition of /asset/bluetoothService
type BluetoothService struct {
	ServiceName           string `json:"serviceName,omitempty" cql:"service_name"`
	FldaccItem            string `json:"fldaccItem,omitempty" cql:"fldacc_item"`
	FldFileBrws           string `json:"fldFileBrws,omitempty" cql:"fld_file_brws"`
	AuthRequired          string `json:"authRequired,omitempty" cql:"auth_required"`
	State                 string `json:"state,omitempty" cql:"state"`
	StateOthrItemsAccept  string `json:"stateOthrItemsAccept,omitempty" cql:"state_othr_items_accept"`
	OthrPimItemsAccept    string `json:"othrPimItemsAccept,omitempty" cql:"othr_pim_items_accept"`
	ReceiveItemAction     string `json:"receiveItemAction,omitempty" cql:"receive_item_action"`
}

// AssetBluetooth is the struct definition of /asset/bluetoothInfo
type AssetBluetooth struct {
	Hardware BluetoothHardware  `json:"hardware,omitempty" cql:"hardware"`
	Services []BluetoothService `json:"services,omitempty" cql:"services"`
	Devices  []BluetoothDevice  `json:"devices,omitempty" cql:"devices"`
}
