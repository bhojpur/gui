//go:build !windows || !ci
// +build !windows !ci

package gl

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
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/stretchr/testify/assert"
)

func TestInnerRect_Stretch(t *testing.T) {
	pos := gui.NewPos(10, 10)
	size := gui.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillStretch, 0.0)

	assert.Equal(t, size, innerSize)
	assert.Equal(t, pos, innerPos)
}

func TestInnerRect_StretchIgnoreRatio(t *testing.T) {
	pos := gui.NewPos(10, 10)
	size := gui.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillStretch, 2.0)

	assert.Equal(t, size, innerSize)
	assert.Equal(t, pos, innerPos)
}

func TestInnerRect_ContainScale(t *testing.T) {
	pos := gui.NewPos(10, 10)
	size := gui.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillContain, 1.0)

	assert.Equal(t, size, innerSize)
	assert.Equal(t, pos, innerPos)
}

func TestInnerRect_ContainPillarbox(t *testing.T) {
	pos := gui.NewPos(10, 10)
	size := gui.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillContain, 0.5)

	assert.Equal(t, gui.NewSize(20, 40), innerSize)
	assert.Equal(t, gui.NewPos(20, 10), innerPos)
}

func TestInnerRect_Original(t *testing.T) {
	// TODO add check for minsize somehow?
	pos := gui.NewPos(10, 10)
	size := gui.NewSize(40, 40)

	innerSize1, innerPos1 := rectInnerCoords(size, pos, canvas.ImageFillOriginal, 0.5)
	innerSize2, innerPos2 := rectInnerCoords(size, pos, canvas.ImageFillContain, 0.5)

	assert.Equal(t, innerSize2, innerSize1)
	assert.Equal(t, innerPos2, innerPos1)
}
