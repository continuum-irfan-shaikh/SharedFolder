package modelMocks

import (
	"context"
	"errors"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"github.com/gocql/gocql"
)

// UserSitesRepoMock stores UserSites mocks
type UserSitesRepoMock struct {
	repo map[string][]entities.UserSites
}

// NewUserSitesRepoMock returns a new UserSitesRepoMock filled with data if needed
func NewUserSitesRepoMock(isFilled bool) UserSitesRepoMock {
	mock := UserSitesRepoMock{}

	mock.repo = make(map[string][]entities.UserSites)
	if isFilled {
		for _, uSite := range DefaultUserSites {
			mock.repo[uSite.PartnerID] = append(mock.repo[uSite.PartnerID], uSite)
		}
	}
	return mock
}

// InsertUserSites inserts UserSites in Cassandra
func (us UserSitesRepoMock) InsertUserSites(ctx context.Context, partnerID, userID string, siteIDs []int64) error {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return errors.New("cassandra is down")
	}
	us.repo[partnerID] = append(us.repo[partnerID], entities.UserSites{
		PartnerID: partnerID,
		UserID:    userID,
		SiteIDs:   siteIDs,
	})
	return nil
}

// GetUserSites returns UserSites found by partnerID and userID
func (us UserSitesRepoMock) Sites(ctx context.Context, partnerID, userID string) (uSites entities.UserSites, err error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return uSites, errors.New("cassandra is down")
	}

	userSitesList := us.repo[partnerID]
	for _, item := range userSitesList {
		if item.UserID == userID {
			return uSites, nil
		}
	}
	return
}

// InsertSitesByTaskID ...
func (us UserSitesRepoMock) InsertSitesByTaskID(ctx context.Context, taskID gocql.UUID, siteIDs []string) error {
	return nil
}

// GetSitesByTaskID ...
func (us UserSitesRepoMock) GetSitesByTaskID(ctx context.Context, taskID gocql.UUID) (siteIDs []string, err error) {
	return nil, nil
}

// Endpoints - select Endpoints by partnerID and sites
func (us UserSitesRepoMock) Endpoints(ctx context.Context, partnerID string, siteIDs []string) ([]entities.Endpoints, error) {
	return nil, nil
}

// EndpointsByPartner ..
func (us UserSitesRepoMock) EndpointsByPartner(ctx context.Context, partnerID string) ([]entities.Endpoints, error) {
	return nil, nil
}

// SaveEndpoints - inserts endpoints to DB
func (us UserSitesRepoMock) SaveEndpoints(ctx context.Context, ep []entities.Endpoints) {
}

// DefaultUserSites is a list of UserSites mocks
var DefaultUserSites = []entities.UserSites{
	{
		PartnerID: "partner1",
		UserID:    "user1-1",
		SiteIDs: []int64{
			1 - 11111,
			1 - 22222,
			1 - 33333,
		},
	},
	{
		PartnerID: "partner1",
		UserID:    "user1-2",
		SiteIDs: []int64{
			1 - 22222,
			1 - 33333,
		},
	},
	{
		PartnerID: "partner2",
		UserID:    "user2-1",
		SiteIDs: []int64{
			2 - 11111,
			2 - 22222,
			2 - 33333,
		},
	},
}
