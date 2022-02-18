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
)

// ShortcutHandler is a default implementation of the shortcut handler
// for the canvasObject
type ShortcutHandler struct {
	entry sync.Map // map[string]func(Shortcut)
}

// TypedShortcut handle the registered shortcut
func (sh *ShortcutHandler) TypedShortcut(shortcut Shortcut) {
	val, ok := sh.entry.Load(shortcut.ShortcutName())
	if !ok {
		return
	}

	f := val.(func(Shortcut))
	f(shortcut)
}

// AddShortcut register a handler to be executed when the shortcut action is triggered
func (sh *ShortcutHandler) AddShortcut(shortcut Shortcut, handler func(shortcut Shortcut)) {
	sh.entry.Store(shortcut.ShortcutName(), handler)
}

// RemoveShortcut removes a registered shortcut
func (sh *ShortcutHandler) RemoveShortcut(shortcut Shortcut) {
	sh.entry.Delete(shortcut.ShortcutName())
}

// Shortcut is the interface used to describe a shortcut action
type Shortcut interface {
	ShortcutName() string
}

// KeyboardShortcut describes a shortcut meant to be triggered by a keyboard action.
type KeyboardShortcut interface {
	Shortcut
	Key() KeyName
	Mod() KeyModifier
}

// ShortcutPaste describes a shortcut paste action.
type ShortcutPaste struct {
	Clipboard Clipboard
}

var _ KeyboardShortcut = (*ShortcutPaste)(nil)

// Key returns the KeyName for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutPaste) Key() KeyName {
	return KeyV
}

// Mod returns the KeyModifier for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutPaste) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutPaste) ShortcutName() string {
	return "Paste"
}

// ShortcutCopy describes a shortcut copy action.
type ShortcutCopy struct {
	Clipboard Clipboard
}

var _ KeyboardShortcut = (*ShortcutCopy)(nil)

// Key returns the KeyName for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutCopy) Key() KeyName {
	return KeyC
}

// Mod returns the KeyModifier for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutCopy) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutCopy) ShortcutName() string {
	return "Copy"
}

// ShortcutCut describes a shortcut cut action.
type ShortcutCut struct {
	Clipboard Clipboard
}

var _ KeyboardShortcut = (*ShortcutCut)(nil)

// Key returns the KeyName for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutCut) Key() KeyName {
	return KeyX
}

// Mod returns the KeyModifier for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutCut) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutCut) ShortcutName() string {
	return "Cut"
}

// ShortcutSelectAll describes a shortcut selectAll action.
type ShortcutSelectAll struct{}

var _ KeyboardShortcut = (*ShortcutSelectAll)(nil)

// Key returns the KeyName for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutSelectAll) Key() KeyName {
	return KeyA
}

// Mod returns the KeyModifier for this shortcut.
//
// Implements: KeyboardShortcut
func (se *ShortcutSelectAll) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutSelectAll) ShortcutName() string {
	return "SelectAll"
}
