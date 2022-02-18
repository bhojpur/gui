package engine

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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func syncMapLen(m *sync.Map) (n int) {
	m.Range(func(_, _ interface{}) bool {
		n++
		return true
	})
	return
}

func TestShortcutHandler_AddShortcut(t *testing.T) {
	handle := &ShortcutHandler{}

	handle.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {})
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.Equal(t, 2, syncMapLen(&handle.entry))
}

func TestShortcutHandler_RemoveShortcut(t *testing.T) {
	handler := &ShortcutHandler{}
	handler.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {})
	handler.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.Equal(t, 2, syncMapLen(&handler.entry))

	handler.RemoveShortcut(&ShortcutCopy{})

	assert.Equal(t, 1, syncMapLen(&handler.entry))

	handler.RemoveShortcut(&ShortcutPaste{})

	assert.Equal(t, 0, syncMapLen(&handler.entry))
}

func TestShortcutHandler_HandleShortcut(t *testing.T) {
	handle := &ShortcutHandler{}
	cutCalled, copyCalled, pasteCalled := false, false, false

	handle.AddShortcut(&ShortcutCut{}, func(shortcut Shortcut) {
		cutCalled = true
	})
	handle.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {
		copyCalled = true
	})
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {
		pasteCalled = true
	})

	handle.TypedShortcut(&ShortcutCut{})
	assert.True(t, cutCalled)
	handle.TypedShortcut(&ShortcutCopy{})
	assert.True(t, copyCalled)
	handle.TypedShortcut(&ShortcutPaste{})
	assert.True(t, pasteCalled)
}
