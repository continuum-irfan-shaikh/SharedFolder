package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/lock"
)

//zkErrorToAction - a map of zk error vs the expected action to be taken by the client
var zkErrorToAction = map[error]lock.ErrorAction{
	zk.ErrConnectionClosed: lock.TryLock,
	zk.ErrSessionExpired:   lock.CreateNewLock,
	zk.ErrDeadlock:         lock.TryUnlock,
	zk.ErrNoNode:           lock.CreateNewLock,
	zk.ErrUnknown:          lock.TryLock,
}

//Lock - Wrapper method to encapsulate error cases for zookeeper lock
func (lw lockWrapper) Lock() error {
	err := lw.zkLock.Lock()
	le := &lock.Error{Code: err}
	switch err {
	case nil:
		return nil
	case zk.ErrConnectionClosed:
		fallthrough
	case zk.ErrSessionExpired:
		fallthrough
	case zk.ErrDeadlock:
		le.Action = zkErrorToAction[le.Code]
	default:
		le.Action = zkErrorToAction[le.Code]
	}

	//if we don't have an action mapped for the error, default to TryLock action
	if le.Action == 0 {
		le.Action = lock.TryLock
	}
	return le
}

//Unlock - Wrapper method to encapsulate error cases for zookeeper unlock
func (lw lockWrapper) Unlock() error {
	err := lw.zkLock.Unlock()
	le := &lock.Error{Code: err}
	switch err {
	case nil:
		return nil
	case zk.ErrConnectionClosed:
		fallthrough
	case zk.ErrNoNode:
		fallthrough
	case zk.ErrSessionExpired:
		le.Action = zkErrorToAction[le.Code]
	default:
		//we don't know what hit us
		le.Action = zkErrorToAction[le.Code]
	}

	//if we don't have an action mapped for the error, default to TryLock action
	if le.Action == 0 {
		le.Action = lock.TryLock
	}
	return le
}
