package goblin

import "fmt"

// Refer to https://go.googlesource.com/proposal/+/master/design/draft-iofs.md for FS implementation

const (
	pathSeparator = "/"
)

type Vault interface {
	StatFS
	ReadDirFS
	GlobFS
	ReadFileFS

	fmt.GoStringer
	fmt.Stringer
}
