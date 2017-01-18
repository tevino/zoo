package test

import (
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/tevino/zoo/enhanced"
)

// ZNodeAssertion contains helpers of ZNode assertion.
type ZNodeAssertion struct {
	t *testing.T
	ClientGetter
}

// AssertZNode asserts given ZNodes exist.
func (a *ZNodeAssertion) AssertZNode(pathes ...string) {
	for _, p := range pathes {
		var _, _, err = a.Client().Get(p)
		assert.NoError(a.t, err, "Error getting data of '%s': %s", p, err)
	}
}

// AssertZNodeWithValue asserts a ZNode with given value.
func (a *ZNodeAssertion) AssertZNodeWithValue(path, value string) {
	var data, _, err = a.Client().Get(path)
	assert.NoError(a.t, err, "Error getting data of '%s'", path)
	assert.Equal(a.t, value, string(data), "Value not match")
}

// AssertNoZNode asserts given ZNodes do not exist.
func (a *ZNodeAssertion) AssertNoZNode(pathes ...string) {
	for _, p := range pathes {
		var _, _, err = a.Client().Get(p)
		assert.Equal(a.t, zk.ErrNoNode, err, "Error getting data of '%s': %v", p, err)
	}
}

// ClientGetter represents a getter of *enhanced.Client.
type ClientGetter interface {
	Client() *enhanced.Client
}
