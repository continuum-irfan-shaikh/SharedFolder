package TasksAndSequencesSelectorSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//TasksAndSequencesSelector : TasksAndSequencesSelector structure
type TasksAndSequencesSelector struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Engine      string    `json:"engine"`
	Categories  []string  `json:"categories"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UpdatedBy   string    `json:"updatedBy"`
	TaskType    TaskType  `json:"taskType"`
}

const (
	//ParametrizedTask : ParametrizedTask type of the TasksAndSequencesSelector structure
	ParametrizedTask = "PARAMETRIZED_TASK"
	//NotParametrizedTask : NotParametrizedTask type of the TasksAndSequencesSelector structure
	NotParametrizedTask = "NOT_PARAMETRIZED_TASK"
	//CustomTaskTemplate : CustomTaskTemplate type of the TasksAndSequencesSelector structure
	CustomTaskTemplate = "CUSTOM_TASK_TEMPLATE"
	//CustomSequenceTemplate : CustomSequenceTemplate type of the TasksAndSequencesSelector structure
	CustomSequenceTemplate = "CUSTOM_SEQUENCE_TEMPLATE"
	//CustomPatchingTemplate : CustomPatchingTemplate type of the TasksAndSequencesSelector structure
	CustomPatchingTemplate = "CUSTOM_PATCHING_TEMPLATE"
	//CustomPatchingkPolicyTemplate : CustomPatchingkPolicyTemplate type of TasksAndSequencesSelector structure
	CustomPatchingkPolicyTemplate = "CUSTOM_PATCHING_POLICY_TEMPLATE"
	powerShell                    = "powershell"
	cmd                           = "cmd"
	bash                          = "bash"
	windows                       = "windows"
	linux                         = "linux"
)

type TaskType int

const (
	_ TaskType = iota
	Action
	Script
	ApplicationUpdate
	OSUpdate
	Sequence
	WebrootAction
)

var taskTypes = map[TaskType]string{
	Action:            "Action",
	Script:            "Script",
	ApplicationUpdate: "Application Update",
	OSUpdate:          "OS Update",
	Sequence:          "Sequence",
	WebrootAction:     "Webroot Action",
}

//TasksAndSequencesSelectorType : TasksAndSequencesSelectorType GraphQL Schema
var TasksAndSequencesSelectorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "selector",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
		"engine": &graphql.Field{
			Type:        graphql.String,
			Description: "Engine of the script",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					if CurData.Engine == powerShell || CurData.Engine == cmd {
						return windows, nil
					}
					if CurData.Engine == bash {
						return linux, nil
					}
					return CurData.Engine, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Description",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
		"categories": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "Categories",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},
		"CreatedAt": &graphql.Field{
			Type:        graphql.String,
			Description: "CreatedAt",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.CreatedAt.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"CreatedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "CreatedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.CreatedBy, nil
				}
				return nil, nil
			},
		},
		"UpdatedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "UpdatedBy for custom tasks",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.UpdatedBy, nil
				}
				return nil, nil
			},
		},
		"UpdatedAt": &graphql.Field{
			Type:        graphql.String,
			Description: "UpdatedAt for custom tasks",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return CurData.UpdatedAt.Format(time.RFC3339), nil
				}
				return nil, nil
			},
		},
		"taskType": &graphql.Field{
			Type:        graphql.String,
			Description: "Task Type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TasksAndSequencesSelector); ok {
					return taskTypes[CurData.TaskType], nil
				}
				return nil, nil
			},
		},
	},
})

//TasksAndSequencesSelectorList : TasksAndSequencesSelectorList connection definition structure
var TasksAndSequencesSelectorList = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TasksAndSequencesSelectorList",
	NodeType: TasksAndSequencesSelectorType,
})

//TasksAndSequencesSelectorListData : List of TasksAndSequencesSelector Structures
type TasksAndSequencesSelectorListData struct {
	TasksAndSequencesSelectorList []TasksAndSequencesSelector `json:"TasksAndSequencesSelectorList"`
}

//TasksAndSequencesSelectorListType : TasksAndSequencesSelectorListType GraphQL Schema
var TasksAndSequencesSelectorListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TasksAndSequencesSelectorList",
	Fields: graphql.Fields{
		"TasksAndSequencesSelectorList": &graphql.Field{
			Type:        TasksAndSequencesSelectorList.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of TasksAndSequencesSelectors",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TasksAndSequencesSelectorListData); ok {
					var resSlice []interface{}
					for _, val := range CurData.TasksAndSequencesSelectorList {
						resSlice = append(resSlice, val)
					}
					return Relay.ConnectionFromArray(resSlice, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})
