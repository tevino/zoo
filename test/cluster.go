package test

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/tevino/zoo/enhanced"
)

// ZkCluster is a managed ZooKeeper cluster for testing purpose.
type ZkCluster struct {
	*zk.TestCluster
	action *NodeAction
}

// Connect starts a client to a single server at given index.
func (c *ZkCluster) Connect(idx int) (*enhanced.Client, error) {
	var servers = []string{strings.Split(c.ConnectionString(), ",")[idx]}
	return enhanced.Connect(servers, time.Second)
}

// ConnectAll starts a client to all servers.
func (c *ZkCluster) ConnectAll() (*enhanced.Client, error) {
	var servers = strings.Split(c.ConnectionString(), ",")
	return enhanced.Connect(servers, time.Second)
}

// ConnectionString returns connection string like: 127.0.0.1:21810,127.0.0.1:21811
func (c *ZkCluster) ConnectionString() string {
	if len(c.Servers) == 0 {
		panic("No Server")
	}
	var buf []byte
	for i, s := range c.Servers {
		buf = append(buf, []byte("127.0.0.1:")...)
		buf = strconv.AppendInt(buf, (int64)(s.Port), 10)
		if i != len(c.Servers)-1 {
			buf = append(buf, byte(','))
		}
	}
	return string(buf)
}

// DoCreate creates znodes specified by given YAML.
func (c *ZkCluster) DoCreate(yml []byte) error {
	return c.applyAction(yml, c.action.DoCreate)
}

// DoDelete deletes znodes specified by given YAML.
func (c *ZkCluster) DoDelete(yml []byte) error {
	return c.applyAction(yml, c.action.DoDelete)
}

// DoUpdate is DoCreate with ErrNodeExists ignored.
func (c *ZkCluster) DoUpdate(yml []byte) error {
	return c.applyAction(yml, c.action.DoUpdate)
}

func (c *ZkCluster) startAction() error {
	var err error
	c.action, err = StartNewNodeAction(c)
	return err
}

// Stop stops the testing cluster.
func (c *ZkCluster) Stop() error {
	var err error
	if c.action != nil {
		c.action.Stop()
	}
	if c.TestCluster != nil {
		err = c.TestCluster.Stop()
	}
	return err
}

func (c *ZkCluster) applyAction(yml []byte, act NodeActionFunc) error {
	root, err := UnmarshalYAML(yml)
	if err != nil {
		return err
	}

	for key, node := range root {
		if err := ForEachNode(key, node, act); err != nil {
			return err
		}
	}
	return nil
}

// StartZkCluster starts a managed ZooKeeper with given size of nodes.
func StartZkCluster(size int, stdout, stderr io.Writer) (*ZkCluster, error) {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	if size < 1 {
		return nil, errors.New("Cluster size must be larger than 0")
	}

	zk, err := zk.StartTestCluster(size, stdout, stderr)
	if err != nil {
		return nil, err
	}
	var cluster = &ZkCluster{TestCluster: zk}
	err = cluster.startAction()
	return cluster, err
}
