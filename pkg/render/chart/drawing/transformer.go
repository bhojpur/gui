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

// Transformer apply the Matrix transformation tr
type Transformer struct {
	Tr        Matrix
	Flattener Flattener
}

// MoveTo implements the path builder interface.
func (t Transformer) MoveTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.MoveTo(u, v)
}

// LineTo implements the path builder interface.
func (t Transformer) LineTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.LineTo(u, v)
}

// LineJoin implements the path builder interface.
func (t Transformer) LineJoin() {
	t.Flattener.LineJoin()
}

// Close implements the path builder interface.
func (t Transformer) Close() {
	t.Flattener.Close()
}

// End implements the path builder interface.
func (t Transformer) End() {
	t.Flattener.End()
}
