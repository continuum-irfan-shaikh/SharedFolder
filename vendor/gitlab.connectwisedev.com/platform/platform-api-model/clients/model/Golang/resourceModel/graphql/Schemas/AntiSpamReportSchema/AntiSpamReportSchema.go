package AntiSpamReportSchema

import (
	"github.com/graphql-go/graphql"
)

// SOCAntiSpamRptData  schema with AntiSpamData array
type SOCAntiSpamRptData struct {
	PartnerID          string `json:"partnerId"`
	SiteID             string `json:"siteId"`
	ClientID           string `json:"clientId"`
	AntiSpamReportData []AntiSpamDomainData
}

// AntiSpamDomainData  schema
type AntiSpamDomainData struct {
	DomainName string `json:"domainName"`
	SPF        bool   `json:"spf"`
	DMARC      bool   `json:"dmarc"`
}

//SOCAntiSpamRptDataType : SOCAntiSpamRptData GraphQL Schema
var SOCAntiSpamRptDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "socAntiSpamRptDataType",
	Fields: graphql.Fields{
		//define the  partnerID field type as string to output correct format data
		"partnerID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCAntiSpamRptData); ok {
					//return the partnerID field when no errors
					return CurData.PartnerID, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},

		//define the  siteID field type as string to output correct format data
		"siteID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCAntiSpamRptData); ok {
					//return the IndustrySector field when no errors
					return CurData.SiteID, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},
		//define the  clientID field type as string to output correct format data
		"clientID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCAntiSpamRptData); ok {
					//return the clientID field when no errors
					return CurData.ClientID, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},
		//define the  AntiSpamReportData array field of type AntiSpamDomainData to output correct format data
		"AntiSpamReportData": &graphql.Field{
			Type: graphql.NewList(AntiSpamDomainDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SOCAntiSpamRptData); ok {
					//return the domainName,spf,dmarc data which is part of array field when no errors
					return CurData.AntiSpamReportData, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},
	},
}) //SOCAntiSpamRpt ends

//AntiSpamDomainDataType : Antispam data for domainname,spf and dmarc GraphQL Schema
var AntiSpamDomainDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AntiSpamDomainData",
	Fields: graphql.Fields{
		"domainName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiSpamDomainData); ok {
					return CurData.DomainName, nil
				}
				return nil, nil
			},
		},
		"spf": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiSpamDomainData); ok {
					return CurData.SPF, nil
				}
				return nil, nil
			},
		},
		"dmarc": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiSpamDomainData); ok {
					return CurData.DMARC, nil
				}
				return nil, nil
			},
		},
	},
}) //AntiSpamDomainDataType ends
