package registryService

import (
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/serverMonitoring/registry"
	"github.com/gocql/gocql"
)

// RegEntry uses as registry service response struct
type RegEntry struct {
	EndpointID gocql.UUID `json:"endpointID"`
	registry.Data
}
