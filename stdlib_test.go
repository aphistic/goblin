package goblin

import (
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalk(t *testing.T) {
	t.Run("walk", func(t *testing.T) {
		v := newTestVault()

		var paths []string
		err := Walk(v, ".", func(path string, info os.FileInfo, err error) error {
			paths = append(paths, path)
			return nil
		})
		require.NoError(t, err)

		sort.Strings(paths)

		assert.Equal(t,
			[]string{
				".",
				"dir1",
				"dir1/dir11",
				"dir1/dir11/file.txt",
				"dir1/file.txt",
				"dir2",
				"dir2/dir21",
				"dir2/dir21/file.txt",
				"dir2/dir22",
				"dir2/dir22/file.txt",
				"file.txt",
			},
			paths,
		)
	})
}
