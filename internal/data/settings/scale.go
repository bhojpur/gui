package settings

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
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

type scaleItems struct {
	scale   float32
	name    string
	preview *canvas.Text
	button  *widget.Button
}

var scales = []*scaleItems{
	{scale: 0.5, name: "Tiny"},
	{scale: 0.8, name: "Small"},
	{scale: 1, name: "Normal"},
	{scale: 1.3, name: "Large"},
	{scale: 1.8, name: "Huge"}}

func (s *Settings) appliedScale(value float32) {
	for _, scale := range scales {
		scale.preview.TextSize = theme.TextSize() * scale.scale / value
	}
}

func (s *Settings) chooseScale(value float32) {
	s.guiSettings.Scale = value

	for _, scale := range scales {
		if scale.scale == value {
			scale.button.Importance = widget.HighImportance
		} else {
			scale.button.Importance = widget.MediumImportance
		}

		scale.button.Refresh()
	}
}

func (s *Settings) makeScaleButtons() []gui.CanvasObject {
	var buttons = make([]gui.CanvasObject, len(scales))
	for i, scale := range scales {
		value := scale.scale
		button := widget.NewButton(scale.name, func() {
			s.chooseScale(value)
		})
		if s.guiSettings.Scale == scale.scale {
			button.Importance = widget.HighImportance
		}

		scale.button = button
		buttons[i] = button
	}

	return buttons
}

func (s *Settings) makeScaleGroup(scale float32) *widget.Card {
	scalePreviewBox := gui.NewContainerWithLayout(layout.NewGridLayout(5), s.makeScalePreviews(scale)...)
	scaleBox := gui.NewContainerWithLayout(layout.NewGridLayout(5), s.makeScaleButtons()...)

	return widget.NewCard("Scale", "", container.NewVBox(scalePreviewBox, scaleBox, newRefreshMonitor(s)))
}

func (s *Settings) makeScalePreviews(value float32) []gui.CanvasObject {
	var previews = make([]gui.CanvasObject, len(scales))
	for i, scale := range scales {
		text := canvas.NewText("A", theme.ForegroundColor())
		text.Alignment = gui.TextAlignCenter
		text.TextSize = theme.TextSize() * scale.scale / value

		scale.preview = text
		previews[i] = text
	}

	return previews
}

func (s *Settings) refreshScalePreviews() {
	for _, scale := range scales {
		scale.preview.Color = theme.ForegroundColor()
	}
}

// refreshMonitor is a simple widget that updates canvas components when the UI is asked to refresh.
// Captures theme and scale changes without the settings monitoring code.
type refreshMonitor struct {
	widget.Label
	settings *Settings
}

func (r *refreshMonitor) Refresh() {
	r.settings.refreshScalePreviews()
	r.Label.Refresh()
}

func newRefreshMonitor(s *Settings) *refreshMonitor {
	r := &refreshMonitor{settings: s}
	r.Hide()
	return r
}
