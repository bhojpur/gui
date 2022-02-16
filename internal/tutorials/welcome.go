package tutorials

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
	"net/url"

	"github.com/bhojpur/gui/internal/data"
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		gui.LogError("Could not parse URL", err)
	}

	return link
}

func welcomeScreen(_ gui.Window) gui.CanvasObject {
	logo := canvas.NewImageFromResource(data.BhojpurScene)
	logo.FillMode = canvas.ImageFillContain
	if gui.CurrentDevice().IsMobile() {
		logo.SetMinSize(gui.NewSize(171, 125))
	} else {
		logo.SetMinSize(gui.NewSize(228, 167))
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to Bhojpur GUI demo application", gui.TextAlignCenter, gui.TextStyle{Bold: true}),
		logo,
		container.NewHBox(
			widget.NewHyperlink("bhojpur.net", parseURL("https://bhojpur.net/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("documentation", parseURL("https://docs.bhojpur.net/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("sponsor", parseURL("https://bhojpur.net/sponsor/")),
		),
	))
}
