package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetMonitorData : AssetMonitorData Structure
type AssetMonitorData struct {
	DeviceID     string
	Name         string
	Manufacturer string
	ScreenHeight int64
	ScreenWidth  int64
	Resolution   string
}

//AssetMonitorType : AssetMonitor GraphQL Schema
var AssetMonitorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetMonitor",
	Fields: graphql.Fields{
		"deviceID": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMonitorData); ok {
					return CurData.DeviceID, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMonitorData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "Manufacturer of the monitor",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMonitorData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"screenHeight": &graphql.Field{
			Type:        graphql.String,
			Description: "The logical height of the display in screen coordinates",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMonitorData); ok {
					return CurData.ScreenHeight, nil
				}
				return nil, nil
			},
		},

		"screenWidth": &graphql.Field{
			Type:        graphql.String,
			Description: "The logical width of the display in screen coordinates",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMonitorData); ok {
					return CurData.ScreenWidth, nil
				}
				return nil, nil
			},
		},

		"resolution": &graphql.Field{
			Type:        graphql.String,
			Description: "The logical width of the display in screen coordinates",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMonitorData); ok {
					return CurData.Resolution, nil
				}
				return nil, nil
			},
		},
	},
})
