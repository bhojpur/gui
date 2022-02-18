package async_test

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal/async"
)

func TestQueue(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q := async.NewCanvasObjectQueue()
		if q.Out() != nil {
			t.Fatalf("dequeue empty queue returns non-nil")
		}
	})

	t.Run("length", func(t *testing.T) {
		q := async.NewCanvasObjectQueue()
		if q.Len() != 0 {
			t.Fatalf("empty queue has non-zero length")
		}

		obj := canvas.NewRectangle(color.Black)

		q.In(obj)
		if q.Len() != 1 {
			t.Fatalf("count of enqueue wrong, want %d, got %d.", 1, q.Len())
		}

		q.Out()
		if q.Len() != 0 {
			t.Fatalf("count of dequeue wrong, want %d, got %d", 0, q.Len())
		}
	})

	t.Run("in-out", func(t *testing.T) {
		q := async.NewCanvasObjectQueue()

		want := []gui.CanvasObject{
			canvas.NewRectangle(color.Black),
			canvas.NewRectangle(color.Black),
			canvas.NewRectangle(color.Black),
		}

		for i := 0; i < len(want); i++ {
			q.In(want[i])
		}

		var x []gui.CanvasObject
		for {
			e := q.Out()
			if e == nil {
				break
			}
			x = append(x, e)
		}

		for i := 0; i < len(want); i++ {
			if x[i] != want[i] {
				t.Fatalf("input does not match output, want %+v, got %+v", want[i], x[i])
			}
		}
	})
}
