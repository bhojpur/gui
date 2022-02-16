//go:build !mobile
// +build !mobile

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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/test"
)

func TestPopUpMenu_KeyboardControl(t *testing.T) {
	var lastTriggered string
	m, w := setupPopUpMenuWithSubmenusTest(func(triggered string) { lastTriggered = triggered })
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()
	m.ShowAtPosition(gui.NewPos(13, 45))

	focused := c.Focused()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	assert.Equal(t, "", lastTriggered)

	// Traverse
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_active.xml", c, "nothing happens when trying open entry without submenu")
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c, "no wrap when reaching end")
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_active.xml", c, "no wrap when reaching end")
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_sub_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_sub_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_sub_active.xml", c, "no wrap when reaching end")
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_sub_active.xml", c, "nothing happens when trying open entry without submenu")
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyLeft})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyUp})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyUp})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_active.xml", c, "no wrap when reaching start")
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyLeft})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyLeft})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c, "nothing happens when no sub-menus are open")

	// trigger actions
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyEnter})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "Option B", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)

	m.Show()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyReturn})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "Sub Option A", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)

	m.Show()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyRight})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeySpace})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "Sub Sub Option A", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)

	// dismiss without triggering action
	lastTriggered = "none"
	m.Show()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	c.Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyEscape})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "none", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)
}
