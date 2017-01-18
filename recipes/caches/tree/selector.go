package tree

// Selector controls the nodes a Cache processes.
// When iterating over the children of a parent node, a given node's children
// are queried only if TraverseChildren returns true.
// When caching the list of nodes for a parent node, a given node is
// stored only if AcceptChild returns true.
type Selector interface {
	TraverseChildren(string) bool
	AcceptChildData(string) bool
}
