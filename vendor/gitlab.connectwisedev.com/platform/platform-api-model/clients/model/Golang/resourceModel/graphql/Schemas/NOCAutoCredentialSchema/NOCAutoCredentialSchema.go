package NOCAutoCredentialSchema

import (
	"github.com/graphql-go/graphql"
)

// NOCAutoCredentialData : NOCAutoCredentialData Structure
type NOCAutoCredentialData struct {
	AutoCredOpt bool   `json:"autoCredOpt"`
	ResType     string `json:"resType"`
	SiteID      int64  `json:"siteID"`
	UserName    string `json:"userName"`
}

// NOCAutoCredentialType : NOCAutoCredential GraphQL Schema
var NOCAutoCredentialType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NOCAutoCredential",
	Fields: graphql.Fields{
		"autoCredOpt": &graphql.Field{
			Type:        graphql.String,
			Description: "autoCredOpt",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAutoCredentialData); ok {
					return CurData.AutoCredOpt, nil
				}
				return nil, nil
			},
		},

		"resType": &graphql.Field{
			Type:        graphql.String,
			Description: "res Type for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAutoCredentialData); ok {
					return CurData.ResType, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type:        graphql.Int,
			Description: "siteID for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAutoCredentialData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"userName": &graphql.Field{
			Type:        graphql.String,
			Description: "User Name for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAutoCredentialData); ok {
					return CurData.UserName, nil
				}
				return nil, nil
			},
		},
	},
})

// NOCAutoCredentialListData : NOCAutoCredentialListData Structure
type NOCAutoCredentialListData struct {
	Status                int64                   `json:"status"`
	NOCAutoCredentialList []NOCAutoCredentialData `json:"outdata"`
}

// NOCAutoCredentialListType : NOCAutoCredentialList GraphQL Schema
var NOCAutoCredentialListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NOCAutoCredentialListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAutoCredentialListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(NOCAutoCredentialType),
			Description: "NOCAutoCredential list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAutoCredentialListData); ok {
					return CurData.NOCAutoCredentialList, nil
				}
				return nil, nil
			},
		},
	},
})
