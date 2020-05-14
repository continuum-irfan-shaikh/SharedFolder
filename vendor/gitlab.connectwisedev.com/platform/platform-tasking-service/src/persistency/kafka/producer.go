package kafka

import (
	"context"
	"time"

	"github.com/afex/hystrix-go/hystrix"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
)

// BrokerCommand is a command for Hystrix Broker command
var BrokerCommand = "brokerCommand"

// Publish is implementation of model.BrokerDal.PushEnvelope
var Publish = func(transaction string, producerType publisher.ProducerType, cfg config.Configuration,
	messages ...*publisher.Message) error {
	producer, err := publisher.SyncProducer(producerType, GetConfig(producerType, cfg))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		(time.Duration(cfg.KafkaProducerSettings.TimeoutInSecond*3+1) * time.Second))
	defer cancel()

	return producer.Publish(ctx, transaction, messages...)
}

type wrapper struct {
	transaction  string
	producerType publisher.ProducerType
	cfg          config.Configuration
	msgs         []*publisher.Message
}

// PostMessage post message to kafka
var PostMessage = func(transaction string, producerType publisher.ProducerType, cfg config.Configuration,
	messages ...*publisher.Message) error {

	w := wrapper{
		transaction:  transaction,
		producerType: producerType,
		cfg:          cfg,
		msgs:         messages,
	}

	return Do(BrokerCommand, cfg.CircuitBreakerSettings.Enabled, w.publishKafka, nil)
}

// Do is a function to be called for Circuit breaker execution
var Do = func(command string, circuitEnabled bool, execute func() error, fallback func(error) error) error {
	if circuitEnabled {
		return hystrix.Do(command, execute, fallback)
	}
	return execute()
}

func (w wrapper) publishKafka() error {
	return Publish(w.transaction, w.producerType, w.cfg, w.msgs...)
}

// GetConfig returns a publisher Config for Kafka
func GetConfig(producerType publisher.ProducerType, cfg config.Configuration) *publisher.Config {
	conf := &publisher.Config{
		// Address : is a Kafka Broker IPs
		Address: cfg.Kafka.Brokers,

		// TimeoutInSecond - timeout for producer connectivity
		// defaults to 3 Second
		TimeoutInSecond: int64(cfg.KafkaProducerSettings.TimeoutInSecond),

		// ReconnectIntervalInSecond - retry interval for Kafka reconnect
		// defaults to TimeoutInSecond * 4
		ReconnectIntervalInSecond: int64(cfg.KafkaProducerSettings.ReconnectIntervalInSecond * 2),

		// MaxReconnectRetry - Max number or retry to connect on Kafka
		// defaults to 20
		MaxReconnectRetry: cfg.KafkaProducerSettings.MaxRetry,

		// CleanupTimeInSecond - wait time before starting cleanup
		// defaults to TimeoutInSecond * 5
		CleanupTimeInSecond: int64(cfg.KafkaProducerSettings.TimeoutInSecond * 3),
	}

	if producerType == publisher.BigKafkaProducer {
		conf.MaxMessageBytes = cfg.KafkaProducerSettings.MaxMessageBytes
		conf.CompressionType = cfg.KafkaProducerSettings.CompressionType
	}

	return conf
}
