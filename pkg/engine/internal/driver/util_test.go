package driver_test

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
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	internal_widget "github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/layout"
	_ "github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
)

func TestAbsolutePositionForObject(t *testing.T) {
	t1r1c1 := widget.NewLabel("row 1 col 1")
	t1r1c2 := widget.NewLabel("row 1 col 2")
	t1r2c1 := widget.NewLabel("row 2 col 1")
	t1r2c2 := widget.NewLabel("row 2 col 2")
	t1r2c2.Hide()
	t1r1 := gui.NewContainer(t1r1c1, t1r1c2)
	t1r2 := gui.NewContainer(t1r2c1, t1r2c2)
	tree1 := gui.NewContainer(t1r1, t1r2)

	t1r1c1.Move(gui.NewPos(111, 111))
	t1r1c2.Move(gui.NewPos(112, 112))
	t1r2c1.Move(gui.NewPos(121, 121))
	t1r2c2.Move(gui.NewPos(122, 122))
	t1r1.Move(gui.NewPos(11, 11))
	t1r2.Move(gui.NewPos(12, 12))
	tree1.Move(gui.NewPos(1, 1))

	t2r1c1 := widget.NewLabel("row 1 col 1")
	t2r1c2 := widget.NewLabel("row 1 col 2")
	t2r2c1 := widget.NewLabel("row 2 col 1")
	t2r2c2 := widget.NewLabel("row 2 col 2")
	t2r1 := gui.NewContainer(t2r1c1, t2r1c2)
	t2r2 := gui.NewContainer(t2r2c1, t2r2c2)
	tree2 := gui.NewContainer(t2r1, t2r2)

	t2r1c1.Move(gui.NewPos(211, 211))
	t2r1c2.Move(gui.NewPos(212, 212))
	t2r2c1.Move(gui.NewPos(221, 221))
	t2r2c2.Move(gui.NewPos(222, 222))
	t2r1.Move(gui.NewPos(21, 21))
	t2r2.Move(gui.NewPos(22, 22))
	tree2.Move(gui.NewPos(2, 2))

	t3r1 := widget.NewLabel("row 1")
	t3r2 := widget.NewLabel("row 2")
	tree3 := gui.NewContainer(t3r1, t3r2)

	t3r1.Move(gui.NewPos(31, 31))
	t3r2.Move(gui.NewPos(32, 32))
	tree3.Move(gui.NewPos(3, 3))

	trees := []gui.CanvasObject{tree1, tree2, tree3}

	outsideTrees := widget.NewLabel("outside trees")
	outsideTrees.Move(gui.NewPos(10, 10))

	tests := map[string]struct {
		object gui.CanvasObject
		want   gui.Position
	}{
		"tree 1: a cell": {
			object: t1r1c2,
			want:   gui.NewPos(124, 124), // 1 (root) + 11 (row 1) + 112 (cell 2)
		},
		"tree 1: a row": {
			object: t1r2,
			want:   gui.NewPos(13, 13), // 1 (root) + 12 (row 2)
		},
		"tree 1: root": {
			object: tree1,
			want:   gui.NewPos(1, 1),
		},
		"tree 1: a hidden element": {
			object: t1r2c2,
			want:   gui.NewPos(0, 0),
		},

		"tree 2: a row": {
			object: t2r2,
			want:   gui.NewPos(24, 24), // 2 (root) + 22 (row 2)
		},
		"tree 2: root": {
			object: tree2,
			want:   gui.NewPos(2, 2),
		},

		"tree 3: a row": {
			object: t3r2,
			want:   gui.NewPos(35, 35), // 3 (root) + 32 (row 2)
		},
		"tree 3: root": {
			object: tree3,
			want:   gui.NewPos(3, 3),
		},

		"an object not inside any tree": {
			object: outsideTrees,
			want:   gui.NewPos(0, 0),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, driver.AbsolutePositionForObject(tt.object, trees))
		})
	}
}

func TestFindObjectAtPositionMatching(t *testing.T) {
	col1cell1 := &objectTree{
		pos:  gui.NewPos(10, 10),
		size: gui.NewSize(15, 15),
	}
	col1cell2 := &objectTree{
		pos:  gui.NewPos(10, 35),
		size: gui.NewSize(15, 15),
	}
	col1cell3 := &objectTree{
		pos:  gui.NewPos(10, 60),
		size: gui.NewSize(15, 15),
	}
	col1 := &objectTree{
		children: []gui.CanvasObject{col1cell1, col1cell2, col1cell3},
		pos:      gui.NewPos(10, 10),
		size:     gui.NewSize(35, 80),
	}
	col2cell1 := &objectTree{
		pos:  gui.NewPos(10, 10),
		size: gui.NewSize(15, 15),
	}
	col2cell2 := &objectTree{
		pos:  gui.NewPos(10, 35),
		size: gui.NewSize(15, 15),
	}
	col2cell3 := &objectTree{
		pos:  gui.NewPos(10, 60),
		size: gui.NewSize(15, 15),
	}
	col2 := &objectTree{
		children: []gui.CanvasObject{col2cell1, col2cell2, col2cell3},
		pos:      gui.NewPos(55, 10),
		size:     gui.NewSize(35, 80),
	}
	colTree := &objectTree{
		children: []gui.CanvasObject{col1, col2},
		pos:      gui.NewPos(10, 10),
		size:     gui.NewSize(100, 100),
	}
	row1cell1 := &objectTree{
		pos:  gui.NewPos(10, 10),
		size: gui.NewSize(15, 15),
	}
	row1cell2 := &objectTree{
		pos:  gui.NewPos(35, 10),
		size: gui.NewSize(15, 15),
	}
	row1cell3 := &objectTree{
		pos:  gui.NewPos(60, 10),
		size: gui.NewSize(15, 15),
	}
	row1 := &objectTree{
		children: []gui.CanvasObject{row1cell1, row1cell2, row1cell3},
		pos:      gui.NewPos(10, 10),
		size:     gui.NewSize(80, 35),
	}
	row2cell1 := &objectTree{
		pos:  gui.NewPos(10, 10),
		size: gui.NewSize(15, 15),
	}
	row2cell2 := &objectTree{
		pos:  gui.NewPos(35, 10),
		size: gui.NewSize(15, 15),
	}
	row2cell3 := &objectTree{
		pos:  gui.NewPos(60, 10),
		size: gui.NewSize(15, 15),
	}
	row2 := &objectTree{
		children: []gui.CanvasObject{row2cell1, row2cell2, row2cell3},
		pos:      gui.NewPos(10, 55),
		size:     gui.NewSize(80, 35),
	}
	rowTree := &objectTree{
		children: []gui.CanvasObject{row1, row2},
		pos:      gui.NewPos(10, 10),
		size:     gui.NewSize(100, 100),
	}
	tree1 := &objectTree{
		pos:  gui.NewPos(100, 100),
		size: gui.NewSize(5, 5),
	}
	tree2 := &objectTree{
		pos:  gui.NewPos(0, 0),
		size: gui.NewSize(5, 5),
	}
	tree3 := &objectTree{
		pos:  gui.NewPos(50, 50),
		size: gui.NewSize(5, 5),
	}
	for name, tt := range map[string]struct {
		matcher    func(object gui.CanvasObject) bool
		overlay    gui.CanvasObject
		pos        gui.Position
		roots      []gui.CanvasObject
		wantObject gui.CanvasObject
		wantPos    gui.Position
		wantLayer  int
	}{
		"match in overlay and roots": {
			matcher:    func(o gui.CanvasObject) bool { return o.Size().Width == 15 },
			overlay:    colTree,
			pos:        gui.NewPos(35, 60),
			roots:      []gui.CanvasObject{rowTree},
			wantObject: col1cell2,
			wantPos:    gui.NewPos(5, 5),
			wantLayer:  0,
		},
		"match in root but overlay without match present": {
			matcher:    func(o gui.CanvasObject) bool { return o.Size().Width == 15 },
			overlay:    tree1,
			pos:        gui.NewPos(35, 60),
			roots:      []gui.CanvasObject{colTree, rowTree},
			wantObject: nil,
			wantPos:    gui.Position{},
			wantLayer:  0,
		},
		"match in multiple roots without overlay": {
			matcher:    func(o gui.CanvasObject) bool { return o.Size().Width == 15 },
			overlay:    nil,
			pos:        gui.NewPos(83, 83),
			roots:      []gui.CanvasObject{tree1, rowTree, tree2, colTree},
			wantObject: row2cell3,
			wantPos:    gui.NewPos(3, 8),
			wantLayer:  2,
		},
		"no match in roots without overlay": {
			matcher:    func(o gui.CanvasObject) bool { return true },
			overlay:    nil,
			pos:        gui.NewPos(66, 66),
			roots:      []gui.CanvasObject{tree1, tree2, tree3},
			wantObject: nil,
			wantPos:    gui.Position{},
			wantLayer:  3,
		},
		"no overlay and no roots": {
			matcher:    func(o gui.CanvasObject) bool { return true },
			overlay:    nil,
			pos:        gui.NewPos(66, 66),
			roots:      nil,
			wantObject: nil,
			wantPos:    gui.Position{},
			wantLayer:  0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			o, p, l := driver.FindObjectAtPositionMatching(tt.pos, tt.matcher, tt.overlay, tt.roots...)
			assert.Equal(t, tt.wantObject, o, "found object")
			assert.Equal(t, tt.wantPos, p, "position of found object")
			assert.Equal(t, tt.wantLayer, l, "layer of found object (0 - overlay, 1, 2, 3… - roots")
		})
	}
}

func TestReverseWalkVisibleObjectTree(t *testing.T) {
	child1 := canvas.NewRectangle(color.White)
	child1.SetMinSize(gui.NewSize(100, 100))
	child2 := canvas.NewRectangle(color.Black)
	child2.Hide()
	child3 := canvas.NewRectangle(color.White)
	base := container.NewHBox(child1, child2, child3)

	var walked []gui.CanvasObject
	driver.ReverseWalkVisibleObjectTree(
		base,
		func(object gui.CanvasObject, position gui.Position, clippingPos gui.Position, clippingSize gui.Size) bool {
			walked = append(walked, object)
			return false
		},
		nil,
	)

	assert.Equal(t, []gui.CanvasObject{base, child3, child1}, walked)
}

func TestReverseWalkVisibleObjectTree_Clip(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(gui.NewSize(100, 100))
	child := canvas.NewRectangle(color.Black)
	base := gui.NewContainerWithLayout(
		layout.NewGridLayout(1),
		rect,
		internal_widget.NewScroll(child),
		gui.NewContainerWithLayout(
			layout.NewGridLayout(2),
			canvas.NewCircle(color.White),
			canvas.NewCircle(color.White),
			canvas.NewCircle(color.White),
			&scrollable{},
		),
	)

	var scClipPos, scrollableClipPos gui.Position
	var scClipSize, scrollableClipSize gui.Size

	driver.ReverseWalkVisibleObjectTree(base, func(object gui.CanvasObject, position gui.Position, clippingPos gui.Position, clippingSize gui.Size) bool {
		if _, ok := object.(*internal_widget.Scroll); ok {
			scClipPos = clippingPos
			scClipSize = clippingSize
		} else if _, ok = object.(gui.Scrollable); ok {
			scrollableClipPos = clippingPos
			scrollableClipSize = clippingSize
		}
		return false
	}, nil)

	// layout:
	// +-------------------------------+
	// | 0,0: rect 100x100             |
	// +-------------------------------+
	// |            padding            |
	// +-------------------------------+
	// | 0,104: scroller 100x100       |
	// +-------------------------------+
	// |            padding            |
	// +--------------+-+--------------+
	// | circle 48x48 |p| circle 48x48 |
	// +--------------+-+--------------+
	// |            padding            |
	// +--------------+-+--------------+
	// | circle 48x48 |p| scrollable   |
	// +--------------+-+--------------+
	assert.Equal(t, gui.NewPos(0, 104), scClipPos)
	assert.Equal(t, gui.NewSize(100, 100), scClipSize)
	assert.Equal(t, gui.NewPos(52, 260), scrollableClipPos)
	assert.Equal(t, gui.NewSize(48, 48), scrollableClipSize)
}

func TestWalkVisibleObjectTree(t *testing.T) {
	child1 := canvas.NewRectangle(color.White)
	child1.SetMinSize(gui.NewSize(100, 100))
	child2 := canvas.NewRectangle(color.Black)
	child2.Hide()
	child3 := canvas.NewRectangle(color.White)
	base := container.NewHBox(child1, child2, child3)

	var walked []gui.CanvasObject
	driver.WalkVisibleObjectTree(base, func(object gui.CanvasObject, position gui.Position, clippingPos gui.Position, clippingSize gui.Size) bool {
		walked = append(walked, object)
		return false
	}, nil)

	assert.Equal(t, []gui.CanvasObject{base, child1, child3}, walked)
}

func TestWalkVisibleObjectTree_Clip(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(gui.NewSize(100, 100))
	child := canvas.NewRectangle(color.Black)
	base := gui.NewContainerWithLayout(
		layout.NewGridLayout(1),
		rect,
		internal_widget.NewScroll(child),
		gui.NewContainerWithLayout(
			layout.NewGridLayout(2),
			canvas.NewCircle(color.White),
			canvas.NewCircle(color.White),
			canvas.NewCircle(color.White),
			&scrollable{},
		),
	)

	var scClipPos, scrollableClipPos gui.Position
	var scClipSize, scrollableClipSize gui.Size

	driver.WalkVisibleObjectTree(base, func(object gui.CanvasObject, position gui.Position, clippingPos gui.Position, clippingSize gui.Size) bool {
		if _, ok := object.(*internal_widget.Scroll); ok {
			scClipPos = clippingPos
			scClipSize = clippingSize
		} else if _, ok = object.(gui.Scrollable); ok {
			scrollableClipPos = clippingPos
			scrollableClipSize = clippingSize
		}
		return false
	}, nil)

	// layout:
	// +-------------------------------+
	// | 0,0: rect 100x100             |
	// +-------------------------------+
	// |            padding            |
	// +-------------------------------+
	// | 0,104: scroller 100x100       |
	// +-------------------------------+
	// |            padding            |
	// +--------------+-+--------------+
	// | circle 48x48 |p| circle 48x48 |
	// +--------------+-+--------------+
	// |            padding            |
	// +--------------+-+--------------+
	// | circle 48x48 |p| scrollable   |
	// +--------------+-+--------------+
	assert.Equal(t, gui.NewPos(0, 104), scClipPos)
	assert.Equal(t, gui.NewSize(100, 100), scClipSize)
	assert.Equal(t, gui.NewPos(52, 260), scrollableClipPos)
	assert.Equal(t, gui.NewSize(48, 48), scrollableClipSize)
}

func TestWalkWholeObjectTree(t *testing.T) {
	child1 := canvas.NewRectangle(color.White)
	child1.SetMinSize(gui.NewSize(100, 100))
	child2 := canvas.NewRectangle(color.Black)
	child2.Hide()
	child3 := canvas.NewRectangle(color.White)
	base := container.NewHBox(child1, child2, child3)

	var walked []gui.CanvasObject
	driver.WalkCompleteObjectTree(base, func(object gui.CanvasObject, position gui.Position, clippingPos gui.Position, clippingSize gui.Size) bool {
		walked = append(walked, object)
		return false
	}, nil)

	assert.Equal(t, []gui.CanvasObject{base, child1, child2, child3}, walked)
}

var _ gui.Widget = (*objectTree)(nil)

type objectTree struct {
	children []gui.CanvasObject
	hidden   bool
	pos      gui.Position
	size     gui.Size
}

func (o *objectTree) Size() gui.Size {
	return o.size
}

func (o *objectTree) Resize(size gui.Size) {
	o.size = size
}

func (o *objectTree) Position() gui.Position {
	return o.pos
}

func (o *objectTree) Move(position gui.Position) {
	o.pos = position
}

func (o *objectTree) MinSize() gui.Size {
	return o.size
}

func (o objectTree) Visible() bool {
	return !o.hidden
}

func (o *objectTree) Show() {
	o.hidden = false
}

func (o *objectTree) Hide() {
	o.hidden = true
}

func (o *objectTree) Refresh() {
}

func (o *objectTree) CreateRenderer() gui.WidgetRenderer {
	r := &objectTreeRenderer{}
	r.SetObjects(o.children)
	return r
}

var _ gui.WidgetRenderer = (*objectTreeRenderer)(nil)

type objectTreeRenderer struct {
	internal_widget.BaseRenderer
}

func (o objectTreeRenderer) Layout(_ gui.Size) {
}

func (o objectTreeRenderer) MinSize() gui.Size {
	return gui.NewSize(0, 0)
}

func (o objectTreeRenderer) Refresh() {
}

type scrollable struct {
	pos  gui.Position
	size gui.Size
}

var _ gui.CanvasObject = (*scrollable)(nil)
var _ gui.Scrollable = (*scrollable)(nil)

func (s *scrollable) Hide() {
	panic("implement me")
}

func (s *scrollable) MinSize() gui.Size {
	return gui.NewSize(1, 1)
}

func (s *scrollable) Move(position gui.Position) {
	s.pos = position
}

func (s *scrollable) Position() gui.Position {
	return s.pos
}

func (s *scrollable) Refresh() {
	panic("implement me")
}

func (s *scrollable) Resize(size gui.Size) {
	s.size = size
}

func (s *scrollable) Scrolled(event *gui.ScrollEvent) {
	panic("implement me")
}

func (s *scrollable) Show() {
	panic("implement me")
}

func (s *scrollable) Size() gui.Size {
	return s.size
}

func (s *scrollable) Visible() bool {
	return true
}
