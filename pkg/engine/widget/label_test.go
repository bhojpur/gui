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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/data/binding"
	"github.com/bhojpur/gui/pkg/engine/internal/painter/software"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestLabel_Binding(t *testing.T) {
	label := NewLabel("Init")
	assert.Equal(t, "Init", label.Text)

	str := binding.NewString()
	label.Bind(str)
	waitForBinding()
	assert.Equal(t, "", label.Text)

	str.Set("Updated")
	waitForBinding()
	assert.Equal(t, "Updated", label.Text)

	label.Unbind()
	waitForBinding()
	assert.Equal(t, "Updated", label.Text)
}

func TestLabel_Hide(t *testing.T) {
	label := NewLabel("Test")
	label.CreateRenderer()
	label.Hide()
	label.Refresh()

	assert.True(t, label.Hidden)
	assert.False(t, label.provider.Hidden) // we don't propagate hide

	label.Show()
	assert.False(t, label.Hidden)
	assert.False(t, label.provider.Hidden)
}

func TestLabel_MinSize(t *testing.T) {
	label := NewLabel("Test")
	minA := label.MinSize()

	assert.Less(t, theme.Padding()*2, minA.Width)

	label.SetText("Longer")
	minB := label.MinSize()
	assert.Less(t, minA.Width, minB.Width)

	label.Text = "."
	label.Refresh()
	minC := label.MinSize()
	assert.Greater(t, minB.Width, minC.Width)
}

func TestLabel_Resize(t *testing.T) {
	label := NewLabel("Test")
	label.CreateRenderer()
	size := gui.NewSize(100, 20)
	label.Resize(size)

	assert.Equal(t, size, label.Size())
	assert.Equal(t, size, label.provider.Size())

	label.SetText("Longer")
	assert.Equal(t, size, label.Size())
	assert.Equal(t, size, label.provider.Size())
}

func TestLabel_Text(t *testing.T) {
	label := &Label{Text: "Test"}
	label.Refresh()

	assert.Equal(t, "Test", label.Text)
	assert.Equal(t, "Test", textRenderTexts(label)[0].Text)
}

func TestLabel_Text_Refresh(t *testing.T) {
	label := &Label{Text: ""}

	assert.Equal(t, "", label.Text)
	assert.Equal(t, "", textRenderTexts(label)[0].Text)

	label.Text = "Test"
	label.Refresh()
	assert.Equal(t, "Test", label.Text)
	assert.Equal(t, "Test", textRenderTexts(label)[0].Text)
}

func TestLabel_SetText(t *testing.T) {
	label := &Label{Text: "Test"}
	label.SetText("Crashy")
	label.Refresh()
	label.SetText("New")

	assert.Equal(t, "New", label.Text)
	assert.Equal(t, "New", textRenderTexts(label)[0].Text)
}

func TestLabel_Alignment(t *testing.T) {
	label := &Label{Text: "Test", Alignment: gui.TextAlignTrailing}
	label.Refresh()

	assert.Equal(t, gui.TextAlignTrailing, textRenderTexts(label)[0].Alignment)
}

func TestLabel_Alignment_Later(t *testing.T) {
	label := &Label{Text: "Test"}
	label.Refresh()
	assert.Equal(t, gui.TextAlignLeading, textRenderTexts(label)[0].Alignment)

	label.Alignment = gui.TextAlignTrailing
	label.Refresh()
	assert.Equal(t, gui.TextAlignTrailing, textRenderTexts(label)[0].Alignment)
}

func TestText_MinSize_MultiLine(t *testing.T) {
	textOneLine := NewLabel("Break")
	min := textOneLine.MinSize()
	textMultiLine := NewLabel("Bre\nak")
	min2 := textMultiLine.MinSize()

	assert.True(t, min2.Width < min.Width)
	assert.True(t, min2.Height > min.Height)

	yPos := float32(-1)
	for _, text := range test.WidgetRenderer(textMultiLine).(*textRenderer).Objects() {
		assert.True(t, text.Size().Height < min2.Height)
		assert.True(t, text.Position().Y > yPos)
		yPos = text.Position().Y
	}
}

func TestText_MinSizeAdjustsWithContent(t *testing.T) {
	text := NewLabel("Line 1\nLine 2\n")
	initialMin := text.MinSize()

	text.SetText("Line 1\nLine 2\nLonger Line\n")
	assert.Greater(t, text.MinSize().Width, initialMin.Width)
	assert.Greater(t, text.MinSize().Height, initialMin.Height)

	text.SetText("Line 1\nLine 2\n")
	assert.Equal(t, initialMin, text.MinSize())
	assert.Equal(t, initialMin, text.provider.MinSize())
}

func TestLabel_ApplyTheme(t *testing.T) {
	text := NewLabel("Line 1")
	text.Hide()

	render := test.WidgetRenderer(text).(*textRenderer)
	assert.Equal(t, theme.ForegroundColor(), render.Objects()[0].(*canvas.Text).Color)
	text.Show()
	assert.Equal(t, theme.ForegroundColor(), render.Objects()[0].(*canvas.Text).Color)
}

func TestLabel_CreateRendererDoesNotAffectSize(t *testing.T) {
	text := NewLabel("Hello")
	text.Resize(text.MinSize())
	size := text.Size()
	assert.NotEqual(t, gui.NewSize(0, 0), size)
	assert.Equal(t, size, text.MinSize())

	r := text.CreateRenderer()
	assert.Equal(t, size, text.Size())
	assert.Equal(t, size, text.MinSize())
	assert.Equal(t, size, r.MinSize())
	r.Layout(size)
	assert.Equal(t, size, text.Size())
	assert.Equal(t, size, text.MinSize())
}

func TestLabel_ChangeTruncate(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetPadded(false)
	text := NewLabel("Hello")
	c.SetContent(text)
	c.Resize(text.MinSize())
	test.AssertRendersToMarkup(t, "label/default.xml", c)

	truncSize := text.MinSize().Subtract(gui.NewSize(10, 0))
	text.Resize(truncSize)
	text.Wrapping = gui.TextTruncate
	text.Refresh()
	test.AssertRendersToMarkup(t, "label/truncate.xml", c)
}

func TestNewLabelWithData(t *testing.T) {
	str := binding.NewString()
	str.Set("Init")

	label := NewLabelWithData(str)
	waitForBinding()
	assert.Equal(t, "Init", label.Text)
}
