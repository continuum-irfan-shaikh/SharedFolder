package zookeeper

import (
	"testing"

	"github.com/samuel/go-zookeeper/zk"
)

func TestLock(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "ErrConnectionClosed",
			expectedErr: zk.ErrConnectionClosed,
		},
		{
			name:        "ErrSessionExpired",
			expectedErr: zk.ErrSessionExpired,
		},
		{
			name:        "ErrDeadlock",
			expectedErr: zk.ErrDeadlock,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lock := &LockMock{}
			lock.When("Lock").Return(test.expectedErr)
			lWrapper := lockWrapper{
				zkLock: lock,
			}
			if err := lWrapper.Lock(); err.Error() != test.expectedErr.Error() {
				t.Fatalf("expected err: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}

func TestUnlock(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "ErrConnectionClosed",
			expectedErr: zk.ErrConnectionClosed,
		},
		{
			name:        "ErrNoNode",
			expectedErr: zk.ErrNoNode,
		},
		{
			name:        "ErrSessionExpired",
			expectedErr: zk.ErrSessionExpired,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lock := &LockMock{}
			lock.When("Unlock").Return(test.expectedErr)
			lWrapper := lockWrapper{
				zkLock: lock,
			}
			if err := lWrapper.Unlock(); err.Error() != test.expectedErr.Error() {
				t.Fatalf("expected err: %s, got: %s", test.expectedErr, err)
			}
		})
	}
}
