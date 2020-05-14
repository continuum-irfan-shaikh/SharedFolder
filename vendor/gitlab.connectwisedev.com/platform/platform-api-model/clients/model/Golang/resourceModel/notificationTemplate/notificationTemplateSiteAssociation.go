package notificationTemplate

// NotificationTemplateSiteAssociationMessage is model of Notification Template Association to Site message
type NotificationTemplateSiteAssociationMessage struct {
	MessageType            int    `json:"messageType"`
	Action                 string `json:"action"`
	NotificationTemplateID string `json:"notificationTemplateId"`
	PartnerID              int    `json:"partnerId"`
	ClientID               int    `json:"clientId"`
	SiteID                 int    `json:"siteId"`
}
