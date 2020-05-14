package SOCUserRiskScoresSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//ProfileScoreUserData : Profile Scores User Data Structure
type ProfileScoreUserData struct {
	ProfileID         string                         `json:"profileID"`
	ProfileName       string                         `json:"profileName"`
	Score             int64                          `json:"score"`
	Categories        []CategoryProfileUserScoreData `json:"categories"`
	LastExecutionTime time.Time                      `json:"lastExecutionTime"`
}

//ProfileScoreUserType : ProfileScore GraphQL Schema
var ProfileScoreUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileScore",
	Fields: graphql.Fields{
		"profileID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreUserData); ok {
					return CurData.ProfileID, nil
				}
				return nil, nil
			},
		},

		"profileName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreUserData); ok {
					return CurData.ProfileName, nil
				}
				return nil, nil
			},
		},

		"score": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreUserData); ok {
					return CurData.Score, nil
				}
				return nil, nil
			},
		},

		"categories": &graphql.Field{
			Type: graphql.NewList(CategoryProfileUserScoreType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreUserData); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},

		"lastExecutionTime": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreUserData); ok {
					return CurData.LastExecutionTime, nil
				}
				return nil, nil
			},
		},
	},
})

//ProfileProtectUserData : Profile Protect User Data
type ProfileProtectUserData struct {
	PartnerID            string                 `json:"partnerID"`
	ClientID             string                 `json:"clientID"`
	SiteID               string                 `json:"siteID"`
	EndpointID           string                 `json:"endpointID"`
	UserID               string                 `json:"userID"`
	UserName             string                 `json:"userName"`
	Domain               string                 `json:"domain"`
	ProfileScoreUserData []ProfileScoreUserData `json:"profileScore"`
}

//ProfileProtectUserType : ProfileProtect GraphQL Schema
var ProfileProtectUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileProtectUser",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectUserData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"clientID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectUserData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},

		"endpointID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectUserData); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectUserData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"userID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectUserData); ok {
					return CurData.UserID, nil
				}
				return nil, nil
			},
		},

		"profileScore": &graphql.Field{
			Type: graphql.NewList(ProfileScoreUserType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectUserData); ok {
					return CurData.ProfileScoreUserData, nil
				}
				return nil, nil
			},
		},
	},
})

//ProfileProtectUserList : ProfileProtectUserList List struct
type ProfileProtectUserList struct {
	ProfileProtectUserData []ProfileProtectUserData `json:"profileProtectData"`
}

//ProfileProtectUserListType : ProfileProtectList GraphQL Schema
var ProfileProtectUserListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileProtectUserList",
	Fields: graphql.Fields{
		"profileProtectUserList": &graphql.Field{
			Type: graphql.NewList(ProfileProtectUserType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectUserList); ok {
					return CurData.ProfileProtectUserData, nil
				}
				return nil, nil
			},
		},
	},
})
