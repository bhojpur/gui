package util

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"io"
	"os"
	"path/filepath"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// Exists will return true if the passed path exists on the current system.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// CopyFile copies the content of a regular file, source, into target path.
func CopyFile(source, target string) error {
	return copyFileMode(source, target, 0644)
}

// CopyExeFile copies the content of an executable file, source, into target path.
func CopyExeFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0755)
}

// EnsureSubDir will make sure a named directory exists within the parent - creating it if not.
func EnsureSubDir(parent, name string) string {
	path := filepath.Join(parent, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			gui.LogError("Failed to create dirrectory", err)
		}
	}
	return path
}

func copyFileMode(src, tgt string, perm os.FileMode) (err error) {
	if _, err := os.Stat(src); err != nil {
		return err
	}
	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer source.Close()
	target, err := os.OpenFile(tgt, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer func() {
		if r := target.Close(); r != nil && err == nil {
			err = r
		}
	}()
	_, err = io.Copy(target, source)
	return err
}
