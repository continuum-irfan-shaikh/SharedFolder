package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetSystemData : AssetSystemData Structure
type AssetSystemData struct {
	Product             string
	Model               string
	TimeZone            string
	TimeZoneDescription string
	SerialNumber        string
	SystemName          string
	Category            string
}

//AssetSystemType : AssetSystem GraphQL Schema
var AssetSystemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetSystem",
	Fields: graphql.Fields{
		"product": &graphql.Field{
			Type:        graphql.String,
			Description: "System product information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSystemData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},

		"model": &graphql.Field{
			Type:        graphql.String,
			Description: "System model information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSystemData); ok {
					return CurData.Model, nil
				}
				return nil, nil
			},
		},

		"timeZone": &graphql.Field{
			Type:        graphql.String,
			Description: "System time zone information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSystemData); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},

		"timeZoneDescription": &graphql.Field{
			Type:        graphql.String,
			Description: "System time zone description",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSystemData); ok {
					return CurData.TimeZoneDescription, nil
				}
				return nil, nil
			},
		},

		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "System serial number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSystemData); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},

		"systemName": &graphql.Field{
			Type:        graphql.String,
			Description: "System serial name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSystemData); ok {
					return CurData.SystemName, nil
				}
				return nil, nil
			},
		},

		"category": &graphql.Field{
			Type:        graphql.String,
			Description: "System category like server, switch etc.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSystemData); ok {
					return CurData.Category, nil
				}
				return nil, nil
			},
		},
	},
})
