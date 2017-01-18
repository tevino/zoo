package tree

import (
	"errors"
	"testing"

	"github.com/bmizerany/assert"
)

func TestErrorListenersCount(t *testing.T) {
	var ls = NewErrorListeners()
	assert.Equal(t, 0, ls.Count())
}

func TestErrorListenersAdd(t *testing.T) {
	var ls = NewErrorListeners()
	var fn = NewErrorListener(func(error) {})
	ls.Add(fn)
	assert.Equal(t, 1, ls.Count())
}

func TestErrorListenerDel(t *testing.T) {
	var ls = NewErrorListeners()
	var fn = NewErrorListener(func(error) {})
	ls.Add(fn)
	ls.Del(fn)
	assert.Equal(t, 0, ls.Count())
}

func TestErrorListenerBroadcast(t *testing.T) {
	var ls = NewErrorListeners()
	var called = 0
	var fn = NewErrorListener(func(error) { called++ })
	ls.Add(fn)
	ls.Broadcast(errors.New("test error"))
	assert.Equal(t, 1, called)
}
