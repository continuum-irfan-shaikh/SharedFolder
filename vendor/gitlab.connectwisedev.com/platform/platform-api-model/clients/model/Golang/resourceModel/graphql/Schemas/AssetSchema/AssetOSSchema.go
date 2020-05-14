package AssetSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//AssetOsData : AssetOsData Structure
type AssetOsData struct {
	Product      string
	Manufacturer string
	OsLanguage   string
	SerialNumber string
	Version      string
	InstallDate  time.Time
	Type         string
	Arch         string
	ServicePack  string
	BuildNumber  string
	ReleaseID    string
}

//AssetOsType : AssetOs GraphQL Schema
var AssetOsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetOs",
	Fields: graphql.Fields{
		"product": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system product name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},

		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system manufacturer name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"osLanguage": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system language",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.OsLanguage, nil
				}
				return nil, nil
			},
		},

		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system serial number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system version",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.Version, nil
				}
				return nil, nil
			},
		},

		"installDate": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "Operating system installation date",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.InstallDate, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},

		"arch": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system architecture",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.Arch, nil
				}
				return nil, nil
			},
		},

		"servicePack": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system service pack",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.ServicePack, nil
				}
				return nil, nil
			},
		},

		"buildNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system build number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.BuildNumber, nil
				}
				return nil, nil
			},
		},
		"releaseID": &graphql.Field{
			Type:        graphql.String,
			Description: "Operating system release ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetOsData); ok {
					return CurData.ReleaseID, nil
				}
				return nil, nil
			},
		},
	},
})
