package DeviceSummaryTicketSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//DeviceDetailsCriticalTicketCount : DeviceDetailsCriticalTicketCount Structure
type DeviceDetailsCriticalTicketCount struct {
	Count int64 `json:"count"`
	RegID int64 `json:"regid"`
}

//DeviceDetailsCriticalTicketCountType : DeviceDetailsCriticalTicketCount GraphQL Schema
var DeviceDetailsCriticalTicketCountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeviceDetailsCriticalTicketCount",
	Fields: graphql.Fields{
		"count": &graphql.Field{
			Type:        graphql.String,
			Description: "count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsCriticalTicketCount); ok {
					return CurData.Count, nil
				}
				return nil, nil
			},
		},

		"regid": &graphql.Field{
			Type:        graphql.String,
			Description: "regid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsCriticalTicketCount); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},
	},
})

//DeviceDetailsCriticalTicketCountConnectionDefinition : DeviceDetailsCriticalTicketCountConnectionDefinition structure
var DeviceDetailsCriticalTicketCountConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "DeviceDetailsCriticalTicketCount",
	NodeType: DeviceDetailsCriticalTicketCountType,
})

//DeviceDetailsCriticalTicketCountList : DeviceDetailsCriticalTicketCountList Structure
type DeviceDetailsCriticalTicketCountList struct {
	Status int64                              `json:"status"`
	Data   []DeviceDetailsCriticalTicketCount `json:"outdata"`
}

//DeviceDetailsCriticalTicketCountListType : DeviceDetailsCriticalTicketCountList GraphQL Schema
var DeviceDetailsCriticalTicketCountListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeviceDetailsCriticalTicketCountList",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsCriticalTicketCountList); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        DeviceDetailsCriticalTicketCountConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "ticket count list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(DeviceDetailsCriticalTicketCountList); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Data {
						arraySliceRet = append(arraySliceRet, CurData.Data[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&DeviceDetailsCriticalTicketCount{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						CountASC := func(p1, p2 interface{}) bool {
							return p1.(DeviceDetailsCriticalTicketCount).Count < p2.(DeviceDetailsCriticalTicketCount).Count
						}
						CountDESC := func(p1, p2 interface{}) bool {
							return p1.(DeviceDetailsCriticalTicketCount).Count > p2.(DeviceDetailsCriticalTicketCount).Count
						}

						RegIDASC := func(p1, p2 interface{}) bool {
							return p1.(DeviceDetailsCriticalTicketCount).RegID < p2.(DeviceDetailsCriticalTicketCount).RegID
						}
						RegIDDESC := func(p1, p2 interface{}) bool {
							return p1.(DeviceDetailsCriticalTicketCount).RegID > p2.(DeviceDetailsCriticalTicketCount).RegID
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "COUNT" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(CountASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(CountDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "REGID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RegIDASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(RegIDDESC).Sort(arraySliceRet)
								}
							} else {
								return nil, errors.New("DeviceDetailsCriticalTicketCount Sort [" + Column + "] No such column exist!!!")
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
