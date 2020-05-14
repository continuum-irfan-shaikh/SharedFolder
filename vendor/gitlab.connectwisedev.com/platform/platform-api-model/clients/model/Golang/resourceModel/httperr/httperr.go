package httperr

//http error codes
const (
	//ErrHTTPReadingRequestBody error for failed to read request body, equivalent to http 400
	ErrHTTPReadingRequestBody = "ErrReadingRequestBody"

	//ErrHTTPInvalidInput error for user input invalid, equivalent to http 400
	ErrHTTPInvalidInput = "ErrInvalidInput"

	//ErrHTTFailedToProcess error for failed to process, equivalent to http 500
	ErrHTTPFailedToProcess = "ErrFailedToProcess"
)

//Response is http/web response body for error
type Response struct {
	ResponseErr ResponseErr `json:"error"`
}

//ResponseErr is the http/web response inner body to have more description about error
type ResponseErr struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
