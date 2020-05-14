package asset

//AssetLocationServices is the struct definition of /resources/asset/assetLocationServices
type AssetLocationServices struct {
	LocDevName       string `json:"locDevName,omitempty" cql:"loc_dev_name"`
	LocBsdName       string `json:"locBsdName,omitempty" cql:"loc_bsd_name"`
	HrdMacAddr       string `json:"hrdMacAddr,omitempty" cql:"hrd_mac_addr"`
	DevType          string `json:"devType,omitempty" cql:"dev_type"`
	IPv4CnfMethod    string `json:"ipv4CnfMethod,omitempty" cql:"ipv4_cnf_method"`
	IPv6CnfMethod    string `json:"ipv6CnfMethod,omitempty" cql:"ipv6_cnf_method"`
	ProxyExcPlist    string `json:"proxyExcPlist,omitempty" cql:"proxy_exc_plist"`
	ProxyFTPMode     string `json:"proxyFTPMode,omitempty" cql:"proxy_ftp_mode"`
	PPPASCPEnabled   string `json:"pppASCPEnabled,omitempty" cql:"ppp_ascp_enabled"`
	RedialCount      int    `json:"redialCount,omitempty" cql:"redial_count"`
	RedialEnabled    string `json:"redialEnabled,omitempty" cql:"redial_enabled"`
	RedialInterval   int    `json:"redialInterval,omitempty" cql:"redial_interval"`
	UseTermScript    string `json:"useTermScript,omitempty" cql:"use_term_script"`
	DialOnDemand     string `json:"dialOnDemand,omitempty" cql:"dial_on_demand"`
	DisctFastSwitch  string `json:"disctFastSwitch,omitempty" cql:"disct_fast_switch"`
	IdleDistconnect  string `json:"idleDistconnect,omitempty" cql:"idle_distconnect"`
	IdleDisctTimer   int    `json:"idleDisctTimer,omitempty" cql:"idle_disct_timer"`
	LogOutDisct      string `json:"logOutDisct,omitempty" cql:"log_out_disct"`
	SleepDisct       string `json:"sleepDisct,omitempty" cql:"sleep_disct"`
	IdleReminder     string `json:"idleReminder,omitempty" cql:"idle_reminder"`
	IdleRemTime      int    `json:"idleRemTime,omitempty" cql:"idle_rem_time"`
	IpcpCmprVJ       string `json:"ipcpCmprVJ,omitempty" cql:"ipcp_cmpr_vj"`
	LcpEchoEnbld     string `json:"lcpEchoEnbld,omitempty" cql:"lcp_echo_enbld"`
	LcpEchoFailure   int    `json:"lcpEchoFailure,omitempty" cql:"lcp_echo_failure"`
	LcpEchoInterval  int    `json:"lcpEchoInterval,omitempty" cql:"lcp_echo_interval"`
	LogFile          string `json:"logFile,omitempty" cql:"log_file"`
	VerboseLog       string `json:"verboseLog,omitempty" cql:"verbose_log"`
	IEEEName         string `json:"ieeeName,omitempty" cql:"ieee_name"`
	IEEEJoinMdFallBk string `json:"ieeeJoinMdFallBk,omitempty" cql:"ieee_join_md_fall_bk"`
	IEEEPowerEnbld   bool   `json:"ieeePowerEnbld,omitempty" cql:"ieee_power_enbld"`
	IEEERemJoinNw    bool   `json:"ieeeRemJoinNw,omitempty" cql:"ieee_rem_join_nw"`
}
