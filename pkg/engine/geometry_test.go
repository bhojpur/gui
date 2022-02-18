package engine_test

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

func TestPosition_Add(t *testing.T) {
	pos1 := gui.NewPos(10, 10)
	pos2 := gui.NewPos(25, 25)

	pos3 := pos1.Add(pos2)

	assert.Equal(t, float32(35), pos3.X)
	assert.Equal(t, float32(35), pos3.Y)
}

func TestPosition_Add_Size(t *testing.T) {
	pos1 := gui.NewPos(10, 10)
	s := gui.NewSize(25, 25)

	pos2 := pos1.Add(s)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_Add_Vector(t *testing.T) {
	pos1 := gui.NewPos(10, 10)
	v := gui.NewDelta(25, 25)

	pos2 := pos1.Add(v)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_AddXY(t *testing.T) {
	pos1 := gui.NewPos(10, 10)

	pos2 := pos1.AddXY(25, 25)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_IsZero(t *testing.T) {
	for name, tt := range map[string]struct {
		p    gui.Position
		want bool
	}{
		"zero value":       {gui.Position{}, true},
		"0,0":              {gui.NewPos(0, 0), true},
		"zero X":           {gui.NewPos(0, 42), false},
		"zero Y":           {gui.NewPos(17, 0), false},
		"non-zero X and Y": {gui.NewPos(6, 9), false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.IsZero())
		})
	}
}

func TestPosition_Subtract(t *testing.T) {
	pos1 := gui.NewPos(25, 25)
	pos2 := gui.NewPos(10, 10)

	pos3 := pos1.Subtract(pos2)

	assert.Equal(t, float32(15), pos3.X)
	assert.Equal(t, float32(15), pos3.Y)
}

func TestPosition_SubtractXY(t *testing.T) {
	pos1 := gui.NewPos(25, 25)

	pos2 := pos1.SubtractXY(10, 10)

	assert.Equal(t, float32(15), pos2.X)
	assert.Equal(t, float32(15), pos2.Y)
}

func TestSize_Add(t *testing.T) {
	size1 := gui.NewSize(10, 10)
	size2 := gui.NewSize(25, 25)

	size3 := size1.Add(size2)

	assert.Equal(t, float32(35), size3.Width)
	assert.Equal(t, float32(35), size3.Height)
}

func TestSize_Add_Position(t *testing.T) {
	size1 := gui.NewSize(10, 10)
	p := gui.NewSize(25, 25)

	size2 := size1.Add(p)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_Add_Vector(t *testing.T) {
	size1 := gui.NewSize(10, 10)
	v := gui.NewDelta(25, 25)

	size2 := size1.Add(v)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_AddWidthHeight(t *testing.T) {
	size1 := gui.NewSize(10, 10)

	size2 := size1.AddWidthHeight(25, 25)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_IsZero(t *testing.T) {
	for name, tt := range map[string]struct {
		s    gui.Size
		want bool
	}{
		"zero value":    {gui.Size{}, true},
		"0x0":           {gui.NewSize(0, 0), true},
		"zero width":    {gui.NewSize(0, 42), false},
		"zero height":   {gui.NewSize(17, 0), false},
		"non-zero area": {gui.NewSize(6, 9), false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.IsZero())
		})
	}
}

func TestSize_Max(t *testing.T) {
	size1 := gui.NewSize(10, 100)
	size2 := gui.NewSize(100, 10)

	size3 := size1.Max(size2)

	assert.Equal(t, float32(100), size3.Width)
	assert.Equal(t, float32(100), size3.Height)
}

func TestSize_Min(t *testing.T) {
	size1 := gui.NewSize(10, 100)
	size2 := gui.NewSize(100, 10)

	size3 := size1.Min(size2)

	assert.Equal(t, float32(10), size3.Width)
	assert.Equal(t, float32(10), size3.Height)
}

func TestSize_Subtract(t *testing.T) {
	size1 := gui.NewSize(25, 25)
	size2 := gui.NewSize(10, 10)

	size3 := size1.Subtract(size2)

	assert.Equal(t, float32(15), size3.Width)
	assert.Equal(t, float32(15), size3.Height)
}

func TestSize_SubtractWidthHeight(t *testing.T) {
	size1 := gui.NewSize(25, 25)

	size2 := size1.SubtractWidthHeight(10, 10)

	assert.Equal(t, float32(15), size2.Width)
	assert.Equal(t, float32(15), size2.Height)
}

func TestVector_IsZero(t *testing.T) {
	v := gui.NewDelta(0, 0)

	assert.True(t, v.IsZero())

	v.DX = 1
	assert.False(t, v.IsZero())
}
