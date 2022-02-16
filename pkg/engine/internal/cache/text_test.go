package cache

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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/stretchr/testify/assert"
)

func TestTextCacheGet(t *testing.T) {
	ResetThemeCaches()
	assert.Equal(t, 0, len(fontSizeCache))

	bound, base := GetFontMetrics("hi", 10, gui.TextStyle{})
	assert.True(t, bound.IsZero())
	assert.Equal(t, float32(0), base)

	SetFontMetrics("hi", 10, gui.TextStyle{}, gui.NewSize(10, 10), 8)
	assert.Equal(t, 1, len(fontSizeCache))

	bound, base = GetFontMetrics("hi", 10, gui.TextStyle{})
	assert.Equal(t, gui.NewSize(10, 10), bound)
	assert.Equal(t, float32(8), base)
}
