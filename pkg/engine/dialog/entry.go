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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// EntryDialog is a variation of a dialog which prompts the user to enter some text.
//
// Deprecated: Use dialog.NewFormDialog() or dialog.ShowFormDialog() with a widget.Entry inside instead.
type EntryDialog struct {
	*formDialog

	entry *widget.Entry

	onClosed func()
}

// SetText changes the current text value of the entry dialog, this can
// be useful for setting a default value.
func (i *EntryDialog) SetText(s string) {
	i.entry.SetText(s)
}

// SetPlaceholder defines the placeholder text for the entry
func (i *EntryDialog) SetPlaceholder(s string) {
	i.entry.SetPlaceHolder(s)
}

// SetOnClosed changes the callback which is run when the dialog is closed,
// which is nil by default.
//
// The callback is called unconditionally whether the user confirms or cancels.
//
// Note that the callback will be called after onConfirm, if both are non-nil.
// This way onConfirm can potential modify state that this callback needs to
// get the user input when the user confirms, while also being able to handle
// the case where the user cancelled.
func (i *EntryDialog) SetOnClosed(callback func()) {
	i.onClosed = callback
}

// NewEntryDialog creates a dialog over the specified window for the user to enter a value.
//
// onConfirm is a callback that runs when the user enters a string of
// text and clicks the "confirm" button. May be nil.
//
// Deprecated: Use dialog.NewFormDialog() with a widget.Entry inside instead.
func NewEntryDialog(title, message string, onConfirm func(string), parent gui.Window) *EntryDialog {
	i := &EntryDialog{entry: widget.NewEntry()}
	items := []*widget.FormItem{widget.NewFormItem(message, i.entry)}
	i.formDialog = NewForm(title, "Ok", "Cancel", items, func(ok bool) {
		// User has confirmed and entered an input
		if ok && onConfirm != nil {
			onConfirm(i.entry.Text)
		}

		if i.onClosed != nil {
			i.onClosed()
		}

		i.entry.Text = ""
		i.win.Hide() // Close directly without executing the callback. This is the callback.
	}, parent).(*formDialog)

	return i
}

// ShowEntryDialog creates a new entry dialog and shows it immediately.
//
// Deprecated: Use dialog.ShowFormDialog() with a widget.Entry inside instead.
func ShowEntryDialog(title, message string, onConfirm func(string), parent gui.Window) {
	NewEntryDialog(title, message, onConfirm, parent).Show()
}
