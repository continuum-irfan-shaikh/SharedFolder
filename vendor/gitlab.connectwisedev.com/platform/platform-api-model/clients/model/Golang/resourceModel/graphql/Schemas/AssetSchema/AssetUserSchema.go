package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetUserData : AssetUserData Structure
type AssetUserData struct {
	UserName		string
	UserType		string
	UserDisabled		bool
	UserLockout		bool
	PasswordRequired	bool
}

//AssetUserType : AssetUser GraphQL Schema
var AssetUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetUser",
	Fields: graphql.Fields{
		"userName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetUserData); ok {
					return CurData.UserName, nil
				}
				return nil, nil
			},
		},

		"userType": &graphql.Field{
			Type:        graphql.String,
			Description: "Whether a guest or Admin account",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetUserData); ok {
					return CurData.UserType, nil
				}
				return nil, nil
			},
		},

		"userDisabled": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Whether Enable/Disabled",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetUserData); ok {
					return CurData.UserDisabled, nil
				}
				return nil, nil
			},
		},

		"userLockout": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Is User account locked out",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetUserData); ok {
					return CurData.UserLockout, nil
				}
				return nil, nil
			},
		},

		"passwordRequired": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Account requires a password or not",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetUserData); ok {
					return CurData.PasswordRequired, nil
				}
				return nil, nil
			},
		},
	},
})
