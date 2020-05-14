package PerformanceSchema

import (
	"time"
	"strings"
	"errors"
	"reflect"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//ProcessorMetricData : ProcessorMetric struct
type ProcessorMetricData struct {
	NumOfProcesses       int64   `json:"numOfProcesses"`
	PercentIOTime        float64 `json:"percentIOTime"`
	PercentIdleTime      float64 `json:"percentIdleTime"`
	PercentInterruptTime float64 `json:"percentInterruptTime"`
	PercentSystemTime    float64 `json:"percentSystemTime"`
	PercentUserTime      float64 `json:"percentUserTime"`
	PercentUtil          float64 `json:"percentUtil"`
	InterruptsPerSec     int64   `json:"interruptsPerSec"`
}

//CoreData : Core struct
type CoreData struct {
	CreateTimeUTC time.Time `json:"createTimeUTC"`
	CreatedBy     string    `json:"createdBy"`
	Index         int32     `json:"index"`
	//	InterruptsPerSec int64
	Name   string          `json:"name"`
	Type   string          `type:"type"`
	Metric ProcessorMetricData `json:"metric"`
}

//ProcessorData : Processor struct
type ProcessorData struct {
	CreateTimeUTC time.Time `json:"createTimeUTC"`
	CreatedBy     string    `json:"createdBy"`
	Index         int32     `json:"index"`
	//	InterruptsPerSec     int64 `json:"createTimeUTC"`
	Name                 string          `json:"name"`
	Type                 string          `json:"type"`
	NumOfCores           int32           `json:"numOfCores"`
	ProcessorQueueLength float64         `json:"processorQueueLength"`
	Metric               ProcessorMetricData `json:"metric"`
	Cores                []CoreData          `json:"cores"`
}

//ProcessorsData : Processors struct
type ProcessorsData struct {
	CreateTimeUTC    time.Time       `json:"createTimeUTC"`
	CreatedBy        string          `json:"createdBy"`
	Index            int32           `json:"index"`
	InterruptsPerSec int64           `json:"interruptsPerSec"`
	Name             string          `json:"name"`
	Type             string          `json:"type"`
	Metric           ProcessorMetricData  `json:"metric"`
	CPUs             []ProcessorData     `json:"cpus"`
}

//ProcessorMetricType : ProcessorMetric graphql object
var ProcessorMetricType = graphql.NewObject(graphql.ObjectConfig{
	Name: "processorMetric",
	Fields: graphql.Fields{
		"numOfProcesses": &graphql.Field{
			Description: "Number of processes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.NumOfProcesses, nil
				}
				return nil, nil
			},
		},

		"percentIOTime": &graphql.Field{
			Description: "PercentIOTime",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.PercentIOTime, nil
				}
				return nil, nil
			},
		},

		"percentIdleTime": &graphql.Field{
			Description: "PercentIdleTime",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.PercentIdleTime, nil
				}
				return nil, nil
			},
		},

		"percentInterruptTime": &graphql.Field{
			Description: "PercentInterruptTime",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.PercentInterruptTime, nil
				}
				return nil, nil
			},
		},

		"percentSystemTime": &graphql.Field{
			Description: "PercentSystemTime",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.PercentSystemTime, nil
				}
				return nil, nil
			},
		},

		"percentUserTime": &graphql.Field{
			Description: "PercentUserTime",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.PercentUserTime, nil
				}
				return nil, nil
			},
		},

		"percentUtil": &graphql.Field{
			Description: "PercentUtil",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.PercentUtil, nil
				}
				return nil, nil
			},
		},

		"interruptsPerSec": &graphql.Field{
			Description: "InterruptsPerSec",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorMetricData); ok {
					return CurrData.InterruptsPerSec, nil
				}
				return nil, nil
			},
		},
	},
})

//CoreType : Core graphql object
var CoreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "core",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(CoreData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(CoreData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(CoreData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		/*	"interruptsPerSec": &graphql.Field{
			Description: "InterruptsPerSec",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(CoreData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},*/

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(CoreData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(CoreData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"metric": &graphql.Field{
			Description: "processor metrics",
			Type:        ProcessorMetricType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(CoreData); ok {
					return CurrData.Metric, nil
				}
				return nil, nil
			},
		},
	},
})

//ProcessorType : Processor graphql object
var ProcessorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "processor",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		/*	"interruptsPerSec": &graphql.Field{
			Description: "InterruptsPerSec",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.InterruptsPerSec, nil
				}
				return nil, nil
			},
		},*/

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"numOfCores": &graphql.Field{
			Description: "NumOfCores",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.NumOfCores, nil
				}
				return nil, nil
			},
		},

		"processorQueueLength": &graphql.Field{
			Description: "ProcessorQueueLength",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.ProcessorQueueLength, nil
				}
				return nil, nil
			},
		},

		"metric": &graphql.Field{
			Description: "processor metrics",
			Type:        ProcessorMetricType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.Metric, nil
				}
				return nil, nil
			},
		},

		"cores": &graphql.Field{
			Description: "array of cores",
			Type:        graphql.NewList(CoreType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorData); ok {
					return CurrData.Cores, nil
				}
				return nil, nil
			},
		},
	},
})

//ProcessorsType : Processors graphql object
var ProcessorsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"interruptsPerSec": &graphql.Field{
			Description: "InterruptsPerSec",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.InterruptsPerSec, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"metric": &graphql.Field{
			Description: "processor metrics",
			Type:        ProcessorMetricType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.Metric, nil
				}
				return nil, nil
			},
		},

		"cpus": &graphql.Field{
			Description: "array of cores",
			Type:        graphql.NewList(ProcessorType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(ProcessorsData); ok {
					return CurrData.CPUs, nil
				}
				return nil, nil
			},
		},
	},
})

//PerformanceProcessorsConnectionDefinition : PerformanceProcessorsConnectionDefinition structure
var PerformanceProcessorsConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PerformanceProcessors",
	NodeType: ProcessorsType,
})

//PerformanceProcessorsListData : PerformanceProcessors List struct
type PerformanceProcessorsListData struct {
	Processors     []ProcessorsData    `json:"processors"`
}

//PerformanceProcessorsListType : PerformanceProcessorsList GraphQL Schema
var PerformanceProcessorsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "performanceProcessorsList",
	Fields: graphql.Fields{
		"processors": &graphql.Field{
			Type:        PerformanceProcessorsConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "processors list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PerformanceProcessorsListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Processors {
						arraySliceRet = append(arraySliceRet, CurData.Processors[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&ProcessorsData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						CreateTimeUTCASC := func(p1, p2 interface{}) bool {
							return p1.(ProcessorsData).CreateTimeUTC.Before(p2.(ProcessorsData).CreateTimeUTC)
						}
						CreateTimeUTCDESC := func(p1, p2 interface{}) bool {
							return p1.(ProcessorsData).CreateTimeUTC.After(p2.(ProcessorsData).CreateTimeUTC)
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == CreateTimeUTCCol {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(CreateTimeUTCASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(CreateTimeUTCDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("PerformanceProcessorsData Sort [" + Column + "] No such column exist!!!")
							}
						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})
