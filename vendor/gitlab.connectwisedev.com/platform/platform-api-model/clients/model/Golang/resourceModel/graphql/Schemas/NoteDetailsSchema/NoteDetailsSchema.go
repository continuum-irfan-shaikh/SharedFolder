package NoteDetails

import (
	"github.com/graphql-go/graphql"
)

//DeviceNotes : DeviceNotes Structure
type DeviceNotes struct {
	NoteID          int64	`json:"noteId"`
	Client       	string 	`json:"client"`
	MspName         string 	`json:"mspName"`
	UserID       	int64 	`json:"userId"`
	CreatedOn      	string 	`json:"createdOn"`
	NotesUpdatedOn 	string 	`json:"notesUpdatedOn"`
	Notes 		string  `json:"notes"`
	CreatedByUser	string 	`json:"createdByUser"`
	UpdatedByUser   string  `json:"updatedByUser"`
	ResourceName    string  `json:"resourceName"`
	ResFriendlyName string  `json:"resFriendlyName"`
	VisibleToMSP    int32   `json:"visibleToMSP"`
}

//DeviceNotesType : DeviceNotes GraphQL Schema
var DeviceNotesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "deviceNotes",
	Fields: graphql.Fields{
		"noteId": &graphql.Field{
			Type:        graphql.String,
			Description: "NoteId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.NoteID, nil
				}
				return nil, nil
			},
		},

		"client": &graphql.Field{
			Type:        graphql.String,
			Description: "Client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.Client, nil
				}
				return nil, nil
			},
		},

		"mspName": &graphql.Field{
			Type:        graphql.String,
			Description: "MspName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.MspName, nil
				}
				return nil, nil
			},
		},

		"userId": &graphql.Field{
			Type:        graphql.String,
			Description: "UserId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.UserID, nil
				}
				return nil, nil
			},
		},

		"createdOn": &graphql.Field{
			Type:        graphql.String,
			Description: "CreatedOn",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.CreatedOn, nil
				}
				return nil, nil
			},
		},

		"notesUpdatedOn": &graphql.Field{
			Type:        graphql.String,
			Description: "NotesUpdatedOn",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.NotesUpdatedOn, nil
				}
				return nil, nil
			},
		},

		"notes": &graphql.Field{
			Type:        graphql.String,
			Description: "Notes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.Notes, nil
				}
				return nil, nil
			},
		},

		"createdByUser": &graphql.Field{
			Type:        graphql.String,
			Description: "CreatedByUser",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.CreatedByUser, nil
				}
				return nil, nil
			},
		},
		
		"updatedByUser": &graphql.Field{
			Type:        graphql.String,
			Description: "UpdatedByUser",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.UpdatedByUser, nil
				}
				return nil, nil
			},
		},
		
		"resourceName": &graphql.Field{
			Type:        graphql.String,
			Description: "resourceName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.ResourceName, nil
				}
				return nil, nil
			},
		},
		
		"resFriendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "resFriendlyName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.ResFriendlyName, nil
				}
				return nil, nil
			},
		},
		
		"visibleToMSP": &graphql.Field{
			Type:        graphql.String,
			Description: "visibleToMSP",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceNotes); ok {
					return CurData.VisibleToMSP, nil
				}
				return nil, nil
			},
		},
	},
})

//GlobalNotes : GlobalNotes Structure
type GlobalNotes struct {
	NoteID          int64	`json:"noteId"`
	MspName         string 	`json:"mspName"`
	CreatedOn      	string 	`json:"createdOn"`
	NotesUpdatedOn  string 	`json:"notesUpdatedOn"`
	Notes 		string  `json:"notes"`
	CreatedByUser	string 	`json:"createdByUser"`
	UpdatedByUser   string  `json:"updatedByUser"`
	VisibleToMSP    int32   `json:"visibleToMSP"`
}

//GlobalNotesType : GlobalNotes GraphQL Schema
var GlobalNotesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "globalNotes",
	Fields: graphql.Fields{
		"noteId": &graphql.Field{
			Type:        graphql.String,
			Description: "NoteId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.NoteID, nil
				}
				return nil, nil
			},
		},

		"mspName": &graphql.Field{
			Type:        graphql.String,
			Description: "MspName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.MspName, nil
				}
				return nil, nil
			},
		},

		"createdOn": &graphql.Field{
			Type:        graphql.String,
			Description: "CreatedOn",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.CreatedOn, nil
				}
				return nil, nil
			},
		},

		"notesUpdatedOn": &graphql.Field{
			Type:        graphql.String,
			Description: "NotesUpdatedOn",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.NotesUpdatedOn, nil
				}
				return nil, nil
			},
		},

		"notes": &graphql.Field{
			Type:        graphql.String,
			Description: "Notes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.Notes, nil
				}
				return nil, nil
			},
		},

		"createdByUser": &graphql.Field{
			Type:        graphql.String,
			Description: "CreatedByUser",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.CreatedByUser, nil
				}
				return nil, nil
			},
		},
		
		"updatedByUser": &graphql.Field{
			Type:        graphql.String,
			Description: "UpdatedByUser",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.UpdatedByUser, nil
				}
				return nil, nil
			},
		},
		
		"visibleToMSP": &graphql.Field{
			Type:        graphql.String,
			Description: "visibleToMSP",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(GlobalNotes); ok {
					return CurData.VisibleToMSP, nil
				}
				return nil, nil
			},
		},
	},
})

//NoteDetails : NoteDetails Structure
type NoteDetails struct {
	RegID       	string		`json:"regId"`
	SiteID      	string 		`json:"siteId"`
	PartnerID   	string 		`json:"partnerId"`
	DNotes		[]DeviceNotes  	`json:"deviceLevelDetails"`
	SNotes 		[]DeviceNotes  	`json:"siteLevelDetails"`
	GNotes 		[]GlobalNotes  	`json:"globalLevelDetails"`
}

//NoteDetailsType : NoteDetails GraphQL Schema
var NoteDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "noteDetails",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "RegId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NoteDetails); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "SiteId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NoteDetails); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "PartnerId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NoteDetails); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"deviceLevelDetails": &graphql.Field{
			Type:        graphql.NewList(DeviceNotesType),
			Description: "deviceLevelDetails",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NoteDetails); ok {
					return CurData.DNotes, nil
				}
				return nil, nil
			},
		},

		"siteLevelDetails": &graphql.Field{
			Type:        graphql.NewList(DeviceNotesType),
			Description: "siteLevelDetails",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NoteDetails); ok {
					return CurData.SNotes, nil
				}
				return nil, nil
			},
		},

		"globalLevelDetails": &graphql.Field{
			Type:        graphql.NewList(GlobalNotesType),
			Description: "globalLevelDetails",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NoteDetails); ok {
					return CurData.GNotes, nil
				}
				return nil, nil
			},
		},
	},
})

//InsertNoteMutationData : InsertNoteMutationData Structure
type InsertNoteMutationData struct {
	PartnerID    int	`json:"partnerId"`
	RegID        int	`json:"RegId"`
	VisibleToMSP int	`json:"VisibleToMSP"`
	NoteLevel    int	`json:"NoteLevel"`
	MemberCode   string	`json:"MemberCode"`
	SiteCode     string	`json:"SiteCode"`
	Notes        string	`json:"Notes"`
	
}

//UpdateNoteMutationData : UpdateNoteMutationData Structure
type UpdateNoteMutationData struct {
	PartnerID    int	`json:"partnerId"`
	RegID        int	`json:"RegId"`
	NoteID       int	`json:"NoteId"`
	VisibleToMSP int	`json:"VisibleToMSP"`
	NoteLevel    int	`json:"NoteLevel"`
	Notes        string	`json:"Notes"`
	
}

//DeleteNoteMutationData : DeleteNoteMutationData Structure
type DeleteNoteMutationData struct {
	PartnerID    int	`json:"partnerId"`
	RegID        int	`json:"RegId"`
	NoteID       int	`json:"NoteId"`
	NoteLevel    int	`json:"NoteLevel"`
	
}
