package TaskingSchema

import (
	"errors"
	"reflect"
	"strings"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//SingleTaskSummaryPageData : SingleTaskSummaryPageData Structure
type SingleTaskSummaryPageData struct {
	SingleTaskData         apiModels.TaskSummaryData   `json:"taskSummary"`
	SingleTaskInstanceList []SingleTaskSummaryInstance `json:"instanceSummary"`
}

//SingleTaskSummaryInstance : SingleTaskSummaryInstance Structure
type SingleTaskSummaryInstance struct {
	ID         string                      `json:"taskInstanceID"`
	RunTime    time.Time                   `json:"runTime"`
	RunStatus  apiModels.LastRunStatusData `json:"runStatus"`
	DeviceList []SingleTaskDevice          `json:"targetSummaries"`
}

//SingleTaskDevice : SingleTaskDevice Structure
type SingleTaskDevice struct {
	EndpointID         string    `json:"endpointID"`
	RunStatus          string    `json:"runStatus"`
	StatusDetails      string    `json:"statusDetails"`
	Output             string    `json:"output"`
	SpecificInstanceID string    `json:"specificInstanceID"`
	OriginID           string    `json:"originID"`
	NextRunTime        time.Time `json:"nextRunTime"`
	CanBePostponed     bool      `json:"canBePostponed"`
	CanBeCanceled      bool      `json:"canBeCanceled"`
	PostponedTime      time.Time `json:"postponedTime"`
	LastRunTime        time.Time `json:"lastRunTime"`
}

//SingleTaskDeviceType : SingleTaskDeviceType GraphQL Schema
var SingleTaskDeviceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SingleTaskDeviceType",
	Fields: graphql.Fields{
		"endpointID": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDevice); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},
		"runStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's RunStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDevice); ok {
					return CurData.RunStatus, nil
				}
				return nil, nil
			},
		},
		"statusDetails": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's StatusDetails",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDevice); ok {
					return CurData.StatusDetails, nil
				}
				return nil, nil
			},
		},
		"output": &graphql.Field{
			Type:        graphql.String,
			Description: "SingleTaskDevice's Output",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDevice); ok {
					return CurData.Output, nil
				}
				return nil, nil
			},
		},
		"specificInstanceID": &graphql.Field{
			Type:        graphql.String,
			Description: "specificInstanceID is ID for MS specific instance for underlining entity",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDevice); ok {
					return CurData.SpecificInstanceID, nil
				}
				return nil, nil
			},
		},
		"originID": &graphql.Field{
			Type:        graphql.String,
			Description: "originID is ID for MS specific entity ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDevice); ok {
					return CurData.OriginID, nil
				}
				return nil, nil
			},
		},
		"nextRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "nextRunTime is time when scheduled task will be running the next time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskDevice); ok {
					return CurData.NextRunTime, nil
				}
				return nil, nil
			},
		},
	},
})

//SingleTaskDeviceListTypeConnDef : SingleTaskDeviceListTypeConnDef structure
var SingleTaskDeviceListTypeConnDef = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "SingleTaskDeviceListTypeConnDef",
	NodeType: SingleTaskDeviceType,
})

//SingleTaskSummaryInstanceType : SingleTaskSummaryInstanceType GraphQL Schema
var SingleTaskSummaryInstanceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SingleTaskSummaryInstanceType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "Task Instance's ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskSummaryInstance); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"runTime": &graphql.Field{
			Type:        graphql.String,
			Description: "Task Instance's run time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskSummaryInstance); ok {
					return CurData.RunTime.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"runStatus": &graphql.Field{
			Type:        LastRunStatusType,
			Description: "Task Instance's global run status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskSummaryInstance); ok {
					return CurData.RunStatus, nil
				}
				return nil, nil
			},
		},
		"deviceList": &graphql.Field{
			Type:        SingleTaskDeviceListTypeConnDef.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SingleTaskSummaryInstance); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.DeviceList {
						arraySliceRet = append(arraySliceRet, CurData.DeviceList[ind])
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})

//SingleTaskInstanceListTypeConnDef : SingleTaskInstanceListTypeConnDef structure
var SingleTaskInstanceListTypeConnDef = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "SingleTaskInstanceListTypeConnDef",
	NodeType: SingleTaskSummaryInstanceType,
})

//SingleTaskSummaryPageType : SingleTaskSummaryPageType GraphQL Schema
var SingleTaskSummaryPageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SingleTaskSummaryPageType",
	Fields: graphql.Fields{
		"singleTaskData": &graphql.Field{
			Type:        TaskingSummaryPageType,
			Description: "Task summary data",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SingleTaskSummaryPageData); ok {
					return CurData.SingleTaskData, nil
				}
				return nil, nil
			},
		},
		"singleTaskInstanceList": &graphql.Field{
			Type:        SingleTaskInstanceListTypeConnDef.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of Task Instances",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SingleTaskSummaryPageData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.SingleTaskInstanceList {
						arraySliceRet = append(arraySliceRet, CurData.SingleTaskInstanceList[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&SingleTaskSummaryInstance{}))
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
							case "RUNTIME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SingleTaskSummaryRTimeASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(SingleTaskSummaryRTimeDESC).Sort(arraySliceRet)
								}
							case "RUNSTATUS":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SingleTaskSummaryRStatusASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(SingleTaskSummaryRStatusDESC).Sort(arraySliceRet)
								}
							default:
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
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

// SingleTaskSummaryRTimeASC ASC sorting function for RunTime column
func SingleTaskSummaryRTimeASC(p1, p2 interface{}) bool {
	return p1.(SingleTaskSummaryInstance).RunTime.Before(p2.(SingleTaskSummaryInstance).RunTime)
}

// SingleTaskSummaryRTimeDESC DESC sorting function for RunTime column
func SingleTaskSummaryRTimeDESC(p1, p2 interface{}) bool {
	return p1.(SingleTaskSummaryInstance).RunTime.After(p2.(SingleTaskSummaryInstance).RunTime)
}

// SingleTaskSummaryRStatusASC ASC sorting function for RunStatus column
func SingleTaskSummaryRStatusASC(p1, p2 interface{}) bool {
	p1Status := getSortableStatus(p1.(SingleTaskSummaryInstance).RunStatus)
	p2Status := getSortableStatus(p2.(SingleTaskSummaryInstance).RunStatus)
	if p1Status == p2Status {
		return SingleTaskSummaryRTimeASC(p1, p2)
	}
	return p1Status > p2Status
}

// SingleTaskSummaryRStatusDESC DESC sorting function for RunStatus column
func SingleTaskSummaryRStatusDESC(p1, p2 interface{}) bool {
	p1Status := getSortableStatus(p1.(SingleTaskSummaryInstance).RunStatus)
	p2Status := getSortableStatus(p2.(SingleTaskSummaryInstance).RunStatus)
	if p1Status == p2Status {
		return SingleTaskSummaryRTimeDESC(p1, p2)
	}
	return p1Status < p2Status
}
