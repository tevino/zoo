package enhanced

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// eventWatcher listens for connection state events.
type eventWatcher struct {
	update    <-chan zk.Event
	listeners *ConnectionStateListeners
	closed    chan struct{}
	Conner
}

func newEventWatcher(update <-chan zk.Event, closed chan struct{}, conner Conner) *eventWatcher {
	return &eventWatcher{
		update:    update,
		listeners: NewConnectionStateListeners(),
		closed:    closed,
		Conner:    conner,
	}
}

// Start starts to listener for events.
func (w *eventWatcher) Start() {
	go w.loop()
}

func (w *eventWatcher) loop() {
	select {
	case e := <-w.update:
		w.listeners.Broadcast(e)
	case <-w.closed:
		return
	}
}

// AddListener adds a ConnectionStateListener.
func (w *eventWatcher) AddListener(listener *ConnectionStateListener) {
	w.listeners.Add(listener)
}

// DelListener deletes a ConnectionStateListener.
func (w *eventWatcher) DelListener(listener *ConnectionStateListener) {
	w.listeners.Del(listener)
}

// IsConnected returns true if the client is connected.
func (w *eventWatcher) IsConnected() bool {
	return isConnectedState(w.Conn().State())
}

// BlockUntilConnected blocks until session is created.
// The returning value indicates whether the session is created.
func (w *eventWatcher) BlockUntilConnected(timeout time.Duration) bool {
	var deadline = time.After(timeout)
	var connected = make(chan struct{}, 1)
	var tmpListener = NewConnectionStateListener(func(e zk.Event) {
		if isConnectedState(e.State) {
			select {
			case connected <- struct{}{}:
			default:
			}
		}
	})
	w.AddListener(tmpListener)
	defer w.DelListener(tmpListener)

	if w.IsConnected() {
		return true
	}
	select {
	case <-connected:
		return true
	case <-deadline:
		return false
	}
}
