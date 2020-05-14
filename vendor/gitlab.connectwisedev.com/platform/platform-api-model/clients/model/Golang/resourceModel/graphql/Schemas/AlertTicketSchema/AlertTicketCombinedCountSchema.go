package AlertTicketSchema

import (
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//AlertTicketCombinedCount : AlertTicketCombinedCount Structure
type AlertTicketCombinedCount struct {
	NOCAlertCount           int64 `json:"nocAlertCount"`
	PartnerAlertCount       int64 `json:"partnerAlertCount"`
	PartnerDefineAlertCount int64 `json:"partnerDefineAlertCount"`
	NOCTicketCount          int64 `json:"nocTicketCount"`
	PartnerTicketCount      int64 `json:"partnerTicketCount"`
	RegID                   int64 `json:"regid"`
}

//AlertTicketCombinedCountType : AlertTicketCombinedCount GraphQL Schema
var AlertTicketCombinedCountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AlertTicketCombinedCount",
	Fields: graphql.Fields{
		"nocAlertCount": &graphql.Field{
			Type:        graphql.String,
			Description: "nocAlertCount",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCombinedCount); ok {
					return CurData.NOCAlertCount, nil
				}
				return nil, nil
			},
		},

		"partnerAlertCount": &graphql.Field{
			Type:        graphql.String,
			Description: "partnerAlertCount",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCombinedCount); ok {
					return CurData.PartnerAlertCount, nil
				}
				return nil, nil
			},
		},

		"partnerDefineAlertCount": &graphql.Field{
			Type:        graphql.String,
			Description: "partnerDefinedAlertCount",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCombinedCount); ok {
					return CurData.PartnerDefineAlertCount, nil
				}
				return nil, nil
			},
		},

		"nocTicketCount": &graphql.Field{
			Type:        graphql.String,
			Description: "nocTicketCount",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCombinedCount); ok {
					return CurData.NOCTicketCount, nil
				}
				return nil, nil
			},
		},

		"partnerTicketCount": &graphql.Field{
			Type:        graphql.String,
			Description: "partnerTicketCount",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCombinedCount); ok {
					return CurData.PartnerTicketCount, nil
				}
				return nil, nil
			},
		},

		"regid": &graphql.Field{
			Type:        graphql.String,
			Description: "regid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCombinedCount); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},
	},
})

//AlertTicketCombinedCountConnectionDefinition : AlertTicketCombinedCountConnectionDefinition structure
var AlertTicketCombinedCountConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "AlertTicketCombinedCount",
	NodeType: AlertTicketCombinedCountType,
})

//AlertTicketCombinedCountList : AlertTicketCombinedCountList Structure
type AlertTicketCombinedCountList struct {
	ErrorMessage []string                   `json:"errorMessage"`
	Data         []AlertTicketCombinedCount `json:"outdata"`
}

//AlertTicketCombinedCountListType : AlertTicketCombinedCountList GraphQL Schema
var AlertTicketCombinedCountListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AlertTicketCombinedCountList",
	Fields: graphql.Fields{

		"errorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "errorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertTicketCombinedCountList); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        AlertTicketCombinedCountConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "alert ticket combined count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(AlertTicketCombinedCountList); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Data {
						arraySliceRet = append(arraySliceRet, CurData.Data[ind])
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})
