package tree

import "sync/atomic"

const (
	// CacheStateLatent indicates Cache.Start() has not yet been called.
	CacheStateLatent CacheState = iota
	// CacheStateStarted indicates Cache.Start() has been called.
	CacheStateStarted
	// CacheStateStopped indicates Cache.Close() has been called.
	CacheStateStopped
)

// CacheState represents the state of a Cache.
type CacheState int32

// SetValueIf does a CAS operation.
func (s *CacheState) SetValueIf(old, new CacheState) bool {
	return atomic.CompareAndSwapInt32((*int32)(s), int32(old), int32(new))
}

// Value returns the state.
func (s *CacheState) Value() CacheState {
	return CacheState(atomic.LoadInt32((*int32)(s)))
}

// EqualTo returns true if value equals to given CacheState.
func (s *CacheState) EqualTo(x CacheState) bool {
	return s.Value() == x
}
