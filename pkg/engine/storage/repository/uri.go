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
	"bufio"
	"mime"
	"path/filepath"
	"strings"
	"unicode/utf8"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// Declare conformance with gui.URI interface.
var _ gui.URI = &uri{}

type uri struct {
	scheme    string
	authority string
	// haveAuthority lets us distinguish between a present-but-empty
	// authority, and having no authority. This is needed because net/url
	// incorrectly handles scheme:/absolute/path URIs.
	haveAuthority bool
	path          string
	query         string
	fragment      string
}

func (u *uri) Extension() string {
	return filepath.Ext(u.path)
}

func (u *uri) Name() string {
	return filepath.Base(u.path)
}

func (u *uri) MimeType() string {

	mimeTypeFull := mime.TypeByExtension(u.Extension())
	if mimeTypeFull == "" {
		mimeTypeFull = "text/plain"

		repo, err := ForURI(u)
		if err != nil {
			return "application/octet-stream"
		}

		readCloser, err := repo.Reader(u)
		if err == nil {
			defer readCloser.Close()
			scanner := bufio.NewScanner(readCloser)
			if scanner.Scan() && !utf8.Valid(scanner.Bytes()) {
				mimeTypeFull = "application/octet-stream"
			}
		}
	}

	return strings.Split(mimeTypeFull, ";")[0]
}

func (u *uri) Scheme() string {
	return u.scheme
}

func (u *uri) String() string {
	// NOTE: this string reconstruction is mandated by IETF RFC3986,
	// section 5.3, pp. 35.

	s := u.scheme + ":"
	if u.haveAuthority {
		s += "//" + u.authority
	}
	s += u.path
	if len(u.query) > 0 {
		s += "?" + u.query
	}
	if len(u.fragment) > 0 {
		s += "#" + u.fragment
	}
	return s
}

func (u *uri) Authority() string {
	return u.authority
}

func (u *uri) Path() string {
	return u.path
}

func (u *uri) Query() string {
	return u.query
}

func (u *uri) Fragment() string {
	return u.fragment
}
