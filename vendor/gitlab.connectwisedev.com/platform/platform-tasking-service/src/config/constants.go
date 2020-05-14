package config

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/types"
)

const (
	UserKeyCTX          types.CTXKeyStr = "user"
	UserEndPointsKeyCTX types.CTXKeyStr = "userEndPoints"
	PartnerIDKeyCTX     types.CTXKeyStr = "partnerID"
	EndpointIDKeyCTX    types.CTXKeyStr = "endpointID"
	SiteIDKeyCTX        types.CTXKeyStr = "siteID"
	ClientIDKeyCTX      types.CTXKeyStr = "clientID"
	TriggerTypeIDKeyCTX types.CTXKeyStr = "triggerType"

	MaxWeekDaysInSchedule int = 4
)
