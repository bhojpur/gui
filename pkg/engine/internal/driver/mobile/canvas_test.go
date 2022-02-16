//go:build !windows || !ci
// +build !windows !ci

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
	"image/color"
	"os"
	"testing"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/driver/mobile"
	"github.com/bhojpur/gui/pkg/engine/layout"
	_ "github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	currentApp := gui.CurrentApp()
	gui.SetCurrentApp(newTestMobileApp())
	ret := m.Run()
	gui.SetCurrentApp(currentApp)
	os.Exit(ret)
}

func TestCanvas_ChildMinSizeChangeAffectsAncestorsUpToRoot(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(gui.NewSize(100, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(gui.NewSize(100, 50))
	leftCol := container.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(gui.NewSize(100, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(gui.NewSize(100, 50))
	rightCol := container.NewVBox(rightObj1, rightObj2)
	content := container.NewHBox(leftCol, rightCol)
	c.SetContent(gui.NewContainerWithLayout(layout.NewCenterLayout(), content))
	c.Resize(gui.NewSize(300, 300))

	oldContentSize := gui.NewSize(200+theme.Padding(), 100+theme.Padding())
	assert.Equal(t, oldContentSize, content.Size())

	leftObj1.SetMinSize(gui.NewSize(110, 60))
	c.EnsureMinSize()

	expectedContentSize := oldContentSize.Add(gui.NewSize(10, 10))
	assert.Equal(t, expectedContentSize, content.Size())
}

func TestCanvas_PixelCoordinateAtPosition(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)

	pos := gui.NewPos(4, 4)
	c.scale = 2.5
	x, y := c.PixelCoordinateForPosition(pos)
	assert.Equal(t, 10, x)
	assert.Equal(t, 10, y)
}

func TestCanvas_Tapped(t *testing.T) {
	tapped := false
	altTapped := false
	buttonTap := false
	var pointEvent *gui.PointEvent
	var tappedObj gui.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(button)
	c.Resize(gui.NewSize(36, 24))
	button.Move(gui.NewPos(3, 3))

	tapPos := gui.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	c.tapUp(tapPos, 0, func(wid gui.Tappable, ev *gui.PointEvent) {
		tapped = true
		tappedObj = wid
		pointEvent = ev
		wid.Tapped(ev)
	}, func(wid gui.SecondaryTappable, ev *gui.PointEvent) {
		altTapped = true
		wid.TappedSecondary(ev)
	}, func(wid gui.DoubleTappable, ev *gui.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid gui.Draggable) {
	})

	assert.True(t, tapped, "tap primary")
	assert.False(t, altTapped, "don't tap secondary")
	assert.True(t, buttonTap, "button should be tapped")
	assert.Equal(t, button, tappedObj)
	if assert.NotNil(t, pointEvent) {
		assert.Equal(t, gui.NewPos(6, 6), pointEvent.AbsolutePosition)
		assert.Equal(t, gui.NewPos(3, 3), pointEvent.Position)
	}
}

func TestCanvas_Tapped_Multi(t *testing.T) {
	buttonTap := false
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(button)
	c.Resize(gui.NewSize(36, 24))
	button.Move(gui.NewPos(3, 3))

	tapPos := gui.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	c.tapUp(tapPos, 1, func(wid gui.Tappable, ev *gui.PointEvent) { // different tapID
		wid.Tapped(ev)
	}, func(wid gui.SecondaryTappable, ev *gui.PointEvent) {
	}, func(wid gui.DoubleTappable, ev *gui.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid gui.Draggable) {
	})

	assert.False(t, buttonTap, "button should not be tapped")
}

func TestCanvas_TappedSecondary(t *testing.T) {
	var pointEvent *gui.PointEvent
	var altTappedObj gui.SecondaryTappable
	obj := &tappableLabel{}
	obj.ExtendBaseWidget(obj)
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(obj)
	c.Resize(gui.NewSize(36, 24))
	obj.Move(gui.NewPos(3, 3))

	tapPos := gui.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	time.Sleep(310 * time.Millisecond)
	c.tapUp(tapPos, 0, func(wid gui.Tappable, ev *gui.PointEvent) {
		obj.tap = true
		wid.Tapped(ev)
	}, func(wid gui.SecondaryTappable, ev *gui.PointEvent) {
		obj.altTap = true
		altTappedObj = wid
		pointEvent = ev
		wid.TappedSecondary(ev)
	}, func(wid gui.DoubleTappable, ev *gui.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid gui.Draggable) {
	})

	assert.False(t, obj.tap, "don't tap primary")
	assert.True(t, obj.altTap, "tap secondary")
	assert.Equal(t, obj, altTappedObj)
	if assert.NotNil(t, pointEvent) {
		assert.Equal(t, gui.NewPos(6, 6), pointEvent.AbsolutePosition)
		assert.Equal(t, gui.NewPos(3, 3), pointEvent.Position)
	}
}

func TestCanvas_Dragged(t *testing.T) {
	dragged := false
	var draggedObj gui.Draggable
	scroll := container.NewScroll(widget.NewLabel("Hi\nHi\nHi"))
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(scroll)
	c.Resize(gui.NewSize(40, 24))
	assert.Equal(t, float32(0), scroll.Offset.Y)

	c.tapDown(gui.NewPos(32, 3), 0)
	c.tapMove(gui.NewPos(32, 10), 0, func(wid gui.Draggable, ev *gui.DragEvent) {
		wid.Dragged(ev)
		dragged = true
		draggedObj = wid
	})

	assert.True(t, dragged)
	assert.Equal(t, scroll, draggedObj)
	dragged = false
	c.tapMove(gui.NewPos(32, 5), 0, func(wid gui.Draggable, ev *gui.DragEvent) {
		wid.Dragged(ev)
		dragged = true
	})
	assert.True(t, dragged)
	assert.Equal(t, gui.NewPos(0, 5), scroll.Offset)
}

func TestCanvas_DraggingOutOfWidget(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	slider := widget.NewSlider(0.0, 100.0)
	c.SetContent(container.NewGridWithRows(2, slider, widget.NewLabel("Outside")))
	c.Resize(gui.NewSize(100, 200))

	assert.Zero(t, slider.Value)
	lastValue := slider.Value

	dragged := false
	c.tapDown(gui.NewPos(23, 13), 0)
	c.tapMove(gui.NewPos(30, 13), 0, func(wid gui.Draggable, ev *gui.DragEvent) {
		assert.Equal(t, slider, wid)
		wid.Dragged(ev)
		dragged = true
	})
	assert.True(t, dragged)
	assert.Greater(t, slider.Value, lastValue)
	lastValue = slider.Value

	dragged = false
	c.tapMove(gui.NewPos(40, 120), 0, func(wid gui.Draggable, ev *gui.DragEvent) {
		assert.Equal(t, slider, wid)
		wid.Dragged(ev)
		dragged = true
	})
	assert.True(t, dragged)
	assert.Greater(t, slider.Value, lastValue)
}

func TestCanvas_Tappable(t *testing.T) {
	content := &touchableLabel{Label: widget.NewLabel("Hi\nHi\nHi")}
	content.ExtendBaseWidget(content)
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(content)
	c.Resize(gui.NewSize(36, 24))
	content.Resize(gui.NewSize(24, 24))

	c.tapDown(gui.NewPos(15, 15), 0)
	assert.True(t, content.down)

	c.tapUp(gui.NewPos(15, 15), 0, func(wid gui.Tappable, ev *gui.PointEvent) {
	}, func(wid gui.SecondaryTappable, ev *gui.PointEvent) {
	}, func(wid gui.DoubleTappable, ev *gui.PointEvent) {
	}, func(wid gui.Draggable) {
	})
	assert.True(t, content.up)

	c.tapDown(gui.NewPos(15, 15), 0)
	c.tapMove(gui.NewPos(35, 15), 0, func(wid gui.Draggable, ev *gui.DragEvent) {
		wid.Dragged(ev)
	})
	assert.True(t, content.cancel)
}

func TestWindow_TappedAndDoubleTapped(t *testing.T) {
	tapped := 0
	but := newDoubleTappableButton()
	but.OnTapped = func() {
		tapped = 1
	}
	but.onDoubleTap = func() {
		tapped = 2
	}

	c := NewCanvas().(*mobileCanvas)
	c.SetContent(gui.NewContainerWithLayout(layout.NewMaxLayout(), but))
	c.Resize(gui.NewSize(36, 24))

	simulateTap(c)
	time.Sleep(700 * time.Millisecond)
	assert.Equal(t, 1, tapped)

	simulateTap(c)
	simulateTap(c)
	time.Sleep(700 * time.Millisecond)
	assert.Equal(t, 2, tapped)
}

func TestGlCanvas_ResizeWithModalPopUpOverlay(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)

	c.SetContent(widget.NewLabel("Content"))

	popup := widget.NewModalPopUp(widget.NewLabel("PopUp"), c)
	popupBgSize := gui.NewSize(200, 200)
	popup.Show()
	popup.Resize(popupBgSize)

	canvasSize := gui.NewSize(600, 700)
	c.Resize(canvasSize)

	// get popup content padding dynamically
	popupContentPadding := popup.MinSize().Subtract(popup.Content.MinSize())

	assert.Equal(t, popupBgSize.Subtract(popupContentPadding), popup.Content.Size())
	assert.Equal(t, canvasSize, popup.Size())
}

func TestCanvas_Focusable(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	content := newFocusableEntry(c)
	c.SetContent(content)
	content.Resize(gui.NewSize(25, 25))

	pos := gui.NewPos(10, 10)
	c.tapDown(pos, 0)
	c.tapUp(pos, 0, func(wid gui.Tappable, ev *gui.PointEvent) {
		wid.Tapped(ev)
	}, nil, nil, nil)
	time.Sleep(time.Millisecond * (doubleClickDelay + 150))
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	c.tapDown(pos, 1)
	c.tapUp(pos, 1, func(wid gui.Tappable, ev *gui.PointEvent) {
		wid.Tapped(ev)
	}, nil, nil, nil)
	time.Sleep(time.Millisecond * (doubleClickDelay + 150))
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	c.Focus(content)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	c.Unfocus()
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)

	content.Disable()
	c.Focus(content)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)

	c.tapDown(gui.NewPos(10, 10), 2)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)
}

type touchableLabel struct {
	*widget.Label
	down, up, cancel bool
}

func (t *touchableLabel) TouchDown(event *mobile.TouchEvent) {
	t.down = true
}

func (t *touchableLabel) TouchUp(event *mobile.TouchEvent) {
	t.up = true
}

func (t *touchableLabel) TouchCancel(event *mobile.TouchEvent) {
	t.cancel = true
}

type tappableLabel struct {
	widget.Label
	tap, altTap bool
}

func (t *tappableLabel) Tapped(_ *gui.PointEvent) {
	t.tap = true
}

func (t *tappableLabel) TappedSecondary(_ *gui.PointEvent) {
	t.altTap = true
}

type focusableEntry struct {
	widget.Entry
	c gui.Canvas

	focusedTimes   int
	unfocusedTimes int
}

func newFocusableEntry(c gui.Canvas) *focusableEntry {
	entry := &focusableEntry{c: c}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (f *focusableEntry) FocusGained() {
	f.focusedTimes++
	f.Entry.FocusGained()
}

func (f *focusableEntry) FocusLost() {
	f.unfocusedTimes++
	f.Entry.FocusLost()
}

func (f *focusableEntry) Tapped(ev *gui.PointEvent) {
	f.c.Focus(f)
}

type doubleTappableButton struct {
	widget.Button

	onDoubleTap func()
}

func (t *doubleTappableButton) DoubleTapped(_ *gui.PointEvent) {
	t.onDoubleTap()
}

func newDoubleTappableButton() *doubleTappableButton {
	but := &doubleTappableButton{}
	but.ExtendBaseWidget(but)

	return but
}

func simulateTap(c *mobileCanvas) {
	c.tapDown(gui.NewPos(15, 15), 0)
	time.Sleep(50 * time.Millisecond)
	c.tapUp(gui.NewPos(15, 15), 0, func(wid gui.Tappable, ev *gui.PointEvent) {
		wid.Tapped(ev)
	}, func(wid gui.SecondaryTappable, ev *gui.PointEvent) {
	}, func(wid gui.DoubleTappable, ev *gui.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid gui.Draggable) {
	})
}

type mobileApp struct {
	gui.App
	driver gui.Driver
}

func (a *mobileApp) Driver() gui.Driver {
	return a.driver
}

func newTestMobileApp() gui.App {
	return &mobileApp{
		App:    gui.CurrentApp(),
		driver: NewGoMobileDriver(),
	}
}
