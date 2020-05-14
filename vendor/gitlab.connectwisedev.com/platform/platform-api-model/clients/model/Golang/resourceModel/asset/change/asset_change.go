package change

import "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/asset"


type AssetCollectionChange struct {
	EndpointID         string                `json:"endpointID"`
	PartnerID          string                `json:"partnerID"`
	EndpointType       string                `json:"endpointType,omitempty"`
	OsType             string                `json:"osType,omitempty"`
	ClientID           string                `json:"clientID"`
	SiteID             string                `json:"siteID"`
	RegID              string                `json:"regID,omitempty"`
	ChangedAttributes  []string              `json:"changedAttributes,omitempty"`
	UpdatedCollection  asset.AssetCollection `json:"updated_collection,omitempty"`
	ExistingCollection asset.AssetCollection `json:"existing_collection,omitempty"`
}
