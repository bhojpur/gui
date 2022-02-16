package layout

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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// Declare conformity with Layout interface
var _ gui.Layout = (*borderLayout)(nil)

type borderLayout struct {
	top, bottom, left, right gui.CanvasObject
}

// NewBorderLayout creates a new BorderLayout instance with top, bottom, left
// and right objects set. All other items in the container will fill the centre
// space
func NewBorderLayout(top, bottom, left, right gui.CanvasObject) gui.Layout {
	return &borderLayout{top, bottom, left, right}
}

// Layout is called to pack all child objects into a specified size.
// For BorderLayout this arranges the top, bottom, left and right widgets at
// the sides and any remaining widgets are maximised in the middle space.
func (b *borderLayout) Layout(objects []gui.CanvasObject, size gui.Size) {
	var topSize, bottomSize, leftSize, rightSize gui.Size
	if b.top != nil && b.top.Visible() {
		b.top.Resize(gui.NewSize(size.Width, b.top.MinSize().Height))
		b.top.Move(gui.NewPos(0, 0))
		topSize = gui.NewSize(size.Width, b.top.MinSize().Height+theme.Padding())
	}
	if b.bottom != nil && b.bottom.Visible() {
		b.bottom.Resize(gui.NewSize(size.Width, b.bottom.MinSize().Height))
		b.bottom.Move(gui.NewPos(0, size.Height-b.bottom.MinSize().Height))
		bottomSize = gui.NewSize(size.Width, b.bottom.MinSize().Height+theme.Padding())
	}
	if b.left != nil && b.left.Visible() {
		b.left.Resize(gui.NewSize(b.left.MinSize().Width, size.Height-topSize.Height-bottomSize.Height))
		b.left.Move(gui.NewPos(0, topSize.Height))
		leftSize = gui.NewSize(b.left.MinSize().Width+theme.Padding(), size.Height-topSize.Height-bottomSize.Height)
	}
	if b.right != nil && b.right.Visible() {
		b.right.Resize(gui.NewSize(b.right.MinSize().Width, size.Height-topSize.Height-bottomSize.Height))
		b.right.Move(gui.NewPos(size.Width-b.right.MinSize().Width, topSize.Height))
		rightSize = gui.NewSize(b.right.MinSize().Width+theme.Padding(), size.Height-topSize.Height-bottomSize.Height)
	}

	middleSize := gui.NewSize(size.Width-leftSize.Width-rightSize.Width, size.Height-topSize.Height-bottomSize.Height)
	middlePos := gui.NewPos(leftSize.Width, topSize.Height)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if child != b.top && child != b.bottom && child != b.left && child != b.right {
			child.Resize(middleSize)
			child.Move(middlePos)
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For BorderLayout this is determined by the MinSize height of the top and
// plus the MinSize width of the left and right, plus any padding needed.
// This is then added to the union of the MinSize for any remaining content.
func (b *borderLayout) MinSize(objects []gui.CanvasObject) gui.Size {
	minSize := gui.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if child != b.top && child != b.bottom && child != b.left && child != b.right {
			minSize = minSize.Max(child.MinSize())
		}
	}

	if b.left != nil && b.left.Visible() {
		minHeight := gui.Max(minSize.Height, b.left.MinSize().Height)
		minSize = gui.NewSize(minSize.Width+b.left.MinSize().Width+theme.Padding(), minHeight)
	}
	if b.right != nil && b.right.Visible() {
		minHeight := gui.Max(minSize.Height, b.right.MinSize().Height)
		minSize = gui.NewSize(minSize.Width+b.right.MinSize().Width+theme.Padding(), minHeight)
	}

	if b.top != nil && b.top.Visible() {
		minWidth := gui.Max(minSize.Width, b.top.MinSize().Width)
		minSize = gui.NewSize(minWidth, minSize.Height+b.top.MinSize().Height+theme.Padding())
	}
	if b.bottom != nil && b.bottom.Visible() {
		minWidth := gui.Max(minSize.Width, b.bottom.MinSize().Width)
		minSize = gui.NewSize(minWidth, minSize.Height+b.bottom.MinSize().Height+theme.Padding())
	}

	return minSize
}
