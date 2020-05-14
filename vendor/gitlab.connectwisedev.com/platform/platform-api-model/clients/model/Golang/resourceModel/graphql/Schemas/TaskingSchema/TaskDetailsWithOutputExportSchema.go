package TaskingSchema

import (
	statuses "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

// TaskDetailsWithOutputExportType : TaskDetailsWithOutputExportType GraphQL Schema
var TaskDetailsWithOutputExportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskDetailsWithOutputExportType",
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
					return CurData.Task.Type, nil
				}
				return nil, nil
			},
		},
		"lastRunTime": &graphql.Field{
			Type:        graphql.String,
			Description: "LastRunTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.TaskInstance.StartedAt.Format(time.RFC3339Nano), nil
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
					// if Task is success returns STD_OUT
					// in other cases - STD_ERR
					status, _ := statuses.TaskInstanceStatusFromText(CurData.ExecutionResult.ExecutionStatus)
					if status == statuses.TaskInstanceSuccess {
						return CurData.ExecutionResult.StdOut, nil
					}
					return CurData.ExecutionResult.StdErr, nil
				}
				return nil, nil
			},
		},
		"deviceName": &graphql.Field{
			Type:        graphql.String,
			Description: "Device name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.AssetData.System.SystemName, nil
				}
				return nil, nil
			},
		},
		"deviceFriendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "Device friendly name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.AssetData.FriendlyName, nil
				}
				return nil, nil
			},
		},
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskDetailsWithOutput); ok {
					return CurData.SiteData.SiteName, nil
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

// TaskDetailsWithOutputExportList : TaskDetailsWithOutputExportList connection definition structure
var TaskDetailsWithOutputExportList = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TaskDetailsWithOutputExportList",
	NodeType: TaskDetailsWithOutputExportType,
})

//TaskDetailsWithOutputExportListType : TaskDetailsWithOutputExportListType GraphQL Schema
var TaskDetailsWithOutputExportListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskDetailsWithOutputExportList",
	Fields: graphql.Fields{
		"TaskDetailsWithOutputExportList": &graphql.Field{
			Type:        TaskDetailsWithOutputExportList.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of TaskDetailsWithOutput for export",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.([]TaskDetailsWithOutput); ok {
					var resSlice []interface{}
					for _, val := range CurData {
						resSlice = append(resSlice, val)
					}
					return Relay.ConnectionFromArray(resSlice, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})
