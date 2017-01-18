package tree

import "errors"

var (
	// ErrNodeNotFound indicates a node can not be found.
	ErrNodeNotFound = errors.New("node not found")
	// ErrRootNotMatch indicates the root path does not match.
	ErrRootNotMatch = errors.New("root path not match")
	// ErrNodeNotLive indicates the state of node is not LIVE.
	ErrNodeNotLive = errors.New("node state is not LIVE")
)
