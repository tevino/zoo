package tree

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestEventListenersCount(t *testing.T) {
	var ls = NewCacheEventListeners()
	assert.Equal(t, 0, ls.Count())
}

func TestEventListenersAdd(t *testing.T) {
	var ls = NewCacheEventListeners()
	var fn = NewCacheEventListener(func(e CacheEvent) {})
	ls.Add(fn)
	assert.Equal(t, 1, ls.Count())
}

func TestEventListenersDel(t *testing.T) {
	var ls = NewCacheEventListeners()
	var fn = NewCacheEventListener(func(e CacheEvent) {})
	ls.Add(fn)
	ls.Del(fn)
	assert.Equal(t, 0, ls.Count())
}

func TestEventListenersBroadcast(t *testing.T) {
	var ls = NewCacheEventListeners()
	var called = 0
	var fn = NewCacheEventListener(func(e CacheEvent) { called++ })
	ls.Add(fn)
	ls.Broadcast(CacheEvent{})
	assert.Equal(t, 1, called)
}
