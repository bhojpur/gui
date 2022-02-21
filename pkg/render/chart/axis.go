package chart

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

// TickPosition is an enumeration of possible tick drawing positions.
type TickPosition int

const (
	// TickPositionUnset means to use the default tick position.
	TickPositionUnset TickPosition = 0
	// TickPositionBetweenTicks draws the labels for a tick between the previous and current tick.
	TickPositionBetweenTicks TickPosition = 1
	// TickPositionUnderTick draws the tick below the tick.
	TickPositionUnderTick TickPosition = 2
)

// YAxisType is a type of y-axis; it can either be primary or secondary.
type YAxisType int

const (
	// YAxisPrimary is the primary axis.
	YAxisPrimary YAxisType = 0
	// YAxisSecondary is the secondary axis.
	YAxisSecondary YAxisType = 1
)

// Axis is a chart feature detailing what values happen where.
type Axis interface {
	GetName() string
	SetName(name string)

	GetStyle() Style
	SetStyle(style Style)

	GetTicks() []Tick
	GenerateTicks(r Renderer, ra Range, vf ValueFormatter) []Tick

	// GenerateGridLines returns the gridlines for the axis.
	GetGridLines(ticks []Tick) []GridLine

	// Measure should return an absolute box for the axis.
	// This is used when auto-fitting the canvas to the background.
	Measure(r Renderer, canvasBox Box, ra Range, style Style, ticks []Tick) Box

	// Render renders the axis.
	Render(r Renderer, canvasBox Box, ra Range, style Style, ticks []Tick)
}
