package tree

import (
	"errors"
	"path"
	"sort"
	"sync"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/tevino/zoo/enhanced"
)

// Node represents a node within a tree of ZNodes.
type Node struct {
	sync.RWMutex
	cache     *Cache
	state     NodeState
	parent    *Node
	path      string
	childData *ChildData
	children  map[string]*Node
	depth     int
}

// NewNode creates a Node with given path and parent.
// NOTE: Parent should be nil if the node is root.
func NewNode(cache *Cache, path string, parent *Node) *Node {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}
	return &Node{
		cache:    cache,
		state:    NodeStatePending,
		parent:   parent,
		path:     path,
		children: make(map[string]*Node),
		depth:    depth,
	}
}

// SwapChildData sets ChildData to given value and returns the old ChildData.
func (n *Node) SwapChildData(d *ChildData) *ChildData {
	n.Lock()
	defer n.Unlock()
	old := n.childData
	n.childData = d
	return old
}

// Children returns the children of current node.
func (n *Node) Children() map[string]*Node {
	n.RLock()
	defer n.RUnlock()
	children := make(map[string]*Node, len(n.children))
	for k, v := range n.children {
		children[k] = v
	}
	return children
}

// FindChild finds a child of current node by its relative path.
// NOTE: path should contain no slash.
func (n *Node) FindChild(path string) (*Node, bool) {
	n.RLock()
	defer n.RUnlock()
	node, ok := n.children[path]
	return node, ok
}

// ChildData returns the ChildData.
func (n *Node) ChildData() *ChildData {
	n.RLock()
	defer n.RUnlock()
	return n.childData
}

// RemoveChild removes child by path.
func (n *Node) RemoveChild(path string) {
	n.Lock()
	defer n.Unlock()
	delete(n.children, path)
}

func (n *Node) refresh() {
	if (n.depth < n.cache.maxDepth) && n.cache.selector.TraverseChildren(n.path) {
		n.cache.incOutstandingOpsBy(2)
		n.doRefreshData()
		n.doRefreshChildren()
	} else {
		n.refreshData()
	}
}

func (n *Node) refreshChildren() {
	if (n.depth < n.cache.maxDepth) && n.cache.selector.TraverseChildren(n.path) {
		n.cache.incOutstandingOpsBy(1)
		n.doRefreshChildren()
	}
}

func (n *Node) refreshData() {
	n.cache.incOutstandingOpsBy(1)
	n.doRefreshData()
}

func (n *Node) doRefreshChildren() {
	n.cache.client.WatchChildren(n.path, n.processChildrenResult, n.processWatchEvent)
	// n.cache.client.GetChildren().UsingWatcher(
	// 	curator.NewWatcher(n.processWatchEvent),
	// ).InBackgroundWithCallback(n.processResult).ForPath(n.path)
}

func (n *Node) doRefreshData() {
	n.cache.client.WatchData(n.path, n.processDataResult, n.processWatchEvent)
	// n.cache.client.GetData().UsingWatcher(
	// 	curator.NewWatcher(n.processWatchEvent),
	// ).InBackgroundWithCallback(n.processResult).ForPath(n.path)
}

func (n *Node) wasReconnected() error {
	n.refresh()
	for _, child := range n.Children() {
		if err := child.wasReconnected(); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) wasCreated() {
	n.refresh()
}

func (n *Node) wasDeleted() {
	oldChildData := n.SwapChildData(nil)
	for _, child := range n.Children() {
		child.wasDeleted()
	}

	if n.cache.state.EqualTo(CacheStateStopped) {
		return
	}

	oldState := n.state.SwapValue(NodeStateDead)
	if oldState == NodeStateLive {
		n.cache.publishEvent(CacheEventNodeRemoved, oldChildData)
	}

	if n.parent == nil {
		// Root node; use an exist query to watch for existence.
		n.cache.client.WatchExist(n.path, n.processExistResult, n.processWatchEvent)
		// n.cache.client.CheckExists().UsingWatcher(
		// 	curator.NewWatcher(n.processWatchEvent),
		// ).InBackgroundWithCallback(n.processResult).ForPath(n.path)
	} else {
		// Remove from parent if we're currently a child
		n.parent.RemoveChild(path.Base(n.path))
	}
}

// processWatchEvent processes watch events.
func (n *Node) processWatchEvent(evt zk.Event) {
	// n.cache.logger.Debugf("ProcessWatchEvent: %v", evt)
	switch evt.Type {
	case zk.EventNodeCreated:
		if n.parent != nil {
			n.cache.handleException(errors.New("unexpected NodeCreated on non-root node"))
			return
		}
		n.wasCreated()
	case zk.EventNodeChildrenChanged:
		n.refreshChildren()
	case zk.EventNodeDataChanged:
		n.refreshData()
	case zk.EventNodeDeleted:
		n.wasDeleted()
	default:
		// Leave other type of events unhandled
		// n.cache.logger.Printf("Event received: %v", evt)
	}
}

func (n *Node) processChildrenResult(result enhanced.ChildrenResult) {
	switch result.Err {
	case zk.ErrNoNode:
		n.wasDeleted()
	case nil:
		n.Lock()
		oldChildData := n.childData
		if oldChildData != nil && oldChildData.Stat().Mzxid == result.Stat.Mzxid {
			// Only update stat if mzxid is same, otherwise we might obscure
			// GET_DATA event updates.
			n.childData.SetStat(result.Stat)
		}
		n.Unlock()

		if len(result.Children) == 0 {
			break
		}
		// Present new children in sorted order for test determinism.
		children := sort.StringSlice(result.Children)
		sort.Sort(children)
		for _, child := range children {
			if accepted := n.cache.selector.AcceptChildData(path.Join(n.path, child)); !accepted {
				continue
			}
			n.Lock()
			if _, exists := n.children[child]; !exists {
				fullPath := path.Join(n.path, child)
				node := NewNode(n.cache, fullPath, n)
				n.children[child] = node
				node.wasCreated()
			}
			n.Unlock()
		}
	}

	n.cache.completeOutstandingOps()
}

func (n *Node) processExistResult(result enhanced.ExistResult) {
	if n.parent != nil {
		n.cache.handleException(errors.New("unexpected EXISTS on non-root node"))
	}
	if result.Err == nil {
		n.state.SetValueIf(NodeStateDead, NodeStatePending)
		n.wasCreated()
	}

	n.cache.completeOutstandingOps()
}

func (n *Node) processDataResult(result enhanced.DataResult) {
	// n.cache.logger.Debugf("ProcessResult: %v", evt)
	var newStat = result.Stat

	// case curator.CHILDREN:
	switch result.Err {
	case zk.ErrNoNode:
		n.wasDeleted()
	case nil:
		newChildData := NewChildData(result.Path, newStat, result.Data)
		oldChildData := n.ChildData()
		if n.cache.cacheData {
			n.SwapChildData(newChildData)
		} else {
			n.SwapChildData(NewChildData(result.Path, newStat, nil))
		}

		var added bool
		if n.parent == nil {
			// We're the singleton root.
			added = n.state.SwapValue(NodeStateLive) != NodeStateLive
		} else {
			added = n.state.SetValueIf(NodeStatePending, NodeStateLive)
			if !added {
				// Ordinary nodes are not allowed to transition from dead -> live;
				// make sure this isn't a delayed response that came in after death.
				if !n.state.EqualTo(NodeStateLive) {
					return
				}
			}
		}

		if added {
			n.cache.publishEvent(CacheEventNodeAdded, newChildData)
		} else {
			if oldChildData == nil || oldChildData.Stat().Mzxid != newStat.Mzxid {
				n.cache.publishEvent(CacheEventNodeUpdated, newChildData)
			}
		}
	default:
		// n.cache.logger.Printf("Unknown GET_DATA event[%v]: %s", evt.Path(), evt.Err())
	}

	n.cache.completeOutstandingOps()
}
