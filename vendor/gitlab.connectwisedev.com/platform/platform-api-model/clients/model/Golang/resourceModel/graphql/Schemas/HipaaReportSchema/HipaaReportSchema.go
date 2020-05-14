package HipaaReportSchema

import (
	"github.com/gocql/gocql"
	"github.com/graphql-go/graphql"
)

//HipaaReportData : HipaaReport Data Structure
type HipaaReportData struct {
	PartnerID            string                 `json:"partnerID"`
	PartnerName          string                 `json:"partnerName"`
	SiteID               string                 `json:"siteID"`
	SiteName             string                 `json:"siteName"`
	ClientID             string                 `json:"clientID"`
	LegalDisclaimerTitle string                 `json:"legalDisclaimerTitle"`
	LegalDisclaimer      string                 `json:"legalDisclaimer"`
	HipaaRequirementData []HipaaRequirementData `json:"requirements"`
}

//HipaaReportType : HipaaReport GraphQL Schema
var HipaaReportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HipaaReport",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"siteID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"siteName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
		"clientID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"legalDisclaimerTitle": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.LegalDisclaimerTitle, nil
				}
				return nil, nil
			},
		},
		"legalDisclaimer": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.LegalDisclaimer, nil
				}
				return nil, nil
			},
		},
		"partnerName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.PartnerName, nil
				}
				return nil, nil
			},
		},
		"requirements": &graphql.Field{
			Type: graphql.NewList(HipaaRequirementType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaReportData); ok {
					return CurData.HipaaRequirementData, nil
				}
				return nil, nil
			},
		},
	},
})

//HipaaRequirementData : Hipaa Requirement Data Structure
type HipaaRequirementData struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	RuleName    string     `json:"ruleName"`
	Type        string     `json:"type"`
	Count       []Count    `json:"counts"`
	Categories  []Category `json:"categories"`
}

// Category :Category Data Structure
type Category struct {
	ID                    string                  `json:"id"`
	Name                  string                  `json:"name"`
	Description           string                  `json:"description"`
	Type                  string                  `json:"type"`
	TotalEntities         int                     `json:"totalEntities"`
	NotRespondingEntities EntitiesObject          `json:"notRespondingEntities"`
	ViolatedEntities      EntitiesObject          `json:"violatedEntities"`
	InActiveEntities      EntitiesObject          `json:"inActiveEntities"`
	ExclusionList         ExclusionEntitiesObject `json:"exclusionList"`
}

// EntitiesObject :EntitiesObject Data Structure
type EntitiesObject struct {
	Entities []Entity `json:"entities"`
}

// Entity :Entity Data Structure
type Entity struct {
	ID           string `json:"id"`
	Type         string `json:"entityType,omitempty"`
	FriendlyName string `json:"friendlyName"`
	SystemName   string `json:"systemName"`
	AccountType  string `json:"accountType"`
}

//HipaaRequirementType : Hipaa Requirement GraphQL Schema
var HipaaRequirementType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HipaaRequirement",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaRequirementData); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaRequirementData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaRequirementData); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},

		"ruleName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaRequirementData); ok {
					return CurData.RuleName, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaRequirementData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},

		"counts": &graphql.Field{
			Type: graphql.NewList(CountType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaRequirementData); ok {
					return CurData.Count, nil
				}
				return nil, nil
			},
		},
		"categories": &graphql.Field{
			Type: graphql.NewList(CategoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaRequirementData); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},
	},
})

//Count : Count Entities Data Structure
type Count struct {
	Entity         string `json:"entity"`
	TotalCount     int64  `json:"totalCount"`
	ViolationCount int64  `json:"violationCount"`
}

//CategoryType : Category GraphQL Schema
var CategoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "categories",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
		"totalEntities": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.TotalEntities, nil
				}
				return nil, nil
			},
		},
		"exclusionList": &graphql.Field{
			Type: ExclusionEntitiesObjType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.ExclusionList, nil
				}
				return nil, nil
			},
		},
		"violatedEntities": &graphql.Field{
			Type: EntitiesType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.ViolatedEntities, nil
				}
				return nil, nil
			},
		},
		"notRespondingEntities": &graphql.Field{
			Type: EntitiesType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.NotRespondingEntities, nil
				}
				return nil, nil
			},
		},
		"inActiveEntities": &graphql.Field{
			Type: EntitiesType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Category); ok {
					return CurData.InActiveEntities, nil
				}
				return nil, nil
			},
		},
	},
})

//EntitiesType is Entities Type
var EntitiesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EntitiesType",
	Fields: graphql.Fields{
		"entities": &graphql.Field{
			Type: graphql.NewList(EntityType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EntitiesObject); ok {
					return CurData.Entities, nil
				}
				return nil, nil
			},
		},
	},
})

//EntityType : Entity GraphQL Schema
var EntityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "entities",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Entity); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"friendlyName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Entity); ok {
					return CurData.FriendlyName, nil
				}
				return nil, nil
			},
		},
		"systemName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Entity); ok {
					return CurData.SystemName, nil
				}
				return nil, nil
			},
		},
		"accountType": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Entity); ok {
					return CurData.AccountType, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Entity); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
	},
})

//CountType : Count Entities GraphQL Schema
var CountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Count",
	Fields: graphql.Fields{
		"entity": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Count); ok {
					return CurData.Entity, nil
				}
				return nil, nil
			},
		},
		"totalCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Count); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
		"violationCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Count); ok {
					return CurData.ViolationCount, nil
				}
				return nil, nil
			},
		},
	},
})

//HipaaExclusion struct to store response of profiling ms api
type HipaaExclusion struct {
	Categories    []HipaaExclusionCategory `json:"categories"`
	TotalEntities []Entity                 `json:"entities"`
}

//HipaaExclusionCategory to store hipaa exclusion category received from profiling api
type HipaaExclusionCategory struct {
	CategoryID    string   `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	ExclusionList []string `json:"exclusionList"`
}

//HipaaExclusionData json returned from graph ql
type HipaaExclusionData struct {
	Categories []HipaaExclusionCategoryData `json:"categories"`
	Entities   []Entity                     `json:"entities"`
}

//HipaaExclusionCategoryData structure returned from graph ql
type HipaaExclusionCategoryData struct {
	CategoryID      string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	SkippedEntities []Entity `json:"skippedEntities"`
}

//HipaaExclusionCategoryType is type for HipaaExclusionCategory
var HipaaExclusionCategoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HipaaExclusionCategory",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaExclusionCategoryData); ok {
					return CurData.CategoryID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaExclusionCategoryData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaExclusionCategoryData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},
		"skippedEntities": &graphql.Field{
			Type: graphql.NewList(EntityType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaExclusionCategoryData); ok {
					return CurData.SkippedEntities, nil
				}
				return nil, nil
			},
		},
	},
})

//HipaaExclusionType is type for HipaaExclusion
var HipaaExclusionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HipaaExclusion",
	Fields: graphql.Fields{
		"categories": &graphql.Field{
			Type: graphql.NewList(HipaaExclusionCategoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaExclusionData); ok {
					return CurData.Categories, nil
				}
				return nil, nil
			},
		},
		"entities": &graphql.Field{
			Type: graphql.NewList(EntityType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(HipaaExclusionData); ok {
					return CurData.Entities, nil
				}
				return nil, nil
			},
		},
	}})

// ExclusionEntitiesObject :Category Data Structure
type ExclusionEntitiesObject struct {
	ExclusionEntities []ExclusionEntities `json:"exclusionEntities"`
}

// ExclusionEntities :ExclusionEntitiesType Data Structure
type ExclusionEntities struct {
	ID       gocql.UUID `json:"id"`
	Name     string     `json:"name"`
	Reason   string     `json:"reason"`
	Entities []Entity   `json:"entities"`
}

//ExclusionEntitiesType : Hipaa Requirement GraphQL Schema
var ExclusionEntitiesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExclusionEntities",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionEntities); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionEntities); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"reason": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionEntities); ok {
					return CurData.Reason, nil
				}
				return nil, nil
			},
		},
		"entities": &graphql.Field{
			Type: graphql.NewList(EntityType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionEntities); ok {
					return CurData.Entities, nil
				}
				return nil, nil
			},
		},
	},
})

//ExclusionEntitiesObjType is Entities Type
var ExclusionEntitiesObjType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExclusionEntitiesObject",
	Fields: graphql.Fields{
		"exclusionEntities": &graphql.Field{
			Type: graphql.NewList(ExclusionEntitiesType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ExclusionEntitiesObject); ok {
					return CurData.ExclusionEntities, nil
				}
				return nil, nil
			},
		},
	},
})
