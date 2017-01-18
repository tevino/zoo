package tree

// CacheEventListener is a handler of cache event.
type CacheEventListener struct {
	fn func(CacheEvent)
}

// Handle calls the function with e.
func (l *CacheEventListener) Handle(e CacheEvent) {
	l.fn(e)
}

// NewCacheEventListener creates CacheEventListener with fn.
func NewCacheEventListener(fn func(CacheEvent)) *CacheEventListener {
	return &CacheEventListener{fn}
}
