package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetProcessorData : AssetProcessorData Structure
type AssetProcessorData struct {
	Product       string
	NumberOfCores int64
	ClockSpeedMhz float64
	Family        int64
	Manufacturer  string
	ProcessorType string
	SerialNumber  string
	Level         int64
}

//AssetProcessorType : AssetProcessor GraphQL Schema
var AssetProcessorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetProcessor",
	Fields: graphql.Fields{
		"product": &graphql.Field{
			Type:        graphql.String,
			Description: "Processor product name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},

		"numberOfCores": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of cores",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.NumberOfCores, nil
				}
				return nil, nil
			},
		},

		"clockSpeedMhz": &graphql.Field{
			Type:        graphql.Float,
			Description: "Processor clock speed in MegaHertz",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.ClockSpeedMhz, nil
				}
				return nil, nil
			},
		},

		"family": &graphql.Field{
			Type:        graphql.String,
			Description: "Processor family",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.Family, nil
				}
				return nil, nil
			},
		},

		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "Processor manufacturer name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"processorType": &graphql.Field{
			Type:        graphql.String,
			Description: "Processor type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.ProcessorType, nil
				}
				return nil, nil
			},
		},

		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "Processor serial number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},

		"level": &graphql.Field{
			Type:        graphql.String,
			Description: "Processor level",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetProcessorData); ok {
					return CurData.Level, nil
				}
				return nil, nil
			},
		},
	},
})
