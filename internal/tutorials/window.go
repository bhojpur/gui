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
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func windowScreen(_ gui.Window) gui.CanvasObject {
	windowGroup := container.NewVBox(
		widget.NewButton("New window", func() {
			w := gui.CurrentApp().NewWindow("Hello")
			w.SetContent(widget.NewLabel("Hello World!"))
			w.Show()
		}),
		widget.NewButton("Fixed size window", func() {
			w := gui.CurrentApp().NewWindow("Fixed")
			w.SetContent(gui.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("Hello World!")))

			w.Resize(gui.NewSize(240, 180))
			w.SetFixedSize(true)
			w.Show()
		}),
		widget.NewButton("Toggle between fixed/not fixed window size", func() {
			w := gui.CurrentApp().NewWindow("Toggle fixed size")
			w.SetContent(gui.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewCheck("Fixed size", func(toggle bool) {
				if toggle {
					w.Resize(gui.NewSize(240, 180))
				}
				w.SetFixedSize(toggle)
			})))
			w.Show()
		}),
		widget.NewButton("Centered window", func() {
			w := gui.CurrentApp().NewWindow("Central")
			w.SetContent(gui.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("Hello World!")))

			w.CenterOnScreen()
			w.Show()
		}))

	drv := gui.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		windowGroup.Objects = append(windowGroup.Objects,
			widget.NewButton("Splash Window (only use on start)", func() {
				w := drv.CreateSplashWindow()
				w.SetContent(widget.NewLabelWithStyle("Hello World!\n\nMake a splash!",
					gui.TextAlignCenter, gui.TextStyle{Bold: true}))
				w.Show()

				go func() {
					time.Sleep(time.Second * 3)
					w.Close()
				}()
			}))
	}

	otherGroup := widget.NewCard("Other", "",
		widget.NewButton("Notification", func() {
			gui.CurrentApp().SendNotification(&gui.Notification{
				Title:   "Bhojpur Demo Application",
				Content: "Testing notifications...",
			})
		}))

	return container.NewVBox(widget.NewCard("Windows", "", windowGroup), otherGroup)
}
