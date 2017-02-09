package enhanced

import "github.com/samuel/go-zookeeper/zk"

type nsBasicOperations struct {
	basicOperations
	*namespace
}

func newNSBasicOperations(conner Conner, ns *namespace) nsBasicOperations {
	return nsBasicOperations{
		basicOperations: newBasicOperations(conner),
		namespace:       ns,
	}
}

// Get fetches value and stat of given znode.
func (nb *nsBasicOperations) Get(p string) ([]byte, *zk.Stat, error) {
	return nb.get(nb.namespaced(p))
}

// Exist returns true and stat of given znode if it exists.
func (nb *nsBasicOperations) Exist(p string) (bool, *zk.Stat, error) {
	return nb.exist(nb.namespaced(p))
}

// GetChildren fetches children of given path.
func (nb *nsBasicOperations) GetChildren(p string) ([]string, *zk.Stat, error) {
	return nb.getChildren(nb.namespaced(p))
}

// Set sets the value on given znode.
func (nb *nsBasicOperations) Set(p string, value []byte, version int32) (*zk.Stat, error) {
	return nb.set(nb.namespaced(p), value, version)
}

// Create creates given znode with value set to nil.
func (nb *nsBasicOperations) Create(p string) error {
	return nb.create(nb.namespaced(p))
}

// CreateValue creates given znode with value.
func (nb *nsBasicOperations) CreateValue(p string, value []byte) error {
	return nb.createValue(nb.namespaced(p), value)
}

// Delete deletes given znode.
func (nb *nsBasicOperations) Delete(p string, version int32) error {
	return nb.delete(nb.namespaced(p), version)
}

// DeleteWithChildren deletes given znode with its children if any.
func (nb *nsBasicOperations) DeleteWithChildren(p string) error {
	return nb.deleteWithChildren(nb.namespaced(p))
}

// CreateWithParents create path with its parents created if missing.
func (nb *nsBasicOperations) CreateWithParents(p string) error {
	return nb.createWithParents(nb.namespaced(p))
}

// CreateValueWithParents create path with value and its parents created if missing.
func (nb *nsBasicOperations) CreateValueWithParents(p string, value []byte) error {
	return nb.createValueWithParents(p, value)
}
