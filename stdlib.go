package goblin

// These are all types or methods that mirror the proposed
// fsys options from the Go standard library.
//
// See the following links for more information:
// Proposal: https://go.googlesource.com/proposal/+/master/design/draft-iofs.md
// Go PR to watch: https://golang.org/s/draft-iofs-code
// io/fs code as of July 21: https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/

import (
	"os"
	"path/filepath"
)

// FS mirrors the proposed io/fs.FS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/fs.go#21
type FS interface {
	Open(name string) (File, error)
}

// StatFS mirrors the proposed io/fs.StatFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/stat.go#8
type StatFS interface {
	FS
	Stat(name string) (os.FileInfo, error)
}

// ReadDirFS mirrors the proposed io/fs.ReadDirFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/readdir.go#14
type ReadDirFS interface {
	FS
	ReadDir(name string) ([]os.FileInfo, error)
}

// GlobFS mirrors the proposed io/fs.GlobFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/glob.go#13
type GlobFS interface {
	FS
	Glob(pattern string) ([]string, error)
}

// ReadFileFS mirrors the proposed io/fs.ReadFileFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/readfile.go#11
type ReadFileFS interface {
	FS
	ReadFile(name string) ([]byte, error)
}

// File mirrors the proposed io/fs.File
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/fs.go#73
type File interface {
	Stat() (os.FileInfo, error)
	Read(buf []byte) (int, error)
	Close() error
}

// ReadDirFile mirrors the proposed io/fs.ReadDirFile
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/fs.go#83
type ReadDirFile interface {
	File
	ReadDir(n int) ([]os.FileInfo, error)
}

// Walk mirrors the proposed io/fs.Walk
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/walk.go#70
func Walk(fs ReadDirFS, path string, f filepath.WalkFunc) error {
	// TODO This needs to match the published method
	// TODO This is probably wrong, it'll live in io/fs stdlib eventually
	err := walk(fs, path, f)
	if err != nil {
		return err
	}
	return nil
}

func walk(fs ReadDirFS, path string, f filepath.WalkFunc) error {
	infos, err := fs.ReadDir(path)
	if err != nil {
		return err
	}

	for _, info := range infos {
		// The root path signifier needs to be special cased so it's not
		// added to the path passed on to the next methods.
		var fullPath string
		if path == filesystemRootPath {
			fullPath = info.Name()
		} else {
			fullPath = path + pathSeparator + info.Name()
		}

		if info.IsDir() {
			err = walk(fs, fullPath, f)
			if err != nil {
				return err
			}
		}

		err = f(fullPath, info, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
