package IntegrationStatusSchema

// IntegrationStatusData integration status data type
type IntegrationStatusData struct {
	EndPointID string              `json:"endpointId"`
	Statuses   []IntegrationStatus `json:"statuses"`
}

// IntegrationStatus integration status type
type IntegrationStatus struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	ProductID string `json:"productId"`
}
