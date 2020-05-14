package TaskingSchema

import (
	"github.com/graphql-go/graphql"
)

//Schedule is structure of Schedule information
type Schedule struct {
	Regularity string  `json:"regularity"`
	StartDate  string  `json:"startDate,omitempty"`
	EndDate    string  `json:"endDate,omitempty"`
	TimeZone   string  `json:"timeZone,omitempty"`
	Repeat     *Repeat `json:"repeat,omitempty"`
}

//SchedulerType : Schedule Data GraphQL Schema
var SchedulerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Schedule",
	Fields: graphql.Fields{
		"regularity": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.Regularity, nil
				}
				return nil, nil
			},
		},

		"startDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.StartDate, nil
				}
				return nil, nil
			},
		},

		"timeZone": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},

		"repeat": &graphql.Field{
			Type: RepeatType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},
	},
})

//Repeat is structure of Schedule information
type Repeat struct {
    Frequency string `json:"frequency,omitempty"`
    Every     int    `json:"every,omitempty"`
    RunTime   string `json:"runTime,omitempty"`
}

//RepeatType : Repeat Data GraphQL Schema
var RepeatType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Repeat",
	Fields: graphql.Fields{
		"frequency": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.Regularity, nil
				}
				return nil, nil
			},
		},

		"every": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.StartDate, nil
				}
				return nil, nil
			},
		},

		"runTime": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},
	},
})

//TaskingData : Tasking Response Structure
type TaskingData struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	CreatedAt          string   `json:"createdAt"`
	CreatedBy          string   `json:"createdBy"`
	PartnerID          string   `json:"partnerId"`
	OriginID           string   `json:"originId"`
	State              string   `json:"state"`
	Trigger            string   `json:"trigger"`
	Type               string   `json:"type"`
	Parameters         string   `json:"parameters"`
	NextRunTime        string   `json:"nextRunTime"`
	ExternalTask       bool     `json:"externalTask"`
	ResultWebhook      string   `json:"resultWebhook"`
	Schedule           Schedule `json:"schedule"`
	IsRequireNOCAccess bool     `json:"isRequireNOCAccess"`
	ModifiedBy         string   `json:"modifiedBy"`
	ModifiedAt         string   `json:"modifiedAt"`
	ScheduleType       string   `json:"scheduleType"`
}

//TaskingDataType : Tasking Data GraphQL Schema
var TaskingDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Tasking",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},

		"createdAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.CreatedAt, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"partnerId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"originId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.OriginID, nil
				}
				return nil, nil
			},
		},
		"state": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.State, nil
				}
				return nil, nil
			},
		},
		"trigger": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.Trigger, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
		"parameters": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.Parameters, nil
				}
				return nil, nil
			},
		},
		"nextRunTime": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.NextRunTime, nil
				}
				return nil, nil
			},
		},
		"externalTask": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.ExternalTask, nil
				}
				return nil, nil
			},
		},
		"resultWebhook": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.ResultWebhook, nil
				}
				return nil, nil
			},
		},
		"schedule": &graphql.Field{
			Type: SchedulerType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.Schedule, nil
				}
				return nil, nil
			},
		},
		"isRequireNOCAccess": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.IsRequireNOCAccess, nil
				}
				return nil, nil
			},
		},
		"modifiedBy": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.ModifiedBy, nil
				}
				return nil, nil
			},
		},
		"modifiedAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingData); ok {
					return CurData.ModifiedAt, nil
				}
				return nil, nil
			},
		},
	},
})

//TaskingDeleteResponse : Tasking Response for delete Structure
type TaskingDeleteResponse struct {
	Response string `json:"response"`
}

//TaskingDeleteResponseType : Delete response GraphQL Schema
var TaskingDeleteResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TaskingDeleteResponse",
	Fields: graphql.Fields{
		"response": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TaskingDeleteResponse); ok {
					return CurData.Response, nil
				}
				return nil, nil
			},
		},

		"every": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.StartDate, nil
				}
				return nil, nil
			},
		},

		"runTime": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Schedule); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},
	},
})
