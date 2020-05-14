package TaskingSchema

import (
	"encoding/json"
	"errors"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
	"reflect"
	"strings"
	"time"
)

const skipRunningVariableName = "skipRunning"

//TaskingSummaryPageType : TaskingSummaryPageType GraphQL Schema
var TaskingSummaryPageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskingSummaryPageData",
	Fields: graphql.Fields{
		"taskID": &graphql.Field{
			Type:        graphql.String,
			Description: "Task ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.TaskID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Task Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Task Type (Task or Sequence)",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return DefineTaskTypeView(CurData.Type), nil
				}
				return nil, nil
			},
		},
		"runOn": &graphql.Field{
			Type:        RunOnType,
			Description: "On how many devices/sites/dynamic groups task is running",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.RunOn, nil
				}
				return nil, nil
			},
		},
		"regularity": &graphql.Field{
			Type:        graphql.String,
			Description: "Schedule type: recurrent/one time/run now",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.Regularity, nil
				}
				return nil, nil
			},
		},
		"initiatedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "Task's author",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.InitiatedBy, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Task's status: active/inactive",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"lastRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "Task's last run time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.LastRunTime.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"lastRunStatus": &graphql.Field{
			Type:        LastRunStatusType,
			Description: "Task's last run status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.LastRunStatus, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.String,
			Description: "Task's created time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.TaskSummaryData); ok {
					return CurData.CreatedAt.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
	},
})

//TaskingSummaryPageConnDef : TaskingSummaryPageConnDef structure
var TaskingSummaryPageConnDef = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TaskingSummaryPage",
	NodeType: TaskingSummaryPageType,
})

//TaskingSummaryPageListData : List of TaskingSummaryPage Structures
type TaskingSummaryPageListData struct {
	TaskingSummaryPage []apiModels.TaskSummaryData `json:"taskingSummaryPageList"`
}

//TaskingSummaryPageListType : TaskingSummaryPageListType GraphQL Schema
var TaskingSummaryPageListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "taskingSummaryPageList",
	Fields: graphql.Fields{
		"taskingSummaryPageList": &graphql.Field{
			Type:        TaskingSummaryPageConnDef.ConnectionType,
			Args:        appendArg(Relay.ConnectionArgs, skipRunningVariableName, graphql.ArgumentConfig{Type: graphql.Int}),
			Description: "Tasking Summary Page list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TaskingSummaryPageListData); ok {
					var arraySliceRet []interface{}

					type WidgetInfo struct {
						Failed       int `json:"failedCount"`
						SomeFailures int `json:"someFailuresCount"`
						Running      int `json:"runningCount"`
						Success      int `json:"successCount"`
						Total        int `json:"total"`
					}

					var widgetInfo = WidgetInfo{}
					var skipRunning int
					var remainder int

					if skipRun, ok := p.Args[skipRunningVariableName]; ok {
						if skipRun, ok := skipRun.(int); ok {
							skipRunning = skipRun
						}
					}

					for ind := range CurData.TaskingSummaryPage {
						if !(skipRunning == 1 && getSortableStatus(CurData.TaskingSummaryPage[ind].LastRunStatus) == runningStatus) {
							arraySliceRet = append(arraySliceRet, CurData.TaskingSummaryPage[ind])
						}

						switch getSortableStatus(CurData.TaskingSummaryPage[ind].LastRunStatus) {
						case failedStatus:
							widgetInfo.Failed++
						case someFailuresStatus:
							widgetInfo.SomeFailures++
						case runningStatus:
							widgetInfo.Running++
						case successStatus:
							widgetInfo.Success++
						default:
							remainder++
						}
					}

					widgetInfo.Total = widgetInfo.Failed + widgetInfo.SomeFailures + widgetInfo.Running + widgetInfo.Success + remainder

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&apiModels.TaskSummaryData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							switch strings.ToUpper(Column) {
							case "NAME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryNameASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryNameDESC).Sort(arraySliceRet)
								}
							case "LASTRUNSTATUS":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryLRStatusASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryLRStatusDESC).Sort(arraySliceRet)
								}
							case "TYPE":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryTypeASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryTypeDESC).Sort(arraySliceRet)
								}
							case "STATUS":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryStatusASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryStatusDESC).Sort(arraySliceRet)
								}
							case "INITIATEDBY":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryInitiatedByASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryInitiatedByDESC).Sort(arraySliceRet)
								}
							case "REGULARITY":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryRegularityASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryRegularityDESC).Sort(arraySliceRet)
								}
							case "LASTRUNTIME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryLRTimeASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryLRTimeDESC).Sort(arraySliceRet)
								}
							case "RUNON":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(TaskSummaryRunOnASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(TaskSummaryRunOnDESC).Sort(arraySliceRet)
								}
							default:
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
							}
						}
					}

					staticInfo, err := json.Marshal(widgetInfo)
					if err != nil {
						return nil, err
					}

					return Relay.ConnectionFromArray(arraySliceRet, args, string(staticInfo)), nil
				}
				return nil, nil
			},
		},
	},
})

// TaskSummaryLRStatusDESC DESC sorting function for LastRunStatus column
func TaskSummaryLRStatusDESC(p1, p2 interface{}) bool {
	p1Status := getSortableStatus(p1.(apiModels.TaskSummaryData).LastRunStatus)
	p2Status := getSortableStatus(p2.(apiModels.TaskSummaryData).LastRunStatus)

	if p1Status == p2Status {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return p1Status < p2Status
}

// TaskSummaryLRStatusASC ASC sorting function for LastRunStatus column
func TaskSummaryLRStatusASC(p1, p2 interface{}) bool {
	p1Status := getSortableStatus(p1.(apiModels.TaskSummaryData).LastRunStatus)
	p2Status := getSortableStatus(p2.(apiModels.TaskSummaryData).LastRunStatus)

	if p1Status == p2Status {
		return TaskSummaryLRTimeDESC(p1, p2)
	}

	return p1Status > p2Status
}

// TaskSummaryNameASC ASC sorting function for Name column
func TaskSummaryNameASC(p1, p2 interface{}) bool {
	p1Name := p1.(apiModels.TaskSummaryData).Name
	p2Name := p2.(apiModels.TaskSummaryData).Name
	if p1Name == p2Name {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

// TaskSummaryNameDESC DESC sorting function for Name column
func TaskSummaryNameDESC(p1, p2 interface{}) bool {
	p1Name := p1.(apiModels.TaskSummaryData).Name
	p2Name := p2.(apiModels.TaskSummaryData).Name
	if p1Name == p2Name {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}

// TaskSummaryTypeASC ASC sorting function for Type column
func TaskSummaryTypeASC(p1, p2 interface{}) bool {
	p1Type := p1.(apiModels.TaskSummaryData).Type
	p2Type := p2.(apiModels.TaskSummaryData).Type
	if p1Type == p2Type {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return Relay.StringLessOp(p1Type, p2Type)
}

// TaskSummaryTypeDESC DESC sorting function for Type column
func TaskSummaryTypeDESC(p1, p2 interface{}) bool {
	p1Type := p1.(apiModels.TaskSummaryData).Type
	p2Type := p2.(apiModels.TaskSummaryData).Type
	if p1Type == p2Type {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return !Relay.StringLessOp(p1Type, p2Type)
}

// TaskSummaryLRTimeASC ASC sorting function for Last Run Time column
func TaskSummaryLRTimeASC(p1, p2 interface{}) bool {
	p1LRTime := p1.(apiModels.TaskSummaryData).LastRunTime.Truncate(time.Minute)
	p2LRTime := p2.(apiModels.TaskSummaryData).LastRunTime.Truncate(time.Minute)
	if p1LRTime == p2LRTime {
		return p1.(apiModels.TaskSummaryData).CreatedAt.Before(p2.(apiModels.TaskSummaryData).CreatedAt)
	}
	return p1LRTime.Before(p2LRTime)
}

// TaskSummaryLRTimeDESC DESC sorting function for Last Run Time column
func TaskSummaryLRTimeDESC(p1, p2 interface{}) bool {
	p1LRTime := p1.(apiModels.TaskSummaryData).LastRunTime.Truncate(time.Minute)
	p2LRTime := p2.(apiModels.TaskSummaryData).LastRunTime.Truncate(time.Minute)
	if p1LRTime == p2LRTime {
		return p1.(apiModels.TaskSummaryData).CreatedAt.After(p2.(apiModels.TaskSummaryData).CreatedAt)
	}
	return p1LRTime.After(p2LRTime)
}

// TaskSummaryInitiatedByASC ASC sorting function for InitiatedBy column
func TaskSummaryInitiatedByASC(p1, p2 interface{}) bool {
	p1InitiatedBy := p1.(apiModels.TaskSummaryData).InitiatedBy
	p2InitiatedBy := p2.(apiModels.TaskSummaryData).InitiatedBy
	if p1InitiatedBy == p2InitiatedBy {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return Relay.StringLessOp(p1InitiatedBy, p2InitiatedBy)
}

// TaskSummaryInitiatedByDESC DESC sorting function for InitiatedBy column
func TaskSummaryInitiatedByDESC(p1, p2 interface{}) bool {
	p1InitiatedBy := p1.(apiModels.TaskSummaryData).InitiatedBy
	p2InitiatedBy := p2.(apiModels.TaskSummaryData).InitiatedBy
	if p1InitiatedBy == p2InitiatedBy {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return !Relay.StringLessOp(p1InitiatedBy, p2InitiatedBy)
}

// TaskSummaryStatusASC ASC sorting function for Status column
func TaskSummaryStatusASC(p1, p2 interface{}) bool {
	p1Status := p1.(apiModels.TaskSummaryData).Status
	p2Status := p2.(apiModels.TaskSummaryData).Status
	if p1Status == p2Status {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return Relay.StringLessOp(p1Status, p2Status)
}

// TaskSummaryStatusDESC DESC sorting function for Status column
func TaskSummaryStatusDESC(p1, p2 interface{}) bool {
	p1Status := p1.(apiModels.TaskSummaryData).Status
	p2Status := p2.(apiModels.TaskSummaryData).Status
	if p1Status == p2Status {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return !Relay.StringLessOp(p1Status, p2Status)
}

// TaskSummaryRegularityASC ASC sorting function for Regularity column
func TaskSummaryRegularityASC(p1, p2 interface{}) bool {
	p1Schedule := p1.(apiModels.TaskSummaryData).Regularity
	p2Schedule := p2.(apiModels.TaskSummaryData).Regularity
	if p1Schedule == p2Schedule {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return Relay.StringLessOp(p1Schedule, p2Schedule)
}

// TaskSummaryRegularityDESC DESC sorting function for Regularity column
func TaskSummaryRegularityDESC(p1, p2 interface{}) bool {
	p1Schedule := p1.(apiModels.TaskSummaryData).Regularity
	p2Schedule := p2.(apiModels.TaskSummaryData).Regularity
	if p1Schedule == p2Schedule {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return !Relay.StringLessOp(p1Schedule, p2Schedule)
}

// TaskSummaryRunOnASC ASC sorting function for RunOn column
func TaskSummaryRunOnASC(p1, p2 interface{}) bool {
	p1RunOn := p1.(apiModels.TaskSummaryData).RunOn.TargetCount
	p2RunOn := p2.(apiModels.TaskSummaryData).RunOn.TargetCount
	if p1RunOn == p2RunOn {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return p1RunOn < p2RunOn
}

// TaskSummaryRunOnDESC DESC sorting function for RunOn column
func TaskSummaryRunOnDESC(p1, p2 interface{}) bool {
	p1RunOn := p1.(apiModels.TaskSummaryData).RunOn.TargetCount
	p2RunOn := p2.(apiModels.TaskSummaryData).RunOn.TargetCount
	if p1RunOn == p2RunOn {
		return TaskSummaryLRTimeDESC(p1, p2)
	}
	return p1RunOn > p2RunOn
}

func appendArg(connectionArgs graphql.FieldConfigArgument, argName string, argConfig graphql.ArgumentConfig) graphql.FieldConfigArgument {
	newConnectionArgs := make(map[string]*graphql.ArgumentConfig)

	for k, v := range connectionArgs {
		newConnectionArgs[k] = v

	}
	newConnectionArgs[argName] = &argConfig
	return newConnectionArgs
}
