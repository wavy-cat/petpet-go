package cache

import "errors"

type BytesCache interface {
	// Push adds an item to the cache.
	Push(key string, value []byte) error
	// Pull returns an element from the cache. If the element is not found, it returns cache.ErrNotExists.
	Pull(key string) ([]byte, error)
	Close() error
}

var ErrNotExists = errors.New("element does not exist")
