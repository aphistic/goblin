package goblin

// LoadMemoryOption is an option used when loading an in-memory vault.
type LoadMemoryOption func(*loadMemoryOptions)

type loadMemoryOptions struct{}

func newLoadMemoryOptions() *loadMemoryOptions {
	return &loadMemoryOptions{}
}

// LoadMemoryVault takes a binary representation of a memory vault and unmarshales it into a vault.
func LoadMemoryVault(vaultData []byte, opts ...LoadMemoryOption) (Vault, error) {
	loadOpts := newLoadMemoryOptions()
	for _, opt := range opts {
		opt(loadOpts)
	}

	v := NewMemoryVault()
	err := v.UnmarshalBinary(vaultData)
	if err != nil {
		return nil, err
	}

	return v, nil
}
