package ExclusionsSchema

import (
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/HipaaReportSchema"
	"github.com/gocql/gocql"
	"github.com/graphql-go/graphql"
)

//ExclusionsResponse for response from the post request
type ExclusionsResponse struct {
	StatusCode  int    `json:"statusCode"`
	Description string `json:"description"`
}

//ExclusionsResponseType for graphQL structure
var ExclusionsResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExclusionsResponse",
	Fields: graphql.Fields{
		"statusCode": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsResponse); ok {
					return CurData.StatusCode, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsResponse); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
	},
})

//ExclusionHistory for response for history from REST API
type ExclusionHistory struct {
	ExclusionHistory []History `json:"history"`
	CategoryType     string    `json:"categoryType"`
}

//History stores single transaction on an exclusion
type History struct {
	Action        string   `json:"action"`
	ExclusionList []string `json:"exclusionList"`
	Reason        string   `json:"reason"`
	UpdatedAt     int64    `json:"updatedAt"`
	UpdatedBy     string   `json:"updatedBy"`
}

//HistoryResponse for history log response
type HistoryResponse struct {
	Action    string `json:"action"`
	Entities  string `json:"entities"`
	Reason    string `json:"reason"`
	UpdatedAt string `json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`
}

//ExclusionHistoryData for response for history
type ExclusionHistoryData struct {
	ExclusionHistory string `json:"history"`
}

//ExclusionHistoryType is graphQL object for history
var ExclusionHistoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExclusionHistory",
	Fields: graphql.Fields{
		"history": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionHistoryData); ok {
					return CurData.ExclusionHistory, nil
				}
				return nil, nil
			},
		},
	}})

//EntityExclusions struct to store response of profiling ms api
type EntityExclusions struct {
	SiteID        string                     `json:"siteID,omitempty"`
	Categories    []ExclusionCategory        `json:"categories"`
	TotalEntities []HipaaReportSchema.Entity `json:"entities,omitempty"`
}

//ModifyExclusions struct to modify/delete Exclusion
type ModifyExclusions struct {
	SiteID          string       `json:"siteID"`
	CategoryID      string       `json:"categoryID"`
	ExclusionsGroup []Exclusions `json:"exclusions"`
}

//DeleteExclusionAPIReq request
type DeleteExclusionAPIReq struct {
	CategoryID  string       `json:"categoryID"`
	ExclusionID []gocql.UUID `json:"exclusionID"`
}

//ExclusionCategory to store  exclusion category received from profiling api
type ExclusionCategory struct {
	CategoryID      string       `json:"id"`
	Name            string       `json:"name,omitempty"`
	Type            string       `json:"type,omitempty"`
	ExclusionsGroup []Exclusions `json:"exclusions"`
}

//Exclusions stores exclusion category
type Exclusions struct {
	ID            gocql.UUID `json:"id"`
	Name          string     `json:"name,omitempty"`
	Reason        string     `json:"reason,omitempty"`
	ExclusionList []string   `json:"exclusionList"`
	Criteria      string     `json:"criteria"`
	UpdatedAt     int64      `json:"updatedAt,omitempty"`
	UpdatedBy     string     `json:"updatedBy,omitempty"`
}

//ExclusionsData for the graphql response
type ExclusionsData struct {
	ID              gocql.UUID                 `json:"id"`
	Name            string                     `json:"name,omitempty"`
	Reason          string                     `json:"reason,omitempty"`
	SkippedEntities []HipaaReportSchema.Entity `json:"skippedEntities,omitempty"`
	Criteria        string                     `json:"criteria"`
	UpdatedAt       int64                      `json:"updatedAt,omitempty"`
	UpdatedBy       string                     `json:"updatedBy"`
}

//ExclusionsDataType for graphQL exclusion groups for categories
var ExclusionsDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExclusionsInfo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsData); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"reason": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsData); ok {
					return CurData.Reason, nil
				}
				return nil, nil
			},
		},
		"criteria": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsData); ok {
					return CurData.Criteria, nil
				}
				return nil, nil
			},
		},
		"skippedEntities": &graphql.Field{
			Type: graphql.NewList(HipaaReportSchema.EntityType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsData); ok {
					return CurData.SkippedEntities, nil
				}
				return nil, nil
			},
		},
		"updatedBy": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsData); ok {
					return CurData.UpdatedBy, nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionsData); ok {
					return CurData.UpdatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

//EntityExclusionsData json returned from graph ql
type EntityExclusionsData struct {
	Categories []ExclusionCategoryData    `json:"categories"`
	Entities   []HipaaReportSchema.Entity `json:"entities,omitempty"`
}

//ExclusionCategoryData structure returned from graph ql
type ExclusionCategoryData struct {
	CategoryID      string           `json:"id"`
	Name            string           `json:"name,omitempty"`
	Type            string           `json:"type,omitempty"`
	ExclusionsGroup []ExclusionsData `json:"exclusions"`
}

//ExclusionCategoryType is type for ExclusionCategory
var ExclusionCategoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExclusionCategory",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionCategoryData); ok {
					return CurData.CategoryID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionCategoryData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionCategoryData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
		"exclusions": &graphql.Field{
			Type: graphql.NewList(ExclusionsDataType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionCategoryData); ok {
					return CurData.ExclusionsGroup, nil
				}
				return nil, nil
			},
		},
	},
})

//EntityExclusionsType is type for Exclusion
var EntityExclusionsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EntityExclusion",
	Fields: graphql.Fields{
		"categories": &graphql.Field{
			Type: graphql.NewList(ExclusionCategoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EntityExclusionsData); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},
		"entities": &graphql.Field{
			Type: graphql.NewList(HipaaReportSchema.EntityType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EntityExclusionsData); ok {
					return CurData.Entities, nil
				}
				return nil, nil
			},
		},
	}})
