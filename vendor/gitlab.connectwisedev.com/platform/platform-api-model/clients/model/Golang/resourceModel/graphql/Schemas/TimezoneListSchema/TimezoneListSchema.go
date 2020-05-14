package TimezoneListSchema

import (
	"github.com/graphql-go/graphql"
)

//TimezoneData : TimezoneData Structure
type TimezoneData struct {
	ZoneID          string `json:"ZoneID"`
	ZoneDisplayName string `json:"ZoneDisplayName"`
	UTCTimediff     string `json:"UTCTimediff"`
}

//TimezoneType : Timezone GraphQL Schema
var TimezoneType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TimezoneDetails",
	Fields: graphql.Fields{
		"ZoneID": &graphql.Field{
			Type:        graphql.String,
			Description: "TimeZone ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TimezoneData); ok {
					return CurData.ZoneID, nil
				}
				return nil, nil
			},
		},

		"ZoneDisplayName": &graphql.Field{
			Type:        graphql.String,
			Description: "Description",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TimezoneData); ok {
					return CurData.ZoneDisplayName, nil
				}
				return nil, nil
			},
		},

		"UTCTimediff": &graphql.Field{
			Type:        graphql.String,
			Description: "UTCTimediff",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TimezoneData); ok {
					return CurData.UTCTimediff, nil
				}
				return nil, nil
			},
		},
	},
})

//TimezoneListData : TimezoneListData Structure
type TimezoneListData struct {
	Status       int64          `json:"status"`
	TimezoneList []TimezoneData `json:"outdata"`
}

//TimezoneListType : TimezoneListType GraphQL Schema
var TimezoneListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TimezoneList",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TimezoneListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        graphql.NewList(TimezoneType),
			Description: "Site Timezones",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TimezoneListData); ok {
					return CurData.TimezoneList, nil
				}
				return nil, nil
			},
		},
	},
})
