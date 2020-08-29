package goblin

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	filesystemRootPath = "."
)

func makeMemoryVaultName(name string) string {
	return fmt.Sprintf("goblinMemoryVaultX%s", name)
}

// MemoryVaultOption is an option used when creating a memory vault.
type MemoryVaultOption func(*memoryVaultOptions)

type memoryVaultOptions struct {
	Root *memoryDir
}

func newMemoryVaultOptions() *memoryVaultOptions {
	return &memoryVaultOptions{
		Root: newMemoryDir(filesystemRootPath),
	}
}

func memoryVaultRoot(root *memoryDir) MemoryVaultOption {
	return func(vaultOpts *memoryVaultOptions) {
		vaultOpts.Root = root
	}
}

// MemoryVault is a vault stored in memory. It can be used as a temporary
// in-memory filesystem or as a way to load a filesystem from binary data,
// such as one embedded in a file.
type MemoryVault struct {
	root *memoryDir
}

var _ Vault = &MemoryVault{}

// NewMemoryVault creates a new memory vault.
func NewMemoryVault(opts ...MemoryVaultOption) *MemoryVault {
	vaultOpts := newMemoryVaultOptions()
	for _, opt := range opts {
		opt(vaultOpts)
	}

	return &MemoryVault{
		root: vaultOpts.Root,
	}
}

func (v *MemoryVault) String() string {
	return "Memory Vault"
}

// WriteFile reads data from the provided io.Reader and then writes it to the memory vault
// at the provided path.
func (v *MemoryVault) WriteFile(path string, r io.Reader, opts ...FileOption) error {
	pathParts := strings.Split(path, string(os.PathSeparator))
	curRoot := v.root
	for idx := 0; idx < len(pathParts)-1; idx++ {
		part := pathParts[idx]
		newRoot, ok := curRoot.nodes[part]
		if !ok {
			// TODO set a real modtime
			curFullPath := strings.Join(pathParts[:idx+1], pathSeparator)
			newRoot = newMemoryDir(curFullPath)
			curRoot.nodes[part] = newRoot
		}

		dirRoot, ok := newRoot.(*memoryDir)
		if !ok {
			return fmt.Errorf("%s is not a directory", part)
		}

		curRoot = dirRoot
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f := newMemoryFile(path, make([]byte, len(data)), opts...)
	copy(f.data, data)
	curRoot.nodes[pathParts[len(pathParts)-1]] = f

	return nil
}

// Open will open the file at the provided path from the in-memory vault.
func (v *MemoryVault) Open(name string) (File, error) {
	tokens, err := splitPath(name)
	if err != nil {
		return nil, err
	}

	f, err := v.root.GetNode(tokens)
	if err != nil {
		return nil, err
	}

	if memFile, ok := f.(*memoryFile); ok {
		return memFile.Open()
	}

	return nil, fmt.Errorf("cannot open directory")
}

// Stat returns file info for the provided path in the in-memory vault.
func (v *MemoryVault) Stat(name string) (os.FileInfo, error) {
	tokens, err := splitPath(name)
	if err != nil {
		return nil, err
	}

	n, err := v.root.GetNode(tokens)
	if err != nil {
		return nil, err
	}

	return n.Stat()
}

// ReadDir returns a slice of file info for the provided directory in the in-memory vault.
func (v *MemoryVault) ReadDir(dirName string) ([]os.FileInfo, error) {
	var res []os.FileInfo

	var pathTokens []string
	var err error
	if strings.TrimSpace(dirName) == filesystemRootPath {
		pathTokens = []string{}
	} else {
		pathTokens, err = splitPath(dirName)
		if err != nil {
			return nil, err
		}
	}

	node, err := v.root.GetNode(pathTokens)
	if err != nil {
		return nil, err
	}

	dirNode, ok := node.(*memoryDir)
	if !ok {
		return nil, fmt.Errorf("not a directory: %s", dirName)
	}

	for _, node := range dirNode.Nodes() {
		fi, err := node.Stat()
		if err != nil {
			return nil, err
		}
		res = append(res, fi)
	}

	// ReadDir returns contents in filename order
	sort.Slice(res, func(i, j int) bool {
		rI := res[i]
		rJ := res[j]

		return strings.Compare(rI.Name(), rJ.Name()) == -1
	})

	return res, nil
}

// Glob returns names of files in the in-memory vault that match the given pattern.
func (v *MemoryVault) Glob(pattern string) ([]string, error) {
	// Naive implementation that just navigates the whole FS tree
	// and runs filepath.Match on all the paths.

	if strings.ContainsRune(pattern, '\\') {
		return nil, fmt.Errorf("backslash is not allowed in glob patterns")
	}

	var glob func(string, *memoryDir) ([]string, error)
	glob = func(pattern string, dir *memoryDir) ([]string, error) {
		var res []string
		for _, node := range dir.nodes {
			if dirNode, ok := node.(*memoryDir); ok {
				dirNodes, err := glob(pattern, dirNode)
				if err != nil {
					return nil, err
				}
				res = append(res, dirNodes...)
			}

			nodePath := node.FullPath()
			match, err := filepath.Match(pattern, nodePath)
			if err != nil {
				return nil, err
			}

			if match {
				res = append(res, node.FullPath())
			}
		}
		return res, nil
	}

	if strings.TrimSpace(pattern) == filesystemRootPath {
		pattern = "*"
	}

	res, err := glob(pattern, v.root)
	if err != nil {
		return nil, err
	}

	// Sort the glob results by name
	sort.Strings(res)

	return res, nil
}

// ReadFile returns the contents of the file at the given path from
// the in-memory vault.
func (v *MemoryVault) ReadFile(name string) ([]byte, error) {
	f, err := v.Open(name)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
