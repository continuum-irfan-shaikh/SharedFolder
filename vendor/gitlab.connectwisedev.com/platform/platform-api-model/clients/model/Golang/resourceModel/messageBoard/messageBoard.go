package messageBoard

// MessageBoardPayload is a message board kafka message
type MessageBoardPayload struct {
	Action int          `json:"operation"`
	Data   MessageBoard `json:"data"`
}

// MessageBoard is a model to hold Message Board
type MessageBoard struct {
	RecID                    string `json:"id,omitempty"`
	LegacyID                 string `json:"legacyId,omitempty"`
	JobID                    string `json:"jobId,omitempty"`
	LastUpdatedByTeam        string `json:"lastUpdatedByTeam,omitempty"`
	ScheduleCreatedBy        string `json:"scheduleCreatedBy,omitempty"`
	SchedulesGeneratedBy     string `json:"schedulesGeneratedBy,omitempty"`
	ScheduleStatusUpdatedBy  string `json:"scheduleStatusUpdatedBy,omitempty"`
	ScheduleLastUpdatedBy    string `json:"scheduleLastUpdatedBy,omitempty"`
	ScheduleStatus           string `json:"scheduleStatus,omitempty"`
	CreatedOn                string `json:"createdOn,omitempty"`
	CreatedByID              string `json:"createdById,omitempty"`
	Message                  string `json:"message,omitempty"`
	ExpiredOn                string `json:"expiredOn,omitempty"`
	Subject                  string `json:"subject,omitempty"`
	Status                   string `json:"status,omitempty"`
	StatusUpdatedOn          string `json:"statusUpdatedOn,omitempty"`
	Zone                     string `json:"zone,omitempty"`
	ZoneDateTime             string `json:"zoneDateTime,omitempty"`
	FromDate                 string `json:"fromDate,omitempty"`
	ToDate                   string `json:"toDate,omitempty"`
	ScheduleType             string `json:"scheduleType,omitempty"`
	IsGeneratedBy            string `json:"isGeneratedBy,omitempty"`
	MessageScheduleCreatedOn string `json:"messageScheduleCreatedOn,omitempty"`
	Type                     string `json:"type,omitempty"`
	LastUpdatedBy            string `json:"lastUpdatedBy,omitempty"`
	CreatedBy                string `json:"createdBy,omitempty"`
	NocActionVariable        string `json:"nocActionVariable,omitempty"`
	StatusUpdatedBy          string `json:"statusUpdatedBy,omitempty"`
	IncidentID               string `json:"incidentId,omitempty"`
	MessageScheduleID        string `json:"messageScheduleId,omitempty"`
	ActionID                 string `json:"actionId,omitempty"`
	MsgPostID                string `json:"msgPostId,omitempty"`
	UpdatedOn                string `json:"updatedOn,omitempty"`
	MsgScheduledOnDays       string `json:"msgScheduledOnDays,omitempty"`
	FromTimeToTime           string `json:"fromTimeToTime,omitempty"`
	SubcategoryID            string `json:"subcategoryId,omitempty"`
	PartnerID                string `json:"partnerId,omitempty"`
	ClientID                 string `json:"clientId,omitempty"`
	SiteID                   string `json:"siteId,omitempty"`
	ResourceID               string `json:"endpointId,omitempty"`
	MessageID                string `json:"messageId,omitempty"`
	IsActive                 bool   `json:"isActive"`
	VisibletoMSP             bool   `json:"visibletoMSP"`
	VisibletoNOC             bool   `json:"visibletoNOC"`
}
