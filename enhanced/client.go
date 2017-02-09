package enhanced

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// Client is an enhanced connection.
type Client struct {
	conn        *zk.Conn
	eventUpdate <-chan zk.Event
	closed      chan struct{}
	namespace
	nsBasicOperations
	watchOperations
}

// Connect establishes a new connection to a pool of zookeeper
// servers. The provided session timeout sets the amount of time for which
// a session is considered valid after losing connection to a server. Within
// the session timeout it's possible to reestablish a connection to a different
// server and keep the same session. This is means any ephemeral nodes and
// watches are maintained.
func Connect(servers []string, sessionTimeout time.Duration) (*Client, error) {
	// TODO:
	// 1. connect timeout
	conn, evt, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return nil, err
	}
	return newClient(conn, evt), nil
}

func newClient(conn *zk.Conn, eventUpdate <-chan zk.Event) *Client {
	var c = &Client{
		conn:        conn,
		eventUpdate: eventUpdate,
		closed:      make(chan struct{}),
	}
	c.nsBasicOperations = newNSBasicOperations(c, &c.namespace)
	c.watchOperations = watchOperations{
		namespace: &c.namespace,
		closed:    c.closed,
		Conner:    c,
	}
	return c
}

// Conn returns internal zk.Conn.
func (c *Client) Conn() *zk.Conn {
	return c.conn
}

// Close closes inner connection then stops all watchers.
func (c *Client) Close() {
	c.conn.Close()
	close(c.closed)
}

// SetDigestAuth sets auth as digest.
func (c *Client) SetDigestAuth(auth []byte) error {
	return c.conn.AddAuth("digest", auth)
}

// Namespace returns namespace used for all operation.
func (c *Client) Namespace() string {
	return c.ns()
}

// SetNamespace sets the namespace used for all operation.
// NOTE: namespace does not start with /.
func (c *Client) SetNamespace(ns string) *Client {
	c.setNS(ns)
	return c
}

// IsConnected returns true if a session is currently created.
func (c *Client) IsConnected() bool {
	return isConnectedState(c.conn.State())
}

func isConnectedState(s zk.State) bool {
	switch s {
	case zk.StateHasSession, zk.StateConnectedReadOnly:
		return true
	default:
		return false
	}
}

// BlockUntilConnected blocks until session is created.
// The returning value indicates whether the session is created.
func (c *Client) BlockUntilConnected(timeout time.Duration) bool {
	var deadline = time.After(timeout)
	for {
		if c.IsConnected() {
			return true
		}
		select {
		case evt := <-c.eventUpdate:
			if isConnectedState(evt.State) {
				return true
			}
		case <-deadline:
			return false
		}
	}
}
