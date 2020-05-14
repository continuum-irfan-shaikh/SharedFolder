package goc

import (
	"testing"
	"time"

	"github.com/gocql/gocql"
)

const (
	version  = 4
	keyspace = "gockle_test"
)

func TestNewSession(t *testing.T) {
	if a, e := NewSession(nil), (session{}); a != e {
		t.Errorf("Actual session %v, expected %v", a, e)
	}

	var c = gocql.NewCluster("localhost")

	c.ProtoVersion = version

	var s, err = c.CreateSession()

	if err != nil {
		t.Skip(err)
	}

	if a, e := NewSession(s), (session{s: s}); a != e {
		t.Errorf("Actual session %v, expected %v", a, e)
	}
}

func TestNewSimpleSession(t *testing.T) {
	hosts := []string{"localhost"}
	timeout := 3 * time.Second
	if s, err := NewSimpleSession(keyspace, hosts, timeout); err == nil {
		t.Error("Actual no error, expected error")
	} else if s != nil {
		t.Errorf("Actual session %v, expected nil", s)
		s.Close()
	}

	a, err := NewSimpleSession(keyspace, hosts, timeout)
	switch {
	case err != nil:
		t.Skip(err)
	case a == nil:
		t.Errorf("Actual session nil, expected not nil")
	default:
		a.Close()
	}
}
