package enhanced

import "github.com/samuel/go-zookeeper/zk"

// ConnGetter represents a getter of *zk.Conn
type ConnGetter interface {
	Conn() *zk.Conn
}

// Watch contains APIs of watching ZNode.
type Watch struct {
	ConnGetter
	closed chan struct{}
}

// WatchChildren watches children of given path.
// processResult will be called after watching.
// processEvent will be called with further event ONCE.
func (w *Watch) WatchChildren(p string, processResult func(ChildrenResult), processEvent func(zk.Event)) {
	go func() {
		children, stat, change, err := w.Conn().ChildrenW(p)
		processResult(ChildrenResult{
			Path:     p,
			Children: children,
			Stat:     stat,
			Err:      err})
		w.waitForEvent(change, processEvent)
	}()
}

// WatchData watches data of given path.
// processResult will be called after watching.
// processEvent will be called with further event ONCE.
func (w *Watch) WatchData(p string, processResult func(DataResult), processEvent func(zk.Event)) {
	go func() {
		data, stat, change, err := w.Conn().GetW(p)
		processResult(DataResult{
			Path: p,
			Data: data,
			Stat: stat,
			Err:  err})
		w.waitForEvent(change, processEvent)
	}()
}

// WatchExist watches existence of given path.
// processResult will be called after watching.
// processEvent will be called with further event ONCE.
func (w *Watch) WatchExist(p string, processResult func(ExistResult), processEvent func(zk.Event)) {
	go func() {
		exist, stat, change, err := w.Conn().ExistsW(p)
		processResult(ExistResult{
			Exist: exist,
			Path:  p,
			Stat:  stat,
			Err:   err})
		w.waitForEvent(change, processEvent)
	}()
}

func (w *Watch) waitForEvent(ch <-chan zk.Event, processEvent func(zk.Event)) {
	select {
	case evt := <-ch:
		processEvent(evt)
	case <-w.closed:
		// TODO: send event to eventCallback?
	}
}
