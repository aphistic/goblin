package goblin

import (
	"fmt"
	"os"
	"strings"
)

func makeVaultName(name string) string {
	return fmt.Sprintf("goblinVaultX%s", name)
}

type Vault interface{}

type vault struct {
	name string

	root *fsDir
}

func newVault(vaultName string) *vault {
	return &vault{
		name: vaultName,
		root: newFsDir(),
	}
}

func (v *vault) SetFile(path string, data []byte) error {
	pathParts := strings.Split(path, string(os.PathSeparator))
	curRoot := v.root
	for idx := 0; idx < len(pathParts)-1; idx++ {
		part := pathParts[idx]
		newRoot, ok := curRoot.nodes[part]
		if !ok {
			newRoot = newFsDir()
			curRoot.nodes[part] = newRoot
		}

		dirRoot, ok := newRoot.(*fsDir)
		if !ok {
			return fmt.Errorf("%s is not a directory", part)
		}

		curRoot = dirRoot
	}

	f := newFsFile()
	f.data = make([]byte, len(data))
	copy(f.data, data)
	curRoot.nodes[pathParts[len(pathParts)-1]] = f

	return nil
}

func (v *vault) Files() map[string][]byte {
	return v.root.Files()
}
