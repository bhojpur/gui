package spatial

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
	"math"
	"net/http"
	"net/url"

	"github.com/nfnt/resize"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"golang.org/x/image/draw"
)

const tileSize = 256

// Map widget renders an interactive map using OpenStreetMap tile data.
type Map struct {
	widget.BaseWidget

	pixels     *image.NRGBA
	w, h       int
	zoom, x, y int

	cl *http.Client
}

// NewMap creates a new instance of the map widget.
func NewMap() *Map {
	m := &Map{cl: &http.Client{}}
	m.ExtendBaseWidget(m)
	return m
}

// MinSize returns the smallest possible size for a widget.
// For our map this is a constant size representing a single tile on a device with
// the highest known DPI (4x).
func (m *Map) MinSize() gui.Size {
	return gui.NewSize(64, 64)
}

// CreateRenderer returns the renderer for this widget.
// A map renderer is simply the map Raster with user interface elements overlaid.
func (m *Map) CreateRenderer() gui.WidgetRenderer {
	license, _ := url.Parse("https://openstreetmap.org")
	copyright := widget.NewHyperlink("OpenStreetMap", license)
	copyright.Alignment = gui.TextAlignTrailing
	zoom := container.NewVBox(
		newMapButton(theme.ZoomInIcon(), func() {
			if m.zoom >= 19 {
				return
			}
			m.zoom++
			m.x *= 2
			m.y *= 2
			m.Refresh()
		}),
		newMapButton(theme.ZoomOutIcon(), func() {
			if m.zoom <= 0 {
				return
			}
			m.zoom--
			m.x /= 2
			m.y /= 2
			m.Refresh()
		}))

	move := container.NewGridWithColumns(3, layout.NewSpacer(),
		newMapButton(theme.MoveUpIcon(), func() {
			m.y--
			m.Refresh()
		}), layout.NewSpacer(),
		newMapButton(theme.NavigateBackIcon(), func() {
			m.x--
			m.Refresh()
		}), layout.NewSpacer(),
		newMapButton(theme.NavigateNextIcon(), func() {
			m.x++
			m.Refresh()
		}), layout.NewSpacer(),
		newMapButton(theme.MoveDownIcon(), func() {
			m.y++
			m.Refresh()
		}), layout.NewSpacer())

	overlay := container.NewBorder(nil, copyright, container.NewVBox(move), zoom)

	c := container.NewMax(canvas.NewRaster(m.draw), overlay)
	return widget.NewSimpleRenderer(c)
}

func (m *Map) draw(w, h int) image.Image {
	scale := 1
	tileSize := tileSize
	// TODO use retina tiles once OSM supports it in their server (text scaling issues)...
	if c := gui.CurrentApp().Driver().CanvasForObject(m); c != nil {
		scale = int(c.Scale())
		if scale < 1 {
			scale = 1
		}
		tileSize = tileSize * scale
	}

	if m.w != w || m.h != h {
		m.pixels = image.NewNRGBA(image.Rect(0, 0, w, h))
	}

	midTileX := (w - tileSize*2) / 2
	midTileY := (h - tileSize*2) / 2
	if m.zoom == 0 {
		midTileX += tileSize / 2
		midTileY += tileSize / 2
	}

	count := 1 << m.zoom
	mx := m.x + int(float32(count)/2-0.5)
	my := m.y + int(float32(count)/2-0.5)
	firstTileX := mx - int(math.Ceil(float64(midTileX)/float64(tileSize)))
	firstTileY := my - int(math.Ceil(float64(midTileY)/float64(tileSize)))

	for x := firstTileX; (x-firstTileX)*tileSize <= w+tileSize; x++ {
		for y := firstTileY; (y-firstTileY)*tileSize <= h+tileSize; y++ {
			if x < 0 || y < 0 || x >= int(count) || y >= int(count) {
				continue
			}

			src, err := getTile(x, y, m.zoom, m.cl)
			if err != nil {
				gui.LogError("tile fetch error", err)
				continue
			}

			pos := image.Pt(midTileX+(x-mx)*tileSize,
				midTileY+(y-my)*tileSize)
			scaled := src
			if scale > 1 {
				scaled = resize.Resize(uint(tileSize), uint(tileSize), src, resize.Lanczos2)
			}
			draw.Copy(m.pixels, pos, scaled, image.Rect(0, 0, tileSize, tileSize), draw.Over, nil)
		}
	}

	return m.pixels
}
