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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericParent(t *testing.T) {
	cases := []struct {
		input  string
		expect string
		err    error
	}{
		{"foo://example.com:8042/over/there?name=ferret#nose", "foo://example.com:8042/over?name=ferret#nose", nil},
		{"file:///", "file:///", ErrURIRoot},
		{"file:///foo", "file:///", nil},
	}

	for i, c := range cases {
		t.Logf("case %d, input='%s', expect='%s', err='%v'\n", i, c.input, c.expect, c.err)

		u, err := ParseURI(c.input)
		assert.Nil(t, err)

		res, err := GenericParent(u)
		assert.Equal(t, c.err, err)

		// In the case where there is a non-nil error, res is defined
		// to be nil, so we cannot call res.String() without causing
		// a panic.
		if err == nil {
			assert.Equal(t, c.expect, res.String())
		}
	}
}

func TestGenericChild(t *testing.T) {
	cases := []struct {
		input     string
		component string
		expect    string
		err       error
	}{
		{"foo://example.com:8042/over/there?name=ferret#nose", "bar", "foo://example.com:8042/over/there/bar?name=ferret#nose", nil},
		{"file:///", "quux", "file:///quux", nil},
		{"file:///foo", "baz", "file:///foo/baz", nil},
	}

	for i, c := range cases {
		t.Logf("case %d, input='%s', component='%s', expect='%s', err='%v'\n", i, c.input, c.component, c.expect, c.err)

		u, err := ParseURI(c.input)
		assert.Nil(t, err)

		res, err := GenericChild(u, c.component)
		assert.Equal(t, c.err, err)

		assert.Equal(t, c.expect, res.String())
	}
}
