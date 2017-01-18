package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ClusterOperation does CUD operations of ZkCluster with NoError assertion.
type ClusterOperation struct {
	t *testing.T
	ClusterGetter
}

// MustCreate calls DoCreate on inner ZkCluster with NoError assertion.
func (o *ClusterOperation) MustCreate(yml string) {
	var err = o.Zk().DoCreate([]byte(yml))
	assert.NoError(o.t, err, "Error creating: \n%s\n", yml)
}

// MustDelete calls DoDelete on inner ZkCluster with NoError assertion.
func (o *ClusterOperation) MustDelete(yml string) {
	var err = o.Zk().DoDelete([]byte(yml))
	assert.NoError(o.t, err, "Error deleting: \n%s\n", yml)
}

// MustUpdate calls DoUpdate on inner ZkCluster with NoError assertion.
func (o *ClusterOperation) MustUpdate(yml string) {
	var err = o.Zk().DoUpdate([]byte(yml))
	assert.NoError(o.t, err, "Error updating: \n%s\n", yml)
}

// ClusterGetter represents a getter of *ZkCluster.
type ClusterGetter interface {
	Zk() *ZkCluster
}
