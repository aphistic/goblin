package goblin

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type vaultSuite struct{}

func (s *vaultSuite) TestNewVault(t sweet.T) {
	v := newVault("foo")
	Expect(v).To(Equal(&vault{
		name: "foo",
		root: newFsDir(),
	}))
}

func (s *vaultSuite) TestSetFile(t sweet.T) {
	v := newVault("foo")

	err := v.SetFile("this/is/a/path", []byte{0x01, 0x02, 0x03, 0x04})
	Expect(err).To(BeNil())
	Expect(v).To(Equal(&vault{
		name: "foo",
		root: &fsDir{
			nodes: map[string]fsNode{
				"this": &fsDir{
					nodes: map[string]fsNode{
						"is": &fsDir{
							nodes: map[string]fsNode{
								"a": &fsDir{
									nodes: map[string]fsNode{
										"path": &fsFile{
											data: []byte{0x01, 0x02, 0x03, 0x04},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}))
	err = v.SetFile("this/is/a/file.txt", []byte{0x04, 0x03, 0x02, 0x01})
	Expect(err).To(BeNil())
	Expect(v).To(Equal(&vault{
		name: "foo",
		root: &fsDir{
			nodes: map[string]fsNode{
				"this": &fsDir{
					nodes: map[string]fsNode{
						"is": &fsDir{
							nodes: map[string]fsNode{
								"a": &fsDir{
									nodes: map[string]fsNode{
										"path": &fsFile{
											data: []byte{0x01, 0x02, 0x03, 0x04},
										},
										"file.txt": &fsFile{
											data: []byte{0x04, 0x03, 0x02, 0x01},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}))
}

func (s *vaultSuite) TestSetFileOverwrite(t sweet.T) {
	v := newVault("foo")

	err := v.SetFile("this/is/a/path", []byte{0x01, 0x02, 0x03, 0x04})
	Expect(err).To(BeNil())
	Expect(v).To(Equal(&vault{
		name: "foo",
		root: &fsDir{
			nodes: map[string]fsNode{
				"this": &fsDir{
					nodes: map[string]fsNode{
						"is": &fsDir{
							nodes: map[string]fsNode{
								"a": &fsDir{
									nodes: map[string]fsNode{
										"path": &fsFile{
											data: []byte{0x01, 0x02, 0x03, 0x04},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}))
	err = v.SetFile("this/is/a/path", []byte{0x04, 0x03, 0x02, 0x01})
	Expect(err).To(BeNil())
	Expect(v).To(Equal(&vault{
		name: "foo",
		root: &fsDir{
			nodes: map[string]fsNode{
				"this": &fsDir{
					nodes: map[string]fsNode{
						"is": &fsDir{
							nodes: map[string]fsNode{
								"a": &fsDir{
									nodes: map[string]fsNode{
										"path": &fsFile{
											data: []byte{0x04, 0x03, 0x02, 0x01},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}))
}

func (s *vaultSuite) TestSetFileInvalidDirectory(t sweet.T) {
	v := newVault("foo")

	err := v.SetFile("dir/file", []byte{0x01, 0x02, 0x03, 0x04})
	Expect(err).To(BeNil())
	Expect(v).To(Equal(&vault{
		name: "foo",
		root: &fsDir{
			nodes: map[string]fsNode{
				"dir": &fsDir{
					nodes: map[string]fsNode{
						"file": &fsFile{
							data: []byte{0x01, 0x02, 0x03, 0x04},
						},
					},
				},
			},
		},
	}))
	err = v.SetFile("dir/file/otherfile", []byte{0x04, 0x03, 0x02, 0x01})
	Expect(err).To(Equal(fmt.Errorf("file is not a directory")))
	Expect(v).To(Equal(&vault{
		name: "foo",
		root: &fsDir{
			nodes: map[string]fsNode{
				"dir": &fsDir{
					nodes: map[string]fsNode{
						"file": &fsFile{
							data: []byte{0x01, 0x02, 0x03, 0x04},
						},
					},
				},
			},
		},
	}))
}

func (s *vaultSuite) TestFiles(t sweet.T) {
	v := newVault("foo")

	err := v.SetFile("this/is/a/file", []byte{0x01})
	Expect(err).To(BeNil())
	err = v.SetFile("this/is/another/file", []byte{0x02})
	Expect(err).To(BeNil())
	err = v.SetFile("file.txt", []byte{0x03})
	Expect(err).To(BeNil())

	Expect(v.Files()).To(Equal(map[string][]byte{
		"this/is/a/file":       []byte{0x01},
		"this/is/another/file": []byte{0x02},
		"file.txt":             []byte{0x03},
	}))
}
