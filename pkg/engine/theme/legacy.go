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

	gui "github.com/bhojpur/gui/pkg/engine"
)

// FromLegacy returns a 2.0 Theme created from the given LegacyTheme data.
// This is a transition path and will be removed in the future (probably version 3.0).
//
// Since: 2.0
func FromLegacy(t gui.LegacyTheme) gui.Theme {
	return &legacyWrapper{old: t}
}

var _ gui.Theme = (*legacyWrapper)(nil)

type legacyWrapper struct {
	old gui.LegacyTheme
}

func (l *legacyWrapper) Color(n gui.ThemeColorName, v gui.ThemeVariant) color.Color {
	switch n {
	case ColorNameBackground:
		return l.old.BackgroundColor()
	case ColorNameForeground:
		return l.old.TextColor()
	case ColorNameButton:
		return l.old.ButtonColor()
	case ColorNameDisabledButton:
		return l.old.DisabledButtonColor()
	case ColorNameDisabled:
		return l.old.DisabledTextColor()
	case ColorNameFocus:
		return l.old.FocusColor()
	case ColorNameHover:
		return l.old.HoverColor()
	case ColorNamePlaceHolder:
		return l.old.PlaceHolderColor()
	case ColorNamePrimary:
		return l.old.PrimaryColor()
	case ColorNameScrollBar:
		return l.old.ScrollBarColor()
	case ColorNameShadow:
		return l.old.ShadowColor()
	default:
		return DefaultTheme().Color(n, v)
	}
}

func (l *legacyWrapper) Font(s gui.TextStyle) gui.Resource {
	if s.Monospace {
		return l.old.TextMonospaceFont()
	}
	if s.Bold {
		if s.Italic {
			return l.old.TextBoldItalicFont()
		}
		return l.old.TextBoldFont()
	}
	if s.Italic {
		return l.old.TextItalicFont()
	}
	return l.old.TextFont()
}

func (l *legacyWrapper) Icon(n gui.ThemeIconName) gui.Resource {
	return DefaultTheme().Icon(n)
}

func (l *legacyWrapper) Size(n gui.ThemeSizeName) float32 {
	switch n {
	case SizeNameInlineIcon:
		return float32(l.old.IconInlineSize())
	case SizeNamePadding:
		return float32(l.old.Padding())
	case SizeNameScrollBar:
		return float32(l.old.ScrollBarSize())
	case SizeNameScrollBarSmall:
		return float32(l.old.ScrollBarSmallSize())
	case SizeNameText:
		return float32(l.old.TextSize())
	default:
		return DefaultTheme().Size(n)
	}
}
