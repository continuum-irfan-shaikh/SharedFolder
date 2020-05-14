package kafka

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/messaging"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/kafka"
)

const (
	// MsgTypeTaskCreated ..
	MsgTypeTaskCreated = `TASK_CREATED`
	// MsgTypeTaskUpdated ..
	MsgTypeTaskUpdated = `TASK_UPDATED`
	// MsgTypeTaskInstanceUpdated ..
	MsgTypeTaskInstanceUpdated = `TASK_INSTANCE_UPDATED`
	// MsgTypeExecutionResultUpdated ..
	MsgTypeExecutionResultUpdated = `EXECUTION_RESULT_UPDATED`
)

// Tasking represent wrapper over kafka producer for pushing tasking events
type Tasking struct {
	config config.Configuration
}

// NewTasking initializes struct for Tasking
func NewTasking(cfg config.Configuration) *Tasking {
	return &Tasking{config: cfg}
}

// Push sends KafkaMessage  with msgType into TaskingEventsTopic of kafka
func (t *Tasking) Push(message interface{}, msgType string) error {
	msgObj, err := publisher.EncodeObject(&message)
	if err != nil {
		return fmt.Errorf("publisher is not able to encode message: %v. err: %s", message, err.Error())
	}

	msg := &publisher.Message{
		Topic: t.config.TaskingEventsTopic,
		Value: msgObj,
	}

	msg.AddHeader(messaging.MessageType, msgType)

	err = kafka.PostMessage("", publisher.RegularKafkaProducer, t.config, msg)

	if err != nil {
		return fmt.Errorf("sending massage to kafka topic %v filed. err: %s", t.config.TaskingEventsTopic, err.Error())
	}

	return nil
}
