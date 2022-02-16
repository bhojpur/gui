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
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/data/binding"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Binding(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Init")
	assert.Equal(t, "Init", entry.Text)

	str := binding.NewString()
	entry.Bind(str)
	waitForBinding()
	assert.Equal(t, "", entry.Text)

	err := str.Set("Updated")
	assert.Nil(t, err)
	waitForBinding()
	assert.Equal(t, "Updated", entry.Text)

	entry.SetText("Typed")
	v, err := str.Get()
	assert.Nil(t, err)
	assert.Equal(t, "Typed", v)

	entry.Unbind()
	waitForBinding()
	assert.Equal(t, "Typed", entry.Text)
}

func TestEntry_Clicked(t *testing.T) {
	entry, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("MMM\nWWW\n")
	test.AssertRendersToMarkup(t, "entry/tapped_initial.xml", c)

	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/tapped_focused.xml", c)

	testCharSize := theme.TextSize()
	pos := gui.NewPos(entryOffset+theme.Padding()+testCharSize*1.5, entryOffset+theme.Padding()+testCharSize/2) // tap in the middle of the 2nd "M"
	clickCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_2nd_m.xml", c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	pos = gui.NewPos(entryOffset+theme.Padding()+testCharSize*2.5, entryOffset+theme.Padding()+testCharSize/2) // tap in the middle of the 3rd "M"
	clickCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_3rd_m.xml", c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	pos = gui.NewPos(entryOffset+theme.Padding()+testCharSize*4, entryOffset+theme.Padding()+testCharSize/2) // tap after text
	clickCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_after_last_col.xml", c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 3, entry.CursorColumn)

	pos = gui.NewPos(entryOffset+testCharSize, entryOffset+testCharSize*4) // tap below rows
	clickCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_after_last_row.xml", c)
	assert.Equal(t, 2, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_CursorColumn(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	// only 0 columns, do nothing
	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorColumn)

	// 1, this should increment
	entry.SetText("a")
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorColumn)

	left := &gui.KeyEvent{Name: gui.KeyLeft}
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorColumn)

	// don't go beyond left
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_CursorColumn_Ends(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hello")
	assert.Equal(t, 0, entry.CursorColumn)

	// down should go to end for last line
	down := &gui.KeyEvent{Name: gui.KeyDown}
	entry.TypedKey(down)
	assert.Equal(t, 5, entry.CursorColumn)
	assert.Equal(t, 0, entry.CursorRow)

	// up should go to start for first line
	up := &gui.KeyEvent{Name: gui.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_CursorColumn_Jump(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("a\nbc")

	// go to end of text
	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)
	entry.TypedKey(right)
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	// go up, to a shorter line
	up := &gui.KeyEvent{Name: gui.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_CursorColumn_Wrap(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("a\nb")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// go to end of line
	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	// wrap to new line
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// and back
	left := &gui.KeyEvent{Name: gui.KeyLeft}
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_CursorColumn_Wrap2(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = gui.TextWrapWord
	entry.Resize(gui.NewSize(72, 64))
	entry.SetText("1234")
	entry.CursorColumn = 3
	test.Type(entry, "a")
	test.Type(entry, "b")
	test.Type(entry, "c")
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
	w := test.NewWindow(entry)
	w.Resize(gui.NewSize(70, 70))
	test.AssertImageMatches(t, "entry/wrap_multi_line_cursor.png", w.Canvas().Capture())
}

func TestEntry_CursorPasswordRevealer(t *testing.T) {
	pr := widget.NewPasswordEntry().ActionItem.(desktop.Cursorable)
	assert.Equal(t, desktop.DefaultCursor, pr.Cursor())
}

func TestEntry_CursorRow(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("test")
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &gui.KeyEvent{Name: gui.KeyDown}
	entry.TypedKey(down)
	assert.Equal(t, 0, entry.CursorRow)

	// 2 lines, this should increment
	entry.SetText("test\nrows")
	entry.TypedKey(down)
	assert.Equal(t, 1, entry.CursorRow)

	up := &gui.KeyEvent{Name: gui.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_Disableable(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_empty.xml", c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_disabled_empty.xml", c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_empty.xml", c)

	entry.SetPlaceHolder("Type!")
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_placeholder.xml", c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_disabled_placeholder.xml", c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_placeholder.xml", c)

	entry.SetText("Hello")
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_custom_value.xml", c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_disabled_custom_value.xml", c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_custom_value.xml", c)
}

func TestEntry_Disabled_TextSelection(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	entry.SetText("Testing")
	entry.Disable()
	c := window.Canvas()

	assert.True(t, entry.Disabled())
	test.DoubleTap(entry)

	test.AssertImageMatches(t, "entry/disabled_text_selected.png", c.Capture())

	entry.FocusLost()
	test.AssertImageMatches(t, "entry/disabled_text_unselected.png", c.Capture())
}

func TestEntry_EmptySelection(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("text")

	// trying to select at the edge
	typeKeys(entry, keyShiftLeftDown, gui.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())

	typeKeys(entry, gui.KeyRight)
	assert.Equal(t, 1, entry.CursorColumn)

	// stop selecting at the edge when nothing is selected
	typeKeys(entry, gui.KeyLeft, keyShiftLeftDown, gui.KeyRight, gui.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())
	assert.Equal(t, 0, entry.CursorColumn)

	// check that the selection has been removed
	typeKeys(entry, gui.KeyRight, keyShiftLeftDown, gui.KeyRight, gui.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_Focus(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/focus_gained.xml", c)

	entry.FocusLost()
	test.AssertRendersToMarkup(t, "entry/focus_lost.xml", c)

	window.Canvas().Focus(entry)
	test.AssertRendersToMarkup(t, "entry/focus_gained.xml", c)
}

func TestEntry_FocusWithPopUp(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.TapSecondaryAt(entry, gui.NewPos(1, 1))

	test.AssertRendersToMarkup(t, "entry/focus_with_popup_initial.xml", c)

	test.TapCanvas(c, gui.NewPos(20, 20))
	test.AssertRendersToMarkup(t, "entry/focus_with_popup_entry_selected.xml", c)

	test.TapSecondaryAt(entry, gui.NewPos(1, 1))
	test.AssertRendersToMarkup(t, "entry/focus_with_popup_initial.xml", c)

	test.TapCanvas(c, gui.NewPos(5, 5))
	test.AssertRendersToMarkup(t, "entry/focus_with_popup_dismissed.xml", c)
}

func TestEntry_HidePopUpOnEntry(t *testing.T) {
	entry := widget.NewEntry()
	tapPos := gui.NewPos(1, 1)
	c := gui.CurrentApp().Driver().CanvasForObject(entry)

	assert.Nil(t, c.Overlays().Top())

	test.TapSecondaryAt(entry, tapPos)
	assert.NotNil(t, c.Overlays().Top())

	test.Type(entry, "KJGFD")
	assert.Nil(t, c.Overlays().Top())
	assert.Equal(t, "KJGFD", entry.Text)
}

func TestEntry_MinSize(t *testing.T) {
	entry := widget.NewEntry()
	min := entry.MinSize()
	entry.SetPlaceHolder("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetText("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetPlaceHolder("Hello")
	assert.Equal(t, entry.MinSize().Width, min.Width)
	assert.Equal(t, entry.MinSize().Height, min.Height)

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	entry.Wrapping = gui.TextWrapOff
	entry.Refresh()
	assert.Greater(t, entry.MinSize().Width, min.Width)

	min = entry.MinSize()
	entry.ActionItem = canvas.NewCircle(color.Black)
	assert.Equal(t, min.Add(gui.NewSize(theme.IconInlineSize()+theme.Padding(), 0)), entry.MinSize())
}

func TestEntryMultiline_MinSize(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	min := entry.MinSize()
	entry.SetText("Hello")
	assert.Equal(t, entry.MinSize().Width, min.Width)
	assert.Equal(t, entry.MinSize().Height, min.Height)

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	entry.Wrapping = gui.TextWrapOff
	entry.Refresh()
	assert.Greater(t, entry.MinSize().Width, min.Width)

	entry.Wrapping = gui.TextWrapBreak
	entry.Refresh()
	assert.Equal(t, entry.MinSize().Width, min.Width)

	min = entry.MinSize()
	entry.ActionItem = canvas.NewCircle(color.Black)
	assert.Equal(t, min.Add(gui.NewSize(theme.IconInlineSize()+theme.Padding(), 0)), entry.MinSize())
}

func TestEntry_MultilineSelect(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// Extend the selection down one row
	typeKeys(e, gui.KeyDown)
	test.AssertRendersToMarkup(t, "entry/selection_add_one_row_down.xml", c)
	assert.Equal(t, "sting\nTesti", e.SelectedText())

	typeKeys(e, gui.KeyUp)
	test.AssertRendersToMarkup(t, "entry/selection_remove_one_row_up.xml", c)
	assert.Equal(t, "sti", e.SelectedText())

	typeKeys(e, gui.KeyUp)
	test.AssertRendersToMarkup(t, "entry/selection_remove_add_one_row_up.xml", c)
	assert.Equal(t, "ng\nTe", e.SelectedText())
}

func TestEntry_MultilineWrapping_DeleteWithBackspace(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = gui.TextWrapWord
	entry.Resize(gui.NewSize(64, 64))
	test.Type(entry, "line1")
	test.Type(entry, "\nline2")
	test.Type(entry, "\nline3")

	assert.Equal(t, 5, entry.CursorColumn)
	assert.Equal(t, 2, entry.CursorRow)

	for i := 0; i < 4; i++ {
		entry.TypedKey(&gui.KeyEvent{Name: gui.KeyBackspace})
		assert.Equal(t, 4-i, entry.CursorColumn)
		assert.Equal(t, 2, entry.CursorRow)
	}

	entry.TypedKey(&gui.KeyEvent{Name: gui.KeyBackspace})
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 2, entry.CursorRow)

	assert.NotPanics(t, func() {
		entry.TypedKey(&gui.KeyEvent{Name: gui.KeyBackspace})
	})
	assert.Equal(t, 5, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
}

func TestEntry_Notify(t *testing.T) {
	entry := widget.NewEntry()
	changed := false

	entry.OnChanged = func(string) {
		changed = true
	}
	entry.SetText("Test")

	assert.True(t, changed)
}

func TestEntry_OnCopy(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, gui.KeyRight, gui.KeyRight, keyShiftLeftDown, gui.KeyRight, gui.KeyRight, gui.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &gui.ShortcutCopy{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnCopy_Password(t *testing.T) {
	e := widget.NewPasswordEntry()
	e.SetText("Testing")
	typeKeys(e, keyShiftLeftDown, gui.KeyRight, gui.KeyRight, gui.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &gui.ShortcutCopy{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnCut(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, gui.KeyRight, gui.KeyRight, keyShiftLeftDown, gui.KeyRight, gui.KeyRight, gui.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &gui.ShortcutCut{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Teng", e.Text)
}

func TestEntry_OnCut_Password(t *testing.T) {
	e := widget.NewPasswordEntry()
	e.SetText("Testing")
	typeKeys(e, keyShiftLeftDown, gui.KeyRight, gui.KeyRight, gui.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &gui.ShortcutCut{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnKeyDown(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "Hi")

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Backspace(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")
	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	key := &gui.KeyEvent{Name: gui.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_BackspaceBeyondText(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")
	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)

	key := &gui.KeyEvent{Name: gui.KeyBackspace}
	entry.TypedKey(key)
	entry.TypedKey(key)
	entry.TypedKey(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceBeyondTextAndNewLine(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &gui.KeyEvent{Name: gui.KeyDown}
	entry.TypedKey(down)
	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)

	key := &gui.KeyEvent{Name: gui.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
	assert.Equal(t, "H", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceNewline(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &gui.KeyEvent{Name: gui.KeyDown}
	entry.TypedKey(down)

	key := &gui.KeyEvent{Name: gui.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceUnicode(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "è")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	bs := &gui.KeyEvent{Name: gui.KeyBackspace}
	entry.TypedKey(bs)
	assert.Equal(t, "", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_OnKeyDown_Delete(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")
	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &gui.KeyEvent{Name: gui.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_DeleteBeyondText(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")

	key := &gui.KeyEvent{Name: gui.KeyDelete}
	entry.TypedKey(key)
	entry.TypedKey(key)
	entry.TypedKey(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_DeleteNewline(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("H\ni")

	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)

	key := &gui.KeyEvent{Name: gui.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_HomeEnd(t *testing.T) {
	entry := &widget.Entry{}
	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	end := &gui.KeyEvent{Name: gui.KeyEnd}
	entry.TypedKey(end)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	home := &gui.KeyEvent{Name: gui.KeyHome}
	entry.TypedKey(home)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_OnKeyDown_Insert(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "Hi")
	assert.Equal(t, "Hi", entry.Text)

	left := &gui.KeyEvent{Name: gui.KeyLeft}
	entry.TypedKey(left)

	test.Type(entry, "o")
	assert.Equal(t, "Hoi", entry.Text)
}

func TestEntry_OnKeyDown_Newline(t *testing.T) {
	entry, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/on_key_down_newline_initial.xml", c)

	right := &gui.KeyEvent{Name: gui.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &gui.KeyEvent{Name: gui.KeyReturn}
	entry.TypedKey(key)

	assert.Equal(t, "H\ni", entry.Text)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "o")
	assert.Equal(t, "H\noi", entry.Text)
	test.AssertRendersToMarkup(t, "entry/on_key_down_newline_typed.xml", c)
}

func TestEntry_OnPaste(t *testing.T) {
	clipboard := test.NewClipboard()
	shortcut := &gui.ShortcutPaste{Clipboard: clipboard}
	tests := []struct {
		name             string
		entry            *widget.Entry
		clipboardContent string
		wantText         string
		wantRow, wantCol int
	}{
		{
			name:             "singleline: empty content",
			entry:            widget.NewEntry(),
			clipboardContent: "",
			wantText:         "",
			wantRow:          0,
			wantCol:          0,
		},
		{
			name:             "singleline: simple text",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard content",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: UTF8 text",
			entry:            widget.NewEntry(),
			clipboardContent: "Hié™שרה",
			wantText:         "Hié™שרה",
			wantRow:          0,
			wantCol:          7,
		},
		{
			name:             "singleline: with new line",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: with tab",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard\tcontent",
			wantText:         "clipboard\tcontent",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "password: with new line",
			entry:            widget.NewPasswordEntry(),
			clipboardContent: "3SB=y+)z\nkHGK(hx6 -e_\"1TZu q^bF3^$u H[:e\"1O.",
			wantText:         `3SB=y+)z kHGK(hx6 -e_"1TZu q^bF3^$u H[:e"1O.`,
			wantRow:          0,
			wantCol:          44,
		},
		{
			name:             "multiline: with new line",
			entry:            widget.NewMultiLineEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard\ncontent",
			wantRow:          1,
			wantCol:          7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clipboard.SetContent(tt.clipboardContent)
			tt.entry.TypedShortcut(shortcut)
			assert.Equal(t, tt.wantText, tt.entry.Text)
			assert.Equal(t, tt.wantRow, tt.entry.CursorRow)
			assert.Equal(t, tt.wantCol, tt.entry.CursorColumn)
		})
	}
}

func TestEntry_PageUpDown(t *testing.T) {
	t.Run("single line", func(*testing.T) {
		e, window := setupImageTest(t, false)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing")
		test.AssertRendersToMarkup(t, "entry/select_initial.xml", c)

		// move right, press & hold shift and pagedown
		typeKeys(e, gui.KeyRight, keyShiftLeftDown, gui.KeyPageDown)
		assert.Equal(t, "esting", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_single_line_shift_pagedown.xml", c)

		// while shift is held press pageup
		typeKeys(e, gui.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_single_line_shift_pageup.xml", c)

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, gui.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_single_line_pagedown.xml", c)
	})

	t.Run("page down single line", func(*testing.T) {
		e, window := setupImageTest(t, true)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing\nTesting\nTesting")
		test.AssertRendersToMarkup(t, "entry/select_multi_line_initial.xml", c)

		// move right, press & hold shift and pagedown
		typeKeys(e, gui.KeyRight, keyShiftLeftDown, gui.KeyPageDown)
		assert.Equal(t, "esting\nTesting\nTesting", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_multi_line_shift_pagedown.xml", c)

		// while shift is held press pageup
		typeKeys(e, gui.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_multi_line_shift_pageup.xml", c)

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, gui.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_multi_line_pagedown.xml", c)
	})
}

func TestEntry_PasteOverSelection(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, gui.KeyRight, gui.KeyRight, keyShiftLeftDown, gui.KeyRight, gui.KeyRight, gui.KeyRight)

	clipboard := test.NewClipboard()
	clipboard.SetContent("Insert")
	shortcut := &gui.ShortcutPaste{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "Insert", clipboard.Content())
	assert.Equal(t, "TeInsertng", e.Text)
}

func TestEntry_PasteUnicode(t *testing.T) {
	e := widget.NewMultiLineEntry()
	e.SetText("line")
	e.CursorColumn = 4

	clipboard := test.NewClipboard()
	clipboard.SetContent("thing {\n\titem: 'val测试'\n}")
	shortcut := &gui.ShortcutPaste{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "thing {\n\titem: 'val测试'\n}", clipboard.Content())
	assert.Equal(t, "linething {\n\titem: 'val测试'\n}", e.Text)

	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 1, e.CursorColumn)
}

func TestEntry_Placeholder(t *testing.T) {
	entry := &widget.Entry{}
	entry.Text = "Text"
	entry.PlaceHolder = "Placehold"

	window := test.NewWindow(entry)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, "Text", entry.Text)
	test.AssertRendersToMarkup(t, "entry/placeholder_with_text.xml", c)

	entry.SetText("")
	assert.Equal(t, "", entry.Text)
	test.AssertRendersToMarkup(t, "entry/placeholder_without_text.xml", c)
}

func TestEntry_Select(t *testing.T) {
	for name, tt := range map[string]struct {
		keys          []gui.KeyName
		text          string
		setupReverse  bool
		wantMarkup    string
		wantSelection string
		wantText      string
	}{
		"delete single-line": {
			keys:       []gui.KeyName{gui.KeyDelete},
			wantText:   "Testing\nTeng\nTesting",
			wantMarkup: "entry/selection_delete_single_line.xml",
		},
		"delete multi-line": {
			keys:       []gui.KeyName{gui.KeyDown, gui.KeyDelete},
			wantText:   "Testing\nTeng",
			wantMarkup: "entry/selection_delete_multi_line.xml",
		},
		"delete reverse multi-line": {
			keys:         []gui.KeyName{keyShiftLeftDown, gui.KeyDown, gui.KeyDelete},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantMarkup:   "entry/selection_delete_reverse_multi_line.xml",
		},
		"delete select down with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyDelete, gui.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"delete reverse select down with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyDelete, gui.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"delete select up with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyDelete, gui.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		"delete reverse select up with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyDelete, gui.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		// The backspace delete behaviour is the same as via delete.
		"backspace single-line": {
			keys:       []gui.KeyName{gui.KeyBackspace},
			wantText:   "Testing\nTeng\nTesting",
			wantMarkup: "entry/selection_delete_single_line.xml",
		},
		"backspace multi-line": {
			keys:       []gui.KeyName{gui.KeyDown, gui.KeyBackspace},
			wantText:   "Testing\nTeng",
			wantMarkup: "entry/selection_delete_multi_line.xml",
		},
		"backspace reverse multi-line": {
			keys:         []gui.KeyName{keyShiftLeftDown, gui.KeyDown, gui.KeyBackspace},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantMarkup:   "entry/selection_delete_reverse_multi_line.xml",
		},
		"backspace select down with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyBackspace, gui.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"backspace reverse select down with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyBackspace, gui.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"backspace select up with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyBackspace, gui.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		"backspace reverse select up with Shift held": {
			keys:          []gui.KeyName{keyShiftLeftDown, gui.KeyBackspace, gui.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		// Erase the selection and add a newline at selection start
		"enter": {
			keys:       []gui.KeyName{gui.KeyEnter},
			wantText:   "Testing\nTe\nng\nTesting",
			wantMarkup: "entry/selection_enter.xml",
		},
		"enter reverse": {
			keys:         []gui.KeyName{gui.KeyEnter},
			setupReverse: true,
			wantText:     "Testing\nTe\nng\nTesting",
			wantMarkup:   "entry/selection_enter.xml",
		},
		"replace": {
			text:       "hello",
			wantText:   "Testing\nTehellong\nTesting",
			wantMarkup: "entry/selection_replace.xml",
		},
		"replace reverse": {
			text:         "hello",
			setupReverse: true,
			wantText:     "Testing\nTehellong\nTesting",
			wantMarkup:   "entry/selection_replace.xml",
		},
		"deselect and delete": {
			keys:       []gui.KeyName{keyShiftLeftUp, gui.KeyLeft, gui.KeyDelete},
			wantText:   "Testing\nTeting\nTesting",
			wantMarkup: "entry/selection_deselect_delete.xml",
		},
		"deselect and delete holding shift": {
			keys:       []gui.KeyName{keyShiftLeftUp, gui.KeyLeft, keyShiftLeftDown, gui.KeyDelete},
			wantText:   "Testing\nTeting\nTesting",
			wantMarkup: "entry/selection_deselect_delete.xml",
		},
		// ensure that backspace doesn't leave a selection start at the old cursor position
		"deselect and backspace holding shift": {
			keys:       []gui.KeyName{keyShiftLeftUp, gui.KeyLeft, keyShiftLeftDown, gui.KeyBackspace},
			wantText:   "Testing\nTsting\nTesting",
			wantMarkup: "entry/selection_deselect_backspace.xml",
		},
		// clear selection, select a character and while holding shift issue two backspaces
		"deselect, select and double backspace": {
			keys:       []gui.KeyName{keyShiftLeftUp, gui.KeyRight, gui.KeyLeft, keyShiftLeftDown, gui.KeyLeft, gui.KeyBackspace, gui.KeyBackspace},
			wantText:   "Testing\nTeing\nTesting",
			wantMarkup: "entry/selection_deselect_select_backspace.xml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			entry, window := setupSelection(t, tt.setupReverse)
			defer teardownImageTest(window)
			c := window.Canvas()

			if tt.text != "" {
				test.Type(entry, tt.text)
			} else {
				typeKeys(entry, tt.keys...)
			}
			assert.Equal(t, tt.wantText, entry.Text)
			assert.Equal(t, tt.wantSelection, entry.SelectedText())
			test.AssertRendersToMarkup(t, tt.wantMarkup, c)
		})
	}
}

func TestEntry_SelectAll(t *testing.T) {
	e, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("First Row\nSecond Row\nThird Row")
	test.AssertRendersToMarkup(t, "entry/select_all_initial.xml", c)

	shortcut := &gui.ShortcutSelectAll{}
	e.TypedShortcut(shortcut)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 9, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/select_all_selected.xml", c)
}

func TestEntry_SelectAll_EmptyEntry(t *testing.T) {
	entry := widget.NewEntry()
	entry.TypedShortcut(&gui.ShortcutSelectAll{})

	assert.Equal(t, "", entry.SelectedText())
}

func TestEntry_SelectEndWithoutShift(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// end after releasing shift
	typeKeys(e, keyShiftLeftUp, gui.KeyEnd)
	test.AssertRendersToMarkup(t, "entry/selection_end.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectHomeEnd(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// Hold shift to continue selection
	typeKeys(e, keyShiftLeftDown)

	// T e[s t i]n g -> end -> // T e[s t i n g]
	typeKeys(e, gui.KeyEnd)
	test.AssertRendersToMarkup(t, "entry/selection_add_to_end.xml", c)
	assert.Equal(t, "sting", e.SelectedText())

	// T e[s t i n g] -> home -> [T e]s t i n g
	typeKeys(e, gui.KeyHome)
	test.AssertRendersToMarkup(t, "entry/selection_add_to_home.xml", c)
	assert.Equal(t, "Te", e.SelectedText())
}

func TestEntry_SelectHomeWithoutShift(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// home after releasing shift
	typeKeys(e, keyShiftLeftUp, gui.KeyHome)
	test.AssertRendersToMarkup(t, "entry/selection_home.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapDown(t *testing.T) {
	// down snaps to end, but it also moves
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, gui.KeyDown)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_down.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapLeft(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, gui.KeyLeft)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 2, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_left.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapRight(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, gui.KeyRight)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_right.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapUp(t *testing.T) {
	// up snaps to start, but it also moves
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, gui.KeyUp)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_up.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectedText(t *testing.T) {
	e, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("Testing")
	test.AssertRendersToMarkup(t, "entry/select_initial.xml", c)
	assert.Equal(t, "", e.SelectedText())

	// move right, press & hold shift and move right
	typeKeys(e, gui.KeyRight, keyShiftLeftDown, gui.KeyRight, gui.KeyRight)
	assert.Equal(t, "es", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_selected.xml", c)

	// release shift
	typeKeys(e, keyShiftLeftUp)
	// press shift and move
	typeKeys(e, keyShiftLeftDown, gui.KeyRight)
	assert.Equal(t, "est", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_add_selection.xml", c)

	// release shift and move right
	typeKeys(e, keyShiftLeftUp, gui.KeyRight)
	assert.Equal(t, "", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_move_wo_shift.xml", c)

	// press shift and move left
	typeKeys(e, keyShiftLeftDown, gui.KeyLeft, gui.KeyLeft)
	assert.Equal(t, "st", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_select_left.xml", c)
}

func TestEntry_SelectionHides(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Unfocus()
	test.AssertRendersToMarkup(t, "entry/selection_focus_lost.xml", c)
	assert.Equal(t, "sti", e.SelectedText())

	c.Focus(e)
	test.AssertRendersToMarkup(t, "entry/selection_focus_gained.xml", c)
	assert.Equal(t, "sti", e.SelectedText())
}

func TestEntry_SetPlaceHolder(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 0, len(entry.Text))

	entry.SetPlaceHolder("Test")
	assert.Equal(t, 0, len(entry.Text))
	test.AssertRendersToMarkup(t, "entry/set_placeholder_set.xml", c)

	entry.SetText("Hi")
	assert.Equal(t, 2, len(entry.Text))
	test.AssertRendersToMarkup(t, "entry/set_placeholder_replaced.xml", c)
}

func TestEntry_SetPlaceHolder_ByField(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 0, len(entry.Text))

	entry.PlaceHolder = "Test"
	entry.Refresh()
	assert.Equal(t, 0, len(entry.Text))
	test.AssertRendersToMarkup(t, "entry/set_placeholder_set.xml", c)

	entry.SetText("Hi")
	assert.Equal(t, 2, len(entry.Text))
	test.AssertRendersToMarkup(t, "entry/set_placeholder_replaced.xml", c)
}

func TestEntry_Disable_KeyDown(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "H")
	entry.Disable()
	test.Type(entry, "i")
	assert.Equal(t, "H", entry.Text)

	entry.Enable()
	test.Type(entry, "i")
	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_Disable_OnFocus(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Disable()
	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/focused_disabled.xml", c)

	entry.Enable()
	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/focused_enabled.xml", c)
}

func TestEntry_SetText_EmptyString(t *testing.T) {
	entry := widget.NewEntry()

	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	entry = widget.NewMultiLineEntry()
	test.Type(entry, "test\ntest")

	down := &gui.KeyEvent{Name: gui.KeyDown}
	entry.TypedKey(down)

	assert.Equal(t, 4, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_SetText_Manual(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Text = "Test"
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_changed.xml", c)
}

func TestEntry_SetText_Overflow(t *testing.T) {
	entry := widget.NewEntry()

	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)

	entry.SetText("x")
	assert.Equal(t, 1, entry.CursorColumn)

	key := &gui.KeyEvent{Name: gui.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, 1, entry.CursorColumn)
	assert.Equal(t, "x", entry.Text)

	key = &gui.KeyEvent{Name: gui.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, "", entry.Text)
}

func TestEntry_SetText_Underflow(t *testing.T) {
	entry := widget.NewEntry()
	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)

	entry.Text = ""
	entry.Refresh()
	assert.Equal(t, 0, entry.CursorColumn)

	key := &gui.KeyEvent{Name: gui.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, "", entry.Text)
}

func TestEntry_SetText_Overflow_Multiline(t *testing.T) {
	entry := widget.NewEntry()
	entry.MultiLine = true

	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 0, entry.CursorRow)

	entry.SetText("ab\ncd\nef")
	typeKeys(entry, gui.KeyDown, gui.KeyDown, gui.KeyRight)
	assert.Equal(t, 1, entry.CursorColumn)
	assert.Equal(t, 2, entry.CursorRow)
	entry.SetText("AB\nAAAA")
	assert.Equal(t, 4, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
}

func TestEntry_SetTextStyle(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Text = "Styled Text"
	entry.TextStyle = gui.TextStyle{Bold: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_style_bold.xml", c)

	entry.TextStyle = gui.TextStyle{Monospace: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_style_monospace.xml", c)

	entry.TextStyle = gui.TextStyle{Italic: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_style_italic.xml", c)
}

func TestEntry_Submit(t *testing.T) {
	t.Run("Callback", func(t *testing.T) {
		var submission string
		entry := &widget.Entry{
			OnSubmitted: func(s string) {
				submission = s
			},
		}
		t.Run("SingleLine_Enter", func(t *testing.T) {
			entry.MultiLine = false
			entry.SetText("a")
			entry.TypedKey(&gui.KeyEvent{Name: gui.KeyEnter})
			assert.Equal(t, "a", entry.Text)
			assert.Equal(t, "a", submission)
		})
		t.Run("SingleLine_Return", func(t *testing.T) {
			entry.MultiLine = false
			entry.SetText("b")
			entry.TypedKey(&gui.KeyEvent{Name: gui.KeyReturn})
			assert.Equal(t, "b", entry.Text)
			assert.Equal(t, "b", submission)
		})
		t.Run("MultiLine_ShiftEnter", func(t *testing.T) {
			entry.MultiLine = true
			entry.SetText("c")
			typeKeys(entry, keyShiftLeftDown, gui.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "c", entry.Text)
			assert.Equal(t, "c", submission)
			entry.SetText("d")
			typeKeys(entry, keyShiftRightDown, gui.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "d", entry.Text)
			assert.Equal(t, "d", submission)
		})
		t.Run("MultiLine_ShiftReturn", func(t *testing.T) {
			entry.MultiLine = true
			entry.SetText("e")
			typeKeys(entry, keyShiftLeftDown, gui.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "e", entry.Text)
			assert.Equal(t, "e", submission)
			entry.SetText("f")
			typeKeys(entry, keyShiftRightDown, gui.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "f", entry.Text)
			assert.Equal(t, "f", submission)
		})
	})
	t.Run("NoCallback", func(t *testing.T) {
		entry := &widget.Entry{}
		resetEntry := func() {
			entry.SetText("")
		}
		t.Run("SingleLine_Enter", func(t *testing.T) {
			resetEntry()
			entry.MultiLine = false
			entry.SetText("a")
			entry.TypedKey(&gui.KeyEvent{Name: gui.KeyEnter})
			assert.Equal(t, "a", entry.Text)
		})
		t.Run("SingleLine_Return", func(t *testing.T) {
			resetEntry()
			entry.MultiLine = false
			entry.SetText("b")
			entry.TypedKey(&gui.KeyEvent{Name: gui.KeyReturn})
			assert.Equal(t, "b", entry.Text)
		})
		t.Run("MultiLine_ShiftEnter", func(t *testing.T) {
			resetEntry()
			entry.MultiLine = true
			entry.SetText("c")
			typeKeys(entry, keyShiftLeftDown, gui.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "\nc", entry.Text)
			entry.SetText("d")
			entry.CursorRow = 0
			entry.CursorColumn = 0
			typeKeys(entry, keyShiftRightDown, gui.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "\nd", entry.Text)
		})
		t.Run("MultiLine_ShiftReturn", func(t *testing.T) {
			resetEntry()
			entry.MultiLine = true
			entry.SetText("e")
			typeKeys(entry, keyShiftLeftDown, gui.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "\ne", entry.Text)
			entry.SetText("f")
			entry.CursorRow = 0
			entry.CursorColumn = 0
			typeKeys(entry, keyShiftRightDown, gui.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "\nf", entry.Text)
		})
	})
}

func TestTabable(t *testing.T) {
	entry := &widget.Entry{}
	t.Run("Multiline_Tab_Default", func(t *testing.T) {
		entry.MultiLine = true
		entry.SetText("a")
		typeKeys(entry, gui.KeyTab)
		assert.Equal(t, "\ta", entry.Text)
	})
	t.Run("Singleline_Tab_Default", func(t *testing.T) {
		entry.MultiLine = false
		assert.False(t, entry.AcceptsTab())
	})
}

func TestEntry_TappedSecondary(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	tapPos := gui.NewPos(20, 10)
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_full_menu.xml", c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Disable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_read_menu.xml", c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Password = true
	entry.Refresh()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_no_password_menu.xml", c)
	assert.Nil(t, c.Overlays().Top(), "No popup for disabled password")

	entry.Enable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_password_menu.xml", c)
	assert.Equal(t, 1, len(c.Overlays().List()))
}

func TestEntry_TextWrap(t *testing.T) {
	for name, tt := range map[string]struct {
		multiLine bool
		want      string
		wrap      gui.TextWrap
	}{
		"single line WrapOff": {
			want: "entry/wrap_single_line_off.xml",
		},
		"single line Truncate": {
			wrap: gui.TextTruncate,
			want: "entry/wrap_single_line_truncate.xml",
		},
		// Disallowed - fallback to TextWrapTruncate (horizontal)
		"single line WrapBreak": {
			wrap: gui.TextWrapBreak,
			want: "entry/wrap_single_line_truncate.xml",
		},
		// Disallowed - fallback to TextWrapTruncate (horizontal)
		"single line WrapWord": {
			wrap: gui.TextWrapWord,
			want: "entry/wrap_single_line_truncate.xml",
		},
		"multi line WrapOff": {
			multiLine: true,
			want:      "entry/wrap_multi_line_off.xml",
		},
		// Disallowed - fallback to TextWrapOff
		"multi line Truncate": {
			multiLine: true,
			wrap:      gui.TextTruncate,
			want:      "entry/wrap_multi_line_truncate.xml",
		},
		"multi line WrapBreak": {
			multiLine: true,
			wrap:      gui.TextWrapBreak,
			want:      "entry/wrap_multi_line_wrap_break.xml",
		},
		"multi line WrapWord": {
			multiLine: true,
			wrap:      gui.TextWrapWord,
			want:      "entry/wrap_multi_line_wrap_word.xml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			e, window := setupImageTest(t, tt.multiLine)
			defer teardownImageTest(window)
			c := window.Canvas()

			c.Focus(e)
			e.Wrapping = tt.wrap
			if tt.multiLine {
				e.SetText("A long text on short words w/o NLs or LFs.")
			} else {
				e.SetText("Testing Wrapping")
			}
			test.AssertRendersToMarkup(t, tt.want, c)
		})
	}
}

func TestEntry_TextWrap_Changed(t *testing.T) {
	e, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.Wrapping = gui.TextWrapOff
	e.SetText("Testing Wrapping")
	test.AssertRendersToMarkup(t, "entry/wrap_single_line_off.xml", c)

	e.Wrapping = gui.TextTruncate
	e.Refresh()
	test.AssertRendersToMarkup(t, "entry/wrap_single_line_truncate.xml", c)

	e.Wrapping = gui.TextWrapOff
	e.Refresh()
	test.AssertRendersToMarkup(t, "entry/wrap_single_line_off.xml", c)
}

func TestMultiLineEntry_MinSize(t *testing.T) {
	entry := widget.NewEntry()
	singleMin := entry.MinSize()

	multi := widget.NewMultiLineEntry()
	multiMin := multi.MinSize()

	assert.Equal(t, singleMin.Width, multiMin.Width)
	assert.True(t, multiMin.Height > singleMin.Height)

	multi.MultiLine = false
	multiMin = multi.MinSize()
	assert.Equal(t, singleMin.Height, multiMin.Height)
}

func TestNewEntryWithData(t *testing.T) {
	str := binding.NewString()
	err := str.Set("Init")
	assert.Nil(t, err)

	entry := widget.NewEntryWithData(str)
	waitForBinding()
	assert.Equal(t, "Init", entry.Text)

	entry.SetText("Typed")
	v, err := str.Get()
	assert.Nil(t, err)
	assert.Equal(t, "Typed", v)
}

func TestPasswordEntry_ActionItemSizeAndPlacement(t *testing.T) {
	e := widget.NewEntry()
	b := widget.NewButton("", func() {})
	b.Icon = theme.CancelIcon()
	e.ActionItem = b
	test.WidgetRenderer(e).Layout(e.MinSize())
	assert.Equal(t, gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize()), b.Size())
	assert.Equal(t, gui.NewPos(e.MinSize().Width-2*theme.Padding()-b.Size().Width, 2*theme.Padding()), b.Position())
}

func TestPasswordEntry_NewlineIgnored(t *testing.T) {
	entry := widget.NewPasswordEntry()
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

func TestPasswordEntry_Obfuscation(t *testing.T) {
	entry, window := setupPasswordTest(t)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertRendersToMarkup(t, "password_entry/obfuscation_typed.xml", c)
}

func TestPasswordEntry_Placeholder(t *testing.T) {
	entry, window := setupPasswordTest(t)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetPlaceHolder("Password")
	test.AssertRendersToMarkup(t, "password_entry/placeholder_initial.xml", c)

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertRendersToMarkup(t, "password_entry/placeholder_typed.xml", c)
}

func TestPasswordEntry_Reveal(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	t.Run("NewPasswordEntry constructor", func(t *testing.T) {
		entry := widget.NewPasswordEntry()
		window := test.NewWindow(entry)
		defer window.Close()
		window.Resize(gui.NewSize(150, 100))
		entry.Resize(entry.MinSize().Max(gui.NewSize(130, 0)))
		entry.Move(gui.NewPos(10, 10))
		c := window.Canvas()

		test.AssertRendersToMarkup(t, "password_entry/initial.xml", c)

		c.Focus(entry)
		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/revealed.xml", c)
		assert.Equal(t, entry, c.Focused())

		// update the Password field
		entry.Password = true
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		tapPos := gui.NewPos(140-theme.Padding()*2-theme.IconInlineSize()/2, 10+entry.Size().Height/2)
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/revealed.xml", c)
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)
		assert.Equal(t, entry, c.Focused())
	})

	// This test cover backward compatibility use case when on an Entry widget
	// the Password field is set to true.
	// In this case the action item will be set when the renderer is created.
	t.Run("Entry with Password field", func(t *testing.T) {
		entry := &widget.Entry{Password: true, Wrapping: gui.TextWrapWord}
		entry.Refresh()
		window := test.NewWindow(entry)
		defer window.Close()
		window.Resize(gui.NewSize(150, 100))
		entry.Resize(entry.MinSize().Max(gui.NewSize(130, 0)))
		entry.Move(gui.NewPos(10, 10))
		c := window.Canvas()

		test.AssertRendersToMarkup(t, "password_entry/initial.xml", c)

		c.Focus(entry)
		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/revealed.xml", c)
		assert.Equal(t, entry, c.Focused())
	})
}

func TestSingleLineEntry_NewlineIgnored(t *testing.T) {
	entry := &widget.Entry{MultiLine: false}
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

const (
	entryOffset = 10

	keyShiftLeftDown  gui.KeyName = "LeftShiftDown"
	keyShiftLeftUp    gui.KeyName = "LeftShiftUp"
	keyShiftRightDown gui.KeyName = "RightShiftDown"
	keyShiftRightUp   gui.KeyName = "RightShiftUp"
)

var typeKeys = func(e *widget.Entry, keys ...gui.KeyName) {
	var keyDown = func(key *gui.KeyEvent) {
		e.KeyDown(key)
		e.TypedKey(key)
	}

	for _, key := range keys {
		switch key {
		case keyShiftLeftDown:
			keyDown(&gui.KeyEvent{Name: desktop.KeyShiftLeft})
		case keyShiftLeftUp:
			e.KeyUp(&gui.KeyEvent{Name: desktop.KeyShiftLeft})
		case keyShiftRightDown:
			keyDown(&gui.KeyEvent{Name: desktop.KeyShiftRight})
		case keyShiftRightUp:
			e.KeyUp(&gui.KeyEvent{Name: desktop.KeyShiftRight})
		default:
			keyDown(&gui.KeyEvent{Name: key})
			e.KeyUp(&gui.KeyEvent{Name: key})
		}
	}
}

func checkNewlineIgnored(t *testing.T, entry *widget.Entry) {
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &gui.KeyEvent{Name: gui.KeyDown}
	entry.TypedKey(down)
	assert.Equal(t, 0, entry.CursorRow)

	// return is ignored, do nothing
	ret := &gui.KeyEvent{Name: gui.KeyReturn}
	entry.TypedKey(ret)
	assert.Equal(t, 0, entry.CursorRow)

	up := &gui.KeyEvent{Name: gui.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func setupImageTest(t *testing.T, multiLine bool) (*widget.Entry, gui.Window) {
	test.NewApp()

	var entry *widget.Entry
	if multiLine {
		entry = &widget.Entry{MultiLine: true, Wrapping: gui.TextWrapWord}
	} else {
		entry = &widget.Entry{Wrapping: gui.TextWrapOff}
	}
	w := test.NewWindow(entry)
	w.Resize(gui.NewSize(150, 200))

	if multiLine {
		entry.Resize(gui.NewSize(120, 100))
	} else {
		entry.Resize(entry.MinSize().Max(gui.NewSize(120, 0)))
	}
	entry.Move(gui.NewPos(10, 10))

	if multiLine {
		test.AssertRendersToMarkup(t, "entry/initial_multiline.xml", w.Canvas())
	} else {
		test.AssertRendersToMarkup(t, "entry/initial.xml", w.Canvas())
	}

	return entry, w
}

func setupPasswordTest(t *testing.T) (*widget.Entry, gui.Window) {
	test.NewApp()

	entry := widget.NewPasswordEntry()
	w := test.NewWindow(entry)
	w.Resize(gui.NewSize(150, 100))

	entry.Resize(entry.MinSize().Max(gui.NewSize(130, 0)))
	entry.Move(gui.NewPos(entryOffset, entryOffset))

	test.AssertRendersToMarkup(t, "password_entry/initial.xml", w.Canvas())

	return entry, w
}

// Selects "sti" on line 2 of a new multiline
// T e s t i n g
// T e[s t i]n g
// T e s t i n g
func setupSelection(t *testing.T, reverse bool) (*widget.Entry, gui.Window) {
	e, window := setupImageTest(t, true)
	e.SetText("Testing\nTesting\nTesting")
	c := window.Canvas()
	c.Focus(e)
	if reverse {
		e.CursorRow = 1
		e.CursorColumn = 5
		typeKeys(e, keyShiftLeftDown, gui.KeyLeft, gui.KeyLeft, gui.KeyLeft)
		test.AssertRendersToMarkup(t, "entry/selection_initial_reverse.xml", c)
		assert.Equal(t, "sti", e.SelectedText())
	} else {
		e.CursorRow = 1
		e.CursorColumn = 2
		typeKeys(e, keyShiftLeftDown, gui.KeyRight, gui.KeyRight, gui.KeyRight)
		test.AssertRendersToMarkup(t, "entry/selection_initial.xml", c)
		assert.Equal(t, "sti", e.SelectedText())
	}

	return e, window
}

func teardownImageTest(w gui.Window) {
	w.Close()
	test.NewApp()
}

func waitForBinding() {
	time.Sleep(time.Millisecond * 100) // data resolves on background thread
}

// clickCanvas is an analogue of test.TapCanvas that also sends MouseDown/MouseUp events
func clickCanvas(c gui.Canvas, pos gui.Position) {
	if o, p := findMouseable(c, pos); o != nil {
		clickPrimary(c, o.(desktop.Mouseable), &gui.PointEvent{AbsolutePosition: pos, Position: p})
	}
}

func findMouseable(c gui.Canvas, pos gui.Position) (o gui.CanvasObject, p gui.Position) {
	matches := func(object gui.CanvasObject) bool {
		_, ok := object.(desktop.Mouseable)
		return ok
	}
	o, p, _ = driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	return
}

func clickPrimary(c gui.Canvas, obj desktop.Mouseable, ev *gui.PointEvent) {
	handleFocusOnTap(c, obj)
	mouseEvent := &desktop.MouseEvent{
		PointEvent: *ev,
		Button:     desktop.MouseButtonPrimary,
	}
	obj.MouseDown(mouseEvent)
	obj.MouseUp(mouseEvent)
	if tap, ok := obj.(gui.Tappable); ok {
		tap.Tapped(ev)
	}
}

func handleFocusOnTap(c gui.Canvas, obj interface{}) {
	if c == nil {
		return
	}
	unfocus := true
	if focus, ok := obj.(gui.Focusable); ok {
		if dis, ok := obj.(gui.Disableable); !ok || !dis.Disabled() {
			unfocus = false
			if focus != c.Focused() {
				unfocus = true
			}
		}
	}
	if unfocus {
		c.Unfocus()
	}
}
