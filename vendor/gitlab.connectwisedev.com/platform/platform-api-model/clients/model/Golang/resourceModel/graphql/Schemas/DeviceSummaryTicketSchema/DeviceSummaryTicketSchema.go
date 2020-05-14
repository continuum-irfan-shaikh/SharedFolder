package DeviceSummaryTicketSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//DeviceSummaryTicketDetailsList : DeviceSummaryTicketDetailsList Structure
type DeviceSummaryTicketDetailsList struct {
	ErrorMessage []string        `json:"errorMessage"`
	Data         []TicketDetails `json:"deviceSummaryTicketDetailsList"`
}

//TicketDetails : TicketDetails Stkructure
type TicketDetails struct {
	ID              int64  `json:"id"`
	CreatedDateTime string `json:"createddatetime"`
	SubjectLine     string `json:"subjectline"`
}

//DeviceSummaryTicketDetailsListType : DeviceSummaryTicketDetailsList GraphQL Schema
var DeviceSummaryTicketDetailsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeviceSummaryTicketDetailsList",
	Fields: graphql.Fields{
		"errorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "errorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceSummaryTicketDetailsList); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},

		"deviceSummaryTicketDetailsList": &graphql.Field{
			Type:        DeviceSummaryTicketDetailsConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "ticket details list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(DeviceSummaryTicketDetailsList); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Data {
						arraySliceRet = append(arraySliceRet, CurData.Data[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&TicketDetails{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")

						IDASC := func(p1, p2 interface{}) bool {
							return p1.(TicketDetails).ID < p2.(TicketDetails).ID
						}
						IDDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketDetails).ID > p2.(TicketDetails).ID
						}

						CreatedDateTimeASC := func(p1, p2 interface{}) bool {
							return p1.(TicketDetails).CreatedDateTime < p2.(TicketDetails).CreatedDateTime
						}
						CreatedDateTimeDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketDetails).CreatedDateTime > p2.(TicketDetails).CreatedDateTime
						}

						SubjectLineASC := func(p1, p2 interface{}) bool {
							return p1.(TicketDetails).SubjectLine < p2.(TicketDetails).SubjectLine
						}
						SubjectLineDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketDetails).SubjectLine > p2.(TicketDetails).SubjectLine
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "ID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(IDASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(IDDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "CREATEDDATETIME" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(CreatedDateTimeASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(CreatedDateTimeDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "SUBJECTLINE" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SubjectLineASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(SubjectLineDESC).Sort(arraySliceRet)
								}
							} else {
								return nil, errors.New("TicketDetails Sort [" + Column + "] No such column exist!!!")
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

//TicketDetailsType : TicketDetailsType GraphQL Schema
var TicketDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ticketDetails",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketDetails); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},

		"createddatetime": &graphql.Field{
			Type:        graphql.String,
			Description: "createddatetime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketDetails); ok {
					return CurData.CreatedDateTime, nil
				}
				return nil, nil
			},
		},

		"subjectline": &graphql.Field{
			Type:        graphql.String,
			Description: "subjectline",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketDetails); ok {
					return CurData.SubjectLine, nil
				}
				return nil, nil
			},
		},
	},
})

//DeviceSummaryTicketDetailsConnectionDefinition : DeviceSummaryTicketDetailsConnectionDefinition structure
var DeviceSummaryTicketDetailsConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TicketDetails",
	NodeType: TicketDetailsType,
})

//TicketDetailsList : TicketDetailsList Structure
type TicketDetailsList struct {
	Status int64           `json:"status"`
	Data   []TicketDetails `json:"outdata"`
}
