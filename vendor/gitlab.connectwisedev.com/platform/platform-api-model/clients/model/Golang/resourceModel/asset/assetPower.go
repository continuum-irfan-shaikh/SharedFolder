package asset

//AssetPower is the struct definition of /resources/asset/assetPower
type AssetPower struct {
	PowerType      string `json:"powerType,omitempty" cql:"power_type"`
	SysSleepTimer  string `json:"sysSleepTimer,omitempty" cql:"sys_sleep_timer"`
	DiscSleepTimer string `json:"discSleepTimer,omitempty" cql:"disc_sleep_timer"`
	DispSleepTimer string `json:"dispSleepTimer,omitempty" cql:"disp_sleep_timer"`
	SleepOnPower   string `json:"sleepOnPower,omitempty" cql:"sleep_on_power"`
	AutoRestart    string `json:"autoRestart,omitempty" cql:"auto_restart"`
	WakeOnLan      string `json:"wakeOnLan,omitempty" cql:"wake_on_lan"`
	CurPower       string `json:"curPower,omitempty" cql:"cur_power"`
	UsesDim        string `json:"usesDim,omitempty" cql:"uses_dim"`
	NWReachability string `json:"nwReachability,omitempty" cql:"nw_reachability"`
	KernalRestart  string `json:"kernalRestart,omitempty" cql:"kernal_restart"`
	UpsInstalled   string `json:"upsInstalled,omitempty" cql:"ups_installed"`
}
