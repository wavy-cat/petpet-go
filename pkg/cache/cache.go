package cache

import "errors"

type BytesCache interface {
	// Push adds an item to the cache.
	Push(key string, value []byte) error
	// Pull returns an element from the cache. If the element is not found, it returns cache.NotExists.
	Pull(key string) ([]byte, error)
	Close() error
}

var NotExists = errors.New("element does not exist")
