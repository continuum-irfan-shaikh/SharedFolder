package SiteOsPatchingSchema

import (
	"github.com/graphql-go/graphql"
)

//OsPatchingData : OsPatchingData Structure
type OsPatchingData struct {
	RegType              string `json:"regType"`
	CurrentCount         int    `json:"current"`
	OldCount             int    `json:"notCurrent"`
	RestartNeededCount   int    `json:"rebootNeeded"`
	ExcludedCount        int    `json:"excluded"`
	NotSeenRecentlyCount int    `json:"notSeenRecently"`
}

//OsPatchingType : OsPatching GraphQL Schema
var OsPatchingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OsPatchingDetails",
	Fields: graphql.Fields{
		"current": &graphql.Field{
			Type:        graphql.String,
			Description: "current",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(OsPatchingData); ok {
					return CurData.CurrentCount, nil
				}
				return nil, nil
			},
		},

		"notCurrent": &graphql.Field{
			Type:        graphql.String,
			Description: "notCurrent",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(OsPatchingData); ok {
					return CurData.OldCount, nil
				}
				return nil, nil
			},
		},

		"rebootNeeded": &graphql.Field{
			Type:        graphql.String,
			Description: "rebootNeeded",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(OsPatchingData); ok {
					return CurData.RestartNeededCount, nil
				}
				return nil, nil
			},
		},

		"excluded": &graphql.Field{
			Type:        graphql.String,
			Description: "excluded",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(OsPatchingData); ok {
					return CurData.ExcludedCount, nil
				}
				return nil, nil
			},
		},

		"notSeenRecently": &graphql.Field{
			Type:        graphql.String,
			Description: "notSeenRecently",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(OsPatchingData); ok {
					return CurData.NotSeenRecentlyCount, nil
				}
				return nil, nil
			},
		},

		"regType": &graphql.Field{
			Type:        graphql.String,
			Description: "regType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(OsPatchingData); ok {
					return CurData.RegType, nil
				}
				return nil, nil
			},
		},
	},
})

//SiteOsPatchingData : SiteOsPatchingData Structure
type SiteOsPatchingData struct {
	OsPatchingList []OsPatchingData `json:"osPatchingList"`
}

//SiteOsPatchingType : SiteOsPatchingType GraphQL Schema
var SiteOsPatchingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SiteOsPatching",
	Fields: graphql.Fields{
		"osPatchingList": &graphql.Field{
			Type:        graphql.NewList(OsPatchingType),
			Description: "OsPatchingList",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteOsPatchingData); ok {
					return CurData.OsPatchingList, nil
				}
				return nil, nil
			},
		},
	},
})
