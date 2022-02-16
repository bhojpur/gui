package app

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
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
)

// ApplyThemeTo ensures that the specified canvasobject and all widgets and themeable objects will
// be updated for the current theme.
func ApplyThemeTo(content gui.CanvasObject, canv gui.Canvas) {
	if content == nil {
		return
	}

	switch o := content.(type) {
	case gui.Widget:
		for _, co := range cache.Renderer(o).Objects() {
			ApplyThemeTo(co, canv)
		}
		cache.Renderer(o).Layout(content.Size()) // theme can cause sizing changes
	case *gui.Container:
		for _, co := range o.Objects {
			ApplyThemeTo(co, canv)
		}
		if l := o.Layout; l != nil {
			l.Layout(o.Objects, o.Size()) // theme can cause sizing changes
		}
	}
	content.Refresh()
}

// ApplySettings ensures that all widgets and themeable objects in an application will be updated for the current theme.
// It also checks that scale changes are reflected if required
func ApplySettings(set gui.Settings, app gui.App) {
	ApplySettingsWithCallback(set, app, nil)
}

// ApplySettingsWithCallback ensures that all widgets and themeable objects in an application will be updated for the current theme.
// It also checks that scale changes are reflected if required. Also it will call `onEveryWindow` on every window
// interaction
func ApplySettingsWithCallback(set gui.Settings, app gui.App, onEveryWindow func(w gui.Window)) {
	for _, window := range app.Driver().AllWindows() {
		ApplyThemeTo(window.Content(), window.Canvas())
		for _, overlay := range window.Canvas().Overlays().List() {
			ApplyThemeTo(overlay, window.Canvas())
		}
		if onEveryWindow != nil {
			onEveryWindow(window)
		}
	}
}
