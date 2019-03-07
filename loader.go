package goblin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
)

func LoadVault(vaultData []byte) (Vault, error) {
	r := bytes.NewReader(vaultData)

	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gr)

	v := newVault("")
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		b := bytes.NewBuffer(nil)
		buf := make([]byte, 4096)
		for {
			n, err := tr.Read(buf)
			if err == io.EOF {
				// Ignore an EOF at this point because we might have
				// read something at the end of the file
			} else if err != nil {
				return nil, err
			}

			curWrite := 0
			for {
				wN, wErr := b.Write(buf[curWrite:n])
				if wErr != nil {
					return nil, wErr
				}
				curWrite += wN

				if curWrite >= n {
					break
				}
			}

			if err == io.EOF {
				break
			}
		}

		v.SetFile(header.Name, b.Bytes())
	}

	return v, nil
}
