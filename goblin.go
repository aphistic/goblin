// Package goblin provides interaction with various types of filesystems, including those
// embedded in the binary.
//
// It's designed to implement the proposed Go standard library filesystem interfaces. Since
// these interfaces do not exist in a stable release yet, they've been included in this package
// for now. Eventually those interfaces will be removed in favor of the stable release. Types
// that are included temporarily are noted in their documentation.
package goblin

import (
	"fmt"
	"time"
)

// Refer to https://go.googlesource.com/proposal/+/master/design/draft-iofs.md for FS implementation

const (
	pathSeparator = "/"
)

// Vault is the interface that provides all the interfaces a Goblin vault must implement.
type Vault interface {
	StatFS
	ReadDirFS
	GlobFS
	ReadFileFS

	fmt.Stringer
}

// FileOption is a common set of options used when creating or
// managing files.
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

// FileModTime specifies the modified time to use for the file.
func FileModTime(modTime time.Time) FileOption {
	return func(fOpts *fileOptions) {
		fOpts.ModTime = modTime
	}
}
