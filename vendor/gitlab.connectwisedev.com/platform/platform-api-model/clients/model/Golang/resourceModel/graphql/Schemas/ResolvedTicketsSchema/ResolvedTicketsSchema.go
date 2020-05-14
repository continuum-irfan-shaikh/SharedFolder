package ResolvedTicketsSchema

import (
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

type ResolvedTicket struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

//ResolvedTicketsType : Resolved Tickets GraphQL Schema
var ResolvedTicketsType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ResolutionsData",
	Description: "Custom Corvire-Freshdesk types in resolved tickets",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Ticket resolution type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ResolvedTicket); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
		"count": &graphql.Field{
			Type:        graphql.Int,
			Description: "Ticket resolution count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ResolvedTicket); ok {
					return CurData.Count, nil
				}
				return nil, nil
			},
		},
	},
})

//ResolvedTicketsList : ResolvedTicketsList structure
var ResolvedTicketsList = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "resolvedTicketsList",
	NodeType: ResolvedTicketsType,
})

//TaskingCountListData : Resolved Tickets Structure
type ResolutionsListData struct {
	Resolutions []ResolvedTicket `json:"resolutions"`
}

// Example: {{graphQL}}/GraphQL/?{{ "query": "query {ResolvedTickets(from:$from,to:$to){ResolvedTicketsList{edges{cursor,node{type,count}}}}}", "variables": { "from": "1528464646830", "to": "1536240646830"}}
var ResolvedTicketsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "resolvedTicketsList",
	Fields: graphql.Fields{
		"resolvedTicketsList": &graphql.Field{
			Type:        ResolvedTicketsList.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Resolved Tickets list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)

				if CurData, ok := p.Source.([]ResolvedTicket); ok {
					var arraySliceRet []interface{}
					for _, val := range CurData {
						arraySliceRet = append(arraySliceRet, val)
					}

					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})
