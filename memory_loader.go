package goblin

type LoadOption func(*loadOptions)

type loadOptions struct{}

func newLoadOptions() *loadOptions {
	return &loadOptions{}
}

// LoadVault takes a binary representation of a Goblin vault and unmarshales it into a vault.
func LoadMemoryVault(vaultData []byte, opts ...LoadOption) (Vault, error) {
	loadOpts := newLoadOptions()
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
