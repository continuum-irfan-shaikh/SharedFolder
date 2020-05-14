package WebrootLicenseDetailsSchema

import (
	"github.com/graphql-go/graphql"
)

// WebrootLicenseDetailsData : WebrootLicenseDetailsData structure
type WebrootLicenseDetailsData struct {
	PartnerID          string `json:"partnerID"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	CompanyName        string `json:"companyName"`
	CustomerEmail      string `json:"customerEmail"`
	Address1           string `json:"address1"`
	Address2           string `json:"address2"`
	Country            string `json:"country"`
	State              string `json:"state"`
	City               string `json:"city"`
	PostalCode         string `json:"postalCode"`
	AccessRequested    bool   `json:"accessRequested"`
	ProvisioningStatus string `json:"provisioningStatus"`
}

// WebrootLicenseDetailsType : WebrootLicenseDetailsType GraphQL Schema
var WebrootLicenseDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootLicenseDetailsData",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"firstName": &graphql.Field{
			Type:        graphql.String,
			Description: "First name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.FirstName, nil
				}
				return nil, nil
			},
		},
		"lastName": &graphql.Field{
			Type:        graphql.String,
			Description: "Last name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.LastName, nil
				}
				return nil, nil
			},
		},
		"companyName": &graphql.Field{
			Type:        graphql.String,
			Description: "Company name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.CompanyName, nil
				}
				return nil, nil
			},
		},
		"customerEmail": &graphql.Field{
			Type:        graphql.String,
			Description: "Customer email",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.CustomerEmail, nil
				}
				return nil, nil
			},
		},
		"address1": &graphql.Field{
			Type:        graphql.String,
			Description: "Address 1",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.Address1, nil
				}
				return nil, nil
			},
		},
		"address2": &graphql.Field{
			Type:        graphql.String,
			Description: "Address 2",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.Address2, nil
				}
				return nil, nil
			},
		},
		"country": &graphql.Field{
			Type:        graphql.String,
			Description: "Country",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.Country, nil
				}
				return nil, nil
			},
		},
		"state": &graphql.Field{
			Type:        graphql.String,
			Description: "State",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.State, nil
				}
				return nil, nil
			},
		},
		"city": &graphql.Field{
			Type:        graphql.String,
			Description: "City",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.City, nil
				}
				return nil, nil
			},
		},
		"postalCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Postal code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.PostalCode, nil
				}
				return nil, nil
			},
		},
		"accessRequested": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Is access admin requested",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.AccessRequested, nil
				}
				return nil, nil
			},
		},
		"provisioningStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Provisioning status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseDetailsData); ok {
					return CurData.ProvisioningStatus, nil
				}
				return nil, nil
			},
		},
	},
})
