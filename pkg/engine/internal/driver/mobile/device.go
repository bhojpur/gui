package mobile

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
	"github.com/bhojpur/gui/pkg/engine/driver/mobile"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/size"

	gui "github.com/bhojpur/gui/pkg/engine"
)

type device struct {
	safeTop, safeLeft, safeWidth, safeHeight int
}

//lint:file-ignore U1000 Var currentDPI is used in other files, but not here
var (
	currentOrientation size.Orientation
	currentDPI         float32
)

// Declare conformity with Device
var _ gui.Device = (*device)(nil)

func (*device) Orientation() gui.DeviceOrientation {
	switch currentOrientation {
	case size.OrientationLandscape:
		return gui.OrientationHorizontalLeft
	default:
		return gui.OrientationVertical
	}
}

func (*device) IsMobile() bool {
	return true
}

func (*device) IsBrowser() bool {
	return false
}

func (*device) HasKeyboard() bool {
	return false
}

func (*device) ShowVirtualKeyboard() {
	showVirtualKeyboard(mobile.DefaultKeyboard)
}

func (*device) ShowVirtualKeyboardType(keyboard mobile.KeyboardType) {
	showVirtualKeyboard(keyboard)
}

func (*device) HideVirtualKeyboard() {
	hideVirtualKeyboard()
}
