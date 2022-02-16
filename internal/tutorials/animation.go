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
	"image/color"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func makeAnimationScreen(_ gui.Window) gui.CanvasObject {
	curves := makeAnimationCurves()
	curves.Move(gui.NewPos(0, 140+theme.Padding()))
	return gui.NewContainerWithoutLayout(makeAnimationCanvas(), curves)
}

func makeAnimationCanvas() gui.CanvasObject {
	rect := canvas.NewRectangle(color.Black)
	rect.Resize(gui.NewSize(410, 140))

	a := canvas.NewColorRGBAAnimation(theme.PrimaryColorNamed(theme.ColorBlue), theme.PrimaryColorNamed(theme.ColorGreen),
		time.Second*3, func(c color.Color) {
			rect.FillColor = c
			canvas.Refresh(rect)
		})
	a.RepeatCount = gui.AnimationRepeatForever
	a.AutoReverse = true
	a.Start()

	var a2 *gui.Animation
	i := widget.NewIcon(theme.CheckButtonCheckedIcon())
	a2 = canvas.NewPositionAnimation(gui.NewPos(0, 0), gui.NewPos(350, 80), time.Second*3, func(p gui.Position) {
		i.Move(p)

		width := 10 + (p.X / 7)
		i.Resize(gui.NewSize(width, width))
	})
	a2.RepeatCount = gui.AnimationRepeatForever
	a2.AutoReverse = true
	a2.Curve = gui.AnimationLinear
	a2.Start()

	running := true
	var toggle *widget.Button
	toggle = widget.NewButton("Stop", func() {
		if running {
			a.Stop()
			a2.Stop()
			toggle.SetText("Start")
		} else {
			a.Start()
			a2.Start()
			toggle.SetText("Stop")
		}
		running = !running
	})
	toggle.Resize(toggle.MinSize())
	toggle.Move(gui.NewPos(152, 54))
	return gui.NewContainerWithoutLayout(rect, i, toggle)
}

func makeAnimationCurves() gui.CanvasObject {
	label1, box1, a1 := makeAnimationCurveItem("EaseInOut", gui.AnimationEaseInOut, 0)
	label2, box2, a2 := makeAnimationCurveItem("EaseIn", gui.AnimationEaseIn, 30+theme.Padding())
	label3, box3, a3 := makeAnimationCurveItem("EaseOut", gui.AnimationEaseOut, 60+theme.Padding()*2)
	label4, box4, a4 := makeAnimationCurveItem("Linear", gui.AnimationLinear, 90+theme.Padding()*3)

	start := widget.NewButton("Compare", func() {
		a1.Start()
		a2.Start()
		a3.Start()
		a4.Start()
	})
	start.Resize(start.MinSize())
	start.Move(gui.NewPos(0, 120+theme.Padding()*4))
	return gui.NewContainerWithoutLayout(label1, label2, label3, label4, box1, box2, box3, box4, start)
}

func makeAnimationCurveItem(label string, curve gui.AnimationCurve, yOff float32) (
	text *widget.Label, box gui.CanvasObject, anim *gui.Animation) {
	text = widget.NewLabel(label)
	text.Alignment = gui.TextAlignCenter
	text.Resize(gui.NewSize(380, 30))
	text.Move(gui.NewPos(0, yOff))
	box = newThemedBox()
	box.Resize(gui.NewSize(30, 30))
	box.Move(gui.NewPos(0, yOff))

	anim = canvas.NewPositionAnimation(
		gui.NewPos(0, yOff), gui.NewPos(380, yOff), time.Second, func(p gui.Position) {
			box.Move(p)
			box.Refresh()
		})
	anim.Curve = curve
	anim.AutoReverse = true
	anim.RepeatCount = 1
	return
}

// themedBox is a simple box that change its background color according
// to the selected theme
type themedBox struct {
	widget.BaseWidget
}

func newThemedBox() *themedBox {
	b := &themedBox{}
	b.ExtendBaseWidget(b)
	return b
}

func (b *themedBox) CreateRenderer() gui.WidgetRenderer {
	b.ExtendBaseWidget(b)
	bg := canvas.NewRectangle(theme.ForegroundColor())
	return &themedBoxRenderer{bg: bg, objects: []gui.CanvasObject{bg}}
}

type themedBoxRenderer struct {
	bg      *canvas.Rectangle
	objects []gui.CanvasObject
}

func (r *themedBoxRenderer) Destroy() {
}

func (r *themedBoxRenderer) Layout(size gui.Size) {
	r.bg.Resize(size)
}

func (r *themedBoxRenderer) MinSize() gui.Size {
	return r.bg.MinSize()
}

func (r *themedBoxRenderer) Objects() []gui.CanvasObject {
	return r.objects
}

func (r *themedBoxRenderer) Refresh() {
	r.bg.FillColor = theme.ForegroundColor()
	r.bg.Refresh()
}
