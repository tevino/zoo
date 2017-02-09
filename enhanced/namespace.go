package enhanced

import "path"

type namespace string

func (n *namespace) namespaced(p string) string {
	return path.Join("/", string(*n), p)
}

func (n *namespace) ns() string {
	return string(*n)
}

func (n *namespace) setNS(ns string) {
	*n = namespace(ns)

}
