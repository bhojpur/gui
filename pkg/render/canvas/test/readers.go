package test

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

import "io"

// PlainReader implements an io.Reader and wraps over an existing io.Reader to
// hide other functions it implements.
type PlainReader struct {
	r io.Reader
}

// NewPlainReader returns a new PlainReader.
func NewPlainReader(r io.Reader) *PlainReader {
	return &PlainReader{r}
}

// Read implements the io.Reader interface.
func (r *PlainReader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

// ErrorReader implements an io.Reader that will do N successive reads before
// it returns ErrPlain.
type ErrorReader struct {
	n int
}

// NewErrorReader returns a new ErrorReader.
func NewErrorReader(n int) *ErrorReader {
	return &ErrorReader{n}
}

// Read implements the io.Reader interface.
func (r *ErrorReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.n == 0 {
		return 0, ErrPlain
	}
	r.n--
	b[0] = '.'
	return 1, nil
}

// InfiniteReader implements an io.Reader that will always read-in one character.
type InfiniteReader struct{}

// NewInfiniteReader returns a new InfiniteReader.
func NewInfiniteReader() *InfiniteReader {
	return &InfiniteReader{}
}

// Read implements the io.Reader interface.
func (r *InfiniteReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	b[0] = '.'
	return 1, nil
}

// EmptyReader implements an io.Reader that will always return 0, nil.
type EmptyReader struct {
}

// NewEmptyReader returns a new EmptyReader.
func NewEmptyReader() *EmptyReader {
	return &EmptyReader{}
}

// Read implements the io.Reader interface.
func (r *EmptyReader) Read(b []byte) (n int, err error) {
	return 0, io.EOF
}
