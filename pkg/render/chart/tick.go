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

import (
	"fmt"
	"math"
	"strings"
)

// TicksProvider is a type that provides ticks.
type TicksProvider interface {
	GetTicks(r Renderer, defaults Style, vf ValueFormatter) []Tick
}

// Tick represents a label on an axis.
type Tick struct {
	Value float64
	Label string
}

// Ticks is an array of ticks.
type Ticks []Tick

// Len returns the length of the ticks set.
func (t Ticks) Len() int {
	return len(t)
}

// Swap swaps two elements.
func (t Ticks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less returns if i's value is less than j's value.
func (t Ticks) Less(i, j int) bool {
	return t[i].Value < t[j].Value
}

// String returns a string representation of the set of ticks.
func (t Ticks) String() string {
	var values []string
	for i, tick := range t {
		values = append(values, fmt.Sprintf("[%d: %s]", i, tick.Label))
	}
	return strings.Join(values, ", ")
}

// GenerateContinuousTicks generates a set of ticks.
func GenerateContinuousTicks(r Renderer, ra Range, isVertical bool, style Style, vf ValueFormatter) []Tick {
	if vf == nil {
		vf = FloatValueFormatter
	}

	var ticks []Tick
	min, max := ra.GetMin(), ra.GetMax()

	if ra.IsDescending() {
		ticks = append(ticks, Tick{
			Value: max,
			Label: vf(max),
		})
	} else {
		ticks = append(ticks, Tick{
			Value: min,
			Label: vf(min),
		})
	}

	minLabel := vf(min)
	style.GetTextOptions().WriteToRenderer(r)
	labelBox := r.MeasureText(minLabel)

	var tickSize float64
	if isVertical {
		tickSize = float64(labelBox.Height() + DefaultMinimumTickVerticalSpacing)
	} else {
		tickSize = float64(labelBox.Width() + DefaultMinimumTickHorizontalSpacing)
	}

	domain := float64(ra.GetDomain())
	domainRemainder := domain - (tickSize * 2)
	intermediateTickCount := int(math.Floor(float64(domainRemainder) / float64(tickSize)))

	rangeDelta := math.Abs(max - min)
	tickStep := rangeDelta / float64(intermediateTickCount)

	roundTo := GetRoundToForDelta(rangeDelta) / 10
	intermediateTickCount = MinInt(intermediateTickCount, DefaultTickCountSanityCheck)

	for x := 1; x < intermediateTickCount; x++ {
		var tickValue float64
		if ra.IsDescending() {
			tickValue = max - RoundUp(tickStep*float64(x), roundTo)
		} else {
			tickValue = min + RoundUp(tickStep*float64(x), roundTo)
		}
		ticks = append(ticks, Tick{
			Value: tickValue,
			Label: vf(tickValue),
		})
	}

	if ra.IsDescending() {
		ticks = append(ticks, Tick{
			Value: min,
			Label: vf(min),
		})
	} else {
		ticks = append(ticks, Tick{
			Value: max,
			Label: vf(max),
		})
	}

	return ticks
}
