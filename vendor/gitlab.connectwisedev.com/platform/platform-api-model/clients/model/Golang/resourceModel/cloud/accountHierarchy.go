package cloud

import "time"

//AccountHierarchySummary represents summary of accounts
type AccountHierarchySummary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	TenantID    string    `json:"tenantid"`
	UserName    string    `json:"username"`
	IsMonitored bool      `json:"ismonitored"`
	State       string    `json:"state"`
	MappedOn    time.Time `json:"mappedon"`
	MappedBy    string    `json:"mappedby"`
}

//MappingStatus represent status of mapping
type MappingStatus struct {
	IsMapped bool `json:"ismapped"`
}
