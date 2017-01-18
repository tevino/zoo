package enhanced

import "github.com/samuel/go-zookeeper/zk"

// ChildrenResult contains the result of a ChildrenWatch.
type ChildrenResult struct {
	Path     string
	Stat     *zk.Stat
	Err      error
	Children []string
}
