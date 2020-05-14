package PerformanceSchema

import (
	"time"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//CreateTimeUTCCol : CreateTimeUTC column name
var CreateTimeUTCCol = "CREATETIMEUTC"

//PerformanceCollectionData : PerformanceCollection struct
type PerformanceCollectionData struct {
	CreateTimeUTC time.Time `json:"createTimeUTC"`
	CreatedBy     string    `json:"createdBy"`
	//	Index         int32                `json:"index"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Processors []ProcessorsData           `json:"processors"`
	Memory     []PerformanceMemoryData    `json:"memory"`
	Storages   []PerformanceStoragesData  `json:"storages"`
	Network    []PerformanceNetworkData   `json:"network"`
	Processes  []PerformanceProcessesData `json:"processes"`
}

//PerformanceCollectionType : PerformanceCollectionType graphql object
var PerformanceCollectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PerformanceCollectionData",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC", 
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"processors": &graphql.Field{
			Description: "processor metrics",
			Type:        graphql.NewList(ProcessorsType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.Processors, nil
				}
				return nil, nil
			},
		},

		"memory": &graphql.Field{
			Description: "memory metrics  ",
			Type:        graphql.NewList(PerformanceMemoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.Memory, nil
				}
				return nil, nil
			},
		},

		"storages": &graphql.Field{
			Description: "network metrics",
			Type:        graphql.NewList(PerformanceStoragesType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.Storages, nil
				}
				return nil, nil
			},
		},

		"network": &graphql.Field{
			Description: "memory metrics  ",
			Type:        graphql.NewList(PerformanceNetworkType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.Network, nil
				}
				return nil, nil
			},
		},

		"processes": &graphql.Field{
			Description: "processesmemory metrics  ",
			Type:        graphql.NewList(PerformanceProcessesType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceCollectionData); ok {
					return CurrData.Processes, nil
				}
				return nil, nil
			},
		},
	},
})
