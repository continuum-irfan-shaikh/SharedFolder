package asset

//AssetLogicalVolume is the struct definition of /resources/asset/assetLogicalVolume
type AssetLogicalVolume struct {
	Name           string `json:"name,omitempty" cql:"name"`
	Description    string `json:"description,omitempty" cql:"description"`
	Capacity       int64  `json:"capacity,omitempty" cql:"capacity"`
	DFSize         int64  `json:"dfsize,omitempty" cql:"dfsize"`
	FreeSpaceBytes int64  `json:"freeSpaceBytes,omitempty" cql:"freespace_bytes"`
	UsedSpaceBytes int64  `json:"usedSpaceBytes,omitempty" cql:"usedspace_bytes"`
}
