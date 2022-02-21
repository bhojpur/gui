package drawing

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

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

type point struct {
	X, Y float64
}

type mockLine struct {
	inner []point
}

func (ml *mockLine) LineTo(x, y float64) {
	ml.inner = append(ml.inner, point{x, y})
}

func (ml mockLine) Len() int {
	return len(ml.inner)
}

func TestTraceQuad(t *testing.T) {
	// replaced new assertions helper

	// Quad
	// x1, y1, cpx1, cpy2, x2, y2 float64
	// do the 9->12 circle segment
	quad := []float64{10, 20, 20, 20, 20, 10}
	liner := &mockLine{}
	TraceQuad(liner, quad, 0.5)
	testutil.AssertNotZero(t, liner.Len())
}
