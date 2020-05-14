package entities

// UserSites describes user sites data for particular partner ID and particular user ID
type UserSites struct {
	PartnerID string
	UserID    string
	SiteIDs   []int64
}
