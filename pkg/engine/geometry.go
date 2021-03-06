package engine

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

var _ Vector2 = (*Delta)(nil)
var _ Vector2 = (*Position)(nil)
var _ Vector2 = (*Size)(nil)

// Vector2 marks geometry types that can operate as a coordinate vector.
type Vector2 interface {
	Components() (float32, float32)
	IsZero() bool
}

// Delta is a generic X, Y coordinate, size or movement representation.
type Delta struct {
	DX, DY float32
}

// NewDelta returns a newly allocated Delta representing a movement in the X and Y axis.
func NewDelta(dx float32, dy float32) Delta {
	return Delta{DX: dx, DY: dy}
}

// Components returns the X and Y elements of this Delta.
func (v Delta) Components() (float32, float32) {
	return v.DX, v.DY
}

// IsZero returns whether the Position is at the zero-point.
func (v Delta) IsZero() bool {
	return v.DX == 0.0 && v.DY == 0.0
}

// Position describes a generic X, Y coordinate relative to a parent Canvas
// or CanvasObject.
type Position struct {
	X float32 // The position from the parent's left edge
	Y float32 // The position from the parent's top edge
}

// NewPos returns a newly allocated Position representing the specified coordinates.
func NewPos(x float32, y float32) Position {
	return Position{x, y}
}

// Add returns a new Position that is the result of offsetting the current
// position by p2 X and Y.
func (p Position) Add(v Vector2) Position {
	x, y := v.Components()
	return Position{p.X + x, p.Y + y}
}

// AddXY returns a new Position by adding x and y to the current one.
func (p Position) AddXY(x, y float32) Position {
	return Position{p.X + x, p.Y + y}
}

// Components returns the X and Y elements of this Position
func (p Position) Components() (float32, float32) {
	return p.X, p.Y
}

// IsZero returns whether the Position is at the zero-point.
func (p Position) IsZero() bool {
	return p.X == 0.0 && p.Y == 0.0
}

// Subtract returns a new Position that is the result of offsetting the current
// position by p2 -X and -Y.
func (p Position) Subtract(v Vector2) Position {
	x, y := v.Components()
	return Position{p.X - x, p.Y - y}
}

// SubtractXY returns a new Position by subtracting x and y from the current one.
func (p Position) SubtractXY(x, y float32) Position {
	return Position{p.X - x, p.Y - y}
}

// Size describes something with width and height.
type Size struct {
	Width  float32 // The number of units along the X axis.
	Height float32 // The number of units along the Y axis.
}

// NewSize returns a newly allocated Size of the specified dimensions.
func NewSize(w float32, h float32) Size {
	return Size{w, h}
}

// Add returns a new Size that is the result of increasing the current size by
// s2 Width and Height.
func (s Size) Add(v Vector2) Size {
	w, h := v.Components()
	return Size{s.Width + w, s.Height + h}
}

// AddWidthHeight returns a new Size by adding width and height to the current one.
func (s Size) AddWidthHeight(width, height float32) Size {
	return Size{s.Width + width, s.Height + height}
}

// IsZero returns whether the Size has zero width and zero height.
func (s Size) IsZero() bool {
	return s.Width == 0.0 && s.Height == 0.0
}

// Max returns a new Size that is the maximum of the current Size and s2.
func (s Size) Max(v Vector2) Size {
	x, y := v.Components()

	maxW := Max(s.Width, x)
	maxH := Max(s.Height, y)

	return NewSize(maxW, maxH)
}

// Min returns a new Size that is the minimum of the current Size and s2.
func (s Size) Min(v Vector2) Size {
	x, y := v.Components()

	minW := Min(s.Width, x)
	minH := Min(s.Height, y)

	return NewSize(minW, minH)
}

// Components returns the Width and Height elements of this Size
func (s Size) Components() (float32, float32) {
	return s.Width, s.Height
}

// Subtract returns a new Size that is the result of decreasing the current size
// by s2 Width and Height.
func (s Size) Subtract(v Vector2) Size {
	w, h := v.Components()
	return Size{s.Width - w, s.Height - h}
}

// SubtractWidthHeight returns a new Size by subtracting width and height from the current one.
func (s Size) SubtractWidthHeight(width, height float32) Size {
	return Size{s.Width - width, s.Height - height}
}
