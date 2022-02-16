package playground

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

// It provides tooling for running Bhojpur GUI applications inside the Go playground

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/software"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

func imageToPlayground(img image.Image) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		gui.LogError("Failed to encode image", err)
		return
	}

	enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println("IMAGE:" + enc)
}

// RenderCanvas takes a canvas and converts it into an inline image for showing in the playground
func RenderCanvas(c gui.Canvas) {
	imageToPlayground(software.RenderCanvas(c, theme.DarkTheme()))
}

// RenderWindow takes a window and converts it's canvas into an inline image for showing in the playground
func RenderWindow(w gui.Window) {
	RenderCanvas(w.Canvas())
}

// Render takes a canvasobject and converts it into an inline image for showing in the playground
func Render(obj gui.CanvasObject) {
	imageToPlayground(software.Render(obj, theme.DarkTheme()))
}
