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
	"strconv"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func scaleString(c gui.Canvas) string {
	return strconv.FormatFloat(float64(c.Scale()), 'f', 2, 32)
}

func texScaleString(c gui.Canvas) string {
	pixels, _ := c.PixelCoordinateForPosition(gui.NewPos(1, 1))
	texScale := float32(pixels) / c.Scale()
	return strconv.FormatFloat(float64(texScale), 'f', 2, 32)
}

func prependTo(g *gui.Container, s string) {
	g.Objects = append([]gui.CanvasObject{widget.NewLabel(s)}, g.Objects...)
	g.Refresh()
}

func setScaleText(scale, tex *widget.Label, win gui.Window) {
	for scale.Visible() {
		scale.SetText(scaleString(win.Canvas()))
		tex.SetText(texScaleString(win.Canvas()))

		time.Sleep(time.Second)
	}
}

// advancedScreen loads a panel that shows details and settings that are a bit
// more detailed than normally needed.
func advancedScreen(win gui.Window) gui.CanvasObject {
	scale := widget.NewLabel("")
	tex := widget.NewLabel("")

	screen := widget.NewCard("Screen info", "", widget.NewForm(
		&widget.FormItem{Text: "Scale", Widget: scale},
		&widget.FormItem{Text: "Texture Scale", Widget: tex},
	))

	go setScaleText(scale, tex, win)

	label := widget.NewLabel("Just type...")
	generic := container.NewVBox()
	desk := container.NewVBox()

	genericCard := widget.NewCard("", "Generic", container.NewVScroll(generic))
	deskCard := widget.NewCard("", "Desktop", container.NewVScroll(desk))

	win.Canvas().SetOnTypedRune(func(r rune) {
		prependTo(generic, "Rune: "+string(r))
	})
	win.Canvas().SetOnTypedKey(func(ev *gui.KeyEvent) {
		prependTo(generic, "Key : "+string(ev.Name))
	})
	if deskCanvas, ok := win.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(ev *gui.KeyEvent) {
			prependTo(desk, "KeyDown: "+string(ev.Name))
		})
		deskCanvas.SetOnKeyUp(func(ev *gui.KeyEvent) {
			prependTo(desk, "KeyUp  : "+string(ev.Name))
		})
	}

	return container.NewHBox(
		container.NewVBox(screen,
			widget.NewButton("Custom Theme", func() {
				gui.CurrentApp().Settings().SetTheme(newCustomTheme())
			}),
			widget.NewButton("Fullscreen", func() {
				win.SetFullScreen(!win.FullScreen())
			}),
		),
		container.NewBorder(label, nil, nil, nil,
			container.NewGridWithColumns(2, genericCard, deskCard),
		),
	)
}
