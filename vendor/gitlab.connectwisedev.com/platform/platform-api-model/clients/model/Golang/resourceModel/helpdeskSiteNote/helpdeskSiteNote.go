package helpdeskSiteNote

// HelpdeskSiteNotePayload is a helpdesk site note kafka message
type HelpdeskSiteNotePayload struct {
	Action int              `json:"messageChangeAction"`
	Data   HelpdeskSiteNote `json:"data"`
}

// HelpdeskSiteNote is a model to hold Helpdesk Site Note
type HelpdeskSiteNote struct {
	LegacyID                 string `json:"legacyID,omitempty"`
	RecID                    string `json:"recID,omitempty"`
	CreatedDateTime          string `json:"createdDateTime,omitempty"`
	CreatedBy                string `json:"createdBy,omitempty"`
	CreatedByID              string `json:"createdByID,omitempty"`
	CreatedCulture           string `json:"createdCulture,omitempty"`
	LastModDateTime          string `json:"lastModifiedDateTime,omitempty"`
	LastModBy                string `json:"lastModBy,omitempty"`
	LastModByID              string `json:"lastModByID,omitempty"`
	OwnedBy                  string `json:"ownedBy,omitempty"`
	OwnedByID                string `json:"ownedByID,omitempty"`
	OwnedByTeam              string `json:"ownedByTeam,omitempty"`
	OwnedByTeamID            string `json:"ownedByTeamID,omitempty"`
	MessageID                string `json:"messageId,omitempty"`
	Message                  string `json:"message,omitempty"`
	ExpiredOn                string `json:"expiredOn,omitempty"`
	Subject                  string `json:"subject,omitempty"`
	Status                   string `json:"status,omitempty"`
	StatusUpdatedOn          string `json:"statusUpdatedOn,omitempty"`
	IncidentID               string `json:"incidentID,omitempty"`
	MessageScheduleID        string `json:"messageScheduleId,omitempty"`
	MsgScheduledOnDays       string `json:"msgScheduledOnDays,omitempty"`
	Zone                     string `json:"zone,omitempty"`
	ZoneDateTime             string `json:"zoneDateTime,omitempty"`
	SubcategoryID            string `json:"subcategoryID,omitempty"`
	FromDate                 string `json:"fromDate,omitempty"`
	ToDate                   string `json:"toDate,omitempty"`
	FromTimeToTime           string `json:"fromTime_ToTime,omitempty"`
	StatusUpdatedBy          string `json:"statusUpdatedBy,omitempty"`
	ScheduleType             string `json:"scheduleType,omitempty"`
	MsgPostID                string `json:"msgPostID,omitempty"`
	NocActionVariable        string `json:"nocActionVariable,omitempty"`
	JobID                    string `json:"job_ID,omitempty"`
	IsGeneratedBy            string `json:"isGeneratedBy,omitempty"`
	LastUpdatedbyTeam        string `json:"lastUpdatedbyTeam,omitempty"`
	ActionID                 string `json:"action_ID,omitempty"`
	MessageScheduleCreatedOn string `json:"messageScheduleCreatedOn,omitempty"`
	ScheduleCreatedBy        string `json:"scheduleCreatedBy,omitempty"`
	SchedulesGeneratedBy     string `json:"schedulesGeneratedBy,omitempty"`
	ScheduleStatus           string `json:"scheduleStatus,omitempty"`
	ScheduleStatusUpdatedBy  string `json:"scheduleStatusUpdatedBy,omitempty"`
	ScheduleLastUpdatedBy    string `json:"scheduleLastUpdatedBy,omitempty"`
	NoteType                 string `json:"noteType,omitempty"`
	SiteID                   int    `json:"siteId,omitempty"`
	ResourceID               int    `json:"resourceID,omitempty"`
	PartnerID                int    `json:"partnerID,omitempty"`
	IsActive                 bool   `json:"isActive,omitempty"`
	VisibletoMSP             bool   `json:"visibletoMSP,omitempty"`
	VisibletoNOC             bool   `json:"visibletoNOC,omitempty"`
}
