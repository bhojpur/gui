// It defines how a Bhojpur GUI application should look when rendered
package theme // import "github.com/bhojpur/gui/pkg/engine/theme"

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
	"os"
	"strings"

	gui "github.com/bhojpur/gui/pkg/engine"
)

const (
	// ColorRed is the red primary color name.
	//
	// Since: 1.4
	ColorRed = "red"
	// ColorOrange is the orange primary color name.
	//
	// Since: 1.4
	ColorOrange = "orange"
	// ColorYellow is the yellow primary color name.
	//
	// Since: 1.4
	ColorYellow = "yellow"
	// ColorGreen is the green primary color name.
	//
	// Since: 1.4
	ColorGreen = "green"
	// ColorBlue is the blue primary color name.
	//
	// Since: 1.4
	ColorBlue = "blue"
	// ColorPurple is the purple primary color name.
	//
	// Since: 1.4
	ColorPurple = "purple"
	// ColorBrown is the brown primary color name.
	//
	// Since: 1.4
	ColorBrown = "brown"
	// ColorGray is the gray primary color name.
	//
	// Since: 1.4
	ColorGray = "gray"

	// ColorNameBackground is the name of theme lookup for background color.
	//
	// Since: 2.0
	ColorNameBackground gui.ThemeColorName = "background"

	// ColorNameButton is the name of theme lookup for button color.
	//
	// Since: 2.0
	ColorNameButton gui.ThemeColorName = "button"

	// ColorNameDisabledButton is the name of theme lookup for disabled button color.
	//
	// Since: 2.0
	ColorNameDisabledButton gui.ThemeColorName = "disabledButton"

	// ColorNameDisabled is the name of theme lookup for disabled foreground color.
	//
	// Since: 2.0
	ColorNameDisabled gui.ThemeColorName = "disabled"

	// ColorNameError is the name of theme lookup for foreground error color.
	//
	// Since: 2.0
	ColorNameError gui.ThemeColorName = "error"

	// ColorNameFocus is the name of theme lookup for focus color.
	//
	// Since: 2.0
	ColorNameFocus gui.ThemeColorName = "focus"

	// ColorNameForeground is the name of theme lookup for foreground color.
	//
	// Since: 2.0
	ColorNameForeground gui.ThemeColorName = "foreground"

	// ColorNameHover is the name of theme lookup for hover color.
	//
	// Since: 2.0
	ColorNameHover gui.ThemeColorName = "hover"

	// ColorNameInputBackground is the name of theme lookup for background color of an input field.
	//
	// Since: 2.0
	ColorNameInputBackground gui.ThemeColorName = "inputBackground"

	// ColorNamePlaceHolder is the name of theme lookup for placeholder text color.
	//
	// Since: 2.0
	ColorNamePlaceHolder gui.ThemeColorName = "placeholder"

	// ColorNamePressed is the name of theme lookup for the tap overlay color.
	//
	// Since: 2.0
	ColorNamePressed gui.ThemeColorName = "pressed"

	// ColorNamePrimary is the name of theme lookup for primary color.
	//
	// Since: 2.0
	ColorNamePrimary gui.ThemeColorName = "primary"

	// ColorNameScrollBar is the name of theme lookup for scrollbar color.
	//
	// Since: 2.0
	ColorNameScrollBar gui.ThemeColorName = "scrollBar"

	// ColorNameSelection is the name of theme lookup for selection color.
	//
	// Since: 2.1
	ColorNameSelection gui.ThemeColorName = "selection"

	// ColorNameShadow is the name of theme lookup for shadow color.
	//
	// Since: 2.0
	ColorNameShadow gui.ThemeColorName = "shadow"

	// SizeNameCaptionText is the name of theme lookup for helper text size, normally smaller than regular text size.
	//
	// Since: 2.0
	SizeNameCaptionText gui.ThemeSizeName = "helperText"

	// SizeNameInlineIcon is the name of theme lookup for inline icons size.
	//
	// Since: 2.0
	SizeNameInlineIcon gui.ThemeSizeName = "iconInline"

	// SizeNamePadding is the name of theme lookup for padding size.
	//
	// Since: 2.0
	SizeNamePadding gui.ThemeSizeName = "padding"

	// SizeNameScrollBar is the name of theme lookup for the scrollbar size.
	//
	// Since: 2.0
	SizeNameScrollBar gui.ThemeSizeName = "scrollBar"

	// SizeNameScrollBarSmall is the name of theme lookup for the shrunk scrollbar size.
	//
	// Since: 2.0
	SizeNameScrollBarSmall gui.ThemeSizeName = "scrollBarSmall"

	// SizeNameSeparatorThickness is the name of theme lookup for the thickness of a separator.
	//
	// Since: 2.0
	SizeNameSeparatorThickness gui.ThemeSizeName = "separator"

	// SizeNameText is the name of theme lookup for text size.
	//
	// Since: 2.0
	SizeNameText gui.ThemeSizeName = "text"

	// SizeNameHeadingText is the name of theme lookup for text size of a heading.
	//
	// Since: 2.1
	SizeNameHeadingText gui.ThemeSizeName = "headingText"

	// SizeNameSubHeadingText is the name of theme lookup for text size of a sub-heading.
	//
	// Since: 2.1
	SizeNameSubHeadingText gui.ThemeSizeName = "subHeadingText"

	// SizeNameInputBorder is the name of theme lookup for input border size.
	//
	// Since: 2.0
	SizeNameInputBorder gui.ThemeSizeName = "inputBorder"

	// VariantDark is the version of a theme that satisfies a user preference for a light look.
	//
	// Since: 2.0
	VariantDark gui.ThemeVariant = 0

	// VariantLight is the version of a theme that satisfies a user preference for a dark look.
	//
	// Since: 2.0
	VariantLight gui.ThemeVariant = 1

	// potential for adding theme types such as high visibility or monochrome...
	variantNameUserPreference gui.ThemeVariant = 2 // locally used in builtinTheme for backward compatibility
)

// BackgroundColor returns the theme's background color.
func BackgroundColor() color.Color {
	return safeColorLookup(ColorNameBackground, currentVariant())
}

// ButtonColor returns the theme's standard button color.
func ButtonColor() color.Color {
	return safeColorLookup(ColorNameButton, currentVariant())
}

// CaptionTextSize returns the size for caption text.
func CaptionTextSize() float32 {
	return current().Size(SizeNameCaptionText)
}

// DarkTheme defines the built-in dark theme colors and sizes.
//
// Deprecated: This method ignores user preference and should not be used, it will be removed in v3.0.
func DarkTheme() gui.Theme {
	theme := &builtinTheme{variant: VariantDark}

	theme.initFonts()
	return theme
}

// DefaultTextBoldFont returns the font resource for the built-in bold font style.
func DefaultTextBoldFont() gui.Resource {
	return bold
}

// DefaultTextBoldItalicFont returns the font resource for the built-in bold and italic font style.
func DefaultTextBoldItalicFont() gui.Resource {
	return bolditalic
}

// DefaultTextFont returns the font resource for the built-in regular font style.
func DefaultTextFont() gui.Resource {
	return regular
}

// DefaultTextItalicFont returns the font resource for the built-in italic font style.
func DefaultTextItalicFont() gui.Resource {
	return italic
}

// DefaultTextMonospaceFont returns the font resource for the built-in monospace font face.
func DefaultTextMonospaceFont() gui.Resource {
	return monospace
}

// DefaultSymbolFont returns the font resource for the built-in symbol font.
//
// Since: 2.2
func DefaultSymbolFont() gui.Resource {
	return symbol
}

// DefaultTheme returns a built-in theme that can adapt to the user preference of light or dark colors.
//
// Since: 2.0
func DefaultTheme() gui.Theme {
	return defaultTheme
}

// DisabledButtonColor returns the theme's disabled button color.
func DisabledButtonColor() color.Color {
	return safeColorLookup(ColorNameDisabledButton, currentVariant())
}

// DisabledColor returns the foreground color for a disabled UI element.
//
// Since: 2.0
func DisabledColor() color.Color {
	return safeColorLookup(ColorNameDisabled, currentVariant())
}

// DisabledTextColor returns the theme's disabled text color - this is actually the disabled color since 1.4.
//
// Deprecated: Use theme.DisabledColor() colour instead.
func DisabledTextColor() color.Color {
	return DisabledColor()
}

// ErrorColor returns the theme's error text color.
//
// Since: 2.0
func ErrorColor() color.Color {
	return safeColorLookup(ColorNameError, currentVariant())
}

// FocusColor returns the color used to highlight a focused widget.
func FocusColor() color.Color {
	return safeColorLookup(ColorNameFocus, currentVariant())
}

// ForegroundColor returns the theme's standard foreground color for text and icons.
//
// Since: 2.0
func ForegroundColor() color.Color {
	return safeColorLookup(ColorNameForeground, currentVariant())
}

// HoverColor returns the color used to highlight interactive elements currently under a cursor.
func HoverColor() color.Color {
	return safeColorLookup(ColorNameHover, currentVariant())
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() float32 {
	return current().Size(SizeNameInlineIcon)
}

// InputBackgroundColor returns the color used to draw underneath input elements.
func InputBackgroundColor() color.Color {
	return current().Color(ColorNameInputBackground, currentVariant())
}

// InputBorderSize returns the input border size (or underline size for an entry).
//
// Since: 2.0
func InputBorderSize() float32 {
	return current().Size(SizeNameInputBorder)
}

// LightTheme defines the built-in light theme colors and sizes.
//
// Deprecated: This method ignores user preference and should not be used, it will be removed in v3.0.
func LightTheme() gui.Theme {
	theme := &builtinTheme{variant: VariantLight}

	theme.initFonts()
	return theme
}

// Padding is the standard gap between elements and the border around interface elements.
func Padding() float32 {
	return current().Size(SizeNamePadding)
}

// PlaceHolderColor returns the theme's standard text color.
func PlaceHolderColor() color.Color {
	return safeColorLookup(ColorNamePlaceHolder, currentVariant())
}

// PressedColor returns the color used to overlap tapped features.
//
// Since: 2.0
func PressedColor() color.Color {
	return safeColorLookup(ColorNamePressed, currentVariant())
}

// PrimaryColor returns the color used to highlight primary features.
func PrimaryColor() color.Color {
	return safeColorLookup(ColorNamePrimary, currentVariant())
}

// PrimaryColorNamed returns a theme specific color value for a named primary color.
//
// Since: 1.4
func PrimaryColorNamed(name string) color.Color {
	col, ok := primaryColors[name]
	if !ok {
		return primaryColors[ColorBlue]
	}
	return col
}

// PrimaryColorNames returns a list of the standard primary color options.
//
// Since: 1.4
func PrimaryColorNames() []string {
	return []string{ColorRed, ColorOrange, ColorYellow, ColorGreen, ColorBlue, ColorPurple, ColorBrown, ColorGray}
}

// ScrollBarColor returns the color (and translucency) for a scrollBar.
func ScrollBarColor() color.Color {
	return safeColorLookup(ColorNameScrollBar, currentVariant())
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer.
func ScrollBarSize() float32 {
	return current().Size(SizeNameScrollBar)
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer.
func ScrollBarSmallSize() float32 {
	return current().Size(SizeNameScrollBarSmall)
}

// SelectionColor returns the color for a selected element.
//
// Since: 2.1
func SelectionColor() color.Color {
	return safeColorLookup(ColorNameSelection, currentVariant())
}

// SeparatorThicknessSize is the standard thickness of the separator widget.
//
// Since: 2.0
func SeparatorThicknessSize() float32 {
	return current().Size(SizeNameSeparatorThickness)
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation.
func ShadowColor() color.Color {
	return safeColorLookup(ColorNameShadow, currentVariant())
}

// TextBoldFont returns the font resource for the bold font style.
func TextBoldFont() gui.Resource {
	return safeFontLookup(gui.TextStyle{Bold: true})
}

// TextBoldItalicFont returns the font resource for the bold and italic font style.
func TextBoldItalicFont() gui.Resource {
	return safeFontLookup(gui.TextStyle{Bold: true, Italic: true})
}

// TextColor returns the theme's standard text color - this is actually the foreground color since 1.4.
//
// Deprecated: Use theme.ForegroundColor() colour instead.
func TextColor() color.Color {
	return safeColorLookup(ColorNameForeground, currentVariant())
}

// TextFont returns the font resource for the regular font style.
func TextFont() gui.Resource {
	return safeFontLookup(gui.TextStyle{})
}

// TextHeadingSize returns the text size for header text.
//
// Since: 2.1
func TextHeadingSize() float32 {
	return current().Size(SizeNameHeadingText)
}

// TextItalicFont returns the font resource for the italic font style.
func TextItalicFont() gui.Resource {
	return safeFontLookup(gui.TextStyle{Italic: true})
}

// TextMonospaceFont returns the font resource for the monospace font face.
func TextMonospaceFont() gui.Resource {
	return safeFontLookup(gui.TextStyle{Monospace: true})
}

// TextSize returns the standard text size.
func TextSize() float32 {
	return current().Size(SizeNameText)
}

// TextSubHeadingSize returns the text size for sub-header text.
//
// Since: 2.1
func TextSubHeadingSize() float32 {
	return current().Size(SizeNameSubHeadingText)
}

var (
	defaultTheme = setupDefaultTheme()

	errorColor  = color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
	focusColors = map[string]color.Color{
		ColorRed:    color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0x7f},
		ColorOrange: color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0x7f},
		ColorYellow: color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0x7f},
		ColorGreen:  color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0x7f},
		ColorBlue:   color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0x7f},
		ColorPurple: color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0x7f},
		ColorBrown:  color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0x7f},
		ColorGray:   color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0x7f},
	}
	primaryColors = map[string]color.Color{
		ColorRed:    color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff},
		ColorOrange: color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff},
		ColorYellow: color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0xff},
		ColorGreen:  color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0xff},
		ColorBlue:   color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0xff},
		ColorPurple: color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0xff},
		ColorBrown:  color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0xff},
		ColorGray:   color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0xff},
	}
	selectionColors = map[string]color.Color{
		ColorRed:    color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0x3f},
		ColorOrange: color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0x3f},
		ColorYellow: color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0x3f},
		ColorGreen:  color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0x3f},
		ColorBlue:   color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0x3f},
		ColorPurple: color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0x3f},
		ColorBrown:  color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0x3f},
		ColorGray:   color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0x3f},
	}

	darkPalette = map[gui.ThemeColorName]color.Color{
		ColorNameBackground:      color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff},
		ColorNameButton:          color.Transparent,
		ColorNameDisabled:        color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x42},
		ColorNameDisabledButton:  color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xff},
		ColorNameError:           errorColor,
		ColorNameForeground:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		ColorNameHover:           color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0f},
		ColorNameInputBackground: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x19},
		ColorNamePlaceHolder:     color.NRGBA{R: 0xb2, G: 0xb2, B: 0xb2, A: 0xff},
		ColorNamePressed:         color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x66},
		ColorNameScrollBar:       color.NRGBA{A: 0x99},
		ColorNameShadow:          color.NRGBA{A: 0x66},
	}

	lightPalette = map[gui.ThemeColorName]color.Color{
		ColorNameBackground:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		ColorNameButton:          color.Transparent,
		ColorNameDisabled:        color.NRGBA{A: 0x42},
		ColorNameDisabledButton:  color.NRGBA{R: 0xe5, G: 0xe5, B: 0xe5, A: 0xff},
		ColorNameError:           errorColor,
		ColorNameForeground:      color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xff},
		ColorNameHover:           color.NRGBA{A: 0x0f},
		ColorNameInputBackground: color.NRGBA{A: 0x19},
		ColorNamePlaceHolder:     color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
		ColorNamePressed:         color.NRGBA{A: 0x19},
		ColorNameScrollBar:       color.NRGBA{A: 0x99},
		ColorNameShadow:          color.NRGBA{A: 0x33},
	}
)

type builtinTheme struct {
	variant gui.ThemeVariant

	regular, bold, italic, boldItalic, monospace gui.Resource
}

func (t *builtinTheme) initFonts() {
	t.regular = regular
	t.bold = bold
	t.italic = italic
	t.boldItalic = bolditalic
	t.monospace = monospace

	font := os.Getenv("BHOJPUR_FONT")
	if font != "" {
		t.regular = loadCustomFont(font, "Regular", regular)
		if t.regular == regular { // failed to load
			t.bold = loadCustomFont(font, "Bold", bold)
			t.italic = loadCustomFont(font, "Italic", italic)
			t.boldItalic = loadCustomFont(font, "BoldItalic", bolditalic)
		} else { // first custom font loaded, fall back to that
			t.bold = loadCustomFont(font, "Bold", t.regular)
			t.italic = loadCustomFont(font, "Italic", t.regular)
			t.boldItalic = loadCustomFont(font, "BoldItalic", t.regular)
		}
	}
	font = os.Getenv("BHOJPUR_FONT_MONOSPACE")
	if font != "" {
		t.monospace = loadCustomFont(font, "Regular", monospace)
	}
}

func (t *builtinTheme) Color(n gui.ThemeColorName, v gui.ThemeVariant) color.Color {
	if t.variant != variantNameUserPreference {
		v = t.variant
	}
	colors := darkPalette
	if v == VariantLight {
		colors = lightPalette
	}

	if n == ColorNamePrimary {
		return PrimaryColorNamed(gui.CurrentApp().Settings().PrimaryColor())
	} else if n == ColorNameFocus {
		return focusColorNamed(gui.CurrentApp().Settings().PrimaryColor())
	} else if n == ColorNameSelection {
		return selectionColorNamed(gui.CurrentApp().Settings().PrimaryColor())
	}

	if c, ok := colors[n]; ok {
		return c
	}
	return color.Transparent
}

func (t *builtinTheme) Font(style gui.TextStyle) gui.Resource {
	if style.Monospace {
		return t.monospace
	}
	if style.Bold {
		if style.Italic {
			return t.boldItalic
		}
		return t.bold
	}
	if style.Italic {
		return t.italic
	}
	return t.regular
}

func (t *builtinTheme) Size(s gui.ThemeSizeName) float32 {
	switch s {
	case SizeNameSeparatorThickness:
		return 1
	case SizeNameInlineIcon:
		return 20
	case SizeNamePadding:
		return 4
	case SizeNameScrollBar:
		return 16
	case SizeNameScrollBarSmall:
		return 3
	case SizeNameText:
		return 14
	case SizeNameHeadingText:
		return 24
	case SizeNameSubHeadingText:
		return 18
	case SizeNameCaptionText:
		return 11
	case SizeNameInputBorder:
		return 2
	default:
		return 0
	}
}

func current() gui.Theme {
	if gui.CurrentApp() == nil || gui.CurrentApp().Settings().Theme() == nil {
		return DarkTheme()
	}

	return gui.CurrentApp().Settings().Theme()
}

func currentVariant() gui.ThemeVariant {
	if std, ok := current().(*builtinTheme); ok {
		if std.variant != variantNameUserPreference {
			return std.variant // override if using the old LightTheme() or DarkTheme() constructor
		}
	}

	return gui.CurrentApp().Settings().ThemeVariant()
}

func focusColorNamed(name string) color.Color {
	col, ok := focusColors[name]
	if !ok {
		return focusColors[ColorBlue]
	}
	return col
}

func loadCustomFont(env, variant string, fallback gui.Resource) gui.Resource {
	variantPath := strings.Replace(env, "Regular", variant, -1)

	res, err := gui.LoadResourceFromPath(variantPath)
	if err != nil {
		gui.LogError("Error loading specified font", err)
		return fallback
	}

	return res
}

func safeColorLookup(n gui.ThemeColorName, v gui.ThemeVariant) color.Color {
	col := current().Color(n, v)
	if col == nil {
		gui.LogError("Loaded theme returned nil color", nil)
		return fallbackColor
	}
	return col
}

func safeFontLookup(s gui.TextStyle) gui.Resource {
	font := current().Font(s)
	if font != nil {
		return font
	}
	gui.LogError("Loaded theme returned nil font", nil)

	if s.Monospace {
		return DefaultTextMonospaceFont()
	}
	if s.Bold {
		if s.Italic {
			return DefaultTextBoldItalicFont()
		}
		return DefaultTextBoldFont()
	}
	if s.Italic {
		return DefaultTextItalicFont()
	}

	return DefaultTextFont()
}

func selectionColorNamed(name string) color.Color {
	col, ok := selectionColors[name]
	if !ok {
		return selectionColors[ColorBlue]
	}
	return col
}

func setupDefaultTheme() gui.Theme {
	theme := &builtinTheme{variant: variantNameUserPreference}

	theme.initFonts()
	return theme
}
