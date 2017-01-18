package tree

// FuncAcceptAll is a function that always returns true.
var FuncAcceptAll = func(string) bool { return true }

// DefaultSelector always returns true.
var DefaultSelector = NewSelector(FuncAcceptAll, FuncAcceptAll)

type generalSelector struct {
	traverseChildren func(string) bool
	acceptChildData  func(string) bool
}

func (s *generalSelector) TraverseChildren(fullPath string) bool {
	return s.traverseChildren(fullPath)
}

func (s *generalSelector) AcceptChildData(fullPath string) bool {
	return s.acceptChildData(fullPath)
}

// NewSelector creates Selector.
func NewSelector(traverseChildren, acceptChild func(string) bool) Selector {
	return &generalSelector{traverseChildren: traverseChildren, acceptChildData: acceptChild}
}
