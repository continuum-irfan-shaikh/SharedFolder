package asset

// PhysicalMemory is the struct definition of /resources/asset/assetPhysicalMemory
type PhysicalMemory struct {
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	SerialNumber string `json:"serialNumber,omitempty" cql:"serial_number"`
	SizeBytes    uint64 `json:"sizeBytes" cql:"size_bytes"`
	ECC          string `json:"ecc,omitempty" cql:"ecc"`
	Slot         string `json:"slot,omitempty" cql:"slot"`
	Type         string `json:"type,omitempty" cql:"type"`
	Speed        string `json:"speed,omitempty" cql:"speed"`
	Status       string `json:"status,omitempty" cql:"status"`
	PartNumber   string `json:"partNumber,omitempty" cql:"part_number"`
	TotalWidth   uint64 `json:"totalwidth,omitempty" cql:"totalwidth"`
	DataWidth    uint64 `json:"datawidth,omitempty" cql:"datawidth"`
}
