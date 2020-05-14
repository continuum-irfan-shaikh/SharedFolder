package WebrootDNSSchema

import (
	"github.com/graphql-go/graphql"
)

type WrDnsSiteData struct {
	Vendor       string `json:"vendor"`
	NodeId       string `json:"nodeId"`
	Level        int    `json:"level"`
	ParentNodeId string `json:"parentNodeId"`
	Attributes   string `json:"attributes"`
}

var WRDnsSiteType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WRDnsSiteDetails",
	Fields: graphql.Fields{
		"vendor": &graphql.Field{
			Type:        graphql.String,
			Description: "vendor",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDnsSiteData); ok {
					return CurData.Vendor, nil
				}
				return nil, nil
			},
		},
		"nodeId": &graphql.Field{
			Type:        graphql.String,
			Description: "nodeId ID for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDnsSiteData); ok {
					return CurData.NodeId, nil
				}
				return nil, nil
			},
		},
		"level": &graphql.Field{
			Type:        graphql.String,
			Description: "level ID for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDnsSiteData); ok {
					return CurData.Level, nil
				}
				return nil, nil
			},
		},
		"parentNodeId": &graphql.Field{
			Type:        graphql.String,
			Description: "parentNodeId for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDnsSiteData); ok {
					return CurData.ParentNodeId, nil
				}
				return nil, nil
			},
		},
		"attributes": &graphql.Field{
			Type:        graphql.String,
			Description: "attributes Desc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDnsSiteData); ok {
					return CurData.Attributes, nil
				}
				return nil, nil
			},
		},
	},
})

type WrDnsSiteListData struct {
	WrDnsSiteList []WrDnsSiteData `json:"dnsSites"`
}

var WrDnsSiteListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WrDnsSiteList",
	Fields: graphql.Fields{
		"dnsSites": &graphql.Field{
			Type:        graphql.NewList(WRDnsSiteType),
			Description: "wrDnsSiteL list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDnsSiteListData); ok {
					return CurData.WrDnsSiteList, nil
				}
				return nil, nil
			},
		},
	},
})

type gsmtype struct {
	Gsm string `json:"gsm-key"`
}

//WrGsmObject : Object of gsmtype
var WrGsmObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "WrGsmObject",
	Fields: graphql.Fields{
		"gsm": &graphql.Field{
			Type:        graphql.String,
			Description: "GSM",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(gsmtype); ok {
					return CurData.Gsm, nil
				}
				return nil, nil
			},
		},
	},
})

//WrDNSParentKeyDataType : Type for Webroot DNS ParentKeyCode Data
type WrDNSParentKeyDataType struct {
	ID         string  `json:"_id"`
	Attributes gsmtype `json:"attributes"`
}

//WrDNSParentKeyDataObject : Object of WrDNSParentKeyDataType
var WrDNSParentKeyDataObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "WrDNSParentKeyDataObject",
	Fields: graphql.Fields{
		"_id": &graphql.Field{
			Type:        graphql.String,
			Description: "ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDNSParentKeyDataType); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"attributes": &graphql.Field{
			Type:        WrGsmObject,
			Description: "Attributes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDNSParentKeyDataType); ok {
					return CurData.Attributes, nil
				}
				return nil, nil
			},
		},
	},
})

//WrDNSParentKeyCodeType : Type for Webroot DNS ParentKeyCode root
type WrDNSParentKeyCodeType struct {
	Webrootparentkeycode []WrDNSParentKeyDataType `json:"wrdnsparentkeycode"`
}

//WrDNSParentKeyCodeObject : Object for WrDNSParentKeyCodeType
var WrDNSParentKeyCodeObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "WrDNSParentKeyCodeObject",
	Fields: graphql.Fields{
		"webrootparentkeycode": &graphql.Field{
			Type:        graphql.NewList(WrDNSParentKeyDataObject),
			Description: "WrDNSParentKeyCodeType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WrDNSParentKeyCodeType); ok {
					return CurData.Webrootparentkeycode, nil
				}
				return nil, nil
			},
		},
	},
})

type WrDnsMutationData struct {
	WRSiteKeyCode string `json:"_id"`
}

type WrDnsProductData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
