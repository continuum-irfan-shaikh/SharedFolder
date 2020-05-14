package RegistrationTokenSchema

import (
	"github.com/graphql-go/graphql"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

type RegistrationToken struct {
	Token  string `json:"token"`
}

//RegistrationTokenType : Registration Token GraphQL Schema
var RegistrationTokenType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "RegistrationTokenData",
	Description: "Custom Corvire-Freshdesk types in resolved tickets",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "Registration Token",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RegistrationToken); ok {
					return CurData.Token, nil
				}
				return nil, nil
			},
		},
	},
})

//RegistrationTokenList : RegistrationTokenList structure
var RegistrationTokenList = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "registrationTokenList",
	NodeType: RegistrationTokenType,
})

//RegistrationTokenListData : RegistrationTokenListData Structure
type RegistrationTokenListData struct {
	RegistrationTokens []RegistrationToken `json:"registrationTokens"`
}

// Example: {{graphQL}}/GraphQL/?{{
//	"query": "query {registrationTokens(clientID:$clientID,siteID:$siteID,endpointID:$endpointID,agentID:$agentID,legacyRegID:$legacyRegID,ipAddress:$ipAddress){registrationTokenList{edges{cursor,node{token}}}}}",
//	"variables": {
//		"clientID": "50109291",
//		"siteID": "50109291",
//		"endpointID": "",
//		"agentID": "",
//		"legacyRegID": "",
//		"ipAddress": "1.1.1.1"
//	}
//}}
var RegistrationTokenListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "registrationTokenList",
	Fields: graphql.Fields{
		"registrationTokenList": &graphql.Field{
			Type:        RegistrationTokenList.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Registration Token List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)

				if CurData, ok := p.Source.([]RegistrationToken); ok {
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