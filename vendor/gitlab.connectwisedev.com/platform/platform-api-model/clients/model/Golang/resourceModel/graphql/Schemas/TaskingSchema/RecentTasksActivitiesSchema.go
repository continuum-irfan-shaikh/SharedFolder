package TaskingSchema

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/gocql/gocql"
	"github.com/graphql-go/graphql"
)

const (
	//ScriptTaskType - Task type - "Script"
	ScriptTaskType = "script"
	//SequenceTaskType - Task type - "Sequence"
	SequenceTaskType = "sequence"
	first            = 1
	second           = 2
	undefined        = -1
)

//RecentTasksActivity : RecentTasksActivity Structure TODO: will be moved to platform-api-model
type RecentTasksActivity struct {
	Task           Task                    `json:"task"`
	TaskInstance   TaskInstance            `json:"taskInstance"`
	Statuses       TaskInstanceStatusCount `json:"statuses"`
	LastRunStatus  string                  `json:"lastRunStatus"`
	CanBePostponed bool                    `json:"canBePostponed"`
	CanBeCanceled  bool                    `json:"canBeCanceled"`
}

//Task struct TODO: will be moved to platform-api-model
type Task struct {
	ID               gocql.UUID                `json:"id"`
	Name             string                    `json:"name"`
	Description      string                    `json:"description"`
	CreatedAt        time.Time                 `json:"createdAt"`
	CreatedBy        string                    `json:"createdBy"`
	ModifiedBy       string                    `json:"modifiedBy"`
	ModifiedAt       time.Time                 `json:"modifiedAt"`
	RunTimeUTC       time.Time                 `json:"nextRunTime"`
	PostponedTime    time.Time                 `json:"postponedTime"`
	PartnerID        string                    `json:"partnerId"`
	OriginID         gocql.UUID                `json:"originId"` // script or patch ID
	Trigger          string                    `json:"trigger"`
	Type             string                    `json:"type"`
	Parameters       string                    `json:"parameters"`
	ExternalTask     bool                      `json:"externalTask"`
	ResultWebhook    string                    `json:"resultWebhook"`
	Targets          Target                    `json:"targets"`
	ManagedEndpoints []ManagedEndpointDetailed `json:"managedEndpoints"`
	Schedule         apiModels.Schedule        `json:"schedule"`
	DefinitionID     gocql.UUID                `json:"definitionID"`
}

//TaskInstance struct TODO: will be moved to platform-api-model
type TaskInstance struct {
	PartnerID   string                                      `json:"partnerId"`
	ID          gocql.UUID                                  `json:"id"`
	TaskID      gocql.UUID                                  `json:"taskId"`
	OriginID    gocql.UUID                                  `json:"originId"` //ScriptID or PatchID
	StartedAt   time.Time                                   `json:"startedAt"`
	LastRunTime time.Time                                   `json:"lastRunTime"`
	Statuses    map[gocql.UUID]apiModels.TaskInstanceStatus `json:"statuses"`
	TriggeredBy string                                      `json:"triggeredBy"`
}

//TaskInstanceStatusCount struct TODO: will be moved to platform-api-model
type TaskInstanceStatusCount struct {
	TaskInstanceID gocql.UUID
	SuccessCount   int
	FailureCount   int
}

//ManagedEndpointDetailed struct TODO: will be moved to platform-api-model
type ManagedEndpointDetailed struct {
	apiModels.ManagedEndpoint
	Location string `json:"location"`
	State    string `json:"state"`
}

//Target struct TODO: will be moved to platform-api-model
type Target struct {
	IDs  []string `json:"ids"`
	Type string   `json:"type"`
}

//TaskFilters struct that represents filter applied by user on RecentTasksActivityList
type TaskFilters struct {
	ByUser        []string `json:"byUser"`
	LastRunStatus []string `json:"status"`
	Type          []string `json:"type"`
}

//FilterInfo struct to be used to represent type of filtering
type FilterInfo struct {
	ByUser        map[string]int `json:"byUser"`
	LastRunStatus map[string]int `json:"status"`
	Type          map[string]int `json:"type"`
}

// RecentTasksActivityType : RecentTasksActivityType GraphQL Schema
// Example: {{graphQL}}/GraphQL/?{RecentTasksActivity{RecentTasksActivityList{edges{cursor,node{name,type,lastRunTime,createdBy,deviceCount,successCount,failureCount,lastRunStatus}}pageInfo{startCursor,endCursor,hasNextPage,hasPreviousPage,totalCount}}}}
var RecentTasksActivityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RecentTasksActivityType",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Task name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.Name, nil
				}
				return nil, nil
			},
		},
		"modifiedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "modifiedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.ModifiedBy, nil
				}
				return nil, nil
			},
		},
		"modifiedAt": &graphql.Field{
			Type:        graphql.String,
			Description: "modifiedAt",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					if CurData.Task.ModifiedAt.IsZero() {
						return "", nil
					}

					return CurData.Task.ModifiedAt.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"taskID": &graphql.Field{
			Type:        graphql.String,
			Description: "Task ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.ID, nil
				}
				return nil, nil
			},
		},
		"instanceID": &graphql.Field{
			Type:        graphql.String,
			Description: "TaskInstance ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.TaskInstance.ID, nil
				}
				return nil, nil
			},
		},
		"originID": &graphql.Field{
			Type:        graphql.String,
			Description: "originID is ID for MS specific entity ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.TaskInstance.OriginID, nil
				}
				return nil, nil
			},
		},
		"definitionID": &graphql.Field{
			Type:        graphql.String,
			Description: "definitionID is ID for MS specific entity ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.DefinitionID, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Task type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.Type, nil
				}
				return nil, nil
			},
		},
		"lastRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "LastRunTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					if CurData.TaskInstance.LastRunTime.IsZero() {
						return "", nil
					}

					return CurData.TaskInstance.LastRunTime.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.String,
			Description: "Creation task time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.CreatedAt.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"nextRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "nextRunTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					if CurData.Task.RunTimeUTC.IsZero() {
						return "", nil
					}

					if CurData.Task.PostponedTime.IsZero() {
						return CurData.Task.RunTimeUTC.Format(time.RFC3339), nil
					}

					return CurData.Task.PostponedTime.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"createdBy": &graphql.Field{
			Type:        graphql.String,
			Description: "Task is CreatedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.CreatedBy, nil
				}
				return nil, nil
			},
		},
		"deviceCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return len(CurData.TaskInstance.Statuses), nil
				}
				return nil, nil
			},
		},
		"successCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of Successful executions",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Statuses.SuccessCount, nil
				}
				return nil, nil
			},
		},
		"failureCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of Failed executions",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Statuses.FailureCount, nil
				}
				return nil, nil
			},
		},
		"lastRunStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "LastRunStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.LastRunStatus, nil
				}
				return nil, nil
			},
		},
		"canBePostponed": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "This field tells us if the task can be postponed or not (if there is task with pending device)",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.CanBePostponed, nil
				}
				return nil, nil
			},
		},
		"canBeCanceled": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "This field tells us if the task can be canceled or not",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.CanBeCanceled, nil
				}
				return nil, nil
			},
		},
		"taskRegularity": &graphql.Field{
			Type:        graphql.String,
			Description: "taskRegularity",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.Task.Schedule.Regularity.String(), nil
				}
				return nil, nil
			},
		},
		"triggerType": &graphql.Field{
			Type:        graphql.String,
			Description: "triggerType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					return CurData.TaskInstance.TriggeredBy, nil
				}
				return nil, nil
			},
		},
		"statuses": &graphql.Field{
			Type:        graphql.String,
			Description: "statuses",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RecentTasksActivity); ok {
					statuses := make(map[string]int)
					for _, status := range CurData.TaskInstance.Statuses {
						statusStr, err := apiModels.TaskInstanceStatusText(status)
						if err != nil {
							return nil, err
						}
						statuses[statusStr] = statuses[statusStr] + 1
					}

					statusesInfo, err := json.Marshal(statuses)
					if err != nil {
						return nil, err
					}

					return string(statusesInfo), nil
				}
				return nil, nil
			},
		},
	},
})

// RecentTasksActivityList : RecentTasksActivityList connection definition structure
var RecentTasksActivityList = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "RecentTasksActivityList",
	NodeType: RecentTasksActivityType,
})

//RecentTasksActivityListType : RecentTasksActivityListType GraphQL Schema
var RecentTasksActivityListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RecentTasksActivityList",
	Fields: graphql.Fields{
		"RecentTasksActivityList": &graphql.Field{
			Type:        RecentTasksActivityList.ConnectionType,
			Args:        appendArg(Relay.ConnectionArgs, "filteringStr", graphql.ArgumentConfig{Type: graphql.String}),
			Description: "List of RecentTasksActivities",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				filteringStr := interfaceToString(p.Args["filteringStr"])
				var filters TaskFilters
				var isFilteringExist bool

				if filteringStr != "" {
					jsonFilteringStr := []byte(filteringStr)
					err := json.Unmarshal(jsonFilteringStr, &filters)
					if err != nil {
						return nil, err
					}
					isFilteringExist = true
				}

				if CurData, ok := p.Source.([]RecentTasksActivity); ok {
					var resSlice []interface{}

					var filterInfo = FilterInfo{
						ByUser:        make(map[string]int),
						LastRunStatus: make(map[string]int),
						Type:          make(map[string]int),
					}

					for _, v := range CurData {
						//changing task type for UI
						v.Task.Type = DefineTaskTypeView(v.Task.Type)

						if _, ok := filterInfo.ByUser[v.Task.CreatedBy]; !ok {
							filterInfo.ByUser[v.Task.CreatedBy] = 0
						}

						isFilteredType := isFiltered(filters.Type, v.Task.Type)
						isFilteredByUser := isFiltered(filters.ByUser, v.Task.CreatedBy) || isFiltered(filters.ByUser, v.Task.ModifiedBy)
						isFilteredLastRunStatus := isFiltered(filters.LastRunStatus, v.LastRunStatus)

						if isFilteredLastRunStatus && isFilteredType {
							filterInfo.ByUser[v.Task.CreatedBy] = filterInfo.ByUser[v.Task.CreatedBy] + 1
							if v.Task.CreatedBy != v.Task.ModifiedBy && v.Task.ModifiedBy != "" {
								filterInfo.ByUser[v.Task.ModifiedBy] = filterInfo.ByUser[v.Task.ModifiedBy] + 1
							}
						}

						if isFilteredType && isFilteredByUser {
							filterInfo.LastRunStatus[v.LastRunStatus] = filterInfo.LastRunStatus[v.LastRunStatus] + 1
						}

						if isFilteredLastRunStatus && isFilteredByUser {
							filterInfo.Type[v.Task.Type] = filterInfo.Type[v.Task.Type] + 1
						}

						if !isFilteringExist {
							resSlice = append(resSlice, v)
							continue
						}

						if !isFilteredTask(&filters, &v) {
							continue
						}

						resSlice = append(resSlice, v)
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&RecentTasksActivity{}))
						resSlice, err = Relay.Filter(string(args.Filter), val, resSlice)
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
							case "LASTRUNTIME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RecentTasksActivityLRTimeASC).Sort(resSlice)
								} else {
									Relay.SortBy(RecentTasksActivityLRTimeDESC).Sort(resSlice)
								}
							case "LASTRUNSTATUS":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RecentTasksActivityLRStatusASC).Sort(resSlice)
								} else {
									Relay.SortBy(RecentTasksActivityLRStatusDESC).Sort(resSlice)
								}
							case "CREATEDBY":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RecentTasksActivityCreatedByASC).Sort(resSlice)
								} else {
									Relay.SortBy(RecentTasksActivityCreatedByDESC).Sort(resSlice)
								}
							case "TYPE":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RecentTasksActivityTypeASC).Sort(resSlice)
								} else {
									Relay.SortBy(RecentTasksActivityTypeDESC).Sort(resSlice)
								}
							case "NAME":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RecentTasksActivityNameASC).Sort(resSlice)
								} else {
									Relay.SortBy(RecentTasksActivityNameDESC).Sort(resSlice)
								}
							case "CREATEDAT":
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RecentTasksActivityCreatedAtASC).Sort(resSlice)
								} else {
									Relay.SortBy(RecentTasksActivityCreatedAtDESC).Sort(resSlice)
								}
							default:
								return nil, errors.New("RecentTasksActivity Sort [" + Column + "] No such column exist!!!")
							}
						}
					}

					staticInfo, err := json.Marshal(filterInfo)
					if err != nil {
						return nil, err
					}
					return Relay.ConnectionFromArray(resSlice, args, string(staticInfo)), nil
				}
				return nil, nil
			},
		},
	},
})

// isFilteredTask - check task to meet the filters conditions
func isFilteredTask(f *TaskFilters, rta *RecentTasksActivity) bool {
	switch {
	case !isFiltered(f.LastRunStatus, rta.LastRunStatus),
		!isFiltered(f.Type, rta.Task.Type),
		!isFiltered(f.ByUser, rta.Task.CreatedBy) && !isFiltered(f.ByUser, rta.Task.ModifiedBy):
		return false
	default:
		return true
	}
}

// isFiltered - if task satisfied the filter
func isFiltered(f []string, param string) bool {
	return len(f) == 0 || strSliceContains(f, param)
}

// RecentTasksActivityCreatedAtASC function sorts data by CreatedAt field in ascending order
// sorting layers: CreatedAt(ASC)  > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityCreatedAtASC > subsortingByNameASC
func RecentTasksActivityCreatedAtASC(p1, p2 interface{}) bool {
	p1CreatedAt := p1.(RecentTasksActivity).Task.CreatedAt.Truncate(time.Minute)
	p2CreatedAt := p2.(RecentTasksActivity).Task.CreatedAt.Truncate(time.Minute)

	if p1CreatedAt == p2CreatedAt {
		return subsortingByNameASC(p1, p2)
	}

	return p1CreatedAt.Before(p2CreatedAt)
}

// RecentTasksActivityCreatedAtDESC function sorts data by CreatedAt field in descending order
// sorting layers: CreatedAt(DESC)  > Name(DESC)
// function chain according to sorting layers: RecentTasksActivityCreatedAtDESC > subsortingByNameASC
func RecentTasksActivityCreatedAtDESC(p1, p2 interface{}) bool {
	p1CreatedAt := p1.(RecentTasksActivity).Task.CreatedAt.Truncate(time.Minute)
	p2CreatedAt := p2.(RecentTasksActivity).Task.CreatedAt.Truncate(time.Minute)

	if p1CreatedAt == p2CreatedAt {
		return subsortingByNameASC(p1, p2)
	}

	return p1CreatedAt.After(p2CreatedAt)
}

// RecentTasksActivityLRTimeASC function sorts data by LastRunTime field in ascending order
// sorting layers: LRTime > LRStatus(ASC) > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityLRTimeASC > LRTimeSubsortingByLRStatusASC > subsortingByNameASC
func RecentTasksActivityLRTimeASC(p1, p2 interface{}) bool {
	p1RTime := p1.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)
	p2RTime := p2.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)

	if p1RTime == p2RTime {
		return LRTimeSubsortingByLRStatusASC(p1, p2)
	}

	return p1RTime.Before(p2RTime)
}

// RecentTasksActivityLRTimeDESC function sorts data by LastRunTime field in descending order
// sorting layers: LRTime > LRStatus(ASC) > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityLRTimeDESC > LRTimeSubsortingByLRStatusASC > subsortingByNameASC
func RecentTasksActivityLRTimeDESC(p1, p2 interface{}) bool {
	p1RTime := p1.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)
	p2RTime := p2.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)

	if p1RTime == p2RTime {
		return LRTimeSubsortingByLRStatusASC(p1, p2)
	}

	return p1RTime.After(p2RTime)
}

// LRTimeSubsortingByLRStatusASC function sorts data by LastRunTime field in ascending order
func LRTimeSubsortingByLRStatusASC(p1, p2 interface{}) bool {
	p1Status := getSortableLastRunStatus(p1.(RecentTasksActivity).LastRunStatus)
	p2Status := getSortableLastRunStatus(p2.(RecentTasksActivity).LastRunStatus)

	if p1Status == p2Status {
		return subsortingByNameASC(p1, p2)
	}

	return p1Status < p2Status
}

func subsortingByNameASC(p1, p2 interface{}) bool {
	p1Name := p1.(RecentTasksActivity).Task.Name
	p2Name := p2.(RecentTasksActivity).Task.Name
	return Relay.StringLessOp(p1Name, p2Name)
}

// RecentTasksActivityLRStatusASC function sorts data by LastRunStatus field in ascending order
// sorting layers: LRStatus > LRTime(DESC)  > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityLRStatusASC > subsortingByLRTimeDESC > subsortingByNameASC
func RecentTasksActivityLRStatusASC(p1, p2 interface{}) bool {
	p1Status := getSortableLastRunStatus(p1.(RecentTasksActivity).LastRunStatus)
	p2Status := getSortableLastRunStatus(p2.(RecentTasksActivity).LastRunStatus)

	if p1Status == p2Status {
		return subsortingByLRTimeDESC(p1, p2)
	}

	return p1Status < p2Status
}

// RecentTasksActivityLRStatusDESC function sorts data by LastRunStatus field in descending order
// sorting layers: LRStatus > LRTime(DESC)  > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityLRStatusDESC > subsortingByLRTimeDESC > subsortingByNameASC
func RecentTasksActivityLRStatusDESC(p1, p2 interface{}) bool {
	p1Status := getSortableLastRunStatus(p1.(RecentTasksActivity).LastRunStatus)
	p2Status := getSortableLastRunStatus(p2.(RecentTasksActivity).LastRunStatus)

	if p1Status == p2Status {
		return subsortingByLRTimeDESC(p1, p2)
	}

	return p1Status > p2Status
}

func subsortingByLRTimeDESC(p1, p2 interface{}) bool {
	p1RTime := p1.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)
	p2RTime := p2.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)

	if p1RTime == p2RTime {
		return subsortingByNameASC(p1, p2)
	}

	return p1RTime.After(p2RTime)
}

// RecentTasksActivityTypeASC function sorts data by Type field in ascending order
// sorting layers: Type > LRStatus(ASC) > LRTime(DESC) > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityTypeASC > RecentTasksActivityLRStatusASC > subsortingByLRTimeDESC > subsortingByNameASC
func RecentTasksActivityTypeASC(p1, p2 interface{}) bool {
	p1Type := getSortableType(p1.(RecentTasksActivity).Task.Type)
	p2Type := getSortableType(p2.(RecentTasksActivity).Task.Type)

	if p1Type == p2Type {
		return RecentTasksActivityLRStatusASC(p1, p2)
	}

	return p1Type < p2Type
}

// RecentTasksActivityTypeDESC function sorts data by Type field in descending order
// sorting layers: Type > LRStatus(ASC) > LRTime(DESC) > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityTypeDESC > RecentTasksActivityLRStatusASC > subsortingByLRTimeDESC > subsortingByNameASC
func RecentTasksActivityTypeDESC(p1, p2 interface{}) bool {
	p1Type := getSortableType(p1.(RecentTasksActivity).Task.Type)
	p2Type := getSortableType(p2.(RecentTasksActivity).Task.Type)

	if p1Type == p2Type {
		return RecentTasksActivityLRStatusASC(p1, p2)
	}

	return p1Type > p2Type
}

// RecentTasksActivityCreatedByASC function sorts data by CreatedBy field in ascending order
// sorting layers: CreateBy > LRStatus(ASC) > LRTime(DESC) > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityCreatedByASC > RecentTasksActivityLRStatusASC > subsortingByLRTimeDESC > subsortingByNameASC
func RecentTasksActivityCreatedByASC(p1, p2 interface{}) bool {
	p1CreatedBy := p1.(RecentTasksActivity).Task.CreatedBy
	p2CreatedBy := p2.(RecentTasksActivity).Task.CreatedBy

	if p1CreatedBy == p2CreatedBy {
		return RecentTasksActivityLRStatusASC(p1, p2)
	}

	return Relay.StringLessOp(p1CreatedBy, p2CreatedBy)
}

// RecentTasksActivityCreatedByDESC function sorts data by CreatedBy field in descending order
// sorting layers: CreateBy > LRStatus(ASC) > LRTime(DESC) > Name(ASC)
// function chain according to sorting layers: RecentTasksActivityCreatedByDESC > RecentTasksActivityLRStatusASC > subsortingByLRTimeDESC > subsortingByNameASC
func RecentTasksActivityCreatedByDESC(p1, p2 interface{}) bool {
	p1CreatedBy := p1.(RecentTasksActivity).Task.CreatedBy
	p2CreatedBy := p2.(RecentTasksActivity).Task.CreatedBy

	if p1CreatedBy == p2CreatedBy {
		return RecentTasksActivityLRStatusASC(p1, p2)
	}

	return !Relay.StringLessOp(p1CreatedBy, p2CreatedBy)
}

// RecentTasksActivityNameASC function sorts data by Name field in ascending order
// sorting layers: Name > LRStatus(ASC) > LRTime(DESC)
// function chain according to sorting layers: RecentTasksActivityNameASC > NameSubsortingByLRStatusASC > LRStatusSubsortingByLRTimeDESC
func RecentTasksActivityNameASC(p1, p2 interface{}) bool {
	p1Name := p1.(RecentTasksActivity).Task.Name
	p2Name := p2.(RecentTasksActivity).Task.Name

	if p1Name == p2Name {
		return NameSubsortingByLRStatusASC(p1, p2)
	}

	return Relay.StringLessOp(p1Name, p2Name)
}

// RecentTasksActivityNameDESC function sorts data by Name field in descending order
// sorting layers: Name > LRStatus(ASC) > LRTime(DESC)
// function chain according to sorting layers: RecentTasksActivityNameDESC > NameSubsortingByLRStatusASC > LRStatusSubsortingByLRTimeDESC
func RecentTasksActivityNameDESC(p1, p2 interface{}) bool {
	p1Name := p1.(RecentTasksActivity).Task.Name
	p2Name := p2.(RecentTasksActivity).Task.Name

	if p1Name == p2Name {
		return NameSubsortingByLRStatusASC(p1, p2)
	}

	return !Relay.StringLessOp(p1Name, p2Name)
}

//NameSubsortingByLRStatusASC function sorts data by LRStatus field in ascending order
func NameSubsortingByLRStatusASC(p1, p2 interface{}) bool {
	p1Status := getSortableLastRunStatus(p1.(RecentTasksActivity).LastRunStatus)
	p2Status := getSortableLastRunStatus(p2.(RecentTasksActivity).LastRunStatus)

	if p1Status == p2Status {
		return LRStatusSubsortingByLRTimeDESC(p1, p2)
	}

	return p1Status < p2Status
}

//LRStatusSubsortingByLRTimeDESC function sorts data by LRTime field in descending order
func LRStatusSubsortingByLRTimeDESC(p1, p2 interface{}) bool {
	p1RTime := p1.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)
	p2RTime := p2.(RecentTasksActivity).TaskInstance.LastRunTime.Truncate(time.Minute)
	return p1RTime.After(p2RTime)
}

func getSortableLastRunStatus(status string) int {
	switch status {
	case "Failed":
		return failedStatus
	case "Some Failures":
		return someFailuresStatus
	case "Running":
		return runningStatus
	case "Success":
		return successStatus
	}
	return undefined
}

func getSortableType(taskType string) int {
	switch taskType {
	case ScriptTaskType:
		return first
	case SequenceTaskType:
		return second
	}
	return undefined
}

//DefineTaskTypeView - fucntion that changes task type specifically for UI (script -> task)
func DefineTaskTypeView(taskType string) (taskTypeView string) {
	if taskType == ScriptTaskType {
		return "task"
	}
	return taskType
}

func strSliceContains(strSrlice []string, candidate string) bool {
	for _, existed := range strSrlice {
		if candidate == existed {
			return true
		}
	}
	return false
}

func interfaceToString(Itr interface{}) (sReturn string) {
	if Itr != nil {
		sReturn = Itr.(string)
	}
	return sReturn
}
