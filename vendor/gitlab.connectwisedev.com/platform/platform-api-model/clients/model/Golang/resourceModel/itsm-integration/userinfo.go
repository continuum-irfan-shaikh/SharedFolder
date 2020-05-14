package itsm_integration

// UserInfo struct for payload to update Status field in UserInfo business object in Cherwell
type UserInfo struct {
	Status string `json:"status"`
}
