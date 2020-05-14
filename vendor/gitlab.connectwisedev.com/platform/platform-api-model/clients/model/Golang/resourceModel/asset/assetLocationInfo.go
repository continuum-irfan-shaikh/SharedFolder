package asset

//AssetLocation is the struct definition of /resources/asset/assetLocation
type AssetLocation struct {
	LocType          string                  `json:"locType,omitempty" cql:"loc_type"`
	ActiveLoc        string                  `json:"activeLoc,omitempty" cql:"active_loc"`
	LocationServices []AssetLocationServices `json:"locationServices,omitempty" cql:"location_services"`
}
