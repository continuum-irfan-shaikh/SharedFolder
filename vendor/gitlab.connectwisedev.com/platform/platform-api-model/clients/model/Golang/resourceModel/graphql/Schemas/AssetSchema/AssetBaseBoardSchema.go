package AssetSchema

import (
	"time"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//AssetBaseBoardData : AssetBaseBoardData Structure
type AssetBaseBoardData struct {
	Product      string
	Manufacturer string
	Model        string
	SerialNumber string
	Name         string
	Version	     string
	InstallDate  time.Time
}

//AssetBaseBoardType : AssetBaseBoard GraphQL Schema
var AssetBaseBoardType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetBaseBoard",
	Fields: graphql.Fields{
		"product": &graphql.Field{
			Type:        graphql.String,
			Description: "BaseBoard product name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBaseBoardData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},

		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "BaseBoard product manufacturer",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBaseBoardData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"model": &graphql.Field{
			Type:        graphql.String,
			Description: "BaseBoard product model",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBaseBoardData); ok {
					return CurData.Model, nil
				}
				return nil, nil
			},
		},

		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "BaseBoard serial number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBaseBoardData); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "BaseBoard name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBaseBoardData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "BaseBoard version",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBaseBoardData); ok {
					return CurData.Version, nil
				}
				return nil, nil
			},
		},

		"installDate": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "BaseBoard installation date",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetBaseBoardData); ok {
					return CurData.InstallDate, nil
				}
				return nil, nil
			},
		},
	},
})
