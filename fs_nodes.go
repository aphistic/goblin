package goblin

import "os"

type openFSFile struct {
	fullPath string
	f        *os.File
}

var _ File = &openFSFile{}

func newOpenFSFile(fullPath string, f *os.File) *openFSFile {
	return &openFSFile{
		fullPath: fullPath,
		f:        f,
	}
}

func (off *openFSFile) Stat() (os.FileInfo, error) {
	return os.Stat(off.fullPath)
}

func (off *openFSFile) Read(buf []byte) (int, error) {
	return off.f.Read(buf)
}

func (off *openFSFile) Close() error {
	return off.f.Close()
}
