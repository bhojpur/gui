//go:build !ci
// +build !ci

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

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
)

func BenchmarkPosition_Add(b *testing.B) {
	b.Run("Add()", benchmarkPositionAdd)
	b.Run("AddXY()", benchmarkPositionAddXY)
}

func BenchmarkPosition_Subtract(b *testing.B) {
	b.Run("Subtract()", benchmarkPositionSubtract)
	b.Run("SubtractXY()", benchmarkPositionSubtractXY)
}

func BenchmarkSize_Add(b *testing.B) {
	b.Run("Add()", benchmarkSizeAdd)
	b.Run("AddWidthHeight()", benchmarkSizeAddWidthHeight)
}

func BenchmarkSize_Subtract(b *testing.B) {
	b.Run("Subtract()", benchmarkSizeSubtract)
	b.Run("SubtractWidthHeight()", benchmarkSizeSubtractWidthHeight)
}

// This test prevents Position.Add to be simplified to `return p.AddXY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkPositionAdd)
	addXY := testing.Benchmark(benchmarkPositionAddXY)
	assert.Less(t, add.NsPerOp()/addXY.NsPerOp(), int64(5))
}

// This test prevents Position.Subtract to be simplified to `return p.SubtractXY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkPositionSubtract)
	subtractXY := testing.Benchmark(benchmarkPositionSubtractXY)
	assert.Less(t, subtract.NsPerOp()/subtractXY.NsPerOp(), int64(5))
}

// This test prevents Size.Add to be simplified to `return s.AddWidthHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkSizeAdd)
	addWidthHeight := testing.Benchmark(benchmarkSizeAddWidthHeight)
	assert.Less(t, add.NsPerOp()/addWidthHeight.NsPerOp(), int64(5))
}

// This test prevents Size.Subtract to be simplified to `return s.SubtractWidthHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkSizeSubtract)
	subtractWidthHeight := testing.Benchmark(benchmarkSizeSubtractWidthHeight)
	assert.Less(t, subtract.NsPerOp()/subtractWidthHeight.NsPerOp(), int64(5))
}

var benchmarkResult interface{}

func benchmarkPositionAdd(b *testing.B) {
	pos := gui.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.Add(gui.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = pos
}

func benchmarkPositionAddXY(b *testing.B) {
	pos := gui.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.AddXY(float32(n), float32(n))
	}
	benchmarkResult = pos
}

func benchmarkPositionSubtract(b *testing.B) {
	pos := gui.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.Subtract(gui.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = pos
}

func benchmarkPositionSubtractXY(b *testing.B) {
	pos := gui.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.SubtractXY(float32(n), float32(n))
	}
	benchmarkResult = pos
}

func benchmarkSizeAdd(b *testing.B) {
	size := gui.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.Add(gui.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = size
}

func benchmarkSizeAddWidthHeight(b *testing.B) {
	size := gui.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.AddWidthHeight(float32(n), float32(n))
	}
	benchmarkResult = size
}

func benchmarkSizeSubtract(b *testing.B) {
	size := gui.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.Subtract(gui.NewSize(float32(n), float32(n)))
	}
	benchmarkResult = size
}

func benchmarkSizeSubtractWidthHeight(b *testing.B) {
	size := gui.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.SubtractWidthHeight(float32(n), float32(n))
	}
	benchmarkResult = size
}
