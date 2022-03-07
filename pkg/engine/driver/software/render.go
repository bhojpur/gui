package software

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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
)

// RenderCanvas takes a canvas and renders it to a regular Go image using the provided Theme.
// This is the same as setting the application theme and then calling Canvas.Capture().
func RenderCanvas(c gui.Canvas, t gui.Theme) image.Image {
	gui.CurrentApp().Settings().SetTheme(t)
	app.ApplyThemeTo(c.Content(), c)

	return c.Capture()
}

// Render takes a canvas object and renders it to a regular Go image using the provided Theme.
// The returned image will be set to the object's minimum size.
// Use the theme.LightTheme() or theme.DarkTheme() to access the builtin themes.
func Render(obj gui.CanvasObject, t gui.Theme) image.Image {
	c := NewCanvas()
	c.SetPadded(false)
	c.SetContent(obj)

	gui.CurrentApp().Settings().SetTheme(t)
	app.ApplyThemeTo(obj, c)
	return c.Capture()
}
