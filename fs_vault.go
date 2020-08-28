package goblin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FSVault struct {
	rootPath string
}

var _ Vault = &FSVault{}

func NewFSVault(rootPath string) *FSVault {
	return &FSVault{
		rootPath: rootPath,
	}
}

func (v *FSVault) makePath(name string) (string, error) {
	if name == filesystemRootPath {
		return v.rootPath, nil
	}

	fullPath := v.rootPath
	if !strings.HasSuffix(v.rootPath, pathSeparator) &&
		!strings.HasPrefix(name, pathSeparator) {
		fullPath += pathSeparator
	}
	fullPath += name

	return fullPath, nil
}

func (v *FSVault) Open(name string) (File, error) {
	fullPath, err := v.makePath(name)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return newOpenFSFile(fullPath, f), nil
}

func (v *FSVault) Stat(name string) (os.FileInfo, error) {
	fullPath, err := v.makePath(name)
	if err != nil {
		return nil, err
	}

	return os.Stat(fullPath)
}

func (v *FSVault) ReadDir(dirName string) ([]os.FileInfo, error) {
	fullPath, err := v.makePath(dirName)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadDir(fullPath)
}

func (v *FSVault) Glob(pattern string) ([]string, error) {
	fullPattern, err := v.makePath(pattern)
	if err != nil {
		return nil, err
	}

	return filepath.Glob(fullPattern)
}

func (v *FSVault) ReadFile(name string) ([]byte, error) {
	fullPath, err := v.makePath(name)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fullPath)
}
