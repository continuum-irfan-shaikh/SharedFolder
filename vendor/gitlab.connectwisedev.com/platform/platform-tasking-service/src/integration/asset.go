package integration

import (
	"context"
	"time"

	"github.com/gocql/gocql"
)

// Asset is an interface to communicate with asset MS
type Asset interface {
	GetLocationByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (location *time.Location, err error)
	GetResourceTypeByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (ResourceType, error)
	GetSiteIDByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (siteID, clientID string, err error)
	GetMachineNameByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (string, error)
}

//go:generate mockgen -destination=../mocks/mocks-gomock/assetRepo_mock.go -package=mocks -source=./asset.go

// ResourceType represents endpoint resource type
type ResourceType string

const (
	// Desktop is a Desktop type
	Desktop ResourceType = "Desktop"
	// Server is a Server type
	Server ResourceType = "Server"
)

// IsAllResources returns if resource is selected to all or not
func (r ResourceType) IsAllResources() bool {
	return string(r) == ""
}
