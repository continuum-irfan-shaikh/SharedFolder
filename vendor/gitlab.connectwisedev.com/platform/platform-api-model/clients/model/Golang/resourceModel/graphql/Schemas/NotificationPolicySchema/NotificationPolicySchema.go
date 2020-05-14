package NotificationPolicySchema

import (
	"github.com/graphql-go/graphql"
)

// NotificationPolicyData : NotificationPolicyData Structure
type NotificationPolicyData struct {
	SiteID       int64  `json:"SiteID"`
	TemplateID   int64  `json:"TemplateID"`
	TemplateName string `json:"TemplateName"`
	CreatedBy    int64  `json:"CreatedBy"`
	TimeZone     string `json:"TimeZone"`
}

// NotificationPolicyType : NotificationPolicy GraphQL Schema
var NotificationPolicyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NotificationPolicy",
	Fields: graphql.Fields{
		"SiteID": &graphql.Field{
			Type:        graphql.Int,
			Description: "SiteID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotificationPolicyData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"TemplateID": &graphql.Field{
			Type:        graphql.Int,
			Description: "TemplateID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotificationPolicyData); ok {
					return CurData.TemplateID, nil
				}
				return nil, nil
			},
		},

		"CreatedBy": &graphql.Field{
			Type:        graphql.Int,
			Description: "CreatedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotificationPolicyData); ok {
					return CurData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"TemplateName": &graphql.Field{
			Type:        graphql.String,
			Description: "TemplateName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotificationPolicyData); ok {
					return CurData.TemplateName, nil
				}
				return nil, nil
			},
		},

		"TimeZone": &graphql.Field{
			Type:        graphql.String,
			Description: "TimeZone",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotificationPolicyData); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},
	},
})

// NotificationPolicyListData : NotificationPolicyListData Structure
type NotificationPolicyListData struct {
	Status                 int64                    `json:"status"`
	NotificationPolicyList []NotificationPolicyData `json:"outdata"`
}

// NotificationPolicyListType : NotificationPolicyList GraphQL Schema
var NotificationPolicyListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NotificationPolicyListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotificationPolicyListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(NotificationPolicyType),
			Description: "NotificationPolicy list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotificationPolicyListData); ok {
					return CurData.NotificationPolicyList, nil
				}
				return nil, nil
			},
		},
	},
})
