package asset

//AssetRaidController is the struct definition of /resources/asset/assetRaidController
type AssetRaidController struct {
	SoftwareRaid string `json:"softwareRaid,omitempty" cql:"software_raid"`
	HardwareRaid string `json:"hardwareRaid,omitempty" cql:"hardware_raid"`
	Vendor       string `json:"vendor,omitempty" cql:"vendor"`
	*HardwareRaidController
}

type HardwareRaidController struct {
	HardwareRaidInfo []HardwareRaidInfo `json:"hardwareRaidInfo,omitempty" cql:"hardware_raid"`
	TotalSpace       int64              `json:"totalSpace,omitempty" cql:"total_space"`
	UsedSpace        int64              `json:"usedSpace,omitempty" cql:"used_space"`
}

type HardwareRaidInfo struct {
	Type         string `json:"type,omitempty" cql:"type"`
	SerialNumber string `json:"serialNumber,omitempty" cql:"serial_number"`
	Capacity     string `json:"capacity,omitempty" cql:"capacity"`
}
