package ticket

// response type and description
const (
	TicketType               = 5001
	TicketDesc               = "Ticket response for trigger table"
	AlertType                = 5002
	AlertDesc                = "Alert response for trigger table"
	TicketNotesType          = 5003
	TicketNotesDesc          = "Ticket notes response for trigger table"
	AlertNotesType           = 5004
	AlertNotesDesc           = "Alert notes response for trigger table"
	ContactTicketMappingType = 5005
	ContactTicketMappingDesc = "ContactTicketMapping response for trigger table"
)

// operation type
const (
	Unknown = 9000
	Create  = 1
	Update  = 2
	Delete  = 3
)
