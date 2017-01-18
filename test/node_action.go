package test

import (
	"fmt"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/tevino/zoo/enhanced"
)

// NodeActionFunc represents an action to a ZNode.
type NodeActionFunc func(fullPath string, node ZNode) error

// NodeAction contains NodeActionFunc(s) that are used while traversing a tree.
type NodeAction struct {
	client *enhanced.Client
}

// StartNewNodeAction create then starts NodeAction on given ZkCluster.
func StartNewNodeAction(c *ZkCluster) (*NodeAction, error) {
	client, err := c.Connect(0)
	if err != nil {
		return nil, fmt.Errorf("Error starting client: %s", err)
	}
	return &NodeAction{client: client}, nil
}

// Stop stops NodeAction.
func (a *NodeAction) Stop() {
	a.client.Close()
}

// DoCreate is a NodeActionFunc which creates a ZNode with its parents.
func (a *NodeAction) DoCreate(fullPath string, node ZNode) error {
	var err error

	if node.HasValue() {
		err = a.client.CreateValueWithParents(fullPath, []byte(*node.Value))
	} else {
		err = a.client.CreateWithParents(fullPath)
	}
	return err
}

// DoDelete is a NodeActionFunc which deletes a ZNode with its children.
func (a *NodeAction) DoDelete(fullPath string, node ZNode) error {
	return a.client.DeleteWithChildren(fullPath)
}

// DoUpdate is DoCreate with ErrNodeExists ignored.
func (a *NodeAction) DoUpdate(fullPath string, node ZNode) error {
	var err = a.DoCreate(fullPath, node)
	if err == zk.ErrNodeExists {
		err = nil
		if node.HasValue() {
			_, err = a.client.Set(fullPath, []byte(*node.Value), -1)
		}
	}
	return err
}
