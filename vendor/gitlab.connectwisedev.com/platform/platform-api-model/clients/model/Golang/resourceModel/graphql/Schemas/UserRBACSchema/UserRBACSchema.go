package UserRBACSchema

import (
	"github.com/graphql-go/graphql"
)

//UserRBACData : UserRBACData Structure
type UserRBACData struct {
	MenuIDs         string `json:"menuIDs"`
	ReadonlyMenuIDs string `json:"readonlyMenuIDs"`
}

//UserRBACDataType : UserRBACData GraphQL Schema
var UserRBACDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserRBACData",
	Fields: graphql.Fields{
		"MenuIDs": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserRBACData); ok {
					return CurData.MenuIDs, nil
				}
				return nil, nil
			},
		},

		"ReadonlyMenuIDs": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserRBACData); ok {
					return CurData.ReadonlyMenuIDs, nil
				}
				return nil, nil
			},
		},
	},
})
