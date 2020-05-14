package NOCAccessLevelSchema

import (
	"github.com/graphql-go/graphql"
)

// NOCAccessLevelData : NOCAccessLevelData Structure
type NOCAccessLevelData struct {
	SiteID          int64  `json:"siteID"`
	TemplateID      int64  `json:"templateID"`
	TemplateName    string `json:"templateName"`
	AccessLevel     int64  `json:"accessLevel"`
	AccessLevelDesc string `json:"accessLevelDesc"`
}

//NOCAccessLevelType : NOCAccessLevel GraphQL Schema
var NOCAccessLevelType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NOCAccessLevel",
	Fields: graphql.Fields{
		"siteID": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"templateID": &graphql.Field{
			Type:        graphql.String,
			Description: "Template ID for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelData); ok {
					return CurData.TemplateID, nil
				}
				return nil, nil
			},
		},

		"templateName": &graphql.Field{
			Type:        graphql.String,
			Description: "Template Name for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelData); ok {
					return CurData.TemplateName, nil
				}
				return nil, nil
			},
		},

		"accessLevel": &graphql.Field{
			Type:        graphql.String,
			Description: "AccessLevel for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelData); ok {
					return CurData.AccessLevel, nil
				}
				return nil, nil
			},
		},

		"accessLevelDesc": &graphql.Field{
			Type:        graphql.String,
			Description: "AccessLevel Desc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelData); ok {
					return CurData.AccessLevelDesc, nil
				}
				return nil, nil
			},
		},
	},
})

//NOCAccessLevelListData : NOCAccessLevelListData Structure
type NOCAccessLevelListData struct {
	Status             int64                `json:"status"`
	NOCAccessLevelList []NOCAccessLevelData `json:"outdata"`
}

//NOCAccessLevelListType : NOCAccessLevelList GraphQL Schema
var NOCAccessLevelListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NOCAccessLevelListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(NOCAccessLevelType),
			Description: "NOCAccessLevel list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelListData); ok {
					return CurData.NOCAccessLevelList, nil
				}
				return nil, nil
			},
		},
	},
})
