package asset

//DrivePartition is the struct definition of /resources/asset/assetDrivePartition
type DrivePartition struct {
	Name           string `json:"name,omitempty" cql:"name"`
	Label          string `json:"label,omitempty" cql:"label"`
	FileSystem     string `json:"fileSystem,omitempty" cql:"file_system"`
	Description    string `json:"description,omitempty" cql:"description"`
	SizeBytes      int64  `json:"sizeBytes" cql:"size_bytes"`
	Writable       string `json:"writable" cql:"writable"`
	MountPoint     string `json:"mountPoint" cql:"mount_point"`
	DFSize         int64  `json:"dfsize" cql:"dfsize"`
	FreeSpaceBytes int64  `json:"freeSpaceBytes" cql:"freespace_bytes"`
	UsedSpaceBytes int64  `json:"usedSpaceBytes" cql:"usedspace_bytes"`
	Version        string `json:"version" cql:"version"`
	Vendor         string `json:"vendor" cql:"vendor"`
}
