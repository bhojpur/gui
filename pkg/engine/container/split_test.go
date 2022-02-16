package container

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
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/stretchr/testify/assert"
)

func TestSplitContainer_MinSize(t *testing.T) {
	rectA := canvas.NewRectangle(color.Black)
	rectA.SetMinSize(gui.NewSize(10, 10))
	rectB := canvas.NewRectangle(color.Black)
	rectB.SetMinSize(gui.NewSize(10, 10))
	t.Run("Horizontal", func(t *testing.T) {
		min := NewHSplit(rectA, rectB).MinSize()
		assert.Equal(t, rectA.MinSize().Width+rectB.MinSize().Width+dividerThickness(), min.Width)
		assert.Equal(t, gui.Max(rectA.MinSize().Height, gui.Max(rectB.MinSize().Height, dividerLength())), min.Height)
	})
	t.Run("Vertical", func(t *testing.T) {
		min := NewVSplit(rectA, rectB).MinSize()
		assert.Equal(t, gui.Max(rectA.MinSize().Width, gui.Max(rectB.MinSize().Width, dividerLength())), min.Width)
		assert.Equal(t, rectA.MinSize().Height+rectB.MinSize().Height+dividerThickness(), min.Height)
	})
}

func TestSplitContainer_Resize(t *testing.T) {
	for name, tt := range map[string]struct {
		horizontal       bool
		size             gui.Size
		wantLeadingPos   gui.Position
		wantLeadingSize  gui.Size
		wantTrailingPos  gui.Position
		wantTrailingSize gui.Size
	}{
		"horizontal": {
			true,
			gui.NewSize(100, 100),
			gui.NewPos(0, 0),
			gui.NewSize(50-dividerThickness()/2, 100),
			gui.NewPos(50+dividerThickness()/2, 0),
			gui.NewSize(50-dividerThickness()/2, 100),
		},
		"vertical": {
			false,
			gui.NewSize(100, 100),
			gui.NewPos(0, 0),
			gui.NewSize(100, 50-dividerThickness()/2),
			gui.NewPos(0, 50+dividerThickness()/2),
			gui.NewSize(100, 50-dividerThickness()/2),
		},
		"horizontal insufficient width": {
			true,
			gui.NewSize(20, 100),
			gui.NewPos(0, 0),
			// minSize of leading is 1/3 of minSize of trailing
			gui.NewSize((20-dividerThickness())/4, 100),
			gui.NewPos((20-dividerThickness())/4+dividerThickness(), 0),
			gui.NewSize((20-dividerThickness())*3/4, 100),
		},
		"vertical insufficient height": {
			false,
			gui.NewSize(100, 20),
			gui.NewPos(0, 0),
			// minSize of leading is 1/3 of minSize of trailing
			gui.NewSize(100, (20-dividerThickness())/4),
			gui.NewPos(0, (20-dividerThickness())/4+dividerThickness()),
			gui.NewSize(100, (20-dividerThickness())*3/4),
		},
		"horizontal zero width": {
			true,
			gui.NewSize(0, 100),
			gui.NewPos(0, 0),
			gui.NewSize(0, 100),
			gui.NewPos(dividerThickness(), 0),
			gui.NewSize(0, 100),
		},
		"horizontal zero height": {
			true,
			gui.NewSize(100, 0),
			gui.NewPos(0, 0),
			gui.NewSize(50-dividerThickness()/2, 0),
			gui.NewPos(50+dividerThickness()/2, 0),
			gui.NewSize(50-dividerThickness()/2, 0),
		},
		"vertical zero width": {
			false,
			gui.NewSize(0, 100),
			gui.NewPos(0, 0),
			gui.NewSize(0, 50-dividerThickness()/2),
			gui.NewPos(0, 50+dividerThickness()/2),
			gui.NewSize(0, 50-dividerThickness()/2),
		},
		"vertical zero height": {
			false,
			gui.NewSize(100, 0),
			gui.NewPos(0, 0),
			gui.NewSize(100, 0),
			gui.NewPos(0, dividerThickness()),
			gui.NewSize(100, 0),
		},
	} {
		t.Run(name, func(t *testing.T) {
			objA := canvas.NewRectangle(color.White)
			objB := canvas.NewRectangle(color.Black)
			objA.SetMinSize(gui.NewSize(10, 10))
			objB.SetMinSize(gui.NewSize(30, 30))
			var c *Split
			if tt.horizontal {
				c = NewHSplit(objA, objB)
			} else {
				c = NewVSplit(objA, objB)
			}
			c.Resize(tt.size)

			assert.Equal(t, tt.wantLeadingPos, objA.Position(), "leading position")
			assert.Equal(t, tt.wantLeadingSize, objA.Size(), "leading size")
			assert.Equal(t, tt.wantTrailingPos, objB.Position(), "trailing position")
			assert.Equal(t, tt.wantTrailingSize, objB.Size(), "trailing size")
		})
	}
}

func TestSplitContainer_SetRatio(t *testing.T) {
	size := gui.NewSize(100, 100)
	usableLength := 100 - float64(dividerThickness())

	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	t.Run("Horizontal", func(t *testing.T) {
		sc := NewHSplit(objA, objB)
		sc.Resize(size)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(0.75)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, float32(0.75*usableLength), sizeA.Width)
			assert.Equal(t, float32(100), sizeA.Height)
			assert.Equal(t, float32(0.25*usableLength), sizeB.Width)
			assert.Equal(t, float32(100), sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.25)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, float32(0.25*usableLength), sizeA.Width)
			assert.Equal(t, float32(100), sizeA.Height)
			assert.Equal(t, float32(0.75*usableLength), sizeB.Width)
			assert.Equal(t, float32(100), sizeB.Height)
		})
	})
	t.Run("Vertical", func(t *testing.T) {
		sc := NewVSplit(objA, objB)
		sc.Resize(size)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(0.75)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, float32(100), sizeA.Width)
			assert.Equal(t, float32(0.75*usableLength), sizeA.Height)
			assert.Equal(t, float32(100), sizeB.Width)
			assert.Equal(t, float32(0.25*usableLength), sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.25)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, float32(100), sizeA.Width)
			assert.Equal(t, float32(0.25*usableLength), sizeA.Height)
			assert.Equal(t, float32(100), sizeB.Width)
			assert.Equal(t, float32(0.75*usableLength), sizeB.Height)
		})
	})
}

func TestSplitContainer_SetRatio_limits(t *testing.T) {
	size := gui.NewSize(50, 50)
	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objA.SetMinSize(size)
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB.SetMinSize(size)
	t.Run("Horizontal", func(t *testing.T) {
		sc := NewHSplit(objA, objB)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(1.0)
			sc.Resize(gui.NewSize(200, 50))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 150-dividerThickness(), sizeA.Width)
			assert.Equal(t, float32(50), sizeA.Height)
			assert.Equal(t, float32(50), sizeB.Width)
			assert.Equal(t, float32(50), sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.0)
			sc.Resize(gui.NewSize(200, 50))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, float32(50), sizeA.Width)
			assert.Equal(t, float32(50), sizeA.Height)
			assert.Equal(t, 150-dividerThickness(), sizeB.Width)
			assert.Equal(t, float32(50), sizeB.Height)
		})
	})
	t.Run("Vertical", func(t *testing.T) {
		sc := NewVSplit(objA, objB)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(1.0)
			sc.Resize(gui.NewSize(50, 200))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, float32(50), sizeA.Width)
			assert.Equal(t, 150-dividerThickness(), sizeA.Height)
			assert.Equal(t, float32(50), sizeB.Width)
			assert.Equal(t, float32(50), sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.0)
			sc.Resize(gui.NewSize(50, 200))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, float32(50), sizeA.Width)
			assert.Equal(t, float32(50), sizeA.Height)
			assert.Equal(t, float32(50), sizeB.Width)
			assert.Equal(t, 150-dividerThickness(), sizeB.Height)
		})
	})
}

func TestSplitContainer_swap_contents(t *testing.T) {
	dl := dividerLength()
	dt := dividerThickness()
	initialWidth := 10 + 10 + dt
	initialHeight := gui.Max(10, dl)
	expectedWidth := 100 + 10 + dt
	expectedHeight := gui.Max(100, dl)

	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objA.SetMinSize(gui.NewSize(10, 10))
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB.SetMinSize(gui.NewSize(10, 10))
	objC := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objC.SetMinSize(gui.NewSize(100, 100))
	t.Run("Leading", func(t *testing.T) {
		sc := NewHSplit(objA, objB)
		min := sc.MinSize()
		assert.Equal(t, float32(initialWidth), min.Width)
		assert.Equal(t, float32(initialHeight), min.Height)
		sc.Leading = objC
		sc.Refresh()
		min = sc.MinSize()
		assert.Equal(t, float32(expectedWidth), min.Width)
		assert.Equal(t, float32(expectedHeight), min.Height)
	})
	t.Run("Trailing", func(t *testing.T) {
		sc := NewHSplit(objA, objB)
		min := sc.MinSize()
		assert.Equal(t, float32(initialWidth), min.Width)
		assert.Equal(t, float32(initialHeight), min.Height)
		sc.Trailing = objC
		sc.Refresh()
		min = sc.MinSize()
		assert.Equal(t, float32(expectedWidth), min.Width)
		assert.Equal(t, float32(expectedHeight), min.Height)
	})
}

func TestSplitContainer_divider_cursor(t *testing.T) {
	t.Run("Horizontal", func(t *testing.T) {
		divider := newDivider(&Split{Horizontal: true})
		assert.Equal(t, desktop.HResizeCursor, divider.Cursor())
	})
	t.Run("Vertical", func(t *testing.T) {
		divider := newDivider(&Split{Horizontal: false})
		assert.Equal(t, desktop.VResizeCursor, divider.Cursor())
	})
}

func TestSplitContainer_divider_drag(t *testing.T) {
	size := gui.NewSize(10, 10)
	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objA.SetMinSize(size)
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB.SetMinSize(size)
	t.Run("Horizontal", func(t *testing.T) {
		split := NewHSplit(objA, objB)
		split.Resize(gui.NewSize(100, 100))
		divider := newDivider(split)
		assert.Equal(t, 0.5, split.Offset)

		divider.Dragged(&gui.DragEvent{
			PointEvent: gui.PointEvent{Position: gui.NewPos(20, 9)},
			Dragged:    gui.NewDelta(10, -1),
		})
		assert.Equal(t, 0.6, split.Offset)

		divider.DragEnd()
		assert.Equal(t, 0.6, split.Offset)
	})
	t.Run("Vertical", func(t *testing.T) {
		split := NewVSplit(objA, objB)
		split.Resize(gui.NewSize(100, 100))
		divider := newDivider(split)
		assert.Equal(t, 0.5, split.Offset)

		divider.Dragged(&gui.DragEvent{
			PointEvent: gui.PointEvent{Position: gui.NewPos(9, 20)},
			Dragged:    gui.NewDelta(-1, 10),
		})
		assert.Equal(t, 0.6, split.Offset)

		divider.DragEnd()
		assert.Equal(t, 0.6, split.Offset)
	})
}

func TestSplitContainer_divider_drag_StartOffsetLessThanMinSize(t *testing.T) {
	size := gui.NewSize(30, 30)
	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objA.SetMinSize(size)
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB.SetMinSize(size)
	t.Run("Horizontal", func(t *testing.T) {
		split := NewHSplit(objA, objB)
		split.Resize(gui.NewSize(100, 100))
		divider := newDivider(split)
		t.Run("Leading", func(t *testing.T) {
			split.SetOffset(0.1)

			divider.Dragged(&gui.DragEvent{
				Dragged: gui.NewDelta(10, 0),
			})
			divider.DragEnd()

			assert.Equal(t, 0.4, split.Offset)
		})
		t.Run("Trailing", func(t *testing.T) {
			split.SetOffset(0.9)

			divider.Dragged(&gui.DragEvent{
				Dragged: gui.NewDelta(-10, 0),
			})
			divider.DragEnd()

			assert.Equal(t, 0.6, split.Offset)
		})
	})
	t.Run("Vertical", func(t *testing.T) {
		split := NewVSplit(objA, objB)
		split.Resize(gui.NewSize(100, 100))
		divider := newDivider(split)
		t.Run("Leading", func(t *testing.T) {
			split.SetOffset(0.1)

			divider.Dragged(&gui.DragEvent{
				Dragged: gui.NewDelta(0, 10),
			})
			divider.DragEnd()

			assert.Equal(t, 0.4, split.Offset)
		})
		t.Run("Trailing", func(t *testing.T) {
			split.SetOffset(0.9)

			divider.Dragged(&gui.DragEvent{
				Dragged: gui.NewDelta(0, -10),
			})
			divider.DragEnd()

			assert.Equal(t, 0.6, split.Offset)
		})
	})
}

func TestSplitContainer_divider_hover(t *testing.T) {
	t.Run("Horizontal", func(t *testing.T) {
		divider := newDivider(&Split{Horizontal: true})
		assert.False(t, divider.hovered)

		divider.MouseIn(&desktop.MouseEvent{})
		assert.True(t, divider.hovered)

		divider.MouseOut()
		assert.False(t, divider.hovered)
	})
	t.Run("Vertical", func(t *testing.T) {
		divider := newDivider(&Split{Horizontal: false})
		assert.False(t, divider.hovered)

		divider.MouseIn(&desktop.MouseEvent{})
		assert.True(t, divider.hovered)

		divider.MouseOut()
		assert.False(t, divider.hovered)
	})
}

func TestSplitContainer_divider_MinSize(t *testing.T) {
	t.Run("Horizontal", func(t *testing.T) {
		divider := newDivider(&Split{Horizontal: true})
		min := divider.MinSize()
		assert.Equal(t, dividerThickness(), min.Width)
		assert.Equal(t, dividerLength(), min.Height)
	})
	t.Run("Vertical", func(t *testing.T) {
		divider := newDivider(&Split{Horizontal: false})
		min := divider.MinSize()
		assert.Equal(t, dividerLength(), min.Width)
		assert.Equal(t, dividerThickness(), min.Height)
	})
}
