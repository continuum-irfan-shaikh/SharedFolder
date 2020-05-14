package SecurityAssessmentToolSchema

import (
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/ProfileProtectSchema"
	"github.com/graphql-go/graphql"
)

//SATSummaryResult type for api response
type SATSummaryResult struct {
	PartnerID             string         `json:"partnerId"`
	ClientID              string         `json:"clientId"`
	SiteID                string         `json:"siteId"`
	Endpoints             []CategoryData `json:"endpoints"`
	TotalMachines         string         `json:"totalmachines"`
	SiteSelfRank          string         `json:"siteselfrank"`
	OtherSitesMarketRank  string         `json:"othersitesmktrank"`
	SimilarRangeSitesRank string         `json:"similarrangesitesrank"`
	IndustryWideRank      string         `json:"industrywiderank"`
}

//CategoryData type for endpoint data
type CategoryData struct {
	Category string `json:"categoryName"`
	Count    string `json:"count"`
	Message  string `json:"message"`
}

//SATSummaryResultType  SATSummaryResult type of Graphql schema
var SATSummaryResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SATSummaryResult",
	Fields: graphql.Fields{
		"partnerId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"siteId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"clientId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"endpoints": &graphql.Field{
			Type: graphql.NewList(CategoryDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.Endpoints, nil
				}
				return nil, nil
			},
		},
		"totalmachines": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.TotalMachines, nil
				}
				return nil, nil
			},
		},
		"siteselfrank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.SiteSelfRank, nil
				}
				return nil, nil
			},
		},
		"othersitesmktrank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.OtherSitesMarketRank, nil
				}
				return nil, nil
			},
		},
		"similarrangesitesrank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.SimilarRangeSitesRank, nil
				}
				return nil, nil
			},
		},
		"industrywiderank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryResult); ok {
					return CurData.IndustryWideRank, nil
				}
				return nil, nil
			},
		},
	},
})

//CategoryDataType : endpoint or user summary report GraphQL Schema
var CategoryDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CategoryData",
	Fields: graphql.Fields{
		"categoryName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Category, nil
				}
				return nil, nil
			},
		},
		"count": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Count, nil
				}
				return nil, nil
			},
		},
		"message": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Message, nil
				}
				return nil, nil
			},
		},
	},
})

// EndpointDataConnectionDefinition : EndpointDataConnectionDefinition structure
var EndpointDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "endpoints",
	NodeType: CategoryDataType,
})

//SATSummaryUserResult type for api response
type SATSummaryUserResult struct {
	PartnerID                 string         `json:"partnerId"`
	ClientID                  string         `json:"clientId"`
	SiteID                    string         `json:"siteId"`
	Users                     []CategoryData `json:"users"`
	TotalUsers                string         `json:"totalusers"`
	TotalRiskUsers            string         `json:"totalriskusers"`
	SiteUserSelfRank          string         `json:"siteuserselfrank"`
	OtherSiteUserMktfRank     string         `json:"othersitesusermktrank"`
	SimilarRangeUserSitesRank string         `json:"similarrangesitesuserrank"`
	IndustryWideUserRank      string         `json:"industrywideuserrank"`
}

//AssessSummaryResult type for api response
type AssessSummaryResult struct {
	PartnerID                    string                        `json:"partnerId"`
	ClientID                     string                        `json:"clientId"`
	SiteID                       string                        `json:"siteId"`
	AssessEndpointSummaryResults []AssessEndpointSummaryResult `json:"endpoints"`
	DarkWebUserDetails           []AssessDarkWebDetail         `json:"darkwebDetails"`
	UserDetails                  []AssessUserDetail            `json:"userDetails"`
	TotalMachines                string                        `json:"totalMachines"`
	Domain                       string                        `json:"domain"`
	Email                        string                        `json:"email"`
	EmailBreachCount             int                           `json:"emailBreachCount"`
}

//AssessDarkWebDetail type for all Dark web user data
type AssessDarkWebDetail struct {
	User         string `json:"user"`
	Credential   string `json:"credential"`
	BreachSource string `json:"breachSource"`
	PublishDate  string `json:"publishDate"`
}

//DarkWebDetails type for all Dark web user data
type DarkWebDetails struct {
	PartnerID          string                             `json:"partnerId"`
	ClientID           string                             `json:"clientId"`
	SiteID             string                             `json:"siteId"`
	Domain             string                             `json:"domain"`
	Email              string                             `json:"email"`
	EmailBreachCount   int                                `json:"emailBreachCount"`
	DarkWebUserDetails []ProfileProtectSchema.SpyCloudRow `json:"userData"`
}

//AssessEndpointSummaryResult type for api response
type AssessEndpointSummaryResult struct {
	UnsecuredNetworkDirectory string `json:"unsecuredNetworkDirectory"`
	RDPEnabled                string `json:"remoteDesktopAccess"`
	MissingAdvanceProtection  string `json:"missingAdvanceProtection"`
	EndpointName              string `json:"endpointName"`
	EndpointID                string `json:"endpointId"`
}

//AssessUserDetail - user details
type AssessUserDetail struct {
	Username              string `json:"userName"`
	PasswordComplexity    string `json:"passwordComplexity"`
	RdpEnabled            string `json:"remoteDesktopAccess"`
	AccountsNotLoggedin   string `json:"accountsNotLoggedIn"`
	AccountsLastLoginTime string `json:"accountsLastLoginTime"`
	StaleAccount          string `json:"staleAccount"`
	Domain                string `json:"domain"`
}

//UserDetail - user details
type UserDetail struct {
	PartnerID         string             `json:"partnerId"`
	ClientID          string             `json:"clientId"`
	SiteID            string             `json:"siteId"`
	AssessUserDetails []AssessUserDetail `json:"users"`
}

//AssetSummaryResultType  AssessSummaryResult type of Graphql schema
var AssetSummaryResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AssetSummaryResult",
	Fields: graphql.Fields{
		"partnerId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"siteId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"clientId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"domain": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.Domain, nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.Email, nil
				}
				return nil, nil
			},
		},
		"endpoints": &graphql.Field{
			Type: graphql.NewList(AssessEndpointSummaryResultType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.AssessEndpointSummaryResults, nil
				}
				return nil, nil
			},
		},
		"darkwebDetails": &graphql.Field{
			Type: graphql.NewList(DarkwebDetailType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.DarkWebUserDetails, nil
				}
				return nil, nil
			},
		},
		"userDetails": &graphql.Field{
			Type: graphql.NewList(UserDetailType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.UserDetails, nil
				}
				return nil, nil
			},
		},
		"totalMachines": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessSummaryResult); ok {
					return CurData.TotalMachines, nil
				}
				return nil, nil
			},
		},
	},
})

//AssessEndpointSummaryResultType : endpoint or user summary report GraphQL Schema
var AssessEndpointSummaryResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AssessEndpointSummaryResult",
	Fields: graphql.Fields{
		"unsecuredNetworkDirectory": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessEndpointSummaryResult); ok {
					return CurData.UnsecuredNetworkDirectory, nil
				}
				return nil, nil
			},
		},
		"remoteDesktopAccess": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessEndpointSummaryResult); ok {
					return CurData.RDPEnabled, nil
				}
				return nil, nil
			},
		},
		"missingAdvanceProtection": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessEndpointSummaryResult); ok {
					return CurData.MissingAdvanceProtection, nil
				}
				return nil, nil
			},
		},
		"endpointName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessEndpointSummaryResult); ok {
					return CurData.EndpointName, nil
				}
				return nil, nil
			},
		},
		"endpointId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessEndpointSummaryResult); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},
	},
})

//DarkwebDetailType : Darkweb summary report GraphQL Schema
var DarkwebDetailType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DarkwebDetailType",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessDarkWebDetail); ok {
					return CurData.User, nil
				}
				return nil, nil
			},
		},
		"credential": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessDarkWebDetail); ok {
					return CurData.Credential, nil
				}
				return nil, nil
			},
		},
		"breachSource": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessDarkWebDetail); ok {
					return CurData.BreachSource, nil
				}
				return nil, nil
			},
		},
		"publishDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessDarkWebDetail); ok {
					return CurData.PublishDate, nil
				}
				return nil, nil
			},
		},
	},
})

//UserDetailType : User summary report GraphQL Schema
var UserDetailType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserDetailType",
	Fields: graphql.Fields{
		"userName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessUserDetail); ok {
					return CurData.Username, nil
				}
				return nil, nil
			},
		},
		"passwordComplexity": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessUserDetail); ok {
					return CurData.PasswordComplexity, nil
				}
				return nil, nil
			},
		},
		"remoteDesktopAccess": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessUserDetail); ok {
					return CurData.RdpEnabled, nil
				}
				return nil, nil
			},
		},
		"accountsNotLoggedIn": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessUserDetail); ok {
					return CurData.AccountsNotLoggedin, nil
				}
				return nil, nil
			},
		},
		"accountsLastLoginTime": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessUserDetail); ok {
					return CurData.AccountsLastLoginTime, nil
				}
				return nil, nil
			},
		},
		"staleAccount": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessUserDetail); ok {
					return CurData.StaleAccount, nil
				}
				return nil, nil
			},
		},
		"domain": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssessUserDetail); ok {
					return CurData.Domain, nil
				}
				return nil, nil
			},
		},
	},
})

//SATSummaryUserResultType  SATSummaryUserResult type of Graphql schema
var SATSummaryUserResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SATSummaryUserResult",
	Fields: graphql.Fields{
		"partnerId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"siteId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"clientId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"users": &graphql.Field{
			Type: graphql.NewList(CategoryDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.Users, nil
				}
				return nil, nil
			},
		},
		"totalusers": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.TotalUsers, nil
				}
				return nil, nil
			},
		},
		"totalriskusers": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.TotalRiskUsers, nil
				}
				return nil, nil
			},
		},
		"siteuserselfrank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.SiteUserSelfRank, nil
				}
				return nil, nil
			},
		},
		"othersitesusermktrank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.OtherSiteUserMktfRank, nil
				}
				return nil, nil
			},
		},
		"similarrangesitesuserrank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.SimilarRangeUserSitesRank, nil
				}
				return nil, nil
			},
		},
		"industrywideuserrank": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SATSummaryUserResult); ok {
					return CurData.IndustryWideUserRank, nil
				}
				return nil, nil
			},
		},
	},
})

// UserDataConnectionDefinition : UserDataConnectionDefinition structure
var UserDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "users",
	NodeType: CategoryDataType,
})

//SatUsageResponse for response from the post request
type SatUsageResponse struct {
	StatusCode  int    `json:"statusCode"`
	Description string `json:"description"`
}

//SatUsageReport struct to store response of profiling ms api
type SatUsageReport struct {
	PartnerID              string `json:"partnerId"`
	SiteID                 string `json:"siteId"`
	ClientID               string `json:"clientId"`
	CreatedBy              string `json:"createdby"`
	IsDarkWebAvailable     bool   `json:"isdarkwebavailable"`
	IsEndpointAvailable    bool   `json:"isendpointavailable"`
	IsUserAccountAvailable bool   `json:"isuseraccountavailable"`
	DomainList             string `json:"domainlist"`
	MaskingPreference      string `json:"masking_preference"`
	ReportType             string `json:"reporttype"`
	DarkwebRecordCount     int64  `json:"darkwebRecordCount"`
}

//SatUsageResponseType for graphQL structure
var SatUsageResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SatUsageResponse",
	Fields: graphql.Fields{
		"statusCode": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SatUsageResponse); ok {
					return CurData.StatusCode, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SatUsageResponse); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
	},
})

//DarkwebUserPreference struct request params
type DarkwebUserPreference struct {
	PartnerID string `json:"partnerId"`
	UserID    string `json:"userId"`
}

//DarkwebUserPreferenceResponse struct for the post request
type DarkwebUserPreferenceResponse struct {
	PartnerID         string `json:"partnerId"`
	UserID            string `json:"userId"`
	CreatedAt         int    `json:"createdAt"`
	MaskingPreference string `json:"masking_preference"`
}

//DarkwebUserPreferenceResponseType for graphQL structure
var DarkwebUserPreferenceResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DarkwebUserPreferenceResponse",
	Fields: graphql.Fields{
		"partnerId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DarkwebUserPreferenceResponse); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"userId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DarkwebUserPreferenceResponse); ok {
					return CurData.UserID, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DarkwebUserPreferenceResponse); ok {
					return CurData.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"masking_preference": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DarkwebUserPreferenceResponse); ok {
					return CurData.MaskingPreference, nil
				}
				return nil, nil
			},
		},
	},
})
