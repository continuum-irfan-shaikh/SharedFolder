package AgentSchema

import (
	"github.com/graphql-go/graphql"
)

//PartnerTokenMapData : PartnerTokenMapData Structure
type PartnerTokenMapData struct {
	PartnerID	string
	ClientID	string
	SiteID		string
	EndpointID	string
	AgentID		string
	LegacyRegID	string
	Ipaddress	string
}

//GenerateTokenData : GenerateTokenData Structure
type GenerateTokenData struct {
	Token	string `json:"token"`
}

//GenerateTokenType : GenerateToken GraphQL Schema
var GenerateTokenType = graphql.NewObject(graphql.ObjectConfig{
	Name: "generateToken",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "token",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GenerateTokenData); ok {
					return CurData.Token, nil
				}
				return nil, nil
			},
		},
	},
})

//EndpointMappingsData : EndpointMappingsData Structure
type EndpointMappingsData struct {
	PartnerID	string	`json:"partner_id"`
	ClientID	string	`json:"client_id"`
	SiteID		string	`json:"site_id"`
	EndpointID	string	`json:"endpoint_id"`
	AgentID		string	`json:"agent_id"`
	LegacyRegID	string	`json:"legacy_reg_id"`
}

//UpdateMappingData : UpdateMappingData Structure
type UpdateMappingData struct {
	EndpointID	string `json:"EndpointID"`
	ErrorMessage	string `json:"ErrorMessage"`
}

//UpdateMappingType : UpdateMapping GraphQL Schema
var UpdateMappingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "updateMapping",
	Fields: graphql.Fields{
		"EndpointID": &graphql.Field{
			Type:        graphql.String,
			Description: "EndpointID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UpdateMappingData); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},

		"ErrorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "ErrorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UpdateMappingData); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},
	},
})

//UpdateMappingListData : UpdateMappingList Structure
type UpdateMappingListData struct {
	UpdateMappingList	[]UpdateMappingData	`json:"updateMappingList"`
}

//UpdateMappingListType : UpdateMappingList GraphQL Schema
var UpdateMappingListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "updateMappingList",
	Fields: graphql.Fields{
		"updateMappingList": &graphql.Field{
			Type:        graphql.NewList(UpdateMappingType),
			Description: "UpdateMappingList",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UpdateMappingListData); ok {
					return CurData.UpdateMappingList, nil
				}
				return nil, nil
			},
		},
	},
})
