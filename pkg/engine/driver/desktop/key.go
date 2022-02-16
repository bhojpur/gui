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
	gui "github.com/bhojpur/gui/pkg/engine"
)

const (
	// KeyNone represents no key
	KeyNone gui.KeyName = ""
	// KeyShiftLeft represents the left shift key
	KeyShiftLeft gui.KeyName = "LeftShift"
	// KeyShiftRight represents the right shift key
	KeyShiftRight gui.KeyName = "RightShift"
	// KeyControlLeft represents the left control key
	KeyControlLeft gui.KeyName = "LeftControl"
	// KeyControlRight represents the right control key
	KeyControlRight gui.KeyName = "RightControl"
	// KeyAltLeft represents the left alt key
	KeyAltLeft gui.KeyName = "LeftAlt"
	// KeyAltRight represents the right alt key
	KeyAltRight gui.KeyName = "RightAlt"
	// KeySuperLeft represents the left "Windows" key (or "Command" key on macOS)
	KeySuperLeft gui.KeyName = "LeftSuper"
	// KeySuperRight represents the right "Windows" key (or "Command" key on macOS)
	KeySuperRight gui.KeyName = "RightSuper"
	// KeyMenu represents the left or right menu / application key
	KeyMenu gui.KeyName = "Menu"
	// KeyPrintScreen represents the key used to cause a screen capture
	KeyPrintScreen gui.KeyName = "PrintScreen"

	// KeyCapsLock represents the caps lock key, tapping once is the down event then again is the up
	KeyCapsLock gui.KeyName = "CapsLock"
)

// Modifier captures any key modifiers (shift etc) pressed during this key event
type Modifier int

const (
	// ShiftModifier represents a shift key being held
	ShiftModifier Modifier = 1 << iota
	// ControlModifier represents the ctrl key being held
	ControlModifier
	// AltModifier represents either alt keys being held
	AltModifier
	// SuperModifier represents either super keys being held
	SuperModifier
)

// Keyable describes any focusable canvas object that can accept desktop key events.
// This is the traditional key down and up event that is not applicable to all devices.
type Keyable interface {
	gui.Focusable

	KeyDown(*gui.KeyEvent)
	KeyUp(*gui.KeyEvent)
}
