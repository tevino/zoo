package tree

import "sync"

// CacheEventListeners is a set of CacheEventListeners.
type CacheEventListeners struct {
	sync.RWMutex
	listeners []*CacheEventListener
}

// NewCacheEventListeners creates empty CacheEventListeners.
func NewCacheEventListeners() *CacheEventListeners {
	var ls = &CacheEventListeners{}
	ls.Clear()
	return ls
}

// Add adds a CacheEventListener.
func (l *CacheEventListeners) Add(listener *CacheEventListener) {
	l.Lock()
	l.listeners = append(l.listeners, listener)
	l.Unlock()
}

// Del deletes the first occurrence of listener.
func (l *CacheEventListeners) Del(listener *CacheEventListener) {
	l.Lock()
	defer l.Unlock()
	for i, lis := range l.listeners {
		if lis == listener {
			l.listeners = append(l.listeners[:i], l.listeners[i+1:]...)
			return
		}
	}
}

// Count returns the number of listeners.
func (l *CacheEventListeners) Count() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.listeners)
}

// Clear deletes all listeners.
func (l *CacheEventListeners) Clear() {
	l.Lock()
	l.listeners = make([]*CacheEventListener, 0)
	l.Unlock()
}

// Broadcast calls every listener with given CacheEvent.
func (l *CacheEventListeners) Broadcast(e CacheEvent) {
	var wg sync.WaitGroup
	l.RLock()
	for _, listener := range l.listeners {
		wg.Add(1)
		go func(listener *CacheEventListener) {
			listener.Handle(e)
			wg.Done()
		}(listener)
	}
	l.RUnlock()
	wg.Wait()
}
