package ProfileProtectSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
	"github.com/graphql-go/graphql"
)

//SpyCloudResult type for api response
type SpyCloudResult struct {
	PartnerID   string        `json:"partnerId"`
	ClientID    string        `json:"clientId"`
	SiteID      string        `json:"siteId"`
	Domain      string        `json:"domain"`
	Email       string        `json:"email"`
	PartnerName string        `json:"partnername"`
	UserData    []SpyCloudRow `json:"userData"`
}

//SpyCloudRow type for api response
type SpyCloudRow struct {
	EmailUserName     string    `json:"email_username"`
	Domain            string    `json:"domain"`
	Password          string    `json:"password"`
	PublishDate       time.Time `json:"spycloud_publish_date"`
	PasswordPlainText string    `json:"password_plaintext"`
	EmailDomain       string    `json:"email_domain"`
	PasswordType      string    `json:"password_type"`
	Email             string    `json:"email"`
	Site              string    `json:"site"`
}

//PartnerDetailsType type for partner details from rmm api
type PartnerDetailsType struct {
	Outdata []PartnerDetails `json:"outdata"`
}

//PartnerDetails struct for partner name
type PartnerDetails struct {
	PartnerName string `json:"Partner_Name"`
	MemberName  string `json:"MemberName"`
}

//SpyCloudResultWrapperType wrapper for spycloud result
var SpyCloudResultWrapperType = graphql.NewObject(graphql.ObjectConfig{
	Name: "userDetailList",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudResult); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"siteID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudResult); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"clientID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudResult); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"partnername": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudResult); ok {
					return CurData.PartnerName, nil
				}
				return nil, nil
			},
		},
		"domain": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudResult); ok {
					return CurData.Domain, nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudResult); ok {
					return CurData.Email, nil
				}
				return nil, nil
			},
		},
		"userDetailList": &graphql.Field{
			Type: DarkWebBreachConnectionDefinition.ConnectionType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)

				if CurData, ok := p.Source.(SpyCloudResult); ok {
					var arrayOut []interface{}
					for ind := range CurData.UserData {
						arrayOut = append(arrayOut, CurData.UserData[ind])
					}
					// return CurData.UserData, nil
					return Relay.ConnectionFromArray(arrayOut, args, ""), nil

				}
				return nil, nil
			},
		},
	},
})

//DarkWebBreachPasswordType : dark web breach password GraphQL Schema
var DarkWebBreachPasswordType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DarkWebPasswordBreach",
	Fields: graphql.Fields{
		"email_username": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.EmailUserName, nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.Email, nil
				}
				return nil, nil
			},
		},
		"email_domain": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.EmailDomain, nil
				}
				return nil, nil
			},
		},
		"password_type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.PasswordType, nil
				}
				return nil, nil
			},
		},
		"password": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.Password, nil
				}
				return nil, nil
			},
		},
		"password_plaintext": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.PasswordPlainText, nil
				}
				return nil, nil
			},
		},
		"spycloud_publish_date": &graphql.Field{
			Type: CustomDataTypes.DateTimeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.PublishDate, nil
				}
				return nil, nil
			},
		},
		"site": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SpyCloudRow); ok {
					return CurData.Site, nil
				}
				return nil, nil
			},
		},
	},
})

// DarkWebBreachConnectionDefinition : DarkWebBreachConnectionDefinition structure
var DarkWebBreachConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "userDetailList",
	NodeType: DarkWebBreachPasswordType,
})
