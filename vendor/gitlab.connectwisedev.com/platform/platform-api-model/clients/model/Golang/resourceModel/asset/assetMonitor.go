package asset

//Monitor is the struct definition of /resources/asset/assetMonitor
type Monitor struct {
	Name         string `json:"name,omitempty" cql:"name"`
	DeviceID     string `json:"deviceID,omitempty" cql:"device_id"`
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	ScreenHeight uint32 `json:"screenHeight,omitempty" cql:"screen_height"`
	ScreenWidth  uint32 `json:"screenWidth,omitempty" cql:"screen_width"`
	Resolution   string `json:"resolution,omitempty" cql:"resolution"`
	PixelDepth   string `json:"pixelDepth,omitempty" cql:"pixel_depth"`
	MainDisplay  string `json:"mainDisplay,omitempty" cql:"main_display"`
	Mirror       string `json:"mirror,omitempty" cql:"mirror"`
	Online       string `json:"online,omitempty" cql:"online"`
	Builtin      string `json:"builtin,omitempty" cql:"builtin"`
	ConnType     string `json:"connType,omitempty" cql:"conn_type"`
	Description  string `json:"description,omitempty" cql:"description"`
	Version      string `json:"version,omitempty" cql:"version"`
}
