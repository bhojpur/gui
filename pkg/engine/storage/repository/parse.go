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
	"path/filepath"
	"runtime"
	"strings"

	uriParser "github.com/fredbi/uri"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// NewFileURI implements the back-end logic to storage.NewFileURI, which you
// should use instead. This is only here because other functions in repository
// need to call it, and it prevents a circular import.
//
// Since: 2.0
func NewFileURI(path string) gui.URI {
	// URIs are supposed to use forward slashes. On Windows, it
	// should be OK to use the platform native filepath with UNIX
	// or NT style paths, with / or \, but when we reconstruct
	// the URI, we want to have / only.
	if runtime.GOOS == "windows" {
		// seems that sometimes we end up with
		// double-backslashes
		path = filepath.ToSlash(path)
	}

	return &uri{
		scheme:        "file",
		haveAuthority: true,
		authority:     "",
		path:          path,
		query:         "",
		fragment:      "",
	}
}

// ParseURI implements the back-end logic for storage.ParseURI, which you
// should use instead. This is only here because other functions in repository
// need to call it, and it prevents a circular import.
//
// Since: 2.0
func ParseURI(s string) (gui.URI, error) {
	// Extract the scheme.
	scheme := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			break
		}
		scheme += string(s[i])
	}
	scheme = strings.ToLower(scheme)

	if scheme == "file" {
		// Does this really deserve to be special? In principle, the
		// purpose of this check is to pass it to NewFileURI, which
		// allows platform path seps in the URI (against the RFC, but
		// easier for people building URIs naively on Windows). Maybe
		// we should punt this to whoever generated the URI in the
		// first place?

		path := s[5:] // everything after file:
		if len(path) > 2 && path[:2] == "//" {
			path = path[2:]
		}

		// Windows files can break authority checks, so just return the parsed file URI
		return NewFileURI(path), nil
	}

	repo, err := ForScheme(scheme)
	if err == nil {
		// If the repository registered for this scheme implements a parser
		if c, ok := repo.(CustomURIRepository); ok {
			return c.ParseURI(s)
		}
	}

	// There was no repository registered, or it did not provide a parser

	// Ugly hack to work around fredbi/uri. Technically, something like
	// foo:/// is invalid because it implies a host, but also has an empty
	// host. However, this is a very common occurrence, so we convert a
	// leading ":///" to "://".
	rest := strings.TrimPrefix(s, scheme+":")
	dummyHost := false
	if len(rest) >= 3 && rest[0:3] == "///" {
		rest = "//" + "TEMP.TEMP/" + strings.TrimPrefix(rest, "///")
		dummyHost = true
	}
	s = scheme + ":" + rest

	l, err := uriParser.Parse(s)
	if err != nil {
		return nil, err
	}

	authority := ""
	if !dummyHost {
		// User info makes no sense without a host, see next comment.
		if userInfo := l.Authority().UserInfo(); len(userInfo) > 0 {
			authority += userInfo + "@"
		}

		// In this case, we had to insert a "host" to make the parser
		// happy, but it isn't really a host, so we can just drop it.
		// If dummyHost isn't set, then we should have a valid host and
		// we can include it as normal.
		authority += l.Authority().Host()

		// Port obviously makes no sense without a host.
		if port := l.Authority().Port(); len(port) > 0 {
			authority += ":" + port
		}
	}

	return &uri{
		scheme:    l.Scheme(),
		authority: authority,
		// workaround for net/url, see type uri struct comments
		haveAuthority: true,
		path:          l.Authority().Path(),
		query:         l.Query().Encode(),
		fragment:      l.Fragment(),
	}, nil
}
