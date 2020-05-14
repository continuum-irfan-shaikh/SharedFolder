package tasking

// Cacher is an interface for implementing cache methods
type Cacher interface {
	Set(key, value []byte, expireSeconds int) (err error)
	Get(key []byte) (value []byte, err error)
}
