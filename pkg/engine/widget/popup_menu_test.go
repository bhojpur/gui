package widget_test

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

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func TestPopUpMenu_Move(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.Show()
	test.AssertRendersToMarkup(t, "popup_menu/shown.xml", c)

	m.Move(gui.NewPos(20, 20))
	test.AssertRendersToMarkup(t, "popup_menu/moved.xml", c)

	m.Move(gui.NewPos(190, 10))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_right.xml", c)

	m.Move(gui.NewPos(10, 190))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_down.xml", c)
}

func TestPopUpMenu_Resize(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.ShowAtPosition(gui.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "popup_menu/shown_at_pos.xml", c)

	m.Resize(m.Size().Add(gui.NewSize(10, 10)))
	test.AssertRendersToMarkup(t, "popup_menu/grown.xml", c)

	largeSize := c.Size().Add(gui.NewSize(10, 10))
	m.Resize(largeSize)
	test.AssertRendersToMarkup(t, "popup_menu/canvas_too_small.xml", c)
	assert.Equal(t, gui.NewSize(largeSize.Width, c.Size().Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func TestPopUpMenu_Show(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.Show()
	test.AssertRendersToMarkup(t, "popup_menu/shown.xml", c)
}

func TestPopUpMenu_ShowAtPosition(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.ShowAtPosition(gui.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "popup_menu/shown_at_pos.xml", c)

	m.Hide()
	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.ShowAtPosition(gui.NewPos(190, 10))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_right.xml", c)

	m.Hide()
	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.ShowAtPosition(gui.NewPos(10, 190))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_down.xml", c)

	m.Hide()
	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)
	menuSize := c.Size().Add(gui.NewSize(10, 10))
	m.Resize(menuSize)

	m.ShowAtPosition(gui.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "popup_menu/canvas_too_small.xml", c)
	assert.Equal(t, gui.NewSize(menuSize.Width, c.Size().Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func setupPopUpMenuTest() (*widget.PopUpMenu, gui.Window) {
	test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.NRGBA{G: 150, B: 150, A: 255}))
	w.Resize(gui.NewSize(200, 200))
	m := widget.NewPopUpMenu(gui.NewMenu(
		"",
		gui.NewMenuItem("Option A", nil),
		gui.NewMenuItem("Option B", nil),
	), w.Canvas())
	return m, w
}

func setupPopUpMenuWithSubmenusTest(callback func(string)) (*widget.PopUpMenu, gui.Window) {
	test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.NRGBA{G: 150, B: 150, A: 255}))
	w.Resize(gui.NewSize(800, 600))
	itemA := gui.NewMenuItem("Option A", func() { callback("Option A") })
	itemB := gui.NewMenuItem("Option B", func() { callback("Option B") })
	itemBA := gui.NewMenuItem("Sub Option A", func() { callback("Sub Option A") })
	itemBB := gui.NewMenuItem("Sub Option B", func() { callback("Sub Option B") })
	itemBBA := gui.NewMenuItem("Sub Sub Option A", func() { callback("Sub Sub Option A") })
	itemBBB := gui.NewMenuItem("Sub Sub Option B", func() { callback("Sub Sub Option B") })
	itemB.ChildMenu = gui.NewMenu("", itemBA, itemBB)
	itemBB.ChildMenu = gui.NewMenu("", itemBBA, itemBBB)
	m := widget.NewPopUpMenu(gui.NewMenu("", itemA, itemB), w.Canvas())
	return m, w
}

func tearDownPopUpMenuTest(w gui.Window) {
	w.Close()
	test.NewApp()
}
