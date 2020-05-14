package PerformanceSchema

import (
	"github.com/graphql-go/graphql"
	"time"
)

//PerformanceProcessData : PerformanceProcess struct
type PerformanceProcessData struct {
	CreateTimeUTC   string  `json:"createTimeUTC"`
	CreatedBy       string  `json:"createdBy"`
	Index           int32   `json:"index"`
	Name            string  `json:"name"`
	ProcessID       int32   `json:"processid"`
	Type            string  `json:"type"`
	PercentCPUUsage float64 `json:"percentCPUUsage"`
	HandleCount     int64   `json:"handleCount"`
	ThreadCount     int64   `json:"threadCount"`
	PrivateBytes    int64   `json:"privateBytes"`
	DiskReadBytes	int64   `json:"diskReadBytes"`
	DiskWriteBytes	int64   `json:"diskWriteBytes"`
	NetSendBytes	int64   `json:"netSendBytes"`
	NetReceiveBytes	int64   `json:"netReceiveBytes"`
	FetchTimeUTC	string  `json:"fetchTimeUTC"`
}

//performanceProcessType : performanceProcess graphql object
var performanceProcessType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PerformanceProcessData",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC", 
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"processId": &graphql.Field{
			Description: "ProcessID",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.ProcessID, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"percentCPUUsage": &graphql.Field{
			Description: "CPU usage in percentage",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.PercentCPUUsage, nil
				}
				return nil, nil
			},
		},

		"handleCount": &graphql.Field{
			Description: "Number of handles",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.HandleCount, nil
				}
				return nil, nil
			},
		},

		"threadCount": &graphql.Field{
			Description: "Number of threads",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.ThreadCount, nil
				}
				return nil, nil
			},
		},

		"privateBytes": &graphql.Field{
			Description: "Private bytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.PrivateBytes, nil
				}
				return nil, nil
			},
		},

		"diskReadBytes": &graphql.Field{
			Description: "Number of bytes read from the disk",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.DiskReadBytes, nil
				}
				return nil, nil
			},
		},

		"diskWriteBytes": &graphql.Field{
			Description: "Number of bytes written to the disk",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.DiskWriteBytes, nil
				}
				return nil, nil
			},
		},

		"netSendBytes": &graphql.Field{
			Description: "Number of bytes sent over the network",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.NetSendBytes, nil
				}
				return nil, nil
			},
		},

		"netReceiveBytes": &graphql.Field{
			Description: "Number of bytes received over the network",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceProcessData); ok {
					return CurrData.NetReceiveBytes, nil
				}
				return nil, nil
			},
		},
		
		"fetchTimeUTC": &graphql.Field{
			Description: "To fetch data request UTC time",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return time.Now().UTC(), nil
			},
		},
	},
})
