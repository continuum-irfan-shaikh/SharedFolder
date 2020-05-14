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

//PerformanceMemoryData : PerformanceMemory struct
type PerformanceMemoryData struct {
	CreateTimeUTC                   time.Time `json:"createTimeUTC"`
	CreatedBy                       string    `json:"createdBy"`
	Index                           int32     `json:"index"`
	Name                            string    `json:"name"`
	Type                            string    `json:"type"`
	PhysicalTotalBytes              int64     `json:"physicalTotalBytes"`
	PhysicalInUseBytes              int64     `json:"physicalInUseBytes"`
	PhysicalAvailableBytes          int64     `json:"physicalAvailableBytes"`
	VirtualTotalBytes               int64     `json:"virtualTotalBytes"`
	VirtualInUseBytes               int64     `json:"virtualInUseBytes"`
	VirtualAvailableBytes           int64     `json:"virtualAvailableBytes"`
	PercentCommittedInUseBytes      int64     `json:"percentCommittedInUseBytes"`
	CommittedBytes                  int64     `json:"committedBytes"`
	FreeSystemPageTableEntriesBytes int64     `json:"freeSystemPageTableEntriesBytes"`
	PoolNonPagedBytes               int64     `json:"poolNonPagedBytes"`
	PagesPerSecondBytes             int64     `json:"pagesPerSecondBytes"`
	SwapinPerSecondBytes            int64     `json:"swapinPerSecondBytes"`
	SwapoutPerSecondBytes           int64     `json:"swapoutPerSecondBytes"`
	PagesOutputPerSec               int64     `json:"pagesOutputPerSec"`
}

//PerformanceMemoryType : PerformanceMemory graphql object
var PerformanceMemoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Description: "CreateTimeUTC",
			Type:        CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Description: "CreatedBy",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"index": &graphql.Field{
			Description: "Index",
			Type:        graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.Index, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Description: "Name",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Description: "Type",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.Type, nil
				}
				return nil, nil
			},
		},

		"physicalTotalBytes": &graphql.Field{
			Description: "Total RAM of physical memory in bytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.PhysicalTotalBytes, nil
				}
				return nil, nil
			},
		},

		"physicalInUseBytes": &graphql.Field{
			Description: "Total RAM in use of physical memory in bytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.PhysicalInUseBytes, nil
				}
				return nil, nil
			},
		},

		"physicalAvailableBytes": &graphql.Field{
			Description: " Total RAM available of physical memory in bytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.PhysicalAvailableBytes, nil
				}
				return nil, nil
			},
		},

		"virtualTotalBytes": &graphql.Field{
			Description: "Total RAM of virtual memory in bytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.VirtualTotalBytes, nil
				}
				return nil, nil
			},
		},

		"virtualInUseBytes": &graphql.Field{
			Description: "Total RAM in use of virtual memory in bytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.VirtualInUseBytes, nil
				}
				return nil, nil
			},
		},

		"virtualAvailableBytes": &graphql.Field{
			Description: "Total RAM available of virtual memory in bytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.VirtualAvailableBytes, nil
				}
				return nil, nil
			},
		},

		"percentCommittedInUseBytes": &graphql.Field{
			Description: "percentCommittedInUseBytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.PercentCommittedInUseBytes, nil
				}
				return nil, nil
			},
		},

		"committedBytes": &graphql.Field{
			Description: "committedBytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.CommittedBytes, nil
				}
				return nil, nil
			},
		},

		"freeSystemPageTableEntriesBytes": &graphql.Field{
			Description: "freeSystemPageTableEntriesBytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.FreeSystemPageTableEntriesBytes, nil
				}
				return nil, nil
			},
		},

		"poolNonPagedBytes": &graphql.Field{
			Description: "poolNonPagedBytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.PoolNonPagedBytes, nil
				}
				return nil, nil
			},
		},

		"pagesPerSecondBytes": &graphql.Field{
			Description: "pagesPerSecondBytes",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.PagesPerSecondBytes, nil
				}
				return nil, nil
			},
		},

		"swapinPerSecondBytes": &graphql.Field{
			Description: "Rate at which memory is swapped from disk into active memory per second",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.SwapinPerSecondBytes, nil
				}
				return nil, nil
			},
		},

		"swapoutPerSecondBytes": &graphql.Field{
			Description: "Rate at which memory is being swapped from active memory to disk per second",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.SwapoutPerSecondBytes, nil
				}
				return nil, nil
			},
		},

		"pagesOutputPerSec": &graphql.Field{
			Description: "Number of pages written to disk to free up space in physical memory",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PerformanceMemoryData); ok {
					return CurrData.PagesOutputPerSec, nil
				}
				return nil, nil
			},
		},
	},
})

//PerformanceMemoryConnectionDefinition : PerformanceMemoryConnectionDefinition structure
var PerformanceMemoryConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PerformanceMemory",
	NodeType: PerformanceMemoryType,
})

//PerformanceMemoryListData : PerformanceMemory List struct
type PerformanceMemoryListData struct {
	Memory     []PerformanceMemoryData    `json:"memory"`
}

//PerformanceMemoryListType : PerformanceMemoryList GraphQL Schema
var PerformanceMemoryListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "performanceMemoryList",
	Fields: graphql.Fields{
		"memory": &graphql.Field{
			Type:        PerformanceMemoryConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "memory list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PerformanceMemoryListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Memory {
						arraySliceRet = append(arraySliceRet, CurData.Memory[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PerformanceMemoryData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						CreateTimeUTCASC := func(p1, p2 interface{}) bool {
							return p1.(PerformanceMemoryData).CreateTimeUTC.Before(p2.(PerformanceMemoryData).CreateTimeUTC)
						}
						CreateTimeUTCDESC := func(p1, p2 interface{}) bool {
							return p1.(PerformanceMemoryData).CreateTimeUTC.After(p2.(PerformanceMemoryData).CreateTimeUTC)
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
								return nil, errors.New("PerformanceMemoryData Sort [" + Column + "] No such column exist!!!")
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

