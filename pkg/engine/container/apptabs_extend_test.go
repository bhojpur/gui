package container

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
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type extendedAppTabs struct {
	AppTabs
}

func newExtendedAppTabs(items ...*TabItem) *extendedAppTabs {
	ret := &extendedAppTabs{}
	ret.ExtendBaseWidget(ret)

	ret.Items = items
	return ret
}

func TestAppTabs_Extended_Tapped(t *testing.T) {
	tabs := newExtendedAppTabs(
		NewTabItem("Test1", widget.NewLabel("Test1")),
		NewTabItem("Test2", widget.NewLabel("Test2")),
	)
	tabs.Resize(gui.NewSize(150, 150)) // Ensure AppTabs is big enough to show both tab buttons
	r := test.WidgetRenderer(tabs).(*appTabsRenderer)

	tab1 := r.bar.Objects[0].(*gui.Container).Objects[0].(*tabButton)
	tab2 := r.bar.Objects[0].(*gui.Container).Objects[1].(*tabButton)
	require.Equal(t, 0, tabs.SelectedIndex())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)

	tab2.Tapped(&gui.PointEvent{})
	assert.Equal(t, 1, tabs.SelectedIndex())
	require.Equal(t, theme.ForegroundColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab2).(*tabButtonRenderer).label.Color)
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())

	tab2.Tapped(&gui.PointEvent{})
	assert.Equal(t, 1, tabs.SelectedIndex())
	require.Equal(t, theme.ForegroundColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab2).(*tabButtonRenderer).label.Color)
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())

	tab1.Tapped(&gui.PointEvent{})
	assert.Equal(t, 0, tabs.SelectedIndex())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.ForegroundColor(), test.WidgetRenderer(tab2).(*tabButtonRenderer).label.Color)
	assert.True(t, tabs.Items[0].Content.Visible())
	assert.False(t, tabs.Items[1].Content.Visible())
}
