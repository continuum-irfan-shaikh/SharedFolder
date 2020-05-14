package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

// Sarma configuration options
var (
	brokers = ""
	version = ""
	group   = ""
	topics  = ""
	oldest  = true
	verbose = false
)

func init() {
	flag.StringVar(&brokers, "brokers", "localhost:9092", "Kafka bootstrap brokers to connect to, as a comma separated list") //nolint
	flag.StringVar(&group, "group", "continuum", "Kafka consumer group definition")
	flag.StringVar(&version, "version", "1.1.1", "Kafka cluster version")
	flag.StringVar(&topics, "topics", "test", "Kafka topics to be consumed, as a comma separated list")
	flag.BoolVar(&oldest, "oldest", true, "Kafka consumer consume initial ofset from oldest")
	flag.BoolVar(&verbose, "verbose", false, "Sarama logging")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topics) == 0 {
		panic("no topics given to be consumed, please set the -topics flag")
	}

	if len(group) == 0 {
		panic("no Kafka consumer group defined, please set the -group flag")
	}
}

func main() {
	log.Println("Starting a new Sarama consumer")

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.ClientID = "continuum"
	config.Consumer.Group.Session.Timeout = time.Minute
	//config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	//config.Consumer.Group.Rebalance.Timeout = time.Minute

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{
		ready: make(chan bool),
	}

	// Start with a client
	// clt, err := sarama.NewClient(strings.Split(brokers, ","), config)
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer func() { _ = clt.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), "continuum", config)
	// Start a new consumer group
	// client, err := sarama.NewConsumerGroupFromClient("continuum", clt)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	go func() {
		for {
			if err1 := client.Consume(ctx, strings.Split(topics, ","), &consumer); err1 != nil {
				log.Panicf("Error from consumer: %+v", err)
			}
			// check if context was canceled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context canceled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	log.Printf("Set-up called")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	log.Printf("Clean-up called %+v", s.Claims())
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE: Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

	//for message := range claim.Messages() {
	message := <-claim.Messages()
	log.Printf("%s %v / %v value = %s, topic = %s", session.MemberID(), message.Partition, message.Offset, string(message.Value), message.Topic)
	// for i, h := range message.Headers {
	// 	log.Printf("Header %v ==> Key = %s, Value = %s", i, string(h.Key), string(h.Value))
	// }
	time.Sleep(30 * time.Second)
	session.MarkMessage(message, "")
	// }
	return nil
}
