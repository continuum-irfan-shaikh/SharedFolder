package kafka_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/kafka"
)

var mockKafkaMessage = "test_kafka_message"

func TestPostMessage(t *testing.T) {
	t.Run("reach_failed_threshold_and_reconnect_after_first_success", func(t *testing.T) {
		var (
			requestVolumeThreshold = 5
			expectedPublishError   = errors.New("test")
			cfg                    = config.Configuration{}
		)

		oldPublish := kafka.Publish
		oldReconnect := publisher.ReConnect
		defer func() {
			kafka.Publish = oldPublish
			publisher.ReConnect = oldReconnect
		}()

		cfg.CircuitBreakerSettings.Enabled = false

		// Mock publish with failed result
		var publishCount int32
		kafka.Publish = func(transaction string, producerType publisher.ProducerType, cfg config.Configuration,
			messages ...*publisher.Message) error {
			atomic.AddInt32(&publishCount, 1)
			return expectedPublishError
		}

		// Mock reconnect
		var reconnectTries int32
		publisher.ReConnect = func(transaction string, cfg *publisher.Config) {
			atomic.AddInt32(&reconnectTries, 1)
		}

		// Send requestVolumeThreshold x 2 messages with failed result
		messageObj, _ := publisher.EncodeObject(&mockKafkaMessage)
		for i := 0; i < requestVolumeThreshold*2; i++ {
			err := kafka.PostMessage(gocql.TimeUUID().String(), publisher.RegularKafkaProducer,
				cfg,
				&publisher.Message{
					Topic: "patching_calculate_status",
					Value: messageObj,
				})
			assert.NotNil(t, err)
			assert.EqualError(t, err, expectedPublishError.Error())
		}

		// Wait some time
		time.Sleep(time.Millisecond * 200)

		// Mock publish with success result
		kafka.Publish = func(transaction string, producerType publisher.ProducerType, cfg config.Configuration,
			messages ...*publisher.Message) error {
			return nil
		}

		// Post should be without errors and should be Closed state and mark as connected
		err := kafka.PostMessage(gocql.TimeUUID().String(), publisher.RegularKafkaProducer,
			cfg,
			&publisher.Message{
				Topic: "patching_calculate_status",
				Value: messageObj,
			})
		assert.Nil(t, err)
	})

	t.Run("all_messages_success_without_reconnect", func(t *testing.T) {
		var (
			requestVolumeThreshold   = 5
			expectedReconnectRetries int32
			cfg                      = config.Configuration{}
		)

		oldPublish := kafka.Publish
		oldReconnect := publisher.ReConnect
		defer func() {
			kafka.Publish = oldPublish
			publisher.ReConnect = oldReconnect
		}()

		cfg.CircuitBreakerSettings.Enabled = false

		// Mock publish with failed result
		var publishCount int32
		kafka.Publish = func(transaction string, producerType publisher.ProducerType, cfg config.Configuration,
			messages ...*publisher.Message) error {
			atomic.AddInt32(&publishCount, 1)
			return nil
		}

		// Mock reconnect
		var reconnectTries int32
		publisher.ReConnect = func(transaction string, cfg *publisher.Config) {
			atomic.AddInt32(&reconnectTries, 1)
		}

		// Send requestVolumeThreshold x 2 messages with failed result
		messageObj, _ := publisher.EncodeObject(&mockKafkaMessage)
		for i := 0; i < requestVolumeThreshold*2; i++ {
			err := kafka.PostMessage(gocql.TimeUUID().String(), publisher.RegularKafkaProducer,
				cfg,
				&publisher.Message{
					Topic: "patching_calculate_status",
					Value: messageObj,
				})
			assert.Nil(t, err)
		}

		assert.Equal(t, int32(requestVolumeThreshold*2), publishCount)
		assert.Equal(t, expectedReconnectRetries, reconnectTries)
	})

	t.Run("Original_publish_failed_disabled_circuit", func(t *testing.T) {
		oldReconnect := publisher.ReConnect
		defer func() {
			publisher.ReConnect = oldReconnect
		}()

		cfg := config.Configuration{}
		cfg.Kafka.Brokers = []string{"wrong.dns"}
		cfg.CircuitBreakerSettings.Enabled = false

		messageObj, _ := publisher.EncodeObject(&mockKafkaMessage)
		err := kafka.PostMessage(gocql.TimeUUID().String(), publisher.RegularKafkaProducer,
			cfg,
			&publisher.Message{
				Topic: "patching_calculate_status",
				Value: messageObj,
			})
		assert.NotNil(t, err)
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("regular_kafka", func(t *testing.T) {
		cfg := config.Configuration{}
		cfg.Kafka.Brokers = []string{"wrong.dns"}

		pcfg := kafka.GetConfig(publisher.RegularKafkaProducer, cfg)
		assert.Equal(t, cfg.Kafka.Brokers, pcfg.Address)
	})

	t.Run("big_kafka", func(t *testing.T) {
		cfg := config.Configuration{}
		cfg.Kafka.Brokers = []string{"wrong.dns"}
		cfg.KafkaProducerSettings.MaxMessageBytes = 10

		pcfg := kafka.GetConfig(publisher.BigKafkaProducer, cfg)
		assert.Equal(t, cfg.Kafka.Brokers, pcfg.Address)
		assert.Equal(t, cfg.KafkaProducerSettings.MaxMessageBytes, pcfg.MaxMessageBytes)
	})
}
