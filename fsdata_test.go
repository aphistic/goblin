package goblin

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type fsDirSuite struct{}

func (s *fsDirSuite) TestFiles(t sweet.T) {
	d := &fsDir{
		nodes: map[string]fsNode{
			"dir1": &fsDir{
				nodes: map[string]fsNode{
					"dir2": &fsDir{
						nodes: map[string]fsNode{
							"file.txt": &fsFile{
								data: []byte{0x01},
							},
						},
					},
					"file.txt": &fsFile{
						data: []byte{0x02},
					},
				},
			},
			"file.txt": &fsFile{
				data: []byte{0x03},
			},
		},
	}
	Expect(d.Files()).To(Equal(map[string][]byte{
		"dir1/dir2/file.txt": []byte{0x01},
		"dir1/file.txt":      []byte{0x02},
		"file.txt":           []byte{0x03},
	}))
}
