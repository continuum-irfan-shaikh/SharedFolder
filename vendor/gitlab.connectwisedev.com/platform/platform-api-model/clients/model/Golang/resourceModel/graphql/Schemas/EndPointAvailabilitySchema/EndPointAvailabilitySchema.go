package EndPointAvailabilitySchema

import (
	"github.com/graphql-go/graphql"
)

//EndPointAvailabilityData : EndPointAvailabilityData Structure
type EndPointAvailabilityData struct {
	RegID        string `json:"regId"`
	Availability int64  `json:"availability"`
}

//EndPointAvailabilityType : EndPointAvailability GraphQL Schema
var EndPointAvailabilityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EndPointAvailability",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointAvailabilityData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"availability": &graphql.Field{
			Type:        graphql.String,
			Description: "Status of end point",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointAvailabilityData); ok {
					return CurData.Availability, nil
				}
				return nil, nil
			},
		},
	},
})

//EndPointAvailabilityListData : EndPointAvailabilityListData Structure
type EndPointAvailabilityListData struct {
	EndPointAvailability []EndPointAvailabilityData `json:"endPointAvailabilityList"`
	TotalCount           int64                      `json:"totalCount"`
}

//EndPointAvailabilityListType : EndPointAvailabilityList GraphQL Schema
var EndPointAvailabilityListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EndPointAvailabilityList",
	Fields: graphql.Fields{
		"endPointAvailabilityList": &graphql.Field{
			Type:        graphql.NewList(EndPointAvailabilityType),
			Description: "endPointAvailability list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointAvailabilityListData); ok {
					return CurData.EndPointAvailability, nil
				}
				return nil, nil
			},
		},

		"totalCount": &graphql.Field{
			Type:        graphql.String,
			Description: "totalCount of patchStatust list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointAvailabilityListData); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
