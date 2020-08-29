package goblin

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aphistic/goblin/internal/logging"
	"github.com/dave/jennifer/jen"
	"github.com/dustin/go-humanize"
)

const (
	goblinImport = "github.com/aphistic/goblin"
)

// MemoryBuilderOption is an option used when creating a memory vault builder
type MemoryBuilderOption func(b *MemoryBuilder)

// MemoryBuilderExportLoader will cause the code for loading the vault to be
// exported. For example, LoadVaultAssets instead of loadVaultAssets.
func MemoryBuilderExportLoader(exportLoader bool) MemoryBuilderOption {
	return func(b *MemoryBuilder) {
		b.exportLoader = exportLoader
	}
}

// MemoryBuilderLogger provides a logger for the builder to use.
func MemoryBuilderLogger(logger logging.Logger) MemoryBuilderOption {
	return func(b *MemoryBuilder) {
		b.logger = logger
	}
}

// MemoryBuilder creates binary or code representations of a memory vault.
type MemoryBuilder struct {
	logger       logging.Logger
	exportLoader bool

	v *MemoryVault
}

// NewMemoryBuilder creates a new memory builder.
func NewMemoryBuilder(opts ...MemoryBuilderOption) *MemoryBuilder {
	b := &MemoryBuilder{
		logger: logging.NewNilLogger(),
		v:      NewMemoryVault(),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Include iterates over all files in the root path, then includes any file matching
// one or more of the provided globs in the memory vault being built.
func (b *MemoryBuilder) Include(rootPath string, globs []string) error {
	for _, glob := range globs {
		if strings.Contains(glob, "..") {
			return fmt.Errorf(".. cannot be used in include paths")
		}

		fullPathGlob := filepath.Join(rootPath, glob)
		globDir := filepath.Dir(fullPathGlob)
		if !strings.HasSuffix(globDir, string(os.PathSeparator)) {
			globDir = globDir + string(os.PathSeparator)
		}

		matches, err := filepath.Glob(fullPathGlob)
		if err != nil {
			b.logger.Printf("error with path '%s': %s\n", fullPathGlob, err)
			return err
		}

		for _, match := range matches {
			fInfo, err := os.Stat(match)
			if err != nil {
				return err
			}

			// TODO Test how this is achieved
			filePath := strings.TrimPrefix(match, rootPath)
			filePath = strings.TrimPrefix(filePath, globDir)
			filePath = strings.TrimPrefix(filePath, pathSeparator)

			b.logger.Printf("Adding: %s... ", filePath)
			data, err := ioutil.ReadFile(match)
			if err != nil {
				return err
			}

			err = b.v.WriteFile(
				filePath,
				bytes.NewBuffer(data),
				FileModTime(fInfo.ModTime()),
			)
			if err != nil {
				return err
			}
			b.logger.Printf("%s\n", humanize.Bytes(uint64(len(data))))
		}
	}

	return nil
}

// WriteBinary writes the binary representation of the memory vault to the provided io.Writer.
func (b *MemoryBuilder) WriteBinary(w io.Writer) error {
	vaultData, err := b.v.MarshalBinary()
	if err != nil {
		return err
	}

	curLoc := 0
	_, err = w.Write(vaultData[curLoc:])
	if err != nil {
		return err
	}

	return nil
}

// WriteLoader writes code and binary data to the provided io.Writer to allow loading the memory
// vault being built at runtime.
func (b *MemoryBuilder) WriteLoader(packageName string, vaultName string, w io.Writer) error {
	vaultData, err := b.v.MarshalBinary()
	if err != nil {
		return err
	}

	genFile := jen.NewFile(packageName)
	fullVaultName := makeMemoryVaultName(vaultName)

	loadPrefix := "loadVault"
	if b.exportLoader {
		loadPrefix = "LoadVault"
	}

	genFile.Func().Id(loadPrefix+strings.Title(vaultName)).
		Params().
		Params(jen.Qual(goblinImport, "Vault"), jen.Id("error")).
		Block(
			jen.Return(
				jen.Qual(goblinImport, "LoadMemoryVault").Params(jen.Id(fullVaultName)),
			),
		)

	var vaultValues []jen.Code
	for _, b := range vaultData {
		vaultValues = append(vaultValues, jen.LitByte(b))
	}
	genFile.Var().Id(fullVaultName).Op("=").Index().Byte().Values(vaultValues...)

	err = genFile.Render(w)
	if err != nil {
		return err
	}

	return nil
}
