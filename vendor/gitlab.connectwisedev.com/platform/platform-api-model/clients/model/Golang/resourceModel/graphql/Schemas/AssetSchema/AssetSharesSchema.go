package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetSharesData : AssetSharesData Structure
type AssetSharesData struct {
	Name		string
	Caption		string
	Description	string
	Path		string
	Access		string
	UserAccess	[]string
	Type		[]string
}

//AssetSharesType : AssetShares GraphQL Schema
var AssetSharesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetShares",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSharesData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"caption": &graphql.Field{
			Type:        graphql.String,
			Description: "Caption",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSharesData); ok {
					return CurData.Caption, nil
				}
				return nil, nil
			},
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Description",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSharesData); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},

		"path": &graphql.Field{
			Type:        graphql.String,
			Description: "Path",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSharesData); ok {
					return CurData.Path, nil
				}
				return nil, nil
			},
		},

		"access": &graphql.Field{
			Type:        graphql.String,
			Description: "Access",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSharesData); ok {
					return CurData.Access, nil
				}
				return nil, nil
			},
		},

		"userAccess": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "UserAccess",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSharesData); ok {
					return CurData.UserAccess, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "Type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetSharesData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
	},
})
