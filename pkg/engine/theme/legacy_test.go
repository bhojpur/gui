package theme

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
	"testing"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
)

var oldTheme = &legacyTheme{}

func TestFromLegacy(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.NotNil(t, newTheme)
	assert.Equal(t, oldTheme, newTheme.(*legacyWrapper).old)
}

func TestLegacyWrapper_Color(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.Equal(t, oldTheme.BackgroundColor(), newTheme.Color(ColorNameBackground, VariantLight))
	assert.Equal(t, oldTheme.ShadowColor(), newTheme.Color(ColorNameShadow, VariantLight))
	assert.Equal(t, oldTheme.TextColor(), newTheme.Color(ColorNameForeground, VariantLight))
}

func TestLegacyWrapper_Font(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.Equal(t, oldTheme.TextFont(), newTheme.Font(gui.TextStyle{}))
	assert.Equal(t, oldTheme.TextBoldFont(), newTheme.Font(gui.TextStyle{Bold: true}))
	assert.Equal(t, oldTheme.TextItalicFont(), newTheme.Font(gui.TextStyle{Italic: true}))
	assert.Equal(t, oldTheme.TextMonospaceFont(), newTheme.Font(gui.TextStyle{Monospace: true}))
}

func TestLegacyWrapper_Size(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.Equal(t, oldTheme.IconInlineSize(), int(newTheme.Size(SizeNameInlineIcon)))
	assert.Equal(t, oldTheme.Padding(), int(newTheme.Size(SizeNamePadding)))
	assert.Equal(t, oldTheme.TextSize(), int(newTheme.Size(SizeNameText)))
}

var _ gui.LegacyTheme = (*legacyTheme)(nil)

type legacyTheme struct {
}

func (t *legacyTheme) BackgroundColor() color.Color {
	return BackgroundColor()
}

func (t *legacyTheme) ButtonColor() color.Color {
	return ButtonColor()
}

func (t *legacyTheme) DisabledButtonColor() color.Color {
	return DisabledButtonColor()
}

func (t *legacyTheme) DisabledTextColor() color.Color {
	return DisabledColor()
}

func (t *legacyTheme) FocusColor() color.Color {
	return FocusColor()
}

func (t *legacyTheme) HoverColor() color.Color {
	return HoverColor()
}

func (t *legacyTheme) PlaceHolderColor() color.Color {
	return PlaceHolderColor()
}

func (t *legacyTheme) PrimaryColor() color.Color {
	return PrimaryColor()
}

func (t *legacyTheme) ScrollBarColor() color.Color {
	return ScrollBarColor()
}

func (t *legacyTheme) ShadowColor() color.Color {
	return ShadowColor()
}

func (t *legacyTheme) TextColor() color.Color {
	return ForegroundColor()
}

func (t *legacyTheme) TextSize() int {
	return int(TextSize())
}

func (t *legacyTheme) TextFont() gui.Resource {
	return TextFont()
}

func (t *legacyTheme) TextBoldFont() gui.Resource {
	return TextBoldFont()
}

func (t *legacyTheme) TextItalicFont() gui.Resource {
	return TextItalicFont()
}

func (t *legacyTheme) TextBoldItalicFont() gui.Resource {
	return TextBoldItalicFont()
}

func (t *legacyTheme) TextMonospaceFont() gui.Resource {
	return TextMonospaceFont()
}

func (t *legacyTheme) Padding() int {
	return int(Padding())
}

func (t *legacyTheme) IconInlineSize() int {
	return int(IconInlineSize())
}

func (t *legacyTheme) ScrollBarSize() int {
	return int(ScrollBarSize())
}

func (t *legacyTheme) ScrollBarSmallSize() int {
	return int(ScrollBarSmallSize())
}
