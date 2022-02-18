package render

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"
)

// Create the test entry with 3 completion items.
func createEntry() *CompletionEntry {
	entry := NewCompletionEntry([]string{"zoo", "boo"})
	entry.OnChanged = func(s string) {
		data := []string{"foo", "bar", "baz"}
		entry.SetOptions(data)
		entry.ShowCompletion()
	}
	return entry
}

// Check if the data is filled with corresponding options.
func TestCompletionEntry(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.SetText("init")
	assert.Equal(t, 3, len(entry.Options))
}

// Show the completion menu
func TestCompletionEntry_ShowMenu(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.SetText("init")
	assert.True(t, entry.popupMenu.Visible())

}

// Navigate in menu and select an entry.
func TestCompletionEntry_Navigate(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.SetText("init")

	// navigate to "bar" and press enter, the entry should contain
	// "bar" in value
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyReturn})

	assert.Equal(t, "bar", entry.Text)
	assert.False(t, entry.popupMenu.Visible())
}

// Ensure the cursor is set to the end of entry after completion.
func TestCompletionEntry_CursorPosition(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.OnChanged = func(s string) {
		entry.SetOptions([]string{"foofoo", "barbar", "bazbaz"})
		entry.ShowCompletion()
	}
	entry.SetText("barb")

	// navigate to "bar" and press enter, the entry should contain
	// "bar" in value
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyReturn})

	assert.Equal(t, 6, entry.CursorColumn)

}

// Hide the menu on Escape key.
func TestCompletionEntry_Escape(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.SetText("init")

	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyEscape})

	assert.False(t, entry.popupMenu.Visible())
}

// Hide the menu on rune pressed.
func TestCompletionEntry_Rune(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.SetText("foobar")
	entry.CursorColumn = 6 // place the cursor after the text

	// type some chars...
	win.Canvas().Focused().TypedRune('x')
	win.Canvas().Focused().TypedRune('y')
	assert.Equal(t, "foobarxy", entry.Text)

	// make a move and type other char
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	win.Canvas().Focused().TypedRune('z')
	assert.Equal(t, "foobarxyz", entry.Text)

	assert.True(t, entry.popupMenu.Visible())
}

// Hide the menu on rune pressed.
func TestCompletionEntry_Rotation(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.SetText("foobar")
	entry.CursorColumn = 6 // place the cursor after the text

	// loop one time (nb items + 1) to go back to the first element
	for i := 0; i <= len(entry.Options); i++ {
		win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown})
	}

	assert.Equal(t, 0, entry.navigableList.selected)

	// Do the same in reverse order, here, onlh one time to go on the last item
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyUp})
	assert.Equal(t, len(entry.Options)-1, entry.navigableList.selected)
}

// Test if completion is hidden when there is no options.
func TestCompletionEntry_WithEmptyOptions(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	entry.OnChanged = func(s string) {
		entry.SetOptions([]string{})
		entry.ShowCompletion()
	}

	entry.SetText("foo")
	assert.Nil(t, entry.popupMenu) // popupMenu should not being created
}

// Test sumbission with opened completion.
func TestCompletionEntry_OnSubmit(t *testing.T) {
	entry := createEntry()
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	submitted := false
	entry.OnSubmitted = func(s string) {
		submitted = true
		entry.HideCompletion()
		assert.True(t, entry.popupMenu.Hidden)
	}
	entry.OnChanged = func(s string) {
		entry.ShowCompletion()
	}

	entry.SetText("foo")
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyReturn})
	assert.True(t, submitted)
}

// Test double submission issue, when the user select an option in list and press "Enter", then
// the "OnSubmitted" method should not be called. It should be called only after the user pressed a
// second time.
func TestCompletionEntry_DoubleSubmissionIssue(t *testing.T) {
	entry := createEntry()
	entry.SetOptions([]string{"foofoo", "bar", "baz"})
	win := test.NewWindow(entry)
	win.Resize(gui.NewSize(500, 300))
	defer win.Close()

	submitted := false
	entry.OnSubmitted = func(s string) {
		submitted = true
	}

	win.Canvas().Focus(entry)
	entry.SetText("foo")
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyDown}) // select foofoo
	assert.False(t, submitted)
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyReturn}) // OnSubmitted should NOT be called
	assert.False(t, submitted)
	assert.False(t, entry.popupMenu.Visible())
	win.Canvas().Focused().TypedKey(&gui.KeyEvent{Name: gui.KeyReturn}) // OnSubmitted should be called
	assert.True(t, submitted)
}
