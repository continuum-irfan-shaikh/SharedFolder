package asset

// AssetDiagnostic is the struct definition of /asset/assetDiagnostic
type AssetDiagnostic struct {
	Name      string `json:"testName,omitempty" cql:"test_name"`
	LastRunOn string `json:"lastRunOn,omitempty" cql:"last_run_on"`
	Result    string `json:"result,omitempty" cql:"result"`
}
