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
	"testing"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
)

func TestButton_MinSize(t *testing.T) {
	button := widget.NewButton("Hi", nil)
	min := button.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestButton_SetText(t *testing.T) {
	button := widget.NewButton("Hi", nil)
	min1 := button.MinSize()

	button.SetText("Longer")
	min2 := button.MinSize()

	assert.True(t, min2.Width > min1.Width)
	assert.Equal(t, min2.Height, min1.Height)
}

func TestButton_MinSize_Icon(t *testing.T) {
	button := widget.NewButton("Hi", nil)
	min1 := button.MinSize()

	button.SetIcon(theme.CancelIcon())
	min2 := button.MinSize()

	assert.True(t, min2.Width > min1.Width)
	assert.Equal(t, min2.Height, min1.Height)
}

func TestButton_Cursor(t *testing.T) {
	button := widget.NewButton("Test", nil)
	assert.Equal(t, desktop.DefaultCursor, button.Cursor())
}

func TestButton_Tapped(t *testing.T) {
	tapped := make(chan bool)
	button := widget.NewButton("Hi", func() {
		tapped <- true
	})

	go test.Tap(button)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for button tap")
		}
	}()
}

func TestButton_Disable(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	tapped := false
	button := widget.NewButtonWithIcon("Test", theme.HomeIcon(), func() {
		tapped = true
	})
	button.Disable()
	w := test.NewWindow(button)
	defer w.Close()

	test.Tap(button)
	assert.False(t, tapped, "Button should have been disabled")
	test.AssertImageMatches(t, "button/disabled.png", w.Canvas().Capture())

	button.Enable()
	test.AssertImageMatches(t, "button/initial.png", w.Canvas().Capture())
}

func TestButton_Enable(t *testing.T) {
	tapped := make(chan bool)
	button := widget.NewButton("Test", func() {
		tapped <- true
	})

	button.Disable()
	go test.Tap(button)
	func() {
		select {
		case <-tapped:
			assert.Fail(t, "Button should have been disabled")
		case <-time.After(1 * time.Second):
		}
	}()

	button.Enable()
	go test.Tap(button)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Button should have been re-enabled")
		}
	}()
}

func TestButton_Disabled(t *testing.T) {
	button := widget.NewButton("Test", func() {})
	assert.False(t, button.Disabled())
	button.Disable()
	assert.True(t, button.Disabled())
	button.Enable()
	assert.False(t, button.Disabled())
}

func TestButton_Hover(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	b := widget.NewButtonWithIcon("Test", theme.HomeIcon(), func() {})
	w := test.NewWindow(b)
	defer w.Close()

	if !gui.CurrentDevice().IsMobile() {
		test.MoveMouse(w.Canvas(), gui.NewPos(5, 5))
		test.AssertImageMatches(t, "button/hovered.png", w.Canvas().Capture())
	}

	b.Importance = widget.HighImportance
	b.Refresh()
	if !gui.CurrentDevice().IsMobile() {
		test.AssertImageMatches(t, "button/high_importance_hovered.png", w.Canvas().Capture())

		test.MoveMouse(w.Canvas(), gui.NewPos(0, 0))
		b.Refresh()
	}
	test.AssertImageMatches(t, "button/high_importance.png", w.Canvas().Capture())

	b.Importance = widget.MediumImportance
	b.Refresh()
	test.AssertImageMatches(t, "button/initial.png", w.Canvas().Capture())
}

func TestButton_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		text      string
		icon      gui.Resource
		alignment widget.ButtonAlign
		placement widget.ButtonIconPlacement
	}{
		"text_only_center_leading": {
			text:      "Test",
			alignment: widget.ButtonAlignCenter,
			placement: widget.ButtonIconLeadingText,
		},
		"text_only_center_trailing": {
			text:      "Test",
			alignment: widget.ButtonAlignCenter,
			placement: widget.ButtonIconTrailingText,
		},
		"text_only_leading_leading": {
			text:      "Test",
			alignment: widget.ButtonAlignLeading,
			placement: widget.ButtonIconLeadingText,
		},
		"text_only_leading_trailing": {
			text:      "Test",
			alignment: widget.ButtonAlignLeading,
			placement: widget.ButtonIconTrailingText,
		},
		"text_only_trailing_leading": {
			text:      "Test",
			alignment: widget.ButtonAlignTrailing,
			placement: widget.ButtonIconLeadingText,
		},
		"text_only_trailing_trailing": {
			text:      "Test",
			alignment: widget.ButtonAlignTrailing,
			placement: widget.ButtonIconTrailingText,
		},
		"text_only_multiline": {
			text:      "Test\nLine2",
			alignment: widget.ButtonAlignCenter,
		},
		"icon_only_center_leading": {
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignCenter,
			placement: widget.ButtonIconLeadingText,
		},
		"icon_only_center_trailing": {
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignCenter,
			placement: widget.ButtonIconTrailingText,
		},
		"icon_only_leading_leading": {
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignLeading,
			placement: widget.ButtonIconLeadingText,
		},
		"icon_only_leading_trailing": {
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignLeading,
			placement: widget.ButtonIconTrailingText,
		},
		"icon_only_trailing_leading": {
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignTrailing,
			placement: widget.ButtonIconLeadingText,
		},
		"icon_only_trailing_trailing": {
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignTrailing,
			placement: widget.ButtonIconTrailingText,
		},
		"text_icon_center_leading": {
			text:      "Test",
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignCenter,
			placement: widget.ButtonIconLeadingText,
		},
		"text_icon_center_trailing": {
			text:      "Test",
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignCenter,
			placement: widget.ButtonIconTrailingText,
		},
		"text_icon_leading_leading": {
			text:      "Test",
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignLeading,
			placement: widget.ButtonIconLeadingText,
		},
		"text_icon_leading_trailing": {
			text:      "Test",
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignLeading,
			placement: widget.ButtonIconTrailingText,
		},
		"text_icon_trailing_leading": {
			text:      "Test",
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignTrailing,
			placement: widget.ButtonIconLeadingText,
		},
		"text_icon_trailing_trailing": {
			text:      "Test",
			icon:      theme.CancelIcon(),
			alignment: widget.ButtonAlignTrailing,
			placement: widget.ButtonIconTrailingText,
		},
	} {
		t.Run(name, func(t *testing.T) {
			button := &widget.Button{
				Text:          tt.text,
				Icon:          tt.icon,
				Alignment:     tt.alignment,
				IconPlacement: tt.placement,
			}

			window := test.NewWindow(button)
			defer window.Close()
			window.Resize(button.MinSize().Max(gui.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, "button/layout_"+name+".xml", window.Canvas())
		})
	}
}

func TestButton_ChangeTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	b := widget.NewButton("Test", func() {})
	w := test.NewWindow(b)
	defer w.Close()
	w.Resize(gui.NewSize(200, 200))
	b.Resize(b.MinSize())
	b.Move(gui.NewPos(10, 10))

	test.AssertImageMatches(t, "button/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		b.Resize(b.MinSize())
		b.Refresh()
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "button/theme_changed.png", w.Canvas().Capture())
	})
}
