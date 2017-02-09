package enhanced

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestClientWatchNS(t *testing.T) {
	var c = newClient(nil, nil).SetNamespace("prefix")
	assert.Equal(t, "/prefix/xxx", c.watchOperations.namespaced("xxx"))
	assert.Equal(t, "/prefix/xxx", c.namespaced("xxx"))
}

func TestEmptyClientWatchNS(t *testing.T) {
	var c = newClient(nil, nil)
	assert.Equal(t, "/xxx", c.watchOperations.namespaced("xxx"))
	assert.Equal(t, "/xxx", c.namespaced("xxx"))
}
