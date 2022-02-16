package widget

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
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// Card widget groups title, subtitle with content and a header image
//
// Since: 1.4
type Card struct {
	BaseWidget
	Title, Subtitle string
	Image           *canvas.Image
	Content         gui.CanvasObject
}

// NewCard creates a new card widget with the specified title, subtitle and content (all optional).
//
// Since: 1.4
func NewCard(title, subtitle string, content gui.CanvasObject) *Card {
	card := &Card{
		Title:    title,
		Subtitle: subtitle,
		Content:  content,
	}

	card.ExtendBaseWidget(card)
	return card
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (c *Card) CreateRenderer() gui.WidgetRenderer {
	c.ExtendBaseWidget(c)

	header := canvas.NewText(c.Title, theme.ForegroundColor())
	header.TextStyle.Bold = true
	subHeader := canvas.NewText(c.Subtitle, theme.ForegroundColor())

	objects := []gui.CanvasObject{header, subHeader}
	if c.Image != nil {
		objects = append(objects, c.Image)
	}
	if c.Content != nil {
		objects = append(objects, c.Content)
	}
	r := &cardRenderer{widget.NewShadowingRenderer(objects, widget.CardLevel),
		header, subHeader, c}
	r.applyTheme()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (c *Card) MinSize() gui.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetContent changes the body of this card to have the specified content.
func (c *Card) SetContent(obj gui.CanvasObject) {
	c.Content = obj

	c.Refresh()
}

// SetImage changes the image displayed above the title for this card.
func (c *Card) SetImage(img *canvas.Image) {
	c.Image = img

	c.Refresh()
}

// SetSubTitle updates the secondary title for this card.
func (c *Card) SetSubTitle(text string) {
	c.Subtitle = text

	c.Refresh()
}

// SetTitle updates the main title for this card.
func (c *Card) SetTitle(text string) {
	c.Title = text

	c.Refresh()
}

type cardRenderer struct {
	*widget.ShadowingRenderer

	header, subHeader *canvas.Text

	card *Card
}

const (
	cardMediaHeight = 128
)

// Layout the components of the card container.
func (c *cardRenderer) Layout(size gui.Size) {
	pos := gui.NewPos(theme.Padding()/2, theme.Padding()/2)
	size = size.Subtract(gui.NewSize(theme.Padding(), theme.Padding()))
	c.LayoutShadow(size, pos)

	if c.card.Image != nil {
		c.card.Image.Move(pos)
		c.card.Image.Resize(gui.NewSize(size.Width, cardMediaHeight))
		pos.Y += cardMediaHeight
	}

	contentPad := theme.Padding()
	if c.card.Title != "" || c.card.Subtitle != "" {
		titlePad := theme.Padding() * 2
		size.Width -= titlePad * 2
		pos.X += titlePad
		pos.Y += titlePad

		if c.card.Title != "" {
			height := c.header.MinSize().Height
			c.header.Move(pos)
			c.header.Resize(gui.NewSize(size.Width, height))
			pos.Y += height + theme.Padding()
		}

		if c.card.Subtitle != "" {
			height := c.subHeader.MinSize().Height
			c.subHeader.Move(pos)
			c.subHeader.Resize(gui.NewSize(size.Width, height))
			pos.Y += height + theme.Padding()
		}

		size.Width = size.Width + titlePad*2
		pos.X = pos.X - titlePad
		pos.Y += titlePad
	}

	size.Width -= contentPad * 2
	pos.X += contentPad
	if c.card.Content != nil {
		height := size.Height - contentPad*2 - (pos.Y - theme.Padding()/2) // adjust for content and initial offset
		if c.card.Title != "" || c.card.Subtitle != "" {
			height += contentPad
			pos.Y -= contentPad
		}
		c.card.Content.Move(pos.Add(gui.NewPos(0, contentPad)))
		c.card.Content.Resize(gui.NewSize(size.Width, height))
	}
}

// MinSize calculates the minimum size of a card.
// This is based on the contained text, image and content.
func (c *cardRenderer) MinSize() gui.Size {
	hasHeader := c.card.Title != ""
	hasSubHeader := c.card.Subtitle != ""
	hasImage := c.card.Image != nil
	hasContent := c.card.Content != nil

	if !hasHeader && !hasSubHeader && !hasContent { // just image, or nothing
		if c.card.Image == nil {
			return gui.NewSize(theme.Padding(), theme.Padding()) // empty, just space for border
		}
		return gui.NewSize(c.card.Image.MinSize().Width+theme.Padding(), cardMediaHeight+theme.Padding())
	}

	contentPad := theme.Padding()
	min := gui.NewSize(theme.Padding(), theme.Padding())
	if hasImage {
		min = gui.NewSize(min.Width, min.Height+cardMediaHeight)
	}

	if hasHeader || hasSubHeader {
		titlePad := theme.Padding() * 2
		min = min.Add(gui.NewSize(0, titlePad*2))
		if hasHeader {
			headerMin := c.header.MinSize()
			min = gui.NewSize(gui.Max(min.Width, headerMin.Width+titlePad*2+theme.Padding()),
				min.Height+headerMin.Height)
			if hasSubHeader {
				min.Height += theme.Padding()
			}
		}
		if hasSubHeader {
			subHeaderMin := c.subHeader.MinSize()
			min = gui.NewSize(gui.Max(min.Width, subHeaderMin.Width+titlePad*2+theme.Padding()),
				min.Height+subHeaderMin.Height)
		}
	}

	if hasContent {
		contentMin := c.card.Content.MinSize()
		min = gui.NewSize(gui.Max(min.Width, contentMin.Width+contentPad*2+theme.Padding()),
			min.Height+contentMin.Height+contentPad*2)
	}

	return min
}

func (c *cardRenderer) Refresh() {
	c.header.Text = c.card.Title
	c.header.Refresh()
	c.subHeader.Text = c.card.Subtitle
	c.subHeader.Refresh()

	objects := []gui.CanvasObject{c.header, c.subHeader}
	if c.card.Image != nil {
		objects = append(objects, c.card.Image)
	}
	if c.card.Content != nil {
		objects = append(objects, c.card.Content)
	}
	c.ShadowingRenderer.SetObjects(objects)

	c.applyTheme()
	c.Layout(c.card.Size())
	c.ShadowingRenderer.RefreshShadow()
	canvas.Refresh(c.card.super())
}

// applyTheme updates this button to match the current theme
func (c *cardRenderer) applyTheme() {
	if c.header != nil {
		c.header.TextSize = theme.TextHeadingSize()
		c.header.Color = theme.ForegroundColor()
	}
	if c.subHeader != nil {
		c.subHeader.TextSize = theme.TextSize()
		c.subHeader.Color = theme.ForegroundColor()
	}
}
