package main

import (
	"github.com/maraino/go-mock"
	"github.com/samuel/go-zookeeper/zk"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/zookeeper"
)

// main contain examples for zookeeper package
func main() {
	zkMockObj, original := zookeeper.InitMock()
	defer zookeeper.Restore(original)

	l := &zookeeper.LockMock{}
	l.When("Lock").Return(nil)
	l.When("Unlock").Return(nil)
	zkMockObj.When("NewLock", mock.Any, mock.Any).Return(l)
	zkMockObj.When("Children", mock.Any).Return([]string{}, &zk.Stat{}, nil)

	// Create new zookeeper lock
	lock := zookeeper.NewLock("some_name_test")
	if err := lock.Lock(); err != nil {
		handleError(err)
	}

	// Unlock with defer
	defer func() {
		err := lock.Unlock()
		if err != nil {
			handleError(err)
		}
	}()

	// Init broadcast
	broadCast, err := zookeeper.InitBroadcast("example", 0)
	if err != nil {
		handleError(err)
	}

	// Add handler
	broadCast.AddHandler("example_handler", func(e *distributed.Event) {})

	// Create Event
	if err := broadCast.CreateEvent(distributed.Event{Type: "example_name", Payload: "example_payload"}); err != nil {
		handleError(err)
	}

}

func handleError(_ error) {
	// actions to handle error
}
