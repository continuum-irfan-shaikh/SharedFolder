package integration

import (
	"context"

	"github.com/gocql/gocql"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
)

// AgentEncryptionService - provide interface for encryption credentials by key stored for particular endpoint
type AgentEncryptionService interface {
	Encrypt(ctx context.Context, endpointID gocql.UUID, credentials agentModels.Credentials) (encrypted agentModels.Credentials, err error)
}
