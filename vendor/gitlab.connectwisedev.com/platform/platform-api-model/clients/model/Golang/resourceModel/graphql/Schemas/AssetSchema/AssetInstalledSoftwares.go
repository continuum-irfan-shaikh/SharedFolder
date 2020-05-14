package AssetSchema

import (
	"time"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//AssetInstalledSoftwares : AssetInstalledSoftwares Structure
type AssetInstalledSoftwares struct {
	Name		string
      	Version		string
      	Publisher 	string
      	InstallDate	time.Time
}

//AssetInstalledSoftwaresType : AssetInstalledSoftwares GraphQL Schema
var AssetInstalledSoftwaresType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetInstalledSoftwares",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the installed software",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetInstalledSoftwares); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "Version of the installed software",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetInstalledSoftwares); ok {
					return CurData.Version, nil
				}
				return nil, nil
			},
		},

		"publisher": &graphql.Field{
			Type:        graphql.String,
			Description: "Publisher of the installed software",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetInstalledSoftwares); ok {
					return CurData.Publisher, nil
				}
				return nil, nil
			},
		},

		"installDate": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "Date on which software was installed",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetInstalledSoftwares); ok {
					return CurData.InstallDate, nil
				}
				return nil, nil
			},
		},
	},
})
