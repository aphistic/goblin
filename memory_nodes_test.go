package goblin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryDirNaming(t *testing.T) {
	t.Run("root special case", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)
		assert.Equal(t, filesystemRootPath, d.FullPath())
		assert.Equal(t, filesystemRootPath, d.Name())
	})

	t.Run("root directory", func(t *testing.T) {
		d := newMemoryDir("dir1")
		assert.Equal(t, "dir1", d.FullPath())
		assert.Equal(t, "dir1", d.Name())
	})

	t.Run("deep directory", func(t *testing.T) {
		d := newMemoryDir("dir1/dir2/dir3")
		assert.Equal(t, "dir1/dir2/dir3", d.FullPath())
		assert.Equal(t, "dir3", d.Name())
	})

}
func TestMemoryDirCreateNode(t *testing.T) {
	t.Run("create root directory", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)
		err := d.CreateNode([]string{"dir1"}, newMemoryDir("dir1"))
		assert.NoError(t, err)

		assert.Contains(t, d.nodes, "dir1")
		assert.IsType(t, &memoryDir{}, d.nodes["dir1"])
	})

	t.Run("create directory tree", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)
		err := d.CreateNode(
			[]string{"dir1", "dir2", "dir3"},
			newMemoryDir("dir3", FileModTime(time.Unix(1234, 5678))),
		)
		require.NoError(t, err)
		nextDir := d

		require.Contains(t, nextDir.nodes, "dir1")
		require.IsType(t, &memoryDir{}, nextDir.nodes["dir1"])
		nextDir = nextDir.nodes["dir1"].(*memoryDir)

		require.Contains(t, nextDir.nodes, "dir2")
		require.IsType(t, &memoryDir{}, nextDir.nodes["dir2"])
		nextDir = nextDir.nodes["dir2"].(*memoryDir)

		require.Contains(t, nextDir.nodes, "dir3")
		require.IsType(t, &memoryDir{}, nextDir.nodes["dir3"])
		nextDir = nextDir.nodes["dir3"].(*memoryDir)

		assert.Equal(t, nextDir.name, "dir3")
		assert.Equal(t, nextDir.modTime, time.Unix(1234, 5678))
	})

	t.Run("create root file", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)
		err := d.CreateNode(
			[]string{"file.txt"},
			newMemoryFile("file.txt", []byte{0x01, 0x02}),
		)
		assert.NoError(t, err)

		assert.Contains(t, d.nodes, "file.txt")
		assert.IsType(t, &memoryFile{}, d.nodes["file.txt"])

		f := d.nodes["file.txt"].(*memoryFile)
		assert.Equal(t, "file.txt", f.name)
		assert.Equal(t, f.data, []byte{0x01, 0x02})
	})

	t.Run("create file in tree", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)
		err := d.CreateNode(
			[]string{"dir1", "dir2", "file.txt"},
			newMemoryFile("dir1/dir2/file.txt", []byte{0x01, 0x02}),
		)
		assert.NoError(t, err)
		nextDir := d

		require.Contains(t, nextDir.nodes, "dir1")
		require.IsType(t, &memoryDir{}, nextDir.nodes["dir1"])
		nextDir = nextDir.nodes["dir1"].(*memoryDir)

		require.Contains(t, nextDir.nodes, "dir2")
		require.IsType(t, &memoryDir{}, nextDir.nodes["dir2"])
		nextDir = nextDir.nodes["dir2"].(*memoryDir)

		require.Contains(t, nextDir.nodes, "file.txt")
		require.IsType(t, &memoryFile{}, nextDir.nodes["file.txt"])

		f := nextDir.nodes["file.txt"].(*memoryFile)
		assert.Equal(t, "file.txt", f.name)
		assert.Equal(t, f.data, []byte{0x01, 0x02})
	})

	t.Run("don't overwrite existing node", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)

		err := d.CreateNode([]string{"file.txt"}, newMemoryFile("file.txt", []byte{0x01}))
		assert.NoError(t, err)

		err = d.CreateNode([]string{"file.txt"}, newMemoryFile("file.txt", []byte{0x02}))
		assert.EqualError(t, err, "file already exists")

		require.Contains(t, d.nodes, "file.txt")
		require.IsType(t, &memoryFile{}, d.nodes["file.txt"])

		f := d.nodes["file.txt"].(*memoryFile)
		assert.Equal(t, f.data, []byte{0x01})
	})

	t.Run("empty path", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)
		err := d.CreateNode([]string{}, &memoryFile{})
		assert.EqualError(t, err, "path cannot be empty")
	})

	t.Run("nil node", func(t *testing.T) {
		d := newMemoryDir(filesystemRootPath)
		err := d.CreateNode([]string{"file.txt"}, nil)
		assert.EqualError(t, err, "node cannot be nil")
	})
}

func TestMemoryDirGetNode(t *testing.T) {
	t.Run("empty path returns current node", func(t *testing.T) {
		d := newMemoryDir("dir1")
		n, err := d.GetNode([]string{})
		require.NoError(t, err)
		assert.Equal(t, d, n)
	})

	t.Run("nil path returns current node", func(t *testing.T) {
		d := newMemoryDir("dir1")
		n, err := d.GetNode(nil)
		require.NoError(t, err)
		assert.Equal(t, d, n)
	})
}

func TestMemoryFileGetNode(t *testing.T) {
	t.Run("empty path returns current node", func(t *testing.T) {
		f := newMemoryFile("file.txt", []byte{0x01})
		n, err := f.GetNode([]string{})
		require.NoError(t, err)
		assert.Equal(t, f, n)
	})

	t.Run("nil path returns current node", func(t *testing.T) {
		f := newMemoryFile("file.txt", []byte{0x01})
		n, err := f.GetNode(nil)
		require.NoError(t, err)
		assert.Equal(t, f, n)
	})

	t.Run("any other path returns an error", func(t *testing.T) {
		f := newMemoryFile("file.txt", []byte{0x01})
		n, err := f.GetNode([]string{"foo"})
		assert.EqualError(t, err, "cannot get deeper nodes from file")
		assert.Nil(t, n)
	})
}
