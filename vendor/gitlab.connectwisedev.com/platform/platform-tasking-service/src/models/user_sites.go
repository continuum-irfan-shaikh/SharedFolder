package models

//go:generate mockgen -destination=../mocks/mocks-gomock/userSitesPersistence_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models UserSitesPersistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	en "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	userSitesFields    = `partner_id, user_id, site_ids`
	insertUserSitesCQL = `INSERT INTO user_sites (` + userSitesFields + `) VALUES (?, ?, ?)`
	selectUserSitesCQL = `SELECT ` + userSitesFields + ` FROM user_sites WHERE partner_id = ? AND user_id = ?`

	insertTaskSitesCQL     = `INSERT INTO task_sites (task_id , sites) VALUES (?, ?)`
	selectSitesByTaskIDCQL = `SELECT sites FROM task_sites WHERE task_id = ?`
)

// UserSitesPersistence is an interface for userSites persistence
type UserSitesPersistence interface {
	InsertUserSites(ctx context.Context, partnerID, userID string, siteIDs []int64) error
	Sites(ctx context.Context, partnerID, userID string) (uSites en.UserSites, err error)

	InsertSitesByTaskID(ctx context.Context, taskID gocql.UUID, siteIDs []string) error
	GetSitesByTaskID(ctx context.Context, taskID gocql.UUID) (siteIDs []string, err error)

	Endpoints(ctx context.Context, partnerID string, siteIDs []string) ([]en.Endpoints, error)
	EndpointsByPartner(ctx context.Context, partnerID string) ([]en.Endpoints, error)
	SaveEndpoints(ctx context.Context, ep []en.Endpoints)
}

// User is userSites persistence for Cassandra
type User struct{}

// UserSitesPersistenceInstance is userSites persistence instance
var UserSitesPersistenceInstance UserSitesPersistence = &User{}

// InsertUserSites inserts UserSites in Cassandra
func (*User) InsertUserSites(ctx context.Context, partnerID, userID string, siteIDs []int64) error {
	return cassandra.QueryCassandra(ctx, insertUserSitesCQL,
		partnerID,
		userID,
		siteIDs,
	).Exec()
}

// Sites returns UserSites found by partnerID and userID
func (*User) Sites(ctx context.Context, partnerID, userID string) (uSites en.UserSites, err error) {
	iter := cassandra.QueryCassandra(ctx, selectUserSitesCQL, partnerID, userID).Iter()

	iter.Scan(
		&uSites.PartnerID,
		&uSites.UserID,
		&uSites.SiteIDs,
	)

	if err = iter.Close(); err != nil {
		return uSites, errors.Errorf("error while working with found entities: %v", err)
	}
	return
}

// InsertSitesByTaskID inserts Sites into the Cassandra
func (*User) InsertSitesByTaskID(ctx context.Context, taskID gocql.UUID, siteIDs []string) error {
	return cassandra.QueryCassandra(ctx, insertTaskSitesCQL, taskID, siteIDs).Exec()
}

// GetSitesByTaskID returns Sites by taskID
func (*User) GetSitesByTaskID(ctx context.Context, taskID gocql.UUID) (siteIDs []string, err error) {
	iter := cassandra.QueryCassandra(ctx, selectSitesByTaskIDCQL, taskID).Iter()

	iter.Scan(
		&siteIDs,
	)

	if err = iter.Close(); err != nil {
		return nil, errors.Errorf("error while working with found entities: %v", err)
	}
	return
}

// Endpoints - select Endpoints by partnerID and sites
func (*User) Endpoints(ctx context.Context, partnerID string, siteIDs []string) ([]en.Endpoints, error) {
	query := `SELECT partner_id, site_id, client_id, endpoints FROM endpoints WHERE partner_id = ? AND site_id IN `
	sites := fmt.Sprintf(`('%s')`, strings.Join(siteIDs, `','`))
	endpoints := make([]en.Endpoints, 0)

	iter := cassandra.QueryCassandra(ctx, query+sites, partnerID).Iter()

	ep := en.Endpoints{}
	for iter.Scan(&ep.PartnerID, &ep.SiteID, &ep.Endpoints, &ep.ClientID) {
		endpoints = append(endpoints, ep)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return endpoints, nil
}

// EndpointsByPartner - select Endpoints by partnerID and sites
func (*User) EndpointsByPartner(ctx context.Context, partnerID string) ([]en.Endpoints, error) {
	query := `SELECT partner_id, site_id, client_id, endpoints FROM endpoints WHERE partner_id = ?`
	endpoints := make([]en.Endpoints, 0)

	iter := cassandra.QueryCassandra(ctx, query, partnerID).Iter()

	ep := en.Endpoints{}
	for iter.Scan(&ep.PartnerID, &ep.SiteID, &ep.ClientID, &ep.Endpoints) {
		endpoints = append(endpoints, ep)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return endpoints, nil
}

// SaveEndpoints - inserts endpoints to DB
func (*User) SaveEndpoints(ctx context.Context, ep []en.Endpoints) {
	query := `INSERT INTO endpoints (partner_id, site_id, client_id, endpoints) VALUES (?, ?, ?, ?)`
	for _, v := range ep {
		err := cassandra.QueryCassandra(ctx, query, v.PartnerID, v.SiteID, v.ClientID, v.Endpoints).Exec()
		if err != nil {
			logger.Log.ErrfCtx(ctx, errorcode.ErrorCantCreateNewTask,"Repo: User.SaveEndpoints: can't save endpoint to DB: %v", err)
		}
	}
}
