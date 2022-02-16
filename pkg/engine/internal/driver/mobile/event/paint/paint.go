// It defines an event for the app being ready to paint.
package paint // import "github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/paint"

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

// Event indicates that the app is ready to paint the next frame of the GUI.
//
//A frame is completed by calling the App's Publish method.
type Event struct {
	// External is true for paint events sent by the screen driver.
	//
	// An external event may be sent at any time in response to an
	// operating system event, for example the window opened, was
	// resized, or the screen memory was lost.
	//
	// Programs actively drawing to the screen as fast as vsync allows
	// should ignore external paint events to avoid a backlog of paint
	// events building up.
	External bool
}
