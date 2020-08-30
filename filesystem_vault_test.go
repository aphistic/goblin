package goblin

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFSVaultStringer(t *testing.T) {
	t.Run("string method", func(t *testing.T) {
		td, err := ioutil.TempDir("", testTempPattern)
		require.NoError(t, err)
		defer os.RemoveAll(td)

		v := NewFilesystemVault(td)
		assert.Equal(t, "Filesystem Vault ("+td+")", v.String())
	})
}

func TestFSVaultMakePath(t *testing.T) {
	t.Run("standard path", func(t *testing.T) {
		td, err := ioutil.TempDir("", testTempPattern)
		require.NoError(t, err)
		defer os.RemoveAll(td)

		v := NewFilesystemVault(td)

		p, err := v.makePath("this/is/a/file.txt")
		require.NoError(t, err)
		assert.Equal(t, path.Join(td, "this/is/a/file.txt"), p)
	})

	t.Run("root path", func(t *testing.T) {
		td, err := ioutil.TempDir("", testTempPattern)
		require.NoError(t, err)
		defer os.RemoveAll(td)

		v := NewFilesystemVault(td)

		p, err := v.makePath(".")
		require.NoError(t, err)
		assert.Equal(t, td, p)
	})
}
