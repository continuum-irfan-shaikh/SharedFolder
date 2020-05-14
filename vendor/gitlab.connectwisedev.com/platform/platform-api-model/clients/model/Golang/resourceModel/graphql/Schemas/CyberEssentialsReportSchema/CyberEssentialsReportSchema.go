package CyberEssentialsReportSchema

import (
	"github.com/gocql/gocql"
	"github.com/graphql-go/graphql"
)

//CyberEssentialReportData : CyberEssentialReportData Data Structure
type CyberEssentialReportData struct {
	PartnerID                string                     `json:"partnerID"`
	PartnerName              string                     `json:"partnerName"`
	SiteID                   string                     `json:"siteID"`
	SiteName                 string                     `json:"siteName"`
	ClientID                 string                     `json:"clientID"`
	LegalDisclaimerTitle     string                     `json:"legalDisclaimerTitle"`
	LegalDisclaimer          string                     `json:"legalDisclaimer"`
	CyberEssentialCategories []CyberEssentialCategories `json:"cyberEssentialCategories"`
}

//CyberEssentialReportDataType : CyberEssentialReportData GraphQL Schema
var CyberEssentialReportDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CyberEssentialReport",
	Fields: graphql.Fields{
		"partnerID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"siteID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"siteName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
		"clientID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"legalDisclaimerTitle": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.LegalDisclaimerTitle, nil
				}
				return nil, nil
			},
		},
		"legalDisclaimer": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.LegalDisclaimer, nil
				}
				return nil, nil
			},
		},
		"partnerName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.PartnerName, nil
				}
				return nil, nil
			},
		},
		"cyberEssentialCategories": &graphql.Field{
			Type: graphql.NewList(CyberEssentialCategoriesType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialReportData); ok {
					return CurData.CyberEssentialCategories, nil
				}
				return nil, nil
			},
		},
	},
})

//CyberEssentialCategories : CyberEssential categories
type CyberEssentialCategories struct {
	Name      string      `json:"name"`
	Questions []Questions `json:"questions"`
}

//CyberEssentialCategoriesType is type for CyberEssentialCategories
var CyberEssentialCategoriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CyberEssentialCategories",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialCategories); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"questions": &graphql.Field{
			Type: graphql.NewList(QuestionsType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CyberEssentialCategories); ok {
					return CurData.Questions, nil
				}
				return nil, nil
			},
		},
	}})

//Source is the information required for graphQL for making further request
type Source struct {
	ServiceName string `json:"serviceName"`
	Attributes  string `json:"attributes"`
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

//Questions : object for cyber essentials questions
type Questions struct {
	ID             string     `json:"id"`
	Number         string     `json:"number"`
	Description    string     `json:"description"`
	Answer         string     `json:"answer"`
	Comment        string     `json:"comment"`
	Source         Source     `json:"source"`
	NotesTableData NotesTable `json:"notesTable"`
	Categories     []Category `json:"categories"`
}

//NotesTable for table structure for notes type question
type NotesTable struct {
	Description string   `json:"description"`
	ColumnNames string   `json:"columnNames"`
	Separator   string   `json:"separator"`
	Result      []string `json:"result"`
}

//NotesTableType is type for NotesTable
var NotesTableType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NotesTable",
	Fields: graphql.Fields{
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesTable); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
		"columnNames": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesTable); ok {
					return CurData.ColumnNames, nil
				}
				return nil, nil
			},
		},
		"separator": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesTable); ok {
					return CurData.Separator, nil
				}
				return nil, nil
			},
		},
		"result": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesTable); ok {
					return CurData.Result, nil
				}
				return nil, nil
			},
		},
	}})

//QuestionsType : Hipaa Requirement GraphQL Schema
var QuestionsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Questions",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Questions); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},

		"number": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Questions); ok {
					return CurData.Number, nil
				}
				return nil, nil
			},
		},

		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Questions); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},
		"answer": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Questions); ok {
					return CurData.Answer, nil
				}
				return nil, nil
			},
		},
		"comment": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Questions); ok {
					return CurData.Comment, nil
				}
				return nil, nil
			},
		},
		"notesTable": &graphql.Field{
			Type: NotesTableType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Questions); ok {
					return CurData.NotesTableData, nil
				}
				return nil, nil
			},
		},

		"categories": &graphql.Field{
			Type: graphql.NewList(CategoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Questions); ok {
					return CurData.Categories, nil
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
