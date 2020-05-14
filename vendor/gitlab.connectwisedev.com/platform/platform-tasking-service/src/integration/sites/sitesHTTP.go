package sites

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
)

const siteDelimiter = "%2C"

type sitesResponse struct {
	EndpointID gocql.UUID `json:"endpointID"`
}

// Client represents sites connector client
type Client struct {
	cli integration.HTTPClient
}

// NewClient returns new sites client
func NewClient(cli integration.HTTPClient) *Client {
	return &Client{cli: cli}
}

// GetEndpointsBySiteIDs returns endpoints list by siteIDs
func (c Client) GetEndpointsBySiteIDs(ctx context.Context, partnerID string, siteIDs []string) (ids []gocql.UUID, err error) {
	url := fmt.Sprintf("%s/partner/%s/sites/%s/summary",
		config.Config.AssetMsURL, partnerID, strings.Join(siteIDs, siteDelimiter))

	var resp []sitesResponse
	if err = integration.GetDataByURL(ctx, &resp, c.cli, url, "", true); err != nil {
		switch err.(type) {
		case integration.NotFound: // got NotFound error (404) from Asset MS. There is no Managed endpoints for partner
			return ids, nil
		default:
			return
		}
	}

	for _, endpoint := range resp {
		ids = append(ids, endpoint.EndpointID)
	}
	return
}
