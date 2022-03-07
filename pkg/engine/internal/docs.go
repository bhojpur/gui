package internal

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
	"errors"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/storage"
)

var errNoAppID = errors.New("storage API requires a unique ID, use app.NewWithID()")

// Docs is an internal implementation of the document features of the Storage interface.
// It is based on top of the current `file` repository and is rooted at RootDocURI.
type Docs struct {
	RootDocURI gui.URI
}

// Create will create a new document ready for writing, you must write something and close the returned writer
// for the create process to complete.
// If the document for this app with that name already exists an error will be returned.
func (d *Docs) Create(name string) (gui.URIWriteCloser, error) {
	if d.RootDocURI == nil {
		return nil, errNoAppID
	}

	err := d.ensureRootExists()
	if err != nil {
		return nil, err
	}

	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return nil, err
	}

	exists, err := storage.Exists(u)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("document with name " + name + " already exists")
	}

	return storage.Writer(u)
}

// List returns all documents that have been saved by the current application.
// Remember to use `app.NewWithID` so that your storage is unique.
func (d *Docs) List() []string {
	if d.RootDocURI == nil {
		return nil
	}

	var ret []string
	uris, err := storage.List(d.RootDocURI)
	if err != nil {
		return ret
	}

	for _, u := range uris {
		ret = append(ret, u.Name())
	}

	return ret
}

// Open will grant access to the contents of the named file. If an error occurs it is returned instead.
func (d *Docs) Open(name string) (gui.URIReadCloser, error) {
	if d.RootDocURI == nil {
		return nil, errNoAppID
	}

	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return nil, err
	}

	return storage.Reader(u)
}

// Remove will delete the document with the specified name, if it exists
func (d *Docs) Remove(name string) error {
	if d.RootDocURI == nil {
		return errNoAppID
	}

	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return err
	}

	return storage.Delete(u)
}

// Save will open a document ready for writing, you close the returned writer for the save to complete.
// If the document for this app with that name does not exist an error will be returned.
func (d *Docs) Save(name string) (gui.URIWriteCloser, error) {
	if d.RootDocURI == nil {
		return nil, errNoAppID
	}

	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return nil, err
	}

	exists, err := storage.Exists(u)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("document with name " + name + " does not exist")
	}

	return storage.Writer(u)
}

func (d *Docs) ensureRootExists() error {
	exists, err := storage.Exists(d.RootDocURI)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return storage.CreateListable(d.RootDocURI)
}
