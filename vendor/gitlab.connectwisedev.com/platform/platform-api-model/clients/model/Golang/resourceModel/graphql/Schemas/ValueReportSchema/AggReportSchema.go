package ValueReportSchema

import (
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

// ValueReportRisk response
type ValueReportRisk struct {
	ProfileID string                  `json:"profileID"`
	Months    []AggregateReportDaily  `json:"month"`
	Years     []AggregateReportRecord `json:"year"`
}

// DailyResponse response for daily rows
type DailyResponse struct {
	ProfileID string                 `json:"profileID"`
	Month     []AggregateReportDaily `json:"month"`
}

// MonthlyResponse response for monthly rows
type MonthlyResponse struct {
	ProfileID string                  `json:"profileID"`
	Year      []AggregateReportRecord `json:"year"`
}

// AggregateDataType : AggregateData GraphQL Schema
var AggregateDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "aggregateData",
	Fields: graphql.Fields{
		"profileID": &graphql.Field{
			Type:        graphql.String,
			Description: "profile id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ValueReportRisk); ok {
					return CurData.ProfileID, nil
				}
				return nil, nil
			},
		},

		"month": &graphql.Field{
			Type:        AggregateMonthConnectionDefinition.ConnectionType,
			Description: "Daily aggregated list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(ValueReportRisk); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Months {
						arraySliceRet = append(arraySliceRet, CurData.Months[ind])
					}
					timeSort := func(p1, p2 interface{}) bool {
						return p1.(AggregateReportDaily).TimeStamp < p2.(AggregateReportDaily).TimeStamp
					}
					Relay.SortBy(timeSort).Sort(arraySliceRet)
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},

		"year": &graphql.Field{
			Type:        AggregateYearConnectionDefinition.ConnectionType,
			Description: "Monthly aggregated list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(ValueReportRisk); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Years {
						arraySliceRet = append(arraySliceRet, CurData.Years[ind])
					}
					timeSort := func(p1, p2 interface{}) bool {
						return p1.(AggregateReportRecord).TimeStamp < p2.(AggregateReportRecord).TimeStamp
					}
					Relay.SortBy(timeSort).Sort(arraySliceRet)
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})

// AggregateMonthConnectionDefinition : AggregateMonthConnectionDefinition structure
var AggregateMonthConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "month",
	NodeType: DailyType,
})

// AggregateYearConnectionDefinition : AggregateYearConnectionDefinition structure
var AggregateYearConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "year",
	NodeType: MonthlyType,
})

//ValueReportRiskListData : ValueReportRisk Structure
type ValueReportRiskListData struct {
	PartnerID   string            `json:"partnerID"`
	SiteID      string            `json:"siteID"`
	ClientID    string            `json:"clientID"`
	ValueReport []ValueReportRisk `json:"reportList"`
}

//ValueReportListType : ValueReportListType GraphQL Schema
var ValueReportListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ValueReportList",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type:        graphql.String,
			Description: "partner id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ValueReportRiskListData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"siteID": &graphql.Field{
			Type:        graphql.String,
			Description: "site id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ValueReportRiskListData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"clientID": &graphql.Field{
			Type:        graphql.String,
			Description: "client id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ValueReportRiskListData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"valueReportList": &graphql.Field{
			Type: graphql.NewList(AggregateDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ValueReportRiskListData); ok {
					return CurData.ValueReport, nil
				}
				return nil, nil
			},
		},
	},
})
