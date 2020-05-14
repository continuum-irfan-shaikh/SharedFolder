package v2

import (
	"time"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
)

// Task represents task v2 struct
type Task struct {
	ID                 string           `json:"id"                     valid:"unsettableByUsers"` // uuid
	PartnerID          string           `json:"partnerId"              valid:"unsettableByUsers"`
	DefinitionID       string           `json:"definitionID"           valid:"optional"`         // uuid
	LastExecutionID    string           `json:"-"                      valid:"-"`                // uuid
	OriginID           string           `json:"originId"               valid:"requiredForUsers"` // script or patch ID uuid
	Name               string           `json:"name"`
	Description        string           `json:"description"            valid:"unsettableByUsers"`
	CreatedAt          time.Time        `json:"createdAt"              valid:"unsettableByUsers"`
	CreatedBy          string           `json:"createdBy"              valid:"unsettableByUsers"`
	Type               string           `json:"type"                   valid:"validType"`
	Parameters         string           `json:"parameters"             valid:"-"`
	External           bool             `json:"external"               valid:"-"`
	ResultWebhook      string           `json:"resultWebhook"          valid:"url"`
	IsRequireNOCAccess bool             `json:"isRequireNOCAccess"     valid:"-"`
	ModifiedBy         string           `json:"modifiedBy"             valid:"unsettableByUsers"`
	ModifiedAt         time.Time        `json:"modifiedAt"             valid:"unsettableByUsers"`
	Schedule           tasking.Schedule `json:"schedule"` // stored as json string
	Targets            []Target
}
