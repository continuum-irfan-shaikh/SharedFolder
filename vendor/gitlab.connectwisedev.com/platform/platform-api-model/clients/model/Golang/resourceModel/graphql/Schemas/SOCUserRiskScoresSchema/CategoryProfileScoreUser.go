package SOCUserRiskScoresSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/ProfileProtectSchema"
	"github.com/graphql-go/graphql"
)

//CategoryProfileUserScoreData : CategoryProfileUserScoreData struct
type CategoryProfileUserScoreData struct {
	Label             string                               `json:"label"`
	Name              string                               `json:"name"`
	CategoryResult    string                               `json:"categoryResult"`
	ExecutionError    string                               `json:"executionError"`
	LastExecutionTime time.Time                            `json:"lastExecutionTime"`
	ExecutionDetails  []ProfileProtectSchema.ExecutionData `json:"executionDetails"`
	RiskStartTime     time.Time                            `json:"riskStartTime"`
}

//CategoryProfileUserScoreType : CategoryProfileUserScoreData Data GraphQL Schema
var CategoryProfileUserScoreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Items",
	Fields: graphql.Fields{
		"label": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileUserScoreData); ok {
					return CurData.Label, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileUserScoreData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"categoryResult": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileUserScoreData); ok {
					return CurData.CategoryResult, nil
				}
				return nil, nil
			},
		},

		"executionError": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileUserScoreData); ok {
					return CurData.ExecutionError, nil
				}
				return nil, nil
			},
		},

		"lastExecutionTime": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileUserScoreData); ok {
					return CurData.LastExecutionTime, nil
				}
				return nil, nil
			},
		},

		"riskStartTime": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileUserScoreData); ok {
					return CurData.RiskStartTime, nil
				}
				return nil, nil
			},
		},

		"executionDetails": &graphql.Field{
			Type: graphql.NewList(ProfileProtectSchema.ExecutionDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ProfileProtectSchema.CategoryProfileScoreData); ok {
					return CurData.ExecutionDetails, nil
				}
				return nil, nil
			},
		},
	},
})
