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
	"github.com/bhojpur/gui/pkg/engine/test"

	"github.com/stretchr/testify/assert"
)

type extendEntry struct {
	Entry
}

func TestEntry_Password_Extended_CreateRenderer(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	w := a.NewWindow("")
	entry := &extendEntry{}
	entry.ExtendBaseWidget(entry)
	entry.Password = true
	entry.Wrapping = gui.TextTruncate
	assert.NotNil(t, test.WidgetRenderer(entry))
	r := test.WidgetRenderer(entry).(*entryRenderer).scroll.Content.(*entryContent)
	p := test.WidgetRenderer(r).(*entryContentRenderer).provider

	w.SetContent(entry)

	test.Type(entry, "Pass")
	texts := test.WidgetRenderer(p).(*textRenderer).Objects()
	assert.Equal(t, passwordChar+passwordChar+passwordChar+passwordChar, texts[0].(*canvas.Text).Text)
	assert.NotNil(t, entry.ActionItem)
	test.Tap(entry.ActionItem.(*passwordRevealer))

	texts = test.WidgetRenderer(p).(*textRenderer).Objects()
	assert.Equal(t, "Pass", texts[0].(*canvas.Text).Text)
	assert.Equal(t, entry, w.Canvas().Focused())
}
