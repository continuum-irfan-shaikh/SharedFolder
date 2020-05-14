# Messaging Services

Messaging is a wrapper on top of [Sarama](https://github.com/Shopify/sarama) wrapper → A Kafka client library

### Third-Party Libraties

- [Sarama](https://github.com/Shopify/sarama)
- **License** [MIT License](https://github.com/Shopify/sarama/blob/master/LICENSE)
- **Description**
  - Sarama is an MIT-licensed Go client library for Apache Kafka version 0.8 (and later).
- [Sarama Cluster](https://github.com/bsm/sarama-cluster)
- **License** [MIT License](https://github.com/bsm/sarama-cluster/blob/master/LICENSE)
- **Description**
  - Cluster extensions for Sarama, the Go client library for Apache Kafka 0.9
  
- **Glide Dependencies**
    ```yaml
    - package: github.com/Shopify/sarama
    version: 35324cf48e33d8260e1c7c18854465a904ade249
    - package: github.com/bsm/sarama-cluster
    version: d5779253526cc8a3129a0e5d7cc429f4b4473ab4
    - package: github.com/eapache/go-resiliency
    version: ea41b0fad31007accc7f806884dcdf3da98b79ce
    subpackages:
    - breaker
    - package: github.com/golang/snappy
    version: 2e65f85255dbc3072edf28d6b5b8efc472979f5a
    - package: github.com/davecgh/go-spew
    version: d8f796af33cc11cb798c1aaeb27a4ebc5099927d
    subpackages:
    - spew
    - package: github.com/eapache/go-xerial-snappy
    version: 776d5712da21bc4762676d614db1d8a64f4238b0
    - package: github.com/eapache/queue
    version: 093482f3f8ce946c05bcba64badd2c82369e084d
    - package: github.com/klauspost/crc32
    version: 22a7f3e6e2308cfd5c10b0512d2bba0a5a7875b2
    - package: github.com/pierrec/lz4
    version: 623b5a2f4d2a41e411730dcdfbfdaeb5c0c4564e
    - package: github.com/pierrec/xxHash
    version: be086f0f67405de2fac6bc563bf8d0f22fa2a6b2
    subpackages:
    - xxHash32
    - package: github.com/rcrowley/go-metrics
    version: 3113b8401b8a98917cde58f8bbd42a1b1c03b1fd
    ```
## Consumer

The goal of this kafka consumer is to provide throttling at consumer level and to enhance reliability for kafka queue management.

Throttling will ensure the number of GO routines spawned never reaches above pre-defined value (SubscriberPerCore \* Number of core) on consumer running machine.

Reliability is achieved and(or) enhanced by OnMessageCompletion CommitMode (described below) which will ensure that only processed messages will be committed to Kafka cluster.
In case of any panic from consumer end, there will be no message loss and uncommitted messages can be processed on next run of consumer.

This consumer needs to be instantiated once and takes parameters as described at [README](https://gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/blob/master/messaging/consumer/README.md)

Features available with this Consumer :

- Pull - This will call Sarama's Pull method and will fetch messages from Kafka Broker
- MarkOffset - This will commit the offset for a given partition and topic
- Close - This will close consumer connection with Kafka cluster
- Health - this will provide Kafka consumer health status

This consumer also provides notification handler which can get notifications when re-balancing in kafka system happens.

This consumer uses [Worker Pool Library](https://github.com/goinggo/work) to implement throtting for Kafka messages.

Any Application or Microservices are expected to handle below scenario in order to use this library efficiently :

- Message duplicate handling - Since message's offset are committed only when it is processed, this is possible that in case of any panic (application box goes down or if it needs a re-start) , uncommitted messages will be read again and hence duplicate messages may arrive. Application or services using this consumer should handle this duplicate scenario

- Error handling - A consumer implements a timeout of 60 seconds (default) for every message processing. In case of any panic or error (system or user error which can be re-triable or non-retriable) , error handling is expected to be done efficiently by application (or service). If no error handling is done, consumer treats the message as processed after time-out. Corresponding GO routines still exist in the worker pool but they will go under natural death after their life-cycles.

This may also lead to complete blocking worker pool (if no of unprocessed messages == size of worker pool) though they would be marked as completed by this Consumer.
Reaching to that case will block further incoming messages from Kafka brokers.

## Publisher

Publisher is a wrapper on top of [Sarama](https://github.com/Shopify/sarama) → A Kafka client library

- [MIT License](https://github.com/Shopify/sarama/blob/master/LICENSE)

**Additional details on top of Sarama**

This package is providing below additional features on top of Sarama

- Reconnect to Kafka
  - Reconnect to another Kafka instance as soon as any one node goes down or stopped responding
- Publish a message using Context
  - Publish a message using to enable timeout for an messaes instead of waiting for longer to avoid increase in memory consumption

**Limitation**

- Today only support Sync producer, Async need to implementation in case needed for any MS

[README](/messaging/publisher)
