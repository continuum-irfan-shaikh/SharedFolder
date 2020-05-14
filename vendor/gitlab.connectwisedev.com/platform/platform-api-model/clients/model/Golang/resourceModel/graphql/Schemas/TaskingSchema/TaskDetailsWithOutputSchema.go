package TaskingSchema

import (
	"encoding/json"
	statuses "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/SitesSchema"
	"time"

	"github.com/gocql/gocql"
	"github.com/graphql-go/graphql"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/AssetSchema"
)

// TaskDetailsWithOutput : TaskDetailsWithOutput structure
type TaskDetailsWithOutput struct {
	Task            Task
	TaskInstance    TaskInstance
	ExecutionResult ExecutionResult
	// AssetData and SiteData are for exporting data
	AssetData AssetSchema.AssetCollectionData
	SiteData  SitesSchema.SitesData
}

// ExecutionResult : ExecutionResult structure
type ExecutionResult struct {
	ManagedEndpointID gocql.UUID `json:"managedEndpointId"`
	TaskInstanceID    gocql.UUID `json:"taskInstanceId"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	ExecutionStatus   string     `json:"executionStatus"`
	StdErr            string     `json:"stdErr"`
	StdOut            string     `json:"stdOut"`
}

// TaskDetailsWithOutputType : TaskDetailsWithOutputType GraphQL Schema
// Example: {{graphQL}}/GraphQL/?{TaskDetailsWithOutput(endpointId:"311b9af1-2c53-11e8-b3d1-02d66c22b7a6",limit:"10"){TaskDetailsWithOutputList{edges{cursor,node{name,type,lastRunTime,createdBy,status,output}}pageInfo{startCursor,endCursor,hasNextPage,hasPreviousPage,totalCount}}}}
var TaskDetailsWithOutputType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskDetailsWithOutputType",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Task name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.Task.Name, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Task type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return DefineTaskTypeView(CurData.Task.Type), nil
				}
				return nil, nil
			},
		},
		"lastRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "LastRunTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.TaskInstance.LastRunTime.Format(time.RFC3339Nano), nil
				}
				return nil, nil
			},
		},
		"createdBy": &graphql.Field{
			Type:        graphql.String,
			Description: "Task is CreatedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.Task.CreatedBy, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Task status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.ExecutionResult.ExecutionStatus, nil
				}
				return nil, nil
			},
		},
		"output": &graphql.Field{
			Type:        graphql.String,
			Description: "LastRunTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					status, err := statuses.TaskInstanceStatusFromText(CurData.ExecutionResult.ExecutionStatus)
					if err != nil {
						return nil, err
					}
					// if Task is success returns STD_OUT
					// in other cases - STD_ERR
					if status == statuses.TaskInstanceSuccess {
						return CurData.ExecutionResult.StdOut, nil
					}
					return CurData.ExecutionResult.StdErr, nil
				}
				return nil, nil
			},
		},
		"taskInstanceId": &graphql.Field{
			Type:        graphql.String,
			Description: "TaskInstance ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.TaskInstance.ID, nil
				}
				return nil, nil
			},
		},
	},
})

// TaskDetailsWithOutputList : TaskDetailsWithOutputList connection definition structure
var TaskDetailsWithOutputList = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TaskDetailsWithOutputList",
	NodeType: TaskDetailsWithOutputType,
})

//TaskDetailsWithOutputListType : TaskDetailsWithOutputListType GraphQL Schema
var TaskDetailsWithOutputListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskDetailsWithOutputList",
	Fields: graphql.Fields{
		"TaskDetailsWithOutputList": &graphql.Field{
			Type:        TaskDetailsWithOutputList.ConnectionType,
			Args:        appendArg(Relay.ConnectionArgs, "failedOnly", graphql.ArgumentConfig{Type: graphql.Boolean}),
			Description: "List of TaskDetailsWithOutput",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				isFilteringExist := InterfaceToBool(p.Args["failedOnly"])

				if CurData, ok := p.Source.([]TaskDetailsWithOutput); ok {
					type StaticInfo struct {
						FailureCount int `json:"failureCount"`
					}

					var statInfo = StaticInfo{FailureCount: 0}
					var resSlice []interface{}

					for _, v := range CurData {
						status, err := statuses.TaskInstanceStatusFromText(v.ExecutionResult.ExecutionStatus)
						if err != nil {
							return nil, err
						}

						if status == statuses.TaskInstanceFailed ||
							status == statuses.TaskInstanceSomeFailures {
							statInfo.FailureCount++
						}

						if !isFilteringExist {
							resSlice = append(resSlice, v)
							continue
						}

						if status == statuses.TaskInstanceFailed ||
							status == statuses.TaskInstanceSomeFailures {
							resSlice = append(resSlice, v)
						}
					}
					staticInfo, err := json.Marshal(statInfo)
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

//InterfaceToBool : Function to Convert Interface To Boolean
func InterfaceToBool(Itr interface{}) bool {
	if Itr != nil {
		if v, ok := Itr.(bool); ok {
			return v
		}
	}
	return false
}
