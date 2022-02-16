// It defines an event for touch input.
package touch // import "github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/touch"

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

// The best source on android input events is the NDK: include/android/input.h
//
// iOS event handling guide:
// https://developer.apple.com/library/ios/documentation/EventHandling/Conceptual/EventHandlingiPhoneOS

import (
	"fmt"
)

// Event is a touch event.
type Event struct {
	// X and Y are the touch location, in pixels.
	X, Y float32

	// Sequence is the sequence number. The same number is shared by all events
	// in a sequence. A sequence begins with a single TypeBegin, is followed by
	// zero or more TypeMoves, and ends with a single TypeEnd. A Sequence
	// distinguishes concurrent sequences but its value is subsequently reused.
	Sequence Sequence

	// Type is the touch type.
	Type Type
}

// Sequence identifies a sequence of touch events.
type Sequence int64

// Type describes the type of a touch event.
type Type byte

const (
	// TypeBegin is a user first touching the device.
	//
	// On Android, this is a AMOTION_EVENT_ACTION_DOWN.
	// On iOS, this is a call to touchesBegan.
	TypeBegin Type = iota

	// TypeMove is a user dragging across the device.
	//
	// A TypeMove is delivered between a TypeBegin and TypeEnd.
	//
	// On Android, this is a AMOTION_EVENT_ACTION_MOVE.
	// On iOS, this is a call to touchesMoved.
	TypeMove

	// TypeEnd is a user no longer touching the device.
	//
	// On Android, this is a AMOTION_EVENT_ACTION_UP.
	// On iOS, this is a call to touchesEnded.
	TypeEnd
)

func (t Type) String() string {
	switch t {
	case TypeBegin:
		return "begin"
	case TypeMove:
		return "move"
	case TypeEnd:
		return "end"
	}
	return fmt.Sprintf("touch.Type(%d)", t)
}
