package dialog

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
	"strconv"
	"sync/atomic"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	internalwidget "github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

var _ gui.Widget = (*colorChannel)(nil)

// colorChannel controls a channel of a color and triggers the callback when changed.
type colorChannel struct {
	widget.BaseWidget
	name      string
	min, max  int
	value     int
	onChanged func(int)
}

// newColorChannel returns a new color channel control for the channel with the given name.
func newColorChannel(name string, min, max, value int, onChanged func(int)) *colorChannel {
	c := &colorChannel{
		name:      name,
		min:       min,
		max:       max,
		value:     clamp(value, min, max),
		onChanged: onChanged,
	}
	c.ExtendBaseWidget(c)
	return c
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (c *colorChannel) CreateRenderer() gui.WidgetRenderer {
	label := widget.NewLabelWithStyle(c.name, gui.TextAlignTrailing, gui.TextStyle{Bold: true})
	entry := newColorChannelEntry(c)
	slider := &widget.Slider{
		Value:       0.0,
		Min:         float64(c.min),
		Max:         float64(c.max),
		Step:        1.0,
		Orientation: widget.Horizontal,
		OnChanged: func(value float64) {
			c.SetValue(int(value))
		},
	}
	r := &colorChannelRenderer{
		BaseRenderer: internalwidget.NewBaseRenderer([]gui.CanvasObject{
			label,
			slider,
			entry,
		}),
		control: c,
		label:   label,
		entry:   entry,
		slider:  slider,
	}
	r.updateObjects()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (c *colorChannel) MinSize() gui.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetValue updates the value in this color widget
func (c *colorChannel) SetValue(value int) {
	value = clamp(value, c.min, c.max)
	if c.value == value {
		return
	}
	c.value = value
	c.Refresh()
	if f := c.onChanged; f != nil {
		f(value)
	}
}

type colorChannelRenderer struct {
	internalwidget.BaseRenderer
	control *colorChannel
	label   *widget.Label
	entry   *colorChannelEntry
	slider  *widget.Slider
}

func (r *colorChannelRenderer) Layout(size gui.Size) {
	lMin := r.label.MinSize()
	eMin := r.entry.MinSize()
	r.label.Move(gui.NewPos(0, (size.Height-lMin.Height)/2))
	r.label.Resize(gui.NewSize(lMin.Width, lMin.Height))
	r.slider.Move(gui.NewPos(lMin.Width, 0))
	r.slider.Resize(gui.NewSize(size.Width-lMin.Width-eMin.Width, size.Height))
	r.entry.Move(gui.NewPos(size.Width-eMin.Width, 0))
	r.entry.Resize(gui.NewSize(eMin.Width, size.Height))
}

func (r *colorChannelRenderer) MinSize() gui.Size {
	lMin := r.label.MinSize()
	sMin := r.slider.MinSize()
	eMin := r.entry.MinSize()
	return gui.NewSize(
		lMin.Width+sMin.Width+eMin.Width,
		gui.Max(lMin.Height, gui.Max(sMin.Height, eMin.Height)),
	)
}

func (r *colorChannelRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.control.Size())
	canvas.Refresh(r.control)
}

func (r *colorChannelRenderer) updateObjects() {
	r.entry.SetText(strconv.Itoa(r.control.value))
	r.slider.Value = float64(r.control.value)
	r.slider.Refresh()
}

type colorChannelEntry struct {
	userChangeEntry
}

func newColorChannelEntry(c *colorChannel) *colorChannelEntry {
	e := &colorChannelEntry{}
	e.Text = "0"
	e.ExtendBaseWidget(e)
	e.setOnChanged(func(text string) {
		value, err := strconv.Atoi(text)
		if err != nil {
			gui.LogError("Couldn't parse: "+text, err)
			return
		}
		c.SetValue(value)
	})
	return e
}

func (e *colorChannelEntry) MinSize() gui.Size {
	// Ensure space for 3 digits
	min := gui.MeasureText("000", theme.TextSize(), gui.TextStyle{})
	min = min.Add(gui.NewSize(theme.Padding()*6, theme.Padding()*4))
	return min.Max(e.Entry.MinSize())
}

type userChangeEntry struct {
	widget.Entry
	userTyped uint32 // atomic, 0 == false, 1 == true
}

func newUserChangeEntry(text string) *userChangeEntry {
	e := &userChangeEntry{}
	e.Entry.Text = text
	e.ExtendBaseWidget(e)
	return e
}

func (e *userChangeEntry) setOnChanged(onChanged func(s string)) {
	e.Entry.OnChanged = func(text string) {
		if !atomic.CompareAndSwapUint32(&e.userTyped, 1, 0) {
			return
		}
		if onChanged != nil {
			onChanged(text)
		}
	}
	e.ExtendBaseWidget(e)
}

func (e *userChangeEntry) TypedRune(r rune) {
	atomic.StoreUint32(&e.userTyped, 1)
	e.Entry.TypedRune(r)
}

func (e *userChangeEntry) TypedKey(ev *gui.KeyEvent) {
	atomic.StoreUint32(&e.userTyped, 1)
	e.Entry.TypedKey(ev)
}
