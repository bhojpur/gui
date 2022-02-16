package binding

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
)

func BenchmarkBoolToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bo := NewBool()
		s := BoolToString(bo)
		s.Get()

		bo.Set(true)
		s.Get()

		s.Set("trap")
		bo.Get()

		s.Set("false")
		bo.Get()
	}
}

func BenchmarkFloatToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := NewFloat()
		s := FloatToString(f)
		s.Get()

		f.Set(0.3)
		s.Get()

		s.Set("wrong")
		f.Get()

		s.Set("5.00")
		f.Get()
	}
}

func BenchmarkIntToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		i := NewInt()
		s := IntToString(i)
		s.Get()

		i.Set(3)
		s.Get()

		s.Set("wrong")
		i.Get()

		s.Set("5")
		i.Get()
	}
}

func BenchmarkStringToBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewString()
		b := StringToBool(s)
		b.Get()

		s.Set("true")
		b.Get()

		s.Set("trap") // bug in fmt.SScanf means "wrong" parses as "false"
		b.Get()

		b.Set(false)
		s.Get()
	}
}

func BenchmarkStringToFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewString()
		f := StringToFloat(s)
		f.Get()

		s.Set("3")
		f.Get()

		s.Set("wrong")
		f.Get()

		f.Set(5)
		s.Get()
	}
}

func BenchmarkStringToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewString()
		i := StringToInt(s)
		i.Get()

		s.Set("3")
		i.Get()

		s.Set("wrong")
		i.Get()

		i.Set(5)
		s.Get()
	}
}
