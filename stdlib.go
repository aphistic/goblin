package goblin

// These are all types or methods that mirror the proposed
// fsys options from the Go standard library.
//
// See the following links for more information:
// Proposal: https://go.googlesource.com/proposal/+/master/design/draft-iofs.md
// Go PR to watch: https://golang.org/s/draft-iofs-code
// io/fs code as of July 21: https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/
//
// Everything in this file is subject to the Go license and is not covered by the
// Goblin license.
//
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// LICENSE:
// Copyright (c) 2009 The Go Authors. All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"errors"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
)

// FS mirrors the proposed io/fs.FS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/fs.go#21
type FS interface {
	Open(name string) (File, error)
}

// StatFS mirrors the proposed io/fs.StatFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/stat.go#8
type StatFS interface {
	FS
	Stat(name string) (os.FileInfo, error)
}

// ReadDirFS mirrors the proposed io/fs.ReadDirFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/readdir.go#14
type ReadDirFS interface {
	FS
	ReadDir(name string) ([]os.FileInfo, error)
}

// ReadDir mirrors the proposed io/fs.ReadDir
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/readdir.go#28
func ReadDir(fsys FS, name string) ([]os.FileInfo, error) {
	if fsys, ok := fsys.(ReadDirFS); ok {
		return fsys.ReadDir(name)
	}
	file, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// TODO: Do we really need the Stat?
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("TODO")
	}
	dir, ok := file.(ReadDirFile)
	if !ok {
		return nil, errors.New("TODO")
	}
	list, err := dir.ReadDir(-1)
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, err
}

// GlobFS mirrors the proposed io/fs.GlobFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/glob.go#13
type GlobFS interface {
	FS
	Glob(pattern string) ([]string, error)
}

// ReadFileFS mirrors the proposed io/fs.ReadFileFS
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/readfile.go#11
type ReadFileFS interface {
	FS
	ReadFile(name string) ([]byte, error)
}

// File mirrors the proposed io/fs.File
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/fs.go#73
type File interface {
	Stat() (os.FileInfo, error)
	Read(buf []byte) (int, error)
	Close() error
}

// Stat mirrors the proposed io/fs.Stat
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/stat.go#20
func Stat(fsys FS, name string) (os.FileInfo, error) {
	if fsys, ok := fsys.(StatFS); ok {
		return fsys.Stat(name)
	}
	file, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file.Stat()
}

// ReadDirFile mirrors the proposed io/fs.ReadDirFile
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/fs.go#83
type ReadDirFile interface {
	File
	ReadDir(n int) ([]os.FileInfo, error)
}

// Walk mirrors the proposed io/fs.Walk
//
// https://go.googlesource.com/go/+/2ad964dc3731dac3ab561ab344042dbe316dbf28/src/io/fs/walk.go#70
func Walk(fsys FS, root string, walkFn filepath.WalkFunc) error {
	info, err := Stat(fsys, root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walk(fsys, root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}
func walk(fsys FS, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}
	infos, err := ReadDir(fsys, path)
	err1 := walkFn(path, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}
	for _, info := range infos {
		filename := pathpkg.Join(path, info.Name())
		err = walk(fsys, filename, info, walkFn)
		if err != nil {
			if !info.IsDir() || err != filepath.SkipDir {
				return err
			}
		}
	}
	return nil
}
