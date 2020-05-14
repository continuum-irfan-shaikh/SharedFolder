package ProfileProtectSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//ProfileScoreData : ProfileScore Data Structure
type ProfileScoreData struct {
	ProfileID   string                     `json:"profileID"`
	ProfileName string                     `json:"profileName"`
	Score       int64                      `json:"score"`
	Categories  []CategoryProfileScoreData `json:"categories"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
}

//ProfileScoreType : ProfileScore GraphQL Schema
var ProfileScoreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileScore",
	Fields: graphql.Fields{
		"profileID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreData); ok {
					return CurData.ProfileID, nil
				}
				return nil, nil
			},
		},

		"profileName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreData); ok {
					return CurData.ProfileName, nil
				}
				return nil, nil
			},
		},

		"score": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreData); ok {
					return CurData.Score, nil
				}
				return nil, nil
			},
		},

		"categories": &graphql.Field{
			Type: graphql.NewList(CategoryProfileScoreType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreData); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},

		"updatedAt": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileScoreData); ok {
					return CurData.UpdatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

//ProfileProtectData : Profile Protect Data
type ProfileProtectData struct {
	PartnerID        string             `json:"partnerID"`
	ClientID         string             `json:"clientID"`
	SiteID           string             `json:"siteID"`
	EndpointID       string             `json:"endpointID"`
	ProfileScoreData []ProfileScoreData `json:"profileScore"`
}

//ProfileProtectType : ProfileProtect GraphQL Schema
var ProfileProtectType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileProtect",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"clientID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"endpointID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectData); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},

		"profileScoreData": &graphql.Field{
			Type: graphql.NewList(ProfileScoreType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectData); ok {
					return CurData.ProfileScoreData, nil
				}
				return nil, nil
			},
		},
	},
})

//ProfileProtectConnectionDefinition : ProfileProtectConnectionDefinition structure
var ProfileProtectConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "ProfileProtect",
	NodeType: ProfileProtectType,
})

//ProfileProtectList : ProfileProtect List struct
type ProfileProtectList struct {
	ProfileProtectData []ProfileProtectData `json:"profileProtectData"`
}

//ProfileProtectListType : ProfileProtectList GraphQL Schema
var ProfileProtectListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileProtectList",
	Fields: graphql.Fields{
		"profileProtectList": &graphql.Field{
			Type: graphql.NewList(ProfileProtectType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectList); ok {
					return CurData.ProfileProtectData, nil
				}
				return nil, nil
			},
		},
	},
})
