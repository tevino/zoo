package enhanced

import "github.com/samuel/go-zookeeper/zk"

// ExistResult contains the result of a ExistWatch.
type ExistResult struct {
	Path  string
	Stat  *zk.Stat
	Err   error
	Exist bool
}
