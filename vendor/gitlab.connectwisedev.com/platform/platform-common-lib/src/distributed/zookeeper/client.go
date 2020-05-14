package zookeeper

import (
	"strings"

	"github.com/samuel/go-zookeeper/zk"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/lock"
)

const zkSeparator = "/"

type (
	// ZKClient describe interface for zookeeper client
	ZKClient interface {
		// State returns the current state of the connection
		State() string
		// Exists checks is item exist
		Exists(path string) (bool, *zk.Stat, error)
		// Get gets data by path
		Get(path string) ([]byte, *zk.Stat, error)
		// Children gets list of item children
		Children(path string) ([]string, *zk.Stat, error)
		// Set sets item data to path
		Set(path string, data []byte, version int32) (*zk.Stat, error)
		// Delete deletes item from zookeeper by its path
		Delete(path string, version int32) error
		// NewLock creates new zookeeper lock
		NewLock(path string, acl []zk.ACL) lock.Locker
		// CreateRecursive creates new zookeeper item recursive
		CreateRecursive(childPath string, data []byte, flag int32, acl []zk.ACL) (string, error)
		// Close closes zookeeper client connection
		Close()
		// Events returns chan of client events
		Events() <-chan zk.Event
	}

	zkClient struct {
		conn   *zk.Conn
		events <-chan zk.Event
	}
)

func (client *zkClient) State() string {
	return client.conn.State().String()
}

func (client *zkClient) Exists(path string) (bool, *zk.Stat, error) {
	return client.conn.Exists(path)
}

func (client *zkClient) Get(path string) ([]byte, *zk.Stat, error) {
	return client.conn.Get(path)
}

func (client *zkClient) Children(path string) ([]string, *zk.Stat, error) {
	return client.conn.Children(path)
}

func (client *zkClient) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	return client.conn.Set(path, data, version)
}

func (client *zkClient) Delete(path string, version int32) error {
	return client.conn.Delete(path, version)
}

func (client *zkClient) NewLock(path string, acl []zk.ACL) lock.Locker {
	return zk.NewLock(client.conn, path, acl)
}

func (client *zkClient) CreateRecursive(childPath string, data []byte, flag int32, acl []zk.ACL) (path string, err error) {
	path, err = client.conn.Create(childPath, data, flag, acl)
	if err != zk.ErrNoNode {
		return path, err
	}

	// Create parent node.
	parts := strings.Split(childPath, zkSeparator)
	// always skip first argument it should be empty string
	for i := range parts[1:] {
		nPath := strings.Join(parts[:i+2], zkSeparator)

		var exists bool
		exists, _, err = client.conn.Exists(nPath)
		if err != nil {
			return path, err
		}

		if exists {
			continue
		}

		// the last one set real data and flag
		if len(parts)-2 == i {
			return client.conn.Create(nPath, data, flag, acl)
		}

		path, err = client.conn.Create(nPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return path, err
		}
	}

	return path, err
}

func (client *zkClient) Close() {
	client.conn.Close()
}

func (client *zkClient) Events() <-chan zk.Event {
	return client.events
}
