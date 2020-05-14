package FreeDiskSpaceSchema

import (
	"github.com/graphql-go/graphql"
)

//FreeDiskSpaceData : FreeDiskSpaceData Structure
type FreeDiskSpaceData struct {
	RegID          string `json:"regId"`
	PartnerID      string `json:"partnerId"`
	SpaceAvailable bool   `json:"spaceAvailable"`
}

//FreeDiskSpaceType : FreeDiskSpace GraphQL Schema
var FreeDiskSpaceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FreeDiskSpace",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(FreeDiskSpaceData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific partner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(FreeDiskSpaceData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"spaceAvailable": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Space Available",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(FreeDiskSpaceData); ok {
					return CurData.SpaceAvailable, nil
				}
				return nil, nil
			},
		},
	},
})

//FreeDiskSpaceListData : FreeDiskSpaceListData Structure
type FreeDiskSpaceListData struct {
	FreeDiskSpace []FreeDiskSpaceData `json:"freeDiskSpaceList"`
	TotalCount    int64               `json:"totalCount"`
}

//FreeDiskSpaceListType : FreeDiskSpaceList GraphQL Schema
var FreeDiskSpaceListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FreeDiskSpace",
	Fields: graphql.Fields{
		"freeDiskSpaceList": &graphql.Field{
			Type:        graphql.NewList(FreeDiskSpaceType),
			Description: "FreeDiskSpace list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(FreeDiskSpaceListData); ok {
					return CurData.FreeDiskSpace, nil
				}
				return nil, nil
			},
		},

		"totalCount": &graphql.Field{
			Type:        graphql.String,
			Description: "totalCount of FreeDiskSpace list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(FreeDiskSpaceListData); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
