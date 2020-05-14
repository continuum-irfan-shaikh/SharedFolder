package ticket

// CherwellJournal represent Cherwell notes object
type CherwellJournal struct {
	NoteID         string `json:"NoteId,omitempty"`
	NOCID          int64  `json:"NOCID,omitempty"`
	NOCType        string `json:"NOCType,omitempty"`
	IncidentID     string `json:"IncidentID"`
	Description    string `json:"Description,omitempty"`
	CreatedBy      string `json:"CreatedBy,omitempty"`
	CreatedOn      string `json:"CreatedOn,omitempty"`
	Source         string `json:"Source,omitempty"`
	Origin         string `json:"Origin,omitempty"`
	Type           string `json:"Type,omitempty"`
	CherwellNoteID string `json:"NoteID,omitempty"`
}
