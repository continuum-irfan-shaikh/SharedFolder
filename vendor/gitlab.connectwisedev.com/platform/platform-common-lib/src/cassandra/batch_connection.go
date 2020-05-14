package cassandra

import (
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
	exc "gitlab.connectwisedev.com/platform/platform-common-lib/src/exception"

	"github.com/gocql/gocql"
)

//batchConnection is responsible of connecting with Cassandra db and creating batch
type batchConnection struct {
	Connection connection
	batch      *gocql.Batch
}

//NewBatchDbConnection is a factory method which returns the struct implementation of BatchDbConnector
func NewBatchDbConnection(conf *DbConfig) (BatchDbConnector, error) {
	return newBatchConnection(conf)
}

// newBatchDbConnection is a constructor of batchDbConnection which will initialize struct and
// will return an open connection object(if no error) of batchDbConnection
func newBatchConnection(conf *DbConfig) (*batchConnection, error) {
	batchdb := &batchConnection{}
	db, err := newConnection(conf)
	if err != nil {
		return nil, err
	}
	batchdb.Connection = *db
	batchdb.batch = db.session.NewBatch(gocql.LoggedBatch)
	return batchdb, err

}

func (d batchConnection) BatchExecution(query string, values [][]interface{}) (err error) {
	length := len(values)
	for i := 0; i < length; i++ {
		d.batch.Query(query, values[i]...)
	}
	return d.executeBatch(d.batch)
}

func (d batchConnection) executeBatch(b *gocql.Batch) error {
	if d.Connection.session == nil {
		return exc.New(ErrDbNoOpenConnection, nil)
	}

	err := circuit.Do(d.Connection.conf.CommandName, d.Connection.conf.CircuitBreaker.Enabled, func() error {
		err := d.Connection.session.ExecuteBatch(b)
		return err
	}, nil)

	if err != nil {
		return err
	}
	return nil
}

//Close function closes the connection and does not return error
func (d batchConnection) Close() {
	if d.Connection.session != nil {
		d.Connection.session.Close()
	}
}

//Closed function to check is session is closed or not
func (d batchConnection) Closed() bool {
	if d.Connection.session != nil {
		return d.Connection.session.Closed()
	}
	return true
}
