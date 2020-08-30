package goblin

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryVaultBinaryMarshalling(t *testing.T) {
	t.Run("marshal and unmarshal are idempotent", func(t *testing.T) {
		mv := NewMemoryVault()

		file1 := "test/path.txt"
		file2 := "different/path.txt"

		err := mv.WriteFile(
			file1, bytes.NewReader([]byte{0x01}),
			FileModTime(time.Unix(1, 0)),
		)
		require.NoError(t, err)

		err = mv.WriteFile(
			file2, bytes.NewReader([]byte{0x02}),
			FileModTime(time.Unix(2, 0)),
		)
		require.NoError(t, err)

		data, err := mv.MarshalBinary()
		require.NoError(t, err)

		mv = NewMemoryVault()
		err = mv.UnmarshalBinary(data)
		require.NoError(t, err)

		data, err = mv.ReadFile(file1)
		require.NoError(t, err)
		fInfo, err := mv.Stat(file1)
		require.NoError(t, err)
		assert.Equal(t, time.Unix(1, 0), fInfo.ModTime())
		assert.Equal(t, []byte{0x01}, data)

		data, err = mv.ReadFile(file2)
		require.NoError(t, err)
		fInfo, err = mv.Stat(file2)
		require.NoError(t, err)
		assert.Equal(t, time.Unix(2, 0), fInfo.ModTime())
		assert.Equal(t, []byte{0x02}, data)
	})
}
