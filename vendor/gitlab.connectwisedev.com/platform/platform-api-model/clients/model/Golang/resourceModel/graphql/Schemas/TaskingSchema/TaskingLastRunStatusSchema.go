package TaskingSchema

import (
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/graphql-go/graphql"
)

const (
	failedStatus       = 1
	someFailuresStatus = 2
	runningStatus      = 3
	successStatus      = 4
	unknownStatus      = -1
)

// LastRunStatusType resolve LastRunStatus field
var LastRunStatusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "lastRunStatus",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Task running summary status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.LastRunStatusData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"deviceCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "Count of target devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.LastRunStatusData); ok {
					return CurData.DeviceCount, nil
				}
				return nil, nil
			},
		},

		"successCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "Count of success runnings",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.LastRunStatusData); ok {
					return CurData.SuccessCount, nil
				}
				return nil, nil
			},
		},

		"failureCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "Count of failed runnings",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(apiModels.LastRunStatusData); ok {
					return CurData.FailureCount, nil
				}
				return nil, nil
			},
		},
	},
})

func getSortableStatus(status apiModels.LastRunStatusData) int {
	if status.DeviceCount <= 0 {
		return unknownStatus
	}

	if status.FailureCount == status.DeviceCount {
		return failedStatus
	}

	if status.SuccessCount == status.DeviceCount {
		return successStatus
	}

	if status.FailureCount+status.SuccessCount == status.DeviceCount {
		return someFailuresStatus
	}
	return runningStatus
}
