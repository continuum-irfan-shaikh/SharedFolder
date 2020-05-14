package SiteHelpDeskSchema

import (
	"github.com/graphql-go/graphql"
)

//HelpDeskSummaryData : HelpDeskSummaryData struct
type HelpDeskSummaryData struct {
	Coverage     string
	SiteID		 string
}

//SiteHelpDeskSummaryDataType : SiteHelpDeskSummaryDataType GraphQL Schema
var SiteHelpDeskSummaryDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SiteHelpDeskSummaryData",
	Fields: graphql.Fields{
		"Coverage": &graphql.Field{
			Type:        graphql.String,
			Description: "help desk coverage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpDeskSummaryData); ok {
					return CurData.Coverage, nil
				}
				return nil, nil
			},
		},
	},
})

//HelpDeskCoverageProduct : HelpDeskCoverageProduct struct
type HelpDeskCoverageProduct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//HelpDeskCoverageProductType : HelpDeskCoverageProductType GraphQL Schema
var HelpDeskCoverageProductType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HelpDeskCoverageProductType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "help desk coverage id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpDeskCoverageProduct); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "help desk coverage id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpDeskCoverageProduct); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
	},
})

//HelpDeskCoverageProductsData : HelpDeskCoverageProductsData struct
type HelpDeskCoverageProductsData struct {
	Outdata     []HelpDeskCoverageProduct
}

//SiteHelpDeskCoverageProductsDataType : SiteHelpDeskCoverageProductsDataType GraphQL Schema
var SiteHelpDeskCoverageProductsDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SiteHelpDeskCoverageProductsDataType",
	Fields: graphql.Fields{
		"outdata": &graphql.Field{
			Type:        graphql.NewList(HelpDeskCoverageProductType),
			Description: "help desk coverage outdata",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HelpDeskCoverageProductsData); ok {
					return CurData.Outdata, nil
				}
				return nil, nil
			},
		},
	},
})
