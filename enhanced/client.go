package enhanced

import (
	"path"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// Client is an enhanced connection.
type Client struct {
	conn        *zk.Conn
	flags       int32
	acl         []zk.ACL
	eventUpdate <-chan zk.Event
	closed      chan struct{}
	Watch
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
	// 2. Namespace
	conn, evt, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return nil, err
	}
	var c = &Client{
		conn:        conn,
		eventUpdate: evt,
		flags:       0,
		acl:         zk.WorldACL(zk.PermAll),
		closed:      make(chan struct{}),
	}
	c.Watch = Watch{closed: c.closed, ConnGetter: c}
	return c, nil
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

// SetFlags sets the flags used for all operation.
func (c *Client) SetFlags(flags int32) {
	c.flags = flags
}

// SetACL sets the ACL used for all operation.
func (c *Client) SetACL(acl []zk.ACL) {
	c.acl = acl
}

// Get fetches value and stat of given znode.
func (c *Client) Get(p string) ([]byte, *zk.Stat, error) {
	return c.conn.Get(p)
}

// Exist returns true and stat of given znode if it exists.
func (c *Client) Exist(p string) (bool, *zk.Stat, error) {
	return c.conn.Exists(p)
}

// GetChildren fetches children of given path.
func (c *Client) GetChildren(p string) ([]string, *zk.Stat, error) {
	return c.conn.Children(p)
}

// Set sets the value on given znode.
func (c *Client) Set(p string, value []byte, version int32) (*zk.Stat, error) {
	return c.conn.Set(p, value, version)
}

// Create creates given znode with value set to nil.
func (c *Client) Create(p string) error {
	_, err := c.conn.Create(p, nil, c.flags, c.acl)
	return err
}

// CreateValue creates given znode with value.
func (c *Client) CreateValue(p string, value []byte) error {
	_, err := c.conn.Create(p, value, c.flags, c.acl)
	return err
}

// Delete deletes given znode.
func (c *Client) Delete(p string, version int32) error {
	return c.conn.Delete(p, version)
}

// DeleteWithChildren deletes given znode with its children if any.
func (c *Client) DeleteWithChildren(p string) error {
	var children, _, err = c.GetChildren(p)
	if err != nil {
		return err
	}
	for _, child := range children {
		err = c.DeleteWithChildren(path.Join(p, child))
		if err != nil {
			return err
		}
	}
	return c.Delete(p, -1)
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

// CreateWithParents create path with its parents created if missing.
func (c *Client) CreateWithParents(p string) error {
	var err = c.Create(p)
	if err == zk.ErrNoNode {
		var parent = path.Dir(p)
		err = c.CreateWithParents(parent)
		if err == nil {
			err = c.Create(p)
		}
	}
	return err
}

// CreateValueWithParents create path with value and its parents created if missing.
func (c *Client) CreateValueWithParents(p string, value []byte) error {
	var err = c.CreateValue(p, value)
	if err == zk.ErrNoNode {
		var parent = path.Dir(p)
		err = c.CreateWithParents(parent)
		if err == nil {
			err = c.CreateValue(p, value)
		}
	}
	return err
}
