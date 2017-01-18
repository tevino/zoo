package test

import (
	"path"

	yaml "gopkg.in/yaml.v2"
)

// ZNode represents a ZooKeeper znode in YAML.
type ZNode struct {
	Value    *string          `yaml:"value,omitempty"`
	Children map[string]ZNode `yaml:"children,omitempty"`
}

// IsLeaf returns true if it has no children.
func (n *ZNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// HasValue returns true if Value is not empty.
func (n *ZNode) HasValue() bool {
	return n.Value != nil
}

// UnmarshalYAML parses bytes into a map of ZNodes.
func UnmarshalYAML(yml []byte) (map[string]ZNode, error) {
	var root map[string]ZNode
	var err = yaml.Unmarshal(yml, &root)
	return root, err
}

// ForEachNode calls do with all children nodes and itself.
func ForEachNode(fullPath string, node ZNode, do func(fullPath string, n ZNode) error) error {
	var err error

	if err = do(fullPath, node); err != nil {
		return err
	}

	for key, child := range node.Children {
		if err = ForEachNode(path.Join(fullPath, key), child, do); err != nil {
			return err
		}
	}
	return err
}
