package goblin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/dave/jennifer/jen"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const goblinImport = "github.com/aphistic/goblin"

type Builder struct {
	vaultName string

	v *vault
}

func NewBuilder(vaultName string) *Builder {
	return &Builder{
		vaultName: vaultName,
		v:         newVault(vaultName),
	}
}

func (b *Builder) Include(globs []string) error {
	for _, glob := range globs {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return err
		}

		for _, match := range matches {
			data, err := ioutil.ReadFile(match)
			if err != nil {
				return err
			}

			err = b.v.SetFile(match, data)
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
