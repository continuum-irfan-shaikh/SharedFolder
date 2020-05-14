package ValueReportSchema

import "github.com/graphql-go/graphql"

// AggregateReportRecord holds single record for a month or week
type AggregateReportRecord struct {
	RecordDate        string  `json:"recordDate"`
	TimeStamp         int64   `json:"timeStamp"`
	AverageRiskScore  float64 `json:"averageRiskScore"`
	MinimumScore      int     `json:"minimumScore"`
	MaximumScore      int     `json:"maximumScore"`
	TotalNoOfEntities int     `json:"totalNoOfEntities"`
	TotalRiskScore    int     `json:"totalRiskScore"`
}

// MonthlyType : MonthlyType GraphQL Schema
var MonthlyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "year",
	Fields: graphql.Fields{
		"recordDate": &graphql.Field{
			Type:        graphql.String,
			Description: "Record Date contain month and year",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportRecord); ok {
					return CurData.RecordDate, nil
				}
				return nil, nil
			},
		},

		"timeStamp": &graphql.Field{
			Type:        graphql.String,
			Description: "Time stamp of first record in month",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportRecord); ok {
					return CurData.TimeStamp, nil
				}
				return nil, nil
			},
		},

		"averageRiskScore": &graphql.Field{
			Type:        graphql.Float,
			Description: "average risk score month basis",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportRecord); ok {
					return CurData.AverageRiskScore, nil
				}
				return nil, nil
			},
		},
		"totalNoOfEntities": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total no of entities",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportRecord); ok {
					return CurData.TotalNoOfEntities, nil
				}
				return nil, nil
			},
		},

		"totalRiskScore": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total Risk score",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportRecord); ok {
					return CurData.TotalRiskScore, nil
				}
				return nil, nil
			},
		},

		"minimumScore": &graphql.Field{
			Type:        graphql.Int,
			Description: "Minimum risk score for a month",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportRecord); ok {
					return CurData.MinimumScore, nil
				}
				return nil, nil
			},
		},

		"maximumScore": &graphql.Field{
			Type:        graphql.Int,
			Description: "Maximum risk score for a month",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportRecord); ok {
					return CurData.MaximumScore, nil
				}
				return nil, nil
			},
		},
	},
})
