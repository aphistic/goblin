package goblin

import (
	"path"
)

type fsNode interface{}

type fsFile struct {
	data []byte
}

func newFsFile() *fsFile {
	return &fsFile{}
}

type fsDir struct {
	nodes map[string]fsNode
}

func newFsDir() *fsDir {
	return &fsDir{
		nodes: map[string]fsNode{},
	}
}

func (d *fsDir) Files() map[string][]byte {
	files := map[string][]byte{}

	for name, node := range d.nodes {
		dir, ok := node.(*fsDir)
		if ok {
			childFiles := dir.Files()
			for fileName, fileData := range childFiles {
				filePath := path.Join(name, fileName)
				files[filePath] = fileData
			}
		}

		f, ok := node.(*fsFile)
		if ok {
			files[name] = f.data
		}
	}

	return files
}
