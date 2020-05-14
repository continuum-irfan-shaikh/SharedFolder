package WebrootDNSSchema

import (
	"github.com/graphql-go/graphql"
)

//WebrootAttributesList : WebrootAttributesList structure
type WebrootAttributesList struct {
	Deactivated bool   `json:"deactivated"`
	Keycode     string `json:"keycode"`
	Name        string `json:"name"`
	SiteID      string `json:"siteId"`
}

//WebrootAtrributesListType : WebrootAttributesList Schema Definition
var WebrootAtrributesListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootAtrributesListType",
	Fields: graphql.Fields{
		"deactivated": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Deactivated",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootAttributesList); ok {
					return CurData.Deactivated, nil
				}
				return nil, nil
			},
		},
		"keycode": &graphql.Field{
			Type:        graphql.String,
			Description: "Keycode",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootAttributesList); ok {
					return CurData.Keycode, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootAttributesList); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootAttributesList); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
	},
})

//WebrootDNSSiteOptList : Structure for WebrootDNSSiteOpt Rest Output
type WebrootDNSSiteOptList struct {
	ID         string                `json:"_id"`
	Attributes WebrootAttributesList `json:"attributes"`
}

//WebrootDNSSiteOptListType : WebrootDNSSiteOptList Schema Definition
var WebrootDNSSiteOptListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootDNSSiteOptData",
	Fields: graphql.Fields{
		"_id": &graphql.Field{
			Type:        graphql.String,
			Description: "ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootDNSSiteOptList); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"attributes": &graphql.Field{
			Type:        WebrootAtrributesListType,
			Description: "Attributes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootDNSSiteOptList); ok {
					return CurData.Attributes, nil
				}
				return nil, nil
			},
		},
	},
})

//WebrootDNSSiteOptListList : Handles Rest API - Top Hierarchy
type WebrootDNSSiteOptListList struct {
	WebrootDNSSiteOptOutput []WebrootDNSSiteOptList `json:"webrootdnssiteoptlist"`
}

//WebrootDNSSiteOptListTypeList : Object of WebrootDNSSiteOptListList
var WebrootDNSSiteOptListTypeList = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootDNSSiteOptOutput",
	Fields: graphql.Fields{
		"webrootdnssiteoptlist": &graphql.Field{
			Type:        graphql.NewList(WebrootDNSSiteOptListType),
			Description: "list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootDNSSiteOptListList); ok {
					return CurData.WebrootDNSSiteOptOutput, nil
				}
				return nil, nil
			},
		},
	},
})
