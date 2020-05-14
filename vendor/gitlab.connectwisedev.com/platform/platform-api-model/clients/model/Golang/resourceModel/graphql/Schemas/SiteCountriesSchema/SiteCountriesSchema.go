package SiteCountriesSchema

import (
	"github.com/graphql-go/graphql"
)

//CountryData : CountryData Structure
type CountryData struct {
	CountryID   int64  `json:"ID"`
	Country     string `json:"Country"`
	CountryCode string `json:"CountryCode"`
}

//CountryType : Country GraphQL Schema
var CountryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CountryDetails",
	Fields: graphql.Fields{
		"CountryCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Country Code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CountryData); ok {
					return CurData.CountryCode, nil
				}
				return nil, nil
			},
		},

		"Country": &graphql.Field{
			Type:        graphql.String,
			Description: "Country",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CountryData); ok {
					return CurData.Country, nil
				}
				return nil, nil
			},
		},

		"ID": &graphql.Field{
			Type:        graphql.String,
			Description: "CountryID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CountryData); ok {
					return CurData.CountryID, nil
				}
				return nil, nil
			},
		},
	},
})

//SiteCountriesData : SiteCountriesData Structure
type SiteCountriesData struct {
	Status      int64         `json:"status"`
	CountryList []CountryData `json:"outdata"`
}

//SiteCountriesType : SiteCountriesType GraphQL Schema
var SiteCountriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SiteCountries",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteCountriesData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        graphql.NewList(CountryType),
			Description: "Site Countries",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SiteCountriesData); ok {
					return CurData.CountryList, nil
				}
				return nil, nil
			},
		},
	},
})
