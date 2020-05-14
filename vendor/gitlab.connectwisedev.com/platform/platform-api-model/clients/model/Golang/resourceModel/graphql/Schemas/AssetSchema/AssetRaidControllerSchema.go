package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetRaidControllerData : AssetRaidControllerData Structure
type AssetRaidControllerData struct {
	SoftwareRaid string
	HardwareRaid string
	Vendor       string
}

//AssetRaidControllerType : AssetRaidController GraphQL Schema
var AssetRaidControllerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetRaidController",
	Fields: graphql.Fields{
		"softwareRaid": &graphql.Field{
			Type:        graphql.String,
			Description: "Software RAID information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetRaidControllerData); ok {
					return CurData.SoftwareRaid, nil
				}
				return nil, nil
			},
		},

		"hardwareRaid": &graphql.Field{
			Type:        graphql.String,
			Description: "Hardware RAID information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetRaidControllerData); ok {
					return CurData.HardwareRaid, nil
				}
				return nil, nil
			},
		},

		"vendor": &graphql.Field{
			Type:        graphql.String,
			Description: "RAID vendor information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetRaidControllerData); ok {
					return CurData.Vendor, nil
				}
				return nil, nil
			},
		},
	},
})
