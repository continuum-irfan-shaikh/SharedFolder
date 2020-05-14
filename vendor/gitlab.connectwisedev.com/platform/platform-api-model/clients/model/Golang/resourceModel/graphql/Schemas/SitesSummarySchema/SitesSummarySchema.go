package SitesSummarySchema

import (
	"github.com/graphql-go/graphql"
	"strconv"
)

//SiteDetailsData : SiteDetailsData Structure
// partner id, Site id, Site name, timezone, code, location info: city, state, country, address, postalCode
type SiteDetailsData struct {
	MemberID    int64  `json:"memberId"`
	SiteID      int64  `json:"siteId"`
	SiteCode    string `json:"siteCode"`
	SiteName    string `json:"siteName"`
	City        string `json:"siteCity"`
	SiteAddress string `json:"siteAddress"`
	TimeZone    string `json:"timeZone"`
	PostalCode  string `json:"sitePostalCode"`
	State       string `json:"siteState"`
	Country     int64  `json:"siteCountry"`
}

//SiteDetailsType : SiteDetails GraphQL Schema
var SiteDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "siteDetails",
	Fields: graphql.Fields{
		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return strconv.FormatInt(CurData.SiteID, 10), nil
				}
				return nil, nil
			},
		},

		"memberId": &graphql.Field{
			Type:        graphql.String,
			Description: "Member ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.MemberID, nil
				}
				return nil, nil
			},
		},

		"siteCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.SiteCode, nil
				}
				return nil, nil
			},
		},

		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},

		"siteAddress": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Address",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.SiteAddress, nil
				}
				return nil, nil
			},
		},

		"sitePostalCode": &graphql.Field{
			Type:        graphql.String,
			Description: "site postal code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.PostalCode, nil
				}
				return nil, nil
			},
		},

		"timeZone": &graphql.Field{
			Type:        graphql.String,
			Description: "time Zone",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},

		"siteCity": &graphql.Field{
			Type:        graphql.String,
			Description: "Site City",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.City, nil
				}
				return nil, nil
			},
		},

		"siteState": &graphql.Field{
			Type:        graphql.String,
			Description: "Site State",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.State, nil
				}
				return nil, nil
			},
		},

		"siteCountry": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Country",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteDetailsData); ok {
					return CurData.Country, nil
				}
				return nil, nil
			},
		},
	},
})

//SitesSummaryData : SitesSummaryData Structure
type SitesSummaryData struct {
	Status   int64             `json:"status"`
	SiteList []SiteDetailsData `json:"outdata"`
}

//SitesSummaryType : SitesSummaryType GraphQL Schema
var SitesSummaryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesSummary",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesSummaryData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        graphql.NewList(SiteDetailsType),
			Description: "Site list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesSummaryData); ok {
					return CurData.SiteList, nil
				}
				return nil, nil
			},
		},
	},
})
