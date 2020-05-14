package itsm_integration

// ServiceCatalog is service catalog
type ServiceCatalog struct {
	// Title is title
	Title string `json:"title"`

	// Description is description
	Description string `json:"description"`

	// SCTCategory is template category
	SCTCategory string `json:"sctCategory"`

	// Service is service
	Service string `json:"service"`

	// Category is category
	Category string `json:"category"`

	// SubCategory is subCategory
	SubCategory string `json:"subCategory"`

	// Status is status
	Status string `json:"status"`

	// MaxOrder is max order
	MaxOrder int `json:"maxOrder"`

	// Company name is company name
	CompanyName string `json:"companyName"`

	// CreatedDateTime is creation date in RFC3339 format
	CreatedDateTime string `json:"createdAt"`

	// CreatedBy is creator name
	CreatedBy string `json:"createdBy"`

	// TemplateID is template ID
	TemplateID string `json:"templateId"`

	// PortalTitle is portal title
	PortalTitle string `json:"portalTitle"`

	// PortalDescription is description
	PortalDescription string `json:"portalDescription"`

	// BusinessOwnerID is business owner's contact id
	BusinessOwnerID string `json:"businessOwnerId"`

	// BusinessOwnerFullName is business owner's full name
	BusinessOwnerFullName string `json:"businessOwnerFullName"`

	// BusinessOwnerEmail is business owner's email
	BusinessOwnerEmail string `json:"businessOwnerEmail"`

	// BusinessOwnerPhone is business owner's phone
	BusinessOwnerPhone string `json:"businessOwnerPhone"`
}
