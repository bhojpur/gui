package binding

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

	"github.com/bhojpur/gui/pkg/engine/storage"
)

func TestStripPrecision(t *testing.T) {
	format := "%2.3f"
	assert.Equal(t, "%2f", stripFormatPrecision(format))
	format = "total=%3.4f%%"
	assert.Equal(t, "total=%3f%%", stripFormatPrecision(format))

	format = "%.2f"
	assert.Equal(t, "%f", stripFormatPrecision(format))

	format = "%4d"
	assert.Equal(t, "%4d", stripFormatPrecision(format))

	format = "%v"
	assert.Equal(t, "%v", stripFormatPrecision(format))
}

func TestURIFromStringHelper(t *testing.T) {
	str := "file:///tmp/test.txt"
	u, err := uriFromString(str)

	assert.Nil(t, err)
	assert.Equal(t, str, u.String())
}

func TestURIToStringHelper(t *testing.T) {
	u := storage.NewFileURI("/tmp/test.txt")
	str, err := uriToString(u)

	assert.Nil(t, err)
	assert.Equal(t, u.String(), str)

	str, err = uriToString(nil)
	assert.Nil(t, err)
	assert.Equal(t, "", str)
}
