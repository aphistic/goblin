package goblin

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestVault() *MemoryVault {
	// Out of order on purpose
	v := NewMemoryVault()

	_ = v.WriteFile("file.txt", bytes.NewBuffer([]byte{0x01}))

	_ = v.WriteFile("dir2/dir22/file.txt", bytes.NewBuffer([]byte{0x05}))
	_ = v.WriteFile("dir2/dir21/file.txt", bytes.NewBuffer([]byte{0x04}))

	_ = v.WriteFile("dir1/file.txt", bytes.NewBuffer([]byte{0x02}))
	_ = v.WriteFile("dir1/dir11/file.txt", bytes.NewBuffer([]byte{0x03}))

	return v
}

func TestMemoryVaultStringers(t *testing.T) {
	t.Run("string method", func(t *testing.T) {
		mv := NewMemoryVault()
		assert.Equal(t, "Memory Vault", mv.String())
	})
}

func TestMemoryVaultOpen(t *testing.T) {
	t.Run("open root directory", func(t *testing.T) {
		v := newTestVault()

		f, err := v.Open(filesystemRootPath)
		require.NoError(t, err)
		assert.NotNil(t, f)
	})
}

func TestMemoryVaultReadDir(t *testing.T) {
	t.Run("read root", func(t *testing.T) {
		v := newTestVault()
		fi, err := v.ReadDir(filesystemRootPath)
		require.NoError(t, err)

		require.Len(t, fi, 3)

		fi0 := fi[0]
		require.Equal(t, "dir1", fi0.Name())
		assert.Equal(t, true, fi0.IsDir())

		fi1 := fi[1]
		require.Equal(t, "dir2", fi1.Name())
		assert.Equal(t, true, fi1.IsDir())

		fi2 := fi[2]
		require.Equal(t, "file.txt", fi2.Name())
		assert.Equal(t, false, fi2.IsDir())
	})

	t.Run("read subdir", func(t *testing.T) {
		v := newTestVault()
		fi, err := v.ReadDir("dir1")
		require.NoError(t, err)

		require.Len(t, fi, 2)

		fi0 := fi[0]
		require.Equal(t, "dir11", fi0.Name())
		assert.Equal(t, true, fi0.IsDir())

		fi1 := fi[1]
		require.Equal(t, "file.txt", fi1.Name())
		assert.Equal(t, false, fi1.IsDir())
	})
}

func TestMemoryVaultGlob(t *testing.T) {
	t.Run("root glob", func(t *testing.T) {
		v := newTestVault()
		names, err := v.Glob(filesystemRootPath)
		require.NoError(t, err)
		require.Len(t, names, 3)
		assert.Equal(t, "dir1", names[0])
		assert.Equal(t, "dir2", names[1])
		assert.Equal(t, "file.txt", names[2])
	})

	t.Run("don't allow backslashes in pattern", func(t *testing.T) {
		v := newTestVault()
		names, err := v.Glob(`/dir1\dir2/*`)
		assert.EqualError(t, err, "backslash is not allowed in glob patterns")
		assert.Nil(t, names)
	})
}
