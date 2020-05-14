package TaskingSchema

import (
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/graphql-go/graphql"
)

type (
	TaskHistory struct {
		ID            string               `json:"id"`
		Name          string               `json:"name"`
		OverallStatus string               `json:"status"`
		LastRunTime   time.Time            `json:"lastRunTime"`
		RunTimeUTC    time.Time            `json:"nextRunTime"`
		Description   string               `json:"description"`
		CreatedBy     string               `json:"createdBy"`
		CreatedAt     time.Time            `json:"createdAt"`
		ModifiedBy    string               `json:"modifiedBy"`
		ModifiedAt    time.Time            `json:"modifiedAt"`
		ExecutionInfo ExecutionInfo        `json:"executionInfo"`
		TaskType      string               `json:"taskType"`
		CanBeCanceled bool                 `json:"canBeCanceled"` // this field tells if the task can be canceled or not
		Regularity    apiModels.Regularity `json:"regularity"`
	}

	ExecutionInfo struct {
		DeviceCount  int `json:"deviceCount"`
		SuccessCount int `json:"successCount"`
		FailedCount  int `json:"failedCount"`
	}

	TaskHistoryListData struct {
		TaskHistoryList []TaskHistory `json:"tasksHistory"`
	}
)

var TasksExecutionHistoryListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "tasksExecutionHistoryList2",
	Fields: graphql.Fields{
		"tasksHistory": &graphql.Field{
			Type: graphql.NewList(TasksExecutionHistoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistoryListData); ok {
					return CurData.TaskHistoryList, nil
				}
				return nil, nil
			},
		},
	},
})

var TasksExecutionHistoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "tasksExecutionHistoryList",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.OverallStatus, nil
				}
				return nil, nil
			},
		},
		"lastRunTime": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.LastRunTime, nil
				}
				return nil, nil
			},
		},
		"nextRunTime": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.RunTimeUTC, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
		"createdBy": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.CreatedBy, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"modifiedBy": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.ModifiedBy, nil
				}
				return nil, nil
			},
		},
		"modifiedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.ModifiedAt, nil
				}
				return nil, nil
			},
		},
		"executionInfo": &graphql.Field{
			Type: ExecutionInfoType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.ExecutionInfo, nil
				}
				return nil, nil
			},
		},
		"taskType": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.TaskType, nil
				}
				return nil, nil
			},
		},
		"canBeCanceled": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.CanBeCanceled, nil
				}
				return nil, nil
			},
		},
		"regularity": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskHistory); ok {
					return CurData.Regularity.String(), nil
				}
				return nil, nil
			},
		},
	},
})

var ExecutionInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExecutionInfo",
	Fields: graphql.Fields{
		"deviceCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExecutionInfo); ok {
					return CurData.DeviceCount, nil
				}
				return nil, nil
			},
		},
		"successCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExecutionInfo); ok {
					return CurData.SuccessCount, nil
				}
				return nil, nil
			},
		},
		"failedCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExecutionInfo); ok {
					return CurData.FailedCount, nil
				}
				return nil, nil
			},
		},
	},
})
