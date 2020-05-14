package NOCAccessLevelSchema

import "github.com/graphql-go/graphql"

// NALPlanData : NALPlanData Structure
type NALPlanData struct {
	AccessLevel     string `json:"AccessLevel"`
	AccessName      string `json:"AccessName"`
	AccessLevelDesc string `json:"AccessDesc"`
}

// NALPlanType : NOCAccessLevelPlan GraphQL Schema
var NALPlanType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NALPlanType",
	Fields: graphql.Fields{
		"AccessLevel": &graphql.Field{
			Type:        graphql.String,
			Description: "NOC accesslevel plan",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NALPlanData); ok {
					return CurData.AccessLevel, nil
				}
				return nil, nil
			},
		},
		"AccessName": &graphql.Field{
			Type:        graphql.String,
			Description: "NOC accesslevel plan",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NALPlanData); ok {
					return CurData.AccessName, nil
				}
				return nil, nil
			},
		},
		"AccessDesc": &graphql.Field{
			Type:        graphql.String,
			Description: "NOC accesslevel plan",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NALPlanData); ok {
					return CurData.AccessLevelDesc, nil
				}
				return nil, nil
			},
		},
	},
})

type NALPlanDataListData struct {
	NALPlanDataList []NALPlanData `json:"nocaccesslevel"`
}

// NALPlanDataListType : NOCAccessLevelPlan GraphQL Schema
var NALPlanDataListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NALPlanDataListType",
	Fields: graphql.Fields{
		"nocaccesslevel": &graphql.Field{
			Type: graphql.NewList(NALPlanType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NALPlanDataListData); ok {
					return CurData.NALPlanDataList, nil
				}
				return nil, nil
			},
		},
	},
})

// NOCAccessLevelPlanListData : NOCAccessLevelPlanListData Structure
type NOCAccessLevelPlanListData struct {
	Status  int64                 `json:"status"`
	Outdata []NALPlanDataListData `json:"outdata"`
}

// NOCAccessLevelPlanListType : NOCAccessLevelPlan GraphQL Schema
var NOCAccessLevelPlanListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NOCAccessLevelplanListData",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelPlanListData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"outdata": &graphql.Field{
			Type:        graphql.NewList(NALPlanDataListType),
			Description: "NOCAccessLevel list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAccessLevelPlanListData); ok {
					return CurData.Outdata, nil
				}
				return nil, nil
			},
		},
	},
})
