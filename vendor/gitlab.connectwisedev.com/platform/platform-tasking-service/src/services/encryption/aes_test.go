package encryption

import (
	"testing"

	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
)

func TestNewService(t *testing.T) {
	aes := NewService(EncPrefix)
	encrypted, _ := aes.Encrypt(agentModels.Credentials{
		UseCurrentUser: false,
		Username:       EncPrefix + "I",
		Domain:         EncPrefix + "Love",
		Password:       EncPrefix + "You",
	})

	aes.Decrypt(encrypted)
}
