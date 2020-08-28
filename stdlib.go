package goblin

import (
	"os"
	"path/filepath"
)

type FS interface {
	Open(name string) (File, error)
}

type StatFS interface {
	FS
	Stat(name string) (os.FileInfo, error)
}

type ReadDirFS interface {
	FS
	ReadDir(name string) ([]os.FileInfo, error)
}

type GlobFS interface {
	FS
	Glob(pattern string) ([]string, error)
}

type ReadFileFS interface {
	FS
	ReadFile(name string) ([]byte, error)
}

type File interface {
	Stat() (os.FileInfo, error)
	Read(buf []byte) (int, error)
	Close() error
}

type ReadDirFile interface {
	File
	ReadDir(n int) ([]os.FileInfo, error)
}

// TODO This is probably wrong, it'll live in io/fs stdlib eventually
func Walk(fs ReadDirFS, path string, f filepath.WalkFunc) error {
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
