package ProfileProtectSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/gocql/gocql"
	"github.com/graphql-go/graphql"
)

//ExecutionData : Execution Data struct
type ExecutionData struct {
	Reason        string   `json:"reason"`
	ReasonDetails []string `json:"reasonDetails"`
}

//ExecutionDataType : ExecutionData Type Schema
var ExecutionDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExecutionData",
	Fields: graphql.Fields{
		"reason": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExecutionData); ok {
					return CurData.Reason, nil
				}
				return nil, nil
			},
		},
		"reasonDetails": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExecutionData); ok {
					return CurData.ReasonDetails, nil
				}
				return nil, nil
			},
		},
	},
})

//ItemData : Item Data struct
type ItemData struct {
	Label            string          `json:"label"`
	Name             string          `json:"name"`
	ExecutionResult  string          `json:"executionResult"`
	ExecutionDetails []ExecutionData `json:"executionDetails"`
	ExecutionStatus  string          `json:"executionStatus"`
	ExecutionError   string          `json:"executionError"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

//ItemType : Item Data GraphQL Schema
var ItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Items",
	Fields: graphql.Fields{
		"label": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ItemData); ok {
					return CurData.Label, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ItemData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"executionResult": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ItemData); ok {
					return CurData.ExecutionResult, nil
				}
				return nil, nil
			},
		},

		"executionDetails": &graphql.Field{
			Type: graphql.NewList(ExecutionDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ItemData); ok {
					return CurData.ExecutionDetails, nil
				}
				return nil, nil
			},
		},

		"executionStatus": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ItemData); ok {
					return CurData.ExecutionStatus, nil
				}
				return nil, nil
			},
		},

		"executionError": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ItemData); ok {
					return CurData.ExecutionError, nil
				}
				return nil, nil
			},
		},

		"updatedAt": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ItemData); ok {
					return CurData.UpdatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

//CategoryProfileScoreData : CategoryProfileScore Data Struct
type CategoryProfileScoreData struct {
	Label             string          `json:"label"`
	Name              string          `json:"name"`
	ExecutionResult   string          `json:"executionResult"`
	ExecutionDetails  []ExecutionData `json:"executionDetails"`
	Items             []ItemData      `json:"items"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	LastExecutionTime time.Time       `json:"lastExecutionTime"`
	ExecutionStatus   string          `json:"executionStatus"`
}

//CategoryProfileScoreType : CategoryProfileScore Data GraphQL Schema
var CategoryProfileScoreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CategoryProfileScore",
	Fields: graphql.Fields{
		"label": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.Label, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"executionResult": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.ExecutionResult, nil
				}
				return nil, nil
			},
		},

		"executionDetails": &graphql.Field{
			Type: graphql.NewList(ExecutionDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.ExecutionDetails, nil
				}
				return nil, nil
			},
		},

		"items": &graphql.Field{
			Type: graphql.NewList(ItemType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.Items, nil
				}
				return nil, nil
			},
		},

		"updatedAt": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.UpdatedAt, nil
				}
				return nil, nil
			},
		},

		"lastExecutionTime": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.LastExecutionTime, nil
				}
				return nil, nil
			},
		},

		"executionStatus": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CategoryProfileScoreData); ok {
					return CurData.ExecutionStatus, nil
				}
				return nil, nil
			},
		},
	},
})

//Script for getting script result
type Script struct {
	ExecutionID            string    `json:"executionid"`
	Name                   string    `json:"name"`
	Description            string    `json:"desc"`
	Content                string    `json:"content"`
	Engine                 string    `json:"executor"`
	Categories             string    `json:"categories"`
	CategoryLabel          string    `json:"categorylabel"`
	SubCategoryLabel       string    `json:"subcategorylabel"`
	EngineMaxVersion       string    `json:"version"`
	ExepectedExecutionTime string    `json:"exectimeinterval"`
	ExecuteNow             bool      `json:"executenow"`
	IsActive               bool      `json:"isactive"`
	IsScriptChanged        bool      `json:"isscriptchanged"`
	CreatedAt              time.Time `json:"createdt"`
	UpdatedAt              time.Time `json:"updatedat"`
	IsMergedScript         bool      `json:"ismergedscript"`
}

//EndpointCatResult for getting data from endpoint_category_result_table
type EndpointCatResult struct {
	PartnerID            string            `json:"partnerID"`
	ClientID             string            `json:"clientID"`
	SiteID               string            `json:"siteID"`
	RegID                string            `json:"regID"`
	EndpointID           string            `json:"endpointID"`
	CategoryID           gocql.UUID        `json:"categoryID"`
	ActualCategoryResult string            `json:"actualCategoryResult"`
	CategoryResult       string            `json:"categoryResult"`
	ExecutionStatus      string            `json:"executionStatus"`
	ErrorDescription     string            `json:"errorDescription"`
	UpdatedAt            int64             `json:"updatedAt"`
	RiskStartTime        int64             `json:"riskStartTime"`
	ExecutionDetails     map[string]string `json:"executionDetails"`
	CategoryLabel        string            `json:"categoryLabel"`
}
