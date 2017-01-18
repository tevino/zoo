package tree

import "sync"

// ErrorListeners is a set of listeners.
type ErrorListeners struct {
	sync.RWMutex
	listeners []*ErrorListener
}

// NewErrorListeners creates empty ErrorListeners.
func NewErrorListeners() *ErrorListeners {
	var l = &ErrorListeners{}
	l.Clear()
	return l
}

// Add adds a ErrorListener.
func (l *ErrorListeners) Add(listener *ErrorListener) {
	l.Lock()
	l.listeners = append(l.listeners, listener)
	l.Unlock()
}

// Del deletes the first occurrence of listener.
func (l *ErrorListeners) Del(listener *ErrorListener) {
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
func (l *ErrorListeners) Count() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.listeners)
}

// Clear removes all listeners.
func (l *ErrorListeners) Clear() {
	l.Lock()
	l.listeners = make([]*ErrorListener, 0)
	l.Unlock()
}

// Broadcast calls every listener with given error.
func (l *ErrorListeners) Broadcast(err error) {
	l.RLock()
	var wg sync.WaitGroup
	for _, listener := range l.listeners {
		wg.Add(1)
		go func(listener *ErrorListener) {
			listener.Handle(err)
			wg.Done()
		}(listener)
	}
	l.RUnlock()
	wg.Wait()
}
