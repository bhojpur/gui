package painter

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
	"image"
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	// DefaultTabWidth is the default width in spaces
	DefaultTabWidth = 4

	// TextDPI is a global constant that determines how text scales to interface sizes
	TextDPI = 78
)

func loadFont(data gui.Resource) *truetype.Font {
	loaded, err := truetype.Parse(data.Content())
	if err != nil {
		gui.LogError("font load error", err)
	}

	return loaded
}

// RenderedTextSize looks up how big a string would be if drawn on screen.
// It also returns the distance from top to the text baseline.
func RenderedTextSize(text string, fontSize float32, style gui.TextStyle) (size gui.Size, baseline float32) {
	size, base := cache.GetFontMetrics(text, fontSize, style)
	if base != 0 {
		return size, base
	}

	size, base = measureText(text, fontSize, style)
	cache.SetFontMetrics(text, fontSize, style, size, base)
	return size, base
}

func measureText(text string, fontSize float32, style gui.TextStyle) (gui.Size, float32) {
	var opts truetype.Options
	opts.Size = float64(fontSize)
	opts.DPI = TextDPI

	face := CachedFontFace(style, &opts)
	advance := MeasureString(face, text, style.TabWidth)

	return gui.NewSize(fixed266ToFloat32(advance), fixed266ToFloat32(face.Metrics().Height)),
		fixed266ToFloat32(face.Metrics().Ascent)
}

type compositeFace struct {
	sync.Mutex

	chosen, fallback         font.Face
	chosenFont, fallbackFont *truetype.Font
}

func (c *compositeFace) containsGlyph(font *truetype.Font, r rune) bool {
	c.Lock()
	defer c.Unlock()

	return font != nil && font.Index(r) != 0
}

func (c *compositeFace) Close() error {
	if c.chosen != nil {
		_ = c.chosen.Close()
	}

	return c.fallback.Close()
}

func (c *compositeFace) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	contains := c.containsGlyph(c.chosenFont, r)

	c.Lock()
	defer c.Unlock()

	if contains {
		return c.chosen.Glyph(dot, r)
	}

	return c.fallback.Glyph(dot, r)
}

func (c *compositeFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	contains := c.containsGlyph(c.chosenFont, r)

	c.Lock()
	defer c.Unlock()

	if contains {
		c.chosen.GlyphBounds(r)
	}
	return c.fallback.GlyphBounds(r)
}

func (c *compositeFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	contains := c.containsGlyph(c.chosenFont, r)

	c.Lock()
	defer c.Unlock()

	if contains {
		return c.chosen.GlyphAdvance(r)
	}
	return c.fallback.GlyphAdvance(r)
}

func (c *compositeFace) Kern(r0, r1 rune) fixed.Int26_6 {
	contains0 := c.containsGlyph(c.chosenFont, r0)
	contains1 := c.containsGlyph(c.chosenFont, r1)

	c.Lock()
	defer c.Unlock()

	if contains0 && contains1 {
		return c.chosen.Kern(r0, r1)
	}
	return c.fallback.Kern(r0, r1)
}

func (c *compositeFace) Metrics() font.Metrics {
	c.Lock()
	defer c.Unlock()

	return c.chosen.Metrics()
}

func newFontWithFallback(chosen, fallback font.Face, chosenFont, fallbackFont *truetype.Font) font.Face {
	return &compositeFace{chosen: chosen, fallback: fallback, chosenFont: chosenFont, fallbackFont: fallbackFont}
}

type fontCacheItem struct {
	font, fallback *truetype.Font
	faces          map[truetype.Options]font.Face
}

var fontCache = make(map[gui.TextStyle]*fontCacheItem)
var fontCacheLock = new(sync.Mutex)

// CachedFontFace returns a font face held in memory. These are loaded from the current theme.
func CachedFontFace(style gui.TextStyle, opts *truetype.Options) font.Face {
	fontCacheLock.Lock()
	defer fontCacheLock.Unlock()
	comp := fontCache[style]

	if comp == nil {
		var f1, f2 *truetype.Font
		switch {
		case style.Monospace:
			f1 = loadFont(theme.TextMonospaceFont())
			f2 = loadFont(theme.DefaultTextMonospaceFont())
		case style.Bold:
			if style.Italic {
				f1 = loadFont(theme.TextBoldItalicFont())
				f2 = loadFont(theme.DefaultTextBoldItalicFont())
			} else {
				f1 = loadFont(theme.TextBoldFont())
				f2 = loadFont(theme.DefaultTextBoldFont())
			}
		case style.Italic:
			f1 = loadFont(theme.TextItalicFont())
			f2 = loadFont(theme.DefaultTextItalicFont())
		default:
			f1 = loadFont(theme.TextFont())
			f2 = loadFont(theme.DefaultTextFont())
		}

		if f1 == nil {
			f1 = f2
		}
		comp = &fontCacheItem{font: f1, fallback: f2, faces: make(map[truetype.Options]font.Face)}
		fontCache[style] = comp
	}

	face := comp.faces[*opts]
	if face == nil {
		f1 := truetype.NewFace(comp.font, opts)
		f2 := truetype.NewFace(comp.fallback, opts)
		face = newFontWithFallback(f1, f2, comp.font, comp.fallback)

		comp.faces[*opts] = face
	}

	return face
}

// ClearFontCache is used to remove cached fonts in the case that we wish to re-load font faces
func ClearFontCache() {
	fontCacheLock.Lock()
	defer fontCacheLock.Unlock()
	for _, item := range fontCache {
		for _, face := range item.faces {
			err := face.Close()

			if err != nil {
				gui.LogError("failed to close font face", err)
				return
			}
		}
	}

	fontCache = make(map[gui.TextStyle]*fontCacheItem)
}

func fixed266ToFloat32(i fixed.Int26_6) float32 {
	return float32(float64(i) / (1 << 6))
}
