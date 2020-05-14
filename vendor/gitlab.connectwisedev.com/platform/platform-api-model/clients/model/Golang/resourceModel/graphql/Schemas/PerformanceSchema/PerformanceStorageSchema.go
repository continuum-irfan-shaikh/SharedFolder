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

//StorageMetricData : StorageMetric struct
type StorageMetricData struct {
	IdleTime        	float64 `json:"idleTime"`
	WriteCompleted  	int64   `json:"writeCompleted"`
	WriteTimeMs     	int64   `json:"writeTimeMs"`
      	AvgDiskWriteQueueLength int64	`json:"avgDiskWriteQueueLength"`
      	AvgDiskSecPerWrite	int64	`json:"avgDiskSecPerWrite"`
	ReadCompleted   	int64   `json:"readCompleted"`
	ReadTimeMs      	int64   `json:"readTimeMs"`
	AvgDiskReadQueueLength	int64   `json:"avgDiskReadQueueLength"`
      	AvgDiskSecPerRead	int64   `json:"avgDiskSecPerRead"`
	FreeSpaceBytes  	int64   `json:"freeSpaceBytes"`
	UsedSpaceBytes  	int64   `json:"usedSpaceBytes"`
	TotalSpaceBytes 	int64   `json:"totalSpaceBytes"`
	DiskTimeTotal   	float64 `json:"diskTimeTotal"`
}


//PerformanceStoragePartitionData : PerformanceStoragePartition struct
type PerformanceStoragePartitionData struct {
	CreateTimeUTC time.Time     `json:"createTimeUTC"`
	CreatedBy     string        `json:"createdBy"`
	Index         int32         `json:"index"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Mounted       bool          `json:"mounted"`
	Metric        StorageMetricData `json:"metric"`
}

//PerformanceStorageData : PerformanceStorage struct
type PerformanceStorageData struct {
	CreateTimeUTC time.Time                     `json:"createTimeUTC"`
	CreatedBy     string                        `json:"createdBy"`
	Index         int32                         `json:"index"`
	Name          string                        `json:"name"`
	Type          string                        `json:"type"`
	Metric        StorageMetricData                 `json:"metric"`
	Partitions    []PerformanceStoragePartitionData `json:"partitions"`
}


//PerformanceStoragesData : PerformanceStorages struct
type PerformanceStoragesData struct {
	CreateTimeUTC time.Time            `json:"createTimeUTC"`
	CreatedBy     string               `json:"createdBy"`
	Index         int32                `json:"index"`
	Name          string               `json:"name"`
	Type          string               `json:"type"`
	Storages      []PerformanceStorageData `json:"storages"`
}

//StorageMetricType : StorageMetric graphql object
var StorageMetricType = graphql.NewObject(graphql.ObjectConfig{
	Name: "storageMetric",
	Fields: graphql.Fields{
		"idleTime": &graphql.Field{
			Description: "Disk Idle Time",
			Type:        graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.IdleTime, nil
				}
				return nil, nil
			},
		},

		"writeCompleted": &graphql.Field{
			Description: "Writes Completed",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.WriteCompleted, nil
				}
				return nil, nil
			},
		},

		"writeTimeMs": &graphql.Field{
			Description: "Write completed in millisecond",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.WriteTimeMs, nil
				}
				return nil, nil
			},
		},

		"avgDiskWriteQueueLength": &graphql.Field{
			Description: "Average number of disk write requests over sample interval",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.AvgDiskWriteQueueLength, nil
				}
				return nil, nil
			},
		},

		"avgDiskSecPerWrite": &graphql.Field{
			Description: "Average time taken for write operations on disk",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.AvgDiskSecPerWrite, nil
				}
				return nil, nil
			},
		},

		"readCompleted": &graphql.Field{
			Description: "Read Completed",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.ReadCompleted, nil
				}
				return nil, nil
			},
		},

		"readTimeMs": &graphql.Field{
			Description: "Read completed in millisecond",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.ReadTimeMs, nil
				}
				return nil, nil
			},
		},

		"avgDiskReadQueueLength": &graphql.Field{
			Description: "Average number of disk read requests over sample interval",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.AvgDiskReadQueueLength, nil
				}
				return nil, nil
			},
		},

		"avgDiskSecPerRead": &graphql.Field{
			Description: "Average time taken for read operations on disk",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.AvgDiskSecPerRead, nil
				}
				return nil, nil
			},
		},

		"freeSpaceBytes": &graphql.Field{
			Description: "Disk free space",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.FreeSpaceBytes, nil
				}
				return nil, nil
			},
		},

		"usedSpaceBytes": &graphql.Field{
			Description: "Disk used space",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.UsedSpaceBytes, nil
				}
				return nil, nil
			},
		},

		"totalSpaceBytes": &graphql.Field{
			Description: "Disk total space ",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.TotalSpaceBytes, nil
				}
				return nil, nil
			},
		},

		"diskTimeTotal": &graphql.Field{
			Description: "Total time spent on read write operation",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(StorageMetricData); ok {
					return CurrData.DiskTimeTotal, nil
				}
				return nil, nil
			},
		},
	},
})

//PerformanceStoragePartitionType : PerformanceStoragePartition graphql object
var PerformanceStoragePartitionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PerformanceStoragePartitionData",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragePartitionData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragePartitionData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragePartitionData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragePartitionData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragePartitionData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"mounted": &graphql.Field{
			Description: "Mounted",
			Type:        graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragePartitionData); ok {
					return CurrData.Mounted, nil
				}
				return nil, nil
			},
		},

		"metric": &graphql.Field{
			Description: "storage metrics  ",
			Type:        StorageMetricType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragePartitionData); ok {
					return CurrData.Metric, nil
				}
				return nil, nil
			},
		},
	},
})

//PerformanceStorageType : PerformanceStorage  graphql object
var PerformanceStorageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PerformanceStorageData",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStorageData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStorageData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStorageData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStorageData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStorageData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"metric": &graphql.Field{
			Description: "storage metrics  ",
			Type:        StorageMetricType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStorageData); ok {
					return CurrData.Metric, nil
				}
				return nil, nil
			},
		},

		"partitions": &graphql.Field{
			Description: "storage partition  ",
			Type:        graphql.NewList(PerformanceStoragePartitionType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStorageData); ok {
					return CurrData.Partitions, nil
				}
				return nil, nil
			},
		},
	},
})

//PerformanceStoragesType : PerformanceStorages graphql object
var PerformanceStoragesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "performanceStorages",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragesData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragesData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragesData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragesData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragesData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"storages": &graphql.Field{
			Description: "Storages",
			Type:        graphql.NewList(PerformanceStorageType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceStoragesData); ok {
					return CurrData.Storages, nil
				}
				return nil, nil
			},
		},
	},
})

//PerformanceStorageConnectionDefinition : PerformanceStorageConnectionDefinition structure
var PerformanceStorageConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PerformanceStorage",
	NodeType: PerformanceStoragesType,
})

//PerformanceStorageListData : PerformanceStorage List struct
type PerformanceStorageListData struct {
	Storage     []PerformanceStoragesData    `json:"storage"`
}

//PerformanceStorageListType : PerformanceStorageList GraphQL Schema
var PerformanceStorageListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "performanceStorageList",
	Fields: graphql.Fields{
		"storage": &graphql.Field{
			Type:        PerformanceStorageConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Storage list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PerformanceStorageListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Storage {
						arraySliceRet = append(arraySliceRet, CurData.Storage[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PerformanceStoragesData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						CreateTimeUTCASC := func(p1, p2 interface{}) bool {
							return p1.(PerformanceStoragesData).CreateTimeUTC.Before(p2.(PerformanceStoragesData).CreateTimeUTC)
						}
						CreateTimeUTCDESC := func(p1, p2 interface{}) bool {
							return p1.(PerformanceStoragesData).CreateTimeUTC.After(p2.(PerformanceStoragesData).CreateTimeUTC)
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
								return nil, errors.New("PerformanceStoragesData Sort [" + Column + "] No such column exist!!!")
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
