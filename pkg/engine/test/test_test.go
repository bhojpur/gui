package test_test

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
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func TestAssertCanvasTappableAt(t *testing.T) {
	c := test.NewCanvas()
	b := widget.NewButton("foo", nil)
	c.SetContent(b)
	c.Resize(gui.NewSize(300, 300))
	b.Resize(gui.NewSize(100, 100))
	b.Move(gui.NewPos(100, 100))

	tt := &testing.T{}
	assert.True(t, test.AssertCanvasTappableAt(tt, c, gui.NewPos(101, 101)), "tappable found")
	assert.False(t, tt.Failed(), "test did not fail")

	tt = &testing.T{}
	assert.False(t, test.AssertCanvasTappableAt(tt, c, gui.NewPos(99, 99)), "tappable not found")
	assert.True(t, tt.Failed(), "test failed")
}

func TestAssertRendersToMarkup(t *testing.T) {
	c := test.NewCanvas()
	c.SetContent(canvas.NewCircle(color.Black))

	markup := "<canvas padded size=\"9x9\">\n" +
		"\t<content>\n" +
		"\t\t<circle fillColor=\"rgba(0,0,0,255)\" pos=\"4,4\" size=\"1x1\"/>\n" +
		"\t</content>\n" +
		"</canvas>\n"

	t.Run("non-existing master", func(t *testing.T) {
		tt := &testing.T{}
		assert.False(t, test.AssertRendersToMarkup(tt, "non_existing_master.xml", c), "non existing master is not equal to rendered markup")
		assert.True(t, tt.Failed(), "test failed")
		assert.Equal(t, markup, readMarkup(t, "testdata/failed/non_existing_master.xml"), "markup was written to disk")
	})

	t.Run("matching master", func(t *testing.T) {
		tt := &testing.T{}
		assert.True(t, test.AssertRendersToMarkup(tt, "markup_master.xml", c), "existing master is equal to rendered markup")
		assert.False(t, tt.Failed(), "test should not fail")
	})

	t.Run("diffing master", func(t *testing.T) {
		tt := &testing.T{}
		assert.False(t, test.AssertRendersToMarkup(tt, "markup_diffing_master.xml", c), "existing master is not equal to rendered markup")
		assert.True(t, tt.Failed(), "test should fail")
		assert.Equal(t, markup, readMarkup(t, "testdata/failed/markup_diffing_master.xml"), "markup was written to disk")
	})

	if !t.Failed() {
		os.RemoveAll("testdata/failed")
	}
}

func TestDrag(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	d := &draggable{}
	c.SetContent(gui.NewContainerWithoutLayout(d))
	c.Resize(gui.NewSize(30, 30))
	d.Resize(gui.NewSize(20, 20))
	d.Move(gui.NewPos(10, 10))

	test.Drag(c, gui.NewPos(5, 5), 10, 10)
	assert.Nil(t, d.event, "nothing happens if no draggable was found at position")
	assert.False(t, d.wasDragged)

	test.Drag(c, gui.NewPos(15, 15), 17, 42)
	assert.Equal(t, &gui.DragEvent{
		PointEvent: gui.PointEvent{Position: gui.Position{X: 5, Y: 5}},
		Dragged:    gui.NewDelta(17, 42),
	}, d.event)
	assert.True(t, d.wasDragged)
}

func TestFocusNext(t *testing.T) {
	c := test.NewCanvas()
	f1 := &focusable{}
	f2 := &focusable{}
	f3 := &focusable{}
	c.SetContent(gui.NewContainerWithoutLayout(f1, f2, f3))

	assert.Nil(t, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f2, c.Focused())
	assert.False(t, f1.focused)
	assert.True(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.True(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)
}

func TestFocusPrevious(t *testing.T) {
	c := test.NewCanvas()
	f1 := &focusable{}
	f2 := &focusable{}
	f3 := &focusable{}
	c.SetContent(gui.NewContainerWithoutLayout(f1, f2, f3))

	assert.Nil(t, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.True(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f2, c.Focused())
	assert.False(t, f1.focused)
	assert.True(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.True(t, f3.focused)
}

func TestScroll(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	s := &scrollable{}
	c.SetContent(gui.NewContainerWithoutLayout(s))
	c.Resize(gui.NewSize(30, 30))
	s.Resize(gui.NewSize(20, 20))
	s.Move(gui.NewPos(10, 10))

	test.Scroll(c, gui.NewPos(5, 5), 10, 10)
	assert.Nil(t, s.event, "nothing happens if no scrollable was found at position")

	test.Scroll(c, gui.NewPos(15, 15), 17, 42)
	assert.Equal(t, &gui.ScrollEvent{Scrolled: gui.NewDelta(17, 42)}, s.event)
}

func readMarkup(t *testing.T, path string) string {
	raw, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return string(raw)
}

var _ gui.Draggable = (*draggable)(nil)

type draggable struct {
	widget.BaseWidget
	event      *gui.DragEvent
	wasDragged bool
}

func (d *draggable) DragEnd() {
	d.wasDragged = true
}

func (d *draggable) Dragged(event *gui.DragEvent) {
	d.event = event
}

var _ gui.Focusable = (*focusable)(nil)

type focusable struct {
	widget.BaseWidget
	focused bool
}

func (f *focusable) FocusGained() {
	f.focused = true
}

func (f *focusable) FocusLost() {
	f.focused = false
}

func (f *focusable) TypedKey(event *gui.KeyEvent) {
}

func (f *focusable) TypedRune(r rune) {
}

var _ gui.Scrollable = (*scrollable)(nil)

type scrollable struct {
	widget.BaseWidget
	event *gui.ScrollEvent
}

func (s *scrollable) Scrolled(event *gui.ScrollEvent) {
	s.event = event
}
