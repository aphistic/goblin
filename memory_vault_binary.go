package goblin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
)

// MarshalBinary encodes the MemoryVault into a binary representation.
func (v *MemoryVault) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	var paths []string
	err := Walk(v, ".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		tokens, err := splitPath(path)
		if err != nil {
			return nil, err
		}

		node, err := v.root.GetNode(tokens)
		if err != nil {
			return nil, err
		}
		fInfo, err := node.Stat()
		if err != nil {
			return nil, err
		}

		err = tw.WriteHeader(&tar.Header{
			Name:    node.FullPath(),
			ModTime: fInfo.ModTime(),
			Size:    fInfo.Size(),
		})
		if err != nil {
			return nil, err
		}

		data, err := v.ReadFile(path)
		if err != nil {
			return nil, err
		}

		totalWritten := 0
		for {
			n, err := tw.Write(data[totalWritten:])
			if err != nil {
				return nil, err
			}
			totalWritten = totalWritten + n

			if totalWritten >= len(data) {
				break
			}
		}
	}

	err = tw.Close()
	if err != nil {
		return nil, err
	}

	err = gw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary decodes the provided data into the MemoryVault.
func (v *MemoryVault) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		b := bytes.NewBuffer(nil)
		buf := make([]byte, 4096)
		for {
			n, err := tr.Read(buf)
			if err == io.EOF {
				// Ignore an EOF at this point because we might have
				// read something at the end of the file
			} else if err != nil {
				return err
			}

			curWrite := 0
			for {
				wN, wErr := b.Write(buf[curWrite:n])
				if wErr != nil {
					return wErr
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

		err = v.WriteFile(
			header.Name,
			b,
			FileModTime(header.ModTime),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
