package persistency

//go:generate mockgen -destination=../mocks/mocks-gomock/cache_mock.go  -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency Cache

//Cache interface
type Cache interface {
	Set(key, value []byte, expireSeconds int) (err error)
	Get(key []byte) (value []byte, err error)
}
