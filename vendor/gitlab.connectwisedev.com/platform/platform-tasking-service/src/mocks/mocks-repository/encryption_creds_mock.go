package mockrepositories

import (
	"context"

	"github.com/gocql/gocql"

	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
)

type EncryptionServiceMock struct {
}

//Encrypt encrypts credentials
func (*EncryptionServiceMock) Encrypt(creds agentModels.Credentials) (agentModels.Credentials, error) {
	return agentModels.Credentials{}, nil
}

//Decrypt decrypts credentials
func (*EncryptionServiceMock) Decrypt(creds agentModels.Credentials) (agentModels.Credentials, error) {
	return agentModels.Credentials{}, nil
}

type AgentEncryptionServiceMock struct {
}

//Encrypt returns encrypted credentials by public key stored for particular endpoint
func (*AgentEncryptionServiceMock) Encrypt(ctx context.Context, endpointID gocql.UUID, credentials agentModels.Credentials) (agentModels.Credentials, error) {
	return agentModels.Credentials{}, nil
}
