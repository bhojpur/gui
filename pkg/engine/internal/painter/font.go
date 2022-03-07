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
	"image/color"
	"image/draw"
	"log"
	"math"
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/goki/freetype"
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

// CachedFontFace returns a font face held in memory. These are loaded from the current theme.
func CachedFontFace(style gui.TextStyle, opts *truetype.Options) font.Face {
	val, ok := fontCache.Load(style)
	if !ok {
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
		case style.Symbol:
			f2 = loadFont(theme.DefaultSymbolFont())
		default:
			f1 = loadFont(theme.TextFont())
			f2 = loadFont(theme.DefaultTextFont())
		}

		if f1 == nil {
			f1 = f2
		}
		val = &fontCacheItem{font: f1, fallback: f2, faces: make(map[truetype.Options]font.Face)}
		fontCache.Store(style, val)
	}

	comp := val.(*fontCacheItem)
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
	fontCache.Range(func(_, val interface{}) bool {
		item := val.(*fontCacheItem)
		for _, face := range item.faces {
			err := face.Close()

			if err != nil {
				gui.LogError("failed to close font face", err)
				return false
			}
		}
		return true
	})

	fontCache = &sync.Map{}
}

// DrawString draws a string into an image.
func DrawString(dst draw.Image, s string, color color.Color, face font.Face, height int, tabWidth int) {
	src := &image.Uniform{C: color}
	dot := freetype.Pt(0, height-face.Metrics().Descent.Ceil())
	walkString(face, s, tabWidth, &dot.X, func(r rune) (fixed.Int26_6, bool) {
		dr, mask, maskp, advance, ok := face.Glyph(dot, r)
		if ok {
			draw.DrawMask(dst, dr, src, image.Point{}, mask, maskp, draw.Over)
		}
		return advance, ok
	})
}

// MeasureString returns how far dot would advance by drawing s with f.
// Tabs are translated into a dot location change.
func MeasureString(f font.Face, s string, tabWidth int) (advance fixed.Int26_6) {
	walkString(f, s, tabWidth, &advance, f.GlyphAdvance)
	return
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

func fixed266ToFloat32(i fixed.Int26_6) float32 {
	return float32(float64(i) / (1 << 6))
}

func loadFont(data gui.Resource) *truetype.Font {
	loaded, err := truetype.Parse(data.Content())
	if err != nil {
		gui.LogError("font load error", err)
	}

	return loaded
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

func newFontWithFallback(chosen, fallback font.Face, chosenFont, fallbackFont ttfFont) font.Face {
	return &compositeFace{chosen: chosen, fallback: fallback, chosenFont: chosenFont, fallbackFont: fallbackFont}
}

func tabStop(f font.Face, x fixed.Int26_6, tabWidth int) fixed.Int26_6 {
	if tabWidth <= 0 {
		tabWidth = DefaultTabWidth
	}
	spacew, ok := f.GlyphAdvance(' ')
	if !ok {
		log.Print("Failed to find space width for tab")
		return x
	}
	tabw := spacew * fixed.Int26_6(tabWidth)
	tabs, _ := math.Modf(float64((x + tabw) / tabw))
	return tabw * fixed.Int26_6(tabs)
}

func walkString(f font.Face, s string, tabWidth int, advance *fixed.Int26_6, cb func(r rune) (fixed.Int26_6, bool)) {
	prevC := rune(-1)
	for _, c := range s {
		if c == '\r' {
			prevC = rune(-1)
			continue
		}
		if prevC >= 0 {
			*advance += f.Kern(prevC, c)
		}
		if c == '\t' {
			*advance = tabStop(f, *advance, tabWidth)
		} else {
			a, ok := cb(c)
			if !ok {
				c = '�' // U+FFFD, the Unicode replacement character
				a, ok = cb(c)
			}
			if !ok {
				// TODO: set prevC = '\ufffd'?
				continue
			}
			*advance += a
		}

		prevC = c
	}
}

type compositeFace struct {
	sync.Mutex

	chosen, fallback         font.Face
	chosenFont, fallbackFont ttfFont
}

func (c *compositeFace) Close() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.chosen != nil {
		err = c.chosen.Close()
	}

	err2 := c.fallback.Close()
	if err2 != nil {
		return err2
	}

	return
}

func (c *compositeFace) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r) {
		return c.chosen.Glyph(dot, r)
	}

	if c.containsGlyph(c.fallbackFont, r) {
		return c.fallback.Glyph(dot, r)
	}

	return
}

func (c *compositeFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r) {
		return c.chosen.GlyphAdvance(r)
	}

	if c.containsGlyph(c.fallbackFont, r) {
		return c.fallback.GlyphAdvance(r)
	}

	return
}

func (c *compositeFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r) {
		return c.chosen.GlyphBounds(r)
	}

	if c.containsGlyph(c.fallbackFont, r) {
		return c.fallback.GlyphBounds(r)
	}

	return
}

func (c *compositeFace) Kern(r0, r1 rune) fixed.Int26_6 {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r0) && c.containsGlyph(c.chosenFont, r1) {
		return c.chosen.Kern(r0, r1)
	}

	return c.fallback.Kern(r0, r1)
}

func (c *compositeFace) Metrics() font.Metrics {
	c.Lock()
	defer c.Unlock()

	return c.chosen.Metrics()
}

func (c *compositeFace) containsGlyph(font ttfFont, r rune) bool {
	return font != nil && font.Index(r) != 0
}

type ttfFont interface {
	Index(rune) truetype.Index
}

type fontCacheItem struct {
	font, fallback *truetype.Font
	faces          map[truetype.Options]font.Face
}

var fontCache = &sync.Map{} // map[gui.TextStyle]*fontCacheItem
