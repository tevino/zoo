package tree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tevino/zoo/test"
)

func TestStartCache(t *testing.T) {
	test.NewZkEnv(t).With(func(zkEnv *test.ZkEnv) {
		var client = zkEnv.NewClientTimeout(time.Second * 2)
		var cache = NewCache(client, "/", nil)
		assert.NoError(t, cache.Start())
	})

}

func TestInitEvent(t *testing.T) {
	test.NewZkEnv(t).With(func(zkEnv *test.ZkEnv) {
		var cache = NewCache(zkEnv.NewClientTimeout(time.Second*2), "/", nil)

		var initReceived = make(chan struct{}, 1)
		func() {
			var listener = NewCacheEventListener(func(e CacheEvent) {
				t.Log(e)
				if CacheEventInitialized == e.Type {
					initReceived <- struct{}{}
				}
			})
			cache.AddEventListener(listener)
		}()

		assert.NoError(t, cache.Start())
		select {
		case <-time.After(time.Second):
			t.Fatal("Waiting for init event timed out")
		case <-initReceived:
		}
	})
}
