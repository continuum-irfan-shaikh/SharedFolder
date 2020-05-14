package cassandra

import (
	"github.com/gocql/gocql"
	"log"
	"strings"
	"time"
)

// Session is a component providing access to Cassandra storage
var Session *gocql.Session

// Load creates a Cassandra Client and populates it with initial data
func Load(cassandraURL, cassandraKeyspace string, cassandraTimeoutSec, cassandraConnNumber int) {

	urls := strings.SplitN(cassandraURL, ",", -1)
	client := gocql.NewCluster(urls...)
	client.ProtoVersion = 4
	client.Consistency = gocql.Quorum
	client.Keyspace = cassandraKeyspace
	client.NumConns = cassandraConnNumber
	client.Timeout = time.Duration(cassandraTimeoutSec) * time.Second

	var err error
	Session, err = client.CreateSession()
	if err != nil {
		log.Fatalf("cannot create Cassandra session: %v", err)
	}
}
