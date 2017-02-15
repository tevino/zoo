package enhanced

import "sync"

// ListenerContainer implements a general container of listeners.
type ListenerContainer struct {
	sync.RWMutex
	listeners []interface{}
}

// NewListenerContainer creates empty ListenerContainer.
func NewListenerContainer() *ListenerContainer {
	return (&ListenerContainer{}).Clear()
}

// Add adds a Listener.
func (l *ListenerContainer) Add(listener interface{}) {
	l.Lock()
	l.listeners = append(l.listeners, listener)
	l.Unlock()
}

// Del deletes the first occurrence of listener.
func (l *ListenerContainer) Del(listener interface{}) {
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
func (l *ListenerContainer) Count() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.listeners)
}

// Clear removes all listeners.
func (l *ListenerContainer) Clear() *ListenerContainer {
	l.Lock()
	l.listeners = make([]interface{}, 0)
	l.Unlock()
	return l
}

// Foreach calls every listener with given function.
func (l *ListenerContainer) Foreach(fn func(interface{})) {
	l.RLock()
	var wg sync.WaitGroup
	for _, listener := range l.listeners {
		wg.Add(1)
		go func(listener interface{}) {
			fn(listener)
			wg.Done()
		}(listener)
	}
	l.RUnlock()
	wg.Wait()
}
