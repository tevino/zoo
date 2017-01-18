package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicYAMLParsing(t *testing.T) {
	for _, s := range []string{
		`/a:`,
		`
/a:
/b:
`, `
/a:
  value: 123
`, `
/a:
  value: 123
  children:
    b:
    c:
`, `
/a:
  value: 123
  children:
    b:
      value: I'm b
    c:
      value: I'm c, I have a child d
      children:
        d:
`, `
/one/two:
  value: asd
  children:
    three:
      value: I'm three
    3rd:
      value: I'm three
`, `
/one:
/two:
/three:
  children:
    four:
    five:
`,
	} {
		root, err := UnmarshalYAML([]byte(s))
		assert.NotNil(t, root)
		assert.NoError(t, err)
	}
}
