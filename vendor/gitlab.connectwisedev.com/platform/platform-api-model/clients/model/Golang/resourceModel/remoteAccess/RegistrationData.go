package remoteAccess

//RegistrationData represents Details for Registered Endpoint in Site
type RegistrationData struct {
	EndpointID string `json:"endpointId"`
	StatusCode int    `json:"statusCode"`
	Details    string `json:"details"`
}
