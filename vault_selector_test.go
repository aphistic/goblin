package goblin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVaultSelectorStringers(t *testing.T) {
	t.Run("string method", func(t *testing.T) {
		vs := NewVaultSelector(
			SelectDefault(NewMemoryVault()),
		)
		assert.Equal(t, "Vault Selector (Memory Vault)", vs.String())
	})
	t.Run("gostring method", func(t *testing.T) {
		vs := NewVaultSelector(
			SelectDefault(NewMemoryVault()),
		)
		assert.Equal(t, "VaultSelector{Vault: MemoryVault{}}", vs.GoString())
	})
}
