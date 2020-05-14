package consumer

import (
	infrastructureConsumer "gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/consumer"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const defaultGoroutineLimit = 25

// Consumer struct representing consumer settings
type Consumer struct {
	Logger        logger.Logger
	Topic         string
	Goroutines    int
	ConsumerGroup string
	HandleFunc    func(_ infrastructureConsumer.Message) error
	Brokers       []string
}

// NewConsumer returns new consumer struct
func NewConsumer(
	log logger.Logger,
	topic string,
	goroutines int,
	group string,
	brokers []string,
	h func(_ infrastructureConsumer.Message) error,
) Consumer {
	return Consumer{
		Logger:        log,
		Topic:         topic,
		Goroutines:    goroutines,
		ConsumerGroup: group,
		HandleFunc:    h,
		Brokers:       brokers,
	}
}

// Consume start consuming messages
func (c Consumer) Consume() {
	cfg := infrastructureConsumer.NewConfig()

	cfg.Topics = []string{c.Topic}
	cfg.SubscriberPerCore = c.Goroutines
	cfg.MessageHandler = c.HandleFunc
	cfg.Address = c.Brokers
	cfg.Group = c.ConsumerGroup

	service, err := infrastructureConsumer.New(cfg)
	if err != nil {
		c.Logger.ErrfCtx(transactionID.NewContext(), errorcode.ErrorCantProcessData, "%s topics: %s", err.Error(), cfg.Topics)
		return
	}

	service.Pull()
}

// GetGoroutineLimits get goroutine limits if present; if not use default value
func GetGoroutineLimits(cfg config.Configuration, topic string) int {
	if value, ok := cfg.Kafka.GoroutineLimits[topic]; ok {
		return value
	}
	return defaultGoroutineLimit
}
