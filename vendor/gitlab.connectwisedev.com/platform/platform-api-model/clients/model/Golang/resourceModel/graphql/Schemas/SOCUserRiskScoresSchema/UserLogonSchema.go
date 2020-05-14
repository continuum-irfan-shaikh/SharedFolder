package SOCUserRiskScoresSchema

import (
	"github.com/graphql-go/graphql"
)

//Events : Event struct
type Events struct {
	LogonType string `json:"logonType"`
	Success   int32  `json:"success"`
	Failure   int32  `json:"failure"`
}

// EventsType : EventsType GraphQL Schema
var EventsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Events",
	Fields: graphql.Fields{
		"logonType": &graphql.Field{
			Type:        graphql.String,
			Description: "logon Type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Events); ok {
					return CurData.LogonType, nil
				}
				return nil, nil
			},
		},

		"success": &graphql.Field{
			Type:        graphql.Int,
			Description: "success",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Events); ok {
					return CurData.Success, nil
				}
				return nil, nil
			},
		},
		"failure": &graphql.Field{
			Type:        graphql.Int,
			Description: "failure",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Events); ok {
					return CurData.Failure, nil
				}
				return nil, nil
			},
		},
	},
})

//LogonHistory : Logon History
type LogonHistory struct {
	Location string   `json:"location"`
	Events   []Events `json:"events"`
}

// LogonHistoryType : LogonHistoryType GraphQL Schema
var LogonHistoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "userLogon",
	Fields: graphql.Fields{
		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LogonHistory); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},

		"events": &graphql.Field{
			Type:        graphql.NewList(EventsType),
			Description: "Events",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LogonHistory); ok {
					return CurData.Events, nil
				}
				return nil, nil
			},
		},
	},
})

// UserLogonData : User logon details
type UserLogonData struct {
	PartnerID   string         `json:"partnerID"`
	ClientID    string         `json:"clientID"`
	SiteID      string         `json:"siteID"`
	UserID      string         `json:"userID"`
	Domain      string         `json:"domain"`
	UserName    string         `json:"userName"`
	LogonHistry []LogonHistory `json:"logonHistory"`
}

// UserLogonType : UserLogonType GraphQL Schema
var UserLogonType = graphql.NewObject(graphql.ObjectConfig{
	Name: "userLogon",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type:        graphql.String,
			Description: "Partner ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserLogonData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"clientID": &graphql.Field{
			Type:        graphql.String,
			Description: "Client ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserLogonData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type:        graphql.String,
			Description: "site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserLogonData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"userID": &graphql.Field{
			Type:        graphql.String,
			Description: "Total Risk score",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserLogonData); ok {
					return CurData.UserID, nil
				}
				return nil, nil
			},
		},

		"domain": &graphql.Field{
			Type:        graphql.String,
			Description: "Average Risk score",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserLogonData); ok {
					return CurData.Domain, nil
				}
				return nil, nil
			},
		},

		"userName": &graphql.Field{
			Type:        graphql.String,
			Description: "User Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserLogonData); ok {
					return CurData.UserName, nil
				}
				return nil, nil
			},
		},

		"logonHistory": &graphql.Field{
			Type:        graphql.NewList(LogonHistoryType),
			Description: "Date in readable format",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserLogonData); ok {
					return CurData.LogonHistry, nil
				}
				return nil, nil
			},
		},
	},
})

//UserLogonList : UserLogon List struct
type UserLogonList struct {
	UserLogonData []UserLogonData `json:"userLogon"`
}
