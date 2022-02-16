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
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/stretchr/testify/assert"
)

func TestToolbarSize(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer(), NewToolbarAction(theme.HomeIcon(), func() {}))
	assert.Equal(t, 2, len(toolbar.Items))
	size := toolbar.MinSize()

	toolbar.Items = append(toolbar.Items, &toolbarLabel{NewLabel("Hi")})
	toolbar.Refresh()
	assert.Equal(t, size.Height, toolbar.MinSize().Height)
	assert.Greater(t, toolbar.MinSize().Width, size.Width)
}

func TestToolbar_Apppend(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer())
	assert.Equal(t, 1, len(toolbar.Items))

	added := NewToolbarAction(theme.ContentCutIcon(), func() {})
	toolbar.Append(added)
	assert.Equal(t, 2, len(toolbar.Items))
	assert.Equal(t, added, toolbar.Items[1])
}

func TestToolbar_Prepend(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer())
	assert.Equal(t, 1, len(toolbar.Items))

	prepend := NewToolbarAction(theme.ContentCutIcon(), func() {})
	toolbar.Prepend(prepend)
	assert.Equal(t, 2, len(toolbar.Items))
	assert.Equal(t, prepend, toolbar.Items[0])
}

func TestToolbar_Replace(t *testing.T) {
	icon := theme.ContentCutIcon()
	toolbar := NewToolbar(NewToolbarAction(icon, func() {}))
	assert.Equal(t, 1, len(toolbar.Items))
	render := test.WidgetRenderer(toolbar)
	assert.Equal(t, icon, render.Objects()[0].(*Button).Icon)

	toolbar.Items[0] = NewToolbarAction(theme.HelpIcon(), func() {})
	toolbar.Refresh()
	assert.NotEqual(t, icon, render.Objects()[0].(*Button).Icon)
}

func TestToolbar_ItemPositioning(t *testing.T) {
	toolbar := &Toolbar{
		Items: []ToolbarItem{
			NewToolbarAction(theme.ContentCopyIcon(), func() {}),
			NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		},
	}
	toolbar.ExtendBaseWidget(toolbar)
	toolbar.Refresh()
	var items []gui.CanvasObject
	for _, o := range test.LaidOutObjects(toolbar) {
		if b, ok := o.(*Button); ok {
			items = append(items, b)
		}
	}
	if assert.Equal(t, 2, len(items)) {
		assert.Equal(t, gui.NewPos(0, 0), items[0].Position())
		assert.Equal(t, gui.NewPos(40, 0), items[1].Position())
	}
}

type toolbarLabel struct {
	*Label
}

func (t *toolbarLabel) ToolbarObject() gui.CanvasObject {
	return t.Label
}
