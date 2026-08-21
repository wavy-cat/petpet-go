package cache

// Noop is a disabled cache implementation.
// Reads always miss, while writes and closing are no-ops.
type Noop struct{}

func NewNoop() BytesCache {
	return Noop{}
}

func (Noop) Push(string, []byte) error {
	return nil
}

func (Noop) Pull(string) ([]byte, error) {
	return nil, ErrNotExists
}

func (Noop) Close() error {
	return nil
}
