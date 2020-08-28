package goblin

import (
	"bytes"
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

type MemoryBuilderOption func(b *MemoryBuilder)

func MemoryBuilderExportLoader(exportLoader bool) MemoryBuilderOption {
	return func(b *MemoryBuilder) {
		b.exportLoader = exportLoader
	}
}

func MemoryBuilderLogger(logger logging.Logger) MemoryBuilderOption {
	return func(b *MemoryBuilder) {
		b.logger = logger
	}
}

type MemoryBuilder struct {
	logger       logging.Logger
	exportLoader bool

	v *MemoryVault
}

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

func (b *MemoryBuilder) Include(globs []string) error {
	for _, glob := range globs {
		globDir := filepath.Dir(glob)
		if !strings.HasSuffix(globDir, string(os.PathSeparator)) {
			globDir = globDir + string(os.PathSeparator)
		}

		matches, err := filepath.Glob(glob)
		if err != nil {
			b.logger.Printf("error with path '%s': %s\n", glob, err)
			return err
		}

		for _, match := range matches {
			fInfo, err := os.Stat(match)
			if err != nil {
				return err
			}
			filePath := strings.TrimPrefix(match, globDir)

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

func (b *MemoryBuilder) WriteBinary(packageName string, vaultName string, w io.Writer) error {
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

func (b *MemoryBuilder) WriteCode(packageName string, vaultName string, w io.Writer) error {
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
