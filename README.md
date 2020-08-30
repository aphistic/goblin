Goblin
======

Goblin is a library and utility for interacting with filesystems in Go, including filesystems
embedded in Go binaries. It's designed around the
[proposed](https://go.googlesource.com/proposal/+/master/design/draft-iofs.md) `io/fs`
interfaces. At some point it will likely also support the
[proposed](https://go.googlesource.com/proposal/+/master/design/draft-embed.md) embedded static
assets support.

**Warning:** Since the interfaces this library is based on are not stable, neither is this API!
It will evolve with the proposal and thus should be considered unstable. The goal is to avoid a
lot of API churn but sometimes it's inevitable!

## Installation

If you'd only like to use Goblin as a library you can install it with:

```bash
$ go get github.com/aphistic/goblin
```

If you'd like to make use of Goblin's capability to embed files in a Go binary, you'll want to
install the `goblin` utility to create those files.

```bash
go get github.com/aphistic/goblin/cmd/goblin
```

## General Usage

The representation of a filesystem in Goblin is known as a vault (get it? goblins? vaults? 😁).
The `Vault` interface in Goblin implements at least a few of the `io/fs` interfaces and adds a
few additional utility methods (some hopefully having `io/fs` interfaces in the future).

```go
// Create a new memory vault
mVault := goblin.NewMemoryVault()

// Create a reader to write the file from with
// some file contents.
fileData := bytes.NewReader([]byte("I'm a file!"))

// Write the file contents to your/file/here.txt with
// a file modified time of April 8, 2020 at Midnight UTC.
_ = mVault.WriteFile(
    "your/file/here.txt", fileData,
    goblin.FileModTime(time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)),
)

// Read all the files in the root of the memory vault. In
// our case it's only the "your" directory.
infos, _ := mVault.ReadDir(".")
fmt.Printf("Root Files:\n")
for _, info := range infos {
    fmt.Printf("  - %s\n", info.Name())
}

// Read the contents of the file we previously wrote.
readData, _ := mVault.ReadFile("your/file/here.txt")
fmt.Printf("File Data: %s\n", readData)

// Output:
// Root Files:
//   - your
// File Data: I'm a file!
```

## Mixing Vaults at Runtime

It's sometimes desired to be able to choose between one or more vaults at runtime. Since this
is a common pattern, Goblin provides a way to do it for you with a `VaultSelector`. The
`VaultSelector` is a `Vault`, so it can be used anywhere a `Vault` would be, and provides a few
options for selecting a `Vault` by default but also supports creating custom selectors.

```go
// Use a FilesystemVault if a path is provided as an environment variable,
// but fall back to using a MemoryVault by default. This, for example,
// could be useful during development when you are iterating on embedded
// files and don't want to have to rebuild for every change.
const assetVaultPathKey = "ASSET_VAULT_PATH"

_ = os.Setenv(assetVaultPathKey, "/")

// Create a filesystem vault and a memory vault to use.
fsysAssetVault := goblin.NewFilesystemVault(os.Getenv(assetVaultPathKey))
memAssetVault := goblin.NewMemoryVault() // Loaded from an embedded vault

// Create a vault selector that will use the filesystem asset vault
// if the environment variable ASSET_VAULT_PATH has a non-empty vault,
// use the embedded asset vault if not.
assetVaultSelector := goblin.NewVaultSelector(
    goblin.SelectEnvNotEmpty(assetVaultPathKey, fsysAssetVault),
    goblin.SelectDefault(memAssetVault),
)

// Since ASSET_VAULT_PATH is not empty, we use the filesystem vault.
v, _ := assetVaultSelector.GetVault()
fmt.Printf("Vault: %s\n", v)

// Output:
// Vault: Filesystem Vault (/)
```

## Embedding Files

To embed files in your binary using Goblin, you'll use the `goblin` utility to generate a Go
code file with the contents of the file and a method to load the contents at runtime.

**Note:** The public interface for the `goblin` utility will be changing and is not stabilized.
It'll be changing to hopefully be cleaner and more intuitive.

There are two common command line arguments that are used when create an embedded vault:

* `--name` or `-n`: The name of the vault to embed. This will be used for any package, vault,
  and file names in the generated code.
* `--include-root` or `-r`: The root path of any included files. Any file paths in the vault
  will be relative to this path. Defaults to the current working directory.
* `--include` or `-i`: A [glob](https://golang.org/pkg/path/filepath/#Match) path to include
  files for. Can be provided more than once.

To include all `.html` files in your project's web directory or a directory below it in a
vault, for example, you would use:

```bash
$ goblin create --name assets --include-root /src/web --include *.html --include **/*.html 
```

This would result in a file called `goblin_assets.go` in the current directory to be created
with contents similar to the following:

```go
package assets

import goblin "github.com/aphistic/goblin"

func loadVaultAssets() (goblin.Vault, error) {
	return goblin.LoadMemoryVault(goblinMemoryVaultXassets)
}

var goblinMemoryVaultXassets = []byte{ /* lots of bytes */ }
```

If you need to specify a package name other than the default (`assets` in our example), you can
use the `--package` or `-p` command line option to provide a different one.

## What's With the Name?

When I was trying to come up with a project name for a utility to embed binary files in Go
binaries one of the names I was working with included `gobin` for **Go** **Bin**aries and when
I read it I misread it as `goblin` at one point and thought that would be a fun name since
goblins love treasure (your files!) and keeping things in vaults.