package HeartBeatSchema

import (
	"time"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
)

//HeartBeatData : HeartBeatData Structure
type HeartBeatData struct {
	EndpointID           	string 		`json:"EndpointID"`
	DCDateTimeUTC       	time.Time 	`json:"DcDateTimeUTC"`
	AgentDateTimeUTC       	time.Time 	`json:"AgentDateTimeUTC"`
	HeartbeatCounter        int64 		`json:"HeartbeatCounter"`
	Installed 		bool  		`json:"Installed"`
	Availability 		bool  		`json:"Availability"`
}

//HeartBeatType : HeartBeat GraphQL Schema
var HeartBeatType = graphql.NewObject(graphql.ObjectConfig{
	Name: "heartBeat",
	Fields: graphql.Fields{
		"endpointId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HeartBeatData); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},

		"dcDateTimeUTC": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "DcDateTimeUTC",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HeartBeatData); ok {
					return CurData.DCDateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"agentDateTimeUTC": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "AgentDateTimeUTC",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HeartBeatData); ok {
					return CurData.AgentDateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"heartbeatCounter": &graphql.Field{
			Description: "HeartbeatCounter",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(HeartBeatData); ok {
					return CurrData.HeartbeatCounter, nil
				}
				return nil, nil
			},
		},

		"installed": &graphql.Field{
			Description: "Installed",
			Type:        graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(HeartBeatData); ok {
					return CurrData.Installed, nil
				}
				return nil, nil
			},
		},

		"availability": &graphql.Field{
			Description: "Availability",
			Type:        graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(HeartBeatData); ok {
					return CurrData.Availability, nil
				}
				return nil, nil
			},
		},
	},
})

//HeartBeatListData : HeartBeatListData Structure
type HeartBeatListData struct {
	HeartBeat []HeartBeatData `json:"heartBeatList"`
}

//HeartBeatListType : HeartBeatList GraphQL Schema
var HeartBeatListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HeartBeatList",
	Fields: graphql.Fields{
		"heartBeatList": &graphql.Field{
			Type:        graphql.NewList(HeartBeatType),
			Description: "HeartBeat List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HeartBeatListData); ok {
					return CurData.HeartBeat, nil
				}
				return nil, nil
			},
		},
	},
})
