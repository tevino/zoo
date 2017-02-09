package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tevino/zoo/enhanced"
)

// ZkEnv should only be used in TESTING.
// It starts a managed ZooKeeper cluster before testing and stops it after.
type ZkEnv struct {
	t         *testing.T
	assert    *assert.Assertions
	zkCluster *ZkCluster
	client    *enhanced.Client
	ZNodeAssertion
	ClusterOperation
}

// NewZkEnv creates a ZkEnv.
func NewZkEnv(t *testing.T) *ZkEnv {
	var env = &ZkEnv{t: t, assert: assert.New(t)}
	env.ZNodeAssertion = ZNodeAssertion{t: t, ClientGetter: env}
	env.ClusterOperation = ClusterOperation{t: t, ClusterGetter: env}
	return env
}

// With calls fn after ".Start()" and "defer .Stop()"
func (z *ZkEnv) With(fn func(env *ZkEnv)) {
	z.Start()
	defer z.Stop()
	fn(z)
}

// Start starts ZooKeeper cluster.
func (z *ZkEnv) Start() *ZkEnv {
	var err error
	z.zkCluster, err = StartZkCluster(1, nil, nil)
	z.assert.NoError(err)
	return z
}

// Stop stops ZooKeeper cluster.
func (z *ZkEnv) Stop() {
	var err = z.zkCluster.Stop()
	z.assert.NoError(err)
	if z.client != nil {
		z.client.Close()
	}
}

// Client returns a client to ZooKeeper cluster.
// NOTE: Don't close the client, it will be closed by .Stop()
func (z *ZkEnv) Client() *enhanced.Client {
	if z.client == nil {
		z.client = z.NewClient()
	}
	return z.client
}

// NewClient creates a client to ZooKeeper cluster.
func (z *ZkEnv) NewClient() *enhanced.Client {
	var client, err = z.zkCluster.ConnectAll()
	z.assert.NoError(err)
	return client
}

// NewClientTimeout creates a client to ZooKeeper cluster and wait until it gets
// connected within given timeout.
func (z *ZkEnv) NewClientTimeout(timeout time.Duration) *enhanced.Client {
	var client = z.NewClient()
	var connected = client.BlockUntilConnected(timeout)
	z.assert.True(connected)
	return client
}

// ConnectionString returns the connection string of ZooKeeper cluster.
func (z *ZkEnv) ConnectionString() string {
	return z.zkCluster.ConnectionString()
}

// Zk returns ZooKeeper cluster.
func (z *ZkEnv) Zk() *ZkCluster {
	return z.zkCluster
}
