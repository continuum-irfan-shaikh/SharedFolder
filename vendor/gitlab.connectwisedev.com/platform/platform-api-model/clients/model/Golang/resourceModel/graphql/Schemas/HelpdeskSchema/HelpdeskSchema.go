package HelpdeskSchema

import "github.com/graphql-go/graphql"

//HelpdeskData : HelpdeskData structure
type HelpdeskData struct {
	ServiceID    int    `json:"ServiceID"`
	TemplateID   int64  `json:"TemplateID"`
	TemplateName string `json:"TemplateName"`
}

//HelpdeskType : HelpdeskType Schema Definition
var HelpdeskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HelpdeskDetails",
	Fields: graphql.Fields{
		"ServiceID": &graphql.Field{
			Type:        graphql.String,
			Description: "Helpdesk ServiceID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpdeskData); ok {
					return CurData.ServiceID, nil
				}
				return nil, nil
			},
		},
		"TemplateID": &graphql.Field{
			Type:        graphql.String,
			Description: "Helpdesk TemplateID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpdeskData); ok {
					return CurData.TemplateID, nil
				}
				return nil, nil
			},
		},
		"TemplateName": &graphql.Field{
			Type:        graphql.String,
			Description: "Helpdesk TemplateName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpdeskData); ok {
					return CurData.TemplateName, nil
				}
				return nil, nil
			},
		},
	},
})

//HelpdeskListData : Structure for HelpdeskData Rest Output
type HelpdeskListData struct {
	Status       int64          `json:"status"`
	HelpdeskList []HelpdeskData `json:"outdata"`
}

//HelpdeskListType : Object of HelpdeskListType
var HelpdeskListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HelpdeskList",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpdeskListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        graphql.NewList(HelpdeskType),
			Description: "Helpdesk Notification",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpdeskListData); ok {
					return CurData.HelpdeskList, nil
				}
				return nil, nil
			},
		},
	},
})

//HelpdeskSiteData : HelpdeskSiteData structure
type HelpdeskSiteData struct {
	SiteID       int64  `json:"SiteID"`
	TemplateID   int64  `json:"TemplateID"`
	TemplateName string `json:"TemplateName"`
}

//HelpdeskSiteListData : Structure for HelpdeskData Rest Output
type HelpdeskSiteListData struct {
	Status           int64              `json:"status"`
	HelpdeskSiteList []HelpdeskSiteData `json:"outdata"`
}
