package internal

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

import gui "github.com/bhojpur/gui/pkg/engine"

// ClipStack keeps track of the areas that should be clipped when drawing a canvas.
// If no clips are present then adding one will be added as-is.
// Subsequent items pushed will be completely within the previous clip.
type ClipStack struct {
	clips []*ClipItem
}

// Pop removes the current top clip and returns it.
func (c *ClipStack) Pop() *ClipItem {
	if len(c.clips) == 0 {
		return nil
	}

	ret := c.clips[len(c.clips)-1]
	c.clips = c.clips[:len(c.clips)-1]
	return ret
}

// Length  returns the number of items in this clip stack. 0 means no clip.
func (c *ClipStack) Length() int {
	return len(c.clips)
}

// Push a new clip onto this stack at position and size specified.
// The returned clip item is the result of calculating the intersection of the requested clip and it's parent.
func (c *ClipStack) Push(p gui.Position, s gui.Size) *ClipItem {
	outer := c.Top()
	inner := outer.Intersect(p, s)

	c.clips = append(c.clips, inner)
	return inner
}

// Top returns the current clip item - it will always be within the bounds of any parent clips.
func (c *ClipStack) Top() *ClipItem {
	if len(c.clips) == 0 {
		return nil
	}

	return c.clips[len(c.clips)-1]
}

// ClipItem represents a single clip in a clip stack, denoted by a size and position.
type ClipItem struct {
	pos  gui.Position
	size gui.Size
}

// Rect returns the position and size parameters of the clip.
func (i *ClipItem) Rect() (gui.Position, gui.Size) {
	return i.pos, i.size
}

// Intersect returns a new clip item that is the intersection of the requested parameters and this clip.
func (i *ClipItem) Intersect(p gui.Position, s gui.Size) *ClipItem {
	ret := &ClipItem{p, s}
	if i == nil {
		return ret
	}

	if ret.pos.X < i.pos.X {
		ret.pos.X = i.pos.X
		ret.size.Width -= i.pos.X - p.X
	}
	if ret.pos.Y < i.pos.Y {
		ret.pos.Y = i.pos.Y
		ret.size.Height -= i.pos.Y - p.Y
	}

	if p.X+s.Width > i.pos.X+i.size.Width {
		ret.size.Width = (i.pos.X + i.size.Width) - ret.pos.X
	}
	if p.Y+s.Height > i.pos.Y+i.size.Height {
		ret.size.Height = (i.pos.Y + i.size.Height) - ret.pos.Y
	}
	return ret
}
