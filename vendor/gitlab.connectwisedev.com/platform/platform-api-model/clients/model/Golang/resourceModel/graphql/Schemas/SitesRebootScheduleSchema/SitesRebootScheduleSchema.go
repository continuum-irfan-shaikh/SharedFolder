package SitesRebootScheduleSchema

import (
	"strconv"

	"github.com/graphql-go/graphql"
)

//RebootScheduleDetailsData : RebootScheduleDetailsData Structure
// Site id, TemplateID, RebootTemplateName
type RebootScheduleDetailsData struct {
	SiteID             int64  `json:"SiteID"`
	TemplateID         int64  `json:"TemplateID"`
	RebootTemplateName string `json:"RebootTemplateName"`
}

//RebootScheduleDetailsType : RebootScheduleDetails GraphQL Schema
var RebootScheduleDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "rebootScheduleDetails",
	Fields: graphql.Fields{
		"SiteID": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleDetailsData); ok {
					return strconv.FormatInt(CurData.SiteID, 10), nil
				}
				return nil, nil
			},
		},

		"RebootTemplateName": &graphql.Field{
			Type:        graphql.String,
			Description: "site postal code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleDetailsData); ok {
					return CurData.RebootTemplateName, nil
				}
				return nil, nil
			},
		},

		"TemplateID": &graphql.Field{
			Type:        graphql.String,
			Description: "Site TemplateID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RebootScheduleDetailsData); ok {
					return strconv.FormatInt(CurData.TemplateID, 10), nil
				}
				return nil, nil
			},
		},
	},
})

//SitesRebootScheduleData : SitesRebootScheduleData Structure
type SitesRebootScheduleData struct {
	Status             int64                       `json:"status"`
	RebootScheduleList []RebootScheduleDetailsData `json:"outdata"`
}

//SitesRebootScheduleType : SitesRebootScheduleType GraphQL Schema
var SitesRebootScheduleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesRebootSchedule",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesRebootScheduleData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        graphql.NewList(RebootScheduleDetailsType),
			Description: "Site list of reboot schedule",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesRebootScheduleData); ok {
					return CurData.RebootScheduleList, nil
				}
				return nil, nil
			},
		},
	},
})
