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
	"github.com/bhojpur/gui/pkg/engine/data/binding"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
)

// Label widget is a label component with appropriate padding and layout.
type Label struct {
	BaseWidget
	Text      string
	Alignment gui.TextAlign // The alignment of the Text
	Wrapping  gui.TextWrap  // The wrapping of the Text
	TextStyle gui.TextStyle // The style of the label text
	provider  *RichText

	binder basicBinder
}

// NewLabel creates a new label widget with the set text content
func NewLabel(text string) *Label {
	return NewLabelWithStyle(text, gui.TextAlignLeading, gui.TextStyle{})
}

// NewLabelWithData returns an Label widget connected to the specified data source.
//
// Since: 2.0
func NewLabelWithData(data binding.String) *Label {
	label := NewLabel("")
	label.Bind(data)

	return label
}

// NewLabelWithStyle creates a new label widget with the set text content
func NewLabelWithStyle(text string, alignment gui.TextAlign, style gui.TextStyle) *Label {
	l := &Label{
		Text:      text,
		Alignment: alignment,
		TextStyle: style,
	}

	l.ExtendBaseWidget(l)
	return l
}

// Bind connects the specified data source to this Label.
// The current value will be displayed and any changes in the data will cause the widget to update.
//
// Since: 2.0
func (l *Label) Bind(data binding.String) {
	l.binder.SetCallback(l.updateFromData) // This could only be done once, maybe in ExtendBaseWidget?
	l.binder.Bind(data)
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (l *Label) CreateRenderer() gui.WidgetRenderer {
	l.provider = NewRichTextWithText(l.Text)
	l.ExtendBaseWidget(l)
	l.syncSegments()

	return l.provider.CreateRenderer()
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (l *Label) ExtendBaseWidget(w gui.Widget) {
	if w == nil {
		w = l
	}
	l.BaseWidget.ExtendBaseWidget(w)
	if l.provider != nil {
		l.provider.ExtendBaseWidget(l.super())
	}
}

// MinSize returns the size that this label should not shrink below.
//
// Implements: gui.Widget
func (l *Label) MinSize() gui.Size {
	if l.provider == nil {
		l.ExtendBaseWidget(l)
		cache.Renderer(l.super())
	}

	return l.provider.MinSize()
}

// Refresh triggers a redraw of the label.
//
// Implements: gui.Widget
func (l *Label) Refresh() {
	if l.provider == nil { // not created until visible
		return
	}
	l.syncSegments()
	l.BaseWidget.Refresh()
	l.provider.Refresh()
}

// Resize sets a new size for the label.
// This should only be called if it is not in a container with a layout manager.
//
// Implements: gui.Widget
func (l *Label) Resize(s gui.Size) {
	l.BaseWidget.Resize(s)
	if l.provider != nil {
		l.provider.Resize(s)
	}
}

// SetText sets the text of the label
func (l *Label) SetText(text string) {
	l.Text = text
	l.Refresh()
}

// Unbind disconnects any configured data source from this Label.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (l *Label) Unbind() {
	l.binder.Unbind()
}

func (l *Label) syncSegments() {
	l.provider.Wrapping = l.Wrapping
	l.provider.Segments = []RichTextSegment{&TextSegment{
		Style: RichTextStyle{
			Alignment: l.Alignment,
			Inline:    true,
			TextStyle: l.TextStyle,
		},
		Text: l.Text,
	}}
}

func (l *Label) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.String)
	if !ok {
		return
	}
	val, err := textSource.Get()
	if err != nil {
		gui.LogError("Error getting current data value", err)
		return
	}
	l.SetText(val)
}
