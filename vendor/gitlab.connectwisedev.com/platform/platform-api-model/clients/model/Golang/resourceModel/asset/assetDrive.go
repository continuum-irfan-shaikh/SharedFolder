package asset

// AssetDrive is the struct definition of /resources/asset/assetDrive
type AssetDrive struct {
	Product            string           `json:"product,omitempty" cql:"product"`
	Manufacturer       string           `json:"manufacturer,omitempty" cql:"manufacturer"`
	MediaType          string           `json:"mediaType,omitempty" cql:"media_type"`
	InterfaceType      string           `json:"interfaceType,omitempty" cql:"interface_type"`
	LogicalName        string           `json:"logicalName,omitempty" cql:"logical_name"`
	SerialNumber       string           `json:"serialNumber,omitempty" cql:"serial_number"`
	Partitions         []string         `json:"partitions,omitempty" cql:"partitions"`
	SizeBytes          int64            `json:"sizeBytes" cql:"size_bytes"`
	NumberOfPartitions int              `json:"numberOfPartitions" cql:"number_of_partitions"`
	PartitionData      []DrivePartition `json:"partitionData,omitempty" cql:"partition_data"`
	LinkSpeed          string           `json:"linkSpeed,omitempty" cql:"link_speed"`
	NLinkSpeed         string           `json:"nLinkSpeed,omitempty" cql:"n_link_speed"`
	Description        string           `json:"description,omitempty" cql:"description"`
	NativeCmdQ         string           `json:"nativeCmdQ,omitempty" cql:"native_cmdq"`
	Model              string           `json:"model,omitempty" cql:"model"`
	DiskRevision       string           `json:"diskRevision,omitempty" cql:"disk_revision"`
	QueueDepth         string           `json:"queueDepth,omitempty" cql:"queue_depth"`
	RemovableMedia     string           `json:"removableMedia,omitempty" cql:"removable_media"`
	DetachableDrive    string           `json:"detachable,omitempty" cql:"detachable"`
	DiskBSDName        string           `json:"diskBSDName,omitempty" cql:"disk_bsd_name"`
	RotationRate       string           `json:"rotationRate,omitempty" cql:"rotation_rate"`
	MediumType         string           `json:"mediumType,omitempty" cql:"medium_type"`
	BayName            string           `json:"bayName,omitempty" cql:"bay_name"`
	PartMapType        string           `json:"partMapType,omitempty" cql:"part_map_type"`
	SmartStatus        string           `json:"smartStatus,omitempty" cql:"smart_status"`
	Version            string           `json:"version,omitempty" cql:"version"`
}
