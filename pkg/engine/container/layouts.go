package container // import gui "github.com/bhojpur/gui/pkg/engine/container"

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
	"github.com/bhojpur/gui/pkg/engine/layout"
)

// NewAdaptiveGrid creates a new container with the specified objects and using the grid layout.
// When in a horizontal arrangement the rowcols parameter will specify the column count, when in vertical
// it will specify the rows. On mobile this will dynamically refresh when device is rotated.
//
// Since: 1.4
func NewAdaptiveGrid(rowcols int, objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewAdaptiveGridLayout(rowcols), objects...)
}

// NewBorder creates a new container with the specified objects and using the border layout.
// The top, bottom, left and right parameters specify the items that should be placed around edges,
// the remaining elements will be in the center. Nil can be used to an edge if it should not be filled.
//
// Since: 1.4
func NewBorder(top, bottom, left, right gui.CanvasObject, objects ...gui.CanvasObject) *gui.Container {
	all := objects
	if top != nil {
		all = append(all, top)
	}
	if bottom != nil {
		all = append(all, bottom)
	}
	if left != nil {
		all = append(all, left)
	}
	if right != nil {
		all = append(all, right)
	}
	return gui.NewContainerWithLayout(layout.NewBorderLayout(top, bottom, left, right), all...)
}

// NewCenter creates a new container with the specified objects centered in the available space.
//
// Since: 1.4
func NewCenter(objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewCenterLayout(), objects...)
}

// NewGridWithColumns creates a new container with the specified objects and using the grid layout with
// a specified number of columns. The number of rows will depend on how many children are in the container.
//
// Since: 1.4
func NewGridWithColumns(cols int, objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewGridLayoutWithColumns(cols), objects...)
}

// NewGridWithRows creates a new container with the specified objects and using the grid layout with
// a specified number of rows. The number of columns will depend on how many children are in the container.
//
// Since: 1.4
func NewGridWithRows(rows int, objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewGridLayoutWithRows(rows), objects...)
}

// NewGridWrap creates a new container with the specified objects and using the gridwrap layout.
// Every element will be resized to the size parameter and the content will arrange along a row and flow to a
// new row if the elements don't fit.
//
// Since: 1.4
func NewGridWrap(size gui.Size, objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewGridWrapLayout(size), objects...)
}

// NewHBox creates a new container with the specified objects and using the HBox layout.
// The objects will be placed in the container from left to right.
//
// Since: 1.4
func NewHBox(objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewHBoxLayout(), objects...)
}

// NewMax creates a new container with the specified objects filling the available space.
//
// Since: 1.4
func NewMax(objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewMaxLayout(), objects...)
}

// NewPadded creates a new container with the specified objects inset by standard padding size.
//
// Since: 1.4
func NewPadded(objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewPaddedLayout(), objects...)
}

// NewVBox creates a new container with the specified objects and using the VBox layout.
// The objects will be stacked in the container from top to bottom.
//
// Since: 1.4
func NewVBox(objects ...gui.CanvasObject) *gui.Container {
	return gui.NewContainerWithLayout(layout.NewVBoxLayout(), objects...)
}
