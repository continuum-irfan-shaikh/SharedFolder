package WebrootSiteMappingSchema

import (
	"github.com/graphql-go/graphql"
)

//WebrootSiteMappingData : WebrootSiteMappingData structure
type WebrootSiteMappingData struct {
	IsDesktop     		 bool   `json:"IsDesktop"`
	IsServer     		 bool   `json:"IsServer"`
	DesktopSiteKeyCode   string `json:"WebRt_SiteKeyCodeDskp"`
	ServerSiteKeyCode    string `json:"WebRt_SiteKeyCodeSvr"`
	DesktopSiteName      string `json:"WebRt_SiteNameDskp"`
	ServerSiteName       string `json:"WebRt_SiteNameSvr"`
}

//WebrootSiteMappingType : WebrootSiteMappingType Schema Definition
var WebrootSiteMappingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootSiteMapping",
	Fields: graphql.Fields{
	    "IsDesktop": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "is desktop",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingData); ok {
					return CurData.IsDesktop, nil
				}
				return nil, nil
			},
		},
	    "IsServer": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Is Server",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingData); ok {
					return CurData.IsServer, nil
				}
				return nil, nil
			},
		},
	    "DesktopSiteKeyCode": &graphql.Field{
			Type:        graphql.String,
			Description: "webroot desktop site key code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingData); ok {
					return CurData.DesktopSiteKeyCode, nil
				}
				return nil, nil
			},
		},
	    "ServerSiteKeyCode": &graphql.Field{
			Type:        graphql.String,
			Description: "webroot server site key code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingData); ok {
					return CurData.ServerSiteKeyCode, nil
				}
				return nil, nil
			},
		},
	    "DesktopSiteName": &graphql.Field{
		Type:        graphql.String,
			Description: "webroot desktop site name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingData); ok {
					return CurData.DesktopSiteName, nil
				}
				return nil, nil
			},
		},
	    "ServerSiteName": &graphql.Field{
			Type:        graphql.String,
			Description: "webroot server site name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingData); ok {
					return CurData.ServerSiteName, nil
				}
				return nil, nil
			},
		},
	},
})

// WebrootSiteMappingListData : WebrootSiteMappingListData data
type WebrootSiteMappingListData struct {
	Status      int64         `json:"status"`
	WebrootSiteMappingList []WebrootSiteMappingData `json:"outdata"`
}

// WebrootSiteMappingListType : WebrootSiteMappingListType GraphQL Query Schema
var WebrootSiteMappingListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootSiteMappingListType",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        graphql.NewList(WebrootSiteMappingType),
			Description: "Webroot Site Mapping",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootSiteMappingListData); ok {
					return CurData.WebrootSiteMappingList, nil
				}
				return nil, nil
			},
		},
	},
})
