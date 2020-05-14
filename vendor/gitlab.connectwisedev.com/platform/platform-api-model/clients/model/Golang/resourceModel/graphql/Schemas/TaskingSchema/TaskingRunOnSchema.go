package TaskingSchema

import (
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/graphql-go/graphql"
)

// RunOnType resolve RunOn field
var RunOnType = graphql.NewObject(graphql.ObjectConfig{
	Name: "runOn",
	Fields: graphql.Fields{
		"count": &graphql.Field{
			Type:        graphql.Int,
			Description: "Count of targets",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.RunOnData); ok {
					return CurData.TargetCount, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of target (Managed Endpoint or Dynamic group)",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.RunOnData); ok {
					return CurData.TargetType, nil
				}
				return nil, nil
			},
		},
	},
})
