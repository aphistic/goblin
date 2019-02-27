package goblin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/aphistic/goblin/internal/logging"
	"github.com/dave/jennifer/jen"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const goblinImport = "github.com/aphistic/goblin"

type BuilderOption func(b *Builder)

func BuilderLogger(logger logging.Logger) BuilderOption{
	return func(b *Builder) {
		b.logger = logger
	}
}

type Builder struct {
	logger logging.Logger

	vaultName string

	v *vault
}

func NewBuilder(vaultName string, opts ...BuilderOption) *Builder {
	b := &Builder{
		logger: logging.NewNilLogger(),
		vaultName: vaultName,
		v:         newVault(vaultName),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Builder) Include(globs []string) error {
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
			filePath := strings.TrimPrefix(match, globDir)

			b.logger.Printf("match: %s\n", filePath)
			data, err := ioutil.ReadFile(match)
			if err != nil {
				return err
			}

			err = b.v.SetFile(filePath, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Builder) Write(packageName string, w io.Writer) error {
	dataBuf := bytes.NewBuffer(nil)
	err := b.writeData(dataBuf)
	if err != nil {
		return err
	}

	genFile := jen.NewFile(packageName)
	vaultName := makeVaultName(b.vaultName)

	genFile.Func().Id("loadVault" + strings.Title(b.vaultName)).
		Params().
		Params(jen.Qual(goblinImport, "Vault"), jen.Id("error")).
		Block(
			jen.Return(
				jen.Qual(goblinImport, "LoadVault").Params(jen.Id(vaultName)),
			),
		)

	var vaultValues []jen.Code
	for _, b := range dataBuf.Bytes() {
		vaultValues = append(vaultValues, jen.LitByte(b))
	}
	genFile.Var().Id(vaultName).Op("=").Index().Byte().Values(vaultValues...)

	err = genFile.Render(w)
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) writeData(w io.Writer) error {
	gw := gzip.NewWriter(w)
	tw := tar.NewWriter(gw)

	for name, data := range b.v.Files() {
		err := tw.WriteHeader(&tar.Header{
			Name: name,
			Size: int64(len(data)),
		})
		if err != nil {
			return err
		}

		totalWritten := 0
		for {
			n, err := tw.Write(data[totalWritten:])
			if err != nil {
				return err
			}
			totalWritten = totalWritten + n

			if totalWritten >= len(data) {
				break
			}
		}
	}

	err := tw.Close()
	if err != nil {
		return err
	}

	err = gw.Close()
	if err != nil {
		return err
	}

	return nil
}
