package engine

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
)

// URIReadCloser represents a cross platform data stream from a file or provider of data.
// It may refer to an item on a filesystem or data in another application that we have access to.
type URIReadCloser interface {
	io.ReadCloser

	URI() URI
}

// URIWriteCloser represents a cross platform data writer for a file resource.
// This will normally refer to a local file resource.
type URIWriteCloser interface {
	io.WriteCloser

	URI() URI
}

// URI represents the identifier of a resource on a target system.  This
// resource may be a file or another data source such as an app or file sharing
// system.
//
// In general, it is expected that URI implementations follow IETF RFC3896.
// Implementations are highly recommended to utilize net/url to implement URI
// parsing methods, especially Scheme(), AUthority(), Path(), Query(), and
// Fragment().
type URI interface {
	fmt.Stringer

	// Extension should return the file extension of the resource
	// referenced by the URI. For example, the Extension() of
	// 'file://foo/bar.baz' is 'baz'. May return an empty string if the
	// referenced resource has none.
	Extension() string

	// Name should return the base name of the item referenced by the URI.
	// For example, the Name() of 'file://foo/bar.baz' is 'bar.baz'.
	Name() string

	// MimeType should return the content type of the resource referenced
	// by the URI. The returned string should be in the format described
	// by Section 5 of RFC2045 ("Content-Type Header Field").
	MimeType() string

	// Scheme should return the URI scheme of the URI as defined by IETF
	// RFC3986. For example, the Scheme() of 'file://foo/bar.baz` is
	// 'file'.
	//
	// Scheme should always return the scheme in all lower-case characters.
	Scheme() string

	// Authority should return the URI authority, as defined by IETF
	// RFC3986.
	//
	// NOTE: the RFC3986 can be obtained by combining the User and Host
	// Fields of net/url's URL structure. Consult IETF RFC3986, section
	// 3.2, pp. 17.
	//
	// Since: 2.0
	Authority() string

	// Path should return the URI path, as defined by IETF RFC3986.
	//
	// Since: 2.0
	Path() string

	// Query should return the URI query, as defined by IETF RFC3986.
	//
	// Since: 2.0
	Query() string

	// Fragment should return the URI fragment, as defined by IETF
	// RFC3986.
	//
	// Since: 2.0
	Fragment() string
}

// ListableURI represents a URI that can have child items, most commonly a
// directory on disk in the native filesystem.
//
// Since: 1.4
type ListableURI interface {
	URI

	// List returns a list of child URIs of this URI.
	List() ([]URI, error)
}
