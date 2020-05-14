package SingleAntiVirusStatusSchema

import (
	"github.com/graphql-go/graphql"
)

type (
	SingleAntiVirusStatusData struct {
		PartnerID       string `json:"partnerId"`
		SiteID          string `json:"siteId"`
		EndpointID      string `json:"endpointId"`
		RegID           int64  `json:"regId"`
		AntiVirus       string `json:"antiVirus"`
		Version         string `json:"version"`
		AntiVirusStatus string `json:"antiVirusStatus"`
		WebrootStatus   string `json:"webrootStatus"`
	}
)

var SingleAntiVirusStatusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SingleAntiVirusStatus",
	Fields: graphql.Fields{
		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific partner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"antiVirus": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of antivirus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.AntiVirus, nil
				}
				return nil, nil
			},
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "Version of antivirus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.Version, nil
				}
				return nil, nil
			},
		},

		"antiVirusStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Status of antivirus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.AntiVirusStatus, nil
				}
				return nil, nil
			},
		},

		"webrootStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Webroot status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.WebrootStatus, nil
				}
				return nil, nil
			},
		},

		"endpointId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleAntiVirusStatusData); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},
	},
})
