package SitesSchema

import (
	"github.com/graphql-go/graphql"
)

//SOCSitesData : SOCSitesData Structure
type SOCSitesData struct {
	PartnerID string `json:"partnerId"`
	SiteID    string `json:"siteId"`
	IsActive  bool   `json:"isactive"`
	Feature   string `json:"feature"`
}

//SOCSiteType : SOC Sites GraphQL Schema
var SOCSiteType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SOCSites",
	Fields: graphql.Fields{
		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "Partner ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCSitesData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCSitesData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"isactive": &graphql.Field{
			Type:        graphql.String,
			Description: "Is Active",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCSitesData); ok {
					return CurData.IsActive, nil
				}
				return nil, nil
			},
		},

		"feature": &graphql.Field{
			Type:        graphql.String,
			Description: "Feature",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCSitesData); ok {
					return CurData.Feature, nil
				}
				return nil, nil
			},
		},
	},
})


//SOCSitesListData : SOCSitesListData Structure
type SOCSitesListData struct {
	SecurityPartnerSites []SOCSitesData `json:"securityPartnerSiteList"`
}

