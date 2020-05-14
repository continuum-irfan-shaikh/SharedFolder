package ticket

// CherwellReporterIncidentMapping represent Cherwell ContactTicketMapping object
type CherwellReporterIncidentMapping struct {
	NOCTicketID int64  `json:"NOCTicketID,omitempty"`
	ContactID   string `json:"ContactID,omitempty"`
	Origin      string `json:"Origin,omitempty"`
	IncidentID  string `json:"IncidentID"`
	CreatedBy   string `json:"CreatedBy,omitempty"`
	CreatedOn   string `json:"CreatedOn,omitempty"`
}
