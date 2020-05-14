package integration

import (
	"context"

	"github.com/gocql/gocql"
)

// SitesConnector  represents interface to communicate with sites API
type SitesConnector interface {
	GetEndpointsBySiteIDs(ctx context.Context,partnerID string, siteIDs []string) ([]gocql.UUID, error)
}
