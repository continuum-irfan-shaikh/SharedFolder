package AlertTicketSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//AlertTicketDetails : AlertTicketDetails Structure
type AlertTicketDetails struct {
	AssigndTo       string `json:"assignd To"`
	Criticality     string `json:"criticality"`
	FamilyName      string `json:"familyname"`
	ID              int64  `json:"id"`
	CreatedDateTime string `json:"createddatetime"`
	Priority        string `json:"priority"`
	QType           string `json:"qtype"`
	StatusName      string `json:"statusname"`
	SubjectLine     string `json:"subjectline"`
	Description     string `json:"description"`
	TimeDiff        int64  `json:"timediff"`
	RefTktID        int64  `json:"refTicketID"`
}

//AlertTicketDetailsType : AlertTicketDetails GraphQL Schema
var AlertTicketDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "alertTicketDetails",
	Fields: graphql.Fields{
		"assigndTo": &graphql.Field{
			Type:        graphql.String,
			Description: "assigndTo",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.AssigndTo, nil
				}
				return nil, nil
			},
		},

		"criticality": &graphql.Field{
			Type:        graphql.String,
			Description: "criticality",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.Criticality, nil
				}
				return nil, nil
			},
		},

		"familyname": &graphql.Field{
			Type:        graphql.String,
			Description: "familyname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.FamilyName, nil
				}
				return nil, nil
			},
		},

		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},

		"createddatetime": &graphql.Field{
			Type:        graphql.String,
			Description: "createddatetime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.CreatedDateTime, nil
				}
				return nil, nil
			},
		},

		"priority": &graphql.Field{
			Type:        graphql.String,
			Description: "priority",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.Priority, nil
				}
				return nil, nil
			},
		},

		"qtype": &graphql.Field{
			Type:        graphql.String,
			Description: "qtype",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.QType, nil
				}
				return nil, nil
			},
		},

		"statusname": &graphql.Field{
			Type:        graphql.String,
			Description: "statusname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.StatusName, nil
				}
				return nil, nil
			},
		},

		"subjectline": &graphql.Field{
			Type:        graphql.String,
			Description: "subjectline",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.SubjectLine, nil
				}
				return nil, nil
			},
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "description",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},

		"timediff": &graphql.Field{
			Type:        graphql.String,
			Description: "timediff",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.TimeDiff, nil
				}
				return nil, nil
			},
		},
		"refTicketID": &graphql.Field{
			Type:        graphql.String,
			Description: "refTicketID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetails); ok {
					return CurData.RefTktID, nil
				}
				return nil, nil
			},
		},
	},
})

//AlertTicketDetailsConnectionDefinition : AlertTicketDetailsConnectionDefinition structure
var AlertTicketDetailsConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "AlertTicketDetails",
	NodeType: AlertTicketDetailsType,
})

//AlertTicketDetailsList : AlertTicketDetailsList Structure
type AlertTicketDetailsList struct {
	Status int64                `json:"status"`
	Data   []AlertTicketDetails `json:"outdata"`
}

//AlertTicketDetailsListType : AlertTicketDetailsList GraphQL Schema
var AlertTicketDetailsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "alertTicketDetailsList",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketDetailsList); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        AlertTicketDetailsConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "alert ticket details list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(AlertTicketDetailsList); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Data {
						arraySliceRet = append(arraySliceRet, CurData.Data[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&AlertTicketDetails{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						AssigndToASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).AssigndTo < p2.(AlertTicketDetails).AssigndTo
						}
						AssigndToDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).AssigndTo > p2.(AlertTicketDetails).AssigndTo
						}

						CriticalityASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).Criticality < p2.(AlertTicketDetails).Criticality
						}
						CriticalityDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).Criticality > p2.(AlertTicketDetails).Criticality
						}

						FamilyNameASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).FamilyName < p2.(AlertTicketDetails).FamilyName
						}
						FamilyNameDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).FamilyName > p2.(AlertTicketDetails).FamilyName
						}

						IDASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).ID < p2.(AlertTicketDetails).ID
						}
						IDDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).ID > p2.(AlertTicketDetails).ID
						}

						CreatedDateTimeASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).CreatedDateTime < p2.(AlertTicketDetails).CreatedDateTime
						}
						CreatedDateTimeDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).CreatedDateTime > p2.(AlertTicketDetails).CreatedDateTime
						}

						PriorityASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).Priority < p2.(AlertTicketDetails).Priority
						}
						PriorityDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).Priority > p2.(AlertTicketDetails).Priority
						}

						QTypeASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).QType < p2.(AlertTicketDetails).QType
						}
						QTypeDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).QType > p2.(AlertTicketDetails).QType
						}

						StatusNameASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).StatusName < p2.(AlertTicketDetails).StatusName
						}
						StatusNameDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).StatusName > p2.(AlertTicketDetails).StatusName
						}

						SubjectLineASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).SubjectLine < p2.(AlertTicketDetails).SubjectLine
						}
						SubjectLineDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).SubjectLine > p2.(AlertTicketDetails).SubjectLine
						}

						DescriptionASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).Description < p2.(AlertTicketDetails).Description
						}
						DescriptionDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).Description > p2.(AlertTicketDetails).Description
						}

						TimeDiffASC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).TimeDiff < p2.(AlertTicketDetails).TimeDiff
						}
						TimeDiffDESC := func(p1, p2 interface{}) bool {
							return p1.(AlertTicketDetails).TimeDiff > p2.(AlertTicketDetails).TimeDiff
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "ASSIGNDTO" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(AssigndToASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(AssigndToDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "CRITICALITY" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(CriticalityASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(CriticalityDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "FAMILYNAME" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(FamilyNameASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(FamilyNameDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "ID" {
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
							} else if strings.ToUpper(Column) == "PRIORITY" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(PriorityASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(PriorityDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "QTYPE" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(QTypeASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(QTypeDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "STATUSNAME" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(StatusNameASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(StatusNameDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "SUBJECTLINE" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SubjectLineASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(SubjectLineDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "DESCRIPTION" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(DescriptionASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(DescriptionDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "TIMEDIFF" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TimeDiffASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(TimeDiffDESC).Sort(arraySliceRet)
								}
							} else {
								return nil, errors.New("AlertTicketDetails Sort [" + Column + "] No such column exist!!!")
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
