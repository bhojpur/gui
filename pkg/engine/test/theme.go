package test

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

var defaultTheme gui.Theme

var _ gui.Theme = (*configurableTheme)(nil)

type configurableTheme struct {
	colors map[gui.ThemeColorName]color.Color
	fonts  map[gui.TextStyle]gui.Resource
	sizes  map[gui.ThemeSizeName]float32
}

// Theme returns a theme useful for image based tests.
func Theme() gui.Theme {
	if defaultTheme == nil {
		defaultTheme = &configurableTheme{
			colors: map[gui.ThemeColorName]color.Color{
				theme.ColorNameBackground:      color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff},
				theme.ColorNameButton:          color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
				theme.ColorNameDisabled:        color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
				theme.ColorNameDisabledButton:  color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.ColorNameError:           color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff},
				theme.ColorNameFocus:           color.NRGBA{R: 0x78, G: 0x3a, B: 0x3a, A: 0xff},
				theme.ColorNameForeground:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
				theme.ColorNameHover:           color.NRGBA{R: 0x88, G: 0xff, B: 0xff, A: 0x22},
				theme.ColorNameInputBackground: color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xff},
				theme.ColorNamePlaceHolder:     color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
				theme.ColorNamePressed:         color.NRGBA{A: 0x33},
				theme.ColorNamePrimary:         color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
				theme.ColorNameScrollBar:       color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
				theme.ColorNameSelection:       color.NRGBA{R: 0x78, G: 0x3a, B: 0x3a, A: 0x99},
				theme.ColorNameShadow:          color.NRGBA{A: 0x88},
			},
			fonts: map[gui.TextStyle]gui.Resource{
				{}:                         theme.DefaultTextFont(),
				{Bold: true}:               theme.DefaultTextBoldFont(),
				{Bold: true, Italic: true}: theme.DefaultTextBoldItalicFont(),
				{Italic: true}:             theme.DefaultTextItalicFont(),
				{Monospace: true}:          theme.DefaultTextMonospaceFont(),
			},
			sizes: map[gui.ThemeSizeName]float32{
				theme.SizeNameInlineIcon:         float32(20),
				theme.SizeNamePadding:            float32(4),
				theme.SizeNameScrollBar:          float32(16),
				theme.SizeNameScrollBarSmall:     float32(3),
				theme.SizeNameSeparatorThickness: float32(1),
				theme.SizeNameText:               float32(14),
				theme.SizeNameHeadingText:        float32(23.8),
				theme.SizeNameSubHeadingText:     float32(18),
				theme.SizeNameCaptionText:        float32(11),
				theme.SizeNameInputBorder:        float32(2),
			},
		}
	}
	return defaultTheme
}

func (t *configurableTheme) Color(n gui.ThemeColorName, _ gui.ThemeVariant) color.Color {
	return t.colors[n]
}

func (t *configurableTheme) Font(style gui.TextStyle) gui.Resource {
	return t.fonts[style]
}

func (t *configurableTheme) Icon(n gui.ThemeIconName) gui.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *configurableTheme) Size(s gui.ThemeSizeName) float32 {
	return t.sizes[s]
}
