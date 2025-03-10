package memory

import (
	"container/list"
	"errors"
	"sync"
)

type entry struct {
	key   string
	value []byte
}

type LRUCache struct {
	capacity uint
	cache    map[string]*list.Element
	ll       *list.List
	mu       sync.Mutex
}

func NewLRUCache(capacity uint) (*LRUCache, error) {
	if capacity <= 0 {
		return nil, errors.New("invalid cache capacity")
	}
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		ll:       list.New(),
	}, nil
}

func (l *LRUCache) Push(key string, value []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// If the key already exists, update its value and move it to the front.
	if el, ok := l.cache[key]; ok {
		l.ll.MoveToFront(el)
		el.Value.(*entry).value = value
		return nil
	}

	// Add new entry
	el := l.ll.PushFront(&entry{key, value})
	l.cache[key] = el

	// If the cache exceeds its capacity, remove the least recently used item.
	if uint(l.ll.Len()) > l.capacity {
		l.removeOldest()
	}
	return nil
}

func (l *LRUCache) Pull(key string) ([]byte, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.cache[key]; ok {
		l.ll.MoveToFront(el)
		return el.Value.(*entry).value, nil
	}
	return nil, errors.New("not exist")
}

func (l *LRUCache) removeOldest() {
	oldest := l.ll.Back()
	if oldest != nil {
		l.ll.Remove(oldest)
		kv := oldest.Value.(*entry)
		delete(l.cache, kv.key)
	}
}
