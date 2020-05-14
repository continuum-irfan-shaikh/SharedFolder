package AlertTicketSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//AlertTicketCount : AlertTicketCount Structure
type AlertTicketCount struct {
	Count	int64  	`json:"count"`
	RegID 	int64  	`json:"regid"`
}

//AlertTicketCountType : AlertTicketCount GraphQL Schema
var AlertTicketCountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "alertTicketCount",
	Fields: graphql.Fields{
		"count": &graphql.Field{
			Type:        graphql.String,
			Description: "count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCount); ok {
					return CurData.Count, nil
				}
				return nil, nil
			},
		},

		"regid": &graphql.Field{
			Type:        graphql.String,
			Description: "regid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCount); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},
	},
})

//AlertTicketCountConnectionDefinition : AlertTicketCountConnectionDefinition structure
var AlertTicketCountConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "AlertTicketCount",
	NodeType: AlertTicketCountType,
})

//AlertTicketCountList : AlertTicketCountList Structure
type AlertTicketCountList struct {
	Status	int64 			`json:"status"`
	Data	[]AlertTicketCount 	`json:"outdata"`
}

//AlertTicketCountListType : AlertTicketCountList GraphQL Schema
var AlertTicketCountListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "alertTicketCountList",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCountList); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        AlertTicketCountConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "alert ticket count list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(AlertTicketCountList); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Data {
						arraySliceRet = append(arraySliceRet, CurData.Data[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&AlertTicketCount{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						CountASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketCount).Count < p2.(AlertTicketCount).Count
						}
						CountDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketCount).Count > p2.(AlertTicketCount).Count
						}

						RegIDASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketCount).RegID < p2.(AlertTicketCount).RegID
						}
						RegIDDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketCount).RegID > p2.(AlertTicketCount).RegID
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "COUNT" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(CountASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(CountDESC).Sort(arraySliceRet)
								}
							}else if strings.ToUpper(Column) == "REGID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RegIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(RegIDDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("AlertTicketCount Sort [" + Column + "] No such column exist!!!")
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
