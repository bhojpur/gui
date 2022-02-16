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
	"fmt"
	"strings"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// repositoryTable stores the mapping of schemes to Repository implementations.
// It should only ever be used by ForURI() and Register().
var repositoryTable map[string]Repository = map[string]Repository{}

// Repository represents a storage repository, which is a set of methods which
// implement specific functions on a URI. Repositories are registered to handle
// specific URI schemes, and the higher-level functions that operate on URIs
// internally look up an appropriate method from the relevant Repository.
//
// The repository interface includes only methods which must be implemented at
// a minimum. Without implementing all of the methods in this interface, a URI
// would not be usable in a useful way. Additional functionality can be exposed
// by using interfaces which extend Repository.
//
// Repositories are registered to handle a specific URI scheme (or schemes)
// using the Register() method. When a higher-level URI function such as
// storage.Copy() is called, the storage package will internally look up
// the repository associated with the scheme of the URI, then it will use
// a type assertion to check if the repository implements CopyableRepository.
// If so, the Copy() function will be run from the repository, otherwise
// storage.Copy() will return NotSupportedError. This works similarly for
// all other methods in repository-related interfaces.
//
// Note that a repository can be registered for multiple URI schemes. In such
// cases, the repository must internally select and implement the correct
// behavior for each URI scheme.
//
// A repository will only ever need to handle URIs with schemes for which it
// was registered, with the exception that functions with more than 1 operand
// such as Copy() and Move(), in which cases only the first operand is
// guaranteed to match a scheme for which the repository is registered.
//
// NOTE: most developers who use Bhojpur GUI should *not* generally attempt to
// call repository methods directly. You should use the methods in the storage
// package, which will automatically detect the scheme of a URI and call into
// the appropriate repository.
//
// Since: 2.0
type Repository interface {

	// Exists will be used to implement calls to storage.Exists() for the
	// registered scheme of this repository.
	//
	// Since: 2.0
	Exists(u gui.URI) (bool, error)

	// Reader will be used to implement calls to storage.Reader()
	// for the registered scheme of this repository.
	//
	// Since: 2.0
	Reader(u gui.URI) (gui.URIReadCloser, error)

	// CanRead will be used to implement calls to storage.CanRead() for the
	// registered scheme of this repository.
	//
	// Since: 2.0
	CanRead(u gui.URI) (bool, error)

	// Destroy is called when the repository is un-registered from a given
	// URI scheme.
	//
	// The string parameter will be the URI scheme that the repository was
	// registered for. This may be useful for repositories that need to
	// handle more than one URI scheme internally.
	//
	// Since: 2.0
	Destroy(string)
}

// CustomURIRepository is an extension of the repository interface which
// allows the behavior of storage.ParseURI to be overridden. This is only
// needed if you wish to generate custom URI types, rather than using
// Bhojpur GUI's URI implementation and net/url based parsing.
//
// NOTE: even for URIs with non-RFC3986-compliant encoding, the URI MUST begin
// with 'scheme:', or storage.ParseURI() will not be able to determine which
// storage repository to delegate to for parsing.
//
// Since: 2.0
type CustomURIRepository interface {
	Repository

	// ParseURI will be used to implement calls to storage.ParseURI()
	// for the registered scheme of this repository.
	ParseURI(string) (gui.URI, error)
}

// WritableRepository is an extension of the Repository interface which also
// supports obtaining a writer for URIs of the scheme it is registered to.
//
// Since: 2.0
type WritableRepository interface {
	Repository

	// Writer will be used to implement calls to storage.WriterTo() for
	// the registered scheme of this repository.
	//
	// Since: 2.0
	Writer(u gui.URI) (gui.URIWriteCloser, error)

	// CanWrite will be used to implement calls to storage.CanWrite() for
	// the registered scheme of this repository.
	//
	// Since: 2.0
	CanWrite(u gui.URI) (bool, error)

	// Delete will be used to implement calls to storage.Delete() for the
	// registered scheme of this repository.
	//
	// Since: 2.0
	Delete(u gui.URI) error
}

// ListableRepository is an extension of the Repository interface which also
// supports obtaining directory listings (generally analogous to a directory
// listing) for URIs of the scheme it is registered to.
//
// Since: 2.0
type ListableRepository interface {
	Repository

	// CanList will be used to implement calls to storage.Listable() for
	// the registered scheme of this repository.
	//
	// Since: 2.0
	CanList(u gui.URI) (bool, error)

	// List will be used to implement calls to storage.List() for the
	// registered scheme of this repository.
	//
	// Since: 2.0
	List(u gui.URI) ([]gui.URI, error)

	// CreateListable will be used to implement calls to
	// storage.CreateListable() for the registered scheme of this
	// repository.
	//
	// Since: 2.0
	CreateListable(u gui.URI) error
}

// HierarchicalRepository is an extension of the Repository interface which
// also supports determining the parent and child items of a URI.
//
// Since: 2.0
type HierarchicalRepository interface {
	Repository

	// Parent will be used to implement calls to storage.Parent() for the
	// registered scheme of this repository.
	//
	// A generic implementation is provided in GenericParent(), which
	// is based on the RFC3986 definition of a URI parent.
	//
	// Since: 2.0
	Parent(gui.URI) (gui.URI, error)

	// Child will be used to implement calls to storage.Child() for
	// the registered scheme of this repository.
	//
	// A generic implementation is provided in GenericParent(), which
	// is based on RFC3986.
	//
	// Since: 2.0
	Child(gui.URI, string) (gui.URI, error)
}

// CopyableRepository is an extension of the Repository interface which also
// supports copying referenced resources from one URI to another.
//
// Since: 2.0
type CopyableRepository interface {
	Repository

	// Copy will be used to implement calls to storage.Copy() for the
	// registered scheme of this repository.
	//
	// A generic implementation is provided by GenericCopy().
	//
	// NOTE: the first parameter is the source, the second is the
	// destination.
	//
	// NOTE: if storage.Copy() is given two URIs of different schemes, it
	// is possible that only the source URI will be of the type this
	// repository is registered to handle. In such cases, implementations
	// are suggested to fail-over to GenericCopy().
	//
	// Since: 2.0
	Copy(gui.URI, gui.URI) error
}

// MovableRepository is an extension of the Repository interface which also
// supports moving referenced resources from one URI to another.
//
// Note: both Moveable and Movable are correct spellings, but Movable is newer
// and more accepted.
//
// Since: 2.0
type MovableRepository interface {
	Repository

	// Move will be used to implement calls to storage.Move() for the
	// registered scheme of this repository.
	//
	// A generic implementation is provided by GenericMove().
	//
	// NOTE: the first parameter is the source, the second is the
	// destination.
	//
	// NOTE: if storage.Move() is given two URIs of different schemes, it
	// is possible that only the source URI will be of the type this
	// repository is registered to handle. In such cases, implementations
	// are suggested to fail-over to GenericMove().
	//
	// Since: 2.0
	Move(gui.URI, gui.URI) error
}

// Register registers a storage repository so that operations on URIs of the
// registered scheme will use methods implemented by the relevant repository
// implementation.
//
// Since: 2.0
func Register(scheme string, repository Repository) {
	scheme = strings.ToLower(scheme)

	prev, ok := repositoryTable[scheme]

	if ok {
		prev.Destroy(scheme)
	}

	repositoryTable[scheme] = repository
}

// ForURI returns the Repository instance which is registered to handle URIs of
// the given scheme. This is a helper method that calls ForScheme() on the
// scheme of the given URI.
//
// NOTE: this function is intended to be used specifically by the storage
// package. It generally should not be used outside of the Bhojpur GUI package -
// instead you should use the methods in the storage package.
//
// Since: 2.0
func ForURI(u gui.URI) (Repository, error) {
	return ForScheme(u.Scheme())
}

// ForScheme returns the Repository instance which is registered to handle URIs
// of the given scheme.
//
// NOTE: this function is intended to be used specifically by the storage
// package. It generally should not be used outside of the Bhojpur GUI package -
// instead you should use the methods in the storage package.
//
// Since: 2.0
func ForScheme(scheme string) (Repository, error) {
	repo, ok := repositoryTable[scheme]

	if !ok {
		return nil, fmt.Errorf("no repository registered for scheme '%s'", scheme)
	}

	return repo, nil
}
