package mobile

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
	"context"
	"image"
	"math"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/mobile"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/common"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

const (
	doubleClickDelay = 500 // ms (maximum interval between clicks for double click detection)
)

var _ gui.Canvas = (*mobileCanvas)(nil)

type mobileCanvas struct {
	common.Canvas

	content          gui.CanvasObject
	windowHead, menu gui.CanvasObject
	scale            float32
	size             gui.Size

	touched map[int]mobile.Touchable
	padded  bool

	onTypedRune func(rune)
	onTypedKey  func(event *gui.KeyEvent)

	inited                bool
	lastTapDown           map[int]time.Time
	lastTapDownPos        map[int]gui.Position
	dragging              gui.Draggable
	dragStart, dragOffset gui.Position

	touchTapCount   int
	touchCancelFunc context.CancelFunc
	touchLastTapped gui.CanvasObject
}

// NewCanvas creates a new gomobile mobileCanvas. This is a mobileCanvas that will render on a mobile device using OpenGL.
func NewCanvas() gui.Canvas {
	ret := &mobileCanvas{padded: true}
	ret.scale = gui.CurrentDevice().SystemScaleForWindow(nil) // we don't need a window parameter on mobile
	ret.touched = make(map[int]mobile.Touchable)
	ret.lastTapDownPos = make(map[int]gui.Position)
	ret.lastTapDown = make(map[int]time.Time)
	ret.Initialize(ret, ret.overlayChanged)
	ret.OnFocus = ret.handleKeyboard
	ret.OnUnfocus = hideVirtualKeyboard

	return ret
}

func (c *mobileCanvas) Capture() image.Image {
	return c.Painter().Capture(c)
}

func (c *mobileCanvas) Content() gui.CanvasObject {
	return c.content
}

func (c *mobileCanvas) InteractiveArea() (gui.Position, gui.Size) {
	scale := gui.CurrentDevice().SystemScaleForWindow(nil) // we don't need a window parameter on mobile

	dev, ok := gui.CurrentDevice().(*device)
	if !ok || dev.safeWidth == 0 || dev.safeHeight == 0 {
		return gui.NewPos(0, 0), c.Size() // running in test mode
	}

	return gui.NewPos(float32(dev.safeLeft)/scale, float32(dev.safeTop)/scale),
		gui.NewSize(float32(dev.safeWidth)/scale, float32(dev.safeHeight)/scale)
}

func (c *mobileCanvas) OnTypedKey() func(*gui.KeyEvent) {
	return c.onTypedKey
}

func (c *mobileCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *mobileCanvas) PixelCoordinateForPosition(pos gui.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *mobileCanvas) Scale() float32 {
	return c.scale
}

func (c *mobileCanvas) SetContent(content gui.CanvasObject) {
	c.setContent(content)
	c.sizeContent(c.Size()) // fixed window size for mobile, cannot stretch to new content
	c.SetDirty()
}

func (c *mobileCanvas) SetOnTypedKey(typed func(*gui.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *mobileCanvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *mobileCanvas) Size() gui.Size {
	return c.size
}

func (c *mobileCanvas) MinSize() gui.Size {
	return c.size // TODO check
}

func (c *mobileCanvas) findObjectAtPositionMatching(pos gui.Position, test func(object gui.CanvasObject) bool) (gui.CanvasObject, gui.Position, int) {
	if c.menu != nil {
		return driver.FindObjectAtPositionMatching(pos, test, c.Overlays().Top(), c.menu)
	}

	return driver.FindObjectAtPositionMatching(pos, test, c.Overlays().Top(), c.windowHead, c.content)
}

func (c *mobileCanvas) handleKeyboard(obj gui.Focusable) {
	isDisabled := false
	if disWid, ok := obj.(gui.Disableable); ok {
		isDisabled = disWid.Disabled()
	}
	if obj != nil && !isDisabled {
		if keyb, ok := obj.(mobile.Keyboardable); ok {
			showVirtualKeyboard(keyb.Keyboard())
		} else {
			showVirtualKeyboard(mobile.DefaultKeyboard)
		}
	} else {
		hideVirtualKeyboard()
	}
}

func (c *mobileCanvas) overlayChanged() {
	c.handleKeyboard(c.Focused())
	c.SetDirty()
}

func (c *mobileCanvas) Resize(size gui.Size) {
	if size == c.size {
		return
	}

	c.sizeContent(size)
}

func (c *mobileCanvas) setContent(content gui.CanvasObject) {
	c.content = content
	c.SetContentTreeAndFocusMgr(content)
}

func (c *mobileCanvas) setMenu(menu gui.CanvasObject) {
	c.menu = menu
	c.SetMenuTreeAndFocusMgr(menu)
}

func (c *mobileCanvas) setWindowHead(head gui.CanvasObject) {
	c.windowHead = head
	c.SetMobileWindowHeadTree(head)
}

func (c *mobileCanvas) applyThemeOutOfTreeObjects() {
	if c.menu != nil {
		app.ApplyThemeTo(c.menu, c) // Ensure our menu gets the theme change message as it's out-of-tree
	}
	if c.windowHead != nil {
		app.ApplyThemeTo(c.windowHead, c) // Ensure our child windows get the theme change message as it's out-of-tree
	}
}

func (c *mobileCanvas) sizeContent(size gui.Size) {
	if c.content == nil { // window may not be configured yet
		return
	}
	c.size = size

	offset := gui.NewPos(0, 0)
	areaPos, areaSize := c.InteractiveArea()

	if c.windowHead != nil {
		topHeight := c.windowHead.MinSize().Height

		if len(c.windowHead.(*gui.Container).Objects) > 1 {
			c.windowHead.Resize(gui.NewSize(areaSize.Width, topHeight))
			offset = gui.NewPos(0, topHeight)
			areaSize = areaSize.Subtract(offset)
		} else {
			c.windowHead.Resize(c.windowHead.MinSize())
		}
		c.windowHead.Move(areaPos)
	}

	topLeft := areaPos.Add(offset)
	for _, overlay := range c.Overlays().List() {
		if p, ok := overlay.(*widget.PopUp); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Refresh()
		} else {
			overlay.Resize(areaSize)
			overlay.Move(topLeft)
		}
	}

	if c.padded {
		c.content.Resize(areaSize.Subtract(gui.NewSize(theme.Padding()*2, theme.Padding()*2)))
		c.content.Move(topLeft.Add(gui.NewPos(theme.Padding(), theme.Padding())))
	} else {
		c.content.Resize(areaSize)
		c.content.Move(topLeft)
	}
}

func (c *mobileCanvas) tapDown(pos gui.Position, tapID int) {
	c.lastTapDown[tapID] = time.Now()
	c.lastTapDownPos[tapID] = pos
	c.dragging = nil

	co, objPos, layer := c.findObjectAtPositionMatching(pos, func(object gui.CanvasObject) bool {
		switch object.(type) {
		case mobile.Touchable, gui.Focusable:
			return true
		}

		return false
	})

	if wid, ok := co.(mobile.Touchable); ok {
		touchEv := &mobile.TouchEvent{}
		touchEv.Position = objPos
		touchEv.AbsolutePosition = pos
		wid.TouchDown(touchEv)
		c.touched[tapID] = wid
	}

	if layer != 1 { // 0 - overlay, 1 - window head / menu, 2 - content
		if wid, ok := co.(gui.Focusable); !ok || wid != c.Focused() {
			c.Unfocus()
		}
	}
}

func (c *mobileCanvas) tapMove(pos gui.Position, tapID int,
	dragCallback func(gui.Draggable, *gui.DragEvent)) {
	previousPos := c.lastTapDownPos[tapID]
	deltaX := pos.X - previousPos.X
	deltaY := pos.Y - previousPos.Y

	if c.dragging == nil && (math.Abs(float64(deltaX)) < tapMoveThreshold && math.Abs(float64(deltaY)) < tapMoveThreshold) {
		return
	}
	c.lastTapDownPos[tapID] = pos

	co, objPos, _ := c.findObjectAtPositionMatching(pos, func(object gui.CanvasObject) bool {
		if _, ok := object.(gui.Draggable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		}

		return false
	})

	if c.touched[tapID] != nil {
		if touch, ok := co.(mobile.Touchable); !ok || c.touched[tapID] != touch {
			touchEv := &mobile.TouchEvent{}
			touchEv.Position = objPos
			touchEv.AbsolutePosition = pos
			c.touched[tapID].TouchCancel(touchEv)
			c.touched[tapID] = nil
		}
	}

	if c.dragging == nil {
		if drag, ok := co.(gui.Draggable); ok {
			c.dragging = drag
			c.dragOffset = previousPos.Subtract(objPos)
			c.dragStart = co.Position()
		} else {
			return
		}
	}

	ev := new(gui.DragEvent)
	draggedObjDelta := c.dragStart.Subtract(c.dragging.(gui.CanvasObject).Position())
	ev.Position = pos.Subtract(c.dragOffset).Add(draggedObjDelta)
	ev.Dragged = gui.Delta{DX: deltaX, DY: deltaY}

	dragCallback(c.dragging, ev)
}

func (c *mobileCanvas) tapUp(pos gui.Position, tapID int,
	tapCallback func(gui.Tappable, *gui.PointEvent),
	tapAltCallback func(gui.SecondaryTappable, *gui.PointEvent),
	doubleTapCallback func(gui.DoubleTappable, *gui.PointEvent),
	dragCallback func(gui.Draggable)) {

	if c.dragging != nil {
		dragCallback(c.dragging)

		c.dragging = nil
		return
	}

	duration := time.Since(c.lastTapDown[tapID])

	if c.menu != nil && c.Overlays().Top() == nil && pos.X > c.menu.Size().Width {
		c.menu.Hide()
		c.menu.Refresh()
		c.setMenu(nil)
		return
	}

	co, objPos, _ := c.findObjectAtPositionMatching(pos, func(object gui.CanvasObject) bool {
		if _, ok := object.(gui.Tappable); ok {
			return true
		} else if _, ok := object.(gui.SecondaryTappable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		} else if _, ok := object.(gui.DoubleTappable); ok {
			return true
		}

		return false
	})

	if wid, ok := co.(mobile.Touchable); ok {
		touchEv := &mobile.TouchEvent{}
		touchEv.Position = objPos
		touchEv.AbsolutePosition = pos
		wid.TouchUp(touchEv)
		c.touched[tapID] = nil
	}

	ev := new(gui.PointEvent)
	ev.Position = objPos
	ev.AbsolutePosition = pos

	if duration < tapSecondaryDelay {
		_, doubleTap := co.(gui.DoubleTappable)
		if doubleTap {
			c.touchTapCount++
			c.touchLastTapped = co
			if c.touchCancelFunc != nil {
				c.touchCancelFunc()
				return
			}
			go c.waitForDoubleTap(co, ev, tapCallback, doubleTapCallback)
		} else {
			if wid, ok := co.(gui.Tappable); ok {
				tapCallback(wid, ev)
			}
		}
	} else {
		if wid, ok := co.(gui.SecondaryTappable); ok {
			tapAltCallback(wid, ev)
		}
	}
}

func (c *mobileCanvas) waitForDoubleTap(co gui.CanvasObject, ev *gui.PointEvent, tapCallback func(gui.Tappable, *gui.PointEvent), doubleTapCallback func(gui.DoubleTappable, *gui.PointEvent)) {
	var ctx context.Context
	ctx, c.touchCancelFunc = context.WithDeadline(context.TODO(), time.Now().Add(time.Millisecond*doubleClickDelay))
	defer c.touchCancelFunc()
	<-ctx.Done()
	if c.touchTapCount == 2 && c.touchLastTapped == co {
		if wid, ok := co.(gui.DoubleTappable); ok {
			doubleTapCallback(wid, ev)
		}
	} else {
		if wid, ok := co.(gui.Tappable); ok {
			tapCallback(wid, ev)
		}
	}
	c.touchTapCount = 0
	c.touchCancelFunc = nil
	c.touchLastTapped = nil
}
