package dynamicGroups

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	m "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/dg"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
)

const (
	// MessageTypeDynamicGroupStartMonitoring represents msg type to start monitoring
	MessageTypeDynamicGroupStartMonitoring = "START-MONITORING"

	// MessageTypeDynamicGroupStopMonitoring represents msg type to start monitoring
	MessageTypeDynamicGroupStopMonitoring = "STOP-MONITORING"

	// TaskingServiceIDPrefix is an unique id for tasking service that required by DG kafka handler
	TaskingServiceIDPrefix = "tasking-"

	dgDelimiter = "%20OR%20"
)

//go:generate mockgen -destination=../../mocks/mocks-integration/pusher_mock.go -package=common -source=./dynamicGroups.go

type Pusher interface {
	Push(ctx context.Context, msg interface{}) error
}

type dgResponse struct {
	// endpoint id
	ID gocql.UUID `json:"id"`
	// site id of DynamicGroup
	SiteID string `json:"site"`
}

// NewDynamicGroupsClient returns new dynamicGroups client
func NewDynamicGroupsClient(client Pusher, http integration.HTTPClient, user models.UserSitesPersistence) *Client {
	return &Client{
		userRepo:   user,
		client:     client,
		httpClient: http,
	}
}

// Client represents client to communicate with dynamicGroups
type Client struct {
	userRepo   models.UserSitesPersistence
	client     Pusher
	httpClient integration.HTTPClient
}

// GetEndpointsByGroupIDs returns endpoints ids by partner and client IDs
func (c *Client) GetEndpointsByGroupIDs(ctx context.Context, targetIDs []string, createdBy, partnerID string, hasNOCAccess bool) (ids []gocql.UUID, err error) {
	var (
		userSites    entities.UserSites
		payload      []dgResponse
		wg           sync.WaitGroup
		userSitesMap = make(map[string]struct{})
		errChan      = make(chan error, 2)
		done         = make(chan int)
	)

	if !hasNOCAccess {
		wg.Add(1)
		go func() {
			defer wg.Done()

			userSites, err = c.userRepo.Sites(ctx, partnerID, createdBy)
			if err != nil {
				errChan <- fmt.Errorf("error while getting user sites fron Cassandra, err: %v", err)
			}

			for _, siteID := range userSites.SiteIDs {
				siteIDstr := strconv.FormatInt(siteID, 10)
				userSitesMap[siteIDstr] = struct{}{}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("%s/partners/%s/dynamic-groups/managed-endpoints/set?expression=%s",
			config.Config.DynamicGroupsMsURL, partnerID, strings.Join(targetIDs, dgDelimiter))

		if err := integration.GetDataByURL(ctx, &payload, c.httpClient, url, "", true); err != nil {
			errChan <- err
		}
	}()

	go func() {
		for e := range errChan {
			err = fmt.Errorf("%v %v", err, e)
		}
		done <- 1
	}()

	wg.Wait()
	close(errChan)
	<-done
	if err != nil {
		return nil, err
	}

	return c.filterResponse(payload, userSitesMap, hasNOCAccess), nil
}

func (c *Client) filterResponse(payload []dgResponse, userSitesMap map[string]struct{}, hasNOCAccess bool) (ids []gocql.UUID) {
	for _, dgEndpoint := range payload {
		if _, ok := userSitesMap[dgEndpoint.SiteID]; ok || hasNOCAccess {
			ids = append(ids, dgEndpoint.ID)
		}
	}
	return ids
}

// StartMonitoringGroups sends message to DG to start monitor group by groupID
func (c *Client) StartMonitoringGroups(ctx context.Context, partnerID string, groupIDs []string, taskID gocql.UUID) error {
	for _, id := range groupIDs {
		payload := m.MonitoringDG{
			PartnerID:      partnerID,
			DynamicGroupID: id,
			ServiceID:      TaskingServiceIDPrefix + taskID.String(),
			Operation:      MessageTypeDynamicGroupStartMonitoring,
		}

		if err := c.client.Push(ctx, payload); err != nil {
			return err
		}
	}
	return nil
}

// StopGroupsMonitoring removes group from monitoring
func (c *Client) StopGroupsMonitoring(ctx context.Context, partnerID string, groupIDs []string, taskID gocql.UUID) error {
	for _, id := range groupIDs {
		payload := m.MonitoringDG{
			PartnerID:      partnerID,
			DynamicGroupID: id,
			ServiceID:      TaskingServiceIDPrefix + taskID.String(),
			Operation:      MessageTypeDynamicGroupStopMonitoring,
		}

		if err := c.client.Push(ctx, payload); err != nil {
			return err
		}
	}
	return nil
}
