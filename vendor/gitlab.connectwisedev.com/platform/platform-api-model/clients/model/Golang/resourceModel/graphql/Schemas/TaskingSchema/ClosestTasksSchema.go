package TaskingSchema

import "github.com/graphql-go/graphql"

//EndpointsClosestTasks - closest tasks by endpoint data representation
type EndpointsClosestTasks map[string]ClosestTasks

//ClosestTasks - closest tasks container
type ClosestTasks struct {
	Previous *ClosestTask `json:"previous,omitempty"`
	Next     *ClosestTask `json:"next,omitempty"`
}

//ClosestTask - closest task details
type ClosestTask struct {
	Name    string `json:"name,omitempty"`
	RunDate int64  `json:"runDate,omitempty"`
	Status  string `json:"status,omitempty"`
}

//ClosestTasksType : ClosestTasks GraphQL Schema
var ClosestTasksType = graphql.NewObject(graphql.ObjectConfig{
	Name: "closestTasks",
	Fields: graphql.Fields{
		"previous": &graphql.Field{
			Type:        ClosestTaskType,
			Description: "previous",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ClosestTasks); ok && CurData.Previous != nil {
					return CurData.Previous, nil
				}
				return nil, nil
			},
		},
		"next": &graphql.Field{
			Type:        ClosestTaskType,
			Description: "next",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ClosestTasks); ok && CurData.Next != nil {
					return CurData.Next, nil
				}
				return nil, nil
			},
		},
	},
})

//ClosestTaskType : ClosestTask GraphQL Schema
var ClosestTaskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "closestTask",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(*ClosestTask); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"runDate": &graphql.Field{
			Type:        graphql.String,
			Description: "runDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(*ClosestTask); ok {
					return CurData.RunDate, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(*ClosestTask); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
	},
})
