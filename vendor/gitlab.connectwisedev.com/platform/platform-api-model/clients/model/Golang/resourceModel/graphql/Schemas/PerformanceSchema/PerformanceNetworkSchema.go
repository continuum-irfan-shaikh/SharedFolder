package PerformanceSchema

import (
	"time"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//PerformanceNetworkInterfaceData : PerformanceNetworkInterface struct
type PerformanceNetworkInterfaceData struct {
	CreateTimeUTC    	time.Time `json:"createTimeUTC"`
	CreatedBy        	string    `json:"createdBy"`
	Index            	int32     `json:"index"`
	Name             	string    `json:"name"`
	Type             	string    `json:"type"`
	TotalBytesPerSec    	int64     `json:"totalBytesPerSec"`
	ReceivedBytes    	int64     `json:"receivedBytes"`
	ReceivedBytesPerSec    	int64     `json:"receivedBytesPerSec"`
	TransmittedBytes 	int64     `json:"transmittedBytes"`
	TransmittedBytesPerSec  int64     `json:"transmittedBytesPerSec"`
	TXQueueLength    	int64     `json:"txQueueLength"`
	RXQueueLength    	int64     `json:"rxQueueLength"`
}

//PerformanceNetworkData : PerformanceNetwork struct
type PerformanceNetworkData struct {
	CreateTimeUTC time.Time                     `json:"createTimeUTC"`
	CreatedBy     string                        `json:"createdBy"`
	Index         int32                         `json:"index"`
	Name          string                        `json:"name"`
	Type          string                        `json:"type"`
	Interface     []PerformanceNetworkInterfaceData `json:"interface"`
}

//PerformanceNetworkInterfaceType : PerformanceNetworkInterface graphql object
var PerformanceNetworkInterfaceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "performanceNetworkInterface",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"totalBytesPerSec": &graphql.Field{
			Description: "Total bytes per second",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.TotalBytesPerSec, nil
				}
				return nil, nil
			},
		},

		"receivedBytes": &graphql.Field{
			Description: "ReceivedBytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.ReceivedBytes, nil
				}
				return nil, nil
			},
		},

		"receivedBytesPerSec": &graphql.Field{
			Description: "Received bytes per second",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.ReceivedBytesPerSec, nil
				}
				return nil, nil
			},
		},

		"transmittedBytes": &graphql.Field{
			Description: "TransmittedBytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.TransmittedBytes, nil
				}
				return nil, nil
			},
		},

		"transmittedBytesPerSec": &graphql.Field{
			Description: "Transmitted bytes per second",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.TransmittedBytesPerSec, nil
				}
				return nil, nil
			},
		},

		"txQueueLength": &graphql.Field{
			Description: "TXQueueLength",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.TXQueueLength, nil
				}
				return nil, nil
			},
		},

		"rxQueueLength": &graphql.Field{
			Description: "RXQueueLength",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkInterfaceData); ok {
					return CurrData.RXQueueLength, nil
				}
				return nil, nil
			},
		},
	},
})

//PerformanceNetworkType : PerformanceNetwork graphql object
var PerformanceNetworkType = graphql.NewObject(graphql.ObjectConfig{
	Name: "performanceNetwork",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"interface": &graphql.Field{
			Description: "Network Interface",
			Type:        graphql.NewList(PerformanceNetworkInterfaceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceNetworkData); ok {
					return CurrData.Interface, nil
				}
				return nil, nil
			},
		},
	},
})
