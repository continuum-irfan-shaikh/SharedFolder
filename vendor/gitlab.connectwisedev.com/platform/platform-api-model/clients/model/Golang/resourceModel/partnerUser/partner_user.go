package partnerUser

import "time"

const (
	// ActionUpdate is record update action
	ActionUpdate = "UPDATE"

	// ActionNew is record create action
	ActionNew = "NEW"

	// ActionDelete is record delete action
	ActionDelete = "DELETE"
)

// PartnerUserChange is partner user change event message
type PartnerUserChange struct {
	// Action is action name
	//
	// Possible values: NEW, DELETE or UPDATE
	Action string `json:"action"`

	// Data contains updated partner user data
	Data PartnerUser `json:"data,omitempty"`
}

// PartnerUser is an API model to hold partner user data
type PartnerUser struct {
	UserID        int    `json:"userId,omitempty"`
	Email         string `json:"email,omitempty"`
	Designation   string `json:"designation,omitempty"`
	FullName      string `json:"fullName,omitempty"`
	City          string `json:"city,omitempty"`
	Address       string `json:"address,omitempty"`
	PartnerID     int    `json:"partnerId,omitempty"`
	IsEnabled     bool   `json:"isEnabled,omitempty"`
	IsNocUser     bool   `json:"isNocUser,omitempty"`
	SendSMS       bool   `json:"sendSms,omitempty"`
	MobilePhone   string `json:"mobilePhone,omitempty"`
	OfficePhone   string `json:"officePhone,omitempty"`
	Extension     string `json:"extension,omitempty"`
	PrimExtension string `json:"primExtension,omitempty"`
	ZipCode       string `json:"zipCode,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"` // CreatedAt is creation date is in RFC3339 format
}

// CreationTime returns parsed creation time from CreatedAt property
func (u *PartnerUser) CreationTime() (time.Time, error) {
	return time.Parse(time.RFC3339, u.CreatedAt)
}
