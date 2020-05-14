package itsm_integration

// NotificationTemplate represents notificationTemplate object in Cherwell
type NotificationTemplate struct {
	RecID                  string            `json:"id"`
	CreatedDateTime        string            `json:"createdDateTime"`
	CreatedBy              string            `json:"createdBy"`
	LastModDateTime        string            `json:"lastModDateTime"`
	CreatedByID            string            `json:"createdById"`
	CreatedCulture         string            `json:"createdCulture"`
	LastModBy              string            `json:"lastModBy"`
	LastModByID            string            `json:"lastModById"`
	OwnedBy                string            `json:"ownedBy"`
	OwnedByID              string            `json:"ownedById"`
	OwnedByTeam            string            `json:"ownedByTeam"`
	OwnedByTeamID          string            `json:"ownedByTeamId"`
	NotificationTemplateID string            `json:"notificationTemplateId"`
	Data                   string            `json:"data"`
	SiteID                 string            `json:"siteId"`
	PartnerID              string            `json:"partnerId"`
	TemplateKey            string            `json:"templateKey"`
	TemplateID             string            `json:"templateId"`
	TemplateType           string            `json:"templateType"`
	UserKey                string            `json:"userKey"`
	ServiceID              string            `json:"serviceId"`
	ClientID               string            `json:"clientId"`
	TemplateName           string            `json:"templateName"`
	PartnerRecID           string            `json:"partnerRecId"`
	SiteTimezone           string            `json:"siteTimezone"`
	ScheduledActions       []ScheduledAction `json:"scheduledActions"`
}
