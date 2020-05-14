package asset

import "time"

//InstalledSoftware is the struct definition of /resources/asset/assetInstalledSoftware
type InstalledSoftware struct {
	Name               string    `json:"name,omitempty" cql:"name"`
	Publisher          string    `json:"publisher,omitempty" cql:"publisher"`
	Version            string    `json:"version,omitempty" cql:"version"`
	InstallDate        time.Time `json:"installDate,omitempty" cql:"install_date"`
	UserName           string    `json:"userName,omitempty" cql:"user_name"`
	LastAccessDateTime time.Time `json:"lastAccessDateTime,omitempty" cql:"last_access_datetime"`
	Kind               string    `json:"kind,omitempty" cql:"kind"`
	Intel64Bit         string    `json:"intel64Bit,omitempty" cql:"intel64bit"`
	GetInfoString      string    `json:"getInfoString,omitempty" cql:"info_string"`
	AppStore           string    `json:"appStore,omitempty" cql:"appstore"`
	Location           string    `json:"location,omitempty" cql:"location"`
	InstalledSizeBytes string    `json:"installedSizeBytes,omitempty" cql:"installedsize_bytes"`
}
