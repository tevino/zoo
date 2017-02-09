package enhanced

import (
	"path"

	"github.com/samuel/go-zookeeper/zk"
)

type basicOperations struct {
	Conner
	flags int32
	acl   []zk.ACL
}

func newBasicOperations(conner Conner) basicOperations {
	return basicOperations{
		Conner: conner,
		flags:  0,
		acl:    zk.WorldACL(zk.PermAll),
	}
}

// SetFlags sets the flags used for all operation.
func (o *basicOperations) SetFlags(flags int32) {
	o.flags = flags
}

// SetACL sets the ACL used for all operation.
func (o *basicOperations) SetACL(acl []zk.ACL) {
	o.acl = acl
}

func (o *basicOperations) get(p string) ([]byte, *zk.Stat, error) {
	return o.Conn().Get(p)
}

func (o *basicOperations) exist(p string) (bool, *zk.Stat, error) {
	return o.Conn().Exists(p)
}

func (o *basicOperations) getChildren(p string) ([]string, *zk.Stat, error) {
	return o.Conn().Children(p)
}

func (o *basicOperations) set(p string, value []byte, version int32) (*zk.Stat, error) {
	return o.Conn().Set(p, value, version)
}

func (o *basicOperations) create(p string) error {
	_, err := o.Conn().Create(p, nil, o.flags, o.acl)
	return err
}

func (o *basicOperations) createValue(p string, value []byte) error {
	_, err := o.Conn().Create(p, value, o.flags, o.acl)
	return err
}

func (o *basicOperations) delete(p string, version int32) error {
	return o.Conn().Delete(p, version)
}

func (o *basicOperations) deleteWithChildren(p string) error {
	var children, _, err = o.getChildren(p)
	if err != nil {
		return err
	}
	for _, child := range children {
		err = o.deleteWithChildren(path.Join(p, child))
		if err != nil {
			return err
		}
	}
	return o.delete(p, -1)
}

func (o *basicOperations) createWithParents(p string) error {
	var err = o.create(p)
	if err == zk.ErrNoNode {
		var parent = path.Dir(p)
		err = o.createWithParents(parent)
		if err == nil {
			err = o.create(p)
		}
	}
	return err
}

func (o *basicOperations) createValueWithParents(p string, value []byte) error {
	var err = o.createValue(p, value)
	if err == zk.ErrNoNode {
		var parent = path.Dir(p)
		err = o.createWithParents(parent)
		if err == nil {
			err = o.createValue(p, value)
		}
	}
	return err
}
