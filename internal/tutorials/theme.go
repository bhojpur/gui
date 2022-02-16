package tutorials

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

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var (
	purple = &color.NRGBA{R: 128, G: 0, B: 128, A: 255}
	orange = &color.NRGBA{R: 198, G: 123, B: 0, A: 255}
	grey   = &color.Gray{Y: 123}
)

// customTheme is a simple demonstration of a bespoke theme loaded by a Bhojpur GUI app.
type customTheme struct {
}

func (customTheme) Color(c gui.ThemeColorName, _ gui.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return purple
	case theme.ColorNameButton, theme.ColorNameDisabled:
		return color.Black
	case theme.ColorNamePlaceHolder, theme.ColorNameScrollBar:
		return grey
	case theme.ColorNamePrimary, theme.ColorNameHover, theme.ColorNameFocus:
		return orange
	case theme.ColorNameShadow:
		return &color.RGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 0xcc}
	default:
		return color.White
	}
}

func (customTheme) Font(style gui.TextStyle) gui.Resource {
	return theme.DarkTheme().Font(style)
}

func (customTheme) Icon(n gui.ThemeIconName) gui.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (customTheme) Size(s gui.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 5
	case theme.SizeNameText:
		return 18
	case theme.SizeNameHeadingText:
		return 30
	case theme.SizeNameSubHeadingText:
		return 25
	case theme.SizeNameCaptionText:
		return 15
	case theme.SizeNameInputBorder:
		return 1
	default:
		return 0
	}
}

func newCustomTheme() gui.Theme {
	return &customTheme{}
}
