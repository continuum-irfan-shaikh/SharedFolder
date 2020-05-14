package cassandra

import (
	"context"
	"fmt"
	"log"
	"strings"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/db"
	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/db/goc"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/web/rest"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
)

// Session is a component providing access to Cassandra storage
var Session ISession

var isTest bool

// Load creates a Cassandra Client and populates it with initial data
func Load() {
	var err error
	session, err := GetSession()
	if err != nil {
		log.Fatalf("cannot create Cassandra session: %v", err)
	}
	Session = NewSession(session)
}

func GetSession() (goc.Session, error) {
	urls := strings.SplitN(config.Config.CassandraURL, ",", -1)
	if isTest {
		return nil, nil
	}

	timeout := fmt.Sprintf("%ds", config.Config.CassandraTimeoutSec)
	if err := db.Load(urls, config.Config.CassandraKeyspace, timeout, logger.CassandraLogger); err != nil {
		return nil, err
	}
	return db.Session, nil
}

// QueryCassandra creates a session with cassandra and return a point on a query
func QueryCassandra(ctx context.Context, cql string, values ...interface{}) IQuery {
	cassandraQuery := Session.Query(cql, values...)
	logger.Log.DebugfCtx(ctx, "Performed cassandra query: \"%v. with values %v", cql, values)
	return cassandraQuery
}

// ConnectionStatus struct for Cassandra status
type ConnectionStatus struct {
	Session ISession
}

// Status used for getting status of Casandra connection
func (c *ConnectionStatus) Status(conn rest.OutboundConnectionStatus) *rest.OutboundConnectionStatus {
	conn.ConnectionType = "Cassandra"
	conn.ConnectionURLs = strings.SplitN(config.Config.CassandraURL, ",", -1)

	if c.Session == nil {
		conn.ConnectionStatus = rest.ConnectionStatusUnavailable
		return &conn
	}

	if c.Session.Closed() {
		conn.ConnectionStatus = rest.ConnectionStatusUnavailable
		return &conn
	}

	conn.ConnectionStatus = rest.ConnectionStatusActive
	return &conn
}
