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

import gui "github.com/bhojpur/gui/pkg/engine"

// MouseButton represents a single button in a desktop MouseEvent
type MouseButton int

const (
	// MouseButtonPrimary is the most common mouse button - on some systems the only one.
	// This will normally be on the left side of a mouse.
	//
	// Since: 2.0
	MouseButtonPrimary MouseButton = 1 << iota

	// MouseButtonSecondary is the secondary button on most mouse input devices.
	// This will normally be on the right side of a mouse.
	//
	// Since: 2.0
	MouseButtonSecondary

	// MouseButtonTertiary is the middle button on the mouse, assuming it has one.
	//
	// Since: 2.0
	MouseButtonTertiary

	// LeftMouseButton is the most common mouse button - on some systems the only one.
	//
	// Deprecated: use MouseButtonPrimary which will adapt to mouse configuration.
	LeftMouseButton = MouseButtonPrimary

	// RightMouseButton is the secondary button on most mouse input devices.
	//
	// Deprecated: use MouseButtonSecondary which will adapt to mouse configuration.
	RightMouseButton = MouseButtonSecondary
)

// MouseEvent contains data relating to desktop mouse events
type MouseEvent struct {
	gui.PointEvent
	Button   MouseButton
	Modifier gui.KeyModifier
}

// Mouseable represents desktop mouse events that can be sent to CanvasObjects
type Mouseable interface {
	MouseDown(*MouseEvent)
	MouseUp(*MouseEvent)
}

// Hoverable is used when a canvas object wishes to know if a pointer device moves over it.
type Hoverable interface {
	// MouseIn is a hook that is called if the mouse pointer enters the element.
	MouseIn(*MouseEvent)
	// MouseMoved is a hook that is called if the mouse pointer moved over the element.
	MouseMoved(*MouseEvent)
	// MouseOut is a hook that is called if the mouse pointer leaves the element.
	MouseOut()
}
