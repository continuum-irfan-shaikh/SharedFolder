package WebrootLicenseSummarySchema

import (
	"time"

	"github.com/graphql-go/graphql"
)

// WebrootLicenseSummaryData : WebrootLicenseSummaryData structure
type WebrootLicenseSummaryData struct {
	ProductName           string    `json:"productName"`
	Keycode               string    `json:"keycode"`
	AccessStatus          string    `json:"accessStatus"`
	AccessStatusGrantedAt time.Time `json:"accessStatusGrantedAt"`
	Email                 string    `json:"email"`
}

// WebrootLicenseSummaryType : WebrootLicenseSummaryType GraphQL Schema
var WebrootLicenseSummaryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootLicenseSummaryData",
	Fields: graphql.Fields{
		"productName": &graphql.Field{
			Type:        graphql.String,
			Description: "Product name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseSummaryData); ok {
					return CurData.ProductName, nil
				}
				return nil, nil
			},
		},
		"keyCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Key code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseSummaryData); ok {
					return CurData.Keycode, nil
				}
				return nil, nil
			},
		},
		"accessStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Access status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseSummaryData); ok {
					return CurData.AccessStatus, nil
				}
				return nil, nil
			},
		},
		"accessStatusGrantedAt": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Date of granting access",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseSummaryData); ok {
					return CurData.AccessStatusGrantedAt, nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type:        graphql.String,
			Description: "Email",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootLicenseSummaryData); ok {
					return CurData.Email, nil
				}
				return nil, nil
			},
		},
	},
})
