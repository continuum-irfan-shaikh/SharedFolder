package models

import "time"

type Schedule struct {
	StartRunTime time.Time `json:"startDate"         valid:"requiredForRecurrendAndOneTime" `
	EndRunTime   time.Time `json:"endDate"           valid:"optionalOnlyForRecurrent"`
}

type TimeData struct {
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"ModifiedAt"`
}

type Data struct {
	PartnerID    string   `json:"startDate"`
	MEID         string   `json:"MEID"`
	TaskName     string   `json:"TaskName"`
	CreatedAt    string   `json:"createdAt"`
	CreatedBy    string   `json:"createdBy"`
	StartRunTime string   `json:"startDate"`
	EndRunTime   string   `json:"endDate"`
	ModifiedAt   string   `json:"ModifiedAt"`
	ModifiedBy   string   `json:"ModifiedBy"`
	State        string   `json:"State"`
	UserSites    []string `json:"UserSites"`
	MachineName  string   `json:"MachineName"`
}

type Sites struct {
	SiteList []struct {
		ID int64 `json:"siteId"`
	} `json:"siteDetailList"`
}

type Asset struct {
	EndpointID   string `json:"endpointID"`
	FriendlyName string `json:"friendlyName"`
}

func StructToMap(asset []Asset) (groupedNames map[string]string) {
	groupedNames = make(map[string]string)
	for _, val := range asset {
		groupedNames[val.EndpointID] = val.FriendlyName
	}

	return groupedNames
}

// GetTaskStateText returns text representation of task state
func GetTaskStateText(taskState string) string {
	switch taskState {
	case "1":
		return "Active"
	case "2":
		return "Inactive"
	case "3":
		return "Disabled"
	default:
		return "unknownStatus"
	}
}
