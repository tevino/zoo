package enhanced

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestEmptyNS(t *testing.T) {
	var ns namespace
	assert.Equal(t, "", ns.ns())
	assert.Equal(t, "/", ns.namespaced(""))
}

func TestNSRedundantSlash(t *testing.T) {
	var ns namespace
	assert.Equal(t, "/x", ns.namespaced("/x"))
	assert.Equal(t, "/x", ns.namespaced("//x"))
	assert.Equal(t, "/x", ns.namespaced("///x"))

	ns.setNS("prefix")
	assert.Equal(t, "/prefix/x", ns.namespaced("/x"))
	assert.Equal(t, "/prefix/x", ns.namespaced("//x"))
	assert.Equal(t, "/prefix/x", ns.namespaced("///x"))
}

func TestNSSet(t *testing.T) {
	var ns namespace
	ns.setNS("prefix")
	assert.Equal(t, "prefix", ns.ns())
}

func TestNSNamespaced(t *testing.T) {
	var ns namespace
	ns.setNS("prefix")
	assert.Equal(t, "/prefix/xxx", ns.namespaced("xxx"))
}

func TestNSGet(t *testing.T) {
	var ns namespace
	assert.Equal(t, "", ns.ns())
	ns.setNS("xxx")
	assert.Equal(t, "xxx", ns.ns())
}
