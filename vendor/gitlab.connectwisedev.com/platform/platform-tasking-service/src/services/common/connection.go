package common

import (
	"fmt"
	"strings"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"github.com/Shopify/sarama"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	dbMessageType = "OutboundConnectionStatus"
	dbNameSuffix  = "OutboundConnectionStatus"

	//ConnectionStatusActive correct connection
	ConnectionStatusActive = "Active"
	//ConnectionStatusUnavailable fail connection
	ConnectionStatusUnavailable = "Unavailable"
)

var (
	// connMethods stores map of all supported connections and their corresponding methods
	connMethods = map[string]func(status OutboundConnectionStatus) OutboundConnectionStatus{
		"Kafka":     getKafkaStatus,
		"Cassandra": getCassandraStatus,
	}
)

// OutboundConnectionStatus represents status of all connections into DG service
type OutboundConnectionStatus struct {
	TimeStampUTC     time.Time `json:"timeStampUTC"`
	Type             string    `json:"type"`
	Name             string    `json:"name"`
	ConnectionType   string    `json:"connectionType"`
	ConnectionURLs   []string  `json:"connectionURLs"`
	ConnectionStatus string    `json:"connectionStatus"`
}

// GetOutboundConnectionStatus used for getting status of all connections in DG service
func GetOutboundConnectionStatus() (connections []OutboundConnectionStatus) {
	baseConn := OutboundConnectionStatus{
		TimeStampUTC: time.Now(),
		Type:         dbMessageType,
		Name:         fmt.Sprintf("%s-%s", config.Config.Version.ServiceName, dbNameSuffix),
	}

	for _, conn := range config.Connections {
		connections = append(connections, connMethods[conn](baseConn))
	}

	return
}

// getCassandraStatus used for getting status of Casandra connection
func getCassandraStatus(conn OutboundConnectionStatus) OutboundConnectionStatus {
	conn.ConnectionType = "Cassandra"
	conn.ConnectionURLs = strings.SplitN(config.Config.CassandraURL, ",", -1)
	conn.ConnectionStatus = ConnectionStatusActive

	return conn
}

// getKafkaStatus used for getting status of Kafka connection
func getKafkaStatus(conn OutboundConnectionStatus) OutboundConnectionStatus {
	conn.ConnectionType = "Kafka"
	conn.ConnectionURLs = strings.SplitN(config.Config.KafkaBrokers, ",", -1)
	ctx := transactionID.NewContext()
	consumer, err := sarama.NewConsumer(strings.Split(config.Config.KafkaBrokers, ","), nil)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "Kafka connection error: %v.", err)
		conn.ConnectionStatus = ConnectionStatusUnavailable
		return conn
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "getKafkaStatus: error while closing consumer: %v", err)
		}
	}()

	conn.ConnectionStatus = ConnectionStatusActive

	return conn
}
