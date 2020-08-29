package goblin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSVaultStringer(t *testing.T) {
	t.Run("string method", func(t *testing.T) {
		v := NewFilesystemVault("/")
		assert.Equal(t, "Filesystem Vault (/)", v.String())
	})
}
