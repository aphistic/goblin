package goblin

import (
	"fmt"
	"os"
	"strings"
)

type SelectOption func(*VaultSelector)

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

type VaultSelector struct {
	selectedVault Vault
	vaultOptions  []SelectVaultMethod
}

var _ Vault = &VaultSelector{}

func NewVaultSelector(opts ...SelectOption) *VaultSelector {
	vs := &VaultSelector{}
	for _, opt := range opts {
		opt(vs)
	}

	return vs
}

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

func (vs *VaultSelector) GoString() string {
	v, err := vs.GetVault()
	if err != nil {
		return `VaultSelector{Vault: Error{Could not select}}`
	}

	return `VaultSelector{Vault: ` + v.GoString() + `}`
}

func (vs *VaultSelector) String() string {
	v, err := vs.GetVault()
	if err != nil {
		return `Vault Selector (Could not select vault)`
	}

	return `Vault Selector (` + v.String() + `)`
}

func (vs *VaultSelector) AppendSelector(svm SelectVaultMethod) {
	vs.vaultOptions = append(vs.vaultOptions, svm)
}

func (vs *VaultSelector) Open(name string) (File, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.Open(name)
}

func (vs *VaultSelector) Stat(name string) (os.FileInfo, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.Stat(name)
}

func (vs *VaultSelector) ReadDir(dirName string) ([]os.FileInfo, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.ReadDir(dirName)
}

func (vs *VaultSelector) Glob(pattern string) ([]string, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.Glob(pattern)
}

func (vs *VaultSelector) ReadFile(name string) ([]byte, error) {
	v, err := vs.GetVault()
	if err != nil {
		return nil, err
	}

	return v.ReadFile(name)
}
