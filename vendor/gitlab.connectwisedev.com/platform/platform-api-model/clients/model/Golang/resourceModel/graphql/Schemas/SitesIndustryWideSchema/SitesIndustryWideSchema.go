package SitesIndustryWideSchema

import (
	"github.com/graphql-go/graphql"
)

//SitesIndustryWideData type for api response
type SitesIndustryWideData struct {
	IndustryID          int    `json:"industryID"`
	IndustryDescription string `json:"description"`
}

//SiteIndustryMapping struct to store request body for Site industry mapping API
type SiteIndustryMapping struct {
	PartnerID  string `json:"MemberID"`
	SiteID     string `json:"SiteID"`
	IndustryID string `json:"IndustryID"`
	UserID     string `json:"CreatedUserId"`
}

//SitesIndustryWideType : SitesIndustryWide GraphQL Schema
var SitesIndustryWideType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesIndustryWide",
	Fields: graphql.Fields{
		//define the  industryid field type as Int to output correct format data
		"industryID": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesIndustryWideData); ok {
					//return the IndustryID field when no errors
					return CurData.IndustryID, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},

		//define the  industryid field type as String to output correct format data
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesIndustryWideData); ok {
					//return the IndustryDescription field when no errors
					return CurData.IndustryDescription, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},
	},
}) //SitesIndustryWideType ends
