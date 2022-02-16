package repository

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
	"strings"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// splitNonEmpty works exactly like strings.Split(), but only returns non-empty
// components.
func splitNonEmpty(str, sep string) []string {
	components := []string{}
	for _, v := range strings.Split(str, sep) {
		if len(v) > 0 {
			components = append(components, v)
		}
	}
	return components
}

// GenericParent can be used as a common-case implementation of
// HierarchicalRepository.Parent(). It will create a parent URI based on
// IETF RFC3986.
//
// In short, the URI is separated into it's component parts, the path component
// is split along instances of '/', and the trailing element is removed. The
// result is concatenated and parsed as a new URI.
//
// If the URI path is empty or '/', then a nil URI is returned, along with
// ErrURIRoot.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericParent(u gui.URI) (gui.URI, error) {
	p := u.Path()

	if p == "" || p == "/" {
		return nil, ErrURIRoot
	}

	components := splitNonEmpty(p, "/")

	newURI := u.Scheme() + "://" + u.Authority()

	// there will be at least one component, since we know we don't have
	// '/' or ''.
	newURI += "/"
	if len(components) > 1 {
		newURI += strings.Join(components[:len(components)-1], "/")
	}

	// stick the query and fragment back on the end
	if q := u.Query(); len(q) > 0 {
		newURI += "?" + q
	}

	if f := u.Fragment(); len(f) > 0 {
		newURI += "#" + f
	}

	// NOTE: we specifically want to use ParseURI, rather than &uri{},
	// since the repository for the URI we just created might be a
	// CustomURIRepository that implements it's own ParseURI.
	return ParseURI(newURI)
}

// GenericChild can be used as a common-case implementation of
// HierarchicalRepository.Child(). It will create a child URI by separating the
// URI into it's component parts as described in IETF RFC 3986, then appending
// "/" + component to the path, then concatenating the result and parsing it as
// a new URI.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericChild(u gui.URI, component string) (gui.URI, error) {

	// split into components and add the new one
	components := splitNonEmpty(u.Path(), "/")
	components = append(components, component)

	// generate the scheme, authority, and path
	newURI := u.Scheme() + "://" + u.Authority()
	newURI += "/" + strings.Join(components, "/")

	// stick the query and fragment back on the end
	if q := u.Query(); len(q) > 0 {
		newURI += "?" + q
	}
	if f := u.Fragment(); len(f) > 0 {
		newURI += "#" + f
	}

	// NOTE: we specifically want to use ParseURI, rather than &uri{},
	// since the repository for the URI we just created might be a
	// CustomURIRepository that implements it's own ParseURI.
	return ParseURI(newURI)
}

// GenericCopy can be used a common-case implementation of
// CopyableRepository.Copy(). It will perform the copy by obtaining a reader
// for the source URI, a writer for the destination URI, then writing the
// contents of the source to the destination.
//
// For obvious reasons, the destination URI must have a registered
// WritableRepository.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericCopy(source gui.URI, destination gui.URI) error {
	// Look up repositories for the source and destination.
	srcrepo, err := ForURI(source)
	if err != nil {
		return err
	}

	dstrepo, err := ForURI(destination)
	if err != nil {
		return err
	}

	// The destination must be writable.
	destwrepo, ok := dstrepo.(WritableRepository)
	if !ok {
		return ErrOperationNotSupported
	}

	// Create a reader and a writer.
	srcReader, err := srcrepo.Reader(source)
	if err != nil {
		return err
	}
	defer srcReader.Close()

	dstWriter, err := destwrepo.Writer(destination)
	if err != nil {
		return err
	}
	defer dstWriter.Close()

	// Perform the copy.
	_, err = io.Copy(dstWriter, srcReader)
	return err
}

// GenericMove can be used a common-case implementation of
// MovableRepository.Move(). It will perform the move by obtaining a reader
// for the source URI, a writer for the destination URI, then writing the
// contents of the source to the destination. Following this, the source
// will be deleted using WritableRepository.Delete.
//
// For obvious reasons, the source and destination URIs must both be writable.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericMove(source gui.URI, destination gui.URI) error {
	// This looks a lot like GenericCopy(), but I duplicated the code
	// to avoid having to look up the repositories more than once.

	// Look up repositories for the source and destination.
	srcrepo, err := ForURI(source)
	if err != nil {
		return err
	}

	dstrepo, err := ForURI(destination)
	if err != nil {
		return err
	}

	// The source and destination must both be writable, since the source
	// is being deleted, which requires WritableRepository.
	destwrepo, ok := dstrepo.(WritableRepository)
	if !ok {
		return ErrOperationNotSupported
	}

	srcwrepo, ok := srcrepo.(WritableRepository)
	if !ok {
		return ErrOperationNotSupported
	}

	// Create the reader and writer to perform the copy operation.
	srcReader, err := srcrepo.Reader(source)
	if err != nil {
		return err
	}

	dstWriter, err := destwrepo.Writer(destination)
	if err != nil {
		return err
	}
	defer dstWriter.Close()

	// Perform the copy.
	_, err = io.Copy(dstWriter, srcReader)
	if err != nil {
		return err
	}

	// Finally, delete the source only if the move finished without error.
	srcReader.Close()
	return srcwrepo.Delete(source)
}
