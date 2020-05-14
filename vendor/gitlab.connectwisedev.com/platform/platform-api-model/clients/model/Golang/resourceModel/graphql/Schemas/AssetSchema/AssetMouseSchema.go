package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetMouseData : AssetMouseData Structure
type AssetMouseData struct {
	Manufacturer	string
	Name		string
	DeviceID	string
	DeviceInterface	int64
	PointingType	int64
	Buttons		int64
}

//AssetMouseType : AssetMouse GraphQL Schema
var AssetMouseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetMouse",
	Fields: graphql.Fields{
		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "Manufacturer of the mouse",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMouseData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name/caption of the mouse",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMouseData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"deviceID": &graphql.Field{
			Type:        graphql.String,
			Description: "Device ID of the mouse",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMouseData); ok {
					return CurData.DeviceID, nil
				}
				return nil, nil
			},
		},

		"deviceInterface": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of interface used for the pointing device. For ex. PS/2,USB etc.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMouseData); ok {
					return CurData.DeviceInterface, nil
				}
				return nil, nil
			},
		},

		"pointingType": &graphql.Field{
			Type:        graphql.String,
			Description: "Pointing Type of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMouseData); ok {
					return CurData.PointingType, nil
				}
				return nil, nil
			},
		},

		"buttons": &graphql.Field{
			Type:        graphql.String,
			Description: "Number of buttons on the pointing device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetMouseData); ok {
					return CurData.Buttons, nil
				}
				return nil, nil
			},
		},
	},
})
