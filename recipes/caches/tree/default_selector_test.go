package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSelector(t *testing.T) {
	assert.True(t, DefaultSelector.AcceptChildData(""))
	assert.True(t, DefaultSelector.TraverseChildren(""))
}
