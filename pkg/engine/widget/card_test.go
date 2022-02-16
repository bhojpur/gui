package widget_test

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
)

func TestCard_SetImage(t *testing.T) {
	c := widget.NewCard("Title", "sub", widget.NewLabel("Content"))
	r := test.WidgetRenderer(c)
	assert.Equal(t, 4, len(r.Objects())) // the 3 above plus shadow

	c.SetImage(canvas.NewImageFromResource(theme.BhojpurLogo()))
	assert.Equal(t, 5, len(r.Objects()))
}

func TestCard_SetContent(t *testing.T) {
	c := widget.NewCard("Title", "sub", widget.NewLabel("Content"))
	r := test.WidgetRenderer(c)
	assert.Equal(t, 4, len(r.Objects())) // the 3 above plus shadow

	newContent := widget.NewLabel("New")
	c.SetContent(newContent)
	assert.Equal(t, 4, len(r.Objects()))
	assert.Equal(t, newContent, r.Objects()[3])
}

func TestCard_Layout(t *testing.T) {
	test.NewApp()

	for name, tt := range map[string]struct {
		title, subtitle string
		icon            *canvas.Image
		content         gui.CanvasObject
	}{
		"title": {
			title:    "Title",
			subtitle: "",
			icon:     nil,
			content:  nil,
		},
		"subtitle": {
			title:    "",
			subtitle: "Subtitle",
			icon:     nil,
			content:  nil,
		},
		"titles": {
			title:    "Title",
			subtitle: "Subtitle",
			icon:     nil,
			content:  nil,
		},
		"titles_image": {
			title:    "Title",
			subtitle: "Subtitle",
			icon:     canvas.NewImageFromResource(theme.BhojpurLogo()),
			content:  nil,
		},
		"just_image": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.BhojpurLogo()),
			content:  nil,
		},
		"just_content": {
			title:    "",
			subtitle: "",
			icon:     nil,
			content:  newContentRect(),
		},
		"title_content": {
			title:    "Hello",
			subtitle: "",
			icon:     nil,
			content:  newContentRect(),
		},
		"image_content": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.BhojpurLogo()),
			content:  newContentRect(),
		},
		"all_items": {
			title:    "Longer title",
			subtitle: "subtitle with length",
			icon:     canvas.NewImageFromResource(theme.BhojpurLogo()),
			content:  newContentRect(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			card := &widget.Card{
				Title:    tt.title,
				Subtitle: tt.subtitle,
				Image:    tt.icon,
				Content:  tt.content,
			}

			window := test.NewWindow(card)
			size := card.MinSize().Max(gui.NewSize(80, 0)) // give a little width for image only tests
			window.Resize(size.Add(gui.NewSize(theme.Padding()*2, theme.Padding()*2)))
			if tt.content != nil {
				assert.Equal(t, float32(10), tt.content.Size().Height)
			}
			test.AssertRendersToMarkup(t, "card/layout_"+name+".xml", window.Canvas())

			window.Close()
		})
	}
}

func TestCard_MinSize(t *testing.T) {
	content := widget.NewLabel("simple")
	card := &widget.Card{Content: content}

	inner := card.MinSize().Subtract(gui.NewSize(theme.Padding()*3, theme.Padding()*3)) // shadow + content pad
	assert.Equal(t, content.MinSize(), inner)
}

func newContentRect() *canvas.Rectangle {
	rect := canvas.NewRectangle(color.Gray{0x66})
	rect.StrokeColor = color.Black
	rect.StrokeWidth = 2
	rect.SetMinSize(gui.NewSize(10, 10))

	return rect
}
