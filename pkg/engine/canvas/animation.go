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
	"image/color"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
)

const (
	// DurationStandard is the time a standard interface animation will run.
	//
	// Since: 2.0
	DurationStandard = time.Millisecond * 300
	// DurationShort is the time a subtle or small transition should use.
	//
	// Since: 2.0
	DurationShort = time.Millisecond * 150
)

// NewColorRGBAAnimation sets up a new animation that will transition from the start to stop Color over
// the specified Duration. The colour transition will move linearly through the RGB colour space.
// The content of fn should apply the color values to an object and refresh it.
// You should call Start() on the returned animation to start it.
//
// Since: 2.0
func NewColorRGBAAnimation(start, stop color.Color, d time.Duration, fn func(color.Color)) *gui.Animation {
	r1, g1, b1, a1 := start.RGBA()
	r2, g2, b2, a2 := stop.RGBA()

	rStart := int(r1 >> 8)
	gStart := int(g1 >> 8)
	bStart := int(b1 >> 8)
	aStart := int(a1 >> 8)
	rDelta := float32(int(r2>>8) - rStart)
	gDelta := float32(int(g2>>8) - gStart)
	bDelta := float32(int(b2>>8) - bStart)
	aDelta := float32(int(a2>>8) - aStart)

	return &gui.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(color.RGBA{R: scaleChannel(rStart, rDelta, done), G: scaleChannel(gStart, gDelta, done),
				B: scaleChannel(bStart, bDelta, done), A: scaleChannel(aStart, aDelta, done)})
		}}
}

// NewPositionAnimation sets up a new animation that will transition from the start to stop Position over
// the specified Duration. The content of fn should apply the position value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
//
// Since: 2.0
func NewPositionAnimation(start, stop gui.Position, d time.Duration, fn func(gui.Position)) *gui.Animation {
	xDelta := float32(stop.X - start.X)
	yDelta := float32(stop.Y - start.Y)

	return &gui.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(gui.NewPos(scaleVal(start.X, xDelta, done), scaleVal(start.Y, yDelta, done)))
		}}
}

// NewSizeAnimation sets up a new animation that will transition from the start to stop Size over
// the specified Duration. The content of fn should apply the size value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
//
// Since: 2.0
func NewSizeAnimation(start, stop gui.Size, d time.Duration, fn func(gui.Size)) *gui.Animation {
	widthDelta := float32(stop.Width - start.Width)
	heightDelta := float32(stop.Height - start.Height)

	return &gui.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(gui.NewSize(scaleVal(start.Width, widthDelta, done), scaleVal(start.Height, heightDelta, done)))
		}}
}

func scaleChannel(start int, diff, done float32) uint8 {
	return uint8(start + int(diff*done))
}

func scaleVal(start float32, delta, done float32) float32 {
	return start + delta*done
}
