package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetBiosData : AssetBiosData Structure
type AssetBiosData struct {
	Product       string
	Manufacturer  string
	Version       string
	SerialNumber  string
	SmbiosVersion string
}

//AssetBiosType : AssetBios GraphQL Schema
var AssetBiosType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetBios",
	Fields: graphql.Fields{
		"product": &graphql.Field{
			Type:        graphql.String,
			Description: "BIOS name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBiosData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},

		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "BIOS manufacturer",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBiosData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "BIOS version",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBiosData); ok {
					return CurData.Version, nil
				}
				return nil, nil
			},
		},

		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "BIOS serial number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBiosData); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},

		"smbiosVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "SMBIOS version",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBiosData); ok {
					return CurData.SmbiosVersion, nil
				}
				return nil, nil
			},
		},
	},
})
