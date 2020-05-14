package consumer

import (
	"fmt"
	"runtime/debug"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/utils"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

const returnErrors = true

//New function to return Consumer Service
func New(cfg Config) (Service, error) {
	if len(cfg.Address) == 0 || cfg.Group == "" || len(cfg.Topics) == 0 || cfg.MessageHandler == nil {
		return nil, fmt.Errorf("some of the values are blank. MessageHandler, Address : %v, Group : %s and Topics : %v",
			cfg.Address, cfg.Group, cfg.Topics)
	}

	pool := workerPool{cfg: cfg}
	err := pool.initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to create worker pool : %+v", err)
	}

	config := cluster.NewConfig()

	config.Consumer.Offsets.Retention = cfg.Retention
	config.Consumer.Return.Errors = returnErrors
	config.Group.Return.Notifications = cfg.NotificationHandler != nil
	config.Consumer.Offsets.Initial = cfg.OffsetsInitial
	config.Metadata.Full = cfg.MetadataFull
	config.ClientID = "Continuum"

	config.Version = sarama.MaxVersion
	if cfg.ConsumerMode == PullOrdered || cfg.ConsumerMode == PullOrderedWithOffsetReplay {
		config.Group.Mode = cluster.ConsumerModePartitions
	}

	consumer, err := cluster.NewConsumer(cfg.Address, cfg.Group, cfg.Topics, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer : %+v", err)
	}

	sc := &saramaConsumer{
		cfg:        cfg,
		consumer:   consumer,
		pool:       pool,
		canConsume: true,
	}
	sc.commitStrategy = getCommitStrategy(cfg, sc.MarkOffset)
	sc.consumerStrategy = getConsumerStrategy(cfg)
	return sc, nil
}

type saramaConsumer struct {
	cfg              Config
	consumer         *cluster.Consumer
	pool             workerPool
	commitStrategy   commitStrategy
	consumerStrategy consumerStrategy
	canConsume       bool
}

//Pull batch of messages from Kafka topic
func (sc *saramaConsumer) Pull() {
	go sc.consumeError()
	go sc.consumeRebalanceNotification()

	sc.consumerStrategy.pull(sc)
}

func (sc *saramaConsumer) MarkOffset(topic string, partition int32, offset int64) {
	sc.consumer.MarkPartitionOffset(topic, partition, offset, "")
}

func (sc *saramaConsumer) Close() error {
	err := sc.consumer.Close()
	if err != nil {
		return fmt.Errorf("failed to close consumer : %+v", err)
	}
	sc.pool.Shutdown()
	return nil
}

func (sc *saramaConsumer) consumeError() {
	for err := range sc.consumer.Errors() {
		invokeErrorHandler(err, nil, sc.cfg)
		sc.canConsume = false
	}
}

//Handles notification in case of rebalance
func (sc *saramaConsumer) consumeRebalanceNotification() {
	if sc.cfg.NotificationHandler != nil {
		for notification := range sc.consumer.Notifications() {
			sc.handleConsumerStrategy(notification)
			invokeNotificationHandler(fmt.Sprintf("%+v", notification), sc.cfg)
		}
	}
}

//Executes Consumer mode specific handles for Rebalancing
func (sc *saramaConsumer) handleConsumerStrategy(notification *cluster.Notification) {
	switch sc.cfg.ConsumerMode {
	case PullOrderedWithOffsetReplay:
		cs := sc.consumerStrategy.(*pullOrderedWithOffsetReplay)
		//refresh the pullOrderedWithOffsetReplay partitionkeystore with the
		//new set of partitions which the consumer needs to handle when rebalance is done successfully.
		cs.refreshPartitionKeyStore(notification)
	}

}

func invokeErrorHandler(err error, message *Message, cfg Config) {
	transaction := utils.GetTransactionID()
	if message != nil {
		transaction = message.GetTransactionID()
	}

	defer func() {
		if r := recover(); r != nil {
			Logger().Error(transaction, "invokeErrorHandlerRecovered", "Panic: While processing %v, trace : %s", r, string(debug.Stack()))
		}
	}()
	if cfg.ErrorHandler != nil {
		cfg.ErrorHandler(err, message)
	} else {
		Logger().Debug(transaction, "Consumer failed to consume with : %+v", string(debug.Stack()))
	}
}

func invokeNotificationHandler(notification string, cfg Config) {
	defer func() {
		if r := recover(); r != nil {
			invokeErrorHandler(
				fmt.Errorf("invokeNotificationHandler.panic: While processing %v, trace : %s", r, string(debug.Stack())),
				nil,
				cfg,
			)
		}
	}()
	cfg.NotificationHandler(notification)
}

func (sc *saramaConsumer) Health() (Health, error) {
	health := newHealth(sc.cfg)
	health.CanCosume = sc.canConsume
	conf := sarama.NewConfig()
	client, err := sarama.NewClient(sc.cfg.Address, conf)
	if err != nil {
		health.CanCosume = false
		return *health, err
	}
	defer client.Close() //nolint

	sc.checkConnectionState(client, health)

	health.State = sc.pool.pool.State()
	sc.populatePartition(client, health)
	return *health, err
}

func (sc *saramaConsumer) populatePartition(client sarama.Client, health *Health) {
	for _, topic := range sc.cfg.Topics {
		partitions, err := client.Partitions(topic)
		if err != nil {
			Logger().Info(sc.cfg.TransactionID, "Error while getting partitions for topic: %s, %s", topic, err)
		} else {
			health.PerTopicPartitions[topic] = partitions
		}
	}
}

func (sc *saramaConsumer) checkConnectionState(client sarama.Client, health *Health) {
	coordinator, err := client.Coordinator(sc.cfg.Group)
	if err != nil {
		Logger().Debug(sc.cfg.TransactionID, "Failed to get Coordinator, error: %+v", err)
		health.CanCosume = false
		return
	}
	connected, err := coordinator.Connected()
	if err != nil {
		Logger().Debug(sc.cfg.TransactionID, "Failed to get Coordinator, error: %+v", err)
		health.CanCosume = false
	}

	if connected {
		for _, address := range sc.cfg.Address {
			health.ConnectionState[address] = true
		}
	}
}
