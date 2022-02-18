package desktop

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
	"runtime"
	"strings"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// Declare conformity with Shortcut interface
var _ gui.Shortcut = (*CustomShortcut)(nil)
var _ gui.KeyboardShortcut = (*CustomShortcut)(nil)

// CustomShortcut describes a shortcut desktop event.
type CustomShortcut struct {
	gui.KeyName
	Modifier gui.KeyModifier
}

// Key returns the key name of this shortcut.
// @implements KeyboardShortcut
func (cs *CustomShortcut) Key() gui.KeyName {
	return cs.KeyName
}

// Mod returns the modifier of this shortcut.
// @implements KeyboardShortcut
func (cs *CustomShortcut) Mod() gui.KeyModifier {
	return cs.Modifier
}

// ShortcutName returns the shortcut name associated to the event
func (cs *CustomShortcut) ShortcutName() string {
	id := &strings.Builder{}
	id.WriteString("CustomDesktop:")
	id.WriteString(modifierToString(cs.Modifier))
	id.WriteString("+")
	id.WriteString(string(cs.KeyName))
	return id.String()
}

func modifierToString(mods gui.KeyModifier) string {
	s := []string{}
	if (mods & gui.KeyModifierShift) != 0 {
		s = append(s, string("Shift"))
	}
	if (mods & gui.KeyModifierControl) != 0 {
		s = append(s, string("Control"))
	}
	if (mods & gui.KeyModifierAlt) != 0 {
		s = append(s, string("Alt"))
	}
	if (mods & gui.KeyModifierSuper) != 0 {
		if runtime.GOOS == "darwin" {
			s = append(s, string("Command"))
		} else {
			s = append(s, string("Super"))
		}
	}
	return strings.Join(s, "+")
}
