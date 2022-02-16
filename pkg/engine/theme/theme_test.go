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

	gui "github.com/bhojpur/gui/pkg/engine"

	"github.com/stretchr/testify/assert"
)

func TestThemeChange(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	bg := BackgroundColor()

	gui.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, bg, BackgroundColor())
}

func TestTheme_Bootstrapping(t *testing.T) {
	current := gui.CurrentApp().Settings().Theme()
	gui.CurrentApp().Settings().SetTheme(nil)

	// this should not crash
	BackgroundColor()

	gui.CurrentApp().Settings().SetTheme(current)
}

func TestBuiltinTheme_ShadowColor(t *testing.T) {
	shadow := ShadowColor()

	_, _, _, a := shadow.RGBA()
	assert.NotEqual(t, 255, a)
}

func TestTheme_Dark_ReturnsCorrectBackground(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	bg := BackgroundColor()
	assert.Equal(t, DarkTheme().Color(ColorNameBackground, VariantDark), bg, "wrong dark theme background color")
}

func TestTheme_Light_ReturnsCorrectBackground(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(LightTheme())
	bg := BackgroundColor()
	assert.Equal(t, LightTheme().Color(ColorNameBackground, VariantLight), bg, "wrong light theme background color")
}

func Test_ButtonColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ButtonColor()
	assert.Equal(t, DarkTheme().Color(ColorNameButton, VariantDark), c, "wrong button color")
}

func Test_TextColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ForegroundColor()
	assert.Equal(t, DarkTheme().Color(ColorNameForeground, VariantDark), c, "wrong text color")
}

func Test_DisabledTextColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := DisabledColor()
	assert.Equal(t, DarkTheme().Color(ColorNameDisabled, VariantDark), c, "wrong disabled text color")
}

func Test_PlaceHolderColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PlaceHolderColor()
	assert.Equal(t, DarkTheme().Color(ColorNamePlaceHolder, VariantDark), c, "wrong placeholder color")
}

func Test_PrimaryColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PrimaryColor()
	assert.Equal(t, DarkTheme().Color(ColorNamePrimary, VariantDark), c, "wrong primary color")
}

func Test_HoverColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := HoverColor()
	assert.Equal(t, DarkTheme().Color(ColorNameHover, VariantDark), c, "wrong hover color")
}

func Test_FocusColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := FocusColor()
	assert.Equal(t, DarkTheme().Color(ColorNameFocus, VariantDark), c, "wrong focus color")
}

func Test_ScrollBarColor(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ScrollBarColor()
	assert.Equal(t, DarkTheme().Color(ColorNameScrollBar, VariantDark), c, "wrong scrollbar color")
}

func Test_TextSize(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNameText), TextSize(), "wrong text size")
}

func Test_TextFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Regular.ttf"
	result := TextFont().Name()
	assert.Equal(t, expect, result, "wrong regular text font")
}

func Test_TextBoldFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Bold.ttf"
	result := TextBoldFont().Name()
	assert.Equal(t, expect, result, "wrong bold text font")
}

func Test_TextItalicFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Italic.ttf"
	result := TextItalicFont().Name()
	assert.Equal(t, expect, result, "wrong italic text font")
}

func Test_TextBoldItalicFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-BoldItalic.ttf"
	result := TextBoldItalicFont().Name()
	assert.Equal(t, expect, result, "wrong bold italic text font")
}

func Test_TextMonospaceFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "DejaVuSansMono-Powerline.ttf"
	result := TextMonospaceFont().Name()
	assert.Equal(t, expect, result, "wrong monospace font")
}

func Test_Padding(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNamePadding), Padding(), "wrong padding")
}

func Test_IconInlineSize(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNameInlineIcon), IconInlineSize(), "wrong inline icon size")
}

func Test_ScrollBarSize(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNameScrollBar), ScrollBarSize(), "wrong inline icon size")
}

func Test_DefaultTextFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Regular.ttf"
	result := DefaultTextFont().Name()
	assert.Equal(t, expect, result, "wrong default text font")
}

func Test_DefaultTextBoldFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Bold.ttf"
	result := DefaultTextBoldFont().Name()
	assert.Equal(t, expect, result, "wrong default text bold font")
}

func Test_DefaultTextBoldItalicFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-BoldItalic.ttf"
	result := DefaultTextBoldItalicFont().Name()
	assert.Equal(t, expect, result, "wrong default text bold italic font")
}

func Test_DefaultTextItalicFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Italic.ttf"
	result := DefaultTextItalicFont().Name()
	assert.Equal(t, expect, result, "wrong default text italic font")
}

func Test_DefaultTextMonospaceFont(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "DejaVuSansMono-Powerline.ttf"
	result := DefaultTextMonospaceFont().Name()
	assert.Equal(t, expect, result, "wrong default monospace font")
}

func TestEmptyTheme(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(&emptyTheme{})
	assert.NotNil(t, ForegroundColor())
	assert.NotNil(t, TextFont())
	assert.NotNil(t, HelpIcon())
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
}

type emptyTheme struct {
}

func (e *emptyTheme) Color(n gui.ThemeColorName, v gui.ThemeVariant) color.Color {
	return nil
}

func (e *emptyTheme) Font(s gui.TextStyle) gui.Resource {
	return nil
}

func (e *emptyTheme) Icon(n gui.ThemeIconName) gui.Resource {
	return nil
}

func (e *emptyTheme) Size(n gui.ThemeSizeName) float32 {
	return 0
}
