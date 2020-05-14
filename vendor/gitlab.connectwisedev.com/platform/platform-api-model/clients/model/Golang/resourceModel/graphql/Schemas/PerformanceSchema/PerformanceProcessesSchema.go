package PerformanceSchema

import (
	"github.com/graphql-go/graphql"
)

//PerformanceProcessesData : PerformanceProcessesData struct
type PerformanceProcessesData struct {
	CreateTimeUTC string               `json:"createTimeUTC"`
	CreatedBy     string               `json:"createdBy"`
	Name          string               `json:"createdBy"`
	Type          string               `json:"type"`
	Processes     []PerformanceProcessData `json:"processes"`
	ProcessesCnt  int64 		      `json:"processesCnt"`
}

//PerformanceProcessesType : PerformanceProcessesType graphql object
var PerformanceProcessesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PerformanceProcessesData",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC", 
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessesData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessesData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessesData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessesData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"processes": &graphql.Field{
			Description: "Performance of process",
			Type:        graphql.NewList(performanceProcessType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessesData); ok {
					return CurrData.Processes, nil
				}
				return nil, nil
			},
		},

		"processesCnt": &graphql.Field{
			Description: "Performance of process Count",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessesData); ok {
					return CurrData.ProcessesCnt, nil
				}
				return nil, nil
			},
		},
	},
})
