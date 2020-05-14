package SitePatchPolicySchema

import (
	"github.com/graphql-go/graphql"
)

// SitePatchPolicyData : SitePatchPolicyData Structure
type SitePatchPolicyData struct {
	DesktopPolicyID   int64  `json:"desktopPolicyID"`
	DesktopPolicyname string `json:"desktopPolicyname"`
	ServerPolicyID    int64  `json:"serverPolicyID"`
	ServerPolicyname  string `json:"serverPolicyname"`
}

// SitePatchPolicyType : SitePatchPolicy GraphQL Schema
var SitePatchPolicyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitePatchPolicy",
	Fields: graphql.Fields{
		"desktopPolicyID": &graphql.Field{
			Type:        graphql.Int,
			Description: "Desktop policyid for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitePatchPolicyData); ok {
					return CurData.DesktopPolicyID, nil
				}
				return nil, nil
			},
		},

		"desktopPolicyname": &graphql.Field{
			Type:        graphql.String,
			Description: "Desktop policyname for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitePatchPolicyData); ok {
					return CurData.DesktopPolicyname, nil
				}
				return nil, nil
			},
		},

		"serverPolicyID": &graphql.Field{
			Type:        graphql.Int,
			Description: "Server policyid for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitePatchPolicyData); ok {
					return CurData.ServerPolicyID, nil
				}
				return nil, nil
			},
		},
		"serverPolicyname": &graphql.Field{
			Type:        graphql.String,
			Description: "Server policyname for site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitePatchPolicyData); ok {
					return CurData.ServerPolicyname, nil
				}
				return nil, nil
			},
		},
	},
})

// SitePatchPolicyListData : SitePatchPolicyListData Structure
type SitePatchPolicyListData struct {
	Status              int64                 `json:"status"`
	SitePatchPolicyList []SitePatchPolicyData `json:"outdata"`
}

// SitePatchPolicyListType : SitePatchPolicyList GraphQL Schema
var SitePatchPolicyListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitePatchPolicyListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitePatchPolicyListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(SitePatchPolicyType),
			Description: "SitePatchPolicy list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitePatchPolicyListData); ok {
					return CurData.SitePatchPolicyList, nil
				}
				return nil, nil
			},
		},
	},
})
