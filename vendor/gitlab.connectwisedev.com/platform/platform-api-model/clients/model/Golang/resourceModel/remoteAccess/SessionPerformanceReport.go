package remoteAccess

import "time"

//SessionPerformanceReportReqBody represents http POST request body structure
type SessionPerformanceReportReqBody struct {
	Vendor        string   `json:"vendor"`
	EndpointsList []string `json:"endpointsList"`
	SessionID     string   `json:"sessionId"`
}

//SessionPerformanceReportNode represents Details for Session Performance Report Record
type SessionPerformanceReportNode struct {
	PartnerID                   string    `json:"partnerId"`
	EndpointID                  string    `json:"endpointId"`
	SessionID                   string    `json:"sessionId"`
	ATotalProcessing            float64   `json:"a_total_processing"`
	B1DatabaseInitialization    float64   `json:"b_1_database_initialization"`
	B2PreProcessing             float64   `json:"b_2_pre_processing"`
	B3PostProcessing            float64   `json:"b_3_post_processing"`
	B4AudidProcessing           float64   `json:"b_4_audid_processing"`
	C11DatabaseConnection       float64   `json:"c_1_1_database_connection"`
	C21UserAccessCheck          float64   `json:"c_2_1_user_access_check"`
	C22ThirdPartyParamFetch     float64   `json:"c_2_2_third_party_param_fetch"`
	C23ConnectionReqBody        float64   `json:"c_2_3_connection_req_body"`
	C24ParamFetch               float64   `json:"c_2_4_param_fetch"`
	C31RemoteConnectionRequest  float64   `json:"c_3_1_remote_connection_request"`
	C32OneClickRequest          float64   `json:"c_3_2_one_click_request"`
	C41SessionAuditEntry        float64   `json:"c_4_1_session_audit_entry"`
	C42StatusAuditEntry         float64   `json:"c_4_2_status_audit_entry"`
	C43EnableManifestAuditEntry float64   `json:"c_4_3_enable_manifest_audit_entry"`
	D241DatabaseFetch           float64   `json:"d_2_4_1_database_fetch"`
	D242FallbackAssetFetch      float64   `json:"d_2_4_2_fallback_asset_fetch"`
	D243FallbackItsapiFetch     float64   `json:"d_2_4_3_fallback_itsapi_fetch"`
	D321OneClickEntitlement     float64   `json:"d_3_2_1_one_click_entitlement"`
	D322OneClickVault           float64   `json:"d_3_2_2_one_click_vault"`
	Dcdtime                     time.Time `json:"dcdtime"`
}

//SessionPerformanceReportData represents Details for Session Performance Report
type SessionPerformanceReportData struct {
	OutData []SessionPerformanceReportNode `json:"outData"`
}
