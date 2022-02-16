package driver

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
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
)

// AbsolutePositionForObject returns the absolute position of an object in a set of object trees.
// If the object is not part of any of the trees, the position (0,0) is returned.
func AbsolutePositionForObject(object gui.CanvasObject, trees []gui.CanvasObject) gui.Position {
	var pos gui.Position
	findPos := func(o gui.CanvasObject, p gui.Position, _ gui.Position, _ gui.Size) bool {
		if o == object {
			pos = p
			return true
		}
		return false
	}
	for _, tree := range trees {
		if WalkVisibleObjectTree(tree, findPos, nil) {
			break
		}
	}
	return pos
}

// FindObjectAtPositionMatching is used to find an object in a canvas at the specified position.
// The matches function determines of the type of object that is found at this position is of a suitable type.
// The various canvas roots and overlays that can be searched are also passed in.
func FindObjectAtPositionMatching(mouse gui.Position, matches func(object gui.CanvasObject) bool, overlay gui.CanvasObject, roots ...gui.CanvasObject) (gui.CanvasObject, gui.Position, int) {
	var found gui.CanvasObject
	var foundPos gui.Position

	findFunc := func(walked gui.CanvasObject, pos gui.Position, clipPos gui.Position, clipSize gui.Size) bool {
		if !walked.Visible() {
			return false
		}

		if mouse.X < clipPos.X || mouse.Y < clipPos.Y {
			return false
		}

		if mouse.X >= clipPos.X+clipSize.Width || mouse.Y >= clipPos.Y+clipSize.Height {
			return false
		}

		if mouse.X < pos.X || mouse.Y < pos.Y {
			return false
		}

		if mouse.X >= pos.X+walked.Size().Width || mouse.Y >= pos.Y+walked.Size().Height {
			return false
		}

		if matches(walked) {
			found = walked
			foundPos = gui.NewPos(mouse.X-pos.X, mouse.Y-pos.Y)
		}
		return false
	}

	layer := 0
	if overlay != nil {
		WalkVisibleObjectTree(overlay, findFunc, nil)
	} else {
		for _, root := range roots {
			layer++
			if root == nil {
				continue
			}
			WalkVisibleObjectTree(root, findFunc, nil)
			if found != nil {
				break
			}
		}
	}

	return found, foundPos, layer
}

// ReverseWalkVisibleObjectTree will walk an object tree in reverse order for all visible objects
// executing the passed functions following the following rules:
// - beforeChildren is called for the start obj before traversing its children
// - the obj's children are traversed by calling walkObjects on each of the visible items
// - afterChildren is called for the obj after traversing the obj's children
// The walk can be aborted by returning true in one of the functions:
// - if beforeChildren returns true, further traversing is stopped immediately, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
func ReverseWalkVisibleObjectTree(
	obj gui.CanvasObject,
	beforeChildren func(gui.CanvasObject, gui.Position, gui.Position, gui.Size) bool,
	afterChildren func(gui.CanvasObject, gui.CanvasObject),
) bool {
	clipSize := gui.NewSize(math.MaxInt32, math.MaxInt32)
	return walkObjectTree(obj, true, nil, gui.NewPos(0, 0), gui.NewPos(0, 0), clipSize, beforeChildren, afterChildren, true)
}

// WalkCompleteObjectTree will walk an object tree for all objects (ignoring visible state) executing the passed
// functions following the following rules:
// - beforeChildren is called for the start obj before traversing its children
// - the obj's children are traversed by calling walkObjects on each of the items
// - afterChildren is called for the obj after traversing the obj's children
// The walk can be aborted by returning true in one of the functions:
// - if beforeChildren returns true, further traversing is stopped immediately, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
func WalkCompleteObjectTree(
	obj gui.CanvasObject,
	beforeChildren func(gui.CanvasObject, gui.Position, gui.Position, gui.Size) bool,
	afterChildren func(gui.CanvasObject, gui.CanvasObject),
) bool {
	clipSize := gui.NewSize(math.MaxInt32, math.MaxInt32)
	return walkObjectTree(obj, false, nil, gui.NewPos(0, 0), gui.NewPos(0, 0), clipSize, beforeChildren, afterChildren, false)
}

// WalkVisibleObjectTree will walk an object tree for all visible objects executing the passed functions following
// the following rules:
// - beforeChildren is called for the start obj before traversing its children
// - the obj's children are traversed by calling walkObjects on each of the visible items
// - afterChildren is called for the obj after traversing the obj's children
// The walk can be aborted by returning true in one of the functions:
// - if beforeChildren returns true, further traversing is stopped immediately, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
func WalkVisibleObjectTree(
	obj gui.CanvasObject,
	beforeChildren func(gui.CanvasObject, gui.Position, gui.Position, gui.Size) bool,
	afterChildren func(gui.CanvasObject, gui.CanvasObject),
) bool {
	clipSize := gui.NewSize(math.MaxInt32, math.MaxInt32)
	return walkObjectTree(obj, false, nil, gui.NewPos(0, 0), gui.NewPos(0, 0), clipSize, beforeChildren, afterChildren, true)
}

func walkObjectTree(
	obj gui.CanvasObject,
	reverse bool,
	parent gui.CanvasObject,
	offset, clipPos gui.Position,
	clipSize gui.Size,
	beforeChildren func(gui.CanvasObject, gui.Position, gui.Position, gui.Size) bool,
	afterChildren func(gui.CanvasObject, gui.CanvasObject),
	requireVisible bool,
) bool {
	if obj == nil {
		return false
	}
	if requireVisible && !obj.Visible() {
		return false
	}
	pos := obj.Position().Add(offset)

	var children []gui.CanvasObject
	switch co := obj.(type) {
	case *gui.Container:
		children = co.Objects
	case gui.Widget:
		children = cache.Renderer(co).Objects()
	}

	if _, ok := obj.(gui.Scrollable); ok {
		clipPos = pos
		clipSize = obj.Size()
	}

	if beforeChildren != nil {
		if beforeChildren(obj, pos, clipPos, clipSize) {
			return true
		}
	}

	cancelled := false
	followChild := func(child gui.CanvasObject) bool {
		if walkObjectTree(child, reverse, obj, pos, clipPos, clipSize, beforeChildren, afterChildren, requireVisible) {
			cancelled = true
			return true
		}
		return false
	}
	if reverse {
		for i := len(children) - 1; i >= 0; i-- {
			if followChild(children[i]) {
				break
			}
		}
	} else {
		for _, child := range children {
			if followChild(child) {
				break
			}
		}
	}

	if afterChildren != nil {
		afterChildren(obj, parent)
	}
	return cancelled
}
