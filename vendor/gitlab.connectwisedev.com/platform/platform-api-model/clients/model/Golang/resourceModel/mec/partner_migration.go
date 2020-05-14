package mec

const (
	//HeaderPartnerMigration is the type of message
	HeaderPartnerMigration string = "PARTNER-MIGRATION"
)

//PartnerMigration is the message that will get published when a site gets migrated from one to another partner
type PartnerMigration struct {
	SiteID        string
	FromPartnerID string
	ToPartnerID   string
}
