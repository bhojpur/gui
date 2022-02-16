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
	"net/url"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHyperlink_MinSize(t *testing.T) {
	u, err := url.Parse("https://bhojpur.net/")
	assert.Nil(t, err)

	hyperlink := NewHyperlink("Test", u)
	hyperlink.CreateRenderer()
	hyperlink.provider.CreateRenderer()
	minA := hyperlink.MinSize()

	assert.Less(t, theme.Padding()*2, minA.Width)

	hyperlink.SetText("Longer")
	minB := hyperlink.MinSize()
	assert.Less(t, minA.Width, minB.Width)

	hyperlink.Text = "."
	hyperlink.Refresh()
	minC := hyperlink.MinSize()
	assert.Greater(t, minB.Width, minC.Width)
}

func TestHyperlink_Cursor(t *testing.T) {
	u, err := url.Parse("https://bhojpur.net/")
	hyperlink := NewHyperlink("Test", u)

	assert.Nil(t, err)
	assert.Equal(t, desktop.PointerCursor, hyperlink.Cursor())
}

func TestHyperlink_Alignment(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test", Alignment: gui.TextAlignTrailing}
	hyperlink.CreateRenderer()
	assert.Equal(t, gui.TextAlignTrailing, textRenderTexts(hyperlink.provider)[0].Alignment)
}

func TestHyperlink_Hide(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test"}
	hyperlink.CreateRenderer()
	hyperlink.Hide()
	hyperlink.Refresh()

	assert.True(t, hyperlink.Hidden)
	assert.False(t, hyperlink.provider.Hidden) // we don't propagate hide

	hyperlink.Show()
	assert.False(t, hyperlink.Hidden)
	assert.False(t, hyperlink.provider.Hidden)
}

func TestHyperlink_Focus(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	hyperlink := &Hyperlink{Text: "Test"}
	w := test.NewWindow(hyperlink)
	w.SetPadded(false)
	defer w.Close()
	w.Resize(hyperlink.MinSize())

	test.AssertImageMatches(t, "hyperlink/initial.png", w.Canvas().Capture())
	hyperlink.FocusGained()
	test.AssertImageMatches(t, "hyperlink/focus.png", w.Canvas().Capture())
	hyperlink.FocusLost()
	test.AssertImageMatches(t, "hyperlink/initial.png", w.Canvas().Capture())
}

func TestHyperlink_Resize(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test"}
	hyperlink.CreateRenderer()
	size := gui.NewSize(100, 20)
	hyperlink.Resize(size)

	assert.Equal(t, size, hyperlink.Size())
	assert.Equal(t, size, hyperlink.provider.Size())
}

func TestHyperlink_SetText(t *testing.T) {
	u, err := url.Parse("https://bhojpur.net/")
	assert.Nil(t, err)

	hyperlink := &Hyperlink{Text: "Test", URL: u}
	hyperlink.CreateRenderer()
	hyperlink.SetText("New")

	assert.Equal(t, "New", hyperlink.Text)
	assert.Equal(t, "New", textRenderTexts(hyperlink.provider)[0].Text)
}

func TestHyperlink_SetUrl(t *testing.T) {
	sURL, err := url.Parse("https://github.com/bhojpur/gui")
	assert.Nil(t, err)

	// test constructor
	hyperlink := NewHyperlink("Test", sURL)
	assert.Equal(t, sURL, hyperlink.URL)

	// test setting functions
	sURL, err = url.Parse("https://bhojpur.net")
	assert.Nil(t, err)
	hyperlink.SetURL(sURL)
	assert.Equal(t, sURL, hyperlink.URL)
}

func TestHyperlink_CreateRendererDoesNotAffectSize(t *testing.T) {
	u, err := url.Parse("https://github.com/bhojpur/gui")
	require.NoError(t, err)
	link := NewHyperlink("Test", u)
	link.Resize(link.MinSize())
	size := link.Size()
	assert.NotEqual(t, gui.NewSize(0, 0), size)
	assert.Equal(t, size, link.MinSize())

	r := link.CreateRenderer()
	link.provider.CreateRenderer()
	assert.Equal(t, size, link.Size())
	assert.Equal(t, size, link.MinSize())
	assert.Equal(t, size, r.MinSize())
	r.Layout(size)
	assert.Equal(t, size, link.Size())
	assert.Equal(t, size, link.MinSize())
}
