package enhanced

import "github.com/samuel/go-zookeeper/zk"

// DataResult contains the result of a DataWatch.
type DataResult struct {
	Path string
	Stat *zk.Stat
	Err  error
	Data []byte
}
