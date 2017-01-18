package tree

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync/atomic"

	"github.com/samuel/go-zookeeper/zk"

	"github.com/tevino/abool"
	"github.com/tevino/zoo/enhanced"
)

// Cache is a a utility that attempts to keep all data from all children of a ZK path locally cached.
// It will watch the ZK path, respond to update/create/delete events, pull down the data, etc.
// You can register a listener that will get notified when changes occur.
//
// NOTE: It's not possible to stay transactionally in sync. Users must be prepared for false-positives
// and false-negatives. Additionally, always use the version number when updating data to avoid overwriting
// another process' change.
type Cache struct {
	// Tracks the number of outstanding background requests in flight. The first time this count reaches 0, we publish the initialized event.
	outstandingOps uint64
	isInitialized  *abool.AtomicBool
	root           *Node
	client         *enhanced.Client
	cacheData      bool
	maxDepth       int
	selector       Selector
	eventListeners *CacheEventListeners
	errorListeners *ErrorListeners
	state          CacheState
	// connectionStateListener curator.ConnectionStateListener
	// logger                  Logger
	createParent bool
}

// NewCache creates a Cache for the given client and path with default options.
//
// If the client is namespaced, all operations on the resulting TreeCache will be in terms of
// the namespace, including all published events.  The given path is the root at which the
// TreeCache will watch and explore.  If no node exists at the given path, the TreeCache will
// be initially empty.
func NewCache(client *enhanced.Client, root string, selector Selector) *Cache {
	if selector == nil {
		selector = DefaultSelector
	}
	var cache = &Cache{
		isInitialized:  abool.New(),
		client:         client,
		maxDepth:       math.MaxInt32,
		cacheData:      true,
		selector:       selector,
		state:          CacheStateLatent,
		eventListeners: NewCacheEventListeners(),
		errorListeners: NewErrorListeners(),
		// logger: &DummyLogger{},
	}
	cache.root = NewNode(cache, root, nil)
	// cache.connectionStateListener = curator.NewConnectionStateListener(
	// 	func(client curator.CuratorFramework, newState curator.ConnectionState) {
	// 		c.handleStateChange(newState)
	// 	})
	return cache
}

// Start starts the TreeCache.
// The cache is not started automatically. You must call this method.
func (c *Cache) Start() error {
	if !c.state.SetValueIf(CacheStateLatent, CacheStateStarted) {
		return errors.New("already started/stopped")
	}
	if c.createParent {
		if err := c.createParentNodes(); err != nil {
			return err
		}
	}
	if !c.client.IsConnected() {
		// TODO: allow disconnected?
		return errors.New("client not connected")
	}
	c.root.wasCreated()
	return nil
}

func (c *Cache) createParentNodes() error {
	var err = c.client.CreateWithParents(c.root.path)
	if err == zk.ErrNodeExists {
		err = nil
	} else if err != nil {
		err = fmt.Errorf("failed to create parents: %s", err)
	}
	return err
}

// SetCacheData sets whether or not to cache byte data per node, default true.
// NOTE: When this set to false, the events still contain data of znode
// but you can't query them by TreeCache.CurrentData/CurrentChildren
func (c *Cache) SetCacheData(cacheData bool) *Cache {
	c.cacheData = cacheData
	return c
}

// SetMaxDepth sets the maximum depth to explore/watch.
// Set to 0 will watch only the root node.
// Set to 1 will watch the root node and its immediate children.
// Default to math.MaxInt32.
func (c *Cache) SetMaxDepth(depth int) *Cache {
	c.maxDepth = depth
	return c
}

// SetCreateParentNodes sets whether to auto-create parent nodes for the cached path.
// By default, TreeCache does not do this.
// Note: Parent nodes is only created when Start() is called.
func (c *Cache) SetCreateParentNodes(yes bool) *Cache {
	c.createParent = yes
	return c
}

// // SetLogger sets the inner Logger of TreeCache.
// func (c *Cache) SetLogger(l Logger) *Cache {
// 	c.logger = l
// 	return c
// }

// Stop stops the cache.
func (c *Cache) Stop() {
	if c.state.SetValueIf(CacheStateStarted, CacheStateStopped) {
		// c.client.ConnectionStateListenable().RemoveListener(c.connectionStateListener)
		// c.listeners.Clear()
		c.root.wasDeleted()
	}
}

// AddEventListener adds a CacheEventListener.
func (c *Cache) AddEventListener(l *CacheEventListener) {
	c.eventListeners.Add(l)
}

// DelEventListener deletes the first occurrence of l.
func (c *Cache) DelEventListener(l *CacheEventListener) {
	c.eventListeners.Del(l)
}

// AddErrorListener adds an error listener.
func (c *Cache) AddErrorListener(l *ErrorListener) {
	c.errorListeners.Add(l)
}

// DelErrorListener deletes the first occurrence of l.
func (c *Cache) DelErrorListener(l *ErrorListener) {
	c.errorListeners.Del(l)
}

// // Listeners returns the cache listeners.
// func (c *Cache) Listeners() EventListeners {
// 	return &c.listeners
// }

// // UnhandledErrorListeners allows catching unhandled errors in asynchornous operations.
// func (c *Cache) UnhandledErrorListenable() ErrorListeners {
// 	return &c.errorListeners
// }

// findNode finds the node which matches the given path.
// ErrRootNotMatch is returned if the given path doesn't share a same root with
// the TreeCache.
// ErrNodeNotFound is returned if the given path can not be found.
func (c *Cache) findNode(path string) (*Node, error) {
	if !strings.HasPrefix(path, c.root.path) {
		return nil, ErrRootNotMatch
	}

	path = strings.TrimPrefix(path, c.root.path)
	current := c.root
	for _, part := range strings.Split(path, "/") {
		if part == "" {
			continue
		}
		next, exists := current.FindChild(part)
		if !exists {
			return nil, ErrNodeNotFound
		}
		current = next
	}

	return current, nil
}

// CurrentChildren returns the current set of children at the given full path, mapped by child name.
// There are no guarantees of accuracy; this is merely the most recent view of the data.
// If there is no node at this path, ErrNodeNotFound is returned.
func (c *Cache) CurrentChildren(fullPath string) (map[string]*ChildData, error) {
	node, err := c.findNode(fullPath)
	if err != nil {
		return nil, err
	}
	if !node.state.EqualTo(NodeStateLive) {
		return nil, ErrNodeNotLive
	}

	children := node.Children()
	m := make(map[string]*ChildData, len(children))
	for child, childNode := range children {
		// Double-check liveness after retreiving data.
		childData := childNode.ChildData()
		if childData != nil && childNode.state.EqualTo(NodeStateLive) {
			m[child] = childData
		}
	}

	// Double-check liveness after retreiving children.
	if !node.state.EqualTo(NodeStateLive) {
		return nil, ErrNodeNotLive
	}
	return m, nil
}

// CurrentData returns the current data for the given full path.
// There are no guarantees of accuracy. This is merely the most recent view of the data.
// If there is no node at the given path, ErrNodeNotFound is returned.
func (c *Cache) CurrentData(fullPath string) (*ChildData, error) {
	node, err := c.findNode(fullPath)
	if err != nil {
		return nil, err
	}
	if !node.state.EqualTo(NodeStateLive) {
		return nil, ErrNodeNotLive
	}

	return node.ChildData(), nil
}

// handleException sends an exception to all listeners, or else log the error if there are none.
func (c *Cache) handleException(e error) {
	if c.errorListeners.Count() == 0 {
		// c.logger.Printf("%s", e)
		return
	}
	c.errorListeners.Broadcast(e)
}

// func (c *Cache) handleStateChange(newState curator.ConnectionState) {
// 	switch newState {
// 	case curator.SUSPENDED:
// 		c.publishEvent(CacheEventConnSuspended, nil)
// 	case curator.LOST:
// 		c.isInitialized.UnSet()
// 		c.publishEvent(CacheEventConnLost, nil)
// 	case curator.CONNECTED:
// 		c.root.wasCreated()
// 	case curator.RECONNECTED:
// 		if err := c.root.wasReconnected(); err == nil {
// 			c.publishEvent(CacheEventConnReconnected, nil)
// 		}
// 	}
// }

// publishEvent publish an event with given type and data to all listeners.
func (c *Cache) publishEvent(tp CacheEventType, data *ChildData) {
	if !c.state.EqualTo(CacheStateStopped) {
		var evt = CacheEvent{Type: tp, Data: data}
		// c.logger.Debugf("publishEvent: %v", evt)
		go c.eventListeners.Broadcast(evt)
	}
}

// // callListeners calls all listeners with given event.
// // Error is handled by handleException().
// func (c *Cache) callListeners(evt TreeCacheEvent) {
// 	c.listeners.ForEach(func(listener interface{}) {
// 		if err := listener.(TreeCacheListener).ChildEvent(c.client, evt); err != nil {
// 			c.handleException(err)
// 		}
// 	})
// }

// incOutstandingOpsBy increases oustandingOps by given value.
func (c *Cache) incOutstandingOpsBy(n int) {
	atomic.AddUint64(&c.outstandingOps, uint64(n))
}

func (c *Cache) completeOutstandingOps() {
	// Decrease by 1
	if atomic.AddUint64(&c.outstandingOps, ^uint64(0)) == 0 {
		if !c.isInitialized.IsSet() {
			c.isInitialized.Set()
			c.publishEvent(CacheEventInitialized, nil)
		}
	}
}
