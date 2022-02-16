package test

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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/storage"
)

type file struct {
	*os.File
	path string
}

type directory struct {
	gui.URI
}

// Declare conformity to the ListableURI interface
var _ gui.ListableURI = (*directory)(nil)

func (f *file) Open() (io.ReadCloser, error) {
	return os.Open(f.path)
}

func (f *file) Save() (io.WriteCloser, error) {
	return os.Open(f.path)
}

func (f *file) ReadOnly() bool {
	return true
}

func (f *file) Name() string {
	return filepath.Base(f.path)
}

func (f *file) URI() gui.URI {
	return storage.NewURI("file://" + f.path)
}

func openFile(uri gui.URI, create bool) (*file, error) {
	if uri.Scheme() != "file" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := uri.String()[7:]
	f, err := os.Open(path)
	if err != nil && create {
		f, err = os.Create(path)
	}
	return &file{File: f, path: path}, err
}

func (d *testDriver) FileReaderForURI(uri gui.URI) (gui.URIReadCloser, error) {
	return openFile(uri, false)
}

func (d *testDriver) FileWriterForURI(uri gui.URI) (gui.URIWriteCloser, error) {
	return openFile(uri, true)
}

func (d *testDriver) ListerForURI(uri gui.URI) (gui.ListableURI, error) {
	if uri.Scheme() != "file" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := uri.String()[len(uri.Scheme())+3 : len(uri.String())]
	s, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !s.IsDir() {
		return nil, fmt.Errorf("path '%s' is not a directory, cannot convert to listable URI", path)
	}

	return &directory{URI: uri}, nil
}

func (d *directory) List() ([]gui.URI, error) {
	if d.Scheme() != "file" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := d.String()[len(d.Scheme())+3 : len(d.String())]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	urilist := []gui.URI{}

	for _, f := range files {
		uri := storage.NewURI("file://" + filepath.Join(path, f.Name()))
		urilist = append(urilist, uri)
	}

	return urilist, nil
}
