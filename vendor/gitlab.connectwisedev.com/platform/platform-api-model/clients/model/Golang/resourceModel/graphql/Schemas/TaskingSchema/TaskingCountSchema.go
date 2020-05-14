package TaskingSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//TaskingCountData : TaskingCountData Structure
type TaskingCountData struct {
	Count    int64    `json:"count"`
	TargetID string   `json:"managedEndpointId"`
}

//TaskingCountType : TaskingCountData GraphQL Schema
var TaskingCountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskingCountData",
	Fields: graphql.Fields{
		"taskSequenceCount": &graphql.Field{
			Type:        graphql.String,
			Description: "Tasking Count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingCountData); ok {
					return CurData.Count, nil
				}
				return nil, nil
			},
		},
		"managedEndpointId": &graphql.Field{
			Type:        graphql.String,
			Description: "TargetID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingCountData); ok {
					return CurData.TargetID, nil
				}
				return nil, nil
			},
		},
	},
})

//TaskingCountConnectionDefinition : TaskingCountConnectionDefinition structure
var TaskingCountConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TaskingCount",
	NodeType: TaskingCountType,
})

//TaskingCountListData : TaskingCountListData Structure
type TaskingCountListData struct {
	TaskingCount []TaskingCountData `json:"taskingCountList"`
}

//TaskingCountListType : TaskingCountList GraphQL Schema
var TaskingCountListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskingCountList",
	Fields: graphql.Fields{
		"taskingCountList": &graphql.Field{
			Type:        TaskingCountConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "TaskingCount list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TaskingCountListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.TaskingCount {
						arraySliceRet = append(arraySliceRet, CurData.TaskingCount[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&TaskingCountData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						TargetIDASC := func(p1, p2 interface{}) bool {
							return p1.(TaskingCountData).TargetID < p2.(TaskingCountData).TargetID
						}
						TargetIDDESC := func(p1, p2 interface{}) bool {
							return p1.(TaskingCountData).TargetID > p2.(TaskingCountData).TargetID
						}

						TaskCountASC := func(p1, p2 interface{}) bool {
							return p1.(TaskingCountData).Count < p2.(TaskingCountData).Count
						}
						TaskCountDESC := func(p1, p2 interface{}) bool {
							return p1.(TaskingCountData).Count > p2.(TaskingCountData).Count
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "TARGETID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TargetIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(TargetIDDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "TASKCOUNT" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskCountASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(TaskCountDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("TaskingCountData Sort [" + Column + "] No such column exist!!!")
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
