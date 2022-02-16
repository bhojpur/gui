package mobile

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
	gui "github.com/bhojpur/gui/pkg/engine"

	"github.com/bhojpur/gui/pkg/engine/storage/repository"
)

// declare conformance with repository types
var _ repository.Repository = (*mobileFileRepo)(nil)
var _ repository.HierarchicalRepository = (*mobileFileRepo)(nil)
var _ repository.ListableRepository = (*mobileFileRepo)(nil)
var _ repository.WritableRepository = (*mobileFileRepo)(nil)

type mobileFileRepo struct {
}

func (m *mobileFileRepo) CanList(u gui.URI) (bool, error) {
	return canListURI(u), nil
}

func (m *mobileFileRepo) CanRead(u gui.URI) (bool, error) {
	return true, nil // TODO check a file can be read
}

func (m *mobileFileRepo) CanWrite(u gui.URI) (bool, error) {
	return true, nil // TODO check a file can be written
}

func (m *mobileFileRepo) Child(u gui.URI, name string) (gui.URI, error) {
	if u == nil || u.Scheme() != "file" {
		return nil, repository.ErrOperationNotSupported
	}

	return repository.GenericChild(u, name)
}

func (m *mobileFileRepo) CreateListable(u gui.URI) error {
	return createListableURI(u)
}

func (m *mobileFileRepo) Delete(u gui.URI) error {
	// TODO: implement this
	return repository.ErrOperationNotSupported
}

func (m *mobileFileRepo) Destroy(string) {
}

func (m *mobileFileRepo) Exists(u gui.URI) (bool, error) {
	return existsURI(u)
}

func (m *mobileFileRepo) List(u gui.URI) ([]gui.URI, error) {
	return listURI(u)
}

func (m *mobileFileRepo) Parent(u gui.URI) (gui.URI, error) {
	if u == nil || u.Scheme() != "file" {
		return nil, repository.ErrOperationNotSupported
	}

	return repository.GenericParent(u)
}

func (m *mobileFileRepo) Reader(u gui.URI) (gui.URIReadCloser, error) {
	return fileReaderForURI(u)
}

func (m *mobileFileRepo) Writer(u gui.URI) (gui.URIWriteCloser, error) {
	return fileWriterForURI(u)
}
