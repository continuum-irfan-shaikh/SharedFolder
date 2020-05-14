package TasksAndSequencesDetailsSchema

import (
	"errors"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/TaskingSchema"
	"github.com/graphql-go/graphql"
)

//TasksAndSequencesDetails : TasksAndSequencesDetails structure
type TasksAndSequencesDetails struct {
	TaskID                    string `json:"taskId"` //uuid
	TaskName                  string `json:"taskName"`
	Description               string `json:"description"`
	Type                      string `json:"type"` //script or seq
	PartnerID                 string `json:"partnerId"`
	ManagedEndpointID         string `json:"managedEndpointId"` //uuid
	ExecutionID               string `json:"executionId"`       //uuid
	OriginID                  string `json:"originId"`          //uuid
	Regularity                string `json:"regularity"`
	InitiatedBy               string `json:"initiatedBy"`
	Status                    string `json:"status"`
	LastRunTime               string `json:"lastRunTime"` //"2018-07-02T12:46:29.889Z"
	NextRunTime               string `json:"nextRunTime"`
	LastRunStatus             string `json:"lastRunStatus"`
	LastRunStdOut             string `json:"lastRunStdOut"`
	LastRunStdErr             string `json:"lastRunStdErr"`
	ResultMessage             string `json:"resultMessage"`
	LastRunSequenceInstanceId string `json:"lastRunSequenceInstanceId"`
	DeviceCount               int    `json:"deviceCount"`
	CanBePostponed            bool   `json:"canBePostponed"`
	CanBeCanceled             bool   `json:"canBeCanceled"`
	ModifiedBy                string `json:"modifiedBy"`
}

const emptyTimeString = "0001-01-01T00:00:00Z"

//TasksAndSequencesDetailsType : TasksAndSequencesDetailsType GraphQL Schema
var TasksAndSequencesDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "details",
	Fields: graphql.Fields{
		"taskId": &graphql.Field{
			Type:        graphql.String,
			Description: "taskId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.TaskID, nil
				}
				return nil, nil
			},
		},
		"taskName": &graphql.Field{
			Type:        graphql.String,
			Description: "taskName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.TaskName, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "description",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return TaskingSchema.DefineTaskTypeView(CurData.Type), nil
				}
				return nil, nil
			},
		},
		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "partnerId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"managedEndpointId": &graphql.Field{
			Type:        graphql.String,
			Description: "managedEndpointId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.ManagedEndpointID, nil
				}
				return nil, nil
			},
		},
		"executionId": &graphql.Field{
			Type:        graphql.String,
			Description: "executionId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.ExecutionID, nil
				}
				return nil, nil
			},
		},
		"originId": &graphql.Field{
			Type:        graphql.String,
			Description: "originId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.OriginID, nil
				}
				return nil, nil
			},
		},
		"regularity": &graphql.Field{
			Type:        graphql.String,
			Description: "regularity",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.Regularity, nil
				}
				return nil, nil
			},
		},
		"initiatedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "initiatedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.InitiatedBy, nil
				}
				return nil, nil
			},
		},
		"modifiedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "modifiedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.ModifiedBy, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"canBePostponed": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "canBePostponed",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.CanBePostponed, nil
				}
				return nil, nil
			},
		},
		"canBeCanceled": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "canBeCanceled",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.CanBeCanceled, nil
				}
				return nil, nil
			},
		},
		"lastRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "lastRunTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.LastRunTime, nil
				}
				return nil, nil
			},
		},
		"nextRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "nextRunTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				CurData, ok := p.Source.(TasksAndSequencesDetails)
				if !ok {
					return nil, nil
				}

				if CurData.NextRunTime == emptyTimeString {
					return "", nil
				}
				return CurData.NextRunTime, nil
			},
		},
		"lastRunStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "lastRunStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.LastRunStatus, nil
				}
				return nil, nil
			},
		},
		"lastRunStdOut": &graphql.Field{
			Type:        graphql.String,
			Description: "lastRunStdOut",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.LastRunStdOut, nil
				}
				return nil, nil
			},
		},
		"lastRunStdErr": &graphql.Field{
			Type:        graphql.String,
			Description: "lastRunStdErr",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.LastRunStdErr, nil
				}
				return nil, nil
			},
		},
		"resultMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "resultMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.ResultMessage, nil
				}
				return nil, nil
			},
		},
		"lastRunSequenceInstanceId": &graphql.Field{
			Type:        graphql.String,
			Description: "lastRunSequenceInstanceId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.LastRunSequenceInstanceId, nil
				}
				return nil, nil
			},
		},

		"deviceCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "deviceCount",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesDetails); ok {
					return CurData.DeviceCount, nil
				}
				return nil, nil
			},
		},
	},
},
)

//TasksAndSequencesDetailsList : TasksAndSequencesDetailsList connection definition structure
var TasksAndSequencesDetailsList = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TasksAndSequencesDetailsList",
	NodeType: TasksAndSequencesDetailsType,
})

//TasksAndSequencesDetailsListData : List of TasksAndSequencesDetails Structures
type TasksAndSequencesDetailsListData struct {
	TasksAndSequencesDetailsList []TasksAndSequencesDetails `json:"TasksAndSequencesSelectorList"`
}

//TasksAndSequencesDetailsListType : TasksAndSequencesDetailsListType GraphQL Schema
var TasksAndSequencesDetailsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TasksAndSequencesDetailsList",
	Fields: graphql.Fields{
		"TasksAndSequencesDetailsList": &graphql.Field{
			Type:        TasksAndSequencesDetailsList.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of TasksAndSequencesDetails",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TasksAndSequencesDetailsListData); ok {
					var resSlice []interface{}
					for _, val := range CurData.TasksAndSequencesDetailsList {
						resSlice = append(resSlice, val)
					}

					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")

						for _, subQ := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQ)
							switch strings.ToUpper(Column) {
							case "NAME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SingleTaskNameASC).Sort(resSlice)
								} else {
									Relay.SortBy(SingleTaskNameDESC).Sort(resSlice)
								}
							case "SCHEDULE":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(ScheduleASC).Sort(resSlice)
								} else {
									Relay.SortBy(ScheduleDESC).Sort(resSlice)
								}
							case "STATUS":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(StatusASC).Sort(resSlice)
								} else {
									Relay.SortBy(StatusDESC).Sort(resSlice)
								}
							case "LASTRUNTIME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(LastRunTimeASC).Sort(resSlice)
								} else {
									Relay.SortBy(LastRunTimeDESC).Sort(resSlice)
								}
							case "LASTRUNSTATUS":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(LastRunStatusASC).Sort(resSlice)
								} else {
									Relay.SortBy(LastRunStatusDESC).Sort(resSlice)
								}
							default:
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
							}

						} //end for loop

					} //end "if args.Sort"
					return Relay.ConnectionFromArray(resSlice, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})

// SingleTaskNameASC ASC sorting function for TaskName column
func SingleTaskNameASC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).TaskName
	p2Name := p2.(TasksAndSequencesDetails).TaskName
	if p1Name == p2Name {
		return Relay.StringLessOp(p1.(TasksAndSequencesDetails).Description, p2.(TasksAndSequencesDetails).Description)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

// SingleTaskNameDESC DESC sorting function for TaskName column
func SingleTaskNameDESC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).TaskName
	p2Name := p2.(TasksAndSequencesDetails).TaskName
	if p1Name == p2Name {
		return !Relay.StringLessOp(p1.(TasksAndSequencesDetails).Description, p2.(TasksAndSequencesDetails).Description)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}

// ScheduleASC ASC sorting function for Schedule column
func ScheduleASC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).Regularity
	p2Name := p2.(TasksAndSequencesDetails).Regularity
	if p1Name == p2Name {
		return Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

// ScheduleDESC DESC sorting function for Schedule column
func ScheduleDESC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).Regularity
	p2Name := p2.(TasksAndSequencesDetails).Regularity
	if p1Name == p2Name {
		return !Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}

// StatusASC ASC sorting function for Status column
func StatusASC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).Status
	p2Name := p2.(TasksAndSequencesDetails).Status
	if p1Name == p2Name {
		return Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

// StatusDESC DESC sorting function for Status column
func StatusDESC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).Status
	p2Name := p2.(TasksAndSequencesDetails).Status
	if p1Name == p2Name {
		return !Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}

// LastRunTimeASC ASC sorting function for LastRunTime column
func LastRunTimeASC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).LastRunTime
	p2Name := p2.(TasksAndSequencesDetails).LastRunTime
	if p1Name == p2Name {
		return Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

// LastRunTimeDESC DESC sorting function for LastRunTime column
func LastRunTimeDESC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).LastRunTime
	p2Name := p2.(TasksAndSequencesDetails).LastRunTime
	if p1Name == p2Name {
		return !Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}

// LastRunStatusASC ASC sorting function for LastRunStatus column
func LastRunStatusASC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).LastRunStatus
	p2Name := p2.(TasksAndSequencesDetails).LastRunStatus
	if p1Name == p2Name {
		return Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

// LastRunStatusDESC DESC sorting function for LastRunStatus column
func LastRunStatusDESC(p1, p2 interface{}) bool {
	p1Name := p1.(TasksAndSequencesDetails).LastRunStatus
	p2Name := p2.(TasksAndSequencesDetails).LastRunStatus
	if p1Name == p2Name {
		return !Relay.StringLessOp(p1.(TasksAndSequencesDetails).InitiatedBy, p2.(TasksAndSequencesDetails).InitiatedBy)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}
