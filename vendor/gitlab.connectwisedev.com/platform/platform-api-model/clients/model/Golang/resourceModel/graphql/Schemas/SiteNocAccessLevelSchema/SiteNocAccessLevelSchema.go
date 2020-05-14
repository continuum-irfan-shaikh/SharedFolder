package SiteNocAccessLevelSchema

import "github.com/graphql-go/graphql"

type SiteNocResourceData struct {
	RegID        string `json:"regID"`
	TemplateID   string `json:"templateID"`
	TemplateName string `json:"templateName"`
	AccessLevel  string `json:"accessLevel"`
}

//SiteNocResourceDataType : SiteNocResourceDataType GraphQL Schema
var SiteNocResourceDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SiteNocResourceData",
	Fields: graphql.Fields{
		"regID": &graphql.Field{
			Type:        graphql.String,
			Description: "Reg ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocResourceData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"templateID": &graphql.Field{
			Type:        graphql.String,
			Description: "template ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocResourceData); ok {
					return CurData.TemplateID, nil
				}
				return nil, nil
			},
		},

		"templateName": &graphql.Field{
			Type:        graphql.String,
			Description: "Nemplate Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocResourceData); ok {
					return CurData.TemplateName, nil
				}
				return nil, nil
			},
		},

		"accessLevel": &graphql.Field{
			Type:        graphql.String,
			Description: "Access Level",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocResourceData); ok {
					return CurData.AccessLevel, nil
				}
				return nil, nil
			},
		},
	},
})

type SiteNocAccessData struct {
	IsNOC                    string                `json:"isNOC"`
	SiteID                   string                `json:"siteID"`
	UpdatedBy                string                `json:"updatedBy"`
	TemplateID               string                `json:"templateID"`
	TemplateName             string                `json:"templateName"`
	AccessLevel              string                `json:"accessLevel"`
	IsAppliedAtResourceLevel string                `json:"isAppliedAtResourceLevel"`
	Resources                []SiteNocResourceData `json:"resources"`
}

//SiteNocAccessDataType : SiteNocAccessDataType GraphQL Schema
var SiteNocAccessDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SiteNocAccessData",
	Fields: graphql.Fields{
		"resources": &graphql.Field{
			Type:        graphql.NewList(SiteNocResourceDataType),
			Description: "resources list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.Resources, nil
				}
				return nil, nil
			},
		},

		"templateID": &graphql.Field{
			Type:        graphql.String,
			Description: "template ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.TemplateID, nil
				}
				return nil, nil
			},
		},

		"templateName": &graphql.Field{
			Type:        graphql.String,
			Description: "Nemplate Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.TemplateName, nil
				}
				return nil, nil
			},
		},

		"accessLevel": &graphql.Field{
			Type:        graphql.String,
			Description: "Access Level",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.AccessLevel, nil
				}
				return nil, nil
			},
		},

		"isNOC": &graphql.Field{
			Type:        graphql.String,
			Description: "isNoc - 1 or 0",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.IsNOC, nil
				}
				return nil, nil
			},
		},

		"isAppliedAtResourceLevel": &graphql.Field{
			Type:        graphql.String,
			Description: "isAppliedAtResourceLevel - 1 or 0",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.IsAppliedAtResourceLevel, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type:        graphql.String,
			Description: "site id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"updatedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "email id - updated by",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteNocAccessData); ok {
					return CurData.UpdatedBy, nil
				}
				return nil, nil
			},
		},
	},
})
