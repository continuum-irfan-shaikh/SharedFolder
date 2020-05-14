package asset

// USBDevice defines the structure of a USB devices as part of new asset model
type USBDevice struct {
	DeviceName   string `json:"deviceName,omitempty" cql:"dev_name"`
	ProductID    string `json:"productID,omitempty" cql:"product_id"`
	VendorID     string `json:"vendorID,omitempty" cql:"vendor_id"`
	Version      string `json:"version,omitempty" cql:"version"`
	SerialNum    string `json:"serialnum,omitempty" cql:"serial_num"`
	Speed        string `json:"speed,omitempty" cql:"speed"`
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	LocID        string `json:"locID,omitempty" cql:"loc_id"`
	CurrAvl      string `json:"currAvl,omitempty" cql:"curr_avl"`
	CurrReq      string `json:"currReq,omitempty" cql:"curr_req"`
}

// USBBus defines the structure of a USB Bus as part of new asset model
type USBBus struct {
	Board           string      `json:"Board,omitempty" cql:"board"`
	HstCtrlLocation string      `json:"hstCtrlLocation,omitempty" cql:"hst_ctrl_location"`
	HstCtrlDrv      string      `json:"hstCtrlDrv,omitempty" cql:"hst_ctrl_drv"`
	PciDeviceID     string      `json:"pciDeviceID,omitempty" cql:"pci_device_id"`
	PciRevisionID   string      `json:"pciRevisionID,omitempty" cql:"pci_revision_id"`
	PciVendorID     string      `json:"pciVendorID,omitempty" cql:"pci_vendor_id"`
	BusNum          string      `json:"busNum,omitempty" cql:"bus_num"`
	Devices         []USBDevice `json:"devices,omitempty" cql:"devices"`
}
