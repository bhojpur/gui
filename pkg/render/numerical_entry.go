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
	"strconv"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/mobile"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// NumericalEntry is an extended entry that only allows numerical input.
// Only integers are allowed by default. Support for floats can be enabled by setting AllowFloat.
type NumericalEntry struct {
	widget.Entry
	AllowFloat bool
}

// NewNumericalEntry returns an extended entry that only allows numerical input.
func NewNumericalEntry() *NumericalEntry {
	entry := &NumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// TypedRune is called when this item receives a char event.
//
// Implements: gui.Focusable
func (e *NumericalEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
		return
	}

	if e.AllowFloat && (r == '.' || r == ',') {
		e.Entry.TypedRune(r)
	}
}

// TypedShortcut handles the registered shortcuts.
//
// Implements: gui.Shortcutable
func (e *NumericalEntry) TypedShortcut(shortcut gui.Shortcut) {
	paste, ok := shortcut.(*gui.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	if e.isNumber(paste.Clipboard.Content()) {
		e.Entry.TypedShortcut(shortcut)
	}
}

// Keyboard sets up the right keyboard to use on mobile.
//
// Implements: mobile.Keyboardable
func (e *NumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func (e *NumericalEntry) isNumber(content string) bool {
	if e.AllowFloat {
		_, err := strconv.ParseFloat(content, 64)
		return err == nil
	}

	_, err := strconv.Atoi(content)
	return err == nil
}
