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
	"image/color"
	"net/url"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var _ gui.Focusable = (*Hyperlink)(nil)
var _ gui.Widget = (*Hyperlink)(nil)

// Hyperlink widget is a text component with appropriate padding and layout.
// When clicked, the default web browser should open with a URL
type Hyperlink struct {
	BaseWidget
	Text      string
	URL       *url.URL
	Alignment gui.TextAlign // The alignment of the Text
	Wrapping  gui.TextWrap  // The wrapping of the Text
	TextStyle gui.TextStyle // The style of the hyperlink text

	focused, hovered bool
	provider         *RichText
}

// NewHyperlink creates a new hyperlink widget with the set text content
func NewHyperlink(text string, url *url.URL) *Hyperlink {
	return NewHyperlinkWithStyle(text, url, gui.TextAlignLeading, gui.TextStyle{})
}

// NewHyperlinkWithStyle creates a new hyperlink widget with the set text content
func NewHyperlinkWithStyle(text string, url *url.URL, alignment gui.TextAlign, style gui.TextStyle) *Hyperlink {
	hl := &Hyperlink{
		Text:      text,
		URL:       url,
		Alignment: alignment,
		TextStyle: style,
	}

	return hl
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (hl *Hyperlink) CreateRenderer() gui.WidgetRenderer {
	hl.ExtendBaseWidget(hl)
	hl.provider = NewRichTextWithText(hl.Text)
	hl.provider.ExtendBaseWidget(hl.provider)
	hl.syncSegments()

	focus := canvas.NewRectangle(color.Transparent)
	focus.StrokeColor = theme.FocusColor()
	focus.StrokeWidth = 2
	focus.Hide()
	under := canvas.NewRectangle(theme.PrimaryColor())
	under.Hide()
	return &hyperlinkRenderer{hl: hl, objects: []gui.CanvasObject{hl.provider, focus, under}, focus: focus, under: under}
}

// Cursor returns the cursor type of this widget
func (hl *Hyperlink) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (hl *Hyperlink) FocusGained() {
	hl.focused = true
	hl.BaseWidget.Refresh()
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (hl *Hyperlink) FocusLost() {
	hl.focused = false
	hl.BaseWidget.Refresh()
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (hl *Hyperlink) MouseIn(*desktop.MouseEvent) {
	hl.hovered = true
	hl.BaseWidget.Refresh()
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (hl *Hyperlink) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (hl *Hyperlink) MouseOut() {
	hl.hovered = false
	hl.BaseWidget.Refresh()
}

// Refresh triggers a redraw of the hyperlink.
//
// Implements: gui.Widget
func (hl *Hyperlink) Refresh() {
	if hl.provider == nil { // not created until visible
		return
	}
	hl.syncSegments()

	hl.provider.Refresh()
	hl.BaseWidget.Refresh()
}

// MinSize returns the smallest size this widget can shrink to
func (hl *Hyperlink) MinSize() gui.Size {
	if hl.provider == nil {
		hl.CreateRenderer()
	}

	return hl.provider.MinSize()
}

// Resize sets a new size for the hyperlink.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (hl *Hyperlink) Resize(size gui.Size) {
	hl.BaseWidget.Resize(size)
	if hl.provider == nil { // not created until visible
		return
	}
	hl.provider.Resize(size)
}

// SetText sets the text of the hyperlink
func (hl *Hyperlink) SetText(text string) {
	hl.Text = text
	if hl.provider == nil { // not created until visible
		return
	}
	hl.syncSegments()
	hl.provider.Refresh()
}

// SetURL sets the URL of the hyperlink, taking in a URL type
func (hl *Hyperlink) SetURL(url *url.URL) {
	hl.URL = url
}

// SetURLFromString sets the URL of the hyperlink, taking in a string type
func (hl *Hyperlink) SetURLFromString(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	hl.URL = u
	return nil
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (hl *Hyperlink) Tapped(*gui.PointEvent) {
	hl.openURL()
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (hl *Hyperlink) TypedRune(rune) {
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (hl *Hyperlink) TypedKey(ev *gui.KeyEvent) {
	if ev.Name == gui.KeySpace {
		hl.openURL()
	}
}

func (hl *Hyperlink) openURL() {
	if hl.URL != nil {
		err := gui.CurrentApp().OpenURL(hl.URL)
		if err != nil {
			gui.LogError("Failed to open url", err)
		}
	}
}

func (hl *Hyperlink) syncSegments() {
	hl.provider.Wrapping = hl.Wrapping
	hl.provider.Segments = []RichTextSegment{&TextSegment{
		Style: RichTextStyle{
			Alignment: hl.Alignment,
			ColorName: theme.ColorNamePrimary,
			Inline:    true,
			TextStyle: hl.TextStyle,
		},
		Text: hl.Text,
	}}
}

var _ gui.WidgetRenderer = (*hyperlinkRenderer)(nil)

type hyperlinkRenderer struct {
	hl    *Hyperlink
	focus *canvas.Rectangle
	under *canvas.Rectangle

	objects []gui.CanvasObject
}

func (r *hyperlinkRenderer) Destroy() {
}

func (r *hyperlinkRenderer) Layout(s gui.Size) {
	r.hl.provider.Resize(s)
	r.focus.Move(gui.NewPos(theme.Padding(), theme.Padding()))
	r.focus.Resize(gui.NewSize(s.Width-theme.Padding()*2, s.Height-theme.Padding()*2))
	r.under.Move(gui.NewPos(theme.Padding()*2, s.Height-theme.Padding()*2))
	r.under.Resize(gui.NewSize(s.Width-theme.Padding()*4, 1))
}

func (r *hyperlinkRenderer) MinSize() gui.Size {
	return r.hl.provider.MinSize()
}

func (r *hyperlinkRenderer) Objects() []gui.CanvasObject {
	return r.objects
}

func (r *hyperlinkRenderer) Refresh() {
	r.hl.provider.Refresh()
	r.focus.StrokeColor = theme.FocusColor()
	r.focus.Hidden = !r.hl.focused
	r.under.StrokeColor = theme.PrimaryColor()
	r.under.Hidden = !r.hl.hovered
}
