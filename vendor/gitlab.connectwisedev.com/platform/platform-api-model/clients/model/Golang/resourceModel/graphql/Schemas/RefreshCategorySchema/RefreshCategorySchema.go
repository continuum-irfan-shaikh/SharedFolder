package RefreshCategorySchema

import (
	"github.com/graphql-go/graphql"
)

//RefreshCategoryData for response from the post request
type RefreshCategoryData struct {
	PartnerID            string `json:"partnerID"`
	SiteID               string `json:"siteID"`
	SystemCurrentTime    int64  `json:"systemCurrentTime"`
	LastRefreshedTime    int64  `json:"lastRefreshedTime"`
	ExecutionStatus      string `json:"executionStatus"`
	RefreshIntervalInMin int    `json:"refreshIntervalInMin"`
}

//RefreshCategoryPostResponse for response from the post request
type RefreshCategoryPostResponse struct {
	StatusCode int `json:"statusCode"`
}

//RefreshCategoryPostRequest for response from the post request
type RefreshCategoryPostRequest struct {
	PartnerID string `json:"partnerId"`
	SiteID    string `json:"siteId"`
}

//RefreshCategoryType for graphQL structure
var RefreshCategoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RefreshCategoryData",
	Fields: graphql.Fields{
		"lastRefreshedTime": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RefreshCategoryData); ok {
					return CurData.LastRefreshedTime, nil
				}
				return nil, nil
			},
		},
		"executionStatus": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RefreshCategoryData); ok {
					return CurData.ExecutionStatus, nil
				}
				return nil, nil
			},
		},
		"systemCurrentTime": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RefreshCategoryData); ok {
					return CurData.SystemCurrentTime, nil
				}
				return nil, nil
			},
		},
		"refreshIntervalInMin": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RefreshCategoryData); ok {
					return CurData.RefreshIntervalInMin, nil
				}
				return nil, nil
			},
		},
	},
})

//RefreshCategoriesPostQueryType for graphQL structure
var RefreshCategoriesPostQueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RefreshCategoryPostResponse",
	Fields: graphql.Fields{
		"statusCode": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(RefreshCategoryPostResponse); ok {
					return CurData.StatusCode, nil
				}
				return nil, nil
			},
		},
	},
})
