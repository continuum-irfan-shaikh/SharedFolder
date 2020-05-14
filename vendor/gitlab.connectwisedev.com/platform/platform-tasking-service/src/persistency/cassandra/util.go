package cassandra

import (
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/db/goc"
)

//go:generate mockgen -destination=../../mocks/mocks-cassandra/cassandra_mock.go  -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra ISession,IQuery,IIter,IBatch

// NewSession returns new cassandra session client
func NewSession(s goc.Session) ISession {
	return &session{s: s}
}

type (
	// ISession - database session interface
	ISession interface {
		Query(string, ...interface{}) IQuery
		NewBatch(gocql.BatchType) IBatch
		ExecuteBatch(IBatch) error
		Closed() bool
	}

	// IQuery gocql query interface
	IQuery interface {
		SetConsistency(c gocql.Consistency)
		Exec() error
		Iter() IIter
		Scan(dest ...interface{}) error
		Release()
		PageState(state []byte) IQuery
		PageSize(n int) IQuery
	}

	// IIter interface
	IIter interface {
		Scan(dest ...interface{}) bool
		Close() error
		PageState() []byte
	}

	// IBatch batch interface
	IBatch interface {
		Query(stmt string, args ...interface{})
		Size() int
	}
)

type session struct {
	s goc.Session
}

// Query implements query method
func (c *session) Query(stmt string, values ...interface{}) IQuery {
	return &query{q: c.s.Query(stmt, values...)}
}

// Closed implements Closed session method to check whether session is closed or not
func (c *session) Closed() bool {
	return c.s.Closed()
}

// NewBatch implements NewBatch method
func (c *session) NewBatch(typ gocql.BatchType) IBatch {
	return &batch{b: c.s.NewBatch(typ)}
}

// ExecuteBatch implements NewBatch method
func (c *session) ExecuteBatch(bch IBatch) error {
	b, ok := bch.(*batch)
	if !ok {
		return errors.New("assertion failed")
	}

	return c.s.ExecuteBatch(b.b)
}

type query struct {
	q *gocql.Query
}

// SetConsistency implements SetConsistency method
func (q *query) SetConsistency(c gocql.Consistency) {
	q.q.Consistency(c)
}

// Exec implements Exec method
func (q *query) Exec() error {
	return q.q.Exec()
}

// Iter implements Iter method
func (q *query) Iter() IIter {
	return q.q.Iter()
}

// Scan implements Iter method
func (q *query) Scan(dest ...interface{}) error {
	return q.q.Scan(dest...)
}

// Scan implements Iter method
func (q *query) String() string {
	return fmt.Sprintf("%v", q.q)
}

// Release implements Release method
func (q *query) Release() {
	q.q.Release()
}

// PageSize implements PageSize method
func (q *query) PageSize(n int) IQuery {
	q.q.PageSize(n)
	return q
}

// PageState implements PageState method
func (q *query) PageState(state []byte) IQuery {
	q.q.PageState(state)
	return q
}

type batch struct {
	b *gocql.Batch
}

// Query implements Query method
func (b *batch) Query(stmt string, args ...interface{}) {
	b.b.Query(stmt, args...)
}

func (b *batch) Size() int {
	return b.b.Size()
}
