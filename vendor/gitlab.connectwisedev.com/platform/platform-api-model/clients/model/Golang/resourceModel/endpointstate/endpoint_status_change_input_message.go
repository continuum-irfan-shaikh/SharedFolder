package endpointstate

import (
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
)

// InputMessage structure that represents Kafka message
type InputMessage struct {
	agent.BrokerEnvelope
	Message NotificationMessage `json:"Message"`
}
