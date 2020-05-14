package entities

import (
	"time"

	"github.com/gocql/gocql"
)

// ExecutionExpiration describes execution expiration data for particular task instance ID
type ExecutionExpiration struct {
	ExpirationTimeUTC  time.Time
	PartnerID          string
	TaskInstanceID     gocql.UUID
	ManagedEndpointIDs []gocql.UUID
}
