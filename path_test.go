package goblin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitPath(t *testing.T) {
	t.Run("split root identifier", func(t *testing.T) {
		tokens, err := splitPath(filesystemRootPath)
		require.NoError(t, err)
		assert.Equal(t, []string{filesystemRootPath}, tokens)
	})

	t.Run("split file in root", func(t *testing.T) {
		tokens, err := splitPath("file.txt")
		require.NoError(t, err)
		assert.Equal(t, []string{"file.txt"}, tokens)
	})

	t.Run("split file in subdir", func(t *testing.T) {
		tokens, err := splitPath("dir1/file.txt")
		require.NoError(t, err)
		assert.Equal(t, []string{"dir1", "file.txt"}, tokens)
	})

	t.Run("split file in deep subdir", func(t *testing.T) {
		tokens, err := splitPath("dir1/dir2/dir3/file.txt")
		require.NoError(t, err)
		assert.Equal(t, []string{"dir1", "dir2", "dir3", "file.txt"}, tokens)
	})

	t.Run("split root dir", func(t *testing.T) {
		tokens, err := splitPath("dir1")
		require.NoError(t, err)
		assert.Equal(t, []string{"dir1"}, tokens)
	})

	t.Run("split deep dir", func(t *testing.T) {
		tokens, err := splitPath("dir1/dir2/dir3/dir4")
		require.NoError(t, err)
		assert.Equal(t, []string{"dir1", "dir2", "dir3", "dir4"}, tokens)
	})
}

func TestValidatePath(t *testing.T) {
	t.Run("allow '.' as the path", func(t *testing.T) {
		err := validatePath([]string{filesystemRootPath})
		assert.NoError(t, err)
	})
	t.Run("disallow . in path root", func(t *testing.T) {
		err := validatePath([]string{".", "dir2", "dir3", "dir4"})
		assert.EqualError(t, err, ". is not allowed in paths")
	})

	t.Run("disallow . in path", func(t *testing.T) {
		err := validatePath([]string{"dir1", "dir2", ".", "dir4"})
		assert.EqualError(t, err, ". is not allowed in paths")
	})

	t.Run("disallow . at end of path", func(t *testing.T) {
		err := validatePath([]string{"dir1", "dir2", "dir3", "."})
		assert.EqualError(t, err, ". is not allowed in paths")
	})

	t.Run("disallow .. in path root", func(t *testing.T) {
		err := validatePath([]string{"..", "dir2", "dir3", "dir4"})
		assert.EqualError(t, err, ".. is not allowed in paths")
	})

	t.Run("disallow .. in path", func(t *testing.T) {
		err := validatePath([]string{"dir1", "dir2", "..", "dir4"})
		assert.EqualError(t, err, ".. is not allowed in paths")
	})

	t.Run("disallow .. at end of path", func(t *testing.T) {
		err := validatePath([]string{"dir1", "dir2", "dir3", ".."})
		assert.EqualError(t, err, ".. is not allowed in paths")
	})

	t.Run("disallow path starting with separator", func(t *testing.T) {
		err := validatePath([]string{pathSeparator, "dir1", "dir2", "dir3", "dir4"})
		assert.EqualError(t, err, "path cannot be an absolute path")
	})

	t.Run("disallow empty path segments", func(t *testing.T) {
		err := validatePath([]string{"dir1", "", "dir3", "dir4"})
		assert.EqualError(t, err, "path cannot contain empty segments: dir1//dir3/dir4")
	})

	t.Run("disallow empty path", func(t *testing.T) {
		err := validatePath([]string{})
		assert.EqualError(t, err, "path cannot be empty")
	})

	t.Run("disallow nil path", func(t *testing.T) {
		err := validatePath(nil)
		assert.EqualError(t, err, "path cannot be empty")
	})
}
