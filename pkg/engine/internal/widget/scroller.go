package widget

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
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// ScrollDirection represents the directions in which a Scroll can scroll its child content.
type ScrollDirection int

// Constants for valid values of ScrollDirection.
const (
	// ScrollBoth supports horizontal and vertical scrolling.
	ScrollBoth ScrollDirection = iota
	// ScrollHorizontalOnly specifies the scrolling should only happen left to right.
	ScrollHorizontalOnly
	// ScrollVerticalOnly specifies the scrolling should only happen top to bottom.
	ScrollVerticalOnly
	// ScrollNone turns off scrolling for this container.
	//
	// Since: 2.0
	ScrollNone
)

type scrollBarOrientation int

// We default to vertical as 0 due to that being the original orientation offered
const (
	scrollBarOrientationVertical   scrollBarOrientation = 0
	scrollBarOrientationHorizontal scrollBarOrientation = 1
	scrollContainerMinSize                              = float32(32) // TODO consider the smallest useful scroll view?
)

type scrollBarRenderer struct {
	BaseRenderer
	scrollBar  *scrollBar
	background *canvas.Rectangle
	minSize    gui.Size
}

func (r *scrollBarRenderer) Layout(size gui.Size) {
	r.background.Resize(size)
}

func (r *scrollBarRenderer) MinSize() gui.Size {
	return r.minSize
}

func (r *scrollBarRenderer) Refresh() {
	r.background.FillColor = theme.ScrollBarColor()
	r.background.Refresh()
}

var _ desktop.Hoverable = (*scrollBar)(nil)
var _ gui.Draggable = (*scrollBar)(nil)

type scrollBar struct {
	Base
	area            *scrollBarArea
	draggedDistance float32
	dragStart       float32
	isDragged       bool
	orientation     scrollBarOrientation
}

func (b *scrollBar) CreateRenderer() gui.WidgetRenderer {
	background := canvas.NewRectangle(theme.ScrollBarColor())
	r := &scrollBarRenderer{
		scrollBar:  b,
		background: background,
	}
	r.SetObjects([]gui.CanvasObject{background})
	return r
}

func (b *scrollBar) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (b *scrollBar) DragEnd() {
	b.isDragged = false
}

func (b *scrollBar) Dragged(e *gui.DragEvent) {
	if !b.isDragged {
		b.isDragged = true
		switch b.orientation {
		case scrollBarOrientationHorizontal:
			b.dragStart = b.Position().X
		case scrollBarOrientationVertical:
			b.dragStart = b.Position().Y
		}
		b.draggedDistance = 0
	}

	switch b.orientation {
	case scrollBarOrientationHorizontal:
		b.draggedDistance += e.Dragged.DX
	case scrollBarOrientationVertical:
		b.draggedDistance += e.Dragged.DY
	}
	b.area.moveBar(b.draggedDistance+b.dragStart, b.Size())
}

func (b *scrollBar) MouseIn(e *desktop.MouseEvent) {
	b.area.MouseIn(e)
}

func (b *scrollBar) MouseMoved(*desktop.MouseEvent) {
}

func (b *scrollBar) MouseOut() {
	b.area.MouseOut()
}

func newScrollBar(area *scrollBarArea) *scrollBar {
	b := &scrollBar{area: area, orientation: area.orientation}
	b.ExtendBaseWidget(b)
	return b
}

type scrollBarAreaRenderer struct {
	BaseRenderer
	area *scrollBarArea
	bar  *scrollBar
}

func (r *scrollBarAreaRenderer) Layout(_ gui.Size) {
	var barHeight, barWidth, barX, barY float32
	switch r.area.orientation {
	case scrollBarOrientationHorizontal:
		barWidth, barHeight, barX, barY = r.barSizeAndOffset(r.area.scroll.Offset.X, r.area.scroll.Content.Size().Width, r.area.scroll.Size().Width)
	default:
		barHeight, barWidth, barY, barX = r.barSizeAndOffset(r.area.scroll.Offset.Y, r.area.scroll.Content.Size().Height, r.area.scroll.Size().Height)
	}
	r.bar.Move(gui.NewPos(barX, barY))
	r.bar.Resize(gui.NewSize(barWidth, barHeight))
}

func (r *scrollBarAreaRenderer) MinSize() gui.Size {
	min := theme.ScrollBarSize()
	if !r.area.isLarge {
		min = theme.ScrollBarSmallSize() * 2
	}
	switch r.area.orientation {
	case scrollBarOrientationHorizontal:
		return gui.NewSize(theme.ScrollBarSize(), min)
	default:
		return gui.NewSize(min, theme.ScrollBarSize())
	}
}

func (r *scrollBarAreaRenderer) Refresh() {
	r.Layout(r.area.Size())
	canvas.Refresh(r.bar)
}

func (r *scrollBarAreaRenderer) barSizeAndOffset(contentOffset, contentLength, scrollLength float32) (length, width, lengthOffset, widthOffset float32) {
	if scrollLength < contentLength {
		portion := scrollLength / contentLength
		length = float32(int(scrollLength)) * portion
		if length < theme.ScrollBarSize() {
			length = theme.ScrollBarSize()
		}
	} else {
		length = scrollLength
	}
	if contentOffset != 0 {
		lengthOffset = (scrollLength - length) * (contentOffset / (contentLength - scrollLength))
	}
	if r.area.isLarge {
		width = theme.ScrollBarSize()
	} else {
		widthOffset = theme.ScrollBarSmallSize()
		width = theme.ScrollBarSmallSize()
	}
	return
}

var _ desktop.Hoverable = (*scrollBarArea)(nil)

type scrollBarArea struct {
	Base

	isLarge     bool
	scroll      *Scroll
	orientation scrollBarOrientation
}

func (a *scrollBarArea) CreateRenderer() gui.WidgetRenderer {
	bar := newScrollBar(a)
	return &scrollBarAreaRenderer{BaseRenderer: NewBaseRenderer([]gui.CanvasObject{bar}), area: a, bar: bar}
}

func (a *scrollBarArea) MouseIn(*desktop.MouseEvent) {
	a.isLarge = true
	a.scroll.Refresh()
}

func (a *scrollBarArea) MouseMoved(*desktop.MouseEvent) {
}

func (a *scrollBarArea) MouseOut() {
	a.isLarge = false
	a.scroll.Refresh()
}

func (a *scrollBarArea) moveBar(offset float32, barSize gui.Size) {
	switch a.orientation {
	case scrollBarOrientationHorizontal:
		a.scroll.Offset.X = a.computeScrollOffset(barSize.Width, offset, a.scroll.Size().Width, a.scroll.Content.Size().Width)
	default:
		a.scroll.Offset.Y = a.computeScrollOffset(barSize.Height, offset, a.scroll.Size().Height, a.scroll.Content.Size().Height)
	}
	if f := a.scroll.OnScrolled; f != nil {
		f(a.scroll.Offset)
	}
	a.scroll.refreshWithoutOffsetUpdate()
}

func (a *scrollBarArea) computeScrollOffset(length, offset, scrollLength, contentLength float32) float32 {
	maxOffset := scrollLength - length
	if offset < 0 {
		offset = 0
	} else if offset > maxOffset {
		offset = maxOffset
	}
	ratio := offset / maxOffset
	scrollOffset := ratio * (contentLength - scrollLength)
	return scrollOffset
}

func newScrollBarArea(scroll *Scroll, orientation scrollBarOrientation) *scrollBarArea {
	a := &scrollBarArea{scroll: scroll, orientation: orientation}
	a.ExtendBaseWidget(a)
	return a
}

type scrollContainerRenderer struct {
	BaseRenderer
	scroll                  *Scroll
	vertArea                *scrollBarArea
	horizArea               *scrollBarArea
	leftShadow, rightShadow *Shadow
	topShadow, bottomShadow *Shadow
	oldMinSize              gui.Size
}

func (r *scrollContainerRenderer) layoutBars(size gui.Size) {
	if r.scroll.Direction == ScrollVerticalOnly || r.scroll.Direction == ScrollBoth {
		r.vertArea.Resize(gui.NewSize(r.vertArea.MinSize().Width, size.Height))
		r.vertArea.Move(gui.NewPos(r.scroll.Size().Width-r.vertArea.Size().Width, 0))
		r.topShadow.Resize(gui.NewSize(size.Width, 0))
		r.bottomShadow.Resize(gui.NewSize(size.Width, 0))
		r.bottomShadow.Move(gui.NewPos(0, r.scroll.size.Height))
	}

	if r.scroll.Direction == ScrollHorizontalOnly || r.scroll.Direction == ScrollBoth {
		r.horizArea.Resize(gui.NewSize(size.Width, r.horizArea.MinSize().Height))
		r.horizArea.Move(gui.NewPos(0, r.scroll.Size().Height-r.horizArea.Size().Height))
		r.leftShadow.Resize(gui.NewSize(0, size.Height))
		r.rightShadow.Resize(gui.NewSize(0, size.Height))
		r.rightShadow.Move(gui.NewPos(r.scroll.size.Width, 0))
	}

	r.updatePosition()
}

func (r *scrollContainerRenderer) Layout(size gui.Size) {
	c := r.scroll.Content
	c.Resize(c.MinSize().Max(size))

	r.layoutBars(size)
}

func (r *scrollContainerRenderer) MinSize() gui.Size {
	return r.scroll.MinSize()
}

func (r *scrollContainerRenderer) Refresh() {
	if len(r.BaseRenderer.Objects()) == 0 || r.BaseRenderer.Objects()[0] != r.scroll.Content {
		// push updated content object to baseRenderer
		r.BaseRenderer.Objects()[0] = r.scroll.Content
	}
	if r.oldMinSize == r.scroll.Content.MinSize() && r.oldMinSize == r.scroll.Content.Size() &&
		(r.scroll.Size().Width <= r.oldMinSize.Width && r.scroll.Size().Height <= r.oldMinSize.Height) {
		r.layoutBars(r.scroll.Size())
		return
	}

	r.oldMinSize = r.scroll.Content.MinSize()
	r.Layout(r.scroll.Size())
}

func (r *scrollContainerRenderer) handleAreaVisibility(contentSize, scrollSize float32, area *scrollBarArea) {
	if contentSize <= scrollSize {
		area.Hide()
	} else if r.scroll.Visible() {
		area.Show()
	}
}

func (r *scrollContainerRenderer) handleShadowVisibility(offset, contentSize, scrollSize float32, shadowStart gui.CanvasObject, shadowEnd gui.CanvasObject) {
	if !r.scroll.Visible() {
		return
	}
	if offset > 0 {
		shadowStart.Show()
	} else {
		shadowStart.Hide()
	}
	if offset < contentSize-scrollSize {
		shadowEnd.Show()
	} else {
		shadowEnd.Hide()
	}
}

func (r *scrollContainerRenderer) updatePosition() {
	if r.scroll.Content == nil {
		return
	}
	scrollSize := r.scroll.Size()
	contentSize := r.scroll.Content.Size()

	r.scroll.Content.Move(gui.NewPos(-r.scroll.Offset.X, -r.scroll.Offset.Y))

	if r.scroll.Direction == ScrollVerticalOnly || r.scroll.Direction == ScrollBoth {
		r.handleAreaVisibility(contentSize.Height, scrollSize.Height, r.vertArea)
		r.handleShadowVisibility(r.scroll.Offset.Y, contentSize.Height, scrollSize.Height, r.topShadow, r.bottomShadow)
		cache.Renderer(r.vertArea).Layout(r.scroll.size)
	} else {
		r.vertArea.Hide()
		r.topShadow.Hide()
		r.bottomShadow.Hide()
	}
	if r.scroll.Direction == ScrollHorizontalOnly || r.scroll.Direction == ScrollBoth {
		r.handleAreaVisibility(contentSize.Width, scrollSize.Width, r.horizArea)
		r.handleShadowVisibility(r.scroll.Offset.X, contentSize.Width, scrollSize.Width, r.leftShadow, r.rightShadow)
		cache.Renderer(r.horizArea).Layout(r.scroll.size)
	} else {
		r.horizArea.Hide()
		r.leftShadow.Hide()
		r.rightShadow.Hide()
	}

	if r.scroll.Direction != ScrollHorizontalOnly {
		canvas.Refresh(r.vertArea) // this is required to force the canvas to update, we have no "Redraw()"
	} else {
		canvas.Refresh(r.horizArea) // this is required like above but if we are horizontal
	}
}

// Scroll defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
type Scroll struct {
	Base
	minSize   gui.Size
	Direction ScrollDirection
	Content   gui.CanvasObject
	Offset    gui.Position
	// OnScrolled can be set to be notified when the Scroll has changed position.
	// You should not update the Scroll.Offset from this method.
	//
	// Since: 2.0
	OnScrolled func(gui.Position)
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (s *Scroll) CreateRenderer() gui.WidgetRenderer {
	scr := &scrollContainerRenderer{
		BaseRenderer: NewBaseRenderer([]gui.CanvasObject{s.Content}),
		scroll:       s,
	}
	scr.vertArea = newScrollBarArea(s, scrollBarOrientationVertical)
	scr.topShadow = NewShadow(ShadowBottom, SubmergedContentLevel)
	scr.bottomShadow = NewShadow(ShadowTop, SubmergedContentLevel)
	scr.horizArea = newScrollBarArea(s, scrollBarOrientationHorizontal)
	scr.leftShadow = NewShadow(ShadowRight, SubmergedContentLevel)
	scr.rightShadow = NewShadow(ShadowLeft, SubmergedContentLevel)
	scr.SetObjects(append(scr.Objects(), scr.vertArea, scr.topShadow, scr.bottomShadow, scr.horizArea,
		scr.leftShadow, scr.rightShadow))
	scr.updatePosition()

	return scr
}

//ScrollToBottom will scroll content to container bottom - to show latest info which end user just added
func (s *Scroll) ScrollToBottom() {
	s.Offset.Y = s.Content.MinSize().Height - s.Size().Height
	s.Refresh()
}

//ScrollToTop will scroll content to container top
func (s *Scroll) ScrollToTop() {
	s.Offset.Y = 0
	s.Refresh()
}

// DragEnd will stop scrolling on mobile has stopped
func (s *Scroll) DragEnd() {
}

// Dragged will scroll on any drag - bar or otherwise - for mobile
func (s *Scroll) Dragged(e *gui.DragEvent) {
	if !gui.CurrentDevice().IsMobile() {
		return
	}

	if s.updateOffset(e.Dragged.DX, e.Dragged.DY) {
		s.refreshWithoutOffsetUpdate()
	}
}

// MinSize returns the smallest size this widget can shrink to
func (s *Scroll) MinSize() gui.Size {
	min := gui.NewSize(scrollContainerMinSize, scrollContainerMinSize).Max(s.minSize)
	switch s.Direction {
	case ScrollHorizontalOnly:
		min.Height = gui.Max(min.Height, s.Content.MinSize().Height)
	case ScrollVerticalOnly:
		min.Width = gui.Max(min.Width, s.Content.MinSize().Width)
	case ScrollNone:
		return s.Content.MinSize()
	}
	return min
}

// SetMinSize specifies a minimum size for this scroll container.
// If the specified size is larger than the content size then scrolling will not be enabled
// This can be helpful to appear larger than default if the layout is collapsing this widget.
func (s *Scroll) SetMinSize(size gui.Size) {
	s.minSize = size
}

// Refresh causes this widget to be redrawn in it's current state
func (s *Scroll) Refresh() {
	s.updateOffset(0, 0)
	s.refreshWithoutOffsetUpdate()
}

// Resize is called when this scroller should change size. We refresh to ensure the scroll bars are updated.
func (s *Scroll) Resize(sz gui.Size) {
	if sz == s.size {
		return
	}

	s.Base.Resize(sz)
	s.Refresh()
}

func (s *Scroll) refreshWithoutOffsetUpdate() {
	s.Base.Refresh()
}

// Scrolled is called when an input device triggers a scroll event
func (s *Scroll) Scrolled(ev *gui.ScrollEvent) {
	dx, dy := ev.Scrolled.DX, ev.Scrolled.DY
	if s.Size().Width < s.Content.MinSize().Width && s.Size().Height >= s.Content.MinSize().Height && dx == 0 {
		dx, dy = dy, dx
	}
	if s.updateOffset(dx, dy) {
		s.refreshWithoutOffsetUpdate()
	}
}

func (s *Scroll) updateOffset(deltaX, deltaY float32) bool {
	if s.Content.Size().Width <= s.Size().Width && s.Content.Size().Height <= s.Size().Height {
		if s.Offset.X != 0 || s.Offset.Y != 0 {
			s.Offset.X = 0
			s.Offset.Y = 0
			return true
		}
		return false
	}
	s.Offset.X = computeOffset(s.Offset.X, -deltaX, s.Size().Width, s.Content.MinSize().Width)
	s.Offset.Y = computeOffset(s.Offset.Y, -deltaY, s.Size().Height, s.Content.MinSize().Height)
	if f := s.OnScrolled; f != nil {
		f(s.Offset)
	}
	return true
}

func computeOffset(start, delta, outerWidth, innerWidth float32) float32 {
	offset := start + delta
	if offset+outerWidth >= innerWidth {
		offset = innerWidth - outerWidth
	}
	if offset < 0 {
		offset = 0
	}
	return offset
}

// NewScroll creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed object.
func NewScroll(content gui.CanvasObject) *Scroll {
	s := newScrollContainerWithDirection(ScrollBoth, content)
	s.ExtendBaseWidget(s)
	return s
}

// NewHScroll create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Width to be smaller than that of the passed object.
func NewHScroll(content gui.CanvasObject) *Scroll {
	s := newScrollContainerWithDirection(ScrollHorizontalOnly, content)
	s.ExtendBaseWidget(s)
	return s
}

// NewVScroll create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Height to be smaller than that of the passed object.
func NewVScroll(content gui.CanvasObject) *Scroll {
	s := newScrollContainerWithDirection(ScrollVerticalOnly, content)
	s.ExtendBaseWidget(s)
	return s
}

func newScrollContainerWithDirection(direction ScrollDirection, content gui.CanvasObject) *Scroll {
	s := &Scroll{
		Direction: direction,
		Content:   content,
	}
	s.ExtendBaseWidget(s)
	return s
}
