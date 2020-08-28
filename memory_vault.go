package goblin

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	filesystemRootPath = "."
)

func makeMemoryVaultName(name string) string {
	return fmt.Sprintf("goblinMemoryVaultX%s", name)
}

type FileOption func(*fileOptions)

type fileOptions struct {
	ModTime time.Time
}

func newFileOptions(opts ...FileOption) *fileOptions {
	fo := &fileOptions{
		ModTime: time.Now(),
	}

	for _, opt := range opts {
		opt(fo)
	}

	return fo
}

func FileModTime(modTime time.Time) FileOption {
	return func(fOpts *fileOptions) {
		fOpts.ModTime = modTime
	}
}

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

type MemoryVault struct {
	root *memoryDir
}

var _ Vault = &MemoryVault{}

func NewMemoryVault(opts ...MemoryVaultOption) *MemoryVault {
	vaultOpts := newMemoryVaultOptions()
	for _, opt := range opts {
		opt(vaultOpts)
	}

	return &MemoryVault{
		// TODO give a real modtime
		root: vaultOpts.Root,
	}
}

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
