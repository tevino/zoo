package tree

import "sync/atomic"

// NodeState represents state of a Node.
type NodeState int32

// Available node states.
const (
	NodeStatePending NodeState = iota
	NodeStateLive
	NodeStateDead
)

// Value returns the current state.
func (s *NodeState) Value() NodeState {
	return NodeState(atomic.LoadInt32((*int32)(s)))
}

// SetValueIf does a CAS operation.
func (s *NodeState) SetValueIf(old, new NodeState) bool {
	return atomic.CompareAndSwapInt32((*int32)(s), int32(old), int32(new))
}

// SetValue sets the state to given value.
func (s *NodeState) SetValue(new NodeState) {
	atomic.StoreInt32((*int32)(s), int32(new))
}

// SwapValue sets the state to new returns the old one.
func (s *NodeState) SwapValue(new NodeState) NodeState {
	return NodeState(atomic.SwapInt32((*int32)(s), int32(new)))
}

// EqualTo returns true if value equals to given NodeState.
func (s *NodeState) EqualTo(x NodeState) bool {
	return s.Value() == x
}
