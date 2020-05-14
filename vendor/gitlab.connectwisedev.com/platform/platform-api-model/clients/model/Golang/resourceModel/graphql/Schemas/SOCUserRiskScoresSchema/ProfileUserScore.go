package SOCUserRiskScoresSchema

import (
	"github.com/graphql-go/graphql"
)

//ProfileUsersScore : Profile Score Data Structure
type ProfileUsersScore struct {
	UserID     string                         `json:"userID"`
	UserName   string                         `json:"userName"`
	Domain     string                         `json:"domain"`
	Score      int64                          `json:"score"`
	PartnerID  string                         `json:"partnerID"`
	SiteID     string                         `json:"siteID"`
	EndpointID string                         `json:"endpointID"`
	Categories []CategoryProfileUserScoreData `json:"categories"`
}

//ProfileUserScoreType : ProfileUserScoreType GraphQL Schema
var ProfileUserScoreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileScore",
	Fields: graphql.Fields{
		"userID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.UserID, nil
				}
				return nil, nil
			},
		},
		"userName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.UserName, nil
				}
				return nil, nil
			},
		},
		"domain": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.Domain, nil
				}
				return nil, nil
			},
		},
		"partnerID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"endpointID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},

		"score": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.Score, nil
				}
				return nil, nil
			},
		},

		"categories": &graphql.Field{
			Type: graphql.NewList(CategoryProfileUserScoreType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileUsersScore); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},
	},
})
