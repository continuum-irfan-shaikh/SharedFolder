package itsm_integration

// JournalNote represents journal note to incident in Cherwell
type JournalNote struct {
	// ID is note ID
	ID string `json:"id,omitempty"`
	// CreatedAt is note creation time
	CreatedAt string `json:"createdAt"`
	// CreatedBy is note creator's name
	CreatedBy string `json:"createdBy"`
	// Details is note content
	Details string `json:"details"`
	// ExternalNotes is external note details
	ExternalNotes string `json:"externalNotes"`
	// IncidentID is note's incident ID
	IncidentID string `json:"incidentId"`
	// ParentBOID is business object ID of parent.
	ParentBOID string `json:"-"`
	//InternalFlag
	InternalFlag bool   `json:"internalFlag"`
	PSANoteID    string `json:"PSANoteID"`
	PSANoteType  string `json:"PSANoteType"`
	Origin       string `json:"origin"`
}
