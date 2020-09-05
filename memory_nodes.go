package goblin

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type memoryFileInfo struct {
	filename string
	modTime  time.Time
	isDir    bool
	mode     os.FileMode
	size     int64

	node fsNode
}

var _ os.FileInfo = &memoryFileInfo{}

func (mfi *memoryFileInfo) Name() string {
	return mfi.filename
}

func (mfi *memoryFileInfo) IsDir() bool {
	return mfi.isDir
}

func (mfi *memoryFileInfo) ModTime() time.Time {
	return mfi.modTime
}

func (mfi *memoryFileInfo) Mode() os.FileMode {
	return mfi.mode
}

func (mfi *memoryFileInfo) Size() int64 {
	return mfi.size
}

func (mfi *memoryFileInfo) Sys() interface{} {
	return mfi.node
}

type fsNode interface {
	Name() string
	FullPath() string

	Stat() (os.FileInfo, error)
	Open() (File, error)

	GetNode(path []string) (fsNode, error)
}

type memoryDir struct {
	fullPath string
	name     string
	modTime  time.Time
	nodes    map[string]fsNode
}

var _ File = &memoryDir{}
var _ fsNode = &memoryDir{}

func newMemoryDir(fullPath string, opts ...FileOption) *memoryDir {
	fOpts := newFileOptions(opts...)

	return &memoryDir{
		fullPath: fullPath,
		name:     path.Base(fullPath),
		modTime:  fOpts.ModTime,
		nodes:    map[string]fsNode{},
	}
}

func (d *memoryDir) CreateNode(path []string, node fsNode) error {
	err := validatePath(path)
	if err != nil {
		return err
	} else if node == nil {
		return fmt.Errorf("node cannot be nil")
	}

	var createNode func(*memoryDir, []string, string, []string, fsNode) error
	createNode = func(
		parentNode *memoryDir,
		parentPath []string,
		curPath string,
		nextPath []string,
		node fsNode,
	) error {
		curNode, ok := parentNode.nodes[curPath]
		if len(nextPath) == 0 {
			// Make sure the node doesn't exist before we try to write to it
			if ok {
				return os.ErrExist
			}

			parentNode.nodes[curPath] = node
			return nil
		}

		fullCurPath := append(parentPath, curPath)
		if !ok {
			// If we don't have a current node, we need to create the next
			// dir to fill the rest of the tree.
			curNode = newMemoryDir(strings.Join(fullCurPath, pathSeparator))
			parentNode.nodes[curPath] = curNode
		}

		// Make sure the node we're trying to add to is a directory. We can't add
		// additional nodes to a file...
		if dirNode, ok := curNode.(*memoryDir); ok {
			curPath, nextPath := nextPath[0], nextPath[1:]
			return createNode(dirNode, fullCurPath, curPath, nextPath, node)
		}

		return fmt.Errorf("attempted to add additional nodes to a file")
	}

	curPath, nextPath := path[0], path[1:]
	return createNode(d, []string{}, curPath, nextPath, node)
}

func (d *memoryDir) GetNode(path []string) (fsNode, error) {
	if len(path) == 0 ||
		(len(path) == 1 && path[0] == filesystemRootPath) {
		return d, nil
	}

	err := validatePath(path)
	if err != nil {
		return nil, err
	}

	curPath, nextPaths := path[0], path[1:]
	if curNode, ok := d.nodes[curPath]; ok {
		// If this is the end of the path, return what we have
		if len(nextPaths) == 0 {
			return curNode, nil
		}

		// If this still has additional path segments, traverse down deeper
		return curNode.GetNode(nextPaths)
	}

	return nil, os.ErrNotExist
}

func (d *memoryDir) Nodes() []fsNode {
	var res []fsNode
	for _, node := range d.nodes {
		res = append(res, node)
	}
	return res
}

func (d *memoryDir) Name() string {
	return d.name
}

func (d *memoryDir) FullPath() string {
	return d.fullPath
}

func (d *memoryDir) Close() error {
	return nil
}

func (d *memoryDir) Read(buf []byte) (int, error) {
	return 0, fmt.Errorf("not a file")
}

func (d *memoryDir) Stat() (os.FileInfo, error) {
	return &memoryFileInfo{
		filename: d.name,
		modTime:  d.modTime,
		isDir:    true,
		mode:     0,
		size:     0,
		node:     d,
	}, nil
}

func (d *memoryDir) Open() (File, error) {
	return &openMemoryDir{}, nil
}

type openMemoryDir struct {
	memoryDir
}

var _ File = &openMemoryDir{}

func newOpenMemoryDir(dir *memoryDir) *openMemoryDir {
	return &openMemoryDir{
		memoryDir: *dir,
	}
}

func (omd *openMemoryDir) Close() error {
	return fmt.Errorf("dir close... what it do?")
}

func (omd *openMemoryDir) Read(data []byte) (int, error) {
	return 0, fmt.Errorf("dir read... what it do?")
}

type memoryFile struct {
	fullPath string
	name     string
	modTime  time.Time
	data     []byte
}

var _ fsNode = &memoryFile{}

func newMemoryFile(fullPath string, data []byte, opts ...FileOption) *memoryFile {
	fOpts := newFileOptions(opts...)

	fileName := path.Base(fullPath)

	return &memoryFile{
		fullPath: fullPath,
		name:     fileName,
		modTime:  fOpts.ModTime,
		data:     data,
	}
}

func (f *memoryFile) GetNode(path []string) (fsNode, error) {
	if len(path) == 0 {
		return f, nil
	}

	return nil, fmt.Errorf("cannot get deeper nodes from file")
}

func (f *memoryFile) Name() string {
	return f.name
}

func (f *memoryFile) FullPath() string {
	return f.fullPath
}

func (f *memoryFile) Stat() (os.FileInfo, error) {
	return &memoryFileInfo{
		filename: f.name,
		modTime:  f.modTime,
		isDir:    false,
		mode:     0,
		size:     int64(len(f.data)),
		node:     f,
	}, nil
}

func (f *memoryFile) Open() (File, error) {
	openFile := &openMemoryFile{
		memoryFile: *f,
		curRead:    0,
		closed:     false,
	}

	return openFile, nil
}

type openMemoryFile struct {
	memoryFile
	curRead int
	closed  bool
}

var _ File = &openMemoryFile{}

func (omf *openMemoryFile) Read(buf []byte) (int, error) {
	if omf.closed {
		return 0, os.ErrClosed
	}

	bufLen := len(buf)
	readLeft := len(omf.data) - omf.curRead
	readLen := bufLen

	var readErr error
	if readLeft < readLen {
		readLen = readLeft
		readErr = io.EOF
	}

	copy(buf, omf.data[omf.curRead:omf.curRead+readLen])

	omf.curRead += readLen

	return readLen, readErr
}

func (omf *openMemoryFile) Close() error {
	omf.closed = true
	return nil
}
