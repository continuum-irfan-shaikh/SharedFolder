package agent

//http error codes
const (
	//ErrHTTPReadingRequestBody error for failed to read request body, equivalent to http 400
	ErrHTTPReadingRequestBody = "ErrReadingRequestBody"

	//ErrHTTPInvalidInput error for user input invalid, equivalent to http 400
	ErrHTTPInvalidInput = "ErrInvalidInput"

	//ErrHTTFailedToProcess error for failed to process, equivalent to http 500
	ErrHTTPFailedToProcess = "ErrFailedToProcess"
	
	//ErrEndpointValidation failed to get validate agent or get its mapping
	ErrEndpointValidation = "ErrEndpointValidation"
	
	//ErrAgentAuthentication failed to verify signature with available agent's public key
	ErrAgentAuthentication = "ErrAgentAuthentication"
)

//ErrResponse is http/web response body for error
type ErrResponse struct {
	ResponseErr ResponseErr `json:"error"`
}

//ResponseErr is the http/web response inner body to have more description about error
type ResponseErr struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
