package goblin

import (
	"fmt"
	"os"
)

type SelectOption func(*VaultSelector)

type SelectVaultMethod func() Vault

func SelectEnvBool(envName string, v Vault) SelectOption {
	return func(vs *VaultSelector) {
		vs.AppendSelector(func() Vault {
			// TODO implement this
			return nil
		})
	}
}

func SelectEnvPath(envName string, v Vault) SelectOption {
	return func(vs *VaultSelector) {
		vs.AppendSelector(func() Vault {
			// TODO implement this
			return nil
		})
	}
}

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
	return &VaultSelector{}
}

func (vs *VaultSelector) AppendSelector(svm SelectVaultMethod) {
	vs.vaultOptions = append(vs.vaultOptions, svm)
}

func (vs *VaultSelector) getVault() (Vault, error) {
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

func (vs *VaultSelector) Open(name string) (File, error) {
	v, err := vs.getVault()
	if err != nil {
		return nil, err
	}

	return v.Open(name)
}

func (vs *VaultSelector) Stat(name string) (os.FileInfo, error) {
	v, err := vs.getVault()
	if err != nil {
		return nil, err
	}

	return v.Stat(name)
}

func (vs *VaultSelector) ReadDir(dirName string) ([]os.FileInfo, error) {
	v, err := vs.getVault()
	if err != nil {
		return nil, err
	}

	return v.ReadDir(dirName)
}

func (vs *VaultSelector) Glob(pattern string) ([]string, error) {
	v, err := vs.getVault()
	if err != nil {
		return nil, err
	}

	return v.Glob(pattern)
}

func (vs *VaultSelector) ReadFile(name string) ([]byte, error) {
	v, err := vs.getVault()
	if err != nil {
		return nil, err
	}

	return v.ReadFile(name)
}
