package approvalTemplate

const (
	// ActionUpdate is record update action
	ActionUpdate = "UPDATE"

	// ActionNew is record create action
	ActionNew = "NEW"

	// ActionDelete is record delete action
	ActionDelete = "DELETE"
)

// ApprovalTemplateMessage is approval template message
type ApprovalTemplateMessage struct {
	Action string           `json:"action"`
	Data   ApprovalTemplate `json:"data"`
}

// ApprovalTemplate is model to hold Approval Template
type ApprovalTemplate struct {
	ID                 string `json:"id,omitempty"`
	PartnerID          int    `json:"partnerId,omitempty"`
	ClientID           int    `json:"clientId,omitempty"`
	SiteID             int    `json:"siteID,omitempty"`
	Category           string `json:"category,omitempty"`
	ContactID          string `json:"contactId,omitempty"`
	ContactType        string `json:"contactType,omitempty"`
	ContactMethod      string `json:"contactMethod,omitempty"`
	ContactID2         string `json:"contactID2,omitempty"`
	ContactType2       string `json:"contactType2,omitempty"`
	ContactMethod2     string `json:"contactMethod2,omitempty"`
	AuthHDU            string `json:"authenticateHelpDeskUser,omitempty"`
	AuthApprover       string `json:"authenticateApprover,omitempty"`
	CallApproversSet   string `json:"callApproversSet,omitempty"`
	UserNameFormat     string `json:"newUserAccountCategory_UserNameFormat,omitempty"`
	SecurityGroup      string `json:"newUserAccountCategory_SecurityGroup,omitempty"`
	EmailAddressFormat string `json:"newMailAccountCategory_EmailAddressFormat,omitempty"`
	DistributionGroup  string `json:"newMailAccountCategory_DistributionGroup,omitempty"`
}
