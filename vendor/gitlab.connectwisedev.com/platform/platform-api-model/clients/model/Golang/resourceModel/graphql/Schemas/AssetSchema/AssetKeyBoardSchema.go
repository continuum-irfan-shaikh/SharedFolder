package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetKeyBoardData : AssetKeyBoardData Structure
type AssetKeyBoardData struct {
	DeviceID	string
	Name		string
	Description	string
}

//AssetKeyBoardType : AssetKeyBoard GraphQL Schema
var AssetKeyBoardType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetKeyBoard",
	Fields: graphql.Fields{
		"deviceID": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetKeyBoardData); ok {
					return CurData.DeviceID, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetKeyBoardData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Detailed textual description of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetKeyBoardData); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
	},
})
