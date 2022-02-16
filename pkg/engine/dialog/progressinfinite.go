package dialog

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
	"image/color"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// ProgressInfiniteDialog is a simple dialog window that displays text and a infinite progress bar.
//
// Deprecated: Create a new custom dialog with a widget.ProgressBarInfinite() inside.
type ProgressInfiniteDialog struct {
	*dialog

	bar *widget.ProgressBarInfinite
}

// NewProgressInfinite creates a infinite progress dialog and returns the handle.
// Using the returned type you should call Show().
//
// Deprecated: Create a new custom dialog with a widget.ProgressBarInfinite() inside.
func NewProgressInfinite(title, message string, parent gui.Window) *ProgressInfiniteDialog {
	d := newDialog(title, message, theme.InfoIcon(), nil /*cancel?*/, parent)
	bar := widget.NewProgressBarInfinite()
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(gui.NewSize(200, 0))

	d.setButtons(container.NewMax(rect, bar))
	return &ProgressInfiniteDialog{d, bar}
}

// Hide this dialog and stop the infinite progress goroutine
func (d *ProgressInfiniteDialog) Hide() {
	d.bar.Hide()
	d.dialog.Hide()
}
