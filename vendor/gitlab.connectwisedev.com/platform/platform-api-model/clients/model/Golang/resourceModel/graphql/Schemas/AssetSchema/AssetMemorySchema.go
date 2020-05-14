package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetMemoryData : AssetMemoryData Structure
type AssetMemoryData struct {
	Manufacturer	string
	SerialNumber	string
	SizeBytes     	int64
}

//AssetMemoryType : AssetMemory GraphQL Schema
var AssetMemoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetMemory",
	Fields: graphql.Fields{
		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "Physical Memory Manufacturer",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMemoryData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "Physical Memory SerialNumber",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMemoryData); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},

		"sizeBytes": &graphql.Field{
			Type:        graphql.String,
			Description: "Size of the Physical Memory in Bytes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMemoryData); ok {
					return CurData.SizeBytes, nil
				}
				return nil, nil
			},
		},
	},
})
