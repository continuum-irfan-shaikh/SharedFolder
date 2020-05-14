package notificationTemplate

// NotificationTemplateSyncMessage is notification template message
type NotificationTemplateSyncMessage struct {
	MessageType                 int                          `json:"messageType"`
	Action                      string                       `json:"action"`
	TimeSlotDay                 string                       `json:"timeSlotDay"`
	NotificationTemplateID      string                       `json:"notificationTemplateId"`
	NotificationTemplateName    string                       `json:"notificationTemplateName"`
	PartnerID                   int                          `json:"partnerId"`
	ServiceID                   string                       `json:"serviceId"`
	UserKey                     string                       `json:"userKey"`
	NotificationTemplateDetails []NotificationTemplateDetail `json:"notificationTemplateDetails"`
}

// NotificationTemplateDetail is model to hold Notification template details
type NotificationTemplateDetail struct {
	MissingInformationIsConfigured string `json:"missingInformationIsConfigured,omitempty"`
	MissingInformationPhoneContact string `json:"missingInformationPhoneContact,omitempty"`
	TemplateType                   string `json:"templateType,omitempty"`
	ContactID                      string `json:"contactId1,omitempty"`
	ContactID2                     string `json:"contactId2,omitempty"`
	ContactID3                     string `json:"contactId3,omitempty"`
	ContactID4                     string `json:"contactId4,omitempty"`
	ContactID5                     string `json:"contactId5,omitempty"`
	ContactID6                     string `json:"contactId6,omitempty"`
	ContactID7                     string `json:"contactId7,omitempty"`
	ContactID8                     string `json:"contactId8,omitempty"`
	ContactID9                     string `json:"contactId9,omitempty"`
	ContactID10                    string `json:"contactId10,omitempty"`
	ContactID11                    string `json:"contactId11,omitempty"`
	ContactID12                    string `json:"contactId12,omitempty"`
	ContactID13                    string `json:"contactId13,omitempty"`
	ContactID14                    string `json:"contactId14,omitempty"`
	ContactID15                    string `json:"contactId15,omitempty"`
	ContactID16                    string `json:"contactId16,omitempty"`
	ContactID17                    string `json:"contactId17,omitempty"`
	ContactID18                    string `json:"contactId18,omitempty"`
	ContactID19                    string `json:"contactId19,omitempty"`
	ContactID20                    string `json:"contactId20,omitempty"`
	TimeSlotStartTime              string `json:"timeSlotStartTime,omitempty"`
	TimeSlotEndTime                string `json:"timeSlotEndTime,omitempty"`
	TimeSlotZone                   string `json:"timeSlotZone,omitempty"`
	ActionID                       string `json:"actionId,omitempty"`
	CategoryID                     string `json:"categoryId,omitempty"`
	ActionDescription              string `json:"actionDescription,omitempty"`
	CategoryDescription            string `json:"categoryDescription,omitempty"`
}
