// It provides storage access and management functionality.
package storage

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
)

// OpenFileFromURI loads a file read stream from a resource identifier.
// This is mostly provided so that file references can be saved using their URI and loaded again later.
//
// Deprecated: this has been replaced by storage.Reader(URI)
func OpenFileFromURI(uri gui.URI) (gui.URIReadCloser, error) {
	return Reader(uri)
}

// SaveFileToURI loads a file write stream to a resource identifier.
// This is mostly provided so that file references can be saved using their URI and written to again later.
//
// Deprecated: this has been replaced by storage.Writer(URI)
func SaveFileToURI(uri gui.URI) (gui.URIWriteCloser, error) {
	return Writer(uri)
}

// ListerForURI will attempt to use the application's driver to convert a
// standard URI into a listable URI.
//
// Since: 1.4
func ListerForURI(uri gui.URI) (gui.ListableURI, error) {
	listable, err := CanList(uri)
	if err != nil {
		return nil, err
	}
	if !listable {
		return nil, errors.New("uri is not listable")
	}

	return &legacyListable{uri}, nil
}

type legacyListable struct {
	gui.URI
}

func (l *legacyListable) List() ([]gui.URI, error) {
	return List(l.URI)
}
