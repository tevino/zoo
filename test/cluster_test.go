package test

import "testing"

func TestCreate(t *testing.T) {
	NewZkEnv(t).With(func(env *ZkEnv) {
		env.MustCreate(`
/one/two/three:
  value: 3
  children:
    x:
      value: ex
    y:
      value: why
    z:
/a/b/c:
`)

		env.AssertZNode(
			"/one/two/three/x",
			"/one/two/three/y",
			"/one/two/three/z",
			"/a/b/c",
		)
		env.AssertZNodeWithValue("/one/two/three", "3")
		env.AssertZNodeWithValue("/one/two/three/x", "ex")
		env.AssertZNodeWithValue("/one/two/three/y", "why")
	})
}

func TestDelete(t *testing.T) {
	NewZkEnv(t).With(func(env *ZkEnv) {
		env.MustCreate(`
/one/two/three:
  children:
    x:
    y:
    z:
/a/b/c:
/a/b/c/d/e/f/g:
`)
		env.MustDelete(`
/a/b/c/d:
/one/two/three/x:
/one/two/three/z:
`)

		env.AssertNoZNode(
			"/a/b/c/d",
			"/one/two/three/x",
			"/one/two/three/z",
		)
		env.AssertZNode(
			"/a/b/c",
			"/one/two/three/y",
		)
	})
}

func TestUpdate(t *testing.T) {
	NewZkEnv(t).With(func(env *ZkEnv) {
		env.MustCreate(`
/one/two/three:
  children:
    x:
    y:
    z:
/a/b/c/d/e/f/g:
`)
		env.MustUpdate(`
/one:
  value: 1
/one/two:
  value: 2
/one/two/three/y:
  value: Y
/one/two/three/z:
  value: Z
`)

		env.AssertZNodeWithValue("/one", "1")
		env.AssertZNodeWithValue("/one/two", "2")
		env.AssertZNodeWithValue("/one/two/three/y", "Y")
		env.AssertZNodeWithValue("/one/two/three/z", "Z")
	})
}
