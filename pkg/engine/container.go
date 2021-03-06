package engine

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

// Declare conformity to CanvasObject
var _ CanvasObject = (*Container)(nil)

// Container is a CanvasObject that contains a collection of child objects.
// The layout of the children is set by the specified Layout.
type Container struct {
	size     Size     // The current size of the Container
	position Position // The current position of the Container
	Hidden   bool     // Is this Container hidden

	Layout  Layout         // The Layout algorithm for arranging child CanvasObjects
	Objects []CanvasObject // The set of CanvasObjects this container holds
}

// NewContainer returns a new Container instance holding the specified CanvasObjects.
func NewContainer(objects ...CanvasObject) *Container {
	return NewContainerWithoutLayout(objects...)
}

// NewContainerWithoutLayout returns a new Container instance holding the specified
// CanvasObjects that are manually arranged.
func NewContainerWithoutLayout(objects ...CanvasObject) *Container {
	ret := &Container{
		Objects: objects,
	}

	ret.size = ret.MinSize()
	return ret
}

// NewContainerWithLayout returns a new Container instance holding the specified
// CanvasObjects which will be laid out according to the specified Layout.
func NewContainerWithLayout(layout Layout, objects ...CanvasObject) *Container {
	ret := &Container{
		Objects: objects,
		Layout:  layout,
	}

	ret.size = layout.MinSize(objects)
	ret.layout()
	return ret
}

// Add appends the specified object to the items this container manages.
//
// Since: 1.4
func (c *Container) Add(add CanvasObject) {
	c.Objects = append(c.Objects, add)
	c.layout()
}

// AddObject adds another CanvasObject to the set this Container holds.
func (c *Container) AddObject(o CanvasObject) {
	c.Add(o)
}

// MinSize calculates the minimum size of a Container.
// This is delegated to the Layout, if specified, otherwise it will mimic MaxLayout.
func (c *Container) MinSize() Size {
	if c.Layout != nil {
		return c.Layout.MinSize(c.Objects)
	}

	minSize := NewSize(1, 1)
	for _, child := range c.Objects {
		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

// Move the container (and all its children) to a new position, relative to its parent.
func (c *Container) Move(pos Position) {
	c.position = pos
}

// Position gets the current position of this Container, relative to its parent.
func (c *Container) Position() Position {
	return c.position
}

// Resize sets a new size for the Container.
func (c *Container) Resize(size Size) {
	if c.size == size {
		return
	}

	c.size = size
	c.layout()
}

// Size returns the current size of this container.
func (c *Container) Size() Size {
	return c.size
}

// Show sets this container, and all its children, to be visible.
func (c *Container) Show() {
	if !c.Hidden {
		return
	}

	c.Hidden = false
}

// Hide sets this container, and all its children, to be not visible.
func (c *Container) Hide() {
	if c.Hidden {
		return
	}

	c.Hidden = true
}

// Refresh causes this object to be redrawn in it's current state
func (c *Container) Refresh() {
	c.layout()

	for _, child := range c.Objects {
		child.Refresh()
	}

	// this is basically just canvas.Refresh(c) without the package loop
	o := CurrentApp().Driver().CanvasForObject(c)
	if o == nil {
		return
	}
	o.Refresh(c)
}

// Remove updates the contents of this container to no longer include the specified object.
func (c *Container) Remove(rem CanvasObject) {
	if len(c.Objects) == 0 {
		return
	}

	for i, o := range c.Objects {
		if o != rem {
			continue
		}

		copy(c.Objects[i:], c.Objects[i+1:])
		c.Objects[len(c.Objects)-1] = nil
		c.Objects = c.Objects[:len(c.Objects)-1]

		c.layout()
		return
	}
}

// Visible returns true if the container is currently visible, false otherwise.
func (c *Container) Visible() bool {
	return !c.Hidden
}

func (c *Container) layout() {
	if c.Layout != nil {
		c.Layout.Layout(c.Objects, c.size)
		return
	}
}
