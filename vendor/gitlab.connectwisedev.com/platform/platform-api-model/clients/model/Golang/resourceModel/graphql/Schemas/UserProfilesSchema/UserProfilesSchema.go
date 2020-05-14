package UserProfilesSchema

import (
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

// CategoryData stores data for profile categories
type CategoryData struct {
	Label      string `json:"label"`
	Name       string `json:"name"`
	Importance int    `json:"importance"`
	Active     bool   `json:"active"`
	ScriptID   string `json:"scriptID"`
	Type       string `json:"type"`
	ShortName  string `json:"short_name"`
}

//CategoryType : CategoryType GraphQL Schema
var CategoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Category",
	Fields: graphql.Fields{
		"label": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Label, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"importance": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Importance, nil
				}
				return nil, nil
			},
		},

		"active": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Active, nil
				}
				return nil, nil
			},
		},

		"scriptID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.ScriptID, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},

		"shortName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryData); ok {
					return CurData.ShortName, nil
				}
				return nil, nil
			},
		},
	},
})

// UserProfile stores data
type UserProfile struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Sites       []string       `json:"sites"`
	AllSites    bool           `json:"all_sites"`
	Default     bool           `json:"default"`
	Active      bool           `json:"active"`
	AlertActive bool           `json:"alert_active"`
	Categories  []CategoryData `json:"categories"`
	Threshold   int            `json:"threshold"`
	Type        string         `json:"type"`
}

//UserProfileType : UserProfile GraphQL Schema
var UserProfileType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProfileType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"sites": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.Sites, nil
				}
				return nil, nil
			},
		},
		"all_sites": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.AllSites, nil
				}
				return nil, nil
			},
		},
		"default": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.Default, nil
				}
				return nil, nil
			},
		},
		"active": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.Active, nil
				}
				return nil, nil
			},
		},
		"alertActive": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.AlertActive, nil
				}
				return nil, nil
			},
		},

		"categories": &graphql.Field{
			Type: graphql.NewList(CategoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},
		"threshold": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.Threshold, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfile); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
	},
})

//UserProfileList : UserProfile List struct
type UserProfileList struct {
	UserProfileData []UserProfile `json:"userProfileData"`
}

//UserProfileListType : ProfileProtectList GraphQL Schema
var UserProfileListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserProfileList",
	Fields: graphql.Fields{
		"userProfileList": &graphql.Field{
			Type: graphql.NewList(UserProfileType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserProfileList); ok {
					return CurData.UserProfileData, nil
				}
				return nil, nil
			},
		},
	},
})
