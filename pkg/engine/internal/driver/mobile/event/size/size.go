// It defines an event for the dimensions, physical resolution and
// orientation of the app's window.
package size // import "github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/size"

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
	"image"
)

// Event holds the dimensions, physical resolution and orientation of the app's
// window.
type Event struct {
	// WidthPx and HeightPx are the window's dimensions in pixels.
	WidthPx, HeightPx int

	// WidthPt and HeightPt are the window's physical dimensions in points
	// (1/72 of an inch).
	//
	// The values are based on PixelsPerPt and are therefore approximate, as
	// per the comment on PixelsPerPt.
	WidthPt, HeightPt float32

	// PixelsPerPt is the window's physical resolution. It is the number of
	// pixels in a single float32.
	//
	// There are a wide variety of pixel densities in existing phones and
	// tablets, so apps should be written to expect various non-integer
	// PixelsPerPt values.
	//
	// The value is approximate, in that the OS, drivers or hardware may report
	// approximate or quantized values. An N x N pixel square should be roughly
	// 1 square inch for N = int(PixelsPerPt * 72), although different square
	// lengths (in pixels) might be closer to 1 inch in practice. Nonetheless,
	// this PixelsPerPt value should be consistent with e.g. the ratio of
	// WidthPx to WidthPt.
	PixelsPerPt float32

	// Orientation is the orientation of the device screen.
	Orientation Orientation

	// InsetTopPx, InsetBottomPx, InsetLeftPx and InsetRightPx define the size of any border area in pixels.
	// These values define how far in from the screen edge any controls should be drawn.
	// The inset can be caused by status bars, button overlays or devices cutouts.
	InsetTopPx, InsetBottomPx, InsetLeftPx, InsetRightPx int

	// DarkMode is set to true if this window is currently shown in the OS configured dark / night mode.
	DarkMode bool
}

// Size returns the window's size in pixels, at the time this size event was
// sent.
func (e Event) Size() image.Point {
	return image.Point{e.WidthPx, e.HeightPx}
}

// Bounds returns the window's bounds in pixels, at the time this size event
// was sent.
//
// The top-left pixel is always (0, 0). The bottom-right pixel is given by the
// width and height.
func (e Event) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Point{e.WidthPx, e.HeightPx}}
}

// Orientation is the orientation of the device screen.
type Orientation int

const (
	// OrientationUnknown means device orientation cannot be determined.
	//
	// Equivalent on Android to Configuration.ORIENTATION_UNKNOWN
	// and on iOS to:
	//	UIDeviceOrientationUnknown
	//	UIDeviceOrientationFaceUp
	//	UIDeviceOrientationFaceDown
	OrientationUnknown Orientation = iota

	// OrientationPortrait is a device oriented so it is tall and thin.
	//
	// Equivalent on Android to Configuration.ORIENTATION_PORTRAIT
	// and on iOS to:
	//	UIDeviceOrientationPortrait
	//	UIDeviceOrientationPortraitUpsideDown
	OrientationPortrait

	// OrientationLandscape is a device oriented so it is short and wide.
	//
	// Equivalent on Android to Configuration.ORIENTATION_LANDSCAPE
	// and on iOS to:
	//	UIDeviceOrientationLandscapeLeft
	//	UIDeviceOrientationLandscapeRight
	OrientationLandscape
)
