package MaintenanceScheduleSchema

import (
	"github.com/graphql-go/graphql"
)

// MaintenanceScheduleData : MaintenanceScheduleData Structure
type MaintenanceScheduleData struct {
	SiteID             int64  `json:"SiteID"`
	TemplateID         int64  `json:"TemplateID"`
	RebootTemplateName string `json:"RebootTemplateName"`
}

// MaintenanceScheduleType : MaintenanceSchedule GraphQL Schema
var MaintenanceScheduleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MaintenanceSchedule",
	Fields: graphql.Fields{
		"SiteID": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MaintenanceScheduleData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"TemplateID": &graphql.Field{
			Type:        graphql.String,
			Description: "Template ID for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MaintenanceScheduleData); ok {
					return CurData.TemplateID, nil
				}
				return nil, nil
			},
		},

		"RebootTemplateName": &graphql.Field{
			Type:        graphql.String,
			Description: "Template Name for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MaintenanceScheduleData); ok {
					return CurData.RebootTemplateName, nil
				}
				return nil, nil
			},
		},
	},
})

//MaintenanceScheduleListData : MaintenanceScheduleListData Structure
type MaintenanceScheduleListData struct {
	Status                  int64                     `json:"status"`
	MaintenanceScheduleList []MaintenanceScheduleData `json:"outdata"`
}

//MaintenanceScheduleListType : MaintenanceScheduleList GraphQL Schema
var MaintenanceScheduleListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MaintenanceScheduleListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MaintenanceScheduleListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(MaintenanceScheduleType),
			Description: "Maintenance  list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MaintenanceScheduleListData); ok {
					return CurData.MaintenanceScheduleList, nil
				}
				return nil, nil
			},
		},
	},
})
