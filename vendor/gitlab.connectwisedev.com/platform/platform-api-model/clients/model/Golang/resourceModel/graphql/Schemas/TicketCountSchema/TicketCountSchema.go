package TicketCountSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//TicketCountData : TicketCountData Structure
type TicketCountData struct {
	RegID       int64	`json:"regId"`
	Count       int64	`json:"count"`
	SiteID      int64	`json:"siteId"`
}

//TicketCountDataType : TicketCountData GraphQL Schema
var TicketCountDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ticketCount",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketCountData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"count": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketCountData); ok {
					return CurData.Count, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketCountData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
	},
})

//TicketCountConnectionDefinition : TicketCountConnectionDefinition structure
var TicketCountConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TicketCount",
	NodeType: TicketCountDataType,
})

//TicketCountListData : TicketCountListData Structure
type TicketCountListData struct {
	TicketCount 	[]TicketCountData 	`json:"ticketCountList"`
	TotalCount      int64                 	`json:"totaCount"`
}

//TicketCountListType : TicketCountList GraphQL Schema
var TicketCountListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TicketCountList",
	Fields: graphql.Fields{
		"ticketCountList": &graphql.Field{
			Type:        TicketCountConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "ticket Count List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TicketCountListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.TicketCount {
						arraySliceRet = append(arraySliceRet, CurData.TicketCount[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&TicketCountData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(TicketCountData).SiteID < p2.(TicketCountData).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketCountData).SiteID > p2.(TicketCountData).SiteID
						}

						RegIDASC := func(p1, p2 interface{}) bool {
							return p1.(TicketCountData).RegID < p2.(TicketCountData).RegID
						}
						RegIDDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketCountData).RegID > p2.(TicketCountData).RegID
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "SITEID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SiteIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(SiteIDDESC).Sort(arraySliceRet)
								}
							}else if strings.ToUpper(Column) == "REGID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RegIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(RegIDDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("TicketCountData Sort [" + Column + "] No such column exist!!!")
							}
						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},

		"totalCount": &graphql.Field{
			Type:        graphql.String,
			Description: "totalCount of TicketCount list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketCountListData); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
