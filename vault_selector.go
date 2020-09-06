package goblin

import (
	"fmt"
	"os"
	"strings"
)

// SelectOption is an option used when creating a vault selector.
type SelectOption func(*VaultSelector)

// SelectVaultMethod is the method signature a vault selector option needs to implement.
type SelectVaultMethod func() Vault

// SelectEnvBool will use the provided vault if the given environment variable is
// true. Accepted truthy values are TRUE, T, YES, Y, and 1. Values are case-insensitive.
func SelectEnvBool(envName string, v Vault) SelectOption {
	return func(vs *VaultSelector) {
		vs.AppendSelector(func() Vault {
			val := strings.TrimSpace(strings.ToUpper(os.Getenv(envName)))
			if val == "TRUE" || val == "T" ||
				val == "YES" || val == "Y" ||
				val == "1" {
				return v
			}

			return nil
		})
	}
}

// SelectEnvNotEmpty will use the provided vault if the given environment variable is
// set with a non-empty value.
func SelectEnvNotEmpty(envName string, v Vault) SelectOption {
	return func(vs *VaultSelector) {
		vs.AppendSelector(func() Vault {
			val := strings.TrimSpace(os.Getenv(envName))
			if val != "" {
				return v
			}

			return nil
		})
	}
}

// SelectPath will use the provided vault if the given path exists on disk.
func SelectPath(path string, v Vault) SelectOption {
	return func(vs *VaultSelector) {
		vs.AppendSelector(func() Vault {
			_, err := os.Stat(path)
			if err != nil {
				return nil
			}

			return v
		})
	}
}

// SelectDefault will always return the provided vault.
func SelectDefault(v Vault) SelectOption {
	return func(vs *VaultSelector) {
		vs.AppendSelector(func() Vault {
			return v
		})
	}
}

// VaultSelector is a vault that will use a vault selected by one of the added
// selectors.
type VaultSelector struct {
	selectedVault Vault
	vaultOptions  []SelectVaultMethod
}

var _ Vault = &VaultSelector{}

// NewVaultSelector creates a new vault selector.
func NewVaultSelector(opts ...SelectOption) *VaultSelector {
	vs := &VaultSelector{}
	for _, opt := range opts {
		opt(vs)
	}

	return vs
}

// GetVault returns the vault to be used by the vault selector using the
// currently provided selectors.
func (vs *VaultSelector) GetVault() (Vault, error) {
	if vs.selectedVault != nil {
		return vs.selectedVault, nil
	}

	for _, selectVault := range vs.vaultOptions {
		v := selectVault()
		if v != nil {
			vs.selectedVault = v
			return v, nil
		}
	}

	return nil, fmt.Errorf("no vault could be selected")
}

func (vs *VaultSelector) String() string {
	v, err := vs.GetVault()
	if err != nil {
		return `Vault Selector (Could not select vault)`
	}

	return `Vault Selector (` + v.String() + `)`
}

// AppendSelector adds an additional vault selector to the end of the vault
// selector list.
func (vs *VaultSelector) AppendSelector(svm SelectVaultMethod) {
	vs.vaultOptions = append(vs.vaultOptions, svm)
	vs.selectedVault = nil
}

// Open will open the file at the provided path from the selected vault.
func (vs *VaultSelector) Open(name string) (File, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.Open(name)
}

// Stat returns file info for the provided path from the selected vault.
func (vs *VaultSelector) Stat(name string) (os.FileInfo, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.Stat(name)
}

// ReadDir returns a slice of file info for the provided directory from the
// selected vault.
func (vs *VaultSelector) ReadDir(dirName string) ([]os.FileInfo, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.ReadDir(dirName)
}

// Glob returns names of files in the selected vault that match the given
// pattern.
func (vs *VaultSelector) Glob(pattern string) ([]string, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	if gv, ok := v.(GlobVault); ok {
		return gv.Glob(pattern)
	}

	return nil, fmt.Errorf("not supported")
}

// ReadFile returns the contents of the file at the given path from the
// selected vault.
func (vs *VaultSelector) ReadFile(name string) ([]byte, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.ReadFile(name)
}
