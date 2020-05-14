package endpointstate

const (
	// DeviceUp : Device Up Notification
	DeviceUp = "DeviceUp"
	// DeviceDown : Device dwon Notification
	DeviceDown = "DeviceDown"
)

//NotificationMessage : Endpoint State Notification Message
type NotificationMessage struct {
	NotificationType string `json:"notificationType"`
}
