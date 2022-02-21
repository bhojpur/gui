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
	"testing"

	"github.com/bhojpur/gui/pkg/render/chart/testutil"
)

func TestBuffer(t *testing.T) {
	// replaced new assertions helper

	buffer := NewValueBuffer()

	buffer.Enqueue(1)
	testutil.AssertEqual(t, 1, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 1, buffer.PeekBack())

	buffer.Enqueue(2)
	testutil.AssertEqual(t, 2, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 2, buffer.PeekBack())

	buffer.Enqueue(3)
	testutil.AssertEqual(t, 3, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 3, buffer.PeekBack())

	buffer.Enqueue(4)
	testutil.AssertEqual(t, 4, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 4, buffer.PeekBack())

	buffer.Enqueue(5)
	testutil.AssertEqual(t, 5, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 5, buffer.PeekBack())

	buffer.Enqueue(6)
	testutil.AssertEqual(t, 6, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 6, buffer.PeekBack())

	buffer.Enqueue(7)
	testutil.AssertEqual(t, 7, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 7, buffer.PeekBack())

	buffer.Enqueue(8)
	testutil.AssertEqual(t, 8, buffer.Len())
	testutil.AssertEqual(t, 1, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value := buffer.Dequeue()
	testutil.AssertEqual(t, 1, value)
	testutil.AssertEqual(t, 7, buffer.Len())
	testutil.AssertEqual(t, 2, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value = buffer.Dequeue()
	testutil.AssertEqual(t, 2, value)
	testutil.AssertEqual(t, 6, buffer.Len())
	testutil.AssertEqual(t, 3, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value = buffer.Dequeue()
	testutil.AssertEqual(t, 3, value)
	testutil.AssertEqual(t, 5, buffer.Len())
	testutil.AssertEqual(t, 4, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value = buffer.Dequeue()
	testutil.AssertEqual(t, 4, value)
	testutil.AssertEqual(t, 4, buffer.Len())
	testutil.AssertEqual(t, 5, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value = buffer.Dequeue()
	testutil.AssertEqual(t, 5, value)
	testutil.AssertEqual(t, 3, buffer.Len())
	testutil.AssertEqual(t, 6, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value = buffer.Dequeue()
	testutil.AssertEqual(t, 6, value)
	testutil.AssertEqual(t, 2, buffer.Len())
	testutil.AssertEqual(t, 7, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value = buffer.Dequeue()
	testutil.AssertEqual(t, 7, value)
	testutil.AssertEqual(t, 1, buffer.Len())
	testutil.AssertEqual(t, 8, buffer.Peek())
	testutil.AssertEqual(t, 8, buffer.PeekBack())

	value = buffer.Dequeue()
	testutil.AssertEqual(t, 8, value)
	testutil.AssertEqual(t, 0, buffer.Len())
	testutil.AssertZero(t, buffer.Peek())
	testutil.AssertZero(t, buffer.PeekBack())
}

func TestBufferClear(t *testing.T) {
	// replaced new assertions helper

	buffer := NewValueBuffer()
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)

	testutil.AssertEqual(t, 8, buffer.Len())

	buffer.Clear()
	testutil.AssertEqual(t, 0, buffer.Len())
	testutil.AssertZero(t, buffer.Peek())
	testutil.AssertZero(t, buffer.PeekBack())
}

func TestBufferArray(t *testing.T) {
	// replaced new assertions helper

	buffer := NewValueBuffer()
	buffer.Enqueue(1)
	buffer.Enqueue(2)
	buffer.Enqueue(3)
	buffer.Enqueue(4)
	buffer.Enqueue(5)

	contents := buffer.Array()
	testutil.AssertLen(t, contents, 5)
	testutil.AssertEqual(t, 1, contents[0])
	testutil.AssertEqual(t, 2, contents[1])
	testutil.AssertEqual(t, 3, contents[2])
	testutil.AssertEqual(t, 4, contents[3])
	testutil.AssertEqual(t, 5, contents[4])
}

func TestBufferEach(t *testing.T) {
	// replaced new assertions helper

	buffer := NewValueBuffer()

	for x := 1; x < 17; x++ {
		buffer.Enqueue(float64(x))
	}

	called := 0
	buffer.Each(func(_ int, v float64) {
		if v == float64(called+1) {
			called++
		}
	})

	testutil.AssertEqual(t, 16, called)
}

func TestNewBuffer(t *testing.T) {
	// replaced new assertions helper

	empty := NewValueBuffer()
	testutil.AssertNotNil(t, empty)
	testutil.AssertZero(t, empty.Len())
	testutil.AssertEqual(t, bufferDefaultCapacity, empty.Capacity())
	testutil.AssertZero(t, empty.Peek())
	testutil.AssertZero(t, empty.PeekBack())
}

func TestNewBufferWithValues(t *testing.T) {
	// replaced new assertions helper

	values := NewValueBuffer(1, 2, 3, 4, 5)
	testutil.AssertNotNil(t, values)
	testutil.AssertEqual(t, 5, values.Len())
	testutil.AssertEqual(t, 1, values.Peek())
	testutil.AssertEqual(t, 5, values.PeekBack())
}

func TestBufferGrowth(t *testing.T) {
	// replaced new assertions helper

	values := NewValueBuffer(1, 2, 3, 4, 5)
	for i := 0; i < 1<<10; i++ {
		values.Enqueue(float64(i))
	}

	testutil.AssertEqual(t, 1<<10-1, values.PeekBack())
}
