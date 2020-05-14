package IndustryVerticalSchema

import (
	"github.com/graphql-go/graphql"
)

// IndustryVerticalData  schema with IndustryData array
type IndustryVerticalData struct {
	IndustryDatas []IndustryData
}

// IndustryData  schema
type IndustryData struct {
	IndustryID          int    `json:"industryID"`
	IndustrySector      string `json:"industrySector"`
	IndustryDescription string `json:"description"`
}

//IndustryDataListType : industryData GraphQL Schema
var IndustryDataListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "IndustryVerticalData",
	Fields: graphql.Fields{
		"industryData": &graphql.Field{
			Type: graphql.NewList(IndustryDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(IndustryVerticalData); ok {
					return CurData.IndustryDatas, nil
				}
				return nil, nil
			},
		},
	},
})

//IndustryDataType : IndustryWideVertical GraphQL Schema
var IndustryDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "industryWideVertical",
	Fields: graphql.Fields{
		//define the  industryID field type as Int to output correct format data
		"industryID": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(IndustryData); ok {
					//return the IndustryID field when no errors
					return CurData.IndustryID, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},

		//define the  industryID field type as Int to output correct format data
		"industrySector": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(IndustryData); ok {
					//return the IndustrySector field when no errors
					return CurData.IndustrySector, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},
		//define the  IndustryDescription field type as Int to output correct format data
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(IndustryData); ok {
					//return the IndustryDescription field when no errors
					return CurData.IndustryDescription, nil
				}
				//return nil in case of errors
				return nil, nil
			},
		},
	},
}) //IndustryVertical ends
