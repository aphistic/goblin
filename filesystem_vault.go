package goblin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FilesystemVault is a vault used to interact with a local filesystem. All paths
// provided to a FilesystemVault are relative to the root path.
type FilesystemVault struct {
	rootPath string
}

var _ Vault = &FilesystemVault{}

// NewFilesystemVault creates a FilesystemVault using the given root path as the
// root of the filesystem.
func NewFilesystemVault(rootPath string) *FilesystemVault {
	return &FilesystemVault{
		rootPath: rootPath,
	}
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

// Open returns a file at the given path relative to the vault's root path.
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

// Stat returns file info at the given path relative to the vault's root path.
func (v *FilesystemVault) Stat(name string) (os.FileInfo, error) {
	fullPath, err := v.makePath(name)
	if err != nil {
		return nil, err
	}

	return os.Stat(fullPath)
}

// ReadDir returns directory contents for the given path relative to the vault's
// root path.
func (v *FilesystemVault) ReadDir(dirName string) ([]os.FileInfo, error) {
	fullPath, err := v.makePath(dirName)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadDir(fullPath)
}

// Glob returns names of files in the filesystem that match the given pattern
// relative to the vault's root path.
func (v *FilesystemVault) Glob(pattern string) ([]string, error) {
	fullPattern, err := v.makePath(pattern)
	if err != nil {
		return nil, err
	}

	return filepath.Glob(fullPattern)
}

// ReadFile returns the contents of the file at the given path relative to the
// vault's root path.
func (v *FilesystemVault) ReadFile(name string) ([]byte, error) {
	fullPath, err := v.makePath(name)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(fullPath)
}
