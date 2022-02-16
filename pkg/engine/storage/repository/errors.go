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
	"errors"
)

var (
	// ErrOperationNotSupported may be thrown by certain functions in the storage
	// or repository packages which operate on URIs if an operation is attempted
	// that is not supported for the scheme relevant to the URI, normally because
	// the underlying repository has either not implemented the relevant function,
	// or has explicitly returned this error.
	//
	// Since: 2.0
	ErrOperationNotSupported = errors.New("operation not supported for this URI")

	// ErrURIRoot should be thrown by gui.URI implementations when the caller
	// attempts to take the parent of the root. This way, downstream code that
	// wants to programmatically walk up a URIs parent's will know when to stop
	// iterating.
	//
	// Since: 2.0
	ErrURIRoot = errors.New("cannot take the parent of the root element in a URI")
)
