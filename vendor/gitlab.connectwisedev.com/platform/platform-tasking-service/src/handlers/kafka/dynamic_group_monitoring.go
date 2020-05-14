package kafka

import (
	"context"
	"fmt"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/kafka"
)

// DynamicGroups struct to communicate with DynamicGroups monitoring topic
type DynamicGroups struct {
	config config.Configuration
}

// NewDynamicGroups initializes struct for DynamicGroupMonitoring
func NewDynamicGroups(cfg config.Configuration) *DynamicGroups {
	return &DynamicGroups{config: cfg}
}

// Push sends KafkaMessage  with msgType into DynamicGroupMonitoringTopic of kafka
func (d *DynamicGroups) Push(ctx context.Context, msg interface{}) error {
	msgObj, err := publisher.EncodeObject(&msg)
	if err != nil {
		return fmt.Errorf("publisher is not able to encode message: %v. err: %v", msg, err)
	}

	err = kafka.PostMessage(transactionID.FromContext(ctx), publisher.RegularKafkaProducer, d.config, &publisher.Message{
		Topic: d.config.DynamicGroupMonitoringTopic,
		Value: msgObj,
	})

	if err != nil {
		return fmt.Errorf("sending massage to kafka topic %v filed. err: %s", d.config.DynamicGroupMonitoringTopic, err.Error())
	}

	return nil
}
