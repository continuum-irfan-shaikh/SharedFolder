package AlertTicketSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//AlertTicketResource : AlertTicketResource Structure
type AlertTicketResource struct {
	RegID 	int64  	`json:"regid"`
}

//AlertTicketResourceType : AlertTicketResource GraphQL Schema
var AlertTicketResourceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "alertTicketResource",
	Fields: graphql.Fields{
		"regid": &graphql.Field{
			Type:        graphql.String,
			Description: "regid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketResource); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},
	},
})

//AlertTicketResourceConnectionDefinition : AlertTicketResourceConnectionDefinition structure
var AlertTicketResourceConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "AlertTicketResource",
	NodeType: AlertTicketResourceType,
})

//AlertTicketResourceList : AlertTicketResourceList Structure
type AlertTicketResourceList struct {
	Status	int64 			`json:"status"`
	Data	[]AlertTicketResource 	`json:"outdata"`
}

//AlertTicketResourceListType : AlertTicketResourceList GraphQL Schema
var AlertTicketResourceListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "alertTicketResourceList",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketResourceList); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        AlertTicketResourceConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "alert ticket Resource list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(AlertTicketResourceList); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Data {
						arraySliceRet = append(arraySliceRet, CurData.Data[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&AlertTicketResource{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						RegIDASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketResource).RegID < p2.(AlertTicketResource).RegID
						}
						RegIDDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketResource).RegID > p2.(AlertTicketResource).RegID
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "REGID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RegIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(RegIDDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("AlertTicketResource Sort [" + Column + "] No such column exist!!!")
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
