package RebootScheduleTemplateSchema

import (
	"github.com/graphql-go/graphql"
)

// RebootScheduleTemplateData : RebootScheduleTemplateData Structure
type RebootScheduleTemplateData struct {
	TemplateName    string `json:"TemplateName"`
	RestartDays     string `json:"RestartDays"`
	RestartFromTime string `json:"RestartFromTime"`
	RestartToTime   string `json:"RestartToTime"`
}

// RebootScheduleTemplateType : RebootScheduleTemplate GraphQL Schema
var RebootScheduleTemplateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RebootScheduleTemplate",
	Fields: graphql.Fields{
		"TemplateName": &graphql.Field{
			Type:        graphql.String,
			Description: "TemplateName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleTemplateData); ok {
					return CurData.TemplateName, nil
				}
				return nil, nil
			},
		},

		"RestartDays": &graphql.Field{
			Type:        graphql.String,
			Description: "RestartDays",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleTemplateData); ok {
					return CurData.RestartDays, nil
				}
				return nil, nil
			},
		},

		"RestartFromTime": &graphql.Field{
			Type:        graphql.String,
			Description: "RestartFromTime for template",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleTemplateData); ok {
					return CurData.RestartFromTime, nil
				}
				return nil, nil
			},
		},

		"RestartToTime": &graphql.Field{
			Type:        graphql.String,
			Description: "RestartToTime for template",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleTemplateData); ok {
					return CurData.RestartToTime, nil
				}
				return nil, nil
			},
		},
	},
})

// RebootScheduleTemplateListData : RebootScheduleTemplateListData Structure
type RebootScheduleTemplateListData struct {
	Status                     int64                        `json:"status"`
	RebootScheduleTemplateList []RebootScheduleTemplateData `json:"outdata"`
}

// RebootScheduleTemplateListType : RebootScheduleTemplateList GraphQL Schema
var RebootScheduleTemplateListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RebootScheduleTemplateListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleTemplateListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(RebootScheduleTemplateType),
			Description: "Maintenance  list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleTemplateListData); ok {
					return CurData.RebootScheduleTemplateList, nil
				}
				return nil, nil
			},
		},
	},
})
