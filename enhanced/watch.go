package enhanced

import "github.com/samuel/go-zookeeper/zk"

// watchOperations contains APIs of watching ZNode.
type watchOperations struct {
	Conner
	*namespace
	closed chan struct{}
}

func newWatchOperations(ns *namespace, closed chan struct{}, conner Conner) watchOperations {
	return watchOperations{
		namespace: ns,
		Conner:    conner,
		closed:    closed,
	}
}

// WatchChildren watches children of given path.
// processResult will be called after watching.
// processEvent will be called with further event ONCE.
func (w *watchOperations) WatchChildren(p string, processResult func(ChildrenResult), processEvent func(zk.Event)) {
	go func() {
		children, stat, change, err := w.Conn().ChildrenW(w.namespaced(p))
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
func (w *watchOperations) WatchData(p string, processResult func(DataResult), processEvent func(zk.Event)) {
	go func() {
		data, stat, change, err := w.Conn().GetW(w.namespaced(p))
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
func (w *watchOperations) WatchExist(p string, processResult func(ExistResult), processEvent func(zk.Event)) {
	go func() {
		exist, stat, change, err := w.Conn().ExistsW(w.namespaced(p))
		processResult(ExistResult{
			Exist: exist,
			Path:  p,
			Stat:  stat,
			Err:   err})
		w.waitForEvent(change, processEvent)
	}()
}

func (w *watchOperations) waitForEvent(ch <-chan zk.Event, processEvent func(zk.Event)) {
	select {
	case evt := <-ch:
		processEvent(evt)
	case <-w.closed:
		// TODO: send event to eventCallback?
	}
}
