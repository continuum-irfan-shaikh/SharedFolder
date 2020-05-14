package ValueReportSchema

import "github.com/graphql-go/graphql"

// AggregateReportDaily holds single record for a day
type AggregateReportDaily struct {
	ClientID          string  `json:"clientID"`
	ProfileID         string  `json:"profileID"`
	TotalNoOfEntities int     `json:"totalNoOfEntities"`
	TotalRiskScore    int     `json:"totalRiskScore"`
	AverageRiskScore  float64 `json:"averageRiskScore"`
	TimeStamp         int64   `json:"timeStamp"`
	RecordDate        string  `json:"recordDate"`
	RecordYearDate    string  `json:"recordYearDate"`
	MinimumScore      int     `json:"minimumScore"`
	MaximumScore      int     `json:"maximumScore"`
}

// DailyType : DailyType GraphQL Schema
var DailyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "month",
	Fields: graphql.Fields{
		"clientID": &graphql.Field{
			Type:        graphql.String,
			Description: "Client ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},

		"profileID": &graphql.Field{
			Type:        graphql.String,
			Description: "Profile ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.ProfileID, nil
				}
				return nil, nil
			},
		},

		"totalNoOfEntities": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total no of entities",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.TotalNoOfEntities, nil
				}
				return nil, nil
			},
		},

		"totalRiskScore": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total Risk score",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.TotalRiskScore, nil
				}
				return nil, nil
			},
		},

		"averageRiskScore": &graphql.Field{
			Type:        graphql.Float,
			Description: "Average Risk score",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.AverageRiskScore, nil
				}
				return nil, nil
			},
		},

		"timeStamp": &graphql.Field{
			Type:        graphql.String,
			Description: "Timestamp of record",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.TimeStamp, nil
				}
				return nil, nil
			},
		},

		"recordDate": &graphql.Field{
			Type:        graphql.String,
			Description: "Date in readable format",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.RecordDate, nil
				}
				return nil, nil
			},
		},

		"minimumScore": &graphql.Field{
			Type:        graphql.Int,
			Description: "Minimum risk score for a day",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.MinimumScore, nil
				}
				return nil, nil
			},
		},

		"maximumScore": &graphql.Field{
			Type:        graphql.Int,
			Description: "Maximum risk score for a day",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AggregateReportDaily); ok {
					return CurData.MaximumScore, nil
				}
				return nil, nil
			},
		},
	},
})
