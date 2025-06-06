package cache

type BytesCache interface {
	Push(key string, value []byte) error
	Pull(key string) ([]byte, error)
}
