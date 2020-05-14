package AllSiteSchema

import (
	"github.com/graphql-go/graphql"
)

// AllSitesData : AllSitesData Structure
type AllSitesData struct {
	DesktopOption    string `json:"desktopOption"`
	DesktopServiceID int64  `json:"desktopServiceId"`
	HelpDesk         string `json:"helpDesk"`
	IsActive         string `json:"isActive"`
	IsEnabled        string `json:"isEnabled"`
	ServerOption     string `json:"serverOption"`
	ServerServiceID  int64  `json:"serverServiceId"`
	SiteCode         string `json:"siteCode"`
	SiteID           int64  `json:"siteId"`
	SiteName         string `json:"siteName"`
}

// AllSitesType : AllSites GraphQL Schema
var AllSitesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AllSites",
	Fields: graphql.Fields{
		"desktopOption": &graphql.Field{
			Type:        graphql.String,
			Description: "desktopOption Desktop plan",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.DesktopOption, nil
				}
				return nil, nil
			},
		},
		"desktopServiceId": &graphql.Field{
			Type:        graphql.Int,
			Description: "desktopServiceId plan id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.DesktopServiceID, nil
				}
				return nil, nil
			},
		},
		"helpDesk": &graphql.Field{
			Type:        graphql.Int,
			Description: "helpDesk on site ",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.HelpDesk, nil
				}
				return nil, nil
			},
		},
		"isActive": &graphql.Field{
			Type:        graphql.String,
			Description: "active site ",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.IsActive, nil
				}
				return nil, nil
			},
		},
		"isEnabled": &graphql.Field{
			Type:        graphql.String,
			Description: "isEnabled site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.IsEnabled, nil
				}
				return nil, nil
			},
		},
		"serverOption": &graphql.Field{
			Type:        graphql.String,
			Description: "serverOption server plan ",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.ServerOption, nil
				}
				return nil, nil
			},
		},
		"serverServiceId": &graphql.Field{
			Type:        graphql.String,
			Description: "Server plan id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.ServerServiceID, nil
				}
				return nil, nil
			},
		},
		"siteCode": &graphql.Field{
			Type:        graphql.Int,
			Description: "siteCode ",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.SiteCode, nil
				}
				return nil, nil
			},
		},
		"siteId": &graphql.Field{
			Type:        graphql.Int,
			Description: "siteId ",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Name for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesData); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
	},
})

// AllSitesListData : AllSitesListData Structure
type AllSitesListData struct {
	Status       int64          `json:"status"`
	AllSitesList []AllSitesData `json:"outdata"`
}

// AllSitesListType : AllSitesList GraphQL Schema
var AllSitesListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AllSitesListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(AllSitesType),
			Description: "AllSites list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AllSitesListData); ok {
					return CurData.AllSitesList, nil
				}
				return nil, nil
			},
		},
	},
})
