package goblin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FilesystemVault struct {
	rootPath string
}

var _ Vault = &FilesystemVault{}

func NewFilesystemVault(rootPath string) *FilesystemVault {
	return &FilesystemVault{
		rootPath: rootPath,
	}
}

func (v *FilesystemVault) GoString() string {
	return `FilesystemVault{RootPath: "` + v.rootPath + `"}`
}

func (v *FilesystemVault) String() string {
	return `Filesystem Vault (` + v.rootPath + `)`
}

func (v *FilesystemVault) makePath(name string) (string, error) {
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

func (v *FilesystemVault) Open(name string) (File, error) {
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

func (v *FilesystemVault) Stat(name string) (os.FileInfo, error) {
	fullPath, err := v.makePath(name)
	if err != nil {
		return nil, err
	}

	return os.Stat(fullPath)
}

func (v *FilesystemVault) ReadDir(dirName string) ([]os.FileInfo, error) {
	fullPath, err := v.makePath(dirName)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadDir(fullPath)
}

func (v *FilesystemVault) Glob(pattern string) ([]string, error) {
	fullPattern, err := v.makePath(pattern)
	if err != nil {
		return nil, err
	}

	return filepath.Glob(fullPattern)
}

func (v *FilesystemVault) ReadFile(name string) ([]byte, error) {
	fullPath, err := v.makePath(name)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fullPath)
}
