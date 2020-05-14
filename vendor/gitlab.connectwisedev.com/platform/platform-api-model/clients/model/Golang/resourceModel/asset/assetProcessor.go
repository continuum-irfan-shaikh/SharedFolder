package asset

//AssetProcessor is the struct definition of /resources/asset/assetProcessor
type AssetProcessor struct {
	Product                   string  `json:"product,omitempty" cql:"product"`
	NumberOfProcessors        int     `json:"numberOfProcessors" cql:"number_of_processors"`
	NumberOfCores             int     `json:"numberOfCores" cql:"number_of_cores"`
	ClockSpeedMhz             float64 `json:"clockSpeedMhz" cql:"clock_speed_mhz"`
	Family                    int     `json:"family" cql:"family"`
	Manufacturer              string  `json:"manufacturer,omitempty" cql:"manufacturer"`
	ProcessorType             string  `json:"processorType,omitempty" cql:"processor_type"`
	SerialNumber              string  `json:"serialNumber,omitempty" cql:"serial_number"`
	Level                     int     `json:"level" cql:"level"`
	L2Cache                   string  `json:"l2Cache,omitempty" cql:"l2_cache"`
	L3Cache                   string  `json:"l3Cache,omitempty" cql:"l3_cache"`
	BootRom                   string  `json:"bootRom,omitempty" cql:"boot_rom"`
	NumberOfLogicalProcessors int     `json:"numberOfLogicalProcessors" cql:"number_of_logical_processors"`
	Version                   string  `json:"version" cql:"version"`
}
