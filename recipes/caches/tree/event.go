package tree

import "fmt"

// CacheEventType represents the type of change to a path.
type CacheEventType int

const (
	// CacheEventNodeAdded indicates a node was added.
	CacheEventNodeAdded CacheEventType = iota
	// CacheEventNodeUpdated indicates a node's data was changed.
	CacheEventNodeUpdated
	// CacheEventNodeRemoved indicates a node was removed from the tree.
	CacheEventNodeRemoved

	// CacheEventConnSuspended is called when the connection has changed to SUSPENDED.
	CacheEventConnSuspended
	// CacheEventConnReconnected is called when the connection has changed to RECONNECTED.
	CacheEventConnReconnected
	// CacheEventConnLost is called when the connection has changed to LOST.
	CacheEventConnLost
	// CacheEventInitialized is posted after the initial cache has been fully populated.
	CacheEventInitialized
)

// String returns the string representation of CacheEventType.
// "Unknown" is returned when event type is unknown
func (et CacheEventType) String() string {
	switch et {
	case CacheEventNodeAdded:
		return "NodeAdded"
	case CacheEventNodeUpdated:
		return "NodeUpdated"
	case CacheEventNodeRemoved:
		return "NodeRemoved"
	case CacheEventConnSuspended:
		return "ConnSuspended"
	case CacheEventConnReconnected:
		return "ConnReconnected"
	case CacheEventConnLost:
		return "ConnLost"
	case CacheEventInitialized:
		return "Initialized"
	default:
		return "Unknown"
	}
}

// CacheEvent represents a change to a path
type CacheEvent struct {
	Type CacheEventType
	Data *ChildData
}

// String returns the string representation of CacheEvent
func (e CacheEvent) String() string {
	var path string
	var data []byte
	if e.Data != nil {
		path = e.Data.Path()
		data = e.Data.Data()
	}
	return fmt.Sprintf("CacheEvent{%s %s '%s'}", e.Type, path, data)
}
