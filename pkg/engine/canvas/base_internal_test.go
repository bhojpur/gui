package canvas

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

func TestBase_MinSize(t *testing.T) {
	base := &baseObject{}
	min := base.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

func TestBase_Move(t *testing.T) {
	base := &baseObject{}
	base.Move(gui.NewPos(10, 10))
	pos := base.Position()

	assert.Equal(t, float32(10), pos.X)
	assert.Equal(t, float32(10), pos.Y)
}

func TestBase_Resize(t *testing.T) {
	base := &baseObject{}
	base.Resize(gui.NewSize(10, 10))
	size := base.Size()

	assert.Equal(t, float32(10), size.Width)
	assert.Equal(t, float32(10), size.Height)
}

func TestBase_HideShow(t *testing.T) {
	base := &baseObject{}
	assert.True(t, base.Visible())

	base.Hide()
	assert.False(t, base.Visible())

	base.Show()
	assert.True(t, base.Visible())
}
