package enhanced

import "github.com/samuel/go-zookeeper/zk"

// ConnectionStateListener is a handler of zk.Event.
type ConnectionStateListener struct {
	fn func(zk.Event)
}

// Handle calls the function with e.
func (l *ConnectionStateListener) Handle(e zk.Event) {
	l.fn(e)

}

// NewConnectionStateListener creates ConnectionStateListener from fn.
func NewConnectionStateListener(fn func(zk.Event)) *ConnectionStateListener {
	return &ConnectionStateListener{fn}

}
