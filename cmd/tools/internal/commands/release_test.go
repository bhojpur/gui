package commands

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

func TestReleaser_nameFromCertInfo(t *testing.T) {
	rel := &Releaser{}
	cert := "CN=Company, O=Company, L=City, S=State, C=Country"
	assert.Equal(t, "Company", rel.nameFromCertInfo(cert))
	assert.Equal(t, "Fallback", rel.nameFromCertInfo("Fallback"))
	assert.Equal(t, "Fallback", rel.nameFromCertInfo("Fallback, extra"))

	badCase := "Cn=Company, O=Company, L=City, S=State, C=Country"
	assert.Equal(t, "Company", rel.nameFromCertInfo(badCase))
}

func TestIsValidMacOSCategory(t *testing.T) {
	assert.True(t, isValidMacOSCategory("games"))
	assert.True(t, isValidMacOSCategory("utilities"))
	assert.True(t, isValidMacOSCategory("Games"))

	assert.False(t, isValidMacOSCategory("sporps"))
	assert.False(t, isValidMacOSCategory("android-games"))
	assert.False(t, isValidMacOSCategory(""))
}
