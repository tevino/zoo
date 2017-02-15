package enhanced

import "github.com/samuel/go-zookeeper/zk"

// ConnectionStateListeners is a container of ConnectionStateListeners.
type ConnectionStateListeners struct {
	*ListenerContainer
}

// NewConnectionStateListeners creates empty ConnectionStateListeners.
func NewConnectionStateListeners() *ConnectionStateListeners {
	return &ConnectionStateListeners{NewListenerContainer()}
}

// Add adds a Listener.
func (l *ConnectionStateListeners) Add(listener *ConnectionStateListener) {
	l.ListenerContainer.Add(listener)
}

// Del deletes a Listener.
func (l *ConnectionStateListeners) Del(listener *ConnectionStateListener) {
	l.ListenerContainer.Del(listener)
}

// Broadcast calls every listener with given event.
func (l *ConnectionStateListeners) Broadcast(e zk.Event) {
	l.ListenerContainer.Foreach(func(listener interface{}) {
		listener.(*ConnectionStateListener).Handle(e)
	})
}
