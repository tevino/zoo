package enhanced

import "github.com/samuel/go-zookeeper/zk"

// Conner represents a getter of *zk.Conn.
type Conner interface {
	Conn() *zk.Conn
}
