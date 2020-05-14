package TaskingSchema

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//SingleTaskInstanceSummaryData : SingleTaskInstanceSummaryData Structure
type SingleTaskInstanceSummaryData struct {
	SingleTaskData apiModels.TaskSummaryData  `json:"taskSummary"`
	DeviceList     []SingleTaskDeviceExtended `json:"deviceList"`
}

//SingleTaskDeviceExtended : SingleTaskDeviceExtended Structure
type SingleTaskDeviceExtended struct {
	EndpointID            string                 `json:"endpointID"`
	MachineName           string                 `json:"machineName"`
	FriendlyName          string                 `json:"friendlyName"`
	SiteName              string                 `json:"siteName"`
	RegID                 string                 `json:"regID"`
	RunStatus             string                 `json:"runStatus"`
	StatusDetails         string                 `json:"statusDetails"`
	Output                string                 `json:"output"`
	SequenceTaskInstances []SequenceInstanceTask `json:"sequenceTaskInstances"`
	Error                 string                 `json:"error"`
	// LastRunTime and InitiatedBy fields are for exporting data
	LastRunTime        time.Time `json:"lastRunTime"`
	LastRunBy          string    `json:"lastRunBy"`
	InitiatedBy        string    `json:"initiatedBy"`
	SpecificInstanceID string    `json:"specificInstanceID"`
	OriginID           string    `json:"originID"`
	//NextRunTime are only for scheduled tasks
	NextRunTime    time.Time `json:"nextRunTime"`
	PostponedTime  time.Time `json:"postponedTime"`
	ModifiedBy     string    `json:"modifiedBy"`
	ModifiedAt     time.Time `json:"modifiedAt"`
	CreatedAt      time.Time `json:"createdAt"`
	CanBePostponed bool      `json:"canBePostponed"`
	CanBeCanceled  bool      `json:"canBeCanceled"`
}

//SequenceInstanceTask - SequenceInstanceTask structure
type SequenceInstanceTask struct {
	ID                 string    `json:"id"                 cql:"id"`
	SequenceID         string    `json:"sequenceId"         cql:"sequence_id"`
	SequenceInstanceID string    `json:"sequenceInstanceId" cql:"sequence_instance_id"`
	SequenceTaskID     string    `json:"sequenceTaskId"     cql:"sequence_task_id"`
	EndpointID         string    `json:"endpointId"         cql:"endpoint_id"`
	OriginID           string    `json:"originId"           cql:"origin_id"`
	TaskingTaskID      string    `json:"taskingTaskId"      cql:"tasking_task_id"`
	Type               string    `json:"type"               cql:"type"`
	Name               string    `json:"name"`
	StdOut             string    `json:"stdOut"             cql:"std_out"`
	StdErr             string    `json:"stdErr"             cql:"std_err"`
	ResultMessage      string    `json:"resultMessage"      cql:"result_message"`
	StartedAt          time.Time `json:"startedAt"          cql:"started_at"`
	Status             string    `json:"status"`
	ExitOnFailure      bool      `json:"exitOnFailure"      cql:"exit_on_failure"`
	CreatedAt          time.Time `json:"createdAt"          cql:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt"          cql:"updated_at"`
}

//SingleTaskDeviceExtendedType : SingleTaskDeviceExtendedType GraphQL Schema
var SingleTaskDeviceExtendedType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SingleTaskDeviceExtendedType",
	Fields: graphql.Fields{
		"deviceId": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},
		"machineName": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's MachineName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.MachineName, nil
				}
				return nil, nil
			},
		},
		"friendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's FriendlyName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.FriendlyName, nil
				}
				return nil, nil
			},
		},
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's SiteName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
		"regID": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's Registration ID of Old Agent",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},
		"runStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's RunStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.RunStatus, nil
				}
				return nil, nil
			},
		},
		"statusDetails": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's StatusDetails",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.StatusDetails, nil
				}
				return nil, nil
			},
		},
		"canBePostponed": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "SingleTaskDevice's canBePostponed",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.CanBePostponed, nil
				}
				return nil, nil
			},
		},
		"canBeCanceled": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "SingleTaskDevice's canBeCanceled",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.CanBeCanceled, nil
				}
				return nil, nil
			},
		},
		"output": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's Output",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.Output, nil
				}
				return nil, nil
			},
		},
		"error": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's Output",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.Error, nil
				}
				return nil, nil
			},
		},
		"sequenceTaskInstances": &graphql.Field{
			Type:        graphql.NewList(SequenceInstanceTaskType),
			Description: "SingleTaskDevice's Output",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.SequenceTaskInstances, nil
				}
				return nil, nil
			},
		},
		"lastRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "Task's last run time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					if CurData.LastRunTime.IsZero() {
						return "", nil
					}

					return CurData.LastRunTime.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"lastRunBy": &graphql.Field{
			Type:        graphql.String,
			Description: "User who run task last time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					status, err := apiModels.TaskInstanceStatusFromText(CurData.RunStatus)
					if status == apiModels.TaskInstancePending || status == apiModels.TaskInstanceScheduled || err != nil {
						// in case of Scheduled task we should not display lastRunBy - RMM-36743
						return "", nil
					}

					if CurData.ModifiedBy == "" || CurData.ModifiedAt.After(CurData.LastRunTime) {
						return CurData.InitiatedBy, nil
					}

					return CurData.ModifiedBy, nil
				}
				return nil, nil
			},
		},
		"initiatedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "Task's author",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.InitiatedBy, nil
				}
				return nil, nil
			},
		},
		"specificInstanceID": &graphql.Field{
			Type:        graphql.String,
			Description: "specificInstanceID is ID for MS specific instance for underlining entity",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.SpecificInstanceID, nil
				}
				return nil, nil
			},
		},
		"originID": &graphql.Field{
			Type:        graphql.String,
			Description: "originID is ID for MS specific entity ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.OriginID, nil
				}
				return nil, nil
			},
		},
		"nextRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "nextRunTime is time when scheduled task will be running the next time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					if CurData.NextRunTime.IsZero() {
						return "", nil
					}

					if CurData.PostponedTime.IsZero() {
						return CurData.NextRunTime.Format(time.RFC3339), nil
					}

					return CurData.PostponedTime.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"modifiedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "modifiedBy represents person who modified task",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.ModifiedBy, nil
				}
				return nil, nil
			},
		},
		"modifiedAt": &graphql.Field{
			Type:        graphql.String,
			Description: "modifiedAt is time when the task was modified",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					if CurData.ModifiedAt.IsZero() {
						return "", nil
					}

					return CurData.ModifiedAt.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.String,
			Description: "createdAt is time when the task was created",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDeviceExtended); ok {
					return CurData.CreatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

//SequenceInstanceTaskType defines the graph ql data type for
var SequenceInstanceTaskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SequenceInstanceTaskType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"sequenceId": &graphql.Field{
			Type:        graphql.String,
			Description: "sequenceId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.SequenceID, nil
				}
				return nil, nil
			},
		},
		"sequenceInstanceId": &graphql.Field{
			Type:        graphql.String,
			Description: "sequenceInstanceId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.SequenceInstanceID, nil
				}
				return nil, nil
			},
		},
		"sequenceTaskId": &graphql.Field{
			Type:        graphql.String,
			Description: "sequenceTaskId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.SequenceTaskID, nil
				}
				return nil, nil
			},
		},
		"endpointID": &graphql.Field{
			Type:        graphql.String,
			Description: "endpointID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},
		"originId": &graphql.Field{
			Type:        graphql.String,
			Description: "originId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.OriginID, nil
				}
				return nil, nil
			},
		},
		"taskingTaskId": &graphql.Field{
			Type:        graphql.String,
			Description: "taskingTaskId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.TaskingTaskID, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"stdOut": &graphql.Field{
			Type:        graphql.String,
			Description: "stdOut",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.StdOut, nil
				}
				return nil, nil
			},
		},
		"stdErr": &graphql.Field{
			Type:        graphql.String,
			Description: "stdErr",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.StdErr, nil
				}
				return nil, nil
			},
		},
		"resultMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "resultMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.ResultMessage, nil
				}
				return nil, nil
			},
		},
		"startedAt": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "startedAt",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.StartedAt, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"exitOnFailure": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "exitOnFailure",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.ExitOnFailure, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "createdAt",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type:        graphql.String,
			Description: "updatedAt",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SequenceInstanceTask); ok {
					return CurData.UpdatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

//DeviceListConnDef : DeviceListConnDef structure
var DeviceListConnDef = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "DeviceListConnDef",
	NodeType: SingleTaskDeviceExtendedType,
})

//SingleTaskInstanceSummaryDataType : SingleTaskInstanceSummaryDataType GraphQL Schema
var SingleTaskInstanceSummaryDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SingleTaskInstanceSummaryDataType",
	Fields: graphql.Fields{
		"deviceList": &graphql.Field{
			Type:        DeviceListConnDef.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of Instance Devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SingleTaskInstanceSummaryData); ok {

					type StatInfo struct {
						Statuses           map[string]int `json:"statuses"`
						CreatedAt          time.Time      `json:"createdAt"`
						InitiatedBy        string         `json:"initiatedBy"`
						ModifiedAt         string         `json:"modifiedAt"`
						ModifiedBy         string         `json:"modifiedBy"`
						NearestNextRunTime time.Time      `json:"nearestNextRunTime"`
						TriggerType        string         `json:"triggerType"`
					}

					const emptyTimeString = "0001-01-01T00:00:00Z"
					triggerType := ""
					if len(CurData.SingleTaskData.TriggerTypes) > 0 {
						triggerType = CurData.SingleTaskData.TriggerTypes[0]
					}
					var (
						statuses      = make(map[string]int)
						arraySliceRet []interface{}

						statInfo = StatInfo{
							Statuses:           statuses,
							CreatedAt:          CurData.SingleTaskData.CreatedAt,
							InitiatedBy:        CurData.SingleTaskData.InitiatedBy,
							ModifiedAt:         CurData.SingleTaskData.ModifiedAt.Format(time.RFC3339),
							ModifiedBy:         CurData.SingleTaskData.ModifiedBy,
							NearestNextRunTime: CurData.SingleTaskData.NearestNextRunTime,
							TriggerType:        triggerType,
						}
					)

					for ind := range CurData.DeviceList {
						statuses[CurData.DeviceList[ind].RunStatus] = statuses[CurData.DeviceList[ind].RunStatus] + 1
						arraySliceRet = append(arraySliceRet, CurData.DeviceList[ind])
					}
					statInfo.Statuses = statuses

					if statInfo.ModifiedAt == emptyTimeString {
						statInfo.ModifiedAt = ""
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&SingleTaskDeviceExtended{}))
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
							case "MACHINENAME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SingleTaskDeviceMNameASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(SingleTaskDeviceMNameDESC).Sort(arraySliceRet)
								}
							case "FRIENDLYNAME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SingleTaskDeviceFNameASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(SingleTaskDeviceFNameDESC).Sort(arraySliceRet)
								}
							case "SITENAME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SingleTaskDeviceSiteNameASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(SingleTaskDeviceSiteNameDESC).Sort(arraySliceRet)
								}
							case "RUNSTATUS":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SingleTaskDeviceRStatusASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(SingleTaskDeviceRStatusDESC).Sort(arraySliceRet)
								}
							default:
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
							}
						}
					}

					staticInfo, err := json.Marshal(statInfo)
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

// SingleTaskDeviceMNameASC function sorts data by MachineName field in ascending order
func SingleTaskDeviceMNameASC(p1, p2 interface{}) bool {
	p1MName := p1.(SingleTaskDeviceExtended).MachineName
	p2MName := p2.(SingleTaskDeviceExtended).MachineName
	return Relay.StringLessOp(p1MName, p2MName)
}

// SingleTaskDeviceMNameDESC function sorts data by MachineName field in descending order
func SingleTaskDeviceMNameDESC(p1, p2 interface{}) bool {
	p1MName := p1.(SingleTaskDeviceExtended).MachineName
	p2MName := p2.(SingleTaskDeviceExtended).MachineName
	return !Relay.StringLessOp(p1MName, p2MName)
}

// SingleTaskDeviceFNameASC function sorts data by FriendlyName field in ascending order
func SingleTaskDeviceFNameASC(p1, p2 interface{}) bool {
	p1FName := p1.(SingleTaskDeviceExtended).FriendlyName
	p2FName := p2.(SingleTaskDeviceExtended).FriendlyName
	return Relay.StringLessOp(p1FName, p2FName)
}

// SingleTaskDeviceFNameDESC function sorts data by FriendlyName field in descending order
func SingleTaskDeviceFNameDESC(p1, p2 interface{}) bool {
	p1FName := p1.(SingleTaskDeviceExtended).FriendlyName
	p2FName := p2.(SingleTaskDeviceExtended).FriendlyName
	return !Relay.StringLessOp(p1FName, p2FName)
}

// SingleTaskDeviceSiteNameASC function sorts data by SiteName field in ascending order
// sorting layers: SiteName > MName(ASC)
// function chain according to sorting layers: SingleTaskDeviceSiteNameASC > SingleTaskDeviceMNameASC
func SingleTaskDeviceSiteNameASC(p1, p2 interface{}) bool {
	p1SiteName := p1.(SingleTaskDeviceExtended).SiteName
	p2SiteName := p2.(SingleTaskDeviceExtended).SiteName
	if p1SiteName == p2SiteName {
		return SingleTaskDeviceMNameASC(p1, p2)
	}
	return Relay.StringLessOp(p1SiteName, p2SiteName)
}

// SingleTaskDeviceSiteNameDESC function sorts data by SiteName field in descending order
// sorting layers: SiteName > MName(ASC)
// function chain according to sorting layers: SingleTaskDeviceSiteNameDESC > SingleTaskDeviceMNameASC
func SingleTaskDeviceSiteNameDESC(p1, p2 interface{}) bool {
	p1SiteName := p1.(SingleTaskDeviceExtended).SiteName
	p2SiteName := p2.(SingleTaskDeviceExtended).SiteName
	if p1SiteName == p2SiteName {
		return SingleTaskDeviceMNameASC(p1, p2)
	}
	return !Relay.StringLessOp(p1SiteName, p2SiteName)
}

// SingleTaskDeviceRStatusASC function sorts data by RunStatus field in ascending order
// sorting layers: RStatus > SiteName(ASC) > MName(ASC)
// function chain according to sorting layers: SingleTaskDeviceRStatusASC > SingleTaskDeviceSiteNameASC > SingleTaskDeviceMNameASC
func SingleTaskDeviceRStatusASC(p1, p2 interface{}) bool {
	p1RStatus := getSortableLastRunStatus(p1.(SingleTaskDeviceExtended).RunStatus)
	p2RStatus := getSortableLastRunStatus(p2.(SingleTaskDeviceExtended).RunStatus)
	if p1RStatus == p2RStatus {
		return SingleTaskDeviceSiteNameASC(p1, p2)
	}
	return p1RStatus < p2RStatus
}

// SingleTaskDeviceRStatusDESC function sorts data by RunStatus field in descending order
// sorting layers: RStatus > SiteName(ASC) > MName(ASC)
// function chain according to sorting layers: SingleTaskDeviceRStatusDESC > SingleTaskDeviceSiteNameASC > SingleTaskDeviceMNameASC
func SingleTaskDeviceRStatusDESC(p1, p2 interface{}) bool {
	p1RStatus := getSortableLastRunStatus(p1.(SingleTaskDeviceExtended).RunStatus)
	p2RStatus := getSortableLastRunStatus(p2.(SingleTaskDeviceExtended).RunStatus)
	if p1RStatus == p2RStatus {
		return SingleTaskDeviceSiteNameASC(p1, p2)
	}
	return p1RStatus > p2RStatus
}
