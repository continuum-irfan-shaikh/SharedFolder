package asset

//AssetFirewall is the struct definition of /resources/asset/assetFireWall
type AssetFirewall struct {
	Mode        string `json:"mode,omitempty" cql:"mode"`
	Logging     string `json:"logging,omitempty" cql:"logging"`
	StealthMode string `json:"stealthMode,omitempty" cql:"stealth_mode"`
}
