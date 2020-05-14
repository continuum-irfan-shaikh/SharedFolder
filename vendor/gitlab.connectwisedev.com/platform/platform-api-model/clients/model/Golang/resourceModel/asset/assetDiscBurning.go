package asset

// DiscBurning defines the structure of a DiscBurning as part of new asset model
type DiscBurning struct {
	ModelName        string `json:"modelName,omitempty" cql:"model_name"`
	FirmwareRevision string `json:"firmwareRevision,omitempty" cql:"firmware_revision"`
	Interconnect     string `json:"interconnect,omitempty" cql:"interconnect"`
	BurSupport       string `json:"burSupport,omitempty" cql:"bur_support"`
	Cache            string `json:"cache,omitempty" cql:"cache"`
	ReadDVD          string `json:"readDVD,omitempty" cql:"read_dvd"`
	CDWrite          string `json:"cdWrite,omitempty" cql:"cd_write"`
	DVDWrite         string `json:"dvdWrite,omitempty" cql:"dvd_write"`
	WriteStrtgy      string `json:"writeStrtgy,omitempty" cql:"write_strtgy"`
	Media            string `json:"media,omitempty" cql:"media"`
}
