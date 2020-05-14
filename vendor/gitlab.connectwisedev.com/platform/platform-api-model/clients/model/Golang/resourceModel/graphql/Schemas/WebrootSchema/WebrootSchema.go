package WebrootSchema

import "github.com/graphql-go/graphql"

//WrSiteData : WrSiteData structure
type WrSiteData struct {
	WebroorSitename string `json:"company_name"`
	WebrootSitecode string `json:"keycode"`
}

//WRSiteType : WRSiteType Schema Definition
var WRSiteType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WRSiteDetails",
	Fields: graphql.Fields{
		"company_name": &graphql.Field{
			Type:        graphql.String,
			Description: "Webroor Sitename",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrSiteData); ok {
					return CurData.WebroorSitename, nil
				}
				return nil, nil
			},
		},
		"keycode": &graphql.Field{
			Type:        graphql.String,
			Description: "Webroot Sitecode",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrSiteData); ok {
					return CurData.WebrootSitecode, nil
				}
				return nil, nil
			},
		},
	},
})

//WrSiteListData : Structure for WrSiteData Rest Output
type WrSiteListData struct {
	WrSiteList []WrSiteData `json:"webrootSites"`
}

//WrSiteListType : Object of WrSiteListData
var WrSiteListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WrSiteList",
	Fields: graphql.Fields{
		"webrootSites": &graphql.Field{
			Type:        graphql.NewList(WRSiteType),
			Description: "webroot Site list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrSiteListData); ok {
					return CurData.WrSiteList, nil
				}
				return nil, nil
			},
		},
	},
})

//WrSitesMap : map of WrSiteData
type WrSitesMap struct {
	Licenses map[string]WrSiteData `json:"child_licenses"`
}

type WrMutationData struct {
	PartnerId       string `json:"PartnerId"`
	IsNewSite       string `json:"IsNewSite"`
	SiteId          string `json:"SiteId"`
	ParentKeyCode   string `json:"ParentKeyCode"`
	IsServer        string `json:"IsServer"`
	SiteSvr         string `json:"SiteSvr"`
	SiteKeyCodeSvr  string `json:"SiteKeyCodeSvr"`
	IsDesktop       string `json:"IsDesktop"`
	SiteNameDskp    string `json:"SiteNameDskp"`
	SiteKeyCodeDskp string `json:"SiteKeyCodeDskp"`
	IsNOC           string `json:"IsNOC"`
	CreatedBy       string `json:"CreatedBy"`
	UsrEmailID      string `json:"UsrEmailID"`
}
