package remoteAccess

//ErrorResponse represents message format for Error Response
type ErrorResponse struct {
	Message string `json:"message"`
	HostID  int64  `json:"hostId"`
}
