<p align="center">
<img height=70px src="docs/images/continuum-logo.png">
<img height=70px src="docs/images/Go-Logo_Blue.png">
</p>

# Zookeeper Client Library

This is a Zookeeper client implementation. 
We can use it to broadcast events and work with scheduler and job listener.
### Third-Party Libraries

- [GoLea](https://github.com/Comcast/go-leaderelection) - 
**License** [Apache License 2.0](https://github.com/Comcast/go-leaderelection/blob/master/LICENSE) - 
**Description** - GoLea provides the capability for a set of distributed processes to compete for leadership for a shared resource. It is implemented using Zookeeper for the underlying support. 
- [Go-ZooKeeper](https://github.com/samuel/go-zookeeper) - 
**License** [3-clause BSD](https://github.com/samuel/go-zookeeper/blob/master/LICENSE) - 
**Description** - Native Go Zookeeper Client Library
- [go-mock](https://github.com/maraino/go-mock) - 
**License** [MIT](https://github.com/maraino/go-mock/blob/master/LICENSE) - 
**Description** - A mocking framework for Go.
- [cron](https://github.com/robfig/cron) - 
**License** [MIT](https://github.com/robfig/cron/blob/master/LICENSE) - 
**Description** - Package cron implements a cron spec parser and job runner.

### [Example](https://gitlab.connectwisedev.com/platform/platform-common-lib/tree/master/src/distributed/example/example.go)

**Import Statement**

```go
import	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed"
```

**Zookeeper Session**

```go
// Create zookeeper session
err := zookeeper.Init("localhost:2181", "/openapi-service", logger.ZKLogger)
```

**Zookeeper Client**

```go
// Get the current state of the connection
state := zookeeper.Client.State()

// Check is item exist in zookeeper
isExist := zookeeper.Client.Exists("some/path")

// Get item data from zookeeper by its path
data, _, err := zookeeper.Client.Get("some/path")

// Get list of item children
children, _, err := zookeeper.Client.Children("some/path")

// Set item data to path
_, err := zookeeper.Client.Set("some/path", []byte(data), version)

// Delete item from zookeeper
err := zookeeper.Client.Delete("some/path", version)

// Lock zookeeper item
zkLock := zookeeper.Client.NewLock("some/path", acl)

// Create new zookeeper item recursive
path, err := zookeeper.Client.CreateRecursive("some/path",[]byte(data),flag, acl)

```