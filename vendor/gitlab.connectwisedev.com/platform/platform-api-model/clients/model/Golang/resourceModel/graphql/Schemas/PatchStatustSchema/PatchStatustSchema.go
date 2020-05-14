package PatchStatustSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//PatchStatusData : PatchStatusData Structure
type PatchStatusData struct {
	RegID    int64  `json:"regId"`
	SiteID   int64  `json:"siteId"`
	Status   string `json:"status"`
	OsStatus string `json:"osStatus"`
	TPStatus string `json:"tpStatus"`
}

//PatchStatusType : PatchStatus GraphQL Schema
var PatchStatusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "patchstatus",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchStatusData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchStatusData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch status of the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchStatusData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"osStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Windows patch status of the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchStatusData); ok {
					return CurData.OsStatus, nil
				}
				return nil, nil
			},
		},

		"tpStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Third party patch status of the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchStatusData); ok {
					return CurData.TPStatus, nil
				}
				return nil, nil
			},
		},
	},
})

//PatchStatusConnectionDefinition : PatchStatusConnectionDefinition structure
var PatchStatusConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PatchStatus",
	NodeType: PatchStatusType,
})

//PatchStatusListData : PatchStatusListData Structure
type PatchStatusListData struct {
	PatchStatust []PatchStatusData `json:"patchStatustList"`
	TotalCount   int64             `json:"totalCount"`
}

//PatchStatusListType : PatchStatusList GraphQL Schema
var PatchStatusListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchStatusList",
	Fields: graphql.Fields{
		"patchstatusList": &graphql.Field{
			Type:        PatchStatusConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "patchStatust list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PatchStatusListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.PatchStatust {
						arraySliceRet = append(arraySliceRet, CurData.PatchStatust[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PatchStatusData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						StatusASC := func(p1, p2 interface{}) bool {
							return p1.(PatchStatusData).Status < p2.(PatchStatusData).Status
						}
						StatusDESC := func(p1, p2 interface{}) bool {
							return p1.(PatchStatusData).Status > p2.(PatchStatusData).Status
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "STATUS" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(StatusASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(StatusDESC).Sort(arraySliceRet)
								}
							} else {
								return nil, errors.New("PatchStatusData Sort [" + Column + "] No such column exist!!!")
							}
						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},

		"totalCount": &graphql.Field{
			Type:        graphql.String,
			Description: "totalCount of patchStatust list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchStatusListData); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})

// TPMissedPatchData : struct of missed third party patch
type TPMissedPatchData struct {
	AppName       string `json:"appName"`
	ExecutionID   string `json:"executionID"`
	Message       string `json:"message"`
	AppVersion    string `json:"appVersion"`
	WhitelistDate string `json:"whitelistDate"`
	FailedDate    string `json:"failedDate"`
}

//TPMissedPatchType : TPMissedPatchType GraphQL Schema
var TPMissedPatchType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TPMissedPatch",
	Fields: graphql.Fields{
		"appName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissedPatchData); ok {
					return CurData.AppName, nil
				}
				return nil, nil
			},
		},

		"executionID": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissedPatchData); ok {
					return CurData.ExecutionID, nil
				}
				return nil, nil
			},
		},

		"message": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch status of the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissedPatchData); ok {
					return CurData.Message, nil
				}
				return nil, nil
			},
		},

		"appVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "Windows patch status of the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissedPatchData); ok {
					return CurData.AppVersion, nil
				}
				return nil, nil
			},
		},

		"whitelistDate": &graphql.Field{
			Type:        graphql.String,
			Description: "Third party patch status of the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissedPatchData); ok {
					return CurData.WhitelistDate, nil
				}
				return nil, nil
			},
		},

		"failedDate": &graphql.Field{
			Type:        graphql.String,
			Description: "Third party patch status of the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissedPatchData); ok {
					return CurData.FailedDate, nil
				}
				return nil, nil
			},
		},
	},
})

type TPMissPatchTimeRangeData struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

//TPMissedPatchType : TPMissedPatchType GraphQL Schema
var TPMissPatchTimeRangeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TPMissPatchTimeRange",
	Fields: graphql.Fields{
		"startTime": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissPatchTimeRangeData); ok {
					return CurData.StartTime, nil
				}
				return nil, nil
			},
		},

		"endTime": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TPMissPatchTimeRangeData); ok {
					return CurData.EndTime, nil
				}
				return nil, nil
			},
		},
	},
})

type PatchTPMissedCollectionData struct {
	Status             string                     `json:"status"`
	State              string                     `json:"state"`
	MissedPatches      []TPMissedPatchData        `json:"missedPatches"`
	ExecutedBy         string                     `json:"executedBy"`
	ExecutedAt         string                     `json:"executedAt"`
	NextDeployment     string                     `json:"nextDeployment"`
	TimeRange          []TPMissPatchTimeRangeData `json:"timeRange"`
	FailedAttemptCount int64                      `json:"failedAttemptCount"`
	PolicyName         string                     `json:"policyName"`
	PolicyID           string                     `json:"policyID"`
}

//TPMissedPatchType : TPMissedPatchType GraphQL Schema
var PatchTPMissedCollectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchTPMissedCollection",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"state": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.State, nil
				}
				return nil, nil
			},
		},
		"missedPatches": &graphql.Field{
			Type:        TPMissedPatchType,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.MissedPatches, nil
				}
				return nil, nil
			},
		},
		"executedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.ExecutedBy, nil
				}
				return nil, nil
			},
		},
		"executedAt": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.ExecutedAt, nil
				}
				return nil, nil
			},
		},
		"nextDeployment": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.NextDeployment, nil
				}
				return nil, nil
			},
		},
		"timeRange": &graphql.Field{
			Type:        TPMissPatchTimeRangeType,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.TimeRange, nil
				}
				return nil, nil
			},
		},
		"failedAttemptCount": &graphql.Field{
			Type:        graphql.Int,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.FailedAttemptCount, nil
				}
				return nil, nil
			},
		},
		"policyName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.PolicyName, nil
				}
				return nil, nil
			},
		},
		"policyID": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPMissedCollectionData); ok {
					return CurData.PolicyID, nil
				}
				return nil, nil
			},
		},
	},
})
